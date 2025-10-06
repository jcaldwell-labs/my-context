# Feature Specification: Installation & Usability Improvements

**Feature Branch**: `002-installation-improvements-and`  
**Created**: 2025-10-05  
**Status**: Draft  
**Input**: User description: "Installation improvements and user-requested features: multi-platform builds, project filters, export command, list enhancements, archive/delete commands"

## Clarifications

### Session 2025-10-05 (Sprint 1 Retrospective)
- Q: Which platforms need pre-built binaries? → A: Windows (amd64), Linux (amd64), macOS (amd64, arm64)
- Q: Export format preference? → A: Markdown format for easy sharing and readability
- Q: Project filter parsing strategy? → A: Parse existing "project: phase - description" naming convention, extract project name before first colon
- Q: Archive vs Delete behavior? → A: Archive hides from default list but preserves data; Delete removes entirely with confirmation
- Q: List default limit? → A: Show last 10 contexts by default, add --all flag for complete list
- Q: Installation method priority? → A: Pre-built binaries > curl installer > package managers (brew/choco) as stretch goal

## Execution Flow (main)
```
1. Parse user description from Input
   → Identified: installation blockers, user-requested features, usability gaps
2. Extract key concepts from description
   → Actors: developers on Windows/Linux/macOS/WSL, new users installing tool
   → Actions: build, install, filter, export, archive, delete, search
   → Data: binaries, contexts, metadata, export files
   → Constraints: cross-platform compatibility, backward compatibility with Sprint 1 data
3. Fill User Scenarios & Testing section
   → Installation scenarios (WSL user, Windows user, macOS user)
   → Feature usage scenarios (project filtering, export, archive)
4. Generate Functional Requirements
   → Build system, installation scripts, new flags, new commands
5. Identify Key Entities
   → Binary artifacts, installation scripts, project metadata, export documents
6. Run Review Checklist
   → All requirements testable and aligned with Sprint 1 retrospective findings
7. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

---

## User Scenarios & Testing

### Primary User Story
A developer downloads my-context for the first time on their WSL environment. They want to start using it immediately without building from source. After installation, they create multiple contexts following their natural workflow of organizing work by project (e.g., "ps-cli: Phase 1 - Foundation", "ps-cli: Phase 2 - Testing"). As contexts accumulate, they need to filter by project, export context summaries to share with their team, and archive completed work to keep their active list manageable.

### Acceptance Scenarios

#### Installation & Distribution
1. **Given** a WSL user visits the releases page, **When** they download the linux-amd64 binary, **Then** it runs without requiring rebuild or Go installation

2. **Given** a Windows user with cmd.exe, **When** they run `install.bat`, **Then** the tool is installed to their PATH and works from any directory

3. **Given** a macOS user with ARM chip, **When** they download the darwin-arm64 binary, **Then** it runs natively without Rosetta translation

4. **Given** a user on any platform, **When** they run a curl one-liner installer, **Then** the correct binary for their platform is automatically detected and installed

5. **Given** existing Sprint 1 users with data in `~/.my-context/`, **When** they upgrade to Sprint 2 binaries, **Then** all existing contexts and data remain intact and accessible

#### Project Filtering (User Request #1)
6. **Given** a user with contexts named "ps-cli: Phase 1", "ps-cli: Phase 2", "garden: Planning", **When** they run `my-context list --project ps-cli`, **Then** only "ps-cli: Phase 1" and "ps-cli: Phase 2" are displayed

7. **Given** a user wants to start a new context for an existing project, **When** they run `my-context start "Phase 3" --project ps-cli`, **Then** the context is created as "ps-cli: Phase 3" following the established convention

8. **Given** a user with contexts that don't follow the "project: phase" convention, **When** they run `my-context list --project <name>`, **Then** contexts without colons are treated as their own project (full name = project name)

#### Export Command (User Request #2)
9. **Given** a user wants to share context details, **When** they run `my-context export "ps-cli: Phase 1"`, **Then** a markdown file is created with context name, duration, notes, files, and timestamps

10. **Given** a user specifies output location, **When** they run `my-context export "Phase 1" --to contexts/phase-1.md`, **Then** the export file is written to the specified path (creating directories if needed)

11. **Given** a user wants to export all contexts, **When** they run `my-context export --all`, **Then** each context is exported to separate markdown files in the current directory or specified output directory

12. **Given** an exported markdown file, **When** a user opens it in any text editor or markdown viewer, **Then** it displays human-readable context information with proper formatting

#### List Enhancements
13. **Given** a user has 50 contexts, **When** they run `my-context list`, **Then** only the last 10 contexts are displayed by default with a message indicating more exist

14. **Given** a user wants to see all contexts, **When** they run `my-context list --all`, **Then** all contexts are displayed regardless of count

15. **Given** a user wants to find specific contexts, **When** they run `my-context list --search "bug fix"`, **Then** only contexts containing "bug fix" in their name are displayed (case-insensitive)

16. **Given** a user wants to see only incomplete work, **When** they run `my-context list --active-only`, **Then** only the currently active context is displayed

17. **Given** a user combines filters, **When** they run `my-context list --project ps-cli --limit 5`, **Then** only the 5 most recent ps-cli contexts are displayed

#### Archive & Delete Commands
18. **Given** a user completes a project phase, **When** they run `my-context archive "ps-cli: Phase 1"`, **Then** the context is marked as archived and hidden from default list views

19. **Given** archived contexts exist, **When** a user runs `my-context list --archived`, **Then** only archived contexts are displayed

20. **Given** a user wants to permanently remove a test context, **When** they run `my-context delete "Test Context"`, **Then** they are prompted for confirmation before the context directory and all data are deleted

21. **Given** a user confirms deletion, **When** the delete command completes, **Then** the context is removed from `~/.my-context/` and no longer appears in any list views

22. **Given** a user tries to delete the active context, **When** they run `my-context delete <active-context>`, **Then** they receive a warning that they must stop the context first

#### Bug Fixes
23. **Given** a user adds a note with special characters, **When** they run `my-context note "Budget: $500-800"`, **Then** the note is stored and displayed with all characters intact (including $ symbol)

24. **Given** a user views history, **When** the output includes start/stop transitions, **Then** "(none)" is displayed instead of "NULL" for empty previous/next context fields

### Edge Cases
- What happens when a user tries to export a non-existent context?
  - System should display clear error: "Context '<name>' not found"
- What happens when export output file already exists?
  - System should prompt for overwrite confirmation (y/N). If user declines, export is cancelled with exit code 2.
- What happens when filtering by project that has no contexts?
  - System should display "No contexts found for project '<name>'"
- What happens when trying to archive an already-archived context?
  - System should display message: "Context '<name>' is already archived"
- What happens when installing on a system where the binary already exists?
  - Installation script should detect existing installation and offer to upgrade or skip
- What happens when curl installer cannot detect platform?
  - Script should display error with manual download instructions
- What happens when --project flag is used but context names don't contain colons?
  - Each context name is treated as its own project (exact match required)
- What happens when deleting a context that's referenced in transitions.log?
  - Context data is deleted but transitions log remains intact (historical record preserved)

## Requirements

### Functional Requirements

#### FR-001: Multi-Platform Binary Distribution
- **FR-001.1**: System MUST provide pre-built binaries for Windows (amd64), Linux (amd64), macOS (amd64), and macOS (arm64)
- **FR-001.2**: Binaries MUST be statically linked to avoid dependency issues
- **FR-001.3**: Release artifacts MUST include checksums (SHA256) for verification
- **FR-001.4**: Binaries MUST have platform-appropriate names (e.g., `my-context.exe`, `my-context-linux`, `my-context-darwin`)

#### FR-002: Installation Scripts
- **FR-002.1**: System MUST provide `install.sh` for Unix-like systems (Linux, macOS, WSL)
- **FR-002.2**: System MUST provide `install.bat` for Windows cmd.exe
- **FR-002.3**: System MUST provide `install.ps1` for Windows PowerShell
- **FR-002.4**: All installation scripts MUST detect existing installations and handle upgrades
- **FR-002.5**: Installation MUST preserve existing `~/.my-context/` data during upgrades

#### FR-003: One-Liner Installer
- **FR-003.1**: System MUST provide a curl-based installer that auto-detects platform (Linux, macOS, Windows via Git Bash/WSL)
- **FR-003.2**: Installer MUST download the appropriate binary for the detected platform
- **FR-003.3**: Installer MUST verify binary checksum before installation
- **FR-003.4**: Installer MUST make binary executable and add to PATH

#### FR-004: Project Filtering
- **FR-004.1**: `list` command MUST support `--project <name>` flag to filter contexts by project
- **FR-004.2**: System MUST parse project name from contexts following "project: description" convention (text before first colon)
- **FR-004.3**: Contexts without colons MUST be treated as standalone projects (full name = project name)
- **FR-004.4**: `start` command MUST support `--project <name>` flag to create contexts with project prefix
- **FR-004.5**: Project filtering MUST be case-insensitive

#### FR-005: Export Command
- **FR-005.1**: System MUST provide `export` command with alias `e`
- **FR-005.2**: Export MUST accept context name as argument to export specific context
- **FR-005.3**: Export MUST support `--to <path>` flag to specify output file location
- **FR-005.4**: Export MUST support `--all` flag to export all contexts
- **FR-005.5**: Export output MUST be markdown format containing: context name, start/end times, duration, notes (with timestamps), file associations, touch events
- **FR-005.6**: Export MUST create parent directories if output path doesn't exist
- **FR-005.7**: Export MUST handle file overwrite with confirmation or timestamp suffix

#### FR-006: List Command Enhancements
- **FR-006.1**: `list` command MUST show only last 10 contexts by default
- **FR-006.2**: `list` command MUST support `--all` flag to show all contexts
- **FR-006.3**: `list` command MUST support `--limit <n>` flag to specify custom result count
- **FR-006.4**: `list` command MUST support `--search <term>` flag for case-insensitive name filtering
- **FR-006.5**: `list` command MUST support `--active-only` flag to show only active context
- **FR-006.6**: `list` command MUST display message when more contexts exist than shown (e.g., "Showing 10 of 50 contexts. Use --all to see all.")
- **FR-006.7**: Multiple filters MUST work together (e.g., `--project ps-cli --limit 5`)

#### FR-007: Archive Command
- **FR-007.1**: System MUST provide `archive` command with alias `a`
- **FR-007.2**: Archive MUST accept context name as argument
- **FR-007.3**: Archive MUST mark context as archived in metadata
- **FR-007.4**: Archived contexts MUST be hidden from default `list` views
- **FR-007.5**: `list --archived` MUST show only archived contexts
- **FR-007.6**: Archive operation MUST preserve all context data (notes, files, timestamps)
- **FR-007.7**: Archive MUST prevent archiving the currently active context

#### FR-008: Delete Command
- **FR-008.1**: System MUST provide `delete` command with alias `d`
- **FR-008.2**: Delete MUST accept context name as argument
- **FR-008.3**: Delete MUST prompt for confirmation before removing context
- **FR-008.4**: Delete MUST support `--force` flag to skip confirmation
- **FR-008.5**: Delete MUST remove entire context directory from `~/.my-context/`
- **FR-008.6**: Delete MUST prevent deleting the currently active context without stopping it first
- **FR-008.7**: Delete MUST preserve transitions.log history (don't remove transition entries)

#### FR-009: Bug Fixes
- **FR-009.1**: Note command MUST preserve all special characters including $ symbol
- **FR-009.2**: History command MUST display "(none)" instead of "NULL" for empty context fields
- **FR-009.3**: Note display MUST show full note text including currency symbols

#### FR-010: Documentation
- **FR-010.1**: README MUST include "Building from Source" section with platform-specific instructions
- **FR-010.2**: README MUST include troubleshooting section for WSL users
- **FR-010.3**: System MUST provide TROUBLESHOOTING.md with common installation issues
- **FR-010.4**: Each installation script MUST have inline comments explaining steps

#### FR-011: Backward Compatibility
- **FR-011.1**: Sprint 2 binaries MUST work with Sprint 1 data structures
- **FR-011.2**: Existing contexts from Sprint 1 MUST remain functional after upgrade
- **FR-011.3**: Archive status MUST be added to existing meta.json without breaking existing fields

### Key Entities

#### Binary Artifact
- Platform identifier (windows-amd64, linux-amd64, darwin-amd64, darwin-arm64)
- Executable name
- SHA256 checksum
- Download URL
- Version number

#### Project Metadata
- Project name (extracted from context name before first colon)
- Associated context names
- Filter criteria

#### Export Document
- Source context name
- Export timestamp
- Formatted markdown content (headers, lists, timestamps)
- Output file path

#### Archive Status
- Is archived (boolean flag in meta.json)
- Archive timestamp
- Archived by user (if tracking needed)

---

## Review & Acceptance Checklist

### Content Quality
- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness
- [x] No [NEEDS CLARIFICATION] markers remain (all resolved in Clarifications section)
- [x] Requirements are testable and unambiguous  
- [x] Success criteria are measurable (22 acceptance scenarios)
- [x] Scope is clearly bounded (Sprint 2 focus: installation + user requests)
- [x] Dependencies and assumptions identified (Sprint 1 data compatibility)

---

## Execution Status

- [x] User description parsed
- [x] Key concepts extracted (installation, filtering, export, archive, bug fixes)
- [x] Ambiguities marked and resolved (Clarifications section)
- [x] User scenarios defined (22 acceptance scenarios + 8 edge cases)
- [x] Requirements generated (11 major requirement categories with sub-requirements)
- [x] Entities identified (Binary Artifact, Project Metadata, Export Document, Archive Status)
- [x] Review checklist passed

---

**Alignment with Constitution v1.1.0**:
- Principle VI (User-Driven Design): ✅ Project filter and export features respond directly to observed user behavior from Sprint 1 retrospective
- Principle II (Cross-Platform Compatibility): ✅ Multi-platform binaries address critical WSL installation blocker
- Principle IV (Minimal Surface Area): ✅ Adds only 3 new commands (export, archive, delete) + flags to existing commands
- Principle V (Data Portability): ✅ Export to markdown maintains plain-text philosophy

**Sprint 1 Retrospective Alignment**:
- Addresses "What Went Wrong" items #1, #2, #3, #5, #7, #9
- Delivers user-requested features ranked as P1 (High Value, Low Complexity)
- Fixes identified bugs ($ character, NULL display)
- Total effort: 5.5 days (aligned with Sprint 2 recommendation)
