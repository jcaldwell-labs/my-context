package models

import (
	"fmt"
	"strings"
	"time"
)

// ContextMetadata represents additional metadata for contexts
type ContextMetadata struct {
	CreatedBy string   `json:"created_by,omitempty"` // User who created the context
	Parent    string   `json:"parent,omitempty"`     // Parent context name for hierarchy
	Labels    []string `json:"labels,omitempty"`     // Labels for categorization and search
}

// Validate checks if metadata is valid
func (m *ContextMetadata) Validate() error {
	if m.CreatedBy != "" && len(m.CreatedBy) > 100 {
		return fmt.Errorf("created_by must be 100 characters or less")
	}

	if m.Parent != "" && len(m.Parent) > 200 {
		return fmt.Errorf("parent context name must be 200 characters or less")
	}

	if len(m.Labels) > 10 {
		return fmt.Errorf("cannot have more than 10 labels")
	}

	for _, label := range m.Labels {
		if label == "" {
			return fmt.Errorf("labels cannot be empty strings")
		}
		if len(label) > 50 {
			return fmt.Errorf("label '%s' must be 50 characters or less", label)
		}
		if strings.ContainsAny(label, " \t\n\r") {
			return fmt.Errorf("label '%s' cannot contain whitespace", label)
		}
	}

	return nil
}

// ContextWithMetadata extends the base Context with metadata fields
type ContextWithMetadata struct {
	Name             string          `json:"name"`
	StartTime        time.Time       `json:"start_time"`
	EndTime          *time.Time      `json:"end_time,omitempty"`
	Status           string          `json:"status"` // "active" or "stopped"
	SubdirectoryPath string          `json:"subdirectory_path,omitempty"`
	IsArchived       bool            `json:"is_archived"`
	Metadata         ContextMetadata `json:"metadata,omitempty"`
}

// NewContextWithMetadata creates a new context with metadata
func NewContextWithMetadata(name string, createdBy string, parent string, labels []string) *ContextWithMetadata {
	now := time.Now()
	metadata := ContextMetadata{
		CreatedBy: createdBy,
		Parent:    parent,
		Labels:    labels,
	}

	return &ContextWithMetadata{
		Name:       name,
		StartTime:  now,
		Status:     "active",
		IsArchived: false,
		Metadata:   metadata,
	}
}

// Validate checks if the context with metadata has valid data
func (c *ContextWithMetadata) Validate() error {
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

	// Validate metadata
	if err := c.Metadata.Validate(); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}

	return nil
}

// Duration calculates the duration of the context
func (c *ContextWithMetadata) Duration() time.Duration {
	if c.EndTime != nil {
		return c.EndTime.Sub(c.StartTime)
	}
	return time.Since(c.StartTime)
}

// IsActive returns true if context is currently active
func (c *ContextWithMetadata) IsActive() bool {
	return c.Status == "active"
}

// HasLabel checks if the context has a specific label
func (c *ContextWithMetadata) HasLabel(label string) bool {
	for _, l := range c.Metadata.Labels {
		if l == label {
			return true
		}
	}
	return false
}
