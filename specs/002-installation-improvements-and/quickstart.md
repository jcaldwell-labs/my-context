# Quickstart: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Date**: 2025-10-09  
**Purpose**: End-to-end validation scenarios for Sprint 2 features

## Prerequisites

- Go 1.21+ installed
- Git repository cloned
- Terminal access (bash/zsh for Unix, cmd.exe or PowerShell for Windows)

---

## Scenario 1: Multi-Platform Installation (WSL User)

**Objective**: Verify WSL user can install pre-built binary without building from source

**Steps**:
1. Navigate to GitHub Releases page
2. Download `my-context-linux-amd64`
3. Download `my-context-linux-amd64.sha256`
4. Verify checksum: `sha256sum -c my-context-linux-amd64.sha256`
5. Run `./install.sh`
6. Verify installation: `my-context --version`
7. Verify PATH: `which my-context`

**Expected Results**:
- Checksum verification passes
- Binary installed to `~/.local/bin/my-context`
- PATH includes `~/.local/bin`
- Version displays: "my-context version v2.0.0"
- Existing ~/.my-context/ data preserved (if upgrading)

**Acceptance Criteria**: FR-001, FR-002, FR-003

---

## Scenario 2: Project-Based Workflow

**Objective**: Verify project filtering and organization

**Steps**:
1. Start contexts with project prefixes:
   ```bash
   my-context start "Phase 1" --project ps-cli
   my-context note "Started Phase 1 work"
   my-context stop
   
   my-context start "Phase 2" --project ps-cli
   my-context note "Continuing with Phase 2"
   my-context stop
   
   my-context start "Planning" --project garden
   my-context stop
   ```

2. Filter by project:
   ```bash
   my-context list --project ps-cli
   ```

3. Verify case-insensitive matching:
   ```bash
   my-context list --project PS-CLI
   my-context list --project Ps-Cli
   ```

**Expected Results**:
- Contexts created as: "ps-cli: Phase 1", "ps-cli: Phase 2", "garden: Planning"
- `list --project ps-cli` shows only 2 ps-cli contexts
- Case variations all return same results
- garden context not shown in ps-cli filter

**Acceptance Criteria**: FR-004

---

## Scenario 3: Export and Share

**Objective**: Verify markdown export functionality

**Steps**:
1. Export single context:
   ```bash
   my-context export "ps-cli: Phase 1"
   ```

2. Verify output file exists: `ls ps-cli-Phase-1.md`

3. Export to custom path:
   ```bash
   my-context export "ps-cli: Phase 2" --to reports/phase2.md
   ```

4. Verify overwrite prompt:
   ```bash
   my-context export "ps-cli: Phase 1"  # Should prompt
   n  # Decline
   echo $?  # Should be 2
   ```

5. Force overwrite:
   ```bash
   my-context export "ps-cli: Phase 1" --force
   ```

6. Verify markdown content:
   - Open exported file in text editor
   - Verify headers, notes, files, timestamps present
   - Verify timestamps are in local timezone

**Expected Results**:
- Export succeeds with clear success message
- Markdown file is valid and human-readable
- Custom paths create parent directories
- Overwrite protection works (prompts user)
- --force bypasses prompt
- Timestamps in local timezone (not UTC)

**Acceptance Criteria**: FR-005

---

## Scenario 4: Context Lifecycle (Archive & Delete)

**Objective**: Verify archive and delete operations

**Steps**:
1. Create and complete a context:
   ```bash
   my-context start "Test Context"
   my-context note "This is a test"
   my-context stop
   ```

2. Archive completed context:
   ```bash
   my-context archive "Test Context"
   ```

3. Verify hidden from default list:
   ```bash
   my-context list  # Should not show "Test Context"
   my-context list --archived  # Should show "Test Context"
   ```

4. Try to archive active context (should fail):
   ```bash
   my-context start "Active Context"
   my-context archive "Active Context"  # Should error
   my-context stop
   ```

5. Delete a context:
   ```bash
   my-context delete "Test Context"
   y  # Confirm
   ```

6. Verify context removed:
   ```bash
   my-context list --archived  # Should not show "Test Context"
   ls ~/.my-context/  # Directory should not exist
   ```

7. Verify transitions.log preserved:
   ```bash
   grep "Test Context" ~/.my-context/transitions.log  # Should still show transitions
   ```

**Expected Results**:
- Archive marks context as archived (hidden from default list)
- Cannot archive active context
- Delete prompts for confirmation
- Delete removes context directory
- transitions.log preserves historical entries

**Acceptance Criteria**: FR-007, FR-008

---

## Scenario 5: List Enhancements (Large Dataset)

**Objective**: Verify pagination and filtering with many contexts

**Setup**: Create 50 contexts (scripted):
```bash
for i in {1..50}; do
  my-context start "Context $i"
  my-context stop
done
```

**Steps**:
1. Default list (should show 10):
   ```bash
   my-context list
   ```
   - Verify shows "Showing 10 of 50 contexts"

2. Show all:
   ```bash
   my-context list --all
   ```
   - Verify shows all 50

3. Custom limit:
   ```bash
   my-context list --limit 20
   ```
   - Verify shows 20 contexts

4. Search:
   ```bash
   my-context list --search "Context 1"
   ```
   - Verify shows Context 1, Context 10-19 (substring match)

5. Combined filters:
   ```bash
   my-context list --search "Context 1" --limit 5
   ```
   - Verify shows first 5 matches only

**Expected Results**:
- Default limit works (10 contexts)
- Truncation message displays correct counts
- --all overrides limit
- --limit accepts custom values
- --search performs case-insensitive substring matching
- Multiple filters combine with AND logic

**Acceptance Criteria**: FR-006

---

## Scenario 6: Bug Fixes Validation

**Objective**: Verify Sprint 1 bugs are fixed

**Steps**:
1. Test $ character in notes:
   ```bash
   my-context start "Budget Test"
   my-context note "Budget: $500-800"
   my-context list --active-only
   ```
   - Verify note displays with $ intact

2. Test history NULL display:
   ```bash
   my-context history
   ```
   - Verify shows "(none)" instead of "NULL" for empty context fields

**Expected Results**:
- Special characters preserved in notes
- History displays user-friendly "(none)" for empty fields

**Acceptance Criteria**: FR-009

---

## Scenario 7: Cross-Platform Installation (Windows)

**Objective**: Verify Windows installation scripts work

**Steps (cmd.exe)**:
1. Download `my-context-windows-amd64.exe`
2. Run `install.bat`
3. Verify PATH updated: `echo %PATH%`
4. Test command: `my-context --version`

**Steps (PowerShell)**:
1. Run `install.ps1`
2. Verify PATH: `$env:PATH`
3. Test command: `my-context --version`

**Expected Results**:
- Binary installed to `%USERPROFILE%\bin\`
- PATH updated (visible in new shell session)
- Command works from any directory
- No admin privileges required

**Acceptance Criteria**: FR-002

---

## Scenario 8: Backward Compatibility (Sprint 1 → Sprint 2)

**Objective**: Verify Sprint 1 data works with Sprint 2 binary

**Steps**:
1. Locate Sprint 1 contexts in ~/.my-context/
2. Install Sprint 2 binary (via install script)
3. List contexts: `my-context list`
4. Start old context: `my-context start "Old Context Name"`
5. Add note to old context: `my-context note "Sprint 2 note"`
6. Export old context: `my-context export "Old Context Name"`

**Expected Results**:
- All Sprint 1 contexts visible in Sprint 2
- Can start, stop, add notes to Sprint 1 contexts
- Export works on Sprint 1 contexts
- No data loss or corruption
- New is_archived field defaults to false for old contexts

**Acceptance Criteria**: FR-011

---

## Scenario 9: JSON Output for Scripting

**Objective**: Verify --json flag works across commands

**Steps**:
1. List as JSON:
   ```bash
   my-context list --json | jq '.contexts | length'
   ```

2. Export as JSON:
   ```bash
   my-context export "Context" --json > context.json
   ```

3. Start with JSON output:
   ```bash
   my-context start "Test" --json
   ```

**Expected Results**:
- All commands support --json flag
- Output is valid JSON (parseable by jq)
- JSON format is consistent across commands
- Enables scripting and automation

**Acceptance Criteria**: General requirement (not explicit FR)

---

## Performance Benchmarks

**Test Environment**: Create 1000 contexts with 500 notes each

**Benchmarks**:
```bash
# List performance
time my-context list --all  # Target: <1s

# Export performance
time my-context export "Large Context"  # Target: <1s (500 notes)

# Search performance
time my-context list --search "Phase"  # Target: <1s (1000 contexts)
```

**Expected Results**:
- List 1000 contexts: <1 second
- Export 500 notes: <1 second
- Search 1000 contexts: <1 second

**Acceptance Criteria**: Performance goals from plan.md

---

## Cleanup

```bash
# Remove test contexts
for i in {1..50}; do
  my-context delete "Context $i" --force
done

my-context delete "Budget Test" --force
my-context delete "Active Context" --force
```

---

**Quickstart Complete**: All 9 scenarios ready for manual validation during implementation phase.
