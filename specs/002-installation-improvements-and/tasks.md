# Tasks: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Input**: Design documents from `/specs/002-installation-improvements-and/`  
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow
```
1. Load plan.md from feature directory
   → ✅ Loaded: Go 1.21+, Cobra (existing), multi-platform builds
2. Load design documents:
   → ✅ research.md: 7 key decisions (static builds, user install, project parsing, etc.)
   → ✅ data-model.md: 5 entities (Context modified, ProjectMetadata, ExportDocument, BinaryArtifact, InstallationMetadata)
   → ✅ contracts/: 5 command contracts (export, archive, delete, list-enhanced, project-filter)
   → ✅ quickstart.md: 9 test scenarios
3. Generate tasks by category:
   → Setup: 4 tasks
   → Tests First (TDD): 8 tasks [P]
   → Models: 3 tasks [P]
   → Core Logic: 5 tasks (sequential dependencies)
   → Output Formatting: 2 tasks [P]
   → Commands: 6 tasks [P]
   → Build & Install: 5 tasks
   → Documentation: 3 tasks [P]
   → Bug Fixes: 2 tasks [P]
   → Integration: 3 tasks
4. Total: 41 tasks
```

## Task Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Single project structure (from plan.md):
- Go source: `cmd/my-context/`, `internal/{commands,core,models,output}/`
- Tests: `tests/{integration,unit}/`
- Scripts: `scripts/`
- Docs: `docs/`, root-level `.md` files

---

## Phase 3.1: Setup & Build Infrastructure

- [x] **T001** Create `.github/workflows/release.yml` - GitHub Actions workflow for multi-platform builds (linux/amd64, windows/amd64, darwin/amd64, darwin/arm64), generate SHA256 checksums, upload to releases

- [x] **T002** Create `scripts/build-all.sh` - Shell script to build all 4 platform binaries locally with CGO_ENABLED=0, output to `bin/` directory with platform suffixes

- [x] **T003** [P] Update `cmd/my-context/main.go` - Add version variables (Version, BuildTime, GitCommit) set via ldflags, update root command to display version info

- [x] **T004** [P] Create `docs/TROUBLESHOOTING.md` - Installation troubleshooting guide with platform-specific issues (WSL, Windows PATH, macOS Gatekeeper, permission errors)

---

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (one per command/enhancement)

- [ ] **T005** [P] Create `tests/integration/export_test.go` - Contract tests for export command:
  - Export single context with default path
  - Export with --to custom path
  - Export --all flag (multiple contexts)
  - Export non-existent context (error)
  - Verify markdown format and content
  - Verify JSON output with --json flag

- [ ] **T006** [P] Create `tests/integration/archive_test.go` - Contract tests for archive command:
  - Archive stopped context (success)
  - Archive active context (error)
  - Archive non-existent context (error)
  - Verify is_archived flag in meta.json
  - Verify archived context hidden from default list
  - Verify --archived flag shows archived contexts

- [ ] **T007** [P] Create `tests/integration/delete_test.go` - Contract tests for delete command:
  - Delete with confirmation prompt (accept)
  - Delete with confirmation prompt (cancel)
  - Delete with --force flag (no prompt)
  - Delete active context (error)
  - Verify directory removal
  - Verify transitions.log preserved

- [ ] **T008** [P] Create `tests/integration/list_enhanced_test.go` - Contract tests for list enhancements:
  - Default list (10 most recent)
  - List --all (no limit)
  - List --limit <n> (custom limit)
  - List --search <term> (substring match)
  - List --archived (only archived)
  - List --active-only (only active)
  - Combined filters (project + search + limit)

- [ ] **T009** [P] Create `tests/integration/project_filter_test.go` - Contract tests for project filtering:
  - List --project <name> (filter by project)
  - Start --project <name> (create with project prefix)
  - Project extraction logic (multiple colons, no colon, whitespace)
  - Case-insensitive matching

- [ ] **T010** [P] Create `tests/integration/bug_fixes_test.go` - Contract tests for Sprint 1 bug fixes:
  - Note with $ character (verify preserved)
  - History command (verify "(none)" instead of "NULL")
  - Special characters in notes (!, @, #, etc.)

### Unit Tests

- [ ] **T011** [P] Create `tests/unit/project_parser_test.go` - Unit tests for ExtractProjectName function:
  - Text before first colon
  - No colon (full name)
  - Multiple colons
  - Whitespace trimming
  - Empty string handling

- [ ] **T012** [P] Create `tests/integration/backward_compat_test.go` - Integration test for Sprint 1 → Sprint 2 compatibility:
  - Load Sprint 1 meta.json (without is_archived)
  - Verify Sprint 1 contexts work with Sprint 2 code
  - Verify new features work on old contexts

---

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Model Layer

- [ ] **T013** [P] Modify `internal/models/context.go` - Add IsArchived field to Context struct with JSON tag `"is_archived,omitempty"` for backward compatibility, update validation to prevent archiving active contexts

- [ ] **T014** [P] Create `internal/core/project.go` - Implement ExtractProjectName function (extract text before first colon, trim whitespace, handle edge cases), add FilterByProject function

- [ ] **T015** [P] Create `internal/output/markdown.go` - Implement markdown export formatter:
  - FormatExport function (context → markdown string)
  - Include: header, start/end times, duration, notes with timestamps, files, touch events
  - Local timezone conversion
  - Footer with export timestamp

### Core Logic Layer

- [ ] **T016** Modify `internal/core/context.go` - Add ArchiveContext function (load context, verify stopped, set is_archived=true, save meta.json), add DeleteContext function (verify not active, prompt confirmation, remove directory), add ExportContext function (read all logs, format markdown, write file)

- [ ] **T017** Modify `internal/core/context.go` (depends on T016) - Update ListContexts function to support filters (project, search, limit, archived, active-only), apply AND logic for combined filters, sort by start_time descending

- [ ] **T018** Modify `internal/core/storage.go` - Add WriteMarkdown function (create parent dirs, handle overwrites, atomic write), add SanitizeFilename function (replace special chars for safe filenames)

- [ ] **T019** Modify `internal/output/human.go` - Update formatHistory to display "(none)" instead of "NULL" for empty context fields in transitions

- [ ] **T020** [P] Modify `internal/core/storage.go` - Add ReadNotes, ReadFiles, ReadTouches helper functions to parse log files and return structured data for export

### Command Layer

- [ ] **T021** [P] Create `internal/commands/export.go` - Implement export command with alias 'e':
  - Accept context name argument
  - Support --to <path> flag
  - Support --all flag
  - Call core.ExportContext
  - Handle errors (context not found, file write errors)

- [ ] **T022** [P] Create `internal/commands/archive.go` - Implement archive command with alias 'a':
  - Accept context name argument
  - Call core.ArchiveContext
  - Handle errors (not found, already archived, active context)
  - Display success message

- [ ] **T023** [P] Create `internal/commands/delete.go` - Implement delete command with alias 'd':
  - Accept context name argument
  - Support --force flag
  - Implement confirmation prompt (unless --force)
  - Call core.DeleteContext
  - Handle errors (not found, active context, user cancellation)

- [ ] **T024** [P] Modify `internal/commands/list.go` - Add flags to list command:
  - --project <name> flag
  - --limit <n> flag (default 10)
  - --search <term> flag
  - --all flag
  - --archived flag
  - --active-only flag
  - Validate flag combinations (archived + active-only = error)
  - Display truncation message when limited

- [ ] **T025** [P] Modify `internal/commands/start.go` - Add --project <name> flag:
  - Combine project + name with ": " separator
  - Trim whitespace from both parts
  - Apply existing duplicate name handling
  - Update help text

- [ ] **T026** [P] Modify `internal/commands/note.go` - Fix $ character escaping bug:
  - Review escaping logic in note storage
  - Ensure special characters preserved ($ ! @ # etc.)
  - Update note display to show raw text

---

## Phase 3.4: Build & Installation Scripts

- [x] **T027** Modify `scripts/install.sh` - Enhance Unix installer:
  - Detect existing installation (backup old binary)
  - Install to ~/.local/bin/ (create if needed)
  - Add to PATH in shell rc file (detect bash/zsh)
  - Verify installation successful
  - Preserve ~/.my-context/ data during upgrade

- [ ] **T028** Create `scripts/install.bat` - Windows cmd.exe installer:
  - Detect existing installation
  - Install to %USERPROFILE%\bin\ (create if needed)
  - Add to user PATH via setx command
  - Verify installation
  - Preserve ~/.my-context/ data

- [ ] **T029** Create `scripts/install.ps1` - Windows PowerShell installer:
  - Detect existing installation
  - Install to $env:USERPROFILE\bin\ (create if needed)
  - Add to user PATH via [Environment]::SetEnvironmentVariable
  - Verify installation
  - Preserve ~/.my-context/ data

- [ ] **T030** Create `scripts/curl-install.sh` - One-liner curl installer:
  - Detect platform (Linux, macOS, Windows/WSL via uname)
  - Download appropriate binary from GitHub releases
  - Verify SHA256 checksum
  - Make executable
  - Call install.sh with downloaded binary

- [x] **T031** Update `scripts/build.sh` (if exists) or create wrapper - Add reference to build-all.sh for local multi-platform builds, document GOOS/GOARCH usage

---

## Phase 3.5: Documentation & Polish

- [ ] **T032** [P] Update `README.md` - Add sections:
  - "Building from Source" with platform-specific instructions
  - "Installation" section with curl one-liner and script options
  - Document new commands (export, archive, delete)
  - Document new flags (--project, --limit, --search, etc.)
  - Add troubleshooting link

- [ ] **T033** [P] Update `cmd/my-context/main.go` - Add help text for all new commands and flags, ensure --help displays usage examples

- [ ] **T034** [P] Create `.github/ISSUE_TEMPLATE/` - Add bug report and feature request templates (optional but good practice for Sprint 2 maturity)

---

## Phase 3.6: Integration & Validation

- [ ] **T035** Run all integration tests on Linux - Execute `go test ./tests/integration/... -v` and verify all tests pass, fix any failures

- [ ] **T036** Run all integration tests on Windows (git-bash) - Execute tests in Windows environment, verify path handling, fix platform-specific issues

- [ ] **T037** Manual quickstart validation - Execute all 9 scenarios from quickstart.md:
  1. Multi-platform installation (WSL)
  2. Project-based workflow
  3. Export and share
  4. Context lifecycle (archive/delete)
  5. List enhancements (large dataset)
  6. Bug fixes validation
  7. Cross-platform installation (Windows)
  8. Backward compatibility (Sprint 1 → Sprint 2)
  9. JSON output for scripting

- [ ] **T038** Performance benchmarks - Verify performance targets met:
  - List with 1000 contexts: <1s
  - Export with 500 notes: <1s
  - Search across 1000 contexts: <1s

- [ ] **T039** Binary size verification - Ensure all 4 platform binaries are <10MB each

- [ ] **T040** Cross-platform smoke test - Install on real VMs/machines:
  - Ubuntu 22.04 (WSL and native)
  - Windows 10/11 (cmd.exe and PowerShell)
  - macOS 13+ (Intel and ARM if available)
  - Verify basic commands work identically

- [ ] **T041** Final constitution compliance check - Review all changes against 6 principles:
  - I. Unix Philosophy ✓
  - II. Cross-Platform Compatibility ✓
  - III. Stateful Context Management ✓
  - IV. Minimal Surface Area ✓ (11 commands justified)
  - V. Data Portability ✓
  - VI. User-Driven Design ✓
  - Document any principle tensions in retrospective

---

## Task Dependencies

### Critical Path (Sequential)
```
T001-T004 (Setup)
  ↓
T005-T012 (Tests - must fail) ← BLOCKING GATE
  ↓
T013-T015 (Models - can run [P])
  ↓
T016 (Core context operations)
  ↓
T017 (Depends on T016 - list filtering)
  ↓
T018-T020 (Storage & output - can run [P])
  ↓
T021-T026 (Commands - can run [P])
  ↓
T027-T031 (Build & install - can run [P])
  ↓
T032-T034 (Docs - can run [P])
  ↓
T035-T041 (Integration & validation)
```

### Parallel Execution Opportunities

**Phase 3.1 (Setup)**: T003-T004 can run in parallel with T001-T002

**Phase 3.2 (Tests)**: T005-T012 all run in parallel (8 test files, independent)

**Phase 3.3 (Models)**: T013-T015 can run in parallel (3 different files)

**Phase 3.3 (Storage/Output)**: T018-T020 can run in parallel after T017 completes

**Phase 3.3 (Commands)**: T021-T026 can run in parallel (6 command files)

**Phase 3.4 (Scripts)**: T027-T031 can run in parallel (5 install scripts)

**Phase 3.5 (Docs)**: T032-T034 can run in parallel

**Maximum Parallelism**: Up to 8 tasks can run simultaneously during test-writing phase

---

## Validation Checklist

### Pre-Implementation
- [x] All contracts have corresponding test tasks? YES (T005-T010)
- [x] All entities have model tasks? YES (Context in T013, helpers in T014-T015)
- [x] All commands have implementation tasks? YES (T021-T026)
- [x] TDD enforcement clear? YES (Phase 3.2 MUST complete before 3.3)

### Post-Implementation
- [ ] All tests passing on Linux?
- [ ] All tests passing on Windows?
- [ ] All quickstart scenarios executed successfully?
- [ ] Performance benchmarks met?
- [ ] Binary sizes acceptable?
- [ ] Documentation updated?
- [ ] Backward compatibility verified?

---

## Estimated Effort

**By Phase**:
- Phase 3.1 (Setup): 0.5 days (4 tasks, mostly config)
- Phase 3.2 (Tests): 1.0 days (8 test files, parallel execution)
- Phase 3.3 (Implementation): 2.0 days (17 tasks, mix of parallel and sequential)
- Phase 3.4 (Build/Install): 1.0 days (5 scripts, testing required)
- Phase 3.5 (Docs): 0.5 days (3 tasks, mostly writing)
- Phase 3.6 (Integration): 1.0 days (7 tasks, cross-platform validation)

**Total Estimated**: 6.0 days (with parallelization, closer to 5.5 days actual calendar time)

**Sprint 2 Target**: 5.5 days (from retrospective recommendation) ✅ ON TRACK

---

## Notes

**From research.md**:
- Use CGO_ENABLED=0 for static binaries (no dependencies)
- Install to user-specific directories (no sudo required)
- Project parsing: text before first colon, case-insensitive
- Archive as metadata flag (not directory move)
- Backward compatibility via optional JSON fields

**From data-model.md**:
- IsArchived field optional (omitempty tag)
- ExportDocument format: markdown with headers, lists, timestamps
- Project extraction: strings.SplitN(name, ":", 2)[0]
- Performance target: <1s for list with 1000 contexts

**From contracts/**:
- 5 new/modified command specifications
- Export: markdown generation, --to and --all flags
- Archive: metadata update, hide from default list
- Delete: confirmation prompt, --force flag
- List: 6 new filter flags
- Start: --project flag for naming convention

**From quickstart.md**:
- 9 end-to-end test scenarios covering all features
- Includes performance benchmarks
- Validates backward compatibility explicitly
- Tests JSON output for scripting

---

**Constitution Alignment**:
All tasks designed to maintain compliance with 6 core principles. Sprint 2 adds 3 commands (total 11) which is justified by distinct lifecycle operations (share via export, complete via archive, remove via delete). All features respond to validated user requests from Sprint 1 retrospective.

---

*Tasks generated: 2025-10-05*  
*Based on: plan.md, research.md, data-model.md, contracts/, quickstart.md*  
*Constitution version: 1.1.0*
