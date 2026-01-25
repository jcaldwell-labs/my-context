package commands

import (
	"fmt"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewFocusCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "focus <name|pattern>",
		Aliases: []string{"f"},
		Short:   "Stop current context and resume another in one operation",
		Long:    `Stop the currently active context and resume another context atomically.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if using database backend
			if core.IsUsingDatabase() {
				return focusWithDatabaseBackend(args, jsonOutput)
			}

			// File-based backend
			return focusWithFileBackend(args, jsonOutput)
		},
	}

	return cmd
}

// focusWithFileBackend handles focus command with file-based storage
func focusWithFileBackend(args []string, jsonOutput *bool) error {
	pattern := args[0]
	
	// Check if target is the currently active context first
	state, err := core.GetActiveContext()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, err.Error())
			fmt.Print(jsonStr)
			return nil
		}
		return err
	}

	if state.HasActiveContext() && state.GetActiveContextName() == pattern {
		// Target is already active - no-op
		if *jsonOutput {
			data := map[string]interface{}{
				"message":      "Context already active",
				"context_name": pattern,
			}
			jsonStr, _ := output.FormatJSON("focus", map[string]interface{}{"data": data})
			fmt.Print(jsonStr)
		} else {
			fmt.Printf("Context %q is already active\n", pattern)
		}
		return nil
	}
	
	// Find the target context to focus on (from stopped contexts)
	targetContext, err := findTargetContext(args, false)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 1, err.Error())
			fmt.Print(jsonStr)
			return nil
		}
		return err
	}

	// Stop current context if any
	var stoppedContext *models.Context
	if state.HasActiveContext() {
		stoppedContext, err = core.StopContext()
		if err != nil {
			if *jsonOutput {
				jsonStr, _ := output.FormatJSONError("focus", 2, err.Error())
				fmt.Print(jsonStr)
				return nil
			}
			return err
		}
	}

	// Update target context metadata to clear end time and set status to active
	metaPath := core.GetMetaJSONPath(targetContext.Name)
	var contextMeta models.Context
	if err := core.ReadJSON(metaPath, &contextMeta); err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to read target context metadata: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to read target context metadata: %w", err)
	}
	
	// Clear end time and set to active
	contextMeta.EndTime = nil
	contextMeta.Status = "active"
	
	if err := core.WriteJSON(metaPath, &contextMeta); err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to update target context metadata: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to update target context metadata: %w", err)
	}

	// Set active context in state
	if err := core.SetActiveContext(targetContext.Name); err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to activate target context: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to activate target context: %w", err)
	}

	// Log the transition
	now := time.Now()
	transition := &models.ContextTransition{
		Timestamp:      now,
		PreviousContext: getContextNamePointer(stoppedContext),
		NewContext:     &targetContext.Name,
		TransitionType: models.TransitionStart,
	}

	if err := core.AppendLog(core.GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		// Log error but don't fail the command - transition is still valid
		_ = err
	}

	// Output
	if *jsonOutput {
		data := map[string]interface{}{
			"resumed_context": targetContext.Name,
		}
		if stoppedContext != nil {
			data["stopped_context"] = stoppedContext.Name
			data["stopped_duration_seconds"] = int(stoppedContext.Duration().Seconds())
		}
		jsonStr, _ := output.FormatJSON("focus", map[string]interface{}{"data": data})
		fmt.Print(jsonStr)
	} else {
		if stoppedContext != nil {
			fmt.Printf("✓ Stopped: %s (%s)\n", stoppedContext.Name, output.FormatDuration(stoppedContext.Duration()))
		}
		fmt.Printf("✓ Resumed: %s\n", targetContext.Name)
	}

	return nil
}

// focusWithDatabaseBackend handles focus command with PostgreSQL backend
func focusWithDatabaseBackend(args []string, jsonOutput *bool) error {
	backend, err := core.GetBackend()
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to get backend: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to get backend: %w", err)
	}
	defer backend.Close()

	targetName := args[0]

	// Get the target context to verify it exists
	targetCtx, err := backend.GetContext(targetName)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 1, fmt.Sprintf("target context not found: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("target context not found: %w", err)
	}

	// Get active context
	activeContextName, _ := backend.GetActiveContext()

	// Check if target is already active
	if activeContextName == targetName {
		if *jsonOutput {
			data := map[string]interface{}{
				"message":      "Context already active",
				"context_name": targetName,
			}
			jsonStr, _ := output.FormatJSON("focus", map[string]interface{}{"data": data})
			fmt.Print(jsonStr)
		} else {
			fmt.Printf("Context %q is already active\n", targetName)
		}
		return nil
	}

	// Stop current context if any
	var stoppedCtx *models.Context
	var stoppedDuration time.Duration
	if activeContextName != "" {
		// Get current active context details
		dbCtx, _ := backend.GetContext(activeContextName)

		// Stop it
		endTime := time.Now()
		dbCtx.Status = "stopped"
		dbCtx.EndTime = &endTime

		err = backend.UpdateContext(dbCtx)
		if err != nil {
			if *jsonOutput {
				jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to stop current context: %v", err))
				fmt.Print(jsonStr)
				return nil
			}
			return fmt.Errorf("failed to stop current context: %w", err)
		}

		// Convert for output
		stoppedCtx = &models.Context{
			Name:      dbCtx.Name,
			StartTime: dbCtx.StartTime,
			EndTime:   &endTime,
			Status:    "stopped",
		}
		stoppedDuration = stoppedCtx.Duration()
	}

	// Resume target context
	targetCtx.Status = "active"
	targetCtx.EndTime = nil

	err = backend.UpdateContext(targetCtx)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to update target context: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to update target context: %w", err)
	}

	err = backend.SetActiveContext(targetName)
	if err != nil {
		if *jsonOutput {
			jsonStr, _ := output.FormatJSONError("focus", 2, fmt.Sprintf("failed to set active context: %v", err))
			fmt.Print(jsonStr)
			return nil
		}
		return fmt.Errorf("failed to set active context: %w", err)
	}

	// Output
	if *jsonOutput {
		data := map[string]interface{}{
			"resumed_context": targetName,
		}
		if stoppedCtx != nil {
			data["stopped_context"] = stoppedCtx.Name
			data["stopped_duration_seconds"] = int(stoppedDuration.Seconds())
		}
		jsonStr, _ := output.FormatJSON("focus", map[string]interface{}{"data": data})
		fmt.Print(jsonStr)
	} else {
		fmt.Printf("Context Home: db\n\n")
		if stoppedCtx != nil {
			fmt.Printf("✓ Stopped: %s (%s)\n", stoppedCtx.Name, output.FormatDuration(stoppedDuration))
		}
		fmt.Printf("✓ Resumed: %s\n", targetName)
	}

	return nil
}

// getContextNamePointer returns a pointer to context name or nil
func getContextNamePointer(ctx *models.Context) *string {
	if ctx == nil {
		return nil
	}
	return &ctx.Name
}
