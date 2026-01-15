package unit

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/stretchr/testify/assert"
)

// TestTimelineEntryBasic tests basic timeline entry creation
func TestTimelineEntryBasic(t *testing.T) {
	now := time.Now()
	entry := output.TimelineEntry{
		ContextName: "test-context",
		StartTime:   now,
		EndTime:     nil,
		Duration:    time.Hour,
		IsGap:       false,
		IsActive:    true,
	}

	assert.Equal(t, "test-context", entry.ContextName)
	assert.Equal(t, now, entry.StartTime)
	assert.Nil(t, entry.EndTime)
	assert.Equal(t, time.Hour, entry.Duration)
	assert.False(t, entry.IsGap)
	assert.True(t, entry.IsActive)
}

// TestTimelineEntryGap tests gap entry creation
func TestTimelineEntryGap(t *testing.T) {
	now := time.Now()
	endTime := now.Add(30 * time.Minute)
	entry := output.TimelineEntry{
		ContextName: "(no context)",
		StartTime:   now,
		EndTime:     &endTime,
		Duration:    30 * time.Minute,
		IsGap:       true,
		IsActive:    false,
	}

	assert.Equal(t, "(no context)", entry.ContextName)
	assert.True(t, entry.IsGap)
	assert.False(t, entry.IsActive)
	assert.Equal(t, 30*time.Minute, entry.Duration)
}

// TestTimelineDataEmpty tests empty timeline data
func TestTimelineDataEmpty(t *testing.T) {
	timeline := output.TimelineData{
		Period:       "Today",
		Entries:      []output.TimelineEntry{},
		TotalTracked: 0,
		TotalGaps:    0,
	}

	assert.Equal(t, 0, len(timeline.Entries))
	assert.Equal(t, time.Duration(0), timeline.TotalTracked)
	assert.Equal(t, time.Duration(0), timeline.TotalGaps)
}

// TestTimelineDataWithEntries tests timeline data with multiple entries
func TestTimelineDataWithEntries(t *testing.T) {
	now := time.Now()
	entry1EndTime := now.Add(2 * time.Hour)
	entry2StartTime := now.Add(3 * time.Hour)
	entry2EndTime := now.Add(4 * time.Hour)

	entries := []output.TimelineEntry{
		{
			ContextName: "context-1",
			StartTime:   now,
			EndTime:     &entry1EndTime,
			Duration:    2 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
		{
			ContextName: "(no context)",
			StartTime:   entry1EndTime,
			EndTime:     &entry2StartTime,
			Duration:    1 * time.Hour,
			IsGap:       true,
			IsActive:    false,
		},
		{
			ContextName: "context-2",
			StartTime:   entry2StartTime,
			EndTime:     &entry2EndTime,
			Duration:    1 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
	}

	timeline := output.TimelineData{
		Period:       "Today",
		Entries:      entries,
		TotalTracked: 3 * time.Hour,
		TotalGaps:    1 * time.Hour,
	}

	assert.Equal(t, 3, len(timeline.Entries))
	assert.Equal(t, 3*time.Hour, timeline.TotalTracked)
	assert.Equal(t, 1*time.Hour, timeline.TotalGaps)
}

// TestDayGroupCreation tests day group creation
func TestDayGroupCreation(t *testing.T) {
	now := time.Now()
	y, m, d := now.Date()
	dayStart := time.Date(y, m, d, 0, 0, 0, 0, now.Location())

	entries := []output.TimelineEntry{
		{
			ContextName: "morning-work",
			StartTime:   dayStart.Add(8 * time.Hour),
			EndTime:     nil,
			Duration:    2 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
	}

	dayGroup := output.DayGroup{
		Date:         dayStart,
		Entries:      entries,
		TotalTracked: 2 * time.Hour,
		TotalGaps:    0,
	}

	assert.Equal(t, dayStart, dayGroup.Date)
	assert.Equal(t, 1, len(dayGroup.Entries))
	assert.Equal(t, 2*time.Hour, dayGroup.TotalTracked)
	assert.Equal(t, time.Duration(0), dayGroup.TotalGaps)
}

// TestFormatTimelineEmpty tests formatting empty timeline
func TestFormatTimelineEmpty(t *testing.T) {
	timeline := &output.TimelineData{
		Period:       "Today",
		Entries:      []output.TimelineEntry{},
		TotalTracked: 0,
		TotalGaps:    0,
	}

	result := output.FormatTimeline(timeline)
	assert.Contains(t, result, "Timeline Visualization")
	assert.Contains(t, result, "No contexts found in this period")
}

// TestFormatTimelineWithEntries tests formatting timeline with entries
func TestFormatTimelineWithEntries(t *testing.T) {
	now := time.Now()
	endTime := now.Add(2 * time.Hour)

	entries := []output.TimelineEntry{
		{
			ContextName: "test-context",
			StartTime:   now,
			EndTime:     &endTime,
			Duration:    2 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
	}

	timeline := &output.TimelineData{
		Period:       "Today",
		Entries:      entries,
		TotalTracked: 2 * time.Hour,
		TotalGaps:    0,
	}

	result := output.FormatTimeline(timeline)
	assert.Contains(t, result, "Timeline Visualization")
	assert.Contains(t, result, "test-context")
	assert.Contains(t, result, "Total tracked: 2h")
}

// TestFormatTimelineJSON tests JSON formatting
func TestFormatTimelineJSON(t *testing.T) {
	now := time.Now()
	endTime := now.Add(1 * time.Hour)

	entries := []output.TimelineEntry{
		{
			ContextName: "test-context",
			StartTime:   now,
			EndTime:     &endTime,
			Duration:    1 * time.Hour,
			IsGap:       false,
			IsActive:    true,
		},
	}

	timeline := &output.TimelineData{
		Period:       "Today",
		Entries:      entries,
		TotalTracked: 1 * time.Hour,
		TotalGaps:    0,
	}

	jsonStr, err := output.FormatTimelineJSON(timeline)
	assert.NoError(t, err)
	assert.Contains(t, jsonStr, "\"status\": \"success\"")
	assert.Contains(t, jsonStr, "\"context_name\": \"test-context\"")
	assert.Contains(t, jsonStr, "\"is_active\": true")
	assert.Contains(t, jsonStr, "\"duration_seconds\": 3600")
}

// TestFormatTimelineJSONWithGap tests JSON formatting with gap
func TestFormatTimelineJSONWithGap(t *testing.T) {
	now := time.Now()
	entry1EndTime := now.Add(1 * time.Hour)
	entry2StartTime := now.Add(2 * time.Hour)
	entry2EndTime := now.Add(3 * time.Hour)

	entries := []output.TimelineEntry{
		{
			ContextName: "context-1",
			StartTime:   now,
			EndTime:     &entry1EndTime,
			Duration:    1 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
		{
			ContextName: "(no context)",
			StartTime:   entry1EndTime,
			EndTime:     &entry2StartTime,
			Duration:    1 * time.Hour,
			IsGap:       true,
			IsActive:    false,
		},
		{
			ContextName: "context-2",
			StartTime:   entry2StartTime,
			EndTime:     &entry2EndTime,
			Duration:    1 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
	}

	timeline := &output.TimelineData{
		Period:       "Today",
		Entries:      entries,
		TotalTracked: 2 * time.Hour,
		TotalGaps:    1 * time.Hour,
	}

	jsonStr, err := output.FormatTimelineJSON(timeline)
	assert.NoError(t, err)
	assert.Contains(t, jsonStr, "\"status\": \"success\"")
	assert.Contains(t, jsonStr, "\"is_gap\": true")
	assert.Contains(t, jsonStr, "\"total_tracked_seconds\": 7200")
	assert.Contains(t, jsonStr, "\"total_gaps_seconds\": 3600")
}

// TestFormatTimelineWithDayGroups tests formatting with day groups
func TestFormatTimelineWithDayGroups(t *testing.T) {
	now := time.Now()
	y, m, d := now.Date()
	dayStart := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	endTime := dayStart.Add(10 * time.Hour)

	entries := []output.TimelineEntry{
		{
			ContextName: "morning-work",
			StartTime:   dayStart.Add(8 * time.Hour),
			EndTime:     &endTime,
			Duration:    2 * time.Hour,
			IsGap:       false,
			IsActive:    false,
		},
	}

	dayGroup := output.DayGroup{
		Date:         dayStart,
		Entries:      entries,
		TotalTracked: 2 * time.Hour,
		TotalGaps:    0,
	}

	timeline := &output.TimelineData{
		Period:       "This Week",
		Entries:      entries,
		TotalTracked: 2 * time.Hour,
		TotalGaps:    0,
		DayGroups:    []output.DayGroup{dayGroup},
	}

	result := output.FormatTimeline(timeline)
	assert.Contains(t, result, "Timeline Visualization")
	assert.Contains(t, result, dayStart.Format("2006-01-02"))
	assert.Contains(t, result, "morning-work")
	assert.Contains(t, result, "Tracked: 2h")
}
