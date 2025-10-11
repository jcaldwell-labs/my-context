# Implementation Plan: Context Signaling Protocol

**Branch**: `005-context-signaling-protocol` | **Date**: 2025-10-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/005-context-signaling-protocol/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement a lightweight signaling protocol for my-context to enable event-driven coordination between development team members and processes. The solution uses file-based semaphores (signal files) for event coordination and polling-based monitoring for context changes, providing automated notifications without requiring daemons or network dependencies. Target v2.2.0 with three main features: signal files for event coordination, watch command for monitoring, and enhanced context metadata.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.21+
**Primary Dependencies**: Cobra (existing CLI framework), fsnotify (optional for inotify support)
**Storage**: JSON files in ~/.my-context/ (existing), new signals/ subdirectory
**Testing**: Go's built-in testing + table-driven tests (existing pattern)
**Target Platform**: Linux, macOS, Windows (cross-platform CLI)
**Project Type**: Single CLI application
**Performance Goals**: <1% CPU usage for watch commands, <5s detection latency
**Constraints**: File-based only (no daemons/network), backward compatible, offline-capable
**Scale/Scope**: ~500 LOC addition, 3 new CLI commands, cross-platform support

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Core Principles Compliance

**✓ I. Unix Philosophy**: Signal and watch commands follow composability principles:
- `signal` commands produce text output suitable for shell scripting
- `watch` can be chained with other commands via `--exec`
- File-based approach enables integration with standard Unix tools

**✓ II. Cross-Platform Compatibility**: Implementation uses polling-based approach initially (works on all platforms), with optional inotify optimization for Linux.

**✓ III. Stateful Context Management**: Watch commands integrate with existing active context semantics, signal files are context-agnostic but can reference contexts.

**⚠ IV. Minimal Surface Area**: Adding 2 new commands (`signal`, `watch`) would exceed 12-command limit (11→13). This requires constitution amendment or command consolidation.

**✓ V. Data Portability**: Signal files are plain text timestamps, metadata enhancements use existing JSON structure.

**✓ VI. User-Driven Design**: Features directly address observed gaps in team handoff workflow (manual notifications, polling overhead).

### Required Constitution Amendment

The signaling protocol requires adding `signal` and `watch` commands, exceeding the 12-command limit. **Amendment Proposal**: Increase command limit from 12 to 14 to accommodate event-driven coordination features that enhance team workflows.

*Alternative Considered*: Combine `signal` and `watch` into a single `monitor` command with subcommands, but this reduces discoverability and violates "one thing well" principle.

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
cmd/
├── signal.go          # Signal file management commands (create, list, wait, clear)
└── watch.go           # Context monitoring commands (watch with polling/inotify)

internal/
├── signal/            # Signal file operations
│   ├── manager.go     # Core signal file CRUD operations
│   └── manager_test.go
└── watch/             # Context monitoring
    ├── monitor.go     # Watch logic with polling and optional inotify
    ├── monitor_test.go
    └── patterns.go    # Pattern matching for note filtering

pkg/
├── models/            # Enhanced context metadata models
│   ├── context.go     # Extended Context struct with metadata fields
│   └── signal.go      # Signal file data structures
└── utils/
    └── fs.go          # Cross-platform filesystem utilities
```

**Structure Decision**: Single CLI application following existing Go project layout. New commands added to `cmd/`, business logic in `internal/`, shared models in `pkg/`. Maintains separation between CLI commands and core logic.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Exceeds 12-command limit (11→13) | Event-driven coordination requires distinct `signal` and `watch` commands for Unix philosophy compliance | Combining into single command reduces discoverability and violates "one thing well" principle. Manual notifications are insufficient for team coordination needs. |
