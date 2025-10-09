package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// FormatExportMarkdown generates a markdown document from context data
func FormatExportMarkdown(ctx *models.Context, notes []models.Note, files []models.FileAssociation, touchCount int) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# Context: %s\n\n", ctx.Name))

	// Metadata section
	sb.WriteString(fmt.Sprintf("**Started**: %s\n", ctx.StartTime.Format(time.RFC3339)))

	if ctx.EndTime != nil {
		sb.WriteString(fmt.Sprintf("**Ended**: %s\n", ctx.EndTime.Format(time.RFC3339)))
		duration := ctx.Duration()
		sb.WriteString(fmt.Sprintf("**Duration**: %s\n", formatDuration(duration)))
	} else {
		sb.WriteString("**Ended**: Active\n")
		sb.WriteString(fmt.Sprintf("**Duration**: %s (ongoing)\n", formatDuration(ctx.Duration())))
	}

	if ctx.IsArchived {
		sb.WriteString("**Status**: Archived\n")
	}

	sb.WriteString("\n---\n\n")

	// Notes section
	sb.WriteString("## Notes\n\n")
	if len(notes) == 0 {
		sb.WriteString("*No notes recorded*\n\n")
	} else {
		for _, note := range notes {
			timestamp := note.Timestamp.Format("15:04")
			sb.WriteString(fmt.Sprintf("- `%s` %s\n", timestamp, note.TextContent))
		}
		sb.WriteString("\n")
	}

	// Associated Files section
	sb.WriteString("## Associated Files\n\n")
	if len(files) == 0 {
		sb.WriteString("*No files associated*\n\n")
	} else {
		for _, file := range files {
			timestamp := file.Timestamp.Format("2006-01-02 15:04")
			sb.WriteString(fmt.Sprintf("- `%s` (Added: %s)\n", file.FilePath, timestamp))
		}
		sb.WriteString("\n")
	}

	// Activity section
	sb.WriteString("## Activity\n\n")
	sb.WriteString(fmt.Sprintf("- %d touch events\n\n", touchCount))

	// Footer
	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("*Exported: %s*\n", time.Now().Format(time.RFC3339)))

	return sb.String()
}

// formatDuration converts a duration to human-readable format
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dm", minutes)
}
