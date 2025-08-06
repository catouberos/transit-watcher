package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wagslane/go-rabbitmq"

	"github.com/catouberos/transit-watcher/internal/crawler"
	"github.com/catouberos/transit-watcher/internal/handler"
)

const (
	GoBusAllDataUrl = "https://api.gobus.vn/transit/data/getAllData"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		cancel()
	}()

	// queue setup
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// crawler setup
	transitDataCrawler := crawler.New(2*time.Hour, 0*time.Second)
	defer transitDataCrawler.Close()

	crawler := crawler.New(30*time.Second, 100*time.Millisecond)
	defer crawler.Close()

	transitDataUrls := []string{GoBusAllDataUrl}

	transitDataCrawler.SetURLs(transitDataUrls)

	urls := make(chan []string)
	defer close(urls)

	go func() {
		if err := handler.GoBusDataHandler(conn, transitDataCrawler.Result(), urls); err != nil {
			slog.Error("Error starting GoBus handler", "error", err)
		}

		cancel()
	}()

	go func() {
		for {
			crawler.SetURLs(<-urls)
		}
	}()

	go func() {
		if err := handler.MultiGoGeolocationHandler(conn, crawler.Result()); err != nil {
			slog.Error("Error starting MultiGo handler", "error", err)
		}

		cancel()
	}()

	<-ctx.Done()
}
