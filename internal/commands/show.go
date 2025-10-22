package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewShowCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [context-name]",
		Aliases: []string{"w"},
		Short:   "Show context details",
		Long:    `Display details about the currently active context including notes, files, and touch events.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var contextName string

			// If context name provided as argument, use it
			if len(args) > 0 {
				contextName = args[0]
			} else {
				// No argument - show active context (backward compatible)
				state, err := core.GetActiveContext()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("show", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				if !state.HasActiveContext() {
					errMsg := "No active context"
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("show", 1, errMsg)
						fmt.Print(jsonStr)
					} else {
						fmt.Println(errMsg)
						fmt.Println("Start one with: my-context start <name>")
					}
					return nil
				}

				contextName = state.GetActiveContextName()
			}

			// Get context details (works for stopped or active)
			context, notes, files, touches, err := core.GetContextWithMetadata(contextName)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("show", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.ContextData{
					Context: context,
					Notes:   notes,
					Files:   files,
					Touches: touches,
				}
				jsonStr, err := output.FormatJSON("show", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				// Print context home header
				output.PrintContextHomeHeader(core.GetContextHomeDisplay(), core.GetContextCount())
				fmt.Print(output.FormatContext(context, notes, files, touches))
			}

			return nil
		},
	}

	return cmd
}

