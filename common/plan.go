package common

import "time"

type Plan struct {
	PlanName string
	PlanUser string
	Events   []Event
}

type Event struct {
	PlaceID string
	Time    time.Time
}
