# Clarification: Context Lifecycle Improvements

**Feature**: 004-implement-5-lifecycle
**Phase**: Clarify
**Date**: 2025-10-10

---

## Open Questions from Spec

### Q1: Smart Resume Force Flag Behavior

**Question**: When `--force` is used with duplicate name, should system:
- (A) Automatically create new context with `_2` suffix
- (B) Prompt for new name (current POC behavior)
- (C) Error and require explicit different name

**Decision**: **B - Prompt for new name**

**Rationale**:
- Prevents automatic `_2` suffix proliferation (what we're trying to avoid)
- Forces user to think about meaningful names
- Consistent with smart resume UX (interactive guidance)
- POC validation shows this works well

**Implementation**: In `start` command with `--force`, if duplicate detected, skip resume prompt but still ask for alternative name.

---

### Q2: Note Warnings - Threshold Progression

**Question**: Should warnings be:
- (A) One-time at each threshold (50, 100, 200)
- (B) Every N notes after threshold (every 25 after 200)
- (C) Configurable frequency

**Decision**: **B - Periodic warnings after 200**

**Rationale**:
- One-time warnings get forgotten in long contexts (77+ notes proven viable)
- Periodic gentle reminders keep option visible without being annoying
- POC implements this pattern successfully
- Frequency (every 25) is reasonable (not spam)

**Implementation**:
```go
if noteCount == 50 {
    showWarning(level1)
} else if noteCount == 100 {
    showWarning(level2)
} else if noteCount >= 200 && noteCount % 25 == 0 {
    showWarning(level3Periodic)
}
```

---

### Q3: Resume Backward Compatibility

**Question**: Should `my-context start` still work for resuming existing contexts?
- (A) Yes (soft migration - both work)
- (B) No (force migration to `resume` command)
- (C) Deprecation warning on `start` for existing context

**Decision**: **A - Soft migration, both work**

**Rationale**:
- `resume` is convenience command, not replacement
- Breaking existing scripts/workflows would be disruptive
- Smart resume (FR-MC-NEW-001) makes `start` behavior better anyway
- Users can adopt `resume` at their own pace

**Implementation**:
- `start` command keeps all current behavior + adds smart resume prompt
- `resume` command is explicit alias that only works on stopped contexts
- Both commands use same underlying activation logic

---

### Q4: Bulk Archive Safety Limits

**Question**: Should there be maximum contexts limit for bulk operations?
- (A) No limit (trust confirmation prompt)
- (B) Limit to 100 contexts (require manual for larger)
- (C) Configurable limit via env var

**Decision**: **C - Default limit 100, configurable via MC_BULK_LIMIT**

**Rationale**:
- Safety guard against accidental mass archive (typo in pattern)
- 100 is reasonable upper bound for most usage (Sprint 006 had 16)
- Configurable allows power users to override
- Confirmation prompt is good but limit adds second safety layer

**Implementation**:
```go
limit := getEnvInt("MC_BULK_LIMIT", 100)
if matchCount > limit {
    return fmt.Errorf("Pattern matches %d contexts (limit: %d). Use MC_BULK_LIMIT env var to increase", matchCount, limit)
}
```

**Override example**:
```bash
MC_BULK_LIMIT=500 my-context archive --pattern "old-*"
```

---

### Q5: Lifecycle Advisor - Completion Keyword Dictionary

**Question**: Should completion keywords be:
- (A) Hardcoded list (simple, fast)
- (B) Configurable via file (user customization)
- (C) Pattern-based (regex support)

**Decision**: **A - Hardcoded list**

**Rationale**:
- Sufficient for 95% of use cases (proven in POC)
- Simple implementation, no config file management
- Fast (no file I/O on every stop command)
- Can enhance later if user demand exists

**Keywords to detect** (case-insensitive):
```go
completionKeywords := []string{
    "complete", "completed", "finished", "done",
    "retrospective", "retro", "merged", "deployed",
    "closing", "wrapping", "wrap-up", "ship", "shipped",
}
```

**Detection**: Check last 5 notes for any keyword match (recent notes only, not entire history).

---

## Clarified Requirements

### Smart Resume (FR-001 through FR-005)
- Prompt for new name when `--force` used (not auto `_2`)
- Display note count and last active time in prompt
- Preserve all data on resume (notes, files, transitions)

### Note Warnings (FR-006 through FR-008)
- First warning at 50 notes (default, MC_WARN_AT env var)
- Second warning at 100 notes (MC_WARN_AT_2)
- Periodic warnings every 25 notes after 200 (MC_WARN_AT_3=200)
- All warnings non-blocking, exit code 0

### Resume Command (FR-009 through FR-012)
- New `resume` subcommand
- `--last` flag resumes most recently stopped
- Pattern matching with selection UI on multiple matches
- `start` command still works (backward compat via smart resume)

### Bulk Archive (FR-013 through FR-018)
- Default limit: 100 contexts (MC_BULK_LIMIT env var)
- `--pattern` for glob matching
- `--completed-before` for date filtering
- `--all-stopped` for all stopped contexts
- `--dry-run` for preview
- Confirmation prompt showing count + first 10 contexts

### Lifecycle Advisor (FR-019 through FR-022)
- Hardcoded completion keyword list
- Check last 5 notes only
- Display summary: name, duration, note count
- Detect related contexts (name prefix matching)
- Suggest resume related, archive if complete, or start new

---

## Implementation Decisions

### File Locations

**Go source**:
- `internal/commands/start.go` - Add smart resume logic
- `internal/commands/note.go` - Add warning logic
- `internal/commands/resume.go` - New file for resume command
- `internal/commands/archive.go` - Add bulk operations
- `internal/commands/stop.go` - Add lifecycle advisor output

**Tests**:
- `tests/integration/smart_resume_test.go`
- `tests/integration/note_warnings_test.go`
- `tests/integration/resume_command_test.go`
- `tests/integration/bulk_archive_test.go`
- `tests/integration/lifecycle_advisor_test.go`

### Configuration

**Environment Variables**:
- `MC_WARN_AT` - First warning threshold (default: 50)
- `MC_WARN_AT_2` - Second warning threshold (default: 100)
- `MC_WARN_AT_3` - Periodic warning start (default: 200)
- `MC_BULK_LIMIT` - Max contexts for bulk operations (default: 100)

**No config files needed** - env vars sufficient for customization.

---

## Success Validation

### Sprint 007 Metrics to Capture

**Before Sprint 007** (baseline):
- Context count for Sprint 006: 16
- Active context total: 69
- Archived context total: 1

**During Sprint 007** (with POC scripts):
- Track context count for Sprint 007
- Track `_2`, `_3` suffix occurrences
- Track resume command usage
- Track bulk archive usage

**After Sprint 007** (with Go implementation):
- Measure context fragmentation reduction (target: <4 contexts)
- Measure active context reduction (target: <20)
- Measure adoption of `resume --last` command

### Test Coverage Goals

- Unit tests: 100% coverage for all FR-001 through FR-024
- Integration tests: All 5 user stories + edge cases
- Performance tests: SC-005, SC-006, SC-007 validated

---

## Next Phase

✅ **Clarification Complete**

All open questions answered with decisions and rationale. Ready for planning phase.

**Next**: Create implementation plan with `/plan` command
