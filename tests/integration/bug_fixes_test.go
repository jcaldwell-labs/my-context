package integration

import (
	"strings"
	"testing"
)

// TestNoteDollarCharacterPreserved tests $ character preservation in notes (FR-009.1)
func TestNoteDollarCharacterPreserved(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create context and add note with $ character
	createTestContext(t, "dollar-test")
	noteContent := "Budget: $500-800"
	err := runCommand("note", noteContent)
	if err != nil {
		t.Fatalf("Failed to add note with $ character: %v", err)
	}

	// Execute: List to see note display
	output, _ := runCommandWithOutput("show")

	// Verify: $ character preserved
	if !strings.Contains(output, "$500-800") {
		t.Error("Dollar sign should be preserved in note display")
	}
	if !strings.Contains(output, noteContent) {
		t.Errorf("Full note content should be preserved: expected %q", noteContent)
	}

	runCommand("stop")
}

// TestHistoryDisplaysNoneInsteadOfNull tests "(none)" vs "NULL" display (FR-009.2)
func TestHistoryDisplaysNoneInsteadOfNull(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	// Create first context (no previous context)
	createTestContext(t, "first-context")
	runCommand("stop")

	// Execute: View history
	output, _ := runCommandWithOutput("history")

	// Verify: Shows "(none)" not "NULL" for empty fields
	if strings.Contains(output, "NULL") {
		t.Error("History should display '(none)' not 'NULL' for empty context fields")
	}
	if !strings.Contains(output, "(none)") {
		t.Error("Expected '(none)' for empty previous context field")
	}
}

// TestSpecialCharactersInNotes tests various special characters (FR-009.3)
func TestSpecialCharactersInNotes(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	createTestContext(t, "special-chars-test")

	// Test various special characters
	specialNotes := []string{
		"Price: $100",
		"Email: user@example.com",
		"Tag: #important",
		"Alert: Warning!",
		"Percentage: 50%",
		"Path: /home/user/file.txt",
		"Expression: value = 5 * (3 + 2)",
		"Quote: \"Hello World\"",
		"Apostrophe: It's working",
	}

	for _, note := range specialNotes {
		err := runCommand("note", note)
		if err != nil {
			t.Errorf("Failed to add note with special characters: %q, error: %v", note, err)
		}
	}

	// Verify all notes are preserved
	output, _ := runCommandWithOutput("show")

	for _, note := range specialNotes {
		if !strings.Contains(output, note) {
			t.Errorf("Note with special characters not preserved: %q", note)
		}
	}

	runCommand("stop")
}

// TestNoteWithBackslash tests backslash character handling
func TestNoteWithBackslash(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	createTestContext(t, "backslash-test")

	// Windows path with backslashes
	noteContent := `Path: C:\Users\username\file.txt`
	err := runCommand("note", noteContent)
	if err != nil {
		t.Fatalf("Failed to add note with backslashes: %v", err)
	}

	// Verify: Backslashes preserved
	output, _ := runCommandWithOutput("show")
	if !strings.Contains(output, `C:\Users\username\file.txt`) {
		t.Error("Backslashes should be preserved in notes")
	}

	runCommand("stop")
}

// TestNoteWithUnicode tests Unicode character handling
func TestNoteWithUnicode(t *testing.T) {
	testDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, testDir)

	createTestContext(t, "unicode-test")

	// Various Unicode characters
	unicodeNotes := []string{
		"Emoji: 🎉 ✅ 🚀",
		"Greek: α β γ δ",
		"Math: √ ∞ ≠ ≈",
		"Currency: € £ ¥",
	}

	for _, note := range unicodeNotes {
		runCommand("note", note)
	}

	// Verify: Unicode preserved
	output, _ := runCommandWithOutput("show")
	for _, note := range unicodeNotes {
		if !strings.Contains(output, note) {
			t.Errorf("Unicode note not preserved: %q", note)
		}
	}

	runCommand("stop")
}
