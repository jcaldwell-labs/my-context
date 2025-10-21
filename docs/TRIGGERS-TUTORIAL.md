# My-Context Triggers Tutorial: Personal Productivity Workflows

**Learn how to automate your work context management with signals and watches**

**Reading time**: ~60 minutes
**Skill level**: Beginner to Advanced
**Prerequisites**: Basic familiarity with my-context commands (start, note, stop)

---

## Table of Contents

1. [Introduction](#introduction)
2. [Part 1: The Basics](#part-1-the-basics)
3. [Part 2: Personal Productivity Patterns](#part-2-personal-productivity-patterns)
4. [Part 3: Advanced Workflows](#part-3-advanced-workflows)
5. [Part 4: Troubleshooting & Tips](#part-4-troubleshooting--tips)

---

## Introduction

### What You'll Learn

This tutorial will teach you how to use my-context's **signals** and **watches** to automate your personal productivity workflows. By the end, you'll be able to:

- Get automatic notifications when you hit milestones
- Receive reminders when contexts get too large
- Automate end-of-day cleanup routines
- Integrate my-context with focus timers (Pomodoro, etc.)
- Create smart context lifecycle automation

### What are Signals and Watches?

**Signals** are simple event notifications - think of them as flags you raise to say "something happened!"

**Watches** are monitors that observe your contexts and react when changes occur.

Together, they enable **event-driven workflows** where my-context can automatically respond to your work patterns instead of requiring manual commands.

### When to Use Triggers vs. Manual Commands

| Use Manual Commands When | Use Triggers When |
|--------------------------|-------------------|
| You're actively working | You want automation |
| One-time actions | Repeated patterns |
| Simple workflows | Multi-step workflows |
| Learning my-context | Optimizing workflows |

---

## Part 1: The Basics

### Your First Signal

Signals are stored as timestamped files in `~/.my-context/signals/`. They're incredibly simple but powerful.

**Create a signal:**
```bash
my-context signal create my-first-signal
```

Output:
```
Signal 'my-first-signal' created
```

**List all signals:**
```bash
my-context signal list
```

Output:
```
Signals:
  my-first-signal (created: 2025-10-21T07:25:06-04:00)
```

**Clear a signal:**
```bash
my-context signal clear my-first-signal
```

Output:
```
Signal 'my-first-signal' cleared
```

### Your First Watch

Watches monitor contexts for changes and can execute commands when changes are detected.

**Watch your active context:**
```bash
# Start a context first
my-context start "my-work"
my-context note "Starting implementation"

# In another terminal, watch for changes
my-context watch
```

The watch command will now monitor your active context. When you add notes or files, it detects the changes.

**Watch for new notes specifically:**
```bash
my-context watch --new-notes
```

**Watch with a pattern (filter for specific content):**
```bash
my-context watch --new-notes --pattern="complete|done|finished"
```

This will only trigger when notes contain words like "complete", "done", or "finished".

**Watch and execute a command:**
```bash
my-context watch --new-notes --exec="echo 'Context updated!'"
```

**Stop watching** with `Ctrl+C`.

### Understanding Watch Options

- `--new-notes`: Only watch for new notes (not other file changes)
- `--pattern <regex>`: Filter notes by content using regular expressions
- `--exec <command>`: Execute this command when changes detected
- `--interval <duration>`: How often to check (default: 5s)
- `--timeout <duration>`: Maximum time to watch (default: infinite)

---

## Part 2: Personal Productivity Patterns

### Pattern 1: Context Completion Notifications

**Use Case**: You want to be notified when you complete major milestones, even if you're not actively at your computer.

**Setup:**

1. Create a completion detection script (`~/bin/notify-completion.sh`):

```bash
#!/bin/bash
# ~/bin/notify-completion.sh

CONTEXT_NAME=$(my-context show --json 2>/dev/null | jq -r '.data.context.name // "unknown"')

# Send notification (Linux)
notify-send "🎯 Milestone!" "Completed work in: $CONTEXT_NAME"

# macOS alternative:
# osascript -e "display notification \"Completed work in: $CONTEXT_NAME\" with title \"🎯 Milestone!\""

# Windows (WSL) alternative:
# powershell.exe -Command "& {Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('Completed work in: $CONTEXT_NAME', '🎯 Milestone!')}"
```

2. Make it executable:
```bash
chmod +x ~/bin/notify-completion.sh
```

3. Start watching your work context:
```bash
my-context start "feature-implementation"

# In another terminal or background
my-context watch --new-notes --pattern="complete|done|finished|implemented" --exec="~/bin/notify-completion.sh" &
```

4. When you add a completion note:
```bash
my-context note "Feature implementation complete!"
```

You'll get an automatic notification!

**Advanced variation - Signal-based:**

```bash
# When you complete work, create a signal
my-context note "Feature complete!"
my-context signal create feature-complete

# Separate process watches for the signal
my-context signal wait feature-complete --timeout=8h && ~/bin/notify-completion.sh
```

---

### Pattern 2: Context Size Warnings

**Use Case**: You want reminders when your context is getting too large and might need splitting or archiving.

**Setup:**

1. Create a note-counter script (`~/bin/check-context-size.sh`):

```bash
#!/bin/bash
# ~/bin/check-context-size.sh

CONTEXT=$(my-context show --json 2>/dev/null)
if [ $? -ne 0 ]; then
    exit 0  # No active context
fi

CONTEXT_NAME=$(echo "$CONTEXT" | jq -r '.data.context.name')
NOTE_COUNT=$(echo "$CONTEXT" | jq -r '.data.context.notes | length')

# Thresholds
WARN_AT=${MC_WARN_AT:-50}
WARN_AT_2=${MC_WARN_AT_2:-100}
WARN_AT_3=${MC_WARN_AT_3:-200}

if [ "$NOTE_COUNT" -ge "$WARN_AT_3" ]; then
    echo "⚠️  Context '$CONTEXT_NAME' has $NOTE_COUNT notes (very large!)"
    echo "   Consider: my-context archive \"$CONTEXT_NAME\" && my-context start \"${CONTEXT_NAME}-phase-2\""
elif [ "$NOTE_COUNT" -ge "$WARN_AT_2" ]; then
    echo "⚠️  Context '$CONTEXT_NAME' has $NOTE_COUNT notes (getting large)"
    echo "   Consider: my-context export \"$CONTEXT_NAME\" for review"
elif [ "$NOTE_COUNT" -ge "$WARN_AT" ]; then
    echo "⚠️  Context '$CONTEXT_NAME' has $NOTE_COUNT notes"
    echo "   Consider: starting a new phase if work has shifted focus"
fi
```

2. Make it executable:
```bash
chmod +x ~/bin/check-context-size.sh
```

3. Run it automatically after every note:

Add to your shell profile (`~/.bashrc` or `~/.zshrc`):
```bash
# Wrap my-context note command
mcn() {
    my-context note "$@"
    ~/bin/check-context-size.sh
}

alias cn='mcn'  # Quick alias
```

4. Now when you add notes:
```bash
cn "Adding more implementation details"
cn "Debugging edge case"
# ... after 50 notes ...
cn "Another fix"
```

Output:
```
Note added to context: feature-implementation
⚠️  Context 'feature-implementation' has 51 notes
   Consider: starting a new phase if work has shifted focus
```

---

### Pattern 3: End-of-Day Cleanup Automation

**Use Case**: At the end of each workday, you want to automatically review, export, and archive completed contexts.

**Setup:**

1. Create cleanup script (`~/bin/end-of-day-cleanup.sh`):

```bash
#!/bin/bash
# ~/bin/end-of-day-cleanup.sh

echo "🌅 End-of-Day Context Cleanup"
echo "================================"

# Export all stopped contexts from today
TODAY=$(date +%Y-%m-%d)
EXPORT_DIR=~/context-exports/$TODAY
mkdir -p "$EXPORT_DIR"

echo "📦 Exporting today's contexts to: $EXPORT_DIR"

# Get all contexts (this assumes my-context list outputs context names)
my-context list --json | jq -r '.data.contexts[] | select(.status=="stopped") | .name' | while read -r context; do
    echo "  ✓ Exporting: $context"
    my-context export "$context" --to "$EXPORT_DIR/${context}.md" 2>/dev/null
done

# Archive completed contexts (ones with "complete" in recent notes)
echo ""
echo "📚 Archiving completed contexts..."

my-context list --json | jq -r '.data.contexts[] | select(.status=="stopped") | .name' | while read -r context; do
    # Check if last note contains completion keywords
    LAST_NOTE=$(my-context show "$context" --json 2>/dev/null | jq -r '.data.context.notes[-1].text // ""')

    if echo "$LAST_NOTE" | grep -iE "complete|done|finished|resolved|archived" > /dev/null; then
        echo "  ✓ Archiving: $context (marked complete)"
        my-context archive "$context"
    fi
done

# Create signal for day complete
my-context signal create workday-$(date +%Y-%m-%d)-complete

echo ""
echo "✨ Cleanup complete! Contexts exported to: $EXPORT_DIR"
echo "🏁 Created signal: workday-$(date +%Y-%m-%d)-complete"
```

2. Make it executable:
```bash
chmod +x ~/bin/end-of-day-cleanup.sh
```

3. Run manually or on a schedule:

**Manual:**
```bash
~/bin/end-of-day-cleanup.sh
```

**Automated with cron (run at 6 PM every weekday):**
```bash
# Edit crontab
crontab -e

# Add this line:
0 18 * * 1-5 /home/yourusername/bin/end-of-day-cleanup.sh >> /home/yourusername/cleanup.log 2>&1
```

**On-demand with alias:**
```bash
# Add to ~/.bashrc or ~/.zshrc
alias eod='~/bin/end-of-day-cleanup.sh'
```

Now just type `eod` when you finish work!

---

### Pattern 4: Pomodoro/Focus Session Integration

**Use Case**: You want to track focus sessions (25-minute Pomodoros) and automatically log them in your context.

**Setup:**

1. Create a Pomodoro timer script (`~/bin/pomodoro.sh`):

```bash
#!/bin/bash
# ~/bin/pomodoro.sh - Simple Pomodoro timer with my-context integration

WORK_DURATION=${1:-25}  # Default 25 minutes
BREAK_DURATION=${2:-5}  # Default 5 minutes

# Check if context is active
if ! my-context show --json 2>/dev/null | jq -e '.data.context' > /dev/null; then
    echo "❌ No active context. Start one first: my-context start <name>"
    exit 1
fi

CONTEXT=$(my-context show --json | jq -r '.data.context.name')

echo "🍅 Pomodoro Timer: $WORK_DURATION minutes"
echo "📋 Context: $CONTEXT"
echo ""

# Start work session
my-context note "🍅 Pomodoro started ($WORK_DURATION min focus session)"
my-context signal create pomodoro-started
START_TIME=$(date +%s)

# Countdown
for ((i=$WORK_DURATION; i>0; i--)); do
    echo -ne "⏱️  Focus time remaining: $i minutes\r"
    sleep 60
    my-context touch  # Heartbeat every minute
done

END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))

# Work session complete
echo ""
echo "✅ Pomodoro complete!"
my-context note "🍅 Pomodoro complete (${WORK_DURATION}m focused work)"
my-context signal create pomodoro-complete

# Break time
echo ""
echo "☕ Break time: $BREAK_DURATION minutes"
sleep $((BREAK_DURATION * 60))

echo "🔔 Break over! Ready for next Pomodoro?"
my-context signal clear pomodoro-started
my-context signal clear pomodoro-complete
```

2. Make it executable:
```bash
chmod +x ~/bin/pomodoro.sh
```

3. Use it:
```bash
# Start a context
my-context start "deep-work-session"

# Run Pomodoro (25 min work, 5 min break)
~/bin/pomodoro.sh

# Custom durations: 50 min work, 10 min break
~/bin/pomodoro.sh 50 10
```

4. Optional - Watch for Pomodoro completion in another terminal:
```bash
my-context signal wait pomodoro-complete && echo "Great job! Take your break." | wall
```

**Advanced - Track daily Pomodoros:**

Add to your shell profile:
```bash
# Count today's Pomodoros
alias pomodoro-count='my-context show --json | jq "[.data.context.notes[] | select(.text | contains(\"🍅 Pomodoro complete\"))] | length"'
```

Check your progress:
```bash
$ pomodoro-count
4  # You've done 4 Pomodoros today!
```

---

### Pattern 5: Context Lifecycle Automation

**Use Case**: Automate the decision of whether to resume an existing context or create a new one when starting work.

**Setup:**

This uses the POC script as inspiration but integrates with signals.

1. Create a smart start script (`~/bin/smart-start.sh`):

```bash
#!/bin/bash
# ~/bin/smart-start.sh - Smart context start with resume detection

CONTEXT_NAME="$1"

if [ -z "$CONTEXT_NAME" ]; then
    echo "Usage: smart-start.sh <context-name>"
    exit 1
fi

# Check if context exists and is stopped
EXISTING=$(my-context list --json | jq -r ".data.contexts[] | select(.name==\"$CONTEXT_NAME\" and .status==\"stopped\") | .name")

if [ -n "$EXISTING" ]; then
    # Context exists and is stopped
    NOTE_COUNT=$(my-context show "$CONTEXT_NAME" --json | jq -r '.data.context.notes | length')
    LAST_ACTIVE=$(my-context show "$CONTEXT_NAME" --json | jq -r '.data.context.ended_at // .data.context.started_at')

    echo "📋 Context '$CONTEXT_NAME' already exists (stopped)"
    echo ""
    echo "   Notes: $NOTE_COUNT"
    echo "   Last active: $LAST_ACTIVE"
    echo ""
    read -p "Resume existing context? [Y/n]: " -n 1 -r
    echo

    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        my-context start "$CONTEXT_NAME"
        echo "✅ Resumed context: $CONTEXT_NAME"
        my-context signal create context-resumed
    else
        read -p "New context name (or Enter to use '${CONTEXT_NAME}_2'): " NEW_NAME
        NEW_NAME=${NEW_NAME:-"${CONTEXT_NAME}_2"}
        my-context start "$NEW_NAME"
        echo "✅ Started new context: $NEW_NAME"
        my-context signal create context-created
    fi
else
    # Context doesn't exist, create it
    my-context start "$CONTEXT_NAME"
    echo "✅ Started new context: $CONTEXT_NAME"
    my-context signal create context-created
fi
```

2. Make it executable:
```bash
chmod +x ~/bin/smart-start.sh
```

3. Add alias to your shell:
```bash
# Add to ~/.bashrc or ~/.zshrc
alias mcs='~/bin/smart-start.sh'
```

4. Use it:
```bash
# First time - creates new context
$ mcs "feature-auth"
✅ Started new context: feature-auth

# Stop and later try to start same name
$ my-context stop
$ mcs "feature-auth"

📋 Context 'feature-auth' already exists (stopped)

   Notes: 5
   Last active: 2025-10-21T08:00:00-04:00

Resume existing context? [Y/n]: y
✅ Resumed context: feature-auth
```

---

## Part 3: Advanced Workflows

### Chaining Signals for Multi-Step Workflows

You can chain signals to create complex automated workflows:

**Example: Code → Build → Test → Deploy Pipeline**

```bash
# Terminal 1: Development
my-context start "pipeline-demo"
my-context note "Implementing feature X"
# ... coding ...
my-context note "Implementation complete"
my-context signal create code-complete

# Terminal 2: Build watcher
my-context signal wait code-complete && {
    echo "🔨 Building..."
    ./build.sh
    if [ $? -eq 0 ]; then
        my-context note "Build successful"
        my-context signal create build-complete
    else
        my-context note "Build failed!"
        exit 1
    fi
}

# Terminal 3: Test watcher
my-context signal wait build-complete && {
    echo "🧪 Testing..."
    ./run-tests.sh
    if [ $? -eq 0 ]; then
        my-context note "All tests passed"
        my-context signal create tests-complete
    else
        my-context note "Tests failed!"
        exit 1
    fi
}

# Terminal 4: Deploy watcher
my-context signal wait tests-complete && {
    echo "🚀 Deploying..."
    ./deploy.sh
    my-context note "Deployed to production"
    my-context signal create deploy-complete
}
```

This creates a fully automated pipeline triggered by a single signal!

### Multi-Context Monitoring

Monitor multiple contexts simultaneously:

```bash
#!/bin/bash
# ~/bin/multi-watch.sh - Watch multiple contexts

CONTEXTS=("frontend-work" "backend-work" "docs-work")

for ctx in "${CONTEXTS[@]}"; do
    {
        echo "👀 Watching: $ctx"
        my-context watch "$ctx" --new-notes --exec="echo '[$ctx] New activity!'" &
    }
done

wait  # Keep script running
```

### Integration with System Notifications

**Desktop notifications (Linux - notify-send):**
```bash
my-context watch --new-notes --pattern="urgent|critical|emergency" --exec="notify-send -u critical 'URGENT' 'Check my-context!'"
```

**macOS notifications:**
```bash
my-context watch --new-notes --exec="osascript -e 'display notification \"Context updated\" with title \"My-Context\"'"
```

**Email notifications:**
```bash
# When critical pattern detected, send email
my-context watch --pattern="critical|emergency" --exec="echo 'Critical event in context' | mail -s 'My-Context Alert' you@example.com"
```

**Slack/Discord webhooks:**
```bash
#!/bin/bash
# webhook-notify.sh
WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
CONTEXT=$(my-context show --json | jq -r '.data.context.name')
LAST_NOTE=$(my-context show --json | jq -r '.data.context.notes[-1].text')

curl -X POST "$WEBHOOK_URL" \
  -H 'Content-Type: application/json' \
  -d "{\"text\": \"📋 Context: $CONTEXT\n$LAST_NOTE\"}"
```

Then watch with:
```bash
my-context watch --new-notes --exec="~/bin/webhook-notify.sh"
```

### Shell Integration Examples

**Add context to your prompt:**

```bash
# Add to ~/.bashrc or ~/.zshrc
my_context_prompt() {
    ACTIVE=$(my-context show --json 2>/dev/null | jq -r '.data.context.name // ""')
    if [ -n "$ACTIVE" ]; then
        echo " [ctx:$ACTIVE]"
    fi
}

# Bash
PS1='$(my_context_prompt)\u@\h:\w\$ '

# Zsh
PROMPT='$(my_context_prompt)%n@%m:%~%# '
```

Now your prompt shows:
```
[ctx:feature-work]user@host:~/project$
```

**Auto-note on directory change:**

```bash
# Add to ~/.bashrc or ~/.zshrc
cd() {
    builtin cd "$@"
    if my-context show --json 2>/dev/null | jq -e '.data.context' > /dev/null; then
        my-context file "$(pwd)" 2>/dev/null
    fi
}
```

Now every time you `cd`, the directory is auto-associated with your active context!

**Auto-touch on git commit:**

```bash
# .git/hooks/post-commit (make executable: chmod +x)
#!/bin/bash

if my-context show --json 2>/dev/null | jq -e '.data.context' > /dev/null; then
    COMMIT_MSG=$(git log -1 --pretty=%B)
    my-context note "Committed: $COMMIT_MSG"
fi
```

Now every git commit automatically adds a note to your active context!

---

## Part 4: Troubleshooting & Tips

### Common Pitfalls

#### Watch command consuming too much CPU

**Problem**: The watch command polls every 5 seconds by default, which can add up if you're watching many contexts.

**Solution 1 - Increase interval:**
```bash
my-context watch --interval=30s  # Check every 30 seconds
```

**Solution 2 - Set a timeout:**
```bash
my-context watch --timeout=2h    # Only watch for 2 hours
```

**Solution 3 - Use specific patterns:**
```bash
# Only trigger on specific patterns (fewer false positives)
my-context watch --new-notes --pattern="^(complete|done|finished)"
```

#### Signals pile up over time

**Problem**: You create many signals but forget to clear them.

**Solution - Periodic cleanup:**
```bash
# Add to your end-of-day cleanup
my-context signal list | grep -E "$(date +%Y-%m-%d)" | while read -r signal; do
    my-context signal clear "$(echo $signal | cut -d' ' -f1)"
done
```

Or create an alias:
```bash
alias signal-cleanup='my-context signal list | cut -d" " -f1 | xargs -I {} my-context signal clear {}'
```

#### Signal wait blocks forever

**Problem**: `my-context signal wait` blocks indefinitely if signal never arrives.

**Solution - Always use timeout:**
```bash
# Timeout after 1 hour
my-context signal wait my-signal --timeout=1h

# Check exit code to detect timeout
if [ $? -ne 0 ]; then
    echo "Signal wait timed out or failed"
fi
```

#### Watch doesn't detect changes immediately

**Problem**: There's a delay between making changes and watch triggering.

**Explanation**: Watch uses polling (checks every N seconds) by default, not real-time file monitoring.

**Solution**:
```bash
# Reduce interval for faster detection (uses more CPU)
my-context watch --interval=1s
```

On Linux, my-context may use inotify for real-time detection automatically, but this isn't guaranteed.

### Performance Considerations

#### Polling Intervals

- **Default (5s)**: Good for most use cases
- **Fast (1s)**: Real-time feel, higher CPU usage
- **Slow (30s-1m)**: Background monitoring, minimal impact
- **Very slow (5m+)**: Periodic checks, almost no CPU

**Recommendation**: Start with defaults, only optimize if you notice performance issues.

#### Number of Concurrent Watches

Each watch command runs in a separate process. Watching too many contexts simultaneously can impact performance.

**Guidelines**:
- 1-5 watches: No problem
- 5-10 watches: Increase intervals to 10-30s
- 10+ watches: Consider alternative approaches (cron jobs, signal-based)

#### Signal File Cleanup

Signals are tiny files (26 bytes), so they don't consume much space. However, for cleanliness:

- Clear signals after they're consumed
- Periodic cleanup (daily/weekly)
- Use descriptive names with dates: `workday-2025-10-21-complete`

### Debugging Signal/Watch Issues

#### Enable verbose output

Signals and watches don't have a verbose flag yet, but you can debug with:

```bash
# Check if signals directory exists
ls -la ~/.my-context/signals/

# Monitor signal creation in real-time
watch -n 1 'ls -lt ~/.my-context/signals/ | head -10'

# Monitor context changes
watch -n 1 'cat ~/.my-context/your-context/notes.log | tail -5'
```

#### Test your exec commands

Before using `--exec`, test the command independently:

```bash
# Test your notification script
~/bin/notify-completion.sh

# Test with sample data
echo "test" | ~/bin/process-note.sh
```

#### Check for permission issues

```bash
# Verify signals directory is writable
touch ~/.my-context/signals/test-permission.signal
rm ~/.my-context/signals/test-permission.signal
```

### Best Practices

#### 1. Use Descriptive Signal Names

**Good**:
```bash
my-context signal create feature-auth-complete
my-context signal create tests-passing-2025-10-21
my-context signal create ready-for-review
```

**Bad**:
```bash
my-context signal create done
my-context signal create s1
my-context signal create x
```

#### 2. Clean Up Signals After Use

```bash
# Pattern: Create → Wait → Clear
my-context signal create task-complete
# ... wait for signal ...
my-context signal wait task-complete --timeout=1h
my-context signal clear task-complete
```

#### 3. Use Timeouts on Waits

Always include `--timeout` to prevent infinite blocking:

```bash
my-context signal wait my-signal --timeout=30m
```

#### 4. Combine Patterns for Precision

```bash
# Too broad - triggers on any note with "test"
my-context watch --pattern="test"

# Better - specific pattern
my-context watch --pattern="^All tests (passed|complete)"
```

#### 5. Document Your Workflows

Create a `~/my-context-workflows.md` documenting your custom scripts and triggers:

```markdown
# My Context Workflows

## End of Day
- Run: `eod` (alias for end-of-day-cleanup.sh)
- Exports all contexts to ~/context-exports/<date>/
- Archives completed work
- Creates signal: workday-YYYY-MM-DD-complete

## Pomodoro Sessions
- Run: `~/bin/pomodoro.sh [work-mins] [break-mins]`
- Default: 25 min work, 5 min break
- Creates signals: pomodoro-started, pomodoro-complete
- Logs touch events every minute
```

### Pro Tips

**Tip 1: Use signals for async communication**

If you're working with a team, signals can replace manual notifications:

```bash
# You
my-context note "Code review complete"
my-context signal create code-review-done-issue-123

# Teammate (in background)
my-context signal wait code-review-done-issue-123 && \
  notify-send "Code review ready!"
```

**Tip 2: Combine with cron for scheduled tasks**

```bash
# Check for stale contexts every morning at 9 AM
0 9 * * * /home/user/bin/check-stale-contexts.sh
```

**Tip 3: Use watch patterns as "smart filters"**

```bash
# Only care about completed phases
my-context watch --pattern="Phase [0-9]+ complete"

# Only care about errors
my-context watch --pattern="ERROR|FATAL|CRITICAL"

# Only care about specific features
my-context watch --pattern="^(auth|payment|checkout):"
```

**Tip 4: Create signal chains for dependencies**

```bash
# Phase 1 triggers Phase 2
my-context signal wait phase-1-complete && {
    my-context note "Starting Phase 2"
    ./phase-2-script.sh
    my-context signal create phase-2-complete
}

# Phase 2 triggers Phase 3
my-context signal wait phase-2-complete && {
    my-context note "Starting Phase 3"
    ./phase-3-script.sh
}
```

---

## Conclusion

You now have the tools to automate your personal productivity workflows with my-context signals and watches!

### Quick Reference Card

```bash
# SIGNALS
my-context signal create <name>          # Create event signal
my-context signal list                   # List all signals
my-context signal wait <name> --timeout=1h  # Wait for signal
my-context signal clear <name>           # Remove signal

# WATCHES
my-context watch                         # Watch active context
my-context watch <context>               # Watch specific context
my-context watch --new-notes             # Only watch for notes
my-context watch --pattern="regex"       # Filter by pattern
my-context watch --exec="command"        # Execute on change
my-context watch --interval=10s          # Set poll interval
my-context watch --timeout=2h            # Set max watch time
```

### Next Steps

1. **Try the basic examples** - Create a signal, set up a watch
2. **Implement one productivity pattern** - Start with Pattern 1 (completion notifications)
3. **Customize to your workflow** - Adapt the scripts to your needs
4. **Share your patterns** - Contribute back to the community!

### Resources

- [My-Context README](../README.md) - Main documentation
- [Signaling Protocol Requirements](SIGNALING-PROTOCOL-REQUIREMENTS.md) - Technical details
- [Troubleshooting Guide](TROUBLESHOOTING.md) - Common issues
- [POC Scripts](../scripts/poc/) - Example implementations

---

**Document Version**: 1.0.0
**Created**: 2025-10-21
**Author**: My-Context Team
**License**: MIT

**Feedback**: Found this tutorial helpful? Have suggestions? Open an issue on GitHub!
