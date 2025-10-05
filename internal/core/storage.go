package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// GetContextHome returns the context home directory path
func GetContextHome() string {
	home := os.Getenv("MY_CONTEXT_HOME")
	if home == "" {
		userHome, _ := os.UserHomeDir()
		home = filepath.Join(userHome, ".my-context")
	}
	return home
}

// EnsureContextHome creates the context home directory if it doesn't exist
func EnsureContextHome() error {
	home := GetContextHome()
	return os.MkdirAll(home, 0700)
}

// CreateDir creates a directory with proper permissions
func CreateDir(path string) error {
	return os.MkdirAll(path, 0700)
}

// ReadJSON reads and unmarshals JSON from a file
func ReadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteJSON marshals and writes JSON to a file atomically
func WriteJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first
	tmpPath := path + ".tmp"
	err = os.WriteFile(tmpPath, data, 0600)
	if err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tmpPath, path)
}

// AppendLog appends a line to a log file
func AppendLog(path string, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line + "\n")
	return err
}

// ReadLog reads all lines from a log file
func ReadLog(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	// Remove empty last line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}

// NormalizePath converts a path to absolute and POSIX format for storage
func NormalizePath(path string) (string, error) {
	// Clean the path
	cleaned := filepath.Clean(path)

	// Convert to absolute
	absolute, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}

	// Convert to POSIX format (forward slashes)
	return filepath.ToSlash(absolute), nil
}

// DenormalizePath converts a stored POSIX path to OS-native format
func DenormalizePath(path string) string {
	return filepath.FromSlash(path)
}

// FileExists checks if a file or directory exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetContextDir returns the directory path for a context name
func GetContextDir(contextName string) string {
	return filepath.Join(GetContextHome(), contextName)
}

// SanitizeContextName sanitizes a context name by replacing invalid characters
func SanitizeContextName(name string) string {
	// Replace spaces with underscores
	sanitized := strings.ReplaceAll(name, " ", "_")
	// Remove any remaining path separators
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	return sanitized
}

// ListContextDirs returns all context directory names
func ListContextDirs() ([]string, error) {
	home := GetContextHome()

	entries, err := os.ReadDir(home)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// GetStateFilePath returns the path to the state.json file
func GetStateFilePath() string {
	return filepath.Join(GetContextHome(), "state.json")
}

// GetTransitionsLogPath returns the path to the transitions.log file
func GetTransitionsLogPath() string {
	return filepath.Join(GetContextHome(), "transitions.log")
}

// GetMetaJSONPath returns the path to a context's meta.json file
func GetMetaJSONPath(contextName string) string {
	return filepath.Join(GetContextDir(contextName), "meta.json")
}

// GetNotesLogPath returns the path to a context's notes.log file
func GetNotesLogPath(contextName string) string {
	return filepath.Join(GetContextDir(contextName), "notes.log")
}

// GetFilesLogPath returns the path to a context's files.log file
func GetFilesLogPath(contextName string) string {
	return filepath.Join(GetContextDir(contextName), "files.log")
}

// GetTouchLogPath returns the path to a context's touch.log file
func GetTouchLogPath(contextName string) string {
	return filepath.Join(GetContextDir(contextName), "touch.log")
}
