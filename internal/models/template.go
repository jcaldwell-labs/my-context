package models

import (
	"fmt"
	"strings"
)

// Template represents a context template with predefined settings
type Template struct {
	Name         string   `yaml:"name" json:"name"`
	Description  string   `yaml:"description" json:"description"`
	Labels       []string `yaml:"labels,omitempty" json:"labels,omitempty"`
	InitialNotes []string `yaml:"initial_notes,omitempty" json:"initial_notes,omitempty"`
}

// Validate checks if the template has valid data
func (t *Template) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if strings.ContainsAny(t.Name, "/\\") {
		return fmt.Errorf("template name cannot contain path separators")
	}

	if len(t.Name) > 100 {
		return fmt.Errorf("template name must be 100 characters or less")
	}

	// Validate labels (no whitespace allowed)
	for _, label := range t.Labels {
		if strings.ContainsAny(label, " \t\n\r") {
			return fmt.Errorf("label '%s' cannot contain whitespace", label)
		}
	}

	return nil
}

// Filename returns the YAML filename for this template
func (t *Template) Filename() string {
	return t.Name + ".yaml"
}
