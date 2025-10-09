#!/bin/bash
# sync-status.sh
# Show sync status - what changes are pending export/import
# Usage: ./sync-status.sh [--detailed]

set -e

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

DETAILED=false

# Parse arguments
if [[ "$1" == "--detailed" || "$1" == "-d" ]]; then
    DETAILED=true
fi

echo -e "${GREEN}=== Repository Sync Status ===${NC}"

# Verify we're in a git repository
REPO_PATH="$(git rev-parse --show-toplevel 2>/dev/null || echo "")"
if [ -z "$REPO_PATH" ]; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

cd "$REPO_PATH"
REPO_NAME=$(basename "$REPO_PATH")
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo -e "${YELLOW}Repository:${NC} $REPO_NAME"
echo -e "${YELLOW}Branch:${NC} $CURRENT_BRANCH"
echo -e "${YELLOW}Last commit:${NC} $(git log -1 --format='%h - %s (%cr)' 2>/dev/null || echo 'none')"
echo ""

# Check for uncommitted changes
echo -e "${BLUE}▶ Uncommitted Changes${NC}"
if git diff --quiet && git diff --cached --quiet; then
    echo -e "${GREEN}  ✓ Working directory clean${NC}"
else
    # Count changed files
    MODIFIED=$(git diff --name-only | wc -l)
    STAGED=$(git diff --cached --name-only | wc -l)

    echo -e "${YELLOW}  Modified files: $MODIFIED${NC}"
    echo -e "${YELLOW}  Staged files: $STAGED${NC}"

    if [ "$DETAILED" = true ]; then
        echo ""
        echo -e "${YELLOW}  Modified:${NC}"
        git diff --name-status | sed 's/^/    /'
        if [ $STAGED -gt 0 ]; then
            echo ""
            echo -e "${YELLOW}  Staged:${NC}"
            git diff --cached --name-status | sed 's/^/    /'
        fi
    fi
fi

echo ""

# Check for unpushed commits (if remote exists)
echo -e "${BLUE}▶ Unpushed Commits${NC}"

# Try to find a remote
REMOTE=$(git remote 2>/dev/null | head -1 || echo "")

if [ -z "$REMOTE" ]; then
    echo -e "${YELLOW}  No remote configured${NC}"
    echo -e "${YELLOW}  Use export-patches.sh to export commits for manual transfer${NC}"
else
    REMOTE_BRANCH="$REMOTE/$CURRENT_BRANCH"

    if git rev-parse "$REMOTE_BRANCH" >/dev/null 2>&1; then
        AHEAD=$(git rev-list --count "$REMOTE_BRANCH".."$CURRENT_BRANCH" 2>/dev/null || echo 0)
        BEHIND=$(git rev-list --count "$CURRENT_BRANCH".."$REMOTE_BRANCH" 2>/dev/null || echo 0)

        if [ $AHEAD -gt 0 ]; then
            echo -e "${YELLOW}  Ahead of $REMOTE_BRANCH: $AHEAD commit(s)${NC}"
            if [ "$DETAILED" = true ]; then
                git log --oneline "$REMOTE_BRANCH".."$CURRENT_BRANCH" | sed 's/^/    /'
            fi
        fi

        if [ $BEHIND -gt 0 ]; then
            echo -e "${YELLOW}  Behind $REMOTE_BRANCH: $BEHIND commit(s)${NC}"
            if [ "$DETAILED" = true ]; then
                git log --oneline "$CURRENT_BRANCH".."$REMOTE_BRANCH" | sed 's/^/    /'
            fi
        fi

        if [ $AHEAD -eq 0 ] && [ $BEHIND -eq 0 ]; then
            echo -e "${GREEN}  ✓ Up to date with $REMOTE_BRANCH${NC}"
        fi
    else
        echo -e "${YELLOW}  Remote branch $REMOTE_BRANCH not found${NC}"
    fi
fi

echo ""

# Check for untracked files
echo -e "${BLUE}▶ Untracked Files${NC}"
UNTRACKED_COUNT=$(git ls-files --others --exclude-standard | wc -l)

if [ $UNTRACKED_COUNT -eq 0 ]; then
    echo -e "${GREEN}  ✓ No untracked files${NC}"
else
    echo -e "${YELLOW}  Untracked files: $UNTRACKED_COUNT${NC}"
    if [ "$DETAILED" = true ]; then
        git ls-files --others --exclude-standard | sed 's/^/    /'
    fi
fi

echo ""

# Summary and recommendations
echo -e "${BLUE}▶ Recommendations${NC}"

HAS_CHANGES=false

if ! git diff --quiet || ! git diff --cached --quiet; then
    echo -e "${YELLOW}  • Commit your changes:${NC} git add . && git commit -m 'message'"
    HAS_CHANGES=true
fi

if [ -z "$REMOTE" ]; then
    AHEAD_LOCAL=$(git rev-list --count HEAD 2>/dev/null || echo 0)
    if [ $AHEAD_LOCAL -gt 0 ]; then
        echo -e "${YELLOW}  • Export patches:${NC} ./scripts/export-patches.sh HEAD~3..HEAD"
        HAS_CHANGES=true
    fi
else
    if [ "${AHEAD:-0}" -gt 0 ]; then
        echo -e "${YELLOW}  • Export patches:${NC} ./scripts/export-patches.sh $REMOTE_BRANCH..HEAD"
        HAS_CHANGES=true
    fi
fi

if [ "$HAS_CHANGES" = false ]; then
    echo -e "${GREEN}  ✓ Repository is clean and synchronized${NC}"
fi

echo ""
echo -e "${YELLOW}For detailed view:${NC} $0 --detailed"
