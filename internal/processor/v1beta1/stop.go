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

type stopProcessor struct {
	config       *config.WorkerConfig
	eventHandler events.EventHandler
	store        store.Store
}

func NewStopProcessor(config *config.WorkerConfig, eventHandler events.EventHandler, store store.Store) processor.StopProcessor {
	return &stopProcessor{
		config:       config,
		eventHandler: eventHandler,
		store:        store,
	}
}

func (p *stopProcessor) Validate(ctx context.Context, stop models.Stop) error {
	// cache validation
	var latest models.Stop
	err := p.store.Get(ctx, CacheKey("stop", stop.ID.Value), &latest)
	if err != nil {
		return err
	}

	if stop.Hash == latest.Hash {
		return processor.ErrStaleData
	}

	return nil
}

func (p *stopProcessor) Publish(ctx context.Context, stop models.Stop) error {
	r, err := v1beta1.MapStop(stop)
	if err != nil {
		return err
	}

	event, err := p.eventHandler.CreateEvent(r)
	if err != nil {
		return err
	}

	return p.eventHandler.Send(ctx,
		p.config.Kafka.Topic.Stop,
		event,
	)
}

func (p *stopProcessor) Memoize(ctx context.Context, stop models.Stop) error {
	return p.store.Set(ctx,
		CacheKey("stop", stop.ID.Value),
		stop,
	)
}
