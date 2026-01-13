---
allowed-tools: Bash(my-context start:*)
description: Start a new work context for tracking this session
---

## Context

- Current active context: !`my-context show --json 2>/dev/null | jq -r '.data.name // "none"'`
- Recent contexts: !`my-context list --limit 5 --json 2>/dev/null | jq -r '.data.contexts[].name' | head -5`

## Your task

Start a new work context using `my-context start "<context-name>"`.

If the user provided a context name, use it. Otherwise, suggest a descriptive name based on the current task or conversation topic.

The context name should be:

- Descriptive of the work being done
- Can use "project: phase" format for organization (e.g., "my-context: MCP implementation")
- Avoid generic names like "test" or "work"

After starting, confirm what context was created.
