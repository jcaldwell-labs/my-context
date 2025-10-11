<!--
Sync Impact Report:
Version: Template → 1.0.0 (MAJOR: Initial production ratification)
Changes:
  - Added: All 6 core design principles (Unix Philosophy, Cross-Platform Compatibility, 
           Stateful Context Management, Minimal Surface Area, Data Portability, User-Driven Design)
  - Added: Version Management section (semantic versioning for v1.x.x production lifecycle)
  - Added: Maintenance Policy section (security, compatibility, performance commitments)
  - Added: Feature Development section (sprint-based workflow, TDD requirements)
  - Added: Quality Gates section (release readiness criteria)
  - Added: Governance section (amendment process, compliance requirements)
Templates Updated:
  ✅ plan-template.md - Constitution Check section already references design principles
  ✅ spec-template.md - Design Principles Alignment section matches constitution principles
  ✅ tasks-template.md - Task categorization aligns with TDD and quality gate requirements
  ✅ All speckit commands - No agent-specific references, generic guidance preserved
Follow-up TODOs: None (all placeholders resolved)
-->

# My-Context Constitution

## Core Principles

### I. Unix Philosophy
Every feature MUST adhere to Unix composability principles:
- Each command does one thing well and does it completely
- Commands accept text input and produce text output (stdin/stdout/stderr)
- Zero coupling between commands beyond shared data directory
- Output can be chained with shell tools (grep, awk, jq, etc.)

**Rationale**: Composability ensures longevity. Users can build workflows without waiting for feature requests. Text I/O enables debugging, scripting, and integration with existing toolchains.

### II. Cross-Platform Compatibility
The tool MUST work identically across all supported platforms:
- Target platforms: Windows (cmd.exe, PowerShell, git-bash), Linux, macOS (x86/ARM), WSL2
- Path handling MUST normalize Windows backslashes and POSIX forward slashes transparently
- Single executable with automatic environment detection (no platform-specific builds requiring conditional logic)
- Installation separates development directory from runtime home directory (`~/.my-context/`)

**Rationale**: Developers switch environments frequently. A tool that "almost works" on a platform is worse than one that doesn't claim support. Cross-platform consistency builds trust.

### III. Stateful Context Management
The tool enforces single-active-context semantics:
- Only one context can be active at any time
- Starting a new context automatically stops the previous one (no orphaned states)
- State persists in dedicated home directory (`~/.my-context/` or `$MY_CONTEXT_HOME`)
- All operations default to the active context unless explicitly targeting another

**Rationale**: Mental model simplicity. Developers already juggle enough state (git branches, terminal tabs, IDE windows). One active context eliminates decision paralysis and accidental cross-contamination.

### IV. Minimal Surface Area
Command count MUST remain ≤14 total commands:
- Each command has a single-letter alias (e.g., `start`/`s`, `note`/`n`)
- Defaults are sensible (no required flags for 90% use cases)
- Zero configuration files required (conventions over configuration)
- Built-in help available via `--help` for every command

**Rationale**: Cognitive load is the enemy of adoption. A tool that requires consulting docs for basic operations will be abandoned. Muscle memory forms around ~10 commands; beyond that, users resort to cheat sheets.

**Amendment (v2.2.0 - Iteration 006)**: Limit increased from 12→14 to accommodate distributed journal coordination capabilities (signal, watch, resume). Justified by: (1) Proven user need (final-completion handoff required manual coordination), (2) Core capability (not optional feature), (3) Still minimal (14 commands manageable).

**Current Command Count**: 13/14 (start, stop, resume, note, file, touch, show, list, export, archive, delete, history, signal, watch)

### V. Data Portability
All data MUST be stored as human-readable plain text:
- Metadata in JSON (`.json` files)
- Logs in newline-delimited text (`.log` files)
- Exports in Markdown (`.md` files)
- No proprietary formats, no binary databases, no lock-in
- Data directory is self-contained and relocatable
- Users can grep, sed, awk, or manually edit data files without tool assistance

**Rationale**: Tools come and go. Data outlives them. Plain text ensures users can migrate, backup, version-control, or rescue their data using standard Unix tools—even if my-context ceases to exist.

### VI. User-Driven Design
Features MUST originate from observed user behavior, not speculation:
- Monitor how users actually name contexts, structure workflows, and encounter friction
- Validate feature requests through retrospectives (High Value + Low Complexity prioritized)
- Automate existing manual workflows rather than inventing new ones
- Adapt to user conventions (e.g., "project: phase" naming pattern)

**Rationale**: Features built on assumptions create maintenance burden without delivering value. User-driven design ensures every addition solves a real problem for real users.

## Version Management

### Semantic Versioning (v1.x.x Production Lifecycle)
The project follows strict semantic versioning for v1.x.x releases:
- **v1.MAJOR.MINOR** format
- **MAJOR**: Backward-incompatible changes (data format, command removal, flag behavior)
  - MUST include migration guide and tooling
  - Requires announcement 2+ weeks before release
  - Breaking changes batched (no more than 1 MAJOR release per quarter)
- **MINOR**: New features, new commands, new flags (backward-compatible)
  - Sprint-based cadence (typically 2-3 weeks)
  - Must pass all quality gates (see Quality Gates section)
- **PATCH**: Bug fixes, documentation, performance improvements (no new features)
  - Released as needed (hotfixes prioritized within 48 hours for critical issues)

### Version Compatibility Promise
- **Data Backward Compatibility**: v1.x MUST read data created by any v1.y where y < x
- **Data Forward Compatibility**: v1.x SHOULD write data readable by v1.y where y < x (use optional fields, not required)
- **Command Stability**: Commands and flags in v1.0.0 MUST exist in all v1.x releases (deprecation allowed, removal forbidden)

## Maintenance Policy

### Security Commitments
- Critical security vulnerabilities patched within 48 hours of discovery
- Security patches backported to previous MINOR version (N-1 support)
- CVE tracking for all dependencies (automated via Dependabot or equivalent)

### Compatibility Maintenance
- Go version support: Current + previous 2 major releases (e.g., Go 1.25, 1.24, 1.23)
- Platform support: Windows 10+, macOS 12+, Linux kernel 4.x+, WSL2
- Test matrix MUST cover all supported platforms for each release

### Performance Commitments
- List command: <1s for 1000 contexts
- Export command: <1s for 500 notes
- Search command: <1s for 1000 contexts
- Binary size: <10MB per platform (static linking with CGO_ENABLED=0)

## Feature Development

### Sprint-Based Workflow
New features follow a structured sprint process:
1. **Specification Phase** (`/speckit.specify`):
   - User scenarios and acceptance criteria defined
   - Ambiguities marked and resolved
   - Stakeholder approval before proceeding
2. **Planning Phase** (`/speckit.plan`):
   - Technical approach and architecture decisions documented
   - Constitution compliance validated
   - Data model and API contracts defined
3. **Task Breakdown Phase** (`/speckit.tasks`):
   - Implementation tasks with TDD approach
   - Integration test scenarios defined
4. **Implementation Phase** (`/speckit.implement`):
   - Test-Driven Development (TDD) enforced (tests written first, implementation follows)
   - Incremental commits with context tracking
5. **Review & Release**:
   - All quality gates passed (see Quality Gates section)
   - Documentation updated (README, TROUBLESHOOTING)
   - Sprint retrospective conducted

### Test-Driven Development (TDD) - NON-NEGOTIABLE
All new features and bug fixes MUST follow TDD:
1. Write failing tests that define expected behavior
2. Get user/stakeholder approval on test scenarios
3. Implement minimum code to make tests pass (Red-Green-Refactor)
4. Integration tests for command contracts (input/output validation)
5. Unit tests for core logic and edge cases
6. Backward compatibility tests for data migrations

**Rationale**: Tests are executable specifications. Writing tests first ensures clarity on requirements and prevents scope creep. TDD catches regressions before they reach users.

## Quality Gates

### Release Readiness Criteria
A release is ready when ALL of the following are true:
- [x] All CRITICAL priority items resolved
- [x] All HIGH priority items resolved or explicitly deferred with rationale
- [x] TDD requirements met (tests written first, all passing)
- [x] Cross-platform testing on at least 3 platforms (Linux, Windows, macOS or WSL)
- [x] Backward compatibility verified (Sprint N data works in Sprint N+1)
- [x] Documentation complete (README, TROUBLESHOOTING, changelog)
- [x] CI/CD pipeline passing (multi-platform builds, SHA256 checksums)
- [x] No known data loss scenarios
- [x] Performance targets met (list, export, search under target thresholds)

### Tech Debt Management
- **CRITICAL**: Blockers for release (destructive operations untested, data loss risks)
- **HIGH**: Should complete before release (backward compatibility, documentation gaps)
- **MEDIUM**: Can defer to next sprint (nice-to-have features, additional platforms)
- **LOW**: Backlog (performance optimizations, community requests)

Tech debt MUST be documented in `TECH-DEBT.md` with priority, effort estimate, and deferral rationale.

## Governance

### Amendment Process
1. Propose amendment in GitHub issue or sprint retrospective
2. Justify with user impact, technical necessity, or evolution of best practices
3. Version bump decision:
   - MAJOR: Remove/redefine existing principle
   - MINOR: Add new principle or expand existing guidance
   - PATCH: Clarifications, wording improvements, typos
4. Update constitution and propagate changes to templates (`plan-template.md`, `spec-template.md`, `tasks-template.md`)
5. Commit with message: `docs: amend constitution to vX.Y.Z (description)`

### Compliance Verification
- All PRs MUST verify compliance with Core Principles (checklist in PR template)
- Sprint planning MUST validate Constitution Check in `plan.md`
- Release retrospectives MUST review adherence to Version Management and Quality Gates

### Constitution Supersedes
This constitution supersedes all other practices, conventions, or prior decisions when conflicts arise. If a practice violates a Core Principle, the practice MUST change or the principle MUST be amended through the governance process.

### Guidance References
- Runtime development guidance: `README.md`, `TROUBLESHOOTING.md`
- Sprint workflow guidance: `.specify/templates/commands/*.md`
- Agent-specific patterns: `.cursorrules`, `.cursor/rules/*.mdc`

**Version**: 1.0.0 | **Ratified**: 2025-10-09 | **Last Amended**: 2025-10-09
