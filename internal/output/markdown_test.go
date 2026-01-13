package output

import (
	"strings"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

func TestFormatExportMarkdown(t *testing.T) {
	now := time.Now()
	endTime := now

	t.Run("active context with data", func(t *testing.T) {
		ctx := &models.Context{
			Name:      "test-context",
			Status:    "active",
			StartTime: now.Add(-1 * time.Hour),
		}
		notes := []models.Note{
			{Timestamp: now, TextContent: "Note 1"},
			{Timestamp: now, TextContent: "Note 2"},
		}
		files := []models.FileAssociation{
			{Timestamp: now, FilePath: "/path/to/file.go"},
		}

		got := FormatExportMarkdown(ctx, notes, files, 5)

		// Check header
		if !strings.Contains(got, "# Context: test-context") {
			t.Error("FormatExportMarkdown() should contain context header")
		}

		// Check active status
		if !strings.Contains(got, "Currently Active") {
			t.Error("FormatExportMarkdown() should show Currently Active for active context")
		}

		// Check notes
		if !strings.Contains(got, "## Notes") {
			t.Error("FormatExportMarkdown() should contain Notes section")
		}
		if !strings.Contains(got, "Note 1") {
			t.Error("FormatExportMarkdown() should contain note text")
		}
		if !strings.Contains(got, "Total: 2 notes") {
			t.Error("FormatExportMarkdown() should show note count")
		}

		// Check files
		if !strings.Contains(got, "## Files") {
			t.Error("FormatExportMarkdown() should contain Files section")
		}
		if !strings.Contains(got, "/path/to/file.go") {
			t.Error("FormatExportMarkdown() should contain file path")
		}

		// Check activity
		if !strings.Contains(got, "## Activity") {
			t.Error("FormatExportMarkdown() should contain Activity section")
		}
		if !strings.Contains(got, "Total: 5 touches") {
			t.Error("FormatExportMarkdown() should show touch count")
		}

		// Check footer
		if !strings.Contains(got, "my-context") {
			t.Error("FormatExportMarkdown() should contain footer")
		}
	})

	t.Run("stopped context", func(t *testing.T) {
		ctx := &models.Context{
			Name:      "stopped-context",
			Status:    "stopped",
			StartTime: now.Add(-2 * time.Hour),
			EndTime:   &endTime,
		}

		got := FormatExportMarkdown(ctx, nil, nil, 0)

		if !strings.Contains(got, "**Ended**:") {
			t.Error("FormatExportMarkdown() should show end time")
		}
		if strings.Contains(got, "Currently Active") {
			t.Error("FormatExportMarkdown() should not show Currently Active for stopped context")
		}
	})

	t.Run("archived context", func(t *testing.T) {
		ctx := &models.Context{
			Name:       "archived-context",
			Status:     "stopped",
			StartTime:  now.Add(-3 * time.Hour),
			EndTime:    &endTime,
			IsArchived: true,
		}

		got := FormatExportMarkdown(ctx, nil, nil, 0)

		if !strings.Contains(got, "archived") {
			t.Error("FormatExportMarkdown() should mention archived status")
		}
	})

	t.Run("empty notes and files", func(t *testing.T) {
		ctx := &models.Context{
			Name:      "empty-context",
			Status:    "active",
			StartTime: now,
		}

		got := FormatExportMarkdown(ctx, []models.Note{}, []models.FileAssociation{}, 0)

		// Check notes section shows (none)
		notesIdx := strings.Index(got, "## Notes")
		filesIdx := strings.Index(got, "## Files")
		notesSection := got[notesIdx:filesIdx]
		if !strings.Contains(notesSection, "(none)") {
			t.Error("FormatExportMarkdown() should show (none) for empty notes")
		}

		// Check files section shows (none)
		activityIdx := strings.Index(got, "## Activity")
		filesSection := got[filesIdx:activityIdx]
		if !strings.Contains(filesSection, "(none)") {
			t.Error("FormatExportMarkdown() should show (none) for empty files")
		}

		// Check activity section shows (none) for 0 touches
		activitySection := got[activityIdx:]
		if !strings.Contains(activitySection, "(none)") {
			t.Error("FormatExportMarkdown() should show (none) for 0 touches")
		}
	})
}

func TestMarkdownFormatLocalTime(t *testing.T) {
	// Test with a known time
	testTime := time.Date(2025, 6, 15, 14, 30, 0, 0, time.UTC)
	got := formatLocalTime(testTime)

	// Should be formatted as human-readable date
	// The exact output depends on local timezone, but should contain key elements
	if !strings.Contains(got, "2025") {
		t.Errorf("formatLocalTime() should contain year, got %q", got)
	}
	if !strings.Contains(got, "June") || !strings.Contains(got, "15") {
		t.Errorf("formatLocalTime() should contain month and day, got %q", got)
	}
}

func TestMarkdownFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"seconds", 30 * time.Second, "30s"},
		{"minutes only", 45 * time.Minute, "45m"},
		{"hours and minutes", 2*time.Hour + 15*time.Minute, "2h 15m"},
		{"zero", 0, "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}
