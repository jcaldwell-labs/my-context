# Context Archival: Implementation Guide
## Step-by-step instructions to organize your 166 contexts

This guide walks you through implementing the archival strategy with copy-paste scripts.

---

## Step 1: Backup Everything (Critical!)

**Time: 5 minutes**

Before we move anything, create a complete backup.

```bash
#!/bin/bash
# BACKUP-BEFORE-ARCHIVAL.sh
# Run this first, no matter what

BACKUP_DIR="$HOME/my-context-backups/$(date +%Y-%m-%d_%H-%M-%S)"
mkdir -p "$BACKUP_DIR"

echo "Backing up ~/.my-context to $BACKUP_DIR"
cp -r ~/.my-context "$BACKUP_DIR/my-context-full"
echo "✅ Full backup created: $BACKUP_DIR/my-context-full"

# Also create compressed backup
tar czf "$BACKUP_DIR/my-context-compressed.tar.gz" ~/.my-context/
echo "✅ Compressed backup created: $BACKUP_DIR/my-context-compressed.tar.gz"

# Store backup location
echo "$BACKUP_DIR" > ~/.my-context/.last-backup-location
echo ""
echo "Backup location saved to: ~/.my-context/.last-backup-location"
```

**Run it:**
```bash
bash BACKUP-BEFORE-ARCHIVAL.sh
```

**Verify:**
```bash
ls -lh ~/my-context-backups/*/my-context-compressed.tar.gz
# Should show your backup
```

---

## Step 2: Create Archive Directory Structure

**Time: 2 minutes**

```bash
#!/bin/bash
# CREATE-ARCHIVE-STRUCTURE.sh

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"

echo "Creating archive directory structure..."

# Create main archive directory
mkdir -p "$ARCHIVE_HOME/.archive"

# Create subdirectories
mkdir -p "$ARCHIVE_HOME/.archive/metadata"
mkdir -p "$ARCHIVE_HOME/.archive/experiments"
mkdir -p "$ARCHIVE_HOME/.archive/cold-storage"

# Create monthly directories for last 12 months
echo "Creating monthly directories..."
for i in {0..11}; do
  MONTH=$(date -d "$i months ago" +%Y-%m)
  mkdir -p "$ARCHIVE_HOME/.archive/$MONTH"
done

# Create metadata files
cat > "$ARCHIVE_HOME/.archive/README.md" << 'EOF'
# My-Context Archive

This directory contains archived contexts organized by date.

## Structure
- `metadata/` - Archive metadata and recovery guides
- `experiments/` - Test and POC contexts
- `cold-storage/` - Compressed old archives
- `YYYY-MM/` - Monthly archives

## How to Restore
See `metadata/RECOVERY-GUIDE.md`

## Archive Index
See `metadata/ARCHIVE-INDEX.md`
EOF

cat > "$ARCHIVE_HOME/.archive/metadata/RECOVERY-GUIDE.md" << 'EOF'
# Recovery Guide

## Restore a Single Context from Monthly Archive

```bash
# Find the context
ls ~/.my-context/.archive/2025-10/ | grep "your-context"

# Restore it
mv ~/.my-context/.archive/2025-10/Your-Context ~/.my-context/

# Verify
my-context show "Your-Context"
```

## Restore from Cold Storage (Compressed)

```bash
# Extract
tar xzf ~/.my-context/.archive/cold-storage/q3-2025.tar.gz -C ~/.my-context/.archive/

# Move to active
mv ~/.my-context/.archive/2025-09/Your-Context ~/.my-context/
```

## Restore Entire Month

```bash
# Copy all contexts from a month back to active
cp -r ~/.my-context/.archive/2025-09/* ~/.my-context/

# This recreates all contexts from that month
my-context list | grep "2025-09"
```

## Backup Current Archive

```bash
# Make a backup before restoration
tar czf ~/archive-backup-$(date +%Y-%m-%d).tar.gz ~/.my-context/.archive/
```
EOF

cat > "$ARCHIVE_HOME/.archive/metadata/ARCHIVAL-LOG.md" << 'EOF'
# Archival Log

Track what was archived and when.

## Format
```
Date: YYYY-MM-DD
Action: [ARCHIVE | DELETE | COMPRESS | RESTORE]
Count: [number of contexts]
Reason: [reason for action]
Details: [what was archived]
```

## Entries

[Will be populated as you archive]
EOF

echo "✅ Archive structure created at: $ARCHIVE_HOME/.archive"
echo ""
echo "Directory structure:"
tree -L 2 "$ARCHIVE_HOME/.archive"
```

**Run it:**
```bash
bash CREATE-ARCHIVE-STRUCTURE.sh
```

**Verify:**
```bash
ls -la ~/.my-context/.archive/
# Should see: metadata, experiments, cold-storage, and YYYY-MM folders
```

---

## Step 3: Analyze Current State

**Time: 5 minutes**

```bash
#!/bin/bash
# ANALYZE-CONTEXTS.sh
# See what you're working with

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"

echo "=== MY-CONTEXT ANALYSIS REPORT ==="
echo "Generated: $(date)"
echo ""

echo "SUMMARY"
echo "------"
TOTAL=$(find "$ARCHIVE_HOME" -maxdepth 1 -type d ! -name ".archive" ! -name ".git" | wc -l)
echo "Total contexts: $((TOTAL - 1))"  # -1 for parent

echo ""
echo "BY AGE"
echo "-----"
RECENT=$(find "$ARCHIVE_HOME" -maxdepth 1 -type d -newermt "7 days ago" ! -name ".archive" ! -name ".git" | wc -l)
echo "<7 days old:      $((RECENT - 1)) (keep active)"

WEEK=$(find "$ARCHIVE_HOME" -maxdepth 1 -type d ! -newermt "7 days ago" -newermt "30 days ago" ! -name ".archive" ! -name ".git" | wc -l)
echo "7-30 days old:    $WEEK (review for archival)"

MONTH=$(find "$ARCHIVE_HOME" -maxdepth 1 -type d ! -newermt "30 days ago" -newermt "90 days ago" ! -name ".archive" ! -name ".git" | wc -l)
echo "30-90 days old:   $MONTH (archive soon)"

OLD=$(find "$ARCHIVE_HOME" -maxdepth 1 -type d ! -newermt "90 days ago" ! -name ".archive" ! -name ".git" | wc -l)
echo ">90 days old:     $((OLD - 1)) (ARCHIVE NOW)"

echo ""
echo "BY PATTERN"
echo "---------"
TEST=$(ls "$ARCHIVE_HOME" 2>/dev/null | grep -E "^test|^worktree|^poc|^spike|^○" | wc -l)
echo "Test/experiment:  $TEST (move to experiments/)"

DUP=$(ls "$ARCHIVE_HOME" 2>/dev/null | grep "_2\|_3\|_4\|_5" | wc -l)
echo "Potential dups:   $DUP (consolidate)"

SPRINT=$(ls "$ARCHIVE_HOME" 2>/dev/null | grep -i sprint | wc -l)
echo "Sprint contexts:  $SPRINT (archive by sprint)"

echo ""
echo "BY SIZE"
echo "------"
echo "Total disk usage:"
du -sh "$ARCHIVE_HOME" --exclude=".archive"

echo ""
echo "Largest contexts:"
du -sh "$ARCHIVE_HOME"/* --exclude=".archive" 2>/dev/null | sort -rh | head -10

echo ""
echo "RECOMMENDATIONS"
echo "---------------"
if [ $((OLD - 1)) -gt 50 ]; then
  echo "⚠️  You have many old contexts (>90 days): $((OLD - 1))"
  echo "    Recommend archiving all of these"
fi

if [ $TEST -gt 20 ]; then
  echo "⚠️  You have many test contexts: $TEST"
  echo "    Consider moving to experiments/"
fi

if [ $DUP -gt 30 ]; then
  echo "⚠️  Many potential duplicates: $DUP"
  echo "    Recommend consolidation"
fi
```

**Run it:**
```bash
bash ANALYZE-CONTEXTS.sh
```

**Expected output:** Shows you exactly what to archive

---

## Step 4: Test Run (Single Context)

**Time: 5 minutes**

Before doing bulk operations, test with one non-critical context.

```bash
#!/bin/bash
# TEST-ARCHIVAL-SINGLE.sh
# Test archival with one context

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"
TEST_CONTEXT="worktree-switch_test1-phase-worktree"  # Pick an old test one

if [ ! -d "$ARCHIVE_HOME/$TEST_CONTEXT" ]; then
  echo "❌ Context not found: $TEST_CONTEXT"
  echo "Choose an actual context name"
  exit 1
fi

echo "Testing archival with: $TEST_CONTEXT"
echo ""

# Step 1: Export before archival
echo "1️⃣  Exporting context..."
mkdir -p ~/test-archive-backup
my-context export "$TEST_CONTEXT" --to ~/test-archive-backup/$TEST_CONTEXT.md
echo "   ✓ Exported to ~/test-archive-backup/$TEST_CONTEXT.md"

# Step 2: Move to archive
echo "2️⃣  Moving to archive..."
mkdir -p "$ARCHIVE_HOME/.archive/experiments"
mv "$ARCHIVE_HOME/$TEST_CONTEXT" "$ARCHIVE_HOME/.archive/experiments/"
echo "   ✓ Moved to archive/experiments/"

# Step 3: Verify not in active list
echo "3️⃣  Verifying removal from active..."
if my-context list 2>/dev/null | grep -q "$TEST_CONTEXT"; then
  echo "   ❌ Still in active list! Something went wrong"
  exit 1
else
  echo "   ✓ Not in active list"
fi

# Step 4: Restore it
echo "4️⃣  Testing restore..."
mv "$ARCHIVE_HOME/.archive/experiments/$TEST_CONTEXT" "$ARCHIVE_HOME/"
echo "   ✓ Restored to active"

# Step 5: Verify restoration
if my-context list 2>/dev/null | grep -q "$TEST_CONTEXT"; then
  echo "   ✓ Back in active list"
else
  echo "   ❌ Restore failed!"
  exit 1
fi

echo ""
echo "✅ TEST PASSED - Archival system works!"
echo ""
echo "Next steps: Run bulk archival scripts"
```

**Run it:**
```bash
bash TEST-ARCHIVAL-SINGLE.sh
```

---

## Step 5: Move Test & Experimental Contexts

**Time: 10 minutes**

```bash
#!/bin/bash
# ARCHIVE-EXPERIMENTS.sh
# Move all test/experimental contexts to experiments/

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"
DRY_RUN=true  # Set to false to actually move

if [ "$DRY_RUN" = true ]; then
  echo "🔍 DRY RUN MODE - No changes will be made"
  echo "Change DRY_RUN=false in script to actually move"
fi

echo "Finding experimental contexts..."
echo ""

COUNT=0

# Test contexts
for ctx in "$ARCHIVE_HOME"/test-* "$ARCHIVE_HOME"/poc-* "$ARCHIVE_HOME"/spike-*; do
  [ -d "$ctx" ] || continue
  ctx_name=$(basename "$ctx")
  COUNT=$((COUNT + 1))

  if [ "$DRY_RUN" = true ]; then
    echo "Would move: $ctx_name"
  else
    mkdir -p "$ARCHIVE_HOME/.archive/experiments"
    mv "$ctx" "$ARCHIVE_HOME/.archive/experiments/"
    echo "✓ Moved: $ctx_name"
  fi
done

# Worktree test contexts
for ctx in "$ARCHIVE_HOME"/worktree-* "$ARCHIVE_HOME"/○_*; do
  [ -d "$ctx" ] || continue
  ctx_name=$(basename "$ctx")
  COUNT=$((COUNT + 1))

  if [ "$DRY_RUN" = true ]; then
    echo "Would move: $ctx_name"
  else
    mkdir -p "$ARCHIVE_HOME/.archive/experiments"
    mv "$ctx" "$ARCHIVE_HOME/.archive/experiments/"
    echo "✓ Moved: $ctx_name"
  fi
done

echo ""
echo "Total to move: $COUNT"

if [ "$DRY_RUN" = true ]; then
  echo ""
  echo "Change DRY_RUN=false and run again to actually move"
fi
```

**Run it first in dry-run mode:**
```bash
bash ARCHIVE-EXPERIMENTS.sh
```

**If output looks good, enable actual moves:**
```bash
# Edit script: DRY_RUN=false
# Then run again
bash ARCHIVE-EXPERIMENTS.sh
```

---

## Step 6: Archive Old Contexts (>90 days)

**Time: 15 minutes**

```bash
#!/bin/bash
# ARCHIVE-OLD-CONTEXTS.sh
# Move contexts >90 days old to monthly folders

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"

echo "Archiving contexts older than 90 days..."
echo ""

COUNT=0
for ctx in "$ARCHIVE_HOME"/*; do
  [ -d "$ctx" ] || continue
  [ "$ctx" != "$ARCHIVE_HOME/.archive" ] || continue
  [ "$ctx" != "$ARCHIVE_HOME/.git" ] || continue

  ctx_name=$(basename "$ctx")

  # Check if >90 days old
  if ! find "$ctx" -maxdepth 0 -newermt "90 days ago" -exec false {} +; then
    # Get the month from the file's date
    file_month=$(stat -c %y "$ctx" | cut -d' ' -f1 | cut -d'-' -f1-2)

    # Ensure archive month directory exists
    mkdir -p "$ARCHIVE_HOME/.archive/$file_month"

    # Move it
    mv "$ctx" "$ARCHIVE_HOME/.archive/$file_month/" 2>/dev/null && {
      echo "✓ $ctx_name → $file_month"
      COUNT=$((COUNT + 1))
    } || {
      echo "⚠ Failed to move: $ctx_name"
    }
  fi
done

echo ""
echo "✅ Archived $COUNT contexts"
echo ""
echo "Active contexts remaining:"
ls -d ~/.my-context/*/ 2>/dev/null | grep -v ".archive" | wc -l
```

**Run it:**
```bash
bash ARCHIVE-OLD-CONTEXTS.sh
```

---

## Step 7: Consolidate Duplicates

**Time: 10 minutes**

```bash
#!/bin/bash
# CONSOLIDATE-DUPLICATES.sh
# Keep only the NEWEST version of duplicated contexts

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"

echo "Finding and consolidating duplicate contexts..."
echo ""

# Find contexts with _2, _3, _4 suffixes
declare -A seen

for ctx in "$ARCHIVE_HOME"/*; do
  [ -d "$ctx" ] || continue

  ctx_name=$(basename "$ctx")
  base_name=$(echo "$ctx_name" | sed 's/_[0-9]*$//')

  if [ -z "${seen[$base_name]}" ]; then
    seen[$base_name]="$ctx_name"
  else
    # We have a duplicate - move the older one to experiments
    older="${seen[$base_name]}"
    newer="$ctx_name"

    older_time=$(stat -c %Y "$ARCHIVE_HOME/$older")
    newer_time=$(stat -c %Y "$ARCHIVE_HOME/$newer")

    if [ "$older_time" -lt "$newer_time" ]; then
      to_archive="$older"
      to_keep="$newer"
    else
      to_archive="$newer"
      to_keep="$older"
    fi

    echo "Archiving: $to_archive (keeping: $to_keep)"

    mkdir -p "$ARCHIVE_HOME/.archive/experiments"
    mv "$ARCHIVE_HOME/$to_archive" "$ARCHIVE_HOME/.archive/experiments/"
  fi
done

echo ""
echo "✅ Duplicates consolidated"
```

**Run it:**
```bash
bash CONSOLIDATE-DUPLICATES.sh
```

---

## Step 8: Setup Monthly Automatic Archival

**Time: 5 minutes**

```bash
#!/bin/bash
# ~/.local/bin/archive-contexts-monthly.sh
# Add to cron: 0 0 1 * * ~/.local/bin/archive-contexts-monthly.sh

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"
LOG_FILE="$ARCHIVE_HOME/.archive/metadata/ARCHIVAL-LOG.md"

{
  echo "## Archive Run: $(date)"
  echo "Contexts before: $(ls -d $ARCHIVE_HOME/*/ 2>/dev/null | grep -v '.archive' | wc -l)"
} >> "$LOG_FILE"

# Archive contexts 30+ days old
for ctx in "$ARCHIVE_HOME"/*; do
  [ -d "$ctx" ] || continue
  [ "$ctx" != "$ARCHIVE_HOME/.archive" ] || continue

  ctx_name=$(basename "$ctx")

  # Skip if <30 days old
  if find "$ctx" -maxdepth 0 -newermt "30 days ago" | grep -q .; then
    continue
  fi

  # Check if it's stopped (not active)
  if [ -f "$ctx/meta.json" ] && grep -q '"active": true' "$ctx/meta.json" 2>/dev/null; then
    continue
  fi

  # Archive it
  file_month=$(stat -c %y "$ctx" | cut -d' ' -f1 | cut -d'-' -f1-2)
  mkdir -p "$ARCHIVE_HOME/.archive/$file_month"
  mv "$ctx" "$ARCHIVE_HOME/.archive/$file_month/" 2>/dev/null && {
    echo "  - Archived: $ctx_name → $file_month" >> "$LOG_FILE"
  }
done

{
  echo "Contexts after: $(ls -d $ARCHIVE_HOME/*/ 2>/dev/null | grep -v '.archive' | wc -l)"
  echo ""
} >> "$LOG_FILE"

# Also show current status
echo "[$(date)] Monthly archive complete" >> "$LOG_FILE"
```

**Install to cron:**
```bash
# Make it executable
chmod +x ~/.local/bin/archive-contexts-monthly.sh

# Add to crontab (1st of each month at midnight)
(crontab -l 2>/dev/null; echo "0 0 1 * * ~/.local/bin/archive-contexts-monthly.sh") | crontab -

# Verify
crontab -l | grep archive
```

---

## Complete Execution Plan

### Your Situation: 166 Contexts

Run these in order:

```bash
# 1. BACKUP EVERYTHING (Do this first!)
bash BACKUP-BEFORE-ARCHIVAL.sh

# 2. Create archive structure
bash CREATE-ARCHIVE-STRUCTURE.sh

# 3. Analyze what you have
bash ANALYZE-CONTEXTS.sh

# 4. Test archival works
bash TEST-ARCHIVAL-SINGLE.sh

# 5. Archive experimental contexts
bash ARCHIVE-EXPERIMENTS.sh

# 6. Archive old contexts (>90 days)
bash ARCHIVE-OLD-CONTEXTS.sh

# 7. Consolidate duplicates
bash CONSOLIDATE-DUPLICATES.sh

# 8. Setup monthly automation
# (See step 8 above)

# 9. Verify the result
my-context list | wc -l  # Should be much smaller
ls -la ~/.my-context/.archive/ | wc -l  # Should have organized archives
```

---

## Expected Results

### Before Archival
```
Total: 166 contexts
Active list: overwhelming
Duplicates: many
Test/POC: 20+
Age: mixed (0-365 days)
```

### After Archival
```
Active: ~20-30 contexts (last 30 days)
~/.my-context/.archive/
  ├── experiments/: test/POC contexts
  ├── YYYY-MM/: monthly archives
  ├── cold-storage/: compressed old
  └── metadata/: recovery guides
```

---

## If Something Goes Wrong

**Lost a context?**
```bash
# Restore from backup
tar xzf ~/my-context-backups/[date]/my-context-compressed.tar.gz -C ~/.my-context/ --strip-components=1
```

**Archive structure messed up?**
```bash
# Start over from backup
rm -rf ~/.my-context/.archive
bash CREATE-ARCHIVE-STRUCTURE.sh
# Re-run archival scripts
```

**Can't find a context?**
```bash
# Search everywhere
find ~/.my-context -name "*part-of-name*" -type d
```

---

## Next: Ongoing Maintenance

Once archival is complete, see **ARCHIVAL-MAINTENANCE.md** for:
- Monthly routine
- Quarterly cold-storage compression
- Recovery procedures
- Monitoring active count

---

**Ready? Start with Step 1!** 🚀
