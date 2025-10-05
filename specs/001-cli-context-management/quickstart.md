# Quickstart Testing: CLI Context Management System

**Feature**: 001-cli-context-management  
**Date**: 2025-10-04  
**Purpose**: Manual testing scenarios to validate core functionality

## Prerequisites

- `my-context` binary built and in PATH
- Clean slate: Remove `~/.my-context/` if exists
- Shell: Run tests in git-bash (primary), then validate in cmd.exe and WSL

---

## Scenario 1: Basic Workflow

**Goal**: Validate complete context lifecycle with notes, files, and touches.

### Steps

```bash
# 1. Start a new context
$ my-context start "Bug fix"
Expected: "Started context: Bug fix"

# 2. Verify show displays active context
$ my-context show
Expected: Context name, start time, status=active, 0 notes/files/touches

# 3. Add notes
$ my-context note "Found authentication bug in login handler"
$ my-context note "Issue occurs when token expires"
Expected: "Note added to context: Bug fix" (x2)

# 4. Associate files
$ my-context file src/auth/login.go
$ my-context file tests/auth/login_test.go
Expected: "File associated with context: Bug fix" (x2)

# 5. Record activity touches
$ my-context touch
$ my-context touch
Expected: "Touch recorded in context: Bug fix" (x2)

# 6. Verify show displays all data
$ my-context show
Expected: 
  - 2 notes with timestamps
  - 2 files with timestamps
  - 2 touches

# 7. Stop context
$ my-context stop
Expected: "Stopped context: Bug fix (duration: Xm)"

# 8. Verify no active context
$ my-context show
Expected: "No active context"
```

### Validation

- [ ] All commands succeeded (exit code 0)
- [ ] `~/.my-context/Bug_fix/` directory created
- [ ] `meta.json` contains correct start/end times, status="stopped"
- [ ] `notes.log` contains 2 lines with timestamps
- [ ] `files.log` contains 2 lines with normalized paths
- [ ] `touch.log` contains 2 timestamps
- [ ] `state.json` shows active_context=null
- [ ] `transitions.log` contains "start" and "stop" entries

---

## Scenario 2: Context Switching

**Goal**: Validate automatic context stopping when starting new context.

### Steps

```bash
# 1. Start first context
$ my-context start "Feature A"
Expected: "Started context: Feature A"

# 2. Add some activity
$ my-context note "Working on feature A"
$ my-context touch

# 3. Start second context (should auto-stop first)
$ my-context start "Feature B"
Expected: "Stopped context: Feature A (duration: Xm)"
Expected: "Started context: Feature B"

# 4. Verify Feature B is active
$ my-context show
Expected: Shows "Feature B" as active

# 5. List all contexts
$ my-context list
Expected: Shows both contexts, Feature B marked active (●), Feature A stopped (○)

# 6. Check history
$ my-context history
Expected: Shows 2 transitions (start Feature A, switch to Feature B)
```

### Validation

- [ ] `Feature_A/meta.json` has end_time and status="stopped"
- [ ] `Feature_B/meta.json` has no end_time and status="active"
- [ ] `state.json` shows active_context="Feature B"
- [ ] `transitions.log` contains "switch" entry with both context names

---

## Scenario 3: Duplicate Name Handling

**Goal**: Validate automatic suffix appending for duplicate names.

### Steps

```bash
# 1. Create first context
$ my-context start "Bug fix"
Expected: "Started context: Bug fix"

# 2. Stop it
$ my-context stop

# 3. Try to create same name again
$ my-context start "Bug fix"
Expected: "Context \"Bug fix\" already exists."
Expected: "Started context: Bug fix_2"

# 4. Stop and try again
$ my-context stop
$ my-context start "Bug fix"
Expected: "Started context: Bug fix_3"

# 5. List contexts
$ my-context list
Expected: Shows "Bug fix", "Bug fix_2", "Bug fix_3"
```

### Validation

- [ ] Three directories: `Bug_fix/`, `Bug_fix_2/`, `Bug_fix_3/`
- [ ] Each has independent meta.json and log files
- [ ] Suffix sequence is correct (_2, then _3, not _1)

---

## Scenario 4: Cross-Shell Persistence

**Goal**: Validate state persists across shell sessions and types.

### Steps

**In git-bash**:
```bash
$ my-context start "Cross shell test"
$ my-context note "Added in git-bash"
```

**Switch to cmd.exe**:
```cmd
C:\> my-context show
Expected: Shows "Cross shell test" with git-bash note
C:\> my-context note "Added in cmd"
```

**Switch to WSL**:
```bash
$ my-context show
Expected: Shows "Cross shell test" with both notes
$ my-context stop
```

**Back to git-bash**:
```bash
$ my-context show
Expected: "No active context"
```

### Validation

- [ ] Same context visible in all shells
- [ ] Notes added in one shell visible in others
- [ ] Stop in one shell affects all shells

---

## Scenario 5: JSON Output Validation

**Goal**: Validate JSON output format is valid and parseable.

### Steps

```bash
# 1. Start with JSON output
$ my-context start "JSON test" --json > start.json
$ cat start.json
Expected: Valid JSON with command, timestamp, data.context_name

# 2. Show with JSON
$ my-context show --json > show.json
$ cat show.json
Expected: Valid JSON with full context details

# 3. Add note with JSON
$ my-context note "Test note" --json > note.json

# 4. List with JSON
$ my-context list --json > list.json

# 5. History with JSON
$ my-context history --json > history.json

# 6. Validate all JSON files
$ jq empty start.json show.json note.json list.json history.json
Expected: No errors (all valid JSON)

# 7. Extract specific fields
$ jq -r '.data.context.name' show.json
Expected: "JSON test"

$ jq '.data.contexts | length' list.json
Expected: Number of contexts
```

### Validation

- [ ] All JSON files are valid (parse without error)
- [ ] All contain required fields: command, timestamp, data
- [ ] Timestamps are ISO 8601 format
- [ ] jq can extract nested fields

---

## Scenario 6: Error Conditions

**Goal**: Validate proper error handling and messages.

### Steps

```bash
# 1. Try operations with no active context
$ my-context note "Should fail"
Expected: Error message, exit code 1

$ my-context file somefile.txt
Expected: Error message, exit code 1

$ my-context touch
Expected: Error message, exit code 1

# 2. Try to start with invalid name
$ my-context start ""
Expected: Error: "Context name required", exit code 1

$ my-context start "bug/fix"
Expected: Error: Invalid character '/', exit code 1

# 3. Stop when no context (should NOT error)
$ my-context stop
Expected: "No active context", exit code 0

# 4. Test with JSON error output
$ my-context note "Fail" --json
Expected: Valid JSON with error object, exit code 1
```

### Validation

- [ ] Error messages are clear and actionable
- [ ] Exit codes correct (0=success, 1=user error, 2=system error)
- [ ] Errors go to stderr (except JSON mode uses stdout)
- [ ] JSON error format matches spec

---

## Scenario 7: Path Normalization

**Goal**: Validate cross-platform path handling.

### Steps

**Windows paths**:
```bash
$ my-context start "Path test"
$ my-context file C:\Users\Dev\file.txt
$ my-context file src/relative/path.js
```

**Check storage**:
```bash
$ cat ~/.my-context/Path_test/files.log
Expected: Both paths in POSIX format (forward slashes)
Expected: Relative path converted to absolute
```

**Check display**:
```bash
$ my-context show
Expected: Native path format for display (backslashes on Windows)

$ my-context show --json
Expected: POSIX format in JSON output
```

### Validation

- [ ] All stored paths use forward slashes
- [ ] Relative paths converted to absolute
- [ ] Display format matches OS convention
- [ ] JSON always uses POSIX format

---

## Scenario 8: Single-Letter Aliases

**Goal**: Validate all command aliases work.

### Steps

```bash
$ my-context s "Alias test"    # start
$ my-context n "Note via alias"    # note
$ my-context f README.md    # file
$ my-context t    # touch
$ my-context w    # show
$ my-context l    # list
$ my-context h    # history
$ my-context p    # stop
```

### Validation

- [ ] All aliases produce same output as full commands
- [ ] Help text shows aliases: `my-context help start`

---

## Scenario 9: Help System

**Goal**: Validate help text is accessible and correct.

### Steps

```bash
$ my-context --help
Expected: Lists all commands with brief descriptions

$ my-context help
Expected: Same as --help

$ my-context start --help
Expected: Detailed help for start command including aliases

$ my-context -h
Expected: Same as --help
```

### Validation

- [ ] Help shows all 8 commands
- [ ] Each command's help shows aliases
- [ ] Usage examples provided
- [ ] Flag documentation included

---

## Scenario 10: Multiline and Special Characters

**Goal**: Validate proper escaping of special characters.

### Steps

```bash
# 1. Multiline note
$ my-context start "Special chars"
$ my-context note "Line 1
> Line 2
> Line 3"
Expected: Note added successfully

# 2. Note with pipe character
$ my-context note "This | has | pipes"

# 3. Check storage
$ cat ~/.my-context/Special_chars/notes.log
Expected: Newlines escaped as \n, pipes escaped as \|

# 4. Check display
$ my-context show
Expected: Proper display of multiline notes and pipes
```

### Validation

- [ ] Multiline notes stored correctly
- [ ] Pipe characters don't break parsing
- [ ] Notes display correctly with special chars

---

## Post-Testing Validation

After completing all scenarios:

### Directory Structure
```bash
$ ls -la ~/.my-context/
Expected:
  - state.json
  - transitions.log
  - Multiple context directories
```

### File Integrity
```bash
# All JSON files should be valid
$ find ~/.my-context -name "*.json" -exec jq empty {} \;
Expected: No errors

# All log files should be readable text
$ find ~/.my-context -name "*.log" -exec file {} \;
Expected: All identified as text files
```

### Performance Check
```bash
$ time my-context start "Performance test"
Expected: < 10ms

$ time my-context note "Quick note"
Expected: < 5ms

$ time my-context show
Expected: < 10ms
```

---

## Cleanup

```bash
# Remove test data
$ rm -rf ~/.my-context/

# Verify clean state
$ my-context list
Expected: "No contexts found"
```

---

## Success Criteria

All scenarios must:
- ✅ Execute without errors
- ✅ Produce expected output
- ✅ Create correct file structures
- ✅ Handle errors gracefully
- ✅ Work across cmd.exe, git-bash, and WSL
- ✅ Complete in < 10ms per command
- ✅ Produce valid JSON when --json flag used

**Ready for Implementation**: Once all quickstart tests pass, feature is complete.
