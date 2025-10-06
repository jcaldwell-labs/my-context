#!/usr/bin/env bash
# Installation script for my-context (Linux, macOS, WSL)
# Installs to ~/.local/bin without requiring sudo

set -e

INSTALL_DIR="$HOME/.local/bin"
BINARY_NAME="my-context"

echo "Installing my-context..."

# Detect existing installation
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
	echo "Existing installation found."
	CURRENT_VERSION=$("$INSTALL_DIR/$BINARY_NAME" --version 2>/dev/null || echo "unknown")
	echo "Current version: $CURRENT_VERSION"

	# Backup old binary
	mv "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME.backup"
	echo "Backed up to $BINARY_NAME.backup"
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

# Copy binary (assume binary provided as $1 or detect from bin/)
if [ -n "$1" ]; then
	# Binary path provided as argument
	cp "$1" "$INSTALL_DIR/$BINARY_NAME"
else
	# Auto-detect platform and use appropriate binary from bin/
	OS=$(uname -s)
	ARCH=$(uname -m)

	case "$OS" in
		Linux*)
			BINARY_SRC="bin/my-context-linux-amd64"
			;;
		Darwin*)
			if [ "$ARCH" == "arm64" ]; then
				BINARY_SRC="bin/my-context-darwin-arm64"
			else
				BINARY_SRC="bin/my-context-darwin-amd64"
			fi
			;;
		MINGW*|MSYS*|CYGWIN*)
			BINARY_SRC="bin/my-context-windows-amd64.exe"
			;;
		*)
			echo "Error: Unsupported platform: $OS"
			exit 1
			;;
	esac

	if [ -f "$BINARY_SRC" ]; then
		cp "$BINARY_SRC" "$INSTALL_DIR/$BINARY_NAME"
	else
		echo "Error: Binary not found at $BINARY_SRC"
		echo "Run './scripts/build-all.sh' first to build binaries"
		exit 1
	fi
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Add to PATH if not already present
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
	# Detect shell and add to appropriate RC file
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
	echo "Run: source $RC_FILE (or restart your terminal)"
fi

# Verify installation
if "$INSTALL_DIR/$BINARY_NAME" --version; then
	echo ""
	echo "✓ Installation complete!"

	# Remove backup if exists
	if [ -f "$INSTALL_DIR/$BINARY_NAME.backup" ]; then
		rm "$INSTALL_DIR/$BINARY_NAME.backup"
	fi
else
	echo ""
	echo "✗ Installation verification failed"
	exit 1
fi

echo ""
echo "Note: ~/.my-context/ data directory is preserved (separate from binary)"
