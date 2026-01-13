---
allowed-tools: Bash(my-context stop:*)
description: Stop the current work context
---

## Context

- Current active context: !`my-context show --json 2>/dev/null | jq -r '.data.name // "none"'`
- Context duration: !`my-context show --json 2>/dev/null | jq -r '.data.duration // "N/A"'`
- Note count: !`my-context show --json 2>/dev/null | jq -r '.data.note_count // 0'`

## Your task

Stop the currently active work context using `my-context stop`.

After stopping, report:

- Which context was stopped
- How long it was active
- How many notes were captured

If no context is active, inform the user.
