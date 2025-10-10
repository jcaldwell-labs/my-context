#!/usr/bin/env bash
#
# github-preflight.sh - Pre-flight checks before GitHub operations
#
# Usage: ./github-preflight.sh [options]
#
# Validates:
# - SSH authentication to GitHub
# - Git remotes configuration
# - Existing tags
# - Worktree health
# - Uncommitted changes
#
# Exit codes:
#   0 = All checks passed, ready for GitHub operations
#   1 = One or more checks failed
#

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

# Print functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}GitHub Pre-Flight Check${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
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

# Check 1: SSH keys exist
check_ssh_keys() {
    print_check "SSH keys in ~/.ssh/"

    if [[ -f ~/.ssh/id_rsa ]]; then
        print_pass "Private key exists (id_rsa)"

        # Check permissions
        PERMS=$(stat -c %a ~/.ssh/id_rsa 2>/dev/null || stat -f %A ~/.ssh/id_rsa 2>/dev/null)
        if [[ "$PERMS" == "600" ]]; then
            print_pass "Private key permissions correct (600)"
        else
            print_fail "Private key permissions incorrect ($PERMS, should be 600)"
            print_info "Fix: chmod 600 ~/.ssh/id_rsa"
        fi
    else
        print_fail "Private key missing (id_rsa)"
        print_info "Copy from Windows: cp /mnt/c/Users/\$USER/.ssh/id_rsa ~/.ssh/"
    fi

    if [[ -f ~/.ssh/id_rsa.pub ]]; then
        print_pass "Public key exists (id_rsa.pub)"
    else
        print_warn "Public key missing (id_rsa.pub)"
    fi

    if [[ -f ~/.ssh/known_hosts ]]; then
        print_pass "known_hosts exists"
    else
        print_warn "known_hosts missing (will be created on first connection)"
    fi

    echo ""
}

# Check 2: GitHub SSH authentication
check_github_auth() {
    print_check "GitHub SSH authentication"

    if ssh -T git@github.com 2>&1 | grep -q "successfully authenticated"; then
        print_pass "GitHub SSH authentication successful"
        USERNAME=$(ssh -T git@github.com 2>&1 | grep "Hi" | awk '{print $2}' | tr -d '!')
        print_info "Authenticated as: $USERNAME"
    else
        print_fail "GitHub SSH authentication failed"
        print_info "Test manually: ssh -T git@github.com"
        print_info "Ensure SSH keys are added to GitHub account"
    fi

    echo ""
}

# Check 3: Git remotes
check_git_remotes() {
    print_check "Git remote configuration"

    if git remote -v &>/dev/null; then
        REMOTES=$(git remote)

        if [[ -z "$REMOTES" ]]; then
            print_warn "No remotes configured"
            print_info "Add GitHub remote: git remote add origin git@github.com:USER/REPO.git"
        else
            print_pass "Remotes configured:"
            git remote -v | while IFS= read -r line; do
                print_info "$line"
            done

            # Check for GitHub origin
            if git remote -v | grep -q "origin.*github.com"; then
                print_pass "GitHub origin remote exists"
            else
                print_warn "No GitHub origin remote found"
                print_info "Add: git remote add origin git@github.com:USER/REPO.git"
            fi
        fi
    else
        print_fail "Not a git repository"
        print_info "Run from within a git repository"
    fi

    echo ""
}

# Check 4: Existing tags
check_tags() {
    print_check "Git tags"

    if git tag &>/dev/null; then
        TAGS=$(git tag)

        if [[ -z "$TAGS" ]]; then
            print_pass "No tags (clean slate)"
        else
            TAG_COUNT=$(git tag | wc -l)
            print_warn "$TAG_COUNT tag(s) exist:"
            git tag | while IFS= read -r tag; do
                COMMIT=$(git rev-parse --short "$tag^{commit}" 2>/dev/null || echo "unknown")
                print_info "$tag → $COMMIT"
            done
            print_info "Verify no conflicts before creating new tags"
        fi
    else
        print_fail "Cannot check tags (not a git repository)"
    fi

    echo ""
}

# Check 5: Worktree health (if deb-sanity available)
check_worktree_health() {
    print_check "Worktree health"

    if command -v deb-sanity &>/dev/null; then
        if OUTPUT=$(deb-sanity --health . 2>&1 | grep -E "(HEALTHY|BROKEN|WARNING)"); then
            if echo "$OUTPUT" | grep -q "HEALTHY"; then
                print_pass "Worktree is healthy"
            elif echo "$OUTPUT" | grep -q "BROKEN"; then
                print_fail "Worktree is broken"
                print_info "Run: deb-sanity --health . --fix"
            else
                print_warn "Worktree has warnings"
                print_info "Run: deb-sanity --health . for details"
            fi
        else
            print_warn "Could not determine worktree health"
        fi
    else
        print_warn "deb-sanity not available (skipping worktree health check)"
        print_info "Install deb-sanity for comprehensive checks"
    fi

    echo ""
}

# Check 6: Uncommitted changes
check_uncommitted_changes() {
    print_check "Working tree status"

    if git status &>/dev/null; then
        if git diff --quiet && git diff --cached --quiet; then
            print_pass "No uncommitted changes"
        else
            MODIFIED=$(git diff --name-only | wc -l)
            STAGED=$(git diff --cached --name-only | wc -l)

            if [[ $MODIFIED -gt 0 ]]; then
                print_warn "$MODIFIED file(s) with uncommitted changes"
            fi

            if [[ $STAGED -gt 0 ]]; then
                print_warn "$STAGED file(s) staged but not committed"
            fi

            print_info "Commit or stash changes before pushing to GitHub"
        fi

        # Check for untracked files
        UNTRACKED=$(git ls-files --others --exclude-standard | wc -l)
        if [[ $UNTRACKED -gt 0 ]]; then
            print_warn "$UNTRACKED untracked file(s)"
            print_info "Review untracked files: git status"
        fi
    else
        print_fail "Cannot check status (not a git repository)"
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
        echo -e "${GREEN}✓ Ready for GitHub operations!${NC}"
        return 0
    else
        echo -e "${RED}✗ Fix failed checks before proceeding${NC}"
        return 1
    fi
}

# Main execution
main() {
    print_header
    check_ssh_keys
    check_github_auth
    check_git_remotes
    check_tags
    check_worktree_health
    check_uncommitted_changes
    print_summary
}

# Run
main
exit $?
