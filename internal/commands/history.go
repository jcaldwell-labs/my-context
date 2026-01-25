package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

const defaultHistoryLimit = 50

// filterTransitionsByPeriod filters transitions by time period
func filterTransitionsByPeriod(transitions []*models.ContextTransition, period TimePeriod) []*models.ContextTransition {
	var filtered []*models.ContextTransition
	for _, t := range transitions {
		if t.Timestamp.After(period.Start) && t.Timestamp.Before(period.End) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func NewHistoryCmd(jsonOutput *bool) *cobra.Command {
	var (
		limit int
		today bool
		week  bool
		month bool
		since string
		until string
	)

	cmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"h"},
		Short:   "Show context transition history",
		Long: `Display the chronological history of all context transitions.

Shows start, stop, and switch events for all contexts. A "switch" event
indicates one context stopped and another started at the same moment.

Supports time-based filtering to show transitions for specific periods.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Calculate time period
			period, err := calculatePeriod(today, week, month, since, until)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("history", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

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
				storageTransitions, err := backend.GetTransitions(limit)
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

				// Apply limit for file-based backend
				if limit > 0 && len(transitions) > limit {
					transitions = transitions[:limit]
				}
			}

			// Apply period filter
			transitions = filterTransitionsByPeriod(transitions, period)

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

	cmd.Flags().IntVarP(&limit, "limit", "n", defaultHistoryLimit, "Maximum number of transitions to show (0 for unlimited)")
	cmd.Flags().BoolVarP(&today, "today", "t", false, "Show transitions from today only")
	cmd.Flags().BoolVarP(&week, "week", "w", false, "Show transitions from this week (Mon-Sun)")
	cmd.Flags().BoolVarP(&month, "month", "m", false, "Show transitions from this month")
	cmd.Flags().StringVar(&since, "since", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "End date (YYYY-MM-DD)")

	return cmd
}
