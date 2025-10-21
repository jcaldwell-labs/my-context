#!/usr/bin/env bash
# build.sh - Quick build for current platform with proper metadata
# Usage: ./scripts/build.sh [--install]

set -e

# Get version from git (tag or branch-based)
if git describe --tags --exact-match 2>/dev/null; then
    VERSION=$(git describe --tags --exact-match)
else
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
fi

BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD)

echo "🔨 Building my-context..."
echo "   Version: $VERSION"
echo "   Build: $BUILD_TIME"
echo "   Commit: $GIT_COMMIT"
echo ""

# Build with ldflags for metadata
go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
    -o my-context cmd/my-context/main.go

echo "✅ Build complete: my-context"
echo ""
./my-context --version
echo ""

# Install if --install flag provided
if [[ "$1" == "--install" ]]; then
    echo "📦 Installing to ~/.local/bin/my-context..."
    cp my-context ~/.local/bin/my-context
    echo "✅ Installed: $(which my-context)"
    ~/.local/bin/my-context --version
fi

echo ""
echo "To install: cp my-context ~/.local/bin/my-context"
echo "Or run: ./scripts/build.sh --install"
