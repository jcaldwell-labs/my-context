package integration

import (
	"strings"
	"testing"
)

// TestListProjectFilter tests --project flag filters contexts
func TestListProjectFilter(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create contexts with different projects
	contexts := []string{
		"ps-cli: Phase 1",
		"ps-cli: Phase 2",
		"garden: Planning",
		"garden: Implementation",
	}
	for _, ctx := range contexts {
		createTestContext(t, testDir, ctx)
		runCommand("stop")
	}

	// Execute: Filter by ps-cli project
	output, _ := runCommandWithOutput("list", "--project", "ps-cli")

	// Verify: Only ps-cli contexts shown
	if !strings.Contains(output, "ps-cli: Phase 1") {
		t.Error("Expected ps-cli Phase 1 in output")
	}
	if !strings.Contains(output, "ps-cli: Phase 2") {
		t.Error("Expected ps-cli Phase 2 in output")
	}
	if strings.Contains(output, "garden") {
		t.Error("Garden contexts should not appear in ps-cli filter")
	}
}

// TestStartWithProject tests --project flag for start command
func TestStartWithProject(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Execute: Start with --project flag
	err := runCommand("start", "Phase 3", "--project", "ps-cli")
	if err != nil {
		t.Fatalf("Start with --project failed: %v", err)
	}

	// Verify: Context created with project prefix
	output, _ := runCommandWithOutput("show")
	expectedName := "ps-cli: Phase 3"
	if !strings.Contains(output, expectedName) {
		t.Errorf("Expected context name to be %q", expectedName)
	}

	runCommand("stop")
}

// TestProjectExtractionMultipleColons tests handling of multiple colons
func TestProjectExtractionMultipleColons(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create context with multiple colons
	contextName := "project: Phase 1: Subphase A"
	createTestContext(t, testDir, contextName)
	runCommand("stop")

	// Execute: Filter by project (should extract "project")
	output, _ := runCommandWithOutput("list", "--project", "project")

	// Verify: Context appears in filtered list
	if !strings.Contains(output, contextName) {
		t.Error("Context with multiple colons should be found by first part")
	}
}

// TestProjectExtractionNoColon tests contexts without colons
func TestProjectExtractionNoColon(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create standalone context (no colon)
	standalone := "StandaloneContext"
	createTestContext(t, testDir, standalone)
	runCommand("stop")

	// Execute: Filter by full name
	output, _ := runCommandWithOutput("list", "--project", standalone)

	// Verify: Context found (full name = project name)
	if !strings.Contains(output, standalone) {
		t.Error("Standalone context should match its full name as project")
	}
}

// TestProjectExtractionWhitespace tests whitespace handling
func TestProjectExtractionWhitespace(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create context with whitespace around colon
	contextName := " project  :  description "
	createTestContext(t, testDir, contextName)
	runCommand("stop")

	// Execute: Filter by trimmed project name
	output, _ := runCommandWithOutput("list", "--project", "project")

	// Verify: Whitespace trimmed, context found
	if !strings.Contains(output, "project") {
		t.Error("Project extraction should trim whitespace")
	}
}

// TestProjectCaseInsensitive tests case-insensitive matching (FR-004.5)
func TestProjectCaseInsensitive(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create contexts with mixed case
	contexts := []string{
		"ps-cli: Phase 1",
		"Ps-Cli: Phase 2",
	}
	for _, ctx := range contexts {
		createTestContext(t, testDir, ctx)
		runCommand("stop")
	}

	// Test different case variations
	cases := []string{"ps-cli", "PS-CLI", "Ps-Cli", "pS-cLi"}
	
	for _, searchCase := range cases {
		t.Run("search_"+searchCase, func(t *testing.T) {
			output, _ := runCommandWithOutput("list", "--project", searchCase)

			// Verify: Both contexts found regardless of case
			if !strings.Contains(strings.ToLower(output), "phase 1") {
				t.Errorf("Case-insensitive search for %q should find Phase 1", searchCase)
			}
			if !strings.Contains(strings.ToLower(output), "phase 2") {
				t.Errorf("Case-insensitive search for %q should find Phase 2", searchCase)
			}
		})
	}
}
