package models

type Route struct {
	Id       string         `json:"_id"`
	Number   string         `json:"number"`
	Name     string         `json:"name"`
	Info     RouteInfo      `json:"info"`
	Variants []RouteVariant `json:"vars"`
}

type RouteInfo struct {
	Id                  string `json:"_id"`
	InboundDescription  string `json:"inBoundDescription"`
	OutboundDescription string `json:"outBoundDescription"`
	OperationTime       string `json:"operationTime"`
	Organization        string `json:"orgs"`
	Ticketing           string `json:"tickets"`
	Duration            string `json:"timeOfTrip"`
	TotalTrip           string `json:"totalTrip"`
	RouteType           string `json:"busType"`
}

type RouteVariant struct {
	Id         string             `json:"_id"`
	RouteId    string             `json:"routeId"`
	Name       string             `json:"name"`
	ShortName  string             `json:"shortName"`
	Distance   float32            `json:"distance"`
	StartStop  string             `json:"startStop"`
	EndStop    string             `json:"endStop"`
	IsOutbound bool               `json:"isOutbound"`
	Duration   int32              `json:"runningTime"`
	Stops      []RouteVariantStop `json:"stops"`
}

type RouteVariantStop struct {
	Id            string  `json:"_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Routes        string  `json:"routes"`
	Type          string  `json:"stopType"`
	AddressNumber string  `json:"addressNo"`
	AddressStreet string  `json:"street"`
	AddressWard   string  `json:"ward"`
	AddressZone   string  `json:"zone"`
	Latitude      float32 `json:"lat"`
	Longitude     float32 `json:"lng"`
}
