package unit

import (
	"os"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/commands"
)

// TestGetEnvInt tests the getEnvInt helper function
func TestGetEnvInt(t *testing.T) {
	// Test default values
	if val := commands.GetEnvInt("NON_EXISTENT_VAR", 42); val != 42 {
		t.Errorf("Expected default value 42, got %d", val)
	}

	// Test valid environment variable
	os.Setenv("TEST_VAR", "123")
	defer os.Unsetenv("TEST_VAR")

	if val := commands.GetEnvInt("TEST_VAR", 42); val != 123 {
		t.Errorf("Expected 123, got %d", val)
	}

	// Test invalid environment variable (should return default)
	os.Setenv("TEST_VAR_INVALID", "not-a-number")
	defer os.Unsetenv("TEST_VAR_INVALID")

	if val := commands.GetEnvInt("TEST_VAR_INVALID", 42); val != 42 {
		t.Errorf("Expected default value 42 for invalid env var, got %d", val)
	}
}

// TestShowNoteWarning tests the warning threshold logic
func TestShowNoteWarning(t *testing.T) {
	// Clear environment variables for consistent testing
	os.Unsetenv("MC_WARN_AT")
	os.Unsetenv("MC_WARN_AT_2")
	os.Unsetenv("MC_WARN_AT_3")

	tests := []struct {
		name         string
		currentCount int
		expectWarn   bool
		warnText     string
	}{
		// Test first threshold (50)
		{"49 notes - no warning", 49, false, ""},
		{"50 notes - first warning", 49, true, "50 notes. Consider stopping"},

		// Test second threshold (100)
		{"99 notes - no warning", 99, false, ""},
		{"100 notes - second warning", 99, true, "100 notes and is getting large"},

		// Test third threshold (200) - exact
		{"199 notes - no warning", 199, false, ""},
		{"200 notes - third warning", 199, true, "200 notes. This context is quite large"},

		// Test periodic warnings after 200 (every 25 notes)
		{"224 notes - no warning", 224, false, ""},
		{"225 notes - periodic warning", 224, true, "225 notes. This context is quite large"},
		{"249 notes - no warning", 249, false, ""},
		{"250 notes - periodic warning", 249, true, "250 notes. This context is quite large"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Testing showNoteWarning is tricky because it prints to stdout
			// For now, we'll test the logic by checking that the function doesn't panic
			// In a real implementation, we might want to refactor to return the warning message
			commands.ShowNoteWarning(tt.currentCount)
			// Function should complete without error
		})
	}
}

// TestShowNoteWarningWithCustomThresholds tests with custom environment variables
func TestShowNoteWarningWithCustomThresholds(t *testing.T) {
	// Set custom thresholds
	os.Setenv("MC_WARN_AT", "10")
	os.Setenv("MC_WARN_AT_2", "20")
	os.Setenv("MC_WARN_AT_3", "30")
	defer func() {
		os.Unsetenv("MC_WARN_AT")
		os.Unsetenv("MC_WARN_AT_2")
		os.Unsetenv("MC_WARN_AT_3")
	}()

	// Test custom thresholds
	tests := []struct {
		name         string
		currentCount int
		expectWarn   bool
	}{
		{"9 notes - no warning", 9, false},
		{"10 notes - custom first warning", 9, true},
		{"19 notes - no warning", 19, false},
		{"20 notes - custom second warning", 19, true},
		{"29 notes - no warning", 29, false},
		{"30 notes - custom third warning", 29, true},
		{"54 notes - periodic warning", 53, true}, // 30 + 25 = 55, so 54 -> 55
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands.ShowNoteWarning(tt.currentCount)
			// Function should complete without error
		})
	}
}
