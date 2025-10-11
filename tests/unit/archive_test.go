package unit

import (
	"os"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/commands"
)

// TestGetEnvInt tests the getEnvInt helper function (already tested in note_test.go)
// This is a duplicate test for completeness
func TestArchiveGetEnvInt(t *testing.T) {
	// Test default values
	if val := commands.GetEnvInt("NON_EXISTENT_VAR", 42); val != 42 {
		t.Errorf("Expected default value 42, got %d", val)
	}

	// Test valid environment variable
	os.Setenv("TEST_ARCHIVE_VAR", "123")
	defer os.Unsetenv("TEST_ARCHIVE_VAR")

	if val := commands.GetEnvInt("TEST_ARCHIVE_VAR", 42); val != 123 {
		t.Errorf("Expected 123, got %d", val)
	}

	// Test invalid environment variable (should return default)
	os.Setenv("TEST_ARCHIVE_INVALID", "not-a-number")
	defer os.Unsetenv("TEST_ARCHIVE_INVALID")

	if val := commands.GetEnvInt("TEST_ARCHIVE_INVALID", 42); val != 42 {
		t.Errorf("Expected default value 42 for invalid env var, got %d", val)
	}
}

// TestMatchesPattern tests the pattern matching functionality (copied from resume tests)
func TestArchiveMatchesPattern(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		contextName string
		expected    bool
	}{
		// Exact matches
		{"exact match", "test-context", "test-context", true},
		{"exact no match", "test-context", "other-context", false},

		// Prefix wildcards
		{"prefix match", "test-*", "test-context", true},
		{"prefix no match", "test-*", "other-context", false},
		{"prefix exact", "test-*", "test-", true},

		// Suffix wildcards
		{"suffix match", "*-context", "test-context", true},
		{"suffix no match", "*-context", "test-other", false},
		{"suffix exact", "*-context", "-context", true},

		// Multiple wildcards
		{"multiple wildcards", "test-*-context", "test-abc-context", true},
		{"multiple wildcards no match", "test-*-context", "test-abc-other", false},

		// Edge cases
		{"single wildcard", "*", "any-context", true},
		{"no wildcards", "no-wildcards", "no-wildcards", true},
		{"no wildcards no match", "no-wildcards", "different", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patternParts := strings.Split(tt.pattern, "*")
			result := commands.MatchesPattern(tt.contextName, patternParts)
			if result != tt.expected {
				t.Errorf("MatchesPattern(%q, %v) = %v, expected %v",
					tt.contextName, patternParts, result, tt.expected)
			}
		})
	}
}

// TestParseDateString tests date parsing functionality
func TestParseDateString(t *testing.T) {
	tests := []struct {
		name        string
		dateStr     string
		expectError bool
		expectedDay int
	}{
		{"valid date", "2024-01-15", false, 15},
		{"leap year", "2024-02-29", false, 29},
		{"invalid format", "2024/01/15", true, 0},
		{"invalid date", "2024-13-01", true, 0},
		{"empty string", "", true, 0},
		{"incomplete date", "2024-01", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := commands.ParseDateString(tt.dateStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for invalid date %q, got nil", tt.dateStr)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for valid date %q: %v", tt.dateStr, err)
				return
			}

			if parsed.Day() != tt.expectedDay {
				t.Errorf("Expected day %d, got %d", tt.expectedDay, parsed.Day())
			}

			// Should be set to end of day (23:59:59)
			if parsed.Hour() != 23 || parsed.Minute() != 59 || parsed.Second() != 59 {
				t.Errorf("Expected end of day (23:59:59), got %02d:%02d:%02d",
					parsed.Hour(), parsed.Minute(), parsed.Second())
			}
		})
	}
}

// TestSafetyLimitEnforcement tests the MC_BULK_LIMIT environment variable
func TestSafetyLimitEnforcement(t *testing.T) {
	// Test default limit
	os.Unsetenv("MC_BULK_LIMIT")
	limit := commands.GetEnvInt("MC_BULK_LIMIT", 100)
	if limit != 100 {
		t.Errorf("Expected default limit 100, got %d", limit)
	}

	// Test custom limit
	os.Setenv("MC_BULK_LIMIT", "50")
	defer os.Unsetenv("MC_BULK_LIMIT")

	limit = commands.GetEnvInt("MC_BULK_LIMIT", 100)
	if limit != 50 {
		t.Errorf("Expected custom limit 50, got %d", limit)
	}

	// Test invalid limit (should return default)
	os.Setenv("MC_BULK_LIMIT", "invalid")
	limit = commands.GetEnvInt("MC_BULK_LIMIT", 100)
	if limit != 100 {
		t.Errorf("Expected default limit 100 for invalid env var, got %d", limit)
	}
}
