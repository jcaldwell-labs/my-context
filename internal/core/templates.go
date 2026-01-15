package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"gopkg.in/yaml.v3"
)

// GetTemplatesDir returns the path to the templates directory
func GetTemplatesDir() string {
	return filepath.Join(GetContextHome(), "templates")
}

// EnsureTemplatesDir creates the templates directory if it doesn't exist
func EnsureTemplatesDir() error {
	templatesDir := GetTemplatesDir()
	if err := CreateDir(templatesDir); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}
	return nil
}

// ListTemplates returns all available templates sorted by name
func ListTemplates() ([]*models.Template, error) {
	templatesDir := GetTemplatesDir()

	// Check if templates directory exists
	if !FileExists(templatesDir) {
		return []*models.Template{}, nil // Return empty list if directory doesn't exist
	}

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	templates := make([]*models.Template, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .yaml files
		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		// Load template
		templateName := entry.Name()[:len(entry.Name())-5] // Remove .yaml extension
		template, err := GetTemplate(templateName)
		if err != nil {
			continue // Skip invalid templates
		}

		templates = append(templates, template)
	}

	// Sort by name
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// GetTemplate loads a template by name
func GetTemplate(name string) (*models.Template, error) {
	if name == "" {
		return nil, fmt.Errorf("template name cannot be empty")
	}

	templatePath := filepath.Join(GetTemplatesDir(), name+".yaml")

	if !FileExists(templatePath) {
		return nil, fmt.Errorf("template %q not found", name)
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	var template models.Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template file: %w", err)
	}

	// Validate template
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	return &template, nil
}

// SaveTemplate saves a template to disk
func SaveTemplate(template *models.Template) error {
	if err := template.Validate(); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	if err := EnsureTemplatesDir(); err != nil {
		return err
	}

	templatePath := filepath.Join(GetTemplatesDir(), template.Filename())

	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to serialize template: %w", err)
	}

	if err := os.WriteFile(templatePath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	return nil
}

// DeleteTemplate removes a template from disk
func DeleteTemplate(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	templatePath := filepath.Join(GetTemplatesDir(), name+".yaml")

	if !FileExists(templatePath) {
		return fmt.Errorf("template %q not found", name)
	}

	if err := os.Remove(templatePath); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}

// CreateDefaultTemplates creates default templates if none exist
func CreateDefaultTemplates() error {
	// Check if any templates already exist
	templates, err := ListTemplates()
	if err != nil {
		return err
	}

	if len(templates) > 0 {
		return nil // Templates already exist, don't overwrite
	}

	// Create default templates
	defaultTemplates := []*models.Template{
		{
			Name:        "debug",
			Description: "Template for debugging sessions",
			Labels:      []string{"debug", "investigation"},
			InitialNotes: []string{
				"Investigation started",
				"Hypothesis:",
				"Findings:",
			},
		},
		{
			Name:        "feature",
			Description: "Template for feature development",
			Labels:      []string{"feature", "development"},
			InitialNotes: []string{
				"Requirements:",
				"Implementation plan:",
			},
		},
		{
			Name:        "meeting",
			Description: "Template for meeting notes",
			Labels:      []string{"meeting"},
			InitialNotes: []string{
				"Attendees:",
				"Agenda:",
				"Notes:",
				"Action items:",
			},
		},
	}

	for _, template := range defaultTemplates {
		if err := SaveTemplate(template); err != nil {
			return fmt.Errorf("failed to create default template %q: %w", template.Name, err)
		}
	}

	return nil
}
