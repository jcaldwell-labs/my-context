package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(jsonOutput *bool) *cobra.Command {
	var (
		force bool
	)

	cmd := &cobra.Command{
		Use:     "delete [context-name]",
		Aliases: []string{"d"},
		Short:   "Delete a context permanently",
		Long: `Delete a context and all its data permanently.

This removes the entire context directory from ~/.my-context/.
The context must be stopped before deletion.
Transition history in transitions.log is preserved.

Examples:
  my-context delete "Test Context"
  my-context delete "ps-cli: Phase 1" --force
  my-context d "Old Work"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate: need context name
			if len(args) == 0 {
				return fmt.Errorf("context name required")
			}

			contextName := args[0]

			// Check if context exists
			ctx, _, _, _, err := core.GetContext(contextName)
			if err != nil {
				return fmt.Errorf("context not found: %s", contextName)
			}

			// Prevent deleting active context
			activeCtx, err := core.GetActiveContext()
			if err == nil && activeCtx != nil && activeCtx.ActiveContext != nil && *activeCtx.ActiveContext == ctx.Name {
				return fmt.Errorf("cannot delete active context %q - stop it first with 'my-context stop'", contextName)
			}

			// Confirmation prompt (unless --force)
			if !force {
				fmt.Printf("⚠️  WARNING: This will permanently delete context %q and all its data.\n", contextName)
				fmt.Printf("Are you sure? (yes/no): ")

				reader := bufio.NewReader(os.Stdin)
				response, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read confirmation: %w", err)
				}

				response = strings.TrimSpace(strings.ToLower(response))
				if response != "yes" && response != "y" {
					fmt.Println("Delete canceled.")
					return nil
				}
			}

			// Delete the context (passing force flag and confirmed=true after prompt)
			if err := core.DeleteContext(contextName, force, true); err != nil {
				return fmt.Errorf("failed to delete context: %w", err)
			}

			fmt.Printf("Deleted context: %s\n", contextName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}
