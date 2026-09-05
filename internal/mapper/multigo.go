package mapper

import (
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"codeberg.org/transit-radar/transit-watcher/provider/multigo"
	"github.com/gohugoio/hashstructure"
)

func MapMultiGoGeolocation(geolocation multigo.Geolocation, variantId string) (models.Geolocation, error) {
	hash, err := hashstructure.Hash(geolocation, nil)
	if err != nil {
		return models.Geolocation{}, err
	}

	degree, err := geolocation.Degree.Float64()
	if err != nil {
		return models.Geolocation{}, err
	}

	latitude, err := geolocation.Latitude.Float64()
	if err != nil {
		return models.Geolocation{}, err
	}

	longitude, err := geolocation.Longitude.Float64()
	if err != nil {
		return models.Geolocation{}, err
	}

	speed, err := geolocation.Speed.Float64()
	if err != nil {
		return models.Geolocation{}, err
	}

	return models.Geolocation{
		Degree: float32(degree),
		Location: models.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		},
		Speed: float32(speed),
		VehicleID: models.Identity{
			Identifier: models.ExternalIdentifierMultiGo,
			Value:      geolocation.LicensePlate,
		},
		RouteID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      geolocation.RouteId.String(),
		},
		VariantID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      variantId,
		},
		Timestamp: geolocation.Timestamp,
		Hash:      hash,
	}, nil
}
