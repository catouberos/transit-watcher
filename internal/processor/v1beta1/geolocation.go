package v1beta1

import (
	"context"

	"codeberg.org/transit-radar/transit-watcher/internal/config"
	"codeberg.org/transit-radar/transit-watcher/internal/events"
	v1beta1 "codeberg.org/transit-radar/transit-watcher/internal/mapper/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"codeberg.org/transit-radar/transit-watcher/internal/processor"
	"codeberg.org/transit-radar/transit-watcher/internal/store"
)

type geolocationProcessor struct {
	config       *config.WorkerConfig
	eventHandler events.EventHandler
	store        store.Store
}

func NewGeolocationProcessor(config *config.WorkerConfig, eventHandler events.EventHandler, store store.Store) processor.GeolocationProcessor {
	return &geolocationProcessor{
		config:       config,
		eventHandler: eventHandler,
		store:        store,
	}
}

func (p *geolocationProcessor) Validate(ctx context.Context, geolocation models.Geolocation) error {
	// cache validation
	var latest models.Geolocation
	err := p.store.Get(ctx, CacheKey("geolocation", geolocation.VariantID.Value, geolocation.VehicleID.Value), &latest)
	if err != nil {
		return err
	}

	if !geolocation.Timestamp.After(latest.Timestamp) {
		return processor.ErrStaleData
	}

	return nil
}

func (p *geolocationProcessor) Publish(ctx context.Context, geolocation models.Geolocation) error {
	r, err := v1beta1.MapGeolocation(geolocation)
	if err != nil {
		return err
	}

	event, err := p.eventHandler.CreateEvent(r)
	if err != nil {
		return err
	}

	return p.eventHandler.Send(ctx,
		p.config.Kafka.Topic.Geolocation,
		event,
	)
}

func (p *geolocationProcessor) Memoize(ctx context.Context, geolocation models.Geolocation) error {
	return p.store.Set(ctx,
		CacheKey("geolocation", geolocation.VariantID.Value, geolocation.VehicleID.Value),
		geolocation,
	)
}
