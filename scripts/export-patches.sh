#!/bin/bash
# export-patches.sh
# Export Git commits as patch files for transfer to tools1
# Usage: ./export-patches.sh [<from>..<to>] [--output-dir <dir>]

set -e

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Git Patch Exporter ===${NC}"

# Default values
RANGE=""
OUTPUT_DIR="$HOME/patches"
REPO_PATH="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [<from>..<to>] [--output-dir <dir>]"
            echo ""
            echo "Examples:"
            echo "  $0                          # Export uncommitted changes"
            echo "  $0 HEAD~3..HEAD             # Export last 3 commits"
            echo "  $0 v1.0..HEAD               # Export commits since tag v1.0"
            echo "  $0 origin/main..HEAD        # Export commits ahead of origin/main"
            echo "  $0 HEAD~5..HEAD --output-dir /tmp/patches"
            echo ""
            exit 0
            ;;
        *)
            RANGE="$1"
            shift
            ;;
    esac
done

# Verify we're in a git repository
if [ ! -d "$REPO_PATH/.git" ]; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

cd "$REPO_PATH"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Get repository name
REPO_NAME=$(basename "$REPO_PATH")

# Create timestamped subdirectory
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
PATCH_DIR="$OUTPUT_DIR/${REPO_NAME}-patches-$TIMESTAMP"
mkdir -p "$PATCH_DIR"

echo -e "${YELLOW}Repository:${NC} $REPO_NAME"
echo -e "${YELLOW}Output directory:${NC} $PATCH_DIR"

# Export patches
if [ -z "$RANGE" ]; then
    # No range specified - export uncommitted changes
    echo -e "${YELLOW}Exporting uncommitted changes...${NC}"

    if git diff --quiet && git diff --cached --quiet; then
        echo -e "${RED}No changes to export${NC}"
        rmdir "$PATCH_DIR"
        exit 0
    fi

    # Export staged and unstaged changes
    git diff HEAD > "$PATCH_DIR/uncommitted-changes.patch"
    PATCH_COUNT=1
else
    # Range specified - export commits
    echo -e "${YELLOW}Range:${NC} $RANGE"

    # Verify range is valid
    if ! git rev-parse "$RANGE" >/dev/null 2>&1; then
        echo -e "${RED}Error: Invalid range: $RANGE${NC}"
        rmdir "$PATCH_DIR"
        exit 1
    fi

    # Count commits in range
    PATCH_COUNT=$(git rev-list --count $RANGE 2>/dev/null || echo 0)

    if [ "$PATCH_COUNT" -eq 0 ]; then
        echo -e "${RED}No commits in range: $RANGE${NC}"
        rmdir "$PATCH_DIR"
        exit 0
    fi

    echo -e "${YELLOW}Commits to export:${NC} $PATCH_COUNT"

    # Export commits as patches
    git format-patch -o "$PATCH_DIR" "$RANGE"
fi

# Create metadata file
cat > "$PATCH_DIR/metadata.txt" << EOF
Repository: $REPO_NAME
Exported: $(date '+%Y-%m-%d %H:%M:%S')
Range: ${RANGE:-"uncommitted changes"}
Patch count: $PATCH_COUNT
Current branch: $(git rev-parse --abbrev-ref HEAD)
Current commit: $(git rev-parse --short HEAD)
EOF

# List exported patches
echo ""
echo -e "${GREEN}✓ Exported $PATCH_COUNT patch(es)${NC}"
echo ""
echo -e "${YELLOW}Patches:${NC}"
ls -1 "$PATCH_DIR"/*.patch 2>/dev/null | while read patch; do
    echo "  - $(basename "$patch")"
done

echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Transfer patches to tools1:"
echo "     WinSCP: Upload $PATCH_DIR/*.patch"
echo "  2. On tools1, apply patches:"
echo "     cd ~/$(basename $REPO_NAME).git"
echo "     git am ~/patches/*.patch"
echo ""
echo -e "${YELLOW}Location:${NC} $PATCH_DIR"
