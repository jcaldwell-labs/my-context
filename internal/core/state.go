package core

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

// GetActiveContext reads the current active context from state.json
func GetActiveContext() (*models.AppState, error) {
	if err := EnsureContextHome(); err != nil {
		return nil, err
	}

	statePath := GetStateFilePath()

	// If state file doesn't exist, return empty state
	if !FileExists(statePath) {
		return &models.AppState{
			ActiveContext: nil,
			LastUpdated:   time.Now(),
		}, nil
	}

	var state models.AppState
	if err := ReadJSON(statePath, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// SetActiveContext updates the active context in state.json
func SetActiveContext(contextName string) error {
	if err := EnsureContextHome(); err != nil {
		return err
	}

	state := &models.AppState{
		ActiveContext: &contextName,
		LastUpdated:   time.Now(),
	}

	return WriteJSON(GetStateFilePath(), state)
}

// ClearActiveContext removes the active context from state.json
func ClearActiveContext() error {
	if err := EnsureContextHome(); err != nil {
		return err
	}

	state := &models.AppState{
		ActiveContext: nil,
		LastUpdated:   time.Now(),
	}

	return WriteJSON(GetStateFilePath(), state)
}

// FindContextByName finds a context by its display name (not directory name)
func FindContextByName(name string) (*models.Context, error) {
	contexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	for _, ctx := range contexts {
		if ctx.Name == name {
			return ctx, nil
		}
	}

	return nil, fmt.Errorf("context %q not found", name)
}

// GetNoteCount returns the number of notes in a context
func GetNoteCount(contextName string) (int, error) {
	notesPath := GetNotesLogPath(contextName)
	lines, err := ReadLog(notesPath)
	if err != nil {
		return 0, err
	}

	// Count non-empty lines (each represents a note)
	count := 0
	for _, line := range lines {
		if line != "" {
			count++
		}
	}

	return count, nil
}

// GetLastActiveTime returns the most recent time when the context was active
func GetLastActiveTime(contextName string) (time.Time, error) {
	transitions, err := GetTransitions()
	if err != nil {
		return time.Time{}, err
	}

	var lastActive time.Time

	// Look through transitions in reverse order (most recent first)
	for i := len(transitions) - 1; i >= 0; i-- {
		transition := transitions[i]

		// Context became active if it was started or switched to
		if transition.NewContext != nil && *transition.NewContext == contextName {
			if transition.TransitionType == models.TransitionStart || transition.TransitionType == models.TransitionSwitch {
				lastActive = transition.Timestamp
				break
			}
		}
	}

	if lastActive.IsZero() {
		return time.Time{}, fmt.Errorf("no active period found for context %q", contextName)
	}

	return lastActive, nil
}

// GetMostRecentStopped returns the most recently stopped context
func GetMostRecentStopped() (*models.Context, error) {
	contexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	// Filter to stopped contexts only
	var stoppedContexts []*models.Context
	for _, ctx := range contexts {
		if ctx.Status == "stopped" {
			stoppedContexts = append(stoppedContexts, ctx)
		}
	}

	if len(stoppedContexts) == 0 {
		return nil, fmt.Errorf("no stopped contexts found")
	}

	// Sort by EndTime (most recent first)
	sort.Slice(stoppedContexts, func(i, j int) bool {
		return stoppedContexts[i].EndTime != nil && stoppedContexts[j].EndTime != nil &&
			stoppedContexts[i].EndTime.After(*stoppedContexts[j].EndTime)
	})

	return stoppedContexts[0], nil
}

// FindContextsByPattern finds contexts matching a pattern (supports glob-style wildcards)
func FindContextsByPattern(pattern string) ([]*models.Context, error) {
	allContexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	// Filter to stopped contexts only (resume command only works with stopped contexts)
	var stoppedContexts []*models.Context
	for _, ctx := range allContexts {
		if ctx.Status == "stopped" {
			stoppedContexts = append(stoppedContexts, ctx)
		}
	}

	if pattern == "" {
		return stoppedContexts, nil
	}

	var matches []*models.Context

	// Simple glob-style pattern matching (* wildcards)
	// Convert glob pattern to a simple matcher
	patternParts := strings.Split(pattern, "*")

	for _, ctx := range stoppedContexts {
		if matchesPattern(ctx.Name, patternParts) {
			matches = append(matches, ctx)
		}
	}

	// Sort by EndTime (most recent first) for consistent ordering
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].EndTime != nil && matches[j].EndTime != nil &&
			matches[i].EndTime.After(*matches[j].EndTime)
	})

	return matches, nil
}

// matchesPattern checks if a context name matches a glob pattern
func matchesPattern(name string, patternParts []string) bool {
	if len(patternParts) == 0 {
		return true
	}

	// Handle simple cases
	if len(patternParts) == 1 {
		// No wildcards
		return name == patternParts[0]
	}

	// Check prefix
	if !strings.HasPrefix(name, patternParts[0]) {
		return false
	}

	remainingName := name[len(patternParts[0]):]

	// Check suffix
	lastPart := patternParts[len(patternParts)-1]
	if !strings.HasSuffix(remainingName, lastPart) {
		return false
	}

	// For multiple wildcards, we do a simple substring check
	// This is not a full glob implementation but covers common use cases
	for i := 1; i < len(patternParts)-1; i++ {
		if !strings.Contains(remainingName, patternParts[i]) {
			return false
		}
	}

	return true
}

// FindRelatedContexts finds contexts with similar names (shared prefixes)
func FindRelatedContexts(contextName string) ([]*models.Context, error) {
	allContexts, err := ListContexts()
	if err != nil {
		return nil, err
	}

	// Filter to stopped contexts only
	var stoppedContexts []*models.Context
	for _, ctx := range allContexts {
		if ctx.Status == "stopped" && ctx.Name != contextName {
			stoppedContexts = append(stoppedContexts, ctx)
		}
	}

	// Extract prefix from the context name
	prefix := extractContextPrefix(contextName)

	if prefix == "" {
		return []*models.Context{}, nil
	}

	// Find contexts that share the same prefix
	var related []*models.Context
	for _, ctx := range stoppedContexts {
		if extractContextPrefix(ctx.Name) == prefix {
			related = append(related, ctx)
		}
	}

	// Sort by most recent EndTime
	sort.Slice(related, func(i, j int) bool {
		return related[i].EndTime != nil && related[j].EndTime != nil &&
			related[i].EndTime.After(*related[j].EndTime)
	})

	// Return up to 3 related contexts
	if len(related) > 3 {
		related = related[:3]
	}

	return related, nil
}

// extractContextPrefix extracts the common prefix from a context name
// Examples:
//
//	"my-project: feature-auth" -> "my-project"
//	"feature-login" -> "feature-login" (no colon, return whole name)
//	"ps-cli: Phase 1" -> "ps-cli"
//	"bug-fix-123" -> "bug-fix-123" (no colon, return whole name)
func extractContextPrefix(contextName string) string {
	// Split on first colon and take the first part
	parts := strings.SplitN(contextName, ":", 2)
	if len(parts) > 1 {
		// Trim whitespace and return the prefix
		return strings.TrimSpace(parts[0])
	}

	// No colon found, return the whole name as prefix
	return contextName
}
