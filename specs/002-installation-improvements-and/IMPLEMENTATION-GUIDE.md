# Sprint 2 Implementation Guide

**Feature**: 002-installation-improvements-and
**Generated**: 2025-10-05
**Total Tasks**: 41
**Estimated Effort**: 5.5 days

---

## Overview

41 tasks organized into 6 phases following TDD approach. This guide provides step-by-step instructions for implementing Sprint 2 features.

---

## 🚀 Phase 3.1: Setup & Build Infrastructure (0.5 days)

**Goal**: Prepare build system and foundational infrastructure

### T001: GitHub Actions Workflow
**File**: `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
            output: my-context-linux-amd64
          - goos: windows
            goarch: amd64
            output: my-context-windows-amd64.exe
          - goos: darwin
            goarch: amd64
            output: my-context-darwin-amd64
          - goos: darwin
            goarch: arm64
            output: my-context-darwin-arm64

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        run: |
          CGO_ENABLED=0 GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} \
          go build -ldflags "-X main.Version=${{ github.ref_name }} \
            -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
            -X main.GitCommit=$(git rev-parse --short HEAD)" \
            -o ${{ matrix.output }} ./cmd/my-context/

      - name: Generate Checksum
        run: sha256sum ${{ matrix.output }} > ${{ matrix.output }}.sha256

      - name: Upload Release Asset
        uses: actions/upload-artifact@v3
        with:
          name: ${{ matrix.output }}
          path: |
            ${{ matrix.output }}
            ${{ matrix.output }}.sha256
```

**After completion**: Mark `[ ]` → `[x]` in `specs/002-installation-improvements-and/tasks.md` line 45

---

### T002: Local Build Script
**File**: `scripts/build-all.sh`

```bash
#!/usr/bin/env bash
# Multi-platform build script for local development

set -e

VERSION=${VERSION:-dev}
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

echo "Building my-context ${VERSION} (${GIT_COMMIT})..."

mkdir -p bin/

# Linux amd64
echo "→ Building linux/amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/my-context-linux-amd64 ./cmd/my-context/

# Windows amd64
echo "→ Building windows/amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/my-context-windows-amd64.exe ./cmd/my-context/

# macOS amd64
echo "→ Building darwin/amd64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/my-context-darwin-amd64 ./cmd/my-context/

# macOS arm64
echo "→ Building darwin/arm64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/my-context-darwin-arm64 ./cmd/my-context/

echo "✓ Build complete. Binaries in bin/"
ls -lh bin/
```

**Test**: Run `chmod +x scripts/build-all.sh && ./scripts/build-all.sh` and verify 4 binaries created in `bin/`

**After completion**: Mark tasks.md line 47 as `[x]`

---

### T003 [P]: Version Info in main.go
**File**: `cmd/my-context/main.go`

**Add these variables at package level** (before rootCmd):
```go
var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)
```

**Update rootCmd** (in init() or where rootCmd is defined):
```go
rootCmd.Version = fmt.Sprintf("%s (built: %s, commit: %s)", Version, BuildTime, GitCommit)
```

**Test**: Build and run `go build -o my-context.exe ./cmd/my-context/ && ./my-context.exe --version`

**After completion**: Mark tasks.md line 49 as `[x]`

---

### T004 [P]: Troubleshooting Documentation
**File**: `docs/TROUBLESHOOTING.md`

**Status**: ✅ Already created (git status shows `AM docs/TROUBLESHOOTING.md`)

**Verify**: File exists and contains platform-specific troubleshooting

**After completion**: Mark tasks.md line 51 as `[x]`

---

## ⚠️ Phase 3.2: Tests First (TDD) - BLOCKING GATE (1.0 day)

**CRITICAL**: All tests in this phase MUST be written and MUST FAIL before proceeding to Phase 3.3

**Parallel execution**: All 8 test files (T005-T012) can be written simultaneously

### Test Writing Strategy

1. Create test file with imports and helper functions
2. Write test cases per contract specification
3. Run tests - **verify they FAIL** (no implementation exists yet)
4. Mark task complete in tasks.md

---

### T005 [P]: Export Command Tests
**File**: `tests/integration/export_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/export_test.go`)

**Verify tests exist**:
- TestExportSingleContext
- TestExportWithCustomPath
- TestExportAllContexts
- TestExportNonExistent
- TestExportMarkdownFormat
- TestExportJSONOutput

**Run**: `go test ./tests/integration/export_test.go -v`
**Expected**: All tests FAIL (no export command exists yet)

**After completion**: Mark tasks.md line 61 as `[x]`

---

### T006 [P]: Archive Command Tests
**File**: `tests/integration/archive_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/archive_test.go`)

**Run**: `go test ./tests/integration/archive_test.go -v`
**Expected**: All tests FAIL

**After completion**: Mark tasks.md line 69 as `[x]`

---

### T007 [P]: Delete Command Tests
**File**: `tests/integration/delete_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/delete_test.go`)

**Run**: `go test ./tests/integration/delete_test.go -v`
**Expected**: All tests FAIL

**After completion**: Mark tasks.md line 77 as `[x]`

---

### T008 [P]: List Enhanced Tests
**File**: `tests/integration/list_enhanced_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/list_enhanced_test.go`)

**Run**: `go test ./tests/integration/list_enhanced_test.go -v`
**Expected**: All tests FAIL

**After completion**: Mark tasks.md line 85 as `[x]`

---

### T009 [P]: Project Filter Tests
**File**: `tests/integration/project_filter_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/project_filter_test.go`)

**Run**: `go test ./tests/integration/project_filter_test.go -v`
**Expected**: All tests FAIL

**After completion**: Mark tasks.md line 94 as `[x]`

---

### T010 [P]: Bug Fixes Tests
**File**: `tests/integration/bug_fixes_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/bug_fixes_test.go`)

**Verify tests cover**:
- Note with $ character preservation
- History showing "(none)" instead of "NULL"
- Other special characters (!, @, #, etc.)

**Run**: `go test ./tests/integration/bug_fixes_test.go -v`
**Expected**: All tests FAIL

**After completion**: Mark tasks.md line 100 as `[x]`

---

### T011 [P]: Project Parser Unit Tests
**File**: `tests/unit/project_parser_test.go`

**Create new file**:
```go
package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/my-context-copilot/internal/core"
)

func TestExtractProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with colon", "ps-cli: Phase 1", "ps-cli"},
		{"no colon", "Standalone", "Standalone"},
		{"multiple colons", "project: sub: detail", "project"},
		{"with whitespace", "  project : phase  ", "project"},
		{"empty string", "", ""},
		{"only colon", ":", ""},
		{"trailing colon", "project:", "project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.ExtractProjectName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterByProject(t *testing.T) {
	// Add test cases for FilterByProject function
	// Test case-insensitive matching
	// Test contexts with and without colons
}
```

**Run**: `go test ./tests/unit/project_parser_test.go -v`
**Expected**: Test FAILS (ExtractProjectName doesn't exist yet)

**After completion**: Mark tasks.md line 107 as `[x]`

---

### T012 [P]: Backward Compatibility Tests
**File**: `tests/integration/backward_compat_test.go`

**Status**: ✅ Already created (git status shows `AM tests/integration/backward_compat_test.go`)

**Verify tests cover**:
- Loading Sprint 1 meta.json (without is_archived field)
- Verifying Sprint 1 contexts work with Sprint 2 code
- New features work on old contexts

**Run**: `go test ./tests/integration/backward_compat_test.go -v`
**Expected**: Tests FAIL

**After completion**: Mark tasks.md line 114 as `[x]`

---

**✅ GATE CHECK - Before proceeding to Phase 3.3**:

Run full test suite:
```bash
go test ./tests/integration/... -v
go test ./tests/unit/... -v
```

**Verify**:
- [ ] All test files created (T005-T012)
- [ ] All tests FAIL (no implementation exists)
- [ ] No compilation errors (tests are valid Go code)

**Do NOT proceed to Phase 3.3 until all tests are failing!**

---

## 🔨 Phase 3.3: Core Implementation (2.0 days)

**Only proceed after Phase 3.2 complete with all tests failing**

### Model Layer (Parallel: T013-T015)

#### T013 [P]: Add IsArchived to Context Model
**File**: `internal/models/context.go`

**Add field to Context struct**:
```go
type Context struct {
	Name             string     `json:"name"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          *time.Time `json:"end_time,omitempty"`
	Status           string     `json:"status"`
	SubdirectoryPath string     `json:"subdirectory_path"`
	IsArchived       bool       `json:"is_archived,omitempty"` // NEW
}
```

**Add validation method**:
```go
// CanArchive checks if context can be archived
func (c *Context) CanArchive() error {
	if c.Status == "active" {
		return fmt.Errorf("cannot archive active context")
	}
	return nil
}
```

**After completion**: Mark tasks.md line 125 as `[x]`

---

#### T014 [P]: Project Extraction Logic
**File**: `internal/core/project.go` (NEW)

```go
package core

import (
	"strings"

	"github.com/yourusername/my-context-copilot/internal/models"
)

// ExtractProjectName parses "project: phase" format
// Returns text before first colon, or full name if no colon
func ExtractProjectName(contextName string) string {
	parts := strings.SplitN(contextName, ":", 2)
	return strings.TrimSpace(parts[0])
}

// FilterByProject returns contexts matching project name (case-insensitive)
func FilterByProject(contexts []models.Context, projectName string) []models.Context {
	filtered := []models.Context{}
	for _, ctx := range contexts {
		if strings.EqualFold(ExtractProjectName(ctx.Name), projectName) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}
```

**Test**: `go test ./tests/unit/project_parser_test.go -v` should now **PASS**

**After completion**: Mark tasks.md line 127 as `[x]`

---

#### T015 [P]: Markdown Export Formatter
**File**: `internal/output/markdown.go` (NEW)

```go
package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/my-context-copilot/internal/models"
)

// Note represents a timestamped note entry
type Note struct {
	Timestamp time.Time
	Text      string
}

// FileAssoc represents a file association entry
type FileAssoc struct {
	Timestamp time.Time
	Path      string
}

// FormatExport generates markdown representation of context
func FormatExport(ctx models.Context, notes []Note, files []FileAssoc, touches []time.Time, version string) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# Context: %s\n\n", ctx.Name))
	sb.WriteString(fmt.Sprintf("**Started**: %s\n", ctx.StartTime.Local().Format(time.RFC3339)))

	if ctx.EndTime != nil {
		sb.WriteString(fmt.Sprintf("**Ended**: %s\n", ctx.EndTime.Local().Format(time.RFC3339)))
		duration := ctx.EndTime.Sub(ctx.StartTime)
		sb.WriteString(fmt.Sprintf("**Duration**: %s\n", formatDuration(duration)))
	} else {
		sb.WriteString("**Ended**: Active\n")
		sb.WriteString("**Duration**: Active\n")
	}

	status := ctx.Status
	if ctx.IsArchived {
		status = "archived"
	}
	sb.WriteString(fmt.Sprintf("**Status**: %s\n\n", status))

	// Notes section
	sb.WriteString(fmt.Sprintf("## Notes (%d)\n\n", len(notes)))
	if len(notes) > 0 {
		for _, note := range notes {
			timeStr := note.Timestamp.Local().Format("15:04")
			sb.WriteString(fmt.Sprintf("- `%s` %s\n", timeStr, note.Text))
		}
	} else {
		sb.WriteString("(No notes)\n")
	}

	// Files section
	sb.WriteString(fmt.Sprintf("\n## Associated Files (%d)\n\n", len(files)))
	if len(files) > 0 {
		for _, file := range files {
			sb.WriteString(fmt.Sprintf("- %s\n", file.Path))
			sb.WriteString(fmt.Sprintf("  Added: %s\n", file.Timestamp.Local().Format(time.RFC3339)))
		}
	} else {
		sb.WriteString("(No files)\n")
	}

	// Activity section
	sb.WriteString(fmt.Sprintf("\n## Activity\n\n- %d touch events\n", len(touches)))
	if len(touches) > 0 {
		lastTouch := touches[len(touches)-1]
		sb.WriteString(fmt.Sprintf("- Last activity: %s\n", lastTouch.Local().Format(time.RFC3339)))
	} else {
		sb.WriteString("- No recorded activity\n")
	}

	// Footer
	sb.WriteString(fmt.Sprintf("\n---\n*Exported: %s by my-context v%s*\n",
		time.Now().Local().Format(time.RFC3339), version))

	return sb.String()
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours >= 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
```

**After completion**: Mark tasks.md line 129 as `[x]`

---

### Core Logic Layer (Sequential: T016 → T017, then parallel T018-T020)

#### T016: Core Context Operations
**File**: `internal/core/context.go`

**Add these functions**:

```go
// ArchiveContext marks a context as archived
func ArchiveContext(contextName string) error {
	ctx, err := LoadContext(contextName)
	if err != nil {
		return err
	}

	if err := ctx.CanArchive(); err != nil {
		return err
	}

	if ctx.IsArchived {
		return fmt.Errorf("already archived")
	}

	ctx.IsArchived = true
	return SaveContext(ctx)
}

// DeleteContext removes a context directory permanently
func DeleteContext(contextName string, force bool) error {
	ctx, err := LoadContext(contextName)
	if err != nil {
		return err
	}

	// Check if active
	activeCtx, _ := GetActiveContext()
	if activeCtx != nil && activeCtx.Name == contextName {
		return fmt.Errorf("cannot delete active context")
	}

	// Confirmation prompt (unless --force)
	if !force {
		if !confirmDeletion(ctx) {
			return fmt.Errorf("deletion cancelled")
		}
	}

	// Remove directory
	return os.RemoveAll(ctx.SubdirectoryPath)
}

// ExportContext generates markdown export
func ExportContext(contextName, outputPath string) (string, error) {
	ctx, err := LoadContext(contextName)
	if err != nil {
		return "", err
	}

	// Read logs
	notes, _ := ReadNotes(ctx.SubdirectoryPath)
	files, _ := ReadFiles(ctx.SubdirectoryPath)
	touches, _ := ReadTouches(ctx.SubdirectoryPath)

	// Generate markdown
	markdown := output.FormatExport(ctx, notes, files, touches, Version)

	// Determine output path
	if outputPath == "" {
		outputPath = SanitizeFilename(ctx.Name) + ".md"
	}

	// Write file
	return outputPath, WriteMarkdown(outputPath, markdown)
}

func confirmDeletion(ctx models.Context) bool {
	noteCount, _ := countLogEntries(filepath.Join(ctx.SubdirectoryPath, "notes.log"))
	fileCount, _ := countLogEntries(filepath.Join(ctx.SubdirectoryPath, "files.log"))
	touchCount, _ := countLogEntries(filepath.Join(ctx.SubdirectoryPath, "touch.log"))

	fmt.Printf("⚠️  This will permanently delete context \"%s\" and all associated data.\n", ctx.Name)
	fmt.Printf("   - Notes: %d\n", noteCount)
	fmt.Printf("   - Files: %d\n", fileCount)
	fmt.Printf("   - Touch events: %d\n\n", touchCount)
	fmt.Printf("Continue? (y/N): ")

	var response string
	fmt.Scanln(&response)

	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func countLogEntries(logPath string) (int, error) {
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0, nil
	}
	return len(lines), nil
}
```

**After completion**: Mark tasks.md line 137 as `[x]`

---

#### T017: Enhanced ListContexts (DEPENDS ON T016)
**File**: `internal/core/context.go`

**Add struct for filters**:
```go
type ListFilters struct {
	Project    string
	Search     string
	Limit      int
	ShowAll    bool
	Archived   bool
	ActiveOnly bool
}
```

**Modify or create ListContexts function**:
```go
func ListContexts(filters ListFilters) ([]models.Context, int, error) {
	allContexts, err := loadAllContexts()
	if err != nil {
		return nil, 0, err
	}

	// Apply filters
	filtered := allContexts

	// Project filter
	if filters.Project != "" {
		filtered = FilterByProject(filtered, filters.Project)
	}

	// Search filter
	if filters.Search != "" {
		filtered = filterBySearch(filtered, filters.Search)
	}

	// Archive status filter
	if filters.ActiveOnly {
		filtered = filterActive(filtered)
	} else if filters.Archived {
		filtered = filterArchived(filtered, true)
	} else {
		// Default: hide archived
		filtered = filterArchived(filtered, false)
	}

	// Sort by start_time descending (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].StartTime.After(filtered[j].StartTime)
	})

	totalCount := len(filtered)

	// Apply limit
	if !filters.ShowAll && filters.Limit > 0 && len(filtered) > filters.Limit {
		filtered = filtered[:filters.Limit]
	}

	return filtered, totalCount, nil
}

func filterBySearch(contexts []models.Context, term string) []models.Context {
	result := []models.Context{}
	lowerTerm := strings.ToLower(term)
	for _, ctx := range contexts {
		if strings.Contains(strings.ToLower(ctx.Name), lowerTerm) {
			result = append(result, ctx)
		}
	}
	return result
}

func filterActive(contexts []models.Context) []models.Context {
	for _, ctx := range contexts {
		if ctx.Status == "active" {
			return []models.Context{ctx}
		}
	}
	return []models.Context{}
}

func filterArchived(contexts []models.Context, showArchived bool) []models.Context {
	result := []models.Context{}
	for _, ctx := range contexts {
		if ctx.IsArchived == showArchived {
			result = append(result, ctx)
		}
	}
	return result
}

func loadAllContexts() ([]models.Context, error) {
	// Implementation depends on existing code structure
	// Load all meta.json files from ~/.my-context/*/meta.json
	// Return array of Context objects
}
```

**After completion**: Mark tasks.md line 139 as `[x]`

---

#### T018 [P]: Storage Helpers for Export
**File**: `internal/core/storage.go`

**Add functions**:
```go
func WriteMarkdown(path, content string) error {
	// Create parent dirs
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		// Prompt for overwrite
		fmt.Printf("File exists: %s\nOverwrite? (y/N): ", path)
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("export cancelled")
		}
	}

	// Atomic write
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func SanitizeFilename(name string) string {
	// Replace special chars with underscores
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}
```

**After completion**: Mark tasks.md line 141 as `[x]`

---

#### T019 [P]: Fix NULL Display in History
**File**: `internal/output/human.go`

**Find where history/transitions are formatted** (look for "NULL" string):

**Add helper function**:
```go
func formatContextName(name string) string {
	if name == "" || name == "NULL" {
		return "(none)"
	}
	return name
}
```

**Use in transition formatting**:
```go
// Example: when printing transitions
previousCtx := formatContextName(transition.PreviousContext)
nextCtx := formatContextName(transition.NextContext)
fmt.Printf("%s → %s\n", previousCtx, nextCtx)
```

**Test**: Run bug_fixes_test.go - history tests should now PASS

**After completion**: Mark tasks.md line 143 as `[x]`

---

#### T020 [P]: Log Reading Helpers
**File**: `internal/core/storage.go`

**Add functions**:
```go
func ReadNotes(contextDir string) ([]output.Note, error) {
	notesPath := filepath.Join(contextDir, "notes.log")
	data, err := os.ReadFile(notesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []output.Note{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	notes := []output.Note{}

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		timestamp, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			continue
		}
		notes = append(notes, output.Note{
			Timestamp: timestamp,
			Text:      parts[1],
		})
	}

	return notes, nil
}

func ReadFiles(contextDir string) ([]output.FileAssoc, error) {
	filesPath := filepath.Join(contextDir, "files.log")
	data, err := os.ReadFile(filesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []output.FileAssoc{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	files := []output.FileAssoc{}

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		timestamp, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			continue
		}
		files = append(files, output.FileAssoc{
			Timestamp: timestamp,
			Path:      parts[1],
		})
	}

	return files, nil
}

func ReadTouches(contextDir string) ([]time.Time, error) {
	touchPath := filepath.Join(contextDir, "touch.log")
	data, err := os.ReadFile(touchPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []time.Time{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	touches := []time.Time{}

	for _, line := range lines {
		if line == "" {
			continue
		}
		timestamp, err := time.Parse(time.RFC3339, line)
		if err != nil {
			continue
		}
		touches = append(touches, timestamp)
	}

	return touches, nil
}
```

**After completion**: Mark tasks.md line 145 as `[x]`

---

### Command Layer (Parallel: T021-T026)

All 6 command files can be implemented in parallel since they don't depend on each other.

#### T021 [P]: Export Command
**File**: `internal/commands/export.go` (NEW)

```go
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/my-context-copilot/internal/core"
)

var (
	exportToPath string
	exportAll    bool
)

var exportCmd = &cobra.Command{
	Use:     "export <context-name>",
	Aliases: []string{"e"},
	Short:   "Export context to markdown",
	Long: `Export context data to a markdown file for sharing or archival.

Examples:
  my-context export "ps-cli: Phase 1"
  my-context export "Phase 1" --to reports/phase-1.md
  my-context export --all --to exports/`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if exportAll {
			return exportAllContexts(exportToPath)
		}

		if len(args) == 0 {
			return fmt.Errorf("context name required or use --all to export all contexts")
		}

		outputPath, err := core.ExportContext(args[0], exportToPath)
		if err != nil {
			return err
		}

		fmt.Printf("Exported context \"%s\" to %s\n", args[0], outputPath)
		return nil
	},
}

func exportAllContexts(outputDir string) error {
	// Get all contexts (including archived)
	filters := core.ListFilters{ShowAll: true}
	contexts, _, err := core.ListContexts(filters)
	if err != nil {
		return err
	}

	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	fmt.Printf("Exporting %d contexts...\n", len(contexts))

	for _, ctx := range contexts {
		outputPath := filepath.Join(outputDir, core.SanitizeFilename(ctx.Name)+".md")
		if _, err := core.ExportContext(ctx.Name, outputPath); err != nil {
			fmt.Printf("  ✗ Failed: %s (%v)\n", ctx.Name, err)
			continue
		}
		fmt.Printf("  ✓ %s\n", filepath.Base(outputPath))
	}

	fmt.Printf("Exported %d contexts to %s\n", len(contexts), outputDir)
	return nil
}

func init() {
	exportCmd.Flags().StringVar(&exportToPath, "to", "", "Output file path")
	exportCmd.Flags().BoolVar(&exportAll, "all", false, "Export all contexts")
	rootCmd.AddCommand(exportCmd)
}
```

**Register in root**: Ensure `internal/commands/export.go` is imported in your root command file

**Test**: `go test ./tests/integration/export_test.go -v` should now **PASS**

**After completion**: Mark tasks.md line 149 as `[x]`

---

#### T022 [P]: Archive Command
**File**: `internal/commands/archive.go` (NEW)

```go
package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/my-context-copilot/internal/core"
)

var archiveCmd = &cobra.Command{
	Use:     "archive <context-name>",
	Aliases: []string{"a"},
	Short:   "Mark context as archived",
	Long: `Mark a context as archived (completed work).

Archived contexts are hidden from default list view but remain accessible
with 'my-context show' and 'my-context list --archived'.

Examples:
  my-context archive "ps-cli: Phase 1"
  my-context a "Completed Project"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := core.ArchiveContext(args[0]); err != nil {
			if strings.Contains(err.Error(), "already archived") {
				fmt.Printf("Context \"%s\" is already archived.\n", args[0])
				return nil
			}
			return err
		}

		fmt.Printf("Archived context: %s\n", args[0])
		fmt.Println("Use 'my-context list --archived' to see archived contexts.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
```

**Test**: Archive tests should pass

**After completion**: Mark tasks.md line 156 as `[x]`

---

#### T023 [P]: Delete Command
**File**: `internal/commands/delete.go` (NEW)

```go
package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/my-context-copilot/internal/core"
)

var (
	deleteForce bool
)

var deleteCmd = &cobra.Command{
	Use:     "delete <context-name>",
	Aliases: []string{"d"},
	Short:   "Permanently delete a context",
	Long: `Permanently delete a context and all associated data.

This action cannot be undone. You will be prompted for confirmation
unless --force flag is used.

Examples:
  my-context delete "Test Context"
  my-context delete "Old Work" --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := core.DeleteContext(args[0], deleteForce)
		if err != nil {
			if strings.Contains(err.Error(), "cancelled") {
				fmt.Println("Deletion cancelled.")
				return nil
			}
			return err
		}

		fmt.Printf("Deleted context: %s\n", args[0])
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}
```

**Test**: Delete tests should pass

**After completion**: Mark tasks.md line 162 as `[x]`

---

#### T024 [P]: Enhanced List Command
**File**: `internal/commands/list.go`

**Add flags**:
```go
var (
	listProject    string
	listLimit      int
	listSearch     string
	listAll        bool
	listArchived   bool
	listActiveOnly bool
)
```

**Update init()**:
```go
func init() {
	listCmd.Flags().StringVar(&listProject, "project", "", "Filter by project name")
	listCmd.Flags().IntVar(&listLimit, "limit", 10, "Max contexts to show")
	listCmd.Flags().StringVar(&listSearch, "search", "", "Search context names")
	listCmd.Flags().BoolVar(&listAll, "all", false, "Show all contexts")
	listCmd.Flags().BoolVar(&listArchived, "archived", false, "Show only archived contexts")
	listCmd.Flags().BoolVar(&listActiveOnly, "active-only", false, "Show only active context")
	rootCmd.AddCommand(listCmd)
}
```

**Modify RunE**:
```go
RunE: func(cmd *cobra.Command, args []string) error {
	// Validate conflicting flags
	if listArchived && listActiveOnly {
		return fmt.Errorf("cannot use --archived and --active-only together")
	}

	if listLimit < 0 {
		return fmt.Errorf("--limit must be a positive number")
	}

	filters := core.ListFilters{
		Project:    listProject,
		Search:     listSearch,
		Limit:      listLimit,
		ShowAll:    listAll,
		Archived:   listArchived,
		ActiveOnly: listActiveOnly,
	}

	contexts, totalCount, err := core.ListContexts(filters)
	if err != nil {
		return err
	}

	if len(contexts) == 0 {
		if filters.Project != "" {
			fmt.Printf("No contexts found for project \"%s\"\n", filters.Project)
		} else if filters.Search != "" {
			fmt.Printf("No contexts matching \"%s\"\n", filters.Search)
		} else if filters.ActiveOnly {
			fmt.Println("No active context")
			fmt.Println("Run 'my-context start <name>' to create one.")
		} else {
			fmt.Println("No contexts found")
		}
		return nil
	}

	// Display contexts (use existing display logic or create new)
	for _, ctx := range contexts {
		// Format and print each context
		// Reuse existing list display code
	}

	// Show truncation message
	if !listAll && len(contexts) < totalCount {
		fmt.Printf("\nShowing %d of %d contexts. Use --all to see all.\n",
			len(contexts), totalCount)
	}

	return nil
},
```

**Test**: List enhanced tests should pass

**After completion**: Mark tasks.md line 169 as `[x]`

---

#### T025 [P]: Start Command with --project Flag
**File**: `internal/commands/start.go`

**Add flag variable**:
```go
var startProject string
```

**Update init()**:
```go
func init() {
	startCmd.Flags().StringVar(&startProject, "project", "", "Project name prefix")
	// ... existing flags
}
```

**Modify RunE** (where context name is constructed):
```go
RunE: func(cmd *cobra.Command, args []string) error {
	contextName := args[0]

	// Combine project + name if --project provided
	if startProject != "" {
		contextName = strings.TrimSpace(startProject) + ": " + strings.TrimSpace(contextName)
	}

	// Continue with existing CreateContext logic...
	return core.CreateContext(contextName)
},
```

**Test**: Project filter tests should pass

**After completion**: Mark tasks.md line 179 as `[x]`

---

#### T026 [P]: Fix $ Character Bug in Notes
**File**: `internal/commands/note.go`

**Review current implementation**: Find where note text is captured from args

**Ensure no shell expansion or escaping**:
```go
RunE: func(cmd *cobra.Command, args []string) error {
	// Correct approach - join args as-is without interpretation
	noteText := strings.Join(args, " ")

	// Store raw text (no fmt.Sprintf that might interpret $)
	return core.AddNote(noteText)
},
```

**In storage layer** (check `internal/core/storage.go` or wherever notes are written):
```go
func AddNote(text string) error {
	timestamp := time.Now().Format(time.RFC3339)
	// Use raw text - no string interpolation
	logEntry := timestamp + "|" + text + "\n"
	// Append to notes.log
	return appendToFile(notesLogPath, logEntry)
}
```

**Test**: Run bug_fixes_test.go - $ character tests should PASS

**After completion**: Mark tasks.md line 185 as `[x]`

---

**✅ Checkpoint**: Run all tests
```bash
go test ./tests/integration/... -v
go test ./tests/unit/... -v
```

**Expected**: All tests should **PASS** now (GREEN)

---

## 📦 Phase 3.4: Build & Installation Scripts (1.0 day)

Can run in parallel (T027-T031)

### T027: Enhanced install.sh
**File**: `scripts/install.sh`

```bash
#!/usr/bin/env bash
# Installation script for Unix-like systems (Linux, macOS, WSL)

set -e

INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="my-context"

echo "Installing my-context..."

# Detect existing installation
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
	echo "Existing installation found."
	CURRENT_VERSION=$("$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null || echo "unknown")
	echo "Current version: $CURRENT_VERSION"

	# Backup old binary
	mv "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME.backup"
	echo "Backed up to $BINARY_NAME.backup"
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

# Copy binary (assume binary provided as $1 or use local build)
if [ -n "$1" ]; then
	cp "$1" "$INSTALL_DIR/$BINARY_NAME"
else
	# Look for platform-appropriate binary in bin/
	if [ "$(uname)" == "Darwin" ]; then
		if [ "$(uname -m)" == "arm64" ]; then
			BINARY_SRC="bin/my-context-darwin-arm64"
		else
			BINARY_SRC="bin/my-context-darwin-amd64"
		fi
	else
		BINARY_SRC="bin/my-context-linux-amd64"
	fi

	if [ -f "$BINARY_SRC" ]; then
		cp "$BINARY_SRC" "$INSTALL_DIR/$BINARY_NAME"
	else
		echo "Error: No binary provided and none found at $BINARY_SRC"
		exit 1
	fi
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Add to PATH if not already present
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
	# Detect shell
	if [ -n "$BASH_VERSION" ]; then
		RC_FILE="$HOME/.bashrc"
	elif [ -n "$ZSH_VERSION" ]; then
		RC_FILE="$HOME/.zshrc"
	else
		RC_FILE="$HOME/.profile"
	fi

	echo "" >> "$RC_FILE"
	echo "# Added by my-context installer" >> "$RC_FILE"
	echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$RC_FILE"
	echo "Added $INSTALL_DIR to PATH in $RC_FILE"
	echo "Run: source $RC_FILE"
fi

# Verify installation
if "$INSTALL_DIR/$BINARY_NAME" --version; then
	echo "✓ Installation complete!"

	# Remove backup if exists
	if [ -f "$INSTALL_DIR/$BINARY_NAME.backup" ]; then
		rm "$INSTALL_DIR/$BINARY_NAME.backup"
	fi
else
	echo "✗ Installation verification failed"
	exit 1
fi

# Note: ~/.my-context/ data is preserved (separate from binary)
```

**Test**: Run on Linux/WSL, verify installation

**After completion**: Mark tasks.md line 194 as `[x]`

---

### T028: Windows cmd.exe Installer
**File**: `scripts/install.bat`

**Status**: ✅ Already created (git status shows `AM scripts/build-all.sh`)

**Verify**: File exists with proper cmd.exe syntax

**After completion**: Mark tasks.md line 201 as `[x]`

---

### T029: Windows PowerShell Installer
**File**: `scripts/install.ps1`

```powershell
# Installation script for Windows PowerShell

$ErrorActionPreference = "Stop"

$InstallDir = "$env:USERPROFILE\bin"
$BinaryName = "my-context.exe"

Write-Host "Installing my-context..."

# Detect existing installation
$ExistingBinary = Join-Path $InstallDir $BinaryName
if (Test-Path $ExistingBinary) {
	Write-Host "Existing installation found."
	$CurrentVersion = & $ExistingBinary --version 2>$null
	if ($CurrentVersion) {
		Write-Host "Current version: $CurrentVersion"
	}

	# Backup
	Move-Item $ExistingBinary "$ExistingBinary.backup" -Force
	Write-Host "Backed up to $BinaryName.backup"
}

# Create install directory
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Copy binary (assume provided as parameter or use local build)
if ($args.Count -gt 0) {
	Copy-Item $args[0] $ExistingBinary
} else {
	$LocalBinary = "bin\my-context-windows-amd64.exe"
	if (Test-Path $LocalBinary) {
		Copy-Item $LocalBinary $ExistingBinary
	} else {
		Write-Error "No binary provided and none found at $LocalBinary"
		exit 1
	}
}

# Add to PATH if not present
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
	$NewPath = "$CurrentPath;$InstallDir"
	[Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
	Write-Host "Added $InstallDir to PATH"
	Write-Host "Restart your terminal for PATH changes to take effect"
}

# Verify installation
$Version = & $ExistingBinary --version
if ($LASTEXITCODE -eq 0) {
	Write-Host "✓ Installation complete!"
	Write-Host $Version

	# Remove backup
	if (Test-Path "$ExistingBinary.backup") {
		Remove-Item "$ExistingBinary.backup"
	}
} else {
	Write-Error "Installation verification failed"
	exit 1
}

Write-Host ""
Write-Host "Note: ~/.my-context/ data is preserved (separate from binary)"
```

**Test**: Run on Windows PowerShell

**After completion**: Mark tasks.md line 208 as `[x]`

---

### T030: Curl One-Liner Installer
**File**: `scripts/curl-install.sh`

```bash
#!/usr/bin/env bash
# One-liner curl installer - auto-detects platform

set -e

REPO="yourusername/my-context-copilot"
VERSION="latest"  # Or specify version

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
	linux*)
		PLATFORM="linux"
		BINARY_NAME="my-context-linux-amd64"
		;;
	darwin*)
		PLATFORM="darwin"
		if [ "$ARCH" == "arm64" ]; then
			BINARY_NAME="my-context-darwin-arm64"
		else
			BINARY_NAME="my-context-darwin-amd64"
		fi
		;;
	mingw*|msys*|cygwin*)
		PLATFORM="windows"
		BINARY_NAME="my-context-windows-amd64.exe"
		;;
	*)
		echo "Error: Unsupported platform: $OS"
		echo "Please download manually from https://github.com/$REPO/releases"
		exit 1
		;;
esac

echo "Detected platform: $PLATFORM ($ARCH)"
echo "Downloading $BINARY_NAME..."

# Download binary and checksum
RELEASE_URL="https://github.com/$REPO/releases/latest/download"
curl -sSLO "$RELEASE_URL/$BINARY_NAME"
curl -sSLO "$RELEASE_URL/$BINARY_NAME.sha256"

# Verify checksum
echo "Verifying checksum..."
if command -v sha256sum &> /dev/null; then
	sha256sum -c "$BINARY_NAME.sha256"
elif command -v shasum &> /dev/null; then
	shasum -a 256 -c "$BINARY_NAME.sha256"
else
	echo "Warning: Cannot verify checksum (sha256sum not found)"
fi

# Make executable
chmod +x "$BINARY_NAME"

# Run installer
if [ "$PLATFORM" == "windows" ]; then
	# For Git Bash on Windows
	./install.bat "$BINARY_NAME"
else
	# Download and run install.sh
	curl -sSL "$RELEASE_URL/install.sh" | bash -s "$BINARY_NAME"
fi

# Cleanup
rm -f "$BINARY_NAME" "$BINARY_NAME.sha256"

echo "✓ Installation complete!"
```

**After completion**: Mark tasks.md line 215 as `[x]`

---

### T031: Build Script Wrapper
**File**: `scripts/build.sh` (modify or create)

```bash
#!/usr/bin/env bash
# Local build wrapper - calls build-all.sh

echo "Building my-context for current platform..."
echo ""
echo "For multi-platform builds, use: ./scripts/build-all.sh"
echo ""

# Build for current platform
go build -o my-context.exe ./cmd/my-context/

echo "✓ Build complete: my-context.exe"
echo ""
echo "Run: ./my-context.exe --version"
```

**After completion**: Mark tasks.md line 222 as `[x]`

---

## 📚 Phase 3.5: Documentation & Polish (0.5 day)

Parallel execution (T032-T034)

### T032 [P]: Update README
**File**: `README.md`

**Add/update these sections**:

1. **Installation** (add at top):
```markdown
## Installation

### Quick Install (Linux, macOS, WSL)
```bash
curl -sSL https://raw.githubusercontent.com/yourusername/my-context-copilot/main/scripts/curl-install.sh | bash
```

### Download Pre-Built Binaries
Download from [Releases](https://github.com/yourusername/my-context-copilot/releases/latest):
- Windows: `my-context-windows-amd64.exe`
- Linux: `my-context-linux-amd64`
- macOS (Intel): `my-context-darwin-amd64`
- macOS (ARM): `my-context-darwin-arm64`

### Manual Installation
**Linux/macOS/WSL**:
```bash
chmod +x my-context-*
./scripts/install.sh my-context-linux-amd64
```

**Windows (PowerShell)**:
```powershell
.\scripts\install.ps1 my-context-windows-amd64.exe
```

**Windows (cmd.exe)**:
```batch
scripts\install.bat my-context-windows-amd64.exe
```
```

2. **Building from Source**:
```markdown
## Building from Source

### Prerequisites
- Go 1.21 or later
- Git

### Build for Current Platform
```bash
go build -o my-context.exe ./cmd/my-context/
```

### Build for All Platforms
```bash
./scripts/build-all.sh
# Creates binaries in bin/ directory
```

### Cross-Platform Build Examples
```bash
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my-context-linux ./cmd/my-context/

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o my-context.exe ./cmd/my-context/

# macOS (ARM)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o my-context-darwin-arm64 ./cmd/my-context/
```
```

3. **New Commands** (in Commands section):
```markdown
### Export Context
```bash
my-context export "Context Name"              # Export to current directory
my-context export "Name" --to path/file.md    # Export to specific path
my-context export --all --to exports/         # Export all contexts
my-context e "Name"                           # Alias
```

### Archive Context
```bash
my-context archive "Completed Work"   # Mark as archived
my-context list --archived            # View archived contexts
my-context a "Name"                   # Alias
```

### Delete Context
```bash
my-context delete "Test Context"      # Delete with confirmation
my-context delete "Name" --force      # Delete without prompt
my-context d "Name"                   # Alias
```

### Enhanced List Filtering
```bash
my-context list                               # Show last 10 (default)
my-context list --all                         # Show all contexts
my-context list --limit 5                     # Show last 5
my-context list --project ps-cli              # Filter by project
my-context list --search "bug fix"            # Search by name
my-context list --active-only                 # Show only active
my-context list --archived                    # Show only archived
my-context list --project ps-cli --limit 3    # Combined filters
```

### Project Organization
```bash
my-context start "Phase 1" --project ps-cli   # Creates "ps-cli: Phase 1"
my-context list --project ps-cli              # Show only ps-cli contexts
```
```

4. **Troubleshooting** (add section):
```markdown
## Troubleshooting

For platform-specific installation issues, see [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md).

Common issues:
- WSL: Use Linux binary, not Windows
- macOS: "App cannot be opened" → Right-click > Open to bypass Gatekeeper
- Windows: PATH not updated → Restart terminal or run `refreshenv`
```

**After completion**: Mark tasks.md line 228 as `[x]`

---

### T033 [P]: Help Text
**File**: `cmd/my-context/main.go`

**Ensure all commands have proper help text** - Already covered in T021-T023 Long descriptions

**Verify**:
```bash
./my-context export --help
./my-context archive --help
./my-context delete --help
./my-context list --help
./my-context start --help
```

**After completion**: Mark tasks.md line 235 as `[x]`

---

### T034 [P]: Issue Templates (Optional)
**Files**: `.github/ISSUE_TEMPLATE/bug_report.md`, `feature_request.md`

**Create standard GitHub templates** (optional for Sprint 2)

**After completion**: Mark tasks.md line 237 as `[x]`

---

## ✅ Phase 3.6: Integration & Validation (1.0 day)

Sequential execution with validation checkpoints.

### T035: Linux Integration Tests
```bash
go test ./tests/integration/... -v
```

**Fix any failures**. Ensure all tests pass on Linux/WSL.

**After completion**: Mark tasks.md line 243 as `[x]`

---

### T036: Windows Integration Tests
```bash
# Run in Windows git-bash
go test ./tests/integration/... -v
```

**Fix platform-specific issues** (path handling, file permissions, etc.)

**After completion**: Mark tasks.md line 245 as `[x]`

---

### T037: Manual Quickstart Validation
**Execute all 9 scenarios** from `specs/002-installation-improvements-and/quickstart.md`:

1. Multi-platform installation (WSL)
2. Project-based workflow
3. Export and share
4. Context lifecycle (archive/delete)
5. List enhancements (large dataset)
6. Bug fixes validation
7. Cross-platform installation (Windows)
8. Backward compatibility (Sprint 1 → Sprint 2)
9. JSON output for scripting

**After completion**: Mark tasks.md line 247 as `[x]`

---

### T038: Performance Benchmarks
**Run performance tests**:

```bash
# Benchmark 1: List with 1000 contexts (create test data first)
time my-context list --all  # Target: < 1 second

# Benchmark 2: Export large context (500 notes)
time my-context export "Large Context"  # Target: < 1 second

# Benchmark 3: Search across 1000 contexts
time my-context list --search "test"  # Target: < 1 second
```

**After completion**: Mark tasks.md line 258 as `[x]`

---

### T039: Binary Size Verification
```bash
ls -lh bin/
```

**Verify**: All 4 platform binaries are <10MB each

**After completion**: Mark tasks.md line 263 as `[x]`

---

### T040: Cross-Platform Smoke Test
**Test on real environments**:
- Ubuntu 22.04 (WSL and/or native)
- Windows 10/11 (cmd.exe and PowerShell)
- macOS 13+ (if available)

**Verify basic commands work identically**:
```bash
my-context start "Test"
my-context note "Test note"
my-context list
my-context export "Test"
my-context stop
```

**After completion**: Mark tasks.md line 265 as `[x]`

---

### T041: Final Constitution Compliance Check
**Review all changes** against 6 principles in `.specify/memory/constitution.md`:

- [x] I. Unix Philosophy
- [x] II. Cross-Platform Compatibility
- [x] III. Stateful Context Management
- [x] IV. Minimal Surface Area (11 commands justified)
- [x] V. Data Portability
- [x] VI. User-Driven Design

**Document any principle tensions** in retrospective

**After completion**: Mark tasks.md line 271 as `[x]`

---

## 📊 Progress Tracking

**Track in tasks.md**: Mark `[ ]` → `[x]` after each task completion

**Phase completion checklist**:
- [ ] Phase 3.1: Setup (4 tasks)
- [ ] Phase 3.2: Tests (8 tasks) - **TDD GATE**
- [ ] Phase 3.3: Implementation (17 tasks)
- [ ] Phase 3.4: Build & Install (5 tasks)
- [ ] Phase 3.5: Documentation (3 tasks)
- [ ] Phase 3.6: Validation (7 tasks)

**Total**: 41 tasks → Sprint 2 complete

---

## 🎯 Recommended Execution Schedule

### Day 1 (2-3 hours)
**Morning**:
- T001-T004: Setup & infrastructure (30 min)
- T005-T012: Write failing tests (90 min)
- Verify all tests FAIL (10 min)

**Afternoon**:
- T013-T015: Model layer (60 min)
- T014 should make project_parser_test.go PASS

---

### Day 2 (4 hours)
**Morning**:
- T016: Core context operations (90 min)
- T017: Enhanced ListContexts (60 min)

**Afternoon**:
- T018-T020: Storage helpers (60 min)
- Run tests - many should be GREEN now

---

### Day 3 (4 hours)
**All day**:
- T021-T026: Commands (parallel, 6 files, ~40 min each)
- Run full test suite - all tests should PASS
- Commit: "feat: implement export, archive, delete commands with enhanced filtering"

---

### Day 4 (3 hours)
**Morning**:
- T027-T031: Build & installation scripts (parallel, ~30 min each)

**Afternoon**:
- T032-T034: Documentation (60 min)
- Commit: "build: add multi-platform build and installation scripts"

---

### Day 5 (3 hours)
**Validation day**:
- T035-T036: Integration tests on Linux & Windows
- T037: Manual quickstart scenarios
- T038-T041: Performance, size, smoke tests, compliance
- Final commit: "test: validate Sprint 2 features across platforms"

---

## 🚨 Troubleshooting During Implementation

### Tests won't compile
- Check import paths (update module name)
- Verify testify package installed: `go get github.com/stretchr/testify`

### Tests pass before implementation
- You've written the implementation first (violates TDD)
- Roll back and write tests first

### Path issues on Windows
- Ensure using `filepath.Join()` not string concatenation
- Use `NormalizePath()` from existing codebase

### Git merge conflicts
- Rebase frequently: `git pull --rebase origin main`
- Resolve conflicts before continuing

---

## 📝 Commit Strategy

**Suggested commits**:

1. After Phase 3.1: `build: add GitHub Actions and build scripts`
2. After Phase 3.2: `test: add integration tests for Sprint 2 features (TDD)`
3. After T013-T015: `feat: add archive support and project filtering to models`
4. After T016-T020: `feat: implement core archive, delete, export logic`
5. After T021-T026: `feat: add export, archive, delete commands with enhanced list`
6. After T027-T031: `build: add installation scripts for all platforms`
7. After T032-T034: `docs: update README and add troubleshooting guide`
8. After T035-T041: `test: validate Sprint 2 features across platforms`

---

## ✅ Definition of Done

Sprint 2 is complete when:

- [x] All 41 tasks marked complete in tasks.md
- [x] All tests passing (`go test ./... -v`)
- [x] All quickstart scenarios validated
- [x] Performance benchmarks met
- [x] Binary sizes <10MB
- [x] Documentation updated
- [x] Backward compatibility verified
- [x] Constitution compliance confirmed
- [x] Code committed to `002-installation-improvements-and` branch

**Final step**: Create PR to main branch with Sprint 2 summary

---

**Note**: This guide assumes incremental development. Each section can be tackled independently. Commit frequently for safe progress tracking.
