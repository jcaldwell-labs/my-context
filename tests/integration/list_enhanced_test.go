package integration

import (
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListDefaultLimit tests default list shows only 10 most recent
func TestListDefaultLimit(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create 15 contexts
	for i := 1; i <= 15; i++ {
		ctx, err := core.CreateContext(fmt.Sprintf("context-%d", i))
		require.NoError(t, err)
		core.StopContext(ctx.Name)
	}

	// List with default options
	contexts, err := core.ListContexts(core.ListOptions{
		Limit: 10,
	})
	require.NoError(t, err)

	// Should return only 10 contexts
	assert.Len(t, contexts, 10)

	// Should be newest first (context-15 through context-6)
	assert.Equal(t, "context-15", contexts[0].Name)
}

// TestListAllFlag tests --all flag shows all contexts
func TestListAllFlag(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create 15 contexts
	for i := 1; i <= 15; i++ {
		ctx, _ := core.CreateContext(fmt.Sprintf("context-%d", i))
		core.StopContext(ctx.Name)
	}

	// List with --all
	contexts, err := core.ListContexts(core.ListOptions{
		ShowAll: true,
	})
	require.NoError(t, err)

	// Should return all 15 contexts
	assert.Len(t, contexts, 15)
}

// TestListCustomLimit tests --limit flag with custom value
func TestListCustomLimit(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	for i := 1; i <= 10; i++ {
		ctx, _ := core.CreateContext(fmt.Sprintf("context-%d", i))
		core.StopContext(ctx.Name)
	}

	// List with custom limit
	contexts, err := core.ListContexts(core.ListOptions{
		Limit: 5,
	})
	require.NoError(t, err)

	assert.Len(t, contexts, 5)
	assert.Equal(t, "context-10", contexts[0].Name) // Newest first
}

// TestListSearchFilter tests --search flag for substring matching
func TestListSearchFilter(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create contexts with different names
	core.CreateContext("bug-fix-123")
	core.CreateContext("feature-work")
	core.CreateContext("bug-fix-456")
	core.CreateContext("testing")

	// Search for "bug fix" (case-insensitive substring)
	contexts, err := core.ListContexts(core.ListOptions{
		Search: "bug-fix",
	})
	require.NoError(t, err)

	names := getContextNames(contexts)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "bug-fix-123")
	assert.Contains(t, names, "bug-fix-456")
	assert.NotContains(t, names, "feature-work")
}

// TestListArchivedOnly tests --archived flag
func TestListArchivedOnly(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create regular and archived contexts
	ctx1, _ := core.CreateContext("normal")
	core.StopContext(ctx1.Name)

	ctx2, _ := core.CreateContext("archived")
	core.StopContext(ctx2.Name)
	core.ArchiveContext(ctx2.Name)

	// List archived only
	contexts, err := core.ListContexts(core.ListOptions{
		IncludeArchived: true,
		ArchivedOnly:    true,
	})
	require.NoError(t, err)

	names := getContextNames(contexts)
	assert.Len(t, names, 1)
	assert.Contains(t, names, "archived")
}

// TestListActiveOnly tests --active-only flag
func TestListActiveOnly(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create stopped and active contexts
	ctx1, _ := core.CreateContext("stopped")
	core.StopContext(ctx1.Name)

	ctx2, _ := core.CreateContext("active")
	// Don't stop ctx2

	// List active only
	contexts, err := core.ListContexts(core.ListOptions{
		ActiveOnly: true,
	})
	require.NoError(t, err)

	assert.Len(t, contexts, 1)
	assert.Equal(t, "active", contexts[0].Name)
	assert.Equal(t, "active", contexts[0].Status)
}

// TestListCombinedFilters tests multiple filters working together
func TestListCombinedFilters(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create contexts with project naming
	core.CreateContext("ps-cli: Phase 1")
	core.CreateContext("ps-cli: Phase 2")
	core.CreateContext("ps-cli: Testing")
	core.CreateContext("garden: Planning")

	// All stopped
	core.StopContext("ps-cli: Phase 1")
	core.StopContext("ps-cli: Phase 2")
	core.StopContext("ps-cli: Testing")
	core.StopContext("garden: Planning")

	// Filter by project + search
	contexts, err := core.ListContexts(core.ListOptions{
		Project: "ps-cli",
		Search:  "Phase",
		Limit:   10,
	})
	require.NoError(t, err)

	names := getContextNames(contexts)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "ps-cli: Phase 1")
	assert.Contains(t, names, "ps-cli: Phase 2")
	assert.NotContains(t, names, "ps-cli: Testing")
	assert.NotContains(t, names, "garden: Planning")
}

// TestListWithInvalidLimit tests error handling for invalid limit
func TestListWithInvalidLimit(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to list with negative limit
	_, err := core.ListContexts(core.ListOptions{
		Limit: -5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "limit")
}

// TestListConflictingFlags tests --archived and --active-only together
func TestListConflictingFlags(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to use conflicting flags
	_, err := core.ListContexts(core.ListOptions{
		ArchivedOnly: true,
		ActiveOnly:   true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot use")
}

// TestListNoMatches tests behavior when no contexts match filters
func TestListNoMatches(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	core.CreateContext("test-context")

	// Search for non-existent term
	contexts, err := core.ListContexts(core.ListOptions{
		Search: "nonexistent-term",
	})
	require.NoError(t, err) // Not an error, just empty result
	assert.Len(t, contexts, 0)
}

// TestListNewestFirst tests sorting by start time descending
func TestListNewestFirst(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create contexts with delays to ensure different timestamps
	ctx1, _ := core.CreateContext("first")
	core.StopContext(ctx1.Name)

	// Small delay
	time.Sleep(100 * time.Millisecond)

	ctx2, _ := core.CreateContext("second")
	core.StopContext(ctx2.Name)

	time.Sleep(100 * time.Millisecond)

	ctx3, _ := core.CreateContext("third")
	core.StopContext(ctx3.Name)

	// List all
	contexts, err := core.ListContexts(core.ListOptions{
		ShowAll: true,
	})
	require.NoError(t, err)

	// Verify newest first
	assert.Equal(t, "third", contexts[0].Name)
	assert.Equal(t, "second", contexts[1].Name)
	assert.Equal(t, "first", contexts[2].Name)
}
