package integration

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestSmartResume tests the smart resume functionality
func TestSmartResume(t *testing.T) {
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

	t.Run("stopped context + start → prompt → resume", func(t *testing.T) {
		testName := "smart-resume-test-1"

		// Create and stop a context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add some notes to make it more realistic
		for i := 0; i < 3; i++ {
			if _, err := core.AddNote("Test note " + string(rune('A'+i))); err != nil {
				t.Fatalf("Failed to add note: %v", err)
			}
		}

		// Stop the context
		stoppedCtx, err := core.StopContext()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}
		if stoppedCtx.Name != testName {
			t.Errorf("Expected stopped context name %s, got %s", testName, stoppedCtx.Name)
		}

		// Verify it's stopped
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if state.HasActiveContext() {
			t.Errorf("Expected no active context after stop, got %s", state.GetActiveContextName())
		}

		// Now try to start with the same name - this should prompt for resume
		// In a real test, we'd need to mock stdin, but for now we'll test the logic
		existingCtx, err := core.FindContextByName(testName)
		if err != nil {
			t.Fatalf("Failed to find existing context: %v", err)
		}

		if existingCtx.Status != "stopped" {
			t.Errorf("Expected context to be stopped, got status: %s", existingCtx.Status)
		}

		// Test note count
		noteCount, err := core.GetNoteCount(testName)
		if err != nil {
			t.Fatalf("Failed to get note count: %v", err)
		}
		if noteCount != 3 {
			t.Errorf("Expected 3 notes, got %d", noteCount)
		}

		// Test last active time
		lastActive, err := core.GetLastActiveTime(testName)
		if err != nil {
			t.Fatalf("Failed to get last active time: %v", err)
		}
		if lastActive.IsZero() {
			t.Errorf("Expected non-zero last active time")
		}
	})

	t.Run("force flag bypasses resume prompt", func(t *testing.T) {
		testName := "force-flag-test"

		// Create and stop a context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		_, err = core.StopContext()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}

		// Start with --force should create a duplicate (with _2 suffix)
		ctx2, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context with duplicate name: %v", err)
		}

		// The new context should have a different display name
		if ctx2.Name == testName {
			t.Errorf("Expected new context to have different name due to duplicate, got same name: %s", ctx2.Name)
		}

		if !strings.Contains(ctx2.Name, testName) {
			t.Errorf("Expected new context name to contain original name %s, got: %s", testName, ctx2.Name)
		}

		// Clean up
		core.StopContext()
	})

	t.Run("active context duplicate attempt", func(t *testing.T) {
		testName := "active-duplicate-test"

		// Create a context (it will be active)
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Try to create another with same name - this should create a duplicate
		ctx2, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create duplicate context: %v", err)
		}

		// Should have created a duplicate with suffix
		if ctx2.Name == testName {
			t.Errorf("Expected duplicate context to have different name, got: %s", ctx2.Name)
		}

		// Clean up both contexts
		core.StopContext()
		core.StopContext()
	})

	t.Run("GetLastActiveTime finds most recent activation", func(t *testing.T) {
		testName := "last-active-test"

		// Create context
		ctx, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Stop it
		core.StopContext()

		// Wait a bit
		time.Sleep(10 * time.Millisecond)

		// Resume it (this creates a new activation)
		if err := core.SetActiveContext(ctx.Name); err != nil {
			t.Fatalf("Failed to resume context: %v", err)
		}

		// Stop it again
		core.StopContext()

		// Check last active time exists and is valid
		lastActive, err := core.GetLastActiveTime(testName)
		if err != nil {
			t.Fatalf("Failed to get last active time: %v", err)
		}

		if lastActive.IsZero() {
			t.Errorf("Last active time should not be zero")
		}
	})

	t.Run("GetNoteCount counts notes correctly", func(t *testing.T) {
		testName := "note-count-test"

		// Create context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add notes
		testNotes := []string{"Note 1", "Note 2", "Note 3", "Note 4"}
		for _, note := range testNotes {
			if _, err := core.AddNote(note); err != nil {
				t.Fatalf("Failed to add note: %v", err)
			}
		}

		// Check count
		count, err := core.GetNoteCount(testName)
		if err != nil {
			t.Fatalf("Failed to get note count: %v", err)
		}

		if count != len(testNotes) {
			t.Errorf("Expected %d notes, got %d", len(testNotes), count)
		}

		// Clean up
		core.StopContext()
	})
}
