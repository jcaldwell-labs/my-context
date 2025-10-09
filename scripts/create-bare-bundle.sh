#!/bin/bash
# create-bare-bundle.sh
# Create a single-file bundle for transferring a Git repository to tools1
# Usage: ./create-bare-bundle.sh <repo-name>

set -e  # Exit on error

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Git Bare Repository Bundle Creator ===${NC}"

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/repos.conf"

# Check if repo name provided
if [ $# -lt 1 ]; then
    echo -e "${RED}Error: Repository name required${NC}"
    echo ""
    echo "Usage: $0 <repo-name>"
    echo ""
    if [ -f "$CONFIG_FILE" ]; then
        echo "Available repositories in repos.conf:"
        grep -v '^#' "$CONFIG_FILE" | grep '=' | sed 's/=.*//' | sed 's/^/  - /'
    fi
    exit 1
fi

REPO_NAME="$1"

# Load repository path from config
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}Error: Configuration file not found: $CONFIG_FILE${NC}"
    echo "Create scripts/repos.conf with format: repo_name=/path/to/repo"
    exit 1
fi

REPO_PATH=$(grep "^${REPO_NAME}=" "$CONFIG_FILE" | cut -d'=' -f2-)

if [ -z "$REPO_PATH" ]; then
    echo -e "${RED}Error: Repository '$REPO_NAME' not found in $CONFIG_FILE${NC}"
    echo ""
    echo "Available repositories:"
    grep -v '^#' "$CONFIG_FILE" | grep '=' | sed 's/=.*//' | sed 's/^/  - /'
    exit 1
fi

# Verify repository exists
if [ ! -d "$REPO_PATH" ]; then
    echo -e "${RED}Error: Repository directory not found: $REPO_PATH${NC}"
    exit 1
fi

if [ ! -d "$REPO_PATH/.git" ]; then
    echo -e "${RED}Error: Not a git repository: $REPO_PATH${NC}"
    exit 1
fi

echo -e "${YELLOW}Repository:${NC} $REPO_NAME"
echo -e "${YELLOW}Path:${NC} $REPO_PATH"

# Create temporary working directory
WORK_DIR=$(mktemp -d)
trap "rm -rf $WORK_DIR" EXIT

echo -e "${GREEN}Creating bundle...${NC}"

# Navigate to repository
cd "$REPO_PATH"

# Get current branch info
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
BRANCH_COUNT=$(git branch -a | wc -l)
COMMIT_COUNT=$(git rev-list --all --count)

echo -e "${YELLOW}Current branch:${NC} $CURRENT_BRANCH"
echo -e "${YELLOW}Total branches:${NC} $BRANCH_COUNT"
echo -e "${YELLOW}Total commits:${NC} $COMMIT_COUNT"

# Create git bundle with all refs
BUNDLE_FILE="$WORK_DIR/repo.bundle"
git bundle create "$BUNDLE_FILE" --all

# Verify bundle
echo -e "${GREEN}Verifying bundle...${NC}"
git bundle verify "$BUNDLE_FILE"

# Copy init script to work directory
cp "$SCRIPT_DIR/init-bare-repo.sh" "$WORK_DIR/"
chmod +x "$WORK_DIR/init-bare-repo.sh"

# Create README for the bundle
cat > "$WORK_DIR/README.txt" << 'EOFREADME'
# Bare Repository Bundle

This package contains a Git repository bundle ready for deployment on tools1.

## Quick Setup on tools1

1. Upload this zip file to tools1 via WinSCP
2. Extract: unzip <filename>.zip
3. Run: ./init-bare-repo.sh
4. Share the clone URL with your team

## What's Included

- repo.bundle: Complete Git repository (all branches, tags, history)
- init-bare-repo.sh: Initialization script
- README.txt: This file

## Team Cloning

After initialization on tools1, team members can clone with:

    git clone jcaldwell@tools1.shared.accessoticketing.com:~/<repo-name>.git

Replace <repo-name> with the actual repository name (shown after init).

## Troubleshooting

- If init fails: Check that git is installed on tools1
- Permission denied: Ensure SSH keys are set up for tools1 access
- Clone fails: Verify the path in the clone URL matches the actual location

For detailed workflow documentation, see scripts/README-BARE-REPO-WORKFLOW.md
in the source repository.

EOFREADME

# Get repository metadata
LAST_COMMIT_DATE=$(git log -1 --format=%ci 2>/dev/null || echo "unknown")
LAST_COMMIT_MSG=$(git log -1 --format=%s 2>/dev/null || echo "unknown")

# Add metadata to README
cat >> "$WORK_DIR/README.txt" << EOF

## Repository Metadata

- Repository: $REPO_NAME
- Bundled: $(date '+%Y-%m-%d %H:%M:%S')
- Current branch: $CURRENT_BRANCH
- Total branches: $BRANCH_COUNT
- Total commits: $COMMIT_COUNT
- Last commit: $LAST_COMMIT_DATE
- Last message: $LAST_COMMIT_MSG

EOF

# Create output filename
OUTPUT_DIR="$HOME/bare-bundles"
mkdir -p "$OUTPUT_DIR"
OUTPUT_FILE="$OUTPUT_DIR/${REPO_NAME}-bare-$(date +%Y%m%d-%H%M%S).zip"

# Create zip bundle
echo -e "${GREEN}Creating zip archive...${NC}"
cd "$WORK_DIR"
zip -q "$OUTPUT_FILE" repo.bundle init-bare-repo.sh README.txt

# Display results
echo ""
echo -e "${GREEN}✓ Bundle created successfully!${NC}"
echo ""
echo -e "${YELLOW}Output file:${NC}"
echo "  $OUTPUT_FILE"
echo ""
echo -e "${YELLOW}File size:${NC} $(du -h "$OUTPUT_FILE" | cut -f1)"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Transfer to tools1: Use WinSCP to upload $OUTPUT_FILE"
echo "  2. On tools1: unzip $(basename "$OUTPUT_FILE")"
echo "  3. On tools1: ./init-bare-repo.sh"
echo ""
echo -e "${GREEN}Contents:${NC}"
unzip -l "$OUTPUT_FILE"
echo ""
