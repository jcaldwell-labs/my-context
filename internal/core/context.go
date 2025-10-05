package core

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// CreateContext creates a new context, handling duplicate names and previous active context
func CreateContext(name string) (*models.Context, string, error) {
	if err := EnsureContextHome(); err != nil {
		return nil, "", err
	}

	// Sanitize the name
	sanitizedName := SanitizeContextName(name)

	// Resolve duplicate name
	finalName := resolveDuplicateName(sanitizedName)

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
		Name:             finalName,
		StartTime:        now,
		EndTime:          nil,
		Status:           "active",
		SubdirectoryPath: GetContextDir(finalName),
	}

	// Create context directory
	if err := CreateDir(context.SubdirectoryPath); err != nil {
		return nil, "", err
	}

	// Write meta.json
	if err := WriteJSON(GetMetaJSONPath(finalName), context); err != nil {
		return nil, "", err
	}

	// Create empty log files
	for _, path := range []string{
		GetNotesLogPath(finalName),
		GetFilesLogPath(finalName),
		GetTouchLogPath(finalName),
	} {
		if err := os.WriteFile(path, []byte{}, 0600); err != nil {
			return nil, "", err
		}
	}

	// Update state
	if err := SetActiveContext(finalName); err != nil {
		return nil, "", err
	}

	// Log transition
	transition := &models.ContextTransition{
		Timestamp:      now,
		NewContext:     &finalName,
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

