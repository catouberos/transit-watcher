package models

import (
	"time"
)

type Geolocation struct {
	Degree    float32    `redis:"-"`
	Location  Coordinate `redis:"-"`
	Speed     float32    `redis:"-"`
	VehicleID Identity   `redis:"-"`
	RouteID   Identity   `redis:"-"`
	VariantID Identity   `redis:"-"`
	Timestamp time.Time  `redis:"timestamp"`
	Hash      uint64     `redis:"hash"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
