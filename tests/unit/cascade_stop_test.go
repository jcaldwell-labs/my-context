package unit

import (
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
)

// TestGetActiveChildren tests the GetActiveChildren function
func TestGetActiveChildren(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	assert.NoError(t, err)

	// Create parent context
	parentName := "parent-context"
	_, _, err = core.CreateContext(parentName)
	assert.NoError(t, err)

	// Stop parent to make children independent
	_, err = core.StopContext()
	assert.NoError(t, err)

	// Create child contexts with parent link
	// Note: Creating a new context stops the previous one, so we need to
	// manually mark them as active after creation
	child1 := "child-1"
	_, _, err = core.CreateContextWithMetadata(child1, "", parentName, nil)
	assert.NoError(t, err)

	child2 := "child-2"
	_, _, err = core.CreateContextWithMetadata(child2, "", parentName, nil)
	assert.NoError(t, err)

	// Create child3 and stop it (this one stays stopped)
	child3 := "child-3"
	_, _, err = core.CreateContextWithMetadata(child3, "", parentName, nil)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	// Manually set child1 and child2 to active (simulating multiple active children)
	// This tests the GetActiveChildren logic, not the context lifecycle
	child1Ctx, _, _, _, err := core.GetContextWithMetadata(child1)
	assert.NoError(t, err)
	child1Ctx.Status = "active"
	child1Ctx.EndTime = nil
	err = core.WriteJSON(core.GetMetaJSONPath(child1), child1Ctx)
	assert.NoError(t, err)

	child2Ctx, _, _, _, err := core.GetContextWithMetadata(child2)
	assert.NoError(t, err)
	child2Ctx.Status = "active"
	child2Ctx.EndTime = nil
	err = core.WriteJSON(core.GetMetaJSONPath(child2), child2Ctx)
	assert.NoError(t, err)

	// Get active children
	activeChildren, err := core.GetActiveChildren(parentName)
	assert.NoError(t, err)
	assert.Len(t, activeChildren, 2, "Should have 2 active children")
	assert.Contains(t, activeChildren, child1)
	assert.Contains(t, activeChildren, child2)
	assert.NotContains(t, activeChildren, child3, "Stopped child should not be in active list")
}

// TestGetAllDescendants tests the GetAllDescendants function
func TestGetAllDescendants(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	assert.NoError(t, err)

	// Create hierarchy: parent -> child1, child2 -> grandchild1, grandchild2
	parentName := "parent"
	_, _, err = core.CreateContext(parentName)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	child1 := "child-1"
	_, _, err = core.CreateContextWithMetadata(child1, "", parentName, nil)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	child2 := "child-2"
	_, _, err = core.CreateContextWithMetadata(child2, "", parentName, nil)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	grandchild1 := "grandchild-1"
	_, _, err = core.CreateContextWithMetadata(grandchild1, "", child1, nil)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	grandchild2 := "grandchild-2"
	_, _, err = core.CreateContextWithMetadata(grandchild2, "", child2, nil)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	// Get all descendants
	descendants, err := core.GetAllDescendants(parentName)
	assert.NoError(t, err)
	assert.Len(t, descendants, 4, "Should have 4 descendants (2 children + 2 grandchildren)")
	assert.Contains(t, descendants, child1)
	assert.Contains(t, descendants, child2)
	assert.Contains(t, descendants, grandchild1)
	assert.Contains(t, descendants, grandchild2)
}

// TestGetAllDescendantsWithCircular tests handling of circular dependencies
func TestGetAllDescendantsWithCircular(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	assert.NoError(t, err)

	// Create two contexts
	ctx1 := "context-1"
	_, _, err = core.CreateContext(ctx1)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	ctx2 := "context-2"
	_, _, err = core.CreateContextWithMetadata(ctx2, "", ctx1, nil)
	assert.NoError(t, err)
	_, err = core.StopContext()
	assert.NoError(t, err)

	// Manually create circular dependency (ctx1 -> ctx2 and ctx2 -> ctx1)
	// This should not happen in normal use, but we test for robustness
	err = core.SetParent(ctx1, ctx2)
	assert.NoError(t, err)

	// Get descendants should not infinite loop
	descendants, err := core.GetAllDescendants(ctx1)
	assert.NoError(t, err)
	// Should have at most both contexts (cycle detection should prevent infinite loop)
	assert.LessOrEqual(t, len(descendants), 2)
}

// TestGetActiveChildrenEmpty tests when there are no children
func TestGetActiveChildrenEmpty(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	assert.NoError(t, err)

	// Create context without children
	parentName := "lonely-parent"
	_, _, err = core.CreateContext(parentName)
	assert.NoError(t, err)

	// Get active children
	activeChildren, err := core.GetActiveChildren(parentName)
	assert.NoError(t, err)
	assert.Empty(t, activeChildren, "Should have no active children")
}

// TestGetActiveChildrenNonExistent tests with non-existent parent
func TestGetActiveChildrenNonExistent(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	assert.NoError(t, err)

	// Try to get children of non-existent parent
	activeChildren, err := core.GetActiveChildren("non-existent-parent")
	assert.NoError(t, err, "Should not error for non-existent parent")
	assert.Empty(t, activeChildren, "Should return empty list for non-existent parent")
}
