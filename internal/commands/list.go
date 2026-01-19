package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
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
			fmt.Fprintf(os.Stderr, "warning: skipping context %q in tag filter: %v\n", ctx.Name, err)
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

// buildContextSummaries builds context summaries for JSON output (file-based mode)
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

// filterDBContextsByProject filters database contexts by project name (single-pass)
func filterDBContextsByProject(dbContexts []*pkgmodels.ContextWithMetadata, projectFilter string) []*pkgmodels.ContextWithMetadata {
	if projectFilter == "" {
		return dbContexts
	}
	projectFilter = strings.TrimSpace(projectFilter)
	var filtered []*pkgmodels.ContextWithMetadata
	for _, ctx := range dbContexts {
		contextProject := core.ExtractProjectName(ctx.Name)
		if strings.EqualFold(contextProject, projectFilter) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

// filterDBContextsBySearch filters database contexts by search term (case-insensitive)
func filterDBContextsBySearch(dbContexts []*pkgmodels.ContextWithMetadata, searchTerm string) []*pkgmodels.ContextWithMetadata {
	var filtered []*pkgmodels.ContextWithMetadata
	searchLower := strings.ToLower(searchTerm)
	for _, ctx := range dbContexts {
		if strings.Contains(strings.ToLower(ctx.Name), searchLower) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

// filterDBContextsByTag filters database contexts by tag/label using database metadata
func filterDBContextsByTag(dbContexts []*pkgmodels.ContextWithMetadata, tagFilter string) []*pkgmodels.ContextWithMetadata {
	var filtered []*pkgmodels.ContextWithMetadata
	for _, ctx := range dbContexts {
		for _, tag := range ctx.Metadata.Labels {
			if strings.EqualFold(tag, tagFilter) {
				filtered = append(filtered, ctx)
				break
			}
		}
	}
	return filtered
}

// filterDBContextsByArchiveStatus filters database contexts by archive status
func filterDBContextsByArchiveStatus(dbContexts []*pkgmodels.ContextWithMetadata, showArchived, activeOnly bool) []*pkgmodels.ContextWithMetadata {
	if showArchived {
		var filtered []*pkgmodels.ContextWithMetadata
		for _, ctx := range dbContexts {
			if ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return filtered
	}
	if !activeOnly {
		var filtered []*pkgmodels.ContextWithMetadata
		for _, ctx := range dbContexts {
			if !ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return filtered
	}
	return dbContexts
}

// filterDBContextsByActive filters to show only the active context
func filterDBContextsByActive(dbContexts []*pkgmodels.ContextWithMetadata, activeContextName string) []*pkgmodels.ContextWithMetadata {
	for _, ctx := range dbContexts {
		if ctx.Name == activeContextName {
			return []*pkgmodels.ContextWithMetadata{ctx}
		}
	}
	return []*pkgmodels.ContextWithMetadata{}
}

// applyDBFilters applies all filters to database context list
func applyDBFilters(dbContexts []*pkgmodels.ContextWithMetadata, projectFilter, searchTerm, tagFilter string, showArchived, activeOnly bool, activeContextName string) []*pkgmodels.ContextWithMetadata {
	if projectFilter != "" {
		dbContexts = filterDBContextsByProject(dbContexts, projectFilter)
	}
	if searchTerm != "" {
		dbContexts = filterDBContextsBySearch(dbContexts, searchTerm)
	}
	if tagFilter != "" {
		dbContexts = filterDBContextsByTag(dbContexts, tagFilter)
	}
	dbContexts = filterDBContextsByArchiveStatus(dbContexts, showArchived, activeOnly)
	if activeOnly {
		dbContexts = filterDBContextsByActive(dbContexts, activeContextName)
	}
	return dbContexts
}

// buildContextSummariesFromDB builds context summaries from database models
func buildContextSummariesFromDB(dbContexts []*pkgmodels.ContextWithMetadata) []*output.ContextSummary {
	summaries := make([]*output.ContextSummary, 0, len(dbContexts))
	for _, ctx := range dbContexts {
		summary := &output.ContextSummary{
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

	// Apply all filters directly on database models (preserves metadata for tag filtering)
	filteredDBContexts := applyDBFilters(dbContexts, projectFilter, searchTerm, tagFilter, showArchived, activeOnly, activeContextName)

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

			// File-based backend (existing code)
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

	// Mark mutually exclusive flags
	cmd.MarkFlagsMutuallyExclusive("archived", "active-only")

	return cmd
}
