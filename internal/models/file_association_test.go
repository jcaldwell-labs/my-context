package models

import (
	"strings"
	"testing"
	"time"
)

func TestFileAssociationValidate(t *testing.T) {
	tests := []struct {
		name    string
		fa      FileAssociation
		wantErr string
	}{
		{
			name:    "valid file",
			fa:      FileAssociation{Timestamp: time.Now(), FilePath: "/path/to/file.txt"},
			wantErr: "",
		},
		{
			name:    "empty path",
			fa:      FileAssociation{Timestamp: time.Now(), FilePath: ""},
			wantErr: "file path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fa.Validate()
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

func TestFileAssociationToLogLine(t *testing.T) {
	timestamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	fa := FileAssociation{
		Timestamp: timestamp,
		FilePath:  "/path/to/file.txt",
	}

	got := fa.ToLogLine()
	want := "2025-01-15T10:30:00Z|/path/to/file.txt"

	if got != want {
		t.Errorf("ToLogLine() = %q, want %q", got, want)
	}
}

func TestParseFileLogLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantErr  bool
		wantPath string
	}{
		{
			name:     "valid line",
			line:     "2025-01-15T10:30:00Z|/path/to/file.txt",
			wantErr:  false,
			wantPath: "/path/to/file.txt",
		},
		{
			name:    "missing pipe",
			line:    "2025-01-15T10:30:00Z/path/to/file.txt",
			wantErr: true,
		},
		{
			name:    "invalid timestamp",
			line:    "not-a-timestamp|/path/to/file.txt",
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
			fa, err := ParseFileLogLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFileLogLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fa.FilePath != tt.wantPath {
				t.Errorf("ParseFileLogLine() path = %q, want %q", fa.FilePath, tt.wantPath)
			}
		})
	}
}

func TestFileAssociationRoundtrip(t *testing.T) {
	timestamp := time.Now().UTC().Truncate(time.Second)
	original := FileAssociation{
		Timestamp: timestamp,
		FilePath:  "/path/to/file.go",
	}

	line := original.ToLogLine()
	parsed, err := ParseFileLogLine(line)

	if err != nil {
		t.Fatalf("Roundtrip failed to parse: %v", err)
	}

	if parsed.FilePath != original.FilePath {
		t.Errorf("Roundtrip path mismatch: %q vs %q", parsed.FilePath, original.FilePath)
	}

	if !parsed.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Roundtrip timestamp mismatch: %v vs %v", parsed.Timestamp, original.Timestamp)
	}
}
