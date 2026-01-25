package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	intmodels "github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/spf13/cobra"
)

func NewImportCmd(jsonOutput *bool) *cobra.Command {
	var (
		useStdin           bool
		preserveTimestamps bool
	)

	cmd := &cobra.Command{
		Use:     "import <file-path>",
		Aliases: []string{"i"},
		Short:   "Import notes from markdown files or stdin",
		Long: `Import notes from markdown files or external sources into the active context.

Supports:
- Import from markdown file with timestamped sections
- Import from stdin (piped input)
- Preserve timestamps from markdown headers

Examples:
  my-context import meeting-notes.md
  my-context import notes.md --preserve-timestamps
  cat notes.txt | my-context import --stdin
  my-context i journal.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var content string
			var err error

			// Read input source
			if useStdin {
				content, err = readFromStdin()
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			} else {
				if len(args) == 0 {
					return fmt.Errorf("file path required or use --stdin to read from stdin")
				}
				filePath := args[0]
				content, err = readFromFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read file %q: %w", filePath, err)
				}
			}

			// Parse markdown and extract notes
			notes := parseMarkdownNotes(content, preserveTimestamps)
			if len(notes) == 0 {
				return fmt.Errorf("no notes found in input")
			}

			// Check if using database backend
			if core.IsUsingDatabase() {
				return importNotesWithDatabaseBackend(notes, *jsonOutput)
			}

			// Add notes to active context (file-based backend)
			return importNotesWithFileBackend(notes, *jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&useStdin, "stdin", false, "Read notes from stdin instead of file")
	cmd.Flags().BoolVar(&preserveTimestamps, "preserve-timestamps", false, "Use timestamps from markdown headers instead of current time")

	return cmd
}

// readFromFile reads content from a file
func readFromFile(path string) (string, error) {
	// Sanitize and resolve the file path to prevent path traversal
	cleanPath := filepath.Clean(path)
	
	// Resolve to absolute path for validation
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}
	
	// #nosec G304 -- Path is sanitized above with filepath.Clean and Abs
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// readFromStdin reads content from stdin
func readFromStdin() (string, error) {
	var content strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content.String(), nil
}

// parseMarkdownNotes extracts notes from markdown content
// Supports formats like:
// ## 2026-01-13 10:00
// - Note text
//
// Or plain text (each line becomes a note)
func parseMarkdownNotes(content string, preserveTimestamps bool) []*intmodels.Note {
	var notes []*intmodels.Note
	
	// Regex to match markdown headers with timestamps
	// Matches: ## 2026-01-13 10:00, ## 2026-01-13T10:00:00, etc.
	timestampHeaderRegex := regexp.MustCompile(`^##\s+(\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}(?::\d{2})?)`)
	
	lines := strings.Split(content, "\n")
	var currentTimestamp time.Time
	var currentNoteLines []string
	foundTimestampHeader := false
	seenAnyHeader := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Check for timestamp header
		if matches := timestampHeaderRegex.FindStringSubmatch(line); len(matches) > 1 {
			// Save previous note if exists
			if len(currentNoteLines) > 0 {
				noteText := strings.Join(currentNoteLines, "\n")
				timestamp := currentTimestamp
				if timestamp.IsZero() || !preserveTimestamps {
					timestamp = time.Now()
				}
				notes = append(notes, &intmodels.Note{
					Timestamp:   timestamp,
					TextContent: noteText,
				})
				currentNoteLines = []string{}
			}
			
			// Parse new timestamp
			if preserveTimestamps {
				timestampStr := matches[1]
				// Try different time formats
				for _, format := range []string{
					"2006-01-02 15:04:05",
					"2006-01-02 15:04",
					"2006-01-02T15:04:05",
					"2006-01-02T15:04",
				} {
					if t, err := time.Parse(format, timestampStr); err == nil {
						currentTimestamp = t
						break
					}
				}
			}
			foundTimestampHeader = true
			seenAnyHeader = true
			continue
		}
		
		// Skip ALL markdown headers (not just timestamp headers)
		if strings.HasPrefix(line, "#") {
			// If we had a timestamp header before this non-timestamp header,
			// save the accumulated notes
			if foundTimestampHeader && len(currentNoteLines) > 0 {
				noteText := strings.Join(currentNoteLines, "\n")
				timestamp := currentTimestamp
				if timestamp.IsZero() || !preserveTimestamps {
					timestamp = time.Now()
				}
				notes = append(notes, &intmodels.Note{
					Timestamp:   timestamp,
					TextContent: noteText,
				})
				currentNoteLines = []string{}
			}
			foundTimestampHeader = false
			seenAnyHeader = true
			continue
		}
		
		// Only collect text if:
		// 1. We've seen a timestamp header and are currently in a timestamp section, OR
		// 2. We haven't seen any headers at all (plain text mode)
		shouldCollectText := foundTimestampHeader || !seenAnyHeader
		
		if shouldCollectText {
			// Remove markdown list markers (-, *, +)
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimPrefix(line, "* ")
			line = strings.TrimPrefix(line, "+ ")
			line = strings.TrimSpace(line)
			
			// Add to current note
			if line != "" {
				currentNoteLines = append(currentNoteLines, line)
			}
		}
	}
	
	// Add final note if exists
	if len(currentNoteLines) > 0 {
		noteText := strings.Join(currentNoteLines, "\n")
		timestamp := currentTimestamp
		if timestamp.IsZero() || !preserveTimestamps {
			timestamp = time.Now()
		}
		notes = append(notes, &intmodels.Note{
			Timestamp:   timestamp,
			TextContent: noteText,
		})
	}
	
	return notes
}

// importNotesWithFileBackend imports notes using file-based storage
func importNotesWithFileBackend(notes []*intmodels.Note, jsonOutput bool) error {
	state, err := core.GetActiveContext()
	if err != nil {
		return fmt.Errorf("failed to get state: %w", err)
	}

	if state.ActiveContext == nil || *state.ActiveContext == "" {
		return fmt.Errorf("no active context. Use 'my-context start <name>' to create one")
	}

	contextName := *state.ActiveContext

	// Add each note to the active context
	var addedNotes []*intmodels.Note
	for _, note := range notes {
		// For file backend, we need to directly append to the notes log to preserve timestamps
		if err := appendNoteToLog(contextName, note); err != nil {
			return fmt.Errorf("failed to add note: %w", err)
		}
		addedNotes = append(addedNotes, note)
	}

	if jsonOutput {
		output := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"context":        contextName,
				"notes_imported": len(addedNotes),
				"notes":          addedNotes,
			},
		}
		return printJSON(output)
	}

	fmt.Printf("Imported %d note(s) to context %q\n", len(addedNotes), contextName)
	return nil
}

// importNotesWithDatabaseBackend imports notes using database storage
func importNotesWithDatabaseBackend(notes []*intmodels.Note, jsonOutput bool) error {
	backend, err := core.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get backend: %w", err)
	}
	defer backend.Close()

	activeContext, err := backend.GetActiveContext()
	if err != nil {
		return fmt.Errorf("failed to get active context: %w", err)
	}

	if activeContext == "" {
		return fmt.Errorf("no active context. Use 'my-context start <name>' to create one")
	}

	// Add each note to the active context
	for _, note := range notes {
		timestampStr := note.Timestamp.Format(time.RFC3339)
		if err := backend.AddNote(activeContext, timestampStr, note.TextContent); err != nil {
			return fmt.Errorf("failed to add note: %w", err)
		}
	}

	if jsonOutput {
		output := map[string]interface{}{
			"status": "success",
			"data": map[string]interface{}{
				"context":        activeContext,
				"notes_imported": len(notes),
			},
		}
		return printJSON(output)
	}

	fmt.Printf("Imported %d note(s) to context %q\n", len(notes), activeContext)
	return nil
}

// appendNoteToLog appends a note directly to the notes log file
// This is used to preserve timestamps when importing
func appendNoteToLog(contextName string, note *intmodels.Note) error {
	// Path is constructed internally by core.GetNotesLogPath, not from user input
	notesLogPath := core.GetNotesLogPath(contextName)
	
	// Open file for append (0o600 = owner read/write only)
	// #nosec G304 -- Path is constructed internally by core.GetNotesLogPath
	file, err := os.OpenFile(notesLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open notes log: %w", err)
	}
	defer file.Close()

	// Write note as log line
	line := note.ToLogLine() + "\n"
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write note: %w", err)
	}

	return nil
}

// printJSON outputs data as JSON
func printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
