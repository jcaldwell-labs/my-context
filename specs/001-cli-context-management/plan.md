# Implementation Plan: CLI Context Management System

**Branch**: `001-cli-context-management` | **Date**: 2025-10-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-cli-context-management/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → ✅ Loaded: spec.md with clarifications present
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → ✅ All clarifications resolved in spec
3. Fill the Constitution Check section
   → ✅ Completed - all principles pass
4. Evaluate Constitution Check section
   → ✅ No violations - proceed
5. Execute Phase 0 → research.md
   → ✅ Completed
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent file
   → ✅ Completed
7. Re-evaluate Constitution Check
   → ✅ Design compliant - no violations
8. Plan Phase 2 → Task generation approach described
   → ✅ Completed
9. STOP - Ready for /tasks command
```

## Summary

Build a cross-platform CLI tool for managing developer work contexts with commands to start/stop contexts, add notes/timestamps/files, and view context history. Each context gets an isolated subdirectory in a home folder. Automatic duplicate name handling with _2, _3 suffixes. JSON output support for all commands. Automatic transition logging for audit trails.

**Key Technical Approach**: Single-binary Go application with shell wrappers for environment detection, plain-text file storage in subdirectories, and Unix-style composable commands.

## Technical Context

**Language/Version**: Go 1.21+ (single-binary compilation, excellent cross-platform support, no runtime dependencies)  
**Primary Dependencies**: 
  - `cobra` (CLI framework with subcommands and flags)
  - `viper` (configuration management for home directory path)
  - Standard library for file I/O and JSON serialization
**Storage**: Plain text files in subdirectories (notes.log, files.log, meta.json, touch.log) + central transitions.log  
**Testing**: Go testing package with table-driven tests, `testify/assert` for assertions  
**Target Platform**: Windows (cmd.exe, PowerShell, git-bash, WSL), Linux, macOS  
**Project Type**: Single project (CLI tool)  
**Performance Goals**: <10ms command execution (most operations are simple file I/O), <1MB binary size  
**Constraints**: Zero external runtime dependencies, works offline, no network calls, cross-shell compatibility  
**Scale/Scope**: ~10 commands, single-user local tool, handles 1000s of contexts, minimal memory footprint

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Unix Philosophy Compliance**
- [x] Each component does one thing well - Each command has single responsibility (start, stop, note, etc.)
- [x] Commands use text-based I/O (stdin/args → stdout, errors → stderr) - All output to stdout, errors to stderr
- [x] Components are composable and chainable with standard shell tools - Plain text output enables grep/awk/pipe usage
- [x] No unnecessary coupling between commands - Each command independently functional via context state file

**II. Cross-Platform Compatibility**
- [x] Design works across cmd.exe, PowerShell, git-bash, and WSL - Go compiles to single native binary
- [x] Path handling normalizes Windows and POSIX conventions - filepath.Clean + ToSlash/FromSlash conversions
- [x] No platform-specific dependencies without abstraction layer - All I/O through Go standard library abstractions

**III. Stateful Context Management**
- [x] Feature respects single active context model - State file stores single active context reference
- [x] Context state persists in home directory (not dev folder) - Dedicated MY_CONTEXT_HOME directory
- [x] Operations default to current context appropriately - All commands read active context from state file

**IV. Minimal Surface Area**
- [x] No new commands unless absolutely necessary - 8 essential commands: start, stop, note, file, touch, show, list, history
- [x] Configuration avoided - sensible defaults used - Single env var (MY_CONTEXT_HOME) with default to ~/.my-context
- [x] Help documentation included for any user-facing changes - Cobra auto-generates help for all commands

**V. Data Portability**
- [x] All data stored as plain text (JSON for structured, newline-delimited for logs) - meta.json + .log files
- [x] No proprietary formats introduced - Standard JSON + newline-delimited text
- [x] Data remains readable/editable with standard tools (grep, cat, text editors) - All files are plain text

**Violations & Justifications**: None. Design fully compliant with all five principles.

## Project Structure

### Documentation (this feature)
```
specs/001-cli-context-management/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Technology decisions
├── data-model.md        # Data structures and storage format
├── quickstart.md        # Manual testing scenarios
└── contracts/           # Command I/O contracts
    ├── start.md
    ├── stop.md
    ├── note.md
    ├── file.md
    ├── touch.md
    ├── show.md
    ├── list.md
    └── history.md
```

### Source Code (repository root)
```
my-context-copilot/
├── cmd/
│   └── my-context/
│       └── main.go           # Entry point
├── internal/
│   ├── commands/             # Cobra command implementations
│   │   ├── start.go
│   │   ├── stop.go
│   │   ├── note.go
│   │   ├── file.go
│   │   ├── touch.go
│   │   ├── show.go
│   │   ├── list.go
│   │   └── history.go
│   ├── core/                 # Business logic
│   │   ├── context.go        # Context operations
│   │   ├── state.go          # State management
│   │   └── storage.go        # File I/O operations
│   ├── models/               # Data structures
│   │   ├── context.go
│   │   ├── note.go
│   │   ├── file_association.go
│   │   ├── touch_event.go
│   │   └── transition.go
│   └── output/               # Output formatting
│       ├── human.go          # Human-readable output
│       └── json.go           # JSON output
├── tests/
│   ├── integration/          # Cross-platform integration tests
│   │   ├── commands_test.go
│   │   ├── paths_test.go
│   │   └── shells_test.go
│   └── unit/                 # Unit tests for core logic
│       ├── context_test.go
│       ├── state_test.go
│       └── storage_test.go
├── scripts/
│   ├── build.sh              # Build script for all platforms
│   └── install.sh            # Installation script
├── go.mod
├── go.sum
└── README.md
```

**Structure Decision**: Single project structure (Option 1) selected. This is a standalone CLI tool with no frontend/backend separation or mobile components. Go's internal/ directory enforces encapsulation while tests/ separates integration and unit tests.

## Phase 0: Research & Technology Decisions

### Decision Log

#### 1. Programming Language: Go
**Decision**: Use Go 1.21+ for implementation  
**Rationale**:
- Single-binary compilation with zero runtime dependencies (critical for cross-platform distribution)
- Excellent cross-platform support (Windows, Linux, macOS) with native path handling
- Strong standard library for file I/O, JSON, and time operations
- Fast compilation and execution (<10ms command response time achievable)
- Built-in testing framework with parallel test support
- Static typing catches errors at compile time

**Alternatives Considered**:
- Python: Requires Python runtime installation, slower startup time (~50-100ms), harder cross-platform distribution
- Rust: Steeper learning curve, longer compile times, overkill for simple file I/O operations
- Shell script: Not truly cross-platform (bash vs cmd differences), harder to test, no structured data handling

#### 2. CLI Framework: Cobra
**Decision**: Use spf13/cobra for command structure  
**Rationale**:
- Industry standard for Go CLIs (used by kubectl, hugo, gh)
- Built-in subcommand support matching our 8 commands
- Automatic help generation and flag parsing
- Single-letter alias support via `Aliases` field
- Persistent flags for --json across all commands

**Alternatives Considered**:
- urfave/cli: Less feature-rich, less idiomatic for complex subcommand trees
- Standard flag package: Would require manual subcommand routing and help text

#### 3. Storage Format: Plain Text Files
**Decision**: Use subdirectory-per-context with structured files  
**Rationale**:
- Constitutional requirement: plain text for portability
- Easy to backup, version control, and inspect manually
- No database dependencies or corruption risk
- grep/awk compatible for power users
- Structure: `~/.my-context/<context-name>/meta.json`, `notes.log`, `files.log`, `touch.log`

**File Format Decisions**:
- **meta.json**: Context metadata (name, start_time, end_time, status) - structured for atomic reads
- **notes.log**: Newline-delimited `timestamp|note_text` - append-only for performance
- **files.log**: Newline-delimited `timestamp|normalized_path` - append-only
- **touch.log**: Newline-delimited `timestamp` - append-only
- **transitions.log**: Central file at `~/.my-context/transitions.log` with `timestamp|prev|new|type` format

**Alternatives Considered**:
- SQLite: Too heavyweight, proprietary binary format violates portability principle
- Single JSON file per context: Poor performance for append operations (notes/files/touch)
- Git branches: Overcomplicated, git dependency not needed, harder cross-platform

#### 4. State Management: Single State File
**Decision**: Use `~/.my-context/state.json` to track active context  
**Rationale**:
- Atomic read/write of small JSON file (<1KB)
- Contains: `{"active_context": "context-name-or-null", "last_updated": "timestamp"}`
- Enables instant "which context am I in?" lookups
- Cross-shell visibility (all shells read same file)

#### 5. Duplicate Name Resolution Algorithm
**Decision**: Suffix sequence `_2`, `_3`, etc. starting at first duplicate  
**Implementation**:
```
1. Check if directory exists for given name
2. If not exists: use original name
3. If exists: iterate _2, _3, _4... until finding non-existent directory
4. Use first available name
```
**Performance**: O(n) where n = number of duplicates (typically <5)

#### 6. Path Normalization Strategy
**Decision**: Store paths in POSIX format (forward slash) internally, convert on output  
**Rationale**:
- `filepath.ToSlash()` converts platform paths to POSIX for storage
- `filepath.FromSlash()` converts back for display
- `filepath.Clean()` normalizes before conversion
- Consistent storage format aids cross-platform context sharing

#### 7. JSON Output Format
**Decision**: Use `--json` flag on all commands, output complete structured data  
**Format**:
```json
{
  "command": "show",
  "timestamp": "2025-10-04T15:30:00Z",
  "data": {
    "context": {
      "name": "Bug fix",
      "start_time": "2025-10-04T14:00:00Z",
      "status": "active",
      "notes": [...],
      "files": [...],
      "touches": [...]
    }
  }
}
```

#### 8. Error Handling Strategy
**Decision**: Exit codes + stderr messages  
- 0: Success
- 1: User error (invalid arguments, no active context)
- 2: System error (I/O failure, permission denied)
- All errors to stderr, never to stdout (preserves piping)

## Phase 1: Design & Contracts

### Data Model

See [data-model.md](./data-model.md) for complete entity definitions, relationships, and storage schema.

**Key Entities**:
- Context (name, start_time, end_time, status, subdirectory_path)
- Note (timestamp, text_content)
- FileAssociation (timestamp, file_path)
- TouchEvent (timestamp)
- ContextTransition (timestamp, previous_context, new_context, transition_type)
- AppState (active_context_name, last_updated)

### Command Contracts

See [contracts/](./contracts/) directory for detailed I/O specifications:
- `start.md`: Context creation and activation
- `stop.md`: Context deactivation
- `note.md`: Note addition
- `file.md`: File association
- `touch.md`: Timestamp recording
- `show.md`: Current context display
- `list.md`: All contexts listing
- `history.md`: Transition log display

### Quickstart Testing

See [quickstart.md](./quickstart.md) for manual testing scenarios covering:
1. Basic workflow (start → note → file → touch → show → stop)
2. Context switching (automatic stop on new start)
3. Duplicate name handling
4. Cross-shell persistence
5. JSON output validation
6. Error conditions

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
1. **Setup Tasks** (3 tasks):
   - Initialize Go module and directory structure
   - Install Cobra dependency
   - Create initial main.go with root command

2. **Model Tasks** (5 tasks, parallel):
   - One task per model file (context.go, note.go, file_association.go, touch_event.go, transition.go)
   - Define structs with JSON tags and validation methods

3. **Core Logic Tasks** (3 tasks, sequential dependencies):
   - storage.go: File I/O operations (read/write/append)
   - state.go: State file management (depends on storage.go)
   - context.go: Business logic (depends on storage.go + state.go)

4. **Command Tasks** (8 tasks, parallel):
   - One task per command file implementing Cobra command interface
   - Each depends on core logic being complete

5. **Output Formatting Tasks** (2 tasks, parallel):
   - human.go: Human-readable formatters
   - json.go: JSON output formatters

6. **Integration Test Tasks** (before implementation):
   - Cross-platform path tests
   - Shell wrapper tests (cmd/bash/WSL)
   - State persistence tests
   - Duplicate name tests

7. **Build & Install Tasks**:
   - Build script for Windows/Linux/macOS
   - Installation script
   - Shell wrapper creation

**Ordering Strategy**:
- Models first (no dependencies)
- Core logic second (depends on models)
- Commands third (depends on core + models)
- Tests throughout (TDD approach)
- Build/install last

**Estimated Task Count**: ~35-40 tasks total

## Complexity Tracking

No violations of constitutional principles. Design is intentionally simple:
- Single-binary Go application (no microservices, no database)
- Plain text file storage (no complex data structures)
- Standard library-heavy (minimal external dependencies)
- Direct file I/O (no ORMs, no caching layers)

## Progress Tracking

**Phase Status**:
- [x] Phase 0: Research complete
- [x] Phase 1: Design complete
- [x] Phase 2: Task planning approach described (tasks.md NOT created)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS (all 5 principles satisfied)
- [x] Post-Design Constitution Check: PASS (no violations introduced)
- [x] All NEEDS CLARIFICATION resolved (spec has clarifications section)
- [x] No complexity deviations (simple single-binary design)

---

**Next Steps**: Run `/tasks 001-cli-context-management` to generate tasks.md

*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
