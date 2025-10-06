# Contract: Delete Command

**Command**: `delete` (alias: `d`)  
**Purpose**: Permanently remove a context and all associated data

---

## Signature

```bash
my-context delete <context-name> [--force] [--json]
my-context d <context-name> [--force] [--json]
```

---

## Arguments

**Positional**:
- `<context-name>`: Name of context to delete (required)

**Flags**:
- `--force`: Skip confirmation prompt (for scripting)
- `--json`: Output JSON format (optional)

---

## Behavior

### Delete with Confirmation

**Input**: `my-context delete "Test Context"`

**Actions**:
1. Verify context exists
2. Verify context is not active
3. Display confirmation prompt
4. Wait for user input (y/N)
5. If confirmed: Remove entire `~/.my-context/Test_Context/` directory
6. Print success message

**Output**:
```
⚠️  This will permanently delete context "Test Context" and all associated data.
   - Notes: 15
   - Files: 3
   - Touch events: 8

Continue? (y/N): y

Deleted context: Test Context
```

**Exit Code**: 0 (success), 2 (cancelled by user)

### Delete with --force

**Input**: `my-context delete "Test Context" --force`

**Actions**:
1. Verify context exists
2. Verify context is not active
3. Skip confirmation
4. Remove directory immediately
5. Print success message

**Output**:
```
Deleted context: Test Context
```

**Exit Code**: 0 (success)

---

## Error Handling

### Context Not Found

**Input**: `my-context delete "NonExistent"`

**Output**:
```
Error: Context "NonExistent" not found
Run 'my-context list --all' to see available contexts.
```

**Exit Code**: 1
# Contract: Export Command
### Context Is Active

**Input**: `my-context delete "CurrentWork"` (when CurrentWork is active)

**Output**:
```
Error: Cannot delete active context "CurrentWork"
Stop the context first: my-context stop
```

**Exit Code**: 3

### User Cancels Deletion

**Input**: `my-context delete "Test"` → User enters "N"

**Output**:
```
⚠️  This will permanently delete context "Test" and all associated data.
Continue? (y/N): N

Deletion cancelled.
```

**Exit Code**: 2

### Filesystem Error

**Input**: `my-context delete "Locked"` (permission denied)

**Output**:
```
Error: Cannot delete context directory: permission denied
Path: /home/user/.my-context/Locked/
```

**Exit Code**: 2

---

## Data Removal

**Deleted**:
- `~/.my-context/{context_name}/` (entire directory)
  - `meta.json`
  - `notes.log`
  - `files.log`
  - `touch.log`
  - Any other files in context directory

**Preserved**:
- `~/.my-context/state.json` (global state)
- `~/.my-context/transitions.log` (historical record)
- Other contexts (unaffected)

**Important**: Deletion does NOT remove entries from `transitions.log`. Historical transitions remain for audit purposes.

---

## Side Effects

1. **List command**: Context no longer appears in any list view
2. **Show command**: Returns "context not found" error
3. **Export command**: Returns "context not found" error
4. **History command**: Context name still visible in transitions (historical record preserved)
5. **State file**: If deleted context was last stopped, no change (state references active context only)

---

## Integration Points

**Reads from**:
- `internal/core/context.go`: LoadContext()
- `internal/core/state.go`: GetActiveContext()

**Writes to**:
- Filesystem: Removes `~/.my-context/{context_name}/` directory

**Depends on**:
- Context existence check
- Active context check
- Confirmation prompt (unless --force)

---

## Validation Rules

1. Context name must be non-empty string
2. Context must exist in `~/.my-context/`
3. Context must not be active (status != "active")
4. Confirmation required unless --force flag provided
5. Deletion must be atomic (all-or-nothing)

---

## Confirmation Prompt Details

**Format**:
```
⚠️  This will permanently delete context "{name}" and all associated data.
   - Notes: {count}
   - Files: {count}
   - Touch events: {count}

Continue? (y/N):
```

**Accepted Inputs**:
- `y`, `Y`, `yes`, `Yes`, `YES` → Proceed with deletion
- `n`, `N`, `no`, `No`, `NO`, `<Enter>` → Cancel deletion
- Any other input → Cancel deletion (default to safe option)

---

## JSON Output Format

**Input**: `my-context delete "Test" --force --json`

**Success**:
```json
{
  "command": "delete",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "context_name": "Test",
    "deleted": true,
    "items_removed": {
      "notes": 15,
      "files": 3,
      "touches": 8
    }
  }
}
```

**Cancelled**:
```json
{
  "command": "delete",
  "timestamp": "2025-10-05T20:00:00Z",
  "error": "Deletion cancelled by user",
  "data": {
    "context_name": "Test",
    "deleted": false
  }
}
```

---

## Testing Checklist

- [ ] Delete stopped context with confirmation (user accepts)
- [ ] Delete stopped context with confirmation (user cancels)
- [ ] Delete stopped context with --force (no prompt)
- [ ] Delete active context (error)
- [ ] Delete non-existent context (error)
- [ ] Verify entire directory removed from filesystem
- [ ] Verify context not in list output after deletion
- [ ] Verify transitions.log preserved (historical data intact)
- [ ] Verify other contexts unaffected
- [ ] Verify JSON output format
- [ ] Test confirmation prompt with various inputs (y, n, enter, garbage)
- [ ] Test deletion with insufficient permissions (error handling)

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05

**Command**: `export` (alias: `e`)  
**Purpose**: Generate markdown representation of context data for sharing and archival

---

## Signature

```bash
my-context export <context-name> [--to <path>] [--all] [--json]
my-context e <context-name> [--to <path>] [--all] [--json]
```

---

## Arguments

**Positional**:
- `<context-name>`: Name of context to export (optional if --all used)

**Flags**:
- `--to <path>`: Output file path (default: `./{context_name}.md`)
- `--all`: Export all contexts to separate files
- `--json`: Output JSON instead of markdown (optional, for programmatic use)

---

## Behavior

### Single Context Export

**Input**: `my-context export "ps-cli: Phase 1"`

**Actions**:
1. Verify context exists in `~/.my-context/`
2. Read `meta.json`, `notes.log`, `files.log`, `touch.log`
3. Generate markdown according to data-model.md ExportDocument format
4. Write to `./ps-cli_Phase_1.md` (sanitize filename)
5. Print success message with file path

**Output**:
```
Exported context "ps-cli: Phase 1" to ./ps-cli_Phase_1.md
```

**Exit Code**: 0 (success)

### Export with Custom Path

**Input**: `my-context export "Phase 1" --to contexts/phase-1-summary.md`

**Actions**:
1. Create parent directories if needed (`mkdir -p contexts/`)
2. Check if file exists at target path
3. If exists: Prompt for overwrite confirmation (y/N)
4. Write markdown to specified path
5. Print success message

**Output**:
```
Exported context "Phase 1" to contexts/phase-1-summary.md
```

**Exit Code**: 0 (success), 2 (user cancelled overwrite)

### Export All Contexts

**Input**: `my-context export --all --to reports/`

**Actions**:
1. List all contexts (including archived)
2. Create output directory if needed
3. For each context:
   - Generate sanitized filename
   - Export to `reports/{context_name}.md`
   - Print progress
4. Print summary

**Output**:
```
Exporting 15 contexts to reports/...
  ✓ ps-cli_Phase_1.md
  ✓ ps-cli_Phase_2.md
  ...
Exported 15 contexts to reports/
```

**Exit Code**: 0 (success)

---

## Error Handling

### Context Not Found

**Input**: `my-context export "NonExistent"`

**Output**:
```
Error: Context "NonExistent" not found
Run 'my-context list --all' to see available contexts.
```

**Exit Code**: 1

### File Write Error

**Input**: `my-context export "Test" --to /root/protected.md` (no permission)

**Output**:
```
Error: Cannot write to /root/protected.md: permission denied
```

**Exit Code**: 2

### No Context Name with --all Missing

**Input**: `my-context export` (no args)

**Output**:
```
Error: Context name required or use --all to export all contexts
Usage: my-context export <context-name> [--to <path>] [--all]
```

**Exit Code**: 1

---

## File Format

### Markdown Output

See `data-model.md` for complete ExportDocument format.

**Key Requirements**:
- Timestamps in local timezone with ISO 8601 format
- Duration in human-readable format ("2h 15m", "3d 4h")
- Notes with time prefixes (`` `HH:MM` ``)
- File paths with "Added:" timestamps
- Footer with export timestamp and tool version

### JSON Output (with --json flag)

```json
{
  "command": "export",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "context": {
      "name": "ps-cli: Phase 1",
      "start_time": "2025-10-04T10:00:00Z",
      "end_time": "2025-10-04T18:00:00Z",
      "status": "stopped",
      "is_archived": false
    },
    "notes": [
      {"timestamp": "2025-10-04T10:15:00Z", "text": "Initial planning"},
      ...
    ],
    "files": [
      {"timestamp": "2025-10-04T11:00:00Z", "path": "/path/to/file.go"},
      ...
    ],
    "activity": {
      "touch_count": 12,
      "last_touch": "2025-10-04T17:45:00Z"
    },
    "export_path": "./ps-cli_Phase_1.md"
  }
}
```

---

## Validation Rules

1. Context name must be non-empty string
2. Output path must be valid filesystem path
3. If --to omitted, use current directory (must be writable)
4. Sanitize context name for filename (replace special chars with underscores)
5. Confirm overwrite if target file exists (unless --force flag added in future)

---

## Integration Points

**Reads from**:
- `internal/core/context.go`: LoadContext()
- `internal/core/storage.go`: ReadLog()

**Writes to**:
- Filesystem: User-specified or default path

**Depends on**:
- Context existence check
- Markdown formatter (new in `internal/output/markdown.go`)

---

## Testing Checklist

- [ ] Export single context with default path
- [ ] Export single context with custom --to path
- [ ] Export all contexts to directory
- [ ] Export non-existent context (error handling)
- [ ] Export with special characters in context name (filename sanitization)
- [ ] Export to existing file (overwrite prompt)
- [ ] Export to read-only directory (permission error)
- [ ] Export with --json flag (JSON output format)
- [ ] Export context with no notes/files (empty sections)
- [ ] Export active vs stopped vs archived contexts (all work)
- [ ] Verify markdown renders correctly in GitHub/VS Code
- [ ] Verify timestamps in export are in local timezone

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05
