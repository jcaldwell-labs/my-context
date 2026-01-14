package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateValidation(t *testing.T) {
	tests := []struct {
		name      string
		template  models.Template
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid template",
			template: models.Template{
				Name:        "test",
				Description: "Test template",
				Labels:      []string{"label1", "label2"},
				InitialNotes: []string{
					"Note 1",
					"Note 2",
				},
			},
			wantError: false,
		},
		{
			name: "empty name",
			template: models.Template{
				Name:        "",
				Description: "Test",
			},
			wantError: true,
			errorMsg:  "template name cannot be empty",
		},
		{
			name: "name with path separator",
			template: models.Template{
				Name:        "test/template",
				Description: "Test",
			},
			wantError: true,
			errorMsg:  "template name cannot contain path separators",
		},
		{
			name: "name too long",
			template: models.Template{
				Name:        string(make([]byte, 101)),
				Description: "Test",
			},
			wantError: true,
			errorMsg:  "template name must be 100 characters or less",
		},
		{
			name: "label with whitespace",
			template: models.Template{
				Name:        "test",
				Description: "Test",
				Labels:      []string{"label with space"},
			},
			wantError: true,
			errorMsg:  "cannot contain whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.template.Validate()
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplateFilename(t *testing.T) {
	template := models.Template{
		Name: "my-template",
	}
	assert.Equal(t, "my-template.yaml", template.Filename())
}

func TestSaveAndLoadTemplate(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Create a template
	template := &models.Template{
		Name:        "test-template",
		Description: "Test description",
		Labels:      []string{"debug", "test"},
		InitialNotes: []string{
			"Initial note 1",
			"Initial note 2",
		},
	}

	// Save template
	err := core.SaveTemplate(template)
	require.NoError(t, err)

	// Verify file exists
	templatePath := filepath.Join(testDir, "templates", "test-template.yaml")
	assert.FileExists(t, templatePath)

	// Load template
	loaded, err := core.GetTemplate("test-template")
	require.NoError(t, err)
	assert.Equal(t, template.Name, loaded.Name)
	assert.Equal(t, template.Description, loaded.Description)
	assert.Equal(t, template.Labels, loaded.Labels)
	assert.Equal(t, template.InitialNotes, loaded.InitialNotes)
}

func TestGetTemplateNotFound(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Try to load non-existent template
	_, err := core.GetTemplate("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestListTemplates(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Initially no templates
	templates, err := core.ListTemplates()
	require.NoError(t, err)
	assert.Empty(t, templates)

	// Create some templates
	template1 := &models.Template{
		Name:        "alpha",
		Description: "First template",
	}
	template2 := &models.Template{
		Name:        "beta",
		Description: "Second template",
	}

	err = core.SaveTemplate(template1)
	require.NoError(t, err)
	err = core.SaveTemplate(template2)
	require.NoError(t, err)

	// List templates
	templates, err = core.ListTemplates()
	require.NoError(t, err)
	assert.Len(t, templates, 2)

	// Should be sorted by name
	assert.Equal(t, "alpha", templates[0].Name)
	assert.Equal(t, "beta", templates[1].Name)
}

func TestListTemplatesSkipsNonYAML(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Create templates directory
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	// Create a valid template
	template := &models.Template{
		Name:        "valid",
		Description: "Valid template",
	}
	err = core.SaveTemplate(template)
	require.NoError(t, err)

	// Create a non-YAML file
	nonYAMLPath := filepath.Join(templatesDir, "readme.txt")
	err = os.WriteFile(nonYAMLPath, []byte("This is not a template"), 0o600)
	require.NoError(t, err)

	// Create a subdirectory
	subDirPath := filepath.Join(templatesDir, "subdir")
	err = os.Mkdir(subDirPath, 0o755)
	require.NoError(t, err)

	// List templates
	templates, err := core.ListTemplates()
	require.NoError(t, err)
	assert.Len(t, templates, 1)
	assert.Equal(t, "valid", templates[0].Name)
}

func TestDeleteTemplate(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Create a template
	template := &models.Template{
		Name:        "to-delete",
		Description: "Will be deleted",
	}
	err := core.SaveTemplate(template)
	require.NoError(t, err)

	// Verify it exists
	_, err = core.GetTemplate("to-delete")
	require.NoError(t, err)

	// Delete it
	err = core.DeleteTemplate("to-delete")
	require.NoError(t, err)

	// Verify it's gone
	_, err = core.GetTemplate("to-delete")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteTemplateNotFound(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Try to delete non-existent template
	err := core.DeleteTemplate("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCreateDefaultTemplates(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Create default templates
	err := core.CreateDefaultTemplates()
	require.NoError(t, err)

	// Verify default templates exist
	templates, err := core.ListTemplates()
	require.NoError(t, err)
	assert.Len(t, templates, 3)

	// Check template names
	templateNames := make(map[string]bool)
	for _, t := range templates {
		templateNames[t.Name] = true
	}
	assert.True(t, templateNames["debug"])
	assert.True(t, templateNames["feature"])
	assert.True(t, templateNames["meeting"])

	// Calling again should not create duplicates
	err = core.CreateDefaultTemplates()
	require.NoError(t, err)
	templates, err = core.ListTemplates()
	require.NoError(t, err)
	assert.Len(t, templates, 3)
}

func TestSaveTemplateInvalidTemplate(t *testing.T) {
	// Setup test environment
	testDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", testDir)

	// Try to save invalid template
	template := &models.Template{
		Name: "", // Invalid: empty name
	}
	err := core.SaveTemplate(template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid template")
}
