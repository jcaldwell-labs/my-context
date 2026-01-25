package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewStopCmd(jsonOutput *bool) *cobra.Command {
	var cascade bool
	
	cmd := &cobra.Command{
		Use:     "stop",
		Aliases: []string{"p"},
		Short:   "Stop the active context",
		Long: `Stop the currently active context without starting a new one.

If the active context has child contexts that are still active, a warning will be
displayed by default. Use the --cascade flag to automatically stop all child
contexts along with the parent.

Examples:
  my-context stop              # Stop active context, warn if children exist
  my-context stop --cascade    # Stop active context and all its children`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if using database backend
			if core.IsUsingDatabase() {
				backend, err := core.GetBackend()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("stop", 2, fmt.Sprintf("failed to get backend: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to get backend: %w", err)
				}
				defer backend.Close()

				// Get active context
				contextName, err := backend.GetActiveContext()
				if err != nil || contextName == "" {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{
							"message": "No active context",
						})
						fmt.Print(jsonStr)
					} else {
						fmt.Println("No active context")
					}
					return nil
				}

				// Get context details
				dbCtx, _ := backend.GetContext(contextName)

				// Check for active children
				activeChildren, err := core.GetActiveChildrenDB(backend, contextName)
				if err != nil {
					activeChildren = []string{} // Continue even if we can't get children
				}

				var stoppedChildren []string
				var warningShown bool

				// Handle children based on cascade flag
				if len(activeChildren) > 0 && !cascade {
					// Warning mode: show warning but don't stop children
					if !*jsonOutput {
						fmt.Printf("⚠️  Warning: %d active child context(s) will remain active:\n", len(activeChildren))
						for _, child := range activeChildren {
							childCtx, _ := backend.GetContext(child)
							if childCtx != nil {
								duration := time.Since(childCtx.StartTime)
								fmt.Printf("  - %s (active, %s)\n", child, output.FormatDuration(duration))
							}
						}
						fmt.Printf("\nUse --cascade to stop them along with the parent.\n\n")
					}
					warningShown = true
				} else if len(activeChildren) > 0 && cascade {
					// Cascade mode: stop all descendants (children, grandchildren, etc.)
					allDescendants, err := core.GetAllDescendantsDB(backend, contextName)
					if err != nil {
						if *jsonOutput {
							jsonStr, _ := output.FormatJSONError("stop", 2, fmt.Sprintf("failed to get descendants: %v", err))
							fmt.Print(jsonStr)
							return nil
						}
						return fmt.Errorf("failed to get descendants: %w", err)
					}

					endTime := time.Now()

					// Stop descendants in reverse order (leaves first, then parents)
					for i := len(allDescendants) - 1; i >= 0; i-- {
						descendantName := allDescendants[i]
						descendantCtx, err := backend.GetContext(descendantName)
						if err != nil || descendantCtx.Status != "active" {
							continue // Skip if not found or not active
						}

						descendantCtx.Status = "stopped"
						descendantCtx.EndTime = &endTime

						if err := backend.UpdateContext(descendantCtx); err != nil {
							// Log error but continue with other children
							if !*jsonOutput {
								fmt.Fprintf(cmd.ErrOrStderr(), "Warning: Failed to stop %s: %v\n", descendantName, err)
							}
							continue
						}

						stoppedChildren = append(stoppedChildren, descendantName)

						if !*jsonOutput {
							duration := endTime.Sub(descendantCtx.StartTime)
							fmt.Printf("✓ Stopped: %s (%s)\n", descendantName, output.FormatDuration(duration))
						}
					}
				}

				// Update parent context to stopped
				endTime := time.Now()
				dbCtx.Status = "stopped"
				dbCtx.EndTime = &endTime

				err = backend.UpdateContext(dbCtx)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("stop", 2, fmt.Sprintf("failed to update context: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to update context: %w", err)
				}

				// Clear active context
				err = backend.ClearActiveContext()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("stop", 2, fmt.Sprintf("failed to clear active context: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to clear active context: %w", err)
				}

				// Convert to internal model for output
				context := &models.Context{
					Name:      dbCtx.Name,
					StartTime: dbCtx.StartTime,
					EndTime:   &endTime,
					Status:    "stopped",
				}

				// Output
				if *jsonOutput {
					durationSeconds := int(context.Duration().Seconds())
					data := output.StopData{
						ContextName:     context.Name,
						StartTime:       context.StartTime,
						EndTime:         endTime,
						DurationSeconds: durationSeconds,
					}
					if len(stoppedChildren) > 0 {
						data.StoppedChildren = stoppedChildren
					}
					if warningShown {
						data.ActiveChildrenCount = len(activeChildren)
						jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{
							"data":    data,
							"warning": fmt.Sprintf("%d active child context(s) will remain active. Use --cascade to stop them.", len(activeChildren)),
						})
						fmt.Print(jsonStr)
					} else {
						jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{"data": data})
						fmt.Print(jsonStr)
					}
				} else {
					fmt.Printf("✓ Stopped: %s (%s)\n",
						context.Name,
						output.FormatDuration(context.Duration()))
				}

				return nil
			}

			// File-based backend
			// Get active context name before stopping
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stop", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			if !state.HasActiveContext() {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{
						"message": "No active context",
					})
					fmt.Print(jsonStr)
				} else {
					fmt.Println("No active context")
				}
				return nil
			}

			contextName := state.GetActiveContextName()

			// Check for active children
			activeChildren, err := core.GetActiveChildren(contextName)
			if err != nil {
				activeChildren = []string{} // Continue even if we can't get children
			}

			var stoppedChildren []string
			var warningShown bool

			// Handle children based on cascade flag
			if len(activeChildren) > 0 && !cascade {
				// Warning mode: show warning but don't stop children
				if !*jsonOutput {
					fmt.Printf("⚠️  Warning: %d active child context(s) will remain active:\n", len(activeChildren))
					for _, child := range activeChildren {
						childCtx, err := core.LoadContext(child)
						if err == nil && childCtx.Status == "active" {
							duration := time.Since(childCtx.StartTime)
							fmt.Printf("  - %s (active, %s)\n", child, output.FormatDuration(duration))
						}
					}
					fmt.Printf("\nUse --cascade to stop them along with the parent.\n\n")
				}
				warningShown = true
			} else if len(activeChildren) > 0 && cascade {
				// Cascade mode: stop all descendants (children, grandchildren, etc.)
				allDescendants, err := core.GetAllDescendants(contextName)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("stop", 2, fmt.Sprintf("failed to get descendants: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to get descendants: %w", err)
				}

				endTime := time.Now()

				// Stop descendants in reverse order (leaves first, then parents)
				for i := len(allDescendants) - 1; i >= 0; i-- {
					descendantName := allDescendants[i]
					descendantCtx, err := core.LoadContext(descendantName)
					if err != nil || descendantCtx.Status != "active" {
						continue // Skip if not found or not active
					}

					// Stop the descendant context
					descendantCtx.EndTime = &endTime
					descendantCtx.Status = "stopped"

					if err := core.WriteJSON(core.GetMetaJSONPath(descendantName), descendantCtx); err != nil {
						// Log error but continue with other children
						if !*jsonOutput {
							fmt.Fprintf(cmd.ErrOrStderr(), "Warning: Failed to stop %s: %v\n", descendantName, err)
						}
						continue
					}

					stoppedChildren = append(stoppedChildren, descendantName)

					if !*jsonOutput {
						duration := endTime.Sub(descendantCtx.StartTime)
						fmt.Printf("✓ Stopped: %s (%s)\n", descendantName, output.FormatDuration(duration))
					}
				}
			}

			// Stop the parent context
			context, err := core.StopContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stop", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// No active context (shouldn't happen, but handle it)
			if context == nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{
						"message": "No active context",
					})
					fmt.Print(jsonStr)
				} else {
					fmt.Println("No active context")
				}
				return nil
			}

			// Output
			if *jsonOutput {
				durationSeconds := int(context.Duration().Seconds())
				data := output.StopData{
					ContextName:     context.Name,
					StartTime:       context.StartTime,
					EndTime:         *context.EndTime,
					DurationSeconds: durationSeconds,
				}
				if len(stoppedChildren) > 0 {
					data.StoppedChildren = stoppedChildren
				}
				if warningShown {
					data.ActiveChildrenCount = len(activeChildren)
					jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{
						"data":    data,
						"warning": fmt.Sprintf("%d active child context(s) will remain active. Use --cascade to stop them.", len(activeChildren)),
					})
					fmt.Print(jsonStr)
				} else {
					jsonStr, err := output.FormatJSON("stop", map[string]interface{}{"data": data})
					if err != nil {
						return err
					}
					fmt.Print(jsonStr)
				}
			} else {
				fmt.Printf("✓ Stopped: %s (%s)\n",
					context.Name,
					output.FormatDuration(context.Duration()))

				// Display lifecycle advisor guidance
				err = displayLifecycleGuidance(context)
				if err != nil {
					// Log error but don't fail the command
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: Could not display lifecycle guidance: %v\n", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&cascade, "cascade", false, "Stop child contexts along with the parent")
	
	return cmd
}

// displayLifecycleGuidance shows helpful suggestions after stopping a context
func displayLifecycleGuidance(context *models.Context) error {
	// Get context summary information
	noteCount, err := core.GetNoteCount(context.Name)
	if err != nil {
		return fmt.Errorf("failed to get note count: %w", err)
	}

	// Display context summary
	fmt.Printf("\n📊 Context Summary:\n")
	fmt.Printf("  Name: %s\n", context.Name)
	fmt.Printf("  Duration: %s\n", output.FormatDuration(context.Duration()))
	fmt.Printf("  Notes: %d\n", noteCount)

	// Get full context data for completion detection
	_, notes, _, _, err := core.GetContext(context.Name)
	if err != nil {
		return fmt.Errorf("failed to get context data: %w", err)
	}

	// Find related contexts
	relatedContexts, err := core.FindRelatedContexts(context.Name)
	if err != nil {
		return fmt.Errorf("failed to find related contexts: %w", err)
	}

	// Detect completion
	isComplete := detectCompletion(notes)

	// Display guidance
	displayGuidance(context, relatedContexts, isComplete)

	return nil
}

// displayGuidance shows suggestions based on context analysis
func displayGuidance(context *models.Context, relatedContexts []*models.Context, isComplete bool) {
	fmt.Printf("\n💡 Next Steps:\n")

	hasSuggestions := false

	// Suggestion 1: Resume related context
	if len(relatedContexts) > 0 {
		fmt.Printf("  🔄 Resume related work:\n")
		for i, related := range relatedContexts {
			if i >= 3 { // Limit to 3 suggestions
				break
			}
			fmt.Printf("    • %s\n", related.Name)
		}
		hasSuggestions = true
	}

	// Suggestion 2: Archive if complete
	if isComplete {
		fmt.Printf("  📦 Archive completed context:\n")
		fmt.Printf("    • my-context archive %q\n", context.Name)
		hasSuggestions = true
	} else {
		fmt.Printf("  📦 Consider archiving when complete:\n")
		fmt.Printf("    • my-context archive %q\n", context.Name)
	}

	// Suggestion 3: Start new work
	fmt.Printf("  ✨ Start new work:\n")
	fmt.Printf("    • my-context start \"new-feature-name\"\n")

	if !hasSuggestions {
		fmt.Printf("  No specific suggestions - you're all set!\n")
	}
}

// detectCompletion checks if the last 5 notes contain completion keywords
func detectCompletion(notes []*models.Note) bool {
	if len(notes) == 0 {
		return false
	}

	// Define completion keywords (case-insensitive)
	completionKeywords := []string{
		"complete", "completed", "done", "finished", "finish",
		"closed", "resolved", "fixed", "implemented",
		"merged", "deployed", "released",
	}

	// Check last 5 notes (or all notes if fewer than 5)
	startIdx := len(notes) - 5
	if startIdx < 0 {
		startIdx = 0
	}

	for i := len(notes) - 1; i >= startIdx; i-- {
		note := notes[i]
		noteText := strings.ToLower(note.TextContent)

		for _, keyword := range completionKeywords {
			if strings.Contains(noteText, keyword) {
				return true
			}
		}
	}

	return false
}
