package integration

import (
	"os"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestNoteWarningsIntegration tests the note warning integration with core functionality
func TestNoteWarningsIntegration(t *testing.T) {
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

	t.Run("note count integration with warnings", func(t *testing.T) {
		testName := "integration-warn-test"

		// Start context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add 49 notes
		for i := 0; i < 49; i++ {
			if _, err := core.AddNote("Integration test note"); err != nil {
				t.Fatalf("Failed to add note %d: %v", i+1, err)
			}
		}

		// Check count
		count, err := core.GetNoteCount(testName)
		if err != nil {
			t.Fatalf("Failed to get note count: %v", err)
		}
		if count != 49 {
			t.Errorf("Expected 49 notes, got %d", count)
		}

		// Add 50th note
		if _, err := core.AddNote("Integration test note 50"); err != nil {
			t.Fatalf("Failed to add 50th note: %v", err)
		}

		// Check final count
		finalCount, err := core.GetNoteCount(testName)
		if err != nil {
			t.Fatalf("Failed to get final note count: %v", err)
		}
		if finalCount != 50 {
			t.Errorf("Expected 50 notes after adding 50th, got %d", finalCount)
		}

		// Clean up
		core.StopContext()
	})

	t.Run("warnings with custom thresholds", func(t *testing.T) {
		// Set custom thresholds
		os.Setenv("MC_WARN_AT", "3")
		os.Setenv("MC_WARN_AT_2", "6")
		os.Setenv("MC_WARN_AT_3", "9")
		defer func() {
			os.Unsetenv("MC_WARN_AT")
			os.Unsetenv("MC_WARN_AT_2")
			os.Unsetenv("MC_WARN_AT_3")
		}()

		testName := "custom-threshold-test"

		// Start context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add 8 notes (should reach custom thresholds)
		for i := 0; i < 8; i++ {
			if _, err := core.AddNote("Custom threshold test note"); err != nil {
				t.Fatalf("Failed to add note %d: %v", i+1, err)
			}
		}

		// Verify count
		count, err := core.GetNoteCount(testName)
		if err != nil {
			t.Fatalf("Failed to get note count: %v", err)
		}
		if count != 8 {
			t.Errorf("Expected 8 notes, got %d", count)
		}

		// Clean up
		core.StopContext()
	})

	t.Run("note warnings don't affect functionality", func(t *testing.T) {
		testName := "functional-test"

		// Start context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add many notes to test warnings don't break anything
		for i := 0; i < 60; i++ {
			note, err := core.AddNote("Functional test note")
			if err != nil {
				t.Fatalf("Failed to add note %d: %v", i+1, err)
			}
			if note.TextContent != "Functional test note" {
				t.Errorf("Note content mismatch")
			}
		}

		// Verify final count
		count, err := core.GetNoteCount(testName)
		if err != nil {
			t.Fatalf("Failed to get final count: %v", err)
		}
		if count != 60 {
			t.Errorf("Expected 60 notes, got %d", count)
		}

		// Clean up
		core.StopContext()
	})
}
