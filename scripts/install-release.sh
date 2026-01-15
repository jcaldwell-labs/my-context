#!/usr/bin/env bash
# Official installation script for my-context
# Downloads latest release from GitHub and installs
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/jcaldwell-labs/my-context/main/scripts/install-release.sh | bash
#   curl -fsSL ... | bash -s -- --system   # Install to /usr/local/bin (requires sudo)
#   curl -fsSL ... | bash -s -- --version v3.3.0  # Install specific version

set -e

REPO="jcaldwell-labs/my-context"
BINARY_NAME="my-context"
INSTALL_MODE="user"  # user or system
VERSION="latest"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --system)
            INSTALL_MODE="system"
            shift
            ;;
        --user)
            INSTALL_MODE="user"
            shift
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: install-release.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --user       Install to ~/.local/bin (default, no sudo required)"
            echo "  --system     Install to /usr/local/bin (requires sudo)"
            echo "  --version    Install specific version (e.g., v3.3.0)"
            echo "  -h, --help   Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Set install directory based on mode
if [ "$INSTALL_MODE" = "system" ]; then
    INSTALL_DIR="/usr/local/bin"
    SUDO_CMD="sudo"
else
    INSTALL_DIR="$HOME/.local/bin"
    SUDO_CMD=""
fi

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux*)
        PLATFORM="linux"
        ;;
    darwin*)
        PLATFORM="darwin"
        ;;
    mingw*|msys*|cygwin*)
        PLATFORM="windows"
        ;;
    *)
        echo "Error: Unsupported OS: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Construct binary name
if [ "$PLATFORM" = "windows" ]; then
    ASSET_NAME="${BINARY_NAME}-${PLATFORM}-${ARCH}.exe"
else
    ASSET_NAME="${BINARY_NAME}-${PLATFORM}-${ARCH}"
fi

echo "========================================"
echo "  my-context Installer"
echo "========================================"
echo ""
echo "Platform:    $PLATFORM-$ARCH"
echo "Install to:  $INSTALL_DIR"
echo "Mode:        $INSTALL_MODE"
echo ""

# Get version to install
if [ "$VERSION" = "latest" ]; then
    echo "Fetching latest release..."
    VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest version"
        exit 1
    fi
fi
echo "Version:     $VERSION"
echo ""

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download binary and checksum
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET_NAME"
CHECKSUM_URL="https://github.com/$REPO/releases/download/$VERSION/${ASSET_NAME}.sha256"

echo "Downloading $ASSET_NAME..."
if ! curl -fsSL -o "$TMP_DIR/$ASSET_NAME" "$DOWNLOAD_URL"; then
    echo "Error: Failed to download binary from $DOWNLOAD_URL"
    exit 1
fi

echo "Downloading checksum..."
if ! curl -fsSL -o "$TMP_DIR/${ASSET_NAME}.sha256" "$CHECKSUM_URL"; then
    echo "Error: Failed to download checksum"
    exit 1
fi

# Verify checksum
echo "Verifying checksum..."
cd "$TMP_DIR"
if ! sha256sum -c "${ASSET_NAME}.sha256"; then
    echo "Error: Checksum verification failed!"
    exit 1
fi
echo "Checksum OK"
echo ""

# Check for existing installation
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    CURRENT_VERSION=$("$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null | head -1 || echo "unknown")
    echo "Existing installation: $CURRENT_VERSION"
fi

# Create install directory if needed
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating $INSTALL_DIR..."
    $SUDO_CMD mkdir -p "$INSTALL_DIR"
fi

# Install binary
echo "Installing to $INSTALL_DIR/$BINARY_NAME..."
$SUDO_CMD cp "$TMP_DIR/$ASSET_NAME" "$INSTALL_DIR/$BINARY_NAME"
$SUDO_CMD chmod +x "$INSTALL_DIR/$BINARY_NAME"

# For user mode, ensure directory ownership
if [ "$INSTALL_MODE" = "user" ]; then
    chown "$(id -u):$(id -g)" "$INSTALL_DIR/$BINARY_NAME" 2>/dev/null || true
fi

# Verify installation
echo ""
if "$INSTALL_DIR/$BINARY_NAME" --version; then
    echo ""
    echo "========================================"
    echo "  Installation successful!"
    echo "========================================"
else
    echo "Error: Installation verification failed"
    exit 1
fi

# Check PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "NOTE: $INSTALL_DIR is not in your PATH."
    echo ""
    if [ "$INSTALL_MODE" = "user" ]; then
        echo "Add to your shell profile:"
        echo "  echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
    else
        echo "/usr/local/bin should already be in PATH. Try opening a new terminal."
    fi
fi

echo ""
echo "Data directory: ~/.my-context/ (or set MY_CONTEXT_HOME)"
echo ""
