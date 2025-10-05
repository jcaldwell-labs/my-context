package models

import (
	"time"
)

// AppState represents the global application state tracking the active context
type AppState struct {
	ActiveContext *string   `json:"active_context,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
}

// HasActiveContext returns true if there is an active context
func (s *AppState) HasActiveContext() bool {
	return s.ActiveContext != nil && *s.ActiveContext != ""
}

// GetActiveContextName returns the active context name or empty string
func (s *AppState) GetActiveContextName() string {
	if s.ActiveContext == nil {
		return ""
	}
	return *s.ActiveContext
}

