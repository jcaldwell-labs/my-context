#!/bin/bash
# Installation script for my-context

set -e

echo "Installing my-context..."

# Determine OS
OS=$(uname -s)
ARCH=$(uname -m)

# Select appropriate binary
BINARY=""
case "$OS" in
    Linux*)
        BINARY="dist/my-context-linux"
        ;;
    Darwin*)
        if [ "$ARCH" = "arm64" ]; then
            BINARY="dist/my-context-macos-arm"
        else
            BINARY="dist/my-context-macos"
        fi
        ;;
    MINGW*|MSYS*|CYGWIN*)
        BINARY="dist/my-context.exe"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: Binary not found: $BINARY"
    echo "Run './scripts/build.sh' first to build binaries"
    exit 1
fi

# Install to /usr/local/bin (requires sudo on Unix)
INSTALL_DIR="/usr/local/bin"

if [ "$OS" = "Linux" ] || [ "$OS" = "Darwin" ]; then
    echo "Installing to $INSTALL_DIR (requires sudo)..."
    sudo cp "$BINARY" "$INSTALL_DIR/my-context"
    sudo chmod +x "$INSTALL_DIR/my-context"
    echo "Installed successfully!"
    echo ""
    echo "Test it with: my-context --version"
else
    # Windows - copy to user's bin or suggest manual installation
    echo "On Windows, please manually copy $BINARY to a directory in your PATH"
    echo "For example: cp $BINARY /c/Windows/System32/my-context.exe"
fi

echo ""
echo "Installation complete!"
#!/bin/bash
# Build script for cross-platform compilation

set -e

echo "Building my-context for multiple platforms..."

# Create dist directory
mkdir -p dist

# Build for Windows
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/my-context.exe ./cmd/my-context/

# Build for Linux
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/my-context-linux ./cmd/my-context/

# Build for macOS Intel
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/my-context-macos ./cmd/my-context/

# Build for macOS ARM (M1/M2)
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/my-context-macos-arm ./cmd/my-context/

echo ""
echo "Build complete! Binaries in dist/ directory:"
ls -lh dist/

echo ""
echo "Binary sizes:"
du -h dist/*

