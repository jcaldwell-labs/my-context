package core

import (
	"testing"
)

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name         string
		contextName  string
		patternParts []string
		want         bool
	}{
		{"empty pattern returns false", "any-context", []string{}, false},
		{"single wildcard matches all", "any-context", []string{"", ""}, true},
		{"no wildcard exact match", "test-context", []string{"test-context"}, true},
		{"no wildcard no match", "test-context", []string{"other-context"}, false},
		{"prefix wildcard", "test-context", []string{"test", ""}, true},
		{"prefix wildcard no match", "other-context", []string{"test", ""}, false},
		{"suffix wildcard", "test-context", []string{"", "context"}, true},
		{"suffix wildcard no match", "test-context", []string{"", "other"}, false},
		{"wildcard in middle", "test-context-name", []string{"test", "name"}, true},
		{"wildcard in middle no match", "test-context-name", []string{"test", "other"}, false},
		{"multiple wildcards", "test-context-full-name", []string{"test", "context", "name"}, true},
		{"prefix only match", "ps-cli: Phase 1", []string{"ps-cli", ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchesPattern(tt.contextName, tt.patternParts)
			if got != tt.want {
				t.Errorf("MatchesPattern(%q, %v) = %v, want %v", tt.contextName, tt.patternParts, got, tt.want)
			}
		})
	}
}

func TestExtractContextPrefix(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		want        string
	}{
		{"with colon and spaces", "my-project: feature-auth", "my-project"},
		{"no colon", "feature-login", "feature-login"},
		{"ps-cli style", "ps-cli: Phase 1", "ps-cli"},
		{"bug fix no colon", "bug-fix-123", "bug-fix-123"},
		{"multiple colons", "project: phase: detail", "project"},
		{"empty prefix", ": phase", ""},
		{"colon at end", "project:", "project"},
		{"unicode prefix", "日本語: タスク", "日本語"},
		{"spaces around colon", " project : phase", "project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractContextPrefix(tt.contextName)
			if got != tt.want {
				t.Errorf("extractContextPrefix(%q) = %q, want %q", tt.contextName, got, tt.want)
			}
		})
	}
}
