package commands

import (
	"errors"
	"fmt"
	"time"

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
			// Check if using database backend
			if core.IsUsingDatabase() {
				backend, err := core.GetBackend()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("touch", 2, fmt.Sprintf("failed to get backend: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to get backend: %w", err)
				}
				defer backend.Close()

				// Get active context
				contextName, err := backend.GetActiveContext()
				if err != nil || contextName == "" {
					errMsg := "No active context. Start a context with: my-context start <name>"
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("touch", 1, errMsg)
						fmt.Print(jsonStr)
						return nil
					}
					return errors.New(errMsg)
				}

				// Add touch to database
				timestamp := time.Now().Format(time.RFC3339)
				err = backend.AddTouch(contextName, timestamp)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("touch", 2, fmt.Sprintf("failed to add touch: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to add touch: %w", err)
				}

				// Output
				if *jsonOutput {
					data := map[string]interface{}{"context": contextName, "timestamp": timestamp}
					jsonStr, _ := output.FormatJSON("touch", map[string]interface{}{"data": data})
					fmt.Print(jsonStr)
				} else {
					fmt.Printf("Touch recorded in context: %s\n", contextName)
				}

				return nil
			}

			// File-based backend (existing code)
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
