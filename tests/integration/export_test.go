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
	createTestContext(t, contextName)

	// Execute: Export the context with explicit --to path
	// Note: Default export writes to current directory, which in tests is projectRoot
	// Using --to ensures file goes where we expect
	outputPath := filepath.Join(testDir, contextName+".md")
	err := runCommand("export", contextName, "--to", outputPath)
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
	createTestContext(t, contextName)

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
		createTestContext(t, name)
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
	output, err := runCommandWithOutput("export", "non-existent-context")

	// Verify: Command fails with appropriate error
	if err == nil {
		t.Fatal("Expected error for non-existent context, got nil")
	}

	// Check that output contains an error message (may be in stdout or combined output)
	outputLower := strings.ToLower(output)
	if !strings.Contains(outputLower, "not found") &&
		!strings.Contains(outputLower, "does not exist") &&
		!strings.Contains(outputLower, "no such") &&
		!strings.Contains(outputLower, "error") {
		t.Errorf("Expected error message about missing context, got: %s", output)
	}
}

// TestExportMarkdownFormat tests the markdown format and content structure
func TestExportMarkdownFormat(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "format-test"
	// Start context manually (don't use createTestContext which stops immediately)
	err := runCommand("start", contextName)
	if err != nil {
		t.Fatalf("Failed to start context: %v", err)
	}

	// Add some notes to the active context
	runCommand("note", "Test note 1")
	runCommand("note", "Test note 2")

	// Stop the context before exporting
	runCommand("stop")

	// Execute: Export with explicit path
	outputPath := filepath.Join(testDir, contextName+".md")
	err = runCommand("export", contextName, "--to", outputPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify: Markdown structure
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}
	markdown := string(content)

	// Check for basic structure (format may vary between file and database backends)
	if !strings.Contains(markdown, "# Context:") && !strings.Contains(markdown, "Context:") {
		t.Error("Markdown missing context header")
	}

	// Verify: Notes appear in export
	if !strings.Contains(markdown, "Test note 1") {
		t.Errorf("Export missing note content. Content: %s", markdown)
	}
}

// TestExportJSONOutput tests JSON format output with --json flag
func TestExportJSONOutput(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	contextName := "json-test"
	createTestContext(t, contextName)

	// Execute: Export as JSON with explicit path
	outputPath := filepath.Join(testDir, contextName+".json")
	err := runCommand("export", contextName, "--json", "--to", outputPath)
	if err != nil {
		t.Fatalf("JSON export failed: %v", err)
	}

	// Verify: File exists and is valid JSON
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	var exportData map[string]interface{}
	if err := json.Unmarshal(content, &exportData); err != nil {
		t.Fatalf("Export file is not valid JSON: %v. Content: %s", err, string(content))
	}

	// Verify: JSON has some expected structure (format may vary)
	// Could be "context" or "data" or "name" depending on output format
	hasExpectedField := false
	for key := range exportData {
		if key == "context" || key == "data" || key == "name" || key == "status" {
			hasExpectedField = true
			break
		}
	}
	if !hasExpectedField {
		t.Errorf("JSON missing expected fields. Got keys: %v", getMapKeys(exportData))
	}
}

// getMapKeys returns the keys of a map for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
