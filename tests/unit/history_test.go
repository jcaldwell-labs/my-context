package unit

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestFilterTransitionsByPeriod tests filtering transitions by time period
func TestFilterTransitionsByPeriod(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	// Create test transitions
	ctx1 := "context-today"
	ctx2 := "context-yesterday"
	ctx3 := "context-lastweek"

	transitions := []*models.ContextTransition{
		{
			Timestamp:      now.Add(-1 * time.Hour),
			NewContext:     &ctx1,
			TransitionType: models.TransitionStart,
		},
		{
			Timestamp:      yesterday,
			NewContext:     &ctx2,
			TransitionType: models.TransitionStart,
		},
		{
			Timestamp:      lastWeek,
			NewContext:     &ctx3,
			TransitionType: models.TransitionStart,
		},
	}

	// Period: only today
	period := periodForTest{
		Start: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		End:   now,
	}

	filtered := filterTransitionsByPeriodForTest(transitions, period)

	assert.Len(t, filtered, 1)
	assert.Equal(t, ctx1, *filtered[0].NewContext)
}

// TestFilterTransitionsByPeriodWeek tests filtering for current week
func TestFilterTransitionsByPeriodWeek(t *testing.T) {
	now := time.Now()
	twoDaysAgo := now.Add(-2 * 24 * time.Hour)
	lastMonth := now.Add(-30 * 24 * time.Hour)

	ctx1 := "recent"
	ctx2 := "old"

	transitions := []*models.ContextTransition{
		{
			Timestamp:      twoDaysAgo,
			NewContext:     &ctx1,
			TransitionType: models.TransitionStart,
		},
		{
			Timestamp:      lastMonth,
			NewContext:     &ctx2,
			TransitionType: models.TransitionStart,
		},
	}

	// Calculate start of this week (Monday)
	y, m, d := now.Date()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	daysBack := weekday - 1 // Days since Monday
	weekStart := time.Date(y, m, d-daysBack, 0, 0, 0, 0, now.Location())

	period := periodForTest{
		Start: weekStart,
		End:   now,
	}

	filtered := filterTransitionsByPeriodForTest(transitions, period)

	// Should only include transition from this week
	assert.Len(t, filtered, 1)
	assert.Equal(t, ctx1, *filtered[0].NewContext)
}

// TestFilterTransitionsByPeriodCustomRange tests --since and --until
func TestFilterTransitionsByPeriodCustomRange(t *testing.T) {
	// Create transitions on specific dates
	jan15 := time.Date(2025, time.January, 15, 12, 0, 0, 0, time.UTC)
	jan20 := time.Date(2025, time.January, 20, 12, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, time.January, 25, 12, 0, 0, 0, time.UTC)

	ctx1 := "context-jan15"
	ctx2 := "context-jan20"
	ctx3 := "context-jan25"

	transitions := []*models.ContextTransition{
		{Timestamp: jan15, NewContext: &ctx1, TransitionType: models.TransitionStart},
		{Timestamp: jan20, NewContext: &ctx2, TransitionType: models.TransitionStart},
		{Timestamp: jan25, NewContext: &ctx3, TransitionType: models.TransitionStart},
	}

	// Period: Jan 18 to Jan 22 (should only include jan20)
	period := periodForTest{
		Start: time.Date(2025, time.January, 18, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2025, time.January, 22, 23, 59, 59, 0, time.UTC),
	}

	filtered := filterTransitionsByPeriodForTest(transitions, period)

	assert.Len(t, filtered, 1)
	assert.Equal(t, ctx2, *filtered[0].NewContext)
}

// TestFilterTransitionsByPeriodEmpty tests filtering with no matches
func TestFilterTransitionsByPeriodEmpty(t *testing.T) {
	lastYear := time.Now().Add(-365 * 24 * time.Hour)
	ctx := "old-context"

	transitions := []*models.ContextTransition{
		{Timestamp: lastYear, NewContext: &ctx, TransitionType: models.TransitionStart},
	}

	// Period: only today
	now := time.Now()
	period := periodForTest{
		Start: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		End:   now,
	}

	filtered := filterTransitionsByPeriodForTest(transitions, period)

	assert.Len(t, filtered, 0)
}

// TestFilterTransitionsByPeriodAllTypes tests filtering with different transition types
func TestFilterTransitionsByPeriodAllTypes(t *testing.T) {
	now := time.Now()
	ctx1 := "context-a"
	ctx2 := "context-b"

	transitions := []*models.ContextTransition{
		{
			Timestamp:      now.Add(-3 * time.Hour),
			NewContext:     &ctx1,
			TransitionType: models.TransitionStart,
		},
		{
			Timestamp:       now.Add(-2 * time.Hour),
			PreviousContext: &ctx1,
			TransitionType:  models.TransitionStop,
		},
		{
			Timestamp:       now.Add(-1 * time.Hour),
			PreviousContext: &ctx1,
			NewContext:      &ctx2,
			TransitionType:  models.TransitionSwitch,
		},
	}

	// Period: today
	period := periodForTest{
		Start: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		End:   now,
	}

	filtered := filterTransitionsByPeriodForTest(transitions, period)

	// All transitions should be included (all happened today)
	assert.Len(t, filtered, 3)
	assert.Equal(t, models.TransitionStart, filtered[0].TransitionType)
	assert.Equal(t, models.TransitionStop, filtered[1].TransitionType)
	assert.Equal(t, models.TransitionSwitch, filtered[2].TransitionType)
}

// --- Test Helper Types ---

type periodForTest struct {
	Start time.Time
	End   time.Time
}

// --- Test Helper Functions ---

func filterTransitionsByPeriodForTest(transitions []*models.ContextTransition, period periodForTest) []*models.ContextTransition {
	var filtered []*models.ContextTransition
	for _, t := range transitions {
		if t.Timestamp.After(period.Start) && t.Timestamp.Before(period.End) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
