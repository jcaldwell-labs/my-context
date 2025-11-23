package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestStopLifecycleAdvisorIntegration tests the full lifecycle advisor functionality
func TestStopLifecycleAdvisorIntegration(t *testing.T) {
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

	t.Run("lifecycle advisor summary display", func(t *testing.T) {
		testName := "advisor-summary-test"

		// Create and work on a context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add some notes
		for i := 0; i < 3; i++ {
			if _, err := core.AddNote("Test note for summary"); err != nil {
				t.Fatalf("Failed to add note: %v", err)
			}
		}

		// Stop the context and check advisor output
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}

		outputStr := string(output)

		// Check for summary display
		if !strings.Contains(outputStr, "📊 Context Summary:") {
			t.Errorf("Expected context summary in output, got: %s", outputStr)
		}

		if !strings.Contains(outputStr, "Name: "+testName) {
			t.Errorf("Expected context name in summary, got: %s", outputStr)
		}

		if !strings.Contains(outputStr, "Notes: 3") {
			t.Errorf("Expected note count in summary, got: %s", outputStr)
		}

		// Check for next steps guidance
		if !strings.Contains(outputStr, "💡 Next Steps:") {
			t.Errorf("Expected next steps guidance, got: %s", outputStr)
		}
	})

	t.Run("related context detection", func(t *testing.T) {
		// Create related contexts with same prefix
		relatedContexts := []string{
			"project-x: feature-one",
			"project-x: feature-two",
			"project-x: feature-three",
		}

		for _, name := range relatedContexts {
			_, _, err := core.CreateContext(name)
			if err != nil {
				t.Fatalf("Failed to create context %s: %v", name, err)
			}
			core.StopContext()
		}

		// Create and stop one more context from the same project
		testName := "project-x: test-feature"
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create test context: %v", err)
		}

		// Stop and check for related context suggestions
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}

		outputStr := string(output)

		// Should suggest resuming related contexts
		if !strings.Contains(outputStr, "🔄 Resume related work:") {
			t.Errorf("Expected related work suggestions, got: %s", outputStr)
		}

		// Should show up to 3 related contexts
		if !strings.Contains(outputStr, "project-x: feature-three") ||
			!strings.Contains(outputStr, "project-x: feature-two") ||
			!strings.Contains(outputStr, "project-x: feature-one") {
			t.Errorf("Expected to find 3 related project-x contexts in output: %s", outputStr)
		}
	})

	t.Run("completion keyword triggers archive suggestion", func(t *testing.T) {
		testName := "completion-test"

		// Create context and add completion note
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add note with completion keyword
		if _, err := core.AddNote("Feature implementation completed successfully"); err != nil {
			t.Fatalf("Failed to add completion note: %v", err)
		}

		// Stop and check for archive suggestion
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}

		outputStr := string(output)

		// Should suggest archiving completed context
		if !strings.Contains(outputStr, "📦 Archive completed context:") {
			t.Errorf("Expected archive suggestion for completed context, got: %s", outputStr)
		}

		if !strings.Contains(outputStr, "my-context archive \""+testName+"\"") {
			t.Errorf("Expected archive command suggestion, got: %s", outputStr)
		}
	})

	t.Run("no completion shows archive consideration", func(t *testing.T) {
		testName := "incomplete-test"

		// Create context and add incomplete note
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Add note without completion keyword
		if _, err := core.AddNote("Still working on this feature"); err != nil {
			t.Fatalf("Failed to add incomplete note: %v", err)
		}

		// Stop and check for archive consideration (not strong suggestion)
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}

		outputStr := string(output)

		// Should suggest considering archive when complete
		if !strings.Contains(outputStr, "📦 Consider archiving when complete:") {
			t.Errorf("Expected archive consideration suggestion, got: %s", outputStr)
		}
	})

	t.Run("new work suggestion always present", func(t *testing.T) {
		testName := "new-work-test"

		// Create and stop context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Stop and capture output
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context: %v", err)
		}

		outputStr := string(output)

		// The new work suggestion should always be present
		if !strings.Contains(outputStr, "✨ Start new work:") {
			t.Errorf("Expected new work suggestion to always be present, got: %s", outputStr)
		}
	})

	t.Run("lifecycle advisor with no related contexts", func(t *testing.T) {
		testName := "isolated-context"

		// Create isolated context
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create isolated context: %v", err)
		}

		// Stop and check output
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop isolated context: %v", err)
		}

		outputStr := string(output)

		// Should still show summary and guidance
		if !strings.Contains(outputStr, "📊 Context Summary:") {
			t.Errorf("Expected context summary even with no related contexts")
		}

		// Should show new work suggestion
		if !strings.Contains(outputStr, "✨ Start new work:") {
			t.Errorf("Expected new work suggestion")
		}
	})

	t.Run("lifecycle advisor respects JSON output mode", func(t *testing.T) {
		testName := "json-advisor-test"

		// Create and stop context with JSON output
		_, _, err := core.CreateContext(testName)
		if err != nil {
			t.Fatalf("Failed to create context: %v", err)
		}

		// Stop with JSON flag - should not show advisor guidance
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "stop", "--json")
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context with JSON: %v", err)
		}

		outputStr := string(output)

		// Should be valid JSON (no advisor guidance in JSON mode)
		if !strings.Contains(outputStr, `"success":true`) && !strings.Contains(outputStr, `"data":`) {
			t.Errorf("Expected valid JSON output, got: %s", outputStr)
		}

		// Should NOT contain advisor emojis (since JSON mode should be clean)
		if strings.Contains(outputStr, "📊") || strings.Contains(outputStr, "💡") {
			t.Errorf("JSON output should not contain advisor guidance, got: %s", outputStr)
		}
	})
}
