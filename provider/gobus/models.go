package gobus

import (
	"encoding/json"

	"github.com/twpayne/go-geom/encoding/geojson"
)

type Route struct {
	Id       json.Number    `json:"_id"`
	Number   string         `json:"number" conform:"trim"`
	Name     string         `json:"name" conform:"trim"`
	Info     RouteInfo      `json:"info"`
	Variants []RouteVariant `json:"vars"`
}

type RouteInfo struct {
	Id                  json.Number `json:"_id"`
	InboundDescription  string      `json:"inBoundDescription" conform:"trim"`
	OutboundDescription string      `json:"outBoundDescription" conform:"trim"`
	OperationTime       string      `json:"operationTime" conform:"trim"`
	Organization        string      `json:"orgs" conform:"trim"`
	Ticketing           string      `json:"tickets" conform:"trim"`
	Duration            string      `json:"timeOfTrip" conform:"trim"`
	TotalTrip           string      `json:"totalTrip" conform:"trim"`
	RouteType           string      `json:"busType" conform:"trim"`
}

type RouteVariant struct {
	Id         json.Number        `json:"_id"`
	RouteId    json.Number        `json:"routeId"`
	Name       string             `json:"name" conform:"trim"`
	ShortName  string             `json:"shortName" conform:"trim"`
	Distance   float32            `json:"distance" conform:"trim"`
	StartStop  string             `json:"startStop" conform:"trim"`
	EndStop    string             `json:"endStop" conform:"trim"`
	IsOutbound bool               `json:"isOutbound" conform:"trim"`
	Duration   int32              `json:"runningTime" conform:"trim"`
	Stops      []RouteVariantStop `json:"stops"`
	GeoJSON    geojson.Feature    `json:"geojson"`
}

type RouteVariantWithDescription struct {
	RouteVariant

	Description string `json:"-" conform:"trim"`
}

type RouteVariantStop struct {
	Id            json.Number `json:"_id"`
	Code          string      `json:"code" conform:"trim"`
	Name          string      `json:"name" conform:"trim"`
	Routes        string      `json:"routes" conform:"trim"`
	Type          string      `json:"stopType" conform:"trim"`
	AddressNumber string      `json:"addressNo" conform:"trim"`
	AddressStreet string      `json:"street" conform:"trim"`
	AddressWard   string      `json:"ward" conform:"trim"`
	AddressZone   string      `json:"zone" conform:"trim"`
	Latitude      float32     `json:"lat" conform:"trim"`
	Longitude     float32     `json:"lng" conform:"trim"`
}

type StopResponse struct {
	// here we dont care about additional geojson features
	Stops []Stop `json:"features"`
}

type Stop struct {
	Geometry StopGeometry `json:"geometry"`
	Property StopProperty `json:"properties"`
}

type StopGeometry struct {
	Coordinates [2]json.Number `json:"coordinates"`
}

type StopProperty struct {
	Id       json.Number `json:"stopId"`
	Code     string      `json:"code" conform:"trim"`
	Name     string      `json:"name" conform:"trim"`
	TypeName string      `json:"stopType" conform:"trim"`
}
