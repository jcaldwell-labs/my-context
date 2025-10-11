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
