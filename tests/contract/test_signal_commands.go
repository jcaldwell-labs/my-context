package contract

import (
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildTestBinary builds the my-context binary for testing
func buildTestBinary(t *testing.T) string {
	cmd := exec.Command("go", "build", "-o", "my-context-test", ".")
	err := cmd.Run()
	require.NoError(t, err, "Failed to build test binary")
	return "./my-context-test"
}

// runCommand executes a my-context command and returns stdout, stderr, and exit code
func runCommand(binary string, args ...string) (string, string, int) {
	cmd := exec.Command(binary, args...)
	stdout, stdoutErr := cmd.StdoutPipe()
	if stdoutErr != nil {
		return "", "Failed to create stdout pipe: " + stdoutErr.Error(), 1
	}

	stderr, stderrErr := cmd.StderrPipe()
	if stderrErr != nil {
		return "", "Failed to create stderr pipe: " + stderrErr.Error(), 1
	}

	// Start the command
	err := cmd.Start()
	if err != nil {
		return "", err.Error(), 1
	}

	// Read output
	outBytes, outErr := io.ReadAll(stdout)
	errBytes, errErr := io.ReadAll(stderr)

	// Wait for completion
	exitErr := cmd.Wait()

	stdoutStr := string(outBytes)
	stderrStr := string(errBytes)

	exitCode := 0
	if exitErr != nil {
		if exit, ok := exitErr.(*exec.ExitError); ok {
			exitCode = exit.ExitCode()
		} else {
			exitCode = 1
		}
	}

	// Include any read errors in stderr
	if outErr != nil {
		stderrStr += "\nRead stdout error: " + outErr.Error()
	}
	if errErr != nil {
		stderrStr += "\nRead stderr error: " + errErr.Error()
	}

	return stdoutStr, stderrStr, exitCode
}

// TestSignalCreateCommand tests the signal create command
func TestSignalCreateCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test creating a signal
	_, stderr, exitCode := runCommand(binary, "signal", "create", "test-contract-signal")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 0, exitCode)
	// assert.Contains(t, stdout, "Signal 'test-contract-signal' created")
}

// TestSignalListCommand tests the signal list command
func TestSignalListCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test listing signals
	_, stderr, exitCode := runCommand(binary, "signal", "list")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 0, exitCode)
	// assert.Contains(t, stdout, "No signals found") // or list of signals
}

// TestSignalWaitCommand tests the signal wait command
func TestSignalWaitCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test waiting for a signal with timeout
	start := time.Now()
	_, stderr, exitCode := runCommand(binary, "signal", "wait", "nonexistent-signal", "--timeout", "500ms")
	duration := time.Since(start)

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// Should complete quickly (not wait the full timeout)
	assert.Less(t, duration, 100*time.Millisecond, "Command should fail quickly, not wait for timeout")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 1, exitCode) // timeout exit code
	// assert.Contains(t, stderr, "timeout")
	// assert.Greater(t, duration, 450*time.Millisecond) // should wait most of the timeout
}

// TestSignalClearCommand tests the signal clear command
func TestSignalClearCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test clearing a signal
	_, stderr, exitCode := runCommand(binary, "signal", "clear", "test-clear-signal")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 1, exitCode) // signal doesn't exist
	// assert.Contains(t, stderr, "does not exist")
}

// TestSignalCommandHelp tests the signal command help
func TestSignalCommandHelp(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test signal command help
	_, stderr, exitCode := runCommand(binary, "signal", "--help")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 0, exitCode)
	// assert.Contains(t, stdout, "signal")
	// assert.Contains(t, stdout, "create")
	// assert.Contains(t, stdout, "list")
	// assert.Contains(t, stdout, "wait")
	// assert.Contains(t, stdout, "clear")
}

// TestSignalCommandInvalidArgs tests invalid arguments
func TestSignalCommandInvalidArgs(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test signal command with no subcommand
	_, stderr, exitCode := runCommand(binary, "signal")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 1, exitCode)
	// assert.Contains(t, stderr, "requires at least 1 arg")
}

// TestSignalCreateWithSpecialChars tests signal names with special characters
func TestSignalCreateWithSpecialChars(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test creating a signal with special characters
	_, stderr, exitCode := runCommand(binary, "signal", "create", "test_signal.with-dashes")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 0, exitCode)
	// assert.Contains(t, stdout, "Signal 'test_signal.with-dashes' created")
}

// TestSignalWaitWithInvalidTimeout tests invalid timeout values
func TestSignalWaitWithInvalidTimeout(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test waiting with invalid timeout
	_, stderr, exitCode := runCommand(binary, "signal", "wait", "test-signal", "--timeout", "invalid")

	// The command doesn't exist yet, so it should fail
	assert.NotEqual(t, 0, exitCode, "Command should fail since signal command is not implemented yet")
	assert.Contains(t, stderr, "unknown command", "Should show unknown command error")

	// TODO: Once implemented, this should pass:
	// assert.Equal(t, 1, exitCode)
	// assert.Contains(t, stderr, "invalid duration")
}
