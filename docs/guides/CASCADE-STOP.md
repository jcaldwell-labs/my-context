# Cascade Stop Feature Guide

## Overview

The cascade stop feature allows you to automatically stop all child contexts when stopping a parent context. This helps prevent orphaned active contexts in hierarchical workflows.

## Parent-Child Context Hierarchies

My-context supports organizing contexts in parent-child relationships:

```bash
# Create a parent context
my-context start "sprint-planning"

# Create child contexts linked to parent
my-context start "story-123"
my-context link story-123 sprint-planning

my-context start "bug-456"
my-context link bug-456 sprint-planning
```

## Default Behavior: Warning Mode

By default, when you stop a parent context that has active children, you'll see a warning:

```bash
$ my-context resume sprint-planning
$ my-context stop

⚠️  Warning: 2 active child context(s) will remain active:
  - story-123 (active, 2h 30m)
  - bug-456 (active, 45m)

Use --cascade to stop them along with the parent.

✓ Stopped: sprint-planning (4h 15m)
```

The warning:
- Lists all active child contexts
- Shows how long each has been active
- Suggests using `--cascade` flag
- Still stops the parent context

## Cascade Mode: Stop All Children

Use the `--cascade` flag to automatically stop all descendant contexts (children, grandchildren, etc.):

```bash
$ my-context stop --cascade

✓ Stopped: bug-456 (45m)
✓ Stopped: story-123 (2h 30m)
✓ Stopped: sprint-planning (4h 15m)
```

Cascade stop:
- Stops descendants in reverse order (leaves first, then parents)
- Handles nested hierarchies automatically
- Works with both file-based and database backends
- Records accurate stop times for all contexts

## Nested Hierarchies

Cascade stop works with any depth of hierarchy:

```bash
# Create nested structure:
# parent -> child1, child2 -> grandchild1, grandchild2

my-context start "parent"
my-context stop

my-context start "child1"
my-context link child1 parent

my-context start "child2"  
my-context link child2 parent

my-context start "grandchild1"
my-context link grandchild1 child1

my-context start "grandchild2"
my-context link grandchild2 child2

# Manually mark them all as active (in normal use, only one is active)
# Then stop parent with cascade
my-context resume parent
my-context stop --cascade

# Result: All descendants stopped in order:
# grandchild2 -> grandchild1 -> child2 -> child1 -> parent
```

## JSON Output

### Warning Mode
```json
{
  "command": "stop",
  "timestamp": "2026-01-25T03:24:25Z",
  "data": {
    "data": {
      "context_name": "parent",
      "start_time": "2026-01-25T03:23:00Z",
      "end_time": "2026-01-25T03:24:25Z",
      "duration_seconds": 85,
      "active_children_count": 2
    },
    "warning": "2 active child context(s) will remain active. Use --cascade to stop them."
  }
}
```

### Cascade Mode
```json
{
  "command": "stop",
  "timestamp": "2026-01-25T03:26:05Z",
  "data": {
    "data": {
      "context_name": "parent",
      "start_time": "2026-01-25T03:23:00Z",
      "end_time": "2026-01-25T03:26:05Z",
      "duration_seconds": 185,
      "stopped_children": [
        "child2",
        "child1"
      ]
    }
  }
}
```

## Use Cases

### Daily Workflow
```bash
# Start daily context as parent
my-context start "daily-2026-01-25"

# Create task-specific children throughout the day
my-context start "review-PR-123"
my-context link review-PR-123 daily-2026-01-25

my-context start "fix-bug-456"
my-context link fix-bug-456 daily-2026-01-25

# End of day: stop everything at once
my-context resume daily-2026-01-25
my-context stop --cascade
```

### Sprint Planning
```bash
# Sprint as parent
my-context start "sprint-25"

# Stories as children
my-context start "story-USER-123"
my-context link story-USER-123 sprint-25

my-context start "story-ADMIN-456"
my-context link story-ADMIN-456 sprint-25

# End of sprint: close all stories
my-context resume sprint-25
my-context stop --cascade
```

### Feature Development
```bash
# Feature as parent
my-context start "feature-payment-integration"

# Sub-tasks as children
my-context start "api-integration"
my-context link api-integration feature-payment-integration

my-context start "ui-updates"
my-context link ui-updates feature-payment-integration

my-context start "testing"
my-context link testing feature-payment-integration

# Feature complete: stop all related work
my-context resume feature-payment-integration
my-context stop --cascade
```

## Related Commands

### View Hierarchy
```bash
# Show tree structure
my-context tree

# Show tree for specific context
my-context tree parent

# List children
my-context down parent

# Show parent
my-context up child
```

### Manage Links
```bash
# Link child to parent
my-context link <child> <parent>

# Remove parent link
my-context unlink <child>
```

## Tips

1. **Warning First**: Always stop without `--cascade` first to see what will be affected
2. **Review Hierarchy**: Use `my-context tree` to visualize your hierarchy before cascade stop
3. **Check Active Children**: Use `my-context down` to see which children are active
4. **JSON for Automation**: Use `--json` flag for scripting and automation
5. **Normal Operation**: In normal use, only one context is active at a time. The cascade stop feature is most useful when multiple contexts appear active due to manual manipulation or system errors.

## Troubleshooting

### No Children Detected
If children aren't detected, ensure they have the parent link:
```bash
my-context tree     # Check hierarchy
my-context down parent   # Check children
my-context link child parent  # Re-establish link if needed
```

### Wrong Context Stopped
Cascade stop stops contexts in reverse order (leaves first). This ensures data integrity. If you see unexpected contexts being stopped, review your hierarchy with `my-context tree`.

### JSON Output Structure
The JSON output includes:
- `stopped_children`: Array of stopped child names (cascade mode)
- `active_children_count`: Number of active children (warning mode)
- `warning`: Warning message when children remain active

## Implementation Notes

- Works with both file-based and PostgreSQL database backends
- Handles circular dependencies (cycle detection prevents infinite loops)
- Stops descendants in reverse order for data integrity
- Preserves accurate timestamps for all stopped contexts
- No data loss - all context data remains accessible after stop
