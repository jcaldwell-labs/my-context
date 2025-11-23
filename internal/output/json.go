package output

import (
	"encoding/json"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// JSONResponse represents a standard JSON response structure
type JSONResponse struct {
	Command   string      `json:"command"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
	Error     *JSONError  `json:"error,omitempty"`
}

// JSONError represents an error in JSON format
type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ContextData represents context data for JSON output
type ContextData struct {
	Context interface{}               `json:"context"` // Can be *models.Context or *pkgmodels.ContextWithMetadata
	Notes   []*models.Note            `json:"notes,omitempty"`
	Files   []*models.FileAssociation `json:"files,omitempty"`
	Touches []*models.TouchEvent      `json:"touches,omitempty"`
}

// StartData represents start command output data
type StartData struct {
	ContextName             string  `json:"context_name"`
	OriginalName            string  `json:"original_name"`
	WasDuplicate            bool    `json:"was_duplicate"`
	PreviousContext         *string `json:"previous_context,omitempty"`
	PreviousDurationSeconds *int    `json:"previous_duration_seconds,omitempty"`
}

// StopData represents stop command output data
type StopData struct {
	ContextName     string    `json:"context_name"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	DurationSeconds int       `json:"duration_seconds"`
}

// NoteData represents note command output data
type NoteData struct {
	ContextName   string    `json:"context_name"`
	NoteTimestamp time.Time `json:"note_timestamp"`
	NoteText      string    `json:"note_text"`
}

// FileData represents file command output data
type FileData struct {
	ContextName   string    `json:"context_name"`
	FileTimestamp time.Time `json:"file_timestamp"`
	FilePath      string    `json:"file_path"`
	OriginalPath  string    `json:"original_path"`
}

// TouchData represents touch command output data
type TouchData struct {
	ContextName    string    `json:"context_name"`
	TouchTimestamp time.Time `json:"touch_timestamp"`
}

// ListData represents list command output data
type ListData struct {
	Contexts []*ContextSummary `json:"contexts"`
}

// ContextSummary represents a context summary for list output
type ContextSummary struct {
	Name            string     `json:"name"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	Status          string     `json:"status"`
	DurationSeconds int        `json:"duration_seconds"`
	NoteCount       int        `json:"note_count"`
	FileCount       int        `json:"file_count"`
	TouchCount      int        `json:"touch_count"`
}

// HistoryData represents history command output data
type HistoryData struct {
	Transitions []*models.ContextTransition `json:"transitions"`
}

// FormatJSON formats any data as JSON
func FormatJSON(command string, data interface{}) (string, error) {
	response := JSONResponse{
		Command:   command,
		Timestamp: time.Now(),
		Data:      data,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData) + "\n", nil
}

// FormatJSONError formats an error as JSON
func FormatJSONError(command string, code int, message string) (string, error) {
	response := JSONResponse{
		Command:   command,
		Timestamp: time.Now(),
		Error: &JSONError{
			Code:    code,
			Message: message,
		},
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData) + "\n", nil
}

// ExportData represents export data structure for JSON output
type ExportData struct {
	Name       string                    `json:"name"`
	StartTime  time.Time                 `json:"start_time"`
	EndTime    *time.Time                `json:"end_time,omitempty"`
	Status     string                    `json:"status"`
	IsArchived bool                      `json:"is_archived"`
	Duration   int                       `json:"duration_seconds"`
	Notes      []models.Note             `json:"notes"`
	Files      []models.FileAssociation  `json:"files"`
	TouchCount int                       `json:"touch_count"`
	ExportTime time.Time                 `json:"export_time"`
}

// FormatExportJSON formats context export data as JSON
func FormatExportJSON(ctx *models.Context, notes []models.Note, files []models.FileAssociation, touchCount int) (string, error) {
	exportData := ExportData{
		Name:       ctx.Name,
		StartTime:  ctx.StartTime,
		EndTime:    ctx.EndTime,
		Status:     ctx.Status,
		IsArchived: ctx.IsArchived,
		Duration:   int(ctx.Duration().Seconds()),
		Notes:      notes,
		Files:      files,
		TouchCount: touchCount,
		ExportTime: time.Now(),
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData) + "\n", nil
}
