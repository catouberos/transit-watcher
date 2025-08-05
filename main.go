package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/wagslane/go-rabbitmq"

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
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	routePub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("route"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

	variantPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("variant"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

	geolocationPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("geolocation"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

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

				err = routePub.PublishWithContext(
					context.Background(),
					data,
					[]string{"route.event.updated"},
					rabbitmq.WithPublishOptionsContentType("application/json"),
					rabbitmq.WithPublishOptionsMandatory,
					rabbitmq.WithPublishOptionsPersistentDelivery,
					rabbitmq.WithPublishOptionsExchange("route"))
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

				err = variantPub.PublishWithContext(
					context.Background(),
					data,
					[]string{"variant.event.updated"},
					rabbitmq.WithPublishOptionsContentType("application/json"),
					rabbitmq.WithPublishOptionsMandatory,
					rabbitmq.WithPublishOptionsPersistentDelivery,
					rabbitmq.WithPublishOptionsExchange("variant"))
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

				err = geolocationPub.PublishWithContext(
					context.Background(),
					data,
					[]string{"geolocation.event.created"},
					rabbitmq.WithPublishOptionsContentType("application/json"),
					rabbitmq.WithPublishOptionsMandatory,
					rabbitmq.WithPublishOptionsPersistentDelivery,
					rabbitmq.WithPublishOptionsExchange("geolocation"))
				if err != nil {
					slog.Error("Cannot publish update to AMQP", "error", err)
				}
			}
		}
	}()

	<-done
}
