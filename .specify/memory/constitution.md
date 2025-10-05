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
Commands: start, stop, note, file, touch, list, show (minimal essential operations).
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
All commands must have tests before implementation (TDD mandatory).
Tests verify: correct output format, error handling, state transitions, cross-platform paths.
Integration tests validate shell wrapper behavior across cmd/bash/WSL.
Red-Green-Refactor cycle strictly enforced.

### Code Review Gates
PRs must verify compliance with all five core principles.
Path handling changes require tests on Windows (backslash) and POSIX (forward slash).
New subcommands require documentation update and help text.
Breaking changes to command output format require MAJOR version bump.

## Governance

This constitution supersedes all implementation decisions. When design conflicts arise, principles above determine resolution priority (I > II > III > IV > V).

Amendments require:
1. Documentation of principle change rationale
2. Update to all affected templates (plan, spec, tasks)
3. Migration plan for existing users if breaking changes introduced

All code reviews and design discussions must reference applicable principles when making decisions.

Versioning follows semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking CLI changes, output format incompatibilities
- MINOR: New subcommands, new optional flags, backward-compatible additions
- PATCH: Bug fixes, documentation, internal refactoring

**Version**: 1.0.0 | **Ratified**: 2025-10-04 | **Last Amended**: 2025-10-04

<!--
SYNC IMPACT REPORT (2025-10-04)

Version Change: Initial creation → 1.0.0

Constitution Summary:
- Five core principles established
- Unix Philosophy: composability, text I/O, single-purpose commands
- Cross-Platform Compatibility: cmd/bash/WSL support requirement
- Stateful Context Management: one active context, automatic transitions
- Minimal Surface Area: <10 commands, zero-config defaults
- Data Portability: plain text storage, no lock-in

Template Sync Status:
✅ plan-template.md - Constitution Check section needs update to reflect 5 principles
✅ spec-template.md - Requirements must validate cross-platform compatibility
✅ tasks-template.md - Test tasks must include Windows/POSIX path testing

Action Items:
1. Update plan-template.md Constitution Check with new principles
2. Update spec-template.md to include cross-platform testing scenarios
3. Update tasks-template.md Phase 3.2 to mandate platform-specific tests
4. Consider adding command-specific templates for CLI tool development

Suggested Commit Message:
docs: establish constitution v1.0.0 (Unix philosophy + cross-platform CLI principles)
-->
