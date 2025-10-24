# My-Context Archival Strategy
## Managing 100+ Contexts Across the Lifecycle

> **Problem**: Your context folder has accumulated 166+ contexts over months of development.
> Many are old, experimental, or duplicates. You need them preserved, but they're cluttering
> your active workspace.
>
> **Solution**: A three-tier archival system with automatic organization, preservation, and recovery.

---

## The Problem We're Solving

### Current State
- **166 contexts** in `~/.my-context/`
- **2.9 MB** of historical data
- **Duplicates** (e.g., `002__Tech_Debt_Resolution`, `002__Tech_Debt_Resolution_2`, etc.)
- **Test contexts** (`worktree-*`, `○_worktree-*`, test phases)
- **Old sprints** (Sprint 6, 7, etc. - pre-Sprint lifecycle system)
- **Experimental** branches and explorations
- **Cluttered workspace** - hard to see current work
- **No clear lifecycle** - when should something be archived?

### Consequences
- `my-context list` shows 166 items (overwhelming)
- Hard to find active work in the noise
- Can't tell what's old vs. current at a glance
- Potential confusion about what's been completed
- No clear deprecation path for failed experiments

---

## The Solution: Three-Tier Archival System

### Tier 1: Active Workspace
**Location**: `~/.my-context/` (current)
**Retention**: Last 30 days of active work
**Properties**:
- Only current sprint + last week
- Visible by default in `list` command
- <20 contexts typically
- Easy to scan and find current work

### Tier 2: Monthly Archive
**Location**: `~/.my-context/.archive/YYYY-MM/`
**Retention**: By month (3-6 months kept)
**Properties**:
- Completed sprints
- Finished projects
- Work from previous months
- Searchable and exportable
- Organized by calendar month

### Tier 3: Long-Term Cold Storage
**Location**: `~/.my-context/.archive/cold-storage/`
**Retention**: Historical reference (indefinite)
**Properties**:
- Very old work (>6 months)
- Consolidated into quarterly bundles
- Compressed for space efficiency
- For historical reference only
- Recovery requires extraction

### Special: Experimental Sandbox
**Location**: `~/.my-context/.archive/experiments/`
**Retention**: As needed during development
**Properties**:
- Test contexts
- Failed attempts
- Proof-of-concepts
- Can be cleaned up after feature ships

---

## Archival Rules by Context Type

### Rule 1: Completed Work (Age-Based)
**When**: Context is stopped and >30 days old
**Where**: Tier 2 (monthly archive)
**How**: Auto-archive monthly or manual cleanup

```
Decision Tree:
├─ Stopped and <7 days old?   → Stay in Tier 1 (active)
├─ Stopped and 7-30 days old? → Stay in Tier 1 (active, review soon)
├─ Stopped and 30-90 days old? → Archive to Tier 2 (monthly folder)
├─ Stopped and 90-180 days old? → Keep in Tier 2
└─ Stopped and >180 days old?  → Move to Tier 3 (cold storage)
```

### Rule 2: Experimental/Test Contexts
**When**: Context name indicates experiment
**Patterns**:
- `test-*`
- `worktree-*`
- `poc-*`
- `spike-*`
- `experiment-*`
- `○_*` (marker for worktrees)

**When to Archive**:
- Immediately if feature is rejected
- After 14 days if feature ships (consolidate into real context)
- Quarterly cleanup if still experimental

**Location**: `~/.my-context/.archive/experiments/`

### Rule 3: Duplicates
**When**: You have `context-name`, `context-name_2`, `context-name_3`, etc.

**Decision**:
1. Keep the MOST RECENT one (usually `_2` or `_3`)
2. Archive others with note referencing consolidated context
3. Document reason for duplication in archive metadata

**Example**:
```
Archived: 002__Tech_Debt_Resolution
Reason: Duplicate of 002__Tech_Debt_Resolution_2_2_2
Note: Historical context - refer to _2_2_2 for complete record
```

### Rule 4: Sprint Work
**When**: Sprint is complete

**Archive by Sprint**:
- Collect all `sprint-N-*` contexts
- Create monthly folder: `2025-10-sprint-6-summary/`
- Include sprint summary + all context exports
- Link to sprint retrospective

**Location**: `~/.my-context/.archive/YYYY-MM/sprint-summaries/`

### Rule 5: Project Completion
**When**: Project shipped/completed

**Archive Actions**:
1. Export all project contexts to single markdown
2. Create project archive folder
3. Move all contexts to archive
4. Keep reference context with summary

---

## Implementation Strategy

### Phase 1: Assessment (Today)
Run analysis on current state:

```bash
#!/bin/bash
# Analyze context distribution

echo "=== Context Analysis ==="
echo ""
echo "Total contexts:"
ls -d ~/.my-context/*/ | wc -l

echo ""
echo "By age:"
find ~/.my-context -maxdepth 1 -type d -newermt "7 days ago" | wc -l
echo "  <7 days old (active)"

find ~/.my-context -maxdepth 1 -type d ! -newermt "7 days ago" -newermt "30 days ago" | wc -l
echo "  7-30 days old (review)"

find ~/.my-context -maxdepth 1 -type d ! -newermt "30 days ago" -newermt "90 days ago" | wc -l
echo "  30-90 days old (archive soon)"

find ~/.my-context -maxdepth 1 -type d ! -newermt "90 days ago" | wc -l
echo "  >90 days old (archive now)"

echo ""
echo "By pattern:"
ls -d ~/.my-context/*test* 2>/dev/null | wc -l
echo "  test/worktree contexts"

ls -d ~/.my-context/*_2* 2>/dev/null | wc -l
echo "  duplicates (_2, _3, etc.)"

ls -d ~/.my-context/*sprint* 2>/dev/null | wc -l
echo "  sprint contexts"

echo ""
echo "Disk usage:"
du -sh ~/.my-context
```

### Phase 2: Setup (Today/Tomorrow)
Create archive directory structure:

```bash
#!/bin/bash
# Create archive structure

mkdir -p ~/.my-context/.archive/{experiments,cold-storage,metadata}

# Create monthly archives for last 12 months
for month in {1..12}; do
  DATE=$(date -d "$month months ago" +%Y-%m)
  mkdir -p ~/.my-context/.archive/$DATE
done

# Create sprint summary template
cat > ~/.my-context/.archive/metadata/ARCHIVE-INDEX.md << 'EOF'
# My-Context Archive Index

## Structure
- `.archive/YYYY-MM/` - Monthly archives
- `.archive/experiments/` - Test and POC contexts
- `.archive/cold-storage/` - >6 month old consolidated archives
- `.archive/metadata/` - Archive metadata and indexes

## Recovery
See RECOVERY-GUIDE.md for how to restore archived contexts.

## Statistics
[Will be updated by archival scripts]

EOF

echo "✅ Archive structure created"
```

### Phase 3: Initial Migration (This Week)
Move 90+ day old contexts to Tier 2 & 3:

```bash
#!/bin/bash
# Archive old contexts (>90 days)

ARCHIVE_HOME="$HOME/.my-context"
OLD_THRESHOLD="90"

echo "Archiving contexts older than $OLD_THRESHOLD days..."

find "$ARCHIVE_HOME" -maxdepth 1 -type d ! -newermt "${OLD_THRESHOLD} days ago" ! -name ".archive" ! -name ".git" | while read ctx; do
  ctx_name=$(basename "$ctx")
  ctx_date=$(stat -c %y "$ctx" | cut -d' ' -f1)
  archive_month=$(date -d "$ctx_date" +%Y-%m)

  echo "Moving: $ctx_name (dated $ctx_date) → $archive_month"
  mv "$ctx" "$ARCHIVE_HOME/.archive/$archive_month/"
done

echo "✅ Initial archival complete"
```

### Phase 4: Ongoing Management
Set up cron job for automatic cleanup:

```bash
#!/bin/bash
# ~/.local/bin/archive-contexts.sh
# Run this monthly via cron: 0 0 1 * * ~/.local/bin/archive-contexts.sh

ARCHIVE_HOME="${MY_CONTEXT_HOME:=$HOME/.my-context}"

echo "=== Monthly Context Archival ==="
echo "Date: $(date)"

# Archive 30+ day old stopped contexts
echo ""
echo "Archiving contexts 30-90 days old..."
find "$ARCHIVE_HOME" -maxdepth 1 -type d \
  ! -newermt "30 days ago" \
  -newermt "90 days ago" \
  ! -name ".archive" \
  ! -name ".git" | while read ctx; do

  ctx_name=$(basename "$ctx")
  archive_month=$(date +%Y-%m)

  # Check if context is stopped
  if [ ! -f "$ctx/meta.json" ] || ! grep -q '"active": true' "$ctx/meta.json" 2>/dev/null; then
    mv "$ctx" "$ARCHIVE_HOME/.archive/$archive_month/"
    echo "  ✓ $ctx_name"
  fi
done

# Move 90+ day old to cold storage
echo ""
echo "Moving contexts >90 days to cold storage..."
find "$ARCHIVE_HOME/.archive" -maxdepth 2 -type d ! -newermt "90 days ago" | while read ctx; do
  ctx_name=$(basename "$ctx")

  # Only archive if not already in cold-storage
  if [[ ! "$ctx" =~ "cold-storage" ]] && [ -d "$ctx" ] && [ -n "$(ls -A $ctx)" ]; then
    tar czf "$ARCHIVE_HOME/.archive/cold-storage/$ctx_name-$(date +%s).tar.gz" "$ctx"
    rm -rf "$ctx"
    echo "  ✓ Compressed: $ctx_name"
  fi
done

# Generate archive index
echo ""
echo "Updating archive index..."
cat > "$ARCHIVE_HOME/.archive/INDEX.md" << 'EOF'
# Archive Index
Generated: $(date)

## Active Contexts
$(my-context list --json 2>/dev/null | jq -r '.data.contexts[].name' | head -20)

## Archive Breakdown
$(find "$ARCHIVE_HOME/.archive" -maxdepth 1 -type d -exec basename {} \; | sort)

For recovery, see RECOVERY-GUIDE.md
EOF

echo "✅ Archival complete"
```

---

## Three-Step Archival Process

### For Individual Contexts

```bash
# Step 1: Verify context is done
my-context show "Context Name"

# Step 2: Export before archiving (backup)
my-context export "Context Name" --to ~/work-backups/$(date +%Y-%m)/context-name.md

# Step 3: Archive it
my-context archive "Context Name"
# OR manually move:
# mv ~/.my-context/"Context Name" ~/.my-context/.archive/$(date +%Y-%m)/
```

### For Bulk Operations

```bash
# Archive all Sprint 6 contexts at once
for ctx in $(ls ~/.my-context | grep -i sprint-6); do
  my-context archive "$ctx"
done

# OR move old test contexts to experiments
for ctx in $(ls ~/.my-context | grep "^test-\|^worktree-"); do
  mv ~/.my-context/"$ctx" ~/.my-context/.archive/experiments/
done
```

---

## Quick Reference: Archive Commands

```bash
# See what would be archived (dry run)
find ~/.my-context -maxdepth 1 -type d ! -newermt "30 days ago" ! -name ".archive"

# Archive contexts by age
my-context list --json | jq -r '.data.contexts[] | select(.created < now - 30*24*3600) | .name'

# Export for archival
my-context export "Context Name" --to backup.md
my-context export --all --search "Sprint 5" --to sprint-5-export/

# Move to archive folder manually
mkdir -p ~/.my-context/.archive/$(date +%Y-%m)
mv ~/.my-context/"Old Context" ~/.my-context/.archive/$(date +%Y-%m)/

# Verify archive
ls ~/.my-context/.archive/$(date +%Y-%m) | wc -l
echo "Contexts archived this month"

# Cold storage (compress old archives)
tar czf ~/.my-context/.archive/cold-storage/q3-2025.tar.gz \
  ~/.my-context/.archive/2025-07 \
  ~/.my-context/.archive/2025-08 \
  ~/.my-context/.archive/2025-09
```

---

## Recovery: Bringing Back Archived Contexts

### From Monthly Archive

```bash
# Restore from monthly folder
mv ~/.my-context/.archive/2025-09/Context-Name ~/.my-context/

# It's now active again in your workspace
my-context list
```

### From Cold Storage

```bash
# Extract from compressed archive
tar xzf ~/.my-context/.archive/cold-storage/q3-2025.tar.gz -C ~/.my-context/.archive/

# Move to active workspace
mv ~/.my-context/.archive/2025-09/Context-Name ~/.my-context/

# Verify it's restored
my-context show "Context-Name"
```

### Batch Recovery

```bash
# Restore all contexts from a month
mkdir -p ~/.my-context/restored-$(date +%Y-%m)
cp -r ~/.my-context/.archive/2025-09/* ~/.my-context/restored-$(date +%Y-%m)/

# Now selective restore from there
```

---

## Your Current Situation: Action Plan

### Immediate Actions (This Week)

**Step 1: Assess** (30 min)
```bash
# Run assessment script above
```

**Step 2: Setup** (15 min)
```bash
# Create archive directory structure
mkdir -p ~/.my-context/.archive/{experiments,cold-storage,metadata}
```

**Step 3: Backup** (5 min)
```bash
# Before we move anything, backup state
cd ~/.my-context
git add . && git commit -m "pre-archival backup"
# or
tar czf ~/my-context-backup-$(date +%Y-%m-%d).tar.gz ~/.my-context/
```

**Step 4: Test Run** (15 min)
```bash
# Move just ONE old context to test
# Choose something very old and non-critical
mv ~/.my-context/test-worktree-something ~/.my-context/.archive/experiments/
# Verify my-context still works
my-context list
```

**Step 5: Initial Migration** (30 min)
```bash
# Archive all >90 day old contexts
# (See Phase 3 script above)
```

### Expected Results

**Before:**
- `my-context list` shows 166 contexts
- Hard to find current work
- Workspace cluttered
- 2.9 MB of storage

**After:**
- `my-context list` shows ~15-20 contexts (last 30 days)
- Current work is obvious
- Organized, clean workspace
- Old work preserved but not cluttering
- Easy monthly/quarterly reviews

### Your New Workflow

```
Each Month:
└─ Manually archive contexts >30 days old
   OR auto-run archive script (cron job)

Each Sprint:
└─ Archive previous sprint contexts
   ├─ Export sprint summary
   ├─ Move contexts to .archive/YYYY-MM/sprint-N/
   └─ Create sprint retrospective

Each Quarter:
└─ Compress older archives to cold-storage
   ├─ Consolidate months 1-3 into Q-bundle
   ├─ Compress to tar.gz
   ├─ Remove originals
   └─ Update index

Ongoing:
└─ Clean experiments weekly
   ├─ Delete failed POCs
   ├─ Keep successful ones as reference
   └─ Move to appropriate tier when done
```

---

## Archive Structure You'll Have

```
~/.my-context/
├── active contexts (15-20 most recent)
├── state.json
├── transitions.log
│
└── .archive/
    ├── metadata/
    │   ├── ARCHIVE-INDEX.md      # Main archive index
    │   ├── RECOVERY-GUIDE.md     # How to restore
    │   └── ARCHIVAL-LOG.md       # When/what was archived
    │
    ├── experiments/              # Test/POC contexts
    │   ├── test-worktree-*
    │   ├── poc-*
    │   └── spike-*
    │
    ├── 2025-10/                  # Monthly archives
    │   ├── sprint-summaries/
    │   │   └── sprint-6-summary.md
    │   ├── completed-feature-1/
    │   └── other-contexts/
    │
    ├── 2025-09/
    │   ├── sprint-summaries/
    │   └── contexts/
    │
    ├── cold-storage/             # Compressed old archives
    │   ├── q3-2025.tar.gz        # Jul-Sep contexts
    │   ├── q2-2025.tar.gz        # Apr-Jun contexts
    │   └── q1-2025.tar.gz        # Jan-Mar contexts
    │
    └── README.md                 # Archive guide
```

---

## Best Practices

### ✅ DO:
- Export contexts before archiving (backup)
- Archive by month, not randomly
- Keep recent (30-day) work active
- Tag important contexts for easy finding
- Create sprint summaries before archiving sprint work
- Document why test/experimental contexts failed

### ❌ DON'T:
- Archive active contexts (ones you're still using)
- Delete without exporting first
- Mix old work with current in workspace
- Forget to test recovery process
- Let experimental contexts accumulate
- Archive without backup

---

## Maintenance Schedule

```
Daily:
└─ Normal context operations (start, note, stop)

Weekly:
└─ Review active contexts
   └─ If >25 active, start planning archives

Monthly:
└─ Run archival script
├─ Archive 30+ day old contexts
├─ Export sprint summary
└─ Update archive index

Quarterly:
└─ Compress old monthly archives
├─ Consolidate to cold-storage
├─ Review and delete failed experiments
└─ Update archive statistics

Annually:
└─ Review cold-storage
├─ Consolidate very old archives
└─ Create yearly summary
```

---

## For Your Specific Situation

You have 166 contexts, many of which are:
- **Duplicates** (002_*, 006_*, 007_* with _2, _3 suffixes)
- **Experiments** (worktree-*, test-*, ○_* markers)
- **Old sprints** (Sprint 6, 7 from before lifecycle system)
- **Explorations** (various one-off investigations)

### Recommended First Action
```bash
# 1. Backup everything
tar czf ~/my-context-full-backup-$(date +%Y-%m-%d).tar.gz ~/.my-context/

# 2. Create archive structure
mkdir -p ~/.my-context/.archive/{experiments,cold-storage,metadata}

# 3. Move ALL worktree test contexts
for ctx in ~/.my-context/worktree-* ~/.my-context/○_*; do
  [ -d "$ctx" ] && mv "$ctx" ~/.my-context/.archive/experiments/
done

# 4. Move ALL old sprints (6, 7 = old system)
for ctx in ~/.my-context/*sprint-[67]*; do
  [ -d "$ctx" ] && mv "$ctx" ~/.my-context/.archive/2025-06/  # or appropriate month
done

# 5. Consolidate duplicates - KEEP ONLY _2_2_2 (the newest)
for ctx in ~/.my-context/002__Tech_Debt_Resolution{,_2,_2_2}; do
  [ -d "$ctx" ] && [ "$ctx" != "$newest" ] && mv "$ctx" ~/.my-context/.archive/experiments/
done

# 6. Verify it worked
my-context list | wc -l  # Should be much smaller now
```

This would immediately reduce noise from 166 → ~50-60 more manageable contexts.

---

## Next Steps

1. **Read this strategy** - understand the three tiers
2. **Run assessment** - see current state
3. **Create structure** - set up archive folders
4. **Test with 5 contexts** - move a few old ones, verify recovery
5. **Do initial migration** - move 90+ day old to archive
6. **Setup monthly job** - automate with cron or manual reminder
7. **Document your specific flows** - adapt to your workflows

You'll go from a cluttered 166-context workspace to an organized 20-context active space with 140+ contexts safely preserved in tiers.

---

## Questions & Answers

**Q: Will archiving break anything?**
A: No. Archived contexts are just moved to different directories. They can be recovered anytime.

**Q: How do I search archived contexts?**
A: Extract from cold-storage first, then use grep on their notes.log files, or use find with grep.

**Q: What if I archive something I still need?**
A: Recover it: `mv ~/.my-context/.archive/YYYY-MM/Context-Name ~/.my-context/`

**Q: How much space will I save?**
A: Depends on how much you archive, but cold-storage compression (~5:1) will help.
Moving 100 contexts to archive saves ~2MB in active workspace.

**Q: Can I automate this?**
A: Yes! See the cron script in Phase 4. Run monthly.

**Q: What about my git history?**
A: Archive at the filesystem level (mv command). Version control `.my-context/.git` separately.

---

**Ready to clean up? Let's make your workspace organized and your history preserved!** 🚀
