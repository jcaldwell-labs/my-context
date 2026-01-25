package integration

import (
	"strings"
	"testing"
)

// TestResumeByIndex tests resuming a context by index number
func TestResumeByIndex(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 3 stopped contexts (createTestContext stops context after creation)
	createTestContext(t, "ctx1")
	createTestContext(t, "ctx2")
	createTestContext(t, "ctx3")

	// Execute: Resume context at index 2
	output, err := runCommandWithOutput("resume", "2")
	if err != nil {
		t.Fatalf("Resume command failed: %v", err)
	}

	// Verify: ctx2 was resumed
	if !strings.Contains(output, "Resumed context: ctx2") {
		t.Errorf("Expected to resume ctx2, got: %s", output)
	}

	// Verify: List shows ctx2 as active
	listOutput, _ := runCommandWithOutput("list")
	if !strings.Contains(listOutput, "2. ● ctx2 (active)") {
		t.Error("Expected ctx2 to be active at index 2")
	}
}

// TestResumeByIndexOutOfRange tests error handling for out of range index
func TestResumeByIndexOutOfRange(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 3 contexts (createTestContext stops context after creation)
	createTestContext(t, "ctx1")
	createTestContext(t, "ctx2")
	createTestContext(t, "ctx3")

	// Execute: Resume with index 5 (out of range)
	output, err := runCommandWithOutput("resume", "5")
	if err == nil {
		t.Error("Expected error for out of range index")
	}

	// Verify: Error message contains range information
	if !strings.Contains(output, "out of range") {
		t.Errorf("Expected 'out of range' error, got: %s", output)
	}
	if !strings.Contains(output, "1-3") {
		t.Errorf("Expected valid range 1-3 in error, got: %s", output)
	}
}

// TestResumeByIndexRejectsResumeAlreadyActive tests that resume rejects attempts to resume an already-active context
func TestResumeByIndexRejectsResumeAlreadyActive(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 2 stopped contexts
	runCommand("start", "ctx1")
	runCommand("stop")
	runCommand("start", "ctx2")
	runCommand("stop")

	// Create and keep one active
	runCommand("start", "ctx3")

	// Get the list to confirm ctx3 is at index 1
	listOutput1, _ := runCommandWithOutput("list")
	if !strings.Contains(listOutput1, "1. ● ctx3") {
		t.Fatalf("Expected ctx3 to be at index 1, got: %s", listOutput1)
	}

	// Execute: Try to resume the same active context (index 1)
	output, err := runCommandWithOutput("resume", "1")
	if err == nil {
		t.Error("Expected error when resuming already active context")
	}

	// Verify: Error message indicates context is already active
	if !strings.Contains(output, "already active") {
		t.Errorf("Expected 'already active' error, got: %s", output)
	}
	if !strings.Contains(output, "ctx3") {
		t.Errorf("Expected error to mention ctx3, got: %s", output)
	}
}

// TestResumeByIndexSwitchesContext tests that resume stops current active context
func TestResumeByIndexSwitchesContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 3 contexts with ctx3 active
	runCommand("start", "ctx1")
	runCommand("stop")
	runCommand("start", "ctx2")
	runCommand("stop")
	runCommand("start", "ctx3")

	// Execute: Resume ctx2 (index 2) while ctx3 is active
	output, err := runCommandWithOutput("resume", "2")
	if err != nil {
		t.Fatalf("Resume command failed: %v", err)
	}

	// Verify: ctx2 was resumed
	if !strings.Contains(output, "Resumed context: ctx2") {
		t.Errorf("Expected to resume ctx2, got: %s", output)
	}

	// Verify: ctx3 is now stopped and ctx2 is active
	listOutput, _ := runCommandWithOutput("list")
	if !strings.Contains(listOutput, "1. ○ ctx3 (stopped)") {
		t.Error("Expected ctx3 to be stopped")
	}
	if !strings.Contains(listOutput, "2. ● ctx2 (active)") {
		t.Error("Expected ctx2 to be active")
	}
}

// TestResumeByIndexZero tests that index 0 is invalid
func TestResumeByIndexZero(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 1 context
	runCommand("start", "ctx1")
	runCommand("stop")

	// Execute: Resume with index 0
	output, err := runCommandWithOutput("resume", "0")
	if err == nil {
		t.Error("Expected error for index 0")
	}

	// Verify: Error message indicates invalid range
	if !strings.Contains(output, "out of range") {
		t.Errorf("Expected 'out of range' error, got: %s", output)
	}
}

// TestResumeByIndexPreservesPattern tests that non-numeric arguments still work as patterns
func TestResumeByIndexPreservesPattern(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create contexts with unique names
	runCommand("start", "work-feature-a")
	runCommand("stop")
	runCommand("start", "debug-payment")
	runCommand("stop")

	// Execute: Resume by exact name (not index)
	output, err := runCommandWithOutput("resume", "work-feature-a")
	if err != nil {
		t.Fatalf("Resume command failed: %v", err)
	}

	// Verify: Correct context was resumed
	if !strings.Contains(output, "Resumed context: work-feature-a") {
		t.Errorf("Expected to resume work-feature-a, got: %s", output)
	}
}
