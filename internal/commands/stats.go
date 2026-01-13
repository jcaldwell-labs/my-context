package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/spf13/cobra"
)

// StatsResult holds aggregated statistics
type StatsResult struct {
	Period         TimePeriod
	TotalContexts  int
	TotalSeconds   int64
	AverageSeconds int64
	ByProject      []ProjectStats
	LongestSession *SessionInfo
	ActiveContext  *SessionInfo
}

// TimePeriod represents the time range for stats
type TimePeriod struct {
	Start time.Time
	End   time.Time
	Label string
}

// ProjectStats holds per-project aggregation
type ProjectStats struct {
	Project string
	Seconds int64
	Count   int
	Percent float64
}

// SessionInfo holds info about a specific session
type SessionInfo struct {
	Name    string
	Seconds int64
}

// calculatePeriod determines the time period based on flags
func calculatePeriod(today, week, month bool, since, until string) (TimePeriod, error) {
	now := time.Now()
	var period TimePeriod

	switch {
	case today:
		// Start of today in local time
		y, m, d := now.Date()
		period.Start = time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "Today"

	case week:
		// Start of week (Monday)
		y, m, d := now.Date()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		daysBack := weekday - 1 // Days since Monday
		period.Start = time.Date(y, m, d-daysBack, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "This Week"

	case month:
		// Start of month
		y, m, _ := now.Date()
		period.Start = time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
		period.End = now
		period.Label = "This Month"

	case since != "" || until != "":
		// Custom date range
		if since != "" {
			t, err := time.Parse("2006-01-02", since)
			if err != nil {
				return period, fmt.Errorf("invalid --since date format (use YYYY-MM-DD): %w", err)
			}
			period.Start = t
		} else {
			// Default to beginning of time (1970)
			period.Start = time.Unix(0, 0)
		}

		if until != "" {
			t, err := time.Parse("2006-01-02", until)
			if err != nil {
				return period, fmt.Errorf("invalid --until date format (use YYYY-MM-DD): %w", err)
			}
			// End of the specified day
			period.End = t.Add(24*time.Hour - time.Second)
		} else {
			period.End = now
		}

		switch {
		case since != "" && until != "":
			period.Label = fmt.Sprintf("%s to %s", since, until)
		case since != "":
			period.Label = fmt.Sprintf("Since %s", since)
		default:
			period.Label = fmt.Sprintf("Until %s", until)
		}

	default:
		// All time
		period.Start = time.Unix(0, 0)
		period.End = now
		period.Label = "All Time"
	}

	return period, nil
}

// filterByPeriod filters contexts by time period
func filterContextsByPeriod(contexts []*models.Context, period TimePeriod) []*models.Context {
	var filtered []*models.Context
	for _, ctx := range contexts {
		if ctx.StartTime.After(period.Start) && ctx.StartTime.Before(period.End) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

// filterDBContextsByPeriod filters database contexts by time period
func filterDBContextsByPeriod(contexts []*pkgmodels.ContextWithMetadata, period TimePeriod) []*pkgmodels.ContextWithMetadata {
	var filtered []*pkgmodels.ContextWithMetadata
	for _, ctx := range contexts {
		if ctx.StartTime.After(period.Start) && ctx.StartTime.Before(period.End) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

// statsContextItem provides a common interface for aggregating stats
type statsContextItem struct {
	Name     string
	Duration time.Duration
}

// aggregateStatsGeneric calculates statistics from context items
func aggregateStatsGeneric(items []statsContextItem, period TimePeriod, activeContextName string) StatsResult {
	result := StatsResult{
		Period:        period,
		TotalContexts: len(items),
	}

	if len(items) == 0 {
		return result
	}

	projectTotals := make(map[string]int64)
	projectCounts := make(map[string]int)
	var longestName string
	var longestSeconds int64

	for _, item := range items {
		seconds := int64(item.Duration.Seconds())
		result.TotalSeconds += seconds

		// Extract project
		project := core.ExtractProjectName(item.Name)
		projectTotals[project] += seconds
		projectCounts[project]++

		// Track longest
		if seconds > longestSeconds {
			longestSeconds = seconds
			longestName = item.Name
		}

		// Track active context
		if item.Name == activeContextName {
			result.ActiveContext = &SessionInfo{
				Name:    item.Name,
				Seconds: seconds,
			}
		}
	}

	// Calculate average
	if result.TotalContexts > 0 {
		result.AverageSeconds = result.TotalSeconds / int64(result.TotalContexts)
	}

	// Set longest session
	if longestName != "" {
		result.LongestSession = &SessionInfo{
			Name:    longestName,
			Seconds: longestSeconds,
		}
	}

	// Build project stats (sorted by time descending)
	for project, seconds := range projectTotals {
		percent := 0.0
		if result.TotalSeconds > 0 {
			percent = float64(seconds) / float64(result.TotalSeconds) * 100
		}
		result.ByProject = append(result.ByProject, ProjectStats{
			Project: project,
			Seconds: seconds,
			Count:   projectCounts[project],
			Percent: percent,
		})
	}

	sort.Slice(result.ByProject, func(i, j int) bool {
		return result.ByProject[i].Seconds > result.ByProject[j].Seconds
	})

	return result
}

// aggregateStats calculates statistics from file-based contexts
func aggregateStats(contexts []*models.Context, period TimePeriod, activeContextName string) StatsResult {
	items := make([]statsContextItem, len(contexts))
	for i, ctx := range contexts {
		items[i] = statsContextItem{Name: ctx.Name, Duration: ctx.Duration()}
	}
	return aggregateStatsGeneric(items, period, activeContextName)
}

// aggregateDBStats calculates statistics from database contexts
func aggregateDBStats(contexts []*pkgmodels.ContextWithMetadata, period TimePeriod, activeContextName string) StatsResult {
	items := make([]statsContextItem, len(contexts))
	for i, ctx := range contexts {
		items[i] = statsContextItem{Name: ctx.Name, Duration: ctx.Duration()}
	}
	return aggregateStatsGeneric(items, period, activeContextName)
}

// formatStatsHuman formats stats for human-readable output
func formatStatsHuman(stats *StatsResult, projectFilter string) string {
	var sb strings.Builder

	sb.WriteString("Time Statistics\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Period
	sb.WriteString(fmt.Sprintf("Period: %s\n", stats.Period.Label))
	if projectFilter != "" {
		sb.WriteString(fmt.Sprintf("Project: %s\n", projectFilter))
	}
	sb.WriteString("\n")

	// Summary
	sb.WriteString(fmt.Sprintf("Total Contexts: %d\n", stats.TotalContexts))
	sb.WriteString(fmt.Sprintf("Total Time: %s\n", output.FormatDuration(time.Duration(stats.TotalSeconds)*time.Second)))
	sb.WriteString(fmt.Sprintf("Average Session: %s\n", output.FormatDuration(time.Duration(stats.AverageSeconds)*time.Second)))

	// Active context
	if stats.ActiveContext != nil {
		sb.WriteString(fmt.Sprintf("\nActive Context: %s (%s)\n",
			stats.ActiveContext.Name,
			output.FormatDuration(time.Duration(stats.ActiveContext.Seconds)*time.Second)))
	}

	// By project (only if no project filter and we have multiple projects)
	if projectFilter == "" && len(stats.ByProject) > 1 {
		sb.WriteString("\nBy Project:\n")

		// Show top 5 projects
		showCount := len(stats.ByProject)
		if showCount > 5 {
			showCount = 5
		}

		for i := 0; i < showCount; i++ {
			p := stats.ByProject[i]
			// Progress bar (20 chars max)
			barLen := int(p.Percent / 5) // 100% = 20 chars
			if barLen > 20 {
				barLen = 20
			}
			bar := strings.Repeat("█", barLen) + strings.Repeat("░", 20-barLen)

			sb.WriteString(fmt.Sprintf("  %-20s %8s (%4.1f%%)  %s\n",
				truncateString(p.Project, 20),
				output.FormatDuration(time.Duration(p.Seconds)*time.Second),
				p.Percent,
				bar))
		}

		if len(stats.ByProject) > 5 {
			sb.WriteString(fmt.Sprintf("  ... and %d more projects\n", len(stats.ByProject)-5))
		}
	}

	// Longest session
	if stats.LongestSession != nil && stats.TotalContexts > 1 {
		sb.WriteString(fmt.Sprintf("\nLongest Session: %s (%s)\n",
			truncateString(stats.LongestSession.Name, 40),
			output.FormatDuration(time.Duration(stats.LongestSession.Seconds)*time.Second)))
	}

	return sb.String()
}

// truncateString truncates a string to maxLen with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// statsFromDatabase handles stats from database backend
func statsFromDatabase(jsonOutput *bool, projectFilter string, period TimePeriod) error {
	backend, err := core.GetBackend()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("stats", 2, fmt.Sprintf("failed to get backend: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to get backend: %w", err)
	}
	defer backend.Close()

	// Get all contexts
	dbContexts, err := backend.ListContexts()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("stats", 2, err.Error())
			fmt.Print(jsonStr)
			return nil
		}
		return err
	}

	// Get active context
	activeContextName, _ := backend.GetActiveContext()

	// Filter by period
	dbContexts = filterDBContextsByPeriod(dbContexts, period)

	// Filter by project if specified
	if projectFilter != "" {
		dbContexts = filterDBContextsByProject(dbContexts, projectFilter)
	}

	// Aggregate stats
	stats := aggregateDBStats(dbContexts, period, activeContextName)

	// Output
	if *jsonOutput {
		jsonStr, err := output.FormatStatsJSON(stats.Period, stats.ByProject, stats.LongestSession, stats.ActiveContext, stats.TotalContexts, stats.TotalSeconds, stats.AverageSeconds)
		if err != nil {
			return err
		}
		fmt.Print(jsonStr)
	} else {
		fmt.Print(formatStatsHuman(&stats, projectFilter))
	}

	return nil
}

// NewStatsCmd creates the stats command
func NewStatsCmd(jsonOutput *bool) *cobra.Command {
	var (
		projectFilter string
		today         bool
		week          bool
		month         bool
		since         string
		until         string
	)

	cmd := &cobra.Command{
		Use:     "stats",
		Aliases: []string{"st"},
		Short:   "Show time statistics across contexts",
		Long: `Aggregate and report time spent across contexts.

Shows total time, average session duration, and breakdown by project.
Supports filtering by time period and project.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Calculate time period
			period, err := calculatePeriod(today, week, month, since, until)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stats", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Check if using database backend
			if core.IsUsingDatabase() {
				return statsFromDatabase(jsonOutput, projectFilter, period)
			}

			// File-based backend
			allContexts, err := core.ListContexts()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stats", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Get active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stats", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}
			activeContextName := state.GetActiveContextName()

			// Filter by period
			contexts := filterContextsByPeriod(allContexts, period)

			// Filter by project if specified
			if projectFilter != "" {
				contexts = filterByProject(contexts, projectFilter)
			}

			// Exclude archived
			contexts = filterByArchiveStatus(contexts, false, false)

			// Aggregate stats
			stats := aggregateStats(contexts, period, activeContextName)

			// Output
			if *jsonOutput {
				jsonStr, err := output.FormatStatsJSON(stats.Period, stats.ByProject, stats.LongestSession, stats.ActiveContext, stats.TotalContexts, stats.TotalSeconds, stats.AverageSeconds)
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Print(formatStatsHuman(&stats, projectFilter))
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&projectFilter, "project", "p", "", "Filter by project name")
	cmd.Flags().BoolVarP(&today, "today", "t", false, "Show stats for today only")
	cmd.Flags().BoolVarP(&week, "week", "w", false, "Show stats for this week (Mon-Sun)")
	cmd.Flags().BoolVarP(&month, "month", "m", false, "Show stats for this month")
	cmd.Flags().StringVar(&since, "since", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "End date (YYYY-MM-DD)")

	return cmd
}
