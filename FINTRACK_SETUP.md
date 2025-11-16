# FinTrack Repository Setup Guide

This guide will help you create a separate GitHub repository for FinTrack.

## Quick Start (Automated)

### Step 1: Create GitHub Repository

1. Go to https://github.com/new
2. **Repository name:** `fintrack`
3. **Description:** Terminal-based personal finance tracking and budgeting CLI
4. **Visibility:** Public or Private (your choice)
5. **Important:** Do NOT initialize with README, .gitignore, or license
6. Click "Create repository"
7. Copy the repository URL (e.g., `https://github.com/yourusername/fintrack.git`)

### Step 2: Run Setup Script

```bash
cd ~/my-context
./setup-fintrack-repo.sh https://github.com/yourusername/fintrack.git
```

This will:
- Create `~/projects/fintrack` directory
- Extract all FinTrack files
- Initialize git repository
- Push to GitHub

## Manual Setup (Alternative)

If you prefer to set up manually:

### Step 1: Create GitHub Repository

Same as above.

### Step 2: Extract and Push

```bash
# Create directory
mkdir -p ~/projects/fintrack
cd ~/projects/fintrack

# Extract files
tar -xzf ~/my-context/fintrack-repository.tar.gz --strip-components=1

# Initialize git (if not already)
git init
git branch -M main

# Add remote
git remote add origin https://github.com/yourusername/fintrack.git

# Push to GitHub
git add .
git commit -m "Initial commit: FinTrack Phase 1 foundation"
git push -u origin main
```

## What's Included

The repository contains:

```
fintrack/
├── cmd/fintrack/              # CLI entry point
├── internal/                  # Core application code
│   ├── commands/             # Command implementations
│   ├── config/               # Configuration management
│   ├── db/                   # Database layer
│   ├── models/               # Data models
│   └── output/               # Output formatters
├── tests/unit/               # Unit tests
├── docs/                     # Documentation
│   ├── FINANCE_TRACKER_PLAN.md
│   ├── FINTRACK_QUICKREF.md
│   └── FINTRACK_ROADMAP.md
├── migrations/               # Database migrations
│   └── 001_initial_schema.sql
├── Makefile                  # Build automation
├── README.md                 # Project README
├── go.mod                    # Go dependencies
├── .gitignore               # Git ignore rules
└── fintrack_config.example.yaml  # Example config
```

## After Setup

Once you've pushed to GitHub:

### 1. Set Up Development Environment

```bash
cd ~/projects/fintrack

# Install Go dependencies
make deps

# Set up PostgreSQL database
createdb fintrack
psql -d fintrack -f migrations/001_initial_schema.sql

# Create config file
mkdir -p ~/.config/fintrack
cp fintrack_config.example.yaml ~/.config/fintrack/config.yaml

# Edit config to set your database connection
nano ~/.config/fintrack/config.yaml
```

### 2. Build and Test

```bash
# Build the binary
make build

# Run tests
make test

# Try it out!
./bin/fintrack --version
./bin/fintrack account add "Test Account" --type checking --balance 1000
./bin/fintrack account list
```

### 3. Install (Optional)

```bash
make install
# Now you can use: fintrack account list
```

## Cleanup (Optional)

After successfully setting up the new repository, you can clean up the `my-context` repository:

```bash
cd ~/my-context

# Remove fintrack subdirectory (now in its own repo)
git rm -r fintrack/

# Remove planning docs (now in docs/ of new repo)
git rm FINANCE_TRACKER_PLAN.md
git rm FINTRACK_QUICKREF.md
git rm FINTRACK_ROADMAP.md
git rm fintrack_config.example.yaml
git rm fintrack_schema.sql

# Remove setup files
rm fintrack-repository.tar.gz
rm setup-fintrack-repo.sh
rm FINTRACK_SETUP.md

# Commit cleanup
git commit -m "chore: Move FinTrack to separate repository

FinTrack now lives at: https://github.com/yourusername/fintrack"

git push
```

## Alternative: Keep Both

If you want to keep FinTrack development visible in my-context:

### Option A: Git Submodule

```bash
cd ~/my-context
git rm -r fintrack/  # Remove tracked files
git submodule add https://github.com/yourusername/fintrack.git fintrack
git commit -m "refactor: Convert fintrack to git submodule"
git push
```

Now `fintrack/` is linked but lives in its own repository.

### Option B: Keep as Documentation

Keep the planning docs in my-context for reference, but develop in the separate repo:

```bash
# In my-context
mkdir -p docs/projects
git mv fintrack-repository.tar.gz docs/projects/
git mv FINANCE_TRACKER_PLAN.md docs/projects/
git mv FINTRACK_*.md docs/projects/
git commit -m "docs: Archive FinTrack planning documentation"
```

## Troubleshooting

### "repository not found" error

Make sure you created the repository on GitHub first and copied the correct URL.

### Permission denied

Make sure you're authenticated with GitHub:
```bash
git config --global user.name "Your Name"
git config --global user.email "your@email.com"

# If using HTTPS, you may need a personal access token
# Settings > Developer settings > Personal access tokens
```

### Can't extract tarball

Make sure you're in the correct directory:
```bash
ls -l ~/my-context/fintrack-repository.tar.gz
```

## Need Help?

- **GitHub Docs:** https://docs.github.com/en/get-started/quickstart/create-a-repo
- **FinTrack README:** In the repository after setup
- **FinTrack Docs:** `docs/` directory in the repository

---

**Ready to build the future of terminal-based personal finance!** 🚀
