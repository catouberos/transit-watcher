package processorv1beta1mapper

import (
	processorv1beta1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/processor/v1beta1"
	radarv1mapper "codeberg.org/transit-radar/transit-watcher/internal/mapper/radar/v1"
	"codeberg.org/transit-radar/transit-watcher/internal/models"
	"github.com/gohugoio/hashstructure"
	"github.com/icza/gox/imagex/colorx"
	"google.golang.org/genproto/googleapis/type/color"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func MapRoute(route models.Route) (*processorv1beta1.Route, error) {
	hash, err := hashstructure.Hash(route, nil)
	if err != nil {
		return nil, err
	}

	id, err := MapIdentity(route.ID)
	if err != nil {
		return nil, err
	}

	routeType, err := radarv1mapper.RouteType(route)
	if err != nil {
		return nil, err
	}

	agencyID, err := MapIdentity(route.AgencyID)
	if err != nil {
		return nil, err
	}

	builder := &processorv1beta1.Route_builder{
		Id:          id,
		AgencyId:    agencyID,
		Number:      route.Number,
		Name:        route.Name,
		ShortName:   route.ShortName,
		Description: route.Description,
		Type:        routeType,
		Hash:        hash,
	}

	if route.Color != nil {
		routeColor, err := colorx.ParseHexColor(*route.Color)
		if err != nil {
			return nil, err
		}
		builder.Color = &color.Color{
			Red:   float32(routeColor.R),
			Green: float32(routeColor.G),
			Blue:  float32(routeColor.B),
			Alpha: &wrapperspb.FloatValue{
				Value: float32(routeColor.A),
			},
		}
	}

	return builder.Build(), nil
}
