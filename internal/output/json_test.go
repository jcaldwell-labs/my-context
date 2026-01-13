package output

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name    string
		command string
		data    any
		wantErr bool
	}{
		{
			name:    "simple data",
			command: "test",
			data:    map[string]string{"key": "value"},
			wantErr: false,
		},
		{
			name:    "nil data",
			command: "test",
			data:    nil,
			wantErr: false,
		},
		{
			name:    "struct data",
			command: "start",
			data: StartData{
				ContextName:  "test-context",
				OriginalName: "test-context",
				WasDuplicate: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatJSON(tt.command, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify it's valid JSON
				var parsed JSONResponse
				if err := json.Unmarshal([]byte(got), &parsed); err != nil {
					t.Errorf("FormatJSON() output is not valid JSON: %v", err)
				}
				// Verify command is set
				if parsed.Command != tt.command {
					t.Errorf("FormatJSON() command = %q, want %q", parsed.Command, tt.command)
				}
				// Verify timestamp is set
				if parsed.Timestamp.IsZero() {
					t.Error("FormatJSON() timestamp should not be zero")
				}
			}
		})
	}
}

func TestFormatJSONError(t *testing.T) {
	tests := []struct {
		name    string
		command string
		code    int
		message string
	}{
		{
			name:    "simple error",
			command: "test",
			code:    1,
			message: "something went wrong",
		},
		{
			name:    "not found error",
			command: "show",
			code:    404,
			message: "context not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatJSONError(tt.command, tt.code, tt.message)
			if err != nil {
				t.Errorf("FormatJSONError() unexpected error: %v", err)
				return
			}

			// Parse the output
			var parsed JSONResponse
			if err := json.Unmarshal([]byte(got), &parsed); err != nil {
				t.Errorf("FormatJSONError() output is not valid JSON: %v", err)
				return
			}

			// Verify command
			if parsed.Command != tt.command {
				t.Errorf("FormatJSONError() command = %q, want %q", parsed.Command, tt.command)
			}

			// Verify error is set
			if parsed.Error == nil {
				t.Error("FormatJSONError() error should not be nil")
				return
			}

			if parsed.Error.Code != tt.code {
				t.Errorf("FormatJSONError() error.code = %d, want %d", parsed.Error.Code, tt.code)
			}
			if parsed.Error.Message != tt.message {
				t.Errorf("FormatJSONError() error.message = %q, want %q", parsed.Error.Message, tt.message)
			}

			// Verify data is nil (not set)
			if parsed.Data != nil {
				t.Error("FormatJSONError() data should be nil")
			}
		})
	}
}

func TestFormatExportJSON(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	ctx := &models.Context{
		Name:      "test-context",
		StartTime: pastTime,
		EndTime:   &now,
		Status:    "stopped",
	}

	notes := []models.Note{
		{Timestamp: now, TextContent: "note 1"},
		{Timestamp: now, TextContent: "note 2"},
	}

	files := []models.FileAssociation{
		{Timestamp: now, FilePath: "/path/to/file.txt"},
	}

	got, err := FormatExportJSON(ctx, notes, files, 5)
	if err != nil {
		t.Errorf("FormatExportJSON() unexpected error: %v", err)
		return
	}

	// Parse and verify
	var parsed ExportData
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Errorf("FormatExportJSON() output is not valid JSON: %v", err)
		return
	}

	if parsed.Name != ctx.Name {
		t.Errorf("FormatExportJSON() name = %q, want %q", parsed.Name, ctx.Name)
	}
	if len(parsed.Notes) != 2 {
		t.Errorf("FormatExportJSON() notes count = %d, want 2", len(parsed.Notes))
	}
	if len(parsed.Files) != 1 {
		t.Errorf("FormatExportJSON() files count = %d, want 1", len(parsed.Files))
	}
	if parsed.TouchCount != 5 {
		t.Errorf("FormatExportJSON() touch_count = %d, want 5", parsed.TouchCount)
	}
}

func TestExtractPeriod(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		got := extractPeriod(nil)
		if !got.Start.IsZero() || !got.End.IsZero() || got.Label != "" {
			t.Error("extractPeriod(nil) should return zero StatsPeriod")
		}
	})

	t.Run("non-struct input", func(t *testing.T) {
		got := extractPeriod("not a struct")
		if !got.Start.IsZero() {
			t.Error("extractPeriod(string) should return zero StatsPeriod")
		}
	})

	t.Run("struct with matching fields", func(t *testing.T) {
		now := time.Now()
		input := struct {
			Start time.Time
			End   time.Time
			Label string
		}{
			Start: now,
			End:   now.Add(time.Hour),
			Label: "test",
		}

		got := extractPeriod(input)
		if !got.Start.Equal(now) {
			t.Errorf("extractPeriod() Start = %v, want %v", got.Start, now)
		}
		if got.Label != "test" {
			t.Errorf("extractPeriod() Label = %q, want %q", got.Label, "test")
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		now := time.Now()
		input := &struct {
			Start time.Time
			End   time.Time
			Label string
		}{
			Start: now,
			Label: "ptr",
		}

		got := extractPeriod(input)
		if !got.Start.Equal(now) {
			t.Error("extractPeriod() should handle pointer input")
		}
	})
}

func TestExtractProjectStats(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		got := extractProjectStats(nil)
		if got != nil {
			t.Error("extractProjectStats(nil) should return nil")
		}
	})

	t.Run("non-slice input", func(t *testing.T) {
		got := extractProjectStats("not a slice")
		if got != nil {
			t.Error("extractProjectStats(string) should return nil")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []struct {
			Project string
			Seconds int64
		}{}
		got := extractProjectStats(input)
		if len(got) != 0 {
			t.Error("extractProjectStats([]) should return empty slice")
		}
	})

	t.Run("slice with items", func(t *testing.T) {
		type ProjectItem struct {
			Project string
			Seconds int64
			Count   int
			Percent float64
		}
		input := []ProjectItem{
			{Project: "ps-cli", Seconds: 3600, Count: 5, Percent: 50.0},
			{Project: "garden", Seconds: 1800, Count: 3, Percent: 25.0},
		}

		got := extractProjectStats(input)
		if len(got) != 2 {
			t.Errorf("extractProjectStats() len = %d, want 2", len(got))
			return
		}
		if got[0].Project != "ps-cli" {
			t.Errorf("extractProjectStats()[0].Project = %q, want %q", got[0].Project, "ps-cli")
		}
		if got[0].Seconds != 3600 {
			t.Errorf("extractProjectStats()[0].Seconds = %d, want 3600", got[0].Seconds)
		}
	})
}

func TestExtractSession(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		got := extractSession(nil)
		if got != nil {
			t.Error("extractSession(nil) should return nil")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var input *struct{ Name string }
		got := extractSession(input)
		if got != nil {
			t.Error("extractSession(nil ptr) should return nil")
		}
	})

	t.Run("non-struct input", func(t *testing.T) {
		got := extractSession("not a struct")
		if got != nil {
			t.Error("extractSession(string) should return nil")
		}
	})

	t.Run("struct with empty name", func(t *testing.T) {
		input := struct {
			Name    string
			Seconds int64
		}{}
		got := extractSession(input)
		if got != nil {
			t.Error("extractSession() with empty name should return nil")
		}
	})

	t.Run("struct with valid data", func(t *testing.T) {
		input := struct {
			Name    string
			Seconds int64
		}{Name: "test-session", Seconds: 3600}

		got := extractSession(input)
		if got == nil {
			t.Fatal("extractSession() should return non-nil")
		}
		if got.Name != "test-session" {
			t.Errorf("extractSession().Name = %q, want %q", got.Name, "test-session")
		}
		if got.Seconds != 3600 {
			t.Errorf("extractSession().Seconds = %d, want 3600", got.Seconds)
		}
	})
}

func TestJSONOutputHasNewline(t *testing.T) {
	got, err := FormatJSON("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Error("FormatJSON() output should end with newline")
	}

	got, err = FormatJSONError("test", 1, "error")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Error("FormatJSONError() output should end with newline")
	}
}
