package commands

import (
	"fmt"
	"os"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewWhichCmd(jsonOutput *bool) *cobra.Command {
	var shortOutput bool

	cmd := &cobra.Command{
		Use:     "which",
		Aliases: []string{"where", "home"},
		Short:   "Show context home location",
		Long: `Display the active MY_CONTEXT_HOME location and context count.

This command helps troubleshoot "I don't see my context" issues by showing
which context home is currently active.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			homePath := core.GetContextHome()
			homeDisplay := core.GetContextHomeDisplay()
			contextCount := core.GetContextCount()

			// Get active context name
			state, _ := core.GetActiveContext()
			activeContext := "(none)"
			if state != nil && state.HasActiveContext() {
				activeContext = state.GetActiveContextName()
			}

			// Get MY_CONTEXT_HOME env value
			envValue := os.Getenv("MY_CONTEXT_HOME")
			envSet := envValue != ""

			// Check if using database backend and get partition info
			usingDatabase := core.IsUsingDatabase()
			var partitionName, schemaName string
			if usingDatabase {
				partitionName = core.ExtractPartition()
				schemaName = core.GetPartitionSchema()
				if partitionName == "" {
					partitionName = "(default)"
				}
			}

			// Output
			switch {
			case *jsonOutput:
				data := map[string]interface{}{
					"context_home":         homePath,
					"context_home_display": homeDisplay,
					"context_count":        contextCount,
					"active_context":       activeContext,
					"env_set":              envSet,
					"env_value":            envValue,
					"using_database":       usingDatabase,
				}
				if usingDatabase {
					data["partition"] = partitionName
					data["schema"] = schemaName
				}
				jsonStr, err := output.FormatJSON("which", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			case shortOutput:
				// Short output: just the path
				fmt.Println(homeDisplay)
			default:
				// Full output
				fmt.Printf("Context Home: %s\n\n", homeDisplay)

				fmt.Println("Environment:")
				if envSet {
					fmt.Printf("  MY_CONTEXT_HOME=%s\n", envValue)
					fmt.Println("  (Set in current shell)")
				} else {
					fmt.Println("  MY_CONTEXT_HOME=(not set)")
					fmt.Println("  (Using default location)")
				}

				fmt.Println()
				fmt.Println("Details:")
				if usingDatabase {
					fmt.Printf("  Backend: PostgreSQL database\n")
					fmt.Printf("  Partition: %s\n", partitionName)
					fmt.Printf("  Schema: %s\n", schemaName)
				} else {
					fmt.Printf("  Backend: File-based storage\n")
					fmt.Printf("  Location: %s\n", homePath)
					fmt.Printf("  State file: %s/state.json\n", homePath)
				}
				fmt.Printf("  Contexts: %d\n", contextCount)
				fmt.Printf("  Active: %s\n", activeContext)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&shortOutput, "short", "s", false, "Short output (path only)")

	return cmd
}
