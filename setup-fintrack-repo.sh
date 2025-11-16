#!/bin/bash
# Setup script for creating a new FinTrack repository

echo "========================================="
echo "FinTrack Repository Setup"
echo "========================================="
echo ""

# Check if GitHub repository URL is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <github-repo-url>"
    echo ""
    echo "Example:"
    echo "  $0 https://github.com/yourusername/fintrack.git"
    echo ""
    echo "Steps:"
    echo "1. Create a new repository on GitHub (https://github.com/new)"
    echo "   - Name: fintrack"
    echo "   - Don't initialize with README"
    echo "2. Copy the repository URL"
    echo "3. Run this script with the URL"
    exit 1
fi

REPO_URL="$1"
WORK_DIR="$HOME/projects/fintrack"

echo "Repository URL: $REPO_URL"
echo "Working directory: $WORK_DIR"
echo ""

# Ask for confirmation
read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# Create working directory
echo "Creating directory..."
mkdir -p "$WORK_DIR"
cd "$WORK_DIR" || exit 1

# Extract the prepared repository
echo "Extracting FinTrack files..."
tar -xzf ~/my-context/fintrack-repository.tar.gz --strip-components=1

# Initialize git if needed
if [ ! -d .git ]; then
    echo "Initializing git repository..."
    git init
    git branch -M main
fi

# Add remote
echo "Adding remote origin..."
git remote add origin "$REPO_URL" 2>/dev/null || git remote set-url origin "$REPO_URL"

# Commit and push
echo "Pushing to GitHub..."
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

Phase 1 MVP in progress - account management complete" || echo "Commit already exists"

git push -u origin main

echo ""
echo "========================================="
echo "✅ FinTrack repository setup complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. cd $WORK_DIR"
echo "2. Set up PostgreSQL: createdb fintrack"
echo "3. Run migrations: psql -d fintrack -f migrations/001_initial_schema.sql"
echo "4. Configure: cp fintrack_config.example.yaml ~/.config/fintrack/config.yaml"
echo "5. Build: make build"
echo "6. Test: make test"
echo "7. Run: ./bin/fintrack account list"
echo ""
echo "Documentation:"
echo "- README.md - Quick start guide"
echo "- docs/FINANCE_TRACKER_PLAN.md - Complete design"
echo "- docs/FINTRACK_QUICKREF.md - Command reference"
echo "- docs/FINTRACK_ROADMAP.md - Implementation timeline"
echo ""
