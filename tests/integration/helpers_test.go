package integration

import (
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// buildTestBinary builds the my-context binary for testing
func buildTestBinary(t *testing.T) string {
	cmd := exec.Command("go", "build", "-o", "my-context-test", "./cmd/my-context/")
	err := cmd.Run()
	require.NoError(t, err, "Failed to build test binary")
	return "./my-context-test"
}

// runCommand executes a command and returns error (for simple test cases)
func runCommand(args ...string) error {
	// This is a placeholder that will be replaced with actual command execution
	// For now, it will fail (which is expected for TDD)
	return os.ErrNotExist
}

// runCommandWithOutput executes a command and returns stdout and error
func runCommandWithOutput(args ...string) (string, error) {
	// Placeholder - will return empty until implementation exists
	return "", os.ErrNotExist
}

// runCommandWithInput executes a command with stdin input
func runCommandWithInput(args ...string) error {
	// Placeholder
	return os.ErrNotExist
}

// runCommandFull executes a my-context command and returns stdout, stderr, and exit code
func runCommandFull(binary string, args ...string) (string, string, int) {
	cmd := exec.Command(binary, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err.Error(), 1
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err.Error(), 1
	}

	// Start the command
	err = cmd.Start()
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

// setupTestEnvironment creates a temporary test directory
func setupTestEnvironment(t *testing.T) string {
	t.Helper()
	testDir, err := os.MkdirTemp("", "my-context-test-*")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	// Set MY_CONTEXT_HOME to test directory
	os.Setenv("MY_CONTEXT_HOME", testDir)
	return testDir
}

// cleanupTestEnvironment removes the temporary test directory
func cleanupTestEnvironment(t *testing.T, testDir string) {
	t.Helper()
	os.RemoveAll(testDir)
	os.Unsetenv("MY_CONTEXT_HOME")
}

// createTestContext creates a test context directory structure
func createTestContext(t *testing.T, testDir, contextName string) {
	t.Helper()
	// This will fail until the actual implementation exists
	err := runCommand("start", contextName)
	if err != nil {
		t.Logf("Note: Context creation failed (expected until implementation): %v", err)
	}
	runCommand("stop")
}
