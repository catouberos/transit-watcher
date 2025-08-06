package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"regexp"
	"time"

	"github.com/catouberos/transit-watcher/internal/crawler"
	"github.com/catouberos/transit-watcher/internal/models"
	"github.com/wagslane/go-rabbitmq"
)

func MultiGoGeolocationHandler(conn *rabbitmq.Conn, responses <-chan *crawler.CrawlerResponse) error {
	r, err := regexp.Compile("([0-9]+)$")
	if err != nil {
		return err
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("geolocation"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	for response := range responses {
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			slog.Info("Updated vehicle location", "location", loc)

			params, err := NewGeolocationInsertData(&loc, isOutbound)
			if err != nil {
				slog.Error("Error parsing geolocation data", "error", err)
				continue
			}

			data, err := json.Marshal(params)
			if err != nil {
				slog.Error("Cannot marshal geolocation data", "error", err)
				continue
			}

			err = publisher.PublishWithContext(
				ctx,
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

	return nil
}
