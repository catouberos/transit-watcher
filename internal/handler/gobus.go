package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/catouberos/transit-watcher/internal/crawler"
	"github.com/catouberos/transit-watcher/internal/models"
	"github.com/wagslane/go-rabbitmq"
)

func GoBusDataHandler(conn *rabbitmq.Conn, responses <-chan *crawler.CrawlerResponse, results chan []string) error {
	routePub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("route"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	variantPub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("variant"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return err
	}

	for response := range responses {
		// unmarshal
		routes := []models.GoBusRoute{}
		variants := []models.GoBusRouteVariantWithDescription{}

		err := json.Unmarshal(response.Body, &routes)
		if err != nil {
			continue
		}

		for _, route := range routes {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logger := slog.New(slog.Default().Handler())
			logger.With("route", route)

			params, err := NewRouteInsertData(&route)
			if err != nil {
				logger.Error("Error parsing route data", "error", err)
			}

			data, err := json.Marshal(params)
			if err != nil {
				logger.Error("Error marshal route data", "error", err)
			}

			err = routePub.PublishWithContext(
				ctx,
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
				var description string

				if variant.IsOutbound {
					description = route.Info.OutboundDescription
				} else {
					description = route.Info.InboundDescription
				}

				variants = append(variants, models.GoBusRouteVariantWithDescription{
					GoBusRouteVariant: variant,
					Description:       description,
				})
			}
		}

		<-time.After(1 * time.Second)

		for _, variant := range variants {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			logger := slog.New(slog.Default().Handler())
			logger.With("variant", variant)

			params, err := NewVariantInsertData(&variant)
			if err != nil {
				logger.Error("Error parsing variant data", "error", err)
			}

			data, err := json.Marshal(params)
			if err != nil {
				logger.Error("Error marshal variant data", "error", err)
			}

			err = variantPub.PublishWithContext(
				ctx,
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

		urls := FilterTransitRoutes(&routes)

		slog.Info("Updated routes/variants from GoBus", "count", len(urls))

		results <- urls
	}

	return nil
}
