# my-context Best Practices for Spec Kit Workflow

**Version**: 1.0.0
**Date**: 2025-10-10
**Purpose**: Guide consistent context management during spec kit execution

---

## Context Size Guidelines

### Optimal Context Size: **50 notes per phase**

**Rationale**:
- 20 notes: ❌ Too small (caused 16 fragmented contexts in Sprint 006)
- 50 notes: ✅ **Sweet spot** - Captures full phase with room for details
- 75-100 notes: ⚠️ Getting large but acceptable for complex phases
- 100+ notes: Consider splitting into sub-phases or exporting

### When to Split Context

**Start new context when**:
- Transitioning between major phases (spec → implementation → testing)
- Switching focus areas (frontend → backend, feature A → feature B)
- Approaching 75 notes (proactive split)
- Starting new sprint or major milestone

**Continue existing context when**:
- Same phase, same focus area
- Closely related work (refactoring same module)
- Under 75 notes total
- Temporary spike in detail (debugging, complex feature)

---

## Context Naming Convention

### Pattern: `{project}: {feature}__{phase}`

**Examples**:
```
my-context: 004-lifecycle__spec-plan-tasks
my-context: 004-lifecycle__implementation-phase1
my-context: 004-lifecycle__testing-integration
deb-sanity: 006-worktrees__spec-plan-tasks
deb-sanity: 006-worktrees__implementation-core
```

**Components**:
- `{project}`: Project name (my-context, deb-sanity, ps-cli)
- `{feature}`: Feature number and short name (004-lifecycle, 006-worktrees)
- `{phase}`: Current phase (spec-plan-tasks, implementation, testing, retro)

**Sub-phases** (if needed):
```
my-context: 004-lifecycle__implementation-phase1  # Smart Resume
my-context: 004-lifecycle__implementation-phase2  # Warnings + Resume
my-context: 004-lifecycle__implementation-phase3  # Bulk + Advisor
```

**Avoid**:
- Numbered suffixes: ❌ `foo_2`, `foo_3` (indicates fragmentation)
- Ambiguous names: ❌ `test`, `setup`, `bugfix`
- Overly long: ❌ `my-context-lifecycle-improvements-implementation-smart-resume-and-warnings`

---

## What to Note

### Phase Milestones
```bash
my-context note "Starting Phase 1: Smart Resume implementation (T001-T017)"
my-context note "Phase 1 complete: Smart Resume working - duplicate detection, resume prompt, tests passing"
my-context note "Blocker: FR-003 needs clarification on edge case XYZ"
my-context note "Decision: Using approach B for interactive prompts (better UX)"
```

### After Completing Tasks
```bash
# After significant task or group of tasks
my-context note "T001-T010 complete: Core state.go methods added, tests passing"

# After checkpoint
my-context note "Checkpoint 1 passed: Smart resume working, 17/17 tasks complete"
```

### Decisions & Blockers
```bash
my-context note "Decision: Hardcoded completion keywords (Q5 answer A) for simplicity"
my-context note "Blocker: Pattern matching library needed - researching options"
my-context note "Resolved blocker: Using filepath.Match() from stdlib"
```

### Session Boundaries
```bash
# At end of work session
my-context note "End of day: 35/121 tasks complete, Phase 1 + 2A done, Phase 2B in progress"

# When splitting phases
my-context note "Phase 1 complete, stopping to start fresh context for Phase 2"
my-context stop
my-context start "my-context: 004-lifecycle__implementation-phase2"
```

---

## What to File (Track)

### Specification Artifacts
```bash
# At start of phase
my-context file "specs/004-implement-5-lifecycle/spec.md"
my-context file "specs/004-implement-5-lifecycle/tasks.md"
my-context file "specs/004-implement-5-lifecycle/plan.md"
```

### Source Files Being Implemented
```bash
# As you create/modify significant files
my-context file "internal/commands/start.go"
my-context file "internal/commands/resume.go"
my-context file "internal/core/state.go"
my-context file "tests/integration/smart_resume_test.go"
```

### Configuration & Documentation
```bash
# Important config or doc changes
my-context file "README.md"
my-context file ".env.example"
my-context file "CHANGELOG.md"
```

**Avoid**:
- Generated files (binaries, build artifacts)
- Temporary test files
- IDE-specific files
- Files you only read (not modified)

**Guideline**: Track files where you make meaningful changes (10+ lines or significant logic)

---

## What to Touch

### Generally NOT needed for implementation

`my-context touch` is for recording timestamps without notes. Rarely needed during spec kit workflow.

**Use touch for**:
- Recording presence without specific note ("checked in on project")
- Timestamp-only events (started meeting, resumed after break)

**Better alternatives**:
- Use `my-context note` instead (provides context)
- Files already track modification times
- Notes capture more valuable information

**Recommendation**: **Skip `touch` command** - use `note` and `file` instead

---

## Spec Kit Phase Mapping

### Specify → Clarify → Plan → Tasks

**Context**: `{project}: {feature}__spec-plan-tasks`
**Target**: 20-30 notes for all planning phases combined
**Files**: spec.md, clarify.md, plan.md, tasks.md

**Example notes**:
```bash
my-context start "my-context: 004-lifecycle__spec-plan-tasks"
my-context note "Spec phase: 5 user stories defined with priorities P1-P3"
my-context note "Clarify phase: 5 questions answered"
my-context note "Plan phase: 5-phase implementation plan created"
my-context note "Tasks phase: 121 tasks generated"
my-context note "Planning complete - ready for implementation"
my-context stop
```

---

### Implementation Phase (TDD)

**Context per major phase** (when exceeding 50 notes):
- `{project}: {feature}__implementation-phase1` (P1 features)
- `{project}: {feature}__implementation-phase2` (P2 features)
- `{project}: {feature}__implementation-phase3` (P3 features)

**OR single context if <75 notes total**:
- `{project}: {feature}__implementation`

**Target**: 50 notes per phase, split if exceeding 75

**Example notes**:
```bash
my-context start "my-context: 004-lifecycle__implementation-phase1"
my-context file "internal/commands/start.go"
my-context note "T001 complete: Added FindContextByName to state.go"
my-context note "T005-T009 complete: Interactive prompt logic implemented"
my-context note "Phase 1 checkpoint: Smart resume working, all tests passing (17/17 tasks)"
my-context stop

my-context start "my-context: 004-lifecycle__implementation-phase2"
my-context note "Starting Phase 2: Note warnings + Resume command"
# ... continue
```

---

### Testing & Integration Phase

**Context**: `{project}: {feature}__testing-integration`
**Target**: 30-40 notes
**Focus**: Integration tests, performance validation, bug fixes

**Example notes**:
```bash
my-context start "my-context: 004-lifecycle__testing-integration"
my-context note "Running full test suite: 121/121 tests passing"
my-context note "Performance benchmarks: all targets met (SC-005→007)"
my-context note "Integration test: All 5 features working together"
my-context note "Bug found: Pattern matching fails on special chars - fixing"
my-context note "Bug fixed: Escaped special chars in glob patterns"
my-context stop
```

---

### Retrospective Phase

**Context**: `{project}: {feature}__retrospective`
**Target**: 10-15 notes
**Focus**: Lessons learned, metrics, improvements

**Example notes**:
```bash
my-context start "my-context: 004-lifecycle__retrospective"
my-context note "Implementation complete: 5 days actual vs 5 days estimated"
my-context note "Success metrics: Context fragmentation reduced 85% (exceeds 80% target)"
my-context note "Lesson: POC validation saved 2 days of rework"
my-context note "Tech debt: Pattern matching could use optimization"
my-context stop
my-context archive "my-context: 004-lifecycle__retrospective"
```

---

## Auto-Detection Patterns (for Spec Kit Commands)

### How to Extract Context Name from FEATURE_DIR

```bash
# From FEATURE_DIR path: /path/to/specs/004-implement-5-lifecycle/
FEATURE_NUM=$(basename "$FEATURE_DIR")           # 004-implement-5-lifecycle
PROJECT_NAME=$(basename "$(git rev-parse --show-toplevel)") # my-context-dev → my-context

# Generate context name
CONTEXT_NAME="${PROJECT_NAME}: ${FEATURE_NUM}__implementation"
# Result: "my-context: 004-lifecycle__implementation"
```

### Smart Resume Integration

```bash
# Check if context exists
if my-context list | grep -q "○.*${CONTEXT_NAME}"; then
    echo "📋 Context '${CONTEXT_NAME}' exists (stopped). Resuming..."
    my-context start "${CONTEXT_NAME}"  # Will use smart resume when FR-MC-NEW-001 implemented!
else
    my-context start "${CONTEXT_NAME}"
fi
```

---

## Context Management During Implementation

### Start of Implementation
```bash
# Extract feature info
FEATURE_DIR="/path/to/specs/004-implement-5-lifecycle"
FEATURE_NAME="004-lifecycle"
TASK_COUNT=$(grep -c "^- \[ \]" "$FEATURE_DIR/tasks.md")
PHASE_COUNT=$(grep -c "^## Phase" "$FEATURE_DIR/tasks.md")

# Start session
my-context start "my-context: ${FEATURE_NAME}__implementation"
my-context file "$FEATURE_DIR/tasks.md"
my-context file "$FEATURE_DIR/spec.md"
my-context note "Starting implementation: $TASK_COUNT tasks across $PHASE_COUNT phases"
```

### During Implementation (After Each Significant Milestone)
```bash
# After completing task group
my-context note "T001-T017 complete: Smart Resume implemented and tested"

# When creating new files
my-context file "internal/commands/resume.go"

# When hitting checkpoint
my-context note "Checkpoint 1: Smart Resume working - all acceptance scenarios passing"
```

### Splitting Phases (When Approaching 50 Notes)
```bash
# Check note count
NOTE_COUNT=$(my-context show | grep -oP '(?<=Notes \()\d+')

if [[ $NOTE_COUNT -ge 45 ]]; then
    # Stop and start new phase
    my-context note "Phase 1 complete - splitting context (approaching 50 notes)"
    my-context stop
    my-context start "my-context: 004-lifecycle__implementation-phase2"
    my-context note "Continuing implementation: Starting Phase 2 (Note Warnings + Resume Command)"
fi
```

### End of Implementation
```bash
# Final checkpoint
my-context note "Implementation complete: All 121 tasks done, tests passing, documentation updated"
my-context note "Success criteria validated: SC-001→010 all met"
my-context note "Ready for merge: Branch 004-implement-5-lifecycle → master"

# Stop session
my-context stop
```

---

## Recommended Context Flow for Sprint 004

### Planning Phase (DONE ✅)
```
Context: "my-context: 004-lifecycle__spec-plan-tasks"
Duration: 15 minutes
Notes: 7 notes (spec → clarify → plan → tasks progression)
```

### Implementation Phase (NEXT)

**Option A: Single Context** (if confident <75 notes total)
```
Context: "my-context: 004-lifecycle__implementation"
Target: 60-75 notes covering all 121 tasks
```

**Option B: Multi-Phase Contexts** (recommended for quality tracking)
```
Phase 1: "my-context: 004-lifecycle__implementation-phase1"  # T001-T017, ~20 notes
Phase 2: "my-context: 004-lifecycle__implementation-phase2"  # T018-T051, ~30 notes
Phase 3: "my-context: 004-lifecycle__implementation-phase3"  # T052-T099, ~40 notes
Phase 4: "my-context: 004-lifecycle__testing-integration"    # T100-T121, ~25 notes
```

### Retrospective (FINAL)
```
Context: "my-context: 004-lifecycle__retrospective"
Target: 10-15 notes
Then: Archive after completion
```

**Total contexts for Sprint 004**: 2-5 (vs Sprint 006's 16!) ✅

---

## Integration with .cursor Composer

When Composer runs `/speckit.implement`:

1. ✅ **Auto-starts my-context session** with feature name
2. ✅ **Tracks spec and tasks files** automatically
3. ✅ **Adds notes after each phase** completion
4. ✅ **Monitors note count** - suggests split at 50 notes
5. ✅ **Stops session** when implementation complete

**User just needs to**:
- Let Composer run
- Review my-context notes for progress
- Archive context when sprint merges to master

---

## Quick Reference

### Starting Implementation
```bash
# Automatic (via /speckit.implement in Composer):
my-context start "my-context: {feature}__implementation"

# Manual:
cd /path/to/project
FEATURE=$(basename specs/004-implement-5-lifecycle | cut -d'-' -f1-2)
my-context start "$(basename $(pwd) | sed 's/-dev//'): ${FEATURE}__implementation"
```

### During Implementation
```bash
# After task group
my-context note "T001-T010 complete: {what was built}"

# When creating file
my-context file "internal/commands/newfile.go"

# At checkpoint
my-context note "Checkpoint: {milestone} complete - tests passing"

# If approaching 50 notes
NOTE_COUNT=$(my-context show | grep -oP '(?<=Notes \()\d+')
echo "Current notes: $NOTE_COUNT"  # If >45, consider phase split
```

### Ending Implementation
```bash
# Final summary
my-context note "Implementation complete: all tasks done, SC-001→010 validated"

# Stop session
my-context stop

# Later, after merge
my-context archive "my-context: 004-lifecycle__implementation"
```

---

**Best Practice Summary**:
- ✅ Target 50 notes per context (sweet spot)
- ✅ Use naming convention: `{project}: {feature}__{phase}`
- ✅ Note milestones, decisions, blockers
- ✅ File significant source changes (10+ lines)
- ✅ Split phases proactively (before hitting 75 notes)
- ✅ Stop at phase boundaries
- ✅ Archive after merge to master

---

## Script and Automation Support (v2.2.0+)

### Non-Interactive Mode (CI/CD, Scripts, Automation)

**my-context v2.2.0+ automatically detects TTY** and adapts behavior for scripts:

**In Non-Interactive Mode** (no TTY - scripts, CI/CD):
```bash
# Auto-resumes existing context (no prompt!)
my-context start "existing-context-name"
# → Resumes silently if exists, creates if new

# Auto-suffixes with _2 if duplicate (no prompt!)
my-context start "existing-name" --force
# → Creates existing-name_2 automatically

# All commands work without prompts
my-context note "observation"
my-context stop
```

**In Interactive Mode** (has TTY - terminal):
```bash
# Prompts for resume
my-context start "existing-context"
# → "Context exists. Resume? [Y/n]:"

# Prompts for new name
my-context start "existing" --force
# → "Enter new name:"
```

**For Automation**: Just use commands normally - TTY detection handles it!

**Example deb-sanity automation**:
```bash
#!/bin/bash
# CI pipeline or automated script
my-context start "ci-build-${BUILD_ID}"  # Works without prompts!
my-context note "Build started: ${BUILD_ID}"
run_tests.sh
my-context note "Tests: $TEST_RESULT"
my-context stop
```

---

## Signal Coordination (v2.2.0+)

### Purpose: Event-Driven Coordination Without Polling

**Signal files** = Simple touch files in `~/.my-context/signals/` for event notification

### Use Cases

**1. Agent Coordination**:
```bash
# Agent A: Complete work, signal next agent
my-context note "Implementation complete: all tests passing"
my-context signal create spec-005-implementation-complete
my-context stop

# Agent B: Wait for signal, then start
my-context signal wait spec-005-implementation-complete --timeout=2h
my-context start "spec-005__testing"
my-context note "Starting validation based on Agent A completion"
```

**2. Binary Update Notification**:
```bash
# After building new my-context binary
cd ~/projects/my-context-dev
./scripts/build.sh --install
my-context signal create my-context-binary-updated

# Team members watching (background script):
my-context signal wait my-context-binary-updated && \
  my-context --version && \
  notify-send "my-context updated: $(my-context --version)"
```

**3. Phase Completion Coordination**:
```bash
# In spec kit implementation
my-context note "✅ Phase 1 complete: Core features implemented"
my-context signal create phase-1-complete

# Supervisor or next agent
my-context signal wait phase-1-complete --timeout=4h
# Proceed to phase 2 or testing
```

### Commands

```bash
# Create signal
my-context signal create <signal-name>

# Wait for signal (blocks until signal exists or timeout)
my-context signal wait <signal-name> --timeout=30m

# List active signals
my-context signal list

# Clear signal (cleanup)
my-context signal clear <signal-name>
```

**Signals are lightweight**: Just timestamp files, no daemons, no network

---

## Context Monitoring (v2.2.0+)

### Purpose: Event-Driven Observation (No Constant Polling!)

**Watch command** monitors contexts and executes commands when changes detected

### Supervisor Pattern (Monitor Implementation Agents)

```bash
# Supervisor (you) starts watch in background
my-context watch "deb-sanity: 007-github__implementation" \
  --new-notes \
  --exec="notify-send 'Agent progress update'" &

# Implementation agent works
my-context start "deb-sanity: 007-github__implementation"
my-context note "Checkpoint: Tasks T010-T020 complete"
# → Supervisor gets instant notification!

# Agent signals completion
my-context note "✅ Implementation complete"
# → Supervisor notified again
```

### Pattern Matching

```bash
# Watch for specific patterns in notes
my-context watch "impl-context" \
  --pattern="Phase.*complete" \
  --exec="./start-next-phase.sh"

# Triggers only when notes match pattern
my-context note "Phase 1 complete: Core done"
# → Executes start-next-phase.sh

my-context note "Working on feature X"
# → No trigger (doesn't match pattern)
```

### Timeout and Control

```bash
# With timeout (auto-stop after duration)
my-context watch "context" --new-notes --timeout=2h

# Stop watch with Ctrl+C (graceful)
^C  # Stops watching cleanly
```

**Use Cases**:
- Supervisor monitoring multiple agents
- Automated phase transitions
- Progress notifications
- Integration test triggers

---

## Context Metadata & Ownership (v2.2.0+)

### Tracking Who Created What

**Metadata fields** in every context (v2.2.0+):
- **created-by**: Who/what created this context
- **parent**: Parent context (for hierarchies)
- **labels**: Tags for categorization

### Setting Metadata on Creation

```bash
# Tool attribution
my-context start "feature-work" --created-by="deb-sanity"

# Parent relationship
my-context start "phase-2" --parent="main-sprint"

# Labels for categorization
my-context start "task" --labels="bugfix,urgent,backend"

# All combined
my-context start "deb-sanity: 007-github__impl-phase1" \
  --created-by="speckit-implement" \
  --parent="deb-sanity: 007-github__planning" \
  --labels="implementation,github,phase1"
```

### Viewing Metadata

```bash
# Context with metadata shows:
my-context show

# Output includes:
Metadata:
  Created by: speckit-implement
  Parent: deb-sanity: 007-github__planning
  Labels: implementation, github, phase1
```

### Querying by Metadata (Future v2.3.0)

```bash
# Find all contexts created by tool
my-context list --created-by="deb-sanity"

# Find by label
my-context list --label="bugfix"

# Find children of parent
my-context list --parent="sprint-007"
```

---

## Context Hierarchies & Relationships

### Parent-Child Pattern (Top-Level Sprint → Phases)

**Example**: deb-sanity Sprint 007 structure

```
deb-sanity: 007-github__planning (TOP LEVEL)
  ├── deb-sanity: 007-github__impl-phase1 (child)
  ├── deb-sanity: 007-github__impl-phase2 (child)
  ├── deb-sanity: 007-github__impl-phase3 (child)
  └── deb-sanity: 007-github__validation (child)
```

**Create hierarchy**:
```bash
# 1. Top-level sprint planning
my-context start "deb-sanity: 007-github__planning"
my-context note "Spec complete: 4 user stories, 56 tasks"
my-context file "specs/007-github/spec.md"
my-context stop

# 2. Implementation phase 1 (child)
my-context start "deb-sanity: 007-github__impl-phase1" \
  --parent="deb-sanity: 007-github__planning" \
  --labels="implementation,phase1"
my-context note "Parent: 007-github__planning"
my-context note "Starting Phase 1: Pre-flight validation (T001-T020)"
# ... 50 notes of implementation ...
my-context stop

# 3. Implementation phase 2 (child)
my-context start "deb-sanity: 007-github__impl-phase2" \
  --parent="deb-sanity: 007-github__planning" \
  --labels="implementation,phase2"
my-context note "Sibling: 007-github__impl-phase1 (Phase 1)"
# ... implementation continues ...
```

**Reference Patterns**:
- **By name**: `"Parent: deb-sanity: 007-github__planning"` in notes
- **By path**: `my-context file "/path/to/parent/spec.md"` (cross-context file ref)
- **By metadata**: `--parent="context-name"` flag
- **By label**: `--labels="sprint-007,github"` for grouping

---

## Spec Kit Integration (Quality Gates Using Contexts)

### How Spec Kit Commands Use my-context

**Each spec kit command = ceremony with context tracking**:

**`/specify` Ceremony**:
```bash
# Auto-starts context
my-context start "{project}: spec-NNN__planning"
my-context note "🚀 Stage 0: draft - Running /specify"
my-context file "specs/spec-NNN/spec.md"
my-context note "✅ Spec generated: 5 user stories, 24 requirements"
# Context continues for /clarify, /plan, /tasks
```

**`/clarify` Ceremony**:
```bash
my-context note "🚀 Stage 1: clarified - Answering questions"
my-context note "Q1: Decision B selected (rationale: ...)"
my-context file "specs/spec-NNN/clarify.md"
my-context note "✅ All clarifications resolved: 0 NEEDS CLARIFICATION remaining"
```

**`/plan` Ceremony**:
```bash
my-context note "🚀 Stage 2: planned - Creating implementation plan"
my-context note "Architecture: Event sourcing + file watching"
my-context file "specs/spec-NNN/plan.md"
my-context note "✅ Plan complete: 3 phases, 5 days estimated"
```

**`/tasks` Ceremony**:
```bash
my-context note "🚀 Stage 3: tasked - Generating task list"
my-context note "Generated 40 tasks across 6 phases"
my-context file "specs/spec-NNN/tasks.md"
my-context note "✅ Tasks ready: 40 numbered, dependencies clear"
my-context stop  # End planning, start implementation
```

**`/implement` Ceremony**:
```bash
# New context for implementation
my-context start "{project}: spec-NNN__impl-phase1" \
  --parent="{project}: spec-NNN__planning" \
  --created-by="speckit-implement" \
  --labels="implementation,phase1"

my-context note "🚀 Stage 5: implementing - Phase 1 starting"
my-context touch  # Heartbeat every 10 minutes
my-context note "📊 Checkpoint: T001-T015 complete"
my-context file "internal/commands/signal.go"
my-context note "✅ Phase 1 complete: Signal commands working"
my-context signal create phase-1-complete  # Coordinate!
my-context stop
```

---

## Quality Gates Using Context Audit Trail

### Gate 1: Planning Completeness

**Check**: Did spec kit ceremonies complete?
```bash
# Verify planning context has all artifacts
my-context show "spec-NNN__planning" | grep "Files ("
# Should show: spec.md, clarify.md, plan.md, tasks.md

# Check for quality notes
my-context show "spec-NNN__planning" | grep "✅.*complete"
# Should have: Spec complete, Clarify complete, Plan complete, Tasks complete
```

**Gate**: If planning context <20 notes or missing artifacts → incomplete planning

---

### Gate 2: Implementation Progress

**Check**: Are checkpoints documented?
```bash
# Count checkpoint notes
my-context show "spec-NNN__impl-phase1" | grep "Checkpoint" | wc -l
# Should have 1 per 10-15 tasks

# Check heartbeats
my-context show "spec-NNN__impl-phase1" | grep "Activity:"
# Should show touches (heartbeats every 10min)
```

**Gate**: If no checkpoints or >30min without activity → review needed

---

### Gate 3: Test Coverage

**Check**: Are tests documented?
```bash
# Look for test completion notes
my-context show "spec-NNN__validation" | grep -i "test.*pass"

# Check for bug notes
my-context show | grep -i "bug\|issue\|fail"
```

**Gate**: Must have test results documented before merge

---

### Gate 4: Constitution Compliance

**Check**: Was constitution validated?
```bash
# Look for constitution check note
my-context show "sprint-integration-v2.2.0" | grep -i "constitution"

# Should see: "Constitution validated: {principles checked}"
```

**Gate**: Constitution validation required before release

---

## Project Lifecycle Journaling: All Participants

### Who Participates in the Journal

**Humans**:
- Developers (implementation notes)
- Reviewers (code review notes)  
- QA testers (validation notes)
- Project managers (status updates)

**AI Agents**:
- Spec writers (spec kit /specify)
- Implementation agents (/implement with Composer)
- Testing agents (automated test runs)
- Analysis agents (/analyze quality checks)

**Tools**:
- deb-sanity (auto-creates contexts with --created-by="deb-sanity")
- Spec kit commands (--created-by="speckit-{command}")
- Build systems (--created-by="ci-pipeline")

**All use same ~/.my-context/ folder** = unified project journal!

---

### Cross-Tool Coordination Example

**Real Pattern** (from final-completion handoff):

```
final-completion (top-level communication context)
  ├── created-by: "manual" (my-context team started it)
  ├── 114 notes from:
  │   ├── my-context team (implementation, fixes, vision)
  │   ├── deb-sanity team (validation, feedback, testing)
  │   └── Composer agent (automated implementation)
  ├── 15 files tracked:
  │   ├── my-context project files (specs, docs)
  │   └── deb-sanity analysis files (cross-project!)
  └── Duration: 17+ hours async communication

Pattern: Multiple participants, one journal, complete audit trail
```

**How it worked**:
1. my-context team: `my-context note "HANDOFF TO DEB-SANITY..."`
2. deb-sanity team: `my-context resume final-completion && my-context note "RECEIVED..."`
3. Back and forth: 114 notes exchanged
4. Result: Coordinated releases, bug fixes, roadmap alignment

---

### Spec Kit Quality Guard Integration

**Pattern**: Each ceremony adds quality gate via context

```bash
# Stage 0: draft
/specify "feature description"
# Creates context: {project}: spec-NNN__planning
# Gate: Must have spec.md with user stories

# Stage 1: clarified  
/clarify
# Adds to context: clarify.md
# Gate: Zero [NEEDS CLARIFICATION] markers

# Stage 2: planned
/plan
# Adds to context: plan.md
# Gate: Architecture + phases documented

# Stage 3: tasked
/tasks
# Adds to context: tasks.md
# Gate: All tasks numbered, files specified

# Stage 4: checklist (Cursor)
/checklist
# Adds to context: checklists/requirements.md  
# Gate: All checklist items [X] passing

# Stage 5: implementing
/implement
# Creates NEW context: spec-NNN__impl-phaseN
# Gate: Checkpoints every 10-15 tasks, heartbeats every 10min
```

**Audit Trail**: Every context = proof of gate completion

---

## Naming Convention (CRITICAL - Don't Ignore!)

### Mandatory Pattern: `{project}: {feature}__{phase}`

**Why project name FIRST**:
- Quick filtering: `my-context list | grep "^deb-sanity"`
- Tool attribution visible immediately
- Cross-project journal organization

**Examples** (Good ✅):
```
deb-sanity: 007-github__spec-plan-tasks
my-context: 005-signaling__impl-phase1
ps-cli: 003-package__validation
```

**Examples** (Bad ❌ - Will cause confusion):
```
spec-007 (missing project!)
implementation (ambiguous!)
007-sprint-007-github_3 (no project, numbered suffix = fragmentation!)
```

### With Metadata (v2.2.0+)

```bash
my-context start "deb-sanity: 007-github__impl-phase1" \
  --created-by="speckit-implement" \
  --parent="deb-sanity: 007-github__planning" \
  --labels="implementation,github,phase1"
```

**Result**: Name is human-readable, metadata is queryable

---

## Context Hierarchy Visualization

### Example: Complete Sprint with Phases

```
deb-sanity: 007-github-publishing (TOP LEVEL SPRINT)
│ created-by: "manual"
│ labels: ["sprint-007", "github", "publishing"]
│ notes: 30 (planning phase)
│ files: spec.md, clarify.md, plan.md, tasks.md
│
├── deb-sanity: 007-github__impl-phase1 (CHILD - Pre-flight)
│   │ parent: "deb-sanity: 007-github-publishing"
│   │ created-by: "speckit-implement"
│   │ labels: ["implementation", "phase1", "preflight"]
│   │ notes: 50 (implementation phase)
│   │ files: deb-sanity.sh, github-preflight.sh
│   │ signal: phase-1-complete (when done)
│   │
│   └── Status: stopped (complete)
│
├── deb-sanity: 007-github__impl-phase2 (CHILD - Prepare Release)
│   │ parent: "deb-sanity: 007-github-publishing"
│   │ sibling: "007-github__impl-phase1"
│   │ notes: 45
│   │ signal: phase-2-complete
│   │
│   └── Status: stopped (complete)
│
├── deb-sanity: 007-github__validation (CHILD - Testing)
│   │ parent: "deb-sanity: 007-github-publishing"
│   │ notes: 40 (UAT phase)
│   │ files: test results, bug reports
│   │
│   └── Status: stopped (UAT done)
│
└── deb-sanity: 007-github__retro (CHILD - Retrospective)
    │ parent: "deb-sanity: 007-github-publishing"
    │ notes: 15 (lessons learned)
    │
    └── Status: stopped → READY TO ARCHIVE ENTIRE HIERARCHY
```

**Query hierarchy**:
```bash
# Find all children
my-context list | grep "deb-sanity: 007-github"

# Or with metadata (v2.3.0):
my-context list --parent="deb-sanity: 007-github-publishing"
```

---

## Cross-Context References

### Patterns for Linking Contexts

**1. By Name** (in notes):
```bash
my-context note "Parent context: deb-sanity: 007-github__planning"
my-context note "Sibling: 007-github__impl-phase1 (Phase 1 complete)"
my-context note "Blocked by: final-completion (waiting for my-context team response)"
```

**2. By File Reference** (shared files):
```bash
# Track files from other contexts
my-context file "/home/be-dev-agent/projects/deb-sanity/specs/007-github/spec.md"
my-context file "/home/be-dev-agent/projects/my-context-dev/docs/SIGNALING-PROTOCOL.md"

# Cross-project references work!
```

**3. By Metadata** (parent field):
```bash
my-context start "child-context" --parent="parent-context-name"
# Establishes formal relationship
```

**4. By Signal** (event coordination):
```bash
my-context note "Waiting for signal: spec-005-implementation-complete"
my-context signal wait spec-005-implementation-complete
my-context note "Signal received: Parent implementation complete, starting validation"
```

---

## Approved Journaling Patterns (All Participants)

### Pattern 1: Tool-Initiated Contexts

**When deb-sanity creates worktree**:
```bash
# deb-sanity.sh automatically does:
my-context start "deb-sanity: ${FEATURE}__worktree-setup" \
  --created-by="deb-sanity" \
  --labels="worktree,setup"
my-context note "Worktree created: ${BRANCH_NAME}"
my-context file "${WORKTREE_PATH}/README.md"
```

**Approved**: Tools can auto-create contexts with proper attribution

---

### Pattern 2: Spec Kit Workflow

**Automatic context management**:
```bash
# /specify creates planning context
my-context start "{project}: spec-NNN__planning" \
  --created-by="speckit-specify"

# /implement creates implementation context(s)
my-context start "{project}: spec-NNN__impl-phase1" \
  --parent="{project}: spec-NNN__planning" \
  --created-by="speckit-implement"
```

**Approved**: Spec kit commands manage contexts automatically

---

### Pattern 3: Human + AI Collaboration

**Same context, multiple contributors**:
```bash
# Human starts work
my-context start "feature-analysis"
my-context note "Investigating performance issue"

# AI agent adds analysis
my-context note "Analysis: Bottleneck in loop at line 245"
my-context file "src/performance-hotspot.go"

# Human resumes
my-context note "Applied AI suggestion: 40% performance improvement"

# Both contribute to same audit trail!
```

**Approved**: Humans and AI share contexts freely

---

### Pattern 4: Cross-Team Handoff

**Async communication via shared context** (proven pattern!):
```bash
# Team A: Handoff
my-context start "cross-team-handoff"
my-context note "🤝 HANDOFF TO TEAM B: Feature complete, needs validation"
my-context file "deliverable.md"
my-context stop

# Team B: Receive
my-context resume "cross-team-handoff"
my-context note "🤝 RECEIVED: Starting validation"
my-context note "✅ Validation complete: Ready to merge"
my-context stop

# Team A: Response
my-context resume "cross-team-handoff"
my-context note "✅ MERGED: Thank you for validation"
my-context signal create handoff-complete
```

**Approved**: Contexts enable async team coordination

---

## Complete Audit Trail Example (Real Usage)

### Iteration 004 (Lifecycle Improvements) Full Trail

**Planning** (1 context):
```
my-context: 004-lifecycle__spec-plan-tasks
├── Duration: 15 minutes
├── Notes: 7 (spec → clarify → plan → tasks)
├── Files: spec.md, clarify.md, plan.md, tasks.md
└── Created-by: manual
```

**Implementation** (3 contexts - multi-phase):
```
my-context: 004-lifecycle__impl-phase1 (Smart Resume)
├── Parent: 004-lifecycle__spec-plan-tasks
├── Duration: 2 days
├── Notes: 50 (checkpoints, heartbeats, completion)
├── Files: start.go, state.go, tests
└── Signal: phase-1-complete

my-context: 004-lifecycle__impl-phase2 (Warnings + Resume)
├── Parent: 004-lifecycle__spec-plan-tasks
├── Duration: 1 day
├── Notes: 45
└── Signal: phase-2-complete

my-context: 004-lifecycle__impl-phase3 (Bulk + Advisor)
├── Parent: 004-lifecycle__spec-plan-tasks
├── Duration: 1.5 days
├── Notes: 40
└── Signal: phase-3-complete
```

**Testing** (1 context):
```
my-context: 004-lifecycle__validation
├── Parent: 004-lifecycle__spec-plan-tasks
├── Duration: 1 day
├── Notes: 40 (UAT results, 10/10 tests passed)
└── Files: test reports, bug fixes
```

**Sprint Integration** (1 context):
```
my-context: sprint-integration-v2.0.0
├── References: All 004-lifecycle contexts
├── Notes: 15 (merge, test, tag, release)
└── Status: archived (release complete)
```

**Total**: 6 contexts, ~200 notes, complete audit trail from idea → release

**Queryable**: Can reconstruct entire iteration timeline from contexts!

---

## Quick Reference Card (Print & Pin!)

### Context Creation (Always Include Project!)

```bash
# Planning
my-context start "{project}: {feature}__planning"

# Implementation  
my-context start "{project}: {feature}__impl-phase1" \
  --parent="{project}: {feature}__planning"

# With tool attribution
my-context start "{project}: {feature}__impl" \
  --created-by="speckit-implement" \
  --labels="implementation"
```

### During Work

```bash
my-context note "📊 Checkpoint: {what completed}"
my-context file "{path/to/significant/file}"
my-context touch  # Heartbeat (every 10min in long sessions)
```

### Coordination

```bash
my-context signal create {event-name}
my-context signal wait {event-name} --timeout=30m
my-context watch "{context}" --new-notes --exec="./notify.sh"
```

### End of Phase

```bash
my-context note "✅ Phase complete: {summary}"
my-context signal create phase-N-complete
my-context stop
```

---

**Version**: 2.0.0 → **3.0.0** (Updated for v2.2.0-beta features)
**Date**: 2025-10-11 (Updated with signal/watch/metadata)
**Purpose**: Complete guide for project lifecycle journaling

