// Package filters provides a unified interface for filtering contexts across different storage backends.
package filters

import (
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// ContextFilter provides a fluent interface for filtering contexts.
// Implementations should support method chaining for composing multiple filters.
type ContextFilter interface {
	ByProject(project string) ContextFilter
	ByTag(tag string) ContextFilter
	ByArchiveStatus(showArchived, activeOnly bool) ContextFilter
	BySearch(term string) ContextFilter
	ByDateRange(start, end time.Time) ContextFilter
	Apply() []interface{}
}

// ContextItem is an adapter interface for working with both file-based and database contexts
type ContextItem interface {
	GetName() string
	GetStartTime() time.Time
	GetEndTime() *time.Time
	GetIsArchived() bool
	GetIsActive() bool
	GetLabels() []string
}

// FileContextItem adapts file-based Context to ContextItem interface
type FileContextItem struct {
	Context  *models.Context
	Metadata *pkgmodels.ContextMetadata
}

func (f *FileContextItem) GetName() string {
	return f.Context.Name
}

func (f *FileContextItem) GetStartTime() time.Time {
	return f.Context.StartTime
}

func (f *FileContextItem) GetEndTime() *time.Time {
	return f.Context.EndTime
}

func (f *FileContextItem) GetIsArchived() bool {
	return f.Context.IsArchived
}

func (f *FileContextItem) GetIsActive() bool {
	return f.Context.IsActive()
}

func (f *FileContextItem) GetLabels() []string {
	if f.Metadata != nil {
		return f.Metadata.Labels
	}
	return []string{}
}

// DBContextItem adapts database ContextWithMetadata to ContextItem interface
type DBContextItem struct {
	Context *pkgmodels.ContextWithMetadata
}

func (d *DBContextItem) GetName() string {
	return d.Context.Name
}

func (d *DBContextItem) GetStartTime() time.Time {
	return d.Context.StartTime
}

func (d *DBContextItem) GetEndTime() *time.Time {
	return d.Context.EndTime
}

func (d *DBContextItem) GetIsArchived() bool {
	return d.Context.IsArchived
}

func (d *DBContextItem) GetIsActive() bool {
	return d.Context.IsActive()
}

func (d *DBContextItem) GetLabels() []string {
	return d.Context.Metadata.Labels
}

// filterByProject extracts project name from context name and matches against filter
func filterByProject(item ContextItem, projectFilter string) bool {
	projectName := core.ExtractProjectName(item.GetName())
	return strings.EqualFold(projectName, projectFilter)
}

// filterByTag checks if context has the specified tag
func filterByTag(item ContextItem, tag string) bool {
	for _, label := range item.GetLabels() {
		if strings.EqualFold(label, tag) {
			return true
		}
	}
	return false
}

// filterByArchiveStatus filters contexts based on archive status flags
func filterByArchiveStatus(item ContextItem, showArchived, activeOnly bool) bool {
	if activeOnly && !item.GetIsActive() {
		return false
	}
	if !showArchived && item.GetIsArchived() {
		return false
	}
	return true
}

// filterBySearch performs case-insensitive substring search on context name
func filterBySearch(item ContextItem, term string) bool {
	return strings.Contains(strings.ToLower(item.GetName()), strings.ToLower(term))
}

// filterByDateRange checks if context start time falls within the specified range
func filterByDateRange(item ContextItem, start, end time.Time) bool {
	startTime := item.GetStartTime()
	return !startTime.Before(start) && !startTime.After(end)
}
