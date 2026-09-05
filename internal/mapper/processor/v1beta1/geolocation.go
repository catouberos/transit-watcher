package processorv1beta1mapper

import (
	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"github.com/gohugoio/hashstructure"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapGeolocation(geolocation models.Geolocation) (*processorv1beta1.Geolocation, error) {
	hash, err := hashstructure.Hash(geolocation, nil)
	if err != nil {
		return nil, err
	}

	vehicleID, err := MapIdentity(geolocation.VehicleID)
	if err != nil {
		return nil, err
	}

	variantID, err := MapIdentity(geolocation.VariantID)
	if err != nil {
		return nil, err
	}

	routeID, err := MapIdentity(geolocation.RouteID)
	if err != nil {
		return nil, err
	}

	builder := &processorv1beta1.Geolocation_builder{
		Degree: float64(geolocation.Degree),
		Location: &latlng.LatLng{
			Latitude:  geolocation.Location.Latitude,
			Longitude: geolocation.Location.Longitude,
		},
		Speed:     float64(geolocation.Speed),
		VehicleId: vehicleID,
		RouteId:   routeID,
		VariantId: variantID,
		Timestamp: timestamppb.New(geolocation.Timestamp),
		Hash:      hash,
	}

	return builder.Build(), nil
}
