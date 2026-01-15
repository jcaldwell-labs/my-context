package contract

import (
	"encoding/json"
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
	statsTestBinaryPath string
)

// buildStatsTestBinary builds the my-context binary for testing
func buildStatsTestBinary(t *testing.T) string {
	if statsTestBinaryPath != "" {
		return statsTestBinaryPath
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
	tmpDir, err := os.MkdirTemp("", "stats-test-*")
	require.NoError(t, err)

	binaryName := "my-context-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	statsTestBinaryPath = filepath.Join(tmpDir, binaryName)

	cmd := exec.Command("go", "build", "-o", statsTestBinaryPath, "./cmd/my-context/")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to build test binary")

	return statsTestBinaryPath
}

// runStatsCommand executes a my-context command with isolated MY_CONTEXT_HOME
func runStatsCommand(binary, testDir string, args ...string) (stdout string, exitCode int) {
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

// TestStatsCommandEmpty tests stats when no contexts exist
func TestStatsCommandEmpty(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runStatsCommand(binary, testDir, "stats")

	assert.Equal(t, 0, exitCode, "stats should succeed even when empty: %s", stdout)
	assert.Contains(t, stdout, "Time Statistics", "Should show statistics header")
	assert.Contains(t, stdout, "Total Contexts: 0", "Should show zero contexts")
}

// TestStatsCommandWithContexts tests stats with existing contexts
func TestStatsCommandWithContexts(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create some contexts with time gaps
	_, exitCode := runStatsCommand(binary, testDir, "start", "project-a: task 1")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Small delay to accumulate some time
	time.Sleep(100 * time.Millisecond)

	_, exitCode = runStatsCommand(binary, testDir, "stop")
	require.Equal(t, 0, exitCode, "stop should succeed")

	_, exitCode = runStatsCommand(binary, testDir, "start", "project-b: task 1")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check stats
	stdout, exitCode := runStatsCommand(binary, testDir, "stats")

	assert.Equal(t, 0, exitCode, "stats should succeed: %s", stdout)
	assert.Contains(t, stdout, "Time Statistics", "Should show statistics header")
	assert.Contains(t, stdout, "Total Contexts:", "Should show total contexts")
}

// TestStatsCommandTodayFlag tests the --today flag
func TestStatsCommandTodayFlag(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create a context
	_, exitCode := runStatsCommand(binary, testDir, "start", "today-test")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check stats for today
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--today")

	assert.Equal(t, 0, exitCode, "stats --today should succeed: %s", stdout)
	assert.Contains(t, stdout, "Period: Today", "Should show 'Today' period")
}

// TestStatsCommandWeekFlag tests the --week flag
func TestStatsCommandWeekFlag(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create a context
	_, exitCode := runStatsCommand(binary, testDir, "start", "week-test")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check stats for this week
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--week")

	assert.Equal(t, 0, exitCode, "stats --week should succeed: %s", stdout)
	assert.Contains(t, stdout, "Period: This Week", "Should show 'This Week' period")
}

// TestStatsCommandMonthFlag tests the --month flag
func TestStatsCommandMonthFlag(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create a context
	_, exitCode := runStatsCommand(binary, testDir, "start", "month-test")
	require.Equal(t, 0, exitCode, "start should succeed")

	// Check stats for this month
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--month")

	assert.Equal(t, 0, exitCode, "stats --month should succeed: %s", stdout)
	assert.Contains(t, stdout, "Period: This Month", "Should show 'This Month' period")
}

// TestStatsCommandProjectFilter tests the --project flag
func TestStatsCommandProjectFilter(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create contexts for different projects
	_, exitCode := runStatsCommand(binary, testDir, "start", "project-a: task 1")
	require.Equal(t, 0, exitCode)
	_, exitCode = runStatsCommand(binary, testDir, "stop")
	require.Equal(t, 0, exitCode)

	_, exitCode = runStatsCommand(binary, testDir, "start", "project-b: task 1")
	require.Equal(t, 0, exitCode)
	_, exitCode = runStatsCommand(binary, testDir, "stop")
	require.Equal(t, 0, exitCode)

	// Filter by project-a
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--project", "project-a")

	assert.Equal(t, 0, exitCode, "stats --project should succeed: %s", stdout)
	assert.Contains(t, stdout, "Project: project-a", "Should show project filter")
	assert.Contains(t, stdout, "Total Contexts: 1", "Should only count project-a contexts")
}

// TestStatsCommandSinceFlag tests the --since flag
func TestStatsCommandSinceFlag(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create a context
	_, exitCode := runStatsCommand(binary, testDir, "start", "since-test")
	require.Equal(t, 0, exitCode)

	// Check stats since a date
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--since", "2020-01-01")

	assert.Equal(t, 0, exitCode, "stats --since should succeed: %s", stdout)
	assert.Contains(t, stdout, "Since 2020-01-01", "Should show 'Since' period")
}

// TestStatsCommandUntilFlag tests the --until flag
func TestStatsCommandUntilFlag(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Check stats until a future date
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--until", "2030-12-31")

	assert.Equal(t, 0, exitCode, "stats --until should succeed: %s", stdout)
	assert.Contains(t, stdout, "Until 2030-12-31", "Should show 'Until' period")
}

// TestStatsCommandSinceUntilFlags tests combined --since and --until flags
func TestStatsCommandSinceUntilFlags(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--since", "2025-01-01", "--until", "2025-12-31")

	assert.Equal(t, 0, exitCode, "stats with date range should succeed: %s", stdout)
	assert.Contains(t, stdout, "2025-01-01 to 2025-12-31", "Should show date range")
}

// TestStatsCommandInvalidSince tests invalid --since date format
func TestStatsCommandInvalidSince(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--since", "invalid")

	assert.NotEqual(t, 0, exitCode, "stats with invalid date should fail")
	assert.Contains(t, stdout, "invalid", "Should mention invalid date")
}

// TestStatsCommandJSONOutput tests JSON output format
func TestStatsCommandJSONOutput(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create a context
	_, exitCode := runStatsCommand(binary, testDir, "start", "json-test")
	require.Equal(t, 0, exitCode)

	// Get stats as JSON
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--json")

	assert.Equal(t, 0, exitCode, "stats --json should succeed: %s", stdout)

	// Verify it's valid JSON
	var result map[string]interface{}
	err := json.Unmarshal([]byte(stdout), &result)
	assert.NoError(t, err, "Output should be valid JSON")

	// Verify structure
	assert.Equal(t, "stats", result["command"], "Should have command field")
	assert.Contains(t, result, "data", "Should have data field")
}

// TestStatsCommandHelp tests the stats command help
func TestStatsCommandHelp(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "--help")

	assert.Equal(t, 0, exitCode, "stats --help should succeed: %s", stdout)
	assert.Contains(t, stdout, "stats", "Help should mention stats")
	assert.Contains(t, stdout, "today", "Help should mention --today flag")
	assert.Contains(t, stdout, "week", "Help should mention --week flag")
	assert.Contains(t, stdout, "month", "Help should mention --month flag")
	assert.Contains(t, stdout, "project", "Help should mention --project flag")
	assert.Contains(t, stdout, "since", "Help should mention --since flag")
	assert.Contains(t, stdout, "until", "Help should mention --until flag")
}

// TestStatsCommandShortAlias tests the 'st' alias
func TestStatsCommandShortAlias(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runStatsCommand(binary, testDir, "st")

	assert.Equal(t, 0, exitCode, "stats alias 'st' should succeed: %s", stdout)
	assert.Contains(t, stdout, "Time Statistics", "Should show statistics output")
}

// TestStatsCommandShortFlags tests short flag variants
func TestStatsCommandShortFlags(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Test -t for --today
	stdout, exitCode := runStatsCommand(binary, testDir, "stats", "-t")
	assert.Equal(t, 0, exitCode, "stats -t should succeed: %s", stdout)
	assert.Contains(t, stdout, "Today", "Should use today period")

	// Test -w for --week
	stdout, exitCode = runStatsCommand(binary, testDir, "stats", "-w")
	assert.Equal(t, 0, exitCode, "stats -w should succeed: %s", stdout)
	assert.Contains(t, stdout, "This Week", "Should use week period")

	// Test -m for --month
	stdout, exitCode = runStatsCommand(binary, testDir, "stats", "-m")
	assert.Equal(t, 0, exitCode, "stats -m should succeed: %s", stdout)
	assert.Contains(t, stdout, "This Month", "Should use month period")

	// Test -p for --project
	stdout, exitCode = runStatsCommand(binary, testDir, "stats", "-p", "test")
	assert.Equal(t, 0, exitCode, "stats -p should succeed: %s", stdout)
	assert.Contains(t, stdout, "Project: test", "Should filter by project")
}

// TestStatsCommandActiveContext tests that active context is shown
func TestStatsCommandActiveContext(t *testing.T) {
	binary := buildStatsTestBinary(t)
	testDir := t.TempDir()

	// Create and keep context active
	_, exitCode := runStatsCommand(binary, testDir, "start", "active-context-test")
	require.Equal(t, 0, exitCode)

	// Check stats shows active context
	stdout, exitCode := runStatsCommand(binary, testDir, "stats")

	assert.Equal(t, 0, exitCode, "stats should succeed: %s", stdout)
	assert.Contains(t, stdout, "Active Context:", "Should show active context section")
	assert.Contains(t, stdout, "active-context-test", "Should show active context name")
}
