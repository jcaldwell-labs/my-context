package models

import (
	"fmt"
	"strings"
	"time"
)

// Context represents a work session with associated metadata
type Context struct {
	Name             string     `json:"name"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	Status           string     `json:"status"` // "active" or "stopped"
	SubdirectoryPath string     `json:"subdirectory_path,omitempty"`
	IsArchived       bool       `json:"is_archived"`
}

// Validate checks if the context has valid data
func (c *Context) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("context name cannot be empty")
	}

	if strings.ContainsAny(c.Name, "/\\") {
		return fmt.Errorf("context name cannot contain path separators")
	}

	if len(c.Name) > 200 {
		return fmt.Errorf("context name must be 200 characters or less")
	}

	if c.Status != "active" && c.Status != "stopped" {
		return fmt.Errorf("status must be 'active' or 'stopped'")
	}

	if c.Status == "stopped" && c.EndTime == nil {
		return fmt.Errorf("stopped context must have an end time")
	}

	if c.Status == "active" && c.EndTime != nil {
		return fmt.Errorf("active context cannot have an end time")
	}

	if c.IsArchived && c.Status == "active" {
		return fmt.Errorf("cannot archive an active context")
	}

	return nil
}

// Duration calculates the duration of the context
func (c *Context) Duration() time.Duration {
	if c.EndTime != nil {
		return c.EndTime.Sub(c.StartTime)
	}
	return time.Since(c.StartTime)
}

// IsActive returns true if context is currently active
func (c *Context) IsActive() bool {
	return c.Status == "active"
}
