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

		close(done)
	}()

	// queue setup
	queue := queues.New("amqp://guest:guest@localhost:5672/")
	defer queue.Close()

	// crawler setup
	transitDataCrawler := crawler.New(2*time.Hour, 0*time.Second)
	defer transitDataCrawler.Close()

	crawler := crawler.New(30*time.Second, 100*time.Millisecond)
	defer crawler.Close()

	transitDataUrls := []string{GoBusAllDataUrl}

	transitDataCrawler.SetURLs(transitDataUrls)

	go func() {
		for response := range transitDataCrawler.Result() {
			// unmarshal
			routes := []models.GoBusRoute{}
			variants := []models.GoBusRouteVariant{}

			err := json.Unmarshal(response.Body, &routes)
			if err != nil {
				continue
			}

			for _, route := range routes {
				logger := slog.New(slog.Default().Handler())
				logger.With("route", route)

				params, err := queues.NewRouteInsertData(&route)
				if err != nil {
					logger.Error("Error parsing route data", "error", err)
				}

				data, err := json.Marshal(params)
				if err != nil {
					logger.Error("Error marshal route data", "error", err)
				}

				err = queue.Push(amqp.Publishing{
					ContentType: "application/json",
					Body:        data,
				}, "route", "route.event.updated")
				if err != nil {
					logger.Error("Cannot publish route update", "error", err)
				}

				for _, variant := range route.Variants {
					variants = append(variants, variant)
				}
			}

			<-time.After(1 * time.Second)

			for _, variant := range variants {
				logger := slog.New(slog.Default().Handler())
				logger.With("variant", variant)

				params, err := queues.NewVariantInsertData(&variant)
				if err != nil {
					logger.Error("Error parsing variant data", "error", err)
				}

				data, err := json.Marshal(params)
				if err != nil {
					logger.Error("Error marshal variant data", "error", err)
				}

				err = queue.Push(amqp.Publishing{
					ContentType: "application/json",
					Body:        data,
				}, "variant", "variant.event.updated")
				if err != nil {
					logger.Error("Cannot publish variant update", "error", err)
				}
			}

			urls := handler.FilterTransitRoutes(&routes)

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
				slog.Info("Updated vehicle location", "location", loc)

				params, err := queues.NewGeolocationInsertData(&loc, isOutbound)
				if err != nil {
					slog.Error("Error parsing geolocation data", "error", err)
					continue
				}

				data, err := json.Marshal(params)
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
