package commands

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/jefferycaldwell/my-context-copilot/internal/watch"
	"github.com/jefferycaldwell/my-context-copilot/pkg/utils"
	"github.com/spf13/cobra"
)

// getContextNameOrActive returns the specified context name or active context
func getContextNameOrActive(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	state, err := core.GetActiveContext()
	if err != nil {
		return "", fmt.Errorf("failed to get active context: %w", err)
	}
	if !state.HasActiveContext() {
		return "", fmt.Errorf("no active context. Start a context with: my-context start <name>")
	}
	return state.GetActiveContextName(), nil
}

// parseWatchInterval parses the interval duration string
func parseWatchInterval(interval string) (time.Duration, error) {
	if interval != "" {
		return time.ParseDuration(interval)
	}
	return 5 * time.Second, nil
}

// parseWatchTimeout parses the timeout duration string
func parseWatchTimeout(timeout string) (time.Duration, error) {
	if timeout == "" || timeout == "infinite" {
		return 0, nil
	}
	return time.ParseDuration(timeout)
}

// timeoutString returns a string representation of the timeout
func timeoutString(d time.Duration) string {
	if d == 0 {
		return "none"
	}
	return d.String()
}

// printWatchStartMessage prints the watch start message
func printWatchStartMessage(contextName string, newNotesOnly bool, pattern, execCommand string, watchInterval, watchTimeout time.Duration) {
	fmt.Printf("Watching context '%s' for changes...\n", contextName)
	if newNotesOnly {
		fmt.Printf("Monitoring: new notes")
		if pattern != "" {
			fmt.Printf(" (pattern: %s)", pattern)
		}
		fmt.Println()
	} else {
		fmt.Println("Monitoring: any file changes")
	}
	if execCommand != "" {
		fmt.Printf("On change, execute: %s\n", execCommand)
	}
	fmt.Printf("Check interval: %s\n", watchInterval)
	if watchTimeout > 0 {
		fmt.Printf("Timeout: %s\n", watchTimeout)
	} else {
		fmt.Println("Timeout: none (run until interrupted)")
	}
	fmt.Println("Press Ctrl+C to stop watching")
}

func NewWatchCmd(jsonOutput *bool) *cobra.Command {
	var (
		pattern      string
		execCommand  string
		timeout      string
		newNotesOnly bool
		interval     string
	)

	cmd := &cobra.Command{
		Use:   "watch <context>",
		Short: "Watch a context for changes and execute commands",
		Long: `Watch a context for changes and optionally execute commands when changes are detected.

By default, watches for any file changes in the context directory. Use --new-notes to watch
specifically for new notes, and --pattern to filter notes by content.

Examples:
  # Watch for any changes in active context
  my-context watch

  # Watch active context for new notes
  my-context watch --new-notes

  # Watch specific context for notes matching pattern
  my-context watch my-context --new-notes --pattern="Phase.*complete"

  # Watch and execute command on changes
  my-context watch --exec="notify-send 'Context updated'"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get context name
			contextName, err := getContextNameOrActive(args)
			if err != nil {
				return err
			}

			// Get context directory
			metaPath := core.GetMetaJSONPath(contextName)
			contextDir := filepath.Dir(metaPath)

			// Check if context exists
			if !utils.FileExists(filepath.Join(contextDir, "meta.json")) {
				return fmt.Errorf("context '%s' does not exist", contextName)
			}

			// Parse interval and timeout
			watchInterval, err := parseWatchInterval(interval)
			if err != nil {
				return fmt.Errorf("invalid interval: %w", err)
			}

			watchTimeout, err := parseWatchTimeout(timeout)
			if err != nil {
				return fmt.Errorf("invalid timeout: %w", err)
			}

			// Create watch options
			options := watch.Options{
				NewNotesOnly: newNotesOnly,
				Pattern:      pattern,
				ExecCommand:  execCommand,
			}

			// Create watcher
			watcher, err := watch.NewWatcher(contextDir, options)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("watch", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Start watching
			if err := watcher.Start(); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("watch", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output start message
			if *jsonOutput {
				data := map[string]interface{}{
					"context":   contextName,
					"new_notes": newNotesOnly,
					"pattern":   pattern,
					"exec":      execCommand,
					"interval":  watchInterval.String(),
					"timeout":   timeoutString(watchTimeout),
				}
				jsonStr, err := output.FormatJSON("watch", data)
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				printWatchStartMessage(contextName, newNotesOnly, pattern, execCommand, watchInterval, watchTimeout)
			}

			// Wait for timeout or interruption
			if watchTimeout > 0 {
				<-time.After(watchTimeout)
				if err := watcher.Stop(); err != nil {
					return fmt.Errorf("error stopping watcher: %w", err)
				}
				if !*jsonOutput {
					fmt.Println("\nWatch timeout reached")
				}
			} else {
				// Wait indefinitely until interrupted
				select {}
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&newNotesOnly, "new-notes", false, "Watch specifically for new notes (default: watch for any changes)")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Regex pattern to match against note content (only applies with --new-notes)")
	cmd.Flags().StringVar(&execCommand, "exec", "", "Command to execute when changes are detected")
	cmd.Flags().StringVar(&interval, "interval", "5s", "How often to check for changes")
	cmd.Flags().StringVar(&timeout, "timeout", "", "Stop watching after this duration (default: run until interrupted)")

	return cmd
}
