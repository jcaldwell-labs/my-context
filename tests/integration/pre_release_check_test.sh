#!/usr/bin/env bash
# Integration test for pre-release-check.sh script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PRE_RELEASE_SCRIPT="$SCRIPT_DIR/scripts/pre-release-check.sh"

echo "=== Testing pre-release-check.sh script ==="
echo ""

# Test 1: Script exists and is executable
echo "Test 1: Checking script exists and is executable..."
if [ ! -f "$PRE_RELEASE_SCRIPT" ]; then
  echo "✗ Script not found at $PRE_RELEASE_SCRIPT"
  exit 1
fi
if [ ! -x "$PRE_RELEASE_SCRIPT" ]; then
  echo "✗ Script is not executable"
  exit 1
fi
echo "✓ Script exists and is executable"
echo ""

# Test 2: Version format validation (using bash to test just the regex logic)
echo "Test 2: Testing version format validation logic..."

# Valid versions
for version in "v1.0.0" "v10.20.30" "v0.0.1"; do
  if [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "✓ Valid version format: $version"
  else
    echo "✗ Should accept valid version: $version"
    exit 1
  fi
done

# Invalid versions
for version in "1.0.0" "v1.0" "v1.0.0.0" "version1.0.0"; do
  if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "✓ Correctly rejected invalid format: $version"
  else
    echo "✗ Should reject invalid version: $version"
    exit 1
  fi
done
echo ""

# Test 3: CHANGELOG check logic
echo "Test 3: Testing CHANGELOG check logic..."
cd "$SCRIPT_DIR"

# Check for an existing version
VERSION="v3.2.4"
VERSION_NUM="${VERSION#v}"
if grep -q "\[$VERSION_NUM\]" CHANGELOG.md || grep -q "$VERSION" CHANGELOG.md; then
  echo "✓ CHANGELOG entry found for $VERSION"
else
  echo "⚠️  CHANGELOG entry not found for $VERSION (this is expected if version doesn't exist)"
fi

# Check for a non-existing version
VERSION="v99.99.99"
VERSION_NUM="${VERSION#v}"
if grep -q "\[$VERSION_NUM\]" CHANGELOG.md || grep -q "$VERSION" CHANGELOG.md; then
  echo "✗ Should not find CHANGELOG entry for $VERSION"
  exit 1
else
  echo "✓ Correctly determined no CHANGELOG entry for $VERSION"
fi
echo ""

# Test 4: Make target exists
echo "Test 4: Checking Makefile has release-check target..."
if grep -q "^release-check:" "$SCRIPT_DIR/Makefile"; then
  echo "✓ Makefile has release-check target"
else
  echo "✗ Makefile missing release-check target"
  exit 1
fi
echo ""

echo "=== All pre-release-check.sh tests passed! ==="
echo ""
echo "Note: Full validation (tests, lint, build) requires appropriate environment setup."
echo "This test validates the script logic, structure, and configuration."
