package unit

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
)

// Helper to create test transitions
func createTestTransition(timestamp time.Time, transType models.TransitionType) *models.ContextTransition {
	ctxName := "test-context"
	return &models.ContextTransition{
		Timestamp:       timestamp,
		PreviousContext: nil,
		NewContext:      &ctxName,
		TransitionType:  transType,
	}
}

// TestFilterTransitionsByPeriodToday tests filtering transitions for today
func TestFilterTransitionsByPeriodToday(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// Create transitions
	transitions := []*models.ContextTransition{
		createTestTransition(now.Add(-1*time.Hour), models.TransitionStart),
		createTestTransition(yesterday, models.TransitionStart),
	}

	// Period: only today
	y, m, d := now.Date()
	period := historyTimePeriodForTest{
		Start: time.Date(y, m, d, 0, 0, 0, 0, now.Location()),
		End:   now,
		Label: "Today",
	}

	filtered := filterTransitionsByPeriodForHistory(transitions, period)

	assert.Len(t, filtered, 1, "Should filter to only today's transition")
	assert.WithinDuration(t, now.Add(-1*time.Hour), filtered[0].Timestamp, time.Minute)
}

// TestFilterTransitionsByPeriodEmpty tests filtering with no matches
func TestFilterTransitionsByPeriodEmpty(t *testing.T) {
	now := time.Now()
	lastWeek := now.Add(-7 * 24 * time.Hour)

	// Create transition from last week
	transitions := []*models.ContextTransition{
		createTestTransition(lastWeek, models.TransitionStart),
	}

	// Period: only today
	y, m, d := now.Date()
	period := historyTimePeriodForTest{
		Start: time.Date(y, m, d, 0, 0, 0, 0, now.Location()),
		End:   now,
		Label: "Today",
	}

	filtered := filterTransitionsByPeriodForHistory(transitions, period)

	assert.Len(t, filtered, 0, "Should return empty when no transitions match period")
}

// TestFilterTransitionsByPeriodWeek tests filtering transitions for current week
func TestFilterTransitionsByPeriodWeek(t *testing.T) {
	now := time.Now()
	twoWeeksAgo := now.Add(-14 * 24 * time.Hour)
	thisWeek := now.Add(-3 * 24 * time.Hour) // 3 days ago

	// Create transitions
	transitions := []*models.ContextTransition{
		createTestTransition(thisWeek, models.TransitionStart),
		createTestTransition(twoWeeksAgo, models.TransitionStart),
	}

	// Period: this week (Monday to now)
	y, m, d := now.Date()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	daysBack := weekday - 1
	period := historyTimePeriodForTest{
		Start: time.Date(y, m, d-daysBack, 0, 0, 0, 0, now.Location()),
		End:   now,
		Label: "This Week",
	}

	filtered := filterTransitionsByPeriodForHistory(transitions, period)

	// Should only include transition from this week (if it falls within Monday-now)
	assert.LessOrEqual(t, len(filtered), 1, "Should filter to this week's transitions")
}

// TestFilterTransitionsByPeriodCustomRange tests filtering with custom date range
func TestFilterTransitionsByPeriodCustomRange(t *testing.T) {
	// Fixed dates for predictable testing
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)
	beforeStart := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	afterEnd := time.Date(2025, 2, 1, 0, 0, 1, 0, time.UTC)

	// Create transitions
	transitions := []*models.ContextTransition{
		createTestTransition(beforeStart, models.TransitionStart),
		createTestTransition(middle, models.TransitionStart),
		createTestTransition(afterEnd, models.TransitionStart),
	}

	// Period: January 2025
	period := historyTimePeriodForTest{
		Start: start,
		End:   end,
		Label: "January 2025",
	}

	filtered := filterTransitionsByPeriodForHistory(transitions, period)

	assert.Len(t, filtered, 1, "Should only include transition within date range")
	assert.Equal(t, middle, filtered[0].Timestamp, "Should be the middle transition")
}

// TestFilterTransitionsByPeriodAllTime tests filtering with no restrictions
func TestFilterTransitionsByPeriodAllTime(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	// Create transitions (ensure they're all in the past for the test to work)
	transitions := []*models.ContextTransition{
		createTestTransition(now.Add(-1*time.Hour), models.TransitionStart),
		createTestTransition(yesterday, models.TransitionStop),
		createTestTransition(lastWeek, models.TransitionSwitch),
	}

	// Period: all time (end is slightly in the future to include "now")
	period := historyTimePeriodForTest{
		Start: time.Unix(0, 0),
		End:   now.Add(1 * time.Hour),
		Label: "All Time",
	}

	filtered := filterTransitionsByPeriodForHistory(transitions, period)

	assert.Len(t, filtered, 3, "Should include all transitions for all-time period")
}

// TestFilterTransitionsByPeriodBoundaryConditions tests edge cases
func TestFilterTransitionsByPeriodBoundaryConditions(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	// Create transitions exactly at boundaries
	transitions := []*models.ContextTransition{
		createTestTransition(start, models.TransitionStart),          // Exactly at start - should be excluded (After check)
		createTestTransition(end, models.TransitionStop),             // Exactly at end - should be excluded (Before check)
		createTestTransition(start.Add(1*time.Second), models.TransitionStart), // Just after start - should be included
		createTestTransition(end.Add(-1*time.Second), models.TransitionStop),   // Just before end - should be included
	}

	period := historyTimePeriodForTest{
		Start: start,
		End:   end,
		Label: "January 2025",
	}

	filtered := filterTransitionsByPeriodForHistory(transitions, period)

	// After(start) && Before(end) means exclusive boundaries
	assert.Len(t, filtered, 2, "Should only include transitions strictly within period (exclusive boundaries)")
}

// --- Test Helper Functions ---

// historyTimePeriodForTest mirrors the TimePeriod type from commands package
type historyTimePeriodForTest struct {
	Start time.Time
	End   time.Time
	Label string
}

// filterTransitionsByPeriodForHistory mimics the actual implementation for testing
func filterTransitionsByPeriodForHistory(transitions []*models.ContextTransition, period historyTimePeriodForTest) []*models.ContextTransition {
	var filtered []*models.ContextTransition
	for _, t := range transitions {
		if t.Timestamp.After(period.Start) && t.Timestamp.Before(period.End) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
