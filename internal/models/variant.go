package models

type Variant struct {
	ID        Identity  `json:"identity" redis:"-"`
	Headsign  string    `json:"headsign" redis:"headsign"`
	ShortName string    `json:"shortName" redis:"shortName"`
	Direction Direction `json:"direction" redis:"-"`

	StopIDs []Identity `json:"stopIDs" redis:"-"`

	Hash uint64 `redis:"hash"`
}
