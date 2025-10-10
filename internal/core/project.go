package core

import (
	"strings"
)

// ExtractProjectName extracts the project name from a context name.
// If the context name contains a colon, returns the text before the first colon (trimmed).
// If no colon, returns the full context name (trimmed).
//
// Examples:
//   "ps-cli: Phase 1" -> "ps-cli"
//   "garden: Planning" -> "garden"
//   "StandaloneContext" -> "StandaloneContext"
//   "project: Phase 1: Subphase A" -> "project" (only first colon used)
func ExtractProjectName(contextName string) string {
	// Split on first colon only
	parts := strings.SplitN(contextName, ":", 2)
	
	if len(parts) > 1 {
		// Has colon - return text before it, trimmed
		return strings.TrimSpace(parts[0])
	}
	
	// No colon - full name is project name, trimmed
	return strings.TrimSpace(contextName)
}

// FilterByProject filters a list of context names by project (case-insensitive)
func FilterByProject(contextNames []string, projectName string) []string {
	return FilterContextsByProject(contextNames, projectName)
}

// FilterContextsByProject filters a list of context names by project (case-insensitive)
func FilterContextsByProject(contextNames []string, projectName string) []string {
	if projectName == "" {
		return contextNames
	}

	var filtered []string
	projectLower := strings.ToLower(strings.TrimSpace(projectName))

	for _, name := range contextNames {
		contextProject := ExtractProjectName(name)
		contextProjectLower := strings.ToLower(contextProject)

		if contextProjectLower == projectLower {
			filtered = append(filtered, name)
		}
	}

	return filtered
}

// ProjectMetadata represents extracted project information
type ProjectMetadata struct {
	ProjectName   string
	ContextNames  []string
	ContextCount  int
}

// ExtractProjectMetadata analyzes a list of contexts and groups them by project
func ExtractProjectMetadata(contextNames []string) map[string]*ProjectMetadata {
	projects := make(map[string]*ProjectMetadata)

	for _, name := range contextNames {
		projectName := ExtractProjectName(name)
		projectKey := strings.ToLower(projectName) // Case-insensitive grouping

		if meta, exists := projects[projectKey]; exists {
			meta.ContextNames = append(meta.ContextNames, name)
			meta.ContextCount++
		} else {
			projects[projectKey] = &ProjectMetadata{
				ProjectName:  projectName, // Keep original casing
				ContextNames: []string{name},
				ContextCount: 1,
			}
		}
	}

	return projects
}
