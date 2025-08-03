package queues

import (
	"github.com/catouberos/transit-watcher/internal/models"

	"github.com/catouberos/geoloc/dto"
)

func NewGeolocationInsertData(g *models.MultiGoGeolocation, variation int64) *dto.GeolocationByPlateInsert {
	return &dto.GeolocationByPlateInsert{
		Degree:       g.Degree,
		Latitude:     g.Latitude,
		Longitude:    g.Longitude,
		Speed:        g.Speed,
		LicensePlate: g.LicensePlate,
		RouteID:      g.RouteId,
		VariationID:  variation,
		Timestamp:    g.Timestamp,
	}
}
