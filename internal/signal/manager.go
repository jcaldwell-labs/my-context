package signal

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/jefferycaldwell/my-context-copilot/pkg/utils"
)

// Manager handles signal file operations
type Manager struct {
	signalsDir string
	mu         sync.RWMutex
}

// NewManager creates a new signal manager
func NewManager(signalsDir string) (*Manager, error) {
	if err := utils.EnsureDir(signalsDir); err != nil {
		return nil, fmt.Errorf("failed to create signals directory: %w", err)
	}

	return &Manager{
		signalsDir: signalsDir,
	}, nil
}

// CreateSignal creates a new signal file
func (m *Manager) CreateSignal(name string) (*models.Signal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	signal := models.NewSignal(name, m.signalsDir)

	if signal.Exists() {
		return nil, fmt.Errorf("signal '%s' already exists", name)
	}

	if err := signal.Create(); err != nil {
		return nil, fmt.Errorf("failed to create signal: %w", err)
	}

	return signal, nil
}

// GetSignal returns information about a specific signal
func (m *Manager) GetSignal(name string) (*models.Signal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	signal := models.NewSignal(name, m.signalsDir)

	if !signal.Exists() {
		return nil, fmt.Errorf("signal '%s' does not exist", name)
	}

	// Load actual file information
	return models.LoadSignalFromFile(signal.Path)
}

// ListSignals returns all existing signals
func (m *Manager) ListSignals() ([]models.SignalInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	files, err := utils.ListFiles(m.signalsDir, "*.signal")
	if err != nil {
		return nil, fmt.Errorf("failed to list signal files: %w", err)
	}

	var signals []models.SignalInfo
	for _, file := range files {
		signal, err := models.LoadSignalFromFile(file)
		if err != nil {
			// Log error but continue with other signals
			continue
		}
		signals = append(signals, signal.ToInfo())
	}

	return signals, nil
}

// WaitForSignal blocks until a signal appears or timeout expires
func (m *Manager) WaitForSignal(name string, timeout time.Duration) error {
	m.mu.RLock()
	signal := models.NewSignal(name, m.signalsDir)
	m.mu.RUnlock()

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond) // Check every 500ms
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.RLock()
			exists := signal.Exists()
			m.mu.RUnlock()

			if exists {
				return nil
			}

			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for signal '%s'", name)
			}

		case <-time.After(time.Until(deadline)):
			return fmt.Errorf("timeout waiting for signal '%s'", name)
		}
	}
}

// ClearSignal removes a signal file
func (m *Manager) ClearSignal(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	signal := models.NewSignal(name, m.signalsDir)

	if !signal.Exists() {
		return fmt.Errorf("signal '%s' does not exist", name)
	}

	if err := signal.Remove(); err != nil {
		return fmt.Errorf("failed to remove signal: %w", err)
	}

	return nil
}

// ClearAllSignals removes all signal files (useful for cleanup)
func (m *Manager) ClearAllSignals() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	files, err := utils.ListFiles(m.signalsDir, "*.signal")
	if err != nil {
		return fmt.Errorf("failed to list signal files: %w", err)
	}

	for _, file := range files {
		signal, err := models.LoadSignalFromFile(file)
		if err != nil {
			continue // Skip invalid files
		}

		if err := signal.Remove(); err != nil {
			return fmt.Errorf("failed to remove signal %s: %w", filepath.Base(file), err)
		}
	}

	return nil
}

// SignalExists checks if a signal exists without loading full details
func (m *Manager) SignalExists(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	signal := models.NewSignal(name, m.signalsDir)
	return signal.Exists()
}
