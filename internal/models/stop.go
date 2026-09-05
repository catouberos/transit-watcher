package models

type StopType int

const (
	StopTypeUnspecified StopType = iota
	StopTypeStopPlatform
	StopTypeStation
	StopTypeEntraceExit
	StopTypeGenericNode
	StopTypeBoardingArea
)

type Stop struct {
	ID       Identity   `json:"id" redis:"-"`
	Code     string     `json:"code" redis:"code"`
	Name     string     `json:"name" redis:"name"`
	Type     StopType   `json:"type" redis:"-"`
	Location Coordinate `json:"location" redis:"-"`

	Hash uint64 `json:"hash" redis:"hash"`
}
