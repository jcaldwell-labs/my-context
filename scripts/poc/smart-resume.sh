#!/usr/bin/env bash
#
# FR-MC-NEW-001: Smart Resume on Start POC
#
# Wraps `my-context start` to detect existing stopped contexts and prompt for resume
#
# Usage: ./smart-resume.sh <context-name> [--force]
#

set -e

show_usage() {
    cat <<EOF
Usage: my-context-start <context-name> [--force]

Start or resume a context with smart duplicate detection.

OPTIONS:
  --force             Force creation of new context even if name exists
  -h, --help          Show this help

EXAMPLES:
  my-context-start "foo"           # Prompts if "foo" exists (stopped)
  my-context-start "foo" --force   # Creates new "foo" regardless

EOF
}

# Parse arguments
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_usage
    exit 0
fi

CONTEXT_NAME="$1"
FORCE=false

if [[ "$2" == "--force" ]]; then
    FORCE=true
fi

# Check if context already exists and is stopped
ALL_CONTEXTS=$(my-context list 2>/dev/null || true)

# Check for exact stopped match
if echo "$ALL_CONTEXTS" | grep -q "○ $CONTEXT_NAME (stopped)"; then
    if [[ "$FORCE" == "true" ]]; then
        # Force new context - suggest different name
        echo "⚠️  Context '$CONTEXT_NAME' exists (stopped) but --force specified"
        echo ""
        read -p "Create new context with different name (e.g., ${CONTEXT_NAME}_2): " NEW_NAME

        if [[ -z "$NEW_NAME" ]]; then
            echo "❌ No name provided, aborting"
            exit 1
        fi

        my-context start "$NEW_NAME"
        exit 0
    fi

    # Get context info
    CONTEXT_DETAILS=$(my-context list | grep -A 5 "○ $CONTEXT_NAME")

    # Extract note count if available
    NOTE_INFO=$(echo "$CONTEXT_DETAILS" | grep -oP '\d+(?= notes)' | head -1 || echo "unknown")

    # Extract last active time if available (simplified - would need more parsing)
    echo "📋 Context '$CONTEXT_NAME' exists (stopped)"
    echo ""
    if [[ "$NOTE_INFO" != "unknown" ]]; then
        echo "   Notes: $NOTE_INFO"
    fi
    echo ""
    read -p "Resume existing context? [Y/n]: " RESUME

    RESUME=${RESUME:-Y}

    if [[ "$RESUME" =~ ^[Yy]$ ]]; then
        my-context start "$CONTEXT_NAME"
        echo "✅ Resumed context: $CONTEXT_NAME"
        exit 0
    else
        echo ""
        read -p "Start new context with different name (e.g., ${CONTEXT_NAME}_2): " NEW_NAME

        if [[ -z "$NEW_NAME" ]]; then
            echo "❌ No name provided, aborting"
            exit 1
        fi

        my-context start "$NEW_NAME"
        exit 0
    fi
fi

# Check if context is active
if echo "$ALL_CONTEXTS" | grep -q "● $CONTEXT_NAME (active)"; then
    echo "❌ Context '$CONTEXT_NAME' is already active"
    exit 1
fi

# Context doesn't exist - create it
my-context start "$CONTEXT_NAME"

exit 0
