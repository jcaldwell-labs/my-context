---
allowed-tools: Bash(my-context export:*)
description: Export a context to markdown for documentation
---

## Context

- Active context: !`my-context show --json 2>/dev/null | jq -r '.data.name // "none"'`
- Recent contexts: !`my-context list --limit 3 --json 2>/dev/null | jq -r '.data.contexts[].name'`

## Your task

Export a context to markdown using `my-context export "<context-name>"`.

If no context name provided, export the active context.

Options:

- `--to <path>` - Write to specific file
- Without `--to` - Output to stdout

After exporting, summarize what was exported (notes count, duration, etc.) and where it was saved.
