# Tasks: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Input**: Design documents from `/specs/002-installation-improvements-and/`  
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow
```
1. Load plan.md from feature directory
   → ✅ Loaded: Go 1.21+, Cobra (existing), multi-platform builds
2. Load design documents:
   → ✅ research.md: 7 key decisions (static builds, user install, project extraction, etc.)
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
   → Documentation: 5 tasks [P] (expanded post-analysis)
   → Bug Fixes: 2 tasks [P]
   → Integration: 3 tasks
4. Total: 43 tasks (41 original + 2 added for FR-010 coverage)
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

- [x] **T005** [P] Create `tests/integration/export_test.go` - Contract tests for export command:
  - Export single context with default path
  - Export with --to custom path
  - Export --all flag (multiple contexts)
  - Export non-existent context (error)
  - Verify markdown format and content
  - Verify JSON output with --json flag

- [x] **T006** [P] Create `tests/integration/archive_test.go` - Contract tests for archive command:
  - Archive stopped context (success)
  - Archive active context (error)
  - Archive non-existent context (error)
  - Verify is_archived flag in meta.json
  - Verify archived context hidden from default list
  - Verify --archived flag shows archived contexts

- [x] **T007** [P] Create `tests/integration/delete_test.go` - Contract tests for delete command:
  - Delete with confirmation prompt (accept)
  - Delete with confirmation prompt (cancel)
  - Delete with --force flag (no prompt)
  - Delete active context (error)
  - Verify directory removal (context directory deleted from ~/.my-context/)
  - Verify transitions.log preserved (per FR-008.7): create context, transition to another, delete first context, verify transitions.log still contains original transition entries

- [x] **T008** [P] Create `tests/integration/list_enhanced_test.go` - Contract tests for list enhancements:
  - Default list (10 most recent)
  - List --all (no limit)
  - List --limit <n> (custom limit)
  - List --search <term> (substring match)
  - List --archived (only archived)
  - List --active-only (only active)
  - Combined filters (project + search + limit)

- [x] **T009** [P] Create `tests/integration/project_filter_test.go` - Contract tests for project filtering:
  - List --project <name> (filter by project)
  - Start --project <name> (create with project prefix)
  - Project extraction logic (multiple colons, no colon, whitespace)
  - Case-insensitive matching (per FR-004.5): test "PS-CLI" matches "ps-cli: Phase 1" and "Ps-Cli: Phase 2"

- [x] **T010** [P] Create `tests/integration/bug_fixes_test.go` - Contract tests for Sprint 1 bug fixes:
  - Note with $ character (verify preserved)
  - History command (verify "(none)" instead of "NULL")
  - Special characters in notes (!, @, #, etc.)

### Unit Tests

- [x] **T011** [P] Create `tests/unit/project_parser_test.go` - Unit tests for ExtractProjectName function:
  - Text before first colon
  - No colon (full name)
  - Multiple colons
  - Whitespace trimming
  - Empty string handling

- [x] **T012** [P] Create `tests/integration/backward_compat_test.go` - Integration test for Sprint 1 → Sprint 2 compatibility:
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
  - Support --to <path> flag (default: ./<context-name>.md)
  - Support --all flag (export all contexts to separate files)
  - Support --force flag to skip overwrite confirmation (per FR-005.8)
  - Implement overwrite confirmation prompt if output file exists (per FR-005.7): "File exists. Overwrite? (y/N)"
  - Exit code 2 if user declines overwrite
  - Call core.ExportContext
  - Handle errors (context not found, file write errors, user cancellation)

- [ ] **T022** [P] Create `internal/commands/archive.go` - Implement archive command with alias 'a':
  - Accept context name argument
  - Call core.ArchiveContext
  - Handle errors (not found, already archived, active context)
  - Display success message

- [ ] **T023** [P] Create `internal/commands/delete.go` - Implement delete command with alias 'd':
  - Accept context name argument
  - Support --force flag to skip confirmation (per FR-008.4)
  - Implement confirmation prompt (per FR-008.3): "Delete context '<name>' permanently? This cannot be undone. (y/N)"
  - Read user input from stdin, accept 'y' or 'yes' (case-insensitive) to proceed
  - Exit code 1 if user cancels (enters anything other than yes)
  - Prevent deletion if context is active (per FR-008.6): "Cannot delete active context. Stop it first with: my-context stop"
  - Call core.DeleteContext after confirmation
  - Handle errors (not found, active context, user cancellation)
  - Display success message: "Context '<name>' deleted"

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

- [ ] **T028** Create `scripts/install.bat` - Windows cmd.exe installer (per FR-002.2, FR-002.4, FR-002.5):
  - Detect existing installation (check if %USERPROFILE%\bin\my-context.exe exists)
  - If exists, prompt: "Existing installation found. Upgrade? (y/N)" or skip if declined
  - Backup old binary to my-context.exe.bak before upgrade
  - Install to %USERPROFILE%\bin\ (create directory if needed with mkdir)
  - Add to user PATH via setx command (only if not already in PATH)
  - Verify installation by running my-context --version
  - Preserve ~/.my-context/ data directory during upgrade (never delete user data)
  - Display success message with version info

- [ ] **T029** Create `scripts/install.ps1` - Windows PowerShell installer (per FR-002.3, FR-002.4, FR-002.5):
  - Detect existing installation (Test-Path for $env:USERPROFILE\bin\my-context.exe)
  - If exists, prompt: "Existing installation found. Upgrade? (y/N)" or skip if declined
  - Backup old binary to my-context.exe.bak before upgrade
  - Install to $env:USERPROFILE\bin\ (create with New-Item if needed)
  - Add to user PATH via [Environment]::SetEnvironmentVariable('Path', ..., 'User') (only if not already in PATH)
  - Verify installation by running my-context --version
  - Preserve ~/.my-context/ data directory during upgrade (never delete user data)
  - Display success message with version info

- [ ] **T030** Create `scripts/curl-install.sh` - One-liner curl installer (per FR-003.1-003.4):
  - Detect platform (Linux, macOS, Windows/WSL via uname -s and uname -m)
  - Determine correct binary name (my-context-linux-amd64, my-context-darwin-arm64, etc.)
  - Download appropriate binary from GitHub releases (latest or specified version)
  - Download corresponding .sha256 checksum file
  - Verify SHA256 checksum matches (fail installation if mismatch per FR-003.3)
  - Make executable with chmod +x (per FR-003.4)
  - Delegate to install.sh for PATH configuration (satisfies FR-003.4 "add to PATH")
  - Display installation success with version info

- [x] **T031** Update `scripts/build.sh` (if exists) or create wrapper - Add reference to build-all.sh for local multi-platform builds, document GOOS/GOARCH usage

---

## Phase 3.5: Documentation & Polish

- [ ] **T032** [P] Update `README.md` - Add sections per FR-010.1-10.2:
  - "Building from Source" with platform-specific build commands (Go 1.21+, CGO_ENABLED=0)
  - "Installation" section with curl one-liner and script options (install.sh, install.bat, install.ps1)
  - Document new commands (export, archive, delete) with examples
  - Document new flags (--project, --limit, --search, --archived, --active-only, --all, --force)
  - Troubleshooting section for WSL users (link to TROUBLESHOOTING.md)
  - Validate completion against FR-010.1 and FR-010.2 requirements

- [ ] **T032a** [P] Validate `docs/TROUBLESHOOTING.md` content - Verify it includes per FR-010.3:
  - WSL-specific issues (path resolution, binary permissions)
  - Windows PATH configuration problems (cmd.exe vs PowerShell)
  - macOS Gatekeeper warnings ("unidentified developer")
  - Permission errors during installation
  - Common "command not found" scenarios with solutions

- [x] **T032b** [P] Review installation script comments - Verify inline documentation per FR-010.4:
  - scripts/install.sh: Comment each step (detection, backup, PATH modification)
  - scripts/install.bat: Comment Windows-specific logic (setx usage, directory creation)
  - scripts/install.ps1: Comment PowerShell-specific patterns (Environment variables)
  - scripts/curl-install.sh: Comment platform detection and checksum verification

- [x] **T033** [P] Update `cmd/my-context/main.go` - Add help text for all new commands and flags, ensure --help displays usage examples

- [x] **T034** [P] Create `.github/ISSUE_TEMPLATE/` - Add bug report and feature request templates for Sprint 2 maturity

---

## Phase 3.6: Integration & Validation

- [x] **T035** Run all integration tests on Linux - Execute `go test ./tests/integration/... -v` and verify all tests pass, fix any failures

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

- [ ] **T041** Final design principles review - Review all changes against core design principles:
  - Unix Philosophy (text I/O, composability) ✓
  - Cross-Platform Compatibility (Windows, WSL, macOS) ✓
  - Stateful Context Management (one active context) ✓
  - Minimal Surface Area (11 total commands - 3 added this sprint) ✓
  - Data Portability (plain text JSON, markdown exports) ✓
  - User-Driven Design (features from Sprint 1 retrospective) ✓
  - Document any principle tensions or trade-offs in retrospective

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

**Design Principles Alignment**:
All tasks designed to maintain compliance with core design principles. Sprint 2 adds 3 commands (total 11) which is justified by distinct lifecycle operations (share via export, complete via archive, remove via delete). All features respond to validated user requests from Sprint 1 retrospective.

---

*Tasks generated: 2025-10-05*  
*Updated: 2025-10-09 (post-analysis remediation)*  
*Based on: plan.md, research.md, data-model.md, contracts/, quickstart.md*
