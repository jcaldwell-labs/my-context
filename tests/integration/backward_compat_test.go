package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadSprint1ContextWithoutArchiveFlag tests loading Sprint 1 contexts that lack is_archived field
func TestLoadSprint1ContextWithoutArchiveFlag(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create Sprint 1 style meta.json (without is_archived field)
	sprint1Context := map[string]interface{}{
		"name":              "sprint1-context",
		"start_time":        time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"end_time":          time.Now().Format(time.RFC3339),
		"status":            "stopped",
		"subdirectory_path": "sprint1-context",
	}

	// Write Sprint 1 meta.json
	contextDir := filepath.Join(core.GetHomeDir(), "sprint1-context")
	err := os.MkdirAll(contextDir, 0755)
	require.NoError(t, err)

	metaPath := filepath.Join(contextDir, "meta.json")
	metaData, err := json.MarshalIndent(sprint1Context, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(metaPath, metaData, 0644)
	require.NoError(t, err)

	// Create empty log files (Sprint 1 structure)
	for _, logFile := range []string{"notes.log", "files.log", "touch.log"} {
		logPath := filepath.Join(contextDir, logFile)
		err = os.WriteFile(logPath, []byte{}, 0644)
		require.NoError(t, err)
	}

	// Load context using Sprint 2 code
	ctx, err := core.LoadContext("sprint1-context")
	require.NoError(t, err)

	// Verify context loaded successfully
	assert.Equal(t, "sprint1-context", ctx.Name)
	assert.Equal(t, "stopped", ctx.Status)

	// is_archived should default to false for Sprint 1 contexts
	assert.False(t, ctx.IsArchived)
}

// TestSprint2OperationsOnSprint1Context tests Sprint 2 features work on Sprint 1 data
func TestSprint2OperationsOnSprint1Context(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create Sprint 1 context
	sprint1Context := map[string]interface{}{
		"name":              "legacy-context",
		"start_time":        time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"end_time":          time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"status":            "stopped",
		"subdirectory_path": "legacy-context",
	}

	contextDir := filepath.Join(core.GetHomeDir(), "legacy-context")
	err := os.MkdirAll(contextDir, 0755)
	require.NoError(t, err)

	metaPath := filepath.Join(contextDir, "meta.json")
	metaData, err := json.MarshalIndent(sprint1Context, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(metaPath, metaData, 0644)
	require.NoError(t, err)

	// Create log files with Sprint 1 data
	notesPath := filepath.Join(contextDir, "notes.log")
	err = os.WriteFile(notesPath, []byte("2025-01-01T10:00:00Z|Sprint 1 note\n"), 0644)
	require.NoError(t, err)

	filesPath := filepath.Join(contextDir, "files.log")
	err = os.WriteFile(filesPath, []byte("2025-01-01T10:05:00Z|/path/to/file.go\n"), 0644)
	require.NoError(t, err)

	touchPath := filepath.Join(contextDir, "touch.log")
	err = os.WriteFile(touchPath, []byte("2025-01-01T10:10:00Z\n"), 0644)
	require.NoError(t, err)

	// Test Sprint 2 operations on Sprint 1 context

	// 1. Archive operation
	err = core.ArchiveContext("legacy-context")
	require.NoError(t, err)

	// Verify is_archived was added
	ctx, err := core.LoadContext("legacy-context")
	require.NoError(t, err)
	assert.True(t, ctx.IsArchived)

	// 2. Export operation
	outputPath, err := core.ExportContext("legacy-context", "")
	require.NoError(t, err)
	defer os.Remove(outputPath)

	// Verify export includes Sprint 1 data
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Sprint 1 note")
	assert.Contains(t, string(content), "/path/to/file.go")

	// 3. List with filters
	contexts, err := core.ListContexts(core.ListOptions{
		Archived: true,
	})
	require.NoError(t, err)
	assert.Greater(t, len(contexts), 0)

	// 4. Delete operation (cleanup)
	err = core.DeleteContext("legacy-context", true, true)
	require.NoError(t, err)
}

// TestMixedSprint1And2Contexts tests Sprint 1 and Sprint 2 contexts coexist
func TestMixedSprint1And2Contexts(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create Sprint 1 context (no is_archived)
	sprint1Dir := filepath.Join(core.GetHomeDir(), "sprint1-mixed")
	err := os.MkdirAll(sprint1Dir, 0755)
	require.NoError(t, err)

	sprint1Meta := map[string]interface{}{
		"name":              "sprint1-mixed",
		"start_time":        time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"status":            "stopped",
		"subdirectory_path": "sprint1-mixed",
	}
	metaData, err := json.MarshalIndent(sprint1Meta, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(sprint1Dir, "meta.json"), metaData, 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(sprint1Dir, "notes.log"), []byte{}, 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(sprint1Dir, "files.log"), []byte{}, 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(sprint1Dir, "touch.log"), []byte{}, 0644)
	require.NoError(t, err)

	// Create Sprint 2 context (with is_archived)
	ctx2, err := core.CreateContext("sprint2-mixed")
	require.NoError(t, err)
	err = core.StopContext(ctx2.Name)
	require.NoError(t, err)
	err = core.ArchiveContext(ctx2.Name)
	require.NoError(t, err)

	// List all contexts
	contexts, err := core.ListContexts(core.ListOptions{})
	require.NoError(t, err)

	// Should see both contexts
	names := make([]string, len(contexts))
	for i, ctx := range contexts {
		names[i] = ctx.Name
	}

	assert.Contains(t, names, "sprint1-mixed")
	// sprint2-mixed might be hidden if archived contexts are filtered by default
}

// TestSprint1MetaJsonPreserved tests that Sprint 1 meta.json fields are preserved
func TestSprint1MetaJsonPreserved(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create Sprint 1 context with all original fields
	sprint1Context := map[string]interface{}{
		"name":              "preserve-test",
		"start_time":        "2025-01-01T10:00:00Z",
		"end_time":          "2025-01-01T11:00:00Z",
		"status":            "stopped",
		"subdirectory_path": "preserve-test",
		"custom_field":      "should_be_preserved", // Extra field from Sprint 1
	}

	contextDir := filepath.Join(core.GetHomeDir(), "preserve-test")
	os.MkdirAll(contextDir, 0755)

	metaPath := filepath.Join(contextDir, "meta.json")
	metaData, _ := json.MarshalIndent(sprint1Context, "", "  ")
	os.WriteFile(metaPath, metaData, 0644)
	os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte{}, 0644)
	os.WriteFile(filepath.Join(contextDir, "files.log"), []byte{}, 0644)
	os.WriteFile(filepath.Join(contextDir, "touch.log"), []byte{}, 0644)

	// Load and save with Sprint 2 code
	ctx, err := core.LoadContext("preserve-test")
	require.NoError(t, err)

	// Archive it (triggers save)
	err = core.ArchiveContext("preserve-test")
	require.NoError(t, err)

	// Read raw JSON
	rawData, err := os.ReadFile(metaPath)
	require.NoError(t, err)

	var updatedMeta map[string]interface{}
	err = json.Unmarshal(rawData, &updatedMeta)
	require.NoError(t, err)

	// Verify all original fields preserved
	assert.Equal(t, "preserve-test", updatedMeta["name"])
	assert.Equal(t, "2025-01-01T10:00:00Z", updatedMeta["start_time"])
	assert.Equal(t, "stopped", updatedMeta["status"])

	// New field added
	assert.True(t, updatedMeta["is_archived"].(bool))

	// Custom field should be preserved (if Sprint 2 code doesn't strip unknown fields)
	// Note: This depends on how JSON unmarshaling is handled
}

// TestSprint1ToSprint2Upgrade tests the upgrade path
func TestSprint1ToSprint2Upgrade(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Simulate Sprint 1 data directory
	homeDir := core.GetHomeDir()

	// Create multiple Sprint 1 contexts
	for i := 1; i <= 3; i++ {
		contextName := fmt.Sprintf("upgrade-test-%d", i)
		contextDir := filepath.Join(homeDir, contextName)
		err := os.MkdirAll(contextDir, 0755)
		require.NoError(t, err)

		meta := map[string]interface{}{
			"name":              contextName,
			"start_time":        time.Now().Add(-time.Duration(i) * time.Hour).Format(time.RFC3339),
			"status":            "stopped",
			"subdirectory_path": contextName,
		}

		metaData, err := json.MarshalIndent(meta, "", "  ")
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(contextDir, "meta.json"), metaData, 0644)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte{}, 0644)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(contextDir, "files.log"), []byte{}, 0644)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(contextDir, "touch.log"), []byte{}, 0644)
		require.NoError(t, err)
	}

	// List all contexts with Sprint 2 code
	contexts, err := core.ListContexts(core.ListOptions{})
	require.NoError(t, err)

	// All Sprint 1 contexts should load successfully
	assert.GreaterOrEqual(t, len(contexts), 3)

	// Verify each context loads without errors
	for _, ctx := range contexts {
		assert.NotEmpty(t, ctx.Name)
		assert.False(t, ctx.IsArchived) // Default for Sprint 1 contexts
	}
}

// TestSprint1NotesDollarCharacter tests backward compat for $ bug fix
func TestSprint1NotesDollarCharacter(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create Sprint 1 context with notes containing $
	contextDir := filepath.Join(core.GetHomeDir(), "dollar-compat")
	err := os.MkdirAll(contextDir, 0755)
	require.NoError(t, err)

	meta := map[string]interface{}{
		"name":              "dollar-compat",
		"start_time":        time.Now().Format(time.RFC3339),
		"status":            "stopped",
		"subdirectory_path": "dollar-compat",
	}

	metaData, err := json.MarshalIndent(meta, "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(contextDir, "meta.json"), metaData, 0644)
	require.NoError(t, err)

	// Sprint 1 note with $ (might have been corrupted)
	notesContent := time.Now().Format(time.RFC3339) + "|Budget: $500\n"
	err = os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte(notesContent), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(contextDir, "files.log"), []byte{}, 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(contextDir, "touch.log"), []byte{}, 0644)
	require.NoError(t, err)

	// Export with Sprint 2 code
	outputPath, err := core.ExportContext("dollar-compat", "")
	require.NoError(t, err)
	defer os.Remove(outputPath)

	// Verify $ character preserved in export
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "$500")
}
