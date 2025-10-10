# Task List: Context Lifecycle Improvements

**Feature**: 004-implement-5-lifecycle
**Phase**: Tasks
**Date**: 2025-10-10

---

## Phase 1: Smart Resume (P1) - Day 1-2

### FR-001→005: Smart Resume on Start

- [X] **T001**: Add `FindContextByName(name string) (*Context, error)` to `internal/core/state.go`
- [X] **T002**: Add `GetNoteCount(contextName string) (int, error)` to `internal/core/state.go`
- [X] **T003**: Add `GetLastActiveTime(contextName string) (time.Time, error)` to `internal/core/state.go`
- [X] **T004**: Add duplicate detection logic in `internal/commands/start.go` before creating context
- [X] **T005**: Implement interactive prompt function `promptResume(ctx *Context) (bool, error)`
- [X] **T006**: Display context summary (name, note count, last active) in prompt
- [X] **T007**: Handle user response: Y (resume), n (prompt for new name), cancel (abort)
- [X] **T008**: Add `--force` flag support to `start` command
- [X] **T009**: Implement new name prompt when resume declined: `promptNewName(original string) (string, error)`
- [X] **T010**: Add error handling for active context duplicate attempt
- [X] **T011**: Write unit tests for duplicate detection logic
- [X] **T012**: Write integration test: stopped context + start → resume (Y path)
- [X] **T013**: Write integration test: stopped context + start → new name (n path)
- [X] **T014**: Write integration test: `--force` flag behavior
- [X] **T015**: Write integration test: active context error case
- [X] **T016**: Validate against POC `smart-resume.sh` for feature parity
- [X] **T017**: Run all existing tests to ensure no regressions

**Deliverable**: Smart resume working, 17 tasks complete

---

## Phase 2A: Note Warnings (P2) - Day 2

### FR-006→008: Note Limit Warnings

- [X] **T018**: Add threshold check logic in `internal/commands/note.go` (before adding note)
- [X] **T019**: Implement `getEnvInt(key string, defaultVal int) int` helper function
- [X] **T020**: Read threshold env vars: MC_WARN_AT (50), MC_WARN_AT_2 (100), MC_WARN_AT_3 (200)
- [X] **T021**: Implement `showNoteWarning(count int, level int)` with messaging
- [X] **T022**: Add warning at first threshold (50 notes): chunking guidance
- [X] **T023**: Add warning at second threshold (100 notes): getting large message
- [X] **T024**: Add periodic warnings after third threshold (every 25 notes after 200)
- [X] **T025**: Ensure exit code remains 0 (non-blocking warnings)
- [X] **T026**: Write unit test for threshold detection (49→50, 99→100, 200, 225)
- [X] **T027**: Write integration test: warnings at default thresholds
- [X] **T028**: Write integration test: custom thresholds via env vars
- [X] **T029**: Write integration test: exit code validation (must be 0)
- [X] **T030**: Validate warning doesn't interfere with JSON output mode
- [X] **T031**: Run all existing tests for regressions

**Deliverable**: Note warnings working, 14 tasks complete

---

## Phase 2B: Resume Command (P2) - Day 2-3

### FR-009→012: Resume Command Alias

- [X] **T032**: Create `internal/commands/resume.go` with cobra command setup
- [X] **T033**: Add positional argument parsing for context name/pattern
- [X] **T034**: Add `--last` flag definition
- [X] **T035**: Implement `GetMostRecentStopped() (*Context, error)` in `internal/core/state.go`
- [X] **T036**: Implement resume --last logic (activate most recent stopped)
- [X] **T037**: Implement `FindContextsByPattern(pattern string) ([]*Context, error)` in state.go
- [X] **T038**: Implement pattern matching for context names (glob-style)
- [X] **T039**: Implement `PromptSelection(contexts []*Context) (*Context, error)` for multiple matches
- [X] **T040**: Display numbered selection UI when multiple contexts match
- [X] **T041**: Handle single match (activate without prompt)
- [X] **T042**: Handle no matches (error with list of available stopped contexts)
- [X] **T043**: Register `resume` command in `cmd/my-context/main.go`
- [X] **T044**: Write unit test: most recent stopped detection
- [X] **T045**: Write unit test: pattern matching logic
- [X] **T046**: Write integration test: resume specific name
- [X] **T047**: Write integration test: resume --last
- [X] **T048**: Write integration test: pattern with single match (no prompt)
- [X] **T049**: Write integration test: pattern with multiple matches (selection UI)
- [X] **T050**: Write integration test: non-existent context error
- [X] **T051**: Validate against POC `resume-alias.sh` for parity

**Deliverable**: Resume command working, 20 tasks complete

---

## Phase 3: Bulk Archive (P3) - Day 3-4

### FR-013→018: Bulk Archive with Patterns

- [X] **T052**: Add `--pattern` string flag to archive command
- [X] **T053**: Add `--dry-run` boolean flag
- [X] **T054**: Add `--completed-before` string flag (date in YYYY-MM-DD)
- [X] **T055**: Add `--all-stopped` boolean flag
- [X] **T056**: Implement `MatchContextsByPattern(pattern string) ([]*Context, error)`
- [X] **T057**: Implement glob-style pattern matching (support * wildcards)
- [X] **T058**: Implement `FilterByStopDate(contexts []*Context, before time.Time) []*Context`
- [X] **T059**: Parse date string to time.Time, validate format
- [X] **T060**: Implement safety limit check (MC_BULK_LIMIT env var, default 100)
- [X] **T061**: Implement `PromptBulkConfirmation(contexts []*Context) (bool, error)`
- [X] **T062**: Display count and first 10 contexts in confirmation prompt
- [X] **T063**: Implement `BulkArchive(contexts []*Context) (success int, failed int, error)`
- [X] **T064**: Skip active contexts with warning message
- [X] **T065**: Handle partial failures (continue on error, report at end)
- [X] **T066**: Implement dry-run mode (show what would be archived)
- [X] **T067**: Write unit test: pattern matching with wildcards
- [X] **T068**: Write unit test: date parsing and filtering
- [X] **T069**: Write unit test: safety limit enforcement
- [X] **T070**: Write integration test: `--pattern` matching
- [X] **T071**: Write integration test: `--dry-run` preview
- [X] **T072**: Write integration test: `--completed-before` date filtering
- [X] **T073**: Write integration test: `--all-stopped` flag
- [X] **T074**: Write integration test: confirmation prompt (yes path)
- [X] **T075**: Write integration test: confirmation prompt (no path)
- [X] **T076**: Write integration test: safety limit exceeded error
- [X] **T077**: Write integration test: partial failure handling
- [X] **T078**: Validate against POC `bulk-archive.sh` for parity

**Deliverable**: Bulk archive working, 27 tasks complete

---

## Phase 4: Lifecycle Advisor (P3) - Day 4-5

### FR-019→022: Lifecycle Advisor

- [X] **T079**: Add context summary display logic in `internal/commands/stop.go` (after stop)
- [X] **T080**: Display name, duration, note count
- [X] **T081**: Implement `FindRelatedContexts(name string) []*Context` in state.go
- [X] **T082**: Extract context prefix for matching (remove suffixes, phase names)
- [X] **T083**: Find stopped contexts with similar prefixes
- [X] **T084**: Display up to 3 related contexts
- [X] **T085**: Implement `DetectCompletion(notes []Note) bool` - check last 5 notes
- [X] **T086**: Define completion keyword list (hardcoded)
- [X] **T087**: Case-insensitive keyword matching
- [X] **T088**: Implement `DisplayLifecycleGuidance(...)` with suggestions
- [X] **T089**: Suggestion 1: Resume related context (if found)
- [X] **T090**: Suggestion 2: Archive if complete (always)
- [X] **T091**: Suggestion 3: Start new work (always)
- [X] **T092**: Additional suggestion: Archive when completion detected
- [X] **T093**: Write unit test: related context prefix matching
- [X] **T094**: Write unit test: completion keyword detection
- [X] **T095**: Write integration test: summary display
- [X] **T096**: Write integration test: related context detection
- [X] **T097**: Write integration test: completion keyword triggers archive suggestion
- [X] **T098**: Write integration test: no related contexts case
- [X] **T099**: Validate against POC `lifecycle-advisor.sh` for parity

**Deliverable**: Lifecycle advisor working, 21 tasks complete

---

## Phase 5: Integration & Testing - Day 5

### End-to-End Testing

- [X] **T100**: Write end-to-end workflow test combining all 5 features
- [X] **T101**: Test: Create context → add notes with warnings → stop with advisor → resume --last
- [X] **T102**: Test: Smart resume prevents duplicates in real workflow
- [X] **T103**: Test: Bulk archive cleans up multiple related contexts
- [X] **T104**: Validate no `_2` suffix contexts created

### Performance Testing

- [X] **T105**: Write benchmark: note warning overhead (<100ms target)
- [X] **T106**: Write benchmark: bulk archive 100 contexts (<5s target)
- [X] **T107**: Write benchmark: resume --last (<500ms target)
- [X] **T108**: Run all benchmarks, validate against success criteria

### Regression Testing

- [X] **T109**: Run full existing test suite (all Sprint 001-002 tests)
- [X] **T110**: Verify backward compatibility (all existing commands unchanged)
- [X] **T111**: Test non-interactive mode (piped input)
- [X] **T112**: Test JSON output mode with all features

### Documentation

- [X] **T113**: Update README.md with lifecycle improvement features
- [X] **T114**: Add examples for all 5 new features
- [X] **T115**: Update help text for modified commands (start, note, archive, stop)
- [X] **T116**: Document environment variables (MC_WARN_AT_*, MC_BULK_LIMIT)
- [X] **T117**: Create CHANGELOG entry for v2.0.0 lifecycle improvements

### Cleanup

- [X] **T118**: Archive POC scripts to `scripts/poc/archive/`
- [X] **T119**: Update POC README to reference Go implementation
- [X] **T120**: Commit all changes with comprehensive message
- [X] **T121**: Merge feature branch to master after validation

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
