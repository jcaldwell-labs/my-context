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
