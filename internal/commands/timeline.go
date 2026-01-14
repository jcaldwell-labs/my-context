package commands

import (
	"fmt"
	"sort"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/spf13/cobra"
)

// buildTimelineFromFileContexts constructs timeline from file-based contexts
func buildTimelineFromFileContexts(contexts []*models.Context, period TimePeriod, activeContextName string) output.TimelineData {
	var entries []output.TimelineEntry

	// Sort contexts by start time
	sortedContexts := make([]*models.Context, len(contexts))
	copy(sortedContexts, contexts)
	sort.Slice(sortedContexts, func(i, j int) bool {
		return sortedContexts[i].StartTime.Before(sortedContexts[j].StartTime)
	})

	// Convert to timeline entries
	for _, ctx := range sortedContexts {
		isActive := ctx.Name == activeContextName && ctx.Status == "active"
		entry := output.TimelineEntry{
			ContextName: ctx.Name,
			StartTime:   ctx.StartTime,
			EndTime:     ctx.EndTime,
			Duration:    ctx.Duration(),
			IsGap:       false,
			IsActive:    isActive,
		}
		entries = append(entries, entry)
	}

	return buildTimelineData(entries, period)
}

// buildTimelineFromDBContexts constructs timeline from database contexts
func buildTimelineFromDBContexts(contexts []*pkgmodels.ContextWithMetadata, period TimePeriod, activeContextName string) output.TimelineData {
	var entries []output.TimelineEntry

	// Sort contexts by start time
	sortedContexts := make([]*pkgmodels.ContextWithMetadata, len(contexts))
	copy(sortedContexts, contexts)
	sort.Slice(sortedContexts, func(i, j int) bool {
		return sortedContexts[i].StartTime.Before(sortedContexts[j].StartTime)
	})

	// Convert to timeline entries
	for _, ctx := range sortedContexts {
		isActive := ctx.Name == activeContextName && ctx.Status == "active"
		entry := output.TimelineEntry{
			ContextName: ctx.Name,
			StartTime:   ctx.StartTime,
			EndTime:     ctx.EndTime,
			Duration:    ctx.Duration(),
			IsGap:       false,
			IsActive:    isActive,
		}
		entries = append(entries, entry)
	}

	return buildTimelineData(entries, period)
}

// buildTimelineData constructs timeline with gap detection
func buildTimelineData(entries []output.TimelineEntry, period TimePeriod) output.TimelineData {
	timeline := output.TimelineData{
		Period:  period,
		Entries: make([]output.TimelineEntry, 0),
	}

	if len(entries) == 0 {
		return timeline
	}

	// Calculate gaps between contexts
	var allEntries []output.TimelineEntry
	for i, entry := range entries {
		// Add the context entry
		allEntries = append(allEntries, entry)
		timeline.TotalTracked += entry.Duration

		// Check for gap to next context
		if i < len(entries)-1 {
			nextEntry := entries[i+1]
			var gapStart time.Time
			if entry.EndTime != nil {
				gapStart = *entry.EndTime
			} else {
				// Active context - gap starts now
				gapStart = time.Now()
			}

			gapDuration := nextEntry.StartTime.Sub(gapStart)
			if gapDuration > 0 {
				gapEnd := nextEntry.StartTime
				gapEntry := output.TimelineEntry{
					ContextName: "(no context)",
					StartTime:   gapStart,
					EndTime:     &gapEnd,
					Duration:    gapDuration,
					IsGap:       true,
					IsActive:    false,
				}
				allEntries = append(allEntries, gapEntry)
				timeline.TotalGaps += gapDuration
			}
		}
	}

	timeline.Entries = allEntries
	return timeline
}

// groupByDay groups timeline entries by day
func groupByDay(timeline output.TimelineData) []output.DayGroup {
	if len(timeline.Entries) == 0 {
		return nil
	}

	dayMap := make(map[string]*output.DayGroup)

	for _, entry := range timeline.Entries {
		// Use start time's date
		y, m, d := entry.StartTime.Date()
		dateKey := fmt.Sprintf("%04d-%02d-%02d", y, m, d)
		dayStart := time.Date(y, m, d, 0, 0, 0, 0, entry.StartTime.Location())

		if dayMap[dateKey] == nil {
			dayMap[dateKey] = &output.DayGroup{
				Date:    dayStart,
				Entries: make([]output.TimelineEntry, 0),
			}
		}

		dayMap[dateKey].Entries = append(dayMap[dateKey].Entries, entry)
		if entry.IsGap {
			dayMap[dateKey].TotalGaps += entry.Duration
		} else {
			dayMap[dateKey].TotalTracked += entry.Duration
		}
	}

	// Convert to sorted slice
	var days []output.DayGroup
	for _, group := range dayMap {
		days = append(days, *group)
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i].Date.Before(days[j].Date)
	})

	return days
}

// timelineFromDatabase handles timeline from database backend
func timelineFromDatabase(jsonOutput *bool, projectFilter string, period TimePeriod) error {
	backend, err := core.GetBackend()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("timeline", 2, fmt.Sprintf("failed to get backend: %v", err))
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
			jsonStr, _ := output.FormatJSONError("timeline", 2, err.Error())
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

	// Build timeline
	timeline := buildTimelineFromDBContexts(dbContexts, period, activeContextName)

	// Group by day if multi-day period
	isMultiDay := period.End.Sub(period.Start) > 24*time.Hour
	if isMultiDay {
		timeline.DayGroups = groupByDay(timeline)
	}

	// Output
	if *jsonOutput {
		jsonStr, err := output.FormatTimelineJSON(&timeline)
		if err != nil {
			return err
		}
		fmt.Print(jsonStr)
	} else {
		fmt.Print(output.FormatTimeline(&timeline))
	}

	return nil
}

// NewTimelineCmd creates the timeline command
func NewTimelineCmd(jsonOutput *bool) *cobra.Command {
	var (
		projectFilter string
		today         bool
		week          bool
		month         bool
		since         string
		until         string
	)

	cmd := &cobra.Command{
		Use:     "timeline",
		Aliases: []string{"tl"},
		Short:   "Show visual timeline of context switches",
		Long: `Display a visual timeline showing when contexts were active.

Shows time blocks for each context, gaps between contexts, and summary statistics.
Supports filtering by time period and project.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Calculate time period
			period, err := calculatePeriod(today, week, month, since, until)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("timeline", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Check if using database backend
			if core.IsUsingDatabase() {
				return timelineFromDatabase(jsonOutput, projectFilter, period)
			}

			// File-based backend
			allContexts, err := core.ListContexts()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("timeline", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Get active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("timeline", 2, err.Error())
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

			// Build timeline
			timeline := buildTimelineFromFileContexts(contexts, period, activeContextName)

			// Group by day if multi-day period
			isMultiDay := period.End.Sub(period.Start) > 24*time.Hour
			if isMultiDay {
				timeline.DayGroups = groupByDay(timeline)
			}

			// Output
			if *jsonOutput {
				jsonStr, err := output.FormatTimelineJSON(&timeline)
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Print(output.FormatTimeline(&timeline))
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&projectFilter, "project", "p", "", "Filter by project name")
	cmd.Flags().BoolVarP(&today, "today", "t", false, "Show timeline for today only")
	cmd.Flags().BoolVarP(&week, "week", "w", false, "Show timeline for this week (Mon-Sun)")
	cmd.Flags().BoolVarP(&month, "month", "m", false, "Show timeline for this month")
	cmd.Flags().StringVar(&since, "since", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "End date (YYYY-MM-DD)")

	return cmd
}
