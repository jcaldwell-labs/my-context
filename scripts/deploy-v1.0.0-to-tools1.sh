#!/bin/bash
# Quick deployment script for v1.0.0 to tools1
# Exports patches ready for WinSCP transfer

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PATCHES_DIR="$HOME/patches-for-tools1-v1.0.0"

echo "================================================"
echo "Deploy v1.0.0 to tools1 - Patch Export"
echo "================================================"
echo ""

# Change to repo root
cd "$REPO_ROOT"

# Verify we're on the right commit
CURRENT_COMMIT=$(git rev-parse HEAD)
V1_TAG=$(git rev-parse v1.0.0)

echo "Current HEAD: $CURRENT_COMMIT"
echo "v1.0.0 tag:   $V1_TAG"
echo ""

# Ask user what to export
echo "What would you like to export?"
echo ""
echo "1) Last 10 commits + v1.0.0 tag (recommended for initial sync)"
echo "2) Last 20 commits + v1.0.0 tag (if tools1 is very behind)"
echo "3) Only v1.0.0 tag commit (minimal)"
echo "4) Custom range (you specify)"
echo ""
read -p "Choice [1-4]: " choice

case $choice in
    1)
        RANGE="HEAD~10..HEAD"
        ;;
    2)
        RANGE="HEAD~20..HEAD"
        ;;
    3)
        RANGE="v1.0.0~1..v1.0.0"
        ;;
    4)
        read -p "Enter range (e.g., abc123..HEAD): " RANGE
        ;;
    *)
        echo "Invalid choice. Using default: HEAD~10..HEAD"
        RANGE="HEAD~10..HEAD"
        ;;
esac

echo ""
echo "Exporting range: $RANGE"
echo ""

# Create patches directory
mkdir -p "$PATCHES_DIR"

# Export patches
echo "Creating patches..."
git format-patch -o "$PATCHES_DIR" "$RANGE"

# Export v1.0.0 tag as a ref
echo "Exporting v1.0.0 tag..."
git show-ref v1.0.0 > "$PATCHES_DIR/tags-v1.0.0.ref"

# Create a bundle as backup (includes everything)
echo "Creating bundle (backup method)..."
git bundle create "$PATCHES_DIR/my-context-copilot-v1.0.0.bundle" --all

# Create README for tools1
cat > "$PATCHES_DIR/README-APPLY-PATCHES.txt" << 'EOF'
# How to Apply These Updates on tools1

## ⭐ Method 1: Bundle Import (RECOMMENDED for Bare Repos)

This is the simplest method for bare repositories like my-context-copilot.git:

1. Upload this entire directory to tools1:
   WinSCP: Upload patches-for-tools1-v1.0.0/ → /home/jcaldwell/patches-incoming/

2. SSH to tools1 and import:
   cd ~/my-context-copilot.git
   git bundle verify ~/patches-incoming/my-context-copilot-v1.0.0.bundle
   git fetch ~/patches-incoming/my-context-copilot-v1.0.0.bundle 'refs/heads/*:refs/heads/*' 'refs/tags/*:refs/tags/*'

3. Verify:
   git tag -l | grep v1.0.0
   git show v1.0.0 --no-patch
   git log --oneline -5

✅ Done! All commits, branches, and tags are now synced.

## Method 2: Patches (Only for Normal Repos with Working Trees)

⚠️  NOTE: This does NOT work on bare repositories (*.git directories).
Only use if you have a normal working copy.

1. Upload this entire directory to tools1

2. SSH to tools1 (must have working tree):
   cd ~/my-context-copilot  # NOT the .git directory!
   git am ~/patches-incoming/*.patch
   
3. Verify:
   git log --oneline -5
   
## Troubleshooting

- If "this operation must be run in a work tree" → You're in a bare repo, use Method 1 (bundle)
- If "Repository does not exist" → Bundle is corrupt, re-export
- If "Tag already exists" → Delete old tag: git tag -d v1.0.0
- If bundle verify fails → Re-download bundle (transfer error)

---
Generated: $(date)
From: $(git config user.name) <$(git config user.email)>
Commit: $(git rev-parse HEAD)
EOF

# Count patches
PATCH_COUNT=$(ls -1 "$PATCHES_DIR"/*.patch 2>/dev/null | wc -l)

echo ""
echo "================================================"
echo "✅ Export Complete!"
echo "================================================"
echo ""
echo "📁 Location: $PATCHES_DIR"
echo "📊 Patches: $PATCH_COUNT files"
echo "📦 Bundle: my-context-copilot-v1.0.0.bundle (backup)"
echo "📋 README: README-APPLY-PATCHES.txt (instructions)"
echo ""
echo "Next Steps:"
echo "1. Open WinSCP and connect to tools1.shared.accessoticketing.com"
echo "2. Upload entire directory: $PATCHES_DIR → /home/jcaldwell/patches-incoming/"
echo "3. SSH to tools1 via PuTTY"
echo "4. Follow instructions in README-APPLY-PATCHES.txt"
echo ""
echo "Full guide: $REPO_ROOT/DEPLOY-v1.0.0-TO-TOOLS1.md"
echo ""

# List patches for review
echo "Patches to be transferred:"
ls -1 "$PATCHES_DIR"/*.patch 2>/dev/null | head -10
if [ $PATCH_COUNT -gt 10 ]; then
    echo "... and $((PATCH_COUNT - 10)) more"
fi
echo ""

# Show recent commits that will be synced
echo "Commits included in this export:"
git log --oneline "$RANGE" | head -10
echo ""

echo "Ready for WinSCP transfer!"

