package integration

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testBinaryOnce sync.Once
	testBinaryPath string
	testBinaryErr  error
	projectRoot    string
)

func init() {
	// Find project root (where go.mod is located)
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Walk up the directory tree until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod)")
		}
		dir = parent
	}
}

// buildTestBinary builds the my-context binary for testing (once)
func buildTestBinary(t *testing.T) string {
	testBinaryOnce.Do(func() {
		// Build in a temporary directory
		tmpDir, err := os.MkdirTemp("", "my-context-build-*")
		if err != nil {
			testBinaryErr = err
			return
		}
		testBinaryPath = filepath.Join(tmpDir, "my-context-test")

		cmd := exec.Command("go", "build", "-o", testBinaryPath, "./cmd/my-context/")
		testBinaryErr = cmd.Run()
	})

	require.NoError(t, testBinaryErr, "Failed to build test binary")
	return testBinaryPath
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	return projectRoot
}

// runCommand executes a command and returns error (for simple test cases)
func runCommand(args ...string) error {
	cmd := exec.Command("go", "run", "./cmd/my-context/main.go")
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir = projectRoot
	return cmd.Run()
}

// runCommandWithOutput executes a command and returns stdout and error
func runCommandWithOutput(args ...string) (string, error) {
	cmd := exec.Command("go", "run", "./cmd/my-context/main.go")
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// runCommandWithInput executes a command with stdin input
func runCommandWithInput(args ...string) error {
	cmd := exec.Command("go", "run", "./cmd/my-context/main.go")
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir = projectRoot
	return cmd.Run()
}

// runCommandFull executes a my-context command and returns stdout, stderr, and exit code
func runCommandFull(binary string, args ...string) (stdoutStr string, stderrStr string, exitCode int) {
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

	stdoutStr = string(outBytes)
	stderrStr = string(errBytes)

	exitCode = 0
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

	return
}

// setupTestEnvironment creates a temporary test directory
func setupTestEnvironment(t *testing.T) string {
	t.Helper()
	testDir, err := os.MkdirTemp("", "my-context-test-*")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	// Set MY_CONTEXT_HOME to test directory (t.Setenv handles cleanup automatically)
	t.Setenv("MY_CONTEXT_HOME", testDir)
	return testDir
}

// cleanupTestEnvironment removes the temporary test directory
func cleanupTestEnvironment(t *testing.T, testDir string) {
	t.Helper()
	os.RemoveAll(testDir)
	// Note: os.Unsetenv not needed - t.Setenv handles cleanup automatically
}

// createTestContext creates a test context directory structure
func createTestContext(t *testing.T, contextName string) {
	t.Helper()
	// This will fail until the actual implementation exists
	err := runCommand("start", contextName)
	if err != nil {
		t.Logf("Note: Context creation failed (expected until implementation): %v", err)
	}
	runCommand("stop")
}
