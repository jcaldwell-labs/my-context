# Tasks: CLI Context Management System

**Feature**: 001-cli-context-management  
**Input**: Design documents from `/specs/001-cli-context-management/`  
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow
```
1. Load plan.md from feature directory
   → ✅ Loaded: Go 1.21+, Cobra, plain-text storage
2. Load design documents:
   → ✅ data-model.md: 6 entities (Context, Note, FileAssociation, TouchEvent, ContextTransition, AppState)
   → ✅ contracts/: 8 command contracts
   → ✅ research.md: Technology decisions
   → ✅ quickstart.md: 10 test scenarios
3. Generate tasks by category:
   → Setup: 3 tasks
   → Tests First (TDD): 10 tasks [P]
   → Models: 6 tasks [P]
   → Core Logic: 3 tasks (sequential)
   → Output Formatting: 2 tasks [P]
   → Commands: 8 tasks [P]
   → Integration: 3 tasks
   → Build & Install: 3 tasks
   → Polish: 3 tasks [P]
4. Total: 41 tasks
```

## Task Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- Single project structure: `cmd/`, `internal/`, `tests/` at repository root
- Go standard layout with `internal/` for encapsulation

---

## Phase 3.1: Setup & Project Initialization

- [ ] **T001** Initialize Go module at repository root: `go mod init github.com/yourusername/my-context-copilot`
- [ ] **T002** Install dependencies: `go get github.com/spf13/cobra@latest github.com/spf13/viper@latest github.com/stretchr/testify@latest`
- [ ] **T003** Create project directory structure: `cmd/my-context/`, `internal/{commands,core,models,output}/`, `tests/{integration,unit}/`, `scripts/`

---

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (one per command)
- [ ] **T004** [P] Contract test for `start` command in `tests/integration/start_test.go` - Verify context creation, duplicate name handling, automatic stop of previous context
- [ ] **T005** [P] Contract test for `stop` command in `tests/integration/stop_test.go` - Verify context deactivation, duration calculation, transition logging
- [ ] **T006** [P] Contract test for `note` command in `tests/integration/note_test.go` - Verify note appending with timestamp, special character escaping
- [ ] **T007** [P] Contract test for `file` command in `tests/integration/file_test.go` - Verify file path normalization and association
- [ ] **T008** [P] Contract test for `touch` command in `tests/integration/touch_test.go` - Verify timestamp recording
- [ ] **T009** [P] Contract test for `show` command in `tests/integration/show_test.go` - Verify context details display, JSON output
- [ ] **T010** [P] Contract test for `list` command in `tests/integration/list_test.go` - Verify all contexts listing with sorting
- [ ] **T011** [P] Contract test for `history` command in `tests/integration/history_test.go` - Verify transition log display

### Cross-Platform Integration Tests
- [ ] **T012** [P] Path normalization test in `tests/integration/paths_test.go` - Test Windows (backslash) vs POSIX (forward slash) path handling, relative to absolute conversion
- [ ] **T013** [P] JSON output validation test in `tests/integration/json_test.go` - Verify all commands produce valid JSON with --json flag, test jq parsing

---

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Model Layer (data structures)
- [ ] **T014** [P] Context model in `internal/models/context.go` - Define Context struct with JSON tags, validation methods (name, timestamps, status)
- [ ] **T015** [P] Note model in `internal/models/note.go` - Define Note struct with timestamp and text_content, escaping methods for pipes and newlines
- [ ] **T016** [P] FileAssociation model in `internal/models/file_association.go` - Define FileAssociation struct with timestamp and normalized file_path
- [ ] **T017** [P] TouchEvent model in `internal/models/touch_event.go` - Define TouchEvent struct with timestamp only
- [ ] **T018** [P] ContextTransition model in `internal/models/transition.go` - Define ContextTransition struct with previous_context, new_context, transition_type enum
- [ ] **T019** [P] AppState model in `internal/models/state.go` - Define AppState struct with active_context_name and last_updated

### Core Logic Layer (business operations)
- [ ] **T020** Storage operations in `internal/core/storage.go` - Implement file I/O: CreateDir, ReadJSON, WriteJSON (atomic), AppendLog, ReadLog, NormalizePath (POSIX conversion)
- [ ] **T021** State management in `internal/core/state.go` - Implement: GetActiveContext, SetActiveContext, ClearActiveContext (depends on storage.go)
- [ ] **T022** Context operations in `internal/core/context.go` - Implement: CreateContext (with duplicate name resolution), StopContext, AddNote, AddFile, AddTouch, ListContexts, GetTransitions (depends on storage.go + state.go)

### Output Formatting Layer
- [ ] **T023** [P] Human-readable formatter in `internal/output/human.go` - Implement formatters for all commands: formatContext, formatList, formatHistory with duration calculations
- [ ] **T024** [P] JSON formatter in `internal/output/json.go` - Implement JSON output structure with command, timestamp, data/error fields

### Command Layer (Cobra commands)
- [ ] **T025** [P] Root command in `cmd/my-context/main.go` - Initialize Cobra root command with persistent --json flag, version info, help text
- [ ] **T026** [P] Start command in `internal/commands/start.go` - Implement start subcommand with alias 's', argument validation, call core.CreateContext
- [ ] **T027** [P] Stop command in `internal/commands/stop.go` - Implement stop subcommand with alias 'p', call core.StopContext
- [ ] **T028** [P] Note command in `internal/commands/note.go` - Implement note subcommand with alias 'n', argument validation, call core.AddNote
- [ ] **T029** [P] File command in `internal/commands/file.go` - Implement file subcommand with alias 'f', path validation, call core.AddFile
- [ ] **T030** [P] Touch command in `internal/commands/touch.go` - Implement touch subcommand with alias 't', call core.AddTouch
- [ ] **T031** [P] Show command in `internal/commands/show.go` - Implement show subcommand with alias 'w', read context and format output
- [ ] **T032** [P] List command in `internal/commands/list.go` - Implement list subcommand with alias 'l', call core.ListContexts, sort by start_time
- [ ] **T033** [P] History command in `internal/commands/history.go` - Implement history subcommand with alias 'h', optional --limit flag, call core.GetTransitions

---

## Phase 3.4: Integration & Configuration

- [ ] **T034** Environment configuration in `cmd/my-context/main.go` - Add viper config for MY_CONTEXT_HOME env var with default ~/.my-context, create home directory if not exists
- [ ] **T035** Error handling in `internal/core/errors.go` - Define custom error types with exit codes: UserError (1), SystemError (2), implement stderr output
- [ ] **T036** Help text in all commands - Add usage examples, alias documentation, flag descriptions to all Cobra commands

---

## Phase 3.5: Build, Install & Distribution

- [ ] **T037** Build script in `scripts/build.sh` - Create cross-platform build script: Windows (amd64), Linux (amd64), macOS (amd64, arm64), strip symbols with -ldflags="-s -w"
- [ ] **T038** Install script in `scripts/install.sh` - Create installation script for copying binary to PATH, setting up shell completions
- [ ] **T039** GitHub Actions workflow in `.github/workflows/build.yml` - Setup CI/CD: test matrix (windows, ubuntu, macos), build artifacts, release automation

---

## Phase 3.6: Polish & Documentation

- [ ] **T040** [P] Unit tests in `tests/unit/` - Add unit tests for storage operations, path normalization, duplicate name resolution, escaping functions
- [ ] **T041** [P] Performance validation - Run quickstart.md Scenario 10 (performance benchmarking), verify <10ms command execution, <5MB binary size
- [ ] **T042** [P] Update README.md - Add installation instructions, usage examples, command reference, link to HERE.md and quickstart.md

---

## Dependencies Graph

```
Setup (T001-T003) → Everything else

Tests (T004-T013) → No dependencies (write first, fail initially)

Models (T014-T019) → Independent, can run in parallel
  ↓
Storage (T020) → Uses models
  ↓
State (T021) → Uses storage
  ↓
Context Core (T022) → Uses storage + state + models
  ↓
Output (T023-T024) → Uses models, can run in parallel
  ↓
Commands (T025-T033) → Use core + output, can run in parallel after dependencies
  ↓
Integration (T034-T036) → Uses all commands
  ↓
Build (T037-T039) → Project complete
  ↓
Polish (T040-T042) → Can run in parallel
```

## Parallel Execution Examples

### After Setup Complete (T001-T003)
```bash
# Launch all contract tests together (they'll fail - expected)
go test ./tests/integration/start_test.go
go test ./tests/integration/stop_test.go
go test ./tests/integration/note_test.go
go test ./tests/integration/file_test.go
go test ./tests/integration/touch_test.go
go test ./tests/integration/show_test.go
go test ./tests/integration/list_test.go
go test ./tests/integration/history_test.go
go test ./tests/integration/paths_test.go
go test ./tests/integration/json_test.go
```

### After Tests Written (T004-T013)
```bash
# Create all model files in parallel
# T014: internal/models/context.go
# T015: internal/models/note.go
# T016: internal/models/file_association.go
# T017: internal/models/touch_event.go
# T018: internal/models/transition.go
# T019: internal/models/state.go
```

### After Core Logic Complete (T020-T022)
```bash
# Create output formatters in parallel
# T023: internal/output/human.go
# T024: internal/output/json.go

# Then create all command files in parallel
# T025: cmd/my-context/main.go
# T026: internal/commands/start.go
# T027: internal/commands/stop.go
# T028: internal/commands/note.go
# T029: internal/commands/file.go
# T030: internal/commands/touch.go
# T031: internal/commands/show.go
# T032: internal/commands/list.go
# T033: internal/commands/history.go
```

### Final Polish Phase
```bash
# Run in parallel
# T040: Write unit tests
# T041: Performance validation
# T042: Update README.md
```

---

## Task Execution Notes

**TDD Enforcement**:
- Tests T004-T013 MUST be written first
- Run `go test ./...` to verify tests fail
- Only then proceed to implementation (T014+)
- After each implementation task, re-run tests to see progress

**Path Handling Reminders**:
- Always use `filepath.Clean()` before `filepath.ToSlash()` for storage
- Use `filepath.FromSlash()` for display output
- Test with Windows paths: `C:\Users\...` and POSIX paths: `/home/...`

**Commit Strategy**:
- Commit after each task completion
- Tag releases after T039 (build system complete)
- Format: `git commit -m "T001: Initialize Go module"`

**Validation Checkpoints**:
- After T013: All tests written and failing ✓
- After T022: Core logic complete, tests start passing ✓
- After T033: All commands implemented, all tests passing ✓
- After T039: Binary builds on all platforms ✓
- After T042: Documentation complete, ready for release ✓

---

## Task Generation Rules Applied

1. **From Contracts**: 8 contract files → 8 contract test tasks (T004-T011) marked [P]
2. **From Data Model**: 6 entities → 6 model creation tasks (T014-T019) marked [P]
3. **From Plan**: 3 core modules → 3 sequential tasks (T020-T022) with dependencies
4. **From Commands**: 8 commands + root → 9 command tasks (T025-T033) marked [P]
5. **From Quickstart**: Cross-platform tests → 2 integration test tasks (T012-T013) marked [P]

**Ordering**: Setup → Tests → Models → Core (sequential) → Output + Commands (parallel) → Integration → Build → Polish

**Total Tasks**: 42 tasks
**Parallel Opportunities**: 30 tasks marked [P] (can run concurrently)
**Sequential Dependencies**: 12 tasks (setup, core logic chain, integration)

---

## Success Criteria

All tasks complete when:
- ✅ All 42 tasks checked off
- ✅ `go test ./...` passes with >80% coverage
- ✅ Binary builds for Windows/Linux/macOS
- ✅ All quickstart.md scenarios pass
- ✅ Performance goals met (<10ms, <5MB)
- ✅ README.md and HERE.md updated

**Ready for Implementation**: Execute tasks T001-T042 in dependency order with TDD discipline.

*Based on plan.md v2025-10-04, constitution.md v1.0.0*
