package models

import "time"

type LogEntry struct {
	Level       string    `json:"level"`
	TimeStamp   time.Time `json:"timestamp"`
	ServiceName string    `json:"service_name"`
	Op          string    `json:"op"`
	Message     string    `json:"string"`
}

type LogFilter struct {
	Level       string
	ServiceName string
	Op          string
	MessageLike string
	From        time.Time
	To          time.Time
	Limit       int
	Offset      int
}

type TimeRangeFilter struct {
	TimeFrom time.Time
	TimeTo   time.Time
}

type ServiceCount struct {
	ServiceName string `json:"service_name"`
	Count       uint64 `json:"count"`
}

type LevelCount struct {
	Level string `json:"level"`
	Count uint64 `json:"count"`
}

type ErrorCount struct {
	Message string `json:"message"`
	Count   uint64 `json:"count"`
}
