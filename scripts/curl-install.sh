#!/usr/bin/env bash
# One-liner curl installer for my-context
# Detects platform and installs appropriate binary
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/USER/REPO/main/scripts/curl-install.sh | bash
#   Or with specific version:
#   curl -sSL https://raw.githubusercontent.com/USER/REPO/main/scripts/curl-install.sh | bash -s -- v2.0.0

set -e

REPO_OWNER="jefferycaldwell"
REPO_NAME="my-context-copilot"
VERSION="${1:-latest}"  # Use latest if not specified

echo "my-context installer"
echo "===================="
echo ""

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
	Linux*)
		PLATFORM="linux-amd64"
		BINARY_NAME="my-context-linux-amd64"
		;;
	Darwin*)
		if [ "$ARCH" == "arm64" ]; then
			PLATFORM="darwin-arm64"
			BINARY_NAME="my-context-darwin-arm64"
		else
			PLATFORM="darwin-amd64"
			BINARY_NAME="my-context-darwin-amd64"
		fi
		;;
	MINGW*|MSYS*|CYGWIN*)
		PLATFORM="windows-amd64"
		BINARY_NAME="my-context-windows-amd64.exe"
		;;
	*)
		echo "Error: Unsupported platform: $OS ($ARCH)"
		echo ""
		echo "Supported platforms:"
		echo "  - Linux (x86_64)"
		echo "  - macOS (x86_64, arm64)"
		echo "  - Windows (x86_64, via Git Bash/WSL)"
		echo ""
		echo "For manual installation, visit:"
		echo "  https://github.com/$REPO_OWNER/$REPO_NAME/releases"
		exit 1
		;;
esac

echo "Detected platform: $PLATFORM"

# Determine download URL
if [ "$VERSION" == "latest" ]; then
	DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/latest/download/$BINARY_NAME"
	CHECKSUM_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/latest/download/$BINARY_NAME.sha256"
else
	DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/$BINARY_NAME"
	CHECKSUM_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/$BINARY_NAME.sha256"
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

cd "$TMP_DIR"

echo "Downloading $BINARY_NAME..."
if command -v curl >/dev/null 2>&1; then
	curl -sSL "$DOWNLOAD_URL" -o "$BINARY_NAME"
	curl -sSL "$CHECKSUM_URL" -o "$BINARY_NAME.sha256"
elif command -v wget >/dev/null 2>&1; then
	wget -q "$DOWNLOAD_URL" -O "$BINARY_NAME"
	wget -q "$CHECKSUM_URL" -O "$BINARY_NAME.sha256"
else
	echo "Error: curl or wget required"
	exit 1
fi

# Verify checksum
echo "Verifying checksum..."
if command -v sha256sum >/dev/null 2>&1; then
	sha256sum -c "$BINARY_NAME.sha256" || {
		echo "Error: Checksum verification failed"
		echo "Downloaded binary may be corrupted or tampered with"
		exit 1
	}
elif command -v shasum >/dev/null 2>&1; then
	shasum -a 256 -c "$BINARY_NAME.sha256" || {
		echo "Error: Checksum verification failed"
		exit 1
	}
else
	echo "Warning: sha256sum/shasum not found, skipping checksum verification"
fi

# Make executable
chmod +x "$BINARY_NAME"

# Detect if we need to use install.bat (pure Windows, not WSL/Git Bash)
if [[ "$OS" == MINGW* || "$OS" == MSYS* ]] && [[ ! -f /proc/version ]]; then
	echo ""
	echo "Windows detected. Please run the following command manually:"
	echo "  scripts\\install.bat $TMP_DIR\\$BINARY_NAME"
	echo ""
	echo "Or copy the binary to your PATH manually:"
	echo "  copy $BINARY_NAME %USERPROFILE%\\bin\\my-context.exe"
	exit 0
fi

# For Unix-like systems (Linux, macOS, WSL), delegate to install.sh
echo ""
echo "Installing to ~/.local/bin..."

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"
cp "$BINARY_NAME" "$INSTALL_DIR/my-context"
chmod +x "$INSTALL_DIR/my-context"

# Add to PATH if needed
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
	# Detect shell
	if [ -n "$BASH_VERSION" ]; then
		RC_FILE="$HOME/.bashrc"
	elif [ -n "$ZSH_VERSION" ]; then
		RC_FILE="$HOME/.zshrc"
	else
		RC_FILE="$HOME/.profile"
	fi

	echo "" >> "$RC_FILE"
	echo "# Added by my-context installer" >> "$RC_FILE"
	echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$RC_FILE"
	echo "Added $INSTALL_DIR to PATH in $RC_FILE"
fi

# Verify installation
if "$INSTALL_DIR/my-context" --version; then
	echo ""
	echo "✓ Installation complete!"
	echo ""
	echo "Run: source ~/.bashrc (or restart your terminal)"
	echo "Then: my-context --help"
else
	echo ""
	echo "✗ Installation verification failed"
	exit 1
fi

