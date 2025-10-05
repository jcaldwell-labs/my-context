<!--
SYNC IMPACT REPORT (2025-10-05)

Version Change: 1.0.0 → 1.0.1

Change Summary:
- PATCH bump: Constitution restored after accidental template reset
- No principle modifications
- No new sections added or removed
- Clarification: Project transitioning from initial development to maintenance phase

Template Sync Status:
✅ plan-template.md - Constitution Check section references 5 principles (verified)
✅ spec-template.md - Cross-platform testing guidance present (verified)
✅ tasks-template.md - Phase 3.2 TDD enforcement present (verified)
✅ agent-file-template.md - Generic guidance, no constitution-specific updates needed

Follow-up Items:
- None (all templates remain aligned with constitution v1.0.0 principles)

Notes:
- Project status: Transitioning to maintenance mode with a few more features planned
- Constitution principles remain stable and applicable
- Next version bump likely MINOR when new maintenance-phase principles emerge

Suggested Commit Message:
docs: restore constitution to v1.0.1 (patch - no principle changes)
-->

# My-Context-Copilot Constitution

## Core Principles

### I. Unix Philosophy
Every component follows Unix principles: Do one thing well, text-based I/O, composable commands.
Each subcommand must be independently useful and chainable with standard shell tools.
Commands read from stdin/args and write to stdout (data) or stderr (errors).
No unnecessary coupling between commands - each must function in isolation.

**Rationale**: Context management is a toolchain problem requiring composability. Users need to integrate context operations into existing workflows (git hooks, CI/CD, shell scripts). Text protocols ensure universal compatibility across environments (cmd, git-bash, WSL).

### II. Cross-Platform Compatibility
Tool must work seamlessly across Windows (cmd.exe, PowerShell), git-bash, and WSL environments.
Single executable with wrapper/shim layer for environment detection.
Path handling must normalize Windows and POSIX conventions automatically.
Installation separates dev folder from runtime home directory.

**Rationale**: Development workflows span multiple shells on Windows. Users switch between cmd, bash, and WSL without changing tools. Context data must persist independently of development environment.

### III. Stateful Context Management
One active context at a time; starting new context stops previous automatically.
Context state persists in dedicated home directory (separate from install path).
All subcommands default to current active context unless explicitly overridden.
Context lifecycle: start → (note/file/touch operations) → stop/start-new.

**Rationale**: Mental context switching requires clear boundaries. Automatic context closure prevents state confusion. Centralized state storage enables cross-shell context sharing and prevents orphaned context data.

### IV. Minimal Surface Area
Commands: start, stop, note, file, touch, list, show, history (minimal essential operations).
Each command has single-letter alias (e.g., `my-context n` for note).
No configuration files unless absolutely necessary - sensible defaults.
Help system built-in and always available.

**Rationale**: Cognitive overhead kills adoption. Users must remember <10 commands total. Single-letter aliases support muscle memory. Zero-config startup removes friction for new users.

### V. Data Portability
All context data stored as plain text (timestamps, notes, file associations).
Export/import must work via standard file operations (cp, cat, grep work directly).
No proprietary formats - JSON for structured data, newline-delimited for logs.
Context home directory is self-contained and relocatable.

**Rationale**: Users must own their data without tool lock-in. Plain text enables git versioning, grep searches, and direct editing. Standard formats support scripting and third-party tool integration.

## Development Workflow

### Test-First Development
All commands MUST have tests before implementation (TDD mandatory).
Tests verify: correct output format, error handling, state transitions, cross-platform paths.
Integration tests validate shell wrapper behavior across cmd/bash/WSL.
Red-Green-Refactor cycle strictly enforced.

**Blocking Gate**: Phase 3.2 (Write Tests) must complete before Phase 3.3 (Implementation).
Cannot proceed to implementation without failing tests proving no code exists yet.

### Code Review Gates
PRs must verify compliance with all five core principles.
Path handling changes require tests on Windows (backslash) and POSIX (forward slash).
New subcommands require documentation update and help text.
Breaking changes to command output format require MAJOR version bump.

### Specification-Driven Development
Every feature starts with a specification (see SDLC.md).
Quality gates enforce: spec → clarify → plan → tasks → implement → review → merge.
Constitution compliance checked at planning stage and before merge.

## Governance

This constitution supersedes all implementation decisions. When design conflicts arise, principles above determine resolution priority (I > II > III > IV > V).

Amendments require:
1. Documentation of principle change rationale
2. Update to all affected templates (plan, spec, tasks)
3. Migration plan for existing users if breaking changes introduced
4. Version bump following semantic versioning rules

All code reviews and design discussions must reference applicable principles when making decisions.

Versioning follows semantic versioning (MAJOR.MINOR.PATCH):
- **MAJOR**: Breaking CLI changes, output format incompatibilities, principle removal/redefinition
- **MINOR**: New subcommands, new optional flags, backward-compatible additions, new principles
- **PATCH**: Bug fixes, documentation, internal refactoring, clarifications

For runtime development guidance, refer to `CLAUDE.md` (Claude Code), `README.md` (user documentation), and `SDLC.md` (development process).

**Version**: 1.0.1 | **Ratified**: 2025-10-04 | **Last Amended**: 2025-10-05
