# Context Signaling Protocol Requirements

**Date**: 2025-10-10
**Status**: Requirements from real-world handoff pattern
**Use Case**: my-context ↔ deb-sanity team communication via "final-completion" context

---

## Real-World Pattern (What We Just Did)

### Scenario: Team Handoff via Shared Context

**Participants**:
- my-context team (us)
- deb-sanity team (them)
- Shared resources: `~/.local/bin/my-context` binary, my-context contexts

**Flow**:
1. ✅ We complete Sprint 004 implementation (all 121 tasks)
2. ✅ We create handoff context "final-completion" with notes + files
3. ✅ **We manually nudge**: "Check my-context show" (Slack/email)
4. ✅ They check, test features, add 20 notes of feedback
5. ✅ We check context, see feedback, fix bugs
6. ✅ We add response notes (3 notes)
7. ✅ We rebuild binary (bdac191)
8. ❌ **Gap**: They don't know binary updated or new notes available
9. ❌ **Manual nudge required again**: "Pull latest, rebuild, check context"

**52 notes exchanged** in async communication - works great! But requires manual notifications.

---

## Signaling Gaps Identified

### Gap 1: Binary Update Notification
**Problem**: We updated binary (be60865 → bdac191), they don't know

**Current Workaround**: Manual Slack message "Binary updated, please check version"

**Ideal**: Automated notification when `~/.local/bin/my-context` changes

---

### Gap 2: New Context Notes Notification
**Problem**: We added notes #31-36 in final-completion, they don't know

**Current Workaround**: Manual nudge "Check my-context show for our response"

**Ideal**: They get notification when new notes added to watched context

---

### Gap 3: Handoff State Changes
**Problem**: Handoff state unclear (are we waiting? did they respond? is it resolved?)

**Current Workaround**: Infer from note patterns ("HANDOFF TO", "FEEDBACK FROM")

**Ideal**: Explicit handoff state machine (pending → in-progress → responded → resolved)

---

### Gap 4: Polling Overhead
**Problem**: If they want automation, must poll `my-context show` repeatedly

**Current Workaround**: Cron job checking every N minutes (wasteful)

**Ideal**: Event-driven notification (only checks when change occurs)

---

### Gap 5: Multi-Context Monitoring
**Problem**: Supervisor wants to monitor multiple contexts (implementation agents, testing agents)

**Current Workaround**: Poll each context individually in loop

**Ideal**: Single watch command monitoring multiple contexts simultaneously

---

## Requirements for Signaling Protocol

### REQ-001: Lightweight Implementation
**Must**: File-based, no daemon processes, no network dependencies
**Rationale**: Simple, reliable, works offline, easy to debug

### REQ-002: Backward Compatible
**Must**: Works without signaling (current manual workflow still viable)
**Rationale**: Signaling is enhancement, not requirement

### REQ-003: Event-Driven Notifications
**Must**: Trigger on changes (not constant polling)
**Rationale**: Efficient, timely, reduces overhead

### REQ-004: Multiple Signal Types
**Must**: Support binary updates, context notes, handoff states, custom signals
**Rationale**: Different use cases need different triggers

### REQ-005: Shared Resource Safety
**Must**: Handle concurrent access to `~/.local/bin/my-context` and contexts
**Rationale**: Avoid race conditions, corrupted state

---

## Proposed Solutions (Ranked by Simplicity)

### Solution 1: File Modification Time Monitoring (Simplest)
**Concept**: Watch file mtime, trigger when changed

```bash
# deb-sanity team runs:
my-context watch binary --exec="echo 'Binary updated!' | notify-send"
# Internally: Polls ~/.local/bin/my-context mtime every 5s, executes on change

my-context watch final-completion --new-notes --exec="./check-handoff.sh"
# Polls ~/.my-context/final-completion/notes.log mtime, executes when modified
```

**Pros**: Dead simple, works today with inotify/stat
**Cons**: Polling overhead (5s intervals), not truly event-driven

---

### Solution 2: Signal Files (File-Based Semaphores)
**Concept**: Touch files to signal events

```bash
# my-context team (after updating binary):
my-context signal create binary-updated-bdac191
# Creates: ~/.my-context/signals/binary-updated-bdac191.signal

# deb-sanity team (background process):
my-context signal wait binary-updated --timeout=1h
# Blocks until ~/.my-context/signals/binary-updated*.signal appears
echo "Binary updated! New version available."

# Or polling mode:
my-context signal list | grep binary-updated && handle_update
```

**Pros**: Explicit, easy to debug (ls signals/), supports multiple signal types
**Cons**: Requires manual signal creation, cleanup needed

---

### Solution 3: Context Metadata with State Machine
**Concept**: Structured handoff states in meta.json

```json
// ~/.my-context/final-completion/meta.json
{
  "handoff": {
    "to": "deb-sanity-team",
    "from": "my-context-team",
    "status": "awaiting-feedback",  // pending|in-progress|responded|resolved
    "last_updated": "2025-10-10T19:42:00Z",
    "updates": [
      {
        "timestamp": "2025-10-10T19:35:00Z",
        "type": "bugfix-applied",
        "notes": "Notes #31-33 added",
        "binary_commit": "bdac191"
      }
    ]
  }
}
```

Commands:
```bash
# my-context team:
my-context handoff final-completion --to=deb-sanity --status=awaiting-feedback
my-context handoff-update final-completion --type=bugfix-applied --binary=bdac191

# deb-sanity team:
my-context handoff-status  # Shows: final-completion awaiting-feedback (1 update)
my-context handoff-accept final-completion  # Changes status to in-progress
```

**Pros**: Structured, queryable, state machine clear
**Cons**: More complex, requires metadata management

---

### Solution 4: Hybrid (Recommended)
**Combine Solutions 1 + 2 for best of both worlds**

**For binary updates** (Solution 2 - Signal files):
```bash
# After installing new binary
my-context signal create binary-updated

# deb-sanity background script
while my-context signal wait binary-updated --timeout=infinite; do
  VERSION=$(my-context --version)
  notify-send "my-context updated: $VERSION"
  my-context signal clear binary-updated
done
```

**For context updates** (Solution 1 - File mtime):
```bash
# deb-sanity watches handoff context
my-context watch final-completion --new-notes --exec="notify-send 'New handoff notes available'"
# Polls notes.log mtime, executes when changed
```

**For handoff states** (Solution 3 - Metadata, optional):
```bash
my-context handoff final-completion --to=deb-sanity --status=awaiting-feedback
# They check: my-context handoff-status
```

---

## Implementation Priorities for Sprint 006

### MVP (Must Have)
1. **Signal files**: create/list/wait/clear commands
2. **Watch command**: monitor file/context changes, execute on change
3. **Binary update pattern**: Specific support for `~/.local/bin/my-context` monitoring

### Nice to Have
4. **Handoff metadata**: Structured state in meta.json
5. **Multi-context watch**: Monitor multiple contexts simultaneously
6. **Signal patterns**: Glob matching for signal wait

### Can Defer
7. **Event sourcing**: Full change log in contexts
8. **Web hooks**: HTTP POST on events (too heavy)
9. **Desktop notifications**: Built-in notify-send integration

---

## Concrete Examples from Our Handoff

### Example 1: Binary Update Signal

**What we did** (manual):
```bash
# my-context team
cd ~/projects/my-context-dev
go build ...
cp my-context ~/.local/bin/
# Then: Slack message "Binary updated to bdac191"
```

**What we want** (automated):
```bash
# my-context team
cd ~/projects/my-context-dev
go build ...
cp my-context ~/.local/bin/
my-context signal create binary-updated  # ← NEW

# deb-sanity team (runs in background)
my-context signal wait binary-updated && {
  my-context --version
  echo "New version: $(my-context --version)" | notify-send
}
```

---

### Example 2: New Notes in Handoff Context

**What we did** (manual):
```bash
# my-context team
my-context note "🔧 BUGFIX APPLIED: ..."
my-context note "✅ SHORTCUT CONFIRMED: ..."
# Then: Slack message "Added response notes, check context"
```

**What we want** (automated):
```bash
# my-context team
my-context note "🔧 BUGFIX APPLIED: ..."
# Automatically creates signal or updates mtime

# deb-sanity team (background daemon)
my-context watch final-completion --new-notes --exec="notify-send 'my-context team responded'"
# Detects new notes, triggers notification
```

---

### Example 3: Handoff State Tracking

**What we did** (manual, inferred from notes):
```
Note #6:  "HANDOFF TO DEB-SANITY TEAM"        → Status: pending
Note #11: "RECEIVED: deb-sanity team reviewing" → Status: in-progress
Note #20: "FEEDBACK FROM deb-sanity TEAM"      → Status: responded
Note #31: "BUGFIX APPLIED"                     → Status: iterating
Note #??:  "VALIDATED - MERGE TO MASTER"       → Status: resolved
```

**What we want** (structured):
```bash
# my-context team
my-context handoff final-completion --to=deb-sanity --status=pending
# Creates metadata, can query later

# deb-sanity team
my-context handoff-status  # Shows: final-completion → deb-sanity (pending)
my-context handoff-accept final-completion --status=in-progress
my-context handoff-update final-completion --status=responded

# Query
my-context handoff-status
# Output:
#   final-completion → deb-sanity: responded (4 updates since handoff)
```

---

## Technical Design Considerations

### File Locations

**Signals**:
```
~/.my-context/signals/
├── binary-updated.signal           # General signal
├── binary-updated-bdac191.signal   # Versioned signal
├── handoff-final-completion.signal # Context-specific
└── implementation-phase1-complete.signal
```

**Metadata**:
```
~/.my-context/final-completion/
├── meta.json          # Context metadata (includes handoff state)
├── notes.log          # Existing
├── files.log          # Existing
└── .last-modified     # Touch file for mtime monitoring
```

### Watch Implementation

**Polling-based** (simple, works everywhere):
```go
func WatchContext(name string, interval time.Duration, pattern string, exec string) {
    lastMtime := getNotesLogMtime(name)
    for {
        time.Sleep(interval)
        currentMtime := getNotesLogMtime(name)
        if currentMtime.After(lastMtime) {
            // New notes detected
            if pattern == "" || notesMatchPattern(name, pattern) {
                executeCommand(exec)
                lastMtime = currentMtime
            }
        }
    }
}
```

**Event-driven** (advanced, Linux-only):
```go
// Use inotify on Linux for zero-overhead event detection
watcher, _ := fsnotify.NewWatcher()
watcher.Add(contextDir + "/notes.log")
for event := range watcher.Events {
    if event.Op&fsnotify.Write != 0 {
        executeCommand(exec)
    }
}
```

**Recommendation**: Start with polling (cross-platform), add inotify optimization later

---

## Integration with Existing Features

### Works With Smart Resume (FR-MC-NEW-001)
```bash
# Signal when context ready to resume
my-context signal create implementation-complete

# Another process waits and resumes
my-context signal wait implementation-complete && my-context resume --last
```

### Works With Lifecycle Advisor (FR-MC-NEW-002)
```bash
# On stop, check for pending handoffs
my-context stop
# Lifecycle advisor could show: "Pending handoff to deb-sanity (awaiting feedback)"
```

### Works With Bulk Archive (FR-MC-NEW-003)
```bash
# After sprint complete, signal ready for cleanup
my-context signal create sprint-004-complete

# Another process archives all sprint contexts
my-context signal wait sprint-004-complete && \
  my-context archive --pattern "004-*" --yes
```

---

## Success Criteria (Draft)

**SC-007-001**: Binary update detected within 10 seconds of install
**SC-007-002**: Context note changes detected within 5 seconds
**SC-007-003**: Watch command CPU usage <1% (polling mode)
**SC-007-004**: Signals support concurrent readers/writers safely
**SC-007-005**: Works on Linux, macOS, Windows (cross-platform)

---

## Next Steps for Sprint 006

1. **Validate pattern**: Use watch/signals during Sprint 005 execution
2. **Answer questions**: 5 questions from outline (signal scope, polling interval, etc.)
3. **Create full spec**: Convert outline → complete spec with acceptance scenarios
4. **Implement**: Signals first (2 days), watch second (2-3 days), handoff optional (1-2 days)

---

## Immediate Use Case for Sprint 005

**Could use signaling during Sprint 005 implementation**:
```bash
# Supervisor (you)
my-context watch "my-context: 005-ux-polish__implementation" \
  --pattern="Phase.*complete" \
  --exec="echo 'Phase complete' | notify-send"

# Implementation agent
my-context note "Phase 1 complete: Timestamps implemented"
# Supervisor's watch triggers, you get notification
```

**Benefit**: Real-world validation of signaling pattern before full Sprint 006 implementation

---

**Document Version**: 1.0.0
**Status**: Requirements captured from real handoff
**Next**: Incorporate into Sprint 006 spec when ready
