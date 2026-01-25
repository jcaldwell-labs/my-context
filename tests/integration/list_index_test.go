package integration

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestListShowsIndexNumbers tests that list output shows index numbers
func TestListShowsIndexNumbers(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 3 contexts (createTestContext stops context after creation)
	createTestContext(t, "ctx1")
	createTestContext(t, "ctx2")
	createTestContext(t, "ctx3")

	// Execute: List command
	output, err := runCommandWithOutput("list")
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	// Verify: Index numbers are shown
	if !strings.Contains(output, "1. ○ ctx3") {
		t.Error("Expected index 1 for ctx3")
	}
	if !strings.Contains(output, "2. ○ ctx2") {
		t.Error("Expected index 2 for ctx2")
	}
	if !strings.Contains(output, "3. ○ ctx1") {
		t.Error("Expected index 3 for ctx1")
	}
}

// TestListJSONContainsIndex tests that JSON output includes index field
func TestListJSONContainsIndex(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 2 contexts (createTestContext stops context after creation)
	createTestContext(t, "ctx1")
	createTestContext(t, "ctx2")

	// Execute: List with JSON output
	output, err := runCommandWithOutput("list", "--json")
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify: Index field exists in contexts
	data := result["data"].(map[string]interface{})
	contexts := data["data"].(map[string]interface{})["contexts"].([]interface{})

	if len(contexts) < 2 {
		t.Fatalf("Expected at least 2 contexts, got %d", len(contexts))
	}

	// First context should have index 1
	firstContext := contexts[0].(map[string]interface{})
	if index, ok := firstContext["index"].(float64); !ok || int(index) != 1 {
		t.Errorf("Expected first context to have index 1, got %v", firstContext["index"])
	}

	// Second context should have index 2
	secondContext := contexts[1].(map[string]interface{})
	if index, ok := secondContext["index"].(float64); !ok || int(index) != 2 {
		t.Errorf("Expected second context to have index 2, got %v", secondContext["index"])
	}
}

// TestListIndexWithActiveContext tests index display with active context
func TestListIndexWithActiveContext(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create 2 stopped contexts
	runCommand("start", "ctx1")
	runCommand("stop")
	runCommand("start", "ctx2")
	runCommand("stop")

	// Create and keep one active
	runCommand("start", "ctx3")

	// Execute: List command
	output, err := runCommandWithOutput("list")
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	// Verify: Active context has correct index
	if !strings.Contains(output, "1. ● ctx3") {
		t.Errorf("Expected index 1 with active indicator for ctx3, got: %s", output)
	}
	if !strings.Contains(output, "2. ○ ctx2") {
		t.Errorf("Expected index 2 for ctx2, got: %s", output)
	}
	if !strings.Contains(output, "3. ○ ctx1") {
		t.Errorf("Expected index 3 for ctx1, got: %s", output)
	}
}
