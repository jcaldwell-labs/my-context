#!/usr/bin/env bash
#
# FR-MC-NEW-003: Bulk Archive with Patterns POC
#
# Archive multiple contexts matching pattern or criteria
#
# Usage:
#   ./bulk-archive.sh --pattern "006-*"
#   ./bulk-archive.sh --completed-before "2025-10-09"
#   ./bulk-archive.sh --all-stopped
#   ./bulk-archive.sh --pattern "foo*" --dry-run
#

set -e

show_usage() {
    cat <<EOF
Usage: my-context-bulk-archive [OPTIONS]

Archive multiple contexts at once.

OPTIONS:
  --pattern PATTERN         Archive contexts matching glob pattern
  --completed-before DATE   Archive contexts stopped before date (YYYY-MM-DD)
  --all-stopped             Archive all stopped contexts
  --dry-run                 Show what would be archived without doing it
  -h, --help                Show this help

EXAMPLES:
  my-context-bulk-archive --pattern "006-*"
  my-context-bulk-archive --completed-before "2025-10-09"
  my-context-bulk-archive --all-stopped
  my-context-bulk-archive --pattern "test*" --dry-run

EOF
}

# Parse arguments
PATTERN=""
DATE_BEFORE=""
ALL_STOPPED=false
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --pattern)
            PATTERN="$2"
            shift 2
            ;;
        --completed-before)
            DATE_BEFORE="$2"
            shift 2
            ;;
        --all-stopped)
            ALL_STOPPED=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "❌ Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Get all contexts
ALL_CONTEXTS=$(my-context list 2>/dev/null || true)

# Filter stopped contexts
STOPPED_CONTEXTS=$(echo "$ALL_CONTEXTS" | grep "○" || true)

if [[ -z "$STOPPED_CONTEXTS" ]]; then
    echo "📭 No stopped contexts found"
    exit 0
fi

# Apply filters
FILTERED_CONTEXTS="$STOPPED_CONTEXTS"

if [[ -n "$PATTERN" ]]; then
    # Filter by pattern (glob-style)
    FILTERED_CONTEXTS=$(echo "$FILTERED_CONTEXTS" | grep "$PATTERN" || true)
fi

if [[ -n "$DATE_BEFORE" ]]; then
    # Simplified date filtering (would need more robust parsing in production)
    echo "⚠️  Date filtering not fully implemented in POC - showing all stopped contexts"
    # In real implementation: parse transition times and filter by date
fi

if [[ "$ALL_STOPPED" == "true" ]]; then
    # Already have all stopped contexts
    :
fi

# Check if any contexts match
if [[ -z "$FILTERED_CONTEXTS" ]]; then
    echo "📭 No contexts match criteria"
    exit 0
fi

# Extract context names
CONTEXT_NAMES=$(echo "$FILTERED_CONTEXTS" | sed 's/.*○ \(.*\) (stopped.*/\1/')
CONTEXT_COUNT=$(echo "$CONTEXT_NAMES" | grep -v "^$" | wc -l)

# Show what would be archived
echo "Found $CONTEXT_COUNT stopped context(s):"
echo ""

while IFS= read -r context; do
    [[ -z "$context" ]] && continue
    # Try to get note count (simplified - would need better parsing)
    echo "  ○ $context"
done <<< "$CONTEXT_NAMES"

echo ""

if [[ "$DRY_RUN" == "true" ]]; then
    echo "🔍 DRY RUN - no changes made"
    echo ""
    echo "Would archive $CONTEXT_COUNT context(s)"
    exit 0
fi

# Confirm before archiving
read -p "Archive all $CONTEXT_COUNT contexts? [y/N]: " CONFIRM
CONFIRM=${CONFIRM:-N}

if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 0
fi

# Archive each context
ARCHIVED_COUNT=0
FAILED_COUNT=0

while IFS= read -r context; do
    [[ -z "$context" ]] && continue

    if my-context archive "$context" 2>/dev/null; then
        echo "  ✓ Archived: $context"
        ARCHIVED_COUNT=$((ARCHIVED_COUNT + 1))
    else
        echo "  ✗ Failed: $context"
        FAILED_COUNT=$((FAILED_COUNT + 1))
    fi
done <<< "$CONTEXT_NAMES"

echo ""
echo "✅ Archived $ARCHIVED_COUNT context(s)"
if [[ $FAILED_COUNT -gt 0 ]]; then
    echo "⚠️  Failed to archive $FAILED_COUNT context(s)"
fi

exit 0
