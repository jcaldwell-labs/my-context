package models

import (
	"testing"
	"time"
)

func TestTouchEventToLogLine(t *testing.T) {
	timestamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	te := TouchEvent{Timestamp: timestamp}

	got := te.ToLogLine()
	want := "2025-01-15T10:30:00Z"

	if got != want {
		t.Errorf("ToLogLine() = %q, want %q", got, want)
	}
}

func TestParseTouchLogLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{
			name:    "valid timestamp",
			line:    "2025-01-15T10:30:00Z",
			wantErr: false,
		},
		{
			name:    "invalid timestamp",
			line:    "not-a-timestamp",
			wantErr: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "wrong format",
			line:    "2025-01-15 10:30:00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te, err := ParseTouchLogLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTouchLogLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && te == nil {
				t.Error("ParseTouchLogLine() returned nil for valid input")
			}
		})
	}
}

func TestTouchEventRoundtrip(t *testing.T) {
	timestamp := time.Now().UTC().Truncate(time.Second)
	original := TouchEvent{Timestamp: timestamp}

	line := original.ToLogLine()
	parsed, err := ParseTouchLogLine(line)

	if err != nil {
		t.Fatalf("Roundtrip failed to parse: %v", err)
	}

	if !parsed.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Roundtrip timestamp mismatch: %v vs %v", parsed.Timestamp, original.Timestamp)
	}
}
