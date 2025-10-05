<!--
SYNC IMPACT REPORT (2025-10-05)

Version Change: 1.0.1 → 1.1.0

Change Summary:
- MINOR bump: New principle added (Principle VI: User-Driven Design)
- Sprint 1 retrospective findings incorporated
- No existing principles modified
- Added guidance on observing organic user patterns and prioritizing user-requested features

New Sections Added:
- Principle VI: User-Driven Design (prioritize observed user workflows over assumed needs)

Modified Sections:
- Development Workflow: Added "Sprint Retrospectives" subsection
- Governance: Updated to reference Sprint 1 learnings

Template Sync Status:
✅ plan-template.md - Constitution Check will need VI added (pending update)
✅ spec-template.md - User Scenarios section aligns with VI (no changes needed)
✅ tasks-template.md - TDD principles remain unchanged (verified)
⚠️ Templates need VI added to Constitution Check gates

Follow-up Items:
- Update plan-template.md Constitution Check to include Principle VI
- Sprint 2 spec should validate against all 6 principles
- Consider adding "user workflow observation" to research.md template

Rationale for MINOR bump:
- New principle expands governance without invalidating existing work
- Sprint 1 delivered successfully under 5 principles
- VI formalizes emergent best practice discovered during Sprint 1
- Backward compatible: existing features comply with VI retroactively

Notes:
- Sprint 1 revealed users adopted "project: phase - description" naming convention organically
- Users manually exported contexts and prefixed notes (PLAN/WORK/OVERHEAD)
- Project filter and export commands are direct responses to observed behavior
- This principle ensures future features validate against real usage patterns

Suggested Commit Message:
docs: amend constitution to v1.1.0 (add Principle VI: User-Driven Design from Sprint 1 learnings)
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

### VI. User-Driven Design
Observe and formalize organic user patterns before imposing structure.
Prioritize features that automate existing manual workflows over assumed needs.
User-requested features (validated through retrospectives) take precedence over speculative enhancements.
Tool adapts to user conventions (e.g., naming patterns) rather than enforcing rigid schemas.

**Rationale**: Sprint 1 revealed users organically adopted `"project: phase - description"` naming conventions and manual export workflows without prompting. High-value features emerge from observing real usage patterns, not predicted requirements. Features that automate existing manual processes (export command, project filters) deliver immediate ROI. User feedback loops (retrospectives, feature requests) prevent feature bloat and ensure development focuses on validated needs.

## Development Workflow

### Test-First Development
All commands MUST have tests before implementation (TDD mandatory).
Tests verify: correct output format, error handling, state transitions, cross-platform paths.
Integration tests validate shell wrapper behavior across cmd/bash/WSL.
Red-Green-Refactor cycle strictly enforced.

**Blocking Gate**: Phase 3.2 (Write Tests) must complete before Phase 3.3 (Implementation).
Cannot proceed to implementation without failing tests proving no code exists yet.

### Code Review Gates
PRs must verify compliance with all six core principles.
Path handling changes require tests on Windows (backslash) and POSIX (forward slash).
New subcommands require documentation update and help text.
Breaking changes to command output format require MAJOR version bump.

### Specification-Driven Development
Every feature starts with a specification (see SDLC.md).
Quality gates enforce: spec → clarify → plan → tasks → implement → review → merge.
Constitution compliance checked at planning stage and before merge.

### Sprint Retrospectives
Each sprint concludes with retrospective ceremony documenting:
- What went well (preserve patterns)
- What went wrong (tech debt, bugs)
- User-requested features with priority rankings
- Observed user patterns and organic conventions
- Constitution compliance assessment

User requests validated through retrospectives must be prioritized over speculative features.
Retrospective findings inform constitution amendments and next sprint planning.

## Governance

This constitution supersedes all implementation decisions. When design conflicts arise, principles above determine resolution priority (I > II > III > IV > V > VI).

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

**Sprint History**:
- Sprint 1 (2025-10-04 to 2025-10-05): Initial release, 8 commands delivered, 100% completion
- Retrospective findings: WSL installation issues, user naming conventions discovered, project filter requested

**Version**: 1.1.0 | **Ratified**: 2025-10-04 | **Last Amended**: 2025-10-05
