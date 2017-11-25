package common

import "time"

type Plan struct {
	Name   string
	Events []Event
}

type Event struct {
	PlaceName string
	Time      time.Time
}
