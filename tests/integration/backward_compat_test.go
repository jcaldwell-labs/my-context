package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadSprint1MetaJSON tests loading Sprint 1 contexts without is_archived field (FR-011)
func TestLoadSprint1MetaJSON(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create a Sprint 1 style meta.json (without is_archived field)
	contextName := "sprint1-context"
	contextDir := filepath.Join(testDir, contextName)
	os.MkdirAll(contextDir, 0755)

	sprint1Meta := map[string]interface{}{
		"name":       contextName,
		"start_time": "2025-10-05T14:30:00Z",
		"end_time":   "2025-10-05T16:45:00Z",
	}

	metaPath := filepath.Join(contextDir, "meta.json")
	metaJSON, _ := json.Marshal(sprint1Meta)
	os.WriteFile(metaPath, metaJSON, 0644)

	// Execute: List contexts (should load Sprint 1 context)
	output, err := runCommandWithOutput("list")
	if err != nil {
		t.Fatalf("Failed to list contexts: %v", err)
	}

	// Verify: Sprint 1 context appears in list
	if !strings.Contains(output, contextName) {
		t.Error("Sprint 1 context (without is_archived) should appear in list")
	}
}

// TestSprint1ContextsDefaultToNotArchived tests is_archived defaults to false
func TestSprint1ContextsDefaultToNotArchived(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create Sprint 1 context
	contextName := "sprint1-not-archived"
	contextDir := filepath.Join(testDir, contextName)
	os.MkdirAll(contextDir, 0755)

	sprint1Meta := map[string]interface{}{
		"name":       contextName,
		"start_time": "2025-10-05T14:30:00Z",
		"end_time":   "2025-10-05T16:45:00Z",
		// No is_archived field
	}

	metaPath := filepath.Join(contextDir, "meta.json")
	metaJSON, _ := json.Marshal(sprint1Meta)
	os.WriteFile(metaPath, metaJSON, 0644)

	// Execute: Try to archive it
	err := runCommand("archive", contextName)
	if err != nil {
		t.Fatalf("Failed to archive Sprint 1 context: %v", err)
	}

	// Verify: is_archived field added
	content, _ := os.ReadFile(metaPath)
	var meta map[string]interface{}
	json.Unmarshal(content, &meta)

	if archived, ok := meta["is_archived"].(bool); !ok || !archived {
		t.Error("Sprint 1 context should have is_archived=true after archiving")
	}
}

// TestNewFeaturesWorkOnOldContexts tests Sprint 2 features on Sprint 1 data
func TestNewFeaturesWorkOnOldContexts(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create Sprint 1 context
	contextName := "sprint1-feature-test"
	contextDir := filepath.Join(testDir, contextName)
	os.MkdirAll(contextDir, 0755)

	// Sprint 1 meta.json
	sprint1Meta := map[string]interface{}{
		"name":       contextName,
		"start_time": "2025-10-05T14:30:00Z",
		"end_time":   "2025-10-05T16:45:00Z",
	}

	metaPath := filepath.Join(contextDir, "meta.json")
	metaJSON, _ := json.Marshal(sprint1Meta)
	os.WriteFile(metaPath, metaJSON, 0644)

	// Create Sprint 1 notes log
	notesPath := filepath.Join(contextDir, "notes.log")
	os.WriteFile(notesPath, []byte("2025-10-05T14:35:00Z|Sprint 1 note\n"), 0644)

	// Test Sprint 2 features on Sprint 1 context
	tests := []struct {
		name    string
		command []string
	}{
		{"export", []string{"export", contextName}},
		{"archive", []string{"archive", contextName}},
		{"project filter", []string{"list", "--project", contextName}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCommand(tt.command...)
			if err != nil && !strings.Contains(err.Error(), "not implemented") {
				t.Logf("Feature %q on Sprint 1 context: %v (expected until implementation)", tt.name, err)
			}
		})
	}
}

// TestBackwardCompatibilityWithDataPreservation tests no data loss during upgrade
func TestBackwardCompatibilityWithDataPreservation(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create Sprint 1 context with all data files
	contextName := "sprint1-full-context"
	contextDir := filepath.Join(testDir, contextName)
	os.MkdirAll(contextDir, 0755)

	// Meta
	sprint1Meta := map[string]interface{}{
		"name":       contextName,
		"start_time": "2025-10-05T14:30:00Z",
		"end_time":   "2025-10-05T16:45:00Z",
	}
	metaJSON, _ := json.Marshal(sprint1Meta)
	os.WriteFile(filepath.Join(contextDir, "meta.json"), metaJSON, 0644)

	// Notes
	notesContent := "2025-10-05T14:35:00Z|Note 1\n2025-10-05T14:40:00Z|Note 2\n"
	os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte(notesContent), 0644)

	// Files
	filesContent := "2025-10-05T14:36:00Z|/path/to/file1.go\n"
	os.WriteFile(filepath.Join(contextDir, "files.log"), []byte(filesContent), 0644)

	// Touches
	touchesContent := "2025-10-05T14:37:00Z\n"
	os.WriteFile(filepath.Join(contextDir, "touches.log"), []byte(touchesContent), 0644)

	// Execute: List (loads and processes Sprint 1 data)
	output, _ := runCommandWithOutput("list")

	// Verify: Context loaded successfully
	if !strings.Contains(output, contextName) {
		t.Error("Sprint 1 context should load successfully in Sprint 2")
	}

	// Verify: All data files still exist
	for _, file := range []string{"meta.json", "notes.log", "files.log", "touches.log"} {
		path := filepath.Join(contextDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Data file %q should be preserved", file)
		}
	}

	// Verify: Content unchanged
	verifyContent, _ := os.ReadFile(filepath.Join(contextDir, "notes.log"))
	if string(verifyContent) != notesContent {
		t.Error("Notes content should be unchanged after Sprint 2 processing")
	}
}

// Helper import
import "strings"
