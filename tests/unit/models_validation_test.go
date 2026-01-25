package unit

import (
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestContextValidate(t *testing.T) {
	now := time.Now()
	later := now.Add(1 * time.Hour)

	tests := []struct {
		name    string
		context models.Context
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid active context",
			context: models.Context{
				Name:      "test-context",
				StartTime: now,
				Status:    "active",
				EndTime:   nil,
			},
			wantErr: false,
		},
		{
			name: "valid stopped context",
			context: models.Context{
				Name:      "test-context",
				StartTime: now,
				EndTime:   &later,
				Status:    "stopped",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			context: models.Context{
				Name:      "",
				StartTime: now,
				Status:    "active",
			},
			wantErr: true,
			errMsg:  "context name cannot be empty",
		},
		{
			name: "name with forward slash",
			context: models.Context{
				Name:      "test/context",
				StartTime: now,
				Status:    "active",
			},
			wantErr: true,
			errMsg:  "context name cannot contain path separators",
		},
		{
			name: "name with backslash",
			context: models.Context{
				Name:      "test\\context",
				StartTime: now,
				Status:    "active",
			},
			wantErr: true,
			errMsg:  "context name cannot contain path separators",
		},
		{
			name: "name too long",
			context: models.Context{
				Name:      string(make([]byte, 201)),
				StartTime: now,
				Status:    "active",
			},
			wantErr: true,
			errMsg:  "context name must be 200 characters or less",
		},
		{
			name: "invalid status",
			context: models.Context{
				Name:      "test-context",
				StartTime: now,
				Status:    "paused",
			},
			wantErr: true,
			errMsg:  "status must be 'active' or 'stopped'",
		},
		{
			name: "stopped without end time",
			context: models.Context{
				Name:      "test-context",
				StartTime: now,
				Status:    "stopped",
				EndTime:   nil,
			},
			wantErr: true,
			errMsg:  "stopped context must have an end time",
		},
		{
			name: "active with end time",
			context: models.Context{
				Name:      "test-context",
				StartTime: now,
				Status:    "active",
				EndTime:   &later,
			},
			wantErr: true,
			errMsg:  "active context cannot have an end time",
		},
		{
			name: "archived active context",
			context: models.Context{
				Name:       "test-context",
				StartTime:  now,
				Status:     "active",
				IsArchived: true,
			},
			wantErr: true,
			errMsg:  "cannot archive an active context",
		},
		{
			name: "archived stopped context",
			context: models.Context{
				Name:       "test-context",
				StartTime:  now,
				EndTime:    &later,
				Status:     "stopped",
				IsArchived: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.context.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileAssociationValidate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		fa      models.FileAssociation
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid file association",
			fa: models.FileAssociation{
				FilePath:  "/path/to/file.txt",
				Timestamp: now,
			},
			wantErr: false,
		},
		{
			name: "empty path",
			fa: models.FileAssociation{
				FilePath:  "",
				Timestamp: now,
			},
			wantErr: true,
			errMsg:  "file path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fa.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContextTransitionValidate(t *testing.T) {
	now := time.Now()
	ctx1 := "context-1"
	ctx2 := "context-2"

	tests := []struct {
		name    string
		tr      models.ContextTransition
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid start transition",
			tr: models.ContextTransition{
				Timestamp:      now,
				NewContext:     &ctx1,
				TransitionType: models.TransitionStart,
			},
			wantErr: false,
		},
		{
			name: "valid stop transition",
			tr: models.ContextTransition{
				Timestamp:       now,
				PreviousContext: &ctx1,
				TransitionType:  models.TransitionStop,
			},
			wantErr: false,
		},
		{
			name: "valid switch transition",
			tr: models.ContextTransition{
				Timestamp:       now,
				PreviousContext: &ctx1,
				NewContext:      &ctx2,
				TransitionType:  models.TransitionSwitch,
			},
			wantErr: false,
		},
		{
			name: "start without new context",
			tr: models.ContextTransition{
				Timestamp:      now,
				TransitionType: models.TransitionStart,
			},
			wantErr: true,
			errMsg:  "start transition must have new_context",
		},
		{
			name: "stop without previous context",
			tr: models.ContextTransition{
				Timestamp:      now,
				TransitionType: models.TransitionStop,
			},
			wantErr: true,
			errMsg:  "stop transition must have previous_context",
		},
		{
			name: "stop with new context",
			tr: models.ContextTransition{
				Timestamp:       now,
				PreviousContext: &ctx1,
				NewContext:      &ctx2,
				TransitionType:  models.TransitionStop,
			},
			wantErr: true,
			errMsg:  "stop transition cannot have new_context",
		},
		{
			name: "switch without previous context",
			tr: models.ContextTransition{
				Timestamp:      now,
				NewContext:     &ctx2,
				TransitionType: models.TransitionSwitch,
			},
			wantErr: true,
			errMsg:  "switch transition must have previous_context",
		},
		{
			name: "switch without new context",
			tr: models.ContextTransition{
				Timestamp:       now,
				PreviousContext: &ctx1,
				TransitionType:  models.TransitionSwitch,
			},
			wantErr: true,
			errMsg:  "switch transition must have new_context",
		},
		{
			name: "invalid transition type",
			tr: models.ContextTransition{
				Timestamp:      now,
				NewContext:     &ctx1,
				TransitionType: "pause",
			},
			wantErr: true,
			errMsg:  "invalid transition type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tr.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTouchEventToLogLine(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	touch := models.TouchEvent{
		Timestamp: timestamp,
	}

	expected := "2024-01-15T10:30:00Z"
	assert.Equal(t, expected, touch.ToLogLine())
}

func TestParseTouchLogLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{
			name:    "valid timestamp",
			line:    "2024-01-15T10:30:00Z",
			wantErr: false,
		},
		{
			name:    "invalid format",
			line:    "not a timestamp",
			wantErr: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			touch, err := models.ParseTouchLogLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, touch)
				assert.False(t, touch.Timestamp.IsZero())
			}
		})
	}
}
