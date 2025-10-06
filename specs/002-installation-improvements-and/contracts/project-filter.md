# Contract: Project Filter (Start Command Enhancement)

**Command**: `start` (alias: `s`) - Enhanced with --project flag  
**Purpose**: Create contexts with project prefix following "project: phase" convention

---

## Signature

```bash
my-context start <name> [--project <project-name>] [--json]
my-context s <name> [--project <project-name>] [--json]
```

---

## Arguments

**Positional**:
- `<name>`: Context name or phase name (required)

**Flags**:
- `--project <project-name>`: Prepend project name to context (optional)
- `--json`: Output JSON format (optional)

---

## Behavior

### Start Without Project Flag (Existing Behavior)

**Input**: `my-context start "Working on feature"`

**Actions**:
1. Stop any active context
2. Create context named exactly "Working on feature"
3. Normalize name for directory (spaces → underscores)
4. Create context directory and meta.json
5. Update state.json
6. Log transition

**Output**:
```
Started context: Working on feature
```

**Exit Code**: 0

### Start With Project Flag

**Input**: `my-context start "Phase 1 - Foundation" --project ps-cli`

**Actions**:
1. Stop any active context
2. Combine project + name: "ps-cli: Phase 1 - Foundation"
3. Check for duplicate name (apply _2 suffix if needed)
4. Create context with combined name
5. Update state.json
6. Log transition

**Output**:
```
Started context: ps-cli: Phase 1 - Foundation
```

**Exit Code**: 0

### Start With Project Flag (Short Phase Name)

**Input**: `my-context start "Phase 2" --project ps-cli`

**Actions**:
1. Combine: "ps-cli: Phase 2"
2. Create context

**Output**:
```
Started context: ps-cli: Phase 2
```

**Exit Code**: 0
# Contract: List Command Enhancements
---

## Name Formatting Rules

### Without --project flag:
- Use name as-is: `"Working on feature"` → `"Working on feature"`

### With --project flag:
- Format: `"{project}: {name}"`
- Example: `"ps-cli"` + `"Phase 1"` → `"ps-cli: Phase 1"`
- Trim whitespace from both parts
- Single space after colon

### Directory Naming:
- Replace spaces with underscores: `"ps-cli: Phase 1"` → `"ps-cli_Phase_1"`
- Preserve colons, hyphens, alphanumerics
- Remove or replace special characters: `/ \ * ? " < > |` → `_`

---

## Duplicate Name Handling

### Scenario: Context "ps-cli: Phase 1" Already Exists

**Input**: `my-context start "Phase 1" --project ps-cli`

**Actions**:
1. Combine to "ps-cli: Phase 1"
2. Check existence
3. Append "_2" suffix: "ps-cli: Phase 1_2"
4. Create with suffixed name

**Output**:
```
Started context: ps-cli: Phase 1_2
```

**Note**: Same duplicate handling as Sprint 1 (existing behavior)

---

## Integration with List --project

The --project flag on `start` ensures contexts are created in a format that works with `list --project` filtering.

**Workflow**:
```bash
# Create contexts with project
my-context start "Phase 1" --project ps-cli
my-context start "Phase 2" --project ps-cli
my-context start "Research" --project garden

# Filter by project
my-context list --project ps-cli
# Shows: ps-cli: Phase 1, ps-cli: Phase 2

my-context list --project garden
# Shows: garden: Research
```

---

## Error Handling

### Empty Name

**Input**: `my-context start "" --project ps-cli`

**Output**:
```
Error: Context name cannot be empty
Usage: my-context start <name> [--project <project-name>]
```

**Exit Code**: 1

### Empty Project Name

**Input**: `my-context start "Phase 1" --project ""`

**Output**:
```
Error: Project name cannot be empty
Usage: my-context start <name> [--project <project-name>]
```

**Exit Code**: 1

---

## JSON Output Format

**Input**: `my-context start "Phase 1" --project ps-cli --json`

**Output**:
```json
{
  "command": "start",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "context_name": "ps-cli: Phase 1",
    "project_name": "ps-cli",
    "phase_name": "Phase 1",
    "previous_context": "Other Context",
    "transition_type": "switch"
  }
}
```

---

## Backward Compatibility

### Sprint 1 Users

Users who don't use --project flag experience no change:
- `my-context start "MyContext"` works exactly as before
- Existing contexts without "project: " format work normally
- List command treats non-colon names as standalone projects

### Migration Path

Users can gradually adopt project convention:
1. Start: Use names like "project: phase" manually (Sprint 1 compatible)
2. Upgrade: Use --project flag for convenience (Sprint 2)
3. Filter: Use list --project to organize work (Sprint 2)

---

## Integration Points

**Calls**:
- `internal/core/context.go`: CreateContext() - Modified to accept project parameter
- `internal/core/state.go`: StopActiveContext(), SetActiveContext()
- `internal/core/storage.go`: WriteJSON(), AppendLog()

**Modified**:
- `internal/commands/start.go`: Add --project flag parsing

**Compatible With**:
- `list --project` filtering
- Archive/delete commands (work with full name)
- Export command (preserves full name)

---

## Validation Rules

1. Name must be non-empty string after trimming
2. Project must be non-empty string after trimming (if provided)
3. Combined name must be ≤200 characters
4. Combined name must not conflict with existing context (or get _2 suffix)
5. Project and name are trimmed of leading/trailing whitespace

---

## Testing Checklist

- [ ] Start without --project flag (existing behavior)
- [ ] Start with --project flag (combined name)
- [ ] Start with --project and short phase name
- [ ] Start with --project and duplicate name (suffix applied)
- [ ] Start with empty name (error)
- [ ] Start with empty project (error)
- [ ] Start with very long combined name (validation)
- [ ] Verify directory naming (special chars handled)
- [ ] Verify list --project filters correctly on created contexts
- [ ] Verify JSON output format
- [ ] Verify backward compatibility (no --project still works)
- [ ] Verify whitespace trimming on project and name

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05

**Command**: `list` (alias: `l`)  
**Purpose**: Display contexts with filtering, pagination, and search capabilities

---

## Signature

```bash
my-context list [--project <name>] [--limit <n>] [--search <term>] [--all] [--archived] [--active-only] [--json]
my-context l [--project <name>] [--limit <n>] [--search <term>] [--all] [--archived] [--active-only] [--json]
```

---

## Arguments

**Flags** (all optional, can be combined):
- `--project <name>`: Filter by project name (case-insensitive)
- `--limit <n>`: Show only N most recent contexts (default: 10)
- `--search <term>`: Filter by context name containing term (case-insensitive)
- `--all`: Show all contexts (overrides --limit default)
- `--archived`: Show only archived contexts (mutually exclusive with --active-only)
- `--active-only`: Show only active context (mutually exclusive with --archived)
- `--json`: Output JSON format

---

## Behavior

### Default Behavior (No Flags)

**Input**: `my-context list`

**Actions**:
1. Load all non-archived contexts
2. Sort by start_time (newest first)
3. Limit to 10 most recent
4. Display with status indicators
5. Show truncation message if more exist

**Output**:
```
● CurrentWork (active)
    Started: 2025-10-05 19:30 (25m ago)

  ○ ps-cli: Phase 2 (stopped)
    Started: 2025-10-05 15:00 (4h 30m ago)
    Duration: 3h 15m

  ○ ps-cli: Phase 1 (stopped)
    Started: 2025-10-04 10:00 (1d 9h ago)
    Duration: 8h

  ... (7 more contexts)

Showing 10 of 25 contexts. Use --all to see all.
```

**Exit Code**: 0

### Project Filter

**Input**: `my-context list --project ps-cli`

**Actions**:
1. Load all non-archived contexts
2. Extract project name from each context (text before first colon)
3. Case-insensitive match against "ps-cli"
4. Apply default limit (10)
5. Display filtered results

**Output**:
```
  ○ ps-cli: Phase 2 (stopped)
    Started: 2025-10-05 15:00 (4h 30m ago)
    Duration: 3h 15m

  ○ ps-cli: Phase 1 (stopped)
    Started: 2025-10-04 10:00 (1d 9h ago)
    Duration: 8h

Found 2 contexts for project "ps-cli"
```

**Exit Code**: 0

### Search Filter

**Input**: `my-context list --search "bug fix"`

**Actions**:
1. Load all non-archived contexts
2. Case-insensitive substring match on context names
3. Apply default limit
4. Display results

**Output**:
```
  ○ Bug fix #123 (stopped)
    Started: 2025-10-03 14:00 (2d 5h ago)
    Duration: 2h 30m

Found 1 context matching "bug fix"
```

**Exit Code**: 0 (even if no matches)

### Combined Filters

**Input**: `my-context list --project ps-cli --limit 5 --search "Phase"`

**Actions**:
1. Load all non-archived contexts
2. Filter by project = "ps-cli"
3. Filter by name contains "Phase"
4. Limit to 5 results
5. Display

**Output**:
```
  ○ ps-cli: Phase 2 (stopped)
    Started: 2025-10-05 15:00 (4h 30m ago)
    Duration: 3h 15m

  ○ ps-cli: Phase 1 (stopped)
    Started: 2025-10-04 10:00 (1d 9h ago)
    Duration: 8h

Found 2 contexts matching filters
```

**Exit Code**: 0

### Archived Contexts

**Input**: `my-context list --archived`

**Actions**:
1. Load all archived contexts (is_archived = true)
2. Sort by start_time (newest first)
3. Apply limit (default 10)
4. Display

**Output**:
```
  ○ garden: Planning (archived)
    Started: 2025-09-15 09:00 (20d ago)
    Duration: 4d 6h

  ○ old-project: Legacy Code (archived)
    Started: 2025-08-01 10:00 (65d ago)
    Duration: 12d 3h

Found 2 archived contexts
```

**Exit Code**: 0

### Active Only

**Input**: `my-context list --active-only`

**Actions**:
1. Load only active context (status = "active")
2. Display

**Output**:
```
● CurrentWork (active)
    Started: 2025-10-05 19:30 (30m ago)

1 active context
```

**Or if no active context**:
```
No active context
Run 'my-context start <name>' to create one.
```

**Exit Code**: 0

---

## Error Handling

### Invalid Limit Value

**Input**: `my-context list --limit -5`

**Output**:
```
Error: --limit must be a positive number
Usage: my-context list [--limit <n>]
```

**Exit Code**: 1

### Conflicting Flags

**Input**: `my-context list --archived --active-only`

**Output**:
```
Error: Cannot use --archived and --active-only together
```

**Exit Code**: 1

### No Matches Found

**Input**: `my-context list --project "nonexistent"`

**Output**:
```
No contexts found for project "nonexistent"
```

**Exit Code**: 0 (not an error, just no results)

---

## Display Format

### Status Indicators
- `●` = Active context (bold/colored in terminal)
- `○` = Stopped context

### Time Display
- **Relative times**: "25m ago", "4h 30m ago", "2d 5h ago"
- **Duration**: "3h 15m", "8h", "4d 6h"
- **Format rules**:
  - < 1 hour: "Xm ago"
  - < 1 day: "Xh Ym ago"
  - ≥ 1 day: "Xd Yh ago"

### Sorting
- **Primary**: Start time (newest first)
- **Secondary**: N/A (unique timestamps)

---

## JSON Output Format

**Input**: `my-context list --project ps-cli --json`

**Output**:
```json
{
  "command": "list",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "contexts": [
      {
        "name": "ps-cli: Phase 2",
        "start_time": "2025-10-05T15:00:00Z",
        "end_time": "2025-10-05T18:15:00Z",
        "status": "stopped",
        "is_archived": false,
        "duration_seconds": 11700
      },
      {
        "name": "ps-cli: Phase 1",
        "start_time": "2025-10-04T10:00:00Z",
        "end_time": "2025-10-04T18:00:00Z",
        "status": "stopped",
        "is_archived": false,
        "duration_seconds": 28800
      }
    ],
    "filters": {
      "project": "ps-cli",
      "limit": 10,
      "archived": false,
      "active_only": false
    },
    "total_count": 2,
    "showing_count": 2
  }
}
```

---

## Integration Points

**Reads from**:
- `internal/core/context.go`: ListContexts()
- `internal/core/storage.go`: LoadAllContextMeta()
- `internal/models/context.go`: Context struct with IsArchived field

**Depends on**:
- Project extraction logic (ExtractProjectName)
- Time formatting utilities
- Status display logic

---

## Validation Rules

1. `--limit` must be positive integer or "all"
2. `--project` accepts any non-empty string
3. `--search` accepts any non-empty string
4. `--archived` and `--active-only` are mutually exclusive
5. Multiple filters combine with AND logic (all must match)

---

## Testing Checklist

- [ ] Default list (10 most recent, non-archived)
- [ ] List --all (no limit)
- [ ] List --limit 5 (custom limit)
- [ ] List --project <name> (filter by project)
- [ ] List --search <term> (substring search)
- [ ] List --archived (only archived contexts)
- [ ] List --active-only (only active context)
- [ ] List with no active context
- [ ] List with combined filters (project + search + limit)
- [ ] List with conflicting flags (archived + active-only)
- [ ] List with invalid limit (negative, zero, non-numeric)
- [ ] List with 0 results (empty message)
- [ ] List with 1000+ contexts (performance)
- [ ] List with truncation message (more than limit)
- [ ] Verify JSON output format
- [ ] Verify sorting (newest first)
- [ ] Verify relative time calculations
- [ ] Verify duration calculations

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05
