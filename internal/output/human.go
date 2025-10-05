package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// FormatContext formats a context with all its data for human-readable output
func FormatContext(ctx *models.Context, notes []*models.Note, files []*models.FileAssociation, touches []*models.TouchEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Context: %s\n", ctx.Name))
	sb.WriteString(fmt.Sprintf("Status: %s\n", ctx.Status))

	// Format start time with relative time
	duration := ctx.Duration()
	sb.WriteString(fmt.Sprintf("Started: %s (%s ago)\n",
		ctx.StartTime.Format("2006-01-02 15:04:05"),
		FormatDuration(duration)))

	// Notes section
	sb.WriteString(fmt.Sprintf("\nNotes (%d):\n", len(notes)))
	if len(notes) == 0 {
		sb.WriteString("  (none)\n")
	} else {
		for _, note := range notes {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n",
				note.Timestamp.Format("15:04"),
				note.TextContent))
		}
	}

	// Files section
	sb.WriteString(fmt.Sprintf("\nFiles (%d):\n", len(files)))
	if len(files) == 0 {
		sb.WriteString("  (none)\n")
	} else {
		for _, file := range files {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n",
				file.Timestamp.Format("15:04"),
				file.FilePath))
		}
	}

	// Activity section
	sb.WriteString(fmt.Sprintf("\nActivity: %d touches", len(touches)))
	if len(touches) > 0 {
		lastTouch := touches[len(touches)-1]
		sb.WriteString(fmt.Sprintf(" (last: %s)", lastTouch.Timestamp.Format("15:04")))
	}
	sb.WriteString("\n")

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
		sb.WriteString(fmt.Sprintf("    Started: %s (%s ago)\n",
			ctx.StartTime.Format("2006-01-02 15:04"),
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

	for _, t := range transitions {
		timestamp := t.Timestamp.Format("2006-01-02 15:04")

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
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
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

