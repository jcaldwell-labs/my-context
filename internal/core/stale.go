package core

import (
	"os"
	"strconv"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// StaleLevel represents the staleness severity
type StaleLevel int

const (
	// StaleLevelNone means context is not stale
	StaleLevelNone StaleLevel = iota
	// StaleLevelWarn means context should show a warning
	StaleLevelWarn
	// StaleLevelStale means context is marked as stale
	StaleLevelStale
	// StaleLevelCritical means context should suggest stopping
	StaleLevelCritical
)

// StaleThresholds holds the time thresholds for staleness
type StaleThresholds struct {
	Warn     time.Duration
	Stale    time.Duration
	Critical time.Duration
}

// GetDefaultStaleThresholds returns the default thresholds
func GetDefaultStaleThresholds() StaleThresholds {
	return StaleThresholds{
		Warn:     4 * time.Hour,
		Stale:    8 * time.Hour,
		Critical: 24 * time.Hour,
	}
}

// GetStaleThresholds returns thresholds from environment or defaults
func GetStaleThresholds() StaleThresholds {
	defaults := GetDefaultStaleThresholds()

	// Check environment variables for custom thresholds
	if warnStr := os.Getenv("MC_STALE_WARN_HOURS"); warnStr != "" {
		if hours, err := strconv.ParseFloat(warnStr, 64); err == nil && hours > 0 {
			defaults.Warn = time.Duration(hours * float64(time.Hour))
		}
	}

	if staleStr := os.Getenv("MC_STALE_HOURS"); staleStr != "" {
		if hours, err := strconv.ParseFloat(staleStr, 64); err == nil && hours > 0 {
			defaults.Stale = time.Duration(hours * float64(time.Hour))
		}
	}

	if criticalStr := os.Getenv("MC_STALE_CRITICAL_HOURS"); criticalStr != "" {
		if hours, err := strconv.ParseFloat(criticalStr, 64); err == nil && hours > 0 {
			defaults.Critical = time.Duration(hours * float64(time.Hour))
		}
	}

	return defaults
}

// GetLastActivityTime returns the most recent activity timestamp for a context
// Activity includes: notes, file associations, touch events
// Returns the context start time if no activity exists
func GetLastActivityTime(ctx *models.Context, notes []*models.Note, files []*models.FileAssociation, touches []*models.TouchEvent) time.Time {
	lastActivity := ctx.StartTime

	// Check notes
	for _, note := range notes {
		if note.Timestamp.After(lastActivity) {
			lastActivity = note.Timestamp
		}
	}

	// Check files
	for _, file := range files {
		if file.Timestamp.After(lastActivity) {
			lastActivity = file.Timestamp
		}
	}

	// Check touches
	for _, touch := range touches {
		if touch.Timestamp.After(lastActivity) {
			lastActivity = touch.Timestamp
		}
	}

	return lastActivity
}

// GetStaleLevel calculates the staleness level based on time since last activity
func GetStaleLevel(lastActivity time.Time, thresholds StaleThresholds) StaleLevel {
	timeSinceActivity := time.Since(lastActivity)

	if timeSinceActivity >= thresholds.Critical {
		return StaleLevelCritical
	}
	if timeSinceActivity >= thresholds.Stale {
		return StaleLevelStale
	}
	if timeSinceActivity >= thresholds.Warn {
		return StaleLevelWarn
	}

	return StaleLevelNone
}

// IsStale returns true if the context is considered stale (at stale level or higher)
func IsStale(lastActivity time.Time, thresholds StaleThresholds) bool {
	level := GetStaleLevel(lastActivity, thresholds)
	return level >= StaleLevelStale
}

// GetStaleLevelString returns a human-readable string for the stale level
func GetStaleLevelString(level StaleLevel) string {
	switch level {
	case StaleLevelCritical:
		return "critical"
	case StaleLevelStale:
		return "stale"
	case StaleLevelWarn:
		return "warn"
	default:
		return "none"
	}
}
