package integration

import (
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListFilterByProject tests --project flag filters correctly
func TestListFilterByProject(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create contexts with project naming convention
	core.CreateContext("ps-cli: Phase 1")
	core.CreateContext("ps-cli: Phase 2")
	core.CreateContext("garden: Planning")
	core.CreateContext("Standalone")

	// Filter by project
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "ps-cli",
	})
	require.NoError(t, err)

	names := getContextNames(contexts)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "ps-cli: Phase 1")
	assert.Contains(t, names, "ps-cli: Phase 2")
	assert.NotContains(t, names, "garden: Planning")
}

// TestStartWithProjectFlag tests --project flag creates correct name
func TestStartWithProjectFlag(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Start context with project flag
	ctx, err := core.CreateContextWithProject("Phase 1 - Foundation", "ps-cli")
	require.NoError(t, err)

	// Verify combined name
	assert.Equal(t, "ps-cli: Phase 1 - Foundation", ctx.Name)

	// Verify it filters correctly
	contexts, _ := core.ListContexts(core.ListOptions{
		Project: "ps-cli",
	})
	names := getContextNames(contexts)
	assert.Contains(t, names, "ps-cli: Phase 1 - Foundation")
}

// TestProjectExtractionMultipleColons tests extraction with multiple colons
func TestProjectExtractionMultipleColons(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context with multiple colons
	core.CreateContext("project: sub: phase")

	// Should extract only text before first colon
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "project",
	})
	require.NoError(t, err)

	names := getContextNames(contexts)
	assert.Contains(t, names, "project: sub: phase")
}

// TestProjectExtractionNoColon tests contexts without colons
func TestProjectExtractionNoColon(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context without colon
	core.CreateContext("StandaloneContext")

	// Full name should be treated as project
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "StandaloneContext",
	})
	require.NoError(t, err)

	assert.Len(t, contexts, 1)
	assert.Equal(t, "StandaloneContext", contexts[0].Name)
}

// TestProjectExtractionWhitespace tests trimming of whitespace
func TestProjectExtractionWhitespace(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context with whitespace around colon
	core.CreateContext("  project  :  phase  ")

	// Should trim whitespace in project extraction
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "project",
	})
	require.NoError(t, err)
	assert.Len(t, contexts, 1)
}

// TestProjectFilterCaseInsensitive tests case-insensitive matching
func TestProjectFilterCaseInsensitive(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	core.CreateContext("PS-CLI: Phase 1")
	core.CreateContext("ps-cli: Phase 2")
	core.CreateContext("Ps-Cli: Phase 3")

	// Filter with lowercase
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "ps-cli",
	})
	require.NoError(t, err)
	assert.Len(t, contexts, 3)

	// Filter with uppercase
	contexts, err = core.ListContexts(core.ListOptions{
		Project: "PS-CLI",
	})
	require.NoError(t, err)
	assert.Len(t, contexts, 3)
}

// TestStartProjectFlagWithEmptyName tests validation
func TestStartProjectFlagWithEmptyName(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to create context with empty phase name
	_, err := core.CreateContextWithProject("", "ps-cli")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

// TestStartProjectFlagWithEmptyProject tests validation
func TestStartProjectFlagWithEmptyProject(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to create context with empty project
	_, err := core.CreateContextWithProject("Phase 1", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "project")
}

// TestProjectFilterNoMatches tests filtering with no results
func TestProjectFilterNoMatches(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	core.CreateContext("ps-cli: Phase 1")

	// Filter for non-existent project
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "nonexistent",
	})
	require.NoError(t, err) // Not an error, just no matches
	assert.Len(t, contexts, 0)
}

// TestProjectExtractionWithStartCommand tests integration
func TestProjectExtractionWithStartCommand(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context using project flag
	ctx, err := core.CreateContextWithProject("Testing Integration", "my-project")
	require.NoError(t, err)

	// Verify name format
	assert.Equal(t, "my-project: Testing Integration", ctx.Name)

	// Verify extraction works
	project := core.ExtractProjectName(ctx.Name)
	assert.Equal(t, "my-project", project)

	// Verify filtering works
	contexts, _ := core.ListContexts(core.ListOptions{
		Project: "my-project",
	})
	assert.Len(t, contexts, 1)
}
