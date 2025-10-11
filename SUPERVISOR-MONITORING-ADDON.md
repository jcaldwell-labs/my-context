# Supervisor Monitoring Addon for /speckit.implement

**Purpose**: Add to implementation agent prompt to enable supervisor monitoring via signaling
**Usage**: Append this to your agent launch prompt when calling /speckit.implement
**Context**: Enables real-time monitoring of implementation progress

---

## Addon Prompt Text

```
SUPERVISOR MONITORING PROTOCOL:

You are being monitored by a supervisor process. Follow these signaling protocols:

1. **Heartbeat Signals** (Required):
   - Run `my-context touch` every 10 minutes during active work
   - Purpose: Proves you're not stuck/blocked
   - Supervisor watches for >30min gaps (triggers intervention check)

2. **Phase Boundary Signals** (Required):
   - BEFORE starting each phase:
     `my-context note "🚀 Starting Phase {N}: {Phase Name} (Tasks {start}-{end})"`
   - AFTER completing each phase:
     `my-context note "✅ Phase {N} complete: {summary} ({completed}/{total} tasks)"`
   - Purpose: Supervisor tracks progress and can pause/resume between phases

3. **Milestone Signals** (Important):
   - After completing significant task groups (every 10-15 tasks):
     `my-context note "📊 Checkpoint: {milestone description} (Tasks {X}-{Y} complete)"`
   - When blocked:
     `my-context note "🚧 BLOCKED: {blocker description} - investigating"`
   - When unblocked:
     `my-context note "✅ UNBLOCKED: {resolution} - continuing"`
   - Purpose: Supervisor can intervene on blockers, celebrate milestones

4. **File Tracking Signals** (Recommended):
   - When creating new source files:
     `my-context file "{absolute/path/to/file.go}"`
   - Track significant files only (>10 lines or core functionality)
   - Purpose: Supervisor can review file changes, understand scope

5. **Completion Signal** (Required):
   - When ALL tasks complete:
     `my-context note "🎉 Implementation complete: All {N} tasks done, tests passing"`
   - Before stopping:
     `my-context note "✅ Quality gate passed: Implementation complete and validated"`
   - Then: `my-context stop`
   - Purpose: Supervisor knows implementation finished, can proceed to next phase

6. **Error/Failure Protocol** (Critical):
   - On test failures:
     `my-context note "❌ Tests failing: {test names} - investigating"`
   - On build errors:
     `my-context note "🔧 Build error: {error summary} - fixing"`
   - On implementation stuck:
     `my-context note "🤔 Need guidance: {question/issue} - awaiting supervisor input"`
   - Purpose: Supervisor can provide help or adjust approach

SUPERVISOR MONITORING COMMANDS (What supervisor will run):

- `my-context show` - Check your progress, latest notes
- `my-context show | grep "Notes ("` - Check note count (watch for 50-note warning)
- `my-context show | tail -20` - See recent 20 notes (progress updates)
- `my-context show | grep -E "Phase.*complete|BLOCKED|ERROR"` - Check milestones/issues

Your notes become the communication channel with the supervisor. Be verbose - notes are cheap, silence is expensive!
```

---

## How Supervisor Uses This

### Supervisor Monitoring Script (Runs in parallel)

```bash
#!/bin/bash
# monitor-implementation.sh

CONTEXT_NAME="deb-sanity: 007-github-publishing__implementation"
CHECK_INTERVAL=60  # Check every 60 seconds

echo "👁️ Monitoring context: $CONTEXT_NAME"

while true; do
    # Get latest note
    LATEST_NOTE=$(my-context show "$CONTEXT_NAME" 2>/dev/null | grep "^\s*\[" | tail -1 || echo "")

    # Check for completion
    if echo "$LATEST_NOTE" | grep -q "Implementation complete\|Quality gate passed"; then
        echo "✅ Implementation complete!"
        notify-send "Agent finished: $CONTEXT_NAME"
        break
    fi

    # Check for blockers
    if echo "$LATEST_NOTE" | grep -q "BLOCKED\|ERROR\|Need guidance"; then
        echo "🚧 Agent blocked - review needed"
        notify-send "Agent needs help: Check my-context show"
        # Optional: Pause monitoring, wait for manual intervention
    fi

    # Check for phase completion
    if echo "$LATEST_NOTE" | grep -q "Phase.*complete"; then
        PHASE=$(echo "$LATEST_NOTE" | grep -oP 'Phase \d+')
        echo "✅ $PHASE complete"
    fi

    # Check for heartbeat (last note/touch within 30min)
    LAST_ACTIVITY=$(my-context show "$CONTEXT_NAME" 2>/dev/null | grep "Activity:" | grep -oP '\d+ touches' || echo "0")
    # If no activity in 30min, flag for review

    sleep $CHECK_INTERVAL
done
```

### Supervisor Dashboard (Simple)

```bash
# Quick status check
watch -n 10 'my-context show | head -20'

# Or extract key info
my-context show | grep -E "Status:|Notes \(|Phase|complete|BLOCKED" | head -20
```

---

## Enhanced Signaling (When Sprint 006 Ships)

**After FR-MC-NEW-007 (Context Signaling) implemented**:

```bash
# Agent creates signals at milestones
my-context signal create phase-1-complete
my-context signal create implementation-complete
my-context signal create tests-passing

# Supervisor waits for signals
my-context signal wait phase-1-complete --timeout=2h && echo "Phase 1 done, reviewing..."
my-context signal wait implementation-complete --timeout=8h && start-validation.sh

# Or watch for note patterns
my-context watch "$CONTEXT_NAME" --pattern="Phase.*complete" --exec="./notify-phase-done.sh"
my-context watch "$CONTEXT_NAME" --pattern="BLOCKED" --exec="./alert-supervisor.sh"
```

---

## Example: Launching Implementation Agent with Monitoring

### Step 1: Supervisor Launches Agent

```bash
# In terminal 1 (supervisor - you)
cd ~/projects/deb-sanity

# Start monitoring script in background
./monitor-implementation.sh &
MONITOR_PID=$!

# Launch implementation agent with monitoring protocol
# Add the SUPERVISOR MONITORING PROTOCOL text above to agent prompt
```

### Step 2: Agent Executes with Signaling

Agent automatically:
- ✅ Uses existing context or starts implementation-specific one
- ✅ Adds notes at phase boundaries
- ✅ Touches every 10 minutes (heartbeat)
- ✅ Tracks files as created
- ✅ Signals completion

### Step 3: Supervisor Monitors

You see in real-time:
```
👁️ Monitoring context: deb-sanity: 007-github-publishing__implementation
✅ Phase 1 complete
✅ Phase 2 complete
🚧 Agent blocked - review needed
  (You check notes, provide guidance)
✅ Phase 3 complete
✅ Implementation complete!
```

### Step 4: Handoff

When complete:
```bash
# Agent stops
my-context stop

# Supervisor checks results
my-context show  # See full implementation log
my-context show | grep "✅.*complete"  # See milestones

# Supervisor archives
my-context archive "deb-sanity: 007-github-publishing__implementation"
```

---

## Key Signaling Patterns for Implementation

### Pattern 1: Regular Heartbeat
**Every 10 minutes**: `my-context touch`
**Supervisor detects**: Agent is alive, making progress
**Triggers intervention if**: >30min gap without note/touch

### Pattern 2: Phase Boundaries
**Phase start**: `my-context note "🚀 Starting Phase {N}: {Name}"`
**Phase end**: `my-context note "✅ Phase {N} complete: {summary}"`
**Supervisor uses**: Track which phase agent is in, predict completion time

### Pattern 3: Checkpoint Progress
**Every 10-15 tasks**: `my-context note "📊 Checkpoint: {milestone}"`
**Supervisor uses**: Verify steady progress, estimate time remaining

### Pattern 4: Block/Unblock Events
**When stuck**: `my-context note "🚧 BLOCKED: {issue}"`
**When resolved**: `my-context note "✅ UNBLOCKED: {resolution}"`
**Supervisor uses**: Provide help, adjust timeline, escalate if needed

### Pattern 5: Completion Signal
**Final**: `my-context note "✅ Quality gate passed: Implementation complete"`
**Supervisor uses**: Know implementation finished, proceed to validation

---

## Configuration for deb-sanity

**Already have** (in .specify/scripts/bash/context-tracking.sh):
- ✅ `context_note` function
- ✅ `context_file` function
- ✅ `context_start` function
- ✅ `my-context touch` integration
- ✅ Journaling policy support

**Just need to add** (for signaling):
- When Sprint 006 ships: Signal creation at milestones
- Watch command for supervisor monitoring
- Handoff protocol for team coordination

**Current workaround** (until Sprint 006):
- Supervisor manually runs `my-context show` every few minutes
- Or use `watch` command: `watch -n 60 'my-context show | tail -20'`

---

**Addon Version**: 1.0.0
**Compatible With**: deb-sanity speckit.implement.md (current version)
**Enhanced By**: Sprint 006 signaling protocol (future)
