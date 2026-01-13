package models

import (
	"strings"
	"testing"
	"time"
)

func TestContextValidate(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	tests := []struct {
		name    string
		context Context
		wantErr string
	}{
		{
			name: "valid active context",
			context: Context{
				Name:      "my-context",
				StartTime: now,
				Status:    "active",
			},
			wantErr: "",
		},
		{
			name: "valid stopped context",
			context: Context{
				Name:      "my-context",
				StartTime: pastTime,
				EndTime:   &now,
				Status:    "stopped",
			},
			wantErr: "",
		},
		{
			name: "valid archived context",
			context: Context{
				Name:       "my-context",
				StartTime:  pastTime,
				EndTime:    &now,
				Status:     "stopped",
				IsArchived: true,
			},
			wantErr: "",
		},
		{
			name: "empty name",
			context: Context{
				Name:      "",
				StartTime: now,
				Status:    "active",
			},
			wantErr: "context name cannot be empty",
		},
		{
			name: "name with forward slash",
			context: Context{
				Name:      "my/context",
				StartTime: now,
				Status:    "active",
			},
			wantErr: "cannot contain path separators",
		},
		{
			name: "name with backslash",
			context: Context{
				Name:      "my\\context",
				StartTime: now,
				Status:    "active",
			},
			wantErr: "cannot contain path separators",
		},
		{
			name: "name too long",
			context: Context{
				Name:      strings.Repeat("a", 201),
				StartTime: now,
				Status:    "active",
			},
			wantErr: "200 characters or less",
		},
		{
			name: "invalid status",
			context: Context{
				Name:      "my-context",
				StartTime: now,
				Status:    "running",
			},
			wantErr: "status must be",
		},
		{
			name: "stopped without end time",
			context: Context{
				Name:      "my-context",
				StartTime: now,
				Status:    "stopped",
				EndTime:   nil,
			},
			wantErr: "stopped context must have an end time",
		},
		{
			name: "active with end time",
			context: Context{
				Name:      "my-context",
				StartTime: now,
				EndTime:   &now,
				Status:    "active",
			},
			wantErr: "active context cannot have an end time",
		},
		{
			name: "archived active context",
			context: Context{
				Name:       "my-context",
				StartTime:  now,
				Status:     "active",
				IsArchived: true,
			},
			wantErr: "cannot archive an active context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.context.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() expected error containing %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Validate() error = %q, want containing %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestContextDuration(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	t.Run("stopped context with end time", func(t *testing.T) {
		ctx := Context{
			Name:      "test",
			StartTime: oneHourAgo,
			EndTime:   &now,
			Status:    "stopped",
		}
		duration := ctx.Duration()
		if duration < 59*time.Minute || duration > 61*time.Minute {
			t.Errorf("Duration() = %v, want ~1h", duration)
		}
	})

	t.Run("active context calculates from now", func(t *testing.T) {
		ctx := Context{
			Name:      "test",
			StartTime: oneHourAgo,
			Status:    "active",
		}
		duration := ctx.Duration()
		if duration < 59*time.Minute {
			t.Errorf("Duration() = %v, want at least ~1h", duration)
		}
	})
}

func TestContextIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"active status", "active", true},
		{"stopped status", "stopped", false},
		{"empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := Context{Status: tt.status}
			if got := ctx.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
