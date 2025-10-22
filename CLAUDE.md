# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

My-Context-Copilot - Cross-platform CLI tool for managing developer work contexts with notes, file associations, and timestamps.

## Quick Commands

### Build and Run
```bash
# Build for current platform
go build -o my-context.exe ./cmd/my-context/

# Run tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Cross-platform build
./scripts/build.sh

# Install locally
./scripts/install.sh
```

### Common Development Tasks
```bash
# Run specific test file
go test ./tests/integration/start_test.go -v

# Check for compilation errors
go build ./...

# Run the CLI
./my-context.exe --help
./my-context.exe start "Test context"
./my-context.exe show
```

## Architecture

### Project Structure
```
cmd/my-context/          # CLI entry point (main.go)
internal/
  commands/              # Command implementations (start, stop, note, etc.)
  core/                  # Business logic (context operations, state, storage)
  models/                # Data structures (Context, Note, FileAssociation, etc.)
  output/                # Output formatters (human-readable, JSON)
tests/
  integration/           # Integration tests
  unit/                  # Unit tests
specs/                   # Feature specifications
  ###-feature-name/     # Numbered feature directories
.specify/                # Spec Kit tooling
  scripts/bash/          # Automation scripts
  templates/             # Document templates
  memory/
    constitution.md      # Project principles
```

### Key Design Patterns

**Plain Text Storage**: All data stored in `~/.my-context/` as human-readable text files (JSON for metadata, plain logs for events). Files can be viewed with standard tools (`cat`, `grep`), version controlled, and manually edited.

**Context Lifecycle**:
- Only one context active at a time
- Starting new context automatically stops previous
- Each context gets isolated directory with meta.json, notes.log, files.log, touch.log
- Global state.json tracks active context
- Global transitions.log tracks all context switches
- **NEW Sprint 2**: Contexts can be archived (is_archived flag in meta.json)
- **NEW Sprint 2**: Archived contexts hidden from default list, visible with --archived flag

**Cross-Platform Path Handling**:
- Internal storage uses POSIX format (forward slashes)
- `NormalizePath()` converts input paths to absolute POSIX
- `DenormalizePath()` converts back to OS-native for display
- Critical for Windows compatibility (backslash vs forward slash)

**Command Architecture** (Cobra framework):
- Root command in `cmd/my-context/main.go`
- Subcommands in `internal/commands/*.go`
- Each command supports single-letter aliases (e.g., `start` → `s`)
- `--json` flag available globally for machine-readable output
- Commands delegate to `internal/core/` for business logic
- **NEW Sprint 2**: export (e), archive (a), delete (d) commands
- **NEW Sprint 2**: list command supports --project, --limit, --search, --all, --archived, --active-only
- **NEW Sprint 2**: start command supports --project flag for "project: phase" naming

**Project Organization** (Sprint 2):
- Users organize contexts with "project: phase - description" naming convention
- `list --project <name>` filters by project (extracts text before first colon)
- `start "Phase 1" --project ps-cli` creates "ps-cli: Phase 1"
- Case-insensitive project matching

**Output Strategy**:
- `internal/output/human.go`: Human-readable formatting
- `internal/output/json.go`: Machine-readable JSON with status/data/error structure
- **NEW Sprint 2**: `internal/output/markdown.go`: Markdown export formatter
- All commands support both formats via `--json` flag

### Data Flow Example (Start Command)
```
User: my-context start "Bug fix"
  ↓
commands.StartCmd validates input
  ↓
core.CreateContext():
  - Sanitizes name (spaces → underscores)
  - Checks for duplicates (adds _2, _3 suffix if needed)
  - Creates ~/.my-context/Bug_fix/ directory
  - Writes meta.json with metadata
  ↓
core.SetActiveContext():
  - Stops previous context if any (updates stop time)
  - Records transition in transitions.log
  - Updates state.json with new active context
  ↓
output.PrintContext() displays result
```

## Development Workflow (Specification-Driven)

### Feature Development Process
The project follows a rigorous SDLC defined in `SDLC.md`:

1. **Specification** (`/specify`): Create `spec.md` - WHAT and WHY, no implementation details
2. **Clarification** (`/clarify`): Resolve ambiguities, document assumptions
3. **Planning** (`/plan`): Create technical design (plan.md, data-model.md, contracts/, etc.)
4. **Task Generation** (`/tasks`): Generate dependency-ordered tasks.md
5. **Implementation** (`/implement`): Execute tasks following TDD
6. **Review & Merge**: Self-review → peer review → merge

### Critical: Test-Driven Development (NON-NEGOTIABLE)
- Phase 3.2 (Write Tests) MUST happen BEFORE Phase 3.3 (Implementation)
- Tests should fail initially (red) before implementation
- Cannot proceed to implementation without failing tests
- See `SDLC.md` sections "TDD Enforcement" and "Blocking Gate"

### Spec Kit Scripts
Located in `.specify/scripts/bash/`:
- `create-new-feature.sh "<description>"` - Start new feature branch and spec
- `setup-plan.sh` - Initialize planning phase
- `check-prerequisites.sh` - Validate environment
- `update-agent-context.sh` - Update AI context

### Feature Directories
Each feature lives in `specs/###-feature-name/`:
- `spec.md` - Requirements (functional, acceptance criteria)
- `plan.md` - Technical design decisions
- `tasks.md` - Numbered implementation tasks
- `data-model.md` - Entity definitions
- `contracts/` - API/command specifications
- `quickstart.md` - Manual test scenarios
- `research.md` - Technology decisions

## Testing Strategy

### Test Organization
- Integration tests: `tests/integration/*_test.go`
- Unit tests: `tests/unit/*_test.go`
- Framework: `github.com/stretchr/testify`

### Cross-Platform Testing Requirements
**Minimum platforms** (from SDLC.md):
- Windows: git-bash (primary) + cmd.exe (fallback)
- Linux OR macOS: One of Ubuntu/macOS/WSL

**Test each platform for**:
- Binary runs
- Commands work correctly
- Paths normalized properly
- State persists across sessions
- JSON output valid

## Environment and Configuration

**Environment Variables**:
- `MY_CONTEXT_HOME`: Override default context storage directory (default: `~/.my-context/`)

**Storage Locations**:
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

## Code Conventions

### Context Names
- Spaces converted to underscores for directory names
- Duplicates get `_2`, `_3` suffix automatically
- Path separators (`/`, `\`) replaced with underscores

### Error Handling
- Use descriptive error messages
- Exit code 0 for success
- Exit code 1 for user errors (bad input, no active context)
- Exit code 2 for system errors (file I/O failures)

### File Operations
- Use atomic writes for JSON files (write to .tmp, then rename)
- Use append-only logs for events
- Always check file existence before reading
- Return empty arrays/slices for missing log files (not errors)

## Common Pitfalls

**Windows Path Issues**: Always use `NormalizePath()` before storing paths and `DenormalizePath()` before displaying. Never store Windows backslash paths directly.

**Context State**: Remember only one context can be active. Starting a new context must stop the previous one.

**Duplicate Names**: Don't error on duplicate context names - automatically append suffix (_2, _3).

**Test-First**: Never implement features without writing failing tests first. This is a hard requirement per constitution.

**Log Files**: Empty log files are valid. Return empty slices, not errors, when log files don't exist.

## Dependencies

Core libraries (see `go.mod`):
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/stretchr/testify` - Testing framework

Go version: 1.25.1

## Related Documentation

- `README.md` - User-facing documentation and usage examples
- `SDLC.md` - Complete software development lifecycle
- `IMPLEMENTATION.md` - Implementation roadmap and task breakdown
- `specs/001-cli-context-management/` - Core feature specification
- `.specify/memory/constitution.md` - Project principles (needs customization)

## Recent Changes

### Sprint 2 (2025-10-05): Installation & Usability Improvements
- **Multi-platform builds**: Static binaries for Windows, Linux, macOS (amd64 + arm64)
- **Installation scripts**: install.sh, install.bat, install.ps1 for all platforms
- **Project filtering**: `list --project <name>` and `start --project <name>` flags
- **Export command**: Generate markdown summaries with `export <context> --to <path>`
- **Archive command**: Mark completed work with `archive <context>`
- **Delete command**: Remove contexts permanently with `delete <context>`
- **List enhancements**: --limit, --search, --all, --archived, --active-only flags
- **Bug fixes**: $ character preserved in notes, NULL replaced with "(none)" in history
- **Backward compatibility**: Sprint 1 data works seamlessly with Sprint 2 binary

### Sprint 1 (2025-10-04 to 2025-10-05): Initial Release
- 8 core commands: start, stop, note, file, touch, show, list, history
- Plain-text storage in ~/.my-context/
- Cross-platform support (Windows, Linux, macOS, WSL)
- Automatic context switching
- JSON output mode

## Constitution Principles

This project follows strict design principles documented in `.specify/memory/constitution.md`:

1. **Unix Philosophy**: Composable commands, text I/O, single purpose
2. **Cross-Platform Compatibility**: Works on Windows, Linux, macOS, WSL
3. **Stateful Context Management**: One active context, automatic transitions
4. **Minimal Surface Area**: <12 commands total, single-letter aliases
5. **Data Portability**: Plain text, no lock-in, greppable
6. **User-Driven Design** (Sprint 2): Observe and formalize organic user patterns

All feature development must align with these principles.
