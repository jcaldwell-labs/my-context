package commands

import (
	"fmt"
	"strings"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(jsonOutput *bool) *cobra.Command {
	var (
		projectFilter string
		searchTerm    string
		tagFilter     string
		limitCount    int
		showAll       bool
		showArchived  bool
		activeOnly    bool
	)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all contexts",
		Long: `List all contexts (active and stopped) with their status and timestamps.

Supports filtering by project, search term, and archive status.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all contexts
			allContexts, err := core.ListContexts()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("list", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Apply filters
			contexts := allContexts

			// Filter by project
			if projectFilter != "" {
				var contextNames []string
				for _, ctx := range contexts {
					contextNames = append(contextNames, ctx.Name)
				}
				filteredNames := core.FilterContextsByProject(contextNames, projectFilter)

				// Keep only matching contexts
				var filtered []*models.Context
				for _, ctx := range contexts {
					for _, name := range filteredNames {
						if ctx.Name == name {
							filtered = append(filtered, ctx)
							break
						}
					}
				}
				contexts = filtered
			}

			// Filter by search term (case-insensitive)
			if searchTerm != "" {
				var filtered []*models.Context
				searchLower := strings.ToLower(searchTerm)
				for _, ctx := range contexts {
					if strings.Contains(strings.ToLower(ctx.Name), searchLower) {
						filtered = append(filtered, ctx)
					}
				}
				contexts = filtered
			}

			// Filter by tag
			if tagFilter != "" {
				var filtered []*models.Context
				for _, ctx := range contexts {
					// Load context with metadata to check tags
					ctxWithMeta, _, _, _, err := core.GetContextWithMetadata(ctx.Name)
					if err != nil {
						continue // Skip if can't load metadata
					}
					// Check if context has the tag
					for _, tag := range ctxWithMeta.Metadata.Labels {
						if strings.EqualFold(tag, tagFilter) {
							filtered = append(filtered, ctx)
							break
						}
					}
				}
				contexts = filtered
			}

			// Filter by archive status
			if showArchived {
				// Show only archived
				var filtered []*models.Context
				for _, ctx := range contexts {
					if ctx.IsArchived {
						filtered = append(filtered, ctx)
					}
				}
				contexts = filtered
			} else if !activeOnly {
				// Default: hide archived contexts
				var filtered []*models.Context
				for _, ctx := range contexts {
					if !ctx.IsArchived {
						filtered = append(filtered, ctx)
					}
				}
				contexts = filtered
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

			// Filter by active only
			if activeOnly {
				var filtered []*models.Context
				for _, ctx := range contexts {
					if ctx.Name == activeContextName {
						filtered = append(filtered, ctx)
						break
					}
				}
				contexts = filtered
			}

			// Apply limit (default 10 unless --all)
			totalCount := len(contexts)
			if !showAll && limitCount > 0 && len(contexts) > limitCount {
				contexts = contexts[:limitCount]
			}

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
				// Print context home header
				output.PrintContextHomeHeader(core.GetContextHomeDisplay(), core.GetContextCount())
				fmt.Print(output.FormatContextList(contexts, activeContextName))

				// Show truncation message if limited
				if !showAll && limitCount > 0 && totalCount > len(contexts) {
					fmt.Printf("\nShowing %d of %d contexts. Use --all to see all.\n", len(contexts), totalCount)
				}
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&projectFilter, "project", "", "Filter by project name")
	cmd.Flags().StringVar(&searchTerm, "search", "", "Search contexts by name (case-insensitive)")
	cmd.Flags().StringVar(&tagFilter, "tag", "", "Filter by tag/label")
	cmd.Flags().IntVar(&limitCount, "limit", 10, "Maximum number of contexts to show")
	cmd.Flags().BoolVar(&showAll, "all", false, "Show all contexts (no limit)")
	cmd.Flags().BoolVar(&showArchived, "archived", false, "Show only archived contexts")
	cmd.Flags().BoolVar(&activeOnly, "active-only", false, "Show only the active context")

	return cmd
}
