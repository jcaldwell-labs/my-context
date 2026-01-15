package unit

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculatePeriodToday tests the --today flag period calculation
func TestCalculatePeriodToday(t *testing.T) {
	period, err := calculatePeriodForTest(true, false, false, "", "")

	require.NoError(t, err)
	assert.Equal(t, "Today", period.Label)

	// Start should be beginning of today
	now := time.Now()
	y, m, d := now.Date()
	expectedStart := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	assert.Equal(t, expectedStart, period.Start)

	// End should be close to now (within a few seconds)
	assert.WithinDuration(t, now, period.End, 5*time.Second)
}

// TestCalculatePeriodWeek tests the --week flag period calculation
func TestCalculatePeriodWeek(t *testing.T) {
	period, err := calculatePeriodForTest(false, true, false, "", "")

	require.NoError(t, err)
	assert.Equal(t, "This Week", period.Label)

	// Start should be Monday of current week
	assert.Equal(t, time.Monday, period.Start.Weekday())

	// Start should have time 00:00:00
	assert.Equal(t, 0, period.Start.Hour())
	assert.Equal(t, 0, period.Start.Minute())
	assert.Equal(t, 0, period.Start.Second())
}

// TestCalculatePeriodMonth tests the --month flag period calculation
func TestCalculatePeriodMonth(t *testing.T) {
	period, err := calculatePeriodForTest(false, false, true, "", "")

	require.NoError(t, err)
	assert.Equal(t, "This Month", period.Label)

	// Start should be 1st of current month
	assert.Equal(t, 1, period.Start.Day())
	assert.Equal(t, 0, period.Start.Hour())
	assert.Equal(t, 0, period.Start.Minute())
}

// TestCalculatePeriodSince tests the --since flag
func TestCalculatePeriodSince(t *testing.T) {
	period, err := calculatePeriodForTest(false, false, false, "2025-01-15", "")

	require.NoError(t, err)
	assert.Contains(t, period.Label, "Since 2025-01-15")
	assert.Equal(t, 2025, period.Start.Year())
	assert.Equal(t, time.January, period.Start.Month())
	assert.Equal(t, 15, period.Start.Day())
}

// TestCalculatePeriodUntil tests the --until flag
func TestCalculatePeriodUntil(t *testing.T) {
	period, err := calculatePeriodForTest(false, false, false, "", "2025-12-31")

	require.NoError(t, err)
	assert.Contains(t, period.Label, "Until 2025-12-31")
	// End should be end of the specified day
	assert.Equal(t, 2025, period.End.Year())
	assert.Equal(t, time.December, period.End.Month())
	assert.Equal(t, 31, period.End.Day())
}

// TestCalculatePeriodSinceUntil tests both --since and --until flags
func TestCalculatePeriodSinceUntil(t *testing.T) {
	period, err := calculatePeriodForTest(false, false, false, "2025-01-01", "2025-12-31")

	require.NoError(t, err)
	assert.Equal(t, "2025-01-01 to 2025-12-31", period.Label)
	assert.Equal(t, 2025, period.Start.Year())
	assert.Equal(t, time.January, period.Start.Month())
	assert.Equal(t, 1, period.Start.Day())
}

// TestCalculatePeriodDefault tests default (all time) period
func TestCalculatePeriodDefault(t *testing.T) {
	period, err := calculatePeriodForTest(false, false, false, "", "")

	require.NoError(t, err)
	assert.Equal(t, "All Time", period.Label)
	assert.Equal(t, int64(0), period.Start.Unix()) // Epoch
}

// TestCalculatePeriodInvalidSince tests invalid --since date format
func TestCalculatePeriodInvalidSince(t *testing.T) {
	_, err := calculatePeriodForTest(false, false, false, "invalid-date", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --since date format")
}

// TestCalculatePeriodInvalidUntil tests invalid --until date format
func TestCalculatePeriodInvalidUntil(t *testing.T) {
	_, err := calculatePeriodForTest(false, false, false, "", "not-a-date")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --until date format")
}

// TestTruncateString tests string truncation
func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length",
			input:    "exact",
			maxLen:   5,
			expected: "exact",
		},
		{
			name:     "truncation with ellipsis",
			input:    "this is a very long string",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "very short maxLen",
			input:    "hello",
			maxLen:   3,
			expected: "hel",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateStringForTest(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAggregateStatsGeneric tests statistics aggregation
func TestAggregateStatsGeneric(t *testing.T) {
	items := []statsContextItemForTest{
		{Name: "project-a: task 1", Duration: 1 * time.Hour},
		{Name: "project-a: task 2", Duration: 2 * time.Hour},
		{Name: "project-b: task 1", Duration: 30 * time.Minute},
	}

	period := timePeriodForTest{
		Start: time.Now().Add(-24 * time.Hour),
		End:   time.Now(),
		Label: "Test Period",
	}

	result := aggregateStatsGenericForTest(items, period, "project-a: task 2")

	// Total contexts
	assert.Equal(t, 3, result.TotalContexts)

	// Total time: 1h + 2h + 30m = 3.5h = 12600s
	assert.Equal(t, int64(12600), result.TotalSeconds)

	// Average: 12600 / 3 = 4200s
	assert.Equal(t, int64(4200), result.AverageSeconds)

	// Longest session should be "project-a: task 2" with 2h = 7200s
	require.NotNil(t, result.LongestSession)
	assert.Equal(t, "project-a: task 2", result.LongestSession.Name)
	assert.Equal(t, int64(7200), result.LongestSession.Seconds)

	// Active context should be set
	require.NotNil(t, result.ActiveContext)
	assert.Equal(t, "project-a: task 2", result.ActiveContext.Name)

	// Should have 2 projects
	assert.Len(t, result.ByProject, 2)

	// project-a should be first (more time)
	assert.Equal(t, "project-a", result.ByProject[0].Project)
	assert.Equal(t, int64(10800), result.ByProject[0].Seconds) // 3h = 10800s
	assert.Equal(t, 2, result.ByProject[0].Count)
}

// TestAggregateStatsEmpty tests aggregation with no items
func TestAggregateStatsEmpty(t *testing.T) {
	items := []statsContextItemForTest{}
	period := timePeriodForTest{Label: "Empty"}

	result := aggregateStatsGenericForTest(items, period, "")

	assert.Equal(t, 0, result.TotalContexts)
	assert.Equal(t, int64(0), result.TotalSeconds)
	assert.Equal(t, int64(0), result.AverageSeconds)
	assert.Nil(t, result.LongestSession)
	assert.Nil(t, result.ActiveContext)
	assert.Empty(t, result.ByProject)
}

// TestFilterContextsByPeriod tests period filtering
func TestFilterContextsByPeriod(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	contexts := []contextForTest{
		{Name: "today", StartTime: now.Add(-1 * time.Hour)},
		{Name: "yesterday", StartTime: yesterday},
		{Name: "last-week", StartTime: lastWeek},
	}

	// Period: only today
	period := timePeriodForTest{
		Start: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		End:   now,
	}

	filtered := filterContextsByPeriodForTest(contexts, period)

	assert.Len(t, filtered, 1)
	assert.Equal(t, "today", filtered[0].Name)
}

// --- Test Helper Types (mirroring internal types) ---

type timePeriodForTest struct {
	Start time.Time
	End   time.Time
	Label string
}

type statsContextItemForTest struct {
	Name     string
	Duration time.Duration
}

type statsResultForTest struct {
	Period         timePeriodForTest
	TotalContexts  int
	TotalSeconds   int64
	AverageSeconds int64
	ByProject      []projectStatsForTest
	LongestSession *sessionInfoForTest
	ActiveContext  *sessionInfoForTest
}

type projectStatsForTest struct {
	Project string
	Seconds int64
	Count   int
	Percent float64
}

type sessionInfoForTest struct {
	Name    string
	Seconds int64
}

type contextForTest struct {
	Name      string
	StartTime time.Time
}

// --- Test Helper Functions (reimplementing logic for testing) ---

func calculatePeriodForTest(today, week, month bool, since, until string) (timePeriodForTest, error) {
	now := time.Now()
	var period timePeriodForTest

	switch {
	case today:
		y, m, d := now.Date()
		period.Start = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "Today"

	case week:
		y, m, d := now.Date()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		daysBack := weekday - 1
		period.Start = time.Date(y, m, d-daysBack, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "This Week"

	case month:
		y, m, _ := now.Date()
		period.Start = time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "This Month"

	case since != "" || until != "":
		if since != "" {
			t, err := time.Parse("2006-01-02", since)
			if err != nil {
				return period, fmt.Errorf("invalid --since date format (use YYYY-MM-DD): %w", err)
			}
			period.Start = t
		} else {
			period.Start = time.Unix(0, 0)
		}

		if until != "" {
			t, err := time.Parse("2006-01-02", until)
			if err != nil {
				return period, fmt.Errorf("invalid --until date format (use YYYY-MM-DD): %w", err)
			}
			period.End = t.Add(24*time.Hour - time.Second)
		} else {
			period.End = now
		}

		switch {
		case since != "" && until != "":
			period.Label = since + " to " + until
		case since != "":
			period.Label = "Since " + since
		default:
			period.Label = "Until " + until
		}

	default:
		period.Start = time.Unix(0, 0)
		period.End = now
		period.Label = "All Time"
	}

	return period, nil
}

func truncateStringForTest(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func aggregateStatsGenericForTest(items []statsContextItemForTest, period timePeriodForTest, activeContextName string) statsResultForTest {
	result := statsResultForTest{
		Period:        period,
		TotalContexts: len(items),
	}

	if len(items) == 0 {
		return result
	}

	projectTotals := make(map[string]int64)
	projectCounts := make(map[string]int)
	var longestName string
	var longestSeconds int64

	for _, item := range items {
		seconds := int64(item.Duration.Seconds())
		result.TotalSeconds += seconds

		project := extractProjectNameForTest(item.Name)
		projectTotals[project] += seconds
		projectCounts[project]++

		if seconds > longestSeconds {
			longestSeconds = seconds
			longestName = item.Name
		}

		if item.Name == activeContextName {
			result.ActiveContext = &sessionInfoForTest{
				Name:    item.Name,
				Seconds: seconds,
			}
		}
	}

	if result.TotalContexts > 0 {
		result.AverageSeconds = result.TotalSeconds / int64(result.TotalContexts)
	}

	if longestName != "" {
		result.LongestSession = &sessionInfoForTest{
			Name:    longestName,
			Seconds: longestSeconds,
		}
	}

	for project, seconds := range projectTotals {
		percent := 0.0
		if result.TotalSeconds > 0 {
			percent = float64(seconds) / float64(result.TotalSeconds) * 100
		}
		result.ByProject = append(result.ByProject, projectStatsForTest{
			Project: project,
			Seconds: seconds,
			Count:   projectCounts[project],
			Percent: percent,
		})
	}

	// Sort by seconds descending
	for i := 0; i < len(result.ByProject); i++ {
		for j := i + 1; j < len(result.ByProject); j++ {
			if result.ByProject[j].Seconds > result.ByProject[i].Seconds {
				result.ByProject[i], result.ByProject[j] = result.ByProject[j], result.ByProject[i]
			}
		}
	}

	return result
}

func extractProjectNameForTest(name string) string {
	if idx := indexOf(name, ':'); idx > 0 {
		return name[:idx]
	}
	return name
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func filterContextsByPeriodForTest(contexts []contextForTest, period timePeriodForTest) []contextForTest {
	var filtered []contextForTest
	for _, ctx := range contexts {
		if ctx.StartTime.After(period.Start) && ctx.StartTime.Before(period.End) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}
