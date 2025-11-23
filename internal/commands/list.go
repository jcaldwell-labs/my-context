package commands

import (
	"fmt"
	"strings"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

// filterByProject filters contexts by project name
func filterByProject(contexts []*models.Context, projectFilter string) []*models.Context {
	contextNames := make([]string, 0, len(contexts))
	for _, ctx := range contexts {
		contextNames = append(contextNames, ctx.Name)
	}
	filteredNames := core.FilterContextsByProject(contextNames, projectFilter)

	filtered := make([]*models.Context, 0, len(contexts))
	for _, ctx := range contexts {
		for _, name := range filteredNames {
			if ctx.Name == name {
				filtered = append(filtered, ctx)
				break
			}
		}
	}
	return filtered
}

// filterBySearch filters contexts by search term (case-insensitive)
func filterBySearch(contexts []*models.Context, searchTerm string) []*models.Context {
	var filtered []*models.Context
	searchLower := strings.ToLower(searchTerm)
	for _, ctx := range contexts {
		if strings.Contains(strings.ToLower(ctx.Name), searchLower) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

// filterByTag filters contexts by tag/label
func filterByTag(contexts []*models.Context, tagFilter string) []*models.Context {
	var filtered []*models.Context
	for _, ctx := range contexts {
		ctxWithMeta, _, _, _, err := core.GetContextWithMetadata(ctx.Name)
		if err != nil {
			continue
		}
		for _, tag := range ctxWithMeta.Metadata.Labels {
			if strings.EqualFold(tag, tagFilter) {
				filtered = append(filtered, ctx)
				break
			}
		}
	}
	return filtered
}

// filterByArchiveStatus filters contexts by archive status
func filterByArchiveStatus(contexts []*models.Context, showArchived, activeOnly bool) []*models.Context {
	if showArchived {
		var filtered []*models.Context
		for _, ctx := range contexts {
			if ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return filtered
	}
	if !activeOnly {
		var filtered []*models.Context
		for _, ctx := range contexts {
			if !ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return filtered
	}
	return contexts
}

// filterByActive filters to show only the active context
func filterByActive(contexts []*models.Context, activeContextName string) []*models.Context {
	for _, ctx := range contexts {
		if ctx.Name == activeContextName {
			return []*models.Context{ctx}
		}
	}
	return []*models.Context{}
}

// applyFilters applies all filters to the context list
func applyFilters(contexts []*models.Context, projectFilter, searchTerm, tagFilter string, showArchived, activeOnly bool, activeContextName string) []*models.Context {
	if projectFilter != "" {
		contexts = filterByProject(contexts, projectFilter)
	}
	if searchTerm != "" {
		contexts = filterBySearch(contexts, searchTerm)
	}
	if tagFilter != "" {
		contexts = filterByTag(contexts, tagFilter)
	}
	contexts = filterByArchiveStatus(contexts, showArchived, activeOnly)
	if activeOnly {
		contexts = filterByActive(contexts, activeContextName)
	}
	return contexts
}

// buildContextSummaries builds context summaries for JSON output
func buildContextSummaries(contexts []*models.Context) []*output.ContextSummary {
	summaries := make([]*output.ContextSummary, 0, len(contexts))
	for _, ctx := range contexts {
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
	return summaries
}

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

			// Apply all filters
			contexts := applyFilters(allContexts, projectFilter, searchTerm, tagFilter, showArchived, activeOnly, activeContextName)

			// Apply limit (default 10 unless --all)
			totalCount := len(contexts)
			if !showAll && limitCount > 0 && len(contexts) > limitCount {
				contexts = contexts[:limitCount]
			}

			// Output
			if *jsonOutput {
				summaries := buildContextSummaries(contexts)
				data := output.ListData{Contexts: summaries}
				jsonStr, err := output.FormatJSON("list", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				output.PrintContextHomeHeader(core.GetContextHomeDisplay(), core.GetContextCount())
				fmt.Print(output.FormatContextList(contexts, activeContextName))

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
