package filters

import (
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// FileContextFilter implements ContextFilter for file-based storage
type FileContextFilter struct {
	contexts      []*models.Context
	items         []ContextItem
	projectFilter string
	tagFilter     string
	searchTerm    string
	showArchived  bool
	activeOnly    bool
	dateStart     time.Time
	dateEnd       time.Time
	hasDateRange  bool
}

// NewFileContextFilter creates a new filter for file-based contexts
func NewFileContextFilter(contexts []*models.Context) *FileContextFilter {
	return &FileContextFilter{
		contexts: contexts,
	}
}

// ByProject filters contexts by project name
func (f *FileContextFilter) ByProject(project string) ContextFilter {
	f.projectFilter = project
	return f
}

// ByTag filters contexts by tag/label
func (f *FileContextFilter) ByTag(tag string) ContextFilter {
	f.tagFilter = tag
	return f
}

// ByArchiveStatus filters contexts based on archive status
func (f *FileContextFilter) ByArchiveStatus(showArchived, activeOnly bool) ContextFilter {
	f.showArchived = showArchived
	f.activeOnly = activeOnly
	return f
}

// BySearch filters contexts by search term in name
func (f *FileContextFilter) BySearch(term string) ContextFilter {
	f.searchTerm = term
	return f
}

// ByDateRange filters contexts by date range
func (f *FileContextFilter) ByDateRange(start, end time.Time) ContextFilter {
	f.dateStart = start
	f.dateEnd = end
	f.hasDateRange = true
	return f
}

// Apply executes all filters and returns the filtered results
func (f *FileContextFilter) Apply() []interface{} {
	// First, load metadata for tag filtering if needed
	if f.tagFilter != "" {
		f.loadMetadata()
	}

	// Apply all filters
	filtered := make([]interface{}, 0, len(f.contexts))
	for i, ctx := range f.contexts {
		var item ContextItem
		
		// Create item with metadata if we loaded it
		if len(f.items) > 0 && i < len(f.items) {
			item = f.items[i]
		} else {
			item = &FileContextItem{Context: ctx}
		}

		// Apply all active filters
		if f.projectFilter != "" && !filterByProject(item, f.projectFilter) {
			continue
		}
		if f.tagFilter != "" && !filterByTag(item, f.tagFilter) {
			continue
		}
		if !filterByArchiveStatus(item, f.showArchived, f.activeOnly) {
			continue
		}
		if f.searchTerm != "" && !filterBySearch(item, f.searchTerm) {
			continue
		}
		if f.hasDateRange && !filterByDateRange(item, f.dateStart, f.dateEnd) {
			continue
		}

		filtered = append(filtered, ctx)
	}

	return filtered
}

// loadMetadata loads metadata for all contexts (needed for tag filtering)
func (f *FileContextFilter) loadMetadata() {
	f.items = make([]ContextItem, len(f.contexts))
	for i, ctx := range f.contexts {
		// Load metadata for this context
		ctxWithMeta, _, _, _, _ := core.GetContextWithMetadata(ctx.Name)
		var metadata *pkgmodels.ContextMetadata
		if ctxWithMeta != nil {
			metadata = &ctxWithMeta.Metadata
		}
		f.items[i] = &FileContextItem{
			Context:  ctx,
			Metadata: metadata,
		}
	}
}
