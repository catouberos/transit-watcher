package crawler

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// acknowledgement: this module is highly inspired by a question on stackoverflow
// see: https://stackoverflow.com/questions/63812394/periodically-crawl-api-in-golang

type PeriodicCrawlerOptions struct {
	Interval time.Duration
	Timeout  time.Duration
	Delay    time.Duration
}

type PeriodicCrawler struct {
	interval time.Duration
	delay    time.Duration
	client   http.Client

	doneChan   chan bool
	resultChan chan []byte

	mu   sync.Mutex
	urls []string
}

func NewPeriodicCrawler(interval, delay time.Duration) *PeriodicCrawler {
	props := &PeriodicCrawlerOptions{
		Interval: interval,
		Delay:    delay,
	}

	return &PeriodicCrawler{
		interval: props.Interval,
		delay:    props.Delay,
		client:   http.Client{},

		doneChan:   make(chan bool, 1),
		resultChan: make(chan []byte),
	}
}

func (c *PeriodicCrawler) Start() error {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), c.interval*2)

			delay := c.delay

			for _, url := range c.urls {
				go c.fetchWithDelay(ctx, url, delay)
				delay += 100 * time.Millisecond
			}

			<-time.After(c.interval)
			cancel()
		}
	}()

	<-c.doneChan

	// cleanup
	close(c.doneChan)
	close(c.resultChan)

	return nil
}

func (c *PeriodicCrawler) SetURLs(urls []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.urls = urls
}

// one-way channels

func (c *PeriodicCrawler) Result() <-chan []byte {
	return c.resultChan
}

func (c *PeriodicCrawler) Done() chan<- bool {
	return c.doneChan
}

// internal helper

func (c *PeriodicCrawler) fetch(ctx context.Context, url string) error {
	log.Printf("[INFO] Fetching %s...", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return err
	}

	err = InjectCredentials(&req.Header)

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
			return err
		}

		c.resultChan <- body
	} else {
		log.Printf("[DEBUG] Status code %d\n", resp.StatusCode)
	}

	return nil
}

func (c *PeriodicCrawler) fetchWithDelay(ctx context.Context, url string, delay time.Duration) error {
	<-time.After(delay)
	return c.fetch(ctx, url)
}
