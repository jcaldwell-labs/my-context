#!/usr/bin/env bash
#
# FR-MC-NEW-004: Note Limit Warning POC
#
# Wraps `my-context note` to add warnings at note count thresholds
#
# Usage: ./note-with-warning.sh "Your note message"
#

set -e

# Configurable thresholds
WARN_THRESHOLD_1=${MC_WARN_AT:-50}
WARN_THRESHOLD_2=${MC_WARN_AT_2:-100}
WARN_THRESHOLD_3=${MC_WARN_AT_3:-200}

# Get current context info
CONTEXT_INFO=$(my-context show 2>/dev/null) || {
    echo "❌ No active context. Start one with: my-context start <name>"
    exit 1
}

# Extract current note count
NOTE_COUNT=$(echo "$CONTEXT_INFO" | grep -oP '(?<=Notes \()\d+' || echo "0")

# Add the note (pass through all arguments)
my-context note "$@"

# Increment count for warning check
NEW_COUNT=$((NOTE_COUNT + 1))

# Check thresholds and show warnings
if [[ $NEW_COUNT -eq $WARN_THRESHOLD_1 ]]; then
    echo ""
    echo "⚠️  Context has $NEW_COUNT notes. Consider:"
    echo "   - Stop and start new phase context if switching focus"
    echo "   - Continue if still same work (100+ notes is fine)"
    echo "   - Export to markdown: my-context export <context-name>"
elif [[ $NEW_COUNT -eq $WARN_THRESHOLD_2 ]]; then
    echo ""
    echo "⚠️  Context has $NEW_COUNT notes (getting large). Consider:"
    echo "   - Export to markdown for review: my-context export <context-name>"
    echo "   - Split into phase contexts if work has shifted"
    echo "   - Or continue - large contexts are OK for complex work"
elif [[ $NEW_COUNT -ge $WARN_THRESHOLD_3 ]] && [[ $(($NEW_COUNT % 25)) -eq 0 ]]; then
    echo ""
    echo "📊 Context has $NEW_COUNT notes (very detailed tracking!)"
    echo "   Consider exporting: my-context export <context-name>"
fi

exit 0
