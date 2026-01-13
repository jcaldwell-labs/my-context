package models

import (
	"strings"
	"testing"
	"time"
)

func TestTransitionValidate(t *testing.T) {
	ctx1 := "context-1"
	ctx2 := "context-2"
	empty := ""

	tests := []struct {
		name    string
		ct      ContextTransition
		wantErr string
	}{
		{
			name: "valid start",
			ct: ContextTransition{
				Timestamp:      time.Now(),
				TransitionType: TransitionStart,
				NewContext:     &ctx1,
			},
			wantErr: "",
		},
		{
			name: "start without new context",
			ct: ContextTransition{
				Timestamp:      time.Now(),
				TransitionType: TransitionStart,
			},
			wantErr: "start transition must have new_context",
		},
		{
			name: "start with empty new context",
			ct: ContextTransition{
				Timestamp:      time.Now(),
				TransitionType: TransitionStart,
				NewContext:     &empty,
			},
			wantErr: "start transition must have new_context",
		},
		{
			name: "valid stop",
			ct: ContextTransition{
				Timestamp:       time.Now(),
				TransitionType:  TransitionStop,
				PreviousContext: &ctx1,
			},
			wantErr: "",
		},
		{
			name: "stop without previous context",
			ct: ContextTransition{
				Timestamp:      time.Now(),
				TransitionType: TransitionStop,
			},
			wantErr: "stop transition must have previous_context",
		},
		{
			name: "stop with new context",
			ct: ContextTransition{
				Timestamp:       time.Now(),
				TransitionType:  TransitionStop,
				PreviousContext: &ctx1,
				NewContext:      &ctx2,
			},
			wantErr: "stop transition cannot have new_context",
		},
		{
			name: "valid switch",
			ct: ContextTransition{
				Timestamp:       time.Now(),
				TransitionType:  TransitionSwitch,
				PreviousContext: &ctx1,
				NewContext:      &ctx2,
			},
			wantErr: "",
		},
		{
			name: "switch without previous",
			ct: ContextTransition{
				Timestamp:      time.Now(),
				TransitionType: TransitionSwitch,
				NewContext:     &ctx2,
			},
			wantErr: "switch transition must have previous_context",
		},
		{
			name: "switch without new",
			ct: ContextTransition{
				Timestamp:       time.Now(),
				TransitionType:  TransitionSwitch,
				PreviousContext: &ctx1,
			},
			wantErr: "switch transition must have new_context",
		},
		{
			name: "invalid transition type",
			ct: ContextTransition{
				Timestamp:      time.Now(),
				TransitionType: "invalid",
			},
			wantErr: "invalid transition type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ct.Validate()
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

func TestTransitionToLogLine(t *testing.T) {
	timestamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	ctx1 := "context-1"
	ctx2 := "context-2"

	tests := []struct {
		name string
		ct   ContextTransition
		want string
	}{
		{
			name: "start transition",
			ct: ContextTransition{
				Timestamp:      timestamp,
				TransitionType: TransitionStart,
				NewContext:     &ctx1,
			},
			want: "2025-01-15T10:30:00Z|NULL|context-1|start",
		},
		{
			name: "stop transition",
			ct: ContextTransition{
				Timestamp:       timestamp,
				TransitionType:  TransitionStop,
				PreviousContext: &ctx1,
			},
			want: "2025-01-15T10:30:00Z|context-1|NULL|stop",
		},
		{
			name: "switch transition",
			ct: ContextTransition{
				Timestamp:       timestamp,
				TransitionType:  TransitionSwitch,
				PreviousContext: &ctx1,
				NewContext:      &ctx2,
			},
			want: "2025-01-15T10:30:00Z|context-1|context-2|switch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ct.ToLogLine()
			if got != tt.want {
				t.Errorf("ToLogLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseTransitionLogLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantErr  bool
		wantType TransitionType
	}{
		{
			name:     "valid start",
			line:     "2025-01-15T10:30:00Z|NULL|context-1|start",
			wantErr:  false,
			wantType: TransitionStart,
		},
		{
			name:     "valid stop",
			line:     "2025-01-15T10:30:00Z|context-1|NULL|stop",
			wantErr:  false,
			wantType: TransitionStop,
		},
		{
			name:     "valid switch",
			line:     "2025-01-15T10:30:00Z|context-1|context-2|switch",
			wantErr:  false,
			wantType: TransitionSwitch,
		},
		{
			name:    "invalid format",
			line:    "2025-01-15T10:30:00Z|context-1|start",
			wantErr: true,
		},
		{
			name:    "invalid timestamp",
			line:    "not-a-timestamp|NULL|context-1|start",
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
			ct, err := ParseTransitionLogLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTransitionLogLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ct.TransitionType != tt.wantType {
				t.Errorf("ParseTransitionLogLine() type = %q, want %q", ct.TransitionType, tt.wantType)
			}
		})
	}
}

func TestTransitionRoundtrip(t *testing.T) {
	timestamp := time.Now().UTC().Truncate(time.Second)
	ctx1 := "context-1"
	ctx2 := "context-2"

	original := ContextTransition{
		Timestamp:       timestamp,
		TransitionType:  TransitionSwitch,
		PreviousContext: &ctx1,
		NewContext:      &ctx2,
	}

	line := original.ToLogLine()
	parsed, err := ParseTransitionLogLine(line)

	if err != nil {
		t.Fatalf("Roundtrip failed to parse: %v", err)
	}

	if parsed.TransitionType != original.TransitionType {
		t.Errorf("Roundtrip type mismatch: %q vs %q", parsed.TransitionType, original.TransitionType)
	}

	if *parsed.PreviousContext != *original.PreviousContext {
		t.Errorf("Roundtrip previous mismatch: %q vs %q", *parsed.PreviousContext, *original.PreviousContext)
	}

	if *parsed.NewContext != *original.NewContext {
		t.Errorf("Roundtrip new mismatch: %q vs %q", *parsed.NewContext, *original.NewContext)
	}
}
