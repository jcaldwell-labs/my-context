# Feature Specification: CLI Context Management System

**Feature Branch**: `001-cli-context-management`  
**Created**: 2025-10-04  
**Status**: Draft  
**Input**: User description: "CLI context management with start/stop/note/timestamp/show commands and git-based storage"

## Clarifications

### Session 2025-10-04
- Q: Duplicate context name handling strategy → A: Allow duplicates with sequence suffix (_2, _3, etc.) starting with first duplicate
- Q: Transition log access method → A: Provide both dedicated `history` command and JSON output support
- Q: JSON output mode support → A: Yes, support --json flag for machine-readable output
- Q: Git branch vs directory structure for context isolation → A: Use subdirectories (not git branches); auto-create new subdirectory per context

## Execution Flow (main)
```
1. Parse user description from Input
   → If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   → Identify: actors, actions, data, constraints
3. For each unclear aspect:
   → Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   → If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   → Each requirement must be testable
   → Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   → If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   → If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

### Section Requirements
- **Mandatory sections**: Must be completed for every feature
- **Optional sections**: Include only when relevant to the feature
- When a section doesn't apply, remove it entirely (don't leave as "N/A")

---

## User Scenarios & Testing *(mandatory)*

### Primary User Story
A developer working on multiple projects needs to track what they're working on, capture notes and decisions in context, and maintain a timeline of their work sessions. They switch between different work contexts throughout the day (e.g., "fixing bug #123", "implementing feature X", "code review session") and need to keep each context's information separate but accessible.

When starting a new context, the previous one should automatically close, preventing confusion. All context data (notes, timestamps, file associations) should be stored in a way that allows easy review of past work sessions and supports audit trails for compliance or productivity tracking.

### Acceptance Scenarios

1. **Given** no active context exists, **When** user runs `my-context start "Working on login feature"`, **Then** a new context is created with that name, becomes the active context, and the start timestamp is recorded

2. **Given** an active context "Bug fixing", **When** user runs `my-context start "Code review"`, **Then** the "Bug fixing" context is automatically stopped with an end timestamp, and "Code review" becomes the new active context

3. **Given** an active context exists, **When** user runs `my-context note "Changed API endpoint to /v2/users"`, **Then** the note is associated with the current context with a timestamp

4. **Given** an active context exists, **When** user runs `my-context file /path/to/modified-file.js`, **Then** the file path is associated with the current context for reference

5. **Given** an active context exists, **When** user runs `my-context touch`, **Then** a timestamp is recorded in the current context log (indicating activity without detailed notes)

6. **Given** an active context exists, **When** user runs `my-context show`, **Then** the system displays the current context name, start time, elapsed duration, recent notes, and associated files

7. **Given** multiple contexts have been created, **When** user runs `my-context list`, **Then** the system displays all contexts (active and stopped) with names, start times, durations, and status

8. **Given** context transitions have occurred, **When** user runs `my-context history`, **Then** they can see a chronological history of all context switches with timestamps and context names

9. **Given** a user works across cmd.exe, git-bash, and WSL, **When** they run any my-context command, **Then** the commands work identically across all environments and access the same context data

10. **Given** a context has associated files and notes, **When** user exports or backs up their context data, **Then** all context information is accessible as plain text files in the home directory

11. **Given** a context named "Bug fix" already exists, **When** user runs `my-context start "Bug fix"`, **Then** a new context "Bug fix_2" is created and becomes active

12. **Given** any command is executed, **When** user adds `--json` flag, **Then** the output is formatted as valid JSON for machine parsing

### Edge Cases
- What happens when user tries to add a note but no context is active?
  - System should prompt user to start a context first or auto-create a default context
- What happens when the same file is associated with multiple contexts?
  - Both associations should be preserved; contexts are independent
- What happens when user provides a very long context name (>200 characters)?
  - System should truncate or reject with helpful error message
- What happens when context data storage becomes corrupted?
  - System should detect corruption and provide recovery options or fail gracefully
- What happens when trying to start a context with the same name as an existing one?
  - System automatically appends sequence suffix (_2, _3, etc.) starting with first duplicate
- What happens when user runs `my-context stop` when no context is active?
  - System should display "No active context" message (not an error)
- What happens when context transitions occur across different shell environments?
  - System must maintain state in shared home directory accessible by all shells
- What happens when creating "Bug fix_2" but "Bug fix" and "Bug fix_2" already exist?
  - System creates "Bug fix_3" by finding the next available sequence number

## Requirements *(mandatory)*

### Functional Requirements

**Core Commands**
- **FR-001**: System MUST provide a `start <name>` command that creates and activates a new context with the given name
- **FR-002**: System MUST automatically stop the previous active context when a new context is started
- **FR-003**: System MUST provide a `stop` command that explicitly ends the current active context
- **FR-004**: System MUST provide a `note <text>` command that adds a timestamped note to the active context
- **FR-005**: System MUST provide a `file <path>` command that associates a file path with the active context
- **FR-006**: System MUST provide a `touch` command that records a timestamp in the active context without additional data
- **FR-007**: System MUST provide a `show` command that displays the current active context details
- **FR-008**: System MUST provide a `list` command that displays all contexts (active and historical)
- **FR-009**: System MUST provide a `history` command that displays the chronological transition log

**Command Aliases**
- **FR-010**: Each command MUST have a single-letter alias (s=start, p=stop, n=note, f=file, t=touch, w=show, l=list, h=history)

**State Management**
- **FR-011**: System MUST maintain exactly one active context at any time (or none)
- **FR-012**: System MUST persist context state across shell sessions and environment switches
- **FR-013**: System MUST record timestamps for context start, context stop, notes, file associations, and touch events
- **FR-014**: System MUST store all context data in a dedicated home directory separate from the installation path

**Duplicate Name Handling**
- **FR-015**: When creating a context with an existing name, system MUST automatically append sequence suffix (_2, _3, etc.)
- **FR-016**: Sequence suffix MUST start with "_2" for the first duplicate (original has no suffix)
- **FR-017**: System MUST find the next available sequence number by checking existing contexts

**Context Transitions**
- **FR-018**: System MUST automatically log all context transitions (start/stop events) to a central transition log
- **FR-019**: System MUST include in transition log: timestamp, previous context name, new context name, transition type
- **FR-020**: Transition log MUST be accessible via the `history` command

**Data Storage**
- **FR-021**: System MUST store all context data as plain text files for portability and direct access
- **FR-022**: Each context MUST have its own isolated subdirectory within the context home directory
- **FR-023**: System MUST automatically create a new subdirectory when a context is created
- **FR-024**: Context subdirectory MUST start empty by default when created
- **FR-025**: System MUST support storing arbitrary file attachments or references within context subdirectory
- **FR-026**: Context home directory MUST be relocatable without breaking functionality

**Cross-Platform**
- **FR-027**: All commands MUST work identically in cmd.exe, PowerShell, git-bash, and WSL environments
- **FR-028**: File path handling MUST automatically normalize Windows (backslash) and POSIX (forward slash) paths
- **FR-029**: System MUST provide a single executable with environment detection and wrapper/shim support

**Output Format**
- **FR-030**: Commands MUST write data output to stdout and errors to stderr
- **FR-031**: Output format MUST be both human-readable and parseable by standard shell tools (grep, awk, etc.)
- **FR-032**: All commands MUST support a `--json` flag that formats output as valid JSON
- **FR-033**: JSON output MUST include all relevant data fields with appropriate type encoding

**Help System**
- **FR-034**: System MUST provide built-in help accessible via `my-context help` or `my-context --help`
- **FR-035**: Each subcommand MUST have help text accessible via `my-context <command> --help`

### Key Entities

- **Context**: Represents a work session with a name, start timestamp, optional end timestamp, collection of notes, collection of file associations, and activity timestamps
  - Attributes: name (with optional sequence suffix), start_time, end_time (nullable), status (active/stopped), subdirectory_path
  - Relationships: has many Notes, has many FileAssociations, has many TouchEvents
  - Naming rules: Duplicate names get automatic _2, _3, etc. suffix

- **Note**: A timestamped text entry associated with a context
  - Attributes: timestamp, text_content
  - Relationships: belongs to one Context

- **FileAssociation**: A file path reference linked to a context
  - Attributes: timestamp, file_path (normalized for cross-platform compatibility)
  - Relationships: belongs to one Context

- **TouchEvent**: A simple timestamp indicating activity in a context
  - Attributes: timestamp
  - Relationships: belongs to one Context

- **ContextTransition**: A log entry recording a change between contexts
  - Attributes: timestamp, previous_context_name (nullable), new_context_name (nullable), transition_type (start/stop/switch)
  - Relationships: references Context entities but stored independently in central log
  - Accessible via: `history` command with optional --json output

- **ContextHome**: The storage location containing all context data and the transition log
  - Attributes: home_directory_path, installation_path (separate)
  - Contains: subdirectories (one per context), transition log file, metadata files
  - Structure: home_directory/context_name/ for each context

---

## Review & Acceptance Checklist
*GATE: Automated checks run during main() execution*

### Content Quality
- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness
- [x] No [NEEDS CLARIFICATION] markers remain - **All 4 clarifications resolved**
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

### Cross-Platform Considerations *(for CLI tools)*
- [x] Path handling scenarios identified (Windows vs POSIX)
- [x] Shell environment compatibility addressed (cmd/PowerShell/bash/WSL)
- [x] Text I/O format specified (for composability with shell tools)
- [x] Platform-specific edge cases documented

---

## Execution Status
*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted: context lifecycle, command set, state management, transition logging, subdirectory-based storage
- [x] Ambiguities marked: 4 [NEEDS CLARIFICATION] items identified
- [x] Clarifications completed: All 4 items resolved (2025-10-04)
- [x] User scenarios defined: 12 acceptance scenarios + 8 edge cases
- [x] Requirements generated: 35 functional requirements
- [x] Entities identified: 6 key entities with clarified attributes
- [x] Review checklist passed: All gates cleared

---

## Next Steps

✅ **Specification complete and ready for planning!**

All clarifications have been integrated. Proceed to the planning phase with:

```bash
/plan 001-cli-context-management
```

The planning phase will design the technical implementation using:
- Subdirectory-based storage (one per context)
- History command for transition log access
- JSON output support across all commands
- Automatic duplicate name handling with sequence suffixes
