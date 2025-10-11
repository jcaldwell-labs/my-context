# Feature Specification: Context Signaling Protocol

**Feature Branch**: `005-context-signaling-protocol`
**Created**: 2025-10-11
**Status**: Draft
**Input**: Three main features: (1) Signal files for event coordination, (2) Watch command with inotify for monitoring, (3) Metadata enhancement (created_by, parent, labels). Target: v2.2.0.

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Signal File Coordination (Priority: P1)

As a development team member, I want to create signal files to notify other team members or processes when important events occur (like binary updates or context state changes), so that I can enable automated coordination without manual notifications.

**Why this priority**: This is the foundation for event-driven coordination, enabling the core use case of automated notifications for binary updates and context changes that we identified in our handoff workflow gaps.

**Independent Test**: Can be fully tested by creating a signal file and verifying it exists, then clearing it and confirming it's removed. This delivers the basic signaling infrastructure without requiring watch functionality.

**Acceptance Scenarios**:

1. **Given** a clean signals directory, **When** I run `my-context signal create binary-updated`, **Then** a file `~/.my-context/signals/binary-updated.signal` is created with current timestamp
2. **Given** multiple signal files exist, **When** I run `my-context signal list`, **Then** all signal files are displayed with their creation times
3. **Given** a signal file exists, **When** I run `my-context signal wait binary-updated --timeout=5s`, **Then** the command returns immediately with success
4. **Given** no signal file exists, **When** I run `my-context signal wait missing-signal --timeout=1s`, **Then** the command times out after 1 second
5. **Given** a signal file exists, **When** I run `my-context signal clear binary-updated`, **Then** the signal file is removed

---

### User Story 2 - Context Monitoring with Watch Command (Priority: P1)

As a team supervisor monitoring multiple development contexts, I want to watch for changes in specific contexts and execute commands when new notes are added, so that I can get notified of progress without constantly checking manually.

**Why this priority**: This addresses the core gap of manual polling for context updates, enabling event-driven notifications for handoff coordination that we experienced in our recent team collaboration.

**Independent Test**: Can be fully tested by watching a test context, adding a note from another terminal, and verifying the watch command detects the change and executes the specified command. This delivers the monitoring capability without requiring signal files.

**Acceptance Scenarios**:

1. **Given** a context with existing notes, **When** I run `my-context watch test-context --new-notes --exec="echo detected"`, **Then** the command starts monitoring and displays "Watching for new notes in test-context"
2. **Given** a watch command is running, **When** I add a note to the watched context from another terminal, **Then** the watch command detects the change within 5 seconds and executes the specified command
3. **Given** a watch command is monitoring with `--pattern="Phase.*complete"`, **When** I add a note matching the pattern, **Then** the command executes; when I add a note not matching, **Then** no execution occurs
4. **Given** a watch command is running, **When** I interrupt it with Ctrl+C, **Then** it stops gracefully without errors

---

### User Story 3 - Enhanced Context Metadata (Priority: P2)

As a context user tracking context relationships and ownership, I want contexts to include metadata like created_by, parent context references, and labels for better organization and querying, so that I can understand context lineage and ownership.

**Why this priority**: This enhances context discoverability and organization, supporting better workflow management as we scale to more contexts and team members.

**Independent Test**: Can be fully tested by creating a context with metadata, querying it, and verifying the metadata is preserved. This delivers the metadata enhancement without requiring the signaling features.

**Acceptance Scenarios**:

1. **Given** I'm creating a new context, **When** I specify `--created-by="alice" --parent="sprint-005" --labels="feature,backend"`, **Then** the context metadata includes these fields
2. **Given** a context exists with metadata, **When** I run `my-context show context-name --metadata`, **Then** all metadata fields are displayed
3. **Given** multiple contexts exist, **When** I run `my-context search --created-by="alice"`, **Then** only contexts created by alice are returned
4. **Given** a context has labels, **When** I run `my-context search --label="feature"`, **Then** only contexts with that label are returned

### Edge Cases

- What happens when signal files are created simultaneously by multiple processes?
- How does the watch command handle context deletion while monitoring?
- What happens when metadata fields contain special characters or very long values?
- How does the system handle watch commands when disk space is full?
- What happens when signal wait timeouts occur during system suspend/resume?
- How does the watch command behave when the executable command fails?
- What happens when context metadata is corrupted or missing required fields?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST support creating signal files for event coordination with `my-context signal create <name>`
- **FR-002**: System MUST support listing existing signal files with `my-context signal list`
- **FR-003**: System MUST support waiting for signal files with timeout using `my-context signal wait <name> --timeout=<duration>`
- **FR-004**: System MUST support clearing/removing signal files with `my-context signal clear <name>`
- **FR-005**: System MUST support watching contexts for new notes with `my-context watch <context> --new-notes --exec="<command>"`
- **FR-006**: System MUST support pattern matching for watched notes with `--pattern` option
- **FR-007**: System MUST support context metadata including created_by, parent, and labels fields
- **FR-008**: System MUST support querying contexts by metadata with `my-context search --created-by=<user>` and `my-context search --label=<label>`
- **FR-009**: System MUST detect context changes within 5 seconds for watch commands
- **FR-010**: System MUST be backward compatible - existing workflows continue to work without signaling
- **FR-011**: System MUST handle concurrent access to signal files safely
- **FR-012**: System MUST work cross-platform (Linux, macOS, Windows) using polling-based approach initially

### Key Entities *(include if feature involves data)*

- **Signal File**: Represents an event notification, stored as `~/.my-context/signals/<name>.signal` with creation timestamp
- **Context Metadata**: Enhanced context information stored in `meta.json` including created_by (string), parent (string), labels (array of strings)
- **Watch Process**: Background monitoring process that polls context files for changes and executes commands

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: Binary update signals detected within 10 seconds of creation
- **SC-002**: Context note changes detected within 5 seconds by watch commands
- **SC-003**: Watch command CPU usage <1% during monitoring (polling mode)
- **SC-004**: Signal files support concurrent readers/writers without corruption
- **SC-005**: All signaling features work on Linux, macOS, and Windows
- **SC-006**: Backward compatibility maintained - existing my-context workflows function unchanged
