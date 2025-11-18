package commands

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewLinkCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link <child> <parent>",
		Short: "Link a child context to a parent context",
		Long: `Create a parent-child relationship between two contexts.
This is useful for organizing contexts hierarchically (e.g., sprint -> task -> subtask).

Examples:
  my-context link "API bugfix" "Sprint 3"
  my-context link "Unit tests" "API bugfix"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			childName := args[0]
			parentName := args[1]

			// Verify both contexts exist
			if _, err := core.LoadContext(childName); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("link", 1, fmt.Sprintf("child context %q not found", childName))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("child context %q not found", childName)
			}

			if _, err := core.LoadContext(parentName); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("link", 1, fmt.Sprintf("parent context %q not found", parentName))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("parent context %q not found", parentName)
			}

			// Prevent circular dependencies
			if childName == parentName {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("link", 1, "cannot link a context to itself")
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("cannot link a context to itself")
			}

			// Set parent relationship
			if err := core.SetParent(childName, parentName); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("link", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"child":  childName,
					"parent": parentName,
				}
				jsonStr, err := output.FormatJSON("link", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("✓ Linked \"%s\" → \"%s\"\n", childName, parentName)
			}

			return nil
		},
	}

	return cmd
}

func NewUnlinkCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "unlink <context>",
		Short: "Remove parent link from a context",
		Long: `Remove the parent-child relationship from a context.

Examples:
  my-context unlink "API bugfix"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			// Verify context exists
			if _, err := core.LoadContext(contextName); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("unlink", 1, fmt.Sprintf("context %q not found", contextName))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("context %q not found", contextName)
			}

			// Remove parent relationship
			if err := core.ClearParent(contextName); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("unlink", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"context": contextName,
				}
				jsonStr, err := output.FormatJSON("unlink", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("✓ Unlinked \"%s\" from its parent\n", contextName)
			}

			return nil
		},
	}
}
