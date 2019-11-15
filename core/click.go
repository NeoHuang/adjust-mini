package core

import "time"

type Click struct {
	Tracker   string
	CreatedAt time.Time
}

func NewClick(tracker string) {
	return &Click{
		Tracker:   tracker,
		CreatedAt: time.Now(),
	}
}
