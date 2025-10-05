# Command Contract: file

**Command**: `my-context file <path>`  
**Alias**: `my-context f <path>`

## Purpose

Associate a file path with the currently active context for tracking which files were modified during this work session.

## Arguments

**Positional**:
- `<path>` (required): File path to associate (absolute or relative)

**Flags**:
- `--json`: Output result as JSON instead of human-readable text

## Input Validation

**Path Requirements**:
- Must not be empty
- Converted to absolute path if relative
- Normalized to POSIX format (forward slashes) for storage

## Behavior

1. Read `state.json` to get active_context_name
2. If null, exit with error "No active context"
3. Normalize path: `filepath.Clean()` → `filepath.Abs()` → `filepath.ToSlash()`
4. Append to `~/.my-context/<name>/files.log`: `<timestamp>|<normalized_path>\n`

## Output

### Human-Readable (stdout)

**Success**:
```
File associated with context: Bug fix
  src/auth/login.go
```

### JSON Output (stdout)

```json
{
  "command": "file",
  "timestamp": "2025-10-04T14:05:00Z",
  "data": {
    "context_name": "Bug fix",
    "file_timestamp": "2025-10-04T14:05:00Z",
    "file_path": "C:/Users/Dev/src/auth/login.go",
    "original_path": "src/auth/login.go"
  }
}
```

## Exit Codes

- `0`: Success
- `1`: User error (no active context, empty path)
- `2`: System error (I/O failure)

## Examples

```bash
$ my-context file src/auth/login.go
File associated with context: Bug fix
  src/auth/login.go

$ my-context f C:\Users\Dev\file.txt
File associated with context: Bug fix
  C:\Users\Dev\file.txt
```

## Test Cases

1. Add relative path → Converts to absolute
2. Add Windows path → Normalizes to POSIX format
3. Add POSIX path → Stores as-is
4. Add with no context → Error code 1
