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

type routeProcessor struct {
	config       *config.WorkerConfig
	eventHandler events.EventHandler
	store        store.Store
}

func NewRouteProcessor(config *config.WorkerConfig, eventHandler events.EventHandler, store store.Store) processor.RouteProcessor {
	return &routeProcessor{
		config:       config,
		eventHandler: eventHandler,
		store:        store,
	}
}

func (p *routeProcessor) Validate(ctx context.Context, route models.Route) error {
	var latest models.Geolocation
	err := p.store.Get(ctx, CacheKey("route", route.ID.Value), &latest)
	if err != nil {
		return err
	}

	if route.Hash == latest.Hash {
		return processor.ErrStaleData
	}

	return nil
}

func (p *routeProcessor) Publish(ctx context.Context, route models.Route) error {
	r, err := v1beta1.MapRoute(route)
	if err != nil {
		return err
	}

	event, err := p.eventHandler.CreateEvent(r)
	if err != nil {
		return err
	}

	return p.eventHandler.Send(ctx,
		p.config.Kafka.Topic.Route,
		event,
	)
}

func (p *routeProcessor) Memoize(ctx context.Context, route models.Route) error {
	if err := p.store.Set(ctx, CacheKey("route", route.ID.Value), route); err != nil {
		return err
	}

	switch route.RouteType {
	case models.RouteTypeBus:
		if err := p.store.Add(ctx, CacheKey("routes", "bus"), route.ID.Value); err != nil {
			return err
		}
	case models.RouteTypeMetro:
		if err := p.store.Add(ctx, CacheKey("routes", "metro"), route.ID.Value); err != nil {
			return err
		}
	}

	if err := p.store.Add(ctx, CacheKey("routes"), route.ID.Value); err != nil {
		return err
	}

	return nil
}
