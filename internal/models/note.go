package models

import (
	"fmt"
	"strings"
	"time"
)

// Note represents a timestamped text entry associated with a context
type Note struct {
	Timestamp   time.Time `json:"timestamp"`
	TextContent string    `json:"text_content"`
}

// Validate checks if the note has valid data
func (n *Note) Validate() error {
	if n.TextContent == "" {
		return fmt.Errorf("note text cannot be empty")
	}

	if len(n.TextContent) > 10000 {
		return fmt.Errorf("note text must be 10,000 characters or less")
	}

	return nil
}

// Escape escapes special characters for log file storage
// Pipes are escaped as \| and newlines as \n
func (n *Note) Escape() string {
	text := n.TextContent
	text = strings.ReplaceAll(text, "\\", "\\\\") // Escape backslashes first
	text = strings.ReplaceAll(text, "|", "\\|")   // Escape pipes
	text = strings.ReplaceAll(text, "\n", "\\n")  // Escape newlines
	return text
}

// Unescape unescapes special characters from log file storage
func UnescapeNote(text string) string {
	result := strings.ReplaceAll(text, "\\n", "\n")   // Unescape newlines
	result = strings.ReplaceAll(result, "\\|", "|")   // Unescape pipes
	result = strings.ReplaceAll(result, "\\\\", "\\") // Unescape backslashes last
	return result
}

// ToLogLine formats the note as a log line: timestamp|text
func (n *Note) ToLogLine() string {
	return fmt.Sprintf("%s|%s", n.Timestamp.Format(time.RFC3339), n.Escape())
}

// ParseNoteLogLine parses a log line into a Note
func ParseNoteLogLine(line string) (*Note, error) {
	parts := strings.SplitN(line, "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid note log line format")
	}

	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	return &Note{
		Timestamp:   timestamp,
		TextContent: UnescapeNote(parts[1]),
	}, nil
}
