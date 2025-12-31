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
	cmd := &cobra.Command{
		Use:     "stop",
		Aliases: []string{"p"},
		Short:   "Stop the active context",
		Long:    `Stop the currently active context without starting a new one.`,
		Args:    cobra.NoArgs,
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

				// Update context to stopped
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

				// Output (simplified, skipping lifecycle guidance for database backend)
				if *jsonOutput {
					durationSeconds := int(context.Duration().Seconds())
					data := output.StopData{
						ContextName:     context.Name,
						StartTime:       context.StartTime,
						EndTime:         endTime,
						DurationSeconds: durationSeconds,
					}
					jsonStr, _ := output.FormatJSON("stop", map[string]interface{}{"data": data})
					fmt.Print(jsonStr)
				} else {
					fmt.Printf("Stopped context: %s (duration: %s)\n",
						context.Name,
						output.FormatDuration(context.Duration()))
				}

				return nil
			}

			// File-based backend (existing code)
			context, err := core.StopContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("stop", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// No active context
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
				jsonStr, err := output.FormatJSON("stop", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Stopped context: %s (duration: %s)\n",
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
