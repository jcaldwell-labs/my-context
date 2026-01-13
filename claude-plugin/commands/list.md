---
allowed-tools: Bash(my-context list:*)
description: List all work contexts
---

## Context

- Context count: !`my-context list --json 2>/dev/null | jq '.data.contexts | length'`
- Active context: !`my-context show --json 2>/dev/null | jq -r '.data.name // "none"'`

## Your task

List work contexts using `my-context list`.

Available flags:

- `--limit N` - Limit results (default 10)
- `--search "text"` - Filter by name
- `--project "name"` - Filter by project
- `--archived` - Show only archived
- `--all` - Show all including archived

Present the list in a clear format showing:

- Context names (highlight active)
- Status and duration
- Note/file counts if available

If user requested filtering, apply the appropriate flags.
