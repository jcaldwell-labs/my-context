# my-context Feature Request POCs

**Purpose**: Proof-of-concept shell script implementations of feature requests identified from real-world usage

**Status**: POC - Validate UX and behavior before implementing in Go

---

## Feature Requests Implemented

### FR-MC-NEW-001: Smart Resume on Start
**Script**: `smart-resume.sh`
**Problem**: Starting context with same name as stopped context causes confusion
**Solution**: Detects duplicate names and prompts to resume or create new

### FR-MC-NEW-002: Context Lifecycle Advisor
**Script**: `lifecycle-advisor.sh`
**Problem**: Unclear when to stop/start/archive contexts
**Solution**: Provides guidance and suggestions after stopping context

### FR-MC-NEW-003: Bulk Archive with Patterns
**Script**: `bulk-archive.sh`
**Problem**: Manual archiving of many contexts is tedious
**Solution**: Archive multiple contexts matching pattern or criteria

### FR-MC-NEW-004: Note Limit Warning
**Script**: `note-with-warning.sh`
**Problem**: No indication when context getting large
**Solution**: Warns at thresholds (50, 100, 200+ notes)

### FR-MC-NEW-005: Resume Command Alias
**Script**: `resume-alias.sh`
**Problem**: `my-context start` for existing context semantically confusing
**Solution**: Explicit `resume` command for resuming stopped contexts

---

## Installation (POC Testing)

### Option 1: Direct Usage
```bash
cd ~/projects/my-context-dev/scripts/poc

# Make scripts executable
chmod +x *.sh

# Use directly
./smart-resume.sh "my-context-name"
./note-with-warning.sh "my note message"
./resume-alias.sh --last
./lifecycle-advisor.sh
./bulk-archive.sh --pattern "test-*" --dry-run
```

### Option 2: Create Aliases (Recommended)
```bash
# Add to ~/.bashrc or ~/.zshrc
alias mcx-start='~/projects/my-context-dev/scripts/poc/smart-resume.sh'
alias mcx-note='~/projects/my-context-dev/scripts/poc/note-with-warning.sh'
alias mcx-resume='~/projects/my-context-dev/scripts/poc/resume-alias.sh'
alias mcx-stop='~/projects/my-context-dev/scripts/poc/lifecycle-advisor.sh'
alias mcx-bulk-archive='~/projects/my-context-dev/scripts/poc/bulk-archive.sh'

# Reload shell
source ~/.bashrc  # or source ~/.zshrc

# Use with short aliases
mcx-start "foo"
mcx-note "implemented feature X"
mcx-resume --last
mcx-stop
mcx-bulk-archive --pattern "006-*"
```

### Option 3: Install to ~/bin (Optional)
```bash
mkdir -p ~/bin
cd ~/projects/my-context-dev/scripts/poc

# Create symlinks with 'mcx-' prefix
ln -sf "$(pwd)/smart-resume.sh" ~/bin/mcx-start
ln -sf "$(pwd)/note-with-warning.sh" ~/bin/mcx-note
ln -sf "$(pwd)/resume-alias.sh" ~/bin/mcx-resume
ln -sf "$(pwd)/lifecycle-advisor.sh" ~/bin/mcx-stop
ln -sf "$(pwd)/bulk-archive.sh" ~/bin/mcx-bulk-archive

# Ensure ~/bin is in PATH
export PATH="$HOME/bin:$PATH"  # Add to ~/.bashrc if not already there

# Use from anywhere
mcx-start "my-context"
mcx-note "my note"
```

---

## Usage Examples

### FR-MC-NEW-001: Smart Resume on Start

**Scenario 1: Resume existing stopped context**
```bash
$ mcx-start "deb-sanity: 007-github"
📋 Context 'deb-sanity: 007-github' exists (stopped)

   Notes: 25

Resume existing context? [Y/n]: y
✅ Resumed context: deb-sanity: 007-github
```

**Scenario 2: Create new context with different name**
```bash
$ mcx-start "deb-sanity: 007-github"
📋 Context 'deb-sanity: 007-github' exists (stopped)

   Notes: 25

Resume existing context? [Y/n]: n

Start new context with different name (e.g., deb-sanity: 007-github_2): deb-sanity: 007-github-phase2
Started new context: deb-sanity: 007-github-phase2
```

**Scenario 3: Force new context (bypass prompt)**
```bash
$ mcx-start "foo" --force
⚠️  Context 'foo' exists (stopped) but --force specified

Create new context with different name (e.g., foo_2): foo-new
Started new context: foo-new
```

---

### FR-MC-NEW-002: Context Lifecycle Advisor

**Usage**:
```bash
$ mcx-stop
Stopped context: deb-sanity: 006-sprint-006 (duration: 2h 15m)

📊 Context Summary:
   Name: deb-sanity: 006-sprint-006
   Duration: 2h 15m
   Notes: 42

📚 Related contexts:
  ○ deb-sanity: 006-sprint-006-phase-1 (stopped)
  ○ deb-sanity: 006-sprint-006-phase-2 (stopped)

💡 Suggestions:
   (1) Resume related: my-context start "deb-sanity: 006-sprint-006-phase-1"
   (2) Archive if complete: my-context archive "deb-sanity: 006-sprint-006"
   (3) Start new work: my-context start <name>

🎯 Completion detected in recent notes!
   Consider archiving: my-context archive "deb-sanity: 006-sprint-006"
```

---

### FR-MC-NEW-003: Bulk Archive with Patterns

**Scenario 1: Dry run (preview)**
```bash
$ mcx-bulk-archive --pattern "006-*" --dry-run
Found 16 stopped context(s):

  ○ 006-sprint-006-large
  ○ 006-sprint-006-large_2
  ○ 006-sprint-006-phase-analysis
  ... (13 more)

🔍 DRY RUN - no changes made

Would archive 16 context(s)
```

**Scenario 2: Archive all matching**
```bash
$ mcx-bulk-archive --pattern "006-*"
Found 16 stopped context(s):

  ○ 006-sprint-006-large
  ○ 006-sprint-006-large_2
  ... (14 more)

Archive all 16 contexts? [y/N]: y
  ✓ Archived: 006-sprint-006-large
  ✓ Archived: 006-sprint-006-large_2
  ... (14 more)

✅ Archived 16 context(s)
```

**Scenario 3: Archive all stopped**
```bash
$ mcx-bulk-archive --all-stopped
Found 48 stopped context(s):

  ○ test-context-1
  ○ old-experiment
  ... (46 more)

Archive all 48 contexts? [y/N]: y
  ✓ Archived: test-context-1
  ... (47 more)

✅ Archived 48 context(s)
```

---

### FR-MC-NEW-004: Note Limit Warning

**Usage**:
```bash
$ mcx-note "Implemented smart resume feature"
Note added to context: deb-sanity: 007-github

# At 50 notes:
$ mcx-note "Completed implementation phase"
Note added to context: deb-sanity: 007-github

⚠️  Context has 50 notes. Consider:
   - Stop and start new phase context if switching focus
   - Continue if still same work (100+ notes is fine)
   - Export to markdown: my-context export <context-name>

# At 100 notes:
$ mcx-note "Another milestone"
Note added to context: deb-sanity: 007-github

⚠️  Context has 100 notes (getting large). Consider:
   - Export to markdown for review: my-context export <context-name>
   - Split into phase contexts if work has shifted
   - Or continue - large contexts are OK for complex work
```

**Configure thresholds**:
```bash
# Custom warning thresholds
export MC_WARN_AT=30      # First warning at 30 notes
export MC_WARN_AT_2=75    # Second warning at 75 notes
export MC_WARN_AT_3=150   # Third warning at 150 notes

mcx-note "your note"
```

---

### FR-MC-NEW-005: Resume Command Alias

**Scenario 1: Resume specific context**
```bash
$ mcx-resume "deb-sanity: 007-github"
Started context: deb-sanity: 007-github
✅ Resumed context: deb-sanity: 007-github
```

**Scenario 2: Resume last stopped**
```bash
$ mcx-resume --last
🔄 Resuming most recent: deb-sanity: 007-github-phase2
Started context: deb-sanity: 007-github-phase2
```

**Scenario 3: Resume with pattern (multiple matches)**
```bash
$ mcx-resume "007-*"
Multiple matching contexts:

  (1) deb-sanity: 007-github
  (2) ps-cli: 007-integration

Select context [1-2]: 1
🔄 Resuming: deb-sanity: 007-github
Started context: deb-sanity: 007-github
```

**Scenario 4: Resume with pattern (single match)**
```bash
$ mcx-resume "007-github*"
🔄 Resuming: deb-sanity: 007-github-phase2
Started context: deb-sanity: 007-github-phase2
```

---

## Validation Checklist

### FR-MC-NEW-001: Smart Resume ✅
- [ ] Detects duplicate stopped context names
- [ ] Shows context summary (notes, last active)
- [ ] Prompts to resume or create new
- [ ] `--force` bypasses prompt
- [ ] Error if context is active

### FR-MC-NEW-002: Lifecycle Advisor ✅
- [ ] Shows context summary after stop
- [ ] Detects related contexts
- [ ] Provides resume/archive/start suggestions
- [ ] Detects completion keywords in notes

### FR-MC-NEW-003: Bulk Archive ✅
- [ ] Pattern matching works
- [ ] Dry-run mode previews changes
- [ ] Confirmation prompt before archiving
- [ ] Counts archived and failed contexts
- [ ] `--all-stopped` archives all

### FR-MC-NEW-004: Note Warning ✅
- [ ] Warns at 50 notes (default)
- [ ] Warns at 100 notes
- [ ] Periodic warnings at 200+
- [ ] Configurable thresholds

### FR-MC-NEW-005: Resume Command ✅
- [ ] Resume specific context
- [ ] `--last` resumes most recent
- [ ] Pattern matching with single result
- [ ] Pattern matching with multiple (user selects)
- [ ] Error handling for non-existent contexts

---

## Next Steps

### Immediate (Sprint 007)
1. **Test POCs in real usage** - Use during Sprint 007 development
2. **Gather feedback** - Note friction points and UX issues
3. **Iterate on scripts** - Improve based on real usage

### Short-Term (Sprint 008)
1. **Validate approach** - Confirm features solve problems
2. **Implement in Go** - Port validated features to my-context codebase
3. **Add tests** - Unit and integration tests for new features

### Long-Term (v2.0.0)
1. **Release as v2.0.0** - All 5 feature requests implemented
2. **Document best practices** - Based on validated usage patterns
3. **Community feedback** - Share with broader my-context users

---

## Known Limitations (POC)

1. **Parsing**: Simplified grep/sed parsing (Go implementation will be more robust)
2. **Date filtering**: Not fully implemented in bulk-archive (would need transition log parsing)
3. **Error handling**: Basic (production would have comprehensive error handling)
4. **Performance**: Shell script overhead (Go will be faster)
5. **Edge cases**: May not handle all my-context output variations

**These are POCs** - Meant to validate UX and behavior, not production-ready implementations.

---

## Troubleshooting

### Scripts not working
```bash
# Ensure scripts are executable
chmod +x ~/projects/my-context-dev/scripts/poc/*.sh

# Check my-context is available
which my-context

# Test basic my-context commands
my-context list
```

### Aliases not found
```bash
# Check if aliases are defined
alias | grep mcx

# Reload shell config
source ~/.bashrc  # or source ~/.zshrc
```

### Pattern matching not working
```bash
# Use quotes around patterns with wildcards
mcx-bulk-archive --pattern "006-*"  # ✓ Correct
mcx-bulk-archive --pattern 006-*    # ✗ Shell expansion may interfere
```

---

**Document Version**: 1.0.0
**Created**: 2025-10-10
**Purpose**: POC validation of my-context feature requests
**Status**: Ready for real-world testing
