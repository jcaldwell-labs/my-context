#!/bin/bash
# Safe cleanup script for tools1 home directory
# Run this on tools1 via SSH/PuTTY
# Usage: bash tools1-cleanup.sh

set -e

echo "================================================"
echo "tools1 Home Directory Cleanup Script"
echo "================================================"
echo ""
echo "This script will organize your home directory by:"
echo "  - Creating archive structure"
echo "  - Moving old/unused directories to archive"
echo "  - Preserving active bare repos and projects"
echo "  - Cleaning up temporary files"
echo ""
echo "⚠️  SAFETY FIRST: This script shows what it will do BEFORE doing it."
echo ""

# Ensure we're in home directory
cd ~
HOMEDIR=$(pwd)
echo "Working in: $HOMEDIR"
echo ""

# Create archive structure (safe, won't overwrite)
echo "Step 1: Creating archive structure..."
mkdir -p archive/old-deployments
mkdir -p archive/logs
mkdir -p archive/misc
mkdir -p archive/old-repos
mkdir -p archive/temp
echo "✅ Archive directories created"
echo ""

# Function to safely move with confirmation
safe_move() {
    local source="$1"
    local dest="$2"
    local description="$3"
    
    if [ -e "$source" ]; then
        echo "📦 Will move: $source → $dest"
        echo "   ($description)"
        return 0
    else
        echo "⏭️  Skipping: $source (doesn't exist)"
        return 1
    fi
}

# Build move plan
echo "Step 2: Planning moves..."
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "OLD LOG ANALYZER DEPLOYMENTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

MOVE_PLAN=()

# Old log analyzer versions
if safe_move "log-analyzer-system-0.0.9" "archive/old-deployments/" "Old version 0.0.9"; then
    MOVE_PLAN+=("mv ~/log-analyzer-system-0.0.9 ~/archive/old-deployments/")
fi

if safe_move "log-analyzer-system-deployed" "archive/old-deployments/" "Old deployed version"; then
    MOVE_PLAN+=("mv ~/log-analyzer-system-deployed ~/archive/old-deployments/")
fi

if safe_move "archive/log-analyzer-system" "archive/old-deployments/" "Old archive copy"; then
    MOVE_PLAN+=("mv ~/archive/log-analyzer-system ~/archive/old-deployments/")
fi

if safe_move "archive/log-analyzer-system.cloned" "archive/old-deployments/" "Old cloned copy"; then
    MOVE_PLAN+=("mv ~/archive/log-analyzer-system.cloned ~/archive/old-deployments/")
fi

if safe_move "archive/log-analyzer-system-old" "archive/old-deployments/" "Very old copy"; then
    MOVE_PLAN+=("mv ~/archive/log-analyzer-system-old ~/archive/old-deployments/")
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "LOGS AND SCRIPTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if safe_move "LOGS" "archive/logs/" "Old log files"; then
    MOVE_PLAN+=("mv ~/LOGS ~/archive/logs/")
fi

if safe_move "log.sh" "archive/logs/" "Old log script"; then
    MOVE_PLAN+=("mv ~/log.sh ~/archive/logs/")
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "MISCELLANEOUS DIRECTORIES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if safe_move "MISC" "archive/misc/" "Miscellaneous files"; then
    MOVE_PLAN+=("mv ~/MISC ~/archive/misc/")
fi

if safe_move "TXT" "archive/misc/" "Text files"; then
    MOVE_PLAN+=("mv ~/TXT ~/archive/misc/")
fi

if safe_move "SQL" "archive/misc/" "SQL files"; then
    MOVE_PLAN+=("mv ~/SQL ~/archive/misc/")
fi

if safe_move "LEVEL3" "archive/misc/" "LEVEL3 directory"; then
    MOVE_PLAN+=("mv ~/LEVEL3 ~/archive/misc/")
fi

if safe_move "cloudformation" "archive/misc/" "CloudFormation templates"; then
    MOVE_PLAN+=("mv ~/cloudformation ~/archive/misc/")
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "OLD/DUPLICATE REPOS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if safe_move "repo.git" "archive/old-repos/" "Old generic bare repo"; then
    MOVE_PLAN+=("mv ~/repo.git ~/archive/old-repos/")
fi

if safe_move "projects/repo.git" "archive/old-repos/" "Duplicate bare repo"; then
    MOVE_PLAN+=("mv ~/projects/repo.git ~/archive/old-repos/")
fi

if safe_move "projects/repo" "archive/old-repos/" "Old working copy"; then
    MOVE_PLAN+=("mv ~/projects/repo ~/archive/old-repos/")
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEMPORARY/CLEANUP"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if safe_move "node_modules" "archive/temp/" "Node modules (orphaned)"; then
    MOVE_PLAN+=("mv ~/node_modules ~/archive/temp/")
fi

if safe_move "data" "archive/temp/" "Old data directory"; then
    MOVE_PLAN+=("mv ~/data ~/archive/temp/")
fi

if safe_move "github" "archive/temp/" "Old GitHub checkout"; then
    MOVE_PLAN+=("mv ~/github ~/archive/temp/")
fi

if safe_move "patches-incoming" "archive/temp/" "Deployment patches (already applied)"; then
    MOVE_PLAN+=("mv ~/patches-incoming ~/archive/temp/")
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "WILL BE PRESERVED (not moved)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ ~/my-context-copilot.git     (production bare repo)"
echo "✅ ~/shell-scripts.git           (bare repo)"
echo "✅ ~/shell-scripts/              (working copy)"
echo "✅ ~/log-analyzer-system.git     (bare repo)"
echo "✅ ~/projects/                   (if contains active work)"
echo "✅ ~/PHD/                        (project directories)"
echo "✅ ~/archive/                    (organized archives)"
echo ""

# Count moves
MOVE_COUNT=${#MOVE_PLAN[@]}

if [ $MOVE_COUNT -eq 0 ]; then
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✨ Nothing to move! Directory already clean."
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 0
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "SUMMARY: $MOVE_COUNT items will be moved"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Show disk space before
echo "Current disk usage:"
du -sh ~/* 2>/dev/null | sort -h | tail -10
echo ""

# Confirmation prompt
read -p "⚠️  Proceed with cleanup? (type 'yes' to confirm): " confirm

if [ "$confirm" != "yes" ]; then
    echo ""
    echo "❌ Cleanup cancelled. No changes made."
    exit 0
fi

echo ""
echo "Step 3: Executing moves..."
echo ""

# Execute moves
MOVED=0
FAILED=0

for cmd in "${MOVE_PLAN[@]}"; do
    echo "Running: $cmd"
    if eval "$cmd" 2>/dev/null; then
        ((MOVED++))
        echo "  ✅ Success"
    else
        ((FAILED++))
        echo "  ❌ Failed"
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "CLEANUP COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Moved: $MOVED items"
if [ $FAILED -gt 0 ]; then
    echo "❌ Failed: $FAILED items"
fi
echo ""

# Show new structure
echo "New home directory structure:"
echo ""
tree -L 1 -d ~ 2>/dev/null || ls -1d ~/*/

echo ""
echo "Archive contents:"
ls -1 ~/archive/

echo ""
echo "Disk usage by directory:"
du -sh ~/* 2>/dev/null | sort -h | tail -15

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "NEXT STEPS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Verify important files are accessible:"
echo "   cd ~/my-context-copilot.git && git log -1"
echo "   cd ~/shell-scripts && ls"
echo ""
echo "2. If everything looks good after a week, delete archives:"
echo "   rm -rf ~/archive/old-deployments"
echo "   rm -rf ~/archive/temp"
echo ""
echo "3. Keep organized structure:"
echo "   ~/[project].git/     → Bare repos"
echo "   ~/projects/          → Working copies"
echo "   ~/archive/           → Old stuff (temporary)"
echo ""
echo "✨ Cleanup complete! Your home directory is now organized."

