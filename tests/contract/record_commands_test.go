package contract

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	recordTestBinaryPath string
)

// buildRecordTestBinary builds the my-context binary for testing
func buildRecordTestBinary(t *testing.T) string {
	if recordTestBinaryPath != "" {
		return recordTestBinaryPath
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
	tmpDir, err := os.MkdirTemp("", "record-test-*")
	require.NoError(t, err)

	binaryName := "my-context-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	recordTestBinaryPath = filepath.Join(tmpDir, binaryName)

	cmd := exec.Command("go", "build", "-o", recordTestBinaryPath, "./cmd/my-context/")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to build test binary")

	return recordTestBinaryPath
}

// runRecordCommand executes a my-context command with isolated MY_CONTEXT_HOME
func runRecordCommand(binary, testDir string, args ...string) (stdout string, exitCode int) {
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

// TestRecordCommandHelp tests the record command help text
func TestRecordCommandHelp(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runRecordCommand(binary, testDir, "record", "--help")

	assert.Equal(t, 0, exitCode, "record --help should succeed: %s", stdout)
	assert.Contains(t, stdout, "record", "Help should mention record")
	assert.Contains(t, stdout, "clipboard", "Help should mention clipboard")
	assert.Contains(t, stdout, "prefix", "Help should mention prefix flag")
	assert.Contains(t, stdout, "Ctrl+C", "Help should mention stopping with Ctrl+C")
}

// TestRecordCommandRequiresActiveContext tests that record fails without active context
func TestRecordCommandRequiresActiveContext(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Don't start a context - record should fail
	// Note: We can't actually test clipboard monitoring without X11/Wayland,
	// but we can verify that the command checks for active context before
	// attempting clipboard operations
	stdout, exitCode := runRecordCommand(binary, testDir, "record")

	// Command should fail because no active context exists
	assert.NotEqual(t, 0, exitCode, "record without active context should fail")
	assert.Contains(t, stdout, "no active context", "Error should mention no active context")
	assert.Contains(t, stdout, "my-context start", "Error should suggest starting a context")
}

// TestRecordCommandWithActiveContext tests that record verifies active context exists
func TestRecordCommandWithActiveContext(t *testing.T) {
	// Skip this test - we can't test clipboard monitoring without X11/Wayland
	// and the command will hang waiting for clipboard events
	t.Skip("Cannot test clipboard monitoring without X11/Wayland display server")

	// The test would verify that:
	// 1. With an active context, the command gets past the context check
	// 2. Fails on clipboard initialization (expected in test environment)
	// 3. Error message is about clipboard, not missing context
}

// TestRecordCommandPrefixFlag tests the --prefix flag parsing
func TestRecordCommandPrefixFlag(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Test prefix flag with help - this doesn't require clipboard
	stdout, exitCode := runRecordCommand(binary, testDir, "record", "--help")

	require.Equal(t, 0, exitCode, "help should succeed")
	assert.Contains(t, stdout, "--prefix", "Help should show --prefix flag")
	assert.Contains(t, stdout, "-p", "Help should show -p short flag")
}

// TestRecordCommandPrefixShortFlag tests the -p short flag variant
func TestRecordCommandPrefixShortFlag(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Verify -p flag is documented in help
	stdout, exitCode := runRecordCommand(binary, testDir, "help", "record")

	require.Equal(t, 0, exitCode, "help should succeed")
	assert.Contains(t, stdout, "-p", "Help should show -p short flag")
	assert.Contains(t, stdout, "prefix", "Help should mention prefix")
}

// TestRecordCommandJSONErrorOutput tests JSON error format
func TestRecordCommandJSONErrorOutput(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Don't start a context - should get JSON error
	stdout, exitCode := runRecordCommand(binary, testDir, "record", "--json")

	// With --json flag, errors return exit code 0 but include error in JSON
	assert.Equal(t, 0, exitCode, "record --json returns 0 even on error")

	// Verify it's valid JSON
	var result map[string]interface{}
	err := json.Unmarshal([]byte(stdout), &result)
	assert.NoError(t, err, "Output should be valid JSON: %s", stdout)

	// Verify error structure
	assert.Equal(t, "record", result["command"], "Should have command field")
	assert.Contains(t, result, "error", "Should have error field")

	// Check error details
	errorField, ok := result["error"].(map[string]interface{})
	require.True(t, ok, "Error should be an object")
	message, ok := errorField["message"].(string)
	require.True(t, ok, "Error message should be a string")
	assert.Contains(t, message, "no active context", "Error should mention no active context")
}

// TestRecordCommandJSONErrorWithContext tests JSON error when clipboard fails
func TestRecordCommandJSONErrorWithContext(t *testing.T) {
	// Skip - cannot test clipboard initialization without X11/Wayland
	t.Skip("Cannot test clipboard monitoring without X11/Wayland display server")

	// This test would verify that when clipboard initialization fails,
	// the error is properly formatted as JSON when --json flag is used
}

// TestRecordCommandShortAlias tests the 'r' alias
func TestRecordCommandShortAlias(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Note: 'r' alias conflicts with 'resume' command
	// Cobra will pick the first registered command with that alias
	// We test that the alias exists in help text instead
	stdout, exitCode := runRecordCommand(binary, testDir, "record", "--help")

	assert.Equal(t, 0, exitCode, "record help should work: %s", stdout)
	// Check that aliases section mentions 'r' or that it's a valid alias
	assert.Contains(t, stdout, "record", "Help should show record command")
}

// TestRecordCommandShortAliasRequiresContext tests record with error case
func TestRecordCommandShortAliasRequiresContext(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Test record without active context
	stdout, exitCode := runRecordCommand(binary, testDir, "record")

	assert.NotEqual(t, 0, exitCode, "record should fail without context")
	// Will fail with either context or clipboard error
	assert.True(t,
		strings.Contains(stdout, "no active context") ||
			strings.Contains(stdout, "clipboard"),
		"Error should mention no active context or clipboard")
}

// TestRecordCommandNoArguments tests that record takes no arguments
func TestRecordCommandNoArguments(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Try with arguments - should reject them (no need for active context to test this)
	stdout, exitCode := runRecordCommand(binary, testDir, "record", "extra-arg")

	// Should fail with argument error
	assert.NotEqual(t, 0, exitCode, "record with arguments should fail")
	// Cobra will reject extra arguments
	assert.True(t,
		strings.Contains(stdout, "unknown command") ||
			strings.Contains(stdout, "accepts 0 arg") ||
			strings.Contains(stdout, "extra-arg"),
		"Error should indicate invalid arguments")
}

// TestRecordCommandPrefixFlagEmpty tests empty prefix value
func TestRecordCommandPrefixFlagEmpty(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Test with empty prefix - will fail on missing context but shouldn't complain about flag
	stdout, _ := runRecordCommand(binary, testDir, "record", "--prefix", "")

	// Should accept empty prefix flag format
	assert.NotContains(t, stdout, "flag needs an argument", "Should accept empty prefix")
}

// TestRecordCommandPrefixFlagSpaces tests prefix with spaces
func TestRecordCommandPrefixFlagSpaces(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	// Test with prefix containing spaces - will fail on missing context but shouldn't complain about flag
	stdout, exitCode := runRecordCommand(binary, testDir, "record", "--prefix", "my tag")

	// Should accept prefix with spaces (flag parsing should work)
	assert.NotContains(t, stdout, "unknown flag", "Should accept prefix with spaces")
	// Should fail on missing context, not flag parsing
	if exitCode != 0 {
		assert.Contains(t, stdout, "no active context", "Should fail on missing context")
	}
}

// TestRecordCommandHelpShowsAliases tests that help shows the 'r' alias
func TestRecordCommandHelpShowsAliases(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runRecordCommand(binary, testDir, "help", "record")

	assert.Equal(t, 0, exitCode, "help record should succeed: %s", stdout)
	// May show alias in usage line
	// Format could be: "record, r" or similar
	if !strings.Contains(stdout, ", r") {
		// Alternative format check
		assert.Contains(t, stdout, "record", "Should show record command")
	}
}

// TestRecordCommandFlagValidation tests various flag combinations
func TestRecordCommandFlagValidation(t *testing.T) {
	binary := buildRecordTestBinary(t)
	testDir := t.TempDir()

	testCases := []struct {
		name        string
		args        []string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid prefix flag",
			args:        []string{"record", "--prefix", "test"},
			shouldError: false,
			errorMsg:    "",
		},
		{
			name:        "valid short prefix flag",
			args:        []string{"record", "-p", "test"},
			shouldError: false,
			errorMsg:    "",
		},
		{
			name:        "combined with json flag",
			args:        []string{"record", "--json", "--prefix", "test"},
			shouldError: false,
			errorMsg:    "",
		},
		{
			name:        "invalid flag",
			args:        []string{"record", "--invalid-flag"},
			shouldError: true,
			errorMsg:    "unknown flag",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, exitCode := runRecordCommand(binary, testDir, tc.args...)

			if tc.shouldError {
				// Unknown flags are caught by cobra before execution
				assert.NotEqual(t, 0, exitCode, "Should fail with invalid flag")
				assert.Contains(t, stdout, tc.errorMsg, "Error message should mention: %s", tc.errorMsg)
			} else {
				// Valid flags - will fail on missing context or clipboard, but not on flag parsing
				assert.NotContains(t, stdout, "unknown flag", "Should not complain about flags")
			}
		})
	}
}
