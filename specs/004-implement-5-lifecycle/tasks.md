# Task List: Context Lifecycle Improvements

**Feature**: 004-implement-5-lifecycle
**Phase**: Tasks
**Date**: 2025-10-10

---

## Phase 1: Smart Resume (P1) - Day 1-2

### FR-001→005: Smart Resume on Start

- [ ] **T001**: Add `FindContextByName(name string) (*Context, error)` to `internal/core/state.go`
- [ ] **T002**: Add `GetNoteCount(contextName string) (int, error)` to `internal/core/state.go`
- [ ] **T003**: Add `GetLastActiveTime(contextName string) (time.Time, error)` to `internal/core/state.go`
- [ ] **T004**: Add duplicate detection logic in `internal/commands/start.go` before creating context
- [ ] **T005**: Implement interactive prompt function `promptResume(ctx *Context) (bool, error)`
- [ ] **T006**: Display context summary (name, note count, last active) in prompt
- [ ] **T007**: Handle user response: Y (resume), n (prompt for new name), cancel (abort)
- [ ] **T008**: Add `--force` flag support to `start` command
- [ ] **T009**: Implement new name prompt when resume declined: `promptNewName(original string) (string, error)`
- [ ] **T010**: Add error handling for active context duplicate attempt
- [ ] **T011**: Write unit tests for duplicate detection logic
- [ ] **T012**: Write integration test: stopped context + start → resume (Y path)
- [ ] **T013**: Write integration test: stopped context + start → new name (n path)
- [ ] **T014**: Write integration test: `--force` flag behavior
- [ ] **T015**: Write integration test: active context error case
- [ ] **T016**: Validate against POC `smart-resume.sh` for feature parity
- [ ] **T017**: Run all existing tests to ensure no regressions

**Deliverable**: Smart resume working, 17 tasks complete

---

## Phase 2A: Note Warnings (P2) - Day 2

### FR-006→008: Note Limit Warnings

- [ ] **T018**: Add threshold check logic in `internal/commands/note.go` (before adding note)
- [ ] **T019**: Implement `getEnvInt(key string, defaultVal int) int` helper function
- [ ] **T020**: Read threshold env vars: MC_WARN_AT (50), MC_WARN_AT_2 (100), MC_WARN_AT_3 (200)
- [ ] **T021**: Implement `showNoteWarning(count int, level int)` with messaging
- [ ] **T022**: Add warning at first threshold (50 notes): chunking guidance
- [ ] **T023**: Add warning at second threshold (100 notes): getting large message
- [ ] **T024**: Add periodic warnings after third threshold (every 25 notes after 200)
- [ ] **T025**: Ensure exit code remains 0 (non-blocking warnings)
- [ ] **T026**: Write unit test for threshold detection (49→50, 99→100, 200, 225)
- [ ] **T027**: Write integration test: warnings at default thresholds
- [ ] **T028**: Write integration test: custom thresholds via env vars
- [ ] **T029**: Write integration test: exit code validation (must be 0)
- [ ] **T030**: Validate warning doesn't interfere with JSON output mode
- [ ] **T031**: Run all existing tests for regressions

**Deliverable**: Note warnings working, 14 tasks complete

---

## Phase 2B: Resume Command (P2) - Day 2-3

### FR-009→012: Resume Command Alias

- [ ] **T032**: Create `internal/commands/resume.go` with cobra command setup
- [ ] **T033**: Add positional argument parsing for context name/pattern
- [ ] **T034**: Add `--last` flag definition
- [ ] **T035**: Implement `GetMostRecentStopped() (*Context, error)` in `internal/core/state.go`
- [ ] **T036**: Implement resume --last logic (activate most recent stopped)
- [ ] **T037**: Implement `FindContextsByPattern(pattern string) ([]*Context, error)` in state.go
- [ ] **T038**: Implement pattern matching for context names (glob-style)
- [ ] **T039**: Implement `PromptSelection(contexts []*Context) (*Context, error)` for multiple matches
- [ ] **T040**: Display numbered selection UI when multiple contexts match
- [ ] **T041**: Handle single match (activate without prompt)
- [ ] **T042**: Handle no matches (error with list of available stopped contexts)
- [ ] **T043**: Register `resume` command in `cmd/my-context/main.go`
- [ ] **T044**: Write unit test: most recent stopped detection
- [ ] **T045**: Write unit test: pattern matching logic
- [ ] **T046**: Write integration test: resume specific name
- [ ] **T047**: Write integration test: resume --last
- [ ] **T048**: Write integration test: pattern with single match (no prompt)
- [ ] **T049**: Write integration test: pattern with multiple matches (selection UI)
- [ ] **T050**: Write integration test: non-existent context error
- [ ] **T051**: Validate against POC `resume-alias.sh` for parity

**Deliverable**: Resume command working, 20 tasks complete

---

## Phase 3: Bulk Archive (P3) - Day 3-4

### FR-013→018: Bulk Archive with Patterns

- [ ] **T052**: Add `--pattern` string flag to archive command
- [ ] **T053**: Add `--dry-run` boolean flag
- [ ] **T054**: Add `--completed-before` string flag (date in YYYY-MM-DD)
- [ ] **T055**: Add `--all-stopped` boolean flag
- [ ] **T056**: Implement `MatchContextsByPattern(pattern string) ([]*Context, error)`
- [ ] **T057**: Implement glob-style pattern matching (support * wildcards)
- [ ] **T058**: Implement `FilterByStopDate(contexts []*Context, before time.Time) []*Context`
- [ ] **T059**: Parse date string to time.Time, validate format
- [ ] **T060**: Implement safety limit check (MC_BULK_LIMIT env var, default 100)
- [ ] **T061**: Implement `PromptBulkConfirmation(contexts []*Context) (bool, error)`
- [ ] **T062**: Display count and first 10 contexts in confirmation prompt
- [ ] **T063**: Implement `BulkArchive(contexts []*Context) (success int, failed int, error)`
- [ ] **T064**: Skip active contexts with warning message
- [ ] **T065**: Handle partial failures (continue on error, report at end)
- [ ] **T066**: Implement dry-run mode (show what would be archived)
- [ ] **T067**: Write unit test: pattern matching with wildcards
- [ ] **T068**: Write unit test: date parsing and filtering
- [ ] **T069**: Write unit test: safety limit enforcement
- [ ] **T070**: Write integration test: `--pattern` matching
- [ ] **T071**: Write integration test: `--dry-run` preview
- [ ] **T072**: Write integration test: `--completed-before` date filtering
- [ ] **T073**: Write integration test: `--all-stopped` flag
- [ ] **T074**: Write integration test: confirmation prompt (yes path)
- [ ] **T075**: Write integration test: confirmation prompt (no path)
- [ ] **T076**: Write integration test: safety limit exceeded error
- [ ] **T077**: Write integration test: partial failure handling
- [ ] **T078**: Validate against POC `bulk-archive.sh` for parity

**Deliverable**: Bulk archive working, 27 tasks complete

---

## Phase 4: Lifecycle Advisor (P3) - Day 4-5

### FR-019→022: Lifecycle Advisor

- [ ] **T079**: Add context summary display logic in `internal/commands/stop.go` (after stop)
- [ ] **T080**: Display name, duration, note count
- [ ] **T081**: Implement `FindRelatedContexts(name string) []*Context` in state.go
- [ ] **T082**: Extract context prefix for matching (remove suffixes, phase names)
- [ ] **T083**: Find stopped contexts with similar prefixes
- [ ] **T084**: Display up to 3 related contexts
- [ ] **T085**: Implement `DetectCompletion(notes []Note) bool` - check last 5 notes
- [ ] **T086**: Define completion keyword list (hardcoded)
- [ ] **T087**: Case-insensitive keyword matching
- [ ] **T088**: Implement `DisplayLifecycleGuidance(...)` with suggestions
- [ ] **T089**: Suggestion 1: Resume related context (if found)
- [ ] **T090**: Suggestion 2: Archive if complete (always)
- [ ] **T091**: Suggestion 3: Start new work (always)
- [ ] **T092**: Additional suggestion: Archive when completion detected
- [ ] **T093**: Write unit test: related context prefix matching
- [ ] **T094**: Write unit test: completion keyword detection
- [ ] **T095**: Write integration test: summary display
- [ ] **T096**: Write integration test: related context detection
- [ ] **T097**: Write integration test: completion keyword triggers archive suggestion
- [ ] **T098**: Write integration test: no related contexts case
- [ ] **T099**: Validate against POC `lifecycle-advisor.sh` for parity

**Deliverable**: Lifecycle advisor working, 21 tasks complete

---

## Phase 5: Integration & Testing - Day 5

### End-to-End Testing

- [ ] **T100**: Write end-to-end workflow test combining all 5 features
- [ ] **T101**: Test: Create context → add notes with warnings → stop with advisor → resume --last
- [ ] **T102**: Test: Smart resume prevents duplicates in real workflow
- [ ] **T103**: Test: Bulk archive cleans up multiple related contexts
- [ ] **T104**: Validate no `_2` suffix contexts created

### Performance Testing

- [ ] **T105**: Write benchmark: note warning overhead (<100ms target)
- [ ] **T106**: Write benchmark: bulk archive 100 contexts (<5s target)
- [ ] **T107**: Write benchmark: resume --last (<500ms target)
- [ ] **T108**: Run all benchmarks, validate against success criteria

### Regression Testing

- [ ] **T109**: Run full existing test suite (all Sprint 001-002 tests)
- [ ] **T110**: Verify backward compatibility (all existing commands unchanged)
- [ ] **T111**: Test non-interactive mode (piped input)
- [ ] **T112**: Test JSON output mode with all features

### Documentation

- [ ] **T113**: Update README.md with lifecycle improvement features
- [ ] **T114**: Add examples for all 5 new features
- [ ] **T115**: Update help text for modified commands (start, note, archive, stop)
- [ ] **T116**: Document environment variables (MC_WARN_AT_*, MC_BULK_LIMIT)
- [ ] **T117**: Create CHANGELOG entry for v2.0.0 lifecycle improvements

### Cleanup

- [ ] **T118**: Archive POC scripts to `scripts/poc/archive/`
- [ ] **T119**: Update POC README to reference Go implementation
- [ ] **T120**: Commit all changes with comprehensive message
- [ ] **T121**: Merge feature branch to master after validation

**Deliverable**: All 5 features complete, tested, documented, 22 tasks

---

## Task Summary

**Total Tasks**: 121
- Phase 1 (Smart Resume): 17 tasks
- Phase 2A (Note Warnings): 14 tasks
- Phase 2B (Resume Command): 20 tasks
- Phase 3 (Bulk Archive): 27 tasks
- Phase 4 (Lifecycle Advisor): 21 tasks
- Phase 5 (Integration & Testing): 22 tasks

**Estimated Effort**:
- Implementation: 89 tasks (Days 1-4)
- Testing: 21 tasks (Day 5, interspersed)
- Documentation: 11 tasks (Day 5)

---

## Execution Checkpoints

### Checkpoint 1: Smart Resume Complete (Day 1-2)
**Verify**:
- [ ] Smart resume prompts on duplicate names
- [ ] Resume preserves all notes and files
- [ ] `--force` flag works correctly
- [ ] All tests passing (integration + unit)

### Checkpoint 2: Quick Wins Complete (Day 2-3)
**Verify**:
- [ ] Note warnings appear at correct thresholds
- [ ] Resume command works with --last and patterns
- [ ] All P1 + P2 features working together
- [ ] No performance degradation

### Checkpoint 3: Bulk Archive Complete (Day 3-4)
**Verify**:
- [ ] Pattern matching works correctly
- [ ] Dry-run previews accurately
- [ ] Safety limits enforced
- [ ] Partial failures handled gracefully

### Checkpoint 4: All Features Complete (Day 4-5)
**Verify**:
- [ ] Lifecycle advisor provides helpful suggestions
- [ ] All 5 features work together
- [ ] Integration tests passing
- [ ] Performance benchmarks met

### Checkpoint 5: Ready for Merge (Day 5)
**Verify**:
- [ ] All 121 tasks complete
- [ ] Documentation updated
- [ ] No regressions in existing functionality
- [ ] Success criteria SC-001 through SC-010 validated

---

## Next Phase

✅ **Tasks Generated**

121 actionable tasks organized into 5 phases over 5 days.

**Next**: Begin implementation with Phase 1 (Smart Resume)

---

**Tasks Version**: 1.0.0
**Status**: Ready for Implementation
**Execution**: Use `/implement` command to begin
