# Deploy v1.0.0 to tools1 - Quick Guide

**Status**: Ready to deploy
**Target**: tools1.shared.accessoticketing.com:/home/jcaldwell/my-context-copilot.git
**Method**: Patch-based sync (bare repo already exists)

---

## Quick Deployment Steps

### 1. Export Patches (WSL - do this now)

```bash
cd /home/be-dev-agent/projects/my-context-copilot-master

# Export all commits since the last sync to tools1
# Adjust the range as needed (e.g., HEAD~10..HEAD if you know the last sync point)
./scripts/export-patches.sh origin/master..HEAD
```

**Output**: Creates `~/patches/my-context-copilot-patches-<timestamp>/`

**Alternative if you want to export specific commits**:
```bash
# Export just v1.0.0 and recent commits
git log --oneline -20  # Review recent commits
./scripts/export-patches.sh <last-synced-commit>..v1.0.0
```

---

### 2. Transfer Patches (WinSCP)

1. **Open WinSCP**
2. **Connect to**: `tools1.shared.accessoticketing.com`
3. **Navigate to**: `/home/jcaldwell/`
4. **Upload**: Local `~/patches/my-context-copilot-patches-<timestamp>/` → tools1 `~/patches-incoming/`

---

### 3. Apply Patches (PuTTY SSH to tools1)

```bash
# SSH via PuTTY to tools1
ssh jcaldwell@tools1.shared.accessoticketing.com

# Navigate to bare repo
cd ~/my-context-copilot.git

# Check current state
git log --oneline -5
git tag -l

# Apply patches
git am ~/patches-incoming/*.patch

# Verify v1.0.0 tag arrived
git tag -l | grep v1.0.0
git log --oneline -1 v1.0.0
```

---

## Alternative: Fresh Bundle (if patches fail)

If patch application fails or you want a clean slate:

### 1. Create Fresh Bundle (WSL)

```bash
cd /home/be-dev-agent/projects/my-context-copilot-master
./scripts/create-bare-bundle.sh my-context-copilot
```

**Output**: `~/bare-bundles/my-context-copilot-bare-<timestamp>.zip`

### 2. Transfer Bundle (WinSCP)

Upload `~/bare-bundles/my-context-copilot-bare-*.zip` → tools1 `~/`

### 3. Reinitialize (PuTTY SSH to tools1)

```bash
ssh jcaldwell@tools1.shared.accessoticketing.com

# Backup old bare repo (just in case)
mv ~/my-context-copilot.git ~/my-context-copilot.git.backup

# Extract and init new bundle
cd ~
unzip my-context-copilot-bare-*.zip
./init-bare-repo.sh

# Verify v1.0.0
cd ~/my-context-copilot.git
git tag -l | grep v1.0.0
```

---

## Verification Checklist

Once deployed, verify on tools1:

```bash
cd ~/my-context-copilot.git

# 1. Check v1.0.0 tag exists
git tag -l | grep v1.0.0

# 2. Verify tag points to correct commit
git show v1.0.0 --no-patch --format="%H %s"

# 3. List all branches
git branch -a

# 4. Check recent commits
git log --oneline -10
```

Expected output for v1.0.0:
- Tag: `v1.0.0`
- Commit: `5a58979` (or similar)
- Message: "docs: update ACTIVE-THREADS-STATUS to reflect Thread 1 completion"

---

## Troubleshooting

### "Cannot apply patches - repository has diverged"

**Solution**: Use fresh bundle method instead.

### "Tag v1.0.0 already exists"

**Solution**: Delete old tag first:
```bash
cd ~/my-context-copilot.git
git tag -d v1.0.0
git am ~/patches-incoming/*.patch
```

### "Permission denied" on git am

**Solution**: Bare repos need special handling:
```bash
# Instead of git am, use git fetch from bundle
cd ~/my-context-copilot.git
git fetch ~/patches-incoming/my-context.bundle 'refs/*:refs/*'
```

---

## Clean Up tools1 Home Directory (Optional)

After successful deployment, consider organizing:

```bash
# On tools1
cd ~
mkdir -p archive/old-deployments

# Move old versions
mv log-analyzer-system-* archive/old-deployments/
mv LOGS archive/
mv log.sh archive/

# Keep organized structure
# ~/my-context-copilot.git     - Active bare repo
# ~/projects/                  - Working copies if needed
# ~/archive/                   - Old stuff
# ~/patches-incoming/          - Temp transfer location
```

---

## Post-Deployment

### Team Members Can Now Clone v1.0.0

```bash
git clone jcaldwell@tools1.shared.accessoticketing.com:~/my-context-copilot.git
cd my-context-copilot
git checkout v1.0.0  # Or stay on master
```

### Future Updates

For subsequent releases, just use patch sync:

```bash
# WSL
cd /home/be-dev-agent/projects/my-context-copilot-master
./scripts/export-patches.sh v1.0.0..v1.1.0

# Transfer → Apply on tools1
```

---

## Quick Reference

| Task | WSL Command | tools1 Command |
|------|-------------|----------------|
| Export patches | `./scripts/export-patches.sh HEAD~10..HEAD` | - |
| Check sync status | `./scripts/sync-status.sh` | - |
| Apply patches | - | `git am ~/patches-incoming/*.patch` |
| Verify version | - | `git tag -l \| grep v1.0.0` |

---

**Full Workflow Documentation**: `scripts/README-BARE-REPO-WORKFLOW.md`
**Created**: 2025-10-10
**For**: v1.0.0 production release

