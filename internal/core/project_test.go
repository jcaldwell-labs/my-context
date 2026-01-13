package core

import (
	"reflect"
	"testing"
)

func TestExtractProjectName(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		want        string
	}{
		{"simple project with phase", "ps-cli: Phase 1", "ps-cli"},
		{"project with spaces", "garden: Planning Session", "garden"},
		{"no colon - full name", "StandaloneContext", "StandaloneContext"},
		{"multiple colons - first only", "project: Phase 1: Subphase A", "project"},
		{"empty string", "", ""},
		{"just colon", ":", ""},
		{"spaces around colon", " project : phase ", "project"},
		{"unicode project name", "日本語プロジェクト: タスク", "日本語プロジェクト"},
		{"colon at end", "project:", "project"},
		{"colon at start", ":phase", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractProjectName(tt.contextName)
			if got != tt.want {
				t.Errorf("ExtractProjectName(%q) = %q, want %q", tt.contextName, got, tt.want)
			}
		})
	}
}

func TestFilterContextsByProject(t *testing.T) {
	allContexts := []string{
		"ps-cli: Phase 1",
		"ps-cli: Phase 2",
		"PS-CLI: Phase 3", // Different case
		"garden: Planning",
		"garden: Implementation",
		"standalone",
	}

	tests := []struct {
		name        string
		projectName string
		want        []string
	}{
		{
			name:        "filter ps-cli (case insensitive)",
			projectName: "ps-cli",
			want:        []string{"ps-cli: Phase 1", "ps-cli: Phase 2", "PS-CLI: Phase 3"},
		},
		{
			name:        "filter PS-CLI uppercase",
			projectName: "PS-CLI",
			want:        []string{"ps-cli: Phase 1", "ps-cli: Phase 2", "PS-CLI: Phase 3"},
		},
		{
			name:        "filter garden",
			projectName: "garden",
			want:        []string{"garden: Planning", "garden: Implementation"},
		},
		{
			name:        "filter standalone",
			projectName: "standalone",
			want:        []string{"standalone"},
		},
		{
			name:        "empty filter returns all",
			projectName: "",
			want:        allContexts,
		},
		{
			name:        "non-existent project",
			projectName: "nonexistent",
			want:        nil,
		},
		{
			name:        "project with spaces in filter",
			projectName: " ps-cli ",
			want:        []string{"ps-cli: Phase 1", "ps-cli: Phase 2", "PS-CLI: Phase 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterContextsByProject(allContexts, tt.projectName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterContextsByProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterByProject(t *testing.T) {
	// FilterByProject is an alias for FilterContextsByProject
	contexts := []string{"a: 1", "a: 2", "b: 1"}
	got := FilterByProject(contexts, "a")
	want := []string{"a: 1", "a: 2"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FilterByProject() = %v, want %v", got, want)
	}
}

func TestExtractProjectMetadata(t *testing.T) {
	contexts := []string{
		"ps-cli: Phase 1",
		"ps-cli: Phase 2",
		"PS-CLI: Phase 3", // Same project, different case
		"garden: Planning",
		"standalone",
	}

	got := ExtractProjectMetadata(contexts)

	// Check ps-cli group (should be lowercase key, includes all 3)
	psCli, ok := got["ps-cli"]
	if !ok {
		t.Fatal("Expected ps-cli key in metadata")
	}
	if psCli.ContextCount != 3 {
		t.Errorf("ps-cli ContextCount = %d, want 3", psCli.ContextCount)
	}
	if len(psCli.ContextNames) != 3 {
		t.Errorf("ps-cli ContextNames length = %d, want 3", len(psCli.ContextNames))
	}

	// Check garden group
	garden, ok := got["garden"]
	if !ok {
		t.Fatal("Expected garden key in metadata")
	}
	if garden.ContextCount != 1 {
		t.Errorf("garden ContextCount = %d, want 1", garden.ContextCount)
	}

	// Check standalone (no colon, full name is project)
	standalone, ok := got["standalone"]
	if !ok {
		t.Fatal("Expected standalone key in metadata")
	}
	if standalone.ContextCount != 1 {
		t.Errorf("standalone ContextCount = %d, want 1", standalone.ContextCount)
	}

	// Total should be 3 unique projects
	if len(got) != 3 {
		t.Errorf("Total projects = %d, want 3", len(got))
	}
}

func TestExtractProjectMetadataEmpty(t *testing.T) {
	got := ExtractProjectMetadata([]string{})
	if len(got) != 0 {
		t.Errorf("Expected empty map for empty input, got %d entries", len(got))
	}
}

func TestExtractProjectMetadataPreservesOriginalCasing(t *testing.T) {
	contexts := []string{"MyProject: Task 1"}
	got := ExtractProjectMetadata(contexts)

	meta, ok := got["myproject"]
	if !ok {
		t.Fatal("Expected myproject key")
	}

	// ProjectName should preserve original casing
	if meta.ProjectName != "MyProject" {
		t.Errorf("ProjectName = %q, want %q (original casing)", meta.ProjectName, "MyProject")
	}
}
