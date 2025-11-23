package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewHistoryCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "Show context transition history",
		Long:    `Display the chronological history of all context transitions.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all transitions
			transitions, err := core.GetTransitions()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("history", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.HistoryData{
					Transitions: transitions,
				}
				jsonStr, err := output.FormatJSON("history", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				// Print context home header
				output.PrintContextHomeHeader(core.GetContextHomeDisplay(), core.GetContextCount())
				fmt.Print(output.FormatTransitionHistory(transitions))
			}

			return nil
		},
	}

	return cmd
}
