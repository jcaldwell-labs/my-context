package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// NormalizePath converts all path separators to the current OS separator
func NormalizePath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, "/", "\\")
	}
	return strings.ReplaceAll(path, "\\", "/")
}

// EnsureDir creates a directory and all necessary parents
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir checks if a path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetModTime returns the modification time of a file
func GetModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// HasFileChanged checks if a file has been modified since the given time
func HasFileChanged(path string, since time.Time) (bool, error) {
	modTime, err := GetModTime(path)
	if err != nil {
		return false, err
	}
	return modTime.After(since), nil
}

// ListFiles returns all files in a directory matching a pattern
func ListFiles(dir, pattern string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		if !IsDir(file) {
			result = append(result, file)
		}
	}

	return result, nil
}

// SafeWriteFile writes to a file atomically by writing to a temp file first
func SafeWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	tempFile := filepath.Join(dir, fmt.Sprintf(".%s.tmp", filepath.Base(path)))

	// Ensure directory exists
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to temp file
	if err := os.WriteFile(tempFile, data, 0o600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, path); err != nil {
		os.Remove(tempFile) // Clean up temp file on failure
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// ReadFileSafe reads a file with error handling
func ReadFileSafe(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}
