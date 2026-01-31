package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/filters"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/jefferycaldwell/my-context-copilot/pkg/storage"
	"github.com/spf13/cobra"
)

// filterStaleContexts filters to show only stale active contexts
func filterStaleContexts(contexts []*models.Context) []*models.Context {
	thresholds := core.GetStaleThresholds()
	var filtered []*models.Context

	for _, ctx := range contexts {
		// Only check active contexts
		if ctx.Status != "active" {
			continue
		}

		// Get last activity time - need to fetch full context data
		_, notes, files, touches, err := core.GetContextWithMetadata(ctx.Name)
		if err != nil {
			continue
		}

		lastActivity := core.GetLastActivityTime(ctx, notes, files, touches)
		if core.IsStale(lastActivity, thresholds) {
			filtered = append(filtered, ctx)
		}
	}

	return filtered
}

// buildContextSummaries builds context summaries for JSON output (file-based mode)
func buildContextSummaries(contexts []*models.Context) []*output.ContextSummary {
	thresholds := core.GetStaleThresholds()
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

		// Add stale information for active contexts
		if ctx.Status == "active" {
			_, notes, files, touches, err := core.GetContextWithMetadata(ctx.Name)
			if err == nil {
				lastActivity := core.GetLastActivityTime(ctx, notes, files, touches)
				staleLevel := core.GetStaleLevel(lastActivity, thresholds)
				summary.IsStale = core.IsStale(lastActivity, thresholds)
				summary.StaleLevel = core.GetStaleLevelString(staleLevel)
			}
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

// filterDBContextsByStale filters to show only stale active contexts
func filterDBContextsByStale(dbContexts []*pkgmodels.ContextWithMetadata, backend storage.Backend) []*pkgmodels.ContextWithMetadata {
	thresholds := core.GetStaleThresholds()
	var filtered []*pkgmodels.ContextWithMetadata

	for _, ctx := range dbContexts {
		// Only check active contexts
		if ctx.Status != "active" {
			continue
		}

		// Get last activity from database
		notes, _ := backend.GetNotes(ctx.Name)
		files, _ := backend.GetFiles(ctx.Name)

		// Convert to internal models
		internalCtx := &models.Context{
			Name:      ctx.Name,
			StartTime: ctx.StartTime,
			EndTime:   ctx.EndTime,
			Status:    ctx.Status,
		}

		var internalNotes []*models.Note
		for _, n := range notes {
			internalNotes = append(internalNotes, &models.Note{
				Timestamp:   n.Timestamp,
				TextContent: n.Content,
			})
		}

		var internalFiles []*models.FileAssociation
		for _, f := range files {
			internalFiles = append(internalFiles, &models.FileAssociation{
				Timestamp: f.Timestamp,
				FilePath:  f.Path,
			})
		}

		// Compute lastActivity as max across notes, files, and optional touch timestamp
		// LastTouchAt is treated as another candidate, not an override
		var touches []*models.TouchEvent
		if ctx.LastTouchAt != nil {
			touches = append(touches, &models.TouchEvent{Timestamp: *ctx.LastTouchAt})
		}
		lastActivity := core.GetLastActivityTime(internalCtx, internalNotes, internalFiles, touches)

		if core.IsStale(lastActivity, thresholds) {
			filtered = append(filtered, ctx)
		}
	}

	return filtered
}

// buildContextSummariesFromDB builds context summaries from database models
func buildContextSummariesFromDB(dbContexts []*pkgmodels.ContextWithMetadata, backend storage.Backend) []*output.ContextSummary {
	thresholds := core.GetStaleThresholds()
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

		// Add stale information for active contexts
		if ctx.Status == "active" {
			notes, _ := backend.GetNotes(ctx.Name)
			files, _ := backend.GetFiles(ctx.Name)

			// Convert to internal models
			internalCtx := &models.Context{
				Name:      ctx.Name,
				StartTime: ctx.StartTime,
				EndTime:   ctx.EndTime,
				Status:    ctx.Status,
			}

			var internalNotes []*models.Note
			for _, n := range notes {
				internalNotes = append(internalNotes, &models.Note{
					Timestamp:   n.Timestamp,
					TextContent: n.Content,
				})
			}

			var internalFiles []*models.FileAssociation
			for _, f := range files {
				internalFiles = append(internalFiles, &models.FileAssociation{
					Timestamp: f.Timestamp,
					FilePath:  f.Path,
				})
			}

			// Compute lastActivity as max across notes, files, and optional touch timestamp
			var touches []*models.TouchEvent
			if ctx.LastTouchAt != nil {
				touches = append(touches, &models.TouchEvent{Timestamp: *ctx.LastTouchAt})
			}
			lastActivity := core.GetLastActivityTime(internalCtx, internalNotes, internalFiles, touches)

			staleLevel := core.GetStaleLevel(lastActivity, thresholds)
			summary.IsStale = core.IsStale(lastActivity, thresholds)
			summary.StaleLevel = core.GetStaleLevelString(staleLevel)
		}

		summaries = append(summaries, summary)
	}
	return summaries
}

// listFromDatabase handles listing contexts from database backend
func listFromDatabase(jsonOutput *bool, projectFilter, searchTerm, tagFilter string, limitCount int, showAll, showArchived, activeOnly, showStale bool) error {
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

	// Apply stale filter if requested
	if showStale {
		filteredDBContexts = filterDBContextsByStale(filteredDBContexts, backend)
	}

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
		summaries := buildContextSummariesFromDB(filteredDBContexts, backend)
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
		showStale     bool
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
				return listFromDatabase(jsonOutput, projectFilter, searchTerm, tagFilter, limitCount, showAll, showArchived, activeOnly, showStale)
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

			// Apply stale filter if requested
			if showStale {
				contexts = filterStaleContexts(contexts)
			}

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
	cmd.Flags().BoolVar(&showStale, "stale", false, "Show only stale contexts (active with no recent activity)")

	// Mark mutually exclusive flags
	cmd.MarkFlagsMutuallyExclusive("archived", "active-only", "stale")

	return cmd
}
