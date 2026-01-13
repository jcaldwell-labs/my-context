package mcp

import (
	"context"
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// StartContextInput defines the input for start_context tool.
type StartContextInput struct {
	Name    string `json:"name" jsonschema_description:"Name for the new context"`
	Project string `json:"project,omitempty" jsonschema_description:"Optional project name (creates 'project: name' format)"`
}

// StartContextOutput defines the output for start_context tool.
type StartContextOutput struct {
	ContextName  string `json:"context_name"`
	Message      string `json:"message"`
	WasDuplicate bool   `json:"was_duplicate,omitempty"`
}

// handleStartContext starts a new work context.
func handleStartContext(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartContextInput,
) (*mcp.CallToolResult, StartContextOutput, error) {
	if input.Name == "" {
		return nil, StartContextOutput{}, fmt.Errorf("context name is required")
	}

	// Build context name with optional project prefix
	contextName := input.Name
	if input.Project != "" {
		contextName = fmt.Sprintf("%s: %s", input.Project, input.Name)
	}

	// Create the context
	createdCtx, actualName, err := core.CreateContext(contextName)
	if err != nil {
		return nil, StartContextOutput{}, fmt.Errorf("failed to start context: %w", err)
	}

	wasDuplicate := actualName != contextName

	return nil, StartContextOutput{
		ContextName:  createdCtx.Name,
		Message:      fmt.Sprintf("Started context: %s", createdCtx.Name),
		WasDuplicate: wasDuplicate,
	}, nil
}

// StopContextInput defines the input for stop_context tool.
type StopContextInput struct {
	// No required input - stops the active context
}

// StopContextOutput defines the output for stop_context tool.
type StopContextOutput struct {
	ContextName string `json:"context_name,omitempty"`
	Message     string `json:"message"`
	Duration    string `json:"duration,omitempty"`
}

// handleStopContext stops the currently active context.
func handleStopContext(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopContextInput,
) (*mcp.CallToolResult, StopContextOutput, error) {
	// Get active context first
	state, err := core.GetActiveContext()
	if err != nil {
		return nil, StopContextOutput{}, fmt.Errorf("failed to get active context: %w", err)
	}

	if !state.HasActiveContext() {
		return nil, StopContextOutput{
			Message: "No active context to stop",
		}, nil
	}

	contextName := state.GetActiveContextName()

	// Stop the context
	stoppedCtx, err := core.StopContext()
	if err != nil {
		return nil, StopContextOutput{}, fmt.Errorf("failed to stop context: %w", err)
	}

	duration := ""
	if stoppedCtx != nil {
		duration = stoppedCtx.Duration().String()
	}

	return nil, StopContextOutput{
		ContextName: contextName,
		Message:     fmt.Sprintf("Stopped context: %s", contextName),
		Duration:    duration,
	}, nil
}

// AddNoteInput defines the input for add_note tool.
type AddNoteInput struct {
	Text string `json:"text" jsonschema_description:"The note text to add"`
}

// AddNoteOutput defines the output for add_note tool.
type AddNoteOutput struct {
	ContextName string `json:"context_name"`
	Message     string `json:"message"`
	NoteCount   int    `json:"note_count"`
}

// handleAddNote adds a timestamped note to the active context.
func handleAddNote(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddNoteInput,
) (*mcp.CallToolResult, AddNoteOutput, error) {
	if input.Text == "" {
		return nil, AddNoteOutput{}, fmt.Errorf("note text is required")
	}

	// Get active context
	state, err := core.GetActiveContext()
	if err != nil {
		return nil, AddNoteOutput{}, fmt.Errorf("failed to get active context: %w", err)
	}

	if !state.HasActiveContext() {
		return nil, AddNoteOutput{}, fmt.Errorf("no active context - start one first with start_context")
	}

	contextName := state.GetActiveContextName()

	// Add the note
	if _, err := core.AddNote(input.Text); err != nil {
		return nil, AddNoteOutput{}, fmt.Errorf("failed to add note: %w", err)
	}

	// Get note count
	noteCount, _ := core.GetNoteCount(contextName)

	return nil, AddNoteOutput{
		ContextName: contextName,
		Message:     "Note added successfully",
		NoteCount:   noteCount,
	}, nil
}

// GetActiveContextInput defines the input for get_active_context tool.
type GetActiveContextInput struct {
	// No required input
}

// GetActiveContextOutput defines the output for get_active_context tool.
type GetActiveContextOutput struct {
	HasActive   bool     `json:"has_active"`
	ContextName string   `json:"context_name,omitempty"`
	Status      string   `json:"status,omitempty"`
	Duration    string   `json:"duration,omitempty"`
	NoteCount   int      `json:"note_count,omitempty"`
	FileCount   int      `json:"file_count,omitempty"`
	Notes       []string `json:"notes,omitempty"`
}

// handleGetActiveContext returns information about the active context.
func handleGetActiveContext(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetActiveContextInput,
) (*mcp.CallToolResult, GetActiveContextOutput, error) {
	state, err := core.GetActiveContext()
	if err != nil {
		return nil, GetActiveContextOutput{}, fmt.Errorf("failed to get active context: %w", err)
	}

	if !state.HasActiveContext() {
		return nil, GetActiveContextOutput{
			HasActive: false,
		}, nil
	}

	contextName := state.GetActiveContextName()

	// Get full context data
	ctxData, notes, files, _, err := core.GetContextWithMetadata(contextName)
	if err != nil {
		return nil, GetActiveContextOutput{}, fmt.Errorf("failed to get context data: %w", err)
	}

	// Extract note texts (last 10)
	var noteTexts []string
	startIdx := 0
	if len(notes) > 10 {
		startIdx = len(notes) - 10
	}
	for i := startIdx; i < len(notes); i++ {
		noteTexts = append(noteTexts, notes[i].TextContent)
	}

	return nil, GetActiveContextOutput{
		HasActive:   true,
		ContextName: ctxData.Name,
		Status:      ctxData.Status,
		Duration:    ctxData.Duration().String(),
		NoteCount:   len(notes),
		FileCount:   len(files),
		Notes:       noteTexts,
	}, nil
}

// ListContextsInput defines the input for list_contexts tool.
type ListContextsInput struct {
	Limit      int    `json:"limit,omitempty" jsonschema_description:"Maximum number of contexts to return (default: 10)"`
	Search     string `json:"search,omitempty" jsonschema_description:"Filter contexts by name substring"`
	Project    string `json:"project,omitempty" jsonschema_description:"Filter contexts by project name"`
	ActiveOnly bool   `json:"active_only,omitempty" jsonschema_description:"Only show active contexts"`
}

// ContextSummary provides a summary of a context.
type ContextSummary struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
	IsActive bool   `json:"is_active"`
}

// ListContextsOutput defines the output for list_contexts tool.
type ListContextsOutput struct {
	Contexts   []ContextSummary `json:"contexts"`
	TotalCount int              `json:"total_count"`
	ActiveName string           `json:"active_name,omitempty"`
	HasMore    bool             `json:"has_more"`
}

// handleListContexts lists all contexts with optional filtering.
func handleListContexts(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ListContextsInput,
) (*mcp.CallToolResult, ListContextsOutput, error) {
	// Set default limit
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	// Build filter
	filter := core.ContextFilter{
		Limit:      limit + 1, // Get one extra to detect "has more"
		Search:     input.Search,
		Project:    input.Project,
		ActiveOnly: input.ActiveOnly,
	}

	// Get filtered contexts
	contexts, err := core.ListContextsFiltered(filter)
	if err != nil {
		return nil, ListContextsOutput{}, fmt.Errorf("failed to list contexts: %w", err)
	}

	// Get active context name
	state, _ := core.GetActiveContext()
	activeName := ""
	if state != nil && state.HasActiveContext() {
		activeName = state.GetActiveContextName()
	}

	// Check if there are more results
	hasMore := len(contexts) > limit
	if hasMore {
		contexts = contexts[:limit]
	}

	// Build summaries
	summaries := make([]ContextSummary, len(contexts))
	for i, c := range contexts {
		summaries[i] = ContextSummary{
			Name:     c.Name,
			Status:   c.Status,
			Duration: c.Duration().String(),
			IsActive: c.Name == activeName,
		}
	}

	return nil, ListContextsOutput{
		Contexts:   summaries,
		TotalCount: len(summaries),
		ActiveName: activeName,
		HasMore:    hasMore,
	}, nil
}

// AddFileInput defines the input for add_file tool.
type AddFileInput struct {
	Path string `json:"path" jsonschema_description:"Absolute or relative path to the file to associate with the context"`
}

// AddFileOutput defines the output for add_file tool.
type AddFileOutput struct {
	ContextName string `json:"context_name"`
	FilePath    string `json:"file_path"`
	Message     string `json:"message"`
	FileCount   int    `json:"file_count"`
}

// handleAddFile associates a file with the active context.
func handleAddFile(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddFileInput,
) (*mcp.CallToolResult, AddFileOutput, error) {
	if input.Path == "" {
		return nil, AddFileOutput{}, fmt.Errorf("file path is required")
	}

	// Get active context
	state, err := core.GetActiveContext()
	if err != nil {
		return nil, AddFileOutput{}, fmt.Errorf("failed to get active context: %w", err)
	}

	if !state.HasActiveContext() {
		return nil, AddFileOutput{}, fmt.Errorf("no active context - start one first with start_context")
	}

	contextName := state.GetActiveContextName()

	// Add the file
	fileAssoc, err := core.AddFile(input.Path)
	if err != nil {
		return nil, AddFileOutput{}, fmt.Errorf("failed to add file: %w", err)
	}

	// Get file count
	_, _, files, _, _ := core.GetContextWithMetadata(contextName)
	fileCount := len(files)

	return nil, AddFileOutput{
		ContextName: contextName,
		FilePath:    fileAssoc.FilePath,
		Message:     "File associated successfully",
		FileCount:   fileCount,
	}, nil
}

// ListFilesInput defines the input for list_files tool.
type ListFilesInput struct {
	ContextName string `json:"context_name,omitempty" jsonschema_description:"Context name to list files for (defaults to active context)"`
	Limit       int    `json:"limit,omitempty" jsonschema_description:"Maximum number of files to return (default: 20)"`
}

// FileSummary provides a summary of a file association.
type FileSummary struct {
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}

// ListFilesOutput defines the output for list_files tool.
type ListFilesOutput struct {
	ContextName string        `json:"context_name"`
	Files       []FileSummary `json:"files"`
	TotalCount  int           `json:"total_count"`
	HasMore     bool          `json:"has_more"`
}

// handleListFiles lists files associated with a context.
func handleListFiles(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ListFilesInput,
) (*mcp.CallToolResult, ListFilesOutput, error) {
	contextName := input.ContextName

	// Default to active context if not specified
	if contextName == "" {
		state, err := core.GetActiveContext()
		if err != nil {
			return nil, ListFilesOutput{}, fmt.Errorf("failed to get active context: %w", err)
		}

		if !state.HasActiveContext() {
			return nil, ListFilesOutput{}, fmt.Errorf("no active context and no context_name specified")
		}

		contextName = state.GetActiveContextName()
	}

	// Get context data with files
	_, _, files, _, err := core.GetContextWithMetadata(contextName)
	if err != nil {
		return nil, ListFilesOutput{}, fmt.Errorf("failed to get context files: %w", err)
	}

	// Set default limit
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	// Check if there are more results
	hasMore := len(files) > limit
	totalCount := len(files)

	// Apply limit
	if len(files) > limit {
		files = files[:limit]
	}

	// Build file summaries
	summaries := make([]FileSummary, len(files))
	for i, f := range files {
		summaries[i] = FileSummary{
			Path:      f.FilePath,
			Timestamp: f.Timestamp.Format("2006-01-02 15:04:05"),
		}
	}

	return nil, ListFilesOutput{
		ContextName: contextName,
		Files:       summaries,
		TotalCount:  totalCount,
		HasMore:     hasMore,
	}, nil
}
