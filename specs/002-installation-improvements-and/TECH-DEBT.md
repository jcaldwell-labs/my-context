# Sprint 2 Technical Debt

**Sprint**: 002-installation-improvements-and
**Created**: 2025-10-06
**Status**: Pending UAT completion

---

## Overview

This document tracks features that were planned for Sprint 2 but deferred to future sprints or tech debt backlog.

---

## Deferred to Sprint 3

### Feature 003: Daily Summary & Sprint Progress

**Status**: Specification created, awaiting clarification phase
**Branch**: `003-daily-summary-feature`
**Spec**: `specs/003-daily-summary-feature/spec.md`

**Why Deferred**:
- Sprint 2 UAT incomplete (export/project filter untested)
- Specification has 12 clarifications needed
- Better to observe 1-2 weeks of Sprint 2 usage first
- Follows SDLC: spec → clarify → plan → test → implement

**Planned Timeline**:
- Week 1: Observation & usage pattern analysis
- Week 2: Clarification & planning
- Week 3: TDD implementation & UAT

**Estimated Effort**: 2-3 weeks

**Priority**: Medium (user-requested, but not blocking)

---

## Installation Scripts (P2 - Build Infrastructure)

### T028: install.bat (Windows cmd.exe)

**Task**: Create `scripts/install.bat` - Windows cmd.exe installer
**Status**: Not started
**Planned Sprint**: Sprint 2 → Deferred

**Requirements**:
- Detect existing installation
- Install to %USERPROFILE%\bin\
- Add to user PATH via setx command
- Verify installation
- Preserve ~/.my-context/ data

**Why Deferred**:
- Users can download binary manually for UAT
- Not blocking P1 feature testing
- Low priority vs P1 feature completion

**Recommendation**:
- Create as Sprint 3 task OR
- Add to backlog for community contribution

**Estimated Effort**: 1-2 hours

---

### T029: install.ps1 (Windows PowerShell)

**Task**: Create `scripts/install.ps1` - Windows PowerShell installer
**Status**: Not started
**Planned Sprint**: Sprint 2 → Deferred

**Requirements**:
- Detect existing installation
- Install to $env:USERPROFILE\bin\
- Add to user PATH via [Environment]::SetEnvironmentVariable
- Verify installation
- Preserve ~/.my-context/ data

**Why Deferred**: Same as T028

**Recommendation**: Combine with T028 as "Windows installation scripts" task

**Estimated Effort**: 1-2 hours

---

### T030: curl-install.sh (One-liner installer)

**Task**: Create `scripts/curl-install.sh` - One-liner curl installer
**Status**: Not started
**Planned Sprint**: Sprint 2 → Deferred

**Requirements**:
- Detect platform (Linux, macOS, Windows/WSL via uname)
- Download appropriate binary from GitHub releases
- Verify SHA256 checksum
- Make executable
- Call install.sh with downloaded binary

**Why Deferred**:
- Depends on GitHub releases (not yet set up)
- Not blocking UAT (users have install.sh)
- Nice-to-have vs must-have

**Recommendation**:
- Defer to Sprint 3 OR
- Add to backlog for v1.0.0 release

**Estimated Effort**: 2-3 hours

---

### T001: GitHub Actions Release Workflow

**Task**: Create `.github/workflows/release.yml` - CI/CD workflow
**Status**: Not started
**Planned Sprint**: Sprint 2 → Deferred

**Requirements**:
- Multi-platform builds (linux/amd64, windows/amd64, darwin/amd64, darwin/arm64)
- Generate SHA256 checksums
- Upload to releases
- Trigger on tag push

**Why Deferred**:
- Can build locally for UAT
- CI/CD nice-to-have for sprint testing
- More important for production releases

**Recommendation**:
- Complete before v2.0.0 final release
- Not blocking Sprint 2 UAT sign-off

**Estimated Effort**: 2-4 hours (includes testing)

---

## Features Implemented but NOT Tested

### List Enhancements (P2)

**Status**: Code complete, untested
**Risk Level**: Low

**Untested Features**:
- `--limit <n>` flag (custom limit)
- `--search <term>` flag (substring search)
- `--all` flag (show all, no limit)
- `--archived` flag (show archived only)
- `--active-only` flag (show only active)
- Combined filters (project + search + limit)

**Why Untested**:
- P2 features (not user-requested)
- Time prioritized for P1 features
- Low risk (filtering logic straightforward)

**Testing Plan**:
- Run Demo scenarios 3.1-3.4 (10 minutes)
- If pass: Document as validated
- If fail: Fix in Sprint 2 or defer to tech debt

**Priority**: Medium (should test before sign-off)

---

### Archive Command (P2)

**Status**: Code complete, CRITICAL to test
**Risk Level**: High (affects data visibility)

**Untested Features**:
- `archive <context>` command
- `is_archived` flag in meta.json
- Hidden from default list
- Visible with `--archived` flag
- Cannot archive active context

**Why Untested**:
- Just completed in Sprint 2
- Requires careful validation (data integrity)

**Testing Plan**:
- MUST test before sign-off (Demo scenario 5.1)
- Verify is_archived flag set
- Verify list visibility behavior
- Verify error handling

**Priority**: HIGH (test before sign-off)

---

### Delete Command (P2)

**Status**: Code complete, CRITICAL to test
**Risk Level**: CRITICAL (destructive operation)

**Untested Features**:
- `delete <context>` command
- Confirmation prompt
- `--force` flag (skip confirmation)
- Cannot delete active context
- Directory removal
- Transitions log preserved

**Why Untested**:
- Just completed in Sprint 2
- DESTRUCTIVE operation - requires thorough testing

**Testing Plan**:
- MUST test before sign-off (Demo scenario 5.2-5.3)
- Test confirmation flow
- Test --force flag
- Test error cases (active context)
- Verify data actually deleted

**Priority**: CRITICAL (BLOCKER for sign-off if not tested)

---

### Export --all Flag (P2)

**Status**: Code complete, untested
**Risk Level**: Low

**Untested Feature**:
- `export --all --to <dir>` (export all contexts to directory)

**Why Untested**:
- P2 feature
- Basic export tested (P1)
- Low risk (loops over export function)

**Testing Plan**:
- Quick test (5 minutes)
- OR defer to tech debt

**Priority**: Low (can defer)

---

### JSON Output for All Commands (P2)

**Status**: Code complete, partially tested
**Risk Level**: Low

**Tested**:
- `start --json`
- `show --json`
- `list --json`

**Untested**:
- `export --json`
- `archive --json`
- `delete --json`
- `note --json`
- `file --json`
- `touch --json`
- `history --json`

**Why Untested**:
- P2 feature (scripting use case)
- Core commands validated
- Low risk (consistent implementation)

**Testing Plan**:
- Add to regression test suite
- OR defer to tech debt

**Priority**: Low (can defer)

---

## Documentation Gaps

### README.md Updates (T032)

**Status**: Partially complete
**Missing**:
- "Building from Source" section
- "Installation" section with curl one-liner
- Document new commands (export, archive, delete)
- Document new flags (--project, --limit, --search, etc.)
- Troubleshooting link

**Why Incomplete**: Prioritized feature completion over docs

**Recommendation**: Complete after UAT sign-off, before v2.0.0 release

**Estimated Effort**: 1-2 hours

---

### TROUBLESHOOTING.md (T004)

**Status**: Complete (created earlier)
**Location**: `docs/TROUBLESHOOTING.md`

✅ No tech debt

---

## Backward Compatibility Testing (T012)

**Status**: Partially tested
**Tested**:
- Sprint 1 contexts visible in Sprint 2
- Basic operations (show, list) work on old contexts

**Untested**:
- New features (export, archive, delete) on Sprint 1 contexts
- is_archived flag backward compatibility (omitempty JSON tag)

**Testing Plan**:
- Run Demo scenario 6.1 (5 minutes)
- Test export on Sprint 1 context
- Verify is_archived defaults to false

**Priority**: Medium (should test before sign-off)

---

## Performance Benchmarks (T038)

**Status**: Not tested
**Targets** (from data-model.md):
- List with 1000 contexts: <1s
- Export with 500 notes: <1s
- Search across 1000 contexts: <1s

**Why Untested**:
- Limited test data (< 50 contexts)
- Performance not reported as issue
- Nice-to-have vs must-have

**Recommendation**: Defer to performance testing phase (post-v2.0.0)

**Priority**: Low

---

## Cross-Platform Testing (T040)

**Status**: Partial (Windows git-bash only)
**Tested Platforms**:
- ✅ Windows (git-bash)

**Untested Platforms**:
- 🔲 Windows (cmd.exe)
- 🔲 Windows (PowerShell)
- 🔲 Ubuntu (WSL)
- 🔲 Ubuntu (native)
- 🔲 macOS (Intel)
- 🔲 macOS (ARM)

**Why Untested**:
- Limited test environments
- Core functionality validated on one platform
- Cross-platform Go builds (should work)

**Recommendation**:
- Community testing OR
- Defer to beta release phase

**Priority**: Medium (important for v2.0.0 final)

---

## Tech Debt Backlog Summary

### CRITICAL (MUST do before sign-off) ✅ COMPLETE
- [x] Test delete command (Demo 5.2-5.3) - **PASSED** (Oct 9, 2025)
- [x] Test archive command (Demo 5.1) - **PASSED** (Oct 9, 2025)

### HIGH (Should do before v2.0.0 release) ✅ COMPLETE
- [x] Test list enhancements (Demo 3.1-3.4) - **PASSED** (Oct 9, 2025)
- [x] Test backward compatibility (Demo 6.1) - **PASSED** (Oct 9, 2025)
- [x] Complete README.md updates (T032) - **VERIFIED** (Oct 9, 2025)
- [x] Create GitHub Actions workflow (T001) - **EXISTS** (Oct 9, 2025)

### MEDIUM (Should do in Sprint 3 or v2.1)
- [ ] Create Windows installers (T028, T029)
- [ ] Create curl installer (T030)
- [ ] Cross-platform testing (T040)
- [ ] JSON output validation for all commands

### LOW (Can defer to backlog)
- [ ] Test export --all flag
- [ ] Performance benchmarks (T038)
- [ ] Issue templates (T034)

---

## Sprint 3 Roadmap

**Planned**:
1. Daily summary feature (003-daily-summary-feature)
2. Address HIGH priority tech debt from above
3. Community feedback from Sprint 2 UAT

**Estimated Duration**: 2-3 weeks

---

## How to Use This Document

**Before UAT Sign-Off**:
1. Complete CRITICAL items (delete/archive testing)
2. Review HIGH items - can any be done quickly?
3. Document decision for MEDIUM/LOW items (defer or do)

**After UAT Sign-Off**:
1. Create GitHub issues for deferred items
2. Label by priority (critical/high/medium/low)
3. Assign to Sprint 3 or backlog
4. Update this document with issue links

**During Sprint 3**:
1. Review deferred items from Sprint 2
2. Prioritize based on user feedback
3. Complete HIGH priority items
4. Re-evaluate MEDIUM/LOW priorities

---

*Tech debt tracker created: 2025-10-06*
*Sprint: 002-installation-improvements-and*
*Next review: After Sprint 2 UAT completion*
