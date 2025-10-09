package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDeleteWithConfirmationAccept tests deleting a context with user confirmation
func TestDeleteWithConfirmationAccept(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "delete-test"
	createTestContext(t, testDir, contextName)
	runCommand("stop")

	// Simulate user accepting confirmation
	err := runCommandWithInput("delete", contextName, "y\n")
	if err != nil {
		t.Fatalf("Delete command failed: %v", err)
	}

	// Verify: Context directory removed
	contextPath := filepath.Join(testDir, contextName)
	if _, err := os.Stat(contextPath); !os.IsNotExist(err) {
		t.Error("Context directory should be removed after deletion")
	}
}

// TestDeleteWithConfirmationCancel tests canceling deletion
func TestDeleteWithConfirmationCancel(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "cancel-delete"
	createTestContext(t, testDir, contextName)
	runCommand("stop")

	// Simulate user declining confirmation
	err := runCommandWithInput("delete", contextName, "n\n")

	// Verify: Command exits with code 1 (user cancellation)
	if err == nil {
		t.Error("Expected non-zero exit when user cancels")
	}

	// Verify: Context still exists
	contextPath := filepath.Join(testDir, contextName)
	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		t.Error("Context should still exist after cancellation")
	}
}

// TestDeleteWithForceFlag tests --force flag skips confirmation
func TestDeleteWithForceFlag(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "force-delete"
	createTestContext(t, testDir, contextName)
	runCommand("stop")

	// Execute: Delete with --force
	err := runCommand("delete", contextName, "--force")
	if err != nil {
		t.Fatalf("Delete --force failed: %v", err)
	}

	// Verify: Context removed without prompting
	contextPath := filepath.Join(testDir, contextName)
	if _, err := os.Stat(contextPath); !os.IsNotExist(err) {
		t.Error("Context should be removed with --force")
	}
}

// TestDeleteActiveContext tests that deleting an active context fails
func TestDeleteActiveContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "active-delete-test"
	createTestContext(t, testDir, contextName)
	// Don't stop - leave active

	// Execute: Try to delete active context
	err := runCommand("delete", contextName, "--force")

	// Verify: Command fails
	if err == nil {
		t.Fatal("Expected error when deleting active context")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "active") {
		t.Errorf("Expected error about active context, got: %v", err)
	}
}

// TestDeletePreservesTransitionsLog tests transitions.log is preserved (FR-008.7)
func TestDeletePreservesTransitionsLog(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create first context
	context1 := "context-1"
	createTestContext(t, testDir, context1)
	runCommand("stop")

	// Create second context (creates transition)
	context2 := "context-2"
	createTestContext(t, testDir, context2)
	runCommand("stop")

	// Read transitions log before deletion
	transitionsPath := filepath.Join(testDir, "transitions.log")
	beforeDelete, _ := os.ReadFile(transitionsPath)

	// Delete first context
	runCommand("delete", context1, "--force")

	// Verify: transitions.log still exists
	if _, err := os.Stat(transitionsPath); os.IsNotExist(err) {
		t.Error("transitions.log should be preserved after deletion")
	}

	// Verify: Original transitions still present
	afterDelete, _ := os.ReadFile(transitionsPath)
	if !strings.Contains(string(afterDelete), context1) {
		t.Error("Deleted context name should still appear in transitions.log history")
	}

	// Verify: Content not reduced (historical record preserved)
	if len(afterDelete) < len(beforeDelete) {
		t.Error("transitions.log content should not be reduced after deletion")
	}
}

// Helper function
func runCommandWithInput(args ...string) error {
	// Placeholder - will simulate stdin input
	return os.ErrNotExist
}
