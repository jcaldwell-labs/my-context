#!/usr/bin/env bash
#
# prepare-public-release.sh - Automate public release preparation
#
# Usage: ./prepare-public-release.sh [branch-name]
#
# This script automates the cleanup and preparation of a public release branch:
# - Removes internal development artifacts
# - Adds public-facing documentation (LICENSE, CONTRIBUTING.md)
# - Updates README for public release
# - Enhances .gitignore
# - Creates a comprehensive commit
#
# Default branch name: public-v1.0.0
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
BRANCH_NAME="${1:-public-v1.0.0}"
DRY_RUN=false

# Print functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Prepare Public Release${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

print_step() {
    echo -e "${BLUE}▸${NC} $1"
}

print_success() {
    echo -e "${GREEN}  ✓${NC} $1"
}

print_info() {
    echo -e "    $1"
}

print_warning() {
    echo -e "${YELLOW}  ⚠${NC} $1"
}

print_error() {
    echo -e "${RED}  ✗${NC} $1"
}

# Check if in git repository
check_git_repo() {
    if ! git rev-parse --is-inside-work-tree &>/dev/null; then
        print_error "Not a git repository"
        exit 1
    fi
}

# Remove internal files
remove_internal_files() {
    print_step "Removing internal development artifacts..."

    local files_removed=0

    # IDE configs
    for dir in .claude .cursor .idea .specify; do
        if [[ -d "$dir" ]]; then
            git rm -r "$dir" &>/dev/null && ((files_removed++))
            print_success "Removed $dir/"
        fi
    done

    # Individual IDE files
    for file in .cursorrules; do
        if [[ -f "$file" ]]; then
            git rm "$file" &>/dev/null && ((files_removed++))
            print_success "Removed $file"
        fi
    done

    # Specs directory
    if [[ -d "specs" ]]; then
        git rm -r specs &>/dev/null && ((files_removed++))
        print_success "Removed specs/"
    fi

    # Internal documentation
    local internal_docs=(
        "CLAUDE.md"
        "CONSTITUTION-*.md"
        "HERE.md"
        "IMPLEMENTATION.md"
        "REMAINING-TECH-DEBT.md"
        "SDLC.md"
        "SETUP.md"
        "SPRINT-*.md"
        "TECH-DEBT*.md"
        "WORKTREE-*.md"
        "demo_*.md"
        "MY-CONTEXT-*.md"
        "TOOLS1-*.md"
        "DEPLOY-*.md"
    )

    for pattern in "${internal_docs[@]}"; do
        for file in $pattern; do
            if [[ -f "$file" ]]; then
                git rm "$file" &>/dev/null && ((files_removed++))
                print_success "Removed $file"
            fi
        done
    done

    # Binaries
    for file in my-context my-context.exe; do
        if [[ -f "$file" ]]; then
            git rm "$file" &>/dev/null && ((files_removed++))
            print_success "Removed $file"
        fi
    done

    # Thread documentation
    if [[ -d "docs" ]]; then
        for file in docs/THREAD-*.md; do
            if [[ -f "$file" ]]; then
                git rm "$file" &>/dev/null && ((files_removed++))
                print_success "Removed $file"
            fi
        done
    fi

    print_info "Total files removed: $files_removed"
    echo ""
}

# Create LICENSE file
create_license() {
    print_step "Creating LICENSE (MIT)..."

    if [[ -f "LICENSE" ]]; then
        print_warning "LICENSE already exists, skipping"
        echo ""
        return
    fi

    cat > LICENSE << 'EOF'
MIT License

Copyright (c) 2025 My-Context Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

    git add LICENSE
    print_success "LICENSE created (21 lines)"
    echo ""
}

# Update .gitignore
update_gitignore() {
    print_step "Updating .gitignore for public release..."

    cat > .gitignore << 'EOF'
# Binaries
/my-context
/my-context.exe
*.exe
/bin/

# IDE and editor configs
.claude/
.cursor/
.cursorrules
.idea/
.vscode/
*.swp
*.swo
*~

# Internal development
.specify/
specs/
CLAUDE.md
HERE.md
IMPLEMENTATION.md
SDLC.md
SETUP.md
SPRINT-*.md
TECH-DEBT*.md
WORKTREE-*.md
CONSTITUTION-*.md
REMAINING-TECH-DEBT.md
demo_*.md
MY-CONTEXT-*.md
TOOLS1-*.md
DEPLOY-*.md

# OS files
.DS_Store
Thumbs.db

# Test coverage
coverage.out
*.test
EOF

    git add .gitignore
    print_success ".gitignore updated (42 lines)"
    echo ""
}

# Update README
update_readme() {
    print_step "Updating README for public release..."

    if [[ ! -f "README.md" ]]; then
        print_warning "README.md not found, skipping"
        echo ""
        return
    fi

    # Update title (if contains "Copilot")
    if grep -q "My-Context-Copilot" README.md; then
        sed -i 's/My-Context-Copilot/My-Context/g' README.md
        print_success "Updated title: My-Context-Copilot → My-Context"
    fi

    # Update GitHub URLs
    if grep -q "USER/REPO" README.md; then
        sed -i 's|USER/REPO|YOUR-USERNAME/my-context|g' README.md
        print_success "Updated GitHub URLs (placeholder: YOUR-USERNAME)"
    fi

    # Comment out curl install if present
    if grep -q "curl -sSL.*curl-install.sh" README.md; then
        sed -i 's|^\(curl -sSL.*\)|# Coming soon - check releases page for now\n# \1|' README.md
        print_success "Commented out curl install (not yet available)"
    fi

    git add README.md
    print_success "README.md updated"
    echo ""
}

# Create commit
create_commit() {
    print_step "Creating commit..."

    # Count changes
    local files_changed=$(git diff --cached --numstat | wc -l)
    local insertions=$(git diff --cached --numstat | awk '{sum+=$1} END {print sum}')
    local deletions=$(git diff --cached --numstat | awk '{sum+=$2} END {print sum}')

    if [[ $files_changed -eq 0 ]]; then
        print_warning "No changes to commit"
        echo ""
        return
    fi

    git commit -m "chore: prepare for public release

- Remove internal development artifacts (IDE configs, specs, internal docs)
- Add LICENSE (MIT)
- Update README.md for public repository
- Enhance .gitignore for public release
- Remove binaries from git tracking

Public-facing repository ready for GitHub.

Changes: $files_changed files, +$insertions -$deletions lines"

    print_success "Commit created"
    print_info "Files changed: $files_changed"
    print_info "Insertions: +$insertions"
    print_info "Deletions: -$deletions"
    echo ""
}

# Summary
print_summary() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}Summary${NC}"
    echo -e "${BLUE}========================================${NC}"

    local commit_hash=$(git rev-parse --short HEAD)
    echo -e "${GREEN}✓ Public release branch prepared!${NC}"
    echo ""
    echo "Branch: $BRANCH_NAME"
    echo "Commit: $commit_hash"
    echo ""
    echo "Next steps:"
    echo "1. Review changes: git show"
    echo "2. Run preflight: ./scripts/github-preflight.sh"
    echo "3. Push to GitHub: git push -u origin $BRANCH_NAME:main"
    echo "4. Create tag: git tag -a v1.0.0 -m \"Release v1.0.0\""
    echo "5. Push tag: git push origin v1.0.0"
    echo ""
}

# Main execution
main() {
    print_header

    # Validate environment
    check_git_repo

    # Perform cleanup and preparation
    remove_internal_files
    create_license
    update_gitignore
    update_readme
    create_commit

    # Summary
    print_summary
}

# Run
main
exit 0
