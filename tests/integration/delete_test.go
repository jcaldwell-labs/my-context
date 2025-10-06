package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeleteWithConfirmationAccept tests deletion with user confirmation
func TestDeleteWithConfirmationAccept(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create and stop a context
	ctx, err := core.CreateContext("test-delete-confirm")
	require.NoError(t, err)
	core.AddNote(ctx.Name, "Test note")
	core.StopContext(ctx.Name)

	// Delete with confirmation (simulate user accepts)
	err = core.DeleteContext(ctx.Name, false, true) // confirmationAccepted=true
	require.NoError(t, err)

	// Verify directory removed
	contextDir := core.GetContextDir(ctx.Name)
	assert.NoDirExists(t, contextDir)

	// Verify context no longer in list
	contexts, _ := core.ListContexts(core.ListOptions{})
	names := getContextNames(contexts)
	assert.NotContains(t, names, ctx.Name)
}

// TestDeleteWithConfirmationCancel tests deletion when user cancels
func TestDeleteWithConfirmationCancel(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("test-delete-cancel")
	require.NoError(t, err)
	core.StopContext(ctx.Name)

	// Delete with confirmation rejected
	err = core.DeleteContext(ctx.Name, false, false) // confirmationAccepted=false
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cancelled")

	// Verify directory still exists
	contextDir := core.GetContextDir(ctx.Name)
	assert.DirExists(t, contextDir)
}

// TestDeleteWithForceFlag tests deletion with --force (no prompt)
func TestDeleteWithForceFlag(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("test-delete-force")
	require.NoError(t, err)
	core.StopContext(ctx.Name)

	// Delete with force flag (no confirmation needed)
	err = core.DeleteContext(ctx.Name, true, false) // force=true
	require.NoError(t, err)

	// Verify directory removed
	contextDir := core.GetContextDir(ctx.Name)
	assert.NoDirExists(t, contextDir)
}

// TestDeleteActiveContext tests that deleting active context fails
func TestDeleteActiveContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create an active context (don't stop it)
	ctx, err := core.CreateContext("active-delete-test")
	require.NoError(t, err)

	// Attempt to delete active context should fail
	err = core.DeleteContext(ctx.Name, true, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "active")
}

// TestDeleteNonExistentContext tests error handling for missing context
func TestDeleteNonExistentContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to delete non-existent context
	err := core.DeleteContext("does-not-exist", true, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestDeletePreservesTransitionsLog tests that transitions.log is preserved
func TestDeletePreservesTransitionsLog(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create contexts and generate transitions
	ctx1, _ := core.CreateContext("first-context")
	ctx2, _ := core.CreateContext("second-context") // Switches from first-context
	core.StopContext(ctx2.Name)

	// Verify transitions exist
	transitionsPath := filepath.Join(core.GetContextHome(), "transitions.log")
	transitionsBefore, err := os.ReadFile(transitionsPath)
	require.NoError(t, err)
	assert.Contains(t, string(transitionsBefore), "first-context")
	assert.Contains(t, string(transitionsBefore), "second-context")

	// Delete first context
	err = core.DeleteContext("first-context", true, false)
	require.NoError(t, err)

	// Verify transitions.log still exists and contains history
	transitionsAfter, err := os.ReadFile(transitionsPath)
	require.NoError(t, err)
	assert.Equal(t, string(transitionsBefore), string(transitionsAfter))

	// Historical record preserved even though context deleted
	assert.Contains(t, string(transitionsAfter), "first-context")
}

// TestDeleteRemovesAllData tests that entire context directory is removed
func TestDeleteRemovesAllData(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context with all data types
	ctx, err := core.CreateContext("comprehensive-delete")
	require.NoError(t, err)

	core.AddNote(ctx.Name, "Note 1")
	core.AddNote(ctx.Name, "Note 2")
	core.AddFile(ctx.Name, "/file1.txt")
	core.AddFile(ctx.Name, "/file2.txt")
	core.AddTouch(ctx.Name)
	core.AddTouch(ctx.Name)
	core.StopContext(ctx.Name)

	// Verify all files exist
	contextDir := core.GetContextDir(ctx.Name)
	assert.DirExists(t, contextDir)
	assert.FileExists(t, filepath.Join(contextDir, "meta.json"))
	assert.FileExists(t, filepath.Join(contextDir, "notes.log"))
	assert.FileExists(t, filepath.Join(contextDir, "files.log"))
	assert.FileExists(t, filepath.Join(contextDir, "touch.log"))

	// Delete context
	err = core.DeleteContext(ctx.Name, true, false)
	require.NoError(t, err)

	// Verify entire directory removed
	assert.NoDirExists(t, contextDir)
	assert.NoFileExists(t, filepath.Join(contextDir, "meta.json"))
	assert.NoFileExists(t, filepath.Join(contextDir, "notes.log"))
}
