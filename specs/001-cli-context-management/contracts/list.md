# Command Contract: list

**Command**: `my-context list`  
**Alias**: `my-context l`

## Purpose

Display all contexts (active and stopped) with names, start times, durations, and status.

## Arguments

**Flags**:
- `--json`: Output result as JSON

## Behavior

1. Scan `~/.my-context/` directory for subdirectories
2. For each subdirectory: Read `meta.json`
3. Count lines in notes.log, files.log, touch.log (optional for display)
4. Sort by start_time (most recent first)
5. Format and output

## Output

### Human-Readable (stdout)

```
Contexts (3):

  ● Bug fix (active)
    Started: 2025-10-04 14:00 (2h 15m ago)
    Notes: 3 | Files: 2 | Touches: 5

  ○ Code review (stopped)
    Started: 2025-10-03 10:00 (1d 6h ago)
    Duration: 1h 23m
    Notes: 5 | Files: 4 | Touches: 8

  ○ Deploy production (stopped)
    Started: 2025-10-02 15:00 (2d 1h ago)
    Duration: 45m
    Notes: 2 | Files: 1 | Touches: 3
```

**No Contexts**:
```
No contexts found
Start one with: my-context start <name>
```

### JSON Output (stdout)

```json
{
  "command": "list",
  "timestamp": "2025-10-04T16:15:00Z",
  "data": {
    "contexts": [
      {
        "name": "Bug fix",
        "start_time": "2025-10-04T14:00:00Z",
        "end_time": null,
        "status": "active",
        "duration_seconds": 8100,
        "note_count": 3,
        "file_count": 2,
        "touch_count": 5
      },
      {
        "name": "Code review",
        "start_time": "2025-10-03T10:00:00Z",
        "end_time": "2025-10-03T11:23:00Z",
        "status": "stopped",
        "duration_seconds": 4980,
        "note_count": 5,
        "file_count": 4,
        "touch_count": 8
      }
    ]
  }
}
```

## Exit Codes

- `0`: Success (even if no contexts)

## Examples

```bash
$ my-context list
Contexts (2):
  ● Bug fix (active)
  ○ Code review (stopped)

$ my-context l --json | jq '.data.contexts | length'
2
```
