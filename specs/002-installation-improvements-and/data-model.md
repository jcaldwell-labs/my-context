# Data Model: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Date**: 2025-10-09  
**Phase**: Phase 1 (Design & Contracts)

## Overview

This document defines the data structures and their relationships for Sprint 2 features. Most structures extend existing Sprint 1 models with backward-compatible additions.

---

## Entity 1: Context (Modified from Sprint 1)

**Purpose**: Represents a work context with enhanced lifecycle management

**Storage**: `~/.my-context/<context-name>/meta.json`

**Structure**:
```json
{
  "name": "ps-cli: Phase 1 - Foundation",
  "start_time": "2025-10-05T14:30:00Z",
  "end_time": null,
  "is_archived": false
}
```

**Fields**:
| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| name | string | ✅ | Context name, may include project prefix | Non-empty, unique |
| start_time | RFC3339 timestamp | ✅ | When context was started | Valid ISO8601 format |
| end_time | RFC3339 timestamp or null | ❌ | When context was stopped, null if active | Must be after start_time if set |
| is_archived | boolean | ❌ | Whether context is archived (hidden from default list) | Defaults to false if omitted |

**Relationships**:
- 1 Context → Many Notes (in notes.log file)
- 1 Context → Many File Associations (in files.log)
- 1 Context → Many Touch Events (in touches.log)

**State Transitions**:
```
Created (end_time=null, is_archived=false)
    ↓ stop
Stopped (end_time set, is_archived=false)
    ↓ archive
Archived (end_time set, is_archived=true)
    ↓ delete
Deleted (directory removed)
```

**Constraints**:
- Cannot archive active context (end_time must be set)
- Cannot delete active context (must stop first)
- Archive status persists across application restarts

**Backward Compatibility**:
- Sprint 1 contexts without `is_archived` field default to false
- Missing field handled via Go's `json:"is_archived,omitempty"` tag

---

## Entity 2: ProjectMetadata (New)

**Purpose**: Extracted project information for filtering and grouping

**Storage**: Ephemeral (derived at runtime, not persisted)

**Structure**:
```go
type ProjectMetadata struct {
    ProjectName   string   // Extracted from context name
    ContextNames  []string // All contexts matching this project
    ContextCount  int      // Number of contexts in project
}
```

**Extraction Logic**:
```go
// "ps-cli: Phase 1" → project="ps-cli"
// "garden: Planning" → project="garden"  
// "Standalone Context" → project="Standalone Context"
```

**Usage**:
- Filter list results: `my-context list --project ps-cli`
- Start with project prefix: `my-context start "Phase 3" --project ps-cli`

**Validation**:
- Project name is case-insensitive (`strings.EqualFold`)
- Empty project names not allowed (minimum 1 character)
- Project names trimmed of leading/trailing whitespace

---

## Entity 3: ExportDocument (New)

**Purpose**: Markdown representation of a context for sharing and archival

**Storage**: User-specified file path or default `./<context-name>.md`

**Structure**:
```markdown
# Context: [Context Name]

**Started**: [Human-readable timestamp with timezone]
**Ended**: [Timestamp or "Currently Active"]
**Duration**: [Human-readable duration, e.g., "2h 15m"]

**Exported**: [Export timestamp]

---

## Notes

- [Local time] Note content
- [Local time] Note content

Total: N notes

---

## Files

- [Local time] /absolute/path/to/file

Total: N files

---

## Activity

- [Local time] Touch event

Total: N touches

---

*Exported from my-context v2.0.0*
```

**Generation Logic**:
1. Read meta.json for context metadata
2. Parse notes.log, files.log, touches.log
3. Convert UTC timestamps to local timezone
4. Format duration as human-readable string
5. Generate markdown sections
6. Write atomically to avoid partial files

**Formatting Rules**:
- Timestamps in local timezone (user-friendly)
- Durations: "Xh Ym" format, "Active" if no end_time
- Notes/Files/Activity sorted chronologically
- Empty sections show "(none)" instead of omitting

**Output Options**:
- `--to <path>`: Custom output path (creates parent directories)
- Default: `./<sanitized-context-name>.md`
- `--force`: Overwrite existing files without confirmation

---

## Entity 4: BinaryArtifact (New)

**Purpose**: Describes a platform-specific pre-built binary

**Storage**: GitHub Releases, referenced by installation scripts

**Structure**:
```json
{
  "platform": "linux-amd64",
  "filename": "my-context-linux-amd64",
  "version": "v2.0.0",
  "download_url": "https://github.com/.../releases/download/v2.0.0/my-context-linux-amd64",
  "checksum_url": "https://github.com/.../releases/download/v2.0.0/my-context-linux-amd64.sha256",
  "size_bytes": 8388608
}
```

**Supported Platforms**:
| Platform | GOOS | GOARCH | Filename |
|----------|------|--------|----------|
| Windows 64-bit | windows | amd64 | my-context-windows-amd64.exe |
| Linux 64-bit | linux | amd64 | my-context-linux-amd64 |
| macOS Intel | darwin | amd64 | my-context-darwin-amd64 |
| macOS ARM (M1/M2) | darwin | arm64 | my-context-darwin-arm64 |

**Checksum Format**:
```
SHA256(my-context-linux-amd64)= 4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d...
```

**Validation**:
- Installation scripts must verify SHA256 before executing
- Mismatch results in installation failure with clear error message

---

## Entity 5: InstallationMetadata (New)

**Purpose**: Tracks where binary is installed for upgrade detection

**Storage**: `~/.my-context/.install-metadata` (hidden file)

**Structure**:
```json
{
  "version": "v2.0.0",
  "binary_path": "/home/user/.local/bin/my-context",
  "installed_at": "2025-10-09T19:00:00Z",
  "install_method": "curl-install.sh"
}
```

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| version | string | Semantic version of installed binary |
| binary_path | string | Absolute path where binary is installed |
| installed_at | RFC3339 | When installation occurred |
| install_method | string | Script used (install.sh, install.bat, install.ps1, curl-install.sh, manual) |

**Usage**:
- Installation scripts check this file to detect upgrades
- Future version: `my-context upgrade` command can use this metadata
- Preserved during upgrades

**Optional**: This metadata file is optional for Sprint 2 (stretch goal). Core installation works without it.

---

## Data Relationships

```
Context
  ├─ Notes (1:N)
  ├─ Files (1:N)
  └─ Touch Events (1:N)

ProjectMetadata (derived)
  └─ Contexts (1:N)

ExportDocument
  └─ Context (1:1) - represents single context export

BinaryArtifact
  └─ Platform (1:1) - one binary per platform

InstallationMetadata
  └─ Binary (1:1) - tracks one installation
```

---

## File System Structure

```
~/.my-context/
├── context-name-1/
│   ├── meta.json              # Context entity (is_archived added)
│   ├── notes.log              # Newline-delimited notes
│   ├── files.log              # Newline-delimited file paths
│   └── touches.log            # Newline-delimited touch timestamps
├── context-name-2/
│   └── ...
├── transitions.log            # Context switch history (preserved on delete)
└── .install-metadata          # Optional: Installation tracking

Exported Files:
./context-name-1.md            # ExportDocument (user-specified location)
```

---

## Migration Strategy

**Sprint 1 → Sprint 2**:
- ✅ No migration required
- ✅ Existing meta.json files work without changes (is_archived defaults to false)
- ✅ Installation preserves entire ~/.my-context/ directory
- ✅ Backward compatibility tested in T012

---

## Performance Considerations

| Operation | Data Access Pattern | Expected Performance |
|-----------|---------------------|---------------------|
| List 1000 contexts | Read 1000 meta.json files | <1s (sequential reads, ~1ms per file) |
| Export 500 notes | Read 1 meta.json + 3 log files | <1s (mostly string formatting) |
| Archive context | Update 1 meta.json file | <50ms (single write) |
| Filter by project | Read all meta.json, filter in memory | <1s (1000 contexts × 1ms) |
| Search by name | Read all meta.json, substring match | <1s (in-memory string operations) |

**Optimization**:
- Context list caching not needed (1s is acceptable for CLI tool)
- Log files read once per export (streaming not required for <500 entries)
- JSON parsing is fast enough for small documents (meta.json ~200 bytes)

---

## Validation Rules

**Context Name**:
- ✅ Non-empty
- ✅ Unique within ~/.my-context/
- ✅ Valid directory name (no `/`, `\`, null bytes)
- ⚠️  Case-sensitive on Unix, case-insensitive on Windows (use exact match to avoid confusion)

**Timestamps**:
- ✅ RFC3339 format (ISO8601 with timezone)
- ✅ end_time must be after start_time if both set
- ✅ Parse failures result in clear error messages

**Project Names**:
- ✅ Extracted before first colon
- ✅ Trimmed of whitespace
- ✅ Case-insensitive matching
- ✅ Contexts without colons use full name as project

---

**Data Model Complete**: All entities defined with fields, relationships, validation rules, and performance targets. Ready for contract generation.
