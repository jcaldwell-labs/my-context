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
		return nil, errors.New("must specify context name/pattern or use --last flag")
	}

	pattern := args[0]
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
		Use:     "resume <name|pattern>",
		Aliases: []string{"r"},
		Short:   "Resume a stopped context",
		Long:    `Resume a previously stopped context by name, pattern, or --last flag.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if using database backend
			if core.IsUsingDatabase() {
				return resumeWithDatabaseBackend(args, resumeLast, jsonOutput)
			}

			// File-based backend (existing code)
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

	// Check if there's already an active context
	activeContext, _ := backend.GetActiveContext()
	if activeContext != "" {
		errMsg := fmt.Sprintf("Cannot resume: context %q is already active", activeContext)
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
			fmt.Print(jsonStr)
			return nil
		}
		return errors.New(errMsg)
	}

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
		errMsg := "Must specify context name or use --last flag"
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
			fmt.Print(jsonStr)
			return nil
		}
		return errors.New(errMsg)
	} else {
		contextName = args[0]
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

	if ctx.Status == "active" {
		errMsg := fmt.Sprintf("Context %q is already active", contextName)
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("resume", 1, errMsg)
			fmt.Print(jsonStr)
			return nil
		}
		return errors.New(errMsg)
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
