# Contract: Project Filter (Start Command Enhancement)

**Command**: `start` (alias: `s`) - Enhanced with --project flag  
**Purpose**: Create contexts with project prefix following "project: phase" convention

---

## Signature

```bash
my-context start <name> [--project <project-name>] [--json]
my-context s <name> [--project <project-name>] [--json]
```

---

## Arguments

**Positional**:
- `<name>`: Context name or phase name (required)

**Flags**:
- `--project <project-name>`: Prepend project name to context (optional)
- `--json`: Output JSON format (optional)

---

## Behavior

### Start Without Project Flag (Existing Behavior)

**Input**: `my-context start "Working on feature"`

**Actions**:
1. Stop any active context
2. Create context named exactly "Working on feature"
3. Normalize name for directory (spaces → underscores)
4. Create context directory and meta.json
5. Update state.json
6. Log transition

**Output**:
```
Started context: Working on feature
```

**Exit Code**: 0

### Start With Project Flag

**Input**: `my-context start "Phase 1 - Foundation" --project ps-cli`

**Actions**:
1. Stop any active context
2. Combine project + name: "ps-cli: Phase 1 - Foundation"
3. Check for duplicate name (apply _2 suffix if needed)
4. Create context with combined name
5. Update state.json
6. Log transition

**Output**:
```
Started context: ps-cli: Phase 1 - Foundation
```

**Exit Code**: 0

### Start With Project Flag (Short Phase Name)

**Input**: `my-context start "Phase 2" --project ps-cli`

**Actions**:
1. Combine: "ps-cli: Phase 2"
2. Create context

**Output**:
```
Started context: ps-cli: Phase 2
```

**Exit Code**: 0

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05
