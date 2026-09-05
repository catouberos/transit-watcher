package radarv1mapper

import (
	"testing"

	radarv1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/radar/v1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRouteType(t *testing.T) {
	tests := []struct {
		name        string
		routeNumber string
		routeType   radarv1.RouteType
	}{{
		name:        "PublicBus",
		routeNumber: "161",
		routeType:   radarv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "PublicBusCharacter",
		routeNumber: "162V",
		routeType:   radarv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "Vinbus",
		routeNumber: "D4",
		routeType:   radarv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "IntercityBus",
		routeNumber: "60-01",
		routeType:   radarv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "TourBus",
		routeNumber: "DL01",
		routeType:   radarv1.RouteType_ROUTE_TYPE_BUS,
	}, {
		name:        "HCMCMetro",
		routeNumber: "MRT1",
		routeType:   radarv1.RouteType_ROUTE_TYPE_METRO,
	}, {
		name:        "SaigonWaterbus",
		routeNumber: "SWB1",
		routeType:   radarv1.RouteType_ROUTE_TYPE_FERRY,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			result, err := RouteType(models.Route{Number: test.routeNumber})
			assert.NoError(err)
			assert.Equal(test.routeType, result)
		})
	}
}
