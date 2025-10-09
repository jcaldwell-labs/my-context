# Command Contract: archive

**Alias**: `a`  
**Purpose**: Mark context as archived (hidden from default list views)

## Syntax

```bash
my-context archive <context-name> [--json]
```

## Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| context-name | string | ✅ | Name of context to archive |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --json | -j | boolean | false | Output result as JSON |

## Behavior

**Input**: `my-context archive "ps-cli: Phase 1"`

**Process**:
1. Validate context exists
2. Load meta.json
3. Check if context is active (end_time is null)
4. If active, error: "Cannot archive active context. Stop it first: my-context stop"
5. Check if already archived
6. If already archived, info message: "Context 'name' is already archived"
7. Set is_archived = true in meta.json
8. Write updated meta.json atomically
9. Print success message

**Output** (stdout):
```
Archived context: ps-cli: Phase 1
```

**Exit Codes**:
- 0: Success (including already archived case)
- 1: Error (not found, active context, write failure)

## Output Format

### Human-Readable (Default)
```
Archived context: ps-cli: Phase 1
```

### JSON (--json)
```json
{
  "success": true,
  "context": "ps-cli: Phase 1",
  "message": "Archived context"
}
```

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Context not found | Error: "Context 'name' not found" (exit 1) |
| Active context | Error: "Cannot archive active context. Stop it first..." (exit 1) |
| Already archived | Info: "Context 'name' is already archived" (exit 0) |
| Write failure | Error: "Failed to archive context: <reason>" (exit 1) |

## Validation

**Pre-Conditions**:
- Context must exist
- Context must be stopped (end_time is set)
- Meta.json must be writable

**Post-Conditions**:
- meta.json contains `"is_archived": true`
- Context hidden from `my-context list` (without --archived flag)
- Context visible with `my-context list --archived`
- All context data preserved (notes, files, touches)

## Examples

```bash
# Archive completed phase
my-context archive "ps-cli: Phase 1"

# JSON output for scripting
my-context archive "Test Context" --json
```

