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
			var contextName string

			// If no context specified, use active context
			if len(args) == 0 {
				state, err := core.GetActiveContext()
				if err != nil {
					return fmt.Errorf("failed to get active context: %w", err)
				}
				if !state.HasActiveContext() {
					return fmt.Errorf("no active context. Start a context with: my-context start <name>")
				}
				contextName = state.GetActiveContextName()
			} else {
				contextName = args[0]
			}

			// Get context directory (sanitized name)
			metaPath := core.GetMetaJSONPath(contextName)
			contextDir := filepath.Dir(metaPath)

			// Check if context exists
			if !utils.FileExists(filepath.Join(contextDir, "meta.json")) {
				return fmt.Errorf("context '%s' does not exist", contextName)
			}

			// Parse interval
			var watchInterval time.Duration
			if interval != "" {
				var err error
				watchInterval, err = time.ParseDuration(interval)
				if err != nil {
					return fmt.Errorf("invalid interval: %w", err)
				}
			} else {
				watchInterval = 5 * time.Second // Default
			}

			// Parse timeout
			var watchTimeout time.Duration
			if timeout != "" {
				if timeout == "infinite" {
					watchTimeout = 0 // No timeout
				} else {
					var err error
					watchTimeout, err = time.ParseDuration(timeout)
					if err != nil {
						return fmt.Errorf("invalid timeout: %w", err)
					}
				}
			} else {
				watchTimeout = 0 // No timeout by default
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

// timeoutString returns a string representation of timeout duration
func timeoutString(timeout time.Duration) string {
	if timeout == 0 {
		return "none"
	}
	return timeout.String()
}
