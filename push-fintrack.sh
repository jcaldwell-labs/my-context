#!/bin/bash
# Push FinTrack to GitHub
# Run this script from your local machine (outside Claude Code)

set -e

REPO_URL="https://github.com/jcaldwell1066/fintrack.git"
SOURCE_DIR="/tmp/fintrack-new"

echo "========================================="
echo "Pushing FinTrack to GitHub"
echo "========================================="
echo ""

# Check if source directory exists
if [ ! -d "$SOURCE_DIR" ]; then
    echo "Error: Source directory not found: $SOURCE_DIR"
    echo ""
    echo "Extracting from tarball..."
    mkdir -p "$SOURCE_DIR"
    cd "$SOURCE_DIR"
    tar -xzf ~/my-context/fintrack-repository.tar.gz --strip-components=1
fi

cd "$SOURCE_DIR"

# Check if git is initialized
if [ ! -d .git ]; then
    echo "Initializing git repository..."
    git init
    git branch -M main
fi

# Check if commit exists
if ! git log --oneline > /dev/null 2>&1; then
    echo "Creating initial commit..."
    git add .
    git commit -m "Initial commit: FinTrack Phase 1 foundation

Terminal-based personal finance tracking and budgeting CLI

Features:
- Account management (CRUD operations)
- PostgreSQL backend with GORM
- Configuration management (YAML/ENV)
- Table and JSON output formats
- Test suite with unit tests
- Cross-platform support

Tech stack: Go 1.21+, PostgreSQL 12+, Cobra, Viper, GORM

Phase 1 MVP in progress - account management complete"
fi

# Set up remote
echo "Setting up remote..."
git remote remove origin 2>/dev/null || true
git remote add origin "$REPO_URL"

# Push to GitHub
echo "Pushing to GitHub..."
git push -u origin main --force

echo ""
echo "========================================="
echo "✅ Successfully pushed to GitHub!"
echo "========================================="
echo ""
echo "Repository: https://github.com/jcaldwell1066/fintrack"
echo ""
echo "Next steps:"
echo "1. Clone to your projects directory:"
echo "   git clone https://github.com/jcaldwell1066/fintrack.git ~/projects/fintrack"
echo ""
echo "2. Set up development environment:"
echo "   cd ~/projects/fintrack"
echo "   make deps"
echo "   createdb fintrack"
echo "   psql -d fintrack -f migrations/001_initial_schema.sql"
echo "   cp fintrack_config.example.yaml ~/.config/fintrack/config.yaml"
echo ""
echo "3. Build and test:"
echo "   make build"
echo "   make test"
echo "   ./bin/fintrack account list"
echo ""
