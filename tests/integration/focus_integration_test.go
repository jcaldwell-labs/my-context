package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestFocusCommandIntegration tests the focus command end-to-end
func TestFocusCommandIntegration(t *testing.T) {
	// Build the test binary once for all subtests
	binaryPath := buildTestBinary(t)
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

	t.Run("focus from one context to another", func(t *testing.T) {
		// Create and activate context-1
		_, _, err := core.CreateContext("focus-test-1")
		if err != nil {
			t.Fatalf("Failed to create context 1: %v", err)
		}

		// Create and activate context-2 (this stops context-1)
		_, _, err = core.CreateContext("focus-test-2")
		if err != nil {
			t.Fatalf("Failed to create context 2: %v", err)
		}

		// Focus back to context-1
		cmd := exec.Command(binaryPath, "focus", "focus-test-1")
		cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to focus context: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Stopped: focus-test-2") {
			t.Errorf("Expected stopped message for context-2, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Resumed: focus-test-1") {
			t.Errorf("Expected resumed message for context-1, got: %s", outputStr)
		}

		// Verify focus-test-1 is now active
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if !state.HasActiveContext() || state.GetActiveContextName() != "focus-test-1" {
			t.Errorf("Expected context focus-test-1 to be active, got %s", state.GetActiveContextName())
		}

		// Clean up
		core.StopContext()
	})

	t.Run("focus with no active context", func(t *testing.T) {
		// Create and stop a context
		_, _, err := core.CreateContext("focus-no-active")
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}
		core.StopContext()

		// Focus should work like resume
		cmd := exec.Command(binaryPath, "focus", "focus-no-active")
		cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to focus context: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if strings.Contains(outputStr, "Stopped:") {
			t.Errorf("Should not show stopped message when no active context, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Resumed: focus-no-active") {
			t.Errorf("Expected resumed message, got: %s", outputStr)
		}

		// Clean up
		core.StopContext()
	})

	t.Run("focus to already active context", func(t *testing.T) {
		// Create and activate a context
		_, _, err := core.CreateContext("focus-already-active")
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Focus to the same context (should be no-op)
		cmd := exec.Command(binaryPath, "focus", "focus-already-active")
		cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to focus context: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "already active") {
			t.Errorf("Expected 'already active' message, got: %s", outputStr)
		}

		// Verify it's still active
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if !state.HasActiveContext() || state.GetActiveContextName() != "focus-already-active" {
			t.Errorf("Expected context to remain active")
		}

		// Clean up
		core.StopContext()
	})

	t.Run("focus to nonexistent context", func(t *testing.T) {
		// Create and activate a context
		_, _, err := core.CreateContext("focus-active-before-error")
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Try to focus on nonexistent context
		cmd := exec.Command(binaryPath, "focus", "nonexistent-context")
		cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
		_, err = cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected error when focusing nonexistent context")
		}

		// Verify original context is still active (atomic operation)
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if !state.HasActiveContext() || state.GetActiveContextName() != "focus-active-before-error" {
			t.Errorf("Expected original context to remain active after failed focus")
		}

		// Clean up
		core.StopContext()
	})

	t.Run("focus with JSON output", func(t *testing.T) {
		// Create and activate context-1
		_, _, err := core.CreateContext("focus-json-1")
		if err != nil {
			t.Fatalf("Failed to create context 1: %v", err)
		}

		// Create and activate context-2
		_, _, err = core.CreateContext("focus-json-2")
		if err != nil {
			t.Fatalf("Failed to create context 2: %v", err)
		}

		// Focus back to context-1 with JSON output
		cmd := exec.Command(binaryPath, "focus", "focus-json-1", "--json")
		cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to focus context: %v\nOutput: %s", err, string(output))
		}

		// Parse JSON output
		var result map[string]interface{}
		if err := json.Unmarshal(output, &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		// Verify JSON structure
		if result["command"] != "focus" {
			t.Errorf("Expected command 'focus', got: %v", result["command"])
		}

		data, ok := result["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data field in JSON output")
		}

		nestedData, ok := data["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected nested data field in JSON output")
		}

		if nestedData["resumed_context"] != "focus-json-1" {
			t.Errorf("Expected resumed_context 'focus-json-1', got: %v", nestedData["resumed_context"])
		}

		if nestedData["stopped_context"] != "focus-json-2" {
			t.Errorf("Expected stopped_context 'focus-json-2', got: %v", nestedData["stopped_context"])
		}

		// Clean up
		core.StopContext()
	})

	t.Run("focus using short alias 'f'", func(t *testing.T) {
		// Create and activate context-1
		_, _, err := core.CreateContext("focus-alias-1")
		if err != nil {
			t.Fatalf("Failed to create context 1: %v", err)
		}

		// Create and activate context-2
		_, _, err = core.CreateContext("focus-alias-2")
		if err != nil {
			t.Fatalf("Failed to create context 2: %v", err)
		}

		// Focus using short alias 'f'
		cmd := exec.Command(binaryPath, "f", "focus-alias-1")
		cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to focus context with alias: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Resumed: focus-alias-1") {
			t.Errorf("Expected resumed message, got: %s", outputStr)
		}

		// Verify it's active
		state, err := core.GetActiveContext()
		if err != nil {
			t.Fatalf("Failed to get active context: %v", err)
		}
		if !state.HasActiveContext() || state.GetActiveContextName() != "focus-alias-1" {
			t.Errorf("Expected context focus-alias-1 to be active")
		}

		// Clean up
		core.StopContext()
	})
}
