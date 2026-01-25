#!/usr/bin/env bash
# pre-release-check.sh - Pre-release validation checks
#
# Usage: ./scripts/pre-release-check.sh [VERSION]
#   VERSION: Optional version tag in format vX.Y.Z (e.g., v3.2.5)
#
# Validates:
# - Clean git state (no uncommitted changes)
# - Full test suite passes with race detector
# - Linter passes
# - All platform builds succeed
# - Version tag format (if provided)
# - CHANGELOG entry exists for version (if provided)
#
# Exit codes:
#   0 = All checks passed, ready for release
#   1 = One or more checks failed

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check counters
CHECKS_PASSED=0
CHECKS_FAILED=0
CHECKS_WARNING=0

# Optional version argument
VERSION_TAG="$1"

# Print functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Pre-Release Validation${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    if [ -n "$VERSION_TAG" ]; then
        echo -e "${BLUE}Target version:${NC} $VERSION_TAG"
        echo ""
    fi
}

print_check() {
    echo -e "${BLUE}[CHECK]${NC} $1"
}

print_pass() {
    echo -e "${GREEN}  ✓${NC} $1"
    ((CHECKS_PASSED++))
}

print_fail() {
    echo -e "${RED}  ✗${NC} $1"
    ((CHECKS_FAILED++))
}

print_warn() {
    echo -e "${YELLOW}  ⚠${NC} $1"
    ((CHECKS_WARNING++))
}

print_info() {
    echo -e "    $1"
}

# Check 1: Clean git state
check_git_state() {
    print_check "Git working directory status"
    
    if [ -n "$(git status --porcelain)" ]; then
        print_fail "Working directory not clean"
        print_info "Uncommitted changes detected:"
        git status --short | while IFS= read -r line; do
            print_info "  $line"
        done
        print_info "Commit or stash changes before release"
        return 1
    else
        print_pass "Git state clean"
    fi
    
    echo ""
}

# Check 2: Run full test suite with race detector
check_tests() {
    print_check "Running full test suite with race detector"
    
    if go test ./... -race > /tmp/pre-release-test.log 2>&1; then
        print_pass "All tests pass"
    else
        print_fail "Tests failed"
        print_info "See /tmp/pre-release-test.log for details"
        print_info "Last 20 lines:"
        tail -n 20 /tmp/pre-release-test.log | while IFS= read -r line; do
            print_info "  $line"
        done
        return 1
    fi
    
    echo ""
}

# Check 3: Run linter
check_lint() {
    print_check "Running golangci-lint"
    
    if ! command -v golangci-lint &> /dev/null; then
        print_fail "golangci-lint not installed"
        print_info "Install: https://golangci-lint.run/usage/install/"
        return 1
    fi
    
    if golangci-lint run --timeout=5m > /tmp/pre-release-lint.log 2>&1; then
        print_pass "Lint passes"
    else
        print_fail "Lint failed"
        print_info "See /tmp/pre-release-lint.log for details"
        print_info "Common issues:"
        head -n 10 /tmp/pre-release-lint.log | while IFS= read -r line; do
            print_info "  $line"
        done
        return 1
    fi
    
    echo ""
}

# Check 4: Build all platforms
check_builds() {
    print_check "Building all platforms"
    
    if [ ! -f "./scripts/build-all.sh" ]; then
        print_fail "build-all.sh script not found"
        return 1
    fi
    
    if ./scripts/build-all.sh > /tmp/pre-release-build.log 2>&1; then
        print_pass "All builds succeed"
        
        # Show built binaries
        if [ -d "bin" ]; then
            BINARY_COUNT=$(ls -1 bin/my-context-* 2>/dev/null | grep -v ".sha256" | wc -l)
            print_info "Built $BINARY_COUNT binaries:"
            ls -lh bin/my-context-* 2>/dev/null | grep -v ".sha256" | while IFS= read -r line; do
                SIZE=$(echo "$line" | awk '{print $5}')
                NAME=$(echo "$line" | awk '{print $9}' | xargs basename)
                print_info "  $NAME ($SIZE)"
            done
        fi
    else
        print_fail "Build failed"
        print_info "See /tmp/pre-release-build.log for details"
        print_info "Last 30 lines:"
        tail -n 30 /tmp/pre-release-build.log | while IFS= read -r line; do
            print_info "  $line"
        done
        return 1
    fi
    
    echo ""
}

# Check 5: Validate version tag format
check_version_format() {
    if [ -z "$VERSION_TAG" ]; then
        print_warn "Version tag not provided (optional)"
        print_info "Usage: $0 vX.Y.Z"
        echo ""
        return 0
    fi
    
    print_check "Version tag format"
    
    if [[ ! "$VERSION_TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_fail "Invalid version format: $VERSION_TAG"
        print_info "Expected format: vX.Y.Z (e.g., v3.2.5)"
        print_info "Current: $VERSION_TAG"
        return 1
    else
        print_pass "Version format valid: $VERSION_TAG"
    fi
    
    echo ""
}

# Check 6: CHANGELOG entry
check_changelog() {
    if [ -z "$VERSION_TAG" ]; then
        print_warn "Version tag not provided, skipping CHANGELOG check"
        echo ""
        return 0
    fi
    
    print_check "CHANGELOG entry for $VERSION_TAG"
    
    if [ ! -f "CHANGELOG.md" ]; then
        print_warn "CHANGELOG.md not found"
        print_info "Create CHANGELOG.md and document release changes"
        echo ""
        return 0
    fi
    
    # Extract version without 'v' prefix for CHANGELOG search
    VERSION_NUMBER="${VERSION_TAG#v}"
    
    # Check for version in CHANGELOG (allow [X.Y.Z] or ## [X.Y.Z] format)
    if grep -q "\[$VERSION_NUMBER\]" CHANGELOG.md || grep -q "## $VERSION_TAG" CHANGELOG.md; then
        print_pass "CHANGELOG entry exists"
        
        # Show the entry
        print_info "Entry preview:"
        if LINE_NUM=$(grep -n "\[$VERSION_NUMBER\]" CHANGELOG.md | head -1 | cut -d: -f1); then
            sed -n "${LINE_NUM},$((LINE_NUM + 5))p" CHANGELOG.md | while IFS= read -r line; do
                print_info "  $line"
            done
        fi
    else
        print_warn "No CHANGELOG entry for $VERSION_TAG"
        print_info "Add release notes to CHANGELOG.md before releasing"
        print_info "Expected format:"
        print_info "  ## [$VERSION_NUMBER] - $(date +%Y-%m-%d)"
        print_info "  "
        print_info "  ### Added/Changed/Fixed"
        print_info "  - Your changes here"
    fi
    
    echo ""
}

# Summary
print_summary() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Summary${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "${GREEN}Passed:${NC}   $CHECKS_PASSED"
    echo -e "${YELLOW}Warnings:${NC} $CHECKS_WARNING"
    echo -e "${RED}Failed:${NC}   $CHECKS_FAILED"
    echo ""
    
    if [[ $CHECKS_FAILED -eq 0 ]]; then
        echo -e "${GREEN}✓ Ready to tag and release!${NC}"
        if [ -n "$VERSION_TAG" ]; then
            echo ""
            echo -e "${YELLOW}Next steps:${NC}"
            echo "  git tag $VERSION_TAG"
            echo "  git push origin $VERSION_TAG"
            echo "  # Then create GitHub release with binaries from bin/"
        fi
        return 0
    else
        echo -e "${RED}✗ Fix failed checks before releasing${NC}"
        return 1
    fi
}

# Main execution
main() {
    print_header
    
    # Run all checks (don't exit on first failure)
    set +e
    
    check_git_state
    check_tests
    check_lint
    check_builds
    check_version_format
    check_changelog
    
    set -e
    
    print_summary
}

# Run
main
exit $?
