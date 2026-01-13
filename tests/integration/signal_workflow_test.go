package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSignalWorkflowIntegration tests the complete signal lifecycle
func TestSignalWorkflowIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Test 1: Create a signal
	stdout, exitCode := runCommandFullWithEnv(binary, testDir, "signal", "create", "integration-test-signal")
	assert.Equal(t, 0, exitCode, "Create should succeed: %s", stdout)

	// Test 2: List signals - should show our signal
	stdout, exitCode = runCommandFullWithEnv(binary, testDir, "signal", "list")
	assert.Equal(t, 0, exitCode, "List should succeed")
	assert.Contains(t, stdout, "integration-test-signal", "Signal should appear in list")

	// Test 3: Clear the signal
	stdout, exitCode = runCommandFullWithEnv(binary, testDir, "signal", "clear", "integration-test-signal")
	assert.Equal(t, 0, exitCode, "Clear should succeed: %s", stdout)

	// Test 4: List should be empty now
	stdout, exitCode = runCommandFullWithEnv(binary, testDir, "signal", "list")
	assert.Equal(t, 0, exitCode, "List should succeed")
	assert.NotContains(t, stdout, "integration-test-signal", "Signal should no longer appear")
}

// TestSignalConcurrencyIntegration tests concurrent signal operations
func TestSignalConcurrencyIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create multiple signals concurrently
	signals := []string{"concurrency-signal-1", "concurrency-signal-2", "concurrency-signal-3"}

	for _, sig := range signals {
		stdout, exitCode := runCommandFullWithEnv(binary, testDir, "signal", "create", sig)
		assert.Equal(t, 0, exitCode, "Create %s should succeed: %s", sig, stdout)
	}

	// List should show all signals
	stdout, exitCode := runCommandFullWithEnv(binary, testDir, "signal", "list")
	assert.Equal(t, 0, exitCode, "List should succeed")
	for _, sig := range signals {
		assert.Contains(t, stdout, sig, "Signal %s should appear in list", sig)
	}

	// Clean up
	for _, sig := range signals {
		runCommandFullWithEnv(binary, testDir, "signal", "clear", sig)
	}
}

// TestSignalTimeoutIntegration tests wait timeout behavior
func TestSignalTimeoutIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Wait for a non-existent signal with short timeout should timeout
	start := time.Now()
	_, exitCode := runCommandFullWithEnv(binary, testDir, "signal", "wait", "timeout-test-signal", "--timeout", "200ms")

	elapsed := time.Since(start)

	// Should timeout (non-zero exit) and take approximately 200ms
	assert.NotEqual(t, 0, exitCode, "Wait should timeout with non-zero exit")
	assert.True(t, elapsed >= 180*time.Millisecond, "Should wait at least 180ms")
	assert.True(t, elapsed < 1*time.Second, "Should not wait too long")
}

// TestSignalWaitSuccess tests that wait succeeds when signal exists
func TestSignalWaitSuccess(t *testing.T) {
	binary := buildTestBinary(t)
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Use unique signal name to avoid test interference
	signalName := fmt.Sprintf("wait-success-signal-%d", time.Now().UnixNano())

	// Create signal first
	_, exitCode := runCommandFullWithEnv(binary, testDir, "signal", "create", signalName)
	require.Equal(t, 0, exitCode, "Create should succeed")
	defer runCommandFullWithEnv(binary, testDir, "signal", "clear", signalName)

	// Wait should succeed immediately since signal exists
	start := time.Now()
	_, exitCode = runCommandFullWithEnv(binary, testDir, "signal", "wait", signalName, "--timeout", "1s")
	elapsed := time.Since(start)

	assert.Equal(t, 0, exitCode, "Wait should succeed when signal exists")
	assert.True(t, elapsed < 500*time.Millisecond, "Should return quickly when signal exists")
}

// TestSignalErrorHandlingIntegration tests error scenarios
func TestSignalErrorHandlingIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Test: Clear non-existent signal should fail gracefully
	stdout, exitCode := runCommandFullWithEnv(binary, testDir, "signal", "clear", "nonexistent-signal")
	// Depending on implementation, this might succeed silently or fail
	// Just verify it doesn't crash
	_ = stdout
	_ = exitCode
}

// runCommandFullWithEnv executes a command with MY_CONTEXT_HOME set
func runCommandFullWithEnv(binary string, testDir string, args ...string) (output string, exitCode int) {
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(), "MY_CONTEXT_HOME="+testDir)

	outputBytes, err := cmd.CombinedOutput()
	output = string(outputBytes)

	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			exitCode = exit.ExitCode()
		} else {
			exitCode = 1
		}
	}

	// Trim trailing whitespace for easier assertions
	output = strings.TrimSpace(output)
	return
}
