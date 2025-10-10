# My-Context: Current Repository State

**Last Updated**: 2025-10-10 16:45
**Session**: setup github repo (extended)
**Status**: ✅ Stable - Ready for development

---

## Repository Overview

This is the **my-context** project - a cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps.

**GitHub Repository**: https://github.com/jcaldwell1066/my-context

---

## Current Working Environment

### Working Directory
```
~/projects/my-context-dev/
```

This is the main development worktree for the my-context project.

### Current Branch
```
master (local) → origin/dev (GitHub)
```

- **Local branch name**: `master`
- **Tracks**: `origin/dev` on GitHub
- **Latest commit**: `3f6a516` "docs: add deb-sanity workflow analysis..."
- **Status**: Clean, up to date with remote

### How to Push Changes
```bash
# Since master tracks origin/dev (different names), use explicit push:
git push origin HEAD:dev

# Or set up simplified push permanently:
git config push.default upstream
git push  # Now just works
```

---

## Branch Relationships

### Development Branch (This Environment)
```
master (local) → origin/dev (GitHub)
```

**Purpose**: Full internal development environment
**Contains**:
- Complete source code
- Internal tooling (`.claude/`, `.cursor/`, `.specify/`)
- Feature specifications (`specs/`)
- SDLC documentation (IMPLEMENTATION.md, SDLC.md, etc.)
- Development automation scripts
- Sprint retrospectives and process docs

**URL**: https://github.com/jcaldwell1066/my-context/tree/dev

**Use this for**:
- Daily development work
- Feature implementation
- Internal documentation
- Process improvements

---

### Public Release Branch
```
public-v1.0.0 (local) → origin/main (GitHub)
main (local) - also points to clean commit
```

**Purpose**: Clean public-facing release
**Contains**:
- Production source code
- Public documentation (README, LICENSE, CONTRIBUTING)
- GitHub Actions workflows
- **Excludes**: 108 internal files (IDE configs, specs, internal docs)

**URL**: https://github.com/jcaldwell1066/my-context (default branch)

**Use this for**:
- Public releases
- Version tagging
- External documentation updates

**How to update**:
```bash
# Switch to public release worktree (if it exists)
cd ~/projects/my-context-public-v1.0.0

# Or create fresh from current state
git worktree add ~/projects/my-context-public-v1.0.0 public-v1.0.0
./scripts/prepare-public-release.sh
git push origin public-v1.0.0:main
```

---

### Feature Branches
```
003-daily-summary-feature (local)
```

**Purpose**: In-progress feature development
**Status**: Intermediate work (not complete)

**How to work on features**:
```bash
# Create new feature worktree
git worktree add ~/projects/my-context-feature-xyz feature/xyz

# Or use deb-sanity (if available)
deb-sanity --worktree-create feature/xyz
```

---

## Remote Configuration

### GitHub Remote
```
origin: git@github.com:jcaldwell1066/my-context.git
```

**Authentication**: SSH keys
**Location**: `~/.ssh/id_rsa`, `~/.ssh/id_rsa.pub`
**Permissions**: 600 (private key), 644 (public key)

**Test authentication**:
```bash
ssh -T git@github.com
# Should show: Hi jcaldwell1066! You've successfully authenticated...
```

### Branch Mappings
| Local Branch | → | Remote Branch | Purpose |
|--------------|---|---------------|---------|
| `master` | → | `origin/dev` | Full development environment |
| `public-v1.0.0` | → | `origin/main` | Clean public release |
| Feature branches | → | `origin/feature/*` | Feature development |

---

## Recent Activity

### Session: "setup github repo" (2h 45m)
**Date**: 2025-10-10 14:16 - 16:45

**Accomplished**:
1. ✅ Created dual-branch GitHub publishing strategy
2. ✅ Published v1.0.0 to GitHub (public + dev branches)
3. ✅ Created automation scripts (github-preflight.sh, prepare-public-release.sh)
4. ✅ Recovered from shell crash (worktree cleanup incident)
5. ✅ Created comprehensive deb-sanity workflow analysis (950 lines)

**Context tracked**: 30 timestamped notes in my-context

**Key documents created**:
- `docs/RETROSPECTIVE-GITHUB-SETUP.md` - Session retrospective
- `docs/DEB-SANITY-WORKFLOW-ANALYSIS.md` - Workflow analysis for deb-sanity team
- `scripts/github-preflight.sh` - Pre-flight checks
- `scripts/prepare-public-release.sh` - Automated public release prep

---

## Next Steps

### Immediate (This Sprint)
1. **Review deb-sanity analysis**
   - Document: `docs/DEB-SANITY-WORKFLOW-ANALYSIS.md`
   - Share with deb-sanity project team
   - Schedule review session to validate findings

2. **Continue feature development**
   - Feature branch: `003-daily-summary-feature`
   - Status: Intermediate work in progress

3. **Maintain GitHub branches**
   - Keep `master → origin/dev` up to date (ongoing development)
   - Update `public-v1.0.0 → origin/main` only for releases

### Short-Term (Next Sprint)
1. **Implement deb-sanity enhancements**
   - Priority 1: GitHub publishing workflows
   - Priority 2: File preparation system
   - Based on analysis recommendations

2. **Extract reusable patterns**
   - Dual-branch publishing workflow
   - Automation scripts for other projects
   - Workflow templates

3. **Process improvements**
   - Document worktree workflows
   - Create runbooks for common operations
   - Automate repetitive tasks

### Long-Term
1. **deb-sanity evolution**
   - Evolve from worktree tool → repository publishing platform
   - Build workflow template library
   - Capture more organic workflows

2. **Project maturity**
   - Continuous integration setup
   - Automated release process
   - Community contribution guidelines

---

## Common Operations

### Daily Development Workflow
```bash
# 1. Work in main dev environment
cd ~/projects/my-context-dev

# 2. Create feature branch
git checkout -b feature/new-feature

# 3. Make changes, test, commit
git add .
git commit -m "feat: implement new feature"

# 4. Push to GitHub dev branch
git push origin HEAD:dev

# 5. Track progress
my-context note "Completed feature X implementation"
```

### Creating a Public Release
```bash
# 1. Ensure working tree is clean
git status

# 2. Create/update public release branch
./scripts/prepare-public-release.sh

# 3. Run pre-flight checks
./scripts/github-preflight.sh

# 4. Push to public main branch
git push origin public-v1.0.0:main

# 5. Tag release
git tag -a v1.x.x -m "Release v1.x.x"
git push origin v1.x.x
```

### Working with Multiple Features
```bash
# Use git worktrees for parallel feature development
git worktree add ~/projects/my-context-feature-a feature/a
git worktree add ~/projects/my-context-feature-b feature/b

# Work in isolation
cd ~/projects/my-context-feature-a
# ... develop feature A ...

cd ~/projects/my-context-feature-b
# ... develop feature B ...

# Clean up when done
git worktree remove ~/projects/my-context-feature-a
git worktree remove ~/projects/my-context-feature-b
```

### Checking Repository Health
```bash
# Git status
git status
git branch -vv
git remote -v

# Worktree health (if deb-sanity available)
deb-sanity --health ~/projects/my-context-dev

# Context status
my-context show
```

---

## Troubleshooting

### Issue: Can't push to GitHub
**Symptoms**: Authentication fails or permission denied

**Solution**:
```bash
# 1. Check SSH keys
ls -la ~/.ssh/
# Should see: id_rsa (600), id_rsa.pub (644)

# 2. Test GitHub authentication
ssh -T git@github.com

# 3. If fails, copy keys from Windows (if applicable)
cp /mnt/c/Users/$USER/.ssh/id_rsa ~/.ssh/
chmod 600 ~/.ssh/id_rsa

# 4. Run pre-flight checks
./scripts/github-preflight.sh
```

### Issue: "upstream branch does not match" error
**Symptoms**: `git push` fails with branch name mismatch

**Solution**:
```bash
# Use explicit push
git push origin HEAD:dev

# Or configure for simplified push
git config push.default upstream
```

### Issue: Worktree directory deleted accidentally
**Symptoms**: Shell can't execute commands, git confused

**Solution**:
```bash
# 1. Navigate to home directory
cd ~

# 2. Check worktree status
git worktree list

# 3. Remove orphaned worktree reference
git worktree prune

# 4. Recreate if needed
git worktree add ~/projects/my-context-dev master

# OR use deb-sanity for safer management
deb-sanity --worktree-create master
```

---

## Key Files and Locations

### Important Documents
- `README.md` - Project overview (public)
- `IMPLEMENTATION.md` - Implementation guide (internal)
- `SDLC.md` - Software development lifecycle (internal)
- `docs/RETROSPECTIVE-GITHUB-SETUP.md` - GitHub setup session retrospective
- `docs/DEB-SANITY-WORKFLOW-ANALYSIS.md` - deb-sanity workflow analysis
- `CURRENT-STATE.md` - This document

### Automation Scripts
- `scripts/github-preflight.sh` - Pre-flight checks before GitHub operations
- `scripts/prepare-public-release.sh` - Automated public release preparation
- `scripts/build.sh` - Build my-context binary
- `scripts/install.sh` - Install my-context locally

### Configuration
- `.gitignore` - Git ignore patterns
- `.github/workflows/release.yml` - GitHub Actions release workflow
- `.claude/` - Claude Code commands (internal)
- `.cursor/` - Cursor IDE commands (internal)

### Source Code
- `cmd/my-context/main.go` - CLI entry point
- `internal/` - Core implementation
  - `commands/` - CLI commands
  - `core/` - Business logic
  - `models/` - Data models
  - `output/` - Output formatting

---

## Project Context

### What is my-context?
A CLI tool for managing developer work contexts with:
- Start/stop/switch between contexts
- Timestamped notes and file associations
- Context transition history
- Export to markdown
- Archive and delete completed contexts

### Why dual-branch strategy?
**Public branch** (`origin/main`):
- Clean codebase for external users
- Professional appearance
- No internal tooling or process docs
- Focus on end-user value

**Dev branch** (`origin/dev`):
- Transparency in development process
- Educational value for contributors
- Shows spec-driven development approach
- Documents decision-making and learnings

### Related Projects
- **deb-sanity**: Environment management tool being developed in parallel
  - Location: `~/projects/deb-sanity`
  - Purpose: Bootstrap CLI for environment awareness and context switching
  - Integration: my-context is used for session tracking in deb-sanity development

---

## Contact and Resources

### GitHub Repository
https://github.com/jcaldwell1066/my-context

### Branches
- Main (public): https://github.com/jcaldwell1066/my-context
- Dev (internal): https://github.com/jcaldwell1066/my-context/tree/dev

### Releases
https://github.com/jcaldwell1066/my-context/releases

### Actions/Builds
https://github.com/jcaldwell1066/my-context/actions

---

## Notes

### Session Tracking
This project uses **my-context** for session tracking. Current session: "setup github repo"

View current session:
```bash
my-context show
```

### Environment
- **Platform**: WSL2 (Ubuntu on Windows)
- **Go Version**: 1.x
- **Git**: 2.x

### Best Practices
1. **Always use pre-flight checks** before pushing to GitHub (`./scripts/github-preflight.sh`)
2. **Track work with my-context** (`my-context note "..."`)
3. **Keep dev branch up to date** (push regularly to `origin/dev`)
4. **Use worktrees for parallel work** (safer than branch switching)
5. **Document decisions** (add notes to my-context session)

---

**Document Version**: 1.0.0
**Created**: 2025-10-10 16:45
**Purpose**: Reference for current repository state and common operations
**Update Frequency**: After major changes or when branch structure evolves
