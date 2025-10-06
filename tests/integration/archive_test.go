package integration

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestArchiveStoppedContext tests archiving a stopped context
func TestArchiveStoppedContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create and stop a context
	ctx, err := core.CreateContext("test-archive")
	require.NoError(t, err)

	err = core.StopContext(ctx.Name)
	require.NoError(t, err)

	// Archive the context
	err = core.ArchiveContext(ctx.Name)
	require.NoError(t, err)

	// Verify is_archived flag in meta.json
	metaPath := core.GetContextMetaPath(ctx.Name)
	data, err := os.ReadFile(metaPath)
	require.NoError(t, err)

	var loadedCtx models.Context
	err = json.Unmarshal(data, &loadedCtx)
	require.NoError(t, err)

	assert.True(t, loadedCtx.IsArchived)
}

// TestArchiveActiveContext tests that archiving active context fails
func TestArchiveActiveContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create an active context (don't stop it)
	ctx, err := core.CreateContext("active-context")
	require.NoError(t, err)

	// Attempt to archive active context should fail
	err = core.ArchiveContext(ctx.Name)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "active")
}

// TestArchiveNonExistentContext tests error handling for missing context
func TestArchiveNonExistentContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to archive non-existent context
	err := core.ArchiveContext("does-not-exist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestArchiveHiddenFromDefaultList tests archived contexts hidden from list
func TestArchiveHiddenFromDefaultList(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create, stop, and archive a context
	ctx1, err := core.CreateContext("visible-context")
	require.NoError(t, err)
	core.StopContext(ctx1.Name)

	ctx2, err := core.CreateContext("archived-context")
	require.NoError(t, err)
	core.StopContext(ctx2.Name)
	core.ArchiveContext(ctx2.Name)

	// List contexts (default - should exclude archived)
	contexts, err := core.ListContexts(core.ListOptions{
		IncludeArchived: false,
	})
	require.NoError(t, err)

	// Verify archived context not in list
	names := getContextNames(contexts)
	assert.Contains(t, names, "visible-context")
	assert.NotContains(t, names, "archived-context")
}

// TestListArchivedFlag tests --archived flag shows only archived contexts
func TestListArchivedFlag(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create archived and non-archived contexts
	ctx1, _ := core.CreateContext("normal-context")
	core.StopContext(ctx1.Name)

	ctx2, _ := core.CreateContext("archived-context")
	core.StopContext(ctx2.Name)
	core.ArchiveContext(ctx2.Name)

	// List only archived contexts
	contexts, err := core.ListContexts(core.ListOptions{
		IncludeArchived: true,
		ArchivedOnly:    true,
	})
	require.NoError(t, err)

	// Verify only archived context in list
	names := getContextNames(contexts)
	assert.NotContains(t, names, "normal-context")
	assert.Contains(t, names, "archived-context")
}

// TestArchiveIdempotent tests archiving already-archived context
func TestArchiveIdempotent(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("test-idempotent")
	require.NoError(t, err)
	core.StopContext(ctx.Name)

	// Archive once
	err = core.ArchiveContext(ctx.Name)
	require.NoError(t, err)

	// Archive again - should not error (idempotent)
	err = core.ArchiveContext(ctx.Name)
	assert.NoError(t, err)
}

// TestArchivePreservesData tests that all context data remains accessible
func TestArchivePreservesData(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context with data
	ctx, err := core.CreateContext("data-preservation-test")
	require.NoError(t, err)

	core.AddNote(ctx.Name, "Important note")
	core.AddFile(ctx.Name, "/test/file.txt")
	core.AddTouch(ctx.Name)
	core.StopContext(ctx.Name)

	// Archive
	err = core.ArchiveContext(ctx.Name)
	require.NoError(t, err)

	// Verify all data still accessible
	loadedCtx, err := core.LoadContext(ctx.Name)
	require.NoError(t, err)

	notes, err := core.GetNotes(ctx.Name)
	require.NoError(t, err)
	assert.Len(t, notes, 1)
	assert.Equal(t, "Important note", notes[0].TextContent)

	files, err := core.GetFiles(ctx.Name)
	require.NoError(t, err)
	assert.Len(t, files, 1)

	touches, err := core.GetTouches(ctx.Name)
	require.NoError(t, err)
	assert.Len(t, touches, 1)
}

// Helper function to extract context names from list
func getContextNames(contexts []models.Context) []string {
	names := make([]string, len(contexts))
	for i, ctx := range contexts {
		names[i] = ctx.Name
	}
	return names
}
