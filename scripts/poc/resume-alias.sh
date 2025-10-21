#!/usr/bin/env bash
#
# FR-MC-NEW-005: Resume Command Alias POC
#
# Provides explicit `resume` command for resuming stopped contexts
#
# Usage:
#   ./resume-alias.sh <context-name>         # Resume specific context
#   ./resume-alias.sh --last                 # Resume most recently stopped
#   ./resume-alias.sh <pattern>              # Resume matching pattern
#

set -e

show_usage() {
    cat <<EOF
Usage: my-context-resume [OPTIONS] [CONTEXT-NAME]

Resume a stopped context.

OPTIONS:
  --last              Resume most recently stopped context
  -h, --help          Show this help

EXAMPLES:
  my-context-resume "foo"           # Resume context named "foo"
  my-context-resume --last          # Resume last stopped context
  my-context-resume "007-*"         # Resume context matching pattern

EOF
}

# Parse arguments
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_usage
    exit 0
fi

if [[ "$1" == "--last" ]]; then
    # Resume most recently stopped context
    LAST_STOPPED=$(my-context list | grep -A 1000 "stopped" | grep "○" | head -1 | sed 's/.*○ \(.*\) (stopped).*/\1/')

    if [[ -z "$LAST_STOPPED" ]]; then
        echo "❌ No stopped contexts found"
        exit 1
    fi

    echo "🔄 Resuming most recent: $LAST_STOPPED"
    my-context start "$LAST_STOPPED"
    exit 0
fi

# Resume specific context or pattern
CONTEXT_NAME="$1"

# Check if context exists and is stopped
ALL_CONTEXTS=$(my-context list)

# Try exact match first
if echo "$ALL_CONTEXTS" | grep -q "○ $CONTEXT_NAME (stopped)"; then
    # Context exists and is stopped - resume it
    my-context start "$CONTEXT_NAME"
    echo "✅ Resumed context: $CONTEXT_NAME"
    exit 0
fi

# Try pattern match
MATCHES=$(echo "$ALL_CONTEXTS" | grep "○.*$CONTEXT_NAME.*(stopped)" | sed 's/.*○ \(.*\) (stopped).*/\1/' || true)
MATCH_COUNT=$(echo "$MATCHES" | grep -v "^$" | wc -l)

if [[ $MATCH_COUNT -eq 0 ]]; then
    echo "❌ No stopped context matching '$CONTEXT_NAME'"
    echo ""
    echo "💡 Available stopped contexts:"
    my-context list | grep "○" | head -5
    exit 1
elif [[ $MATCH_COUNT -eq 1 ]]; then
    # Single match - resume it
    CONTEXT_TO_RESUME=$(echo "$MATCHES" | head -1)
    echo "🔄 Resuming: $CONTEXT_TO_RESUME"
    my-context start "$CONTEXT_TO_RESUME"
    exit 0
else
    # Multiple matches - let user choose
    echo "Multiple matching contexts:"
    echo ""

    i=1
    while IFS= read -r context; do
        [[ -z "$context" ]] && continue
        echo "  ($i) $context"
        i=$((i + 1))
    done <<< "$MATCHES"

    echo ""
    read -p "Select context [1-$((i-1))]: " selection

    SELECTED_CONTEXT=$(echo "$MATCHES" | sed -n "${selection}p")

    if [[ -z "$SELECTED_CONTEXT" ]]; then
        echo "❌ Invalid selection"
        exit 1
    fi

    echo "🔄 Resuming: $SELECTED_CONTEXT"
    my-context start "$SELECTED_CONTEXT"
fi

exit 0
