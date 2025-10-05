package models

import (
	"fmt"
	"strings"
	"time"
)

// TransitionType represents the type of context transition
type TransitionType string

const (
	TransitionStart  TransitionType = "start"
	TransitionStop   TransitionType = "stop"
	TransitionSwitch TransitionType = "switch"
)

// ContextTransition represents a log entry recording a change between contexts
type ContextTransition struct {
	Timestamp       time.Time      `json:"timestamp"`
	PreviousContext *string        `json:"previous_context,omitempty"`
	NewContext      *string        `json:"new_context,omitempty"`
	TransitionType  TransitionType `json:"transition_type"`
}

// Validate checks if the transition has valid data
func (ct *ContextTransition) Validate() error {
	switch ct.TransitionType {
	case TransitionStart:
		if ct.NewContext == nil || *ct.NewContext == "" {
			return fmt.Errorf("start transition must have new_context")
		}
	case TransitionStop:
		if ct.PreviousContext == nil || *ct.PreviousContext == "" {
			return fmt.Errorf("stop transition must have previous_context")
		}
		if ct.NewContext != nil {
			return fmt.Errorf("stop transition cannot have new_context")
		}
	case TransitionSwitch:
		if ct.PreviousContext == nil || *ct.PreviousContext == "" {
			return fmt.Errorf("switch transition must have previous_context")
		}
		if ct.NewContext == nil || *ct.NewContext == "" {
			return fmt.Errorf("switch transition must have new_context")
		}
	default:
		return fmt.Errorf("invalid transition type: %s", ct.TransitionType)
	}

	return nil
}

// ToLogLine formats the transition as a log line: timestamp|prev|new|type
func (ct *ContextTransition) ToLogLine() string {
	prev := "NULL"
	if ct.PreviousContext != nil {
		prev = *ct.PreviousContext
	}

	new := "NULL"
	if ct.NewContext != nil {
		new = *ct.NewContext
	}

	return fmt.Sprintf("%s|%s|%s|%s",
		ct.Timestamp.Format(time.RFC3339),
		prev,
		new,
		ct.TransitionType)
}

// ParseTransitionLogLine parses a log line into a ContextTransition
func ParseTransitionLogLine(line string) (*ContextTransition, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid transition log line format")
	}

	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	var prevContext *string
	if parts[1] != "NULL" {
		prevContext = &parts[1]
	}

	var newContext *string
	if parts[2] != "NULL" {
		newContext = &parts[2]
	}

	transitionType := TransitionType(parts[3])

	return &ContextTransition{
		Timestamp:       timestamp,
		PreviousContext: prevContext,
		NewContext:      newContext,
		TransitionType:  transitionType,
	}, nil
}

