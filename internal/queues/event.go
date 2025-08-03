package queues

import (
	"github.com/catouberos/transit-watcher/internal/models"

	"github.com/catouberos/geoloc/dto"
)

func NewGeolocationInsertData(g *models.MultiGoGeolocation, isOutbound bool) *dto.GeolocationByPlateAndBoundInsert {
	return &dto.GeolocationByPlateAndBoundInsert{
		Degree:       g.Degree,
		Latitude:     g.Latitude,
		Longitude:    g.Longitude,
		Speed:        g.Speed,
		LicensePlate: g.LicensePlate,
		RouteID:      g.RouteId,
		IsOutbound:   isOutbound,
		Timestamp:    g.Timestamp,
	}
}
