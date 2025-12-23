# Software Development Life Cycle (SDLC)

**Project**: my-context-copilot  
**Version**: 1.0.0  
**Last Updated**: 2025-10-05

---

## Overview

This document defines the complete software development lifecycle for my-context-copilot, from feature conception through merge and release. The process is based on specification-driven development with mandatory quality gates.

### Principles

1. **Specification-First**: Every feature starts with a clear spec
2. **Test-Driven Development**: Tests before implementation (non-negotiable)
3. **Constitution-Aligned**: All features comply with 5 core principles
4. **Quality Gates**: Each stage has exit criteria that MUST be met
5. **Incremental Delivery**: Small, frequent commits over large batches

### Tools

- **Spec Kit**: `.specify/scripts/bash/` - Automation scripts
- **Templates**: `.specify/templates/` - Standard document formats
- **Constitution**: `.specify/memory/constitution.md` - Guiding principles

---

## Stage 1: Specification (/specify)

### Purpose
Define WHAT the feature does and WHY it's needed, without implementation details.

### Entry Criteria
- User request, bug report, or identified need
- Problem statement clear

### Activities
1. Run `bash .specify/scripts/bash/create-new-feature.sh "<description>"`
2. Creates feature branch: `###-feature-name`
3. Creates `specs/###-feature-name/spec.md` from template
4. Fill out spec sections:
   - User scenarios (minimum 5)
   - Acceptance criteria
   - Functional requirements (minimum 10)
   - Key entities (if data involved)
   - Edge cases

### Exit Criteria (Quality Gate)
- [ ] All mandatory sections complete
- [ ] Requirements are testable (Given/When/Then format)
- [ ] No implementation details (no languages, frameworks, APIs)
- [ ] Review checklist passed
- [ ] No [NEEDS CLARIFICATION] markers OR clarification planned

### Artifacts
- `specs/###-feature-name/spec.md`
- Feature branch `###-feature-name`

### Time Estimate
Simple feature: 30-60 minutes  
Complex feature: 2-4 hours

---

## Stage 2: Clarification (/clarify)

### Purpose
Resolve all ambiguities and uncertainties in the specification.

### Entry Criteria
- spec.md contains [NEEDS CLARIFICATION] markers
- OR spec has implicit ambiguities

### Activities
1. Review spec for assumptions and gaps
2. Ask clarifying questions
3. Document answers in Clarifications section
4. Update requirements with resolved details
5. Update acceptance scenarios with clarifications

### Exit Criteria (Quality Gate)
- [ ] Zero [NEEDS CLARIFICATION] markers remain
- [ ] All ambiguities documented in Clarifications section
- [ ] Requirements updated with clarified details
- [ ] Spec review checklist passed

### Artifacts
- Updated `spec.md` with Clarifications section

### Time Estimate
Few questions: 15-30 minutes  
Many ambiguities: 1-2 hours

---

## Stage 3: Planning (/plan)

### Purpose
Define HOW the feature will be implemented with technical design decisions.

### Entry Criteria
- spec.md complete and clarified
- No [NEEDS CLARIFICATION] remaining

### Activities
1. Run planning workflow
2. Identify technology stack
3. Make architecture decisions
4. Design data model
5. Define command contracts
6. Create manual test scenarios
7. Check constitution compliance

### Exit Criteria (Quality Gate)
- [ ] plan.md complete with technical context
- [ ] Constitution check passed (all 5 principles)
- [ ] research.md documents all technology choices
- [ ] data-model.md defines entities and storage
- [ ] contracts/ directory has specs for all commands/APIs
- [ ] quickstart.md has 5-10 manual test scenarios

### Artifacts
- `plan.md`
- `research.md`
- `data-model.md`
- `contracts/*.md`
- `quickstart.md`

### Time Estimate
Small feature: 1-2 hours  
Large feature: 4-6 hours

---

## Stage 4: Task Generation (/tasks)

### Purpose
Break down implementation into numbered, ordered, dependency-aware tasks.

### Entry Criteria
- plan.md complete with technical design

### Activities
1. Load design documents
2. Generate tasks from:
   - Models → model creation tasks
   - Contracts → test tasks + implementation tasks
   - Quickstart → validation tasks
3. Number tasks sequentially (T001, T002, ...)
4. Mark parallel tasks [P]
5. Order by dependencies (setup → tests → models → core → commands → polish)

### Exit Criteria (Quality Gate)
- [ ] tasks.md complete with 30-50 tasks
- [ ] All tasks have file paths specified
- [ ] Tasks in dependency order
- [ ] Phase 3.2 (Tests) comes BEFORE Phase 3.3 (Implementation)
- [ ] Parallel tasks truly independent

### Artifacts
- `tasks.md`

### Time Estimate
30-45 minutes (mostly automated generation)

---

## Stage 5: Implementation (/implement)

### Purpose
Execute tasks to build the feature following TDD discipline.

### Entry Criteria
- tasks.md complete

### Activities

#### Phase 3.1: Setup
1. Initialize project structure
2. Install dependencies
3. Configure tools

#### Phase 3.2: Tests First (TDD) 🚨 CRITICAL
1. Write ALL test files from Phase 3.2 tasks
2. Verify tests compile: `go test -c ./...`
3. Run tests and verify they FAIL: `go test ./...`
4. Commit failing tests: `git commit -m "test: add failing tests for [feature]"`

**⛔ BLOCKING GATE**: Cannot proceed to Phase 3.3 without:
- [ ] All test files exist
- [ ] All tests compile
- [ ] All tests fail (proving no implementation exists)

#### Phase 3.3: Core Implementation
1. Implement models, core logic, commands
2. Run tests after each module
3. Watch tests turn green
4. Commit per task or small module

#### Phase 3.4-3.6: Integration, Build, Polish
1. Integration work
2. Build scripts
3. Documentation updates
4. Performance validation

### Exit Criteria (Quality Gate)
- [ ] All tasks from tasks.md complete
- [ ] `go build ./...` succeeds
- [ ] Binary builds: `go build -o <binary>`
- [ ] All tests pass: `go test ./...`
- [ ] Test coverage >70%: `go test -cover ./...`
- [ ] All quickstart scenarios validated (checkboxes in quickstart-validation.md)
- [ ] Cross-platform tested (Windows: git-bash + cmd.exe, Linux OR macOS)
- [ ] Performance goals met (measure and document)

### Artifacts
- Source code
- Test files
- Build scripts
- Binary

### Time Estimate
Small feature: 4-8 hours  
Large feature: 16-40 hours

**Commit Strategy**:
- Commit after each task or logical module
- Format: `type(scope): description`
- Types: feat, fix, test, docs, refactor, chore

---

## Stage 6: Self-Review

### Purpose
Developer validates their own work before requesting review.

### Entry Criteria
- Implementation complete
- All tests passing

### Activities
1. Review against `.specify/checklists/self-review-checklist.md`
2. Verify constitution compliance
3. Check documentation completeness
4. Validate cross-platform
5. Run performance tests

### Exit Criteria (Quality Gate)
- [ ] Self-review checklist 100% complete
- [ ] No known issues or all documented
- [ ] Code comments added for complex logic
- [ ] No TODOs without linked issues

### Artifacts
- Completed self-review checklist

### Time Estimate
30-60 minutes

---

## Stage 7: Code Review & Approval

### Purpose
Peer or maintainer validates quality and adherence to standards.

### Entry Criteria
- Pull request created
- Self-review complete

### Activities
1. Reviewer uses `.specify/checklists/code-review-checklist.md`
2. Verify all PR template items checked
3. Review code quality and architecture
4. Verify tests exist and pass
5. Check constitution compliance
6. Validate documentation

### Exit Criteria (Quality Gate)
- [ ] Code review checklist complete
- [ ] All reviewer concerns addressed
- [ ] Approval granted

### Artifacts
- Review comments
- Approval

### Time Estimate
30-90 minutes per reviewer

---

## Stage 8: Merge & Release

### Purpose
Integrate feature into master and create release if needed.

### Entry Criteria
- Code review approved
- All quality gates passed

### Activities
1. Update CHANGELOG.md
2. Ensure branch is up-to-date with master
3. Resolve any conflicts
4. Merge to master
5. Tag release (if applicable)
6. Generate release artifacts
7. Archive feature branch

### Merge Criteria (FINAL GATE)
- [ ] All stages 1-7 complete
- [ ] No merge conflicts
- [ ] CI passes (if exists)
- [ ] Binary builds successfully
- [ ] Tests pass on master after merge

### Artifacts
- Merge commit
- Release tag (optional)
- Binary artifacts

---

## Definition of Done

A feature is considered "Done" when:

### Functional
✅ All functional requirements from spec.md implemented  
✅ Binary builds without errors  
✅ All commands work as specified  

### Testing (NON-NEGOTIABLE)
✅ Unit tests written for core logic (>70% coverage)  
✅ Integration tests written for all user-facing features  
✅ All tests passing: `go test ./...`  
✅ Manual quickstart scenarios validated and documented  
✅ Cross-platform tested (minimum 2 shells on Windows)  

### Documentation
✅ README.md updated if user-facing changes  
✅ Help text added for new commands  
✅ contracts/ specs exist for new commands  
✅ CHANGELOG.md entry added  

### Quality
✅ Constitution compliance verified (all 5 principles)  
✅ Performance goals met (document measurements)  
✅ No regressions in existing features  
✅ Code reviewed and approved  

### Process
✅ All tasks from tasks.md complete  
✅ Self-review checklist complete  
✅ PR template checklist complete  
✅ Commit messages follow convention  

---

## Quality Gate Enforcement

### Soft Gates (Can be waived with justification)
- CI/CD automation (T039)
- Performance benchmarks (T041)
- Additional polish tasks

### Hard Gates (CANNOT be waived)
🚫 **Phase 3.2 → Phase 3.3**: Tests MUST exist and FAIL first  
🚫 **Merge**: Tests MUST pass  
🚫 **Merge**: Constitution MUST be complied with  
🚫 **Merge**: Quickstart scenarios MUST be validated  

### Enforcement Mechanism
Add to `.specify/scripts/bash/check-merge-readiness.sh`:
```bash
#!/bin/bash
# Check if feature is ready to merge

# 1. Check tests exist
test_count=$(find tests -name "*_test.go" 2>/dev/null | wc -l)
if [ "$test_count" -eq 0 ]; then
    echo "ERROR: No test files found. TDD required."
    exit 1
fi

# 2. Check tests pass
if ! go test ./...; then
    echo "ERROR: Tests failing. Fix before merge."
    exit 1
fi

# 3. Check quickstart validation
if [ ! -f "specs/###-feature/quickstart-validation.md" ]; then
    echo "ERROR: Quickstart scenarios not validated."
    exit 1
fi

echo "✅ Merge readiness: PASS"
```

---

## Exception Process

Sometimes exceptions to the SDLC are necessary.

### When Exceptions Are Allowed
- Hotfixes for critical bugs (skip spec/plan, go straight to fix)
- Documentation-only changes
- Dependency updates with no functional changes

### Exception Approval
- Document in PR: "Exception: [reason]"
- Maintainer must explicitly approve
- Create follow-up issue if technical debt incurred

### Example Exception from Feature 001
```
Exception: Tests deferred to feature 003
Reason: Core functionality working, tests can be added incrementally
Approval: Self (maintainer decision)
Follow-up: Issue #3 - Add test coverage for feature 001
```

---

## Commit Message Convention

### Format
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types
- **feat**: New feature
- **fix**: Bug fix
- **test**: Adding missing tests
- **docs**: Documentation only
- **refactor**: Code restructuring
- **chore**: Build, dependencies, tooling
- **perf**: Performance improvements

### Examples from Feature 001
```
feat(cli): implement start command with duplicate handling
test(integration): add start command contract tests
docs(readme): add usage examples and quick start guide
chore(build): add cross-platform build script
```

---

## Release Process

### Versioning (Semantic Versioning)

**MAJOR.MINOR.PATCH** (e.g., 1.0.0)

- **MAJOR**: Breaking CLI changes (command removal, flag changes, output format breaking)
- **MINOR**: New commands, new optional flags, backward-compatible additions
- **PATCH**: Bug fixes, documentation, internal refactoring

### Release Checklist
- [ ] All features for release merged
- [ ] CHANGELOG.md updated
- [ ] Version bumped in `cmd/my-context/main.go`
- [ ] Tag created: `git tag v1.0.0`
- [ ] Binaries built for all platforms
- [ ] Release notes written
- [ ] Artifacts uploaded

---

## Metrics & Tracking

### Feature Health Metrics
- **Functional Coverage**: Requirements delivered / Requirements specified
- **Task Completion**: Tasks done / Total tasks
- **Test Coverage**: Lines covered / Total lines
- **Constitution Compliance**: Principles satisfied / 5

### Project Health Indicators
- Time from /specify to merge
- Defect rate post-merge
- Technical debt items open
- Test coverage trend

### Feature 001 Baseline
- Functional: 100% (35/35)
- Tasks: 78.6% (33/42)
- Tests: 0% (0 test files)
- Constitution: 90% (4.5/5 - cross-shell gap)

**Target for Future Features**: 
- Functional: 100%
- Tasks: 95%+ (allow 1-2 polish tasks to defer)
- Tests: >70% coverage
- Constitution: 100%

---

## Workflow Diagram

```
┌─────────────┐
│   /specify  │ Create spec.md
└──────┬──────┘
       │
       ├─[NEEDS CLARIFICATION]?
       │   YES↓              NO↓
       │  ┌──────────┐       │
       │  │ /clarify │       │
       │  └────┬─────┘       │
       │       ↓             │
       └───────┴─────────────┘
               ↓
       ┌───────────┐
       │   /plan   │ Create plan.md, research.md, data-model.md, contracts/
       └─────┬─────┘
             │
       ┌─────┴──────┐
       │Constitution│ GATE: Must pass
       │   Check    │
       └─────┬──────┘
             ↓
       ┌───────────┐
       │  /tasks   │ Create tasks.md
       └─────┬─────┘
             ↓
       ┌──────────────┐
       │  /implement  │ Execute tasks
       └──────┬───────┘
              │
       ┌──────┴─────────┐
       │  Phase 3.2:    │ 🚨 BLOCKING GATE
       │  Write Tests   │ Tests MUST exist and FAIL
       │  (TDD)         │
       └──────┬─────────┘
              ↓
       ┌──────┴─────────┐
       │  Phase 3.3:    │ Make tests pass
       │  Implement     │
       └──────┬─────────┘
              ↓
       ┌──────┴─────────┐
       │  Phase 3.4-3.6 │ Integration, build, polish
       └──────┬─────────┘
              │
       ┌──────┴──────┐
       │ Self-Review │ Complete checklist
       └──────┬──────┘
              ↓
       ┌──────┴───────┐
       │ Code Review  │ Peer approval
       └──────┬───────┘
              │
       ┌──────┴──────┐
       │   Merge to  │ Feature complete
       │   master    │
       └─────────────┘
```

---

## TDD Enforcement (Critical)

### The Problem (Feature 001)
Tests were specified in Phase 3.2 but NOT written. Implementation proceeded anyway, violating TDD principle.

### The Solution (Feature 002+)

**BLOCKING GATE at Phase 3.2 → 3.3 transition**:

```bash
#!/bin/bash
# .specify/scripts/bash/verify-tdd-gate.sh

FEATURE_DIR=$1

# Count test files
test_files=$(find tests -name "*_test.go" 2>/dev/null | wc -l)

if [ "$test_files" -eq 0 ]; then
    echo "❌ TDD GATE FAILED: No test files found"
    echo ""
    echo "Phase 3.2 requires tests to be written FIRST."
    echo "Cannot proceed to Phase 3.3 until tests exist and fail."
    echo ""
    echo "Required: Create test files from tasks.md Phase 3.2"
    exit 1
fi

# Run tests (should fail initially)
if go test ./... >/dev/null 2>&1; then
    echo "⚠️  WARNING: All tests passing. Expected failures in TDD."
fi

echo "✅ TDD GATE PASSED: $test_files test files found"
exit 0
```

**Usage**:
```bash
# Before starting Phase 3.3
bash .specify/scripts/bash/verify-tdd-gate.sh

# Only proceed if exit code 0
```

---

## Cross-Platform Testing Requirements

### Minimum Platform Matrix

**Windows**:
- ✅ git-bash (primary development shell)
- ✅ cmd.exe (test fallback compatibility)

**Linux OR macOS**:
- ✅ One of: Ubuntu, macOS, or WSL

### Testing Checklist

For each platform:
- [ ] Binary runs: `./my-context --version`
- [ ] Commands work: `./my-context start "Test"`
- [ ] Paths normalized correctly
- [ ] State persists across sessions
- [ ] JSON output valid

### Documentation

Create `specs/###-feature/cross-platform-validation.md`:
```markdown
## git-bash (Windows)
- [x] Binary runs
- [x] Start/stop working
- [x] Paths normalized
- [x] State persists
- [x] JSON valid

## cmd.exe (Windows)
- [ ] Binary runs
- [ ] Start/stop working
- [ ] Paths normalized
...
```

---

## Merge Decision Framework

### Can Merge If:
✅ All functional requirements delivered  
✅ Tests written and passing  
✅ Documentation complete  
✅ Constitution compliant  
✅ Quickstart scenarios validated  
✅ Cross-platform tested (min 2 shells)  
✅ Code reviewed and approved  

### Can Merge With Exceptions If:
⚠️ Minor polish tasks deferred (T041, T042)  
⚠️ CI/CD deferred (T039) with follow-up issue  
⚠️ Performance "good enough" but not benchmarked  

**Requires**: 
- Explicit exception approval
- Follow-up issue created
- Technical debt documented

### Cannot Merge If:
❌ Core functionality incomplete  
❌ Tests missing or failing  
❌ Constitution violated without justification  
❌ Breaking changes undocumented  
❌ Major bugs unresolved  

### Feature 001 Decision
✅ **MERGED** with exception:
- Exception: Tests deferred (T004-T013)
- Justification: Core functionality complete, tests can be added incrementally
- Follow-up: Feature 003 - Add test coverage
- Approved: Maintainer decision (2025-10-05)

---

## Bug Fix & Hotfix Process

### Standard Bug Fix
1. Create issue: "Bug: [description]"
2. Create branch: `bugfix-###-description`
3. Write failing test demonstrating bug
4. Fix bug
5. Verify test passes
6. Submit PR with test + fix

### Critical Hotfix (Production Down)
1. Create branch from master: `hotfix-description`
2. Implement minimal fix
3. Test manually
4. Merge immediately
5. Create follow-up issue for proper test coverage

---

## Continuous Improvement

### After Each Feature
1. Update gap-analysis.md
2. Document lessons learned
3. Propose SDLC improvements
4. Update checklists if needed

### Quarterly Review
- Review all features
- Calculate health metrics
- Identify patterns (good and bad)
- Refine process

---

## Tools & Automation

### Existing Scripts
- `create-new-feature.sh` - Start new feature
- `setup-plan.sh` - Initialize planning
- `check-prerequisites.sh` - Validate environment
- `update-agent-context.sh` - Update AI context

### Needed Scripts (Feature 002)
- `verify-tdd-gate.sh` - Block Phase 3.3 without tests
- `check-merge-readiness.sh` - Validate merge criteria
- `run-quickstart.sh` - Automated quickstart validation

---

## Summary

This SDLC provides:
1. ✅ Clear stages with entry/exit criteria
2. ✅ Mandatory quality gates (especially TDD)
3. ✅ Definition of Done
4. ✅ Exception process
5. ✅ Metrics for tracking
6. ✅ Continuous improvement

**Feature 001 taught us**: Process must have enforcement mechanisms, not just suggestions.

**Feature 002+ benefit**: Consistent quality, predictable delivery, less technical debt.

---

*Version 1.0.0 - Based on lessons from feature 001 implementation*

