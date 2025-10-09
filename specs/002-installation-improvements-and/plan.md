# Implementation Plan: Installation & Usability Improvements

**Branch**: `002-installation-improvements-and` | **Date**: 2025-10-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-installation-improvements-and/spec.md`

## Summary

Sprint 2 delivers multi-platform distribution and user-requested CLI enhancements for the my-context tool. Primary requirements include pre-built binaries for Windows/Linux/macOS (x86/ARM), automated installation scripts, project-based context filtering, markdown export functionality, and context lifecycle management (archive/delete). This sprint addresses Sprint 1 retrospective feedback focusing on installation friction and workflow efficiency improvements.

## Technical Context

**Language/Version**: Go 1.25.1 (requires 1.21+)  
**Primary Dependencies**: Cobra 1.10.1 (CLI framework), Viper 1.21.0 (configuration), standard library  
**Storage**: JSON files in `~/.my-context/` directory (plain text, greppable)  
**Testing**: Go's built-in testing framework (`testing` package), table-driven tests  
**Target Platform**: Cross-platform CLI (Windows amd64, Linux amd64, macOS amd64/arm64), WSL2 support  
**Project Type**: Single project (CLI tool)  
**Performance Goals**: List command <1s for 1000 contexts, Export <1s for 500 notes, Search <1s for 1000 contexts  
**Constraints**: Static binary <10MB per platform, zero runtime dependencies (CGO_ENABLED=0), backward compatible with Sprint 1 data  
**Scale/Scope**: Expected usage: 50-100 contexts per developer, 500 notes per context max, single-user tool

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Note**: Formal constitution is pending (template only). Using design principles from spec.md for compliance review.

**Unix Philosophy** (composability, text I/O, single-purpose)
- [x] Each component does one thing well - Export handles export, Archive handles archiving, Delete handles deletion
- [x] Commands accept text input and produce text output (stdin/stdout/stderr for all commands)
- [x] No unnecessary coupling - New commands integrate via existing core layer
- [x] Can be chained with shell tools - List output supports grep/awk, export produces markdown files

**Cross-Platform Compatibility** (Windows, git-bash, WSL)
- [x] Works identically across cmd.exe, PowerShell, git-bash, and WSL - Multi-platform builds + install scripts
- [x] Path handling normalizes Windows backslashes and POSIX forward slashes - Already handled in Sprint 1
- [x] Single executable with environment detection - Go's os package provides cross-platform compatibility
- [x] Installation separates dev folder from runtime home directory - Binaries install to user bin, data in ~/.my-context/

**Stateful Context Management** (one active context, automatic transitions)
- [x] Only one active context at a time - Preserved from Sprint 1
- [x] Starting new context automatically stops previous - No changes to this behavior
- [x] State persists in dedicated home directory - Enhanced with archive flag in meta.json
- [x] All operations default to active context - Export/archive/delete accept context name as argument

**Minimal Surface Area** (<=10 commands, single-letter aliases, zero config)
- [x] Feature adds minimal new commands/flags - Adds 3 commands (export/archive/delete), total 11 commands
- [x] Each command has single-letter alias - export (e), archive (a), delete (d)
- [x] Defaults are sensible - List shows last 10, archive hides from default list, delete prompts for confirmation
- [x] Built-in help available - Cobra provides automatic help generation

**Data Portability** (plain text, no lock-in, greppable)
- [x] All data stored as plain text - JSON for metadata, newline-delimited logs, markdown exports
- [x] Can export/import with standard file operations - Export command generates portable markdown
- [x] No proprietary formats - All formats are standard (JSON, markdown, text logs)
- [x] Data directory is self-contained and relocatable - ~/.my-context/ remains self-contained

**User-Driven Design** (observe patterns, automate workflows, validate requests)
- [x] Feature responds to observed user behavior - All features from Sprint 1 retrospective user feedback
- [x] Automates existing manual workflow - Export automates context sharing, project filter automates grouping
- [x] User request validated through retrospective - All 6 features ranked P1 (High Value, Low Complexity)
- [x] Adapts to user conventions - Project parsing respects "project: phase" naming pattern observed in usage

**Violations**: None. All design principles maintained. Command count increase (8→11) justified by distinct lifecycle operations (export for sharing, archive for completion, delete for cleanup).

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
cmd/
└── my-context/           # Main entry point (main.go, root command setup)

internal/
├── commands/             # Cobra command implementations (start, stop, list, note, etc.)
├── core/                 # Business logic (context management, storage, transitions)
├── models/               # Data structures (Context, Transition, metadata)
└── output/               # Formatters (human-readable, JSON, markdown)

tests/
├── integration/          # End-to-end command tests, multi-command workflows
└── unit/                 # Pure function tests, model validation, parser logic

scripts/                  # Build and installation scripts (NEW in Sprint 2)
├── build-all.sh         # Multi-platform build script
├── install.sh           # Unix/macOS/WSL installer
├── install.bat          # Windows cmd.exe installer
├── install.ps1          # Windows PowerShell installer
└── curl-install.sh      # One-liner platform-detection installer

docs/                     # User documentation
└── TROUBLESHOOTING.md   # Platform-specific installation issues (NEW)

specs/                    # Feature specifications (managed by speckit)
└── 002-installation-improvements-and/
    ├── spec.md
    ├── plan.md           # This file
    ├── research.md       # Phase 0 output
    ├── data-model.md     # Phase 1 output
    ├── quickstart.md     # Phase 1 output
    ├── contracts/        # Phase 1 output
    └── tasks.md          # Phase 2 output (from /tasks command)
```

**Structure Decision**: Single Go project with clean separation of concerns. Commands layer (Cobra CLI) → Core layer (business logic) → Models layer (data structures). Tests organized by scope (unit vs integration). Sprint 2 adds `scripts/` directory for distribution and `docs/TROUBLESHOOTING.md` for user support. Existing structure from Sprint 1 requires no changes - new features integrate via new command files and core extensions.

## Complexity Tracking

*No violations*. All design principles maintained. Command count increase (8→11) justified by distinct lifecycle operations.

## Phase 0: Outline & Research

**Output**: research.md

**Key Decisions**:
1. Static binary build strategy (CGO_ENABLED=0 for zero dependencies)
2. User-level installation (no sudo/admin required)
3. Project name extraction pattern (text before first colon)
4. Archive vs Delete implementation (metadata flag vs physical removal)
5. Export format (Markdown for human-readability)
6. List default limit (10 contexts with --all escape hatch)
7. Backward compatibility strategy (additive JSON fields with omitempty)
8. GitHub Actions for release builds (automated multi-platform)
9. Installation script language choices (POSIX sh, batch, PowerShell)

**Status**: ✅ Complete - All technical unknowns resolved, documented in research.md

---

## Phase 1: Design & Contracts

**Output**: data-model.md, contracts/, quickstart.md, .cursor/rules/specify-rules.mdc

**Entities Defined** (5):
1. Context (modified): Added is_archived field
2. ProjectMetadata (new): Runtime-derived project grouping
3. ExportDocument (new): Markdown representation of context
4. BinaryArtifact (new): Platform-specific binary metadata
5. InstallationMetadata (new): Upgrade detection metadata

**Command Contracts** (5):
1. export-command.md: Single/all export, markdown/JSON, overwrite protection
2. archive-command.md: Mark context archived, hide from default list
3. delete-command.md: Permanent removal with confirmation
4. list-enhanced.md: Filters (project, search, archived, active-only), pagination
5. start-with-project.md: --project flag for automatic prefixing

**Test Scenarios** (9 in quickstart.md):
- Multi-platform installation
- Project-based workflow
- Export and share
- Context lifecycle (archive/delete)
- List enhancements with large datasets
- Bug fixes validation
- Cross-platform installation (Windows)
- Backward compatibility (Sprint 1 → Sprint 2)
- JSON output for scripting

**Agent Context**: Updated .cursor/rules/specify-rules.mdc with Go 1.25.1, Cobra/Viper stack, JSON storage

**Status**: ✅ Complete - Ready for task generation

---

## Phase 2: Task Planning Approach

**Task Generation Strategy** (executed by /tasks command):
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs
- Each contract → integration test task [P]
- Each entity → model creation/modification task [P]
- Each user scenario → validation task
- Implementation tasks to make tests pass

**Ordering Strategy**:
- TDD order: Tests before implementation
- Dependency order: Models → Core → Commands → Scripts
- Mark [P] for parallel execution (independent files)

**Estimated Output**: ~40-45 tasks organized by phase:
- Phase 3.1: Setup & Build Infrastructure (4 tasks)
- Phase 3.2: Tests First - TDD (8 contract/unit tests) [P]
- Phase 3.3: Core Implementation (models, core logic, commands)
- Phase 3.4: Build & Installation Scripts (5 scripts)
- Phase 3.5: Documentation & Polish (README, TROUBLESHOOTING)
- Phase 3.6: Integration & Validation (quickstart scenarios, benchmarks)

**IMPORTANT**: This phase description is for planning only. Execute `/tasks` command to generate actual tasks.md.

---

## Progress Tracking

**Phase Status**:
- [x] Phase 0: Research complete
- [x] Phase 1: Design complete
- [x] Phase 2: Task planning approach documented (tasks.md generation ready)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS (no violations)
- [x] Post-Design Constitution Check: PASS (no violations)
- [x] All NEEDS CLARIFICATION resolved (none existed - spec was complete)
- [x] Complexity deviations documented (none required)

**Artifacts Generated**:
- [x] plan.md (this file)
- [x] research.md (9 technical decisions)
- [x] data-model.md (5 entities with validation rules)
- [x] contracts/ (5 command specifications)
- [x] quickstart.md (9 end-to-end scenarios)
- [x] .cursor/rules/specify-rules.mdc (agent context)
- [ ] tasks.md (pending /tasks command)

---

**Next Steps**:
1. Execute `/tasks` command to generate tasks.md
2. Review tasks for completeness and ordering
3. Begin implementation following TDD approach (Phase 3.2 tests first)

---

*Plan completed: 2025-10-09*  
*Based on spec.md (complete), Sprint 1 retrospective findings, and design principles*  
*Ready for task generation*
