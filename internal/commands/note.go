package commands

import (
	"fmt"

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
				return fmt.Errorf(errMsg)
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

