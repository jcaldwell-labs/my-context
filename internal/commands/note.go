package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewNoteCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "note <text>",
		Aliases: []string{"n"},
		Short:   "Add a note to the active context",
		Long:    `Add a timestamped note to the currently active context.`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Join all args as the note text
			noteText := ""
			for i, arg := range args {
				if i > 0 {
					noteText += " "
				}
				noteText += arg
			}

			// Get active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("note", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			if !state.HasActiveContext() {
				errMsg := "No active context. Start a context with: my-context start <name>"
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("note", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			}

			// Check note count thresholds and show warnings (before adding note)
			contextName := state.GetActiveContextName()
			currentCount, err := core.GetNoteCount(contextName)
			if err == nil { // Continue even if we can't get count
				ShowNoteWarning(currentCount)
			}

			// Add the note
			note, err := core.AddNote(noteText)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("note", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.NoteData{
					ContextName:   state.GetActiveContextName(),
					NoteTimestamp: note.Timestamp,
					NoteText:      note.TextContent,
				}
				jsonStr, err := output.FormatJSON("note", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Note added to context: %s\n", state.GetActiveContextName())
			}

			return nil
		},
	}

	return cmd
}

// GetEnvInt reads an integer environment variable with a default value
func GetEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// ShowNoteWarning displays warnings when note count reaches thresholds
func ShowNoteWarning(currentCount int) {
	// Read threshold environment variables
	warnAt1 := GetEnvInt("MC_WARN_AT", 50)
	warnAt2 := GetEnvInt("MC_WARN_AT_2", 100)
	warnAt3 := GetEnvInt("MC_WARN_AT_3", 200)

	newCount := currentCount + 1 // Count after adding this note

	// Display appropriate warning based on new count
	switch {
	case newCount == warnAt1:
		fmt.Printf("⚠️  Context now has %d notes. Consider stopping and starting a new context for better organization.\n", newCount)
	case newCount == warnAt2:
		fmt.Printf("⚠️  Context now has %d notes and is getting large. Consider chunking your work into smaller contexts.\n", newCount)
	case newCount >= warnAt3:
		// Periodic warnings every 25 notes after threshold 3
		if (newCount-warnAt3)%25 == 0 {
			fmt.Printf("⚠️  Context now has %d notes. This context is quite large - consider stopping and creating a new one.\n", newCount)
		}
	}
}
