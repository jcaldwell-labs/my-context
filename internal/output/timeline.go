package output

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

// TimelineEntry represents a single entry in the timeline
type TimelineEntry struct {
	ContextName string        `json:"context_name"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	Duration    time.Duration `json:"duration_seconds"`
	IsGap       bool          `json:"is_gap"`
	IsActive    bool          `json:"is_active"`
}

// TimelineData holds all timeline information
type TimelineData struct {
	Period       interface{}     `json:"period"`
	Entries      []TimelineEntry `json:"entries"`
	TotalTracked time.Duration   `json:"total_tracked_seconds"`
	TotalGaps    time.Duration   `json:"total_gaps_seconds"`
	DayGroups    []DayGroup      `json:"day_groups,omitempty"`
}

// DayGroup represents timeline entries grouped by day
type DayGroup struct {
	Date         time.Time       `json:"date"`
	Entries      []TimelineEntry `json:"entries"`
	TotalTracked time.Duration   `json:"total_tracked_seconds"`
	TotalGaps    time.Duration   `json:"total_gaps_seconds"`
}

// FormatTimeline formats timeline data for human-readable output
func FormatTimeline(timeline *TimelineData) string {
	var sb strings.Builder

	sb.WriteString("Timeline Visualization\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Period header
	periodLabel := getPeriodLabel(timeline.Period)
	sb.WriteString(fmt.Sprintf("Period: %s\n", periodLabel))
	sb.WriteString("\n")

	if len(timeline.Entries) == 0 {
		sb.WriteString("No contexts found in this period\n")
		return sb.String()
	}

	// If multi-day, show grouped by day
	if len(timeline.DayGroups) > 0 {
		for _, dayGroup := range timeline.DayGroups {
			sb.WriteString(formatDayGroup(&dayGroup))
			sb.WriteString("\n")
		}
	} else {
		// Single day - show full timeline
		sb.WriteString(formatSingleDayTimeline(timeline.Entries))
		sb.WriteString("\n")
	}

	// Summary
	sb.WriteString(fmt.Sprintf("Total tracked: %s\n", FormatDuration(timeline.TotalTracked)))
	if timeline.TotalGaps > 0 {
		sb.WriteString(fmt.Sprintf("Total gaps: %s\n", FormatDuration(timeline.TotalGaps)))
	}

	return sb.String()
}

// formatDayGroup formats a single day's timeline
func formatDayGroup(dayGroup *DayGroup) string {
	var sb strings.Builder

	// Date header
	sb.WriteString(fmt.Sprintf("%s\n", dayGroup.Date.Format("2006-01-02 (Monday)")))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Visual timeline bar
	sb.WriteString(formatVisualBar(dayGroup.Entries, dayGroup.Date))
	sb.WriteString("\n")

	// Entry list
	for _, entry := range dayGroup.Entries {
		sb.WriteString(formatTimelineEntry(&entry))
	}

	// Day summary
	sb.WriteString(fmt.Sprintf("\nTracked: %s", FormatDuration(dayGroup.TotalTracked)))
	if dayGroup.TotalGaps > 0 {
		sb.WriteString(fmt.Sprintf(" | Gaps: %s", FormatDuration(dayGroup.TotalGaps)))
	}
	sb.WriteString("\n")

	return sb.String()
}

// formatSingleDayTimeline formats a single day timeline
func formatSingleDayTimeline(entries []TimelineEntry) string {
	var sb strings.Builder

	if len(entries) == 0 {
		return ""
	}

	// Determine the day from first entry
	dayStart := entries[0].StartTime
	y, m, d := dayStart.Date()
	dayDate := time.Date(y, m, d, 0, 0, 0, 0, dayStart.Location())

	// Visual timeline bar
	sb.WriteString(formatVisualBar(entries, dayDate))
	sb.WriteString("\n")

	// Entry list
	for _, entry := range entries {
		sb.WriteString(formatTimelineEntry(&entry))
	}

	return sb.String()
}

// formatVisualBar creates the ASCII timeline bar
func formatVisualBar(entries []TimelineEntry, dayDate time.Time) string {
	const barWidth = 48 // Characters for the bar
	const hoursInDay = 24.0

	var sb strings.Builder

	// Time scale header
	sb.WriteString("08:00 ┃")
	sb.WriteString(strings.Repeat("─", barWidth))
	sb.WriteString("┃ 18:00\n")

	// Calculate bar positions
	var bar strings.Builder
	bar.WriteString("      ┃")

	// Initialize bar with spaces
	barChars := make([]rune, barWidth)
	for i := range barChars {
		barChars[i] = '░'
	}

	// Fill in the bar based on entries
	for _, entry := range entries {
		startHour := float64(entry.StartTime.Hour()) + float64(entry.StartTime.Minute())/60.0

		var endHour float64
		if entry.EndTime != nil {
			endHour = float64(entry.EndTime.Hour()) + float64(entry.EndTime.Minute())/60.0
			// Handle entries that span midnight
			if entry.EndTime.Day() != entry.StartTime.Day() {
				endHour = 24.0
			}
		} else {
			// Active context - extends to now or end of day
			now := time.Now()
			if now.Day() == entry.StartTime.Day() && now.Month() == entry.StartTime.Month() && now.Year() == entry.StartTime.Year() {
				endHour = float64(now.Hour()) + float64(now.Minute())/60.0
			} else {
				endHour = 24.0
			}
		}

		// Clamp to visible range (08:00 - 18:00)
		const visibleStart = 8.0
		const visibleEnd = 18.0
		const visibleRange = visibleEnd - visibleStart

		if endHour <= visibleStart || startHour >= visibleEnd {
			continue // Outside visible range
		}

		// Clamp to visible window
		displayStart := math.Max(startHour, visibleStart)
		displayEnd := math.Min(endHour, visibleEnd)

		// Calculate positions in bar (0 to barWidth)
		startPos := int((displayStart - visibleStart) / visibleRange * float64(barWidth))
		endPos := int((displayEnd - visibleStart) / visibleRange * float64(barWidth))

		if startPos < 0 {
			startPos = 0
		}
		if endPos > barWidth {
			endPos = barWidth
		}

		// Fill the bar
		char := '█'
		if entry.IsGap {
			char = '░'
		}

		for i := startPos; i < endPos && i < barWidth; i++ {
			barChars[i] = char
		}
	}

	bar.WriteString(string(barChars))
	bar.WriteString("┃\n")

	sb.WriteString(bar.String())

	return sb.String()
}

// formatTimelineEntry formats a single timeline entry
func formatTimelineEntry(entry *TimelineEntry) string {
	startTime := entry.StartTime.Format("15:04")

	var endTime string
	if entry.EndTime != nil {
		endTime = entry.EndTime.Format("15:04")
	} else {
		endTime = "now  "
	}

	// Format the bar
	var bar string
	if entry.IsGap {
		bar = "░░░░"
	} else {
		bar = "████"
	}

	// Active indicator
	activeIndicator := ""
	if entry.IsActive {
		activeIndicator = " ●"
	}

	// Truncate context name if too long
	contextName := entry.ContextName
	if len(contextName) > 30 {
		contextName = contextName[:27] + "..."
	}

	return fmt.Sprintf("%s-%s %s %s (%s)%s\n",
		startTime,
		endTime,
		bar,
		contextName,
		FormatDuration(entry.Duration),
		activeIndicator)
}

// getPeriodLabel extracts the label from the period
func getPeriodLabel(period interface{}) string {
	// Handle different period types
	if p, ok := period.(map[string]interface{}); ok {
		if label, ok := p["Label"].(string); ok {
			return label
		}
		if label, ok := p["label"].(string); ok {
			return label
		}
	}

	// Try type assertion for struct with Label field
	type periodWithLabel interface {
		GetLabel() string
	}
	if p, ok := period.(periodWithLabel); ok {
		return p.GetLabel()
	}

	// Fallback - try reflection for Label field
	return "Custom Period"
}

// FormatTimelineJSON formats timeline data as JSON
func FormatTimelineJSON(timeline *TimelineData) (string, error) {
	// Convert durations to seconds for JSON
	type jsonEntry struct {
		ContextName     string  `json:"context_name"`
		StartTime       string  `json:"start_time"`
		EndTime         *string `json:"end_time,omitempty"`
		DurationSeconds float64 `json:"duration_seconds"`
		IsGap           bool    `json:"is_gap"`
		IsActive        bool    `json:"is_active"`
	}

	type jsonDayGroup struct {
		Date                string      `json:"date"`
		Entries             []jsonEntry `json:"entries"`
		TotalTrackedSeconds float64     `json:"total_tracked_seconds"`
		TotalGapsSeconds    float64     `json:"total_gaps_seconds"`
	}

	type jsonTimeline struct {
		Status              string         `json:"status"`
		Period              interface{}    `json:"period"`
		Entries             []jsonEntry    `json:"entries"`
		TotalTrackedSeconds float64        `json:"total_tracked_seconds"`
		TotalGapsSeconds    float64        `json:"total_gaps_seconds"`
		DayGroups           []jsonDayGroup `json:"day_groups,omitempty"`
	}

	// Convert entries
	var jsonEntries []jsonEntry
	for _, e := range timeline.Entries {
		var endTimeStr *string
		if e.EndTime != nil {
			s := e.EndTime.Format(time.RFC3339)
			endTimeStr = &s
		}

		jsonEntries = append(jsonEntries, jsonEntry{
			ContextName:     e.ContextName,
			StartTime:       e.StartTime.Format(time.RFC3339),
			EndTime:         endTimeStr,
			DurationSeconds: e.Duration.Seconds(),
			IsGap:           e.IsGap,
			IsActive:        e.IsActive,
		})
	}

	// Convert day groups
	var jsonDayGroups []jsonDayGroup
	for _, dg := range timeline.DayGroups {
		var dgEntries []jsonEntry
		for _, e := range dg.Entries {
			var endTimeStr *string
			if e.EndTime != nil {
				s := e.EndTime.Format(time.RFC3339)
				endTimeStr = &s
			}

			dgEntries = append(dgEntries, jsonEntry{
				ContextName:     e.ContextName,
				StartTime:       e.StartTime.Format(time.RFC3339),
				EndTime:         endTimeStr,
				DurationSeconds: e.Duration.Seconds(),
				IsGap:           e.IsGap,
				IsActive:        e.IsActive,
			})
		}

		jsonDayGroups = append(jsonDayGroups, jsonDayGroup{
			Date:                dg.Date.Format("2006-01-02"),
			Entries:             dgEntries,
			TotalTrackedSeconds: dg.TotalTracked.Seconds(),
			TotalGapsSeconds:    dg.TotalGaps.Seconds(),
		})
	}

	result := jsonTimeline{
		Status:              "success",
		Period:              timeline.Period,
		Entries:             jsonEntries,
		TotalTrackedSeconds: timeline.TotalTracked.Seconds(),
		TotalGapsSeconds:    timeline.TotalGaps.Seconds(),
		DayGroups:           jsonDayGroups,
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal timeline JSON: %w", err)
	}

	return string(bytes) + "\n", nil
}
