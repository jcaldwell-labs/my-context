package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/jefferycaldwell/my-context-copilot/pkg/storage"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

const (
	// recordDebounceDelay prevents duplicate captures from rapid clipboard changes
	recordDebounceDelay = 100 * time.Millisecond
	// maxPreviewLength is the maximum length of text shown in the capture preview
	maxPreviewLength = 50
)

// truncatePreview truncates text for display preview
func truncatePreview(s string, maxLen int) string {
	// Remove newlines for preview
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.TrimSpace(s)

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// recordFromClipboard watches clipboard and records changes as notes
func recordFromClipboard(prefix string, jsonOutput bool) error {
	// Verify active context exists and get backend if using database
	var activeContextName string
	var backend storage.Backend

	if core.IsUsingDatabase() {
		var err error
		backend, err = core.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}
		defer backend.Close()

		activeContextName, err = backend.GetActiveContext()
		if err != nil || activeContextName == "" {
			return fmt.Errorf("no active context. Start one with: my-context start <name>")
		}
	} else {
		state, err := core.GetActiveContext()
		if err != nil {
			return fmt.Errorf("failed to get active context: %w", err)
		}
		activeContextName = state.GetActiveContextName()
		if activeContextName == "" {
			return fmt.Errorf("no active context. Start one with: my-context start <name>")
		}
	}

	// Initialize clipboard
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to initialize clipboard: %w (clipboard access may require X11/Wayland on Linux)", err)
	}

	// Set up signal handling for graceful exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Print header
	if !jsonOutput {
		fmt.Printf("Recording to context: \"%s\"\n", activeContextName)
		fmt.Println("Paste to capture (Ctrl+C to stop)")
		fmt.Println("───────────────────────────────────")
	}

	// Watch clipboard
	clipChan := clipboard.Watch(ctx, clipboard.FmtText)

	captureCount := 0
	var lastContent string
	var lastCaptureTime time.Time

	// Main loop
	for {
		select {
		case <-sigChan:
			// Graceful exit
			if !jsonOutput {
				fmt.Printf("\n\nCaptured %d notes. Run 'my-context show' to review.\n", captureCount)
			} else {
				jsonStr, _ := output.FormatJSON("record", map[string]interface{}{
					"context_name":  activeContextName,
					"capture_count": captureCount,
					"status":        "stopped",
				})
				fmt.Print(jsonStr)
			}
			return nil

		case data := <-clipChan:
			content := string(data)

			// Skip empty content
			if strings.TrimSpace(content) == "" {
				continue
			}

			// Debounce: skip if same content within debounce delay
			if content == lastContent && time.Since(lastCaptureTime) < recordDebounceDelay {
				continue
			}

			// Update tracking
			lastContent = content
			lastCaptureTime = time.Now()

			// Build note text
			noteText := strings.TrimSpace(content)
			if prefix != "" {
				noteText = fmt.Sprintf("[%s] %s", prefix, noteText)
			}

			// Add note (reuse backend connection if using database)
			var err error
			if core.IsUsingDatabase() {
				timestamp := time.Now().Format(time.RFC3339)
				err = backend.AddNote(activeContextName, timestamp, noteText)
			} else {
				_, err = core.AddNote(noteText)
			}

			if err != nil {
				if !jsonOutput {
					fmt.Printf("Error adding note: %v\n", err)
				}
				continue
			}

			captureCount++

			// Display preview
			if !jsonOutput {
				preview := truncatePreview(content, maxPreviewLength)
				fmt.Printf("[%d] %s  ✓\n", captureCount, preview)
			}
		}
	}
}

// NewRecordCmd creates the record command
func NewRecordCmd(jsonOutput *bool) *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:     "record",
		Aliases: []string{"r"},
		Short:   "Record clipboard changes as notes",
		Long: `Start clipboard watching mode, automatically capturing each clipboard
change as a note to the active context.

Useful for quickly capturing URLs, snippets, code blocks, or other text
from multiple sources without manual 'my-context note' calls.

Press Ctrl+C to stop recording.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := recordFromClipboard(prefix, *jsonOutput)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("record", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix tag for captured notes (e.g., 'url', 'clip')")

	return cmd
}
