#!/usr/bin/env bash
# build.sh - Quick build for current platform
# For multi-platform builds, use: ./scripts/build-all.sh

set -e

echo "Building my-context for current platform..."

go build -o my-context.exe ./cmd/my-context/

echo " Build complete: my-context.exe"
echo ""
echo "Run: ./my-context.exe --version"
echo ""
echo "For multi-platform builds: ./scripts/build-all.sh"
