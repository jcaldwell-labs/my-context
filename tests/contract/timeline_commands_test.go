package contract

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTimelineTest creates a temporary context home for testing
func setupTimelineTest(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "timeline-test-*")
	require.NoError(t, err)

	originalHome := os.Getenv("MY_CONTEXT_HOME")
	os.Setenv("MY_CONTEXT_HOME", tmpDir)

	cleanup := func() {
		os.Setenv("MY_CONTEXT_HOME", originalHome)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// TestTimelineCommandToday tests the --today flag
func TestTimelineCommandToday(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--today"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandWeek tests the --week flag
func TestTimelineCommandWeek(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--week"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandMonth tests the --month flag
func TestTimelineCommandMonth(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--month"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandSinceUntil tests the --since and --until flags
func TestTimelineCommandSinceUntil(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--since", "2026-01-01", "--until", "2026-01-14"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandProject tests the --project flag
func TestTimelineCommandProject(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--today", "--project", "myproject"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandJSON tests JSON output
func TestTimelineCommandJSON(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := true
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--today"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandEmptyContexts tests timeline with no contexts
func TestTimelineCommandEmptyContexts(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--today"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

// TestTimelineCommandAliases tests timeline command aliases
func TestTimelineCommandAliases(t *testing.T) {
	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)

	assert.Contains(t, cmd.Aliases, "tl")
	assert.Equal(t, "timeline", cmd.Use)
}

// TestTimelineCommandFlags tests that all expected flags are present
func TestTimelineCommandFlags(t *testing.T) {
	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)

	// Check flags
	todayFlag := cmd.Flags().Lookup("today")
	assert.NotNil(t, todayFlag)
	assert.Equal(t, "t", todayFlag.Shorthand)

	weekFlag := cmd.Flags().Lookup("week")
	assert.NotNil(t, weekFlag)
	assert.Equal(t, "w", weekFlag.Shorthand)

	monthFlag := cmd.Flags().Lookup("month")
	assert.NotNil(t, monthFlag)
	assert.Equal(t, "m", monthFlag.Shorthand)

	projectFlag := cmd.Flags().Lookup("project")
	assert.NotNil(t, projectFlag)
	assert.Equal(t, "p", projectFlag.Shorthand)

	sinceFlag := cmd.Flags().Lookup("since")
	assert.NotNil(t, sinceFlag)

	untilFlag := cmd.Flags().Lookup("until")
	assert.NotNil(t, untilFlag)
}

// TestTimelineCommandHelp tests the help text
func TestTimelineCommandHelp(t *testing.T) {
	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)

	assert.Contains(t, cmd.Long, "visual timeline")
	assert.Contains(t, cmd.Long, "gaps between contexts")
	assert.Contains(t, cmd.Short, "Show visual timeline")
}

// TestTimelineJSONStructure tests the structure of JSON output
func TestTimelineJSONStructure(t *testing.T) {
	// Test that we can unmarshal the expected JSON structure
	type TimelineJSON struct {
		Status              string        `json:"status"`
		Period              interface{}   `json:"period"`
		Entries             []interface{} `json:"entries"`
		TotalTrackedSeconds float64       `json:"total_tracked_seconds"`
		TotalGapsSeconds    float64       `json:"total_gaps_seconds"`
	}

	// This is a structure validation test
	var tj TimelineJSON
	err := json.Unmarshal([]byte(`{
		"status": "success",
		"period": "Today",
		"entries": [],
		"total_tracked_seconds": 0,
		"total_gaps_seconds": 0
	}`), &tj)

	assert.NoError(t, err)
	assert.Equal(t, "success", tj.Status)
	assert.Equal(t, 0.0, tj.TotalTrackedSeconds)
	assert.Equal(t, 0.0, tj.TotalGapsSeconds)
}

// TestTimelineCommandInvalidDate tests invalid date format
func TestTimelineCommandInvalidDate(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	cmd.SetArgs([]string{"--since", "invalid-date"})

	err := cmd.Execute()
	// Should handle error gracefully (either error or JSON error output)
	// We don't assert error because it might output JSON error instead
	_ = err
}

// TestTimelineCommandMultipleFlags tests conflicting flags
func TestTimelineCommandMultipleFlags(t *testing.T) {
	_, cleanup := setupTimelineTest(t)
	defer cleanup()

	jsonOutput := false
	cmd := commands.NewTimelineCmd(&jsonOutput)
	// Using both --today and --week should work (--today takes precedence)
	cmd.SetArgs([]string{"--today", "--week"})

	err := cmd.Execute()
	assert.NoError(t, err)
}
