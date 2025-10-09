package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/spf13/cobra"
)

func NewArchiveCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "archive [context-name]",
		Aliases: []string{"a"},
		Short:   "Archive a completed context",
		Long: `Archive a context to hide it from default list views while preserving all data.

Archived contexts are hidden from default 'list' output but can be viewed with 'list --archived'.
The context must be stopped before archiving.
All notes, files, and activity are preserved.

Examples:
  my-context archive "ps-cli: Phase 1"
  my-context archive "Completed Work"
  my-context a "Old Project"`,
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

			// Check if already archived
			if ctx.IsArchived {
				return fmt.Errorf("context %q is already archived", contextName)
			}

			// Prevent archiving active context
			activeCtx, err := core.GetActiveContext()
			if err == nil && activeCtx != nil && activeCtx.ActiveContext != nil && *activeCtx.ActiveContext == ctx.Name {
				return fmt.Errorf("cannot archive active context %q - stop it first with 'my-context stop'", contextName)
			}

			// Archive the context
			if err := core.ArchiveContext(contextName); err != nil {
				return fmt.Errorf("failed to archive context: %w", err)
			}

			fmt.Printf("Archived context: %s\n", contextName)
			return nil
		},
	}

	return cmd
}
