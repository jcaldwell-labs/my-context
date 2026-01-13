package core

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
)

func TestMatchesFilter(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		ctx    *models.Context
		filter ContextFilter
		want   bool
	}{
		{
			name:   "empty filter matches all",
			ctx:    &models.Context{Name: "test", Status: "active", StartTime: now},
			filter: ContextFilter{},
			want:   true,
		},
		{
			name:   "archived context hidden by default",
			ctx:    &models.Context{Name: "test", Status: "stopped", StartTime: now, IsArchived: true},
			filter: ContextFilter{},
			want:   false,
		},
		{
			name:   "archived shown with ShowArchived",
			ctx:    &models.Context{Name: "test", Status: "stopped", StartTime: now, IsArchived: true},
			filter: ContextFilter{ShowArchived: true},
			want:   true,
		},
		{
			name:   "archived flag shows only archived",
			ctx:    &models.Context{Name: "test", Status: "stopped", StartTime: now, IsArchived: true},
			filter: ContextFilter{Archived: true},
			want:   true,
		},
		{
			name:   "archived flag hides non-archived",
			ctx:    &models.Context{Name: "test", Status: "active", StartTime: now, IsArchived: false},
			filter: ContextFilter{Archived: true},
			want:   false,
		},
		{
			name:   "active only shows active",
			ctx:    &models.Context{Name: "test", Status: "active", StartTime: now},
			filter: ContextFilter{ActiveOnly: true},
			want:   true,
		},
		{
			name:   "active only hides stopped",
			ctx:    &models.Context{Name: "test", Status: "stopped", StartTime: now},
			filter: ContextFilter{ActiveOnly: true},
			want:   false,
		},
		{
			name:   "project filter matches",
			ctx:    &models.Context{Name: "ps-cli: Phase 1", Status: "active", StartTime: now},
			filter: ContextFilter{Project: "ps-cli"},
			want:   true,
		},
		{
			name:   "project filter case insensitive",
			ctx:    &models.Context{Name: "PS-CLI: Phase 1", Status: "active", StartTime: now},
			filter: ContextFilter{Project: "ps-cli"},
			want:   true,
		},
		{
			name:   "project filter no match",
			ctx:    &models.Context{Name: "garden: Planning", Status: "active", StartTime: now},
			filter: ContextFilter{Project: "ps-cli"},
			want:   false,
		},
		{
			name:   "search filter matches",
			ctx:    &models.Context{Name: "feature-auth-implementation", Status: "active", StartTime: now},
			filter: ContextFilter{Search: "auth"},
			want:   true,
		},
		{
			name:   "search filter case insensitive",
			ctx:    &models.Context{Name: "Feature-Auth-Implementation", Status: "active", StartTime: now},
			filter: ContextFilter{Search: "auth"},
			want:   true,
		},
		{
			name:   "search filter no match",
			ctx:    &models.Context{Name: "feature-login", Status: "active", StartTime: now},
			filter: ContextFilter{Search: "auth"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesFilter(tt.ctx, tt.filter)
			if got != tt.want {
				t.Errorf("matchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyLimit(t *testing.T) {
	now := time.Now()
	contexts := []*models.Context{
		{Name: "ctx-1", Status: "active", StartTime: now},
		{Name: "ctx-2", Status: "active", StartTime: now},
		{Name: "ctx-3", Status: "active", StartTime: now},
		{Name: "ctx-4", Status: "active", StartTime: now},
		{Name: "ctx-5", Status: "active", StartTime: now},
	}

	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{"zero limit returns all", 0, 5},
		{"negative limit returns all", -1, 5},
		{"limit less than count", 3, 3},
		{"limit equal to count", 5, 5},
		{"limit greater than count", 10, 5},
		{"limit of 1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyLimit(contexts, tt.limit)
			if len(got) != tt.want {
				t.Errorf("applyLimit() len = %d, want %d", len(got), tt.want)
			}
		})
	}
}
