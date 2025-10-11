# Tasks: Context Signaling Protocol

**Input**: Design documents from `/specs/005-context-signaling-protocol/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Following TDD approach, tests will be written first where applicable.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `cmd/`, `internal/`, `pkg/` at repository root
- Paths follow Go project layout as specified in plan.md

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for signaling protocol

- [X] T001 Create Go package structure per implementation plan (cmd/, internal/, pkg/)
- [X] T002 [P] Add fsnotify dependency for optional inotify support
- [X] T003 Create signals directory structure in ~/.my-context/signals/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Extend Context model with metadata fields (created_by, parent, labels) in pkg/models/context.go
- [X] T005 Create Signal model and data structures in pkg/models/signal.go
- [X] T006 Create cross-platform filesystem utilities in pkg/utils/fs.go
- [X] T007 Initialize signal manager package in internal/signal/manager.go
- [X] T008 Initialize watch monitor package in internal/watch/monitor.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Signal File Coordination (Priority: P1) 🎯 MVP

**Goal**: Enable file-based event coordination through signal files for team notifications

**Independent Test**: Can be fully tested by creating/listing/waiting for/clearing signal files independently of watch functionality

### Tests for User Story 1 (TDD - Write FIRST) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T009 [P] [US1] Unit tests for signal manager CRUD operations in internal/signal/manager_test.go
- [X] T010 [P] [US1] Contract tests for signal CLI commands in tests/contract/test_signal_commands.go
- [X] T011 [P] [US1] Integration tests for signal file lifecycle in tests/integration/test_signal_workflow.go

### Implementation for User Story 1

- [X] T012 [US1] Implement signal create command in cmd/signal.go
- [X] T013 [US1] Implement signal list command in cmd/signal.go
- [X] T014 [US1] Implement signal wait command with timeout in cmd/signal.go
- [X] T015 [US1] Implement signal clear command in cmd/signal.go
- [X] T016 [US1] Add signal command to main CLI router
- [X] T017 [US1] Add concurrent access safety to signal operations

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently - users can create, list, wait for, and clear signal files

---

## Phase 4: User Story 2 - Context Monitoring with Watch Command (Priority: P1)

**Goal**: Enable automated monitoring of context changes with command execution

**Independent Test**: Can be fully tested by watching contexts and verifying command execution on note changes, independent of signal files

### Tests for User Story 2 (TDD - Write FIRST) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T018 [P] [US2] Unit tests for watch monitor polling logic in internal/watch/monitor_test.go
- [X] T019 [P] [US2] Unit tests for pattern matching in internal/watch/patterns_test.go
- [X] T020 [P] [US2] Contract tests for watch CLI command in tests/contract/test_watch_commands.go
- [X] T021 [P] [US2] Integration tests for watch command execution in tests/integration/test_watch_workflow.go

### Implementation for User Story 2

- [X] T022 [US2] Implement watch command with polling in cmd/watch.go
- [X] T023 [US2] Add pattern matching for note filtering in internal/watch/patterns.go
- [X] T024 [US2] Implement --exec flag for command execution on changes
- [X] T025 [US2] Add optional inotify support for Linux systems
- [X] T026 [US2] Add watch command to main CLI router
- [X] T027 [US2] Implement graceful interruption handling (Ctrl+C)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - users can monitor contexts and execute commands on changes

---

## Phase 5: User Story 3 - Enhanced Context Metadata (Priority: P2)

**Goal**: Add metadata fields for better context organization and querying

**Independent Test**: Can be fully tested by creating contexts with metadata and querying them, independent of signaling features

### Tests for User Story 3 (TDD - Write FIRST) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T028 [P] [US3] Unit tests for enhanced context metadata in pkg/models/context_test.go
- [X] T029 [P] [US3] Contract tests for metadata CLI options in tests/contract/test_metadata_commands.go

### Implementation for User Story 3

- [X] T030 [US3] Extend start command with --created-by, --parent, --labels flags
- [X] T031 [US3] Update show command to display metadata fields
- [ ] T032 [US3] Implement search command with --created-by and --label filters
- [X] T033 [US3] Add metadata validation and error handling
- [X] T034 [US3] Update context serialization to include metadata in meta.json

**Checkpoint**: All user stories should now be independently functional - metadata enhances context management without affecting signaling

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T035 [P] Update CLI help text and command documentation
- [X] T036 [P] Add performance optimizations for watch polling (<1% CPU)
- [X] T037 [P] Implement comprehensive error handling across all commands
- [X] T038 [P] Add cross-platform testing validation
- [X] T039 [P] Update README.md and TROUBLESHOOTING.md with signaling features
- [X] T040 Run quickstart.md validation for signaling examples

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - US1 and US2 can proceed in parallel (both P1 priority)
  - US3 can start after foundational but can run parallel to US1/US2 completion
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 3 (P2)**: Can start after Foundational - No dependencies on signaling stories

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Models before CLI commands
- Core logic before integration features
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, US1 and US2 can start in parallel
- US3 can run parallel to US1/US2 completion
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Signal files)
4. **STOP and VALIDATE**: Test signal file operations independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (Signal MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo (Watch MVP!)
4. Add User Story 3 → Test independently → Deploy/Demo (Metadata enhancement)
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Signal files)
   - Developer B: User Story 2 (Watch command)
   - Developer C: User Story 3 (Metadata) or help with US1/US2
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- TDD: Write failing tests first, then implement to make them pass
- Target: v2.2.0 release
