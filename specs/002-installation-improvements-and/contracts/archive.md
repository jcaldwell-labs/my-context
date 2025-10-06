# Contract: Archive Command

**Command**: `archive` (alias: `a`)  
**Purpose**: Mark a context as archived (completed work, hidden from default list view)

---

## Signature

```bash
my-context archive <context-name> [--json]
my-context a <context-name> [--json]
```

---

## Arguments

**Positional**:
- `<context-name>`: Name of context to archive (required)

**Flags**:
- `--json`: Output JSON format (optional)

---

## Behavior

### Archive Stopped Context

**Input**: `my-context archive "ps-cli: Phase 1"`

**Actions**:
1. Verify context exists
2. Verify context is stopped (status != "active")
3. Check if already archived
4. Update `meta.json`: Set `"is_archived": true`
5. Print success message

**Output**:
```
Archived context: ps-cli: Phase 1
Use 'my-context list --archived' to see archived contexts.
```

**Exit Code**: 0 (success)

---

## Error Handling

### Context Not Found

**Input**: `my-context archive "NonExistent"`

**Output**:
```
Error: Context "NonExistent" not found
Run 'my-context list --all' to see available contexts.
```

**Exit Code**: 1

### Context Is Active

**Input**: `my-context archive "CurrentWork"` (when CurrentWork is active)

**Output**:
```
Error: Cannot archive active context "CurrentWork"
Stop the context first: my-context stop
```

**Exit Code**: 3

### Already Archived

**Input**: `my-context archive "ps-cli: Phase 1"` (when already archived)

**Output**:
```
Context "ps-cli: Phase 1" is already archived.
```

**Exit Code**: 0 (idempotent operation, not an error)

---

## State Transition

**Before**:
```json
{
  "name": "ps-cli: Phase 1",
  "status": "stopped",
  "is_archived": false
}
```

**After**:
```json
{
  "name": "ps-cli: Phase 1",
  "status": "stopped",
  "is_archived": true
}
```

---

## Side Effects

1. **List command**: Context hidden from default `my-context list` output
2. **List --archived**: Context visible in archived-only view
3. **List --all**: Context visible (all contexts regardless of archive status)
4. **Show/Export**: Context remains accessible (archiving doesn't hide data)
5. **Transitions log**: No entry added (archive is metadata change, not lifecycle event)

---

## Integration Points

**Reads from**:
- `internal/core/context.go`: LoadContext()
- `internal/core/state.go`: GetActiveContext()

**Writes to**:
- `~/.my-context/{context_name}/meta.json`: Update is_archived field

**Depends on**:
- Context existence check
- Active context check

---

## Validation Rules

1. Context name must be non-empty string
2. Context must exist in `~/.my-context/`
3. Context status must be "stopped" (not "active")
4. Archive operation is idempotent (safe to call on already-archived context)

---

## JSON Output Format

**Input**: `my-context archive "Phase 1" --json`

**Output**:
```json
{
  "command": "archive",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "context_name": "Phase 1",
    "was_archived": false,
    "now_archived": true
  }
}
```

**If already archived**:
```json
{
  "command": "archive",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "context_name": "Phase 1",
    "was_archived": true,
    "now_archived": true,
    "message": "Already archived"
  }
}
```

---

## Testing Checklist

- [ ] Archive stopped context (success)
- [ ] Archive active context (error)
- [ ] Archive non-existent context (error)
- [ ] Archive already-archived context (idempotent)
- [ ] Verify archived context hidden from default list
- [ ] Verify archived context visible with --archived flag
- [ ] Verify archived context visible with --all flag
- [ ] Verify meta.json updated correctly
- [ ] Verify show/export still work on archived contexts
- [ ] Verify JSON output format
- [ ] Verify backward compatibility (Sprint 1 contexts work)

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05
