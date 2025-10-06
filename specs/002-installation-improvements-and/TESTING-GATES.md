# Testing Gates & Sign-Offs

**Sprint**: Sprint 2 (002-installation-improvements-and)
**Purpose**: Define clear checkpoints, verification criteria, and sign-off requirements throughout SDLC
**Constitutional Authority**: Test-First Development (constitution.md:100-107)

---

## Overview

This document enforces **proof of testing** at critical SDLC stages. Each gate requires:
1. **Specific actions** (commands to run, artifacts to create)
2. **Pass/fail criteria** (measurable outcomes)
3. **Evidence artifacts** (logs, screenshots, recordings)
4. **Sign-off** (who approves, what they verify)

---

## Gate 1: TDD Verification (BLOCKING)

**Location**: After Phase 3.2 (T005-T012 complete) → Before Phase 3.3 (T013+)

**Trigger**: All test files created and marked `[x]` in tasks.md

**Required Actions**:
```bash
# 1. Compile all tests (verify syntax)
go test ./tests/... -run=^$ 2>&1 | tee test-compile-check.log

# 2. Run full test suite
go test ./tests/integration/... -v 2>&1 | tee test-suite-pre-impl.log
go test ./tests/unit/... -v 2>&1 | tee test-suite-unit-pre-impl.log

# 3. Verify 100% FAIL rate
grep -c "PASS" test-suite-pre-impl.log  # Should be 0
grep -c "FAIL" test-suite-pre-impl.log  # Should be >0
```

**Pass Criteria**:
- ✅ All tests compile without syntax errors
- ✅ 0 tests PASS (no implementation exists)
- ✅ All tests FAIL with meaningful error messages (not panics)
- ✅ Test names match contract specifications

**Evidence Required**:
- `specs/002-.../validation-evidence/test-suite-pre-impl.log`
- Summary: `X tests FAIL, 0 tests PASS`

**Sign-Off**:
- [ ] Developer confirms tests written
- [ ] Reviewer confirms tests FAIL for correct reasons (not syntax errors)
- [ ] Both parties verify no implementation code exists

**Next Step After Gate 1 Passes**:
→ Proceed to **Phase 3.3 Task T013** (Model Layer)
→ Commit tests: `git commit -m "test: add Sprint 2 integration tests (TDD - all failing)"`

**If Gate 1 Fails** (tests pass or don't compile):
→ HALT implementation
→ Fix tests or remove accidental implementation
→ Re-run Gate 1

---

## Gate 2: Feature-Level Validation (PER P1 FEATURE)

### Feature 2A: Project Filter (User Request #1)

**Location**: After T014 (project.go), T025 (start --project) complete

**Required Actions**:
```bash
# 1. Run project filter tests
go test ./tests/integration/project_filter_test.go -v

# 2. Run unit tests
go test ./tests/unit/project_parser_test.go -v

# 3. Manual smoke test with real data
my-context start "Test Phase 1" --project test-proj
my-context start "Test Phase 2" --project test-proj
my-context start "Other Work" --project other
my-context list --project test-proj  # Should show only 2 contexts
```

**Pass Criteria**:
- ✅ T009 (project_filter_test.go) all tests PASS
- ✅ T011 (project_parser_test.go) all tests PASS
- ✅ Manual test shows correct filtering
- ✅ Edge cases work (no colon, multiple colons, case-insensitive)

**Evidence Required**:
- Test logs: `project-filter-validation.log`
- Screenshot: Manual test showing filtered list
- Edge case matrix:
  | Input | Expected Project | Actual | Pass? |
  |-------|------------------|--------|-------|
  | "ps-cli: Phase 1" | "ps-cli" | ... | ✅ |
  | "Standalone" | "Standalone" | ... | ✅ |
  | "a: b: c" | "a" | ... | ✅ |

**User Acceptance** (P1 Feature):
- [ ] Demo to original requester (show contexts/phase-1.md:241-244)
- [ ] User confirms it solves stated problem
- [ ] User validates "low implementation effort, high organization value"

**Sign-Off**:
- [ ] Implementer verifies tests pass
- [ ] Reviewer spot-checks edge cases
- [ ] **Original requester** confirms feature meets need

**Next Step**: → Proceed to Feature 2B (Export) or continue with other T013-T020 tasks

---

### Feature 2B: Export Command (User Request #2)

**Location**: After T015 (markdown.go), T021 (export.go) complete

**Required Actions**:
```bash
# 1. Run export tests
go test ./tests/integration/export_test.go -v

# 2. Manual smoke test
my-context start "Export Test Context"
my-context note "Test note with $pecial ch@rs"
my-context file ~/test.txt
my-context stop
my-context export "Export Test Context"  # Default path
cat Export_Test_Context.md  # Verify format

# 3. Test custom path and --all flag
my-context export "Export Test Context" --to docs/test-export.md
my-context export --all --to exports/
```

**Pass Criteria**:
- ✅ T005 (export_test.go) all tests PASS
- ✅ Markdown renders correctly in GitHub/VS Code
- ✅ All data present (notes, files, timestamps, duration)
- ✅ Overwrite prompt works
- ✅ --to flag creates parent directories

**Evidence Required**:
- Test log: `export-validation.log`
- Sample export: `sample-export.md` (rendered markdown)
- Screenshot: Markdown rendered in GitHub

**User Acceptance** (P1 Feature):
- [ ] Demo export workflow to original requester
- [ ] User confirms it "automates manual process" (their words)
- [ ] Verify exported markdown is "easy sharing" format

**Sign-Off**:
- [ ] Implementer verifies markdown format correct
- [ ] Reviewer checks edge cases (empty context, special chars)
- [ ] **Original requester** validates automation value

**Next Step**: → Continue to Archive/Delete features (T016, T022, T023)

---

## Gate 3: Cross-Platform Validation (MANDATORY)

**Location**: After all implementation (T005-T026) complete, before T035

**Required Actions**:
```bash
# Platform Matrix Test
# For EACH platform: Windows, Linux (WSL), macOS (if available)

# 1. Build for platform
./scripts/build-all.sh

# 2. Install
./scripts/install.sh  # Linux/macOS/WSL
# OR scripts/install.bat  # Windows cmd
# OR scripts/install.ps1  # Windows PowerShell

# 3. Run core workflow
my-context start "Cross-Platform Test" --project validation
my-context note "Testing on [platform]"
my-context list --project validation
my-context export "Cross-Platform Test"
my-context stop

# 4. Run full test suite ON THAT PLATFORM
go test ./tests/... -v
```

**Pass Criteria**:
- ✅ Binary runs on platform without errors
- ✅ Installation completes (PATH updated)
- ✅ Core commands work identically across platforms
- ✅ Test suite passes on platform
- ✅ Paths normalized correctly (Windows backslash vs POSIX)

**Evidence Required**:
```
Cross-Platform Test Matrix:
| Platform | Binary Works | Install Works | Commands Work | Tests Pass | Notes |
|----------|--------------|---------------|---------------|------------|-------|
| Windows 10 (cmd) | ✅ | ✅ | ✅ | ✅ | PATH issue resolved |
| Windows 10 (PS) | ✅ | ✅ | ✅ | ✅ | - |
| WSL (Debian) | ✅ | ✅ | ✅ | ✅ | User isolation working |
| macOS (Intel) | ⚠️ N/A | - | - | - | No test machine |
| macOS (ARM) | ⚠️ N/A | - | - | - | No test machine |
```

**Sign-Off**:
- [ ] Tested on Windows (both cmd and PowerShell)
- [ ] Tested on Linux/WSL
- [ ] Documented macOS status (tested or not available)

**Next Step**: → Proceed to Gate 4 (Integration Validation)

---

## Gate 4: Integration & Performance Validation

**Location**: Phase 3.6 (T035-T041)

**4A: Integration Tests (T035-T036)**

**Actions**:
```bash
# Run on Linux
go test ./tests/integration/... -v -timeout 30s | tee integration-linux.log

# Run on Windows
go test ./tests/integration/... -v -timeout 30s | tee integration-windows.log
```

**Pass Criteria**:
- ✅ All integration tests PASS on both platforms
- ✅ No platform-specific failures
- ✅ Test execution time <30s total

---

**4B: Quickstart Scenarios (T037)**

**Actions**: Execute all 9 scenarios from `quickstart.md`

**Evidence Required**: Checklist completion

```markdown
### Quickstart Validation Checklist

- [ ] Scenario 1: Multi-platform installation (WSL)
  - Binary downloaded and runs
  - install.sh completes without sudo
  - my-context --version shows correct info

- [ ] Scenario 2: Project-based workflow
  - `start --project` creates correct name
  - `list --project` filters correctly

- [ ] Scenario 3: Export and share
  - Export creates markdown file
  - File renders in GitHub/VS Code

- [ ] Scenario 4: Context lifecycle (archive/delete)
  - Archive hides from default list
  - Delete prompts for confirmation
  - Transitions.log preserved

- [ ] Scenario 5: List enhancements
  - Default shows 10 limit
  - --search works
  - --all shows everything

- [ ] Scenario 6: Bug fixes validation
  - $ character preserved in notes
  - History shows "(none)" not "NULL"

- [ ] Scenario 7: Cross-platform (Windows)
  - install.bat works (or install.ps1)
  - Commands identical to Linux

- [ ] Scenario 8: Backward compatibility
  - Sprint 1 contexts still work
  - New features work on old contexts

- [ ] Scenario 9: JSON output
  - All commands support --json
  - JSON is valid (jq parses)
```

**Sign-Off**: [ ] All 9 scenarios executed successfully

---

**4C: Performance Benchmarks (T038)**

**Actions**:
```bash
# Create test data (1000 contexts)
for i in {1..1000}; do
  my-context start "Context $i"
  my-context note "Note $i"
  my-context stop
  sleep 0.1
done

# Benchmark list command
time my-context list --all  # Target: <1s

# Benchmark export
# (Create context with 500 notes first)
time my-context export "Large Context"  # Target: <1s

# Benchmark search
time my-context list --search "test"  # Target: <1s
```

**Pass Criteria**:
- ✅ List 1000 contexts: <1 second
- ✅ Export 500 notes: <1 second
- ✅ Search 1000 contexts: <1 second

**Evidence**: `performance-benchmarks.txt` with timing data

---

## Gate 5: User Acceptance (P1 Features Only)

**Location**: After all features complete, before declaring sprint done

**For Each P1 Feature**:

### Project Filter Acceptance
**Requester**: [Original user who provided feedback]
**Acceptance Test**:
```bash
# User runs with their real data
my-context list --project ps-cli-retrofit
# Expected: Shows only ps-cli contexts

my-context start "New Phase" --project ps-cli-retrofit
# Expected: Creates "ps-cli-retrofit: New Phase"
```

**User Feedback Template**:
```
Feature: Project Filter
- Does it solve your stated problem? [ ] Yes [ ] No
- Is the syntax intuitive? [ ] Yes [ ] No
- Does it match your mental model? [ ] Yes [ ] No
- Would you use this daily? [ ] Yes [ ] No
- Any issues or suggestions? ___________________

Sign-Off: ________________ Date: _______
```

### Export Command Acceptance
**Requester**: [Same user]
**Acceptance Test**:
```bash
# User exports their real context
my-context export "ps-cli-retrofit: Phase 1" --to contexts/phase-1.md
# User opens file and reviews

# User tests sharing workflow
cat contexts/phase-1.md  # Is this shareable as-is?
```

**User Feedback**: Same template

---

## Gate 6: Constitution Compliance (Final)

**Location**: T041 - Before sprint close

**Review Checklist**:
```markdown
### Principle I: Unix Philosophy
- [ ] Each command does one thing well
- [ ] Text I/O maintained
- [ ] Composable with shell tools (test: `my-context list | grep project`)

### Principle II: Cross-Platform Compatibility
- [ ] Tested on Windows (cmd + PowerShell)
- [ ] Tested on Linux/WSL
- [ ] Path handling works (POSIX ↔ Windows)

### Principle III: Stateful Context Management
- [ ] One active context enforced
- [ ] Automatic context stopping works
- [ ] State persists correctly

### Principle IV: Minimal Surface Area
- [ ] Added 3 commands (export, archive, delete) - justified
- [ ] Total command count: 11 (acceptable)
- [ ] All have single-letter aliases
- [ ] No config files required

### Principle V: Data Portability
- [ ] All data still plain text
- [ ] Can grep/cat context files
- [ ] Export generates standard markdown

### Principle VI: User-Driven Design ⭐ NEW
- [ ] Project filter addresses observed pattern
- [ ] Export automates documented manual workflow
- [ ] Archive/delete solve real problem (#7)
- [ ] No speculative features added
- [ ] User acceptance obtained for P1 features
```

**Evidence**: `constitution-compliance-checklist.md`

**Sign-Off**: [ ] Tech lead confirms all 6 principles upheld

---

## Clear Next-Step Transitions

### Current State → Next Action Map

**RIGHT NOW** (Phase 3.1 complete):
```
✅ Phase 3.1 done (T001-T004, plus T027, T031)
   ↓
📍 YOU ARE HERE
   ↓
NEXT: Gate 1 Pre-Check
   ↓
Action: Verify test files exist
Command: ls tests/integration/*.go tests/unit/*.go
Expected: 7 test files (export, archive, delete, list_enhanced, project_filter, bug_fixes, backward_compat)
   ↓
If files exist → Mark T005-T012 as [x] in tasks.md
If files missing → Write tests (follow IMPLEMENTATION-GUIDE.md Phase 3.2)
```

---

**After T005-T012 marked complete**:
```
✅ Phase 3.2 done (all tests written)
   ↓
🚨 GATE 1: TDD Verification
   ↓
Action: Run test suite, verify 100% FAIL
Command: go test ./tests/... -v 2>&1 | tee test-suite-pre-impl.log
Expected: All tests FAIL, 0 PASS
   ↓
Create evidence: specs/002-.../validation-evidence/test-suite-pre-impl.log
   ↓
Sign-off: Developer + Reviewer confirm TDD order preserved
   ↓
✅ GATE 1 PASSED → Proceed to Phase 3.3
   ↓
NEXT: Task T013 (Add IsArchived to Context model)
```

---

**After T013-T015 complete** (Model layer):
```
✅ Model layer done
   ↓
NEXT: Task T016 (Core context operations: Archive, Delete, Export)
Note: This is SEQUENTIAL - T017 depends on T016
```

---

**After T016-T020 complete** (Core logic):
```
✅ Core logic done
   ↓
OPTIONAL CHECKPOINT: Run tests for completed features
Command: go test ./tests/integration/archive_test.go -v
Command: go test ./tests/integration/delete_test.go -v
Expected: Tests should start PASSING now
   ↓
If tests still FAIL → Debug implementation
If tests PASS → Continue to command layer
   ↓
NEXT: Task T021 (Export command) - can run in parallel with T022-T026
```

---

**After T021-T026 complete** (Commands):
```
✅ All commands implemented
   ↓
🚨 GATE 2: Feature-Level Validation
   ↓
Action: Validate P1 features with original requester
   ↓
For Project Filter:
  - User tests with their data
  - User provides feedback
  - Sign-off obtained
   ↓
For Export:
  - User tests export workflow
  - User confirms automation value
  - Sign-off obtained
   ↓
✅ GATE 2 PASSED → Proceed to Phase 3.4
   ↓
NEXT: Task T028 (install.bat) - T027/T031 already done
```

---

**After T028-T031 complete** (Installation scripts):
```
✅ Installation scripts done
   ↓
🚨 GATE 3: Cross-Platform Validation
   ↓
Action: Test on each available platform
   ↓
Create: Cross-platform test matrix (see Gate 3 above)
   ↓
If all platforms PASS → GATE 3 PASSED
If any platform FAILS → Fix platform-specific issues
   ↓
NEXT: Phase 3.5 (Documentation)
```

---

**After T032-T034 complete** (Documentation):
```
✅ Documentation updated
   ↓
NEXT: Phase 3.6 (Integration & Validation)
   ↓
Task T035: Run integration tests on Linux
Task T036: Run integration tests on Windows
Task T037: Execute quickstart scenarios (all 9)
Task T038: Performance benchmarks
Task T039: Binary size check (<10MB each)
Task T040: Cross-platform smoke test
Task T041: Constitution compliance check
```

---

**After T035-T041 complete**:
```
✅ All validation tasks done
   ↓
🚨 FINAL GATES: 4, 5, 6
   ↓
Gate 4: Integration validation ✅
Gate 5: User acceptance (P1) ✅
Gate 6: Constitution compliance ✅
   ↓
All evidence collected in validation-evidence/
   ↓
SPRINT 2 COMPLETE ✅
   ↓
NEXT: Create PR to main branch
```

---

## Evidence Artifacts Directory Structure

```
specs/002-installation-improvements-and/
├── validation-evidence/
│   ├── gate-1-tdd/
│   │   ├── test-suite-pre-impl.log
│   │   └── test-compile-check.log
│   ├── gate-2-features/
│   │   ├── project-filter-validation.log
│   │   ├── project-filter-edge-cases.md
│   │   ├── export-validation.log
│   │   ├── sample-export.md
│   │   └── user-acceptance-forms.md
│   ├── gate-3-cross-platform/
│   │   ├── windows-cmd-test.log
│   │   ├── windows-ps-test.log
│   │   ├── wsl-debian-test.log
│   │   └── cross-platform-matrix.md
│   ├── gate-4-integration/
│   │   ├── integration-linux.log
│   │   ├── integration-windows.log
│   │   ├── quickstart-checklist.md
│   │   └── quickstart-demo.gif (optional)
│   ├── gate-5-performance/
│   │   └── performance-benchmarks.txt
│   └── gate-6-constitution/
│       └── compliance-checklist.md
```

---

## Enforcement Mechanism

**How to ensure gates are followed**:

1. **Git hooks** (optional): Pre-commit hook checks for test evidence before allowing commits to certain branches

2. **PR template**: Require evidence links in PR description
   ```
   ## Testing Evidence
   - [ ] Gate 1 (TDD): [link to test-suite-pre-impl.log]
   - [ ] Gate 2 (Features): [link to user-acceptance-forms.md]
   - [ ] Gate 3 (Cross-platform): [link to cross-platform-matrix.md]
   - [ ] Gate 4-6: [links to remaining evidence]
   ```

3. **tasks.md updates**: Cannot mark task `[x]` without corresponding evidence

4. **Sprint review**: Present evidence artifacts during retrospective

---

## Failure Handling

**If Gate 1 fails** (tests pass before implementation):
1. HALT all implementation work immediately
2. Identify which code was written early
3. Options:
   - A. Remove implementation, rewrite tests to fail
   - B. Accept TDD violation, document in retrospective as lesson learned
4. Do NOT proceed to next phase without resolving

**If Gate 2 fails** (user rejects feature):
1. Gather detailed feedback
2. Options:
   - A. Iterate on feature (within sprint if minor)
   - B. Defer to Sprint 3 for redesign
   - C. Remove feature if fundamentally wrong
3. Update USER-REQUEST-VALIDATION.md with failure learnings

**If Gate 3 fails** (platform doesn't work):
1. Document broken platform in cross-platform-matrix.md
2. Fix if blocker (e.g., Windows critical)
3. Defer if edge case (e.g., macOS arm64 unavailable for testing)
4. Update TROUBLESHOOTING.md with known issues

---

## Integration with Existing Documents

**SDLC.md** (constitution.md:100-107):
- Already has TDD blocking gate
- Should reference this document for details

**tasks.md**:
- Add note after Phase 3.2: "⚠️ GATE 1 CHECKPOINT - See TESTING-GATES.md"
- Add note after T026: "⚠️ GATE 2 CHECKPOINT - Validate P1 features"
- Add note before T035: "⚠️ GATE 3 CHECKPOINT - Cross-platform validation"

**IMPLEMENTATION-GUIDE.md**:
- Add "Checkpoint" sections at gate locations
- Link to this document for full criteria

---

## Success Metrics

**Sprint 2 is NOT complete until**:
- [ ] All 6 gates passed
- [ ] All evidence artifacts created
- [ ] All sign-offs obtained
- [ ] validation-evidence/ directory committed to repo

**This ensures**:
- TDD order preserved (prevents technical debt)
- User needs validated (Principle VI compliance)
- Cross-platform quality (Principle II compliance)
- Constitution principles upheld (governance)

---

**Last Updated**: 2025-10-05
**Version**: 1.0 (initial draft for Sprint 2)
