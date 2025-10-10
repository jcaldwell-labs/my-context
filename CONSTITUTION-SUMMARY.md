# Constitution v1.0.0 Summary

**Ratified**: October 9, 2025  
**Status**: Production Ready for v1.0.0 Go-Live  
**File**: `.specify/memory/constitution.md`

---

## Overview

The my-context constitution establishes formal governance for the v1.x.x production lifecycle, codifying the design principles that have guided Sprint 1 and Sprint 2 development into a binding framework for maintenance and ongoing feature development.

---

## Core Principles (6)

### I. Unix Philosophy
- Composability: Each command does one thing well
- Text I/O: stdin/stdout/stderr for all operations
- Zero coupling beyond shared data directory
- Chainable with shell tools (grep, awk, jq)

### II. Cross-Platform Compatibility
- Target platforms: Windows (cmd.exe, PowerShell, git-bash), Linux, macOS (x86/ARM), WSL2
- Transparent path normalization (Windows ↔ POSIX)
- Single executable with environment detection
- Separated dev/runtime directories

### III. Stateful Context Management
- Single active context enforced
- Automatic transitions (no orphaned states)
- State persists in `~/.my-context/`
- Operations default to active context

### IV. Minimal Surface Area
- Command count: ≤12 total (currently 11/12)
- Single-letter aliases for all commands
- Sensible defaults (no required flags for 90% use cases)
- Zero configuration files

### V. Data Portability
- Plain text only: JSON, logs, markdown
- No proprietary formats or binary databases
- Self-contained, relocatable data directory
- Manual editing with standard tools supported

### VI. User-Driven Design
- Features from observed behavior, not speculation
- Retrospective validation (High Value + Low Complexity)
- Automate existing workflows
- Adapt to user conventions

---

## Version Management (v1.x.x)

### Semantic Versioning
- **v1.MAJOR.MINOR** format
- **MAJOR**: Breaking changes (data format, command removal, flag behavior)
  - Migration guide + tooling required
  - 2+ week announcement period
  - Max 1 per quarter
- **MINOR**: New features (commands, flags) - backward compatible
  - Sprint-based cadence (2-3 weeks)
  - Must pass all quality gates
- **PATCH**: Bug fixes, docs, performance - no new features
  - Released as needed (48hr SLA for critical)

### Compatibility Promises
- **Data Backward Compatibility**: v1.x reads all v1.y data (y < x)
- **Data Forward Compatibility**: v1.x writes data readable by older v1.y
- **Command Stability**: v1.0.0 commands exist in all v1.x (deprecation allowed, removal forbidden)

---

## Maintenance Policy

### Security Commitments
- **Critical vulnerabilities**: Patched within 48 hours
- **Backport support**: N-1 MINOR version (current + previous)
- **CVE tracking**: Automated via Dependabot

### Compatibility Maintenance
- **Go versions**: Current + previous 2 major releases (e.g., 1.25, 1.24, 1.23)
- **Platform support**: Windows 10+, macOS 12+, Linux 4.x+, WSL2
- **Test matrix**: All platforms for each release

### Performance Commitments
- List command: <1s for 1000 contexts
- Export command: <1s for 500 notes
- Search command: <1s for 1000 contexts
- Binary size: <10MB per platform

---

## Feature Development

### Sprint-Based Workflow
1. **Specification Phase** (`/speckit.specify`)
   - User scenarios + acceptance criteria
   - Ambiguities resolved
2. **Planning Phase** (`/speckit.plan`)
   - Technical approach + architecture
   - Constitution compliance validation
3. **Task Breakdown** (`/speckit.tasks`)
   - TDD approach with integration tests
4. **Implementation** (`/speckit.implement`)
   - Test-Driven Development enforced
   - Incremental commits
5. **Review & Release**
   - Quality gates passed
   - Documentation updated
   - Retrospective conducted

### Test-Driven Development (TDD) - NON-NEGOTIABLE
1. Write failing tests defining expected behavior
2. Get user/stakeholder approval on test scenarios
3. Implement minimum code to pass tests (Red-Green-Refactor)
4. Integration tests for command contracts
5. Unit tests for core logic + edge cases
6. Backward compatibility tests for migrations

---

## Quality Gates

### Release Readiness Checklist
- [x] All CRITICAL priority items resolved
- [x] All HIGH priority items resolved or deferred with rationale
- [x] TDD requirements met (tests written first, all passing)
- [x] Cross-platform testing (≥3 platforms)
- [x] Backward compatibility verified
- [x] Documentation complete (README, TROUBLESHOOTING, changelog)
- [x] CI/CD pipeline passing
- [x] No known data loss scenarios
- [x] Performance targets met

### Tech Debt Classification
- **CRITICAL**: Blockers (destructive operations untested, data loss risks)
- **HIGH**: Should complete before release (compatibility, docs)
- **MEDIUM**: Defer to next sprint (nice-to-have, additional platforms)
- **LOW**: Backlog (optimizations, community requests)

---

## Governance

### Amendment Process
1. Propose in GitHub issue or retrospective
2. Justify with user impact or technical necessity
3. Version bump decision (MAJOR/MINOR/PATCH)
4. Update constitution + propagate to templates
5. Commit: `docs: amend constitution to vX.Y.Z (description)`

### Compliance Verification
- All PRs verify Core Principles compliance
- Sprint planning validates Constitution Check
- Release retrospectives review adherence

### Constitution Supersedes
This constitution overrides all other practices when conflicts arise. Violating practices must change or the principle must be amended.

---

## Template Consistency Status

✅ **All templates verified compatible**:
- `.specify/templates/plan-template.md` - Constitution Check section generic
- `.specify/templates/spec-template.md` - No constitution references (correct)
- `.specify/templates/tasks-template.md` - TDD alignment verified
- `.cursor/commands/*.md` - No agent-specific references

---

## Next Steps for v1.0.0 Release

1. ✅ **Constitution ratified** - Formal governance established
2. **Tag release**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0: Production-ready with formal constitution"
   git push origin v1.0.0
   ```
3. **Monitor CI/CD** - GitHub Actions builds all platforms
4. **Community testing** - Beta on Windows/macOS
5. **Post-release** - Monitor, collect feedback, plan Sprint 3

---

## Impact

### What Changed
- **From**: Informal design principles in spec.md
- **To**: Formal constitution with version management, maintenance policies, and governance

### What This Means
- **Stability**: v1.x.x guarantees for users (no surprise breaking changes)
- **Clarity**: TDD and quality gates are non-negotiable, not suggestions
- **Accountability**: Performance commitments + security SLAs
- **Evolution**: Clear amendment process for governance changes

### What Stays The Same
- Core design principles (codified from Sprint 1/2 practice)
- Sprint-based workflow
- TDD requirements
- Quality-first approach

---

**Version**: 1.0.0  
**Ratified**: 2025-10-09  
**Next Amendment**: When user feedback or technical evolution requires governance changes  
**Governance File**: `.specify/memory/constitution.md`

