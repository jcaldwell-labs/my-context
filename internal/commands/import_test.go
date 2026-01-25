package commands

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMarkdownNotes_WithTimestamps(t *testing.T) {
	markdown := `# Meeting Notes

## 2026-01-13 10:00
- Discussed payment integration
- Action item: Review PR #123

## 2026-01-13 11:30
- Follow-up on database migration
- Performance issues
`

	notes := parseMarkdownNotes(markdown, true)

	require.Len(t, notes, 2, "Should parse 2 notes")

	// Check first note
	assert.Equal(t, "2026-01-13 10:00:00 +0000 UTC", notes[0].Timestamp.UTC().String())
	assert.Contains(t, notes[0].TextContent, "Discussed payment integration")
	assert.Contains(t, notes[0].TextContent, "Action item: Review PR #123")

	// Check second note
	assert.Equal(t, "2026-01-13 11:30:00 +0000 UTC", notes[1].Timestamp.UTC().String())
	assert.Contains(t, notes[1].TextContent, "Follow-up on database migration")
	assert.Contains(t, notes[1].TextContent, "Performance issues")
}

func TestParseMarkdownNotes_WithoutPreserveTimestamps(t *testing.T) {
	markdown := `## 2026-01-13 10:00
- Note with timestamp header

## 2026-01-13 11:30
- Another note
`

	before := time.Now()
	notes := parseMarkdownNotes(markdown, false)
	after := time.Now()

	require.Len(t, notes, 2, "Should parse 2 notes")

	// Timestamps should be current time, not from markdown
	assert.True(t, notes[0].Timestamp.After(before) || notes[0].Timestamp.Equal(before))
	assert.True(t, notes[0].Timestamp.Before(after) || notes[0].Timestamp.Equal(after))

	assert.True(t, notes[1].Timestamp.After(before) || notes[1].Timestamp.Equal(before))
	assert.True(t, notes[1].Timestamp.Before(after) || notes[1].Timestamp.Equal(after))
}

func TestParseMarkdownNotes_PlainText(t *testing.T) {
	plainText := `This is a simple note
Another line of text
- Bullet point 1
- Bullet point 2

Yet another note entry`

	notes := parseMarkdownNotes(plainText, false)

	require.Len(t, notes, 1, "Should parse 1 note from plain text")
	assert.Contains(t, notes[0].TextContent, "This is a simple note")
	assert.Contains(t, notes[0].TextContent, "Another line of text")
	assert.Contains(t, notes[0].TextContent, "Bullet point 1")
}

func TestParseMarkdownNotes_EmptyInput(t *testing.T) {
	notes := parseMarkdownNotes("", false)
	assert.Len(t, notes, 0, "Should return empty slice for empty input")
}

func TestParseMarkdownNotes_OnlyWhitespace(t *testing.T) {
	notes := parseMarkdownNotes("   \n\n   \n", false)
	assert.Len(t, notes, 0, "Should return empty slice for whitespace-only input")
}

func TestParseMarkdownNotes_MultipleTimestampFormats(t *testing.T) {
	markdown := `## 2026-01-13 10:00
- Note with short timestamp

## 2026-01-13 11:30:45
- Note with seconds

## 2026-01-13T14:00:00
- Note with ISO format

## 2026-01-13T15:30
- Note with ISO format without seconds
`

	notes := parseMarkdownNotes(markdown, true)

	require.Len(t, notes, 4, "Should parse all timestamp formats")

	assert.Equal(t, "2026-01-13 10:00:00 +0000 UTC", notes[0].Timestamp.UTC().String())
	assert.Equal(t, "2026-01-13 11:30:45 +0000 UTC", notes[1].Timestamp.UTC().String())
	assert.Equal(t, "2026-01-13 14:00:00 +0000 UTC", notes[2].Timestamp.UTC().String())
	assert.Equal(t, "2026-01-13 15:30:00 +0000 UTC", notes[3].Timestamp.UTC().String())
}

func TestParseMarkdownNotes_IgnoresNonTimestampHeaders(t *testing.T) {
	markdown := `# Main Title

## Introduction
This should be ignored

## 2026-01-13 10:00
- This should be captured

### Subheading
Also ignored

## Summary
Also ignored
`

	notes := parseMarkdownNotes(markdown, true)

	require.Len(t, notes, 1, "Should only parse notes under timestamp headers")
	assert.Contains(t, notes[0].TextContent, "This should be captured")
}

func TestParseMarkdownNotes_HandlesMultilineNotes(t *testing.T) {
	markdown := `## 2026-01-13 10:00
- First line
Second line of same note
Third line

- New bullet point
More text here

Another paragraph
`

	notes := parseMarkdownNotes(markdown, true)

	require.Len(t, notes, 1, "Should group all text under one timestamp")
	assert.Contains(t, notes[0].TextContent, "First line")
	assert.Contains(t, notes[0].TextContent, "Second line")
	assert.Contains(t, notes[0].TextContent, "New bullet point")
}

func TestParseMarkdownNotes_RemovesBulletMarkers(t *testing.T) {
	markdown := `## 2026-01-13 10:00
- Dash bullet
* Asterisk bullet
+ Plus bullet
Regular text
`

	notes := parseMarkdownNotes(markdown, true)

	require.Len(t, notes, 1)
	// Check that bullet markers are removed
	assert.Contains(t, notes[0].TextContent, "Dash bullet")
	assert.Contains(t, notes[0].TextContent, "Asterisk bullet")
	assert.Contains(t, notes[0].TextContent, "Plus bullet")
	assert.Contains(t, notes[0].TextContent, "Regular text")
}
