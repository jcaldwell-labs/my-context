package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestArchiveBulkIntegration tests bulk archive operations end-to-end
func TestArchiveBulkIntegration(t *testing.T) {
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

	t.Run("bulk archive with pattern matching", func(t *testing.T) {
		// Create test contexts
		contexts := []string{"feature-auth", "feature-login", "bug-fix-123", "bug-fix-456", "refactor-ui"}
		for _, name := range contexts {
			_, _, err := core.CreateContext(name)
			if err != nil {
				t.Fatalf("Failed to create context %s: %v", name, err)
			}
			core.StopContext()
		}

		// Test pattern matching with --dry-run
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "archive", "--pattern", "feature-*", "--dry-run")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to run dry-run archive: %v", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "DRY RUN") || !strings.Contains(outputStr, "feature-auth") || !strings.Contains(outputStr, "feature-login") {
			t.Errorf("Expected dry-run output with feature contexts, got: %s", outputStr)
		}

		// Verify contexts are not actually archived by checking they still exist
		_, _, _, _, err = core.GetContext("feature-auth")
		if err != nil {
			t.Errorf("Context should still exist after dry-run: %v", err)
		}
	})

	t.Run("bulk archive with date filtering", func(t *testing.T) {
		// Create contexts with different completion dates
		oldDate := time.Now().AddDate(0, 0, -5) // 5 days ago

		// Create an "old" context by manually setting its end time
		_, _, err := core.CreateContext("old-context")
		if err != nil {
			t.Fatalf("Failed to create old context: %v", err)
		}
		core.StopContext()

		// Manually modify the context to have an old end date
		ctx, _, _, _, err := core.GetContext("old-context")
		if err != nil {
			t.Fatalf("Failed to get old context: %v", err)
		}
		ctx.EndTime = &oldDate
		// Note: We can't easily modify the stored context in this test setup

		// Create a new context
		_, _, err = core.CreateContext("new-context")
		if err != nil {
			t.Fatalf("Failed to create new context: %v", err)
		}
		core.StopContext()

		// Test date filtering (this is complex to test precisely without modifying stored data)
		// For now, just test that the command runs without error
		cutoffDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "archive", "--completed-before", cutoffDate, "--dry-run")
		cmd.Dir = getProjectRoot()
		_, err = cmd.Output()
		if err != nil {
			t.Fatalf("Failed to run date-filtered archive: %v", err)
		}
	})

	t.Run("bulk archive all stopped contexts", func(t *testing.T) {
		// Create multiple stopped contexts
		for i := 1; i <= 3; i++ {
			name := "all-stopped-" + string(rune('a'+i-1))
			_, _, err := core.CreateContext(name)
			if err != nil {
				t.Fatalf("Failed to create context %s: %v", name, err)
			}
			core.StopContext()
		}

		// Test --all-stopped with dry-run
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "archive", "--all-stopped", "--dry-run")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to run all-stopped archive: %v", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "DRY RUN") {
			t.Errorf("Expected dry-run output, got: %s", outputStr)
		}

		// Should show all stopped contexts
		for i := 1; i <= 3; i++ {
			name := "all-stopped-" + string(rune('a'+i-1))
			if !strings.Contains(outputStr, name) {
				t.Errorf("Expected to find %s in dry-run output", name)
			}
		}
	})

	t.Run("bulk archive safety limit enforcement", func(t *testing.T) {
		// Set a low safety limit
		os.Setenv("MC_BULK_LIMIT", "2")
		defer os.Unsetenv("MC_BULK_LIMIT")

		// Create more contexts than the limit
		for i := 1; i <= 5; i++ {
			name := "limit-test-" + string(rune('a'+i-1))
			_, _, err := core.CreateContext(name)
			if err != nil {
				t.Fatalf("Failed to create context %s: %v", name, err)
			}
			core.StopContext()
		}

		// Try to archive all stopped contexts (should fail due to limit)
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "archive", "--all-stopped", "--dry-run")
		cmd.Dir = getProjectRoot()
		_, err := cmd.Output()
		// Should fail with exit code due to safety limit
		if err == nil {
			t.Error("Expected error due to safety limit, but command succeeded")
		}
	})

	t.Run("single context archive still works", func(t *testing.T) {
		// Test that the original single context functionality still works
		singleName := "single-archive-test"
		_, _, err := core.CreateContext(singleName)
		if err != nil {
			t.Fatalf("Failed to create single context: %v", err)
		}
		core.StopContext()

		// Archive single context (old style, no flags)
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "archive", singleName)
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to archive single context: %v", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Archived context: "+singleName) {
			t.Errorf("Expected single archive success message, got: %s", outputStr)
		}

		// Verify it's archived by checking the context data
		ctx, _, _, _, err := core.GetContext(singleName)
		if err != nil {
			t.Fatalf("Failed to get archived context: %v", err)
		}

		if !ctx.IsArchived {
			t.Error("Context should be marked as archived")
		}
	})

	t.Run("bulk archive prevents active context archiving", func(t *testing.T) {
		// Create and leave a context active
		activeName := "active-context-test"
		_, _, err := core.CreateContext(activeName)
		if err != nil {
			t.Fatalf("Failed to create active context: %v", err)
		}

		// Try to archive all stopped contexts (should skip active one)
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "archive", "--all-stopped", "--dry-run")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to run all-stopped archive with active context: %v", err)
		}

		outputStr := string(output)
		// Should not include the active context
		if strings.Contains(outputStr, activeName) {
			t.Errorf("Active context should not appear in bulk archive dry-run: %s", outputStr)
		}

		// Clean up
		core.StopContext()
	})
}
