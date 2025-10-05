# Command Contract: history

**Command**: `my-context history`  
**Alias**: `my-context h`

## Purpose

Display the chronological transition log showing all context switches with timestamps and transition types.

## Arguments

**Flags**:
- `--json`: Output result as JSON
- `--limit <n>`: Show only last N transitions (default: all)

## Behavior

1. Read `~/.my-context/transitions.log`
2. Parse lines: `timestamp|prev|new|type`
3. If --limit specified, take last N lines
4. Format and output

## Output

### Human-Readable (stdout)

```
Context History:
# Command Contract: note
  2025-10-04 14:00 │ START     │ Bug fix
  2025-10-04 16:30 │ SWITCH    │ Bug fix → Code review
  2025-10-04 17:00 │ STOP      │ Code review
  2025-10-04 17:15 │ START     │ Deploy production
  2025-10-04 18:00 │ STOP      │ Deploy production
```

### JSON Output (stdout)

```json
{
  "command": "history",
  "timestamp": "2025-10-04T18:05:00Z",
  "data": {
    "transitions": [
      {
        "timestamp": "2025-10-04T14:00:00Z",
        "previous_context": null,
        "new_context": "Bug fix",
        "transition_type": "start"
      },
      {
        "timestamp": "2025-10-04T16:30:00Z",
        "previous_context": "Bug fix",
        "new_context": "Code review",
        "transition_type": "switch"
      },
      {
        "timestamp": "2025-10-04T17:00:00Z",
        "previous_context": "Code review",
        "new_context": null,
        "transition_type": "stop"
      }
    ]
  }
}
```

## Exit Codes

- `0`: Success (even if no transitions)

## Examples

```bash
$ my-context history
Context History:
  2025-10-04 14:00 │ START │ Bug fix
  ...

$ my-context h --limit 10
Context History (last 10):
  ...

$ my-context history --json | jq '.data.transitions | length'
5
```

**Command**: `my-context note <text>`  
**Alias**: `my-context n <text>`

## Purpose

Add a timestamped note to the currently active context.

## Arguments

**Positional**:
- `<text>` (required): The note text to add (supports spaces, multiline via quotes)

**Flags**:
- `--json`: Output result as JSON instead of human-readable text

## Input Validation

**Text Requirements**:
- Must not be empty
- Maximum length: 10,000 characters
- Newlines and pipe characters are automatically escaped for storage

## Behavior

1. Read `state.json` to get active_context_name
2. If null, exit with error "No active context"
3. Append to `~/.my-context/<name>/notes.log`: `<timestamp>|<text>\n`
4. Escape special characters: `\n` for newlines, `\|` for pipes

## Output

### Human-Readable (stdout)

**Success**:
```
Note added to context: Bug fix
```

### JSON Output (stdout)

```json
{
  "command": "note",
  "timestamp": "2025-10-04T14:05:00Z",
  "data": {
    "context_name": "Bug fix",
    "note_timestamp": "2025-10-04T14:05:00Z",
    "note_text": "Fixed authentication bug in login handler"
  }
}
```

### Error Output (stderr)

**No Active Context**:
```
Error: No active context
Start a context with: my-context start <name>
```

**Empty Note**:
```
Error: Note text required
Usage: my-context note <text>
```

## Exit Codes

- `0`: Success
- `1`: User error (no active context, empty note)
- `2`: System error (I/O failure)

## Examples

```bash
$ my-context note "Fixed authentication bug"
Note added to context: Bug fix

$ my-context n "Quick update"
Note added to context: Bug fix

$ my-context note "Multiline note
> with details
> across lines"
Note added to context: Bug fix
```

## Test Cases

1. Add note to active context → Appends to notes.log
2. Add note with no context → Error code 1
3. Add empty note → Error code 1
4. Add note with pipe character → Escaped as `\|`
5. Add multiline note → Escaped as `\n`
6. JSON output → Valid JSON
