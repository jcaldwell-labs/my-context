package filters

import (
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// ContextFilter provides a chainable interface for filtering contexts
type ContextFilter interface {
	ByProject(project string) ContextFilter
	BySearch(term string) ContextFilter
	ByTag(tag string) ContextFilter
	ByArchiveStatus(showArchived, activeOnly bool) ContextFilter
	ByActive(activeContextName string) ContextFilter
	ByPeriod(start, end time.Time) ContextFilter
	ByPattern(pattern string) ContextFilter
	ByStopDate(before time.Time) ContextFilter
	GetFileContexts() []*models.Context
	GetDBContexts() []*pkgmodels.ContextWithMetadata
}

// FileContextFilter implements ContextFilter for file-based contexts
type FileContextFilter struct {
	contexts []*models.Context
}

// NewFileContextFilter creates a new filter for file-based contexts
func NewFileContextFilter(contexts []*models.Context) *FileContextFilter {
	return &FileContextFilter{contexts: contexts}
}

// ByProject filters contexts by project name
func (f *FileContextFilter) ByProject(project string) ContextFilter {
	if project == "" {
		return f
	}
	
	// Extract context names
	contextNames := make([]string, 0, len(f.contexts))
	for _, ctx := range f.contexts {
		contextNames = append(contextNames, ctx.Name)
	}
	
	// Use core filtering logic
	filteredNames := core.FilterContextsByProject(contextNames, project)
	
	// Build filtered list
	filtered := make([]*models.Context, 0, len(f.contexts))
	for _, ctx := range f.contexts {
		for _, name := range filteredNames {
			if ctx.Name == name {
				filtered = append(filtered, ctx)
				break
			}
		}
	}
	
	return &FileContextFilter{contexts: filtered}
}

// BySearch filters contexts by search term (case-insensitive)
func (f *FileContextFilter) BySearch(term string) ContextFilter {
	if term == "" {
		return f
	}
	
	searchLower := strings.ToLower(term)
	filtered := make([]*models.Context, 0, len(f.contexts))
	for _, ctx := range f.contexts {
		if strings.Contains(strings.ToLower(ctx.Name), searchLower) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &FileContextFilter{contexts: filtered}
}

// ByTag filters contexts by tag/label
func (f *FileContextFilter) ByTag(tag string) ContextFilter {
	if tag == "" {
		return f
	}
	
	filtered := make([]*models.Context, 0, len(f.contexts))
	for _, ctx := range f.contexts {
		ctxWithMeta, _, _, _, err := core.GetContextWithMetadata(ctx.Name)
		if err != nil {
			continue
		}
		for _, t := range ctxWithMeta.Metadata.Labels {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, ctx)
				break
			}
		}
	}
	
	return &FileContextFilter{contexts: filtered}
}

// ByArchiveStatus filters contexts by archive status
func (f *FileContextFilter) ByArchiveStatus(showArchived, activeOnly bool) ContextFilter {
	if showArchived {
		filtered := make([]*models.Context, 0, len(f.contexts))
		for _, ctx := range f.contexts {
			if ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return &FileContextFilter{contexts: filtered}
	}
	
	if !activeOnly {
		filtered := make([]*models.Context, 0, len(f.contexts))
		for _, ctx := range f.contexts {
			if !ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return &FileContextFilter{contexts: filtered}
	}
	
	return f
}

// ByActive filters to show only the active context
func (f *FileContextFilter) ByActive(activeContextName string) ContextFilter {
	if activeContextName == "" {
		return &FileContextFilter{contexts: []*models.Context{}}
	}
	
	for _, ctx := range f.contexts {
		if ctx.Name == activeContextName {
			return &FileContextFilter{contexts: []*models.Context{ctx}}
		}
	}
	
	return &FileContextFilter{contexts: []*models.Context{}}
}

// ByPeriod filters contexts by time period
func (f *FileContextFilter) ByPeriod(start, end time.Time) ContextFilter {
	filtered := make([]*models.Context, 0, len(f.contexts))
	for _, ctx := range f.contexts {
		if ctx.StartTime.After(start) && ctx.StartTime.Before(end) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &FileContextFilter{contexts: filtered}
}

// ByPattern filters contexts using glob-style pattern matching
func (f *FileContextFilter) ByPattern(pattern string) ContextFilter {
	if pattern == "" {
		return f
	}
	
	matches := make([]*models.Context, 0, len(f.contexts))
	patternParts := strings.Split(pattern, "*")
	
	for _, ctx := range f.contexts {
		if core.MatchesPattern(ctx.Name, patternParts) {
			matches = append(matches, ctx)
		}
	}
	
	return &FileContextFilter{contexts: matches}
}

// ByStopDate filters contexts stopped before the given time
func (f *FileContextFilter) ByStopDate(before time.Time) ContextFilter {
	filtered := make([]*models.Context, 0, len(f.contexts))
	for _, ctx := range f.contexts {
		if ctx.EndTime != nil && ctx.EndTime.Before(before) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &FileContextFilter{contexts: filtered}
}

// GetFileContexts returns the filtered file-based contexts
func (f *FileContextFilter) GetFileContexts() []*models.Context {
	return f.contexts
}

// GetDBContexts returns empty for file-based filter (not applicable)
func (f *FileContextFilter) GetDBContexts() []*pkgmodels.ContextWithMetadata {
	return []*pkgmodels.ContextWithMetadata{}
}

// DBContextFilter implements ContextFilter for database contexts
type DBContextFilter struct {
	contexts []*pkgmodels.ContextWithMetadata
}

// NewDBContextFilter creates a new filter for database contexts
func NewDBContextFilter(contexts []*pkgmodels.ContextWithMetadata) *DBContextFilter {
	return &DBContextFilter{contexts: contexts}
}

// ByProject filters database contexts by project name
func (d *DBContextFilter) ByProject(project string) ContextFilter {
	if project == "" {
		return d
	}
	
	projectFilter := strings.TrimSpace(project)
	filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
	for _, ctx := range d.contexts {
		contextProject := core.ExtractProjectName(ctx.Name)
		if strings.EqualFold(contextProject, projectFilter) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &DBContextFilter{contexts: filtered}
}

// BySearch filters database contexts by search term (case-insensitive)
func (d *DBContextFilter) BySearch(term string) ContextFilter {
	if term == "" {
		return d
	}
	
	searchLower := strings.ToLower(term)
	filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
	for _, ctx := range d.contexts {
		if strings.Contains(strings.ToLower(ctx.Name), searchLower) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &DBContextFilter{contexts: filtered}
}

// ByTag filters database contexts by tag/label
func (d *DBContextFilter) ByTag(tag string) ContextFilter {
	if tag == "" {
		return d
	}
	
	filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
	for _, ctx := range d.contexts {
		for _, t := range ctx.Metadata.Labels {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, ctx)
				break
			}
		}
	}
	
	return &DBContextFilter{contexts: filtered}
}

// ByArchiveStatus filters database contexts by archive status
func (d *DBContextFilter) ByArchiveStatus(showArchived, activeOnly bool) ContextFilter {
	if showArchived {
		filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
		for _, ctx := range d.contexts {
			if ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return &DBContextFilter{contexts: filtered}
	}
	
	if !activeOnly {
		filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
		for _, ctx := range d.contexts {
			if !ctx.IsArchived {
				filtered = append(filtered, ctx)
			}
		}
		return &DBContextFilter{contexts: filtered}
	}
	
	return d
}

// ByActive filters to show only the active context
func (d *DBContextFilter) ByActive(activeContextName string) ContextFilter {
	if activeContextName == "" {
		return &DBContextFilter{contexts: []*pkgmodels.ContextWithMetadata{}}
	}
	
	for _, ctx := range d.contexts {
		if ctx.Name == activeContextName {
			return &DBContextFilter{contexts: []*pkgmodels.ContextWithMetadata{ctx}}
		}
	}
	
	return &DBContextFilter{contexts: []*pkgmodels.ContextWithMetadata{}}
}

// ByPeriod filters database contexts by time period
func (d *DBContextFilter) ByPeriod(start, end time.Time) ContextFilter {
	filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
	for _, ctx := range d.contexts {
		if ctx.StartTime.After(start) && ctx.StartTime.Before(end) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &DBContextFilter{contexts: filtered}
}

// ByPattern filters database contexts using glob-style pattern matching
func (d *DBContextFilter) ByPattern(pattern string) ContextFilter {
	if pattern == "" {
		return d
	}
	
	matches := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
	patternParts := strings.Split(pattern, "*")
	
	for _, ctx := range d.contexts {
		if core.MatchesPattern(ctx.Name, patternParts) {
			matches = append(matches, ctx)
		}
	}
	
	return &DBContextFilter{contexts: matches}
}

// ByStopDate filters database contexts stopped before the given time
func (d *DBContextFilter) ByStopDate(before time.Time) ContextFilter {
	filtered := make([]*pkgmodels.ContextWithMetadata, 0, len(d.contexts))
	for _, ctx := range d.contexts {
		if ctx.EndTime != nil && ctx.EndTime.Before(before) {
			filtered = append(filtered, ctx)
		}
	}
	
	return &DBContextFilter{contexts: filtered}
}

// GetFileContexts returns empty for database filter (not applicable)
func (d *DBContextFilter) GetFileContexts() []*models.Context {
	return []*models.Context{}
}

// GetDBContexts returns the filtered database contexts
func (d *DBContextFilter) GetDBContexts() []*pkgmodels.ContextWithMetadata {
	return d.contexts
}
