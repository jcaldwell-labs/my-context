package output

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"seconds", 30 * time.Second, "30s"},
		{"one minute", 1 * time.Minute, "1m"},
		{"minutes", 45 * time.Minute, "45m"},
		{"one hour", 1 * time.Hour, "1h 0m"},
		{"hours and minutes", 2*time.Hour + 30*time.Minute, "2h 30m"},
		{"one day", 24 * time.Hour, "1d 0h"},
		{"days and hours", 48*time.Hour + 5*time.Hour, "2d 5h"},
		{"zero duration", 0, "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("FormatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

func TestFormatSimpleMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{"simple", "Hello", "Hello\n"},
		{"empty", "", "\n"},
		{"with spaces", "Hello World", "Hello World\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSimpleMessage(tt.message)
			if got != tt.want {
				t.Errorf("FormatSimpleMessage(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	err := fmt.Errorf("no active context")
	got := FormatError(err)
	want := "Error: no active context\n"

	if got != want {
		t.Errorf("FormatError() = %q, want %q", got, want)
	}
}

func TestGetTimestampFormat(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		want   string
	}{
		{"default empty", "", time.RFC3339},
		{"iso", "iso", time.RFC3339},
		{"short", "short", "15:04"},
		{"medium", "medium", "15:04:05"},
		{"long", "long", "2006-01-02 15:04:05"},
		{"custom format", "Jan 02", "Jan 02"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MC_TIMESTAMP_FORMAT", tt.envVal)
			got := getTimestampFormat()
			if got != tt.want {
				t.Errorf("getTimestampFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractContextInfo(t *testing.T) {
	now := time.Now()

	t.Run("pkg models context with metadata", func(t *testing.T) {
		ctx := &pkgmodels.ContextWithMetadata{
			Name:      "test-context",
			Status:    "active",
			StartTime: now,
			Metadata: pkgmodels.ContextMetadata{
				CreatedBy: "test",
			},
		}

		info, err := extractContextInfo(ctx)
		if err != nil {
			t.Errorf("extractContextInfo() unexpected error: %v", err)
		}
		if info.name != "test-context" {
			t.Errorf("extractContextInfo() name = %q, want %q", info.name, "test-context")
		}
		if !info.hasMetadata {
			t.Error("extractContextInfo() hasMetadata should be true")
		}
	})

	t.Run("models context", func(t *testing.T) {
		ctx := &models.Context{
			Name:      "test-context",
			Status:    "stopped",
			StartTime: now,
		}

		info, err := extractContextInfo(ctx)
		if err != nil {
			t.Errorf("extractContextInfo() unexpected error: %v", err)
		}
		if info.name != "test-context" {
			t.Errorf("extractContextInfo() name = %q, want %q", info.name, "test-context")
		}
		if info.hasMetadata {
			t.Error("extractContextInfo() hasMetadata should be false")
		}
	})

	t.Run("unknown type returns error", func(t *testing.T) {
		_, err := extractContextInfo("not a context")
		if err == nil {
			t.Error("extractContextInfo() expected error for unknown type")
		}
	})
}

func TestFormatMetadataSection(t *testing.T) {
	t.Run("no metadata", func(t *testing.T) {
		got := formatMetadataSection(false, nil)
		if got != "" {
			t.Errorf("formatMetadataSection(false) = %q, want empty", got)
		}
	})

	t.Run("empty metadata", func(t *testing.T) {
		meta := pkgmodels.ContextMetadata{}
		got := formatMetadataSection(true, meta)
		if got != "" {
			t.Errorf("formatMetadataSection(empty) = %q, want empty", got)
		}
	})

	t.Run("full metadata", func(t *testing.T) {
		meta := pkgmodels.ContextMetadata{
			CreatedBy: "user",
			Parent:    "parent-ctx",
			Labels:    []string{"tag1", "tag2"},
		}
		got := formatMetadataSection(true, meta)

		if !strings.Contains(got, "Created by: user") {
			t.Error("formatMetadataSection() should contain created by")
		}
		if !strings.Contains(got, "Parent: parent-ctx") {
			t.Error("formatMetadataSection() should contain parent")
		}
		if !strings.Contains(got, "Labels: tag1, tag2") {
			t.Error("formatMetadataSection() should contain labels")
		}
	})
}

func TestFormatNotesSection(t *testing.T) {
	t.Run("empty notes", func(t *testing.T) {
		got := formatNotesSection([]*models.Note{})
		if !strings.Contains(got, "(none)") {
			t.Error("formatNotesSection([]) should contain (none)")
		}
		if !strings.Contains(got, "Notes (0)") {
			t.Error("formatNotesSection([]) should show count 0")
		}
	})

	t.Run("with notes", func(t *testing.T) {
		notes := []*models.Note{
			{Timestamp: time.Now(), TextContent: "First note"},
			{Timestamp: time.Now(), TextContent: "Second note"},
		}
		got := formatNotesSection(notes)

		if !strings.Contains(got, "Notes (2)") {
			t.Error("formatNotesSection() should show count 2")
		}
		if !strings.Contains(got, "First note") {
			t.Error("formatNotesSection() should contain first note")
		}
		if !strings.Contains(got, "Second note") {
			t.Error("formatNotesSection() should contain second note")
		}
	})
}

func TestFormatFilesSection(t *testing.T) {
	t.Run("empty files", func(t *testing.T) {
		got := formatFilesSection([]*models.FileAssociation{})
		if !strings.Contains(got, "(none)") {
			t.Error("formatFilesSection([]) should contain (none)")
		}
		if !strings.Contains(got, "Files (0)") {
			t.Error("formatFilesSection([]) should show count 0")
		}
	})

	t.Run("with files", func(t *testing.T) {
		files := []*models.FileAssociation{
			{Timestamp: time.Now(), FilePath: "/path/to/file.txt"},
		}
		got := formatFilesSection(files)

		if !strings.Contains(got, "Files (1)") {
			t.Error("formatFilesSection() should show count 1")
		}
		if !strings.Contains(got, "/path/to/file.txt") {
			t.Error("formatFilesSection() should contain file path")
		}
	})
}

func TestFormatActivitySection(t *testing.T) {
	t.Run("no touches", func(t *testing.T) {
		got := formatActivitySection([]*models.TouchEvent{})
		if !strings.Contains(got, "Activity: 0 touches") {
			t.Error("formatActivitySection([]) should show 0 touches")
		}
	})

	t.Run("with touches", func(t *testing.T) {
		touches := []*models.TouchEvent{
			{Timestamp: time.Now()},
			{Timestamp: time.Now()},
			{Timestamp: time.Now()},
		}
		got := formatActivitySection(touches)

		if !strings.Contains(got, "Activity: 3 touches") {
			t.Error("formatActivitySection() should show 3 touches")
		}
		if !strings.Contains(got, "last:") {
			t.Error("formatActivitySection() should contain last touch info")
		}
	})
}

func TestFormatContext(t *testing.T) {
	now := time.Now()

	t.Run("models context", func(t *testing.T) {
		ctx := &models.Context{
			Name:      "test-context",
			Status:    "active",
			StartTime: now.Add(-1 * time.Hour),
		}
		notes := []*models.Note{{Timestamp: now, TextContent: "test"}}
		files := []*models.FileAssociation{}
		touches := []*models.TouchEvent{}

		got := FormatContext(ctx, notes, files, touches)

		if !strings.Contains(got, "Context: test-context") {
			t.Error("FormatContext() should contain context name")
		}
		if !strings.Contains(got, "Status: active") {
			t.Error("FormatContext() should contain status")
		}
		if !strings.Contains(got, "Notes (1)") {
			t.Error("FormatContext() should contain notes section")
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		got := FormatContext("invalid", nil, nil, nil)
		if !strings.Contains(got, "Error: Unknown context type") {
			t.Error("FormatContext() should show error for unknown type")
		}
	})
}

func TestFormatContextList(t *testing.T) {
	now := time.Now()

	t.Run("empty list", func(t *testing.T) {
		got := FormatContextList([]*models.Context{}, "")

		if !strings.Contains(got, "Contexts (0)") {
			t.Error("FormatContextList([]) should show 0 count")
		}
		if !strings.Contains(got, "No contexts found") {
			t.Error("FormatContextList([]) should show no contexts message")
		}
		if !strings.Contains(got, "my-context start") {
			t.Error("FormatContextList([]) should show start hint")
		}
	})

	t.Run("with contexts", func(t *testing.T) {
		endTime := now
		contexts := []*models.Context{
			{Name: "active-ctx", Status: "active", StartTime: now.Add(-1 * time.Hour)},
			{Name: "stopped-ctx", Status: "stopped", StartTime: now.Add(-2 * time.Hour), EndTime: &endTime},
		}

		got := FormatContextList(contexts, "active-ctx")

		if !strings.Contains(got, "Contexts (2)") {
			t.Error("FormatContextList() should show count 2")
		}
		if !strings.Contains(got, "● active-ctx") {
			t.Error("FormatContextList() should mark active context with filled circle")
		}
		if !strings.Contains(got, "○ stopped-ctx") {
			t.Error("FormatContextList() should mark inactive context with empty circle")
		}
		if !strings.Contains(got, "Duration:") {
			t.Error("FormatContextList() should show duration for stopped contexts")
		}
	})
}

func TestFormatTransitionHistory(t *testing.T) {
	now := time.Now()

	t.Run("empty history", func(t *testing.T) {
		got := FormatTransitionHistory([]*models.ContextTransition{})

		if !strings.Contains(got, "Context History") {
			t.Error("FormatTransitionHistory() should have header")
		}
		if !strings.Contains(got, "No transitions recorded") {
			t.Error("FormatTransitionHistory([]) should show no transitions message")
		}
	})

	t.Run("with transitions", func(t *testing.T) {
		ctx1 := "context-1"
		ctx2 := "context-2"
		transitions := []*models.ContextTransition{
			{Timestamp: now, TransitionType: models.TransitionStart, NewContext: &ctx1},
			{Timestamp: now, TransitionType: models.TransitionSwitch, PreviousContext: &ctx1, NewContext: &ctx2},
			{Timestamp: now, TransitionType: models.TransitionStop, PreviousContext: &ctx2},
		}

		got := FormatTransitionHistory(transitions)

		if !strings.Contains(got, "START") {
			t.Error("FormatTransitionHistory() should contain START")
		}
		if !strings.Contains(got, "SWITCH") {
			t.Error("FormatTransitionHistory() should contain SWITCH")
		}
		if !strings.Contains(got, "STOP") {
			t.Error("FormatTransitionHistory() should contain STOP")
		}
		if !strings.Contains(got, "context-1 → context-2") {
			t.Error("FormatTransitionHistory() should show switch arrow")
		}
	})
}
