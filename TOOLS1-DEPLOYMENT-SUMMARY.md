# tools1 Deployment - Quick Start

## Current Situation

**Your tools1 Home Directory**: `/home/jcaldwell/`
```
~/my-context-copilot.git    ← Bare repo (already exists)
~/log-analyzer-system.git   ← Another bare repo
~/archive/                   ← Old stuff
~/LOGS/                      ← Various logs
... (many other directories)
```

**Access Methods**:
- WinSCP: File transfers
- PuTTY: SSH command execution
- No direct Git push/pull

**Goal**: Get v1.0.0 onto the existing `~/my-context-copilot.git` bare repo

---

## Fastest Path: One Command

Since `my-context-copilot.git` already exists on tools1, you just need to sync recent commits:

### Step 1: Run Export Script (WSL)

```bash
cd /home/be-dev-agent/projects/my-context-copilot-master
./scripts/deploy-v1.0.0-to-tools1.sh
```

This interactive script will:
- Ask what commit range to export (choose option 1 or 2)
- Create `~/patches-for-tools1-v1.0.0/` with:
  - Patch files (*.patch)
  - Bundle file (backup method)
  - README with instructions
  
**Output**: Ready-to-transfer directory with everything you need

### Step 2: Transfer via WinSCP

1. Open WinSCP
2. Connect to `tools1.shared.accessoticketing.com`
3. Drag and drop: `~/patches-for-tools1-v1.0.0/` → `/home/jcaldwell/patches-incoming/`

### Step 3: Apply via PuTTY

```bash
# SSH to tools1
cd ~/my-context-copilot.git

# Apply patches
git am ~/patches-incoming/*.patch

# Import v1.0.0 tag
git tag v1.0.0 $(cat ~/patches-incoming/tags-v1.0.0.ref | cut -d' ' -f1)

# Verify
git log --oneline -5
git tag -l | grep v1.0.0
```

✅ **Done!** v1.0.0 is now on tools1

---

## Alternative: Bundle Method (If Patches Fail)

If `git am` fails due to divergent history:

### On tools1:
```bash
cd ~/my-context-copilot.git
git bundle verify ~/patches-incoming/my-context-copilot-v1.0.0.bundle
git fetch ~/patches-incoming/my-context-copilot-v1.0.0.bundle 'refs/*:refs/*'
git tag -l | grep v1.0.0
```

---

## Verification

Once deployed, confirm on tools1:

```bash
cd ~/my-context-copilot.git

# Check v1.0.0 exists
git tag -l | grep v1.0.0

# Verify commit
git show v1.0.0 --no-patch

# Should show:
# tag v1.0.0
# Production Release 1.0.0
# ...
```

---

## Clean Up tools1 (Optional)

Your home directory is cluttered. After successful deployment:

```bash
# On tools1
cd ~

# Create archive structure
mkdir -p archive/old-deployments
mkdir -p archive/logs

# Move old stuff (examples)
mv log-analyzer-system-0.0.9 archive/old-deployments/
mv log-analyzer-system-deployed archive/old-deployments/
mv LOGS archive/logs/
mv log.sh archive/

# Keep organized
# ~/my-context-copilot.git     ← Production bare repo
# ~/shell-scripts.git           ← Keep if active
# ~/projects/                   ← Working copies
# ~/archive/                    ← Everything else
# ~/patches-incoming/           ← Temp transfer location
```

**Warning**: Don't delete anything until you verify what's in use!

Check what's active first:
```bash
# See what's being used
ls -ltu ~ | head -20  # Recently accessed files/dirs
```

---

## Documentation Reference

| Document | Purpose |
|----------|---------|
| `TOOLS1-DEPLOYMENT-SUMMARY.md` | This file - quick start |
| `DEPLOY-v1.0.0-TO-TOOLS1.md` | Detailed deployment guide |
| `scripts/README-BARE-REPO-WORKFLOW.md` | Full bare repo workflow |
| `scripts/deploy-v1.0.0-to-tools1.sh` | Automated export script |

---

## Troubleshooting

### "git am fails with conflicts"
→ Use bundle method instead (see Alternative section)

### "Tag v1.0.0 already exists"
→ Delete old tag: `git tag -d v1.0.0` then reapply

### "Bundle does not verify"
→ Re-run export script on WSL and re-transfer

### "Permission denied"
→ Check SSH access: `ssh jcaldwell@tools1.shared.accessoticketing.com "echo OK"`

---

## What You Have Now

✅ Interactive export script: `scripts/deploy-v1.0.0-to-tools1.sh`
✅ Detailed deployment guide: `DEPLOY-v1.0.0-TO-TOOLS1.md`
✅ Full workflow reference: `scripts/README-BARE-REPO-WORKFLOW.md`
✅ Existing bare repo on tools1: `~/my-context-copilot.git`

## What To Do Next

1. **Run**: `./scripts/deploy-v1.0.0-to-tools1.sh`
2. **Transfer**: Use WinSCP to upload patches
3. **Apply**: Use PuTTY to run `git am`
4. **Verify**: Check `git tag -l` on tools1

---

**Estimated Time**: 5-10 minutes
**Created**: 2025-10-10
**Version**: v1.0.0

