package filters

import (
	"time"

	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// DBContextFilter implements ContextFilter for database storage
type DBContextFilter struct {
	contexts      []*pkgmodels.ContextWithMetadata
	projectFilter string
	tagFilter     string
	searchTerm    string
	showArchived  bool
	activeOnly    bool
	dateStart     time.Time
	dateEnd       time.Time
	hasDateRange  bool
}

// NewDBContextFilter creates a new filter for database contexts
func NewDBContextFilter(contexts []*pkgmodels.ContextWithMetadata) *DBContextFilter {
	return &DBContextFilter{
		contexts: contexts,
	}
}

// ByProject filters contexts by project name
func (d *DBContextFilter) ByProject(project string) ContextFilter {
	d.projectFilter = project
	return d
}

// ByTag filters contexts by tag/label
func (d *DBContextFilter) ByTag(tag string) ContextFilter {
	d.tagFilter = tag
	return d
}

// ByArchiveStatus filters contexts based on archive status
func (d *DBContextFilter) ByArchiveStatus(showArchived, activeOnly bool) ContextFilter {
	d.showArchived = showArchived
	d.activeOnly = activeOnly
	return d
}

// BySearch filters contexts by search term in name
func (d *DBContextFilter) BySearch(term string) ContextFilter {
	d.searchTerm = term
	return d
}

// ByDateRange filters contexts by date range
func (d *DBContextFilter) ByDateRange(start, end time.Time) ContextFilter {
	d.dateStart = start
	d.dateEnd = end
	d.hasDateRange = true
	return d
}

// Apply executes all filters and returns the filtered results
func (d *DBContextFilter) Apply() []interface{} {
	filtered := make([]interface{}, 0, len(d.contexts))
	
	for _, ctx := range d.contexts {
		item := &DBContextItem{Context: ctx}

		// Apply all active filters
		if d.projectFilter != "" && !filterByProject(item, d.projectFilter) {
			continue
		}
		if d.tagFilter != "" && !filterByTag(item, d.tagFilter) {
			continue
		}
		if !filterByArchiveStatus(item, d.showArchived, d.activeOnly) {
			continue
		}
		if d.searchTerm != "" && !filterBySearch(item, d.searchTerm) {
			continue
		}
		if d.hasDateRange && !filterByDateRange(item, d.dateStart, d.dateEnd) {
			continue
		}

		filtered = append(filtered, ctx)
	}

	return filtered
}
