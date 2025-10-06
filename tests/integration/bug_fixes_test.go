package integration

import (
	"os"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNoteDollarCharacterPreserved tests $ character is not stripped
func TestNoteDollarCharacterPreserved(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context
	ctx, err := core.CreateContext("dollar-test")
	require.NoError(t, err)

	// Add note with dollar sign
	err = core.AddNote(ctx.Name, "Budget: $500-800 for equipment")
	require.NoError(t, err)

	// Read note back
	notes, err := core.GetNotes(ctx.Name)
	require.NoError(t, err)
	require.Len(t, notes, 1)

	// Verify $ character preserved
	assert.Equal(t, "Budget: $500-800 for equipment", notes[0].TextContent)
	assert.Contains(t, notes[0].TextContent, "$500")
	assert.Contains(t, notes[0].TextContent, "$")
}

// TestNoteSpecialCharactersPreserved tests various special characters
func TestNoteSpecialCharactersPreserved(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("special-chars-test")
	require.NoError(t, err)

	specialChars := []string{
		"Price: $100",
		"Email: user@example.com",
		"Tag: #important",
		"Percentage: 50%",
		"Attention: this is important!",
		"Math: 2 + 2 = 4",
		"Regex: ^[a-z]*$",
		"Template: ${var}",
	}

	// Add notes with special characters
	for _, text := range specialChars {
		err = core.AddNote(ctx.Name, text)
		require.NoError(t, err)
	}

	// Read all notes back
	notes, err := core.GetNotes(ctx.Name)
	require.NoError(t, err)
	require.Len(t, notes, len(specialChars))

	// Verify all special characters preserved
	for i, note := range notes {
		assert.Equal(t, specialChars[i], note.TextContent)
	}
}

// TestHistoryShowsNoneInsteadOfNull tests history output format
func TestHistoryShowsNoneInsteadOfNull(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create first context (no previous context)
	ctx1, err := core.CreateContext("first-context")
	require.NoError(t, err)

	// Stop it (no next context)
	err = core.StopContext(ctx1.Name)
	require.NoError(t, err)

	// Read transitions log
	transitions, err := core.GetTransitions()
	require.NoError(t, err)
	require.Greater(t, len(transitions), 0)

	// Find START transition (should have empty previous context)
	var startTransition *models.ContextTransition
	for i := range transitions {
		if transitions[i].TransitionType == "start" {
			startTransition = &transitions[i]
			break
		}
	}
	require.NotNil(t, startTransition)

	// Verify empty context is represented as "(none)" not "NULL"
	assert.NotEqual(t, "NULL", startTransition.PreviousContext)
	// Note: Implementation should use "(none)" or empty string, not "NULL"
}

// TestNoteDisplayInShow tests note display preserves special chars
func TestNoteDisplayInShow(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("display-test")
	require.NoError(t, err)

	core.AddNote(ctx.Name, "Cost: $1,234.56")

	// Get context for display
	loadedCtx, err := core.LoadContext(ctx.Name)
	require.NoError(t, err)

	notes, err := core.GetNotes(loadedCtx.Name)
	require.NoError(t, err)
	require.Len(t, notes, 1)

	// Verify display would show correct text
	assert.Contains(t, notes[0].TextContent, "$1,234.56")
}

// TestNoteStorageFormat tests that notes.log preserves special characters
func TestNoteStorageFormat(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("storage-test")
	require.NoError(t, err)

	testNote := "Budget: $500 | Cost: $800"
	err = core.AddNote(ctx.Name, testNote)
	require.NoError(t, err)

	// Read notes.log directly
	notesLogPath := core.GetNotesLogPath(ctx.Name)
	content, err := os.ReadFile(notesLogPath)
	require.NoError(t, err)

	// Verify dollar signs stored correctly
	assert.Contains(t, string(content), "$500")
	assert.Contains(t, string(content), "$800")
	assert.Contains(t, string(content), testNote)
}
package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExportSingleContextDefaultPath tests exporting a context to default path
func TestExportSingleContextDefaultPath(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create a test context
	ctx, err := core.CreateContext("test-export-context")
	require.NoError(t, err)

	err = core.AddNote(ctx.Name, "Test note 1")
	require.NoError(t, err)

	err = core.AddNote(ctx.Name, "Test note 2")
	require.NoError(t, err)

	err = core.StopContext(ctx.Name)
	require.NoError(t, err)

	// Export the context
	outputPath, err := core.ExportContext(ctx.Name, "")
	require.NoError(t, err)

	// Verify file created with sanitized name
	assert.FileExists(t, outputPath)
	assert.Contains(t, outputPath, "test-export-context.md")

	// Verify markdown content
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	markdown := string(content)
	assert.Contains(t, markdown, "# Context: test-export-context")
	assert.Contains(t, markdown, "Test note 1")
	assert.Contains(t, markdown, "Test note 2")
	assert.Contains(t, markdown, "## Notes")

	// Cleanup
	os.Remove(outputPath)
}

// TestExportWithCustomPath tests exporting to a specified path
func TestExportWithCustomPath(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("test-custom-path")
	require.NoError(t, err)

	err = core.AddNote(ctx.Name, "Custom export test")
	require.NoError(t, err)

	err = core.StopContext(ctx.Name)
	require.NoError(t, err)

	// Export to custom path
	customPath := filepath.Join(t.TempDir(), "exports", "custom-export.md")
	outputPath, err := core.ExportContext(ctx.Name, customPath)
	require.NoError(t, err)

	// Verify parent directory created
	assert.FileExists(t, outputPath)
	assert.Equal(t, customPath, outputPath)

	// Verify content
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Custom export test")
}

// TestExportAllContexts tests exporting all contexts at once
func TestExportAllContexts(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create multiple contexts
	_, err := core.CreateContext("context-1")
	require.NoError(t, err)
	core.StopContext("context-1")

	_, err = core.CreateContext("context-2")
	require.NoError(t, err)
	core.StopContext("context-2")

	_, err = core.CreateContext("context-3")
	require.NoError(t, err)

	// Export all to directory
	outputDir := t.TempDir()
	exportedFiles, err := core.ExportAllContexts(outputDir)
	require.NoError(t, err)

	// Verify all contexts exported
	assert.Len(t, exportedFiles, 3)

	for _, filePath := range exportedFiles {
		assert.FileExists(t, filePath)
		assert.Contains(t, filePath, outputDir)
	}
}

// TestExportNonExistentContext tests error handling for missing context
func TestExportNonExistentContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Attempt to export non-existent context
	_, err := core.ExportContext("does-not-exist", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestExportWithSpecialCharacters tests filename sanitization
func TestExportWithSpecialCharacters(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create context with special characters
	ctx, err := core.CreateContext("test: with / special \\ chars")
	require.NoError(t, err)

	err = core.StopContext(ctx.Name)
	require.NoError(t, err)

	// Export should sanitize filename
	outputPath, err := core.ExportContext(ctx.Name, "")
	require.NoError(t, err)

	assert.FileExists(t, outputPath)
	// Verify special chars replaced with safe alternatives
	assert.NotContains(t, filepath.Base(outputPath), "/")
	assert.NotContains(t, filepath.Base(outputPath), "\\")

	os.Remove(outputPath)
}

// TestExportMarkdownFormat tests markdown structure and formatting
func TestExportMarkdownFormat(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("format-test")
	require.NoError(t, err)

	// Add various data types
	core.AddNote(ctx.Name, "First note")
	core.AddNote(ctx.Name, "Second note with $special chars!")
	core.AddFile(ctx.Name, "/path/to/file.txt")
	core.AddTouch(ctx.Name)
	core.StopContext(ctx.Name)

	// Export
	outputPath, err := core.ExportContext(ctx.Name, "")
	require.NoError(t, err)
	defer os.Remove(outputPath)

	// Verify markdown structure
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	markdown := string(content)

	// Check required sections
	assert.Contains(t, markdown, "# Context: format-test")
	assert.Contains(t, markdown, "**Started**:")
	assert.Contains(t, markdown, "**Ended**:")
	assert.Contains(t, markdown, "**Duration**:")
	assert.Contains(t, markdown, "## Notes")
	assert.Contains(t, markdown, "## Associated Files")
	assert.Contains(t, markdown, "## Activity")

	// Check content preservation
	assert.Contains(t, markdown, "First note")
	assert.Contains(t, markdown, "Second note with $special chars!")
	assert.Contains(t, markdown, "/path/to/file.txt")
	assert.Contains(t, markdown, "touch events")
}

// TestExportActiveContext tests exporting an active context
func TestExportActiveContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("active-export")
	require.NoError(t, err)

	core.AddNote(ctx.Name, "Active context note")

	// Should be able to export active context
	outputPath, err := core.ExportContext(ctx.Name, "")
	require.NoError(t, err)
	defer os.Remove(outputPath)

	// Verify "Active" status in export
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Active")
}

// TestExportArchivedContext tests exporting an archived context
func TestExportArchivedContext(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	ctx, err := core.CreateContext("archived-export")
	require.NoError(t, err)

	core.AddNote(ctx.Name, "Archived context note")
	core.StopContext(ctx.Name)

	// Archive the context
	err = core.ArchiveContext(ctx.Name)
	require.NoError(t, err)

	// Should still be able to export archived context
	outputPath, err := core.ExportContext(ctx.Name, "")
	require.NoError(t, err)
	defer os.Remove(outputPath)

	// Verify export includes archived status
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "archived")
}

