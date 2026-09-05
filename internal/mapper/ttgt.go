package mapper

import (
	"time"

	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"codeberg.org/transit-radar/transit-watcher/provider/ttgt"
)

func MapTTGTGeolocation(vehicle ttgt.TransitVehicle, routeID, variantID string) models.Geolocation {
	return models.Geolocation{
		Degree: vehicle.Angle,
		Location: models.Coordinate{
			Latitude:  vehicle.Coordinate[1],
			Longitude: vehicle.Coordinate[0],
		},
		RouteID: models.Identity{
			Value:      routeID,
			Identifier: models.ExternalIdentifierEBMS,
		},
		VariantID: models.Identity{
			Value:      variantID,
			Identifier: models.ExternalIdentifierEBMS,
		},
		VehicleID: models.Identity{
			Value:      "",
			Identifier: models.ExternalIdentifierEBMS,
		},
		Timestamp: time.Now(),
	}
}
