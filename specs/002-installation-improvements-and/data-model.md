# Data Model: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Date**: 2025-10-05  
**Status**: Complete

---

## Entity Changes

### 1. Context (Modified)

**Purpose**: Represents a work session with associated notes, files, and activity tracking.

**Storage**: `~/.my-context/{context_name}/meta.json`

**Schema**:
```json
{
  "name": "string (required)",
  "start_time": "RFC3339 timestamp (required)",
  "end_time": "RFC3339 timestamp (optional, null if active)",
  "status": "string enum: active|stopped (required)",
  "subdirectory_path": "string (required)",
  "is_archived": "boolean (optional, defaults to false)"
}
```

**New Field**:
- `is_archived`: Boolean flag indicating context is archived (hidden from default list view)
  - Default: `false`
  - Optional for backward compatibility with Sprint 1 data
  - Does not affect data accessibility (files remain in place)

**Validation Rules**:
- `name`: Non-empty string, max 200 characters
- `start_time`: Must be valid RFC3339 timestamp
- `end_time`: Must be after `start_time` if present
- `status`: Must be "active" or "stopped"
- `is_archived`: Cannot be true if status is "active"

**State Transitions**:
```
Created (active, is_archived=false)
  ↓ stop command
Stopped (stopped, is_archived=false)
  ↓ archive command
Archived (stopped, is_archived=true)
  ↓ delete command
Deleted (removed from filesystem)
```

---

## 2. Project Metadata (Derived)

**Purpose**: Logical grouping of contexts by project name, extracted from context naming convention.

**Storage**: Not persisted (computed on-demand from context names)

**Schema**:
```go
type ProjectMetadata struct {
    Name            string   // Extracted project name
    ContextNames    []string // Contexts belonging to this project
    ActiveCount     int      // Number of active contexts
    ArchivedCount   int      // Number of archived contexts
}
```

**Extraction Logic**:
```go
// ExtractProjectName parses "project: phase - description" format
// Returns text before first colon, or full name if no colon
func ExtractProjectName(contextName string) string {
    parts := strings.SplitN(contextName, ":", 2)
    return strings.TrimSpace(parts[0])
}
```

**Examples**:
- `"ps-cli: Phase 1 - Foundation"` → `"ps-cli"`
- `"garden: Planning"` → `"garden"`
- `"Standalone Context"` → `"Standalone Context"`
- `"project: sub: detail"` → `"project"` (only first colon)

**Filtering Algorithm**:
```go
func FilterByProject(contexts []Context, projectName string) []Context {
    filtered := []Context{}
    for _, ctx := range contexts {
        if strings.EqualFold(ExtractProjectName(ctx.Name), projectName) {
            filtered = append(filtered, ctx)
        }
    }
    return filtered
}
```

---

## 3. Export Document (Output Format)

**Purpose**: Human-readable markdown representation of context data for sharing and archival.

**Storage**: User-specified path (default: `{context_name}.md` in current directory)

**Format**:
```markdown
# Context: {context_name}

**Started**: {start_time in local timezone}
**Ended**: {end_time in local timezone or "Active"}
**Duration**: {human_readable duration or "Active"}
**Status**: {active | stopped | archived}

## Notes ({count})

{if count > 0:}
- `{HH:MM}` {note_text}
- `{HH:MM}` {note_text}
...
{else:}
(No notes)

## Associated Files ({count})

{if count > 0:}
- {file_path}
  Added: {timestamp}
...
{else:}
(No files)

## Activity

- {touch_count} touch events
{if touch_count > 0:}
- Last activity: {most_recent_touch_timestamp}
{else:}
- No recorded activity

---
*Exported: {export_timestamp} by my-context v{version}*
```

**Data Sources**:
- Context metadata: `meta.json`
- Notes: `notes.log` (pipe-delimited: `timestamp|text`)
- Files: `files.log` (pipe-delimited: `timestamp|path`)
- Touch events: `touch.log` (one timestamp per line)

**Processing Rules**:
1. Timestamps converted to local timezone for readability
2. Duration calculated from start/end times (format: "2h 15m" or "3d 4h")
3. Notes displayed in chronological order
4. Files displayed with relative paths if under user home directory
5. Touch events counted, most recent timestamp extracted

---

## 4. Binary Artifact (Build Output)

**Purpose**: Platform-specific executable binaries for distribution.

**Storage**: GitHub Releases (artifacts) and local `bin/` directory during builds

**Naming Convention**:
```
my-context-{os}-{arch}[.exe]

Examples:
- my-context-linux-amd64
- my-context-windows-amd64.exe
- my-context-darwin-amd64
- my-context-darwin-arm64
```

**Metadata**:
```yaml
Platform: linux | windows | darwin
Architecture: amd64 | arm64
Version: {semantic version, e.g., 1.1.0}
BuildTime: {RFC3339 timestamp}
GitCommit: {short SHA, e.g., a1b2c3d}
GoVersion: {e.g., go1.21.5}
Checksum: {SHA256 hash}
```

**Embedded Version Info**:
```go
// Set at build time with ldflags
var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)
```

**Build Command**:
```bash
go build -ldflags "\
  -X main.Version=1.1.0 \
  -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X main.GitCommit=$(git rev-parse --short HEAD)" \
  -o my-context-{os}-{arch} ./cmd/my-context/
```

---

## 5. Installation Script Metadata

**Purpose**: Track installation state for upgrade detection.

**Storage**: `~/.my-context/.install-info` (optional metadata file)

**Schema**:
```json
{
  "version": "1.1.0",
  "installed_at": "2025-10-05T19:30:00Z",
  "binary_path": "/home/user/.local/bin/my-context",
  "install_method": "curl-install | install.sh | manual"
}
```

**Usage**:
- Installation scripts check for this file to detect upgrades
- Provides rollback information if upgrade fails
- Enables telemetry/analytics (future consideration)

**Optional**: Sprint 2 may create this file but not require it. Absence indicates fresh install.

---

## Relationships

```
Context 1:1 → meta.json (storage)
Context 1:N → Note (via notes.log)
Context 1:N → FileAssociation (via files.log)
Context 1:N → TouchEvent (via touch.log)
Context 1:1 → ProjectMetadata (derived, not stored)
Context 1:1 → ExportDocument (generated on demand)

ProjectMetadata 1:N → Context (logical grouping)

BinaryArtifact 1:N → Platform (one binary per platform)
BinaryArtifact N:1 → Version (all platforms share version)
```

---

## Data Migration Strategy

### Sprint 1 → Sprint 2

**Goal**: Zero downtime, zero data loss, zero user intervention required.

**Changes**:
1. Add optional `is_archived` field to `meta.json`
2. No changes to `notes.log`, `files.log`, `touch.log` formats
3. No changes to `state.json` or `transitions.log` formats

**Compatibility Strategy**:
- Use Go's `omitempty` JSON tag for `is_archived` field
- On read: Missing field defaults to `false` (Go zero value)
- On write: Field included if `true`, omitted if `false` (space optimization)

**Validation**:
```go
// Test: Load Sprint 1 meta.json without is_archived field
input := `{"name":"Test","start_time":"2025-10-04T12:00:00Z","status":"active","subdirectory_path":"..."}`
var ctx Context
json.Unmarshal([]byte(input), &ctx)
assert.False(ctx.IsArchived) // Defaults to false

// Test: Write Sprint 2 meta.json with is_archived=false
ctx.IsArchived = false
output, _ := json.MarshalIndent(ctx, "", "  ")
assert.NotContains(string(output), "is_archived") // Omitted

// Test: Write Sprint 2 meta.json with is_archived=true
ctx.IsArchived = true
output, _ := json.MarshalIndent(ctx, "", "  ")
assert.Contains(string(output), `"is_archived": true`) // Included
```

**Rollback Plan**:
If user downgrades to Sprint 1 binary:
- Sprint 1 code ignores unknown JSON fields (standard Go behavior)
- Archived contexts appear as normal stopped contexts (is_archived ignored)
- No data corruption, just loss of archive functionality

---

## Performance Considerations

### List Command with 1000+ Contexts

**Challenge**: Loading and sorting 1000 meta.json files may be slow.

**Optimization Strategy**:
1. **Lazy loading**: Only load meta.json for contexts in result set (after filtering)
2. **Caching**: Store modified times, only reload changed files (future Sprint)
3. **Pagination**: Default limit of 10 reduces I/O by 99% for large datasets

**Benchmark Target**: <1 second for `my-context list` with 1000 contexts on HDD.

**Implementation**:
```go
// Phase 1: List all context directories (fast, no I/O)
dirs := listContextDirectories() // O(1000) directory reads

// Phase 2: Load only meta.json for stat info (name, start time, status)
contexts := loadLightweightMetadata(dirs) // O(1000) small file reads

// Phase 3: Sort by start time (in-memory)
sort.Slice(contexts, ...) // O(n log n)

// Phase 4: Apply filters and limit
filtered := applyFilters(contexts, filters) // O(n)
limited := filtered[:limit] // O(1)

// Phase 5: Load full details only for displayed contexts (if needed)
// For list view, lightweight metadata is sufficient
```

### Export Command Performance

**Challenge**: Reading and formatting multiple log files for large contexts.

**Optimization**:
- Stream log files line-by-line (don't load entire file into memory)
- Use buffered I/O for markdown output
- Format timestamps lazily (only when writing to output)

**Benchmark Target**: <1 second for context with 500 notes + 100 files.

---

## Summary

**Modified Entities**: 1 (Context - added is_archived field)  
**New Derived Entities**: 2 (ProjectMetadata, ExportDocument)  
**New Artifact Entities**: 2 (BinaryArtifact, InstallationMetadata)

**Total Fields Added**: 1 (is_archived)  
**Total Schema Breakages**: 0 (backward compatible)

**Data Migration Required**: No (graceful degradation)  
**Performance Impact**: Minimal (lazy loading, pagination)

---

*Data model complete. Ready for contract generation.*
