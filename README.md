# My-Context

> **🔧 DEVELOPMENT BRANCH**
> This is the **development reference branch** showing the full development process and internal tooling.
> **Looking for the clean public version?** → See [main branch](https://github.com/jcaldwell1066/my-context/tree/main)

A cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps.

## About This Branch

This `dev` branch contains the complete development environment including:

- **Spec Kit Workflow** (`.specify/`) - Feature specification and planning tools
- **SDLC Documentation** (`SDLC.md`, `IMPLEMENTATION.md`) - Development process
- **Feature Specifications** (`specs/`) - Complete spec-driven development artifacts
- **Development Automation** (`scripts/`) - Build, deploy, and release automation
- **Internal Tooling** (`.claude/`, `.cursor/`) - AI-assisted development commands
- **Sprint Retrospectives** - Lessons learned and process improvements
- **Constitution & Governance** - Project principles and decision-making

**Use this branch to:**
- Understand how the project was developed
- See Spec Kit workflow in action
- Reference SDLC best practices
- Learn from sprint retrospectives
- Use automation scripts for your own projects

## Features

- 🚀 **Simple Context Management**: Start, stop, and switch between work contexts effortlessly
- 📝 **Notes & Timestamps**: Capture decisions and track activity within each context
- 📁 **File Associations**: Link relevant files to contexts for easy reference
- 🔄 **Automatic Transitions**: Previous context stops automatically when starting a new one
- 📊 **History Tracking**: Full audit trail of all context switches
- 🚨 **Signal Coordination**: Event-driven coordination between team members and processes
- 👀 **Context Watching**: Monitor contexts for changes and execute commands automatically
- 🏷️ **Context Metadata**: Enhanced context organization with created-by, parent, and labels
- 💾 **Plain Text Storage**: All data stored as human-readable text files
- 🌐 **Cross-Platform**: Works on Windows, Linux, and macOS
- 🔧 **Unix Philosophy**: Composable commands with text I/O

## Installation

### Quick Install (Recommended)

**One-liner install** (Linux, macOS, WSL):
```bash
# Coming soon - check releases page for now
# curl -sSL https://raw.githubusercontent.com/YOUR-USERNAME/my-context/main/scripts/curl-install.sh | bash
```

### Pre-built Binaries

Download binaries from the [releases page](https://github.com/YOUR-USERNAME/my-context/releases):

- **Windows**: `my-context-windows-amd64.exe`
- **Linux**: `my-context-linux-amd64`
- **macOS Intel**: `my-context-darwin-amd64`
- **macOS ARM (M1/M2)**: `my-context-darwin-arm64`

**Installation scripts**:
- Linux/macOS/WSL: `./scripts/install.sh`
- Windows (cmd.exe): `scripts\install.bat`
- Windows (PowerShell): `.\scripts\install.ps1`

### Building from Source

**Requirements**: Go 1.21 or later

```bash
# Clone the repository
git clone https://github.com/YOUR-USERNAME/my-context.git
cd my-context

# Build for current platform
go build -o my-context ./cmd/my-context/

# Or build all platforms
./scripts/build-all.sh

# Install to ~/.local/bin (Unix) or %USERPROFILE%\bin (Windows)
./scripts/install.sh
```

**Build options**:
- Single platform: `go build -o my-context ./cmd/my-context/`
- All platforms: `./scripts/build-all.sh` (outputs to `bin/`)
- Static linking: Binaries are statically linked (CGO_ENABLED=0) with zero runtime dependencies

### Troubleshooting

If you encounter installation issues, please open an issue on GitHub.

## Quick Start

```bash
# Start a new context
my-context start "Working on login feature"

# Add notes
my-context note "Fixed authentication bug"
my-context note "Updated tests"

# Associate files
my-context file src/auth/login.go
my-context file tests/auth/login_test.go

# Record activity
my-context touch

# View current context
my-context show

# List all contexts
my-context list

# View transition history
my-context history

# Stop current context
my-context stop
```

## Commands

All commands support both full names and single-letter aliases.

### `start <name>` (alias: `s`)
Create and activate a new context. If a context is already active, it will be stopped automatically.

```bash
my-context start "Bug fix #123"
my-context s "Quick patch"
```

**Duplicate Names**: If a context with the same name exists, a suffix (_2, _3, etc.) is automatically appended.

### `start <name>` - Enhanced (Sprint 2)

**New: `--project` flag**
```bash
# Create context with project prefix
my-context start "Phase 1" --project ps-cli
# Creates: "ps-cli: Phase 1"

# Quickly group related contexts
my-context start "Planning" --project garden
my-context start "Implementation" --project garden
```

### `stop` (alias: `p`)
Stop the currently active context without starting a new one.

```bash
my-context stop
```

### `note <text>` (alias: `n`)
Add a timestamped note to the active context.

```bash
my-context note "Identified the root cause"
my-context n "Refactored error handling"
```

### `file <path>` (alias: `f`)
Associate a file path with the active context. Paths are normalized automatically.

```bash
my-context file src/main.go
my-context f /absolute/path/to/file.txt
```

### `touch` (alias: `t`)
Record a timestamp indicating activity without detailed notes.

```bash
my-context touch
```

### `show` (alias: `w`)
Display details about the currently active context.

```bash
my-context show
```

Output includes:
- Context name and status
- Start time and duration
- All notes with timestamps
- Associated files
- Touch event count

### `list` (alias: `l`) - Enhanced (Sprint 2)

**New filtering options**:
```bash
# Default: show 10 most recent
my-context list

# Show all contexts
my-context list --all

# Custom limit
my-context list --limit 5

# Filter by project (case-insensitive)
my-context list --project ps-cli

# Search by substring
my-context list --search "Phase"

# Show only archived contexts
my-context list --archived

# Show only active context
my-context list --active-only

# Combine filters
my-context list --project garden --limit 3
```

### `export <name>` (alias: `e`) - New in Sprint 2

Export context data to markdown or JSON for sharing and documentation.

```bash
# Export single context (creates {name}.md)
my-context export "ps-cli: Phase 1"

# Export to custom path
my-context export "Phase 1" --to reports/phase-1.md

# Export all contexts
my-context export --all --to exports/

# Export as JSON
my-context export "Phase 1" --json --to data.json
```

**Export format** (Markdown):
- Context metadata (start/end times, duration)
- All notes with timestamps
- File associations
- Touch activity
- Archive status

### `archive <name>` - New in Sprint 2

Mark a context as archived. Archived contexts are hidden from default `list` output but preserved.

```bash
# Archive a stopped context
my-context archive "old-project: Phase 1"

# View archived contexts
my-context list --archived
```

**Requirements**:
- Context must be stopped (cannot archive active context)
- Data is preserved (only metadata flag changed)

### `delete <name>` - New in Sprint 2

Permanently delete a context directory.

```bash
# Delete with confirmation prompt
my-context delete "old-context"

# Skip confirmation
my-context delete "old-context" --force
```

**Safety features**:
- Cannot delete active context
- Confirmation prompt (unless `--force`)
- Exit code 2 on user cancellation
- `transitions.log` is preserved (historical transitions remain)

### `list` (legacy alias)
List all contexts (active and stopped) with their status.

```bash
my-context list
```

### `history` (alias: `h`)
Show the chronological transition log of all context switches.

```bash
my-context history
```

## JSON Output

All commands support `--json` flag for machine-readable output:

```bash
my-context show --json | jq .
my-context list --json | jq '.data.contexts | length'
```

## Data Storage

All context data is stored in `~/.my-context/` (or `$MY_CONTEXT_HOME` if set):

```
~/.my-context/
├── state.json              # Active context pointer
├── transitions.log         # Transition history
└── Context_Name/           # Per-context directory
    ├── meta.json           # Context metadata
    ├── notes.log           # Timestamped notes
    ├── files.log           # File associations
    └── touch.log           # Activity timestamps
```

All files are plain text and can be:
- ✅ Viewed with standard tools (`cat`, `grep`, `less`)
- ✅ Version controlled with git
- ✅ Backed up with simple file copy
- ✅ Edited manually if needed

## Environment Variables

- `MY_CONTEXT_HOME`: Override the default context storage directory

```bash
export MY_CONTEXT_HOME=/custom/path
my-context start "Test"
```

## Use Cases

### Track Work Sessions
```bash
my-context start "Sprint planning"
my-context note "Discussed user stories for Q1"
my-context note "Priority: authentication refactor"
my-context stop
```

### Context Switching
```bash
my-context start "Feature development"
# ... work on feature ...

# Urgent bug report comes in
my-context start "Hotfix bug #456"
# Previous context automatically stopped
```

### Git Integration
```bash
# In git hooks (.git/hooks/post-commit)
my-context note "Committed: $(git log -1 --pretty=%B)"
```

### End-of-Day Review
```bash
# See what you worked on today
my-context list
my-context history
```

## Tutorials & Guides

### **📚 [Personal Productivity with Triggers Tutorial](docs/TRIGGERS-TUTORIAL.md)**

Learn how to automate your workflow with signals and watches:
- Get notified when you complete milestones
- Receive warnings when contexts get too large
- Automate end-of-day cleanup
- Integrate with Pomodoro timers
- Create smart context lifecycle automation

**Perfect for**: Anyone wanting to level up their productivity with automation

## Advanced Usage

### Backup Your Contexts
```bash
tar czf contexts-backup.tar.gz ~/.my-context/
```

### Version Control Your Context Data
```bash
cd ~/.my-context
git init
git add .
git commit -m "Context snapshot"
```

### Search Notes
```bash
grep -r "bug" ~/.my-context/*/notes.log
```

### Shell Integration
Add to `.bashrc` or `.zshrc`:

```bash
# Show current context in prompt
export PS1='[$(my-context show --json 2>/dev/null | jq -r ".data.context.name // \"no context\"")]$ '

# Quick aliases
alias cn='my-context note'
alias cf='my-context file'
alias cs='my-context show'
```

## Development

### Build
```bash
go build -o my-context ./cmd/my-context/
```

### Test
```bash
go test ./...
```

### Cross-Platform Build
```bash
./scripts/build.sh
```

## Architecture

- **Models** (`internal/models/`): Data structures for contexts, notes, files, etc.
- **Core** (`internal/core/`): Business logic for context operations
- **Commands** (`internal/commands/`): CLI command implementations
- **Output** (`internal/output/`): Human-readable and JSON formatters

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Key principles:
1. Follow the Unix philosophy: Do one thing well
2. Keep it simple: No unnecessary complexity
3. Test first: TDD approach required
4. Plain text: All storage must remain human-readable

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
