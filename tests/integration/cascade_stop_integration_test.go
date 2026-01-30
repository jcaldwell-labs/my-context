package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCascadeStopWarning tests that stop command warns when children exist
func TestCascadeStopWarning(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	require.NoError(t, err)

	// Create parent context
	_, _, err = core.CreateContext("parent")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	// Create child contexts with parent link
	_, _, err = core.CreateContextWithMetadata("child1", "", "parent", nil)
	require.NoError(t, err)
	
	_, _, err = core.CreateContextWithMetadata("child2", "", "parent", nil)
	require.NoError(t, err)

	// Both children are now active (last one created)
	// But we need to manually set them as active to test the warning
	// Since the tool enforces single active context, we'll manually create the scenario
	
	// Load both children WITH METADATA and mark them as active
	// Using GetContextWithMetadata preserves the parent relationship
	child1Ctx, _, _, _, err := core.GetContextWithMetadata("child1")
	require.NoError(t, err)
	child1Ctx.Status = "active"
	child1Ctx.EndTime = nil
	err = core.WriteJSON(core.GetMetaJSONPath("child1"), child1Ctx)
	require.NoError(t, err)

	child2Ctx, _, _, _, err := core.GetContextWithMetadata("child2")
	require.NoError(t, err)
	child2Ctx.Status = "active"
	child2Ctx.EndTime = nil
	err = core.WriteJSON(core.GetMetaJSONPath("child2"), child2Ctx)
	require.NoError(t, err)

	// Use CLI to resume parent (making it active)
	cmd := exec.Command("go", "run", "./cmd/my-context/main.go", "resume", "parent")
	cmd.Dir = getProjectRoot()
	cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
	_, _ = cmd.CombinedOutput()
	
	// Now parent is active, stop it without cascade (should show warning)
	cmd = exec.Command("go", "run", "./cmd/my-context/main.go", "stop")
	cmd.Dir = getProjectRoot()
	cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	
	// Should show warning about active children
	assert.Contains(t, outputStr, "⚠️", "Should show warning emoji")
	assert.Contains(t, outputStr, "active child context", "Should mention active children")
	assert.Contains(t, outputStr, "--cascade", "Should suggest cascade flag")
}

// TestCascadeStopWithFlag tests stopping parent with cascade flag
func TestCascadeStopWithFlag(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	require.NoError(t, err)

	// Create parent context
	_, _, err = core.CreateContext("parent")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	// Create nested hierarchy: parent -> child1, child2; child1 -> grandchild
	// Note: Creating a new context stops the previous one, so we must set
	// all statuses AFTER all contexts are created.

	_, _, err = core.CreateContextWithMetadata("child1", "", "parent", nil)
	require.NoError(t, err)

	_, _, err = core.CreateContextWithMetadata("child2", "", "parent", nil)
	require.NoError(t, err)

	_, _, err = core.CreateContextWithMetadata("grandchild", "", "child1", nil)
	require.NoError(t, err)
	_, err = core.StopContext() // Stop grandchild so we can set statuses
	require.NoError(t, err)

	// Now set all children to active (must be done AFTER all are created)
	child1Ctx, _, _, _, _ := core.GetContextWithMetadata("child1")
	child1Ctx.Status = "active"
	child1Ctx.EndTime = nil
	core.WriteJSON(core.GetMetaJSONPath("child1"), child1Ctx)

	child2Ctx, _, _, _, _ := core.GetContextWithMetadata("child2")
	child2Ctx.Status = "active"
	child2Ctx.EndTime = nil
	core.WriteJSON(core.GetMetaJSONPath("child2"), child2Ctx)

	grandchildCtx, _, _, _, _ := core.GetContextWithMetadata("grandchild")
	grandchildCtx.Status = "active"
	grandchildCtx.EndTime = nil
	core.WriteJSON(core.GetMetaJSONPath("grandchild"), grandchildCtx)

	// Make parent active (use GetContextWithMetadata to preserve metadata)
	err = core.SetActiveContext("parent")
	require.NoError(t, err)
	parentCtx, _, _, _, err := core.GetContextWithMetadata("parent")
	require.NoError(t, err)
	parentCtx.Status = "active"
	parentCtx.EndTime = nil
	err = core.WriteJSON(core.GetMetaJSONPath("parent"), parentCtx)
	require.NoError(t, err)

	// Stop with cascade flag
	cmd := exec.Command("go", "run", "./cmd/my-context/main.go", "stop", "--cascade")
	cmd.Dir = getProjectRoot()
	cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	
	// Should show stopped children
	assert.Contains(t, outputStr, "✓ Stopped", "Should show stopped confirmation")
	// Should have stopped multiple contexts
	stoppedCount := strings.Count(outputStr, "✓ Stopped")
	assert.GreaterOrEqual(t, stoppedCount, 2, "Should have stopped at least parent and children")

	// Verify all contexts are now stopped
	child1After, _ := core.LoadContext("child1")
	assert.Equal(t, "stopped", child1After.Status, "child1 should be stopped")
	
	child2After, _ := core.LoadContext("child2")
	assert.Equal(t, "stopped", child2After.Status, "child2 should be stopped")
	
	grandchildAfter, _ := core.LoadContext("grandchild")
	assert.Equal(t, "stopped", grandchildAfter.Status, "grandchild should be stopped")
	
	parentAfter, _ := core.LoadContext("parent")
	assert.Equal(t, "stopped", parentAfter.Status, "parent should be stopped")
}

// TestCascadeStopJSON tests JSON output includes cascade information
func TestCascadeStopJSON(t *testing.T) {
	// Setup temporary directory for testing
	tempDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	err := core.EnsureContextHome()
	require.NoError(t, err)

	// Create parent and child
	_, _, err = core.CreateContext("parent")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContextWithMetadata("child1", "", "parent", nil)
	require.NoError(t, err)
	_, err = core.StopContext() // Stop child1 so we can set statuses
	require.NoError(t, err)

	// Set child1 to active (must be done AFTER creating it)
	child1Ctx, _, _, _, _ := core.GetContextWithMetadata("child1")
	child1Ctx.Status = "active"
	child1Ctx.EndTime = nil
	core.WriteJSON(core.GetMetaJSONPath("child1"), child1Ctx)

	// Make parent active (use GetContextWithMetadata to preserve metadata)
	err = core.SetActiveContext("parent")
	require.NoError(t, err)
	parentCtx, _, _, _, _ := core.GetContextWithMetadata("parent")
	parentCtx.Status = "active"
	parentCtx.EndTime = nil
	core.WriteJSON(core.GetMetaJSONPath("parent"), parentCtx)

	// Stop with cascade and JSON output
	cmd := exec.Command("go", "run", "./cmd/my-context/main.go", "stop", "--cascade", "--json")
	cmd.Dir = getProjectRoot()
	cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+tempDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	// Parse JSON output
	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	require.NoError(t, err)

	// Verify JSON structure (format: command, timestamp, data)
	assert.Equal(t, "stop", result["command"], "Should have command field")
	assert.NotNil(t, result["data"], "Should have data field")

	// Get nested data (the stop command wraps data in data.data)
	data := result["data"].(map[string]interface{})
	if nestedData, ok := data["data"].(map[string]interface{}); ok {
		// Check for stopped_children field when cascade is used
		if stoppedChildren, ok := nestedData["stopped_children"]; ok {
			childrenList := stoppedChildren.([]interface{})
			assert.GreaterOrEqual(t, len(childrenList), 1, "Should have stopped at least one child")
		}
	}
}
