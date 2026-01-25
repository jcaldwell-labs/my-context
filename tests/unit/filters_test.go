package filters_test

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/filters"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/stretchr/testify/assert"
)

// Helper functions to create test contexts
func createFileContext(name string, startTime time.Time, isArchived bool) *models.Context {
	endTime := startTime.Add(1 * time.Hour)
	return &models.Context{
		Name:      name,
		StartTime: startTime,
		EndTime:   &endTime,
		Status:    "stopped",
		IsArchived: isArchived,
	}
}

func createDBContext(name string, startTime time.Time, isArchived bool, labels []string) *pkgmodels.ContextWithMetadata {
	endTime := startTime.Add(1 * time.Hour)
	return &pkgmodels.ContextWithMetadata{
		Name:       name,
		StartTime:  startTime,
		EndTime:    &endTime,
		Status:     "stopped",
		IsArchived: isArchived,
		Metadata: pkgmodels.ContextMetadata{
			Labels: labels,
		},
	}
}

func TestFileContextFilter_ByProject(t *testing.T) {
	now := time.Now()
	contexts := []*models.Context{
		createFileContext("myproject: Phase 1", now, false),
		createFileContext("myproject: Phase 2", now, false),
		createFileContext("otherproject: Phase 1", now, false),
		createFileContext("Standalone task", now, false),
	}

	filter := filters.NewFileContextFilter(contexts)
	results := filter.ByProject("myproject").Apply()

	assert.Len(t, results, 2, "Should filter to 2 contexts with myproject")
	assert.Equal(t, "myproject: Phase 1", results[0].(*models.Context).Name)
	assert.Equal(t, "myproject: Phase 2", results[1].(*models.Context).Name)
}

func TestDBContextFilter_ByProject(t *testing.T) {
	now := time.Now()
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("myproject: Phase 1", now, false, nil),
		createDBContext("myproject: Phase 2", now, false, nil),
		createDBContext("otherproject: Phase 1", now, false, nil),
		createDBContext("Standalone task", now, false, nil),
	}

	filter := filters.NewDBContextFilter(contexts)
	results := filter.ByProject("myproject").Apply()

	assert.Len(t, results, 2, "Should filter to 2 contexts with myproject")
	assert.Equal(t, "myproject: Phase 1", results[0].(*pkgmodels.ContextWithMetadata).Name)
	assert.Equal(t, "myproject: Phase 2", results[1].(*pkgmodels.ContextWithMetadata).Name)
}

func TestDBContextFilter_ByTag(t *testing.T) {
	now := time.Now()
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("Context 1", now, false, []string{"bug", "urgent"}),
		createDBContext("Context 2", now, false, []string{"feature"}),
		createDBContext("Context 3", now, false, []string{"bug"}),
		createDBContext("Context 4", now, false, nil),
	}

	filter := filters.NewDBContextFilter(contexts)
	results := filter.ByTag("bug").Apply()

	assert.Len(t, results, 2, "Should filter to 2 contexts with 'bug' tag")
	assert.Equal(t, "Context 1", results[0].(*pkgmodels.ContextWithMetadata).Name)
	assert.Equal(t, "Context 3", results[1].(*pkgmodels.ContextWithMetadata).Name)
}

func TestFileContextFilter_ByArchiveStatus(t *testing.T) {
	now := time.Now()
	contexts := []*models.Context{
		createFileContext("Active context", now, false),
		createFileContext("Archived context 1", now, true),
		createFileContext("Archived context 2", now, true),
	}

	t.Run("Hide archived by default", func(t *testing.T) {
		filter := filters.NewFileContextFilter(contexts)
		results := filter.ByArchiveStatus(false, false).Apply()
		assert.Len(t, results, 1, "Should only show non-archived contexts")
	})

	t.Run("Show archived when requested", func(t *testing.T) {
		filter := filters.NewFileContextFilter(contexts)
		results := filter.ByArchiveStatus(true, false).Apply()
		assert.Len(t, results, 3, "Should show all contexts including archived")
	})
}

func TestDBContextFilter_ByArchiveStatus(t *testing.T) {
	now := time.Now()
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("Active context", now, false, nil),
		createDBContext("Archived context 1", now, true, nil),
		createDBContext("Archived context 2", now, true, nil),
	}

	t.Run("Hide archived by default", func(t *testing.T) {
		filter := filters.NewDBContextFilter(contexts)
		results := filter.ByArchiveStatus(false, false).Apply()
		assert.Len(t, results, 1, "Should only show non-archived contexts")
	})

	t.Run("Show archived when requested", func(t *testing.T) {
		filter := filters.NewDBContextFilter(contexts)
		results := filter.ByArchiveStatus(true, false).Apply()
		assert.Len(t, results, 3, "Should show all contexts including archived")
	})
}

func TestFileContextFilter_BySearch(t *testing.T) {
	now := time.Now()
	contexts := []*models.Context{
		createFileContext("Fix payment bug", now, false),
		createFileContext("Add payment feature", now, false),
		createFileContext("Update documentation", now, false),
	}

	filter := filters.NewFileContextFilter(contexts)
	results := filter.BySearch("payment").Apply()

	assert.Len(t, results, 2, "Should find 2 contexts with 'payment' in name")
}

func TestDBContextFilter_BySearch(t *testing.T) {
	now := time.Now()
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("Fix payment bug", now, false, nil),
		createDBContext("Add payment feature", now, false, nil),
		createDBContext("Update documentation", now, false, nil),
	}

	filter := filters.NewDBContextFilter(contexts)
	results := filter.BySearch("payment").Apply()

	assert.Len(t, results, 2, "Should find 2 contexts with 'payment' in name")
}

func TestFileContextFilter_ByDateRange(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	contexts := []*models.Context{
		createFileContext("Context 1", baseTime.AddDate(0, 0, -10), false), // Jan 5
		createFileContext("Context 2", baseTime.AddDate(0, 0, -5), false),  // Jan 10
		createFileContext("Context 3", baseTime, false),                     // Jan 15
		createFileContext("Context 4", baseTime.AddDate(0, 0, 5), false),   // Jan 20
	}

	start := baseTime.AddDate(0, 0, -7) // Jan 8
	end := baseTime.AddDate(0, 0, 2)    // Jan 17

	filter := filters.NewFileContextFilter(contexts)
	results := filter.ByDateRange(start, end).Apply()

	assert.Len(t, results, 2, "Should find contexts within date range")
	assert.Equal(t, "Context 2", results[0].(*models.Context).Name)
	assert.Equal(t, "Context 3", results[1].(*models.Context).Name)
}

func TestDBContextFilter_ByDateRange(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("Context 1", baseTime.AddDate(0, 0, -10), false, nil), // Jan 5
		createDBContext("Context 2", baseTime.AddDate(0, 0, -5), false, nil),  // Jan 10
		createDBContext("Context 3", baseTime, false, nil),                     // Jan 15
		createDBContext("Context 4", baseTime.AddDate(0, 0, 5), false, nil),   // Jan 20
	}

	start := baseTime.AddDate(0, 0, -7) // Jan 8
	end := baseTime.AddDate(0, 0, 2)    // Jan 17

	filter := filters.NewDBContextFilter(contexts)
	results := filter.ByDateRange(start, end).Apply()

	assert.Len(t, results, 2, "Should find contexts within date range")
	assert.Equal(t, "Context 2", results[0].(*pkgmodels.ContextWithMetadata).Name)
	assert.Equal(t, "Context 3", results[1].(*pkgmodels.ContextWithMetadata).Name)
}

func TestFileContextFilter_ChainedFilters(t *testing.T) {
	now := time.Now()
	contexts := []*models.Context{
		createFileContext("myproject: Fix bug", now, false),
		createFileContext("myproject: Add feature", now, false),
		createFileContext("myproject: Old task", now.AddDate(0, 0, -100), false),
		createFileContext("otherproject: Fix bug", now, false),
	}

	start := now.AddDate(0, 0, -10)
	end := now.AddDate(0, 0, 10)

	filter := filters.NewFileContextFilter(contexts)
	results := filter.
		ByProject("myproject").
		BySearch("bug").
		ByDateRange(start, end).
		Apply()

	assert.Len(t, results, 1, "Should apply all filters")
	assert.Equal(t, "myproject: Fix bug", results[0].(*models.Context).Name)
}

func TestDBContextFilter_ChainedFilters(t *testing.T) {
	now := time.Now()
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("myproject: Fix bug", now, false, []string{"bug", "urgent"}),
		createDBContext("myproject: Add feature", now, false, []string{"feature"}),
		createDBContext("myproject: Old bug", now.AddDate(0, 0, -100), false, []string{"bug"}),
		createDBContext("otherproject: Fix bug", now, false, []string{"bug"}),
	}

	start := now.AddDate(0, 0, -10)
	end := now.AddDate(0, 0, 10)

	filter := filters.NewDBContextFilter(contexts)
	results := filter.
		ByProject("myproject").
		ByTag("bug").
		ByDateRange(start, end).
		Apply()

	assert.Len(t, results, 1, "Should apply all filters")
	assert.Equal(t, "myproject: Fix bug", results[0].(*pkgmodels.ContextWithMetadata).Name)
}

func TestFileContextFilter_EmptyResults(t *testing.T) {
	now := time.Now()
	contexts := []*models.Context{
		createFileContext("Context 1", now, false),
		createFileContext("Context 2", now, false),
	}

	filter := filters.NewFileContextFilter(contexts)
	results := filter.BySearch("nonexistent").Apply()

	assert.Len(t, results, 0, "Should return empty results when no matches")
}

func TestDBContextFilter_EmptyResults(t *testing.T) {
	now := time.Now()
	contexts := []*pkgmodels.ContextWithMetadata{
		createDBContext("Context 1", now, false, nil),
		createDBContext("Context 2", now, false, nil),
	}

	filter := filters.NewDBContextFilter(contexts)
	results := filter.BySearch("nonexistent").Apply()

	assert.Len(t, results, 0, "Should return empty results when no matches")
}

func TestContextItem_Adapters(t *testing.T) {
	now := time.Now()
	endTime := now.Add(1 * time.Hour)

	t.Run("filters.FileContextItem adapter", func(t *testing.T) {
		ctx := &models.Context{
			Name:       "Test Context",
			StartTime:  now,
			EndTime:    &endTime,
			Status:     "stopped",
			IsArchived: true,
		}
		metadata := &pkgmodels.ContextMetadata{
			Labels: []string{"tag1", "tag2"},
		}

		item := &filters.FileContextItem{Context: ctx, Metadata: metadata}

		assert.Equal(t, "Test Context", item.GetName())
		assert.Equal(t, now, item.GetStartTime())
		assert.Equal(t, &endTime, item.GetEndTime())
		assert.True(t, item.GetIsArchived())
		assert.False(t, item.GetIsActive())
		assert.Equal(t, []string{"tag1", "tag2"}, item.GetLabels())
	})

	t.Run("filters.DBContextItem adapter", func(t *testing.T) {
		ctx := &pkgmodels.ContextWithMetadata{
			Name:       "Test Context",
			StartTime:  now,
			EndTime:    &endTime,
			Status:     "stopped",
			IsArchived: true,
			Metadata: pkgmodels.ContextMetadata{
				Labels: []string{"tag1", "tag2"},
			},
		}

		item := &filters.DBContextItem{Context: ctx}

		assert.Equal(t, "Test Context", item.GetName())
		assert.Equal(t, now, item.GetStartTime())
		assert.Equal(t, &endTime, item.GetEndTime())
		assert.True(t, item.GetIsArchived())
		assert.False(t, item.GetIsActive())
		assert.Equal(t, []string{"tag1", "tag2"}, item.GetLabels())
	})
}
