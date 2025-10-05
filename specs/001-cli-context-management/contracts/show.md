# Command Contract: show

**Command**: `my-context show`  
**Alias**: `my-context w`

## Purpose

Display details about the currently active context including name, start time, duration, recent notes, associated files, and touch events.

## Arguments

**Flags**:
- `--json`: Output result as JSON instead of human-readable text

## Behavior

1. Read `state.json` to get active_context_name
2. If null, exit with "No active context" message
3. Read `~/.my-context/<name>/meta.json`
4. Read all lines from `notes.log`, `files.log`, `touch.log`
5. Calculate duration (now - start_time)
6. Format and output

## Output

### Human-Readable (stdout)

**Active Context**:
```
Context: Bug fix
Status: active
Started: 2025-10-04 14:00:00 (2h 15m ago)

Notes (3):
  [14:05] Fixed authentication bug in login handler
  [14:30] Updated tests to cover edge case
  [15:00] Refactored error handling

Files (2):
  [14:05] src/auth/login.go
  [14:30] tests/auth/login_test.go

Activity: 5 touches (last: 15:45)
```

**No Active Context**:
```
No active context
Start one with: my-context start <name>
```

### JSON Output (stdout)

**Active Context**:
```json
{
  "command": "show",
  "timestamp": "2025-10-04T16:15:00Z",
  "data": {
    "context": {
      "name": "Bug fix",
      "start_time": "2025-10-04T14:00:00Z",
      "end_time": null,
      "status": "active",
      "duration_seconds": 8100,
      "notes": [
        {
          "timestamp": "2025-10-04T14:05:00Z",
          "text": "Fixed authentication bug in login handler"
        },
        {
          "timestamp": "2025-10-04T14:30:00Z",
          "text": "Updated tests to cover edge case"
        },
        {
          "timestamp": "2025-10-04T15:00:00Z",
          "text": "Refactored error handling"
        }
      ],
      "files": [
        {
          "timestamp": "2025-10-04T14:05:00Z",
          "path": "src/auth/login.go"
        },
        {
          "timestamp": "2025-10-04T14:30:00Z",
          "path": "tests/auth/login_test.go"
        }
      ],
      "touches": [
        {
          "timestamp": "2025-10-04T14:15:00Z"
        },
        {
          "timestamp": "2025-10-04T15:00:00Z"
        },
        {
          "timestamp": "2025-10-04T15:45:00Z"
        }
      ]
    }
  }
}
```
# Command Contract: start
**No Active Context**:
```json
{
  "command": "show",
  "timestamp": "2025-10-04T16:15:00Z",
  "error": {
    "code": 1,
    "message": "No active context"
  }
}
```

### Error Output (stderr)

**No Active Context**:
```
No active context
Start one with: my-context start <name>
```

**Context Directory Missing** (corrupted state):
```
Error: Active context "Bug fix" not found
Context directory missing: ~/.my-context/Bug fix/
Run 'my-context list' to see available contexts
```

## Exit Codes

- `0`: Success (context displayed)
- `1`: No active context (not an error condition, just informational)
- `2`: System error (I/O failure, corrupted state)

## Side Effects

**None** - Read-only operation

## Examples

### Basic Usage
```bash
$ my-context show
Context: Bug fix
Status: active
Started: 2025-10-04 14:00:00 (2h 15m ago)

Notes (3):
  [14:05] Fixed authentication bug in login handler
  [14:30] Updated tests to cover edge case
  [15:00] Refactored error handling

Files (2):
  [14:05] src/auth/login.go
  [14:30] tests/auth/login_test.go

Activity: 5 touches (last: 15:45)
```

### With Alias
```bash
$ my-context w
Context: Bug fix
Status: active
Started: 2025-10-04 14:00:00 (2h 15m ago)
...
```

### JSON Output
```bash
$ my-context show --json
{"command":"show","timestamp":"2025-10-04T16:15:00Z","data":{"context":{...}}}
```

### Piping
```bash
# Get just the context name
$ my-context show | head -1 | cut -d: -f2 | xargs
Bug fix

# With JSON and jq
$ my-context show --json | jq -r '.data.context.name'
Bug fix

# Count notes
$ my-context show --json | jq '.data.context.notes | length'
3

# Get latest note
$ my-context show --json | jq -r '.data.context.notes[-1].text'
Refactored error handling
```

### Cross-Platform Path Display
```bash
# Windows (cmd) - shows native paths
$ my-context show
Files (1):
  [14:05] C:\Users\Dev\src\auth\login.go

# WSL/Linux - shows native paths
$ my-context show
Files (1):
  [14:05] /home/dev/src/auth/login.go

# JSON always uses POSIX format
$ my-context show --json | jq -r '.data.context.files[0].path'
C:/Users/Dev/src/auth/login.go
```

## Test Cases

**Must Pass**:
1. Show active context → Displays all details
2. Show when no context → Exit code 1, helpful message
3. Show with empty notes → "Notes (0)"
4. Show with empty files → "Files (0)"
5. Show with no touches → "Activity: 0 touches"
6. JSON output → Valid JSON, all fields present
7. Duration calculation → Correct elapsed time
8. Cross-platform paths → Native format in human output, POSIX in JSON

**Command**: `my-context start <name>`  
**Alias**: `my-context s <name>`

## Purpose

Create and activate a new context with the given name. If a context is already active, automatically stop it before starting the new one.

## Arguments

**Positional**:
- `<name>` (required): The name of the context to create

**Flags**:
- `--json`: Output result as JSON instead of human-readable text

## Input Validation

**Name Requirements**:
- Must not be empty
- Must not contain path separators (`/` or `\`)
- Must not contain control characters
- Maximum length: 200 characters (before automatic suffix)

**Invalid Examples**:
```bash
my-context start ""                    # Error: Name required
my-context start "bug/fix"             # Error: Invalid character '/'
my-context start "very-long-name..."   # Error: Name too long (>200 chars)
```

## Behavior

### No Active Context
1. Check if context directory exists
2. If exists, resolve duplicate name (append _2, _3, etc.)
3. Create context directory: `~/.my-context/<name>/`
4. Write `meta.json` with start_time, status="active"
5. Create empty log files (notes.log, files.log, touch.log)
6. Update `state.json` with active_context=<name>
7. Append to `transitions.log`: `<timestamp>|NULL|<name>|start`

### Active Context Exists
1. Check if context directory exists for new name
2. Resolve duplicate name if needed
3. Read current active context from `state.json`
4. Update current context's `meta.json`: end_time=now, status="stopped"
5. Create new context directory and files (as above)
6. Update `state.json` with new active_context
7. Append to `transitions.log`: `<timestamp>|<old>|<new>|switch`

## Output

### Human-Readable (stdout)

**Success - New context**:
```
Started context: Bug fix
```

**Success - Duplicate name handled**:
```
Context "Bug fix" already exists.
Started context: Bug fix_2
```

**Success - Switched from previous**:
```
Stopped context: Code review (duration: 1h 23m)
Started context: Bug fix
```

### JSON Output (stdout)

**Success**:
```json
{
  "command": "start",
  "timestamp": "2025-10-04T14:00:00Z",
  "data": {
    "context_name": "Bug fix",
    "original_name": "Bug fix",
    "was_duplicate": false,
    "previous_context": null,
    "previous_duration_seconds": null
  }
}
```

**Success with duplicate resolution**:
```json
{
  "command": "start",
  "timestamp": "2025-10-04T14:00:00Z",
  "data": {
    "context_name": "Bug fix_2",
    "original_name": "Bug fix",
    "was_duplicate": true,
    "previous_context": "Code review",
    "previous_duration_seconds": 4980
  }
}
```

### Error Output (stderr)

**Missing name**:
```
Error: Context name required
Usage: my-context start <name>
```

**Invalid name**:
```
Error: Invalid context name "bug/fix"
Context names cannot contain path separators (/ or \)
```

**System error (I/O failure)**:
```
Error: Failed to create context directory
Permission denied: ~/.my-context/Bug fix/
```

### JSON Error Output (stdout when --json)

```json
{
  "command": "start",
  "timestamp": "2025-10-04T14:00:00Z",
  "error": {
    "code": 1,
    "message": "Context name required"
  }
}
```

## Exit Codes

- `0`: Success
- `1`: User error (invalid name, missing argument)
- `2`: System error (I/O failure, permissions)

## Side Effects

**Filesystem Changes**:
- Creates `~/.my-context/<name>/` directory
- Creates 4 files: `meta.json`, `notes.log`, `files.log`, `touch.log`
- Updates `~/.my-context/state.json`
- Appends line to `~/.my-context/transitions.log`
- If previous context existed: Updates previous context's `meta.json`

**State Changes**:
- Global active context changes to new context
- Previous context (if any) changes to stopped status

## Examples

### Basic Usage
```bash
$ my-context start "Fixing login bug"
Started context: Fixing login bug

$ my-context start "Code review"
Stopped context: Fixing login bug (duration: 45m 12s)
Started context: Code review
```

### With Alias
```bash
$ my-context s "Quick fix"
Started context: Quick fix
```

### Duplicate Name Handling
```bash
$ my-context start "Bug fix"
Started context: Bug fix

$ my-context start "Bug fix"
Context "Bug fix" already exists.
Started context: Bug fix_2

$ my-context start "Bug fix"
Context "Bug fix" already exists.
Context "Bug fix_2" already exists.
Started context: Bug fix_3
```

### JSON Output
```bash
$ my-context start "New feature" --json
{"command":"start","timestamp":"2025-10-04T14:00:00Z","data":{"context_name":"New feature","original_name":"New feature","was_duplicate":false,"previous_context":null,"previous_duration_seconds":null}}
```

### Piping
```bash
# Extract just the context name
$ my-context start "Deploy" | grep -oP '(?<=: ).*'
Deploy

# With JSON and jq
$ my-context start "Deploy" --json | jq -r '.data.context_name'
Deploy
```

## Cross-Platform Notes

**Path Handling**:
- Context name sanitized same way on all platforms
- Directory created using native path separators internally
- User sees consistent behavior across cmd/bash/WSL

**Shell Quoting**:
- Names with spaces require quotes in all shells:
  - cmd: `my-context start "Bug fix"`
  - bash: `my-context start "Bug fix"` or `my-context start 'Bug fix'`
  - PowerShell: `my-context start "Bug fix"`

## Test Cases

**Must Pass**:
1. Start context with simple name → Creates directory, updates state
2. Start context when one active → Stops previous, starts new
3. Start duplicate name → Appends _2 suffix
4. Start with name containing spaces → Handles correctly
5. Start with empty name → Rejects with error code 1
6. Start with name containing `/` → Rejects with error code 1
7. JSON output → Valid JSON parseable by jq
8. Cross-shell → Same behavior in cmd/bash/WSL
