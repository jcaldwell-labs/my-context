package unit

import (
	"os"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// TestFindRelatedContexts tests the related context discovery functionality
func TestFindRelatedContexts(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	originalHome := os.Getenv("MY_CONTEXT_HOME")
	defer func() {
		os.Setenv("MY_CONTEXT_HOME", originalHome)
	}()
	os.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	if err := core.EnsureContextHome(); err != nil {
		t.Fatalf("Failed to create context home: %v", err)
	}

	// Create test contexts with different prefixes
	testContexts := []string{
		"project-a: feature-one",
		"project-a: feature-two",
		"project-a: bug-fix",
		"project-b: task-alpha",
		"standalone-context",
		"another-standalone",
	}

	for _, name := range testContexts {
		_, _, err := core.CreateContext(name)
		if err != nil {
			t.Fatalf("Failed to create context %s: %v", name, err)
		}
		core.StopContext()
	}

	tests := []struct {
		contextName    string
		expectedPrefix string
		expectedCount  int
		description    string
	}{
		{
			contextName:    "project-a: feature-one",
			expectedPrefix: "project-a",
			expectedCount:  2, // Should find "project-a: feature-two" and "project-a: bug-fix"
			description:    "project with colon prefix",
		},
		{
			contextName:    "project-b: task-alpha",
			expectedPrefix: "project-b",
			expectedCount:  0, // No other project-b contexts
			description:    "project with only one context",
		},
		{
			contextName:    "standalone-context",
			expectedPrefix: "standalone-context",
			expectedCount:  0, // No related contexts with same name
			description:    "standalone context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			related, err := core.FindRelatedContexts(tt.contextName)
			if err != nil {
				t.Fatalf("FindRelatedContexts failed: %v", err)
			}

			if len(related) != tt.expectedCount {
				t.Errorf("Expected %d related contexts, got %d", tt.expectedCount, len(related))
				for i, ctx := range related {
					t.Logf("  Related[%d]: %s", i, ctx.Name)
				}
			}

			// Verify all related contexts have the expected prefix
			for _, ctx := range related {
				prefix := extractContextPrefix(ctx.Name)
				if prefix != tt.expectedPrefix {
					t.Errorf("Related context %s has prefix %s, expected %s", ctx.Name, prefix, tt.expectedPrefix)
				}
			}
		})
	}
}

// TestExtractContextPrefix tests the prefix extraction logic
func TestExtractContextPrefix(t *testing.T) {
	tests := []struct {
		contextName string
		expected    string
		description string
	}{
		{"project: feature-name", "project", "colon with space"},
		{"my-app:phase-1", "my-app", "colon without space"},
		{"standalone", "standalone", "no colon"},
		{"multi: colon: name", "multi", "multiple colons"},
		{"", "", "empty string"},
		{":only-colon", "", "only colon prefix"},
		{"only-colon:", "only-colon", "only colon suffix"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := extractContextPrefix(tt.contextName)
			if result != tt.expected {
				t.Errorf("extractContextPrefix(%q) = %q, expected %q", tt.contextName, result, tt.expected)
			}
		})
	}
}

// TestDetectCompletion tests completion keyword detection
func TestDetectCompletion(t *testing.T) {
	// Create mock notes for testing
	createNote := func(text string) *models.Note {
		return &models.Note{
			TextContent: text,
		}
	}

	tests := []struct {
		notes       []*models.Note
		expected    bool
		description string
	}{
		{
			notes:       []*models.Note{},
			expected:    false,
			description: "no notes",
		},
		{
			notes: []*models.Note{
				createNote("Working on feature"),
			},
			expected:    false,
			description: "single note without completion",
		},
		{
			notes: []*models.Note{
				createNote("Feature is complete"),
			},
			expected:    true,
			description: "single note with completion keyword",
		},
		{
			notes: []*models.Note{
				createNote("Started work"),
				createNote("Made progress"),
				createNote("Feature is completed"),
			},
			expected:    true,
			description: "completion in recent notes",
		},
		{
			notes: []*models.Note{
				createNote("Old note without completion"),
				createNote("Another old note"),
				createNote("Very old note"),
				createNote("Even older note"),
				createNote("Ancient note"),
				createNote("Feature is done"),
			},
			expected:    true,
			description: "completion in last 5 notes",
		},
		{
			notes: []*models.Note{
				createNote("Old note: feature is complete"),
				createNote("Another old note"),
				createNote("Very old note"),
				createNote("Even older note"),
				createNote("Ancient note"),
				createNote("Recent note without completion"),
			},
			expected:    false,
			description: "completion only in notes beyond last 5",
		},
		{
			notes: []*models.Note{
				createNote("Case insensitive COMPLETED"),
			},
			expected:    true,
			description: "case insensitive matching",
		},
		{
			notes: []*models.Note{
				createNote("This feature is finished and deployed"),
			},
			expected:    true,
			description: "multiple keywords",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := detectCompletion(tt.notes)
			if result != tt.expected {
				t.Errorf("detectCompletion() = %v, expected %v", result, tt.expected)
				for i, note := range tt.notes {
					t.Logf("  Note[%d]: %s", i, note.TextContent)
				}
			}
		})
	}
}

// Helper functions (duplicated from stop.go for testing)
func extractContextPrefix(contextName string) string {
	parts := strings.SplitN(contextName, ":", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[0])
	}
	return contextName
}

func detectCompletion(notes []*models.Note) bool {
	if len(notes) == 0 {
		return false
	}

	completionKeywords := []string{
		"complete", "completed", "done", "finished", "finish",
		"closed", "resolved", "fixed", "implemented",
		"merged", "deployed", "released",
	}

	startIdx := len(notes) - 5
	if startIdx < 0 {
		startIdx = 0
	}

	for i := len(notes) - 1; i >= startIdx; i-- {
		note := notes[i]
		noteText := strings.ToLower(note.TextContent)

		for _, keyword := range completionKeywords {
			if strings.Contains(noteText, keyword) {
				return true
			}
		}
	}

	return false
}
