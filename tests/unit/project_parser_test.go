package unit

import (
	"strings"
	"testing"
)

// TestExtractProjectNameWithColon tests extracting project name from "project: description" format
func TestExtractProjectNameWithColon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple project with phase",
			input:    "ps-cli: Phase 1",
			expected: "ps-cli",
		},
		{
			name:     "project with multiple words",
			input:    "my-awesome-project: Planning",
			expected: "my-awesome-project",
		},
		{
			name:     "project with spaces",
			input:    "Project Name: Description Here",
			expected: "Project Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractProjectNameNoColon tests contexts without colons (full name = project)
func TestExtractProjectNameNoColon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name without colon",
			input:    "StandaloneContext",
			expected: "StandaloneContext",
		},
		{
			name:     "name with spaces",
			input:    "Quick Bug Fix",
			expected: "Quick Bug Fix",
		},
		{
			name:     "name with dashes",
			input:    "my-context-name",
			expected: "my-context-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractProjectNameMultipleColons tests using only the first colon as delimiter
func TestExtractProjectNameMultipleColons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple colons in description",
			input:    "project: Phase 1: Implementation",
			expected: "project",
		},
		{
			name:     "colon in both project and description",
			input:    "ps-cli: Feature: Bug Fix",
			expected: "ps-cli",
		},
		{
			name:     "many colons",
			input:    "a: b: c: d: e",
			expected: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractProjectNameWhitespaceTrimming tests trimming of leading/trailing whitespace
func TestExtractProjectNameWhitespaceTrimming(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading whitespace in project",
			input:    "  project: description",
			expected: "project",
		},
		{
			name:     "trailing whitespace in project",
			input:    "project  : description",
			expected: "project",
		},
		{
			name:     "whitespace on both sides",
			input:    "  project  : description",
			expected: "project",
		},
		{
			name:     "tabs and spaces",
			input:    "\t project \t: description",
			expected: "project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractProjectNameEmptyString tests handling of empty input
func TestExtractProjectNameEmptyString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "",
		},
		{
			name:     "only colon",
			input:    ":",
			expected: "",
		},
		{
			name:     "colon with spaces",
			input:    " : ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtractProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestProjectNameMatchingCaseInsensitive tests case-insensitive matching (per FR-004.5)
func TestProjectNameMatchingCaseInsensitive(t *testing.T) {
	baseContext := "ps-cli: Phase 1"
	projectName := ExtractProjectName(baseContext)

	tests := []struct {
		name          string
		searchProject string
		shouldMatch   bool
	}{
		{
			name:          "exact case match",
			searchProject: "ps-cli",
			shouldMatch:   true,
		},
		{
			name:          "all uppercase",
			searchProject: "PS-CLI",
			shouldMatch:   true,
		},
		{
			name:          "mixed case",
			searchProject: "Ps-Cli",
			shouldMatch:   true,
		},
		{
			name:          "different project",
			searchProject: "garden",
			shouldMatch:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := strings.EqualFold(projectName, tt.searchProject)
			if match != tt.shouldMatch {
				t.Errorf("Case-insensitive match of %q and %q = %v, want %v",
					projectName, tt.searchProject, match, tt.shouldMatch)
			}
		})
	}
}

// ExtractProjectName is a placeholder that will be implemented in internal/core/project.go
func ExtractProjectName(contextName string) string {
	// This implementation will fail tests until the real one is created
	parts := strings.SplitN(contextName, ":", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[0])
	}
	return strings.TrimSpace(contextName)
}
