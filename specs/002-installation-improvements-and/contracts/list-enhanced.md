# Command Contract: list (Enhanced)

**Alias**: `ls`, `l`  
**Purpose**: Display contexts with filtering and pagination options

## Syntax

```bash
my-context list [--all] [--limit N] [--search TERM] [--project NAME] [--archived] [--active-only] [--json]
```

## Flags (New in Sprint 2)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --all | | boolean | false | Show all contexts (no limit) |
| --limit | -n | int | 10 | Number of contexts to show |
| --search | -s | string | "" | Filter by context name (case-insensitive substring) |
| --project | -p | string | "" | Filter by project name |
| --archived | | boolean | false | Show only archived contexts |
| --active-only | | boolean | false | Show only currently active context |
| --json | -j | boolean | false | Output as JSON |

## Behavior

### Default List (Sprint 2)

**Input**: `my-context list`

**Process**:
1. Scan ~/.my-context/ for context directories
2. Filter out archived contexts (is_archived == true)
3. Sort by start_time descending (most recent first)
4. Take first 10 contexts
5. Display with start time, duration, active indicator

**Output**:
```
Contexts (showing 10 of 47):
  * ps-cli: Phase 2                2025-10-09 14:30  [Active]
    ps-cli: Phase 1                2025-10-05 10:15  (3h 45m)
    garden: Planning               2025-10-03 09:00  (1h 20m)
    ...
    
Showing 10 of 47 contexts. Use --all to see all.
```

### Combined Filters

**Input**: `my-context list --project ps-cli --limit 5 --search Phase`

**Process**:
1. Filter by project (case-insensitive): "ps-cli"
2. Filter by search term (case-insensitive): "Phase"
3. Apply both filters (AND logic)
4. Sort and limit to 5

**Output**:
```
Contexts (showing 5 of 12 matching):
  * ps-cli: Phase 3                2025-10-09 16:00  [Active]
    ps-cli: Phase 2                2025-10-08 10:00  (2h 30m)
    ps-cli: Phase 1                2025-10-05 10:15  (3h 45m)
```

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| No contexts exist | "No contexts found." (exit 0) |
| No contexts match filters | "No contexts found for project 'name'" or "No contexts found matching 'term'" (exit 0) |
| --archived + --active-only | Error: "Cannot combine --archived and --active-only" (exit 1) |
| --limit with --all | --all takes precedence (show all) |
| No active context with --active-only | "No active context" (exit 0) |

## Filter Logic

**Project Extraction**:
- "ps-cli: Phase 1" → project = "ps-cli"
- "Standalone" → project = "Standalone"
- Case-insensitive matching

**Search Matching**:
- Substring match in context name
- Case-insensitive
- Example: "Phase" matches "ps-cli: Phase 1", "Planning Phase", "phase-2"

**Multiple Filters**:
- Applied with AND logic (all conditions must match)
- Order doesn't matter

## Examples

```bash
# Default (last 10, non-archived)
my-context list

# Show all contexts
my-context list --all

# Custom limit
my-context list --limit 20

# Filter by project
my-context list --project ps-cli

# Search by name
my-context list --search "bug fix"

# Show archived contexts
my-context list --archived

# Show only active
my-context list --active-only

# Combined filters
my-context list --project ps-cli --search Phase --limit 5

# JSON output
my-context list --json
```
