# Command Contract: delete

**Alias**: `d`  
**Purpose**: Permanently remove a context and all its data

## Syntax

```bash
my-context delete <context-name> [--force] [--json]
```

## Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| context-name | string | ✅ | Name of context to delete |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --force | -f | boolean | false | Skip confirmation prompt |
| --json | -j | boolean | false | Output result as JSON |

## Behavior

**Input**: `my-context delete "Test Context"`

**Process**:
1. Validate context exists
2. Check if context is active
3. If active, error: "Cannot delete active context. Stop it first: my-context stop"
4. If not --force, prompt: "Delete context 'name' permanently? This cannot be undone. (y/N)"
5. Read user input from stdin
6. If not 'y' or 'yes' (case-insensitive), exit with code 1
7. Remove entire context directory from ~/.my-context/
8. Preserve transitions.log (do not remove transition entries)
9. Print success message

**Output** (stdout):
```
Delete context 'Test Context' permanently? This cannot be undone. (y/N) y
Context 'Test Context' deleted
```

**Exit Codes**:
- 0: Success
- 1: Error (not found, active context, user cancellation, delete failure)

## Output Format

### Human-Readable (Default)
```
Context 'Test Context' deleted
```

### JSON (--json, with --force)
```json
{
  "success": true,
  "context": "Test Context",
  "message": "Context deleted"
}
```

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Context not found | Error: "Context 'name' not found" (exit 1) |
| Active context | Error: "Cannot delete active context..." (exit 1) |
| User cancels | No changes, exit 1 |
| Directory delete fails | Error: "Failed to delete context: <reason>" (exit 1) |
| Context referenced in transitions.log | Delete context, preserve log entries |

## Validation

**Pre-Conditions**:
- Context must exist
- Context must not be active
- User must confirm (unless --force)

**Post-Conditions**:
- Context directory removed from ~/.my-context/
- Context no longer appears in any list view
- transitions.log preserved with historical entries intact
- No orphaned data remains

## Examples

```bash
# Delete with confirmation
my-context delete "Test Context"

# Force delete (no prompt)
my-context delete "Old Context" --force

# JSON output for scripting
my-context delete "Temp" --force --json
```

## Implementation Notes

- Confirmation prompt reads from stdin (supports piping: `echo y | my-context delete "name"`)
- Accept 'y', 'Y', 'yes', 'Yes', 'YES' as affirmative responses
- Any other input (including Enter) is treated as cancellation
- Remove directory recursively (all subdirectories and files)
- transitions.log is at ~/.my-context/transitions.log (not in context directory, so automatically preserved)

