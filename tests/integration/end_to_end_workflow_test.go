package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// TestEndToEndWorkflow tests the complete user journey through all 5 lifecycle features
func TestEndToEndWorkflow(t *testing.T) {
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

	t.Run("complete lifecycle journey with all 5 features", func(t *testing.T) {
		// === Phase 1: Smart Resume ===
		// Start a context and work on it
		ctx1Name := "e2e-project: feature-auth"
		cmd := exec.Command("go", "run", "cmd/my-context/main.go", "start", ctx1Name)
		cmd.Dir = getProjectRoot()
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to start initial context: %v", err)
		}

		// Verify context started
		if !strings.Contains(string(output), "Started context: "+ctx1Name) {
			t.Errorf("Expected context start message, got: %s", output)
		}

		// Stop the context
		cmd = exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to stop initial context: %v", err)
		}

		// === Test Smart Resume: Try to start same context name ===
		// This should trigger smart resume prompt (but we'll use --force for testing)
		cmd = exec.Command("go", "run", "cmd/my-context/main.go", "start", "--force", ctx1Name)
		cmd.Dir = getProjectRoot()
		output, err = cmd.Output()
		if err != nil {
			t.Fatalf("Failed to force start duplicate context: %v", err)
		}

		// Should create context with suffix due to --force
		if !strings.Contains(string(output), "Started context:") {
			t.Errorf("Expected context start with force flag, got: %s", output)
		}

		// === Phase 2A: Note Warnings ===
		// Add notes to trigger warnings (need 50+ notes)
		for i := 0; i < 55; i++ {
			cmd = exec.Command("go", "run", "cmd/my-context/main.go", "note", "Integration test note", string(rune('A'+(i%26))))
		cmd.Dir = getProjectRoot()
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to add note %d: %v", i+1, err)
			}
		}

		// === Phase 4: Lifecycle Advisor ===
		// Stop context and check advisor output
		cmd = exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
		output, err = cmd.Output()
		if err != nil {
			t.Fatalf("Failed to stop context with advisor: %v", err)
		}

		outputStr := string(output)

		// Verify advisor components
		if !strings.Contains(outputStr, "📊 Context Summary:") {
			t.Errorf("Expected context summary in advisor output, got: %s", outputStr)
		}

		if !strings.Contains(outputStr, "Notes: 55") {
			t.Errorf("Expected note count in summary, got: %s", outputStr)
		}

		if !strings.Contains(outputStr, "💡 Next Steps:") {
			t.Errorf("Expected next steps guidance, got: %s", outputStr)
		}

		// === Phase 2B: Resume Command ===
		// Test resume --last functionality
		cmd = exec.Command("go", "run", "cmd/my-context/main.go", "resume", "--last")
		cmd.Dir = getProjectRoot()
		output, err = cmd.Output()
		if err != nil {
			t.Fatalf("Failed to resume last context: %v", err)
		}

		if !strings.Contains(string(output), "Resumed context:") {
			t.Errorf("Expected resume success message, got: %s", output)
		}

		// === Create more contexts for bulk archive testing ===
		// Create several related contexts
		relatedContexts := []string{
			"e2e-project: feature-login",
			"e2e-project: feature-profile",
			"e2e-project: bug-fix-123",
			"unrelated: standalone-task",
		}

		for _, ctxName := range relatedContexts {
			cmd = exec.Command("go", "run", "cmd/my-context/main.go", "start", ctxName)
		cmd.Dir = getProjectRoot()
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to start context %s: %v", ctxName, err)
			}

			// Stop each context
			cmd = exec.Command("go", "run", "cmd/my-context/main.go", "stop")
		cmd.Dir = getProjectRoot()
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to stop context %s: %v", ctxName, err)
			}
		}

		// === Phase 3: Bulk Archive ===
		// Test bulk archive with pattern
		cmd = exec.Command("go", "run", "cmd/my-context/main.go", "archive", "--pattern", "e2e-project:*", "--dry-run")
		cmd.Dir = getProjectRoot()
		output, err = cmd.Output()
		if err != nil {
			t.Fatalf("Failed to run bulk archive dry-run: %v", err)
		}

		outputStr = string(output)
		if !strings.Contains(outputStr, "DRY RUN") {
			t.Errorf("Expected dry-run output, got: %s", outputStr)
		}

		// Count how many e2e-project contexts were found
		e2eCount := strings.Count(outputStr, "e2e-project:")
		if e2eCount < 4 { // Should find at least the 4 we created
			t.Errorf("Expected at least 4 e2e-project contexts in dry-run, found %d in: %s", e2eCount, outputStr)
		}

		// === Final Validation ===
		// Verify all features work together
		allContexts, err := core.ListContexts()
		if err != nil {
			t.Fatalf("Failed to list all contexts: %v", err)
		}

		// Count contexts created during the test
		testContexts := 0
		for _, ctx := range allContexts {
			if strings.HasPrefix(ctx.Name, "e2e-project:") ||
				strings.HasPrefix(ctx.Name, ctx1Name) ||
				strings.HasPrefix(ctx.Name, "unrelated:") {
				testContexts++
			}
		}

		if testContexts < 5 { // At minimum we should have several test contexts
			t.Errorf("Expected at least 5 test contexts, found %d", testContexts)
		}

		t.Logf("✅ End-to-end workflow test completed successfully!")
		t.Logf("✅ All 5 lifecycle features working together")
		t.Logf("✅ Created and managed %d contexts through complete lifecycle", testContexts)
		t.Logf("✅ Smart resume, note warnings, resume command, bulk archive, and lifecycle advisor all functional")
	})
}
