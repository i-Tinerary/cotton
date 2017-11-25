package common

import "time"

type Plan struct {
	PlanName string
	PlanUser string
	Start    time.Time
	Events   []Event
}

type Event struct {
	PlaceID int
	Time    time.Time
}
