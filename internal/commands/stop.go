package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewStopCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Aliases: []string{"p"},
		Short:   "Stop the active context",
		Long:    `Stop the currently active context without starting a new one.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := core.StopContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stop", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// No active context
			if context == nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{
						"message": "No active context",
					})
					fmt.Print(jsonStr)
				} else {
					fmt.Println("No active context")
				}
				return nil
			}

			// Output
			if *jsonOutput {
				durationSeconds := int(context.Duration().Seconds())
				data := output.StopData{
					ContextName:     context.Name,
					StartTime:       context.StartTime,
					EndTime:         *context.EndTime,
					DurationSeconds: durationSeconds,
				}
				jsonStr, err := output.FormatJSON("stop", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Stopped context: %s (duration: %s)\n",
					context.Name,
					output.FormatDuration(context.Duration()))
			}

			return nil
		},
	}

	return cmd
}

