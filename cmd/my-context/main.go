package main

import (
	"fmt"
	"os"

	"github.com/jefferycaldwell/my-context-copilot/internal/commands"
	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool

	// Version information - set via ldflags during build
	Version   = "2.0.0-dev+debian"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "my-context",
	Short: "Context management tool for developers",
	Long: `my-context is a CLI tool for managing work contexts.
Track your work sessions with notes, file associations, and timestamps.`,
	Version: formatVersion(),
}

// formatVersion returns a formatted version string with build metadata
func formatVersion() string {
	return fmt.Sprintf("%s (build: %s, commit: %s)", Version, BuildTime, GitCommit)
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")

	// Add all subcommands
	rootCmd.AddCommand(commands.NewStartCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewStopCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewNoteCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewFileCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewTouchCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewShowCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewListCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewHistoryCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewExportCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewArchiveCmd(&jsonOutput))
	rootCmd.AddCommand(commands.NewDeleteCmd(&jsonOutput))
}

func main() {
	// Ensure context home directory exists
	if err := core.EnsureContextHome(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize context home: %v\n", err)
		os.Exit(2)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
