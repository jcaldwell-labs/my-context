package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all contexts",
		Long:    `List all contexts (active and stopped) with their status and timestamps.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all contexts
			contexts, err := core.ListContexts()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("list", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Get active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("list", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			activeContextName := state.GetActiveContextName()

			// Output
			if *jsonOutput {
				// Build context summaries with counts
				var summaries []*output.ContextSummary
				for _, ctx := range contexts {
					// Get counts from log files
					notesLines, _ := core.ReadLog(core.GetNotesLogPath(ctx.Name))
					filesLines, _ := core.ReadLog(core.GetFilesLogPath(ctx.Name))
					touchesLines, _ := core.ReadLog(core.GetTouchLogPath(ctx.Name))

					summary := &output.ContextSummary{
						Name:            ctx.Name,
						StartTime:       ctx.StartTime,
						EndTime:         ctx.EndTime,
						Status:          ctx.Status,
						DurationSeconds: int(ctx.Duration().Seconds()),
						NoteCount:       len(notesLines),
						FileCount:       len(filesLines),
						TouchCount:      len(touchesLines),
					}
					summaries = append(summaries, summary)
				}

				data := output.ListData{
					Contexts: summaries,
				}
				jsonStr, err := output.FormatJSON("list", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Print(output.FormatContextList(contexts, activeContextName))
			}

			return nil
		},
	}

	return cmd
}

