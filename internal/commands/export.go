package commands

import (
	"fmt"
	"os"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/spf13/cobra"
)

func NewExportCmd(jsonOutput *bool) *cobra.Command {
	var (
		exportToPath string
		exportAll    bool
		exportForce  bool
		exportAsJSON bool
	)

	cmd := &cobra.Command{
		Use:     "export [context-name]",
		Aliases: []string{"e"},
		Short:   "Export context data to markdown file",
		Long: `Export a context's notes, files, and activity to a markdown file for sharing.

Examples:
  my-context export "ps-cli: Phase 1"
  my-context export "Phase 1" --to reports/phase-1.md
  my-context export --all --to exports/
  my-context e "Phase 1"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate: need context name or --all flag
			if len(args) == 0 && !exportAll {
				return fmt.Errorf("context name required or use --all to export all contexts")
			}

			// Check if using database backend
			if core.IsUsingDatabase() {
				return exportWithDatabaseBackend(args, exportAll, exportToPath, exportAsJSON)
			}

			// File-based backend (existing code)
			// Export all contexts
			if exportAll {
				outputDir := exportToPath
				if outputDir == "" {
					outputDir = "."
				}

				exportedPaths, err := core.ExportAllContexts(outputDir, exportAsJSON)
				if err != nil {
					return fmt.Errorf("export failed: %w", err)
				}

				format := "markdown"
				if exportAsJSON {
					format = "JSON"
				}
				fmt.Printf("Exporting %d contexts to %s as %s...\n", len(exportedPaths), outputDir, format)
				for _, path := range exportedPaths {
					fmt.Printf("  ✓ %s\n", path)
				}
				fmt.Printf("Exported %d contexts to %s\n", len(exportedPaths), outputDir)
				return nil
			}

			// Export single context
			contextName := args[0]

			outputPath, err := core.ExportContext(contextName, exportToPath, exportAsJSON)
			if err != nil {
				return err
			}

			format := "markdown"
			if exportAsJSON {
				format = "JSON"
			}
			fmt.Printf("Exported context %q to %s (%s)\n", contextName, outputPath, format)
			return nil
		},
	}

	cmd.Flags().StringVar(&exportToPath, "to", "", "Output file path (default: ./{context_name}.md)")
	cmd.Flags().BoolVar(&exportAll, "all", false, "Export all contexts to separate files")
	cmd.Flags().BoolVar(&exportForce, "force", false, "Overwrite existing files without confirmation")
	cmd.Flags().BoolVar(&exportAsJSON, "json", false, "Output as JSON instead of markdown")

	return cmd
}

// exportWithDatabaseBackend exports context from PostgreSQL backend
func exportWithDatabaseBackend(args []string, exportAll bool, exportToPath string, asJSON bool) error {
	backend, err := core.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get backend: %w", err)
	}
	defer backend.Close()

	if exportAll {
		return fmt.Errorf("--all not yet implemented for database backend")
	}

	// Export single context
	contextName := args[0]

	// Get context from database
	ctx, err := backend.GetContext(contextName)
	if err != nil {
		return fmt.Errorf("context not found: %w", err)
	}

	// Get notes and files
	notes, _ := backend.GetNotes(contextName)
	files, _ := backend.GetFiles(contextName)

	// Determine output path
	outputPath := exportToPath
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s.md", contextName)
	}

	// Create markdown content
	content := fmt.Sprintf("# Context: %s\n\n", ctx.Name)
	content += fmt.Sprintf("**Status**: %s\n", ctx.Status)
	content += fmt.Sprintf("**Started**: %s\n", ctx.StartTime.Format("2006-01-02 15:04:05"))
	if ctx.EndTime != nil {
		content += fmt.Sprintf("**Ended**: %s\n", ctx.EndTime.Format("2006-01-02 15:04:05"))
		duration := ctx.EndTime.Sub(ctx.StartTime)
		content += fmt.Sprintf("**Duration**: %s\n", duration.String())
	}
	content += "\n---\n\n"

	// Add metadata if present
	if len(ctx.Metadata.Labels) > 0 || ctx.Metadata.Parent != "" || ctx.Metadata.CreatedBy != "" {
		content += "## Metadata\n\n"
		if ctx.Metadata.CreatedBy != "" {
			content += fmt.Sprintf("- **Created By**: %s\n", ctx.Metadata.CreatedBy)
		}
		if ctx.Metadata.Parent != "" {
			content += fmt.Sprintf("- **Parent**: %s\n", ctx.Metadata.Parent)
		}
		if len(ctx.Metadata.Labels) > 0 {
			content += fmt.Sprintf("- **Labels**: %v\n", ctx.Metadata.Labels)
		}
		content += "\n"
	}

	// Add notes
	content += fmt.Sprintf("## Notes (%d)\n\n", len(notes))
	if len(notes) > 0 {
		for _, note := range notes {
			content += fmt.Sprintf("**[%s]**\n%s\n\n", note.Timestamp.Format("2006-01-02 15:04:05"), note.Content)
		}
	} else {
		content += "(none)\n\n"
	}

	// Add files
	content += fmt.Sprintf("## Files (%d)\n\n", len(files))
	if len(files) > 0 {
		for _, file := range files {
			content += fmt.Sprintf("- `%s` (added: %s)\n", file.Path, file.Timestamp.Format("2006-01-02 15:04:05"))
		}
	} else {
		content += "(none)\n"
	}

	// Write to file
	err = os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	fmt.Printf("Exported context %q to %s\n", contextName, outputPath)
	return nil
}
