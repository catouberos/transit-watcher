package crawler

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// acknowledgement: this module is highly inspired by a question on stackoverflow
// see: https://stackoverflow.com/questions/63812394/periodically-crawl-api-in-golang

type PeriodicCrawler struct {
	interval time.Duration
	delay    time.Duration
	client   http.Client

	done      chan bool
	responses chan *CrawlerResponse

	mu   sync.Mutex
	urls []string
}

type CrawlerResponse struct {
	Url  string
	Body []byte
}

func New(interval, delay time.Duration) *PeriodicCrawler {
	crawler := PeriodicCrawler{
		interval: interval,
		delay:    delay,
		client:   http.Client{},

		done:      make(chan bool),
		responses: make(chan *CrawlerResponse),
	}

	go crawler.start()

	return &crawler
}

func (c *PeriodicCrawler) start() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), c.interval*2)
		defer cancel()

		delay := c.delay

		for _, url := range c.urls {
			go c.fetchWithDelay(ctx, url, delay)
			delay += 100 * time.Millisecond
		}

		select {
		case <-c.done:
			return
		case <-time.After(c.interval):
		}
	}
}

func (c *PeriodicCrawler) SetURLs(urls []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.urls = urls
}

// one-way channels

func (c *PeriodicCrawler) Result() <-chan *CrawlerResponse {
	return c.responses
}

// internal helper

func (c *PeriodicCrawler) fetch(ctx context.Context, url string) error {
	log.Printf("[INFO] Fetching %s...", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return err
	}

	err = injectCredentials(&req.Header)

	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)

	if err != nil {
		log.Printf("[ERROR] Something happened: %s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			slog.Error("An error occurred while reading response", "error", err)
			return err
		}

		c.responses <- &CrawlerResponse{
			Url:  url,
			Body: body,
		}
	} else {
		log.Printf("[DEBUG] Status code %d\n", resp.StatusCode)
	}

	return nil
}

func (c *PeriodicCrawler) fetchWithDelay(ctx context.Context, url string, delay time.Duration) error {
	<-time.After(delay)
	return c.fetch(ctx, url)
}

func (c *PeriodicCrawler) Close() {
	close(c.done)
	close(c.responses)
}
