# Feature Specification: Context Lifecycle Improvements

**Feature Branch**: `004-implement-5-lifecycle`
**Created**: 2025-10-10
**Status**: Draft
**Input**: User description: "Implement 5 lifecycle improvement features validated via POC testing"

## Context & Rationale

### Problem Statement
Real-world usage during deb-sanity Sprint 006 revealed context lifecycle management gaps:
- **16 fragmented contexts** for single sprint (should have been 1-3)
- **69 active contexts** with only 1 archived (cleanup inefficiency)
- **Unclear lifecycle guidance** (when to start/stop/archive)

### POC Validation
All 5 features implemented as shell scripts in `scripts/poc/` and validated:
- ✅ UX patterns proven effective
- ✅ Addresses real fragmentation issues
- ✅ Ready for Go implementation

### Source
- Analysis: `docs/upstream-tracking/my-context-usage-improvements.md` (deb-sanity)
- Evidence: Sprint 006 usage data, POC testing feedback

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Smart Resume on Duplicate Names (Priority: P1)

Developer stops work on "feature-auth" context. Next day, types `my-context start "feature-auth"` intending to resume but doesn't remember if context exists. System detects duplicate name and prompts to resume existing or create new with different name.

**Why this priority**: **Prevents 80% of context fragmentation**. Sprint 006 evidence: 16 contexts with `_2`, `_3` suffixes from unintentional duplication. Root cause of fragmentation.

**Independent Test**: Create stopped context, attempt start with same name, verify prompt shows context summary and resume option. Delivers immediate fragmentation prevention.

**Acceptance Scenarios**:

1. **Given** context "foo" exists (stopped, 15 notes), **When** `my-context start "foo"`, **Then** prompt: "Context 'foo' exists (stopped, 15 notes). Resume? [Y/n]"
2. **Given** user confirms resume (Y), **When** context activates, **Then** all previous notes and files preserved
3. **Given** user declines (n), **When** prompted for new name, **Then** user provides different name, new context created
4. **Given** `--force` flag used, **When** duplicate detected, **Then** prompt for different name (no resume option)
5. **Given** context "foo" is active, **When** `my-context start "foo"`, **Then** error: "Context already active"

---

### User Story 2 - [Brief Title] (Priority: P2)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 3 - [Brief Title] (Priority: P3)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST [specific capability, e.g., "allow users to create accounts"]
- **FR-002**: System MUST [specific capability, e.g., "validate email addresses"]  
- **FR-003**: Users MUST be able to [key interaction, e.g., "reset their password"]
- **FR-004**: System MUST [data requirement, e.g., "persist user preferences"]
- **FR-005**: System MUST [behavior, e.g., "log all security events"]

*Example of marking unclear requirements:*

- **FR-006**: System MUST authenticate users via [NEEDS CLARIFICATION: auth method not specified - email/password, SSO, OAuth?]
- **FR-007**: System MUST retain user data for [NEEDS CLARIFICATION: retention period not specified]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: [Measurable metric, e.g., "Users can complete account creation in under 2 minutes"]
- **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
- **SC-003**: [User satisfaction metric, e.g., "90% of users successfully complete primary task on first attempt"]
- **SC-004**: [Business metric, e.g., "Reduce support tickets related to [X] by 50%"]

### User Story 2 - Note Count Warnings (Priority: P2)

Developer adds notes throughout complex sprint. At 50 notes, system warns: "Consider chunking if switching focus, or continue if same work (100+ notes fine)."

**Why this priority**: Proactive guidance with zero friction. Trivial (<1 day), high value. Non-blocking warnings guide without disrupting workflow.

**Independent Test**: Add 50 notes, verify warning appears with helpful guidance. Warning doesn't block addition.

**Acceptance Scenarios**:

1. **Given** 49 notes, **When** add 50th, **Then** warning with chunking guidance
2. **Given** 99 notes, **When** add 100th, **Then** "Context has 100 notes (getting large)"  
3. **Given** 200+ notes, **When** add at 225, **Then** periodic warning (every 25)
4. **Given** `MC_WARN_AT=30`, **When** add 30th, **Then** custom threshold warning
5. **Given** warning shown, **When** continue, **Then** exit code 0 (no breakage)

---

### User Story 3 - Explicit Resume Command (Priority: P2)

Developer uses `my-context resume <name>` for semantic clarity. `my-context resume --last` quickly resumes most recent.

**Why this priority**: Semantic clarity + workflow efficiency. `resume --last` common ("pick up where left off"). Complements P1 smart resume.

**Independent Test**: Stop context, `resume` to reactivate. Test `--last`. Pattern matching.

**Acceptance Scenarios**:

1. **Given** "foo" stopped, **When** `my-context resume "foo"`, **Then** "Resumed context: foo (25 notes)"
2. **Given** 3 stopped, **When** `resume --last`, **Then** most recent activates
3. **Given** "007-github", "007-integration" stopped, **When** `resume "007-*"`, **Then** selection prompt
4. **Given** single match, **When** `resume "007-github*"`, **Then** activates without prompt
5. **Given** doesn't exist, **When** `resume "nonexistent"`, **Then** error + available list

---

### User Story 4 - Bulk Archive by Pattern (Priority: P3)

After Sprint 006, 16 contexts "006-*" need archiving. `my-context archive --pattern "006-*"` with dry-run first, then confirm.

**Why this priority**: Cleanup efficiency. Evidence: 69 active (1 archived) - manual archiving abandoned due to tedium.

**Independent Test**: Create test contexts with common prefix, bulk archive with --dry-run, verify selection, archive all.

**Acceptance Scenarios**:

1. **Given** 16 "006-*" stopped, **When** `archive --pattern "006-*" --dry-run`, **Then** lists 16: "DRY RUN"
2. **Given** dry-run verified, **When** `archive --pattern "006-*"`, **Then** "Archive all 16? [y/N]"
3. **Given** confirm, **When** proceeds, **Then** all archived, "Archived 16 context(s)"
4. **Given** pattern matches active, **When** runs, **Then** skip active with warning
5. **Given** `--completed-before "2025-10-09"`, **When** runs, **Then** only before date archived
6. **Given** `--all-stopped`, **When** runs, **Then** all stopped archived after confirm

---

### User Story 5 - Lifecycle Advisor (Priority: P3)

After stopping, system analyzes (42 notes, "retrospective complete"), finds related contexts, suggests: resume related, archive if complete, or start new.

**Why this priority**: Helpful but not critical. P1-P3 address problem more directly. Adds polish, lower ROI.

**Independent Test**: Stop context with completion keywords, verify suggestions + related detection + archive recommendation.

**Acceptance Scenarios**:

1. **Given** 42 notes, **When** stop, **Then** summary: name, duration, count
2. **Given** related exist (similar prefix), **When** stop, **Then** list up to 3 related
3. **Given** "foo-phase-2" related, **When** stop, **Then** "Resume related: my-context start 'foo-phase-2'"
4. **Given** completion keywords in recent notes, **When** stop, **Then** "Completion detected. Consider archiving."
5. **Given** no related, **When** stop, **Then** archive if complete, start new

---

### Edge Cases

- Smart Resume: Names with special chars/spaces - handle quoted names
- Smart Resume: Start, stop, try again - should prompt resume
- Note Warnings: 200+ from imports - warnings informational
- Resume: Multiple exact matches - selection UI
- Bulk Archive: Hundreds match - show count + first 10, confirm
- Bulk Archive: One fails - continue with others, report failures
- Lifecycle Advisor: Stopped <1min after start - skip completion logic

---

## Requirements

### Functional Requirements

**FR-001**: Detect `my-context start` with name matching stopped context
**FR-002**: Display interactive prompt with context summary when duplicate
**FR-003**: Allow resume (Y), new name (n), or cancel
**FR-004**: Support `--force` to bypass resume prompt
**FR-005**: Preserve notes/files when resuming  
**FR-006**: Warn at configurable thresholds (50, 100, 200+)
**FR-007**: Warnings non-blocking (exit 0)
**FR-008**: Support env vars for custom thresholds
**FR-009**: Provide `resume` command for stopped contexts
**FR-010**: Support `resume --last` for most recent
**FR-011**: Pattern matching in `resume`
**FR-012**: Selection UI when multiple matches
**FR-013**: Provide `archive --pattern` for glob matching
**FR-014**: Support `--dry-run` to preview
**FR-015**: Confirmation before bulk archive
**FR-016**: Support `--completed-before` date filtering
**FR-017**: Support `--all-stopped` flag
**FR-018**: Skip active in bulk archive with warning
**FR-019**: Display summary after `stop`
**FR-020**: Detect related contexts, suggest resuming
**FR-021**: Detect completion keywords in recent notes
**FR-022**: Suggest archiving when completion detected
**FR-023**: Handle names with spaces/special chars/unicode
**FR-024**: Return appropriate exit codes

### Key Entities

**Context**: Work session (name, status, notes, files, timestamps). No schema changes.

---

## Success Criteria

### Measurable Outcomes

**SC-001**: Context fragmentation reduced 80% (16→<4 for same scope)
**SC-002**: Active contexts reduced 69→<20 via bulk archive
**SC-003**: 90% users resume successfully on first attempt
**SC-004**: Zero breaking changes (all additive)
**SC-005**: Note warnings <100ms overhead
**SC-006**: Bulk archive <5s for 100 contexts
**SC-007**: Resume --last <500ms
**SC-008**: Lifecycle advisor 95% accuracy on completion keywords
**SC-009**: POC shell functionality replicated in Go with parity
**SC-010**: Comprehensive test coverage (unit + integration)

---

## Implementation Notes

### POC Validation
All 5 features validated in `scripts/poc/`:
- ✅ smart-resume.sh
- ✅ note-with-warning.sh  
- ✅ resume-alias.sh
- ✅ bulk-archive.sh
- ✅ lifecycle-advisor.sh

See `scripts/poc/README.md` for usage examples, installation, validation checklists.

### Priority
1. P1: FR-001→005 (Smart Resume) - Prevents fragmentation
2. P2: FR-006→008 (Warnings) - Trivial, immediate value
3. P2: FR-009→012 (Resume) - Workflow efficiency
4. P3: FR-013→018 (Bulk Archive) - Cleanup enabler
5. P3: FR-019→022 (Lifecycle Advisor) - Polish

### Technology
- Language: Go (existing codebase)
- CLI: cobra (in use)
- Storage: Plain text `~/.my-context/` (no changes)

### Out of Scope
- Automatic archiving, AI suggestions, context merging, cloud sync, GUI

---

## Open Questions

**Q1**: Force flag - prompt for new name (B) vs auto `_2` (A)? **Recommend: B**
**Q2**: Warnings - one-time vs periodic after 200? **Recommend: Periodic**
**Q3**: `start` backward compat - still works (A) vs force `resume` (B)? **Recommend: A**
**Q4**: Bulk limit - default 100? **Recommend: Yes, MC_BULK_LIMIT env var**
**Q5**: Completion keywords - hardcoded (A) vs configurable (B)? **Recommend: A**

---

## Next Steps

1. **Clarify**: Address Q1-Q5 (use `/clarify`)
2. **Plan**: Create implementation plan (use `/plan`)
3. **Tasks**: Generate task list (use `/tasks`)
4. **Implement**: Execute with TDD
5. **Validate**: Test against SC-001→010
6. **Retro**: Compare Sprint 007 to 006

---

**Version**: 1.0.0  
**Status**: Ready for Clarification
**Effort**: 3-5 days implementation + 2 days testing
