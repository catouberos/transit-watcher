package processorv1beta1mapper

import (
	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"github.com/gohugoio/hashstructure"
)

func MapTrip(route models.Route, variant models.Variant) (*processorv1beta1.Trip, error) {
	hash, err := hashstructure.Hash(variant, nil)
	if err != nil {
		return nil, err
	}

	id, err := MapIdentity(variant.ID)
	if err != nil {
		return nil, err
	}

	routeID, err := MapIdentity(route.ID)
	if err != nil {
		return nil, err
	}

	builder := processorv1beta1.Trip_builder{
		Id:        id,
		RouteId:   routeID,
		Headsign:  variant.Headsign,
		ShortName: &variant.ShortName,
		Direction: new(MapDirection(variant.Direction)),
		Hash:      hash,
	}

	return builder.Build(), nil
}

func MapDirection(direction models.Direction) processorv1beta1.Direction {
	switch direction {
	case models.DirectionOneDirection:
		return processorv1beta1.Direction_DIRECTION_ONE_DIRECTION
	case models.DirectionOppositeDirection:
		return processorv1beta1.Direction_DIRECTION_OPPOSITE_DIRECTION
	default:
		return processorv1beta1.Direction_DIRECTION_UNSPECIFIED
	}
}

func MapTripStop(routeID, variantID, stopID models.Identity, orderScore int32) (*processorv1beta1.TripStop, error) {
	pbRouteID, err := MapIdentity(routeID)
	if err != nil {
		return nil, err
	}

	pbTripID, err := MapIdentity(variantID)
	if err != nil {
		return nil, err
	}

	pbStopID, err := MapIdentity(stopID)
	if err != nil {
		return nil, err
	}

	builder := processorv1beta1.TripStop_builder{
		RouteId:    pbRouteID,
		TripId:     pbTripID,
		StopId:     pbStopID,
		OrderScore: orderScore,
	}

	return builder.Build(), nil
}
