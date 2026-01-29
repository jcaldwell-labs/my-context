package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestImportFromFile tests importing notes from a markdown file
func TestImportFromFile(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "test-import-context"
	err := runCommand("start", contextName)
	if err != nil {
		t.Fatalf("Failed to start context: %v", err)
	}

	// Create a test markdown file with timestamped notes
	markdownPath := filepath.Join(testDir, "test-notes.md")
	markdownContent := `# Meeting Notes

## 2026-01-13 10:00
- Discussed payment integration
- Action item: Review PR #123

## 2026-01-13 11:30
- Follow-up on database migration
`
	err = os.WriteFile(markdownPath, []byte(markdownContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test markdown file: %v", err)
	}

	// Execute: Import the markdown file with timestamp preservation
	output, err := runCommandWithOutput("import", markdownPath, "--preserve-timestamps")
	if err != nil {
		t.Logf("Import output: %s", output)
		t.Fatalf("Import command failed: %v", err)
	}

	// Verify: Context has the imported notes
	output, err = runCommandWithOutput("show")
	if err != nil {
		t.Fatalf("Show command failed: %v", err)
	}

	// Check that both notes were imported
	if !strings.Contains(output, "Discussed payment integration") {
		t.Errorf("First note not found in context")
	}
	if !strings.Contains(output, "Follow-up on database migration") {
		t.Errorf("Second note not found in context")
	}

	// Check that timestamps were preserved
	if !strings.Contains(output, "2026-01-13") {
		t.Errorf("Timestamps were not preserved")
	}
}

// TestImportPlainText tests importing plain text notes
func TestImportPlainText(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "test-plain-import"
	err := runCommand("start", contextName)
	if err != nil {
		t.Fatalf("Failed to start context: %v", err)
	}

	// Create a plain text file
	textPath := filepath.Join(testDir, "plain-notes.txt")
	textContent := `This is a simple note
Another line of text
- Bullet point 1
- Bullet point 2`
	err = os.WriteFile(textPath, []byte(textContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test text file: %v", err)
	}

	// Execute: Import without preserve-timestamps
	output, err := runCommandWithOutput("import", textPath)
	if err != nil {
		t.Logf("Import output: %s", output)
		t.Fatalf("Import command failed: %v", err)
	}

	// Verify: Notes were imported
	output, err = runCommandWithOutput("show")
	if err != nil {
		t.Fatalf("Show command failed: %v", err)
	}

	if !strings.Contains(output, "This is a simple note") {
		t.Errorf("Plain text note not found in context")
	}
	if !strings.Contains(output, "Bullet point 1") {
		t.Errorf("Bullet points not imported correctly")
	}
}

// TestImportWithoutActiveContext tests that import fails without an active context
func TestImportWithoutActiveContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create a test file but don't create a context
	textPath := filepath.Join(testDir, "test.txt")
	err := os.WriteFile(textPath, []byte("Test note"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Execute: Try to import without an active context
	err = runCommand("import", textPath)

	// Verify: Command should fail with appropriate error
	if err == nil {
		t.Errorf("Expected import to fail without active context")
	}
}

// TestImportNonexistentFile tests that import fails for non-existent files
func TestImportNonexistentFile(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "test-nonexistent"
	err := runCommand("start", contextName)
	if err != nil {
		t.Fatalf("Failed to start context: %v", err)
	}

	// Execute: Try to import a file that doesn't exist
	err = runCommand("import", "/nonexistent/file.md")

	// Verify: Command should fail
	if err == nil {
		t.Errorf("Expected import to fail for nonexistent file")
	}
}

// TestImportEmptyFile tests that import handles empty files gracefully
func TestImportEmptyFile(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "test-empty"
	err := runCommand("start", contextName)
	if err != nil {
		t.Fatalf("Failed to start context: %v", err)
	}

	// Create an empty file
	emptyPath := filepath.Join(testDir, "empty.txt")
	err = os.WriteFile(emptyPath, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Execute: Try to import empty file
	err = runCommand("import", emptyPath)

	// Verify: Command should fail with "no notes found"
	if err == nil {
		t.Errorf("Expected import to fail for empty file")
	}
}
