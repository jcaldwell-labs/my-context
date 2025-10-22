# Git Worktree Organization Complete

**Date**: 2025-10-22
**Purpose**: Document the new worktree-based repository structure

---

## New Structure

The my-context repository now uses a professional bare repository + worktree setup:

```
~/projects/
├── my-context.git/              # Bare repository (git database)
├── my-context-main/             # Main branch worktree
├── my-context-feature/          # Feature branch worktree (current work)
├── my-context-archive/          # Historical documentation (not git)
└── my-context/                  # OLD folder (can be removed)
```

### Bare Repository
**Location**: `~/projects/my-context.git`
- Central git database for all worktrees
- Contains all branches (9 total)
- Never work directly in this folder
- Commands: `git worktree add/list/remove`

### Main Worktree
**Location**: `~/projects/my-context-main`
**Branch**: `main`
- Primary development workspace
- Stable codebase
- Contains all committed documentation
- Use for: releases, stable features, documentation updates

### Feature Worktree
**Location**: `~/projects/my-context-feature`
**Branch**: `claude/review-context-protocol-011CULDrf4FxttSmyYupm5WL`
- Active feature development
- Working changes from previous workspace
- Contains uncommitted work and BUG-FIX-WATCH-EXEC.md

### Archive Folder
**Location**: `~/projects/my-context-archive`
- Historical documentation (15 files)
- Not part of git repository
- Reference only
- See `my-context-archive/README.md` for details

---

## Documentation Migrated

### From ~/projects/my-context-dev to ~/projects/my-context-main/docs/:
- ✅ WORKFLOW-METHODOLOGY.md
- ✅ BRANCH-STRATEGY-CLARIFICATION.md
- ✅ SDLC.md
- ✅ IMPLEMENTATION.md
- ✅ SETUP.md
- ✅ WINDOWS-BUILD-GUIDE.md
- ✅ TECH-DEBT-RESOLUTION-REPORT.md
- ✅ REMAINING-TECH-DEBT.md

### From ~/projects/my-context-dev to ~/projects/my-context-main/:
- ✅ CLAUDE.md (comprehensive developer guide)

### Archived to ~/projects/my-context-archive/:
- 13 historical documents (see archive README.md)

---

## Working with Worktrees

### List all worktrees
```bash
cd ~/projects/my-context.git
git worktree list
```

### Create new worktree for a feature
```bash
cd ~/projects/my-context.git
git worktree add ~/projects/my-context-<feature> <branch-name>
```

### Switch between worktrees
```bash
# Just cd to the worktree folder
cd ~/projects/my-context-main      # Work on main
cd ~/projects/my-context-feature   # Work on feature
```

### Remove a worktree
```bash
cd ~/projects/my-context.git
git worktree remove <path-to-worktree>
```

---

## Benefits

1. **Multiple branches simultaneously**: Work on main and feature without stashing
2. **Preserved IDE state**: Each worktree has independent .idea, .vscode, etc.
3. **Disk space savings**: Shared git database across all worktrees
4. **Clean separation**: No accidental commits to wrong branch
5. **Easy experimentation**: Spin up worktree for spike, remove when done

---

## Cleanup

The old folder can be safely removed:
```bash
# WARNING: Make sure you're not in this folder first!
# cd somewhere else, then:
rm -rf ~/projects/my-context
```

This folder was the original clone before worktree reorganization. All important:
- Code is in the new worktrees
- Documentation is committed to main worktree
- Historical docs are in archive folder
- All branches are in the bare repository

---

## .claude Configuration

Each worktree has independent `.claude/settings.local.json`:

**Main worktree**:
- Permissions for building, testing, running my-context commands
- Clean state for documentation work

**Feature worktree**:
- Same permissions as main
- Additional read access to main worktree and archive
- For cross-referencing during development

---

## Git Remote

All worktrees share the same remote:
```
origin: git@github.com:jcaldwell1066/my-context.git
```

Pushing/pulling works the same in any worktree, affecting only that worktree's branch.

---

## Next Steps

1. **Optional**: Remove old `~/projects/my-context` folder (original clone)
2. **Continue work** in `~/projects/my-context-feature`
3. **Stable updates** to `~/projects/my-context-main`
4. **Create new worktrees** for different features as needed
5. **Reference archive** at `~/projects/my-context-archive` when needed

---

**Setup completed**: 2025-10-22
**Documentation committed**: main branch (commit 857db93)
**Feature work committed**: claude/review-context-protocol-011CULDrf4FxttSmyYupm5WL (commit ad74b06)
