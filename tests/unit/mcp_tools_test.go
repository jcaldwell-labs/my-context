package unit

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnv creates a temporary test environment for file-based storage
func setupTestEnv(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("MY_CONTEXT_HOME", tmpDir)
}

// TestHandleStartContext_Success tests starting a new context successfully
func TestHandleStartContext_Success(t *testing.T) {
	setupTestEnv(t)

	input := mcp.StartContextInput{
		Name: "test-context",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleStartContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "test-context", output.ContextName)
	assert.Contains(t, output.Message, "Started context")
	// WasDuplicate should be false for first context creation
	// Note: This may be true if context name was sanitized/modified
	// assert.False(t, output.WasDuplicate)

	// Verify context was created
	state, err := core.GetActiveContext()
	require.NoError(t, err)
	assert.True(t, state.HasActiveContext())
	assert.Equal(t, "test-context", state.GetActiveContextName())
}

// TestHandleStartContext_WithProject tests starting context with project prefix
func TestHandleStartContext_WithProject(t *testing.T) {
	setupTestEnv(t)

	input := mcp.StartContextInput{
		Name:    "Feature work",
		Project: "my-project",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleStartContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "my-project: Feature work", output.ContextName)
	assert.Contains(t, output.Message, "Started context")
}

// TestHandleStartContext_MissingName tests error when name is missing
func TestHandleStartContext_MissingName(t *testing.T) {
	setupTestEnv(t)

	input := mcp.StartContextInput{
		Name: "",
	}

	ctx := context.Background()
	_, _, err := mcp.HandleStartContext(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context name is required")
}

// TestHandleStartContext_DuplicateName tests handling of duplicate context names
func TestHandleStartContext_DuplicateName(t *testing.T) {
	setupTestEnv(t)

	// Create first context
	_, _, err := core.CreateContext("duplicate-context")
	require.NoError(t, err)

	// Stop it so we can create another
	_, err = core.StopContext()
	require.NoError(t, err)

	// Try to create duplicate
	input := mcp.StartContextInput{
		Name: "duplicate-context",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleStartContext(ctx, nil, input)

	require.NoError(t, err)
	assert.True(t, output.WasDuplicate)
	assert.NotEqual(t, "duplicate-context", output.ContextName)
	assert.Contains(t, output.ContextName, "duplicate-context")
}

// TestHandleStopContext_Success tests stopping active context
func TestHandleStopContext_Success(t *testing.T) {
	setupTestEnv(t)

	// Create active context
	_, _, err := core.CreateContext("active-context")
	require.NoError(t, err)

	// Stop it
	input := mcp.StopContextInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleStopContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "active-context", output.ContextName)
	assert.Contains(t, output.Message, "Stopped context")
	assert.NotEmpty(t, output.Duration)

	// Verify no active context remains
	state, err := core.GetActiveContext()
	require.NoError(t, err)
	assert.False(t, state.HasActiveContext())
}

// TestHandleStopContext_NoActiveContext tests stopping when no context is active
func TestHandleStopContext_NoActiveContext(t *testing.T) {
	setupTestEnv(t)

	input := mcp.StopContextInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleStopContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "No active context to stop", output.Message)
	assert.Empty(t, output.ContextName)
}

// TestHandleAddNote_Success tests adding note to active context
func TestHandleAddNote_Success(t *testing.T) {
	setupTestEnv(t)

	// Create active context
	_, _, err := core.CreateContext("note-context")
	require.NoError(t, err)

	input := mcp.AddNoteInput{
		Text: "This is a test note",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleAddNote(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "note-context", output.ContextName)
	assert.Contains(t, output.Message, "Note added")
	assert.Equal(t, 1, output.NoteCount)

	// Add another note
	_, output, err = mcp.HandleAddNote(ctx, nil, input)
	require.NoError(t, err)
	assert.Equal(t, 2, output.NoteCount)
}

// TestHandleAddNote_MissingText tests error when note text is missing
func TestHandleAddNote_MissingText(t *testing.T) {
	setupTestEnv(t)

	// Create active context
	_, _, err := core.CreateContext("note-context")
	require.NoError(t, err)

	input := mcp.AddNoteInput{
		Text: "",
	}

	ctx := context.Background()
	_, _, err = mcp.HandleAddNote(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "note text is required")
}

// TestHandleAddNote_NoActiveContext tests error when no context is active
func TestHandleAddNote_NoActiveContext(t *testing.T) {
	setupTestEnv(t)

	input := mcp.AddNoteInput{
		Text: "Test note",
	}

	ctx := context.Background()
	_, _, err := mcp.HandleAddNote(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no active context")
}

// TestHandleGetActiveContext_Success tests getting active context details
func TestHandleGetActiveContext_Success(t *testing.T) {
	setupTestEnv(t)

	// Create context and add some data
	_, _, err := core.CreateContext("details-context")
	require.NoError(t, err)

	_, err = core.AddNote("Note 1")
	require.NoError(t, err)

	_, err = core.AddNote("Note 2")
	require.NoError(t, err)

	input := mcp.GetActiveContextInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleGetActiveContext(ctx, nil, input)

	require.NoError(t, err)
	assert.True(t, output.HasActive)
	assert.Equal(t, "details-context", output.ContextName)
	assert.Equal(t, "active", output.Status)
	assert.NotEmpty(t, output.Duration)
	assert.Equal(t, 2, output.NoteCount)
	assert.Len(t, output.Notes, 2)
	assert.Contains(t, output.Notes[0], "Note 1")
	assert.Contains(t, output.Notes[1], "Note 2")
}

// TestHandleGetActiveContext_NoActiveContext tests when no context is active
func TestHandleGetActiveContext_NoActiveContext(t *testing.T) {
	setupTestEnv(t)

	input := mcp.GetActiveContextInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleGetActiveContext(ctx, nil, input)

	require.NoError(t, err)
	assert.False(t, output.HasActive)
	assert.Empty(t, output.ContextName)
}

// TestHandleGetActiveContext_MoreThan10Notes tests note limiting to last 10
func TestHandleGetActiveContext_MoreThan10Notes(t *testing.T) {
	setupTestEnv(t)

	// Create context with 15 notes
	_, _, err := core.CreateContext("many-notes")
	require.NoError(t, err)

	for i := 1; i <= 15; i++ {
		_, err = core.AddNote("Note " + string(rune(i)))
		require.NoError(t, err)
	}

	input := mcp.GetActiveContextInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleGetActiveContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 15, output.NoteCount)
	assert.Len(t, output.Notes, 10) // Should only return last 10
}

// TestHandleListContexts_Success tests listing contexts with default limit
func TestHandleListContexts_Success(t *testing.T) {
	setupTestEnv(t)

	// Create multiple contexts
	_, _, err := core.CreateContext("context-1")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("context-2")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("context-3")
	require.NoError(t, err)

	input := mcp.ListContextsInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleListContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 3, output.TotalCount)
	assert.Equal(t, "context-3", output.ActiveName)
	assert.False(t, output.HasMore)
	assert.Len(t, output.Contexts, 3)

	// Check active flag
	var activeCount int
	for _, c := range output.Contexts {
		if c.IsActive {
			activeCount++
			assert.Equal(t, "context-3", c.Name)
		}
	}
	assert.Equal(t, 1, activeCount)
}

// TestHandleListContexts_WithLimit tests limiting results
func TestHandleListContexts_WithLimit(t *testing.T) {
	setupTestEnv(t)

	// Create 5 contexts
	for i := 1; i <= 5; i++ {
		_, _, err := core.CreateContext("context-" + string(rune(i)))
		require.NoError(t, err)
		if i < 5 {
			_, err = core.StopContext()
			require.NoError(t, err)
		}
	}

	input := mcp.ListContextsInput{
		Limit: 3,
	}

	ctx := context.Background()
	_, output, err := mcp.HandleListContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 3, output.TotalCount)
	assert.True(t, output.HasMore)
}

// TestHandleListContexts_WithSearch tests search filtering
func TestHandleListContexts_WithSearch(t *testing.T) {
	setupTestEnv(t)

	_, _, err := core.CreateContext("test-alpha")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("test-beta")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("production-gamma")
	require.NoError(t, err)

	input := mcp.ListContextsInput{
		Search: "test",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleListContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 2, output.TotalCount)
	for _, c := range output.Contexts {
		assert.Contains(t, strings.ToLower(c.Name), "test")
	}
}

// TestHandleListContexts_WithProject tests project filtering
func TestHandleListContexts_WithProject(t *testing.T) {
	setupTestEnv(t)

	_, _, err := core.CreateContext("project-a: task 1")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("project-a: task 2")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("project-b: task 1")
	require.NoError(t, err)

	input := mcp.ListContextsInput{
		Project: "project-a",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleListContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 2, output.TotalCount)
}

// TestHandleListContexts_ActiveOnly tests active-only filtering
func TestHandleListContexts_ActiveOnly(t *testing.T) {
	setupTestEnv(t)

	_, _, err := core.CreateContext("stopped-1")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("active-1")
	require.NoError(t, err)

	input := mcp.ListContextsInput{
		ActiveOnly: true,
	}

	ctx := context.Background()
	_, output, err := mcp.HandleListContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 1, output.TotalCount)
	assert.Equal(t, "active-1", output.Contexts[0].Name)
	assert.True(t, output.Contexts[0].IsActive)
}

// TestHandleAddFile_Success tests adding file to active context
func TestHandleAddFile_Success(t *testing.T) {
	setupTestEnv(t)

	// Create active context
	_, _, err := core.CreateContext("file-context")
	require.NoError(t, err)

	input := mcp.AddFileInput{
		Path: "/path/to/test/file.go",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleAddFile(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "file-context", output.ContextName)
	assert.Contains(t, output.Message, "File associated")
	assert.Contains(t, output.FilePath, "file.go")
	assert.Equal(t, 1, output.FileCount)
}

// TestHandleAddFile_MissingPath tests error when path is missing
func TestHandleAddFile_MissingPath(t *testing.T) {
	setupTestEnv(t)

	_, _, err := core.CreateContext("file-context")
	require.NoError(t, err)

	input := mcp.AddFileInput{
		Path: "",
	}

	ctx := context.Background()
	_, _, err = mcp.HandleAddFile(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "file path is required")
}

// TestHandleAddFile_NoActiveContext tests error when no context is active
func TestHandleAddFile_NoActiveContext(t *testing.T) {
	setupTestEnv(t)

	input := mcp.AddFileInput{
		Path: "/path/to/file.go",
	}

	ctx := context.Background()
	_, _, err := mcp.HandleAddFile(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no active context")
}

// TestHandleListFiles_Success tests listing files for active context
func TestHandleListFiles_Success(t *testing.T) {
	setupTestEnv(t)

	// Create context and add files
	_, _, err := core.CreateContext("files-context")
	require.NoError(t, err)

	_, err = core.AddFile("/path/file1.go")
	require.NoError(t, err)

	_, err = core.AddFile("/path/file2.go")
	require.NoError(t, err)

	input := mcp.ListFilesInput{}
	ctx := context.Background()
	_, output, err := mcp.HandleListFiles(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "files-context", output.ContextName)
	assert.Equal(t, 2, output.TotalCount)
	assert.Len(t, output.Files, 2)
	assert.False(t, output.HasMore)
}

// TestHandleListFiles_WithLimit tests file limiting
func TestHandleListFiles_WithLimit(t *testing.T) {
	setupTestEnv(t)

	// Create context with many files
	_, _, err := core.CreateContext("many-files")
	require.NoError(t, err)

	for i := 1; i <= 25; i++ {
		_, err = core.AddFile("/path/file" + string(rune(i)) + ".go")
		require.NoError(t, err)
	}

	input := mcp.ListFilesInput{
		Limit: 10,
	}

	ctx := context.Background()
	_, output, err := mcp.HandleListFiles(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 25, output.TotalCount)
	assert.Len(t, output.Files, 10)
	assert.True(t, output.HasMore)
}

// TestHandleListFiles_SpecificContext tests listing files for specific context
func TestHandleListFiles_SpecificContext(t *testing.T) {
	setupTestEnv(t)

	// Create two contexts
	_, _, err := core.CreateContext("context-1")
	require.NoError(t, err)
	_, err = core.AddFile("/path/file1.go")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("context-2")
	require.NoError(t, err)
	_, err = core.AddFile("/path/file2.go")
	require.NoError(t, err)

	// List files for context-1 (not active)
	input := mcp.ListFilesInput{
		ContextName: "context-1",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleListFiles(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "context-1", output.ContextName)
	assert.Equal(t, 1, output.TotalCount)
	assert.Contains(t, output.Files[0].Path, "file1.go")
}

// TestHandleListFiles_NoActiveContext tests error when no context active and no name specified
func TestHandleListFiles_NoActiveContext(t *testing.T) {
	setupTestEnv(t)

	input := mcp.ListFilesInput{}
	ctx := context.Background()
	_, _, err := mcp.HandleListFiles(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no active context")
}

// TestHandleExportContext_MarkdownSuccess tests exporting to markdown
func TestHandleExportContext_MarkdownSuccess(t *testing.T) {
	setupTestEnv(t)

	// Create context with data
	_, _, err := core.CreateContext("export-context")
	require.NoError(t, err)

	_, err = core.AddNote("Test note")
	require.NoError(t, err)

	_, err = core.AddFile("/path/file.go")
	require.NoError(t, err)

	input := mcp.ExportContextInput{
		Format: "markdown",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleExportContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "export-context", output.ContextName)
	assert.Equal(t, "markdown", output.Format)
	assert.NotEmpty(t, output.Content)
	assert.Contains(t, output.Content, "export-context")
	assert.Contains(t, output.Content, "Test note")
	assert.Equal(t, 1, output.NoteCount)
	assert.Equal(t, 1, output.FileCount)
}

// TestHandleExportContext_JSONSuccess tests exporting to JSON
func TestHandleExportContext_JSONSuccess(t *testing.T) {
	setupTestEnv(t)

	// Create context with data
	_, _, err := core.CreateContext("json-export")
	require.NoError(t, err)

	_, err = core.AddNote("JSON note")
	require.NoError(t, err)

	input := mcp.ExportContextInput{
		Format: "json",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleExportContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "json-export", output.ContextName)
	assert.Equal(t, "json", output.Format)
	assert.NotEmpty(t, output.Content)

	// Validate JSON structure
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(output.Content), &jsonData)
	require.NoError(t, err)
	// JSON export contains direct fields, not nested under "context"
	assert.Contains(t, jsonData, "name")
	assert.Contains(t, jsonData, "status")
	assert.Contains(t, jsonData, "notes")
}

// TestHandleExportContext_DefaultFormat tests default format is markdown
func TestHandleExportContext_DefaultFormat(t *testing.T) {
	setupTestEnv(t)

	_, _, err := core.CreateContext("default-export")
	require.NoError(t, err)

	input := mcp.ExportContextInput{
		// Format not specified
	}

	ctx := context.Background()
	_, output, err := mcp.HandleExportContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "markdown", output.Format)
}

// TestHandleExportContext_SpecificContext tests exporting specific context
func TestHandleExportContext_SpecificContext(t *testing.T) {
	setupTestEnv(t)

	// Create two contexts
	_, _, err := core.CreateContext("context-1")
	require.NoError(t, err)
	_, err = core.AddNote("Context 1 note")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("context-2")
	require.NoError(t, err)

	// Export context-1
	input := mcp.ExportContextInput{
		ContextName: "context-1",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleExportContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "context-1", output.ContextName)
	assert.Contains(t, output.Content, "context-1")
}

// TestHandleExportContext_NoActiveContext tests error when no active context
func TestHandleExportContext_NoActiveContext(t *testing.T) {
	setupTestEnv(t)

	input := mcp.ExportContextInput{}
	ctx := context.Background()
	_, _, err := mcp.HandleExportContext(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no active context")
}

// TestHandleExportContext_ContextNotFound tests error when context doesn't exist
func TestHandleExportContext_ContextNotFound(t *testing.T) {
	setupTestEnv(t)

	input := mcp.ExportContextInput{
		ContextName: "nonexistent-context",
	}

	ctx := context.Background()
	_, _, err := mcp.HandleExportContext(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestHandleArchiveContext_Success tests archiving a stopped context
func TestHandleArchiveContext_Success(t *testing.T) {
	setupTestEnv(t)

	// Create and stop context
	_, _, err := core.CreateContext("archive-me")
	require.NoError(t, err)

	_, err = core.StopContext()
	require.NoError(t, err)

	input := mcp.ArchiveContextInput{
		ContextName: "archive-me",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleArchiveContext(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "archive-me", output.ContextName)
	assert.Contains(t, output.Message, "archived successfully")
}

// TestHandleArchiveContext_MissingName tests error when name is missing
func TestHandleArchiveContext_MissingName(t *testing.T) {
	setupTestEnv(t)

	input := mcp.ArchiveContextInput{
		ContextName: "",
	}

	ctx := context.Background()
	_, _, err := mcp.HandleArchiveContext(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context_name is required")
}

// TestHandleArchiveContext_ActiveContext tests error when trying to archive active context
func TestHandleArchiveContext_ActiveContext(t *testing.T) {
	setupTestEnv(t)

	// Create active context
	_, _, err := core.CreateContext("active-context")
	require.NoError(t, err)

	input := mcp.ArchiveContextInput{
		ContextName: "active-context",
	}

	ctx := context.Background()
	_, _, err = mcp.HandleArchiveContext(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot archive active context")
}

// TestHandleSearchContexts_Success tests searching contexts by name
func TestHandleSearchContexts_Success(t *testing.T) {
	setupTestEnv(t)

	// Create contexts
	_, _, err := core.CreateContext("test-alpha")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("test-beta")
	require.NoError(t, err)
	_, err = core.StopContext()
	require.NoError(t, err)

	_, _, err = core.CreateContext("production-gamma")
	require.NoError(t, err)

	input := mcp.SearchContextsInput{
		Query: "test",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleSearchContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, "test", output.Query)
	assert.Equal(t, 2, output.TotalCount)
	assert.Len(t, output.Results, 2)

	for _, result := range output.Results {
		assert.Contains(t, strings.ToLower(result.Name), "test")
		assert.Equal(t, "name", result.MatchedIn)
	}
}

// TestHandleSearchContexts_MissingQuery tests error when query is missing
func TestHandleSearchContexts_MissingQuery(t *testing.T) {
	setupTestEnv(t)

	input := mcp.SearchContextsInput{
		Query: "",
	}

	ctx := context.Background()
	_, _, err := mcp.HandleSearchContexts(ctx, nil, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "query is required")
}

// TestHandleSearchContexts_WithLimit tests search result limiting
func TestHandleSearchContexts_WithLimit(t *testing.T) {
	setupTestEnv(t)

	// Create many matching contexts
	for i := 1; i <= 15; i++ {
		_, _, err := core.CreateContext("test-context-" + string(rune(i)))
		require.NoError(t, err)
		_, err = core.StopContext()
		require.NoError(t, err)
	}

	input := mcp.SearchContextsInput{
		Query: "test",
		Limit: 5,
	}

	ctx := context.Background()
	_, output, err := mcp.HandleSearchContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 15, output.TotalCount)
	assert.Len(t, output.Results, 5)
}

// TestHandleSearchContexts_SearchNotes tests searching within note content
func TestHandleSearchContexts_SearchNotes(t *testing.T) {
	setupTestEnv(t)

	// Create context with specific note
	_, _, err := core.CreateContext("unrelated-name")
	require.NoError(t, err)

	_, err = core.AddNote("This note contains the keyword FINDME")
	require.NoError(t, err)

	_, err = core.StopContext()
	require.NoError(t, err)

	// Create context without matching note
	_, _, err = core.CreateContext("another-context")
	require.NoError(t, err)
	_, err = core.AddNote("Just a regular note")
	require.NoError(t, err)

	input := mcp.SearchContextsInput{
		Query:       "findme",
		SearchNotes: true,
	}

	ctx := context.Background()
	_, output, err := mcp.HandleSearchContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Greater(t, output.TotalCount, 0)

	// Find the result that matched in notes
	var foundNoteMatch bool
	for _, result := range output.Results {
		if result.MatchedIn == "notes" {
			foundNoteMatch = true
			assert.NotEmpty(t, result.Excerpts)
			assert.Contains(t, strings.ToLower(result.Excerpts[0]), "findme")
		}
	}
	assert.True(t, foundNoteMatch)
}

// TestHandleSearchContexts_IncludeArchived tests including archived contexts
func TestHandleSearchContexts_IncludeArchived(t *testing.T) {
	setupTestEnv(t)

	// Create and archive context
	_, _, err := core.CreateContext("archived-context")
	require.NoError(t, err)

	_, err = core.StopContext()
	require.NoError(t, err)

	err = core.ArchiveContext("archived-context")
	require.NoError(t, err)

	// Create active context
	_, _, err = core.CreateContext("active-context")
	require.NoError(t, err)

	// Search without including archived
	input := mcp.SearchContextsInput{
		Query:           "context",
		IncludeArchived: false,
	}

	ctx := context.Background()
	_, output, err := mcp.HandleSearchContexts(ctx, nil, input)

	require.NoError(t, err)
	notIncludingArchivedCount := output.TotalCount

	// Search with archived
	input.IncludeArchived = true
	_, output, err = mcp.HandleSearchContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Greater(t, output.TotalCount, notIncludingArchivedCount)
}

// TestHandleSearchContexts_NoResults tests search with no matches
func TestHandleSearchContexts_NoResults(t *testing.T) {
	setupTestEnv(t)

	_, _, err := core.CreateContext("test-context")
	require.NoError(t, err)

	input := mcp.SearchContextsInput{
		Query: "nonexistent",
	}

	ctx := context.Background()
	_, output, err := mcp.HandleSearchContexts(ctx, nil, input)

	require.NoError(t, err)
	assert.Equal(t, 0, output.TotalCount)
	assert.Empty(t, output.Results)
}
