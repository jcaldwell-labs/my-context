package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/filters"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/spf13/cobra"
)

// buildContextSummaries builds context summaries for JSON output (file-based mode)
func buildContextSummaries(contexts []*models.Context) []*output.ContextSummary {
	summaries := make([]*output.ContextSummary, 0, len(contexts))
	for i, ctx := range contexts {
		notesLines, _ := core.ReadLog(core.GetNotesLogPath(ctx.Name))
		filesLines, _ := core.ReadLog(core.GetFilesLogPath(ctx.Name))
		touchesLines, _ := core.ReadLog(core.GetTouchLogPath(ctx.Name))

		summary := &output.ContextSummary{
			Index:           i + 1,
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

// convertDBContextToInternal converts a single database model to internal model
func convertDBContextToInternal(dbCtx *pkgmodels.ContextWithMetadata) *models.Context {
	return &models.Context{
		Name:       dbCtx.Name,
		StartTime:  dbCtx.StartTime,
		EndTime:    dbCtx.EndTime,
		Status:     dbCtx.Status,
		IsArchived: dbCtx.IsArchived,
	}
}

// buildContextSummariesFromDB builds context summaries from database models
func buildContextSummariesFromDB(dbContexts []*pkgmodels.ContextWithMetadata) []*output.ContextSummary {
	summaries := make([]*output.ContextSummary, 0, len(dbContexts))
	for i, ctx := range dbContexts {
		summary := &output.ContextSummary{
			Index:           i + 1,
			Name:            ctx.Name,
			StartTime:       ctx.StartTime,
			EndTime:         ctx.EndTime,
			Status:          ctx.Status,
			DurationSeconds: int(ctx.Duration().Seconds()),
			TouchCount:      ctx.TouchCount,
		}
		summaries = append(summaries, summary)
	}
	return summaries
}

// listFromDatabase handles listing contexts from database backend
func listFromDatabase(jsonOutput *bool, projectFilter, searchTerm, tagFilter string, limitCount int, showAll, showArchived, activeOnly bool) error {
	backend, err := core.GetBackend()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("list", 2, fmt.Sprintf("failed to get backend: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to get backend: %w", err)
	}
	defer backend.Close()

	// Get all contexts from database
	dbContexts, err := backend.ListContexts()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("list", 2, err.Error())
			fmt.Print(jsonStr)
			return nil
		}
		return err
	}

	// Get active context from database
	activeContextName, _ := backend.GetActiveContext()

	// Apply filters using filter package
	filter := filters.NewDBContextFilter(dbContexts)
	if projectFilter != "" {
		filter = filter.ByProject(projectFilter).(*filters.DBContextFilter)
	}
	if searchTerm != "" {
		filter = filter.BySearch(searchTerm).(*filters.DBContextFilter)
	}
	if tagFilter != "" {
		filter = filter.ByTag(tagFilter).(*filters.DBContextFilter)
	}
	filter = filter.ByArchiveStatus(showArchived, activeOnly).(*filters.DBContextFilter)
	if activeOnly {
		filter = filter.ByActive(activeContextName).(*filters.DBContextFilter)
	}
	filteredDBContexts := filter.GetDBContexts()

	// Apply limit (default 10 unless --all)
	totalCount := len(filteredDBContexts)
	if !showAll && limitCount > 0 && len(filteredDBContexts) > limitCount {
		filteredDBContexts = filteredDBContexts[:limitCount]
	}

	// Convert filtered contexts to internal models for display
	contexts := make([]*models.Context, 0, len(filteredDBContexts))
	for _, dbCtx := range filteredDBContexts {
		contexts = append(contexts, convertDBContextToInternal(dbCtx))
	}

	// Output
	if *jsonOutput {
		summaries := buildContextSummariesFromDB(filteredDBContexts)
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
			// Check if using database backend
			if core.IsUsingDatabase() {
				return listFromDatabase(jsonOutput, projectFilter, searchTerm, tagFilter, limitCount, showAll, showArchived, activeOnly)
			}

			// File-based backend
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

			// Apply filters using filter package
			filter := filters.NewFileContextFilter(allContexts)
			if projectFilter != "" {
				filter = filter.ByProject(projectFilter).(*filters.FileContextFilter)
			}
			if searchTerm != "" {
				filter = filter.BySearch(searchTerm).(*filters.FileContextFilter)
			}
			if tagFilter != "" {
				filter = filter.ByTag(tagFilter).(*filters.FileContextFilter)
			}
			filter = filter.ByArchiveStatus(showArchived, activeOnly).(*filters.FileContextFilter)
			if activeOnly {
				filter = filter.ByActive(activeContextName).(*filters.FileContextFilter)
			}
			contexts := filter.GetFileContexts()

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

	// Mark mutually exclusive flags
	cmd.MarkFlagsMutuallyExclusive("archived", "active-only")

	return cmd
}
