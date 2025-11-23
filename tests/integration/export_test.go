package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExportSingleContextDefaultPath tests exporting a context to the default path
func TestExportSingleContextDefaultPath(t *testing.T) {
	// Setup: Create a test context
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "test-export-context"
	createTestContext(t, testDir, contextName)

	// Execute: Export the context
	outputPath := filepath.Join(testDir, contextName+".md")
	err := runCommand("export", contextName)
	if err != nil {
		t.Fatalf("Export command failed: %v", err)
	}

	// Verify: Output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected export file at %s, but it doesn't exist", outputPath)
	}

	// Verify: File contains markdown content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	if !strings.Contains(string(content), "# Context:") {
		t.Errorf("Export file missing markdown header")
	}
	if !strings.Contains(string(content), contextName) {
		t.Errorf("Export file missing context name")
	}
}

// TestExportWithCustomPath tests exporting to a custom path using --to flag
func TestExportWithCustomPath(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "test-custom-path"
	createTestContext(t, testDir, contextName)

	// Execute: Export with custom path
	customPath := filepath.Join(testDir, "exports", "my-export.md")
	err := runCommand("export", contextName, "--to", customPath)
	if err != nil {
		t.Fatalf("Export with --to failed: %v", err)
	}

	// Verify: File exists at custom path
	if _, err := os.Stat(customPath); os.IsNotExist(err) {
		t.Errorf("Expected export file at %s", customPath)
	}

	// Verify: Parent directory was created
	parentDir := filepath.Dir(customPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		t.Errorf("Expected parent directory %s to be created", parentDir)
	}
}

// TestExportAllFlag tests exporting all contexts with --all flag
func TestExportAllFlag(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create multiple contexts
	contexts := []string{"context-1", "context-2", "context-3"}
	for _, name := range contexts {
		createTestContext(t, testDir, name)
	}

	// Execute: Export all
	outputDir := filepath.Join(testDir, "all-exports")
	err := runCommand("export", "--all", "--to", outputDir)
	if err != nil {
		t.Fatalf("Export --all failed: %v", err)
	}

	// Verify: All context files exist
	for _, name := range contexts {
		expectedPath := filepath.Join(outputDir, name+".md")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected export file for %s at %s", name, expectedPath)
		}
	}
}

// TestExportNonExistentContext tests error handling for non-existent context
func TestExportNonExistentContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Execute: Try to export non-existent context
	err := runCommand("export", "non-existent-context")

	// Verify: Command fails with appropriate error
	if err == nil {
		t.Fatal("Expected error for non-existent context, got nil")
	}

	expectedError := "not found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got: %v", expectedError, err)
	}
}

// TestExportMarkdownFormat tests the markdown format and content structure
func TestExportMarkdownFormat(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "format-test"
	createTestContext(t, testDir, contextName)

	// Add some notes and files to the context
	runCommand("note", "Test note 1")
	runCommand("note", "Test note 2")

	// Execute: Export
	outputPath := filepath.Join(testDir, contextName+".md")
	err := runCommand("export", contextName)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify: Markdown structure
	content, _ := os.ReadFile(outputPath)
	markdown := string(content)

	requiredSections := []string{
		"# Context:",
		"**Started**:",
		"## Notes",
		"## Files",
		"## Activity",
	}

	for _, section := range requiredSections {
		if !strings.Contains(markdown, section) {
			t.Errorf("Markdown missing required section: %s", section)
		}
	}

	// Verify: Notes appear in export
	if !strings.Contains(markdown, "Test note 1") {
		t.Errorf("Export missing note content")
	}
}

// TestExportJSONOutput tests JSON format output with --json flag
func TestExportJSONOutput(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "json-test"
	createTestContext(t, testDir, contextName)

	// Execute: Export as JSON
	outputPath := filepath.Join(testDir, contextName+".json")
	err := runCommand("export", contextName, "--json", "--to", outputPath)
	if err != nil {
		t.Fatalf("JSON export failed: %v", err)
	}

	// Verify: File is valid JSON
	content, _ := os.ReadFile(outputPath)
	var exportData map[string]interface{}
	if err := json.Unmarshal(content, &exportData); err != nil {
		t.Fatalf("Export file is not valid JSON: %v", err)
	}

	// Verify: JSON structure
	if _, ok := exportData["context"]; !ok {
		t.Error("JSON missing 'context' field")
	}
	if _, ok := exportData["notes"]; !ok {
		t.Error("JSON missing 'notes' field")
	}
	if _, ok := exportData["files"]; !ok {
		t.Error("JSON missing 'files' field")
	}
}

// Helper functions moved to helpers_test.go
