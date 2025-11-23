package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewTreeCmd(jsonOutput *bool) *cobra.Command {
	var showAll bool

	cmd := &cobra.Command{
		Use:   "tree [context]",
		Short: "Show context hierarchy as a tree",
		Long: `Display the parent-child relationships of contexts as a tree structure.
If no context is specified, shows all root contexts (contexts without parents).

Examples:
  my-context tree              # Show all root contexts and their children
  my-context tree "Sprint 3"   # Show tree starting from specific context`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				// Show tree for specific context
				contextName := args[0]

				tree, err := core.GetContextTree(contextName)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("tree", 1, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				// Output
				if *jsonOutput {
					jsonStr, err := json.MarshalIndent(tree, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(jsonStr))
				} else {
					fmt.Printf("Context hierarchy for \"%s\":\n\n", contextName)
					printTree(tree, "", true)
				}
			} else {
				// Show all root contexts
				roots, err := core.GetRootContexts()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("tree", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				if len(roots) == 0 {
					if *jsonOutput {
						data := map[string]interface{}{
							"message": "No contexts found",
						}
						jsonStr, _ := output.FormatJSON("tree", map[string]interface{}{"data": data})
						fmt.Print(jsonStr)
					} else {
						fmt.Println("No contexts found")
					}
					return nil
				}

				// Build trees for all roots
				var trees []*core.ContextTreeNode
				for _, rootName := range roots {
					tree, err := core.GetContextTree(rootName)
					if err != nil {
						continue // Skip roots that can't be built
					}
					trees = append(trees, tree)
				}

				// Output
				if *jsonOutput {
					jsonStr, err := json.MarshalIndent(trees, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(jsonStr))
				} else {
					fmt.Println("Context hierarchy (all root contexts):")
					for i, tree := range trees {
						printTree(tree, "", true)
						if i < len(trees)-1 {
							fmt.Println()
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all contexts including orphans")

	return cmd
}

// printTree prints a tree structure with proper indentation and branch characters
func printTree(node *core.ContextTreeNode, prefix string, isLast bool) {
	// Print current node
	if prefix == "" {
		// Root node
		fmt.Printf("└─ %s\n", node.Name)
	} else {
		// Child node
		marker := "├─"
		if isLast {
			marker = "└─"
		}
		fmt.Printf("%s%s %s\n", prefix, marker, node.Name)
	}

	// Print children
	if len(node.Children) > 0 {
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1

			// Calculate new prefix for children
			var newPrefix string
			if prefix == "" {
				newPrefix = "   "
			} else {
				if isLast {
					newPrefix = prefix + "   "
				} else {
					newPrefix = prefix + "│  "
				}
			}

			printTree(child, newPrefix, isLastChild)
		}
	}
}

func NewUpCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:     "up [context]",
		Aliases: []string{"parent"},
		Short:   "Show or switch to parent context",
		Long: `Show the parent context of the current context, or switch to it.
If a context name is provided, shows the parent of that context.

Examples:
  my-context up              # Show parent of active context
  my-context up "Bug fix"    # Show parent of specific context`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var contextName string

			if len(args) == 1 {
				contextName = args[0]
			} else {
				// Get active context
				state, err := core.GetActiveContext()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("up", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				if !state.HasActiveContext() {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("up", 1, "no active context")
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("no active context")
				}

				contextName = state.GetActiveContextName()
			}

			// Get parent
			ctx, _, _, _, err := core.GetContextWithMetadata(contextName)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("up", 1, fmt.Sprintf("context %q not found", contextName))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("context %q not found", contextName)
			}

			if ctx.Metadata.Parent == "" {
				if *jsonOutput {
					data := map[string]interface{}{
						"context": contextName,
						"message": "no parent",
					}
					jsonStr, _ := output.FormatJSON("up", map[string]interface{}{"data": data})
					fmt.Print(jsonStr)
					return nil
				}
				fmt.Printf("\"%s\" has no parent context\n", contextName)
				return nil
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"context": contextName,
					"parent":  ctx.Metadata.Parent,
				}
				jsonStr, err := output.FormatJSON("up", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Parent of \"%s\": %s\n", contextName, ctx.Metadata.Parent)
			}

			return nil
		},
	}
}

func NewDownCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:     "down [context]",
		Aliases: []string{"children"},
		Short:   "List child contexts",
		Long: `List all child contexts of the current context.
If a context name is provided, lists children of that context.

Examples:
  my-context down              # List children of active context
  my-context down "Sprint 3"   # List children of specific context`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var contextName string

			if len(args) == 1 {
				contextName = args[0]
			} else {
				// Get active context
				state, err := core.GetActiveContext()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("down", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				if !state.HasActiveContext() {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("down", 1, "no active context")
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("no active context")
				}

				contextName = state.GetActiveContextName()
			}

			// Verify context exists
			if _, err := core.LoadContext(contextName); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("down", 1, fmt.Sprintf("context %q not found", contextName))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("context %q not found", contextName)
			}

			// Get children
			children, err := core.GetChildren(contextName)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("down", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"context":  contextName,
					"children": children,
				}
				jsonStr, err := output.FormatJSON("down", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				if len(children) == 0 {
					fmt.Printf("\"%s\" has no child contexts\n", contextName)
				} else {
					fmt.Printf("Children of \"%s\":\n", contextName)
					for _, child := range children {
						fmt.Printf("  • %s\n", child)
					}
				}
			}

			return nil
		},
	}
}
