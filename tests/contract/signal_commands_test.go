package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testBinaryPath string
)

// buildSignalTestBinary builds the my-context binary for testing
func buildSignalTestBinary(t *testing.T) string {
	if testBinaryPath != "" {
		return testBinaryPath
	}

	// Find project root
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}

	// Build binary
	tmpDir, err := os.MkdirTemp("", "signal-test-*")
	require.NoError(t, err)

	binaryName := "my-context-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	testBinaryPath = filepath.Join(tmpDir, binaryName)

	cmd := exec.Command("go", "build", "-o", testBinaryPath, "./cmd/my-context/")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to build test binary")

	return testBinaryPath
}

// runSignalCommand executes a my-context command with isolated MY_CONTEXT_HOME
func runSignalCommand(binary, testDir string, args ...string) (stdout string, exitCode int) {
	cmd := exec.Command(binary, args...)

	// Set up isolated environment
	env := make([]string, 0, len(os.Environ())+1)
	for _, e := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(e), "MY_CONTEXT_HOME=") {
			env = append(env, e)
		}
	}
	env = append(env, "MY_CONTEXT_HOME="+testDir)
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	stdout = strings.TrimSpace(string(output))

	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			exitCode = exit.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return
}

// TestSignalCreateCommand tests the signal create command
func TestSignalCreateCommand(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test creating a signal
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "create", "test-contract-signal")

	assert.Equal(t, 0, exitCode, "signal create should succeed: %s", stdout)
	assert.Contains(t, stdout, "Signal 'test-contract-signal' created")
}

// TestSignalListCommand tests the signal list command
func TestSignalListCommand(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Create a signal first
	_, exitCode := runSignalCommand(binary, testDir, "signal", "create", "list-test-signal")
	require.Equal(t, 0, exitCode, "setup: signal create should succeed")

	// Test listing signals
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "list")

	assert.Equal(t, 0, exitCode, "signal list should succeed: %s", stdout)
	assert.Contains(t, stdout, "list-test-signal", "Created signal should appear in list")
}

// TestSignalListEmpty tests listing when no signals exist
func TestSignalListEmpty(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test listing signals when none exist
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "list")

	assert.Equal(t, 0, exitCode, "signal list should succeed even when empty: %s", stdout)
	assert.Contains(t, stdout, "No signals found")
}

// TestSignalWaitCommand tests the signal wait command with timeout
func TestSignalWaitCommand(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test waiting for a non-existent signal with short timeout
	start := time.Now()
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "wait", "nonexistent-signal", "--timeout", "200ms")
	duration := time.Since(start)

	// Should timeout with non-zero exit
	assert.NotEqual(t, 0, exitCode, "signal wait should timeout: %s", stdout)
	assert.Contains(t, stdout, "timeout", "Should mention timeout in output")

	// Should wait approximately the timeout duration
	assert.GreaterOrEqual(t, duration.Milliseconds(), int64(150), "Should wait at least 150ms")
	assert.Less(t, duration.Milliseconds(), int64(1000), "Should not wait too long")
}

// TestSignalWaitSuccess tests wait succeeds when signal exists
func TestSignalWaitSuccess(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Create the signal first
	_, exitCode := runSignalCommand(binary, testDir, "signal", "create", "wait-success-signal")
	require.Equal(t, 0, exitCode, "setup: signal create should succeed")

	// Wait should succeed immediately
	start := time.Now()
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "wait", "wait-success-signal", "--timeout", "1s")
	duration := time.Since(start)

	assert.Equal(t, 0, exitCode, "signal wait should succeed when signal exists: %s", stdout)
	assert.Contains(t, stdout, "already exists", "Should indicate signal already exists")
	assert.Less(t, duration.Milliseconds(), int64(500), "Should return quickly when signal exists")
}

// TestSignalClearCommand tests the signal clear command
func TestSignalClearCommand(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Create a signal first
	_, exitCode := runSignalCommand(binary, testDir, "signal", "create", "test-clear-signal")
	require.Equal(t, 0, exitCode, "setup: signal create should succeed")

	// Test clearing the signal
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "clear", "test-clear-signal")

	assert.Equal(t, 0, exitCode, "signal clear should succeed: %s", stdout)
	assert.Contains(t, stdout, "Signal 'test-clear-signal' cleared")

	// Verify signal is gone
	stdout, exitCode = runSignalCommand(binary, testDir, "signal", "list")
	assert.Equal(t, 0, exitCode)
	assert.NotContains(t, stdout, "test-clear-signal", "Cleared signal should not appear in list")
}

// TestSignalClearNonExistent tests clearing a non-existent signal
func TestSignalClearNonExistent(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test clearing a signal that doesn't exist
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "clear", "nonexistent-signal")

	assert.NotEqual(t, 0, exitCode, "signal clear should fail for non-existent signal")
	assert.Contains(t, stdout, "does not exist", "Should indicate signal doesn't exist: %s", stdout)
}

// TestSignalCommandHelp tests the signal command help
func TestSignalCommandHelp(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test signal command help
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "--help")

	assert.Equal(t, 0, exitCode, "signal --help should succeed: %s", stdout)
	assert.Contains(t, stdout, "signal", "Help should mention signal")
	assert.Contains(t, stdout, "create", "Help should mention create subcommand")
	assert.Contains(t, stdout, "list", "Help should mention list subcommand")
	assert.Contains(t, stdout, "wait", "Help should mention wait subcommand")
	assert.Contains(t, stdout, "clear", "Help should mention clear subcommand")
}

// TestSignalCreateWithSpecialChars tests signal names with special characters
func TestSignalCreateWithSpecialChars(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test creating a signal with special characters
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "create", "test_signal.with-dashes")

	assert.Equal(t, 0, exitCode, "signal create with special chars should succeed: %s", stdout)
	assert.Contains(t, stdout, "Signal 'test_signal.with-dashes' created")
}

// TestSignalWaitWithInvalidTimeout tests invalid timeout values
func TestSignalWaitWithInvalidTimeout(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Test waiting with invalid timeout
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "wait", "test-signal", "--timeout", "invalid")

	assert.NotEqual(t, 0, exitCode, "signal wait with invalid timeout should fail")
	assert.Contains(t, stdout, "invalid", "Should mention invalid duration: %s", stdout)
}

// TestSignalCreateDuplicate tests creating a duplicate signal
func TestSignalCreateDuplicate(t *testing.T) {
	binary := buildSignalTestBinary(t)
	testDir := t.TempDir()

	// Create a signal
	_, exitCode := runSignalCommand(binary, testDir, "signal", "create", "duplicate-signal")
	require.Equal(t, 0, exitCode, "first create should succeed")

	// Try to create the same signal again
	stdout, exitCode := runSignalCommand(binary, testDir, "signal", "create", "duplicate-signal")

	assert.NotEqual(t, 0, exitCode, "duplicate signal create should fail")
	assert.Contains(t, stdout, "already exists", "Should indicate signal already exists: %s", stdout)
}
