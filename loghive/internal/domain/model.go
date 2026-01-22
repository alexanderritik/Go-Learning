package domain

import "time"

// LogEntry is the structured representation of a log line.
type LogEntry struct {
	Timestamp time.Time
	Level     string // e.g., INFO, ERROR
	Message   string
}
