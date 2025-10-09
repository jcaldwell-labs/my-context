# Bare Repository Bundle Workflow for tools1

Complete guide for sharing Git repositories on tools1 via bare repository bundles and patch-based synchronization.

## Overview

This workflow enables Git repository sharing on tools1.shared.accessoticketing.com when direct push/pull isn't available. It consists of:

1. **Bundle Creation**: Package repository as a single zip file
2. **Bare Repository Setup**: Initialize on tools1 for team cloning
3. **Patch-Based Sync**: Bidirectional updates via patch files

## Prerequisites

### On WSL/Local Machine
- Git installed
- WinSCP or scp for file transfers
- SSH access to tools1 (VPN required)

### On tools1
- Git installed (verify: `git --version`)
- SSH key configured for team access
- Home directory: `/home/jcaldwell/`

## Quick Start

### First-Time Setup (my-context-copilot)

1. **Create Bundle** (WSL):
   ```bash
   cd ~/projects/my-context-copilot
   ./scripts/create-bare-bundle.sh my-context-copilot
   ```

2. **Transfer to tools1** (WinSCP):
   - Connect to: tools1.shared.accessoticketing.com
   - Navigate to: `/home/jcaldwell/`
   - Upload: `~/bare-bundles/my-context-copilot-bare-YYYYMMDD-HHMMSS.zip`

3. **Initialize on tools1** (SSH):
   ```bash
   ssh jcaldwell@tools1.shared.accessoticketing.com
   cd ~
   unzip my-context-copilot-bare-YYYYMMDD-HHMMSS.zip
   ./init-bare-repo.sh
   ```

4. **Team Clones** (Team members):
   ```bash
   git clone jcaldwell@tools1.shared.accessoticketing.com:~/my-context-copilot.git
   ```

## Scripts Reference

### repos.conf

Configuration file listing available repositories:

```bash
# Format: repo_name=absolute_path
my-context-copilot=/home/be-dev-agent/projects/my-context-copilot
deb-sanity=/home/be-dev-agent/projects/deb-sanity
```

**Usage**: Edit to add new repositories.

### create-bare-bundle.sh

Creates a single-file bundle for transfer.

**Usage**:
```bash
./scripts/create-bare-bundle.sh <repo-name>

# Example
./scripts/create-bare-bundle.sh my-context-copilot
```

**Output**: `~/bare-bundles/<repo-name>-bare-<timestamp>.zip`

**Contains**:
- `repo.bundle` - Complete Git history (all branches, tags)
- `init-bare-repo.sh` - Initialization script
- `README.txt` - Setup instructions

### init-bare-repo.sh

Initializes bare repository on tools1 from bundle.

**Usage** (on tools1):
```bash
./init-bare-repo.sh [bundle-file]

# Auto-detects repo.bundle in current directory
./init-bare-repo.sh

# Or specify bundle file
./init-bare-repo.sh repo.bundle
```

**Creates**: `<repo-name>.git/` directory structure matching:
```
my-context-copilot.git/
├── branches
├── hooks
├── info
├── objects
└── refs
```

### export-patches.sh

Export commits as patch files for manual transfer.

**Usage**:
```bash
# Export uncommitted changes
./scripts/export-patches.sh

# Export last 3 commits
./scripts/export-patches.sh HEAD~3..HEAD

# Export commits since tag
./scripts/export-patches.sh v1.0..HEAD

# Custom output directory
./scripts/export-patches.sh HEAD~5..HEAD --output-dir /tmp/patches
```

**Output**: `~/patches/<repo-name>-patches-<timestamp>/`

### import-patches.sh

Apply patch files from tools1 or other sources.

**Usage**:
```bash
# Apply all patches from directory
./scripts/import-patches.sh ~/patches/

# Apply specific patches
./scripts/import-patches.sh ~/patches/*.patch

# Apply single patch
./scripts/import-patches.sh ~/patches/0001-feature.patch
```

**Interactive**: Prompts for confirmation and handles conflicts.

### sync-status.sh

Show current sync status and pending changes.

**Usage**:
```bash
# Summary view
./scripts/sync-status.sh

# Detailed view
./scripts/sync-status.sh --detailed
```

**Shows**:
- Uncommitted changes
- Unpushed commits
- Untracked files
- Recommended actions

## Common Workflows

### Adding a New Repository (deb-sanity)

1. **Update Configuration**:
   ```bash
   # Edit scripts/repos.conf
   echo "deb-sanity=/home/be-dev-agent/projects/deb-sanity" >> scripts/repos.conf
   ```

2. **Create and Transfer Bundle**:
   ```bash
   cd ~/projects/deb-sanity
   ~/projects/my-context-copilot/scripts/create-bare-bundle.sh deb-sanity
   # Transfer ~/bare-bundles/deb-sanity-bare-*.zip via WinSCP
   ```

3. **Initialize on tools1**:
   ```bash
   ssh jcaldwell@tools1.shared.accessoticketing.com
   unzip deb-sanity-bare-*.zip
   ./init-bare-repo.sh
   ```

### Syncing Changes: WSL → tools1

**Scenario**: You made commits on WSL and want to push to tools1.

1. **Check Status** (WSL):
   ```bash
   cd ~/projects/my-context-copilot
   ./scripts/sync-status.sh
   ```

2. **Export Patches** (WSL):
   ```bash
   # Export commits since last sync
   ./scripts/export-patches.sh HEAD~5..HEAD
   # Or export uncommitted changes
   ./scripts/export-patches.sh
   ```

3. **Transfer Patches** (WinSCP):
   - Upload: `~/patches/my-context-copilot-patches-*/` → tools1

4. **Apply Patches** (tools1):
   ```bash
   ssh jcaldwell@tools1.shared.accessoticketing.com
   cd ~/my-context-copilot.git
   git am ~/patches/my-context-copilot-patches-*/*.patch
   ```

### Syncing Changes: tools1 → WSL

**Scenario**: Team member pushed to tools1, you want to pull changes.

1. **Export Patches** (tools1):
   ```bash
   ssh jcaldwell@tools1.shared.accessoticketing.com
   cd ~/my-context-copilot.git
   git format-patch -o ~/patches-outgoing/ <last-synced-commit>..HEAD
   ```

2. **Download Patches** (WinSCP):
   - Download: tools1:`~/patches-outgoing/` → WSL:`~/patches-incoming/`

3. **Apply Patches** (WSL):
   ```bash
   cd ~/projects/my-context-copilot
   ./scripts/import-patches.sh ~/patches-incoming/
   ```

### Team Member Workflow

**Scenario**: Team member wants to clone and contribute.

1. **Clone from tools1**:
   ```bash
   git clone jcaldwell@tools1.shared.accessoticketing.com:~/my-context-copilot.git
   cd my-context-copilot
   ```

2. **Make Changes**:
   ```bash
   git checkout -b feature-branch
   # ... make changes ...
   git add .
   git commit -m "Add feature"
   ```

3. **Push to tools1**:
   ```bash
   git push origin feature-branch
   ```

4. **Create Merge Request** (via patches if needed):
   ```bash
   git format-patch origin/main..feature-branch -o ~/patches/
   # Email patches to maintainer
   ```

## Troubleshooting

### "Bundle does not exist" on tools1

**Problem**: `init-bare-repo.sh` can't find bundle file.

**Solution**:
```bash
# Check file is in current directory
ls -la *.bundle

# Or specify explicitly
./init-bare-repo.sh /path/to/repo.bundle
```

### "Permission denied" when team clones

**Problem**: SSH permissions not set correctly.

**Solution** (on tools1):
```bash
chmod 755 ~/my-context-copilot.git
chmod -R 755 ~/my-context-copilot.git/hooks
```

### Patch application fails with conflicts

**Problem**: `git am` fails due to merge conflicts.

**Solution**:
```bash
# Option 1: Resolve manually
git am --show-current-patch
# Edit conflicted files
git add <resolved-files>
git am --continue

# Option 2: Skip patch
git am --skip

# Option 3: Abort and try later
git am --abort
```

### "Not a git repository" error

**Problem**: Running scripts outside a Git repository.

**Solution**:
```bash
# Navigate to repository first
cd ~/projects/my-context-copilot
./scripts/sync-status.sh
```

### Bundle is too large

**Problem**: Repository history is huge (>100MB).

**Solution**:
```bash
# Create shallow bundle (recent commits only)
git bundle create repo-shallow.bundle HEAD~100..HEAD

# Or use sparse clone (requires Git 2.25+)
git clone --depth=50 --bare original-repo/.git shallow.git
```

## Advanced Topics

### Automating Transfers with scp

Instead of WinSCP, use scp in scripts:

```bash
# Upload bundle
scp ~/bare-bundles/my-context-copilot-bare-*.zip \
    jcaldwell@tools1.shared.accessoticketing.com:~/

# Download patches
scp -r jcaldwell@tools1.shared.accessoticketing.com:~/patches-outgoing/ \
    ~/patches-incoming/
```

### Setting up SSH Keys

For passwordless transfers:

```bash
# Generate key (WSL)
ssh-keygen -t ed25519 -C "jcaldwell@tools1"

# Copy to tools1
ssh-copy-id jcaldwell@tools1.shared.accessoticketing.com

# Test
ssh jcaldwell@tools1.shared.accessoticketing.com "echo 'Connected!'"
```

### Bare Repository Structure

Understanding the bare repo layout:

```
my-context-copilot.git/
├── HEAD                 # Current branch pointer
├── config              # Repository configuration
├── description         # Repository description
├── hooks/              # Server-side hooks
├── info/               # Additional info
├── objects/            # Git objects (commits, trees, blobs)
│   ├── pack/           # Packed objects
│   └── info/
└── refs/               # Branch and tag references
    ├── heads/          # Branches
    └── tags/           # Tags
```

### Pre-receive Hooks (Optional)

Add validation on tools1 before accepting pushes:

```bash
# On tools1
cd ~/my-context-copilot.git/hooks
cat > pre-receive << 'EOF'
#!/bin/bash
# Reject pushes to main branch
while read oldrev newrev refname; do
    if [ "$refname" = "refs/heads/main" ]; then
        echo "Error: Direct pushes to main are not allowed"
        exit 1
    fi
done
EOF
chmod +x pre-receive
```

## FAQ

**Q: Can I use this for private repositories?**
A: Yes, as long as you have SSH access to tools1 and proper permissions.

**Q: How do I delete a bare repository on tools1?**
A: `rm -rf ~/my-context-copilot.git` (warning: irreversible!)

**Q: Can multiple people push simultaneously?**
A: Yes, Git handles concurrent pushes. Conflicts are resolved on pull.

**Q: What if tools1 goes down?**
A: Your local repository has full history. Recreate bare repo when tools1 is back.

**Q: How do I update an existing bare repo?**
A: Don't recreate. Use patches or push directly if SSH access allows.

## See Also

- [Git Bundle Documentation](https://git-scm.com/docs/git-bundle)
- [Git Bare Repositories](https://git-scm.com/book/en/v2/Git-on-the-Server-Getting-Git-on-a-Server)
- [Git Format-Patch](https://git-scm.com/docs/git-format-patch)
- [Git AM (Apply Patches)](https://git-scm.com/docs/git-am)

## Support

For issues or questions:
1. Check this README and troubleshooting section
2. Run `./scripts/sync-status.sh --detailed` for diagnostics
3. Review script output for error messages
4. Contact repository maintainer

---

**Version**: 1.0
**Last Updated**: 2025-10-09
**Maintainer**: jcaldwell@tools1
