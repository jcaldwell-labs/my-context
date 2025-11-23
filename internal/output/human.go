package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// PrintContextHomeHeader prints the context home location with separator
func PrintContextHomeHeader(homeDisplay string, contextCount int) {
	fmt.Printf("Context Home: %s", homeDisplay)
	if contextCount > 0 {
		fmt.Printf(" (%d contexts)", contextCount)
	}
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

// getTimestampFormat returns the timestamp format based on MC_TIMESTAMP_FORMAT env var
func getTimestampFormat() string {
	format := os.Getenv("MC_TIMESTAMP_FORMAT")
	switch format {
	case "short":
		return "15:04" // HH:MM
	case "medium":
		return "15:04:05" // HH:MM:SS
	case "long":
		return "2006-01-02 15:04:05" // Full datetime
	case "iso", "":
		return time.RFC3339 // ISO8601 with timezone (default)
	default:
		// Allow custom Go time format strings
		return format
	}
}

// FormatContext formats a context with all its data for human-readable output
// contextInfo holds extracted context information
type contextInfo struct {
	name        string
	status      string
	startTime   time.Time
	duration    time.Duration
	hasMetadata bool
	metadata    interface{}
}

// extractContextInfo extracts common fields from different context types
func extractContextInfo(ctx interface{}) (contextInfo, error) {
	var info contextInfo

	switch c := ctx.(type) {
	case *pkgmodels.ContextWithMetadata:
		info.name = c.Name
		info.status = c.Status
		info.startTime = c.StartTime
		info.duration = c.Duration()
		info.hasMetadata = true
		info.metadata = c.Metadata
	case *models.Context:
		info.name = c.Name
		info.status = c.Status
		info.startTime = c.StartTime
		info.duration = c.Duration()
		info.hasMetadata = false
	default:
		return info, fmt.Errorf("unknown context type")
	}
	return info, nil
}

// formatMetadataSection formats the metadata section if available
func formatMetadataSection(hasMetadata bool, metadata interface{}) string {
	if !hasMetadata {
		return ""
	}

	meta, ok := metadata.(pkgmodels.ContextMetadata)
	if !ok || (meta.CreatedBy == "" && meta.Parent == "" && len(meta.Labels) == 0) {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\nMetadata:\n")
	if meta.CreatedBy != "" {
		sb.WriteString(fmt.Sprintf("  Created by: %s\n", meta.CreatedBy))
	}
	if meta.Parent != "" {
		sb.WriteString(fmt.Sprintf("  Parent: %s\n", meta.Parent))
	}
	if len(meta.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("  Labels: %s\n", strings.Join(meta.Labels, ", ")))
	}
	return sb.String()
}

// formatNotesSection formats the notes section
func formatNotesSection(notes []*models.Note) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nNotes (%d):\n", len(notes)))

	if len(notes) == 0 {
		sb.WriteString("  (none)\n")
		return sb.String()
	}

	timestampFormat := getTimestampFormat()
	for _, note := range notes {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n",
			note.Timestamp.Format(timestampFormat),
			note.TextContent))
	}
	return sb.String()
}

// formatFilesSection formats the files section
func formatFilesSection(files []*models.FileAssociation) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nFiles (%d):\n", len(files)))

	if len(files) == 0 {
		sb.WriteString("  (none)\n")
		return sb.String()
	}

	timestampFormat := getTimestampFormat()
	for _, file := range files {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n",
			file.Timestamp.Format(timestampFormat),
			file.FilePath))
	}
	return sb.String()
}

// formatActivitySection formats the activity section
func formatActivitySection(touches []*models.TouchEvent) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nActivity: %d touches", len(touches)))

	if len(touches) > 0 {
		timestampFormat := getTimestampFormat()
		lastTouch := touches[len(touches)-1]
		sb.WriteString(fmt.Sprintf(" (last: %s)", lastTouch.Timestamp.Format(timestampFormat)))
	}
	sb.WriteString("\n")
	return sb.String()
}

func FormatContext(ctx interface{}, notes []*models.Note, files []*models.FileAssociation, touches []*models.TouchEvent) string {
	var sb strings.Builder

	// Extract context information
	info, err := extractContextInfo(ctx)
	if err != nil {
		sb.WriteString("Error: Unknown context type\n")
		return sb.String()
	}

	// Format header
	sb.WriteString(fmt.Sprintf("Context: %s\n", info.name))
	sb.WriteString(fmt.Sprintf("Status: %s\n", info.status))
	sb.WriteString(fmt.Sprintf("Started: %s (%s ago)\n",
		info.startTime.Format("2006-01-02 15:04:05"),
		FormatDuration(info.duration)))

	// Format sections
	sb.WriteString(formatMetadataSection(info.hasMetadata, info.metadata))
	sb.WriteString(formatNotesSection(notes))
	sb.WriteString(formatFilesSection(files))
	sb.WriteString(formatActivitySection(touches))

	return sb.String()
}

// FormatContextList formats a list of contexts for human-readable output
func FormatContextList(contexts []*models.Context, activeContextName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Contexts (%d):\n\n", len(contexts)))

	if len(contexts) == 0 {
		sb.WriteString("No contexts found\n")
		sb.WriteString("Start one with: my-context start <name>\n")
		return sb.String()
	}

	for _, ctx := range contexts {
		// Active indicator
		indicator := "○"
		if ctx.Name == activeContextName {
			indicator = "●"
		}

		// Status line
		statusLine := fmt.Sprintf("  %s %s (%s)\n", indicator, ctx.Name, ctx.Status)
		sb.WriteString(statusLine)

		// Start time line
		duration := ctx.Duration()
		timestampFormat := getTimestampFormat()
		// For list view, always show date if not already in format
		if !strings.Contains(timestampFormat, "2006") {
			timestampFormat = "2006-01-02 " + timestampFormat
		}
		sb.WriteString(fmt.Sprintf("    Started: %s (%s ago)\n",
			ctx.StartTime.Format(timestampFormat),
			FormatDuration(duration)))

		// Duration line if stopped
		if ctx.Status == "stopped" && ctx.EndTime != nil {
			actualDuration := ctx.EndTime.Sub(ctx.StartTime)
			sb.WriteString(fmt.Sprintf("    Duration: %s\n", FormatDuration(actualDuration)))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatTransitionHistory formats transition history for human-readable output
func FormatTransitionHistory(transitions []*models.ContextTransition) string {
	var sb strings.Builder

	sb.WriteString("Context History:\n\n")

	if len(transitions) == 0 {
		sb.WriteString("No transitions recorded\n")
		return sb.String()
	}

	timestampFormat := getTimestampFormat()
	// For history, always show date if not already in format
	if !strings.Contains(timestampFormat, "2006") {
		timestampFormat = "2006-01-02 " + timestampFormat
	}

	for _, t := range transitions {
		timestamp := t.Timestamp.Format(timestampFormat)

		switch t.TransitionType {
		case models.TransitionStart:
			sb.WriteString(fmt.Sprintf("  %s │ START     │ %s\n",
				timestamp, *t.NewContext))
		case models.TransitionStop:
			sb.WriteString(fmt.Sprintf("  %s │ STOP      │ %s\n",
				timestamp, *t.PreviousContext))
		case models.TransitionSwitch:
			sb.WriteString(fmt.Sprintf("  %s │ SWITCH    │ %s → %s\n",
				timestamp, *t.PreviousContext, *t.NewContext))
		}
	}

	return sb.String()
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	default:
		days := int(d.Hours()) / 24
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}
}

// FormatSimpleMessage formats a simple message
func FormatSimpleMessage(message string) string {
	return message + "\n"
}

// FormatError formats an error message
func FormatError(err error) string {
	return "Error: " + err.Error() + "\n"
}
