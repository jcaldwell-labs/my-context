package models

import (
	"testing"
	"time"
)

func TestAppStateHasActiveContext(t *testing.T) {
	ctx := "test-context"
	empty := ""

	tests := []struct {
		name  string
		state AppState
		want  bool
	}{
		{
			name:  "has active context",
			state: AppState{ActiveContext: &ctx, LastUpdated: time.Now()},
			want:  true,
		},
		{
			name:  "nil active context",
			state: AppState{ActiveContext: nil, LastUpdated: time.Now()},
			want:  false,
		},
		{
			name:  "empty string context",
			state: AppState{ActiveContext: &empty, LastUpdated: time.Now()},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.state.HasActiveContext()
			if got != tt.want {
				t.Errorf("HasActiveContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppStateGetActiveContextName(t *testing.T) {
	ctx := "test-context"

	tests := []struct {
		name  string
		state AppState
		want  string
	}{
		{
			name:  "has active context",
			state: AppState{ActiveContext: &ctx, LastUpdated: time.Now()},
			want:  "test-context",
		},
		{
			name:  "nil active context",
			state: AppState{ActiveContext: nil, LastUpdated: time.Now()},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.state.GetActiveContextName()
			if got != tt.want {
				t.Errorf("GetActiveContextName() = %q, want %q", got, tt.want)
			}
		})
	}
}
