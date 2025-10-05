package models

import (
	"fmt"
	"time"
)

// TouchEvent represents a simple timestamp indicating activity in a context
type TouchEvent struct {
	Timestamp time.Time `json:"timestamp"`
}

// ToLogLine formats the touch event as a log line: timestamp
func (t *TouchEvent) ToLogLine() string {
	return t.Timestamp.Format(time.RFC3339)
}

// ParseTouchLogLine parses a log line into a TouchEvent
func ParseTouchLogLine(line string) (*TouchEvent, error) {
	timestamp, err := time.Parse(time.RFC3339, line)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	return &TouchEvent{
		Timestamp: timestamp,
	}, nil
}

