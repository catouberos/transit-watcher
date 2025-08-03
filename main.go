package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/catouberos/transit-watcher/internal/crawler"
	"github.com/catouberos/transit-watcher/internal/handler"
	"github.com/catouberos/transit-watcher/internal/models"
	"github.com/catouberos/transit-watcher/internal/queues"
)

const (
	GoBusAllDataUrl = "https://api.gobus.vn/transit/data/getAllData"
)

func main() {
	// done
	done := make(chan bool)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		done <- true
	}()

	// queue setup
	queue := queues.New("amqp://guest:guest@localhost:5672/")
	defer queue.Close()

	// crawler setup
	transitDataCrawler := crawler.New(1*time.Hour, 0*time.Second)
	defer transitDataCrawler.Close()

	crawler := crawler.New(30*time.Second, 100*time.Millisecond)
	defer crawler.Close()

	transitDataUrls := []string{GoBusAllDataUrl}

	transitDataCrawler.SetURLs(transitDataUrls)

	go func() {
		for response := range transitDataCrawler.Result() {
			// unmarshal
			routes := &[]models.GoBusRoute{}

			err := json.Unmarshal(response.Body, routes)
			if err != nil {
				continue
			}

			urls := handler.FilterTransitRoutes(routes)

			slog.Info("Updated routes/variants from GoBus", "count", len(urls))

			crawler.SetURLs(urls)
		}
	}()

	go func() {
		r, err := regexp.Compile("([0-9]+)$")
		if err != nil {
			slog.Error("Error compiling regex", "error", err)
			return
		}

		for response := range crawler.Result() {
			geolocation := []models.MultiGoGeolocation{}

			slog.Info("Updated vehicle location", "data", response)

			err := json.Unmarshal(response.Body, &geolocation)
			if err != nil {
				slog.Error("Cannot unmarshal geolocation data", "error", err)
				continue
			}

			direction := r.FindString(response.Url)

			isOutbound := true
			if direction == "1" {
				isOutbound = false
			}

			for _, loc := range geolocation {
				slog.Info("Updated vehicle location", "data", loc)

				data, err := json.Marshal(queues.NewGeolocationInsertData(&loc, isOutbound))
				if err != nil {
					slog.Error("Cannot marshal geolocation data", "error", err)
					continue
				}

				err = queue.Push(
					amqp.Publishing{
						ContentType: "application/json",
						Body:        data,
					},
					"geolocation",               // exchange
					"geolocation.event.created", // routing key
				)
				if err != nil {
					slog.Error("Cannot publish update to AMQP", "error", err)
				}
			}
		}
	}()

	<-done
}
