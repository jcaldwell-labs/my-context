package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var startProject string
var startForce bool
var startCreatedBy string
var startParent string
var startLabels string

func NewStartCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start <name>",
		Aliases: []string{"s"},
		Short:   "Start a new context",
		Long:    `Start a new context with the given name. If a context is already active, it will be automatically stopped.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			if contextName == "" {
				return fmt.Errorf("context name cannot be empty")
			}

			// Apply project prefix if --project flag is provided
			if startProject != "" {
				contextName = strings.TrimSpace(startProject) + ": " + strings.TrimSpace(contextName)
			}

			// Check for duplicate context name (smart resume)
			existingContext, err := core.FindContextByName(contextName)
			if err == nil && existingContext.Status == "stopped" {
				if startForce {
					// Force flag - never prompt, auto-resolve duplicates
					// Let CreateContext handle auto-suffixing (_2, _3, etc.)
					// This makes --force truly script-friendly by bypassing ALL prompts
				} else {
					// Normal flow - check if interactive
					isInteractive := term.IsTerminal(int(os.Stdin.Fd()))

					if isInteractive {
						// Interactive mode - prompt for resume
						resume, err := promptResume(existingContext)
						if err != nil {
							if *jsonOutput {
								jsonStr, _ := output.FormatJSONError("start", 3, err.Error())
								fmt.Print(jsonStr)
								return nil
							}
							return err
						}

						if resume {
							// Resume the existing context
							return resumeExistingContext(existingContext, jsonOutput)
						} else {
							// User declined resume - prompt for new name
							newName, err := promptNewName(contextName)
							if err != nil {
								if *jsonOutput {
									jsonStr, _ := output.FormatJSONError("start", 3, err.Error())
									fmt.Print(jsonStr)
									return nil
								}
								return err
							}
							contextName = newName
						}
					} else {
						// Non-interactive mode - auto-resume if duplicate
						// Preserves script behavior: duplicate name resumes existing context
						return resumeExistingContext(existingContext, jsonOutput)
					}
				}
			}

			// Parse labels
			var labels []string
			if startLabels != "" {
				labels = strings.Split(strings.ReplaceAll(startLabels, " ", ""), ",")
			}

			// Create the context
			context, previousContext, err := core.CreateContextWithMetadata(contextName, startCreatedBy, startParent, labels)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("start", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.StartData{
					ContextName:  context.Name,
					OriginalName: contextName,
					WasDuplicate: context.Name != core.SanitizeContextName(contextName),
				}
				if previousContext != "" {
					data.PreviousContext = &previousContext
				}
				jsonStr, err := output.FormatJSON("start", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				// Print context home
				fmt.Printf("Context Home: %s\n\n", core.GetContextHomeDisplay())

				// Check if name was modified due to duplicate (e.g., "Bug fix" → "Bug fix_2")
				if context.Name != contextName {
					fmt.Printf("Context \"%s\" already exists.\n", contextName)
				}

				// Show previous context stop message
				if previousContext != "" {
					fmt.Printf("Stopped context: %s\n", previousContext)
				}

				fmt.Printf("✓ Started: %s\n", context.Name)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&startProject, "project", "", "Project name prefix (creates \"project: name\" format)")
	cmd.Flags().BoolVar(&startForce, "force", false, "Force creation of new context without resume prompt")
	cmd.Flags().StringVar(&startCreatedBy, "created-by", "", "User who created this context")
	cmd.Flags().StringVar(&startParent, "parent", "", "Parent context name for hierarchy")
	cmd.Flags().StringVar(&startLabels, "labels", "", "Comma-separated labels for categorization")

	return cmd
}

// promptResume displays context summary and prompts user to resume or create new
func promptResume(ctx *models.Context) (bool, error) {
	noteCount, err := core.GetNoteCount(ctx.Name)
	if err != nil {
		noteCount = 0 // Continue even if we can't get note count
	}

	lastActive, err := core.GetLastActiveTime(ctx.Name)
	if err != nil {
		lastActive = ctx.StartTime // Fallback to start time
	}

	fmt.Printf("Context \"%s\" already exists:\n", ctx.Name)
	fmt.Printf("  Notes: %d\n", noteCount)
	fmt.Printf("  Last active: %s\n", lastActive.Format("2006-01-02 15:04:05"))
	fmt.Print("Resume existing context? [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	switch response {
	case "y", "yes", "":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid response: %s (expected Y/n)", response)
	}
}

// promptNewName prompts for a new context name when user declines resume
func promptNewName(original string) (string, error) {
	fmt.Printf("Enter new name for context (original: %s): ", original)

	reader := bufio.NewReader(os.Stdin)
	newName, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}

	newName = strings.TrimSpace(newName)
	if newName == "" {
		return "", fmt.Errorf("context name cannot be empty")
	}

	return newName, nil
}

// resumeExistingContext resumes an existing stopped context
func resumeExistingContext(ctx *models.Context, jsonOutput *bool) error {
	// Stop any currently active context
	state, err := core.GetActiveContext()
	if err != nil {
		return err
	}

	var previousContext string
	if state.HasActiveContext() {
		previousContext = state.GetActiveContextName()
		_, err := core.StopContext()
		if err != nil {
			return fmt.Errorf("failed to stop previous context: %w", err)
		}
	}

	// Set the existing context as active
	if err := core.SetActiveContext(ctx.Name); err != nil {
		return err
	}

	// Log the transition
	now := time.Now()
	var transitionType models.TransitionType
	if previousContext != "" {
		transitionType = models.TransitionSwitch
	} else {
		transitionType = models.TransitionStart
	}

	transition := &models.ContextTransition{
		Timestamp:      now,
		NewContext:     &ctx.Name,
		TransitionType: transitionType,
	}
	if transitionType == models.TransitionSwitch {
		transition.PreviousContext = &previousContext
	}

	if err := core.AppendLog(core.GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		return err
	}

	// Output
	if *jsonOutput {
		data := output.StartData{
			ContextName:  ctx.Name,
			OriginalName: ctx.Name,
			WasDuplicate: false,
		}
		if previousContext != "" {
			data.PreviousContext = &previousContext
		}
		jsonStr, err := output.FormatJSON("start", map[string]interface{}{"data": data})
		if err != nil {
			return err
		}
		fmt.Print(jsonStr)
	} else {
		if previousContext != "" {
			fmt.Printf("Stopped context: %s\n", previousContext)
		}
		fmt.Printf("Resumed context: %s\n", ctx.Name)
	}

	return nil
}
