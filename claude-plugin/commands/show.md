---
allowed-tools: Bash(my-context show:*)
description: Show the current active context with notes and files
---

## Context

- Active context check: !`my-context show --json 2>/dev/null | jq -r '.status'`

## Your task

Display the current work context using `my-context show`.

Present the information in a clear, readable format including:

- Context name and status
- Duration (how long active)
- Notes captured
- Files associated

If no context is active, inform the user and suggest starting one.
