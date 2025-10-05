# HERE.md - Development Scratchpad

> **Quick reference scratchpad** for developers working on my-context-copilot.  
> Companion to README.md with executable snippets, FAQs, and development shortcuts.

**Last Updated**: 2025-10-04

---

## 🚀 Quick Start

```bash
# Clone and setup
git clone <repo-url> my-context-copilot
cd my-context-copilot

# Initialize Go module
go mod init github.com/yourusername/my-context-copilot
go mod tidy

# Install dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/stretchr/testify@latest

# Build binary
go build -o my-context ./cmd/my-context/

# Install locally
cp my-context /usr/local/bin/  # macOS/Linux
# Or add to PATH on Windows
```

---

## 📁 Project Structure Quick Reference

```bash
# Navigate to key areas
cd cmd/my-context/              # Main entry point
cd internal/commands/           # Command implementations
cd internal/core/               # Business logic
cd internal/models/             # Data structures
cd tests/integration/           # Integration tests
cd specs/001-cli-context-management/  # Feature docs
```

---

## 🔧 Common Development Tasks

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test file
go test ./tests/integration/commands_test.go

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Build for Different Platforms

```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o my-context.exe ./cmd/my-context/

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o my-context-linux ./cmd/my-context/

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o my-context-macos ./cmd/my-context/

# Build for macOS ARM (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o my-context-macos-arm ./cmd/my-context/
```

### Debugging

```bash
# Run with verbose logging (when implemented)
MY_CONTEXT_DEBUG=1 my-context start "Debug test"

# Check what files are being created
ls -la ~/.my-context/
find ~/.my-context/ -type f -exec ls -lh {} \;

# Watch file changes in real-time
watch -n 1 'ls -lh ~/.my-context/'

# Inspect context storage
cat ~/.my-context/state.json | jq .
cat ~/.my-context/transitions.log
cat ~/.my-context/Bug_fix/meta.json | jq .
```

---

## 🧪 Manual Testing Shortcuts

### Quick Smoke Test

```bash
# Clean slate
rm -rf ~/.my-context/

# Basic flow
my-context start "Test context"
my-context note "Test note"
my-context file README.md
my-context touch
my-context show
my-context list
my-context history
my-context stop

# Verify cleanup
my-context show  # Should show "No active context"
```

### Test Duplicate Names

```bash
my-context start "Duplicate"
my-context stop
my-context start "Duplicate"  # Creates Duplicate_2
my-context stop
my-context start "Duplicate"  # Creates Duplicate_3
my-context list | grep Duplicate
```

### Test JSON Output

```bash
# Start and capture JSON
my-context start "JSON test" --json | jq .

# Show with JSON
my-context show --json | jq .

# Extract specific fields
my-context show --json | jq -r '.data.context.name'
my-context list --json | jq '.data.contexts | length'
my-context history --json | jq '.data.transitions[-1]'
```

### Test Cross-Shell

```bash
# In git-bash
my-context start "Cross shell test"
my-context note "From git-bash"

# Switch to cmd.exe (Windows)
my-context show
my-context note "From cmd"

# Back to git-bash
my-context show  # Should show both notes
```

---

## 🐛 Troubleshooting

### Context Home Directory Issues

```bash
# Check where contexts are stored
echo $MY_CONTEXT_HOME  # Should show path or be empty (defaults to ~/.my-context)

# Override default location
export MY_CONTEXT_HOME=/custom/path
my-context start "Test"

# Reset to default
unset MY_CONTEXT_HOME
```

### Permission Issues

```bash
# Check directory permissions
ls -ld ~/.my-context/
ls -la ~/.my-context/

# Fix permissions (Unix/Linux/macOS)
chmod 700 ~/.my-context/
chmod 600 ~/.my-context/*.json
chmod 600 ~/.my-context/*.log
```

### Corrupted State Recovery

```bash
# Backup current state
cp -r ~/.my-context/ ~/.my-context.backup/

# Manually inspect state
cat ~/.my-context/state.json | jq .

# Fix corrupted state.json (reset to no active context)
echo '{"active_context": null, "last_updated": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}' > ~/.my-context/state.json

# Verify contexts are still accessible
my-context list
```

---

## 📊 Data Inspection Queries

### Find All Contexts

```bash
# List all context directories
ls -1d ~/.my-context/*/ | xargs -n1 basename

# Count total contexts
ls -1d ~/.my-context/*/ | wc -l

# Find contexts with most notes
for dir in ~/.my-context/*/; do
  count=$(wc -l < "$dir/notes.log" 2>/dev/null || echo 0)
  echo "$count $(basename "$dir")"
done | sort -rn
```

### Search Notes

```bash
# Find all notes containing "bug"
grep -r "bug" ~/.my-context/*/notes.log

# Find notes from today
today=$(date +%Y-%m-%d)
grep "^${today}" ~/.my-context/*/notes.log

# Find most recent note across all contexts
find ~/.my-context -name notes.log -exec tail -1 {} \; | sort
```

### Analyze Activity

```bash
# Count total touches across all contexts
find ~/.my-context -name touch.log -exec wc -l {} \; | awk '{sum+=$1} END {print sum}'

# Find most active context (by touches)
for dir in ~/.my-context/*/; do
  count=$(wc -l < "$dir/touch.log" 2>/dev/null || echo 0)
  echo "$count $(basename "$dir")"
done | sort -rn | head -5

# Show transition timeline
cat ~/.my-context/transitions.log | awk -F'|' '{print $1, $4, $2, "→", $3}' | tail -20
```

---

## 🔍 Development Queries

### Check Go Module Dependencies

```bash
# List all dependencies
go list -m all

# Check for updates
go list -u -m all

# Update dependencies
go get -u ./...
go mod tidy
```

### Code Statistics

```bash
# Count lines of code
find internal cmd -name '*.go' | xargs wc -l | tail -1

# Count lines by directory
find internal -name '*.go' | xargs wc -l | sort -n

# Find TODOs
grep -rn "TODO" internal/ cmd/

# Find FIXMEs
grep -rn "FIXME" internal/ cmd/
```

---

## 📝 Specification Quick Reference

### View Current Feature Docs

```bash
# Open spec in default editor
code specs/001-cli-context-management/spec.md

# View plan
cat specs/001-cli-context-management/plan.md | less

# View contracts
ls specs/001-cli-context-management/contracts/

# Run quickstart tests
cat specs/001-cli-context-management/quickstart.md
```

### Check Constitution Compliance

```bash
# View constitution
cat .specify/memory/constitution.md | less

# Check if design violates principles (manual review)
# I. Unix Philosophy - Each command single-purpose? ✓
# II. Cross-Platform - Works in cmd/bash/WSL? ✓
# III. Stateful Context - One active at a time? ✓
# IV. Minimal Surface Area - ≤10 commands? ✓
# V. Data Portability - Plain text? ✓
```

---

## 🎯 Performance Benchmarking

### Measure Command Speed

```bash
# Time a single command
time my-context start "Benchmark"

# Average over 100 runs
for i in {1..100}; do
  time my-context note "Test $i" 2>&1
done | grep real | awk '{sum+=$2} END {print sum/NR}'

# Benchmark with many contexts
for i in {1..1000}; do
  my-context start "Context $i"
  my-context stop
done
time my-context list  # Should be <50ms even with 1000 contexts
```

### Check Binary Size

```bash
# Check compiled binary size
ls -lh my-context

# Strip symbols to reduce size
go build -ldflags="-s -w" -o my-context ./cmd/my-context/
ls -lh my-context
```

---

## 🧹 Cleanup & Reset

### Clean Test Data

```bash
# Remove all contexts
rm -rf ~/.my-context/

# Remove specific context
rm -rf ~/.my-context/Context_name/

# Keep contexts but reset state
echo '{"active_context": null, "last_updated": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"}' > ~/.my-context/state.json
```

### Clean Build Artifacts

```bash
# Remove binary
rm -f my-context my-context.exe

# Clean Go cache
go clean -cache -modcache -testcache

# Remove test coverage files
rm -f coverage.out
```

---

## 📚 FAQ

### Q: Where is context data stored?

**A**: By default in `~/.my-context/` (or `$MY_CONTEXT_HOME` if set).

```bash
# View current location
my-context show --json | jq -r '.data.context.subdirectory_path' | xargs dirname | xargs dirname

# Or check environment
echo ${MY_CONTEXT_HOME:-~/.my-context}
```

### Q: How do I backup my contexts?

**A**: Copy the entire home directory.

```bash
# Backup
cp -r ~/.my-context/ ~/my-context-backup-$(date +%Y%m%d)/

# Or use tar
tar czf my-context-backup.tar.gz ~/.my-context/

# Restore
tar xzf my-context-backup.tar.gz -C ~/
```

### Q: Can I version control my contexts with git?

**A**: Yes! The plain text format is git-friendly.

```bash
cd ~/.my-context/
git init
git add .
git commit -m "Initial context snapshot"

# After working
git add .
git commit -m "End of day checkpoint"
```

### Q: How do I migrate to a different machine?

**A**: Copy the home directory or use version control.

```bash
# On old machine
tar czf contexts.tar.gz ~/.my-context/

# Transfer file to new machine, then:
tar xzf contexts.tar.gz -C ~/
```

### Q: What if I want different home directories for work vs personal?

**A**: Use the `MY_CONTEXT_HOME` environment variable.

```bash
# Work contexts
export MY_CONTEXT_HOME=~/work-contexts
my-context start "Work project"

# Personal contexts
export MY_CONTEXT_HOME=~/personal-contexts
my-context start "Side project"

# Or use shell aliases
alias work-context='MY_CONTEXT_HOME=~/work-contexts my-context'
alias personal-context='MY_CONTEXT_HOME=~/personal-contexts my-context'
```

### Q: How do I integrate with git hooks?

**A**: Add to your git hooks.

```bash
# .git/hooks/post-commit
#!/bin/bash
my-context note "Committed: $(git log -1 --pretty=%B)"
my-context file "$(git diff-tree --no-commit-id --name-only -r HEAD | head -1)"

# Make executable
chmod +x .git/hooks/post-commit
```

### Q: Can I export context data to CSV or JSON?

**A**: Yes, use the JSON output and jq.

```bash
# Export all contexts to JSON
my-context list --json > contexts.json

# Convert to CSV (simple)
my-context list --json | jq -r '.data.contexts[] | [.name, .start_time, .status, .note_count] | @csv'

# Export specific context with all details
my-context show --json > context-snapshot.json
```

---

## 🔗 Useful Links

- [Feature Spec](specs/001-cli-context-management/spec.md)
- [Implementation Plan](specs/001-cli-context-management/plan.md)
- [Quickstart Tests](specs/001-cli-context-management/quickstart.md)
- [Constitution](.specify/memory/constitution.md)
- [Cobra Docs](https://github.com/spf13/cobra)
- [Go Documentation](https://golang.org/doc/)

---

## 💡 Tips & Tricks

```bash
# Create a context for every git branch
git checkout -b feature/new-feature
my-context start "$(git branch --show-current)"

# Auto-start context when entering directory (add to .bashrc)
function cd() {
  builtin cd "$@"
  if [ -f ".context" ]; then
    my-context start "$(cat .context)"
  fi
}

# Quick note alias
alias cn='my-context note'
alias cf='my-context file'

# Show context in prompt (add to PS1)
export PS1='[$(my-context show --json 2>/dev/null | jq -r ".data.context.name // \"no context\"")]$ '
```

---

**Note**: This is a living document. Add your own snippets and shortcuts as you discover them!
