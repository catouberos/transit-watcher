package models

import "time"

type Geolocation struct {
	Id               int64     `json:"Id"`
	Degree           float32   `json:"Deg"`
	Latitude         float32   `json:"Lat"`
	Longitude        float32   `json:"Lng"`
	Speed            float32   `json:"Speed"`
	LicensePlate     string    `json:"VehicleNumber"`
	RouteId          int64     `json:"RouteId"`
	CurrentStationId int64     `json:"currentStationId"`
	NextStationId    int64     `json:"nextStationId"`
	Timestamp        time.Time `json:"lastUpdateTime"`
}
