package commands

import (
	"errors"
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	"github.com/spf13/cobra"
)

func NewFileCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "file <path>",
		Aliases: []string{"f"},
		Short:   "Associate a file with the active context",
		Long:    `Associate a file path with the currently active context. Paths are normalized to absolute POSIX format.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			// Get active context
			state, err := core.GetActiveContext()
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("file", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			if !state.HasActiveContext() {
				errMsg := "No active context. Start a context with: my-context start <name>"
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("file", 1, errMsg)
					fmt.Print(jsonStr)
					return nil
				}
				return errors.New(errMsg)
			}

			// Add the file
			file, err := core.AddFile(filePath)
			if err != nil {
				if *jsonOutput {
					jsonStr, _ := output.FormatJSONError("file", 2, err.Error())
					fmt.Print(jsonStr)
					return nil
				}
				return err
			}

			// Output
			if *jsonOutput {
				data := output.FileData{
					ContextName:   state.GetActiveContextName(),
					FileTimestamp: file.Timestamp,
					FilePath:      file.FilePath,
					OriginalPath:  filePath,
				}
				jsonStr, err := output.FormatJSON("file", map[string]interface{}{"data": data})
				if err != nil {
					return err
				}
				fmt.Print(jsonStr)
			} else {
				fmt.Printf("File associated with context: %s\n", state.GetActiveContextName())
				fmt.Printf("  %s\n", core.DenormalizePath(file.FilePath))
			}

			return nil
		},
	}

	return cmd
}

