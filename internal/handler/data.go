package handler

import (
	"regexp"
	"strconv"

	"github.com/catouberos/transit-watcher/internal/models"

	"github.com/catouberos/transit-radar/dto"
)

func NewGeolocationInsertData(g *models.MultiGoGeolocation, isOutbound bool) (*dto.GeolocationByRouteIDAndPlateAndBoundInsert, error) {
	// remove special characters from license plate
	r, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		return nil, err
	}

	licensePlate := r.ReplaceAll([]byte(g.LicensePlate), []byte{})

	return &dto.GeolocationByRouteIDAndPlateAndBoundInsert{
		Degree:       g.Degree,
		Latitude:     g.Latitude,
		Longitude:    g.Longitude,
		Speed:        g.Speed,
		LicensePlate: string(licensePlate),
		RouteID:      g.RouteId,
		IsOutbound:   isOutbound,
		Timestamp:    g.Timestamp,
	}, nil
}

func NewRouteInsertData(g *models.GoBusRoute) (*dto.RouteUpsert, error) {
	ebmsID, err := strconv.ParseInt(g.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	return &dto.RouteUpsert{
		Number:        g.Number,
		Name:          g.Name,
		EbmsID:        ebmsID,
		Active:        true,
		OperationTime: g.Info.OperationTime,
		Ticketing:     g.Info.Ticketing,
		Organization:  g.Info.Organization,
		RouteType:     g.Info.RouteType,
	}, nil
}

func NewVariantInsertData(g *models.GoBusRouteVariantWithDescription) (*dto.VariantByRouteEbmsIDUpsert, error) {
	ebmsID, err := strconv.ParseInt(g.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	routeEbmsID, err := strconv.ParseInt(g.RouteId, 10, 64)
	if err != nil {
		return nil, err
	}

	return &dto.VariantByRouteEbmsIDUpsert{
		Name:          g.Name,
		EbmsID:        ebmsID,
		IsOutbound:    g.IsOutbound,
		RouteEbmsID:   routeEbmsID,
		ShortName:     g.ShortName,
		Description:   g.Description,
		Distance:      g.Distance,
		Duration:      g.Duration,
		StartStopName: g.StartStop,
		EndStopName:   g.EndStop,
	}, nil
}
