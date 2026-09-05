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

type variantProcessor struct {
	config       *config.WorkerConfig
	eventHandler events.EventHandler
	store        store.Store
}

func NewVariantProcessor(config *config.WorkerConfig, eventHandler events.EventHandler, store store.Store) processor.VariantProcessor {
	return &variantProcessor{
		config:       config,
		eventHandler: eventHandler,
		store:        store,
	}
}

func (p *variantProcessor) Validate(ctx context.Context, route models.Route, variant models.Variant) error {
	var latest models.Geolocation
	err := p.store.Get(ctx, CacheKey("variant", variant.ID.Value), &latest)
	if err != nil {
		return err
	}

	if variant.Hash == latest.Hash {
		return processor.ErrStaleData
	}

	return nil
}

func (p *variantProcessor) Publish(ctx context.Context, route models.Route, variant models.Variant) error {
	r, err := v1beta1.MapTrip(route, variant)
	if err != nil {
		return err
	}

	event, err := p.eventHandler.CreateEvent(r)
	if err != nil {
		return err
	}

	if err := p.eventHandler.Send(ctx, p.config.Kafka.Topic.Variant, event); err != nil {
		return err
	}

	for i, stop := range variant.StopIDs {
		stops, err := v1beta1.MapTripStop(route.ID, variant.ID, stop, int32(i))
		if err != nil {
			return err
		}

		stopsEvent, err := p.eventHandler.CreateEvent(stops)
		if err != nil {
			return err
		}

		if err := p.eventHandler.Send(ctx, p.config.Kafka.Topic.VariantStops, stopsEvent); err != nil {
			return err
		}
	}

	return nil
}

func (p *variantProcessor) Memoize(ctx context.Context, route models.Route, variant models.Variant) error {
	if err := p.store.Set(ctx, CacheKey("variant", variant.ID.Value), variant); err != nil {
		return err
	}

	if err := p.store.Add(ctx, CacheKey("route", route.ID.Value, "variants"), variant.ID.Value); err != nil {
		return err
	}

	return nil
}
