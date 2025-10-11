# Implementation Plan: Context Lifecycle Improvements

**Feature**: 004-implement-5-lifecycle
**Phase**: Plan
**Date**: 2025-10-10
**Effort**: 3-5 days implementation + 2 days testing

---

## Overview

Implement 5 validated lifecycle improvements in Go, guided by POC shell scripts in `scripts/poc/`.

**Strategy**: Incremental delivery by priority (P1 → P2 → P3) with TDD approach.

---

## Phase 1: P1 Features - Smart Resume (Day 1-2)

### FR-001 through FR-005: Smart Resume on Start

**Goal**: Prevent context fragmentation (80% reduction target)

**Implementation Steps**:

1. **Modify `internal/commands/start.go`**
   - Add duplicate detection: Check if context name exists in stopped state
   - Add interactive prompt function
   - Display context summary (name, note count, last active time)
   - Handle user response (Y/n/cancel)
   - Support `--force` flag with new name prompt

2. **Update `internal/core/state.go`**
   - Add `FindContextByName(name) *Context` method
   - Add `GetLastActiveTime(contextName) time.Time` method
   - Add `GetNoteCount(contextName) int` method

3. **Add tests**
   - Unit tests: Duplicate detection logic
   - Integration test: `tests/integration/smart_resume_test.go`
     - Test scenario: stopped context + start → prompt → resume
     - Test scenario: stopped context + start → prompt → new name
     - Test scenario: --force flag behavior
     - Test scenario: active context error

**Validation**: Run against POC `smart-resume.sh` behavior for feature parity

**Files Changed**:
- `internal/commands/start.go` (~50 lines added)
- `internal/core/state.go` (~30 lines added)
- `tests/integration/smart_resume_test.go` (new, ~200 lines)

**Deliverable**: Smart resume working, tests passing

---

## Phase 2: P2 Quick Wins - Warnings & Resume Command (Day 2-3)

### FR-006 through FR-008: Note Limit Warnings

**Goal**: Proactive guidance with zero friction

**Implementation Steps**:

1. **Modify `internal/commands/note.go`**
   - Get current note count before adding note
   - Add warning display logic with thresholds
   - Read env vars (MC_WARN_AT, MC_WARN_AT_2, MC_WARN_AT_3)
   - Display appropriate warning based on count

2. **Add helper function**
   - `showNoteWarning(count int, threshold int)` with messaging

3. **Add tests**
   - Unit tests: Threshold detection logic
   - Integration test: `tests/integration/note_warnings_test.go`
     - Test warnings at 50, 100, 200, 225
     - Test custom thresholds via env vars
     - Test exit code remains 0

**Files Changed**:
- `internal/commands/note.go` (~40 lines added)
- `tests/integration/note_warnings_test.go` (new, ~150 lines)

**Effort**: <4 hours (trivial implementation)

---

### FR-009 through FR-012: Resume Command Alias

**Goal**: Semantic clarity + workflow efficiency

**Implementation Steps**:

1. **Create `internal/commands/resume.go`**
   - New cobra command for `resume`
   - Support positional argument (context name or pattern)
   - Support `--last` flag
   - Pattern matching logic
   - Selection UI for multiple matches

2. **Add to `cmd/my-context/main.go`**
   - Register `resume` command

3. **Add helper functions**
   - `GetMostRecentStopped() *Context`
   - `FindContextsByPattern(pattern string) []*Context`
   - `PromptSelection(contexts []*Context) *Context`

4. **Add tests**
   - Unit tests: Pattern matching, last stopped detection
   - Integration test: `tests/integration/resume_command_test.go`
     - Test `resume` specific name
     - Test `resume --last`
     - Test pattern with single match
     - Test pattern with multiple matches (selection UI)

**Files Changed**:
- `internal/commands/resume.go` (new, ~150 lines)
- `cmd/my-context/main.go` (~5 lines added)
- `tests/integration/resume_command_test.go` (new, ~200 lines)

**Effort**: 4-6 hours

**Deliverable**: Day 2-3 complete with warnings + resume command working

---

## Phase 3: P3 Features - Bulk Archive (Day 3-4)

### FR-013 through FR-018: Bulk Archive with Patterns

**Goal**: Efficient cleanup (69 active → <20)

**Implementation Steps**:

1. **Modify `internal/commands/archive.go`**
   - Add `--pattern` flag
   - Add `--dry-run` flag
   - Add `--completed-before` flag
   - Add `--all-stopped` flag
   - Add bulk operation logic with confirmation
   - Implement safety limit (MC_BULK_LIMIT, default 100)

2. **Add helper functions**
   - `MatchContextsByPattern(pattern string) []*Context`
   - `FilterByStopDate(contexts []*Context, before time.Time) []*Context`
   - `PromptBulkConfirmation(contexts []*Context) bool`
   - `BulkArchive(contexts []*Context) (success int, failed int)`

3. **Add tests**
   - Unit tests: Pattern matching, date filtering, limit checking
   - Integration test: `tests/integration/bulk_archive_test.go`
     - Test `--pattern` matching
     - Test `--dry-run` preview
     - Test `--all-stopped`
     - Test confirmation prompt
     - Test safety limit
     - Test partial failure handling

**Files Changed**:
- `internal/commands/archive.go` (~150 lines added)
- `tests/integration/bulk_archive_test.go` (new, ~250 lines)

**Effort**: 8-12 hours

**Deliverable**: Bulk archive working with all flags

---

## Phase 4: P3 Features - Lifecycle Advisor (Day 4-5)

### FR-019 through FR-022: Lifecycle Advisor

**Goal**: Helpful post-stop guidance

**Implementation Steps**:

1. **Modify `internal/commands/stop.go`**
   - After stopping, display context summary
   - Detect related contexts (name prefix matching)
   - Check last 5 notes for completion keywords
   - Display suggestions based on analysis

2. **Add helper functions**
   - `FindRelatedContexts(contextName string) []*Context`
   - `DetectCompletion(notes []Note) bool` - check last 5 for keywords
   - `DisplayLifecycleGuidance(ctx *Context, related []*Context, completion bool)`

3. **Add tests**
   - Unit tests: Related context detection, keyword matching
   - Integration test: `tests/integration/lifecycle_advisor_test.go`
     - Test summary display
     - Test related context detection
     - Test completion keyword detection
     - Test suggestion output

**Files Changed**:
- `internal/commands/stop.go` (~80 lines added)
- `tests/integration/lifecycle_advisor_test.go` (new, ~150 lines)

**Effort**: 4-6 hours

**Deliverable**: All 5 features complete

---

## Phase 5: Integration & Testing (Day 5)

### Integration Testing

**Test all features together**:
1. Create context with smart resume
2. Add notes with warnings
3. Stop with lifecycle advisor
4. Resume with `--last`
5. Bulk archive related contexts

**End-to-end workflow test**:
- Simulate Sprint 007 usage
- Measure context count
- Verify no `_2` suffixes created
- Validate cleanup efficiency

### Performance Testing

**Benchmarks**:
- Note warning overhead: <100ms (SC-005)
- Bulk archive 100 contexts: <5s (SC-006)
- Resume --last: <500ms (SC-007)

**Files**:
- `tests/benchmarks/lifecycle_bench_test.go` (new)

---

## Phase 6: Documentation & Cleanup (Day 5)

### Update Documentation

**Files to update**:
- `README.md` - Add lifecycle improvement features
- `cmd/my-context/main.go` - Update version to reflect features
- `CHANGELOG.md` - Document new features

**POC Migration**:
- Archive POC scripts (move to `scripts/poc/archive/`)
- Update POC README to reference Go implementation
- Keep POCs for historical reference

---

## Technical Approach

### Code Architecture

**Pattern**: Follow existing command structure
- Each feature = enhancement to existing command OR new command file
- Shared logic in `internal/core/`
- Tests mirror source structure

**Interactive Prompts**:
```go
// Use bufio.Reader for user input
reader := bufio.NewReader(os.Stdin)
fmt.Print("Resume? [Y/n]: ")
response, _ := reader.ReadString('\n')
```

**Environment Variables**:
```go
func getEnvInt(key string, defaultVal int) int {
    if val := os.Getenv(key); val != "" {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return defaultVal
}
```

### Data Access

**No schema changes needed** - all features use existing:
- `Context` struct
- `Note` struct
- `state.json`
- `transitions.log`

**New methods in `internal/core/state.go`**:
- `FindContextByName(name string) *Context`
- `GetMostRecentStopped() *Context`
- `FindContextsByPattern(pattern string) []*Context`
- `FindRelatedContexts(name string) []*Context`

---

## Testing Strategy

### TDD Approach

**For each feature**:
1. Write failing integration test
2. Implement minimum code to pass
3. Refactor for clarity
4. Add edge case tests
5. Validate against POC behavior

### Test Organization

```
tests/
├── integration/
│   ├── smart_resume_test.go      # FR-001→005 (200 lines)
│   ├── note_warnings_test.go     # FR-006→008 (150 lines)
│   ├── resume_command_test.go    # FR-009→012 (200 lines)
│   ├── bulk_archive_test.go      # FR-013→018 (250 lines)
│   └── lifecycle_advisor_test.go # FR-019→022 (150 lines)
└── benchmarks/
    └── lifecycle_bench_test.go   # Performance validation
```

**Total test code**: ~1000 lines

---

## Timeline & Milestones

### Day 1: Smart Resume (P1)
- Morning: Implement smart resume logic
- Afternoon: Tests + edge cases
- **Milestone**: Smart resume working, no more `_2` suffixes

### Day 2: Quick Wins (P2)
- Morning: Note warnings implementation
- Early afternoon: Resume command implementation
- Late afternoon: All tests passing
- **Milestone**: P1 + P2 complete (3/5 features)

### Day 3: Bulk Archive (P3)
- Morning: Pattern matching + dry-run
- Afternoon: Date filtering + safety limits
- Evening: Tests + edge cases
- **Milestone**: Bulk archive working

### Day 4: Lifecycle Advisor (P3)
- Morning: Related context detection
- Afternoon: Completion keyword matching + suggestions
- Evening: Tests + integration
- **Milestone**: All 5 features complete

### Day 5: Integration & Polish
- Morning: End-to-end testing
- Afternoon: Performance benchmarks
- Evening: Documentation updates
- **Milestone**: Ready for merge

---

## Dependencies

**None** - All features are additive enhancements to existing commands.

**Build dependencies**:
- Go 1.19+ (existing requirement)
- cobra CLI framework (already in use)
- Standard library only (no new dependencies)

---

## Risks & Mitigation

### Risk 1: Interactive Prompts Break Scripting
**Mitigation**: Ensure all prompts respect pipes/non-TTY environments
```go
if !isatty.IsTerminal(os.Stdin.Fd()) {
    // Non-interactive mode - use defaults
}
```

### Risk 2: Pattern Matching Edge Cases
**Mitigation**: Comprehensive test suite for glob patterns, special characters
- POC already validated basic patterns
- Add tests for edge cases (empty pattern, wildcards, special chars)

### Risk 3: Performance Degradation
**Mitigation**: Benchmark all operations, ensure <100ms overhead
- Note warnings must be fast (checked on every note add)
- Bulk operations acceptable to be slower (run infrequently)

### Risk 4: Backward Compatibility
**Mitigation**: All features are additive
- No changes to existing command behavior
- New flags are opt-in
- Smart resume enhances `start`, doesn't replace it

---

## Next Phase

✅ **Planning Complete**

Implementation plan defined with 5 phases over 5 days.

**Next**: Generate actionable task list with `/tasks` command

---

**Plan Version**: 1.0.0
**Status**: Ready for Task Generation
