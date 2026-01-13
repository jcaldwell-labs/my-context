package models

import (
	"strings"
	"testing"
	"time"
)

func TestNoteValidate(t *testing.T) {
	tests := []struct {
		name    string
		note    Note
		wantErr string
	}{
		{
			name: "valid note",
			note: Note{
				Timestamp:   time.Now(),
				TextContent: "This is a valid note",
			},
			wantErr: "",
		},
		{
			name: "empty text",
			note: Note{
				Timestamp:   time.Now(),
				TextContent: "",
			},
			wantErr: "note text cannot be empty",
		},
		{
			name: "text too long",
			note: Note{
				Timestamp:   time.Now(),
				TextContent: strings.Repeat("a", 10001),
			},
			wantErr: "10,000 characters or less",
		},
		{
			name: "max length valid",
			note: Note{
				Timestamp:   time.Now(),
				TextContent: strings.Repeat("a", 10000),
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.note.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Validate() error = %q, want containing %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestNoteEscape(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"no escaping needed", "simple text", "simple text"},
		{"escape pipe", "before|after", "before\\|after"},
		{"escape newline", "line1\nline2", "line1\\nline2"},
		{"escape backslash", "path\\to\\file", "path\\\\to\\\\file"},
		{"multiple pipes", "a|b|c", "a\\|b\\|c"},
		{"mixed escapes", "a|b\nc\\d", "a\\|b\\nc\\\\d"},
		{"backslash before pipe", "\\|", "\\\\\\|"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{TextContent: tt.text}
			got := note.Escape()
			if got != tt.want {
				t.Errorf("Escape() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnescapeNote(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"no unescaping needed", "simple text", "simple text"},
		{"unescape pipe", "before\\|after", "before|after"},
		{"unescape newline", "line1\\nline2", "line1\nline2"},
		{"unescape backslash", "path\\\\to\\\\file", "path\\to\\file"},
		{"multiple pipes", "a\\|b\\|c", "a|b|c"},
		{"mixed escapes", "a\\|b\\nc\\\\d", "a|b\nc\\d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnescapeNote(tt.text)
			if got != tt.want {
				t.Errorf("UnescapeNote() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEscapeUnescapeRoundtrip(t *testing.T) {
	// Note: Windows paths like "C:\Users\name" have a known limitation
	// where "\n" (backslash-n) is ambiguous with escaped newline.
	// This is an edge case in the simple escape scheme.
	texts := []string{
		"simple text",
		"text with | pipe",
		"text with\nnewline",
		"text with \\ backslash",
		"complex: |, \n, \\, all",
		"Special chars: $100 @email #tag",
	}

	for _, text := range texts {
		t.Run(text[:min(20, len(text))], func(t *testing.T) {
			note := Note{TextContent: text}
			escaped := note.Escape()
			unescaped := UnescapeNote(escaped)

			if unescaped != text {
				t.Errorf("Roundtrip failed:\n  original:  %q\n  escaped:   %q\n  unescaped: %q", text, escaped, unescaped)
			}
		})
	}
}

func TestNoteToLogLine(t *testing.T) {
	timestamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	note := Note{
		Timestamp:   timestamp,
		TextContent: "test note",
	}

	got := note.ToLogLine()
	want := "2025-01-15T10:30:00Z|test note"

	if got != want {
		t.Errorf("ToLogLine() = %q, want %q", got, want)
	}
}

func TestParseNoteLogLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
		wantTxt string
	}{
		{
			name:    "valid line",
			line:    "2025-01-15T10:30:00Z|test note",
			wantErr: false,
			wantTxt: "test note",
		},
		{
			name:    "line with escaped pipe",
			line:    "2025-01-15T10:30:00Z|before\\|after",
			wantErr: false,
			wantTxt: "before|after",
		},
		{
			name:    "missing pipe",
			line:    "2025-01-15T10:30:00Z-no pipe",
			wantErr: true,
		},
		{
			name:    "invalid timestamp",
			line:    "not-a-timestamp|text",
			wantErr: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note, err := ParseNoteLogLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNoteLogLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && note.TextContent != tt.wantTxt {
				t.Errorf("ParseNoteLogLine() text = %q, want %q", note.TextContent, tt.wantTxt)
			}
		})
	}
}

func TestToLogLineParsRoundtrip(t *testing.T) {
	texts := []string{
		"simple text",
		"text with | pipe",
		"text with\nnewline",
		"text with \\ backslash",
	}

	timestamp := time.Now().UTC().Truncate(time.Second)

	for _, text := range texts {
		t.Run(text[:min(20, len(text))], func(t *testing.T) {
			original := Note{
				Timestamp:   timestamp,
				TextContent: text,
			}

			logLine := original.ToLogLine()
			parsed, err := ParseNoteLogLine(logLine)

			if err != nil {
				t.Errorf("Roundtrip failed to parse: %v", err)
				return
			}

			if parsed.TextContent != original.TextContent {
				t.Errorf("Roundtrip text mismatch:\n  original: %q\n  parsed:   %q", original.TextContent, parsed.TextContent)
			}

			if !parsed.Timestamp.Equal(original.Timestamp) {
				t.Errorf("Roundtrip timestamp mismatch:\n  original: %v\n  parsed:   %v", original.Timestamp, parsed.Timestamp)
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
