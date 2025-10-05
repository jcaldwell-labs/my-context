package core

import (
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
