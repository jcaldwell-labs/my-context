# Command Contract: start (Enhanced with --project)

**Alias**: `s`  
**Purpose**: Start a new context, optionally with project prefix

## Syntax

```bash
my-context start <context-name> [--project PROJECT_NAME] [--json]
```

## Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| context-name | string | ✅ | Name/phase of the context |

## Flags (New in Sprint 2)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --project | -p | string | "" | Project name to prefix (creates "project: name") |
| --json | -j | boolean | false | Output result as JSON |

## Behavior

### With Project Flag

**Input**: `my-context start "Phase 3" --project ps-cli`

**Process**:
1. Combine as: "ps-cli: Phase 3" (project + ": " + name)
2. Trim whitespace from both parts
3. Validate combined name (unique, valid directory name)
4. Stop currently active context if exists
5. Create context directory
6. Write meta.json with combined name
7. Record transition in transitions.log
8. Print success message

**Output**:
```
Started context: ps-cli: Phase 3
```

### Without Project Flag (Sprint 1 Behavior Preserved)

**Input**: `my-context start "Standalone Context"`

**Process**:
- Use name as-is
- No automatic prefix added

**Output**:
```
Started context: Standalone Context
```

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Both project and name contain colons | Combined: "project: name:with:colons" (valid) |
| Project flag with empty name | Error: "Context name required" (exit 1) |
| Duplicate name after combination | Error: "Context 'name' already exists" (exit 1) |
| Project name with trailing colon | Auto-trim: "--project ps-cli:" → "ps-cli: Phase 1" |

## Validation

**Pre-Conditions**:
- Combined name must be unique
- Combined name must be valid directory name
- If active context exists, it can be stopped

**Post-Conditions**:
- Context created with combined name if --project used
- Previous context stopped (if existed)
- New context is active
- Transition recorded

## Examples

```bash
# Start with project prefix
my-context start "Phase 2" --project ps-cli
# Creates: "ps-cli: Phase 2"

# Start with existing project (user types manually)
my-context start "ps-cli: Phase 2"
# Creates: "ps-cli: Phase 2" (same result)

# Standalone context
my-context start "Quick Fix"
# Creates: "Quick Fix"

# JSON output
my-context start "Phase 3" --project garden --json
```

## Implementation Notes

- Whitespace handling: trim both project and name before combining
- Colon separator: always use ": " (colon + space)
- Backward compatible: --project flag is optional
- User can still type full name manually (Sprint 1 style)

