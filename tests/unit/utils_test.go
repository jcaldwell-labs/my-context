package unit

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "forward slashes on unix",
			input:    "/home/user/file.txt",
			expected: "/home/user/file.txt",
		},
		{
			name:     "backslashes on unix",
			input:    "C:\\Users\\file.txt",
			expected: "C:/Users/file.txt",
		},
		{
			name:     "mixed slashes",
			input:    "/home/user\\subfolder/file.txt",
			expected: "/home/user/subfolder/file.txt",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.NormalizePath(tt.input)
			
			if runtime.GOOS == "windows" {
				// On Windows, forward slashes should become backslashes
				expectedWindows := tt.input
				expectedWindows = filepath.FromSlash(tt.input)
				assert.Equal(t, expectedWindows, result)
			} else {
				// On Unix, backslashes should become forward slashes
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create simple directory",
			path:    filepath.Join(t.TempDir(), "testdir"),
			wantErr: false,
		},
		{
			name:    "create nested directories",
			path:    filepath.Join(t.TempDir(), "parent", "child", "grandchild"),
			wantErr: false,
		},
		{
			name:    "directory already exists",
			path:    t.TempDir(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.EnsureDir(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, utils.IsDir(tt.path), "Directory should exist")
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	err := os.WriteFile(existingFile, []byte("test"), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     existingFile,
			expected: true,
		},
		{
			name:     "non-existent file",
			path:     filepath.Join(tempDir, "notexists.txt"),
			expected: false,
		},
		{
			name:     "existing directory",
			path:     tempDir,
			expected: true,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FileExists(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDir(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "file.txt")
	err := os.WriteFile(testFile, []byte("test"), 0o600)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "is directory",
			path:     tempDir,
			expected: true,
		},
		{
			name:     "is file",
			path:     testFile,
			expected: false,
		},
		{
			name:     "non-existent path",
			path:     filepath.Join(tempDir, "notexists"),
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsDir(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetModTime(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "file.txt")
	
	beforeWrite := time.Now()
	time.Sleep(10 * time.Millisecond) // Ensure different timestamp
	
	err := os.WriteFile(testFile, []byte("test"), 0o600)
	require.NoError(t, err)
	
	time.Sleep(10 * time.Millisecond) // Ensure different timestamp
	afterWrite := time.Now()

	t.Run("existing file", func(t *testing.T) {
		modTime, err := utils.GetModTime(testFile)
		assert.NoError(t, err)
		assert.True(t, modTime.After(beforeWrite), "Mod time should be after write started")
		assert.True(t, modTime.Before(afterWrite), "Mod time should be before write completed")
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := utils.GetModTime(filepath.Join(tempDir, "notexists.txt"))
		assert.Error(t, err)
	})

	t.Run("directory", func(t *testing.T) {
		modTime, err := utils.GetModTime(tempDir)
		assert.NoError(t, err)
		assert.False(t, modTime.IsZero())
	})
}

func TestHasFileChanged(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "file.txt")
	
	// Initial write
	err := os.WriteFile(testFile, []byte("initial"), 0o600)
	require.NoError(t, err)
	
	initialModTime, err := utils.GetModTime(testFile)
	require.NoError(t, err)
	
	pastTime := initialModTime.Add(-1 * time.Hour)
	futureTime := initialModTime.Add(1 * time.Hour)

	t.Run("changed since past time", func(t *testing.T) {
		changed, err := utils.HasFileChanged(testFile, pastTime)
		assert.NoError(t, err)
		assert.True(t, changed, "File should be changed compared to past time")
	})

	t.Run("not changed since future time", func(t *testing.T) {
		changed, err := utils.HasFileChanged(testFile, futureTime)
		assert.NoError(t, err)
		assert.False(t, changed, "File should not be changed compared to future time")
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := utils.HasFileChanged(filepath.Join(tempDir, "notexists.txt"), pastTime)
		assert.Error(t, err)
	})
}

func TestListFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files and directories
	files := []string{
		"test1.txt",
		"test2.txt",
		"data.json",
		"readme.md",
	}
	
	for _, file := range files {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("content"), 0o600)
		require.NoError(t, err)
	}
	
	// Create a subdirectory (should be excluded from results)
	subdir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subdir, 0o755)
	require.NoError(t, err)

	tests := []struct {
		name          string
		pattern       string
		expectedCount int
		expectedFiles []string
	}{
		{
			name:          "all txt files",
			pattern:       "*.txt",
			expectedCount: 2,
			expectedFiles: []string{"test1.txt", "test2.txt"},
		},
		{
			name:          "all json files",
			pattern:       "*.json",
			expectedCount: 1,
			expectedFiles: []string{"data.json"},
		},
		{
			name:          "all files",
			pattern:       "*",
			expectedCount: 4,
			expectedFiles: files,
		},
		{
			name:          "no matches",
			pattern:       "*.yaml",
			expectedCount: 0,
			expectedFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ListFiles(tempDir, tt.pattern)
			assert.NoError(t, err)
			assert.Len(t, result, tt.expectedCount)
			
			// Check that all expected files are in results (by basename)
			for _, expectedFile := range tt.expectedFiles {
				found := false
				for _, resultFile := range result {
					if filepath.Base(resultFile) == expectedFile {
						found = true
						break
					}
				}
				if tt.expectedCount > 0 {
					assert.True(t, found, "Expected file %s not found in results", expectedFile)
				}
			}
		})
	}
}

func TestSafeWriteFile(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("write new file", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "new.txt")
		content := []byte("test content")
		
		err := utils.SafeWriteFile(testFile, content)
		assert.NoError(t, err)
		
		// Verify file exists and has correct content
		assert.True(t, utils.FileExists(testFile))
		readContent, err := os.ReadFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "existing.txt")
		
		// Write initial content
		err := os.WriteFile(testFile, []byte("old content"), 0o600)
		require.NoError(t, err)
		
		// Overwrite with new content
		newContent := []byte("new content")
		err = utils.SafeWriteFile(testFile, newContent)
		assert.NoError(t, err)
		
		// Verify new content
		readContent, err := os.ReadFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, newContent, readContent)
	})

	t.Run("create nested directory", func(t *testing.T) {
		nestedFile := filepath.Join(tempDir, "nested", "subdir", "file.txt")
		content := []byte("nested content")
		
		err := utils.SafeWriteFile(nestedFile, content)
		assert.NoError(t, err)
		
		// Verify file exists
		assert.True(t, utils.FileExists(nestedFile))
		readContent, err := os.ReadFile(nestedFile)
		assert.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("temp file cleanup on failure", func(t *testing.T) {
		// This test ensures temp files are cleaned up, but it's hard to force a rename failure
		// in a portable way. We'll just verify the normal case doesn't leave temp files.
		testFile := filepath.Join(tempDir, "cleanup.txt")
		content := []byte("test")
		
		err := utils.SafeWriteFile(testFile, content)
		assert.NoError(t, err)
		
		// Check that no .tmp files are left behind
		tempFiles, err := filepath.Glob(filepath.Join(tempDir, ".*.tmp"))
		assert.NoError(t, err)
		assert.Empty(t, tempFiles, "No temporary files should remain")
	})
}

func TestReadFileSafe(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := []byte("test content")
	err := os.WriteFile(testFile, content, 0o600)
	require.NoError(t, err)

	t.Run("read existing file", func(t *testing.T) {
		data, err := utils.ReadFileSafe(testFile)
		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		_, err := utils.ReadFileSafe(filepath.Join(tempDir, "notexists.txt"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")
	})

	t.Run("read empty file", func(t *testing.T) {
		emptyFile := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(emptyFile, []byte(""), 0o600)
		require.NoError(t, err)
		
		data, err := utils.ReadFileSafe(emptyFile)
		assert.NoError(t, err)
		assert.Empty(t, data)
	})
}

// TestIsValidPostgresIdentifier is already defined in partition_test.go
// Removed duplicate test to avoid build errors
