package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestArchiveStoppedContext tests archiving a stopped context successfully
func TestArchiveStoppedContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "stopped-context"
	createTestContext(t, testDir, contextName)
	runCommand("stop") // Ensure context is stopped

	// Execute: Archive the context
	err := runCommand("archive", contextName)
	if err != nil {
		t.Fatalf("Archive command failed: %v", err)
	}

	// Verify: is_archived flag in meta.json
	metaPath := filepath.Join(testDir, contextName, "meta.json")
	content, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("Failed to read meta.json: %v", err)
	}

	var meta map[string]interface{}
	json.Unmarshal(content, &meta)

	if archived, ok := meta["is_archived"].(bool); !ok || !archived {
		t.Error("Expected is_archived to be true in meta.json")
	}
}

// TestArchiveActiveContext tests that archiving an active context fails
func TestArchiveActiveContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "active-context"
	createTestContext(t, testDir, contextName)
	// Don't stop - leave it active

	// Execute: Try to archive active context
	err := runCommand("archive", contextName)

	// Verify: Command fails with appropriate error
	if err == nil {
		t.Fatal("Expected error when archiving active context")
	}

	expectedError := "active context"
	if !strings.Contains(strings.ToLower(err.Error()), expectedError) {
		t.Errorf("Expected error about active context, got: %v", err)
	}
}

// TestArchiveNonExistentContext tests error handling for non-existent context
func TestArchiveNonExistentContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Execute: Try to archive non-existent context
	err := runCommand("archive", "non-existent")

	// Verify: Command fails
	if err == nil {
		t.Fatal("Expected error for non-existent context")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

// TestArchiveHidesFromDefaultList tests that archived contexts are hidden
func TestArchiveHidesFromDefaultList(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create and archive a context
	contextName := "to-be-archived"
	createTestContext(t, testDir, contextName)
	runCommand("stop")
	runCommand("archive", contextName)

	// Execute: List contexts (default, no --archived flag)
	output, err := runCommandWithOutput("list")
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	// Verify: Archived context not in output
	if strings.Contains(output, contextName) {
		t.Error("Archived context should not appear in default list")
	}
}

// TestArchivedFlagShowsArchivedContexts tests --archived flag
func TestArchivedFlagShowsArchivedContexts(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create and archive a context
	contextName := "archived-context"
	createTestContext(t, testDir, contextName)
	runCommand("stop")
	runCommand("archive", contextName)

	// Execute: List with --archived flag
	output, err := runCommandWithOutput("list", "--archived")
	if err != nil {
		t.Fatalf("List --archived failed: %v", err)
	}

	// Verify: Archived context appears in output
	if !strings.Contains(output, contextName) {
		t.Error("Archived context should appear with --archived flag")
	}
}

// TestArchivePreservesData tests that all context data is preserved
func TestArchivePreservesData(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "data-preservation-test"
	createTestContext(t, testDir, contextName)

	// Add notes and files
	runCommand("note", "Important note")
	testFile := filepath.Join(testDir, "test-file.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	runCommand("stop")

	// Get pre-archive data
	notesPath := filepath.Join(testDir, contextName, "notes.log")
	preArchiveNotes, _ := os.ReadFile(notesPath)

	// Execute: Archive
	runCommand("archive", contextName)

	// Verify: Data still exists
	postArchiveNotes, err := os.ReadFile(notesPath)
	if err != nil {
		t.Error("Notes file should still exist after archiving")
	}

	if string(preArchiveNotes) != string(postArchiveNotes) {
		t.Error("Notes content should be preserved after archiving")
	}
}

// TestArchiveAlreadyArchived tests archiving an already-archived context
func TestArchiveAlreadyArchived(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "already-archived"
	createTestContext(t, testDir, contextName)
	runCommand("stop")

	// Archive once
	runCommand("archive", contextName)

	// Execute: Try to archive again
	err := runCommand("archive", contextName)

	// Verify: Either succeeds with message or returns specific error
	// (implementation may handle idempotently)
	if err != nil {
		if !strings.Contains(err.Error(), "already archived") {
			t.Logf("Archive again returned error: %v", err)
		}
	}
	// Success case: command is idempotent
}

// Helper functions moved to helpers_test.go
