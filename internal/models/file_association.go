package models

import (
	"fmt"
	"strings"
	"time"
)

// FileAssociation represents a file path reference linked to a context
type FileAssociation struct {
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"` // Stored in POSIX format
}

// Validate checks if the file association has valid data
func (f *FileAssociation) Validate() error {
	if f.FilePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	return nil
}

// ToLogLine formats the file association as a log line: timestamp|path
func (f *FileAssociation) ToLogLine() string {
	return fmt.Sprintf("%s|%s", f.Timestamp.Format(time.RFC3339), f.FilePath)
}

// ParseFileLogLine parses a log line into a FileAssociation
func ParseFileLogLine(line string) (*FileAssociation, error) {
	parts := strings.SplitN(line, "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid file log line format")
	}

	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	return &FileAssociation{
		Timestamp: timestamp,
		FilePath:  parts[1],
	}, nil
}
