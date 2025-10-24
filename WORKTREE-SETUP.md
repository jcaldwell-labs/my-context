# Git Worktree Organization Complete

**Date**: 2025-10-24 (Updated)
**Purpose**: Document the worktree-based repository structure for parallel version development

---

## Current Structure (Updated 2025-10-24)

The my-context repository uses git worktrees for parallel version development:

```
~/projects/
├── my-context.git/              # Bare repository (git database)
├── my-context/                  # Main worktree (stable, main branch)
├── my-context-main/             # Legacy main worktree
├── my-context-feature/          # Legacy feature worktree
├── my-context-v2.5.1/           # Patch release worktree (bug fixes)
├── my-context-v2.6.0/           # Minor release worktree (new features)
├── my-context-v3.0.0/           # Major release worktree (breaking changes)
├── my-context-skunkworks/       # Experimental worktree (research)
└── my-context-archive/          # Historical documentation (not git)
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

### Version Worktrees (New - 2025-10-24)

**v2.5.1 Patch Release**
**Location**: `~/projects/my-context-v2.5.1`
**Branch**: `release/v2.5.1`
- Bug fixes and completeness work
- Signal command implementation
- Cross-platform testing
- JSON output validation
- **Timeline**: 1-2 weeks
- **Breaking Changes**: None

**v2.6.0 Minor Release**
**Location**: `~/projects/my-context-v2.6.0`
**Branch**: `release/v2.6.0`
- Performance optimizations
- Enhanced filtering/search
- Quality of life improvements
- Developer experience features
- **Timeline**: 2-4 weeks
- **Breaking Changes**: None

**v3.0.0 Major Release**
**Location**: `~/projects/my-context-v3.0.0`
**Branch**: `release/v3.0.0`
- Multi-channel architecture (foreground/background/experiment/shared)
- Detached contexts
- Branch/merge/discard pattern
- **Timeline**: 6-10 weeks
- **Breaking Changes**: YES - state.json → channels.json migration

**Skunkworks (Experimental)**
**Location**: `~/projects/my-context-skunkworks`
**Branch**: `experimental/distributed-journal`
- Distributed ecosystem journal research
- Event sourcing architecture
- Team collaboration primitives
- **Timeline**: Continuous research
- **Breaking Changes**: EVERYTHING - may never merge

### Legacy Worktrees

**Feature Worktree**
**Location**: `~/projects/my-context-feature`
**Branch**: `claude/review-context-protocol-011CULDrf4FxttSmyYupm5WL`
- Legacy active feature development
- Can be removed once work merged

**Main Worktree (Legacy)**
**Location**: `~/projects/my-context-main`
**Branch**: `main`
- Legacy main branch worktree
- Can be removed (use ~/projects/my-context instead)

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
cd ~/projects/my-context  # Or any worktree
git worktree list
```

Output:
```
~/projects/my-context             ee1b7f2 [main]
~/projects/my-context-v2.5.1      ee1b7f2 [release/v2.5.1]
~/projects/my-context-v2.6.0      ee1b7f2 [release/v2.6.0]
~/projects/my-context-v3.0.0      ee1b7f2 [release/v3.0.0]
~/projects/my-context-skunkworks  ee1b7f2 [experimental/distributed-journal]
```

### Switch between worktrees
```bash
# Just cd to the worktree folder
cd ~/projects/my-context              # Stable main branch
cd ~/projects/my-context-v2.5.1       # Work on patches
cd ~/projects/my-context-v2.6.0       # Work on minor features
cd ~/projects/my-context-v3.0.0       # Work on major changes
cd ~/projects/my-context-skunkworks   # Experiment freely
```

### Create new worktree (if needed)
```bash
cd ~/projects/my-context  # Or any worktree
git worktree add ../my-context-<feature> <branch-name>
```

### Remove a worktree
```bash
cd ~/projects/my-context
git worktree remove ../my-context-<feature>
```

### Version-Specific Workflows

**Working on v2.5.1 (Patches)**
```bash
cd ~/projects/my-context-v2.5.1
# Fix signal command implementation
# Add tests
# Cross-platform validation
git commit -m "fix: complete signal command"
git push origin release/v2.5.1
```

**Working on v2.6.0 (Minor Features)**
```bash
cd ~/projects/my-context-v2.6.0
# Add performance optimizations
# Implement enhanced filtering
git commit -m "feat: add regex search"
git push origin release/v2.6.0
```

**Working on v3.0.0 (Major Changes)**
```bash
cd ~/projects/my-context-v3.0.0
# Implement multi-channel architecture
# Breaking changes allowed
git commit -m "feat!: add multi-channel support"
git push origin release/v3.0.0
```

**Experimenting in Skunkworks**
```bash
cd ~/projects/my-context-skunkworks
# Try crazy ideas
# No merge pressure
# Cherry-pick successful experiments to version branches
git commit -m "experiment: distributed journal prototype"
# May never push or merge - that's okay!
```

---

## Benefits

1. **Parallel Version Development**: Work on v2.5.1 patches, v2.6.0 features, and v3.0.0 architecture simultaneously
2. **No Context Switching**: Each version has isolated workspace - no git stash/checkout needed
3. **Preserved IDE State**: Each worktree has independent .idea, .vscode, build artifacts, etc.
4. **Disk Space Savings**: Shared git database across all worktrees (4 worktrees ≈ 1.2x space, not 4x)
5. **Clean Separation**: Impossible to accidentally commit v3.0.0 breaking changes to v2.5.1
6. **Safe Experimentation**: Skunkworks branch for research - no merge pressure
7. **Independent Testing**: Test each version's binaries without rebuilding
8. **Clear Release Pipeline**: v2.5.1 → v2.6.0 → v3.0.0 progression visible in filesystem

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

### Immediate (v2.5.1 - Week 1-2)
1. **cd ~/projects/my-context-v2.5.1**
2. **Complete signal command** implementation (12 TODO tests)
3. **Cross-platform testing** (Windows, macOS, Linux)
4. **JSON output validation** for all commands
5. **Release v2.5.1** when all tests pass

### Short-term (v2.6.0 - Week 2-6)
1. **cd ~/projects/my-context-v2.6.0**
2. **Performance benchmarks** (test with 1000+ contexts)
3. **Enhanced filtering** (regex search, date ranges)
4. **Quality of life** features (templates, bulk operations)
5. **Developer experience** (shell completions, config files)
6. **Release v2.6.0** when features complete

### Long-term (v3.0.0 - Week 4-14)
1. **cd ~/projects/my-context-v3.0.0**
2. **Multi-channel architecture** design and implementation
3. **Detached contexts** support
4. **Migration tooling** (v2.x → v3.x)
5. **Breaking changes** documentation
6. **Release v3.0.0** when stable and tested

### Continuous (Skunkworks)
1. **cd ~/projects/my-context-skunkworks**
2. **Experiment** with distributed journal ideas
3. **Prototype** event sourcing architecture
4. **Research** team collaboration patterns
5. **Cherry-pick** successful experiments to version branches

### Cleanup (Optional)
1. **Remove legacy worktrees**: `~/projects/my-context-main`, `~/projects/my-context-feature`
2. **Keep archive**: `~/projects/my-context-archive` (historical reference)

---

**Initial Setup**: 2025-10-22
**Version Worktrees Added**: 2025-10-24
**Documentation Updated**: 2025-10-24
**Current Stable Version**: v2.5.0 (main branch)
