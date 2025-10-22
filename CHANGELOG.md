# Changelog

All notable changes to my-context will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
