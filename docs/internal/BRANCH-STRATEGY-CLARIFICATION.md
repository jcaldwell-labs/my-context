# Branch Strategy Clarification: my-context Development

**Date**: 2025-10-11
**Purpose**: Prevent branch sprawl, clarify GitHub usage, document clean merge/release workflow

---

## Current Context Confusion

### What You're Seeing

**Command**: `my-context show`
**Result**: `sprint-006: sprint-006: uat-happy-days (status: stopped)`

**This is NOT our my-context development context!**

This is **deb-sanity team's UAT testing context**:
- **App**: deb-sanity project
- **Purpose**: UAT validation of my-context features
- **Owner**: deb-sanity team
- **Notes**: 16 (their integration test observations)
- **Status**: Stopped (they finished testing)

---

### Our Actual Communication Context

**Context**: `final-completion` (status: stopped)
**App**: my-context project
**Purpose**: Cross-team async communication (my-context ↔ deb-sanity)
**Notes**: 107 (implementation → validation → bugfix → UAT → vision)
**Files**: 14 (specs, analysis docs, responses)

**To resume OUR handoff**:
```bash
my-context resume final-completion
```

---

## Context Status Values (Official)

**From internal/models/context.go**:

1. **"active"** - Context currently running
   - Can add notes, files, touches
   - Shows in `my-context show` (no arguments)
   - Only ONE context can be active at a time

2. **"stopped"** - Context paused
   - Cannot add notes (must resume first)
   - Must use `my-context resume <name>` to continue
   - Can have MULTIPLE stopped contexts

3. **"archived"** - Boolean flag (is_archived: true/false)
   - Stopped contexts can be archived
   - Moves to archive storage (implementation detail)
   - Can still export, but won't show in default list

**State Transitions**:
```
[new] --start--> active --stop--> stopped --archive--> stopped (archived=true)
                    ↑                 |
                    └----resume-------┘
```

---

## Git Branch Strategy (Current State)

### Local Branches (6 total)

| Branch | Purpose | Commit | Status |
|--------|---------|--------|--------|
| **master** | Local development HEAD | 813e8b0 | Main dev branch |
| **main** | Public release pointer | bace3a8 | Points to v1.0.0 |
| **public-v1.0.0** | Public release (duplicate) | bace3a8 | **CAN DELETE** |
| **003-daily-summary-feature** | Paused feature | 540b6dc | Needs decision |
| **004-implement-5-lifecycle** ⭐ | Active iteration | f87f76d | **Current work** |
| **005-ux-polish-and** | Planning only | b66dc90 | Spec phase |

---

### GitHub Branches (4 total)

| Remote Branch | Maps From | Purpose | Status |
|---------------|-----------|---------|--------|
| **origin/main** | main | Public releases | PERMANENT |
| **origin/dev** | master | Development reference | PERMANENT |
| **origin/004-implement-5-lifecycle** | 004-* | Feature branch | TEMPORARY |
| **origin/005-ux-polish-and** | 005-* | Feature branch | TEMPORARY |

---

### Branch Relationships (Visual)

```
GitHub (Remote):
  origin/main (bace3a8) ←────┐
                             │ Public releases only
  origin/dev (813e8b0) ←─────┼─── master (local development)
                             │
  origin/004-* (f87f76d) ←───┼─── 004-implement-5-lifecycle (feature)
                             │
  origin/005-* (b66dc90) ←───┴─── 005-ux-polish-and (feature)

Local:
  master (813e8b0) ────────> origin/dev (development)
  main (bace3a8) ──────────> origin/main (public)
  public-v1.0.0 (bace3a8) ─> [CAN DELETE - duplicate of main]

  Feature branches:
  └─ 004-implement-5-lifecycle (f87f76d) ─> origin/004-*
  └─ 005-ux-polish-and (b66dc90) ────────> origin/005-*
  └─ 003-daily-summary (540b6dc) ────────> [NO REMOTE - local only]
```

---

## Iteration 004 Summary (Current Work)

**Branch**: `004-implement-5-lifecycle` (local + GitHub)
**App**: my-context project
**Target Release**: **v2.0.0** (after merge to master)

**What's in it** (12 commits ahead of master):
```
f87f76d - Strategic vision + terminology (Sprint→Iteration)
0d7f7f6 - UAT response (8 feature requests addressed)
29d01a2 - ISO8601 timestamps + supervisor monitoring
13bda2b - ISO8601 default change
101bfbf - Configurable timestamps
9956d3d - Signaling protocol requirements
bdac191 - Build version info
be60865 - --force bugfix
867acbe - deb-sanity notification
4c50251 - Implementation complete (121 tasks)
1c97e2c - .cursor enhancements
34c0fca - Spec/plan/tasks complete
813e8b0 - POC scripts (← master is here)
```

**Features Delivered**:
1. Smart resume (prevents fragmentation)
2. Note warnings (50, 100 note thresholds)
3. Resume command ('r' alias)
4. Bulk archive (pattern-based)
5. Lifecycle advisor (post-stop guidance)
6. Granular timestamps (ISO8601 configurable) ⭐ BONUS

**Status**: ✅ Implementation complete, UAT validated (10/10 tests passed)

---

## Branch Sprawl Prevention Strategy

### Cleanup Plan (Prevent Mess)

**Immediate** (When Iteration 004 Validated):
```bash
# 1. Merge to master
git checkout master
git merge 004-implement-5-lifecycle --no-ff  # Keep history
git push origin HEAD:dev  # Update GitHub dev branch

# 2. Tag release
git tag -a v2.0.0 -m "Release v2.0.0: Lifecycle improvements + timestamps"
git push origin v2.0.0

# 3. Delete feature branches
git branch -d 004-implement-5-lifecycle  # Local
git push origin --delete 004-implement-5-lifecycle  # GitHub

# 4. Delete redundant branch
git branch -d public-v1.0.0  # Points to same commit as main
```

**After cleanup**:
- Local: master, main, 003-*, 005-* (4 branches)
- GitHub: origin/main, origin/dev, origin/005-* (3 branches if 005 active)

---

**For 003-daily-summary**:
```bash
# Option A: Close (not ready)
git branch -D 003-daily-summary-feature

# Option B: Keep (will revisit)
# Leave it, address later when ready to implement
```

**For 005-ux-polish**:
```bash
# Keep until implemented, then merge and delete
# Or: Delete if absorbing into Iteration 004 (already did timestamps)
```

---

### GitHub Branch Policy (Going Forward)

**PERMANENT Branches** (2 only):
- **main**: Public releases (v1.0.0, v2.0.0, etc.)
- **dev**: Development branch (shows full process)

**TEMPORARY Feature Branches**:
- **XXX-feature-name**: Create when starting implementation
- **Merge**: When feature complete and validated
- **DELETE**: Immediately after merge (keep GitHub clean)

**Workflow**:
```
1. Create feature branch: git checkout -b 006-signaling master
2. Push to GitHub: git push -u origin 006-signaling
3. Develop, commit, push
4. When complete: Merge to master
5. Push to dev: git push origin master:dev
6. DELETE feature branch: git push origin --delete 006-signaling
7. Tag if release: git tag v2.2.0 && git push origin v2.2.0
```

**Max branches on GitHub**: 2 permanent + 2-3 active features = 4-5 total

---

## Worktree Arrangement (Current)

**Single worktree**: `/home/be-dev-agent/projects/my-context-dev`
**No additional worktrees**

**This is fine!** Worktrees not needed for:
- Feature development (branch switching works)
- Single developer workflow
- CI/CD (GitHub Actions uses checkout)

**Could use worktrees for**:
- Testing features without switching (minor benefit)
- Running multiple builds simultaneously (rare)

**Recommendation**: **No worktrees needed** - Branch switching sufficient

---

## Clear Path Forward

### Iteration 004 (Current - Ready to Close)

**Status**: ✅ Complete, validated, ready to merge
**Branch**: 004-implement-5-lifecycle (f87f76d)
**Action Plan**:
1. Wait for final deb-sanity confirmation (or proceed now if satisfied)
2. Merge to master
3. Push to origin/dev
4. Tag v2.0.0
5. Delete 004 branches (local + GitHub)
6. Update binary version to v2.0.0 (remove -dev)

---

### Iteration 005 (Future - Absorb or Abandon?)

**Status**: Spec only (b66dc90)
**Branch**: 005-ux-polish-and
**Contains**: Roadmap, multi-feature docs, Sprint 005 spec

**Decision Needed**:
- **Option A**: Absorb into Iteration 004 (already did timestamps, other items minor)
- **Option B**: Keep separate, implement after 004 merges
- **Option C**: Delete (features already addressed elsewhere)

**Recommendation**: **Option A** - Delete 005 branch, timestamps already in 004

---

### GitHub Branch Health (Target State)

**After Iteration 004 Cleanup**:
```
GitHub Branches (2 permanent):
├── main (v2.0.0 tag) - Public release
└── dev (master) - Development

Local Branches (2-3):
├── master - Active development
├── main - Public release pointer
└── 003-daily-summary (optional - if keeping)
```

**Clean, minimal, clear purpose for each branch** ✅

---

## Recommendations

1. **Resume final-completion** for handoff (not sprint-006 UAT context)
2. **Merge Iteration 004** to master (ready now or after final confirmation)
3. **Delete redundant branches** (public-v1.0.0, possibly 005)
4. **Adopt branch policy**: Feature branches are temporary, delete after merge
5. **Keep GitHub clean**: 2 permanent + active features only

**Create detailed merge/cleanup checklist?**