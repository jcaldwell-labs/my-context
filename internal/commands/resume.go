package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

var resumeLast bool

func NewResumeCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resume <name|pattern>",
		Aliases: []string{"r"},
		Short:   "Resume a stopped context",
		Long:    `Resume a previously stopped context by name, pattern, or --last flag.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if we have an active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			if state.HasActiveContext() {
				errMsg := fmt.Sprintf("Cannot resume: context %q is already active", state.GetActiveContextName())
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			}

			var targetContext *models.Context

			// Handle --last flag
			if resumeLast {
				targetContext, err = core.GetMostRecentStopped()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("resume", 1, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}
			} else if len(args) == 0 {
				// No arguments and no --last flag
				errMsg := "Must specify context name/pattern or use --last flag"
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			} else {
				// Handle name/pattern argument
				pattern := args[0]
				contexts, err := core.FindContextsByPattern(pattern)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("resume", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				if len(contexts) == 0 {
					// No matches found
					availableContexts, listErr := getAvailableStoppedContexts()
					errMsg := fmt.Sprintf("No stopped contexts match pattern %q", pattern)
					if listErr == nil && len(availableContexts) > 0 {
						errMsg += fmt.Sprintf(". Available stopped contexts: %s", strings.Join(availableContexts, ", "))
					}

					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
						fmt.Print(jsonStr)
						return nil
					}
					return errors.New(errMsg)
				} else if len(contexts) == 1 {
					// Single match - use it directly
					targetContext = contexts[0]
				} else {
					// Multiple matches - prompt for selection
					targetContext, err = PromptSelection(contexts)
					if err != nil {
						if *jsonOutput {
							jsonStr, _ := output.FormatJSONError("resume", 3, err.Error())
							fmt.Print(jsonStr)
							return nil
						}
						return err
					}
				}
			}

			// Resume the selected context
			return resumeContext(targetContext, jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&resumeLast, "last", false, "Resume the most recently stopped context")

	return cmd
}

// resumeContext resumes a specific context
func resumeContext(ctx *models.Context, jsonOutput *bool) error {
	// Set the context as active
	if err := core.SetActiveContext(ctx.Name); err != nil {
		return fmt.Errorf("failed to activate context: %w", err)
	}

	// Log the transition
	now := time.Now()
	transition := &models.ContextTransition{
		Timestamp:       now,
		NewContext:      &ctx.Name,
		TransitionType:  models.TransitionStart,
	}

	if err := core.AppendLog(core.GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		return fmt.Errorf("failed to log transition: %w", err)
	}

	// Output
	if *jsonOutput {
		data := output.StartData{
			ContextName:  ctx.Name,
			OriginalName: ctx.Name,
			WasDuplicate: false,
		}
		jsonStr, err := output.FormatJSON("resume", map[string]interface{}{"data": data})
		if err != nil {
			return err
		}
		fmt.Print(jsonStr)
	} else {
		fmt.Printf("Resumed context: %s\n", ctx.Name)
	}

	return nil
}

// PromptSelection displays a numbered list of contexts and prompts user to select one
func PromptSelection(contexts []*models.Context) (*models.Context, error) {
	fmt.Println("Multiple contexts match:")
	for i, ctx := range contexts {
		fmt.Printf("  %d. %s\n", i+1, ctx.Name)
	}
	fmt.Print("Select context (1-" + strconv.Itoa(len(contexts)) + "): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(response)
	selection, err := strconv.Atoi(response)
	if err != nil || selection < 1 || selection > len(contexts) {
		return nil, fmt.Errorf("invalid selection: %s (expected 1-%d)", response, len(contexts))
	}

	return contexts[selection-1], nil
}

// getAvailableStoppedContexts returns a list of available stopped context names for error messages
func getAvailableStoppedContexts() ([]string, error) {
	contexts, err := core.ListContexts()
	if err != nil {
		return nil, err
	}

	var stopped []string
	for _, ctx := range contexts {
		if ctx.Status == "stopped" {
			stopped = append(stopped, ctx.Name)
		}
	}

	return stopped, nil
}
