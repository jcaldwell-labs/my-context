# My-Context-Copilot

A cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps.

## Features

- 🚀 **Simple Context Management**: Start, stop, and switch between work contexts effortlessly
- 📝 **Notes & Timestamps**: Capture decisions and track activity within each context
- 📁 **File Associations**: Link relevant files to contexts for easy reference
- 🔄 **Automatic Transitions**: Previous context stops automatically when starting a new one
- 📊 **History Tracking**: Full audit trail of all context switches
- 💾 **Plain Text Storage**: All data stored as human-readable text files
- 🌐 **Cross-Platform**: Works on Windows, Linux, and macOS
- 🔧 **Unix Philosophy**: Composable commands with text I/O

## Installation

### From Source

```bash
# Clone the repository
git clone <repo-url>
cd my-context-copilot

# Build
go build -o my-context.exe ./cmd/my-context/

# Install (optional)
./scripts/install.sh
```

### Pre-built Binaries

Download pre-built binaries from the [releases page](releases) for your platform.

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

### `list` (alias: `l`)
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

See [IMPLEMENTATION.md](IMPLEMENTATION.md) for implementation details and [HERE.md](HERE.md) for development shortcuts.

### Build
```bash
go build -o my-context.exe ./cmd/my-context/
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

1. Follow the Unix philosophy: Do one thing well
2. Keep it simple: No unnecessary complexity
3. Test first: TDD approach required
4. Plain text: All storage must remain human-readable

## License

[Add your license here]

## See Also

- [SETUP.md](SETUP.md) - Installation and setup guide
- [IMPLEMENTATION.md](IMPLEMENTATION.md) - Implementation roadmap
- [HERE.md](HERE.md) - Development scratchpad
- [specs/](specs/) - Feature specifications
