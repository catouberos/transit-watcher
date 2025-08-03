package handler

import (
	"fmt"
	"time"

	"github.com/catouberos/transit-watcher/internal/models"
)

const (
	MultiGoGeolocationUrl = "https://multipass-api.golabs.vn/v2/public/busmap/route_bus_gps?regionCode=hcm&routeId=%s&direction=%d"
)

func FilterTransitRoutes(routes *[]models.GoBusRoute) []string {
	urls := []string{}

	for _, route := range *routes {
		from, to, err := ParseOperationTime(route.Info.OperationTime)

		if err != nil {
			continue
		}

		// skip non-operating routes
		if time.Now().Before(from) || to.After(time.Now().Add(-2*time.Hour)) {
			continue
		}

		for i := 0; i < len(route.Variants); i++ {
			urls = append(urls, fmt.Sprintf(MultiGoGeolocationUrl, route.Id, i))
		}
	}

	return urls
}
