# My-Context

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

> *Track what you're working on, not just what you've done.*

A cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps.

## Why My-Context?

Most developers context-switch constantly but lose track of *why* they made decisions, *what* files were involved, and *when* they stopped mid-task. My-context solves this by creating lightweight, timestamped work journals that follow you across sessions.

**Key benefits:**
- **Instant context capture** — Start tracking in one command, add notes as you go
- **Automatic transitions** — Switching contexts stops the previous one automatically
- **Plain text storage** — All data stored as human-readable files you can grep, version control, or edit
- **Zero lock-in** — No database, no cloud dependency, no vendor tie-in
- **Cross-platform** — Works identically on Windows, Linux, macOS, and WSL

**Perfect for:**
- Developers who context-switch frequently between tasks
- Teams wanting lightweight decision documentation
- Anyone who's asked "what was I doing on that ticket last week?"
- Pomodoro/time-boxing practitioners tracking work sessions
- AI-assisted coding workflows that need session memory

## Demo

```bash
# Start working on a feature
$ my-context start "Implement user auth"
Context started: Implement_user_auth

# Capture decisions as you go
$ my-context note "Using JWT tokens for stateless auth"
$ my-context note "Added refresh token rotation"

# Associate relevant files
$ my-context file src/auth/jwt.go
$ my-context file tests/auth_test.go

# Urgent bug? Just start a new context
$ my-context start "Hotfix: login timeout"
# Previous context automatically stopped

# Review what you worked on
$ my-context show
```

## Quick Start

### Install

**Linux/macOS/WSL (one-liner):**
```bash
curl -sSL https://raw.githubusercontent.com/jcaldwell-labs/my-context/main/scripts/curl-install.sh | bash
```

**Or download binaries** from [Releases](https://github.com/jcaldwell-labs/my-context/releases):
- `my-context-linux-amd64` — Linux (x86_64)
- `my-context-darwin-amd64` — macOS Intel
- `my-context-darwin-arm64` — macOS Apple Silicon
- `my-context-windows-amd64.exe` — Windows

**Build from source:**
```bash
git clone https://github.com/jcaldwell-labs/my-context.git
cd my-context
go build -o my-context ./cmd/my-context/
```

### Your First Context

```bash
# 1. Start a context
my-context start "My first task"

# 2. Add a note
my-context note "Getting started with my-context"

# 3. See your context
my-context show

# 4. List all contexts
my-context list
```

## Features

### Core Context Management
| Command | Alias | Description |
|---------|-------|-------------|
| `start <name>` | `s` | Create and activate a new context |
| `stop` | `p` | Stop the active context |
| `show` | `w` | Display current context details |
| `list` | `l` | List all contexts with filters |
| `history` | `h` | Show context transition history |

### Notes & Files
| Command | Alias | Description |
|---------|-------|-------------|
| `note <text>` | `n` | Add timestamped note to active context |
| `file <path>` | `f` | Associate file with active context |
| `touch` | `t` | Record activity timestamp |

### Organization (Sprint 2+)
| Command | Alias | Description |
|---------|-------|-------------|
| `export <name>` | `e` | Export context to markdown/JSON |
| `archive <name>` | `a` | Archive completed contexts |
| `delete <name>` | `d` | Permanently remove a context |

### Advanced Features

**Project Grouping:**
```bash
my-context start "Phase 1" --project myapp
my-context list --project myapp
```

**Filtering & Search:**
```bash
my-context list --limit 5
my-context list --search "auth"
my-context list --archived
```

**JSON Output:**
```bash
my-context show --json | jq .
```

## Use Cases

### Daily Development Workflow
```bash
# Morning: start your day
my-context start "Sprint 42 - User dashboard"
my-context note "Planning: add charts, fix pagination"

# As you work
my-context note "Charts library: chose recharts over chart.js"
my-context file src/components/Dashboard.tsx

# Context switch happens
my-context start "Code review: PR #234"
# Previous context auto-stopped with timestamp
```

### Git Integration
```bash
# In .git/hooks/post-commit
my-context note "Committed: $(git log -1 --pretty=%B)"
```

### End-of-Day Review
```bash
my-context list --limit 10
my-context export "Sprint 42 - User dashboard"
```

## Data Storage

All data stored in `~/.my-context/` as plain text:

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

Override location: `export MY_CONTEXT_HOME=/custom/path`

## Documentation

- **[Getting Started Guide](docs/guides/GETTING-STARTED.md)** — Full installation and first steps
- **[Triggers Tutorial](docs/tutorials/TRIGGERS-TUTORIAL.md)** — Automation with signals and watches
- **[Troubleshooting](docs/guides/TROUBLESHOOTING.md)** — Common issues and solutions
- **[CLAUDE.md](CLAUDE.md)** — Architecture and development guide

## Comparison

| Feature | my-context | git stash | tmux sessions | note apps |
|---------|-----------|-----------|---------------|-----------|
| Context tracking | Yes | No | No | Manual |
| Automatic transitions | Yes | No | No | No |
| File associations | Yes | Yes (staged) | No | Manual |
| Timestamped notes | Yes | No | No | Yes |
| Plain text storage | Yes | Binary | N/A | Varies |
| Cross-platform | Yes | Yes | Unix only | Varies |
| CLI-first | Yes | Yes | Yes | No |

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Development:**
```bash
go build -o my-context ./cmd/my-context/
go test ./...
```

**Key principles:**
1. Unix philosophy — Do one thing well
2. Plain text — All storage human-readable
3. Test-first — TDD approach required
4. Cross-platform — Windows, Linux, macOS

## Roadmap

See [.github/planning/ROADMAP.md](.github/planning/ROADMAP.md) for upcoming features.

## Community

- **Issues:** [GitHub Issues](https://github.com/jcaldwell-labs/my-context/issues)
- **Discussions:** [GitHub Discussions](https://github.com/jcaldwell-labs/my-context/discussions)

## License

[MIT License](LICENSE) — Free for personal and commercial use.
