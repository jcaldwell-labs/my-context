# Command Contract: touch

**Command**: `my-context touch`  
**Alias**: `my-context t`

## Purpose

Record a timestamp in the current context to indicate activity without adding detailed notes.

## Arguments

**Flags**:
- `--json`: Output result as JSON

## Behavior

1. Read `state.json` to get active_context_name
2. If null, exit with error
3. Append to `~/.my-context/<name>/touch.log`: `<timestamp>\n`

## Output

### Human-Readable (stdout)

```
Touch recorded in context: Bug fix
```

### JSON Output (stdout)

```json
{
  "command": "touch",
  "timestamp": "2025-10-04T14:15:00Z",
  "data": {
    "context_name": "Bug fix",
    "touch_timestamp": "2025-10-04T14:15:00Z"
  }
}
```

## Exit Codes

- `0`: Success
- `1`: No active context

## Examples

```bash
$ my-context touch
Touch recorded in context: Bug fix

$ my-context t
Touch recorded in context: Bug fix
```
