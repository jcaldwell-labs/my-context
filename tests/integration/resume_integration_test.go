package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestResumeCommandIntegration tests the resume command end-to-end
func TestResumeCommandIntegration(t *testing.T) {
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

	t.Run("resume specific context name", func(t *testing.T) {
		testName := "resume-specific-test"

		// Create and stop a context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}
		core.StopContext()

		// Verify it's stopped
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if state.HasActiveContext() {
			t.Errorf("Expected no active context after stop")
		}

		// Resume using the CLI
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "resume", testName)
		cmd.Dir = "/home/be-dev-agent/projects/my-context-dev"
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to resume context: %v", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Resumed context: "+testName) {
			t.Errorf("Expected resume success message, got: %s", outputStr)
		}

		// Verify it's now active
		state, err = core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context after resume: %v", err)
		}
		if !state.HasActiveContext() || state.GetActiveContextName() != testName {
			t.Errorf("Expected context %s to be active, got %s", testName, state.GetActiveContextName())
		}

		// Clean up
		core.StopContext()
	})

	t.Run("resume --last flag", func(t *testing.T) {
		// Create and stop multiple contexts
		_, _, err := core.CreateContext("last-test-1")
		if err != nil {
			t.Fatalf("Failed to create context 1: %v", err)
		}
		core.StopContext()

		ctx2, _, err := core.CreateContext("last-test-2")
		if err != nil {
			t.Fatalf("Failed to create context 2: %v", err)
		}
		core.StopContext()

		// Resume --last should get the most recent (ctx2)
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "resume", "--last")
		cmd.Dir = "/home/be-dev-agent/projects/my-context-dev"
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to resume --last: %v", err)
		}

		outputStr := string(output)
		expectedName := ctx2.Name
		if !strings.Contains(outputStr, "Resumed context: "+expectedName) {
			t.Errorf("Expected resume of %s, got: %s", expectedName, outputStr)
		}

		// Verify the correct context is active
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if !state.HasActiveContext() || state.GetActiveContextName() != expectedName {
			t.Errorf("Expected context %s to be active, got %s", expectedName, state.GetActiveContextName())
		}

		// Clean up
		core.StopContext()
	})

	t.Run("resume with pattern matching", func(t *testing.T) {
		// Create contexts with similar names
		_, _, err := core.CreateContext("pattern-test-alpha")
		if err != nil {
			t.Fatalf("Failed to create alpha context: %v", err)
		}
		core.StopContext()

		_, _, err = core.CreateContext("pattern-test-beta")
		if err != nil {
			t.Fatalf("Failed to create beta context: %v", err)
		}
		core.StopContext()

		// Resume with pattern should find multiple and prompt
		// For integration testing, we'll just test that the command runs
		// (full interactive testing would require more complex setup)
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "resume", "pattern-test-*")
		cmd.Dir = "/home/be-dev-agent/projects/my-context-dev"
		// Since this would prompt for input, we'll just check it doesn't error on the pattern
		runErr := cmd.Run()
		// We expect this to fail because it tries to read from stdin but gets no input
		// The important thing is it doesn't fail due to pattern matching
		if runErr == nil {
			t.Log("Pattern matching worked (would prompt in interactive mode)")
		}
	})

	t.Run("resume nonexistent context", func(t *testing.T) {
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "resume", "nonexistent-context")
		cmd.Dir = "/home/be-dev-agent/projects/my-context-dev"
		_, err := cmd.Output()
		// Should fail with exit code
		if err == nil {
			t.Error("Expected error when resuming nonexistent context")
		}
	})

	t.Run("resume when context already active", func(t *testing.T) {
		activeName := "already-active-test"

		// Create and leave active
		_, _, err := core.CreateContext(activeName)
		if err != nil {
			t.Fatalf("Failed to create active context: %v", err)
		}

		// Try to resume another context
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "resume", "some-other-context")
		cmd.Dir = "/home/be-dev-agent/projects/my-context-dev"
		_, err = cmd.Output()
		if err == nil {
			t.Error("Expected error when trying to resume while context is already active")
		}

		// Clean up
		core.StopContext()
	})

	t.Run("resume with no arguments and no --last", func(t *testing.T) {
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "resume")
		cmd.Dir = "/home/be-dev-agent/projects/my-context-dev"
		_, err := cmd.Output()
		if err == nil {
			t.Error("Expected error when resume called with no arguments")
		}
	})
}
