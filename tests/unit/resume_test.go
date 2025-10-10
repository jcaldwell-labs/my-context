package unit

import (
	"os"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestGetMostRecentStopped tests the GetMostRecentStopped function
func TestGetMostRecentStopped(t *testing.T) {
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

	// Test with no contexts
	_, err := core.GetMostRecentStopped()
	if err == nil {
		t.Error("Expected error when no stopped contexts exist, got nil")
	}

	// Create and stop contexts one by one
	_, _, err = core.CreateContext("context-1")
	if err != nil {
		t.Fatalf("Failed to create context-1: %v", err)
	}
	core.StopContext() // stop context-1
	time.Sleep(10 * time.Millisecond)

	_, _, err = core.CreateContext("context-2")
	if err != nil {
		t.Fatalf("Failed to create context-2: %v", err)
	}
	core.StopContext() // stop context-2
	time.Sleep(10 * time.Millisecond)

	_, _, err = core.CreateContext("context-3")
	if err != nil {
		t.Fatalf("Failed to create context-3: %v", err)
	}
	core.StopContext() // stop context-3 last (most recent)

	// Debug: Check all contexts and their end times
	allContexts, _ := core.ListContexts()
	t.Logf("All contexts:")
	for _, ctx := range allContexts {
		if ctx.Status == "stopped" && ctx.EndTime != nil {
			t.Logf("  %s: stopped at %v", ctx.Name, *ctx.EndTime)
		}
	}

	// Test GetMostRecentStopped returns the most recently stopped context
	recent, err := core.GetMostRecentStopped()
	if err != nil {
		t.Fatalf("Failed to get most recent stopped: %v", err)
	}

	t.Logf("Most recent stopped: %s", recent.Name)
	if recent.Name != "context-3" {
		t.Errorf("Expected most recent stopped to be context-3, got %s", recent.Name)
	}

	// Verify we can get the most recent stopped context
	recent2, err := core.GetMostRecentStopped()
	if err != nil {
		t.Fatalf("Failed to get most recent stopped on second call: %v", err)
	}

	if recent2.Name != recent.Name {
		t.Errorf("Expected consistent results, got %s then %s", recent.Name, recent2.Name)
	}
}

// TestFindContextsByPattern tests the pattern matching functionality
func TestFindContextsByPattern(t *testing.T) {
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

	// Create some test contexts
	contexts := []string{
		"feature-auth",
		"feature-login",
		"bug-fix-123",
		"bug-fix-456",
		"refactor-ui",
		"test-coverage",
	}

	for _, name := range contexts {
		_, _, err := core.CreateContext(name)
		if err != nil {
			t.Fatalf("Failed to create context %s: %v", name, err)
		}
		core.StopContext()
	}

	tests := []struct {
		pattern         string
		expectedCount   int
		expectedMatches []string
	}{
		{"", 6, contexts}, // Empty pattern returns all
		{"feature-auth", 1, []string{"feature-auth"}},
		{"feature-*", 2, []string{"feature-auth", "feature-login"}},
		{"bug-*", 2, []string{"bug-fix-123", "bug-fix-456"}},
		{"*-ui", 1, []string{"refactor-ui"}},
		{"test-*", 1, []string{"test-coverage"}},
		{"nonexistent", 0, []string{}},
		{"*-fix-*", 2, []string{"bug-fix-123", "bug-fix-456"}},
	}

	for _, tt := range tests {
		t.Run("pattern_"+tt.pattern, func(t *testing.T) {
			matches, err := core.FindContextsByPattern(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to find contexts by pattern %q: %v", tt.pattern, err)
			}

			if len(matches) != tt.expectedCount {
				t.Errorf("Pattern %q: expected %d matches, got %d", tt.pattern, tt.expectedCount, len(matches))
			}

			// Check that expected matches are present
			matchNames := make([]string, len(matches))
			for i, ctx := range matches {
				matchNames[i] = ctx.Name
			}

			for _, expected := range tt.expectedMatches {
				found := false
				for _, match := range matchNames {
					if match == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Pattern %q: expected to find %q in matches %v", tt.pattern, expected, matchNames)
				}
			}
		})
	}
}

// TestMatchesPattern tests the internal pattern matching logic
func TestMatchesPattern(t *testing.T) {
	// We can't directly test the private matchesPattern function,
	// but we can test it indirectly through FindContextsByPattern
	// The unit tests above already cover this functionality
	t.Skip("Pattern matching is tested indirectly through FindContextsByPattern")
}
