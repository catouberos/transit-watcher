package processorv1beta1mapper

import (
	"errors"

	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	radarv1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/radar/v1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"google.golang.org/genproto/googleapis/type/latlng"
)

func MapStop(stop models.Stop) (*processorv1beta1.Stop, error) {
	id, err := MapIdentity(stop.ID)
	if err != nil {
		return nil, err
	}

	stopType, err := MapStopType(stop.Type)
	if err != nil {
		return nil, err
	}

	builder := processorv1beta1.Stop_builder{
		Id:   id,
		Code: stop.Code,
		Name: stop.Name,
		Type: stopType,
		Location: &latlng.LatLng{
			Latitude:  stop.Location.Latitude,
			Longitude: stop.Location.Longitude,
		},
	}

	return builder.Build(), nil
}

func MapStopType(stopType models.StopType) (radarv1.StopType, error) {
	switch stopType {
	case models.StopTypeStopPlatform:
		return radarv1.StopType_STOP_TYPE_STOP_PLATFORM, nil
	case models.StopTypeStation:
		return radarv1.StopType_STOP_TYPE_STATION, nil
	}

	return radarv1.StopType_STOP_TYPE_UNSPECIFIED, errors.New("unhandled stop type")
}
