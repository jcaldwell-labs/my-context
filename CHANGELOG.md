# Changelog

All notable changes to my-context will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.2.0] - 2026-01-13

### Added

- **Stats Command** (`my-context stats`): Time tracking aggregation with filters
  - `--today`, `--week`, `--month` for quick date ranges
  - `--since`, `--until` for custom date ranges
  - `--project` to filter by project prefix
  - Shows total duration, context counts, averages, breakdown by project

- **Record Command** (`my-context record`): Clipboard monitoring mode
  - Auto-capture clipboard changes as notes to active context
  - `--tag` flag to prefix captured notes with a tag
  - Runs until Ctrl+C, displays capture count summary

- **MCP Server** (`my-context mcp-server`): Model Context Protocol integration
  - 10 tools for AI agent integration: start_context, stop_context, add_note, get_active_context, list_contexts, add_file, list_files, export_context, archive_context, search_contexts
  - JSON-RPC over stdio transport
  - Compatible with Claude Desktop, VS Code Copilot, and other MCP clients

- **Claude Code Plugin** (`claude-plugin/`): Native Claude Code integration
  - 6 slash commands: `/my-context:start`, `/my-context:stop`, `/my-context:note`, `/my-context:show`, `/my-context:list`, `/my-context:export`
  - Direct bash calls optimized for Claude Code's environment

- **BATS Integration Testing**: Comprehensive CLI test suite
  - 44 automated tests across file mode, database mode, and synthetic workflows
  - `make test-bats` for file mode, `make test-bats-all` for all tests

### Fixed

- **Database Mode Support**: `tag`, `tree`, `history`, `watch` commands now work with PostgreSQL backend
- Lint warnings cleaned up (cyclomatic complexity, unused code, spelling)

## [3.1.3] - 2026-01-13

### Fixed

- `list` command now works correctly in database mode (`MY_CONTEXT_HOME=db`)
- Touch fields properly included in database list operations

## [3.1.2] - 2026-01-12

### Fixed

- Touch count displays correctly in database mode (was showing 0)
- `scanContextRow` includes touch fields for list operations

## [3.1.1] - 2026-01-12

### Fixed

- Minor bug fixes for database backend

## [3.1.0] - 2025-11-18

### Added

- **Context Tags/Labels System**
  - `tag add/remove/list` commands for organizing contexts
  - `list --tag <tag>` filter
  - `start --labels <tags>` to tag on creation

- **Parent-Child Context Relationships**
  - `link <child> <parent>` and `unlink` commands
  - `up` (parent) and `down` (children) navigation
  - `start --parent <parent>` flag

- **Context Tree Visualization**
  - `tree [context]` command with ASCII rendering
  - Circular dependency detection
  - JSON tree output support

### Changed

- Enhanced metadata system with labels, parent, created_by fields
- 100% backward compatible with v2.x contexts

## [3.0.0] - 2025-11-01

### Added

- **PostgreSQL Database Backend**: 10-400x performance improvement
  - `MY_CONTEXT_HOME=db` for single database mode
  - `MY_CONTEXT_HOME=db:partition` for multi-partition support
  - `partitions` command to list all partitions
  - `which` command shows backend and partition info
  - `list --all-partitions` for cross-partition queries

- **Externalized DATABASE_URL**: Configure database connection via environment variable
- **Security hardening**: Command injection prevention in watch --exec

### Changed

- Database backend is optional; file-based storage remains default
- Connection pool configuration for production use

## [2.3.0] - 2025-10-22

### Added

- **Context Home Visibility (MCF-001)**: All commands now display which context home is active
  - `show`, `list`, `history` commands display "Context Home: <path>" header
  - `start` command shows context home before creating context
  - New `which` command to check active context home location
  - `which --short` flag for path-only output (scripting)
  - `which --json` flag for machine-readable output
  - Helps troubleshoot "I don't see my context" issues
  - Makes multi-context-home workflows transparent

### Changed

- `start` command output now shows "✓ Started:" instead of "Started context:"
- Context home path abbreviated with `~` for brevity (e.g., `~/.my-context/` instead of `/home/user/.my-context/`)

### Fixed

- Context home confusion - users now always know which MY_CONTEXT_HOME they're operating on

## [2.2.0] - 2025-10-21

### Added

- Context Signaling Protocol v2.2.0
- Lifecycle improvements
- Watch command with `--exec` flag for executing commands on context changes
- Watch command `--new-notes` and `--pattern` flags

### Fixed

- Watch command now properly monitors notes.log file instead of directory
- Watch `--exec` command properly handles shell syntax (quotes, pipes, variables)
- Watch command now displays command output to user

## [2.0.0] - 2025-10-11

### Added

- Initial public release
- Core commands: start, stop, resume, note, file, touch, show, list, history
- Export command for markdown and JSON output
- Archive command for completed contexts
- Delete command for permanent removal
- Project filtering with `--project` flag
- Labels support for context categorization
- Cross-platform support (Windows, Linux, macOS, WSL)
- Plain-text storage in ~/.my-context/

### Changed

- Improved cross-platform path handling
- Better Windows compatibility

## [1.0.0] - 2025-10-05

### Added

- Initial release
- Basic context management (start, stop, note, file, show, list)
- JSON output support
- Context lifecycle management
