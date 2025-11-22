package commands

import (
	"fmt"
	"strings"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewTagCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tag",
		Aliases: []string{"t"},
		Short:   "Manage context tags/labels",
		Long: `Manage tags/labels for contexts. Tags help organize and filter contexts.

Subcommands:
  add     - Add tags to a context
  remove  - Remove tags from a context
  list    - List all tags or tags for a specific context`,
	}

	cmd.AddCommand(newTagAddCmd(jsonOutput))
	cmd.AddCommand(newTagRemoveCmd(jsonOutput))
	cmd.AddCommand(newTagListCmd(jsonOutput))

	return cmd
}

func newTagAddCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:     "add <context> <tag1> [tag2] ...",
		Aliases: []string{"a"},
		Short:   "Add one or more tags to a context",
		Long: `Add one or more tags to a context. Tags cannot contain whitespace.

Examples:
  my-context tag add "Bug fix" bug urgent
  my-context tag add "Sprint 3" feature frontend`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]
			tags := args[1:]

			// Validate tags
			for _, tag := range tags {
				if strings.ContainsAny(tag, " \t\n\r") {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("tag_add", 1, fmt.Sprintf("tag '%s' cannot contain whitespace", tag))
						fmt.Print(jsonStr)
						return nil
					}
					return fmt.Errorf("tag '%s' cannot contain whitespace", tag)
				}
			}

			// Add tags
			added, err := core.AddTags(contextName, tags)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("tag_add", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"context":    contextName,
					"added_tags": added,
				}
				jsonStr, err := output.FormatJSON("tag_add", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				if len(added) == 0 {
					fmt.Printf("No new tags added to \"%s\" (all tags already exist)\n", contextName)
				} else {
					fmt.Printf("✓ Added %d tag(s) to \"%s\": %s\n", len(added), contextName, strings.Join(added, ", "))
				}
			}

			return nil
		},
	}
}

func newTagRemoveCmd(jsonOutput *bool) *cobra.Command {
	var removeAll bool

	cmd := &cobra.Command{
		Use:     "remove <context> <tag1> [tag2] ...",
		Aliases: []string{"rm", "r"},
		Short:   "Remove one or more tags from a context",
		Long: `Remove one or more tags from a context.

Examples:
  my-context tag remove "Bug fix" urgent
  my-context tag rm "Sprint 3" feature
  my-context tag remove "Sprint 3" --all`,
		Args: func(cmd *cobra.Command, args []string) error {
			if removeAll {
				return cobra.ExactArgs(1)(cmd, args)
			}
			return cobra.MinimumNArgs(2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			contextName := args[0]

			var tags []string
			if removeAll {
				// Get all tags for the context
				allTags, err := core.GetContextTags(contextName)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("tag_remove", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}
				tags = allTags
			} else {
				tags = args[1:]
			}

			// Remove tags
			removed, err := core.RemoveTags(contextName, tags)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("tag_remove", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"context":      contextName,
					"removed_tags": removed,
				}
				jsonStr, err := output.FormatJSON("tag_remove", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				if len(removed) == 0 {
					fmt.Printf("No tags removed from \"%s\" (tags not found)\n", contextName)
				} else {
					fmt.Printf("✓ Removed %d tag(s) from \"%s\": %s\n", len(removed), contextName, strings.Join(removed, ", "))
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&removeAll, "all", false, "Remove all tags from the context")

	return cmd
}

func newTagListCmd(jsonOutput *bool) *cobra.Command {
	var showCount bool

	cmd := &cobra.Command{
		Use:     "list [context]",
		Aliases: []string{"ls", "l"},
		Short:   "List all tags or tags for a specific context",
		Long: `List all tags in use across all contexts, or list tags for a specific context.

Examples:
  my-context tag list              # List all tags with usage counts
  my-context tag list "Bug fix"    # List tags for specific context`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				// List tags for specific context
				contextName := args[0]
				tags, err := core.GetContextTags(contextName)
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("tag_list", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				// Output
				if *jsonOutput {
					data := map[string]interface{}{
						"context": contextName,
						"tags":    tags,
					}
					jsonStr, err := output.FormatJSON("tag_list", map[string]interface{}{"data": data})
					if err != nil {
						return err
					}
					fmt.Print(jsonStr)
				} else {
					if len(tags) == 0 {
						fmt.Printf("Context \"%s\" has no tags\n", contextName)
					} else {
						fmt.Printf("Tags for \"%s\": %s\n", contextName, strings.Join(tags, ", "))
					}
				}
			} else {
				// List all tags across all contexts
				tagCounts, err := core.GetAllTags()
				if err != nil {
					if *jsonOutput {
						jsonStr, _ := output.FormatJSONError("tag_list", 2, err.Error())
						fmt.Print(jsonStr)
						return nil
					}
					return err
				}

				// Output
				if *jsonOutput {
					data := map[string]interface{}{
						"tags": tagCounts,
					}
					jsonStr, err := output.FormatJSON("tag_list", map[string]interface{}{"data": data})
					if err != nil {
						return err
					}
					fmt.Print(jsonStr)
				} else {
					if len(tagCounts) == 0 {
						fmt.Println("No tags in use")
					} else {
						fmt.Println("All tags in use:")
						for tag, count := range tagCounts {
							if showCount {
								fmt.Printf("  %s (%d)\n", tag, count)
							} else {
								fmt.Printf("  %s\n", tag)
							}
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&showCount, "count", "c", true, "Show usage count for each tag")

	return cmd
}
