package commands

import (
	"fmt"

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
