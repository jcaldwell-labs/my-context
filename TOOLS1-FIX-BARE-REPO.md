# Fix: Applying Patches to Bare Repository

## The Problem

```bash
git am ~/patches-incoming/*.patch
fatal: this operation must be run in a work tree
```

**Why**: Bare repositories don't have a working tree, so `git am` doesn't work directly.

---

## Solution: Use the Bundle Instead

The export script included a bundle file specifically for this case!

### On tools1 (PuTTY):

```bash
cd ~/my-context-copilot.git

# Verify the bundle
git bundle verify ~/patches-incoming/my-context-copilot-v1.0.0.bundle

# Fetch all refs from the bundle (includes commits + tags)
git fetch ~/patches-incoming/my-context-copilot-v1.0.0.bundle 'refs/heads/*:refs/heads/*' 'refs/tags/*:refs/tags/*'

# Verify v1.0.0 tag is now present
git tag -l | grep v1.0.0

# Check recent commits
git log --oneline -10 master

# Verify the tag points to the right commit
git show v1.0.0 --no-patch
```

✅ **Done!** v1.0.0 is now on tools1

---

## Alternative: Use a Temporary Worktree

If you need to apply patches in the future, you can use a temporary worktree:

```bash
cd ~/my-context-copilot.git

# Create temporary worktree
git worktree add /tmp/my-context-temp master

# Apply patches in the worktree
cd /tmp/my-context-temp
git am ~/patches-incoming/*.patch

# Push changes back to bare repo
git push origin master
git push origin --tags

# Clean up
cd ~
git worktree remove /tmp/my-context-temp
```

---

## Verification

Confirm everything is correct:

```bash
cd ~/my-context-copilot.git

# 1. Check tag exists
git tag -l

# 2. Verify tag commit
git show v1.0.0 --no-patch --format="%H %s"

# Expected output:
# tag v1.0.0
# Production Release 1.0.0
# ...
# commit 8bf20d1...
# docs: v1.0.0 release complete! 🎉

# 3. Check all branches are up to date
git branch -a

# 4. Verify latest commits
git log --oneline --all --graph -10
```

---

## Why This Happens

**Bare Repository**: `my-context-copilot.git/`
- Contains only Git database (objects, refs, etc.)
- No working files
- Used for sharing/cloning
- Cannot use: `git am`, `git apply`, `git checkout`

**Normal Repository**: `my-context-copilot/`
- Has working tree (actual files)
- Has .git/ directory
- Can use all Git commands

---

## Next Time: Bundle is Easier

For bare repos, bundles are simpler than patches:

**WSL Export**:
```bash
cd ~/projects/my-context-copilot-master
git bundle create my-context-v1.1.0.bundle --all
# Transfer via WinSCP
```

**tools1 Import**:
```bash
cd ~/my-context-copilot.git
git fetch ~/my-context-v1.1.0.bundle 'refs/*:refs/*'
```

No worktree needed!

---

**Summary**: Use `git fetch` with the bundle, not `git am` with patches.

