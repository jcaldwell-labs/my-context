# Sprint 2 Demo - Quick Reference Card

**Version**: 2.0.0-sprint2-uat
**Date**: 2025-10-06
**Time Required**: 30-45 minutes

---

## Pre-Flight Check

```bash
my-context --version  # Should show: 2.0.0-sprint2-uat
my-context list --all  # See current state
```

---

## Critical Path Demo (15 min - MUST PASS)

### 1. Project Filtering (P1)

```bash
# Create with project prefix
my-context start "Demo Session" --project sprint-2
# ✅ Expected: "Started context: sprint-2: Demo Session"

my-context note "Testing project feature"
my-context start "Bug Fix" --project sprint-2
my-context start "Home Task" --project personal

# Filter by project
my-context list --project sprint-2
# ✅ Expected: Shows only 2 sprint-2 contexts

my-context list --project SPRINT-2  # Case insensitive
# ✅ Expected: Same as above
```

**Result**: ✅ PASS / ❌ FAIL

---

### 2. Export Command (P1)

```bash
my-context start "Export Demo" --project sprint-2
my-context note "First note"
my-context note "Budget: $500-$800"
my-context file "CLAUDE.md"

# Basic export
my-context export "sprint-2: Export Demo"
# ✅ Expected: "Exported context ... to sprint-2__Export_Demo.md"

cat sprint-2__Export_Demo.md
# ✅ Expected: Markdown with notes, files, $ preserved

# Custom path
my-context export "sprint-2: Export Demo" --to test-export.md
# ✅ Expected: File created at test-export.md
```

**Result**: ✅ PASS / ❌ FAIL

---

### 3. Bug Fixes (P1)

```bash
# Dollar sign preserved
my-context note "Cost: $1000"
my-context show
# ✅ Expected: Note shows "Cost: $1000" (not escaped)

# History shows "(none)" not "NULL"
my-context history
# ✅ Expected: No "NULL" anywhere, shows "(none)"
```

**Result**: ✅ PASS / ❌ FAIL

---

## Safety-Critical Tests (10 min - BLOCKERS)

### 4. Delete Command (HIGH RISK)

```bash
my-context stop
my-context start "Delete Test" --project temp

# Try delete active (should fail)
my-context delete "temp: Delete Test" --force
# ✅ Expected: Error - cannot delete active context

my-context stop

# Delete with confirmation
my-context delete "temp: Delete Test"
# ✅ Expected: Prompts for confirmation

# Delete with force
my-context delete "temp: Delete Test" --force
# ✅ Expected: "Deleted context: temp: Delete Test"

my-context list --all | grep "Delete Test"
# ✅ Expected: No results (context gone)
```

**Result**: ✅ PASS / ❌ FAIL / 🔲 SKIPPED

---

### 5. Archive Command

```bash
my-context start "Archive Test" --project temp
my-context stop

my-context archive "temp: Archive Test"
# ✅ Expected: "Archived context: temp: Archive Test"

my-context list --project temp
# ✅ Expected: Context NOT shown (hidden)

my-context list --archived
# ✅ Expected: Shows "Archive Test" with archived flag
```

**Result**: ✅ PASS / ❌ FAIL / 🔲 SKIPPED

---

## Nice-to-Have Tests (10 min - OPTIONAL)

### 6. List Enhancements

```bash
my-context list  # Default limit 10
my-context list --limit 3  # Show only 3
my-context list --search demo  # Search filter
my-context list --project sprint-2 --limit 5  # Combined
```

**Result**: ✅ PASS / ❌ FAIL / 🔲 SKIPPED

---

### 7. Backward Compatibility

```bash
# Show old contexts (created before Sprint 2)
my-context list --all
# ✅ Expected: All contexts visible, no errors

# Operate on Sprint 1 context
my-context export "<some-sprint-1-context>"
# ✅ Expected: Works without errors
```

**Result**: ✅ PASS / ❌ FAIL / 🔲 SKIPPED

---

## Score Summary

**Critical Path** (Must Pass):
- [ ] Project filtering (start + list)
- [ ] Export command
- [ ] Bug fixes

**Safety Critical** (Blockers if fail):
- [ ] Delete command
- [ ] Archive command

**Nice-to-Have**:
- [ ] List enhancements
- [ ] Backward compatibility

---

## Decision Matrix

| Scenario | Action |
|----------|--------|
| All Critical Pass + Safety Pass | ✅ **SIGN OFF SPRINT 2** |
| All Critical Pass + Safety Fail | ⚠️ **FIX SAFETY → RETEST** |
| Any Critical Fail | ❌ **FIX CRITICAL → FULL RETEST** |
| Critical Pass + Safety Skipped | ⚠️ **DOCUMENT RISK → DECIDE** |

---

## Quick Commands Reference

```bash
# Setup
my-context --version
my-context list --all

# Project workflow
my-context start "<name>" --project <project>
my-context list --project <project>

# Export
my-context export "<context-name>"
my-context export "<context-name>" --to <path>

# Archive/Delete
my-context archive "<context-name>"
my-context delete "<context-name>"
my-context delete "<context-name>" --force

# List options
my-context list --limit <n>
my-context list --search <term>
my-context list --archived
my-context list --all

# JSON output
my-context <command> --json
```

---

## Time Estimates

- **Minimum Demo** (Critical only): 15 min
- **Safety Testing** (Delete/Archive): +10 min
- **Full Demo** (All features): 30-45 min

---

## Post-Demo Actions

1. **Calculate pass rate**: ___ / 7 tests passed
2. **Identify blockers**: List any FAIL results
3. **Make decision**: Sign off / Fix / Defer
4. **Create tech debt issues** for:
   - Installation scripts (install.bat, install.ps1, curl-install.sh)
   - GitHub Actions workflow
   - Any untested P2 features

---

**Ready to demo? Start with Pre-Flight Check above! ☝️**
