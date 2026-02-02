# Global CLAUDE.md Amendment Log

Tracks all amendments to `~/.claude/CLAUDE.md` for audit and review purposes.

## Amendment Registry

| ID | Title | Effective | Status | Review Due |
|----|-------|-----------|--------|------------|
| A-001 | Context Management with my-context | 2026-02-02 | Active | 2026-05-02 |

---

## A-001: Context Management with my-context

**Effective:** 2026-02-02
**Author:** jcaldwell
**Status:** Active

### Summary

Establishes minimum commitments for using `my-context` CLI to track work contexts, capture decisions, and maintain session continuity.

### Key Requirements

1. Check context state at session start
2. Capture decisions, blockers, discoveries as notes
3. Associate modified files with contexts
4. Add summary note before session end
5. Weekly maintenance (stats review, archiving)

### Rationale

- Work continuity across sessions is lost without explicit tracking
- Decision history enables better retrospectives
- Time tracking via stats informs project planning
- Consistent patterns reduce cognitive overhead

### Success Metrics

- [ ] Contexts started for 80%+ of development sessions
- [ ] Average 3+ notes per context
- [ ] Stale contexts addressed within 48h
- [ ] Weekly stats review completed

### Review Notes

_To be updated at quarterly review._

---

## Amendment Process

### Proposing Amendments

1. Draft amendment with: ID, title, purpose, requirements
2. Discuss with user for approval
3. Add to `~/.claude/CLAUDE.md`
4. Log in this file
5. Commit to my-context repo

### Deprecating Amendments

1. Mark status as "Deprecated" in registry
2. Add deprecation note with reason
3. Remove from active CLAUDE.md or mark as deprecated inline

### Quarterly Review

Review all active amendments for:
- Relevance: Is this still useful?
- Compliance: Are commitments being met?
- Evolution: Should requirements change?
