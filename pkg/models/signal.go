package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Signal represents an event notification signal
type Signal struct {
	Name      string    `json:"name"`       // Signal name/identifier
	CreatedAt time.Time `json:"created_at"` // When the signal was created
	Path      string    `json:"path"`       // Full path to the signal file
}

// NewSignal creates a new signal instance
func NewSignal(name string, signalsDir string) *Signal {
	now := time.Now()
	filename := fmt.Sprintf("%s.signal", name)
	path := filepath.Join(signalsDir, filename)

	return &Signal{
		Name:      name,
		CreatedAt: now,
		Path:      path,
	}
}

// Exists checks if the signal file exists on disk
func (s *Signal) Exists() bool {
	_, err := os.Stat(s.Path)
	return !os.IsNotExist(err)
}

// Create writes the signal file to disk
func (s *Signal) Create() error {
	// Ensure directory exists
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create signals directory: %w", err)
	}

	// Create the signal file with timestamp as content
	content := s.CreatedAt.Format(time.RFC3339) + "\n"
	if err := os.WriteFile(s.Path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create signal file: %w", err)
	}

	return nil
}

// Remove deletes the signal file from disk
func (s *Signal) Remove() error {
	if err := os.Remove(s.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove signal file: %w", err)
	}
	return nil
}

// LoadFromFile creates a Signal instance from an existing file
func LoadSignalFromFile(path string) (*Signal, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("signal file does not exist: %s", path)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stat signal file: %w", err)
	}

	// Extract name from filename
	filename := filepath.Base(path)
	if !strings.HasSuffix(filename, ".signal") {
		return nil, fmt.Errorf("invalid signal filename format: %s", filename)
	}

	name := filename[:len(filename)-7] // Remove .signal extension

	return &Signal{
		Name:      name,
		CreatedAt: info.ModTime(),
		Path:      path,
	}, nil
}

// SignalInfo represents summary information about a signal
type SignalInfo struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"` // Formatted timestamp
	Path      string `json:"path"`
}

// ToInfo converts a Signal to SignalInfo
func (s *Signal) ToInfo() SignalInfo {
	return SignalInfo{
		Name:      s.Name,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		Path:      s.Path,
	}
}
