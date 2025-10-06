# Quickstart: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Date**: 2025-10-05  
**Purpose**: End-to-end testing scenarios for Sprint 2 features

---

## Scenario 1: Multi-Platform Installation (WSL User)

**Goal**: Install my-context on WSL without building from source

**Steps**:
```bash
# 1. Download Linux binary from releases page
wget https://github.com/user/my-context-copilot/releases/download/v1.1.0/my-context-linux-amd64
chmod +x my-context-linux-amd64

# 2. Verify binary works
./my-context-linux-amd64 --version
# Expected: my-context version 1.1.0 (build: 2025-10-05, commit: abc123)

# 3. Install using install script
curl -sSL https://raw.githubusercontent.com/user/my-context-copilot/main/scripts/install.sh | bash

# 4. Verify installation
my-context --version
which my-context
# Expected: ~/.local/bin/my-context

# 5. Verify existing data preserved (if upgrading from Sprint 1)
my-context list
# Expected: All Sprint 1 contexts visible
```

**Success Criteria**:
- Binary runs without "Go not found" error
- Installation completes without sudo
- Binary added to PATH automatically
- Existing contexts accessible after upgrade

---

## Scenario 2: Project-Based Workflow

**Goal**: Organize work by project using new project filtering

**Steps**:
```bash
# 1. Create contexts for ps-cli project
my-context start "Phase 1 - Foundation" --project ps-cli
my-context note "PLAN: Initial setup and requirements gathering"
my-context note "WORK: Created spec and plan documents"

# 2. Switch to Phase 2
my-context start "Phase 2 - Testing" --project ps-cli
my-context note "PLAN: Test strategy defined"

# 3. Work on different project
my-context start "Research" --project garden
my-context note "Looking into vegetable garden layouts"

# 4. Filter by project
my-context list --project ps-cli
# Expected: Shows only "ps-cli: Phase 1 - Foundation" and "ps-cli: Phase 2 - Testing"

my-context list --project garden
# Expected: Shows only "garden: Research"

# 5. List all contexts
my-context list --all
# Expected: Shows all 3 contexts
```

**Success Criteria**:
- Contexts created with "project: phase" naming
- Project filter shows only matching contexts (case-insensitive)
- Contexts without colons treated as standalone projects
- Multiple projects can coexist

---

## Scenario 3: Export and Share Context

**Goal**: Export context data to markdown for sharing with team

**Steps**:
```bash
# 1. Create and populate a context
my-context start "Sprint 1 Retrospective"
my-context note "What went well: TDD approach caught bugs early"
my-context note "What went wrong: WSL installation issues"
my-context file ~/projects/my-context-copilot/SPRINT-01-RETROSPECTIVE.md
my-context touch
my-context stop

# 2. Export to default location
my-context export "Sprint 1 Retrospective"
# Expected: Creates ./Sprint_1_Retrospective.md

# 3. Verify markdown content
cat Sprint_1_Retrospective.md
# Expected: Human-readable markdown with:
#   - Context name as header
#   - Start/end times and duration
#   - All notes with timestamps
#   - File associations
#   - Touch event count

# 4. Export to specific location
my-context export "Sprint 1 Retrospective" --to docs/retrospectives/sprint-1.md
# Expected: Creates docs/retrospectives/sprint-1.md (with parent dirs)

# 5. Export all contexts
mkdir -p exports
my-context export --all --to exports/
# Expected: Creates separate .md file for each context in exports/
```

**Success Criteria**:
- Markdown file generated with all context data
- File renders correctly in GitHub/VS Code
- Timestamps in local timezone
- Duration calculated correctly
- Parent directories created if needed

---

## Scenario 4: Context Lifecycle (Archive and Delete)

**Goal**: Manage completed and test contexts

**Steps**:
```bash
# 1. Create several contexts
my-context start "Test Context 1"
my-context note "Just testing"
my-context stop

my-context start "Test Context 2"
my-context note "Another test"
my-context stop

my-context start "ps-cli: Phase 1"
my-context note "Real work"
my-context stop

# 2. Archive completed work
my-context archive "ps-cli: Phase 1"
# Expected: "Archived context: ps-cli: Phase 1"

# 3. Verify archived hidden from default list
my-context list
# Expected: Shows Test Context 1 and 2, NOT ps-cli: Phase 1

# 4. View archived contexts
my-context list --archived
# Expected: Shows only ps-cli: Phase 1

# 5. Verify archived context still accessible
my-context show "ps-cli: Phase 1"
# Expected: Shows context details (archiving doesn't delete data)

# 6. Delete test contexts
my-context delete "Test Context 1"
# Expected: Confirmation prompt
# User enters: y
# Expected: "Deleted context: Test Context 1"

my-context delete "Test Context 2" --force
# Expected: No prompt, immediate deletion

# 7. Verify deletions
my-context list --all
# Expected: Shows only ps-cli: Phase 1 (archived)

# 8. Attempt to delete active context (should fail)
my-context start "Active Work"
my-context delete "Active Work"
# Expected: "Error: Cannot delete active context. Stop it first."
```

**Success Criteria**:
- Archive hides context from default list
- Archived contexts visible with --archived flag
- Archived data remains accessible (show, export work)
- Delete prompts for confirmation (unless --force)
- Delete removes entire context directory
- Cannot delete active context
- Transitions log preserved after deletion

---

## Scenario 5: List Enhancements (Large Dataset)

**Goal**: Test pagination and filtering with many contexts

**Steps**:
```bash
# 1. Create 25 contexts (simulating real usage)
for i in {1..25}; do
  my-context start "Context $i"
  my-context note "Note for context $i"
  my-context stop
  sleep 1  # Ensure unique timestamps
done

# 2. Default list (should show only 10)
my-context list
# Expected: Shows 10 most recent contexts
# Expected: "Showing 10 of 25 contexts. Use --all to see all."

# 3. Custom limit
my-context list --limit 5
# Expected: Shows 5 most recent contexts

# 4. Show all contexts
my-context list --all
# Expected: Shows all 25 contexts

# 5. Search by name
my-context list --search "Context 1"
# Expected: Shows Context 1, Context 10-19 (substring match)

# 6. Search with limit
my-context list --search "Context" --limit 3
# Expected: Shows 3 most recent contexts containing "Context"

# 7. Active-only filter
my-context start "Current Work"
my-context list --active-only
# Expected: Shows only "Current Work"
```

**Success Criteria**:
- Default limit is 10
- Truncation message displayed when more exist
- Custom limits work correctly
- --all shows everything
- Search is case-insensitive substring match
- --active-only shows single active context
- All filters work correctly together

---

## Scenario 6: Bug Fixes Validation

**Goal**: Verify Sprint 1 bugs are fixed

**Steps**:
```bash
# 1. Test $ character in notes (Bug #3)
my-context start "Budget Planning"
my-context note "Budget: $500-800 for equipment"
my-context note "Estimated cost: $1,200 total"
my-context show
# Expected: Both notes display with $ symbols intact
# Old behavior: Would show "Budget: 00-800" ($ stripped)

# 2. Test NULL display in history (Bug #5)
my-context start "First Context"
my-context stop
my-context history
# Expected: Shows "(none)" for empty previous/next context
# Old behavior: Would show "NULL"

# Example output:
# 2025-10-05 20:00 │ START     │ (none) → First Context
# 2025-10-05 20:05 │ STOP      │ First Context → (none)
```

**Success Criteria**:
- Special characters ($ , ! @ # etc.) preserved in notes
- History shows "(none)" instead of "NULL"
- Notes display correctly in show and export commands

---

## Scenario 7: Cross-Platform Installation (Windows)

**Goal**: Verify Windows installation works with cmd.exe and PowerShell

**Steps (cmd.exe)**:
```batch
REM 1. Download Windows binary
curl -O https://github.com/user/my-context-copilot/releases/download/v1.1.0/my-context-windows-amd64.exe

REM 2. Run installer
install.bat

REM 3. Verify installation
my-context --version
where my-context
REM Expected: C:\Users\Username\bin\my-context.exe

REM 4. Test basic functionality
my-context start "Windows Test"
my-context list
```

**Steps (PowerShell)**:
```powershell
# 1. Download binary
Invoke-WebRequest -Uri "https://..." -OutFile my-context.exe

# 2. Run installer
.\install.ps1

# 3. Verify installation
my-context --version
Get-Command my-context

# 4. Test basic functionality
my-context start "PowerShell Test"
my-context list
```

**Success Criteria**:
- Binary works on Windows without WSL
- Installation adds to user PATH (no admin required)
- Existing contexts preserved if upgrading
- Commands work identically to Linux/macOS

---

## Scenario 8: Backward Compatibility (Sprint 1 → Sprint 2)

**Goal**: Verify Sprint 1 data works with Sprint 2 binary

**Steps**:
```bash
# Prerequisites: Have Sprint 1 contexts in ~/.my-context/

# 1. Check Sprint 1 contexts
ls ~/.my-context/
# Expected: See existing context directories

# 2. Install Sprint 2 binary
./my-context-linux-amd64 --version
# Expected: Version 1.1.0

# 3. List Sprint 1 contexts
./my-context-linux-amd64 list
# Expected: All Sprint 1 contexts visible and functional

# 4. Show Sprint 1 context details
./my-context-linux-amd64 show "My_First_Context"
# Expected: Notes, files, timestamps all display correctly

# 5. Try new Sprint 2 features on Sprint 1 context
./my-context-linux-amd64 archive "My_First_Context"
# Expected: Archive succeeds (is_archived added to meta.json)

./my-context-linux-amd64 export "My_First_Context"
# Expected: Export succeeds (reads Sprint 1 data correctly)

# 6. Verify meta.json updated gracefully
cat ~/.my-context/My_First_Context/meta.json
# Expected: is_archived field present, other fields unchanged
```

**Success Criteria**:
- Sprint 1 contexts load without errors
- Notes, files, timestamps accessible
- New features work on old contexts
- meta.json updated non-destructively
- No data loss or corruption

---

## Scenario 9: JSON Output for Scripting

**Goal**: Verify all commands support --json for automation

**Steps**:
```bash
# 1. Start with JSON output
my-context start "API Test" --json | jq .
# Expected: Valid JSON with context_name, timestamp, etc.

# 2. List with JSON output
my-context list --project ps-cli --json | jq '.data.contexts | length'
# Expected: Number of ps-cli contexts

# 3. Export with JSON output
my-context export "API Test" --json | jq '.data.export_path'
# Expected: Path to exported markdown file

# 4. Archive with JSON output
my-context stop
my-context archive "API Test" --json | jq '.data.now_archived'
# Expected: true

# 5. Parse JSON for automation
CONTEXT_COUNT=$(my-context list --all --json | jq '.data.showing_count')
echo "Total contexts: $CONTEXT_COUNT"
```

**Success Criteria**:
- All commands support --json flag
- JSON output is valid (jq parsing succeeds)
- JSON structure documented in contracts
- Suitable for shell scripting and CI/CD

---

## Performance Benchmarks

**Goal**: Verify performance targets are met

```bash
# Benchmark 1: List command with 1000 contexts
time my-context list --all  # Target: < 1 second

# Benchmark 2: Export large context
# (Create context with 500 notes and 100 files first)
time my-context export "Large Context"  # Target: < 1 second

# Benchmark 3: Search across 1000 contexts
time my-context list --search "test"  # Target: < 1 second
```

**Success Criteria**:
- List with 1000 contexts: <1s on HDD
- Export with 500 notes: <1s
- Search across 1000 contexts: <1s

---

## Cleanup

```bash
# Remove test contexts
for i in {1..25}; do
  my-context delete "Context $i" --force
done

my-context delete "Test Context 1" --force 2>/dev/null
my-context delete "Test Context 2" --force 2>/dev/null
my-context delete "Sprint 1 Retrospective" --force 2>/dev/null
# etc.

# Or nuclear option:
# rm -rf ~/.my-context/*
# (preserves directory structure, removes all contexts)
```

---

**End of Quickstart**

These scenarios cover all Sprint 2 features and validate:
- ✅ Multi-platform installation
- ✅ Project filtering
- ✅ Export command
- ✅ List enhancements
- ✅ Archive/delete commands
- ✅ Bug fixes
- ✅ Cross-platform compatibility
- ✅ Backward compatibility
- ✅ JSON output
- ✅ Performance targets
