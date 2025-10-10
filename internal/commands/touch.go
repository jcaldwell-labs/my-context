package commands

import (
	"errors"
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewTouchCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "touch",
		Aliases: []string{"t"},
		Short:   "Record a timestamp in the active context",
		Long:    `Record a timestamp in the currently active context to indicate activity without detailed notes.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("touch", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			if !state.HasActiveContext() {
				errMsg := "No active context. Start a context with: my-context start <name>"
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("touch", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			}

			// Add the touch
			touch, err := core.AddTouch()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("touch", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.TouchData{
					ContextName:    state.GetActiveContextName(),
					TouchTimestamp: touch.Timestamp,
				}
				jsonStr, err := output.FormatJSON("touch", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Touch recorded in context: %s\n", state.GetActiveContextName())
			}

			return nil
		},
	}

	return cmd
}

