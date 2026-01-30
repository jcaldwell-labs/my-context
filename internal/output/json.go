package output

import (
	"encoding/json"
	"reflect"
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
	ContextName       string    `json:"context_name"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	DurationSeconds   int       `json:"duration_seconds"`
	StoppedChildren   []string  `json:"stopped_children,omitempty"`
	ActiveChildrenCount int     `json:"active_children_count,omitempty"`
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
	Index           int        `json:"index"`
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
	Name       string                   `json:"name"`
	StartTime  time.Time                `json:"start_time"`
	EndTime    *time.Time               `json:"end_time,omitempty"`
	Status     string                   `json:"status"`
	IsArchived bool                     `json:"is_archived"`
	Duration   int                      `json:"duration_seconds"`
	Notes      []models.Note            `json:"notes"`
	Files      []models.FileAssociation `json:"files"`
	TouchCount int                      `json:"touch_count"`
	ExportTime time.Time                `json:"export_time"`
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

// StatsData represents stats command output data
type StatsData struct {
	Period         StatsPeriod        `json:"period"`
	TotalContexts  int                `json:"total_contexts"`
	TotalSeconds   int64              `json:"total_seconds"`
	AverageSeconds int64              `json:"average_seconds"`
	ByProject      []StatsProjectData `json:"by_project,omitempty"`
	LongestSession *StatsSessionData  `json:"longest_session,omitempty"`
	ActiveContext  *StatsSessionData  `json:"active_context,omitempty"`
}

// StatsPeriod represents the time period for stats
type StatsPeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Label string    `json:"label"`
}

// StatsProjectData represents per-project stats
type StatsProjectData struct {
	Project string  `json:"project"`
	Seconds int64   `json:"seconds"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

// StatsSessionData represents session info
type StatsSessionData struct {
	Name    string `json:"name"`
	Seconds int64  `json:"seconds"`
}

// TimePeriod interface for stats command
type TimePeriod interface {
	GetStart() time.Time
	GetEnd() time.Time
	GetLabel() string
}

// ProjectStats interface for stats command
type ProjectStats interface {
	GetProject() string
	GetSeconds() int64
	GetCount() int
	GetPercent() float64
}

// FormatStatsJSON formats stats data as JSON using reflection for type conversion
func FormatStatsJSON(period, byProject, longest, active any, totalContexts int, totalSeconds, avgSeconds int64) (string, error) {
	// Extract period using reflection
	statsPeriod := extractPeriod(period)

	// Convert project stats using reflection
	projectStats := extractProjectStats(byProject)

	// Convert session info using reflection
	longestSession := extractSession(longest)
	activeContext := extractSession(active)

	data := StatsData{
		Period:         statsPeriod,
		TotalContexts:  totalContexts,
		TotalSeconds:   totalSeconds,
		AverageSeconds: avgSeconds,
		ByProject:      projectStats,
		LongestSession: longestSession,
		ActiveContext:  activeContext,
	}

	return FormatJSON("stats", data)
}

// extractPeriod extracts period fields using reflection
func extractPeriod(v any) StatsPeriod {
	if v == nil {
		return StatsPeriod{}
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return StatsPeriod{}
	}

	var period StatsPeriod
	if f := val.FieldByName("Start"); f.IsValid() && f.Type() == reflect.TypeOf(time.Time{}) {
		period.Start = f.Interface().(time.Time)
	}
	if f := val.FieldByName("End"); f.IsValid() && f.Type() == reflect.TypeOf(time.Time{}) {
		period.End = f.Interface().(time.Time)
	}
	if f := val.FieldByName("Label"); f.IsValid() && f.Kind() == reflect.String {
		period.Label = f.String()
	}
	return period
}

// extractProjectStats converts a slice of project stats using reflection
func extractProjectStats(v any) []StatsProjectData {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Slice {
		return nil
	}

	result := make([]StatsProjectData, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		if item.Kind() != reflect.Struct {
			continue
		}

		var ps StatsProjectData
		if f := item.FieldByName("Project"); f.IsValid() && f.Kind() == reflect.String {
			ps.Project = f.String()
		}
		if f := item.FieldByName("Seconds"); f.IsValid() && f.Kind() == reflect.Int64 {
			ps.Seconds = f.Int()
		}
		if f := item.FieldByName("Count"); f.IsValid() && f.Kind() == reflect.Int {
			ps.Count = int(f.Int())
		}
		if f := item.FieldByName("Percent"); f.IsValid() && f.Kind() == reflect.Float64 {
			ps.Percent = f.Float()
		}
		result = append(result, ps)
	}
	return result
}

// extractSession extracts session info using reflection
func extractSession(v any) *StatsSessionData {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	var name string
	var seconds int64
	if f := val.FieldByName("Name"); f.IsValid() && f.Kind() == reflect.String {
		name = f.String()
	}
	if f := val.FieldByName("Seconds"); f.IsValid() && f.Kind() == reflect.Int64 {
		seconds = f.Int()
	}

	if name == "" {
		return nil
	}
	return &StatsSessionData{Name: name, Seconds: seconds}
}
