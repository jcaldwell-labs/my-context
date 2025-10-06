# Contract: List Command Enhancements

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
