package unit

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultStaleThresholds(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()

	assert.Equal(t, 4*time.Hour, thresholds.Warn, "Warn threshold should be 4 hours")
	assert.Equal(t, 8*time.Hour, thresholds.Stale, "Stale threshold should be 8 hours")
	assert.Equal(t, 24*time.Hour, thresholds.Critical, "Critical threshold should be 24 hours")
}

func TestGetStaleThresholds_WithEnvironmentVariables(t *testing.T) {
	// Use t.Setenv to avoid race conditions
	t.Setenv("MC_STALE_WARN_HOURS", "2")
	t.Setenv("MC_STALE_HOURS", "6")
	t.Setenv("MC_STALE_CRITICAL_HOURS", "12")

	thresholds := core.GetStaleThresholds()

	assert.Equal(t, 2*time.Hour, thresholds.Warn)
	assert.Equal(t, 6*time.Hour, thresholds.Stale)
	assert.Equal(t, 12*time.Hour, thresholds.Critical)
}

func TestGetStaleThresholds_WithInvalidEnvironmentVariables(t *testing.T) {
	t.Setenv("MC_STALE_WARN_HOURS", "invalid")
	t.Setenv("MC_STALE_HOURS", "-5")
	t.Setenv("MC_STALE_CRITICAL_HOURS", "0")

	thresholds := core.GetStaleThresholds()

	// Should fall back to defaults
	assert.Equal(t, 4*time.Hour, thresholds.Warn)
	assert.Equal(t, 8*time.Hour, thresholds.Stale)
	assert.Equal(t, 24*time.Hour, thresholds.Critical)
}

func TestGetLastActivityTime_NoActivity(t *testing.T) {
	now := time.Now()
	ctx := &models.Context{
		Name:      "test-context",
		StartTime: now,
		Status:    "active",
	}

	lastActivity := core.GetLastActivityTime(ctx, nil, nil, nil)

	assert.Equal(t, now, lastActivity, "Should return start time when no activity exists")
}

func TestGetLastActivityTime_WithNotes(t *testing.T) {
	startTime := time.Now().Add(-2 * time.Hour)
	noteTime := time.Now().Add(-1 * time.Hour)

	ctx := &models.Context{
		Name:      "test-context",
		StartTime: startTime,
		Status:    "active",
	}

	notes := []*models.Note{
		{Timestamp: noteTime, TextContent: "A note"},
	}

	lastActivity := core.GetLastActivityTime(ctx, notes, nil, nil)

	assert.Equal(t, noteTime, lastActivity, "Should return most recent note time")
}

func TestGetLastActivityTime_WithFiles(t *testing.T) {
	startTime := time.Now().Add(-2 * time.Hour)
	fileTime := time.Now().Add(-30 * time.Minute)

	ctx := &models.Context{
		Name:      "test-context",
		StartTime: startTime,
		Status:    "active",
	}

	files := []*models.FileAssociation{
		{Timestamp: fileTime, FilePath: "/test/file.txt"},
	}

	lastActivity := core.GetLastActivityTime(ctx, nil, files, nil)

	assert.Equal(t, fileTime, lastActivity, "Should return most recent file time")
}

func TestGetLastActivityTime_WithTouches(t *testing.T) {
	startTime := time.Now().Add(-2 * time.Hour)
	touchTime := time.Now().Add(-15 * time.Minute)

	ctx := &models.Context{
		Name:      "test-context",
		StartTime: startTime,
		Status:    "active",
	}

	touches := []*models.TouchEvent{
		{Timestamp: touchTime},
	}

	lastActivity := core.GetLastActivityTime(ctx, nil, nil, touches)

	assert.Equal(t, touchTime, lastActivity, "Should return most recent touch time")
}

func TestGetLastActivityTime_WithMixedActivity(t *testing.T) {
	startTime := time.Now().Add(-5 * time.Hour)
	noteTime := time.Now().Add(-3 * time.Hour)
	fileTime := time.Now().Add(-2 * time.Hour)
	touchTime := time.Now().Add(-1 * time.Hour)

	ctx := &models.Context{
		Name:      "test-context",
		StartTime: startTime,
		Status:    "active",
	}

	notes := []*models.Note{
		{Timestamp: noteTime, TextContent: "A note"},
	}
	files := []*models.FileAssociation{
		{Timestamp: fileTime, FilePath: "/test/file.txt"},
	}
	touches := []*models.TouchEvent{
		{Timestamp: touchTime},
	}

	lastActivity := core.GetLastActivityTime(ctx, notes, files, touches)

	assert.Equal(t, touchTime, lastActivity, "Should return most recent activity across all types")
}

func TestGetStaleLevel_None(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	recentActivity := time.Now().Add(-2 * time.Hour) // 2 hours ago (< 4h warn threshold)

	level := core.GetStaleLevel(recentActivity, thresholds)

	assert.Equal(t, core.StaleLevelNone, level, "Should not be stale with recent activity")
}

func TestGetStaleLevel_Warn(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	warnActivity := time.Now().Add(-5 * time.Hour) // 5 hours ago (>= 4h, < 8h)

	level := core.GetStaleLevel(warnActivity, thresholds)

	assert.Equal(t, core.StaleLevelWarn, level, "Should be at warn level")
}

func TestGetStaleLevel_Stale(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	staleActivity := time.Now().Add(-10 * time.Hour) // 10 hours ago (>= 8h, < 24h)

	level := core.GetStaleLevel(staleActivity, thresholds)

	assert.Equal(t, core.StaleLevelStale, level, "Should be at stale level")
}

func TestGetStaleLevel_Critical(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	criticalActivity := time.Now().Add(-25 * time.Hour) // 25 hours ago (>= 24h)

	level := core.GetStaleLevel(criticalActivity, thresholds)

	assert.Equal(t, core.StaleLevelCritical, level, "Should be at critical level")
}

func TestIsStale_NotStale(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	recentActivity := time.Now().Add(-2 * time.Hour)

	isStale := core.IsStale(recentActivity, thresholds)

	assert.False(t, isStale, "Should not be stale with recent activity")
}

func TestIsStale_Warn(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	warnActivity := time.Now().Add(-5 * time.Hour)

	isStale := core.IsStale(warnActivity, thresholds)

	assert.False(t, isStale, "Warn level should not be considered stale")
}

func TestIsStale_Stale(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	staleActivity := time.Now().Add(-10 * time.Hour)

	isStale := core.IsStale(staleActivity, thresholds)

	assert.True(t, isStale, "Should be stale at stale level")
}

func TestIsStale_Critical(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()
	criticalActivity := time.Now().Add(-25 * time.Hour)

	isStale := core.IsStale(criticalActivity, thresholds)

	assert.True(t, isStale, "Should be stale at critical level")
}

func TestGetStaleLevelString(t *testing.T) {
	tests := []struct {
		level    core.StaleLevel
		expected string
	}{
		{core.StaleLevelNone, "none"},
		{core.StaleLevelWarn, "warn"},
		{core.StaleLevelStale, "stale"},
		{core.StaleLevelCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := core.GetStaleLevelString(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStaleDetection_BoundaryConditions(t *testing.T) {
	thresholds := core.GetDefaultStaleThresholds()

	// Test exact threshold boundaries
	tests := []struct {
		name          string
		timeSince     time.Duration
		expectedLevel core.StaleLevel
	}{
		{"Just before warn", 4*time.Hour - 1*time.Second, core.StaleLevelNone},
		{"Exactly at warn", 4 * time.Hour, core.StaleLevelWarn},
		{"Just after warn", 4*time.Hour + 1*time.Second, core.StaleLevelWarn},
		{"Just before stale", 8*time.Hour - 1*time.Second, core.StaleLevelWarn},
		{"Exactly at stale", 8 * time.Hour, core.StaleLevelStale},
		{"Just after stale", 8*time.Hour + 1*time.Second, core.StaleLevelStale},
		{"Just before critical", 24*time.Hour - 1*time.Second, core.StaleLevelStale},
		{"Exactly at critical", 24 * time.Hour, core.StaleLevelCritical},
		{"Just after critical", 24*time.Hour + 1*time.Second, core.StaleLevelCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activityTime := time.Now().Add(-tt.timeSince)
			level := core.GetStaleLevel(activityTime, thresholds)
			assert.Equal(t, tt.expectedLevel, level)
		})
	}
}

func TestStaleDetection_WithCustomThresholds(t *testing.T) {
	customThresholds := core.StaleThresholds{
		Warn:     1 * time.Hour,
		Stale:    2 * time.Hour,
		Critical: 4 * time.Hour,
	}

	tests := []struct {
		name          string
		timeSince     time.Duration
		expectedLevel core.StaleLevel
	}{
		{"Recent", 30 * time.Minute, core.StaleLevelNone},
		{"Warn level", 90 * time.Minute, core.StaleLevelWarn},
		{"Stale level", 3 * time.Hour, core.StaleLevelStale},
		{"Critical level", 5 * time.Hour, core.StaleLevelCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activityTime := time.Now().Add(-tt.timeSince)
			level := core.GetStaleLevel(activityTime, customThresholds)
			assert.Equal(t, tt.expectedLevel, level)
		})
	}
}

func TestGetStaleThresholds_DefaultsWhenNoEnvVars(t *testing.T) {
	// Use t.Setenv with empty strings to clear env vars; auto-restores after test
	t.Setenv("MC_STALE_WARN_HOURS", "")
	t.Setenv("MC_STALE_HOURS", "")
	t.Setenv("MC_STALE_CRITICAL_HOURS", "")

	thresholds := core.GetStaleThresholds()

	assert.Equal(t, 4*time.Hour, thresholds.Warn)
	assert.Equal(t, 8*time.Hour, thresholds.Stale)
	assert.Equal(t, 24*time.Hour, thresholds.Critical)
}
