#!/bin/bash
# import-patches.sh
# Import and apply Git patches from tools1 or other sources
# Usage: ./import-patches.sh <patch-directory-or-files>

set -e

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Git Patch Importer ===${NC}"

# Check arguments
if [ $# -lt 1 ]; then
    echo -e "${RED}Error: Patch directory or files required${NC}"
    echo ""
    echo "Usage: $0 <patch-dir-or-files>"
    echo ""
    echo "Examples:"
    echo "  $0 ~/patches                      # Apply all patches in directory"
    echo "  $0 ~/patches/*.patch              # Apply specific patches"
    echo "  $0 ~/patches/0001-feature.patch   # Apply single patch"
    exit 1
fi

# Verify we're in a git repository
REPO_PATH="$(git rev-parse --show-toplevel 2>/dev/null || echo "")"
if [ -z "$REPO_PATH" ]; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    echo "Run this script from within a git repository"
    exit 1
fi

cd "$REPO_PATH"
REPO_NAME=$(basename "$REPO_PATH")

echo -e "${YELLOW}Repository:${NC} $REPO_NAME"
echo -e "${YELLOW}Current branch:${NC} $(git rev-parse --abbrev-ref HEAD)"

# Check for uncommitted changes
if ! git diff --quiet || ! git diff --cached --quiet; then
    echo -e "${YELLOW}Warning: You have uncommitted changes${NC}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted"
        exit 1
    fi
fi

# Collect patch files
PATCH_FILES=()

for arg in "$@"; do
    if [ -d "$arg" ]; then
        # Directory - find all .patch files
        while IFS= read -r -d '' file; do
            PATCH_FILES+=("$file")
        done < <(find "$arg" -maxdepth 1 -name "*.patch" -type f -print0 | sort -z)
    elif [ -f "$arg" ]; then
        # Single file
        PATCH_FILES+=("$arg")
    else
        echo -e "${RED}Warning: Not found: $arg${NC}"
    fi
done

# Check if we found any patches
if [ ${#PATCH_FILES[@]} -eq 0 ]; then
    echo -e "${RED}Error: No patch files found${NC}"
    exit 1
fi

echo -e "${YELLOW}Found ${#PATCH_FILES[@]} patch file(s)${NC}"

# Preview patches
echo ""
echo -e "${YELLOW}Patches to apply:${NC}"
for patch in "${PATCH_FILES[@]}"; do
    echo "  - $(basename "$patch")"
done

# Confirm before applying
echo ""
read -p "Apply these patches? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted"
    exit 0
fi

# Apply patches
echo ""
echo -e "${GREEN}Applying patches...${NC}"

SUCCESS_COUNT=0
FAIL_COUNT=0
FAILED_PATCHES=()

for patch in "${PATCH_FILES[@]}"; do
    echo -e "${YELLOW}Applying:${NC} $(basename "$patch")"

    if git am "$patch" 2>&1; then
        ((SUCCESS_COUNT++))
        echo -e "${GREEN}  ✓ Success${NC}"
    else
        ((FAIL_COUNT++))
        FAILED_PATCHES+=("$patch")
        echo -e "${RED}  ✗ Failed${NC}"

        # Ask user what to do
        echo ""
        echo "Options:"
        echo "  1) Skip this patch (git am --skip)"
        echo "  2) Abort all remaining patches (git am --abort)"
        echo "  3) Resolve manually (drop to shell)"
        read -p "Choose (1/2/3): " -n 1 -r choice
        echo ""

        case $choice in
            1)
                git am --skip
                echo "Skipped patch"
                ;;
            2)
                git am --abort
                echo -e "${RED}Aborted patch application${NC}"
                break
                ;;
            3)
                echo ""
                echo -e "${YELLOW}Resolve conflicts manually, then:${NC}"
                echo "  git add <resolved-files>"
                echo "  git am --continue"
                echo ""
                echo "Or skip: git am --skip"
                echo "Or abort: git am --abort"
                exit 1
                ;;
        esac
    fi
done

# Summary
echo ""
echo -e "${GREEN}=== Import Summary ===${NC}"
echo -e "${YELLOW}Total patches:${NC} ${#PATCH_FILES[@]}"
echo -e "${GREEN}Successful:${NC} $SUCCESS_COUNT"
echo -e "${RED}Failed:${NC} $FAIL_COUNT"

if [ $FAIL_COUNT -gt 0 ]; then
    echo ""
    echo -e "${RED}Failed patches:${NC}"
    for patch in "${FAILED_PATCHES[@]}"; do
        echo "  - $(basename "$patch")"
    done
fi

echo ""
echo -e "${GREEN}Current status:${NC}"
git log --oneline -5

echo ""
if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ All patches applied successfully!${NC}"
else
    echo -e "${YELLOW}⚠ Some patches failed. Review the output above.${NC}"
fi
