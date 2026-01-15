package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewTemplateCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Aliases: []string{"tmpl"},
		Short:   "Manage context templates",
		Long: `Manage context templates with predefined settings.

Templates allow you to start contexts with predefined labels and initial notes.

Subcommands:
  list   - List all available templates
  show   - Show details of a specific template
  create - Create a new template interactively`,
	}

	cmd.AddCommand(newTemplateListCmd(jsonOutput))
	cmd.AddCommand(newTemplateShowCmd(jsonOutput))
	cmd.AddCommand(newTemplateCreateCmd(jsonOutput))

	return cmd
}

func newTemplateListCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all available templates",
		Long:    `List all available context templates with their descriptions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			templates, err := core.ListTemplates()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("template_list", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"templates": templates,
				}
				jsonStr, err := output.FormatJSON("template_list", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				if len(templates) == 0 {
					fmt.Println("No templates found.")
					fmt.Println("\nCreate a template with: my-context template create <name>")
					return nil
				}

				fmt.Println("Available templates:")
				for _, t := range templates {
					fmt.Printf("\n  %s\n", t.Name)
					if t.Description != "" {
						fmt.Printf("    %s\n", t.Description)
					}
					if len(t.Labels) > 0 {
						fmt.Printf("    Labels: %s\n", strings.Join(t.Labels, ", "))
					}
					if len(t.InitialNotes) > 0 {
						fmt.Printf("    Initial notes: %d\n", len(t.InitialNotes))
					}
				}
			}

			return nil
		},
	}
}

func newTemplateShowCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show details of a specific template",
		Long:  `Display the full details of a template including all labels and initial notes.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]

			template, err := core.GetTemplate(templateName)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("template_show", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"template": template,
				}
				jsonStr, err := output.FormatJSON("template_show", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("Template: %s\n", template.Name)
				if template.Description != "" {
					fmt.Printf("Description: %s\n", template.Description)
				}

				if len(template.Labels) > 0 {
					fmt.Printf("\nLabels:\n")
					for _, label := range template.Labels {
						fmt.Printf("  - %s\n", label)
					}
				}

				if len(template.InitialNotes) > 0 {
					fmt.Printf("\nInitial notes:\n")
					for _, note := range template.InitialNotes {
						fmt.Printf("  - %s\n", note)
					}
				}
			}

			return nil
		},
	}
}

func newTemplateCreateCmd(jsonOutput *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new template interactively",
		Long: `Create a new context template by providing labels and initial notes interactively.

The template will be saved to MY_CONTEXT_HOME/templates/ as a YAML file.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]

			// Check if template already exists
			if _, err := core.GetTemplate(templateName); err == nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("template_create", 1, fmt.Sprintf("template %q already exists", templateName))
					fmt.Print(jsonStr)
					return nil
				}
				return fmt.Errorf("template %q already exists", templateName)
			}

			// Interactive prompts
			reader := bufio.NewReader(os.Stdin)

			// Description
			fmt.Print("Description (optional): ")
			description, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read description: %w", err)
			}
			description = strings.TrimSpace(description)

			// Labels
			fmt.Print("Labels (comma-separated, optional): ")
			labelsInput, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read labels: %w", err)
			}
			labelsInput = strings.TrimSpace(labelsInput)

			var labels []string
			if labelsInput != "" {
				labels = strings.Split(strings.ReplaceAll(labelsInput, " ", ""), ",")
			}

			// Initial notes
			fmt.Println("Initial notes (one per line, empty line to finish):")
			var initialNotes []string
			for {
				note, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read note: %w", err)
				}
				note = strings.TrimSpace(note)
				if note == "" {
					break
				}
				initialNotes = append(initialNotes, note)
			}

			// Create template
			template := &models.Template{
				Name:         templateName,
				Description:  description,
				Labels:       labels,
				InitialNotes: initialNotes,
			}

			// Validate
			if err := template.Validate(); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("template_create", 1, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Save template
			if err := core.SaveTemplate(template); err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("template_create", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := map[string]interface{}{
					"template": template,
				}
				jsonStr, err := output.FormatJSON("template_create", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("\n✓ Template %q created successfully\n", templateName)
				fmt.Printf("Location: %s\n", core.GetTemplatesDir())
			}

			return nil
		},
	}
}
