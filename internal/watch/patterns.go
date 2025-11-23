package watch

import (
	"fmt"
	"regexp"
	"strings"
)

// PatternMatcher handles pattern matching for notes
type PatternMatcher struct {
	pattern    *regexp.Regexp
	rawPattern string
}

// NewPatternMatcher creates a new pattern matcher
func NewPatternMatcher(pattern string) (*PatternMatcher, error) {
	if pattern == "" {
		return &PatternMatcher{rawPattern: pattern}, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	return &PatternMatcher{
		pattern:    regex,
		rawPattern: pattern,
	}, nil
}

// Matches checks if the given text matches the pattern
func (pm *PatternMatcher) Matches(text string) bool {
	if pm.pattern == nil {
		return true // Empty pattern matches everything
	}
	return pm.pattern.MatchString(text)
}

// IsEmpty returns true if no pattern is set
func (pm *PatternMatcher) IsEmpty() bool {
	return pm.rawPattern == ""
}

// String returns the string representation of the pattern
func (pm *PatternMatcher) String() string {
	return pm.rawPattern
}

// NotePatternMatcher specifically handles note content pattern matching
type NotePatternMatcher struct {
	*PatternMatcher
}

// NewNotePatternMatcher creates a pattern matcher for notes
func NewNotePatternMatcher(pattern string) (*NotePatternMatcher, error) {
	pm, err := NewPatternMatcher(pattern)
	if err != nil {
		return nil, err
	}
	return &NotePatternMatcher{PatternMatcher: pm}, nil
}

// MatchesNote checks if a note matches the pattern
// For notes, we typically want to match against the note content
func (npm *NotePatternMatcher) MatchesNote(noteContent string) bool {
	return npm.Matches(noteContent)
}

// MatchesNoteLines checks if any line in a multi-line note matches the pattern
func (npm *NotePatternMatcher) MatchesNoteLines(noteContent string) bool {
	if npm.IsEmpty() {
		return true
	}

	lines := strings.Split(noteContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && npm.Matches(line) {
			return true
		}
	}
	return false
}

// Common patterns for context monitoring
var (
	// PhasePattern matches notes about phase completion
	PhasePattern = regexp.MustCompile(`(?i)phase.*(?:complete|done|finished)`)

	// BugfixPattern matches notes about bug fixes
	BugfixPattern = regexp.MustCompile(`(?i)(?:bug|fix|issue).*`)

	// DeployPattern matches notes about deployments
	DeployPattern = regexp.MustCompile(`(?i)(?:deploy|release|ship)`)

	// BlockPattern matches notes about blockers or issues
	BlockPattern = regexp.MustCompile(`(?i)(?:block|stuck|issue|problem)`)
)

// PredefinedMatchers returns a map of common pattern matchers
func PredefinedMatchers() map[string]*NotePatternMatcher {
	matchers := make(map[string]*NotePatternMatcher)

	// Create matchers for common patterns
	patterns := map[string]string{
		"phase-complete": `(?i)phase.*(?:complete|done|finished)`,
		"bugfix":         `(?i)(?:bug|fix|issue).*`,
		"deploy":         `(?i)(?:deploy|release|ship)`,
		"blocker":        `(?i)(?:block|stuck|issue|problem)`,
	}

	for name, pattern := range patterns {
		if matcher, err := NewNotePatternMatcher(pattern); err == nil {
			matchers[name] = matcher
		}
	}

	return matchers
}

// ValidatePattern validates a regex pattern without creating a matcher
func ValidatePattern(pattern string) error {
	if pattern == "" {
		return nil
	}
	_, err := regexp.Compile(pattern)
	return err
}

// EscapePatternForLiteral escapes special regex characters for literal string matching
func EscapePatternForLiteral(literal string) string {
	return regexp.QuoteMeta(literal)
}
