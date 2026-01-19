package core

import "strings"

// MatchesPattern checks if a context name matches a glob pattern.
// patternParts should be the result of strings.Split(pattern, "*").
// Returns false if patternParts is empty (no pattern provided).
// Returns true if pattern was just "*" (matches everything).
func MatchesPattern(name string, patternParts []string) bool {
	if len(patternParts) == 0 {
		return false // Empty pattern parts means no pattern was split, so no match
	}

	// Handle the special case of just "*" which should match everything
	if len(patternParts) == 2 && patternParts[0] == "" && patternParts[1] == "" {
		return true
	}

	// Handle simple cases
	if len(patternParts) == 1 {
		// No wildcards - exact match required
		return name == patternParts[0]
	}

	// Check prefix
	if !strings.HasPrefix(name, patternParts[0]) {
		return false
	}

	remainingName := name[len(patternParts[0]):]

	// Check suffix
	lastPart := patternParts[len(patternParts)-1]
	if !strings.HasSuffix(remainingName, lastPart) {
		return false
	}

	// For multiple wildcards, we do a simple substring check
	for i := 1; i < len(patternParts)-1; i++ {
		if !strings.Contains(remainingName, patternParts[i]) {
			return false
		}
	}

	return true
}
