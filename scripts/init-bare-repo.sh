#!/bin/bash
# init-bare-repo.sh
# Initialize a bare Git repository on tools1 from a bundle file
# This script is designed to run on the remote server (tools1.shared.accessoticketing.com)

set -e  # Exit on error

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Bare Repository Initialization ===${NC}"

# Check if bundle file is provided or find it
BUNDLE_FILE=""
if [ $# -eq 1 ]; then
    BUNDLE_FILE="$1"
elif [ -f "repo.bundle" ]; then
    BUNDLE_FILE="repo.bundle"
else
    echo -e "${RED}Error: No bundle file found${NC}"
    echo "Usage: $0 <bundle-file>"
    echo "   or: Place repo.bundle in current directory"
    exit 1
fi

# Verify bundle file exists
if [ ! -f "$BUNDLE_FILE" ]; then
    echo -e "${RED}Error: Bundle file not found: $BUNDLE_FILE${NC}"
    exit 1
fi

# Extract repository name from bundle filename
REPO_NAME=$(basename "$BUNDLE_FILE" .bundle)
BARE_REPO_DIR="${REPO_NAME}.git"

echo -e "${YELLOW}Bundle file:${NC} $BUNDLE_FILE"
echo -e "${YELLOW}Target directory:${NC} $BARE_REPO_DIR"

# Check if bare repo already exists
if [ -d "$BARE_REPO_DIR" ]; then
    echo -e "${RED}Error: Directory $BARE_REPO_DIR already exists${NC}"
    read -p "Do you want to delete it and recreate? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$BARE_REPO_DIR"
        echo -e "${YELLOW}Removed existing directory${NC}"
    else
        echo "Aborted"
        exit 1
    fi
fi

# Clone the bundle as a bare repository
echo -e "${GREEN}Creating bare repository...${NC}"
git clone --bare "$BUNDLE_FILE" "$BARE_REPO_DIR"

# Set proper permissions
echo -e "${GREEN}Setting permissions...${NC}"
chmod -R 755 "$BARE_REPO_DIR"

# Verify bare repository configuration
cd "$BARE_REPO_DIR"
git config core.bare true

# Get current username and hostname for clone URL
CURRENT_USER=$(whoami)
CURRENT_HOST=$(hostname)

# Display success message with clone instructions
echo ""
echo -e "${GREEN}✓ Bare repository created successfully!${NC}"
echo ""
echo -e "${YELLOW}Repository location:${NC}"
echo "  $(pwd)"
echo ""
echo -e "${YELLOW}Team members can clone with:${NC}"
echo "  git clone ${CURRENT_USER}@${CURRENT_HOST}:~/$(basename $BARE_REPO_DIR)"
echo ""
echo -e "${YELLOW}Or with full path:${NC}"
echo "  git clone ${CURRENT_USER}@${CURRENT_HOST}:$(pwd)"
echo ""
echo -e "${GREEN}Repository structure:${NC}"
tree -L 1 -d . 2>/dev/null || ls -la
echo ""
echo -e "${GREEN}Setup complete!${NC}"
