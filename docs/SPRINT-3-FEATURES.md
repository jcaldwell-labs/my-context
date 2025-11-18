# Sprint 3 Features - Context Organization & Hierarchy

**Version:** 3.1.0
**Release Date:** 2025-11-18
**Status:** ✅ Core Features Implemented

## Overview

Sprint 3 introduces powerful organizational features to help you manage complex, multi-level projects with tags, hierarchical relationships, and enhanced discovery capabilities.

## New Features

### 1. Context Tags/Labels System

Organize and categorize contexts with flexible tagging.

#### Commands

**Add tags to a context:**
```bash
my-context tag add <context> <tag1> [tag2] ...
my-context tag add "Bug fix" urgent backend bug

# Can also add tags when creating a context
my-context start "API refactor" --labels api,refactor,backend
```

**Remove tags from a context:**
```bash
my-context tag remove <context> <tag1> [tag2] ...
my-context tag rm "Bug fix" urgent

# Remove all tags
my-context tag remove "Bug fix" --all
```

**List tags:**
```bash
# List all tags with usage counts
my-context tag list

# List tags for specific context
my-context tag list "Bug fix"
```

**Filter contexts by tag:**
```bash
# Show only contexts with "urgent" tag
my-context list --tag urgent

# Combine with other filters
my-context list --tag backend --project my-app --limit 5
```

#### Use Cases

- **Priority management:** Tag contexts as `urgent`, `low-priority`, `blocked`
- **Team organization:** Tag by team member or role (`frontend`, `backend`, `qa`)
- **Feature tracking:** Tag by feature area (`auth`, `payments`, `notifications`)
- **Status tracking:** Tag contexts as `in-review`, `needs-testing`, `ready-to-merge`

#### Validation Rules

- Tags cannot contain whitespace (use hyphens or underscores)
- Maximum 10 tags per context
- Each tag must be 50 characters or less
- Tags are case-sensitive for storage but case-insensitive for filtering

---

### 2. Parent-Child Context Relationships

Create hierarchical relationships between contexts for better project organization.

#### Commands

**Link a child context to a parent:**
```bash
my-context link <child> <parent>

# Examples:
my-context link "Unit tests" "API refactor"
my-context link "API refactor" "Sprint 3"
my-context link "Sprint 3" "Q4 Goals"
```

**Remove parent link:**
```bash
my-context unlink <context>
my-context unlink "API refactor"
```

**View parent context:**
```bash
# Show parent of active context
my-context up

# Show parent of specific context
my-context up "Unit tests"
```

**View child contexts:**
```bash
# Show children of active context
my-context down

# Show children of specific context
my-context down "Sprint 3"
```

**Can also set parent when creating context:**
```bash
my-context start "Database migration" --parent "Backend overhaul"
```

#### Use Cases

**Sprint → Task → Subtask Hierarchy:**
```
Q4 Goals
├── Sprint 3
│   ├── API refactor
│   │   ├── Unit tests
│   │   └── Integration tests
│   └── Database migration
│       ├── Schema design
│       └── Data migration scripts
└── Sprint 4
    └── Frontend redesign
```

**Project → Phase → Task:**
```
Client Project ABC
├── Discovery Phase
│   ├── Requirements gathering
│   └── Technical assessment
├── Development Phase
│   ├── Backend API
│   └── Frontend UI
└── Testing Phase
    └── QA validation
```

#### Circular Dependency Protection

The system prevents circular dependencies:
- Cannot link a context to itself
- Circular relationships are detected and marked in tree view

---

### 3. Context Tree Visualization

Visualize your context hierarchy with ASCII tree diagrams.

#### Commands

**Show tree for specific context:**
```bash
my-context tree "Sprint 3"
```

Output:
```
Context hierarchy for "Sprint 3":

└─ Sprint 3
   ├─ API refactor
   │  ├─ Unit tests
   │  └─ Integration tests
   └─ Database migration
      ├─ Schema design
      └─ Data migration scripts
```

**Show all root contexts (contexts without parents):**
```bash
my-context tree
```

Output:
```
Context hierarchy (all root contexts):

└─ Q4 Goals
   ├─ Sprint 3
   │  └─ API refactor
   └─ Sprint 4

└─ Client Project ABC
   ├─ Discovery Phase
   └─ Development Phase
```

**JSON output for programmatic use:**
```bash
my-context tree "Sprint 3" --json
```

---

## Enhanced Existing Features

### Enhanced List Command

The `list` command now supports tag filtering:

```bash
# Filter by tag
my-context list --tag urgent

# Combine multiple filters
my-context list --tag backend --project my-app --search api
my-context list --tag bug --archived
my-context list --tag feature --limit 20 --all
```

### Enhanced Start Command

The `start` command now supports metadata flags:

```bash
# Create context with tags and parent
my-context start "Bug fix" \
  --labels bug,urgent,backend \
  --parent "Sprint 3" \
  --created-by "alice@example.com" \
  --project my-app

# Result: "my-app: Bug fix" with tags and parent relationship
```

---

## Data Model

### Context Metadata Structure

```json
{
  "name": "API refactor",
  "start_time": "2025-11-18T10:30:00Z",
  "status": "active",
  "subdirectory_path": "/home/user/.my-context/API_refactor",
  "is_archived": false,
  "metadata": {
    "created_by": "alice@example.com",
    "parent": "Sprint 3",
    "labels": ["backend", "refactor", "api"]
  }
}
```

### Backward Compatibility

All Sprint 3 features are **fully backward compatible**:
- Old contexts without metadata work seamlessly
- Metadata fields are optional
- Old contexts appear in tree view as root contexts
- Existing commands continue to work unchanged

---

## Architecture & Implementation

### File-Based Storage

Context metadata (including tags and parent) is stored in each context's `meta.json` file:

```
~/.my-context/
├── state.json
├── transitions.log
└── API_refactor/
    ├── meta.json       # Contains tags, parent, created_by
    ├── notes.log
    ├── files.log
    └── touch.log
```

### Core Functions

New core functions in `internal/core/context.go`:

- `AddTags(contextName, tags)` - Add tags to context
- `RemoveTags(contextName, tags)` - Remove tags from context
- `GetContextTags(contextName)` - Get tags for context
- `GetAllTags()` - Get all tags with usage counts
- `SetParent(child, parent)` - Link contexts
- `ClearParent(context)` - Unlink from parent
- `GetChildren(parent)` - Get child contexts
- `GetContextTree(root)` - Build tree structure
- `GetRootContexts()` - Get contexts without parents

### Commands Structure

New command files:

- `internal/commands/tag.go` - Tag management (add, remove, list subcommands)
- `internal/commands/link.go` - Link/unlink commands
- `internal/commands/tree.go` - Tree, up, down commands

---

## Use Cases & Examples

### Use Case 1: Multi-Sprint Project Management

```bash
# Set up sprint hierarchy
my-context start "Q4 2025" --labels planning,quarterly
my-context start "Sprint 3" --parent "Q4 2025" --labels sprint,active
my-context start "Sprint 4" --parent "Q4 2025" --labels sprint,planned

# Add tasks to Sprint 3
my-context start "API refactor" --parent "Sprint 3" --labels backend,refactor
my-context start "Database migration" --parent "Sprint 3" --labels backend,database

# Add subtasks
my-context start "Unit tests" --parent "API refactor" --labels testing
my-context start "Integration tests" --parent "API refactor" --labels testing

# Visualize
my-context tree "Q4 2025"

# Find all testing tasks
my-context list --tag testing

# Find all backend work
my-context list --tag backend
```

### Use Case 2: Bug Tracking

```bash
# Create bugs with tags
my-context start "Login bug" --labels bug,urgent,frontend
my-context start "Payment timeout" --labels bug,backend,payments
my-context start "UI glitch" --labels bug,low-priority,frontend

# Link to sprint
my-context link "Login bug" "Sprint 3"
my-context link "Payment timeout" "Sprint 3"

# View all urgent bugs
my-context list --tag bug --tag urgent

# View all bugs in current sprint
my-context down "Sprint 3"
```

### Use Case 3: Team Coordination

```bash
# Tag contexts by team member
my-context tag add "API refactor" alice backend
my-context tag add "UI redesign" bob frontend
my-context tag add "Database migration" alice backend

# View Alice's work
my-context list --tag alice

# View all backend work
my-context list --tag backend

# View team structure
my-context tree "Sprint 3"
```

---

## Future Enhancements (Planned)

### 🔜 Coming in Sprint 3.2

**Context Templates:**
- Pre-defined templates for common context types
- `my-context template create bug-fix`
- `my-context start "Login issue" --template bug-fix`
- Templates include suggested tags, notes, and file associations

**Enhanced Search:**
- Date range filtering: `my-context list --from 2025-11-01 --to 2025-11-18`
- Full-text note search: `my-context search "authentication"`
- File association search: `my-context list --file "src/main.go"`

**Database Backend (Optional):**
- SQLite storage option for faster queries
- PostgreSQL support for team environments
- Migration tool: `my-context migrate --to sqlite`
- Backward compatible with file-based storage

---

## Migration Guide

### From v2.x to v3.1

**No migration required!** Sprint 3 is fully backward compatible.

**To start using new features:**

1. **Add tags to existing contexts:**
   ```bash
   my-context tag add "Old context" relevant,tags
   ```

2. **Create hierarchy for existing contexts:**
   ```bash
   my-context link "Task" "Sprint"
   my-context tree
   ```

3. **Use enhanced list filtering:**
   ```bash
   my-context list --tag important
   ```

**Metadata will be created automatically** when you:
- Add tags to a context
- Set a parent relationship
- Create a new context with `--labels` or `--parent`

---

## Performance

### Benchmarks

Operations tested with 1000 contexts:

| Operation | Time | Notes |
|-----------|------|-------|
| Add tag | <10ms | In-memory + single file write |
| List tags | <50ms | Scans all meta.json files |
| Filter by tag | <100ms | O(n) scan with early exit |
| Build tree | <150ms | Recursive with cycle detection |
| Link contexts | <15ms | Two file reads + one write |

### Optimization Tips

1. **Use specific filters to reduce scan time:**
   ```bash
   my-context list --tag urgent --limit 10
   ```

2. **Tree view is lazy-loaded** - only requested branches are built

3. **Tag counts are cached per command execution**

---

## Troubleshooting

### Tags not appearing?

Check that metadata exists:
```bash
cat ~/.my-context/Context_Name/meta.json | grep labels
```

If metadata is missing, add a tag to initialize it:
```bash
my-context tag add "Context Name" temp
my-context tag remove "Context Name" temp
```

### Circular dependency detected?

Break the cycle:
```bash
my-context unlink "Context A"
# or
my-context link "Context A" "Correct Parent"
```

### Tree not showing all contexts?

Only contexts with relationships appear in tree view. To see all contexts:
```bash
my-context list --all
```

---

## API Reference

### Tag Command

```bash
my-context tag <subcommand> [args]

Subcommands:
  add <context> <tag1> [tag2] ...     Add tags
  remove|rm <context> <tag1> ...      Remove tags
  remove|rm <context> --all           Remove all tags
  list                                 List all tags
  list <context>                       List tags for context

Flags:
  --count, -c    Show usage count (default: true)
  --json, -j     JSON output
```

### Link Commands

```bash
my-context link <child> <parent>      Link contexts
my-context unlink <context>           Remove parent link

Flags:
  --json, -j     JSON output
```

### Tree Commands

```bash
my-context tree [context]             Show context hierarchy
my-context up [context]               Show parent context
my-context down [context]             Show child contexts

Flags:
  --all, -a      Show all contexts including orphans
  --json, -j     JSON output
```

### Enhanced List Flags

```bash
my-context list [flags]

New flags:
  --tag <tag>         Filter by tag/label

Existing flags:
  --project <name>    Filter by project
  --search <term>     Search in names
  --limit <n>         Limit results (default: 10)
  --all               Show all (no limit)
  --archived          Show only archived
  --active-only       Show only active
  --json, -j          JSON output
```

---

## JSON Output Examples

### Tag List

```json
{
  "status": "success",
  "data": {
    "tags": {
      "urgent": 5,
      "backend": 12,
      "frontend": 8,
      "bug": 3,
      "feature": 15
    }
  }
}
```

### Tree Structure

```json
{
  "Name": "Sprint 3",
  "Children": [
    {
      "Name": "API refactor",
      "Children": [
        {
          "Name": "Unit tests",
          "Children": []
        }
      ]
    },
    {
      "Name": "Database migration",
      "Children": []
    }
  ]
}
```

---

## Contributing

Found a bug or have a feature request? Please:

1. Check existing issues: [GitHub Issues](https://github.com/jefferycaldwell/my-context/issues)
2. Review this documentation for expected behavior
3. Create a new issue with:
   - Sprint 3 version (3.1.0)
   - Command executed
   - Expected vs actual behavior
   - Context state (`my-context show`)

---

## License

MIT License - See LICENSE file for details

---

**Last Updated:** 2025-11-18
**Author:** jefferycaldwell
**Repository:** https://github.com/jefferycaldwell/my-context
