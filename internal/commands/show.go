package commands

import (
	"fmt"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

// buildShowJSONOutput creates JSON output for the show command with stale information
func buildShowJSONOutput(ctx interface{}, internalCtx *models.Context, notes []*models.Note, files []*models.FileAssociation, touches []*models.TouchEvent) (string, error) {
	lastActivity := core.GetLastActivityTime(internalCtx, notes, files, touches)
	thresholds := core.GetStaleThresholds()
	staleLevel := core.GetStaleLevel(lastActivity, thresholds)
	isStale := core.IsStale(lastActivity, thresholds)
	inactiveSince := time.Since(lastActivity).Hours()

	data := output.ContextData{
		Context:        ctx,
		Notes:          notes,
		Files:          files,
		Touches:        touches,
		IsStale:        isStale,
		StaleLevel:     core.GetStaleLevelString(staleLevel),
		LastActivity:   &lastActivity,
		InactiveSinceH: inactiveSince,
	}
	return output.FormatJSON("show", map[string]interface{}{"data": data})
}

// buildStaleInfo calculates stale information for active contexts
func buildStaleInfo(internalCtx *models.Context, notes []*models.Note, files []*models.FileAssociation, touches []*models.TouchEvent) *output.StaleInfo {
	if internalCtx.Status != "active" {
		return nil
	}
	lastActivity := core.GetLastActivityTime(internalCtx, notes, files, touches)
	thresholds := core.GetStaleThresholds()
	staleLevel := core.GetStaleLevel(lastActivity, thresholds)
	return &output.StaleInfo{
		IsActive:          true,
		LastActivity:      lastActivity,
		StaleLevel:        core.GetStaleLevelString(staleLevel),
		TimeSinceActivity: time.Since(lastActivity),
		ContextStartTime:  internalCtx.StartTime,
	}
}

func NewShowCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [context-name]",
		Aliases: []string{"w"},
		Short:   "Show context details",
		Long:    `Display details about the currently active context including notes, files, and touch events.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var contextName string

			// Check if using database backend
			if core.IsUsingDatabase() {
				backend, err := core.GetBackend()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("show", 2, fmt.Sprintf("failed to get backend: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to get backend: %w", err)
				}
				defer backend.Close()

				// Get context name
				if len(args) > 0 {
					contextName = args[0]
				} else {
					// Get active context from database
					contextName, err = backend.GetActiveContext()
					if err != nil || contextName == "" {
						errMsg := "No active context"
						if *jsonOutput {
							jsonStr, _ := output.FormatJSONError("show", 1, errMsg)
							fmt.Print(jsonStr)
						} else {
							fmt.Println(errMsg)
							fmt.Println("Start one with: my-context start <name>")
						}
						return nil
					}
				}

				// Get context from database
				dbCtx, err := backend.GetContext(contextName)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("show", 2, fmt.Sprintf("context not found: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("context not found: %w", err)
				}

				// Get notes and files from database
				dbNotes, _ := backend.GetNotes(contextName)
				dbFiles, _ := backend.GetFiles(contextName)

				// Convert to internal models
				context := &models.Context{
					Name:      dbCtx.Name,
					StartTime: dbCtx.StartTime,
					EndTime:   dbCtx.EndTime,
					Status:    dbCtx.Status,
				}

				var notes []*models.Note
				for _, n := range dbNotes {
					notes = append(notes, &models.Note{
						Timestamp:   n.Timestamp,
						TextContent: n.Content,
					})
				}

				var files []*models.FileAssociation
				for _, f := range dbFiles {
					files = append(files, &models.FileAssociation{
						Timestamp: f.Timestamp,
						FilePath:  f.Path,
					})
				}

				// Build touch events from database touch count and last touch time
				// We create synthetic touch events to match the expected display format
				var touches []*models.TouchEvent
				if dbCtx.TouchCount > 0 && dbCtx.LastTouchAt != nil {
					// Create synthetic touch events - one per touch count
					// All with the last touch timestamp (we don't track individual times in DB)
					for i := 0; i < dbCtx.TouchCount; i++ {
						touches = append(touches, &models.TouchEvent{
							Timestamp: *dbCtx.LastTouchAt,
						})
					}
				}

				// Output
				if *jsonOutput {
					jsonStr, err := buildShowJSONOutput(dbCtx, context, notes, files, touches)
					if err != nil {
						return err
					}
					fmt.Print(jsonStr)
				} else {
					// Build display string for database mode with partition info
					partition := core.ExtractPartition()
					var homeDisplay string
					if partition != "" {
						homeDisplay = fmt.Sprintf("db:%s", partition)
					} else {
						homeDisplay = "db"
					}

					// Get context count for this partition
					contextCount, _ := backend.GetContextCount()

					// Print header with partition and count
					output.PrintContextHomeHeader(homeDisplay, contextCount)

					// Calculate stale info for display
					staleInfo := buildStaleInfo(context, notes, files, touches)
					fmt.Print(output.FormatContextWithStaleInfo(context, notes, files, touches, staleInfo))
				}

				return nil
			}

			// File-based backend (existing code)
			// If context name provided as argument, use it
			if len(args) > 0 {
				contextName = args[0]
			} else {
				// No argument - show active context (backward compatible)
				state, err := core.GetActiveContext()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("show", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				if !state.HasActiveContext() {
					errMsg := "No active context"
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("show", 1, errMsg)
						fmt.Print(jsonStr)
					} else {
						fmt.Println(errMsg)
						fmt.Println("Start one with: my-context start <name>")
					}
					return nil
				}

				contextName = state.GetActiveContextName()
			}

			// Get context details (works for stopped or active)
			context, notes, files, touches, err := core.GetContextWithMetadata(contextName)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("show", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Convert to internal context model for stale calculation
			internalCtx := &models.Context{
				Name:      context.Name,
				StartTime: context.StartTime,
				EndTime:   context.EndTime,
				Status:    context.Status,
			}

			// Output
			if *jsonOutput {
				jsonStr, err := buildShowJSONOutput(context, internalCtx, notes, files, touches)
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				// Print context home header
				output.PrintContextHomeHeader(core.GetContextHomeDisplay(), core.GetContextCount())

				// Calculate stale info for display
				staleInfo := buildStaleInfo(internalCtx, notes, files, touches)
				fmt.Print(output.FormatContextWithStaleInfo(context, notes, files, touches, staleInfo))
			}

			return nil
		},
	}

	return cmd
}
