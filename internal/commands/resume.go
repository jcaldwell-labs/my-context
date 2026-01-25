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

// findTargetContext finds the context to resume based on arguments and flags
func findTargetContext(args []string, useLast bool) (*models.Context, error) {
	if useLast {
		return core.GetMostRecentStopped()
	}

	if len(args) == 0 {
		return nil, errors.New("must specify context name/pattern/index or use --last flag")
	}

	pattern := args[0]

	// Check if the argument is a numeric index
	if index, err := strconv.Atoi(pattern); err == nil {
		// Get all contexts in list order
		allContexts, err := core.ListContexts()
		if err != nil {
			return nil, fmt.Errorf("failed to list contexts: %w", err)
		}

		// Check if there are any contexts
		if len(allContexts) == 0 {
			return nil, fmt.Errorf("no contexts available to resume")
		}

		// Validate index range (1-based indexing)
		if index < 1 || index > len(allContexts) {
			return nil, fmt.Errorf("index %d is out of range (valid range: 1-%d)", index, len(allContexts))
		}

		// Get the context at the specified index (convert to 0-based)
		targetContext := allContexts[index-1]

		return targetContext, nil
	}

	// Not a numeric index, treat as pattern
	contexts, err := core.FindContextsByPattern(pattern)
	if err != nil {
		return nil, err
	}

	switch len(contexts) {
	case 0:
		availableContexts, listErr := getAvailableStoppedContexts()
		errMsg := fmt.Sprintf("No stopped contexts match pattern %q", pattern)
		if listErr == nil && len(availableContexts) > 0 {
			errMsg += fmt.Sprintf(". Available stopped contexts: %s", strings.Join(availableContexts, ", "))
		}
		return nil, errors.New(errMsg)
	case 1:
		return contexts[0], nil
	default:
		return PromptSelection(contexts)
	}
}

func NewResumeCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resume <name|pattern|index>",
		Aliases: []string{"r"},
		Short:   "Resume a stopped context",
		Long:    `Resume a previously stopped context by name, pattern, index, or --last flag.

Examples:
  my-context resume work-2026-01-15    # Resume by name
  my-context resume work               # Resume by pattern
  my-context resume 2                  # Resume by index from list
  my-context resume --last             # Resume most recent stopped`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if using database backend
			if core.IsUsingDatabase() {
				return resumeWithDatabaseBackend(args, resumeLast, jsonOutput)
			}

			// File-based backend (existing code)
			// Find the target context to resume
			targetContext, err := findTargetContext(args, resumeLast)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Check if we have an active context and stop it if needed
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
				activeContextName := state.GetActiveContextName()
				// Don't allow resuming the already active context
				if targetContext.Name == activeContextName {
					errMsg := fmt.Sprintf("Cannot resume: context %q is already active", activeContextName)
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
						fmt.Print(jsonStr)
						return nil
					}
					return errors.New(errMsg)
				}

				// Stop the active context before resuming
				_, err := core.StopContext()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("resume", 2, fmt.Sprintf("failed to stop active context: %v", err))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("failed to stop active context: %w", err)
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
	// Update context metadata to set status to "active" and clear end time
	ctx.Status = "active"
	ctx.EndTime = nil
	if err := core.WriteJSON(core.GetMetaJSONPath(ctx.Name), ctx); err != nil {
		return fmt.Errorf("failed to update context metadata: %w", err)
	}

	// Set the context as active in state
	if err := core.SetActiveContext(ctx.Name); err != nil {
		return fmt.Errorf("failed to activate context: %w", err)
	}

	// Log the transition
	now := time.Now()
	transition := &models.ContextTransition{
		Timestamp:      now,
		NewContext:     &ctx.Name,
		TransitionType: models.TransitionStart,
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

// resumeWithDatabaseBackend resumes a context using the PostgreSQL backend
func resumeWithDatabaseBackend(args []string, useLast bool, jsonOutput *bool) error {
	backend, err := core.GetBackend()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 2, fmt.Sprintf("failed to get backend: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to get backend: %w", err)
	}
	defer backend.Close()

	var contextName string

	if useLast {
		// Get most recent stopped context
		contexts, err := backend.ListContexts()
		if err != nil || len(contexts) == 0 {
			errMsg := "No stopped contexts found"
			if *jsonOutput {
				jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
				fmt.Print(jsonStr)
				return nil
			}
			return errors.New(errMsg)
		}

		// Find first stopped context
		for _, ctx := range contexts {
			if ctx.Status == "stopped" {
				contextName = ctx.Name
				break
			}
		}

		if contextName == "" {
			errMsg := "No stopped contexts found"
			if *jsonOutput {
				jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
				fmt.Print(jsonStr)
				return nil
			}
			return errors.New(errMsg)
		}
	} else if len(args) == 0 {
		errMsg := "Must specify context name/index or use --last flag"
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
			fmt.Print(jsonStr)
			return nil
		}
		return errors.New(errMsg)
	} else {
		// Check if the argument is a numeric index
		if index, err := strconv.Atoi(args[0]); err == nil {
			// Get all contexts in list order
			allContexts, err := backend.ListContexts()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 2, fmt.Sprintf("failed to list contexts: %v", err))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("failed to list contexts: %w", err)
			}

			// Check if there are any contexts
			if len(allContexts) == 0 {
				errMsg := "no contexts available to resume"
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			}

			// Validate index range (1-based indexing)
			if index < 1 || index > len(allContexts) {
				errMsg := fmt.Sprintf("index %d is out of range (valid range: 1-%d)", index, len(allContexts))
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			}

			// Get the context at the specified index (convert to 0-based)
			targetContext := allContexts[index-1]
			contextName = targetContext.Name
		} else {
			// Not a numeric index, treat as context name
			contextName = args[0]
		}
	}

	// Check if there's already an active context
	activeContext, _ := backend.GetActiveContext()
	if activeContext != "" {
		// Don't allow resuming the already active context
		if contextName == activeContext {
			errMsg := fmt.Sprintf("Cannot resume: context %q is already active", activeContext)
			if *jsonOutput {
				jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
				fmt.Print(jsonStr)
				return nil
			}
			return errors.New(errMsg)
		}

		// Stop the active context before resuming
		ctx, err := backend.GetContext(activeContext)
		if err == nil && ctx.Status == "active" {
			now := time.Now()
			ctx.Status = "stopped"
			ctx.EndTime = &now
			_ = backend.UpdateContext(ctx)
		}
		_ = backend.ClearActiveContext()
	}

	// Get the context
	ctx, err := backend.GetContext(contextName)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 2, fmt.Sprintf("context not found: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("context not found: %w", err)
	}

	// Resume the context (clear end time, set to active)
	ctx.Status = "active"
	ctx.EndTime = nil

	err = backend.UpdateContext(ctx)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 2, fmt.Sprintf("failed to update context: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to update context: %w", err)
	}

	err = backend.SetActiveContext(contextName)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 2, fmt.Sprintf("failed to set active: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to set active: %w", err)
	}

	// Output
	if *jsonOutput {
		data := map[string]interface{}{"context_name": contextName, "status": "resumed"}
		jsonStr, _ := output.FormatJSON("resume", map[string]interface{}{"data": data})
		fmt.Print(jsonStr)
	} else {
		fmt.Printf("Context Home: db\n\n")
		fmt.Printf("✓ Resumed: %s\n", contextName)
	}

	return nil
}
