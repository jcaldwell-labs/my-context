package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// FormatExportMarkdown generates markdown export for a context
func FormatExportMarkdown(ctx *models.Context, notes []models.Note, files []models.FileAssociation, touchCount int) string {
	var sb strings.Builder

	exportTime := time.Now()

	// Header
	sb.WriteString(fmt.Sprintf("# Context: %s\n\n", ctx.Name))

	// Metadata
	sb.WriteString(fmt.Sprintf("**Started**: %s\n", formatLocalTime(ctx.StartTime)))

	if ctx.EndTime != nil {
		sb.WriteString(fmt.Sprintf("**Ended**: %s\n", formatLocalTime(*ctx.EndTime)))
	} else {
		sb.WriteString("**Ended**: Currently Active\n")
	}

	sb.WriteString(fmt.Sprintf("**Duration**: %s\n\n", formatDuration(ctx.Duration())))
	sb.WriteString(fmt.Sprintf("**Exported**: %s\n\n", formatLocalTime(exportTime)))

	if ctx.IsArchived {
		sb.WriteString("*This context is archived*\n\n")
	}

	sb.WriteString("---\n\n")

	// Notes section
	sb.WriteString("## Notes\n\n")
	if len(notes) == 0 {
		sb.WriteString("(none)\n\n")
	} else {
		for _, note := range notes {
			sb.WriteString(fmt.Sprintf("- **%s** %s\n", formatLocalTime(note.Timestamp), note.TextContent))
		}
		sb.WriteString(fmt.Sprintf("\nTotal: %d notes\n\n", len(notes)))
	}

	sb.WriteString("---\n\n")

	// Files section
	sb.WriteString("## Files\n\n")
	if len(files) == 0 {
		sb.WriteString("(none)\n\n")
	} else {
		for _, file := range files {
			sb.WriteString(fmt.Sprintf("- **%s** %s\n", formatLocalTime(file.Timestamp), file.FilePath))
		}
		sb.WriteString(fmt.Sprintf("\nTotal: %d files\n\n", len(files)))
	}

	sb.WriteString("---\n\n")

	// Activity section
	sb.WriteString("## Activity\n\n")
	if touchCount == 0 {
		sb.WriteString("(none)\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("Total: %d touches\n\n", touchCount))
	}

	sb.WriteString("---\n\n")

	// Footer
	sb.WriteString("*Exported from my-context v2.0.0*\n")

	return sb.String()
}

// formatLocalTime converts UTC time to local timezone and formats human-readable
func formatLocalTime(t time.Time) string {
	local := t.Local()
	return local.Format("January 2, 2006 at 3:04 PM MST")
}

// formatDuration formats a duration in human-readable form
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
