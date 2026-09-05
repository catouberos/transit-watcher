package radarv1mapper

import (
	"fmt"
	"regexp"
	"strings"

	radarv1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/radar/v1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
)

func RouteType(route models.Route) (radarv1.RouteType, error) {
	// matches HCMC Metro
	if strings.HasPrefix(route.Number, "MRT") {
		return radarv1.RouteType_ROUTE_TYPE_METRO, nil
	}

	// matches Saigon Waterbus
	if strings.HasPrefix(route.Number, "SWB") {
		return radarv1.RouteType_ROUTE_TYPE_FERRY, nil
	}

	// matches Public Tour Bus
	if strings.HasPrefix(route.Number, "DL") {
		return radarv1.RouteType_ROUTE_TYPE_BUS, nil
	}

	// exclusively match Vinbus "D-" Bus
	if strings.HasPrefix(route.Number, "D") {
		return radarv1.RouteType_ROUTE_TYPE_BUS, nil
	}

	// matches Public Bus
	re, err := regexp.Compile(`[0-9]+(-[0-9]+)?[A-Z]?`)
	if err != nil {
		return radarv1.RouteType_ROUTE_TYPE_UNSPECIFIED, err
	}
	if re.MatchString(route.Number) {
		return radarv1.RouteType_ROUTE_TYPE_BUS, nil
	}

	return radarv1.RouteType_ROUTE_TYPE_UNSPECIFIED, fmt.Errorf("unhandled route type for: %s", route.Number)
}
