#!/usr/bin/env bash
# build-all.sh - Build my-context binaries for all supported platforms
# Usage: ./scripts/build-all.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
OUTPUT_DIR="bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Platform configurations: GOOS GOARCH OUTPUT_NAME
PLATFORMS=(
  "linux amd64 my-context-linux-amd64"
  "windows amd64 my-context-windows-amd64.exe"
  "darwin amd64 my-context-darwin-amd64"
  "darwin arm64 my-context-darwin-arm64"
)

echo -e "${GREEN}Building my-context for multiple platforms...${NC}"
echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"
echo "Git Commit: $GIT_COMMIT"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build for each platform
for platform_config in "${PLATFORMS[@]}"; do
  read -r GOOS GOARCH OUTPUT_NAME <<< "$platform_config"

  echo -e "${YELLOW}Building ${GOOS}/${GOARCH}...${NC}"

  CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build \
    -ldflags "\
      -X main.Version=$VERSION \
      -X main.BuildTime=$BUILD_TIME \
      -X main.GitCommit=$GIT_COMMIT" \
    -o "$OUTPUT_DIR/$OUTPUT_NAME" \
    ./cmd/my-context/

  if [ $? -eq 0 ]; then
    # Generate checksum
    if command -v sha256sum &> /dev/null; then
      (cd "$OUTPUT_DIR" && sha256sum "$OUTPUT_NAME" > "$OUTPUT_NAME.sha256")
    elif command -v shasum &> /dev/null; then
      (cd "$OUTPUT_DIR" && shasum -a 256 "$OUTPUT_NAME" > "$OUTPUT_NAME.sha256")
    else
      echo -e "${RED}Warning: No checksum tool available (sha256sum or shasum)${NC}"
    fi

    # Display binary info
    SIZE=$(ls -lh "$OUTPUT_DIR/$OUTPUT_NAME" | awk '{print $5}')
    echo -e "${GREEN}✓ Built: $OUTPUT_NAME ($SIZE)${NC}"
  else
    echo -e "${RED}✗ Failed to build ${GOOS}/${GOARCH}${NC}"
    exit 1
  fi
  echo ""
done

echo -e "${GREEN}All builds completed successfully!${NC}"
echo "Output directory: $OUTPUT_DIR/"
echo ""
echo "Binaries:"
ls -lh "$OUTPUT_DIR"/*.{exe,,} 2>/dev/null | grep -v ".sha256"

echo ""
echo "Checksums:"
cat "$OUTPUT_DIR"/*.sha256

echo ""
echo -e "${YELLOW}To test a binary:${NC}"
echo "  ./$OUTPUT_DIR/my-context-linux-amd64 --version"
echo "  ./$OUTPUT_DIR/my-context-windows-amd64.exe --version"
echo ""
echo -e "${YELLOW}To install locally:${NC}"
echo "  ./scripts/install.sh"
name: Release

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

permissions:
  contents: write

jobs:
  build:
    name: Build Multi-Platform Binaries
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
            output: my-context-linux-amd64
          - goos: windows
            goarch: amd64
            output: my-context-windows-amd64.exe
          - goos: darwin
            goarch: amd64
            output: my-context-darwin-amd64
          - goos: darwin
            goarch: arm64
            output: my-context-darwin-arm64

    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Get version info
        id: version
        run: |
          VERSION=${GITHUB_REF#refs/tags/v}
          if [ -z "$VERSION" ]; then
            VERSION=$(git describe --tags --always --dirty)
          fi
          BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
          GIT_COMMIT=$(git rev-parse --short HEAD)
          echo "version=$VERSION" >> $GITHUB_OUTPUT
          echo "build_time=$BUILD_TIME" >> $GITHUB_OUTPUT
          echo "git_commit=$GIT_COMMIT" >> $GITHUB_OUTPUT

      - name: Build binary for ${{ matrix.goos }}/${{ matrix.goarch }}
        env:
          CGO_ENABLED: 0
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          go build -ldflags "\
            -X main.Version=${{ steps.version.outputs.version }} \
            -X main.BuildTime=${{ steps.version.outputs.build_time }} \
            -X main.GitCommit=${{ steps.version.outputs.git_commit }}" \
            -o ${{ matrix.output }} ./cmd/my-context/

      - name: Generate SHA256 checksum
        run: |
          sha256sum ${{ matrix.output }} > ${{ matrix.output }}.sha256

      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: binaries
          path: |
            ${{ matrix.output }}
            ${{ matrix.output }}.sha256

  release:
    name: Create Release
    needs: build
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Download artifacts
        uses: actions/download-artifact@v3
        with:
          name: binaries
          path: ./binaries

      - name: List artifacts
        run: ls -lh ./binaries/

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: ./binaries/*
          draft: false
          prerelease: false
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

