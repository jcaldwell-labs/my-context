# Pushing FinTrack to GitHub

Your GitHub repository `jcaldwell1066/fintrack` has been created! Here's how to push the FinTrack code to it.

## Quick Method (Recommended)

Since Claude Code has authentication limitations, the easiest way is to extract and push from your **local terminal** (outside of Claude Code):

### Step 1: Extract the Repository

Open a terminal on your local machine and run:

```bash
# Create directory for FinTrack
mkdir -p ~/projects/fintrack
cd ~/projects/fintrack

# Extract the prepared repository
# (Assuming you've synced my-context from Claude Code)
tar -xzf ~/my-context/fintrack-repository.tar.gz --strip-components=1
```

### Step 2: Push to GitHub

```bash
cd ~/projects/fintrack

# If git isn't initialized (it should be)
git init
git branch -M main

# Add your GitHub repository
git remote add origin https://github.com/jcaldwell1066/fintrack.git

# Push to GitHub
git push -u origin main
```

If the repository already has a README, you may need to force push:

```bash
git push -u origin main --force
```

### Step 3: Verify

Visit https://github.com/jcaldwell1066/fintrack to see your code!

---

## Alternative: Use the Automated Script

I've created a script that automates these steps:

```bash
# From your local terminal (outside Claude Code)
cd ~/my-context
./push-fintrack.sh
```

This script will:
- Extract the FinTrack repository
- Initialize git if needed
- Add the GitHub remote
- Push to origin/main

---

## What Gets Pushed

Your repository will contain:

```
fintrack/
├── cmd/fintrack/              # ✅ CLI application
├── internal/                  # ✅ Core code
│   ├── commands/             # Account management (working!)
│   ├── config/               # Configuration
│   ├── db/                   # Database layer
│   ├── models/               # Data models
│   └── output/               # Formatters
├── tests/unit/               # ✅ Unit tests
├── docs/                     # ✅ Documentation
│   ├── FINANCE_TRACKER_PLAN.md
│   ├── FINTRACK_QUICKREF.md
│   └── FINTRACK_ROADMAP.md
├── migrations/               # ✅ Database schema
│   └── 001_initial_schema.sql
├── Makefile                  # ✅ Build system
├── README.md                 # ✅ Project README
├── go.mod                    # ✅ Dependencies
└── fintrack_config.example.yaml
```

**Total:** 18 files, 5,417 lines of code and documentation

---

## After Pushing

Once the code is on GitHub, set up your development environment:

### 1. Clone the Repository

```bash
# Clone to a clean location
git clone https://github.com/jcaldwell1066/fintrack.git ~/projects/fintrack
cd ~/projects/fintrack
```

### 2. Install Dependencies

```bash
# Download Go dependencies
make deps

# Or manually:
go mod download
```

### 3. Set Up Database

```bash
# Create PostgreSQL database
createdb fintrack

# Run initial schema migration
psql -d fintrack -f migrations/001_initial_schema.sql
```

### 4. Configure FinTrack

```bash
# Create config directory
mkdir -p ~/.config/fintrack

# Copy example config
cp fintrack_config.example.yaml ~/.config/fintrack/config.yaml

# Edit with your database settings
nano ~/.config/fintrack/config.yaml
```

Update the database URL:
```yaml
database:
  url: "postgresql://localhost:5432/fintrack?sslmode=disable"
```

### 5. Build and Test

```bash
# Build the binary
make build

# Run tests
make test

# View test coverage
make test-coverage
open coverage.html
```

### 6. Try It Out!

```bash
# Run from build directory
./bin/fintrack --version
./bin/fintrack account add "Test Account" --type checking --balance 1000
./bin/fintrack account list
./bin/fintrack account show 1

# Or install globally
make install
fintrack account list

# Try JSON output
fintrack account list --json
```

---

## Troubleshooting

### "Permission denied" when pushing

Make sure you're authenticated with GitHub:

```bash
# Check your git config
git config --global user.name
git config --global user.email

# If using HTTPS, you may need a Personal Access Token
# GitHub Settings > Developer settings > Personal access tokens > Tokens (classic)
# Generate new token with 'repo' scope
```

### "Updates were rejected" when pushing

The GitHub repo might have a README. Force push to replace it:

```bash
git push -u origin main --force
```

### Can't extract the tarball

Make sure the file exists:

```bash
ls -lh ~/my-context/fintrack-repository.tar.gz
```

If it's not there, you may need to sync your my-context repository first.

---

## Next Steps After Setup

1. **Continue Phase 1 Development:**
   - Implement transaction management
   - Add category hierarchy
   - CSV import functionality
   - Basic reporting

2. **Set Up GitHub Features:**
   - Add topics/tags to repository
   - Create GitHub issues for Phase 1 tasks
   - Set up branch protection rules
   - Configure GitHub Actions (CI/CD)

3. **Community:**
   - Add LICENSE file (MIT recommended)
   - Create CONTRIBUTING.md
   - Set up issue templates
   - Add GitHub badges to README

---

## Useful Links

- **Repository:** https://github.com/jcaldwell1066/fintrack
- **Documentation:** See `docs/` directory
- **Issues:** https://github.com/jcaldwell1066/fintrack/issues
- **Project Board:** https://github.com/jcaldwell1066/fintrack/projects

---

**Ready to build!** 🚀

Once you push to GitHub and set up your development environment, you'll have a fully functional personal finance CLI with account management, ready for Phase 1 expansion.
