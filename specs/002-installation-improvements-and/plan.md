# Implementation Plan: Installation & Usability Improvements

**Branch**: `002-installation-improvements-and` | **Date**: 2025-10-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-installation-improvements-and/spec.md`

## Summary

Sprint 2 addresses critical installation blockers and delivers high-value user-requested features identified in the Sprint 1 retrospective. Primary focus: multi-platform binary distribution (Windows, Linux, macOS including ARM), installation scripts for all environments, project filtering based on organic user naming conventions ("project: phase - description"), markdown export for sharing contexts, list command enhancements (pagination, search, filters), and archive/delete commands for context lifecycle management. Also fixes bugs identified in Sprint 1 ($ character in notes, NULL display in history).

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: 
- github.com/spf13/cobra (CLI framework, existing)
- github.com/spf13/viper (config, existing)
- No new external dependencies required

**Storage**: Plain-text files in `~/.my-context/` (existing structure, extended with archive flag in meta.json)  
**Testing**: Go testing package with integration tests (existing framework)  
**Target Platform**: Windows (amd64), Linux (amd64), macOS (amd64, arm64) - all via single Go codebase  
**Project Type**: Single (CLI tool with `cmd/`, `internal/`, `tests/` structure)  

**Performance Goals**: 
- Build time: <30 seconds per platform
- Binary size: <10MB per platform
- List command: Handle 1000+ contexts with <1s response time
- Export: Generate markdown file in <1s for typical context (50 notes)

**Constraints**: 
- Must maintain backward compatibility with Sprint 1 data structures
- Cannot break existing CLI interface (only additions allowed)
- Installation scripts must work on clean systems (no Go required)
- Archive flag must be optional in meta.json (existing files lack it)

**Scale/Scope**: 
- 4 platform builds per release
- 3 new commands (export, archive, delete)
- 4 new flags on existing commands (--project, --limit, --search, --all, --archived, --active-only, --force)
- 3 installation scripts (install.sh, install.bat, install.ps1)
- 1 curl-based installer

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Unix Philosophy** (composability, text I/O, single-purpose commands)
- [x] Each component does one thing well - Export creates markdown, archive marks status, delete removes data
- [x] Commands accept text input and produce text output (stdin/stdout/stderr) - All new commands follow existing pattern
- [x] No unnecessary coupling between components - Export/archive/delete are independent operations
- [x] Can be chained with standard shell tools (pipes, grep, etc.) - Export output is markdown (greppable), list filters work with standard tools

**II. Cross-Platform Compatibility** (Windows, git-bash, WSL)
- [x] Works identically across cmd.exe, PowerShell, git-bash, and WSL - Multi-platform binaries ensure this
- [x] Path handling normalizes Windows backslashes and POSIX forward slashes - Existing normalization in storage.go applies to new commands
- [x] Single executable with environment detection - Go's GOOS/GOARCH handles this automatically
- [x] Installation separates dev folder from runtime home directory - Installation scripts place binary in PATH, data stays in `~/.my-context/`

**III. Stateful Context Management** (one active context, automatic transitions)
- [x] Only one active context at a time - No changes to this model
- [x] Starting new context automatically stops previous - No changes to this behavior
- [x] State persists in dedicated home directory - Archive flag added to existing meta.json
- [x] All operations default to active context - Export/archive/delete operate on named contexts (explicit selection)

**IV. Minimal Surface Area** (<=10 commands, single-letter aliases, zero config)
- [x] Feature adds minimal new commands/flags (justify if >1 new command) - Adding 3 commands (export, archive, delete) brings total to 11, but all are essential lifecycle operations
- [x] Each command has single-letter alias - export (e), archive (a), delete (d)
- [x] Defaults are sensible (no config required) - List shows 10 by default, export outputs to current dir, delete requires confirmation
- [x] Built-in help available - All new commands will have --help text

**Justification for 3 new commands**: Export, archive, and delete represent distinct lifecycle stages (share, complete, remove) that cannot be overloaded onto existing commands without violating Unix philosophy. Total command count (11) remains reasonable for tool complexity.

**V. Data Portability** (plain text, no lock-in, greppable)
- [x] All data stored as plain text (JSON or newline-delimited logs) - Archive adds boolean to meta.json, export generates markdown
- [x] Can export/import with standard file operations (cp, cat, grep) - Export explicitly enables this, archive flag is readable JSON
- [x] No proprietary formats - Markdown export is universally readable
- [x] Data directory is self-contained and relocatable - No changes to `~/.my-context/` structure

**VI. User-Driven Design** (observe patterns, automate workflows, validate requests)
- [x] Feature responds to observed user behavior (not speculative) - Project filter responds to "project: phase" naming convention observed in Sprint 1
- [x] Automates existing manual workflow (if applicable) - Export automates manual copy to contexts/ directory, project filter formalizes manual parsing
- [x] User request validated through retrospective or feedback (if applicable) - All features from Sprint 1 retrospective P1 user requests
- [x] Adapts to user conventions rather than enforcing rigid structure - Project filter parses existing names, doesn't require schema change

**Violations (if any)**:
- None - All principles satisfied. The addition of 3 commands (bringing total to 11) is justified by distinct lifecycle operations and remains within reasonable cognitive load.

## Project Structure

### Documentation (this feature)
```
specs/002-installation-improvements-and/
├── plan.md              # This file
├── spec.md              # Feature specification (exists)
├── research.md          # Phase 0 output (to be generated)
├── data-model.md        # Phase 1 output (to be generated)
├── quickstart.md        # Phase 1 output (to be generated)
├── contracts/           # Phase 1 output (to be generated)
│   ├── export.md
│   ├── archive.md
│   ├── delete.md
│   ├── list-enhanced.md
│   └── project-filter.md
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
cmd/
└── my-context/
    └── main.go          # Add version info for binaries

internal/
├── commands/
│   ├── export.go        # NEW: Export command implementation
│   ├── archive.go       # NEW: Archive command implementation
│   ├── delete.go        # NEW: Delete command implementation
│   ├── list.go          # MODIFY: Add --project, --limit, --search, --all, --archived, --active-only flags
│   ├── start.go         # MODIFY: Add --project flag
│   ├── note.go          # MODIFY: Fix $ character escaping bug
│   └── history.go       # MODIFY: Replace NULL with "(none)"
├── core/
│   ├── context.go       # MODIFY: Add archive status handling, project extraction, delete operation
│   └── storage.go       # MODIFY: Add markdown export formatter
├── models/
│   └── context.go       # MODIFY: Add IsArchived field to Context struct
└── output/
    ├── human.go         # MODIFY: Update NULL display logic
    └── markdown.go      # NEW: Markdown export formatter

tests/
├── integration/
│   ├── export_test.go   # NEW: Export command tests
│   ├── archive_test.go  # NEW: Archive command tests
│   ├── delete_test.go   # NEW: Delete command tests
│   ├── list_enhanced_test.go  # NEW: List enhancements tests
│   ├── project_filter_test.go # NEW: Project filtering tests
│   └── bug_fixes_test.go      # NEW: Bug fix validation tests
└── unit/
    └── project_parser_test.go  # NEW: Unit tests for project name extraction

scripts/
├── build-all.sh         # NEW: Multi-platform build script
├── install.sh           # MODIFY: Enhanced with upgrade detection
├── install.bat          # NEW: Windows cmd.exe installer
├── install.ps1          # NEW: Windows PowerShell installer
└── curl-install.sh      # NEW: One-liner curl installer

.github/
└── workflows/
    └── release.yml      # NEW: GitHub Actions workflow for multi-platform builds

docs/
└── TROUBLESHOOTING.md   # NEW: Installation troubleshooting guide
```

**Structure Decision**: Single project structure maintained. All new functionality integrates into existing `internal/` packages following established patterns. Build scripts added at repository root under `scripts/`. GitHub Actions workflow added for automated releases.

## Phase 0: Outline & Research

### Research Tasks

1. **Multi-platform Go builds**
   - Research: Go cross-compilation best practices (CGO_ENABLED=0 for static linking)
   - Decision needed: Build automation (Makefile vs shell scripts vs GitHub Actions)
   - Investigate: Binary signing and notarization for macOS (future consideration)

2. **Installation patterns for CLI tools**
   - Research: Standard installation locations per platform (Unix: /usr/local/bin, Windows: user-specific PATH)
   - Decision needed: Upgrade strategy (in-place replacement vs backup-and-replace)
   - Investigate: PATH modification across different shells (bash profile, zsh rc, PowerShell profile)

3. **Project name parsing strategy**
   - Research: Edge cases in "project: phase" pattern (multiple colons, no colons, special characters)
   - Decision needed: Case sensitivity for project filtering
   - Investigate: Unicode handling in project names

4. **Markdown export format**
   - Research: Standard markdown conventions for structured data (tables vs lists vs headers)
   - Decision needed: Include metadata (timestamps, duration) vs content-only export
   - Investigate: Compatibility with popular markdown viewers (GitHub, VS Code, Obsidian)

5. **Archive vs Delete semantics**
   - Research: Industry patterns (soft delete, trash, archive folders)
   - Decision needed: Archive as metadata flag vs moving to subdirectory
   - Investigate: Recovery mechanisms (unarchive command for Sprint 3?)

6. **List pagination approaches**
   - Research: CLI pagination patterns (limit/offset vs cursor-based)
   - Decision needed: Sort order when limit applied (newest first vs oldest first)
   - Investigate: Performance of filtering 1000+ contexts

7. **Backward compatibility testing**
   - Research: Strategies for testing old data with new code
   - Decision needed: Migration script vs graceful degradation
   - Investigate: Semantic versioning for data format changes

**Output**: research.md with all decisions documented

## Phase 1: Design & Contracts

### Data Model Changes

**Context Entity** (existing, modified):
```
{
  "name": string,
  "start_time": timestamp,
  "end_time": timestamp (optional),
  "status": "active" | "stopped",
  "subdirectory_path": string,
  "is_archived": boolean (NEW, optional for backward compat)
}
```

**Project Metadata** (derived, not stored):
- project_name: extracted from context name before first colon
- Extraction logic: `strings.Split(contextName, ":")[0]` with trim

**Export Document** (output format):
```markdown
# Context: {name}

**Started**: {start_time}
**Ended**: {end_time or "Active"}
**Duration**: {calculated duration}

## Notes

- [{timestamp}] {note_text}
...

## Associated Files

- {file_path} (added {timestamp})
...

## Activity

- {touch_count} touch events
```

### API Contracts

**Export Command**:
```
Input: my-context export <context-name> [--to <path>] [--all]
Output: Markdown file(s) with context data
Exit codes: 0 (success), 1 (context not found), 2 (file write error)
```

**Archive Command**:
```
Input: my-context archive <context-name>
Output: "Archived context: {name}"
Side effect: Sets is_archived=true in meta.json
Exit codes: 0 (success), 1 (context not found), 2 (already archived), 3 (context is active)
```

**Delete Command**:
```
Input: my-context delete <context-name> [--force]
Output: Confirmation prompt (unless --force), then "Deleted context: {name}"
Side effect: Removes entire context directory
Exit codes: 0 (success), 1 (context not found), 2 (deletion cancelled), 3 (context is active)
```

**List Command Enhancements**:
```
Input: my-context list [--project <name>] [--limit <n>] [--search <term>] [--all] [--archived] [--active-only]
Output: Filtered list of contexts
Exit codes: 0 (success, even if no matches)
```

**Project Filter (Start Command)**:
```
Input: my-context start "<phase>" --project <project-name>
Output: Creates context named "{project}: {phase}"
Exit codes: 0 (success)
```

### Contract Tests (Phase 1)

Generate failing tests in `tests/integration/`:
- export_test.go: Verify markdown generation, file creation, --all flag
- archive_test.go: Verify meta.json update, list filtering, active context rejection
- delete_test.go: Verify directory removal, confirmation prompt, active context rejection
- list_enhanced_test.go: Verify --limit, --search, --all flags
- project_filter_test.go: Verify project extraction, filtering, start with --project
- bug_fixes_test.go: Verify $ character preservation, NULL replacement

### Quickstart Scenarios

1. **Multi-platform installation**:
   ```bash
   # Linux/WSL
   curl -sSL https://raw.githubusercontent.com/.../install.sh | bash
   my-context --version
   
   # Windows (PowerShell)
   Invoke-WebRequest -Uri "..." -OutFile install.ps1
   .\install.ps1
   my-context --version
   ```

2. **Project-based workflow**:
   ```bash
   my-context start "Phase 1 - Foundation" --project ps-cli
   my-context note "PLAN: Initial setup"
   my-context list --project ps-cli
   # Shows only ps-cli contexts
   ```

3. **Export and share**:
   ```bash
   my-context export "ps-cli: Phase 1" --to ./docs/phase-1-summary.md
   cat docs/phase-1-summary.md
   # Human-readable markdown
   ```

4. **Context lifecycle**:
   ```bash
   my-context archive "ps-cli: Phase 1"
   my-context list
   # Phase 1 hidden
   my-context list --archived
   # Phase 1 visible
   my-context delete "Test Context"
   # Confirm: y
   # Context removed permanently
   ```

### Agent Context Update

Since we're using Claude Code, update `CLAUDE.md` incrementally:
- Add: Multi-platform build commands (go build with GOOS/GOARCH)
- Add: New commands (export, archive, delete) with brief descriptions
- Add: Project filtering pattern ("project: phase" convention)
- Update: Recent changes section with Sprint 2 additions
- Keep under 150 lines for token efficiency

**Output**: data-model.md, contracts/, quickstart.md, failing tests, updated CLAUDE.md

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
1. Load `.specify/templates/tasks-template.md` as base
2. Generate tasks from Phase 1 contracts:
   - Each contract → contract test task [P]
   - Each new command → implementation task
   - Each modified command → enhancement task
   - Build/install infrastructure → setup tasks
3. Order by TDD principles:
   - Phase 3.1: Setup (build scripts, GitHub Actions)
   - Phase 3.2: Tests First (all contract tests must fail)
   - Phase 3.3: Core Implementation (make tests pass)
   - Phase 3.4: Integration (installation scripts)
   - Phase 3.5: Documentation

**Ordering Strategy**:
- Tests before implementation (TDD mandatory)
- Independent tasks marked [P] for parallel execution:
  - All contract tests can run in parallel
  - Export, archive, delete implementations are independent
  - Build scripts for different platforms are independent
- Sequential tasks:
  - Model changes before command implementations
  - Core functions before commands that use them
  - Bug fixes before integration tests

**Estimated Output**: ~35-40 tasks
- Setup: 3 tasks (build scripts, CI/CD, tools)
- Tests: 8 tasks (6 integration test files + 2 unit tests)
- Models: 2 tasks (Context model update, project extraction utility)
- Core: 4 tasks (archive/delete/export logic, markdown formatter)
- Commands: 6 tasks (3 new commands + 3 modified commands)
- Installation: 4 tasks (3 scripts + 1 curl installer)
- Documentation: 3 tasks (README, TROUBLESHOOTING, inline help)
- Bug fixes: 2 tasks ($ character, NULL display)
- Integration: 3 tasks (backward compat tests, upgrade tests, multi-platform validation)

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, cross-platform validation)

## Complexity Tracking
*No violations - Constitution Check passed*

No additional complexity beyond Sprint 1. The 3 new commands (export, archive, delete) represent natural lifecycle operations rather than scope creep. All features directly address Sprint 1 retrospective findings and user requests.

## Progress Tracking

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS (all 6 principles satisfied)
- [x] Post-Design Constitution Check: PASS (re-verified after Phase 1)
- [x] All NEEDS CLARIFICATION resolved (Sprint 1 retrospective provided all answers)
- [x] Complexity deviations documented (none)

**Artifacts Generated**:
- [x] plan.md (this file)
- [x] research.md (7 key decisions documented)
- [x] data-model.md (5 entities defined)
- [x] contracts/ directory (5 contract files)
  - [x] export.md
  - [x] archive.md
  - [x] delete.md
  - [x] list-enhanced.md
  - [x] project-filter.md
- [x] quickstart.md (9 end-to-end scenarios)
- [x] CLAUDE.md updated (Sprint 2 additions)

---
*Based on Constitution v1.1.0 - See `/.specify/memory/constitution.md`*
