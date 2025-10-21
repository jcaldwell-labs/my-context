#!/usr/bin/env bash
#
# FR-MC-NEW-002: Context Lifecycle Advisor POC
#
# Wraps `my-context stop` to provide lifecycle guidance
#
# Usage: ./lifecycle-advisor.sh
#

set -e

# Get current context info
CONTEXT_INFO=$(my-context show 2>/dev/null) || {
    echo "❌ No active context"
    exit 1
}

# Extract context name
CONTEXT_NAME=$(echo "$CONTEXT_INFO" | grep "^Context:" | sed 's/Context: //')

# Stop the context
my-context stop

# Get context details after stopping
STOPPED_INFO=$(my-context list | grep -A 3 "○ $CONTEXT_NAME" || true)

# Extract duration and notes (simplified parsing)
DURATION=$(echo "$STOPPED_INFO" | grep -oP '(?<=Duration: ).*' || echo "unknown")
NOTE_COUNT=$(echo "$STOPPED_INFO" | grep -oP '\d+(?= notes)' || echo "0")

echo ""
echo "📊 Context Summary:"
echo "   Name: $CONTEXT_NAME"
[[ -n "$DURATION" ]] && echo "   Duration: $DURATION"
[[ -n "$NOTE_COUNT" ]] && echo "   Notes: $NOTE_COUNT"
echo ""

# Find related contexts (similar names)
CONTEXT_PREFIX=$(echo "$CONTEXT_NAME" | sed 's/_[0-9]*$//' | sed 's/-phase.*//' | sed 's/__.*$//')
RELATED=$(my-context list | grep "○.*$CONTEXT_PREFIX" | grep -v "$CONTEXT_NAME" || true)

if [[ -n "$RELATED" ]]; then
    echo "📚 Related contexts:"
    echo "$RELATED" | head -3
    echo ""
fi

# Provide suggestions
echo "💡 Suggestions:"

# Check for related stopped contexts
if [[ -n "$RELATED" ]]; then
    FIRST_RELATED=$(echo "$RELATED" | head -1 | sed 's/.*○ \(.*\) (stopped.*/\1/')
    echo "   (1) Resume related: my-context start \"$FIRST_RELATED\""
fi

# Suggest archiving
echo "   (2) Archive if complete: my-context archive \"$CONTEXT_NAME\""

# Suggest starting new
echo "   (3) Start new work: my-context start <name>"

echo ""

# Check for completion indicators in recent notes
RECENT_NOTES=$(echo "$CONTEXT_INFO" | grep -A 20 "Notes (" | tail -5)

COMPLETION_KEYWORDS="complete|finished|done|retrospective|merged|deployed|closing|wrapping"

if echo "$RECENT_NOTES" | grep -qi "$COMPLETION_KEYWORDS"; then
    echo "🎯 Completion detected in recent notes!"
    echo "   Consider archiving: my-context archive \"$CONTEXT_NAME\""
    echo ""
fi

exit 0
