package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/jefferycaldwell/my-context-copilot/internal/signal"
	"github.com/spf13/cobra"
)

func NewSignalCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal",
		Short: "Manage signal files for event coordination",
		Long: `Manage signal files for event coordination between processes and team members.

Signal files are simple timestamped files stored in ~/.my-context/signals/ that can be used
to coordinate events like binary updates, context changes, or custom notifications.`,
	}

	// Add subcommands
	cmd.AddCommand(newSignalCreateCmd(jsonOutput))
	cmd.AddCommand(newSignalListCmd(jsonOutput))
	cmd.AddCommand(newSignalWaitCmd(jsonOutput))
	cmd.AddCommand(newSignalClearCmd(jsonOutput))

	return cmd
}

func newSignalCreateCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a signal file",
		Long:  `Create a signal file with the given name. The signal file will be stored as ~/.my-context/signals/<name>.signal`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Get signal directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			signalsDir := filepath.Join(homeDir, ".my-context", "signals")

			// Create signal manager
			manager, err := signal.NewManager(signalsDir)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal create", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Create signal
			signal, err := manager.CreateSignal(name)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal create", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"name":       signal.Name,
					"created_at": signal.CreatedAt.Format(time.RFC3339),
					"path":       signal.Path,
				}
				jsonStr, err := output.FormatJSON("signal create", data)
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Signal '%s' created\n", name)
			}

			return nil
		},
	}

	return cmd
}

func newSignalListCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all signal files",
		Long:  `List all existing signal files with their creation timestamps.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get signal directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			signalsDir := filepath.Join(homeDir, ".my-context", "signals")

			// Create signal manager
			manager, err := signal.NewManager(signalsDir)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal list", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// List signals
			signals, err := manager.ListSignals()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal list", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				jsonStr, err := output.FormatJSON("signal list", map[string]interface{}{"signals": signals})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				if len(signals) == 0 {
					fmt.Println("No signals found")
					return nil
				}

				fmt.Println("Signals:")
				for _, sig := range signals {
					fmt.Printf("  %s (created: %s)\n", sig.Name, sig.CreatedAt)
				}
			}

			return nil
		},
	}

	return cmd
}

func newSignalWaitCmd(jsonOutput *bool) *cobra.Command {
	var timeout string

	cmd := &cobra.Command{
		Use:   "wait <name>",
		Short: "Wait for a signal file to appear",
		Long: `Wait for a signal file with the given name to appear. Blocks until the signal is created or timeout expires.

Timeout can be specified as a duration (e.g., "30s", "5m", "1h") or "infinite" to wait indefinitely.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Parse timeout
			var timeoutDuration time.Duration
			if timeout == "infinite" {
				timeoutDuration = 0 // Will be handled as no timeout in manager
			} else {
				var err error
				timeoutDuration, err = time.ParseDuration(timeout)
				if err != nil {
					return fmt.Errorf("invalid timeout duration '%s': %w", timeout, err)
				}
			}

			// Get signal directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			signalsDir := filepath.Join(homeDir, ".my-context", "signals")

			// Create signal manager
			manager, err := signal.NewManager(signalsDir)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal wait", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// If signal already exists, return immediately
			if manager.SignalExists(name) {
				if *jsonOutput {
					jsonStr, err := output.FormatJSON("signal wait", map[string]interface{}{
						"name":    name,
						"existed": true,
					})
					if err != nil {
						return err
					}
					fmt.Print(jsonStr)
				} else {
					fmt.Printf("Signal '%s' already exists\n", name)
				}
				return nil
			}

			// Set a reasonable timeout if infinite was requested
			if timeoutDuration == 0 {
				timeoutDuration = 24 * time.Hour // Default to 24 hours for "infinite"
			}

			// Wait for signal
			err = manager.WaitForSignal(name, timeoutDuration)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal wait", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				jsonStr, err := output.FormatJSON("signal wait", map[string]interface{}{
					"name":    name,
					"existed": false,
				})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Signal '%s' detected\n", name)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&timeout, "timeout", "5m", "Timeout duration (e.g., '30s', '5m', '1h') or 'infinite'")

	return cmd
}

func newSignalClearCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear <name>",
		Short: "Remove a signal file",
		Long:  `Remove the signal file with the given name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Get signal directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			signalsDir := filepath.Join(homeDir, ".my-context", "signals")

			// Create signal manager
			manager, err := signal.NewManager(signalsDir)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal clear", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Clear signal
			err = manager.ClearSignal(name)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("signal clear", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				jsonStr, err := output.FormatJSON("signal clear", map[string]interface{}{"name": name})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Signal '%s' cleared\n", name)
			}

			return nil
		},
	}

	return cmd
}
