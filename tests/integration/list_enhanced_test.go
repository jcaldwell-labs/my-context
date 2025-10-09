package integration

import (
	"strings"
	"testing"
)

// TestListDefaultLimit tests default list shows 10 most recent contexts
func TestListDefaultLimit(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 15 contexts
	for i := 1; i <= 15; i++ {
		createTestContext(t, testDir, fmt.Sprintf("context-%d", i))
		runCommand("stop")
	}

	// Execute: List without flags
	output, _ := runCommandWithOutput("list")

	// Verify: Shows 10 contexts and truncation message
	if !strings.Contains(output, "Showing 10 of 15") {
		t.Error("Expected truncation message for 10 of 15 contexts")
	}
}

// TestListAllFlag tests --all flag shows all contexts
func TestListAllFlag(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 15 contexts
	for i := 1; i <= 15; i++ {
		createTestContext(t, testDir, fmt.Sprintf("context-%d", i))
		runCommand("stop")
	}

	// Execute: List with --all
	output, _ := runCommandWithOutput("list", "--all")

	// Verify: All 15 contexts shown
	for i := 1; i <= 15; i++ {
		contextName := fmt.Sprintf("context-%d", i)
		if !strings.Contains(output, contextName) {
			t.Errorf("Expected to find %s in --all output", contextName)
		}
	}
}

// TestListCustomLimit tests --limit flag
func TestListCustomLimit(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 20 contexts
	for i := 1; i <= 20; i++ {
		createTestContext(t, testDir, fmt.Sprintf("context-%d", i))
		runCommand("stop")
	}

	// Execute: List with --limit 5
	output, _ := runCommandWithOutput("list", "--limit", "5")

	// Verify: Shows 5 contexts
	if !strings.Contains(output, "Showing 5 of 20") {
		t.Error("Expected to show 5 of 20 contexts")
	}
}

// TestListSearch tests --search flag for substring matching
func TestListSearch(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create contexts with different names
	createTestContext(t, testDir, "ps-cli: Phase 1")
	runCommand("stop")
	createTestContext(t, testDir, "ps-cli: Phase 2")
	runCommand("stop")
	createTestContext(t, testDir, "garden: Planning")
	runCommand("stop")

	// Execute: Search for "Phase"
	output, _ := runCommandWithOutput("list", "--search", "Phase")

	// Verify: Only Phase contexts shown
	if !strings.Contains(output, "Phase 1") || !strings.Contains(output, "Phase 2") {
		t.Error("Expected to find Phase contexts")
	}
	if strings.Contains(output, "garden") {
		t.Error("Should not show garden context in Phase search")
	}
}

// TestListArchivedFlag tests --archived shows only archived contexts
func TestListArchivedFlag(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create and archive one context
	archived := "archived-context"
	createTestContext(t, testDir, archived)
	runCommand("stop")
	runCommand("archive", archived)

	// Create normal context
	normal := "normal-context"
	createTestContext(t, testDir, normal)
	runCommand("stop")

	// Execute: List --archived
	output, _ := runCommandWithOutput("list", "--archived")

	// Verify: Only archived context shown
	if !strings.Contains(output, archived) {
		t.Error("Expected archived context in --archived output")
	}
	if strings.Contains(output, normal) {
		t.Error("Normal context should not appear in --archived output")
	}
}

// TestListActiveOnly tests --active-only flag
func TestListActiveOnly(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create stopped contexts
	createTestContext(t, testDir, "stopped-1")
	runCommand("stop")

	// Create active context
	active := "active-context"
	createTestContext(t, testDir, active)

	// Execute: List --active-only
	output, _ := runCommandWithOutput("list", "--active-only")

	// Verify: Only active context shown
	if !strings.Contains(output, active) {
		t.Error("Expected active context in output")
	}
	if strings.Contains(output, "stopped-1") {
		t.Error("Stopped contexts should not appear with --active-only")
	}
}

// TestListCombinedFilters tests combining multiple flags
func TestListCombinedFilters(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create project contexts
	for i := 1; i <= 10; i++ {
		createTestContext(t, testDir, fmt.Sprintf("ps-cli: Phase %d", i))
		runCommand("stop")
	}

	// Execute: List with --project, --search, and --limit
	output, _ := runCommandWithOutput("list", "--project", "ps-cli", "--search", "Phase", "--limit", "3")

	// Verify: Combined filters applied (AND logic)
	lines := strings.Split(output, "\n")
	contextLines := 0
	for _, line := range lines {
		if strings.Contains(line, "ps-cli") && strings.Contains(line, "Phase") {
			contextLines++
		}
	}

	if contextLines > 3 {
		t.Errorf("Expected at most 3 contexts, got %d", contextLines)
	}
}

// Helper (minimal implementation for fmt)
import "fmt"
