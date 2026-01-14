package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

const defaultHistoryLimit = 100

func NewHistoryCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "Show context transition history",
		Long:    `Display the chronological history of all context transitions.`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var transitions []*models.ContextTransition

			// Check if using database backend
			if core.IsUsingDatabase() {
				backend, err := core.GetBackend()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("history", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to get backend: %w", err)
				}
				defer backend.Close()

				// Get transitions from database
				storageTransitions, err := backend.GetTransitions(defaultHistoryLimit)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("history", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to get transitions: %w", err)
				}

				// Convert storage.Transition to models.ContextTransition
				for _, st := range storageTransitions {
					transitions = append(transitions, &models.ContextTransition{
						Timestamp:       st.Timestamp,
						PreviousContext: st.PreviousContext,
						NewContext:      st.NewContext,
						TransitionType:  models.TransitionType(st.TransitionType),
					})
				}
			} else {
				// Get all transitions (file-based backend)
				var err error
				transitions, err = core.GetTransitions()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("history", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}
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
