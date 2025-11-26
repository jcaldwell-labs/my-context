# GitHub Copilot Instructions

This file provides guidance to GitHub Copilot when working with code in this repository.

## Project Overview

My-Context is a cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps. It follows Unix philosophy: composable commands, text I/O, single purpose.

**Status**: Production (Sprint 2)
**Language**: Go 1.24+
**Dependencies**: Cobra (CLI), Viper (config), Testify (testing)

Key features:
- Context lifecycle management (start, stop, resume, show)
- Note-taking with timestamps
- File associations
- Export to markdown/JSON
- Archive/delete contexts
- Project organization with `--project` flag
- Plain-text storage in `~/.my-context/`

## Build System

```bash
# Build
make build          # Build for current platform
make build-all      # Cross-platform builds (Linux, macOS, Windows)
make install        # Install locally

# Testing
make test           # Run all tests with race detector + coverage
make test-short     # Run tests without race detector
make test-integration # Run integration tests only
make test-unit      # Run unit tests only
make benchmark      # Run benchmarks
make coverage       # Generate HTML coverage report

# Code Quality
make fmt            # Format code with gofmt/goimports
make lint           # Run golangci-lint
make vet            # Run go vet
make check          # Run all checks (lint, vet, test)

# Development
make ci             # Full CI pipeline (deps, checks, build)
make watch          # Watch mode for tests (requires entr)
make dev-setup      # Setup development environment

# Dependencies
make deps           # Download dependencies
make tidy           # Tidy go.mod/go.sum

make help           # Show all targets
```

## Architecture

### Directory Structure

```
my-context/
├── cmd/my-context/          # CLI entry point
│   └── main.go              # Root command, version info
├── internal/
│   ├── commands/            # Cobra command implementations (17+ files)
│   │   ├── start.go         # Create & activate context
│   │   ├── stop.go          # Stop active context
│   │   ├── resume.go        # Resume previous context
│   │   ├── note.go          # Add timestamped notes
│   │   ├── file.go          # Associate files
│   │   ├── export.go        # Export to markdown/JSON
│   │   ├── archive.go       # Mark as archived
│   │   ├── delete.go        # Permanently delete
│   │   ├── list.go          # List contexts with filters
│   │   ├── show.go          # Display active context
│   │   └── ...              # Other commands
│   ├── core/                # Business logic
│   │   └── [context operations, state, storage]
│   ├── models/              # Data structures
│   │   └── [Context, Note, FileAssociation]
│   └── output/              # Output formatters
│       ├── human.go         # Human-readable
│       ├── json.go          # JSON with status/data/error
│       └── markdown.go      # Markdown export
├── tests/
│   ├── integration/         # Integration tests (20+ files)
│   ├── unit/                # Unit tests
│   ├── contract/            # Contract tests
│   └── benchmarks/          # Performance benchmarks
├── scripts/                 # Build and install scripts
└── specs/                   # Feature specifications
```

### Key Design Patterns

**Plain Text Storage**:
- All data in `~/.my-context/` (or `$MY_CONTEXT_HOME`)
- JSON for metadata (`meta.json`, `state.json`)
- Plain logs for events (`notes.log`, `files.log`, `touch.log`)
- Files viewable with standard tools, version controllable

**Context Lifecycle**:
- Only one context active at a time
- Starting new context automatically stops previous
- Each context gets isolated directory
- Global `state.json` tracks active context
- Global `transitions.log` tracks all switches

**Cross-Platform Path Handling**:
- Internal storage uses POSIX format (forward slashes)
- `NormalizePath()` converts input to absolute POSIX
- `DenormalizePath()` converts back to OS-native for display

**Project Organization**:
- Convention: `project: phase - description` naming
- `list --project <name>` filters by project
- `start "Phase 1" --project ps-cli` creates "ps-cli: Phase 1"

## Code Style and Conventions

- **Go conventions**: gofmt, effective go guidelines
- **Testing**: Table-driven tests, testify assertions
- **TDD**: Tests MUST be written BEFORE implementation (non-negotiable)
- **Error handling**: Exit 0 success, exit 1 user error, exit 2 system error

**Context Names**:
- Spaces → underscores for directory names
- Duplicates get `_2`, `_3` suffix automatically
- Path separators (`/`, `\`) replaced with underscores

**File Operations**:
- Atomic writes (write to .tmp, then rename)
- Append-only logs for events
- Empty log files return empty slices (not errors)

## Before Committing (Required Steps)

Run these commands before every commit:

1. **Format**: `make fmt` - Format code
2. **Lint**: `make lint` - Check for issues
3. **Vet**: `make vet` - Run go vet
4. **Test**: `make test` - Run tests with race detector

```bash
# Quick pre-commit check (all in one)
make check
```

## Common Development Tasks

### Adding New Commands
1. Write failing tests first (TDD is mandatory)
2. Create command file in `internal/commands/`
3. Define Cobra command with single-letter alias
4. Implement business logic in `internal/core/`
5. Support `--json` flag for machine-readable output
6. Add to root command
7. Run tests: `make test`

### Specification-Driven Development
Feature development follows SDLC.md:
1. **Specification** (`/specify`): Create `spec.md`
2. **Clarification** (`/clarify`): Resolve ambiguities
3. **Planning** (`/plan`): Create technical design
4. **Task Generation** (`/tasks`): Generate tasks.md
5. **Implementation** (`/implement`): Execute with TDD
6. **Review & Merge**

### Cross-Platform Testing
Test on minimum platforms:
- Windows: git-bash (primary) + cmd.exe (fallback)
- Linux OR macOS: One of Ubuntu/macOS/WSL

## Pull Request Standards

When creating PRs, follow these rules:

1. **Always link the issue**: Use `Fixes #N` or `Closes #N`
2. **Fill in all sections**: Never leave placeholder text

**Required PR format:**
```markdown
## Summary
[2-3 sentences describing what and why]

Fixes #[issue-number]

## Changes
- [Actual change 1]
- [Actual change 2]

## Testing
- [x] Tests written BEFORE implementation (TDD)
- [x] All tests pass (`make test`)
- [x] Cross-platform tested

## Type
- [x] New feature | Bug fix | Refactor | Docs | CI
```

## Storage Structure

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

## Common Pitfalls

**Windows Path Issues**: Always use `NormalizePath()` before storing, `DenormalizePath()` before displaying. Never store backslash paths.

**Context State**: Only one context active. Starting new context stops previous.

**Duplicate Names**: Don't error - automatically append suffix (_2, _3).

**Test-First**: Never implement without failing tests first (hard requirement).

**Log Files**: Empty log files are valid - return empty slices, not errors.

## Additional Documentation

- **README.md** - User documentation and examples
- **SDLC.md** - Complete development lifecycle
- **IMPLEMENTATION.md** - Roadmap and task breakdown
- **specs/** - Feature specifications
- **.specify/memory/constitution.md** - Project principles
