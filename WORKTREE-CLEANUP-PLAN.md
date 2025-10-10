# Worktree Cleanup Plan

**Date**: October 9, 2025  
**Status**: All branches merged to master  

---

## Current Worktrees

| Worktree Path | Branch | Latest Commit | Status |
|--------------|--------|---------------|--------|
| `.git` (bare) | 002-installation-improvements-and | b3c38ea | ✅ Merged |
| `my-context-copilot-001-cli-context-management` | 001-cli-context-management | 14beb4b | ✅ Merged (Sprint 1) |
| `my-context-copilot-002-gap-analysis-and` | 002-gap-analysis-and | 2c32f9b | ✅ Merged (Abandoned?) |
| `my-context-copilot-002-installation-improvements-and` | 002-installation-improvements-and | b3c38ea | ✅ Merged (Sprint 2) |
| `my-context-copilot-003-daily-summary-feature` | 003-daily-summary-feature | 540b6dc | ✅ In master (Future work) |
| `my-context-copilot-master` | master | c6e8bcc | ⭐ Current |

---

## Verification

All branches show as merged to master:
```bash
$ git branch --merged master
  001-cli-context-management
  002-gap-analysis-and
  002-installation-improvements-and
  003-daily-summary-feature
  master
```

All commits found in master's history:
```bash
$ git log --oneline master | grep -E "14beb4b|2c32f9b|540b6dc"
540b6dc intermediate
2c32f9b manual
14beb4b feat: implement CLI context management system (v1.0.0)
```

---

## Cleanup Recommendations

### ✅ SAFE TO DELETE (Merged & Complete)

#### 1. `001-cli-context-management` (Sprint 1)
- **Commit**: 14beb4b "feat: implement CLI context management system (v1.0.0)"
- **Status**: Sprint 1 complete, fully merged
- **Recommendation**: **DELETE** worktree and branch
- **Reason**: Historical, no active work

#### 2. `002-installation-improvements-and` (Sprint 2)
- **Commit**: b3c38ea "docs: revise integration analysis based on peer review"
- **Status**: Sprint 2 complete, merged, constitution ratified
- **Recommendation**: **DELETE** worktree and branch
- **Reason**: Sprint 2 is production-ready (v1.0.0), no further work needed

#### 3. `002-gap-analysis-and` (Gap Analysis - Abandoned?)
- **Commit**: 2c32f9b "manual"
- **Status**: Unclear purpose, minimal commit message
- **Recommendation**: **DELETE** worktree and branch
- **Reason**: Appears to be exploratory/abandoned work, already merged

### ⚠️ REVIEW BEFORE DELETING

#### 4. `003-daily-summary-feature` (Sprint 3 - Future Work)
- **Commit**: 540b6dc "intermediate"
- **Status**: Sprint 3 planning/specification phase
- **Files**: `specs/003-daily-summary-feature/spec.md`, `specs/003-daily-summary-feature/plan.md`
- **Recommendation**: **KEEP** if active planning, **DELETE** if abandoned
- **Question**: Is Sprint 3 still planned?

---

## Cleanup Commands

### Step 1: Verify No Uncommitted Work

```bash
# Check each worktree for uncommitted changes
for dir in /home/be-dev-agent/projects/my-context-copilot-*; do
  echo "=== $(basename $dir) ==="
  cd "$dir" && git status --short
done
```

### Step 2: Remove Completed Sprint Worktrees

```bash
cd /home/be-dev-agent/projects/my-context-copilot-master

# Remove Sprint 1 (001)
git worktree remove /home/be-dev-agent/projects/my-context-copilot-001-cli-context-management
git branch -d 001-cli-context-management

# Remove Sprint 2 (002-installation-improvements-and)
git worktree remove /home/be-dev-agent/projects/my-context-copilot-002-installation-improvements-and
git branch -d 002-installation-improvements-and

# Remove gap analysis (002-gap-analysis-and)
git worktree remove /home/be-dev-agent/projects/my-context-copilot-002-gap-analysis-and
git branch -d 002-gap-analysis-and
```

### Step 3: Handle Sprint 3 (Decision Required)

**Option A: Keep Sprint 3** (if still planned)
```bash
# Leave worktree as-is for future work
# No action needed
```

**Option B: Delete Sprint 3** (if abandoned)
```bash
# Archive spec files first
cp -r /home/be-dev-agent/projects/my-context-copilot-003-daily-summary-feature/specs/003-daily-summary-feature \
     /home/be-dev-agent/projects/my-context-copilot-master/specs/

# Remove worktree
git worktree remove /home/be-dev-agent/projects/my-context-copilot-003-daily-summary-feature
git branch -d 003-daily-summary-feature
```

### Step 4: Clean Up Bare Repo Worktree

The `.git` bare repo at `/home/be-dev-agent/projects/my-context-copilot.git` shows branch `002-installation-improvements-and`. This is unusual - bare repos typically shouldn't have a checked-out branch.

```bash
# Check if it's truly a bare repo
cd /home/be-dev-agent/projects/my-context-copilot.git
git config --get core.bare

# If bare=true, this might just be a display artifact from git worktree list
# No action needed unless there are actual files checked out
```

---

## Post-Cleanup State

After cleanup, you should have:

```
/home/be-dev-agent/projects/
├── my-context-copilot.git/                   # Bare repo
├── my-context-copilot-master/                 # Master worktree (KEEP)
└── my-context-copilot-003-daily-summary-feature/  # Sprint 3 (DECISION REQUIRED)
```

---

## Disk Space Recovery

Estimated space to be freed:
- `001-cli-context-management`: ~50MB (source code + git metadata)
- `002-gap-analysis-and`: ~50MB
- `002-installation-improvements-and`: ~50MB
- **Total**: ~150MB (if all 3 deleted)

Additional cleanup after worktree removal:
```bash
cd /home/be-dev-agent/projects/my-context-copilot-master
git gc --aggressive --prune=now
```

---

## Recommendations Summary

| Worktree | Action | Reason |
|----------|--------|--------|
| 001-cli-context-management | ✅ DELETE | Sprint 1 complete |
| 002-gap-analysis-and | ✅ DELETE | Abandoned/merged |
| 002-installation-improvements-and | ✅ DELETE | Sprint 2 complete (v1.0.0 ready) |
| 003-daily-summary-feature | ⚠️ REVIEW | Future work - keep if planned, delete if abandoned |
| master | ✅ KEEP | Primary worktree |
| .git (bare) | ✅ KEEP | Repository base |

---

## Decision Required

**Question for Human**: Is Sprint 3 (003-daily-summary-feature) still planned?

- **If YES**: Keep the worktree for future planning/implementation
- **If NO**: Archive the specs to master and delete the worktree

Once decided, I can execute the cleanup commands.

