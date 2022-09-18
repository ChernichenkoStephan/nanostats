package tg

import "time"

const (
	AUTO = iota
	MAMUAL
)

type ShotType int

type Shot struct {
	Amount  int
	Type    ShotType
	Created time.Time
}
