package mapper

import (
	"fmt"
	"regexp"
	"strings"

	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"codeberg.org/transit-radar/transit-watcher/provider/gobus"
	"github.com/gohugoio/hashstructure"
)

var (
	re = regexp.MustCompile(`[0-9]+(-[0-9]+)?[A-Z]?`)
)

func MapGoBusStop(stop gobus.Stop) (models.Stop, error) {
	hash, err := hashstructure.Hash(stop, nil)
	if err != nil {
		return models.Stop{}, err
	}

	var lat, lng float64
	if lat, err = stop.Geometry.Coordinates[1].Float64(); err != nil {
		return models.Stop{}, err
	}
	if lng, err = stop.Geometry.Coordinates[0].Float64(); err != nil {
		return models.Stop{}, err
	}

	stopType, err := MapGoBusStopType(stop)
	if err != nil {
		return models.Stop{}, err
	}

	return models.Stop{
		ID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      stop.Property.Id.String(),
		},
		Code:     stop.Property.Code,
		Name:     stop.Property.Name,
		Location: models.Coordinate{Latitude: lat, Longitude: lng},
		Type:     stopType,
		Hash:     hash,
	}, nil
}

func MapGoBusRoute(route gobus.Route) (models.Route, error) {
	hash, err := hashstructure.Hash(route, nil)
	if err != nil {
		return models.Route{}, err
	}

	routeType, err := MapGoBusRouteType(route)
	if err != nil {
		return models.Route{}, err
	}

	return models.Route{
		ID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      route.Id.String(),
		},
		Number: route.Number,
		Name:   route.Name,
		AgencyID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      route.Info.Organization,
		},
		RouteType: routeType,

		Hash: hash,
	}, nil
}

func MapGoBusVariant(variant gobus.RouteVariant) (models.Variant, error) {
	hash, err := hashstructure.Hash(variant, nil)
	if err != nil {
		return models.Variant{}, err
	}

	stopIDs := make([]models.Identity, 0, len(variant.Stops))
	for _, stop := range variant.Stops {
		stopIDs = append(stopIDs, models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      stop.Id.String(),
		})
	}

	return models.Variant{
		ID: models.Identity{
			Identifier: models.ExternalIdentifierEBMS,
			Value:      variant.Id.String(),
		},
		Headsign:  variant.Name,
		ShortName: variant.ShortName,
		StopIDs:   stopIDs,
		Hash:      hash,
	}, nil
}

func MapGoBusRouteType(route gobus.Route) (models.RouteType, error) {
	// matches HCMC Metro
	if strings.HasPrefix(route.Number, "MRT") {
		return models.RouteTypeMetro, nil
	}

	// matches Saigon Waterbus
	if strings.HasPrefix(route.Number, "SWB") {
		return models.RouteTypeFerry, nil
	}

	// matches Public Tour Bus
	if strings.HasPrefix(route.Number, "DL") {
		return models.RouteTypeBus, nil
	}

	// exclusively match Vinbus "D-" Bus
	if strings.HasPrefix(route.Number, "D") {
		return models.RouteTypeBus, nil
	}

	// matches Public Bus
	if re.MatchString(route.Number) {
		return models.RouteTypeBus, nil
	}

	return models.RouteTypeUnspecified, fmt.Errorf("unhandled route type: %s", route.Number)
}

func MapGoBusStopType(stop gobus.Stop) (models.StopType, error) {
	// unique values of stop type, identified by
	// jq .features.[].properties.stopType | sort | uniq
	switch stop.Property.TypeName {
	case "Bãi hậu cần", "Bến Bãi QH 568", "Bến xe", "Ga Metro Số 1":
		return models.StopTypeStation, nil
	case "Biển treo", "Nhà chờ", "Ô sơn", "Trạm tạm", "Trụ - hộp thông tin", "Trụ dừng":
		return models.StopTypeStopPlatform, nil
	default:
		return models.StopTypeUnspecified, fmt.Errorf("unhandled stop type: %s", stop.Property.TypeName)
	}
}
