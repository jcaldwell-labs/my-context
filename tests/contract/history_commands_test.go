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
	historyTestBinaryPath string
)

// buildHistoryTestBinary builds the my-context binary for testing
func buildHistoryTestBinary(t *testing.T) string {
	if historyTestBinaryPath != "" {
		return historyTestBinaryPath
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
	tmpDir, err := os.MkdirTemp("", "history-test-*")
	require.NoError(t, err)

	binaryName := "my-context-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	historyTestBinaryPath = filepath.Join(tmpDir, binaryName)

	cmd := exec.Command("go", "build", "-o", historyTestBinaryPath, "./cmd/my-context/")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to build test binary")

	return historyTestBinaryPath
}

// runHistoryCommand executes a my-context command with isolated MY_CONTEXT_HOME
func runHistoryCommand(binary, testDir string, args ...string) (stdout string, exitCode int) {
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

// TestHistoryCommandEmpty tests history when no contexts exist
func TestHistoryCommandEmpty(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runHistoryCommand(binary, testDir, "history")

	assert.Equal(t, 0, exitCode, "history should succeed even when empty: %s", stdout)
	assert.Contains(t, stdout, "Context History", "Should show history header")
}

// TestHistoryCommandWithContexts tests history with context transitions
func TestHistoryCommandWithContexts(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create some contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "first-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	_, exitCode = runHistoryCommand(binary, testDir, "stop")
	require.Equal(t, 0, exitCode, "stop should succeed")

	_, exitCode = runHistoryCommand(binary, testDir, "start", "second-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history
	stdout, exitCode := runHistoryCommand(binary, testDir, "history")

	assert.Equal(t, 0, exitCode, "history should succeed: %s", stdout)
	assert.Contains(t, stdout, "first-context", "Should show first context")
	assert.Contains(t, stdout, "second-context", "Should show second context")
	assert.Contains(t, stdout, "START", "Should show START transitions")
	assert.Contains(t, stdout, "STOP", "Should show STOP transitions")
}

// TestHistoryCommandSwitchDetection tests that switch transitions are detected
func TestHistoryCommandSwitchDetection(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts with automatic switching (start while another is active)
	_, exitCode := runHistoryCommand(binary, testDir, "start", "context-a")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Starting a new context while one is active should create a switch
	_, exitCode = runHistoryCommand(binary, testDir, "start", "context-b")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history
	stdout, exitCode := runHistoryCommand(binary, testDir, "history")

	assert.Equal(t, 0, exitCode, "history should succeed: %s", stdout)
	// In file mode, switch is detected from transitions.log
	// The output should show the transition between contexts
	assert.Contains(t, stdout, "context-a", "Should show context-a")
	assert.Contains(t, stdout, "context-b", "Should show context-b")
}

// TestHistoryCommandLimitFlag tests the --limit flag
func TestHistoryCommandLimitFlag(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create multiple contexts
	for i := 0; i < 5; i++ {
		_, exitCode := runHistoryCommand(binary, testDir, "start", "context-"+string(rune('a'+i)))
		require.Equal(t, 0, exitCode, "start should succeed")
		_, exitCode = runHistoryCommand(binary, testDir, "stop")
		require.Equal(t, 0, exitCode, "stop should succeed")
	}

	// Test with limit
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--limit", "3")

	assert.Equal(t, 0, exitCode, "history with limit should succeed: %s", stdout)
	// Should have limited output
	lines := strings.Split(stdout, "\n")
	transitionCount := 0
	for _, line := range lines {
		if strings.Contains(line, "START") || strings.Contains(line, "STOP") || strings.Contains(line, "SWITCH") {
			transitionCount++
		}
	}
	assert.LessOrEqual(t, transitionCount, 3, "Should respect limit flag")
}

// TestHistoryCommandJSONOutput tests JSON output format
func TestHistoryCommandJSONOutput(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create a context
	_, exitCode := runHistoryCommand(binary, testDir, "start", "json-test-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Get history as JSON
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--json")

	assert.Equal(t, 0, exitCode, "history --json should succeed: %s", stdout)

	// Verify it's valid JSON
	var result map[string]interface{}
	err := json.Unmarshal([]byte(stdout), &result)
	assert.NoError(t, err, "Output should be valid JSON")

	// Verify structure
	assert.Equal(t, "history", result["command"], "Should have command field")
	assert.Contains(t, result, "data", "Should have data field")
}

// TestHistoryCommandHelp tests the history command help
func TestHistoryCommandHelp(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--help")

	assert.Equal(t, 0, exitCode, "history --help should succeed: %s", stdout)
	assert.Contains(t, stdout, "history", "Help should mention history")
	assert.Contains(t, stdout, "limit", "Help should mention limit flag")
	assert.Contains(t, stdout, "switch", "Help should mention switch events")
}

// TestHistoryCommandShortAlias tests the 'h' alias
func TestHistoryCommandShortAlias(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runHistoryCommand(binary, testDir, "h")

	assert.Equal(t, 0, exitCode, "history alias 'h' should succeed: %s", stdout)
	assert.Contains(t, stdout, "Context History", "Should show history output")
}

// TestHistoryCommandLimitZero tests unlimited output with --limit 0
func TestHistoryCommandLimitZero(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create some contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "unlimited-test")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Test with limit 0 (unlimited)
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--limit", "0")

	assert.Equal(t, 0, exitCode, "history with limit 0 should succeed: %s", stdout)
	assert.Contains(t, stdout, "unlimited-test", "Should show context")
}

// TestHistoryCommandTodayFlag tests the --today flag
func TestHistoryCommandTodayFlag(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "today-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with --today filter
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--today")

	assert.Equal(t, 0, exitCode, "history --today should succeed: %s", stdout)
	assert.Contains(t, stdout, "today-context", "Should show today's context")
}

// TestHistoryCommandWeekFlag tests the --week flag
func TestHistoryCommandWeekFlag(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "week-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with --week filter
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--week")

	assert.Equal(t, 0, exitCode, "history --week should succeed: %s", stdout)
	assert.Contains(t, stdout, "week-context", "Should show this week's context")
}

// TestHistoryCommandMonthFlag tests the --month flag
func TestHistoryCommandMonthFlag(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "month-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with --month filter
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--month")

	assert.Equal(t, 0, exitCode, "history --month should succeed: %s", stdout)
	assert.Contains(t, stdout, "month-context", "Should show this month's context")
}

// TestHistoryCommandSinceFlag tests the --since flag
func TestHistoryCommandSinceFlag(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "recent-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with --since filter (should include all contexts from beginning of time)
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--since", "2020-01-01")

	assert.Equal(t, 0, exitCode, "history --since should succeed: %s", stdout)
	assert.Contains(t, stdout, "recent-context", "Should show context after since date")
}

// TestHistoryCommandUntilFlag tests the --until flag
func TestHistoryCommandUntilFlag(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "until-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with --until filter (far future should include everything)
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--until", "2099-12-31")

	assert.Equal(t, 0, exitCode, "history --until should succeed: %s", stdout)
	assert.Contains(t, stdout, "until-context", "Should show context before until date")
}

// TestHistoryCommandSinceUntilFlags tests both --since and --until flags
func TestHistoryCommandSinceUntilFlags(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "range-context")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with both filters
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--since", "2020-01-01", "--until", "2099-12-31")

	assert.Equal(t, 0, exitCode, "history --since --until should succeed: %s", stdout)
	assert.Contains(t, stdout, "range-context", "Should show context in date range")
}

// TestHistoryCommandInvalidSinceDate tests error handling for invalid --since date
func TestHistoryCommandInvalidSinceDate(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Try to use invalid date format
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--since", "invalid-date")

	assert.NotEqual(t, 0, exitCode, "history with invalid --since should fail")
	assert.Contains(t, stdout, "invalid --since date format", "Should show error message")
}

// TestHistoryCommandInvalidUntilDate tests error handling for invalid --until date
func TestHistoryCommandInvalidUntilDate(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Try to use invalid date format
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--until", "not-a-date")

	assert.NotEqual(t, 0, exitCode, "history with invalid --until should fail")
	assert.Contains(t, stdout, "invalid --until date format", "Should show error message")
}

// TestHistoryCommandDateFlagsWithJSON tests JSON output with date flags
func TestHistoryCommandDateFlagsWithJSON(t *testing.T) {
	binary := buildHistoryTestBinary(t)
	testDir := t.TempDir()

	// Create contexts
	_, exitCode := runHistoryCommand(binary, testDir, "start", "json-date-test")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check history with --today and --json
	stdout, exitCode := runHistoryCommand(binary, testDir, "history", "--today", "--json")

	assert.Equal(t, 0, exitCode, "history --today --json should succeed: %s", stdout)

	// Verify it's valid JSON
	var result map[string]interface{}
	err := json.Unmarshal([]byte(stdout), &result)
	assert.NoError(t, err, "Output should be valid JSON")
	assert.Equal(t, "history", result["command"], "Should have command field")
}
