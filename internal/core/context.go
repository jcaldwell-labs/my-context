package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	intmodels "github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/jefferycaldwell/my-context-copilot/internal/output"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// CreateContextWithMetadata creates a new context with metadata, handling duplicate names and previous active context
func CreateContextWithMetadata(name string, createdBy string, parent string, labels []string) (*pkgmodels.ContextWithMetadata, string, error) {
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
	var transitionType intmodels.TransitionType

	// Stop previous context if active
	if state.HasActiveContext() {
		previousContext = state.GetActiveContextName()
		if err := stopContextInternal(previousContext); err != nil {
			return nil, "", fmt.Errorf("failed to stop previous context: %w", err)
		}
		transitionType = intmodels.TransitionSwitch
	} else {
		transitionType = intmodels.TransitionStart
	}

	// Create new context with metadata
	context := pkgmodels.NewContextWithMetadata(displayName, createdBy, parent, labels)

	// Override the subdirectory path
	context.SubdirectoryPath = GetContextDir(finalDirName)

	// Create context directory
	if err := CreateDir(context.SubdirectoryPath); err != nil {
		return nil, "", err
	}

	// Write meta.json
	if err := WriteJSON(GetMetaJSONPath(displayName), context); err != nil {
		return nil, "", fmt.Errorf("failed to write context metadata: %w", err)
	}

	// Create empty logs
	now := time.Now()
	if err := AppendLog(GetNotesLogPath(displayName), fmt.Sprintf("[%s] Context started", now.Format(time.RFC3339))); err != nil {
		return nil, "", fmt.Errorf("failed to initialize notes log: %w", err)
	}

	if err := AppendLog(GetFilesLogPath(displayName), fmt.Sprintf("[%s] Context started", now.Format(time.RFC3339))); err != nil {
		return nil, "", fmt.Errorf("failed to initialize files log: %w", err)
	}

	if err := AppendLog(GetTouchLogPath(displayName), fmt.Sprintf("[%s] Context started", now.Format(time.RFC3339))); err != nil {
		return nil, "", fmt.Errorf("failed to initialize touches log: %w", err)
	}

	// Update state
	newState := &intmodels.AppState{
		ActiveContext: &displayName,
		LastUpdated:   now,
	}

	if err := WriteJSON(GetStateFilePath(), newState); err != nil {
		return nil, "", fmt.Errorf("failed to update state: %w", err)
	}

	// Log transition
	transition := &intmodels.ContextTransition{
		Timestamp:      now,
		NewContext:     &displayName,
		TransitionType: transitionType,
	}

	if transitionType == intmodels.TransitionSwitch {
		transition.PreviousContext = &previousContext
	}

	if err := AppendLog(GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		return nil, "", fmt.Errorf("failed to log transition: %w", err)
	}

	return context, previousContext, nil
}

// CreateContext creates a new context, handling duplicate names and previous active context
func CreateContext(name string) (*intmodels.Context, string, error) {
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
	var transitionType intmodels.TransitionType

	// Stop previous context if active
	if state.HasActiveContext() {
		previousContext = state.GetActiveContextName()
		if err := stopContextInternal(previousContext); err != nil {
			return nil, "", fmt.Errorf("failed to stop previous context: %w", err)
		}
		transitionType = intmodels.TransitionSwitch
	} else {
		transitionType = intmodels.TransitionStart
	}

	// Create new context
	now := time.Now()
	context := &intmodels.Context{
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
	transition := &intmodels.ContextTransition{
		Timestamp:      now,
		NewContext:     &displayName,
		TransitionType: transitionType,
	}
	if transitionType == intmodels.TransitionSwitch {
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
func StopContext() (*intmodels.Context, error) {
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
	var context intmodels.Context
	if err := ReadJSON(GetMetaJSONPath(contextName), &context); err != nil {
		return nil, err
	}

	// Clear active context
	if err := ClearActiveContext(); err != nil {
		return nil, err
	}

	// Log transition
	now := time.Now()
	transition := &intmodels.ContextTransition{
		Timestamp:       now,
		PreviousContext: &contextName,
		NewContext:      nil,
		TransitionType:  intmodels.TransitionStop,
	}

	if err := AppendLog(GetTransitionsLogPath(), transition.ToLogLine()); err != nil {
		return nil, err
	}

	return &context, nil
}

// stopContextInternal stops a context without clearing state (used internally)
func stopContextInternal(contextName string) error {
	// Read current context
	var context intmodels.Context
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
func AddNote(text string) (*intmodels.Note, error) {
	state, err := GetActiveContext()
	if err != nil {
		return nil, err
	}

	if !state.HasActiveContext() {
		return nil, fmt.Errorf("no active context")
	}

	note := &intmodels.Note{
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
func AddFile(filePath string) (*intmodels.FileAssociation, error) {
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

	file := &intmodels.FileAssociation{
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
func AddTouch() (*intmodels.TouchEvent, error) {
	state, err := GetActiveContext()
	if err != nil {
		return nil, err
	}

	if !state.HasActiveContext() {
		return nil, fmt.Errorf("no active context")
	}

	touch := &intmodels.TouchEvent{
		Timestamp: time.Now(),
	}

	contextName := state.GetActiveContextName()
	logPath := GetTouchLogPath(contextName)

	if err := AppendLog(logPath, touch.ToLogLine()); err != nil {
		return nil, err
	}

	return touch, nil
}

// GetContextWithMetadata reads a context with metadata and all associated data
func GetContextWithMetadata(contextName string) (*pkgmodels.ContextWithMetadata, []*intmodels.Note, []*intmodels.FileAssociation, []*intmodels.TouchEvent, error) {
	// Read meta.json
	var context pkgmodels.ContextWithMetadata
	if err := ReadJSON(GetMetaJSONPath(contextName), &context); err != nil {
		// If reading as extended context fails, try reading as old context and convert
		var oldContext intmodels.Context
		if err := ReadJSON(GetMetaJSONPath(contextName), &oldContext); err != nil {
			return nil, nil, nil, nil, err
		}
		// Convert old context to extended context with empty metadata
		context = pkgmodels.ContextWithMetadata{
			Name:             oldContext.Name,
			StartTime:        oldContext.StartTime,
			EndTime:          oldContext.EndTime,
			Status:           oldContext.Status,
			SubdirectoryPath: oldContext.SubdirectoryPath,
			IsArchived:       oldContext.IsArchived,
			Metadata:         pkgmodels.ContextMetadata{}, // Empty metadata
		}
	}

	// Read notes, files, touches (same as GetContext)
	notesLines, err := ReadLog(GetNotesLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var notes []*intmodels.Note
	for _, line := range notesLines {
		if line == "" {
			continue
		}
		note, err := intmodels.ParseNoteLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		notes = append(notes, note)
	}

	filesLines, err := ReadLog(GetFilesLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var files []*intmodels.FileAssociation
	for _, line := range filesLines {
		if line == "" {
			continue
		}
		file, err := intmodels.ParseFileLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		files = append(files, file)
	}

	touchesLines, err := ReadLog(GetTouchLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var touches []*intmodels.TouchEvent
	for _, line := range touchesLines {
		if line == "" {
			continue
		}
		touch, err := intmodels.ParseTouchLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		touches = append(touches, touch)
	}

	return &context, notes, files, touches, nil
}

// GetContext reads a context with all its associated data
func GetContext(contextName string) (*intmodels.Context, []*intmodels.Note, []*intmodels.FileAssociation, []*intmodels.TouchEvent, error) {
	// Read meta.json
	var context intmodels.Context
	if err := ReadJSON(GetMetaJSONPath(contextName), &context); err != nil {
		return nil, nil, nil, nil, err
	}

	// Read notes
	notesLines, err := ReadLog(GetNotesLogPath(contextName))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var notes []*intmodels.Note
	for _, line := range notesLines {
		if line == "" {
			continue
		}
		note, err := intmodels.ParseNoteLogLine(line)
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

	var files []*intmodels.FileAssociation
	for _, line := range filesLines {
		if line == "" {
			continue
		}
		file, err := intmodels.ParseFileLogLine(line)
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

	var touches []*intmodels.TouchEvent
	for _, line := range touchesLines {
		if line == "" {
			continue
		}
		touch, err := intmodels.ParseTouchLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		touches = append(touches, touch)
	}

	return &context, notes, files, touches, nil
}

// ListContexts returns all contexts sorted by start time (most recent first)
func ListContexts() ([]*intmodels.Context, error) {
	dirs, err := ListContextDirs()
	if err != nil {
		return nil, err
	}

	var contexts []*intmodels.Context
	for _, dir := range dirs {
		var context intmodels.Context
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
func GetTransitions() ([]*intmodels.ContextTransition, error) {
	lines, err := ReadLog(GetTransitionsLogPath())
	if err != nil {
		return nil, err
	}

	var transitions []*intmodels.ContextTransition
	for _, line := range lines {
		if line == "" {
			continue
		}
		transition, err := intmodels.ParseTransitionLogLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		transitions = append(transitions, transition)
	}

	return transitions, nil
}

// ExportContext exports a context to a markdown file
func ExportContext(contextName string, outputPath string, asJSON bool) (string, error) {
	// Load context data
	ctx, notes, files, touches, err := GetContext(contextName)
	if err != nil {
		return "", fmt.Errorf("context %q not found", contextName)
	}

	// Convert to model types for export
	var noteModels []intmodels.Note
	for _, n := range notes {
		noteModels = append(noteModels, *n)
	}

	var fileModels []intmodels.FileAssociation
	for _, f := range files {
		fileModels = append(fileModels, *f)
	}

	// Generate content based on format
	var content string
	var defaultExt string
	if asJSON {
		jsonContent, err := output.FormatExportJSON(ctx, noteModels, fileModels, len(touches))
		if err != nil {
			return "", fmt.Errorf("failed to generate JSON export: %w", err)
		}
		content = jsonContent
		defaultExt = ".json"
	} else {
		content = output.FormatExportMarkdown(ctx, noteModels, fileModels, len(touches))
		defaultExt = ".md"
	}

	// Determine output file path
	if outputPath == "" {
		// Default: sanitized context name in current directory
		sanitized := SanitizeFilename(contextName)
		outputPath = sanitized + defaultExt
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

// ExportAllContexts exports all contexts to separate files in a directory
func ExportAllContexts(outputDir string, asJSON bool) ([]string, error) {
	contexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	if err := CreateDir(outputDir); err != nil {
		return nil, err
	}

	ext := ".md"
	if asJSON {
		ext = ".json"
	}

	var exportedPaths []string
	for _, ctx := range contexts {
		sanitized := SanitizeFilename(ctx.Name)
		outputPath := filepath.Join(outputDir, sanitized+ext)

		path, err := ExportContext(ctx.Name, outputPath, asJSON)
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
	var ctx intmodels.Context
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
	var ctx intmodels.Context
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
func LoadContext(contextName string) (*intmodels.Context, error) {
	var ctx intmodels.Context
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
func ListContextsFiltered(filter ContextFilter) ([]*intmodels.Context, error) {
	// Get all contexts
	allContexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	var filtered []*intmodels.Context

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

// AddTags adds tags to a context and returns the list of newly added tags
func AddTags(contextName string, tags []string) ([]string, error) {
	// Load context with metadata
	ctx, _, _, _, err := GetContextWithMetadata(contextName)
	if err != nil {
		return nil, fmt.Errorf("context %q not found", contextName)
	}

	// Track which tags were actually added (not duplicates)
	var added []string
	existingTags := make(map[string]bool)
	for _, tag := range ctx.Metadata.Labels {
		existingTags[tag] = true
	}

	// Add new tags
	for _, tag := range tags {
		if !existingTags[tag] {
			ctx.Metadata.Labels = append(ctx.Metadata.Labels, tag)
			added = append(added, tag)
			existingTags[tag] = true
		}
	}

	// Validate updated context
	if err := ctx.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Write updated meta.json
	metaPath := GetMetaJSONPath(contextName)
	if err := WriteJSON(metaPath, ctx); err != nil {
		return nil, fmt.Errorf("failed to update context: %w", err)
	}

	return added, nil
}

// RemoveTags removes tags from a context and returns the list of removed tags
func RemoveTags(contextName string, tags []string) ([]string, error) {
	// Load context with metadata
	ctx, _, _, _, err := GetContextWithMetadata(contextName)
	if err != nil {
		return nil, fmt.Errorf("context %q not found", contextName)
	}

	// Build set of tags to remove
	toRemove := make(map[string]bool)
	for _, tag := range tags {
		toRemove[tag] = true
	}

	// Filter out tags to remove
	var newLabels []string
	var removed []string
	for _, tag := range ctx.Metadata.Labels {
		if toRemove[tag] {
			removed = append(removed, tag)
		} else {
			newLabels = append(newLabels, tag)
		}
	}

	// Update labels
	ctx.Metadata.Labels = newLabels

	// Write updated meta.json
	metaPath := GetMetaJSONPath(contextName)
	if err := WriteJSON(metaPath, ctx); err != nil {
		return nil, fmt.Errorf("failed to update context: %w", err)
	}

	return removed, nil
}

// GetContextTags returns all tags for a specific context
func GetContextTags(contextName string) ([]string, error) {
	// Load context with metadata
	ctx, _, _, _, err := GetContextWithMetadata(contextName)
	if err != nil {
		return nil, fmt.Errorf("context %q not found", contextName)
	}

	return ctx.Metadata.Labels, nil
}

// GetAllTags returns a map of all tags used across all contexts with their usage counts
func GetAllTags() (map[string]int, error) {
	dirs, err := ListContextDirs()
	if err != nil {
		return nil, err
	}

	tagCounts := make(map[string]int)

	for _, dir := range dirs {
		// Try to load as ContextWithMetadata
		ctx, _, _, _, err := GetContextWithMetadata(dir)
		if err != nil {
			continue // Skip contexts that can't be loaded
		}

		// Count each tag
		for _, tag := range ctx.Metadata.Labels {
			tagCounts[tag]++
		}
	}

	return tagCounts, nil
}

// SetParent sets the parent context for a child context
func SetParent(childName string, parentName string) error {
	// Load child context with metadata
	ctx, _, _, _, err := GetContextWithMetadata(childName)
	if err != nil {
		return fmt.Errorf("child context %q not found", childName)
	}

	// Update parent
	ctx.Metadata.Parent = parentName

	// Validate updated context
	if err := ctx.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Write updated meta.json
	metaPath := GetMetaJSONPath(childName)
	if err := WriteJSON(metaPath, ctx); err != nil {
		return fmt.Errorf("failed to update context: %w", err)
	}

	return nil
}

// ClearParent removes the parent relationship from a context
func ClearParent(contextName string) error {
	// Load context with metadata
	ctx, _, _, _, err := GetContextWithMetadata(contextName)
	if err != nil {
		return fmt.Errorf("context %q not found", contextName)
	}

	// Clear parent
	ctx.Metadata.Parent = ""

	// Write updated meta.json
	metaPath := GetMetaJSONPath(contextName)
	if err := WriteJSON(metaPath, ctx); err != nil {
		return fmt.Errorf("failed to update context: %w", err)
	}

	return nil
}

// GetChildren returns all contexts that have the given context as their parent
func GetChildren(parentName string) ([]string, error) {
	dirs, err := ListContextDirs()
	if err != nil {
		return nil, err
	}

	var children []string

	for _, dir := range dirs {
		// Try to load as ContextWithMetadata
		ctx, _, _, _, err := GetContextWithMetadata(dir)
		if err != nil {
			continue // Skip contexts that can't be loaded
		}

		// Check if this context's parent matches
		if ctx.Metadata.Parent == parentName {
			children = append(children, ctx.Name)
		}
	}

	return children, nil
}

// GetContextTree builds a tree structure starting from a root context
type ContextTreeNode struct {
	Name     string
	Children []*ContextTreeNode
}

func GetContextTree(rootName string) (*ContextTreeNode, error) {
	// Verify root context exists
	if _, err := LoadContext(rootName); err != nil {
		return nil, fmt.Errorf("context %q not found", rootName)
	}

	// Build tree recursively
	return buildTreeNode(rootName, make(map[string]bool))
}

func buildTreeNode(name string, visited map[string]bool) (*ContextTreeNode, error) {
	// Prevent infinite loops in case of circular dependencies
	if visited[name] {
		return &ContextTreeNode{Name: name + " (circular)"}, nil
	}
	visited[name] = true

	node := &ContextTreeNode{Name: name}

	// Get children
	children, err := GetChildren(name)
	if err != nil {
		return nil, err
	}

	// Recursively build child nodes
	for _, childName := range children {
		childNode, err := buildTreeNode(childName, visited)
		if err != nil {
			continue // Skip children that can't be built
		}
		node.Children = append(node.Children, childNode)
	}

	return node, nil
}

// GetRootContexts returns all contexts that don't have a parent
func GetRootContexts() ([]string, error) {
	dirs, err := ListContextDirs()
	if err != nil {
		return nil, err
	}

	var roots []string

	for _, dir := range dirs {
		// Try to load as ContextWithMetadata
		ctx, _, _, _, err := GetContextWithMetadata(dir)
		if err != nil {
			continue // Skip contexts that can't be loaded
		}

		// Check if this context has no parent
		if ctx.Metadata.Parent == "" {
			roots = append(roots, ctx.Name)
		}
	}

	return roots, nil
}
