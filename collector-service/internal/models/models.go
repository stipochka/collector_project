package models

import "time"

type TelemetryRecord struct {
	Timestamp   time.Time
	Level       string
	ServiceName string
	Message     string
	Op          string
}
