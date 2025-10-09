package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// FormatExport generates markdown export for a context
func FormatExport(ctx *models.Context, notes []Note, files []File, touches []Touch, exportTime time.Time) string {
	var sb strings.Builder

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
			sb.WriteString(fmt.Sprintf("- **%s** %s\n", formatLocalTime(note.Timestamp), note.Content))
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
			sb.WriteString(fmt.Sprintf("- **%s** %s\n", formatLocalTime(file.Timestamp), file.Path))
		}
		sb.WriteString(fmt.Sprintf("\nTotal: %d files\n\n", len(files)))
	}

	sb.WriteString("---\n\n")

	// Activity section
	sb.WriteString("## Activity\n\n")
	if len(touches) == 0 {
		sb.WriteString("(none)\n\n")
	} else {
		for _, touch := range touches {
			sb.WriteString(fmt.Sprintf("- **%s** Touch\n", formatLocalTime(touch.Timestamp)))
		}
		sb.WriteString(fmt.Sprintf("\nTotal: %d touches\n\n", len(touches)))
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

// Note represents a timestamped note entry
type Note struct {
	Timestamp time.Time
	Content   string
}

// File represents a file association
type File struct {
	Timestamp time.Time
	Path      string
}

// Touch represents a touch event
type Touch struct {
	Timestamp time.Time
}

// FormatExportJSON generates JSON export for a context
func FormatExportJSON(ctx *models.Context, notes []Note, files []File, touches []Touch, exportTime time.Time) map[string]interface{} {
	export := map[string]interface{}{
		"context": map[string]interface{}{
			"name":             ctx.Name,
			"start_time":       ctx.StartTime.Format(time.RFC3339),
			"duration_seconds": int(ctx.Duration().Seconds()),
			"is_archived":      ctx.IsArchived,
		},
		"notes":  formatNotesJSON(notes),
		"files":  formatFilesJSON(files),
		"touches": formatTouchesJSON(touches),
		"export_metadata": map[string]interface{}{
			"exported_at": exportTime.Format(time.RFC3339),
			"version":     "v2.0.0",
		},
	}

	if ctx.EndTime != nil {
		export["context"].(map[string]interface{})["end_time"] = ctx.EndTime.Format(time.RFC3339)
	} else {
		export["context"].(map[string]interface{})["end_time"] = nil
	}

	return export
}

func formatNotesJSON(notes []Note) []map[string]interface{} {
	result := make([]map[string]interface{}, len(notes))
	for i, note := range notes {
		result[i] = map[string]interface{}{
			"timestamp": note.Timestamp.Format(time.RFC3339),
			"content":   note.Content,
		}
	}
	return result
}

func formatFilesJSON(files []File) []map[string]interface{} {
	result := make([]map[string]interface{}, len(files))
	for i, file := range files {
		result[i] = map[string]interface{}{
			"timestamp": file.Timestamp.Format(time.RFC3339),
			"path":      file.Path,
		}
	}
	return result
}

func formatTouchesJSON(touches []Touch) []map[string]interface{} {
	result := make([]map[string]interface{}, len(touches))
	for i, touch := range touches {
		result[i] = map[string]interface{}{
			"timestamp": touch.Timestamp.Format(time.RFC3339),
		}
	}
	return result
}
