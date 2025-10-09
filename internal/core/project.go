package core

import (
	"strings"
)

// ExtractProjectName extracts the project name from a context name
// Contexts following "project: phase - description" convention return the text before the first colon
// Contexts without colons return the full name as the project name
func ExtractProjectName(contextName string) string {
	// Trim leading/trailing whitespace
	contextName = strings.TrimSpace(contextName)

	// If no colon, the entire name is the project
	if !strings.Contains(contextName, ":") {
		return contextName
	}

	// Extract text before first colon
	parts := strings.SplitN(contextName, ":", 2)
	projectName := strings.TrimSpace(parts[0])

	return projectName
}

// FilterContextsByProject filters a list of context names by project name (case-insensitive)
func FilterContextsByProject(contextNames []string, projectFilter string) []string {
	if projectFilter == "" {
		return contextNames
	}

	projectFilter = strings.ToLower(strings.TrimSpace(projectFilter))
	var filtered []string

	for _, name := range contextNames {
		projectName := ExtractProjectName(name)

		// Case-insensitive comparison
		if strings.ToLower(projectName) == projectFilter {
			filtered = append(filtered, name)
		}
	}

	return filtered
}
