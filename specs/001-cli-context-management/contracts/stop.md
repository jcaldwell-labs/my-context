# Command Contract: stop

**Command**: `my-context stop`  
**Alias**: `my-context p`

## Purpose

Explicitly stop the currently active context without starting a new one.

## Arguments

**Flags**:
- `--json`: Output result as JSON

## Behavior

1. Read `state.json` to get active_context_name
2. If null, display "No active context" (exit 0, not an error)
3. Read `~/.my-context/<name>/meta.json`
4. Update: end_time=now, status="stopped"
5. Write `meta.json` atomically
6. Update `state.json`: active_context=null
7. Append to `transitions.log`: `<timestamp>|<name>|NULL|stop`

## Output

### Human-Readable (stdout)

**Success**:
```
Stopped context: Bug fix (duration: 2h 15m)
```

**No Active Context**:
```
No active context
```

### JSON Output (stdout)

```json
{
  "command": "stop",
  "timestamp": "2025-10-04T16:15:00Z",
  "data": {
    "context_name": "Bug fix",
    "start_time": "2025-10-04T14:00:00Z",
    "end_time": "2025-10-04T16:15:00Z",
    "duration_seconds": 8100
  }
}
```

## Exit Codes

- `0`: Success (even if no context)
- `2`: System error (I/O failure)

## Examples

```bash
$ my-context stop
Stopped context: Bug fix (duration: 2h 15m)

$ my-context stop
No active context
```
