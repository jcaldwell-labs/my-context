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

## Database Backend and Partitioning (v3.0.0)

### Overview

Version 3.0.0 introduces PostgreSQL backend support with partition capabilities, providing 10-400x performance improvement over file-based storage.

### Backend Selection

Set via `MY_CONTEXT_HOME` environment variable:

**File-Based (Default)**:

```bash
export MY_CONTEXT_HOME=~/.my-context  # Or any directory path
```

**Database (Single Partition)**:

```bash
export MY_CONTEXT_HOME=db              # Default PostgreSQL (public schema)
export MY_CONTEXT_HOME=database        # Same as 'db'
export MY_CONTEXT_HOME=pg              # Same as 'db'
```

**Database (Multiple Partitions)** - NEW:

```bash
export MY_CONTEXT_HOME=db:adventure-engine    # Partition for adventure-engine project
export MY_CONTEXT_HOME=db:payment-service     # Partition for payment-service project
export MY_CONTEXT_HOME=db:scrum               # Partition for scrum contexts
```

### Partition Architecture

**Design Philosophy**:

- Each partition = isolated workspace with its own contexts, notes, files, and state
- Partition names are auto-sanitized to valid PostgreSQL schema names
- Each partition uses a separate PostgreSQL schema in the same database
- Backward compatible: `db` alone uses `public` schema

**Partition Naming Rules**:

- Input: Any string (e.g., "Adventure Engine!", "payment-service-v2")
- Auto-sanitized to PostgreSQL schema names:
  - Lowercase conversion
  - Hyphens/spaces → underscores
  - Special chars removed
  - Prefix `p_` if starts with number
- Examples:
  - `adventure-engine` → `adventure_engine`
  - `My Project!` → `my_project`
  - `123test` → `p_123test`

**Per-Partition Isolation**:

- Independent active context (each partition tracks its own active context)
- Separate context counts
- Isolated state management
- Cross-partition queries available via flags

### Partition-Aware Commands

**View Current Partition**:

```bash
my-context which
# Output shows:
#   Backend: PostgreSQL database
#   Partition: adventure-engine
#   Schema: adventure_engine
```

**List All Partitions**:

```bash
my-context partitions               # Show all partitions with stats
my-context partitions --json        # JSON output
my-context p                        # Alias
```

**Cross-Partition Queries**:

```bash
my-context list --all-partitions            # List contexts across all partitions
my-context list --partition=scrum           # List contexts from specific partition
```

### Code Architecture

**Backend Detection** (`internal/core/backend.go`):

- `DetectBackendType()` - Recognizes `db:partition` syntax
- `ExtractPartition()` - Extracts partition name from env var
- `SanitizePartitionName()` - Converts to valid PostgreSQL schema name
- `GetPartitionSchema()` - Returns current schema name
- `GetPostgresConnectionString()` - Sets `search_path` parameter

**PostgreSQL Backend** (`pkg/storage/postgres/postgres_backend.go`):

- `schema` field - Tracks active schema (partition)
- `partition` field - Friendly partition name for display
- `createSchema()` - Creates partition schema dynamically (CREATE SCHEMA IF NOT EXISTS)
- `GetDB()` - Exposes database connection for cross-partition queries
- Cross-partition methods:
  - `ListAllPartitions()` - Returns all my-context schemas
  - `ListContextsAcrossPartitions()` - Queries all partitions
  - `SearchContextsAcrossPartitions()` - Full-text search across partitions

**Storage Interface** (`pkg/storage/storage_interface.go`):

- Core interface remains unchanged
- Cross-partition methods are PostgreSQL-specific (not in interface)
- Allows graceful feature detection via type assertion

**Commands**:

- `internal/commands/partitions.go` - New command for partition management
- `internal/commands/which.go` - Updated to show partition info
- `internal/commands/list.go` - Extended with `--all-partitions`, `--partition` flags

### Database Schema Structure

```
Database: dev_state
├── Schema: public (MY_CONTEXT_HOME=db)
│   ├── contexts
│   ├── context_notes
│   ├── context_files
│   └── state
├── Schema: adventure_engine (MY_CONTEXT_HOME=db:adventure-engine)
│   ├── contexts
│   ├── context_notes
│   ├── context_files
│   └── state
├── Schema: payment_service (MY_CONTEXT_HOME=db:payment-service)
│   ├── contexts
│   ├── context_notes
│   ├── context_files
│   └── state
└── Schema: scrum (MY_CONTEXT_HOME=db:scrum)
    ├── contexts
    ├── context_notes
    ├── context_files
    └── state
```

### Testing

**Unit Tests** (`tests/unit/partition_test.go`):

- `TestExtractPartition()` - Partition name extraction
- `TestSanitizePartitionName()` - Name sanitization rules
- `TestGetPartitionSchema()` - Schema name resolution
- `TestDetectBackendType()` - Backend type detection
- `TestGetPostgresConnectionString()` - Connection string generation

**Run Tests**:

```bash
go test ./tests/unit/partition_test.go -v
```

### Migration from Single Partition

**Existing Data**: Contexts created with `MY_CONTEXT_HOME=db` remain in `public` schema and continue to work.

**Organizing into Partitions**:

1. Export contexts from public schema: `my-context export <context> --to <file>.md`
2. Switch partition: `export MY_CONTEXT_HOME=db:new-partition`
3. Re-create contexts in new partition
4. (Optional) Delete from public schema if no longer needed

### Common Use Cases

**Multi-Project Developer**:

```bash
# Work on adventure-engine
export MY_CONTEXT_HOME=db:adventure-engine
my-context start "Feature X"

# Switch to payment-service
export MY_CONTEXT_HOME=db:payment-service
my-context start "Bug fix Y"

# View all partitions
my-context partitions
```

**Scrum Master**:

```bash
# Separate partition for scrum ceremonies
export MY_CONTEXT_HOME=db:scrum
my-context start "Sprint 25 Planning"
```

**Search Across All Projects**:

```bash
export MY_CONTEXT_HOME=db  # Any partition
my-context list --all-partitions --search "payment"
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

**Duplicate Names**: Don't error on duplicate context names - automatically append suffix (\_2, \_3).

**Test-First**: Never implement features without writing failing tests first. This is a hard requirement per constitution.

**Log Files**: Empty log files are valid. Return empty slices, not errors, when log files don't exist.

## Security Review Checklist (Database Code)

Before creating PRs that involve database operations, complete this checklist:

### SQL Injection Prevention

- [ ] **No dynamic SQL with user input**: Use parameterized queries (`$1`, `$2`) for values
- [ ] **Validate identifiers**: Use `utils.IsValidPostgresIdentifier()` before using schema/table names in `fmt.Sprintf`
- [ ] **Grep check**: Run `grep -n "fmt.Sprintf.*SELECT\|fmt.Sprintf.*INSERT\|fmt.Sprintf.*UPDATE\|fmt.Sprintf.*DELETE\|fmt.Sprintf.*CREATE" pkg/ internal/` to find dynamic SQL

### Resource Management

- [ ] **Error path cleanup**: If `Init()` or similar allocates resources, ensure error paths clean them up
- [ ] **Check for leaks**: Every `sql.Open()` should have corresponding `Close()` in error and success paths
- [ ] **Connection pools**: Set `SetMaxOpenConns`, `SetMaxIdleConns` for production use

### Credentials & Secrets

- [ ] **No hardcoded production credentials**: Dev-only defaults must be documented with warnings
- [ ] **Environment variable priority**: External config (DATABASE_URL) should override defaults
- [ ] **Grep check**: Run `grep -rn "password\|secret\|token\|key" --include="*.go" | grep -v "_test.go"` to audit

### Files & Artifacts

- [ ] **No backup files**: Check `git status` for `.bak`, `.tmp`, `.orig` files before commit
- [ ] **No test artifacts**: Ensure test databases, temp files aren't committed

### Documentation

- [ ] **Interface godocs**: Public interfaces must have contract documentation
- [ ] **Security notes**: Document any security-sensitive code paths

### Testing

- [ ] **Unit tests for new features**: Especially environment variable handling
- [ ] **Negative tests**: Test invalid inputs, SQL injection attempts
- [ ] **Run tests**: `go test ./tests/unit/... -v` before creating PR

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

## Code Quality Checklist (Lessons from PR #6)

PR #6 (Externalize DATABASE_URL) received 40+ review comments. To prevent similar feedback volume, apply this checklist before submitting PRs:

### Security (CRITICAL)

- [ ] **No hardcoded credentials** - Require env vars, never embed passwords in code
- [ ] **SQL injection prevention** - Use `pq.QuoteIdentifier()` for dynamic identifiers, parameterized queries for values
- [ ] **Input validation** - Validate ALL user input before use (use `utils.IsValidPostgresIdentifier()` for schema names)
- [ ] **Defense in depth** - Sanitize AND validate (e.g., `SanitizePartitionName()` + `IsValidPostgresIdentifier()`)

### Resource Management

- [ ] **Clean up on error paths** - Close DB connections, file handles on ALL error returns
- [ ] **Connection pool config** - Set `SetMaxOpenConns()`, `SetMaxIdleConns()`, `SetConnMaxLifetime()`
- [ ] **Context timeouts** - Use `context.WithTimeout()` for DB operations

### Error Handling

- [ ] **No silent failures** - Never use `_ = potentialError()`, always handle or log
- [ ] **Descriptive errors** - Include context: what failed, what was expected, how to fix
- [ ] **TODO stubs must error** - Unimplemented interface methods must return explicit errors, not nil

### Testing

- [ ] **Use `t.Setenv()` not `os.Setenv()`** - Prevents race conditions in parallel tests
- [ ] **Test actual implementations** - Don't copy functions locally, import from the package
- [ ] **Test the feature you're adding** - If PR adds DATABASE_URL support, test DATABASE_URL behavior
- [ ] **No test artifacts in repo** - Don't commit generated test files (add to .gitignore)

### Code Quality

- [ ] **Comments must be accurate** - Double-check examples in comments match actual behavior
- [ ] **URL encoding** - Use `url.QueryEscape()` when building URL query parameters
- [ ] **Document breaking changes** - If bumping major version, explain what breaks and migration path
- [ ] **Interface documentation** - Add godoc explaining contract, error patterns, thread-safety

### Dependencies

- [ ] **Explain dependency changes** - If go.mod changes unexpectedly, verify it's intentional
- [ ] **Run `go mod tidy`** - Ensure dependencies are in sync before committing

### Pre-PR Validation

```bash
# Run locally before pushing
go build ./...                    # Must compile
go test ./...                     # All tests pass
golangci-lint run --timeout=5m    # No lint errors
go mod tidy && git diff go.mod   # No unexpected changes
```

### Common Patterns in This Codebase

**Safe SQL with dynamic schema:**

```go
if !utils.IsValidPostgresIdentifier(schema) {
    return fmt.Errorf("invalid schema name %q", schema)
}
query := fmt.Sprintf("SELECT * FROM %s.table", pq.QuoteIdentifier(schema))
```

**Resource cleanup on error:**

```go
db, err := sql.Open("postgres", connStr)
if err != nil {
    return err
}

if err := db.Ping(); err != nil {
    db.Close()  // Clean up before returning
    return err
}
```

**Environment variables in tests:**

```go
t.Run("test name", func(t *testing.T) {
    t.Setenv("MY_VAR", "value")  // Auto-cleaned after test
    // ... test code
})
```
