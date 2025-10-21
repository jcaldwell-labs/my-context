package integration

import (
	"os/exec"
	"strings"
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
	stdout, stderr := cmd.StdoutPipe(), cmd.StderrPipe()

	// Start the command
	err := cmd.Start()
	if err != nil {
		return "", err.Error(), 1
	}

	// Read output
	outBytes, outErr := stdout.ReadAll()
	errBytes, errErr := stderr.ReadAll()

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

// TestSignalWorkflowIntegration tests the complete signal lifecycle
func TestSignalWorkflowIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// The signal command doesn't exist yet, so all operations should fail
	// This tests that the integration points are ready for when the command is implemented

	// Test 1: Try to create a signal
	stdout, stderr, exitCode := runCommand(binary, "signal", "create", "integration-test-signal")
	assert.NotEqual(t, 0, exitCode, "Create should fail - command not implemented")
	assert.Contains(t, stderr, "unknown command")

	// Test 2: Try to list signals
	stdout, stderr, exitCode = runCommand(binary, "signal", "list")
	assert.NotEqual(t, 0, exitCode, "List should fail - command not implemented")
	assert.Contains(t, stderr, "unknown command")

	// Test 3: Try to wait for a signal
	stdout, stderr, exitCode = runCommand(binary, "signal", "wait", "integration-test-signal", "--timeout", "1s")
	assert.NotEqual(t, 0, exitCode, "Wait should fail - command not implemented")
	assert.Contains(t, stderr, "unknown command")

	// TODO: Once implemented, this full workflow should work:

	/*
		// Create a signal
		stdout, stderr, exitCode = runCommand(binary, "signal", "create", "integration-test-signal")
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, stdout, "created")

		// List signals - should show our signal
		stdout, stderr, exitCode = runCommand(binary, "signal", "list")
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, stdout, "integration-test-signal")

		// Wait for the signal should return immediately
		start := time.Now()
		stdout, stderr, exitCode = runCommand(binary, "signal", "wait", "integration-test-signal", "--timeout", "5s")
		duration := time.Since(start)
		assert.Equal(t, 0, exitCode)
		assert.Less(t, duration, 100*time.Millisecond) // Should return quickly

		// Clear the signal
		stdout, stderr, exitCode = runCommand(binary, "signal", "clear", "integration-test-signal")
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, stdout, "cleared")

		// List signals - should be empty now
		stdout, stderr, exitCode = runCommand(binary, "signal", "list")
		assert.Equal(t, 0, exitCode)
		assert.NotContains(t, stdout, "integration-test-signal")
	*/
}

// TestSignalConcurrencyIntegration tests concurrent signal operations
func TestSignalConcurrencyIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test concurrent signal operations (when implemented)
	// For now, just verify the commands fail as expected

	commands := []struct {
		name string
		args []string
	}{
		{"create1", []string{"signal", "create", "concurrency-signal-1"}},
		{"create2", []string{"signal", "create", "concurrency-signal-2"}},
		{"list", []string{"signal", "list"}},
	}

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			stdout, stderr, exitCode := runCommand(binary, cmd.args...)
			assert.NotEqual(t, 0, exitCode, "Command should fail - not implemented yet")
			assert.Contains(t, stderr, "unknown command")
		})
	}

	// TODO: Once implemented, test concurrent operations:

	/*
		// Create signals concurrently
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				stdout, stderr, exitCode := runCommand(binary, "signal", "create", fmt.Sprintf("concurrent-%d", id))
				assert.Equal(t, 0, exitCode)
			}(i)
		}
		wg.Wait()

		// List should show all signals
		stdout, stderr, exitCode := runCommand(binary, "signal", "list")
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, stdout, "concurrent-")
	*/
}

// TestSignalTimeoutIntegration tests timeout behavior
func TestSignalTimeoutIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test timeout behavior (when implemented)
	stdout, stderr, exitCode := runCommand(binary, "signal", "wait", "timeout-test-signal", "--timeout", "200ms")

	// Should fail because command doesn't exist
	assert.NotEqual(t, 0, exitCode)
	assert.Contains(t, stderr, "unknown command")

	// TODO: Once implemented, test timeout:

	/*
		start := time.Now()
		stdout, stderr, exitCode = runCommand(binary, "signal", "wait", "timeout-test-signal", "--timeout", "200ms")
		duration := time.Since(start)

		assert.Equal(t, 1, exitCode) // timeout exit code
		assert.Contains(t, stderr, "timeout")
		assert.Greater(t, duration, 150*time.Millisecond) // should wait most of timeout
		assert.Less(t, duration, 300*time.Millisecond) // but not much longer
	*/
}

// TestSignalErrorHandlingIntegration tests error scenarios
func TestSignalErrorHandlingIntegration(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() { _ = exec.Command("rm", binary).Run() }()

	// Test various error scenarios (when implemented)
	errorTests := []struct {
		name     string
		args     []string
		errorMsg string
	}{
		{"no args", []string{"signal"}, "requires at least 1 arg"},
		{"invalid subcommand", []string{"signal", "invalid"}, "unknown subcommand"},
		{"create no name", []string{"signal", "create"}, "accepts 1 arg"},
		{"wait no name", []string{"signal", "wait"}, "accepts 1 arg"},
		{"clear no name", []string{"signal", "clear"}, "accepts 1 arg"},
		{"create duplicate", []string{"signal", "create", "duplicate", "&&", "signal", "create", "duplicate"}, "already exists"},
		{"clear nonexistent", []string{"signal", "clear", "nonexistent"}, "does not exist"},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			// For now, all should fail with "unknown command"
			stdout, stderr, exitCode := runCommand(binary, test.args...)
			assert.NotEqual(t, 0, exitCode, "Command should fail - not implemented yet")
			assert.Contains(t, stderr, "unknown command")
		})
	}

	// TODO: Once implemented, test specific error messages:

	/*
		for _, test := range errorTests {
			t.Run(test.name, func(t *testing.T) {
				stdout, stderr, exitCode := runCommand(binary, strings.Fields(test.args[0])...)
				assert.Equal(t, 1, exitCode)
				assert.Contains(t, stderr, test.errorMsg)
			})
		}
	*/
}
