# My-Context

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

> _Track what you're working on, not just what you've done._

A cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps.

## Why My-Context?

Most developers context-switch constantly but lose track of _why_ they made decisions, _what_ files were involved, and _when_ they stopped mid-task. My-context solves this by creating lightweight, timestamped work journals that follow you across sessions.

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

| Command        | Alias | Description                       |
| -------------- | ----- | --------------------------------- |
| `start <name>` | `s`   | Create and activate a new context |
| `stop`         | `p`   | Stop the active context           |
| `show`         | `w`   | Display current context details   |
| `list`         | `l`   | List all contexts with filters    |
| `history`      | `h`   | Show context transition history   |

### Notes & Files

| Command       | Alias | Description                            |
| ------------- | ----- | -------------------------------------- |
| `note <text>` | `n`   | Add timestamped note to active context |
| `file <path>` | `f`   | Associate file with active context     |
| `touch`       | `t`   | Record activity timestamp              |

### Organization (Sprint 2+)

| Command          | Alias | Description                     |
| ---------------- | ----- | ------------------------------- |
| `export <name>`  | `e`   | Export context to markdown/JSON |
| `archive <name>` | `a`   | Archive completed contexts      |
| `delete <name>`  | `d`   | Permanently remove a context    |

### Analytics & Automation

| Command  | Alias | Description                                             |
| -------- | ----- | ------------------------------------------------------- |
| `stats`  |       | Time tracking aggregation (--today, --week, --month)    |
| `record` | `r`   | Clipboard recording mode - auto-capture pastes as notes |

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
my-context list --stale  # Show only stale contexts (inactive 8h+)
```

**Stale Context Detection:**

Automatically detects when active contexts have been inactive for extended periods:

```bash
# View context with stale warnings
my-context show
# ⚠️ Warning: Context active for 5h with no recent activity
#    Last activity: 5h ago

# List only stale contexts (inactive 8h+ by default)
my-context list --stale

# Configure thresholds via environment variables
export MC_STALE_WARN_HOURS=4      # Show warning (default: 4h)
export MC_STALE_HOURS=8            # Mark as stale (default: 8h)
export MC_STALE_CRITICAL_HOURS=24  # Show critical warning (default: 24h)
```

**Warning Levels:**
- **Warn** (4h): Shows warning message on `show` command
- **Stale** (8h): Marked as stale, filterable with `--stale` flag
- **Critical** (24h): Shows critical warning suggesting to stop context

Last activity is calculated from the most recent:
- Note addition
- File association
- Touch event

Stopped contexts never show stale warnings.

**JSON Output:**

```bash
my-context show --json | jq .
# Includes: is_stale, stale_level, last_activity, inactive_since_hours
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

### Database Backend (PostgreSQL)

For teams, high-volume usage, or centralized context storage, my-context supports PostgreSQL:

```bash
# Enable database backend
export MY_CONTEXT_HOME=db
export DATABASE_URL="host=localhost port=5432 user=myuser password=mypass dbname=dev_state sslmode=disable"

# Or use connection URL format
export DATABASE_URL="postgres://myuser:mypass@localhost:5432/dev_state?sslmode=disable"
```

**Partitions** — Isolate contexts by project:

```bash
# Each partition gets its own schema
export MY_CONTEXT_HOME=db:project-alpha
export MY_CONTEXT_HOME=db:project-beta

# List all partitions
my-context partitions

# Query across partitions
my-context list --all-partitions
```

**When to use database backend:**

| Use Case          | File Backend   | Database Backend  |
| ----------------- | -------------- | ----------------- |
| Single developer  | ✅ Recommended | Works             |
| Team sharing      | ❌ Manual sync | ✅ Recommended    |
| 100+ contexts     | Slower         | ✅ 10-400x faster |
| Backup/restore    | File copy      | pg_dump           |
| Cross-machine     | ❌ Manual      | ✅ Built-in       |
| External tool API | ❌ Not available | ✅ Direct SQL access |

**Migration:** Existing file-based contexts remain in `~/.my-context/`. Database mode creates new contexts in PostgreSQL. There's no automatic migration—use `export` to move contexts if needed.

**External Integrations:** The PostgreSQL backend enables external tools to read/write my-context data directly. See [External Integrations](docs/guides/EXTERNAL-INTEGRATIONS.md) for details on integrating with tools like TimeTracker.

## Documentation

- **[Getting Started Guide](docs/guides/GETTING-STARTED.md)** — Full installation and first steps
- **[Triggers Tutorial](docs/tutorials/TRIGGERS-TUTORIAL.md)** — Automation with signals and watches
- **[Troubleshooting](docs/guides/TROUBLESHOOTING.md)** — Common issues and solutions
- **[External Integrations](docs/guides/EXTERNAL-INTEGRATIONS.md)** — Third-party tools and database integrations
- **[CLAUDE.md](CLAUDE.md)** — Architecture and development guide

## External Integrations

My-context supports external tool integrations through the PostgreSQL database backend. External tools can read/write contexts, notes, and files directly to the shared database.

**Known Integrations:**

- **[TimeTracker](docs/guides/EXTERNAL-INTEGRATIONS.md#timetracker-integration)** — C/ncurses TUI for time tracking with bidirectional sync via `tt-context-bridge` script

**Key Features:**
- Direct PostgreSQL access for performance
- Extended metadata support via `note_type` and `tags` columns
- Concurrent write handling via PostgreSQL row-level locking

**Important:** External integrations bypass my-context validation. See the [External Integrations Guide](docs/guides/EXTERNAL-INTEGRATIONS.md) for schema coordination requirements.

## Comparison

| Feature               | my-context | git stash    | tmux sessions | note apps |
| --------------------- | ---------- | ------------ | ------------- | --------- |
| Context tracking      | Yes        | No           | No            | Manual    |
| Automatic transitions | Yes        | No           | No            | No        |
| File associations     | Yes        | Yes (staged) | No            | Manual    |
| Timestamped notes     | Yes        | No           | No            | Yes       |
| Plain text storage    | Yes        | Binary       | N/A           | Varies    |
| Cross-platform        | Yes        | Yes          | Unix only     | Varies    |
| CLI-first             | Yes        | Yes          | Yes           | No        |

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Development:**

```bash
go build -o my-context ./cmd/my-context/
go test ./...
```

**Integration Testing with BATS:**

```bash
# Install BATS (https://bats-core.readthedocs.io/)
# macOS: brew install bats-core
# Debian/Ubuntu: apt install bats

# Run file mode tests
make test-bats

# Run all tests (database tests require DATABASE_URL_TEST)
DATABASE_URL_TEST="postgres://..." make test-bats-all

# Or use the runner script directly
./scripts/run-integration-tests.sh file    # File mode only
./scripts/run-integration-tests.sh db      # Database mode only
./scripts/run-integration-tests.sh all     # Both suites
```

Test suites include:

- **file_mode.bats** — Core commands with file-based storage (17 tests)
- **database_mode.bats** — Database backend commands (19 tests)
- **synthetic_workflows.bats** — Real developer workflow simulations (8 tests)

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

## Related Projects

my-context is part of the [jcaldwell-labs](https://github.com/jcaldwell-labs) organization. Sister projects you may find useful:

**Terminal/TUI Applications:**
| Project | Description | Synergy with my-context |
|---------|-------------|-------------------------|
| [my-grid](https://github.com/jcaldwell-labs/my-grid) | ASCII canvas editor with vim-style navigation | Track grid editing sessions and design decisions |
| [boxes-live](https://github.com/jcaldwell-labs/boxes-live) | Real-time ASCII box drawing with joystick support | Document box layout iterations during prototyping |
| [terminal-stars](https://github.com/jcaldwell-labs/terminal-stars) | Starfield animation for terminals | — |
| [atari-style](https://github.com/jcaldwell-labs/atari-style) | Retro visual effects and shaders for terminal apps | — |
| [smartterm-prototype](https://github.com/jcaldwell-labs/smartterm-prototype) | Smart terminal with readline-like features | Integrate my-context commands into smart workflows |

**CLI Tools:**
| Project | Description | Synergy with my-context |
|---------|-------------|-------------------------|
| [fintrack](https://github.com/jcaldwell-labs/fintrack) | Personal finance tracking CLI (Go) | Shared CLI patterns and Go architecture |
| [tario](https://github.com/jcaldwell-labs/tario) | Terminal-based platformer game (Go) | — |

**Game Engines & Agents:**
| Project | Description | Synergy with my-context |
|---------|-------------|-------------------------|
| [adventure-engine-v2](https://github.com/jcaldwell-labs/adventure-engine-v2) | Multiplayer text adventure engine (C) | Track game development sessions and world-building notes |
| [capability-catalog](https://github.com/jcaldwell-labs/capability-catalog) | Skill/capability definitions for AI agents | my-context provides session memory for agent workflows |

## License

[MIT License](LICENSE) — Free for personal and commercial use.
