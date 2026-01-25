package filters

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	pkgmodels "github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/stretchr/testify/assert"
)

// Test data helpers
func createFileContexts() []*models.Context {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	return []*models.Context{
		{Name: "project-a: feature-1", StartTime: now, Status: "active", IsArchived: false},
		{Name: "project-a: feature-2", StartTime: yesterday, Status: "stopped", EndTime: &yesterday, IsArchived: false},
		{Name: "project-b: bug-fix", StartTime: lastWeek, Status: "stopped", EndTime: &lastWeek, IsArchived: true},
		{Name: "standalone-task", StartTime: now, Status: "stopped", EndTime: &now, IsArchived: false},
	}
}

func createDBContexts() []*pkgmodels.ContextWithMetadata {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	return []*pkgmodels.ContextWithMetadata{
		{
			Name:       "project-a: feature-1",
			StartTime:  now,
			Status:     "active",
			IsArchived: false,
			Metadata:   pkgmodels.ContextMetadata{Labels: []string{"frontend", "urgent"}},
		},
		{
			Name:       "project-a: feature-2",
			StartTime:  yesterday,
			EndTime:    &yesterday,
			Status:     "stopped",
			IsArchived: false,
			Metadata:   pkgmodels.ContextMetadata{Labels: []string{"backend"}},
		},
		{
			Name:       "project-b: bug-fix",
			StartTime:  lastWeek,
			EndTime:    &lastWeek,
			Status:     "stopped",
			IsArchived: true,
			Metadata:   pkgmodels.ContextMetadata{Labels: []string{"urgent", "bug"}},
		},
		{
			Name:       "standalone-task",
			StartTime:  now,
			EndTime:    &now,
			Status:     "stopped",
			IsArchived: false,
			Metadata:   pkgmodels.ContextMetadata{Labels: []string{}},
		},
	}
}

// FileContextFilter tests

func TestFileContextFilter_ByProject(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("filter by project-a", func(t *testing.T) {
		result := filter.ByProject("project-a").GetFileContexts()
		assert.Len(t, result, 2)
		assert.Equal(t, "project-a: feature-1", result[0].Name)
		assert.Equal(t, "project-a: feature-2", result[1].Name)
	})

	t.Run("filter by project-b", func(t *testing.T) {
		result := filter.ByProject("project-b").GetFileContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-b: bug-fix", result[0].Name)
	})

	t.Run("empty project returns all", func(t *testing.T) {
		result := filter.ByProject("").GetFileContexts()
		assert.Len(t, result, 4)
	})

	t.Run("non-existent project returns empty", func(t *testing.T) {
		result := filter.ByProject("nonexistent").GetFileContexts()
		assert.Len(t, result, 0)
	})
}

func TestFileContextFilter_BySearch(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("search for 'feature'", func(t *testing.T) {
		result := filter.BySearch("feature").GetFileContexts()
		assert.Len(t, result, 2)
	})

	t.Run("search for 'bug'", func(t *testing.T) {
		result := filter.BySearch("bug").GetFileContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-b: bug-fix", result[0].Name)
	})

	t.Run("case insensitive search", func(t *testing.T) {
		result := filter.BySearch("FEATURE").GetFileContexts()
		assert.Len(t, result, 2)
	})

	t.Run("empty search returns all", func(t *testing.T) {
		result := filter.BySearch("").GetFileContexts()
		assert.Len(t, result, 4)
	})
}

func TestFileContextFilter_ByArchiveStatus(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("show only archived", func(t *testing.T) {
		result := filter.ByArchiveStatus(true, false).GetFileContexts()
		assert.Len(t, result, 1)
		assert.True(t, result[0].IsArchived)
	})

	t.Run("show only non-archived", func(t *testing.T) {
		result := filter.ByArchiveStatus(false, false).GetFileContexts()
		assert.Len(t, result, 3)
		for _, ctx := range result {
			assert.False(t, ctx.IsArchived)
		}
	})

	t.Run("activeOnly true returns all", func(t *testing.T) {
		result := filter.ByArchiveStatus(false, true).GetFileContexts()
		assert.Len(t, result, 4)
	})
}

func TestFileContextFilter_ByActive(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("filter by active context", func(t *testing.T) {
		result := filter.ByActive("project-a: feature-1").GetFileContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-a: feature-1", result[0].Name)
	})

	t.Run("non-existent active context", func(t *testing.T) {
		result := filter.ByActive("nonexistent").GetFileContexts()
		assert.Len(t, result, 0)
	})

	t.Run("empty active context", func(t *testing.T) {
		result := filter.ByActive("").GetFileContexts()
		assert.Len(t, result, 0)
	})
}

func TestFileContextFilter_ByPeriod(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("filter by recent period", func(t *testing.T) {
		start := time.Now().Add(-2 * 24 * time.Hour)
		end := time.Now().Add(1 * time.Hour)
		result := filter.ByPeriod(start, end).GetFileContexts()
		assert.Len(t, result, 3)
	})

	t.Run("filter by old period", func(t *testing.T) {
		start := time.Now().Add(-30 * 24 * time.Hour)
		end := time.Now().Add(-6 * 24 * time.Hour)
		result := filter.ByPeriod(start, end).GetFileContexts()
		assert.Len(t, result, 1)
	})
}

func TestFileContextFilter_ByPattern(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("pattern with wildcard", func(t *testing.T) {
		result := filter.ByPattern("project-a*").GetFileContexts()
		assert.Len(t, result, 2)
	})

	t.Run("pattern with multiple wildcards", func(t *testing.T) {
		result := filter.ByPattern("*bug*").GetFileContexts()
		assert.Len(t, result, 1)
	})

	t.Run("empty pattern returns all", func(t *testing.T) {
		result := filter.ByPattern("").GetFileContexts()
		assert.Len(t, result, 4)
	})
}

func TestFileContextFilter_ByStopDate(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("filter by stop date", func(t *testing.T) {
		cutoff := time.Now().Add(-3 * 24 * time.Hour)
		result := filter.ByStopDate(cutoff).GetFileContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-b: bug-fix", result[0].Name)
	})

	t.Run("recent cutoff returns none", func(t *testing.T) {
		cutoff := time.Now().Add(-10 * 24 * time.Hour)
		result := filter.ByStopDate(cutoff).GetFileContexts()
		assert.Len(t, result, 0)
	})
}

func TestFileContextFilter_Chaining(t *testing.T) {
	contexts := createFileContexts()
	filter := NewFileContextFilter(contexts)

	t.Run("chain project and archive filters", func(t *testing.T) {
		result := filter.
			ByProject("project-a").
			ByArchiveStatus(false, false).
			GetFileContexts()
		assert.Len(t, result, 2)
	})

	t.Run("chain search and project filters", func(t *testing.T) {
		result := filter.
			BySearch("feature").
			ByProject("project-a").
			GetFileContexts()
		assert.Len(t, result, 2)
	})
}

// DBContextFilter tests

func TestDBContextFilter_ByProject(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("filter by project-a", func(t *testing.T) {
		result := filter.ByProject("project-a").GetDBContexts()
		assert.Len(t, result, 2)
		assert.Equal(t, "project-a: feature-1", result[0].Name)
	})

	t.Run("case insensitive project", func(t *testing.T) {
		result := filter.ByProject("PROJECT-A").GetDBContexts()
		assert.Len(t, result, 2)
	})

	t.Run("empty project returns all", func(t *testing.T) {
		result := filter.ByProject("").GetDBContexts()
		assert.Len(t, result, 4)
	})
}

func TestDBContextFilter_BySearch(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("search for 'feature'", func(t *testing.T) {
		result := filter.BySearch("feature").GetDBContexts()
		assert.Len(t, result, 2)
	})

	t.Run("case insensitive search", func(t *testing.T) {
		result := filter.BySearch("BUG").GetDBContexts()
		assert.Len(t, result, 1)
	})
}

func TestDBContextFilter_ByTag(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("filter by tag 'urgent'", func(t *testing.T) {
		result := filter.ByTag("urgent").GetDBContexts()
		assert.Len(t, result, 2)
	})

	t.Run("filter by tag 'backend'", func(t *testing.T) {
		result := filter.ByTag("backend").GetDBContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-a: feature-2", result[0].Name)
	})

	t.Run("case insensitive tag", func(t *testing.T) {
		result := filter.ByTag("URGENT").GetDBContexts()
		assert.Len(t, result, 2)
	})

	t.Run("non-existent tag", func(t *testing.T) {
		result := filter.ByTag("nonexistent").GetDBContexts()
		assert.Len(t, result, 0)
	})

	t.Run("empty tag returns all", func(t *testing.T) {
		result := filter.ByTag("").GetDBContexts()
		assert.Len(t, result, 4)
	})
}

func TestDBContextFilter_ByArchiveStatus(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("show only archived", func(t *testing.T) {
		result := filter.ByArchiveStatus(true, false).GetDBContexts()
		assert.Len(t, result, 1)
		assert.True(t, result[0].IsArchived)
	})

	t.Run("show only non-archived", func(t *testing.T) {
		result := filter.ByArchiveStatus(false, false).GetDBContexts()
		assert.Len(t, result, 3)
	})
}

func TestDBContextFilter_ByActive(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("filter by active context", func(t *testing.T) {
		result := filter.ByActive("project-a: feature-1").GetDBContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-a: feature-1", result[0].Name)
	})

	t.Run("empty active context", func(t *testing.T) {
		result := filter.ByActive("").GetDBContexts()
		assert.Len(t, result, 0)
	})
}

func TestDBContextFilter_ByPeriod(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("filter by recent period", func(t *testing.T) {
		start := time.Now().Add(-2 * 24 * time.Hour)
		end := time.Now().Add(1 * time.Hour)
		result := filter.ByPeriod(start, end).GetDBContexts()
		assert.Len(t, result, 3)
	})
}

func TestDBContextFilter_ByPattern(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("pattern with wildcard", func(t *testing.T) {
		result := filter.ByPattern("project-a*").GetDBContexts()
		assert.Len(t, result, 2)
	})

	t.Run("empty pattern returns all", func(t *testing.T) {
		result := filter.ByPattern("").GetDBContexts()
		assert.Len(t, result, 4)
	})
}

func TestDBContextFilter_ByStopDate(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("filter by stop date", func(t *testing.T) {
		cutoff := time.Now().Add(-3 * 24 * time.Hour)
		result := filter.ByStopDate(cutoff).GetDBContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-b: bug-fix", result[0].Name)
	})
}

func TestDBContextFilter_Chaining(t *testing.T) {
	contexts := createDBContexts()
	filter := NewDBContextFilter(contexts)

	t.Run("chain project, tag, and archive filters", func(t *testing.T) {
		result := filter.
			ByProject("project-a").
			ByTag("frontend").
			ByArchiveStatus(false, false).
			GetDBContexts()
		assert.Len(t, result, 1)
		assert.Equal(t, "project-a: feature-1", result[0].Name)
	})

	t.Run("chain search and period filters", func(t *testing.T) {
		start := time.Now().Add(-25 * time.Hour)
		end := time.Now().Add(1 * time.Hour)
		result := filter.
			BySearch("feature").
			ByPeriod(start, end).
			GetDBContexts()
		assert.Len(t, result, 2)
	})
}

func TestFilter_CrossTypeAccess(t *testing.T) {
	t.Run("FileFilter GetDBContexts returns empty", func(t *testing.T) {
		contexts := createFileContexts()
		filter := NewFileContextFilter(contexts)
		result := filter.GetDBContexts()
		assert.Len(t, result, 0)
	})

	t.Run("DBFilter GetFileContexts returns empty", func(t *testing.T) {
		contexts := createDBContexts()
		filter := NewDBContextFilter(contexts)
		result := filter.GetFileContexts()
		assert.Len(t, result, 0)
	})
}
