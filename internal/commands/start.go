package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewStartCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start <name>",
		Aliases: []string{"s"},
		Short:   "Start a new context",
		Long:    `Start a new context with the given name. If a context is already active, it will be automatically stopped.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			if contextName == "" {
				return fmt.Errorf("context name cannot be empty")
			}

			// Create the context
			context, previousContext, err := core.CreateContext(contextName)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("start", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.StartData{
					ContextName:  context.Name,
					OriginalName: contextName,
					WasDuplicate: context.Name != core.SanitizeContextName(contextName),
				}
				if previousContext != "" {
					data.PreviousContext = &previousContext
				}
				jsonStr, err := output.FormatJSON("start", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				// Check if name was modified
				if context.Name != core.SanitizeContextName(contextName) {
					fmt.Printf("Context \"%s\" already exists.\n", core.SanitizeContextName(contextName))
				}

				// Show previous context stop message
				if previousContext != "" {
					fmt.Printf("Stopped context: %s\n", previousContext)
				}

				fmt.Printf("Started context: %s\n", context.Name)
			}

			return nil
		},
	}

	return cmd
}
