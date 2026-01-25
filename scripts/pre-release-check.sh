#!/usr/bin/env bash
# pre-release-check.sh - Validate release readiness before tagging
# Usage: ./scripts/pre-release-check.sh [version]
#   version: Optional version tag to validate (e.g., v3.2.5)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== Pre-Release Validation ==="
echo ""

# 1. Check clean git state
echo "Checking git state..."
if [ -n "$(git status --porcelain)" ]; then
  echo -e "${RED}❌ Working directory not clean${NC}"
  git status --short
  exit 1
fi
echo -e "${GREEN}✓ Git state clean${NC}"
echo ""

# 2. Run full test suite
echo "Running tests..."
if go test ./... -race -coverprofile=coverage.out; then
  echo -e "${GREEN}✓ All tests pass${NC}"
else
  echo -e "${RED}❌ Tests failed${NC}"
  exit 1
fi
echo ""

# 3. Run linter
echo "Running linter..."
if command -v golangci-lint > /dev/null; then
  if golangci-lint run --timeout=5m; then
    echo -e "${GREEN}✓ Lint passes${NC}"
  else
    echo -e "${RED}❌ Linting failed${NC}"
    exit 1
  fi
else
  echo -e "${YELLOW}⚠️  Warning: golangci-lint not installed${NC}"
fi
echo ""

# 4. Build all platforms
echo "Building all platforms..."
if ./scripts/build-all.sh; then
  echo -e "${GREEN}✓ All builds succeed${NC}"
else
  echo -e "${RED}❌ Build failed${NC}"
  exit 1
fi
echo ""

# 5. Check version tag format
if [ -n "$1" ]; then
  echo "Validating version format: $1"
  if [[ ! "$1" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ Invalid version format: $1 (expected vX.Y.Z)${NC}"
    exit 1
  fi
  echo -e "${GREEN}✓ Version format valid: $1${NC}"
  echo ""
fi

# 6. Check CHANGELOG has entry for version
if [ -n "$1" ]; then
  echo "Checking CHANGELOG entry..."
  # Extract version number without 'v' prefix (e.g., v3.2.5 -> 3.2.5)
  VERSION_NUM="${1#v}"
  if grep -q "\[$VERSION_NUM\]" CHANGELOG.md || grep -q "$1" CHANGELOG.md; then
    echo -e "${GREEN}✓ CHANGELOG entry exists${NC}"
  else
    echo -e "${YELLOW}⚠️  Warning: No CHANGELOG entry for $1 (checked [$VERSION_NUM] and $1)${NC}"
  fi
  echo ""
fi

echo ""
echo -e "${GREEN}=== Pre-Release Validation Complete ===${NC}"
echo "Ready to tag and release!"
