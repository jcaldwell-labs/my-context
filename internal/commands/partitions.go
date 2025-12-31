package commands

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

// PartitionInfo holds information about a partition
type PartitionInfo struct {
	Schema        string    `json:"schema"`
	ContextCount  int       `json:"context_count"`
	ActiveContext string    `json:"active_context,omitempty"`
	LastActivity  time.Time `json:"last_activity,omitempty"`
}

func NewPartitionsCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "partitions",
		Aliases: []string{"p"},
		Short:   "List all database partitions",
		Long: `List all database partitions (schemas) with context counts and active contexts.

This command shows all partitions that contain my-context data, along with:
- Number of contexts in each partition
- Currently active context (if any)
- Last activity timestamp

Example:
  my-context partitions
  MY_CONTEXT_HOME=db my-context partitions --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Only works with database backend
			if !core.IsUsingDatabase() {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("partitions", 1, "partitions command only works with database backend (MY_CONTEXT_HOME=db)")
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("partitions command only works with database backend\nHint: Set MY_CONTEXT_HOME=db or MY_CONTEXT_HOME=db:partition-name")
			}

			backend, err := core.GetBackend()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("partitions", 2, fmt.Sprintf("failed to get backend: %v", err))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("failed to get backend: %w", err)
			}
			defer backend.Close()

			// Get list of all schemas that have my-context tables
			partitions, err := listPartitions(backend)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("partitions", 2, fmt.Sprintf("failed to list partitions: %v", err))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("failed to list partitions: %w", err)
			}

			if *jsonOutput {
				jsonStr, err := output.FormatJSON("partitions", partitions)
				if err != nil {
					jsonStr, _ := output.FormatJSONError("partitions", 2, fmt.Sprintf("failed to format output: %v", err))
					fmt.Print(jsonStr)
					return nil
				}
				fmt.Print(jsonStr)
				return nil
			}

			// Human-readable output
			if len(partitions) == 0 {
				fmt.Println("No partitions found")
				return nil
			}

			fmt.Printf("Found %d partition(s):\n\n", len(partitions))
			for i, p := range partitions {
				fmt.Printf("%d. Partition: %s\n", i+1, p.Schema)
				fmt.Printf("   Contexts: %d\n", p.ContextCount)
				if p.ActiveContext != "" {
					fmt.Printf("   Active:   %s\n", p.ActiveContext)
				} else {
					fmt.Printf("   Active:   (none)\n")
				}
				if !p.LastActivity.IsZero() {
					fmt.Printf("   Last Activity: %s\n", p.LastActivity.Format("2006-01-02 15:04:05"))
				}
				fmt.Println()
			}

			// Show usage hint for current partition
			currentSchema := core.GetPartitionSchema()
			fmt.Printf("Current partition: %s (MY_CONTEXT_HOME=%s)\n", currentSchema, getCurrentMYContextHome())

			return nil
		},
	}

	return cmd
}

// getCurrentMYContextHome returns the current MY_CONTEXT_HOME value for display
func getCurrentMYContextHome() string {
	partition := core.ExtractPartition()
	if partition == "" {
		return "db"
	}
	return fmt.Sprintf("db:%s", partition)
}

// listPartitions queries the database for all schemas with my-context tables
func listPartitions(backend interface{}) ([]PartitionInfo, error) {
	// Type assert to get access to the underlying database connection
	// We need to query across schemas, which requires direct database access
	type dbProvider interface {
		GetDB() *sql.DB
	}

	provider, ok := backend.(dbProvider)
	if !ok {
		return nil, fmt.Errorf("backend does not provide database access")
	}

	db := provider.GetDB()

	// Query all schemas that have a 'contexts' table
	query := `
		SELECT
			schemaname,
			COUNT(*) as table_count
		FROM pg_tables
		WHERE tablename = 'contexts'
		  AND schemaname NOT IN ('pg_catalog', 'information_schema')
		GROUP BY schemaname
		ORDER BY schemaname
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query partitions: %w", err)
	}
	defer rows.Close()

	var partitions []PartitionInfo
	for rows.Next() {
		var schema string
		var tableCount int

		if err := rows.Scan(&schema, &tableCount); err != nil {
			return nil, fmt.Errorf("failed to scan partition row: %w", err)
		}

		// For each schema, get context count, active context, and last activity
		info := PartitionInfo{
			Schema: schema,
		}

		// Get context count
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s.contexts", schema)
		if err := db.QueryRow(countQuery).Scan(&info.ContextCount); err != nil {
			// Skip if we can't read this schema
			continue
		}

		// Get active context
		activeQuery := fmt.Sprintf(`
			SELECT value->>'name'
			FROM %s.state
			WHERE key = 'active_context'
		`, schema)
		var activeContext sql.NullString
		if err := db.QueryRow(activeQuery).Scan(&activeContext); err == nil && activeContext.Valid {
			info.ActiveContext = activeContext.String
		}

		// Get last activity (most recent context update or touch)
		activityQuery := fmt.Sprintf(`
			SELECT GREATEST(
				COALESCE(MAX(updated_at), '1970-01-01'::timestamptz),
				COALESCE(MAX(last_touch_at), '1970-01-01'::timestamptz)
			) as last_activity
			FROM %s.contexts
		`, schema)
		var lastActivity sql.NullTime
		if err := db.QueryRow(activityQuery).Scan(&lastActivity); err == nil && lastActivity.Valid {
			info.LastActivity = lastActivity.Time
		}

		partitions = append(partitions, info)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating partitions: %w", err)
	}

	return partitions, nil
}
