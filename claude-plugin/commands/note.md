---
allowed-tools: Bash(my-context note:*)
description: Add a timestamped note to the current context
---

## Context

- Current active context: !`my-context show --json 2>/dev/null | jq -r '.data.name // "none"'`
- Recent notes: !`my-context show --json 2>/dev/null | jq -r '.data.notes[-3:][]? // "none"'`

## Your task

Add a note to the current context using `my-context note "<note-text>"`.

If the user provided specific text, use it. Otherwise, summarize the recent work or decision made in this conversation.

Good notes should:

- Be concise but informative
- Capture decisions, discoveries, or key actions
- Be useful for future reference

If no context is active, suggest starting one first.
