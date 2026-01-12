package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

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
					data := map[string]interface{}{
						"context":     dbCtx,
						"notes":       notes,
						"files":       files,
						"touches":     touches,
						"touch_count": dbCtx.TouchCount,
					}
					jsonStr, err := output.FormatJSON("show", map[string]interface{}{"data": data})
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
					fmt.Print(output.FormatContext(context, notes, files, touches))
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

			// Output
			if *jsonOutput {
				data := output.ContextData{
					Context: context,
					Notes:   notes,
					Files:   files,
					Touches: touches,
				}
				jsonStr, err := output.FormatJSON("show", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				// Print context home header
				output.PrintContextHomeHeader(core.GetContextHomeDisplay(), core.GetContextCount())
				fmt.Print(output.FormatContext(context, notes, files, touches))
			}

			return nil
		},
	}

	return cmd
}
