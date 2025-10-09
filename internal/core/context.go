package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
)

// CreateContext creates a new context, handling duplicate names and previous active context
func CreateContext(name string) (*models.Context, string, error) {
	if err := EnsureContextHome(); err != nil {
		return nil, "", err
	}

	// Sanitize the name for directory use
	sanitizedName := SanitizeContextName(name)

	// Resolve duplicate directory name
	finalDirName := resolveDuplicateName(sanitizedName)

	// Preserve original name for display, append suffix if duplicate
	displayName := name
	if finalDirName != sanitizedName {
		// If duplicate detected (e.g., sanitizedName_2), append to display name
		suffix := finalDirName[len(sanitizedName):]
		displayName = name + suffix
	}

	// Get current active context (if any)
	state, err := GetActiveContext()
	if err != nil {
		return nil, "", err
	}

	var previousContext string
	var transitionType models.TransitionType

	// Stop previous context if active
	if state.HasActiveContext() {
		previousContext = state.GetActiveContextName()
		if err := stopContextInternal(previousContext); err != nil {
			return nil, "", fmt.Errorf("failed to stop previous context: %w", err)
		}
		transitionType = models.TransitionSwitch
	} else {
		transitionType = models.TransitionStart
	}

	// Create new context
	now := time.Now()
	context := &models.Context{
		Name:             displayName,
		StartTime:        now,
		EndTime:          nil,
		Status:           "active",
		SubdirectoryPath: GetContextDir(finalDirName),
	}

	// Create context directory
	if err := CreateDir(context.SubdirectoryPath); err != nil {
		return nil, "", err
	}

	// Write meta.json
	if err := WriteJSON(GetMetaJSONPath(finalDirName), context); err != nil {
		return nil, "", err
	}

	// Create empty log files
	for _, path := range []string{
		GetNotesLogPath(finalDirName),
		GetFilesLogPath(finalDirName),
		GetTouchLogPath(finalDirName),
	} {
		if err := os.WriteFile(path, []byte{}, 0600); err != nil {
			return nil, "", err
		}
	}

	// Update state (store display name)
	if err := SetActiveContext(displayName); err != nil {
		return nil, "", err
	}

	// Log transition (use display name)
	transition := &models.ContextTransition{
		Timestamp:      now,
		NewContext:     &displayName,
		TransitionType: transitionType,
	}
	if transitionType == models.TransitionSwitch {
		transition.PreviousContext = &previousContext
	}

	if err := AppendLog(GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		return nil, "", err
	}

	return context, previousContext, nil
}

// resolveDuplicateName finds an available name by appending _2, _3, etc.
func resolveDuplicateName(name string) string {
	contextDir := GetContextDir(name)
	if !FileExists(contextDir) {
		return name
	}

	// Try _2, _3, _4, etc.
	for i := 2; ; i++ {
		candidateName := fmt.Sprintf("%s_%d", name, i)
		contextDir = GetContextDir(candidateName)
		if !FileExists(contextDir) {
			return candidateName
		}
	}
}

// StopContext stops the currently active context
func StopContext() (*models.Context, error) {
	state, err := GetActiveContext()
	if err != nil {
		return nil, err
	}

	if !state.HasActiveContext() {
		return nil, nil // No active context, not an error
	}

	contextName := state.GetActiveContextName()
	if err := stopContextInternal(contextName); err != nil {
		return nil, err
	}

	// Read the stopped context
	var context models.Context
	if err := ReadJSON(GetMetaJSONPath(contextName), &context); err != nil {
		return nil, err
	}

	// Clear active context
	if err := ClearActiveContext(); err != nil {
		return nil, err
	}

	// Log transition
	now := time.Now()
	transition := &models.ContextTransition{
		Timestamp:       now,
		PreviousContext: &contextName,
		NewContext:      nil,
		TransitionType:  models.TransitionStop,
	}

	if err := AppendLog(GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		return nil, err
	}

	return &context, nil
}

// stopContextInternal stops a context without clearing state (used internally)
func stopContextInternal(contextName string) error {
	// Read current context
	var context models.Context
	if err := ReadJSON(GetMetaJSONPath(contextName), &context); err != nil {
		return err
	}

	// Update to stopped
	now := time.Now()
	context.EndTime = &now
	context.Status = "stopped"

	// Write back
	return WriteJSON(GetMetaJSONPath(contextName), &context)
}

// AddNote adds a note to the active context
func AddNote(text string) (*models.Note, error) {
	state, err := GetActiveContext()
	if err != nil {
		return nil, err
	}

	if !state.HasActiveContext() {
		return nil, fmt.Errorf("no active context")
	}

	note := &models.Note{
		Timestamp:   time.Now(),
		TextContent: text,
	}

	if err := note.Validate(); err != nil {
		return nil, err
	}

	contextName := state.GetActiveContextName()
	logPath := GetNotesLogPath(contextName)

	if err := AppendLog(logPath, note.ToLogLine()); err != nil {
		return nil, err
	}

	return note, nil
}

// AddFile associates a file with the active context
func AddFile(filePath string) (*models.FileAssociation, error) {
	state, err := GetActiveContext()
	if err != nil {
		return nil, err
	}

	if !state.HasActiveContext() {
		return nil, fmt.Errorf("no active context")
	}

	// Normalize path
	normalizedPath, err := NormalizePath(filePath)
	if err != nil {
		return nil, err
	}

	file := &models.FileAssociation{
		Timestamp: time.Now(),
		FilePath:  normalizedPath,
	}

	contextName := state.GetActiveContextName()
	logPath := GetFilesLogPath(contextName)

	if err := AppendLog(logPath, file.ToLogLine()); err != nil {
		return nil, err
	}

	return file, nil
}

// AddTouch records a touch event in the active context
func AddTouch() (*models.TouchEvent, error) {
	state, err := GetActiveContext()
	if err != nil {
		return nil, err
	}

	if !state.HasActiveContext() {
		return nil, fmt.Errorf("no active context")
	}

	touch := &models.TouchEvent{
		Timestamp: time.Now(),
	}

	contextName := state.GetActiveContextName()
	logPath := GetTouchLogPath(contextName)

	if err := AppendLog(logPath, touch.ToLogLine()); err != nil {
		return nil, err
	}

	return touch, nil
}

// GetContext reads a context with all its associated data
func GetContext(contextName string) (*models.Context, []*models.Note, []*models.FileAssociation, []*models.TouchEvent, error) {
	// Read meta.json
	var context models.Context
	if err := ReadJSON(GetMetaJSONPath(contextName), &context); err != nil {
		return nil, nil, nil, nil, err
	}

	// Read notes
	notesLines, err := ReadLog(GetNotesLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var notes []*models.Note
	for _, line := range notesLines {
		if line == "" {
			continue
		}
		note, err := models.ParseNoteLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		notes = append(notes, note)
	}

	// Read files
	filesLines, err := ReadLog(GetFilesLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var files []*models.FileAssociation
	for _, line := range filesLines {
		if line == "" {
			continue
		}
		file, err := models.ParseFileLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		files = append(files, file)
	}

	// Read touches
	touchesLines, err := ReadLog(GetTouchLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var touches []*models.TouchEvent
	for _, line := range touchesLines {
		if line == "" {
			continue
		}
		touch, err := models.ParseTouchLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		touches = append(touches, touch)
	}

	return &context, notes, files, touches, nil
}

// ListContexts returns all contexts sorted by start time (most recent first)
func ListContexts() ([]*models.Context, error) {
	dirs, err := ListContextDirs()
	if err != nil {
		return nil, err
	}

	var contexts []*models.Context
	for _, dir := range dirs {
		var context models.Context
		if err := ReadJSON(GetMetaJSONPath(dir), &context); err != nil {
			continue // Skip contexts with invalid meta.json
		}
		contexts = append(contexts, &context)
	}

	// Sort by start time, most recent first
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].StartTime.After(contexts[j].StartTime)
	})

	return contexts, nil
}

// GetTransitions returns all context transitions
func GetTransitions() ([]*models.ContextTransition, error) {
	lines, err := ReadLog(GetTransitionsLogPath())
	if err != nil {
		return nil, err
	}

	var transitions []*models.ContextTransition
	for _, line := range lines {
		if line == "" {
			continue
		}
		transition, err := models.ParseTransitionLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		transitions = append(transitions, transition)
	}

	return transitions, nil
}

// ExportContext exports a context to a markdown file
func ExportContext(contextName string, outputPath string) (string, error) {
	// Load context data
	ctx, notes, files, touches, err := GetContext(contextName)
	if err != nil {
		return "", fmt.Errorf("context %q not found", contextName)
	}

	// Convert to model types for export
	var noteModels []models.Note
	for _, n := range notes {
		noteModels = append(noteModels, *n)
	}

	var fileModels []models.FileAssociation
	for _, f := range files {
		fileModels = append(fileModels, *f)
	}

	// Generate markdown content
	content := output.FormatExportMarkdown(ctx, noteModels, fileModels, len(touches))

	// Determine output file path
	if outputPath == "" {
		// Default: sanitized context name in current directory
		sanitized := SanitizeFilename(contextName)
		outputPath = sanitized + ".md"
	}

	// Create parent directories if needed
	if err := CreateParentDirs(outputPath); err != nil {
		return "", err
	}

	// Write file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write export file: %w", err)
	}

	return outputPath, nil
}

// ExportAllContexts exports all contexts to separate markdown files in a directory
func ExportAllContexts(outputDir string) ([]string, error) {
	contexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	if err := CreateDir(outputDir); err != nil {
		return nil, err
	}

	var exportedPaths []string
	for _, ctx := range contexts {
		sanitized := SanitizeFilename(ctx.Name)
		outputPath := filepath.Join(outputDir, sanitized+".md")

		path, err := ExportContext(ctx.Name, outputPath)
		if err != nil {
			continue // Skip failed exports
		}
		exportedPaths = append(exportedPaths, path)
	}

	return exportedPaths, nil
}

// ArchiveContext marks a context as archived
func ArchiveContext(contextName string) error {
	// Load context
	var ctx models.Context
	metaPath := GetMetaJSONPath(contextName)
	if err := ReadJSON(metaPath, &ctx); err != nil {
		return fmt.Errorf("context %q not found", contextName)
	}

	// Validate: cannot archive active context
	if ctx.Status == "active" {
		return fmt.Errorf("cannot archive active context %q - stop it first", contextName)
	}

	// Check if already archived
	if ctx.IsArchived {
		return fmt.Errorf("context %q is already archived", contextName)
	}

	// Set archived flag
	ctx.IsArchived = true

	// Write updated meta.json
	if err := WriteJSON(metaPath, &ctx); err != nil {
		return fmt.Errorf("failed to update context: %w", err)
	}

	return nil
}

// DeleteContext permanently removes a context and all its data
func DeleteContext(contextName string, force bool, confirmed bool) error {
	// Load context
	var ctx models.Context
	metaPath := GetMetaJSONPath(contextName)
	if err := ReadJSON(metaPath, &ctx); err != nil {
		return fmt.Errorf("context %q not found", contextName)
	}

	// Validate: cannot delete active context
	if ctx.Status == "active" {
		return fmt.Errorf("cannot delete active context %q - stop it first", contextName)
	}

	// Require confirmation unless force or already confirmed
	if !force && !confirmed {
		return fmt.Errorf("deletion requires confirmation")
	}

	// Remove context directory
	// Use sanitized name to get the actual directory path (e.g., "demo: Test" → "demo__Test")
	sanitizedName := SanitizeContextName(contextName)
	actualDir := filepath.Join(GetContextHome(), sanitizedName)

	if err := os.RemoveAll(actualDir); err != nil {
		return fmt.Errorf("failed to delete context directory %q: %w", actualDir, err)
	}

	return nil
}

// LoadContext reads a context by name (Sprint 2: for backward compatibility testing)
func LoadContext(contextName string) (*models.Context, error) {
	var ctx models.Context
	if err := ReadJSON(GetMetaJSONPath(contextName), &ctx); err != nil {
		return nil, fmt.Errorf("context %q not found", contextName)
	}
	return &ctx, nil
}

// GetHomeDir returns the context home directory (for tests)
func GetHomeDir() string {
	return GetContextHome()
}

// ContextFilter defines filtering options for listing contexts
type ContextFilter struct {
	Project      string // Filter by project name (case-insensitive)
	Search       string // Filter by substring in name (case-insensitive)
	Limit        int    // Limit number of results (0 = no limit)
	Archived     bool   // Show only archived contexts
	ActiveOnly   bool   // Show only active context
	ShowArchived bool   // Include archived in results (default: exclude)
}

// ListContextsFiltered returns contexts matching the filter criteria
func ListContextsFiltered(filter ContextFilter) ([]*models.Context, error) {
	// Get all contexts
	allContexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Context

	// Apply filters
	for _, ctx := range allContexts {
		// Skip archived contexts unless explicitly requested
		if ctx.IsArchived && !filter.ShowArchived && !filter.Archived {
			continue
		}

		// If --archived flag, show ONLY archived
		if filter.Archived && !ctx.IsArchived {
			continue
		}

		// If --active-only, show ONLY active
		if filter.ActiveOnly && !ctx.IsActive() {
			continue
		}

		// Project filter (case-insensitive)
		if filter.Project != "" {
			projectName := ExtractProjectName(ctx.Name)
			if !strings.EqualFold(projectName, filter.Project) {
				continue
			}
		}

		// Search filter (case-insensitive substring)
		if filter.Search != "" {
			if !strings.Contains(strings.ToLower(ctx.Name), strings.ToLower(filter.Search)) {
				continue
			}
		}

		filtered = append(filtered, ctx)
	}

	// Apply limit
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}

	return filtered, nil
}
