# Feature Specification: Context Home Visibility

**Feature ID**: MCF-001
**Priority**: P0 (Critical)
**Status**: Proposed
**Created**: 2025-10-22
**Category**: User Experience / Troubleshooting

---

## Problem Statement

Users operate on wrong context sets without realizing it because there is no visual indication of which `MY_CONTEXT_HOME` is currently active.

### Evidence

- User has 168 contexts in `~/.my-context/`, 24 in `~/.my-context-runtime/`
- `my-context show` returns different results depending on MY_CONTEXT_HOME
- Support questions: "Why don't I see my context?" (context exists, wrong home)
- Silent failures when expecting context in one home but operating on another

### Impact

- **Severity**: Critical
- **Frequency**: Every multi-context-home user, every session
- **User Frustration**: High (invisible problem, hard to debug)
- **Workaround**: None (users must manually `echo $MY_CONTEXT_HOME`)

---

## Goals

1. **Visibility**: Users always know which context home they're operating on
2. **Discoverability**: Easy to check context home without memorizing commands
3. **Consistency**: Show context home in all relevant commands
4. **Education**: Help users understand context home concept

### Non-Goals

- Automatic context home switching (out of scope)
- Context home migration tools (separate feature)
- Multi-home listing (complex, defer to future)

---

## User Stories

### Story 1: Checking Active Context Home
```
As a developer
When I run any my-context command
Then I should see which context home is active
So that I know which context set I'm operating on
```

### Story 2: Troubleshooting "Missing" Context
```
As a developer
When my context doesn't appear in `my-context list`
Then I can easily check if I'm using the wrong context home
So that I can switch to the correct home
```

### Story 3: Multi-Home Awareness
```
As a developer
When I work across multiple projects with different context homes
Then I can quickly verify I'm in the right home before operations
So that I don't create contexts in the wrong location
```

---

## Proposed Solution

### Solution A: Context Home Header (Recommended)

Add context home information to command outputs that display context data.

#### Commands Affected

| Command | Current Output | Proposed Output |
|---------|---------------|-----------------|
| `my-context show` | `Context: feat-auth` | `Context Home: ~/.my-context/`<br>`Context: feat-auth` |
| `my-context list` | `Contexts (152):` | `Context Home: ~/.my-context/ (152 contexts)`<br>`Contexts:` |
| `my-context start` | `Started: feat-auth` | `Context Home: ~/.my-context/`<br>`Started: feat-auth` |
| `my-context history` | `Context Transitions:` | `Context Home: ~/.my-context/`<br>`Context Transitions:` |

#### Example: `my-context show`

**Before**:
```
Context: feat-user-auth
Status: active
Started: 2025-10-22 09:00:00 (2h ago)

Notes (5):
  [2025-10-22T09:05:00] Implemented JWT validation
  ...
```

**After**:
```
Context Home: ~/.my-context/
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Context: feat-user-auth
Status: active
Started: 2025-10-22 09:00:00 (2h ago)

Notes (5):
  [2025-10-22T09:05:00] Implemented JWT validation
  ...
```

#### Example: `my-context list`

**Before**:
```
Contexts (152):

  ● first-pomo-doro-sess (active)
    Started: 2025-10-21T10:56:19-04:00 (19h ago)

  ○ bash-review (stopped)
    ...
```

**After**:
```
Context Home: ~/.my-context/ (152 contexts)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ● first-pomo-doro-sess (active)
    Started: 2025-10-21T10:56:19-04:00 (19h ago)

  ○ bash-review (stopped)
    ...
```

#### Example: `my-context start`

**Before**:
```
Started: feat-oauth-integration
```

**After**:
```
Context Home: ~/.my-context/

✓ Started: feat-oauth-integration
```

### Solution B: New `which` Command (Supplementary)

Add dedicated command to show context home location.

#### Command: `my-context which`

**Aliases**: `where`, `home`

**Output**:
```bash
$ my-context which

Context Home: ~/.my-context/

Environment:
  MY_CONTEXT_HOME=/home/be-dev-agent/.my-context/
  (Set in current shell)

Details:
  Location: /home/be-dev-agent/.my-context/
  Contexts: 152
  Active: first-pomo-doro-sess
  State file: /home/be-dev-agent/.my-context/state.json
```

**With `--short` flag**:
```bash
$ my-context which --short
~/.my-context/
```

**Use Cases**:
- Shell scripting: `CONTEXT_HOME=$(my-context which --short)`
- Quick verification: `my-context which`
- Debugging: `my-context which` shows full details

### Solution C: Environment Variable Helper (Optional)

Add shell function to show context home in prompt.

#### Bash/Zsh Function

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Show context home in prompt (optional)
function my_context_ps1() {
  if [ -n "$MY_CONTEXT_HOME" ]; then
    echo "[cx:$(basename $MY_CONTEXT_HOME)]"
  fi
}

# Add to PS1
PS1='$(my_context_ps1) \u@\h:\w\$ '
```

**Result**:
```
[cx:.my-context-project] user@host:~/projects/my-proj$
[cx:.my-context-runtime] user@host:~/projects$
user@host:~$  # No MY_CONTEXT_HOME set
```

**Documentation Location**: References section of my-context-workflow skill

---

## Implementation Details

### Code Changes

#### 1. Helper Function: `GetContextHomePath()`

**Location**: `internal/core/context.go` (or new `internal/core/env.go`)

```go
// GetContextHomePath returns the active context home path
// Returns default if MY_CONTEXT_HOME not set
func GetContextHomePath() string {
	home := os.Getenv("MY_CONTEXT_HOME")
	if home == "" {
		home = getDefaultContextHome()
	}
	return home
}

// GetContextHomeDisplay returns display-friendly path
// Replaces home directory with ~ for brevity
func GetContextHomeDisplay() string {
	path := GetContextHomePath()
	userHome, _ := os.UserHomeDir()
	if strings.HasPrefix(path, userHome) {
		return strings.Replace(path, userHome, "~", 1)
	}
	return path
}
```

#### 2. Output Function: `PrintContextHomeHeader()`

**Location**: `internal/output/human.go`

```go
// PrintContextHomeHeader prints context home with separator
func PrintContextHomeHeader() {
	home := core.GetContextHomeDisplay()
	contextCount := core.GetContextCount()

	fmt.Printf("Context Home: %s", home)
	if contextCount > 0 {
		fmt.Printf(" (%d contexts)", contextCount)
	}
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}
```

#### 3. Update Commands

**Commands to update**:
- `internal/commands/show.go` - Add header before context display
- `internal/commands/list.go` - Add header before context list
- `internal/commands/start.go` - Show context home in success message
- `internal/commands/history.go` - Add header before transitions

**Example: `show.go`**:
```go
func runShow(cmd *cobra.Command, args []string) error {
	// ... existing code ...

	// Add context home header
	if !jsonOutput {
		output.PrintContextHomeHeader()
	}

	// ... existing output logic ...
}
```

#### 4. New `which` Command

**Location**: `internal/commands/which.go`

```go
var WhichCmd = &cobra.Command{
	Use:     "which",
	Aliases: []string{"where", "home"},
	Short:   "Show context home location",
	Long:    "Display the active MY_CONTEXT_HOME location and context count.",
	Run:     runWhich,
}

func runWhich(cmd *cobra.Command, args []string) error {
	short, _ := cmd.Flags().GetBool("short")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	homePath := core.GetContextHomePath()
	homeDisplay := core.GetContextHomeDisplay()
	contextCount := core.GetContextCount()
	activeContext := core.GetActiveContextName()

	if jsonOutput {
		return outputJSON(homePath, contextCount, activeContext)
	}

	if short {
		fmt.Println(homeDisplay)
		return nil
	}

	// Full output
	fmt.Printf("Context Home: %s\n\n", homeDisplay)
	fmt.Printf("Environment:\n")
	fmt.Printf("  MY_CONTEXT_HOME=%s\n", os.Getenv("MY_CONTEXT_HOME"))
	if os.Getenv("MY_CONTEXT_HOME") == "" {
		fmt.Printf("  (Using default location)\n")
	} else {
		fmt.Printf("  (Set in current shell)\n")
	}
	fmt.Println()
	fmt.Printf("Details:\n")
	fmt.Printf("  Location: %s\n", homePath)
	fmt.Printf("  Contexts: %d\n", contextCount)
	fmt.Printf("  Active: %s\n", activeContext)
	fmt.Printf("  State file: %s/state.json\n", homePath)

	return nil
}

func init() {
	WhichCmd.Flags().BoolP("short", "s", false, "Short output (path only)")
}
```

### JSON Output

For `--json` flag compatibility:

**`my-context show --json`**:
```json
{
  "context_home": "/home/be-dev-agent/.my-context/",
  "context": {
    "name": "feat-user-auth",
    "status": "active",
    "started": "2025-10-22T09:00:00-04:00",
    ...
  }
}
```

**`my-context which --json`**:
```json
{
  "context_home": "/home/be-dev-agent/.my-context/",
  "context_home_display": "~/.my-context/",
  "context_count": 152,
  "active_context": "first-pomo-doro-sess",
  "env_set": true,
  "env_value": "/home/be-dev-agent/.my-context/"
}
```

---

## Testing Plan

### Unit Tests

#### Test: `GetContextHomeDisplay()`
```go
func TestGetContextHomeDisplay(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{"With MY_CONTEXT_HOME set to home", "$HOME/.my-context/", "~/.my-context/"},
		{"With MY_CONTEXT_HOME absolute", "/tmp/contexts", "/tmp/contexts"},
		{"With MY_CONTEXT_HOME relative", ".my-context-project", ".my-context-project"},
		{"Without MY_CONTEXT_HOME", "", "~/.my-context/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("MY_CONTEXT_HOME", tt.envValue)
			} else {
				os.Unsetenv("MY_CONTEXT_HOME")
			}

			result := GetContextHomeDisplay()
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

### Integration Tests

#### Test: Show Command with Header
```bash
#!/bin/bash
# Test: my-context show displays context home header

export MY_CONTEXT_HOME=~/.my-context-test
my-context start "test-context"

output=$(my-context show)

# Check for context home in output
if echo "$output" | grep -q "Context Home: ~/.my-context-test"; then
  echo "✓ Context home displayed"
else
  echo "✗ Context home NOT displayed"
  exit 1
fi
```

#### Test: Which Command
```bash
#!/bin/bash
# Test: my-context which shows correct home

export MY_CONTEXT_HOME=/tmp/test-contexts
result=$(my-context which --short)

if [ "$result" = "/tmp/test-contexts" ]; then
  echo "✓ Which command works"
else
  echo "✗ Which command failed"
  exit 1
fi
```

### Manual Testing Scenarios

1. **Cross-Home Verification**
   ```bash
   # Terminal 1
   export MY_CONTEXT_HOME=~/.my-context/
   my-context show  # Should show ~/.my-context/

   # Terminal 2
   export MY_CONTEXT_HOME=~/.my-context-runtime/
   my-context show  # Should show ~/.my-context-runtime/
   ```

2. **Default vs Custom**
   ```bash
   unset MY_CONTEXT_HOME
   my-context which  # Should show default

   export MY_CONTEXT_HOME=.my-context-project
   my-context which  # Should show .my-context-project
   ```

3. **List with Count**
   ```bash
   my-context list  # Should show "Context Home: ~/.my-context/ (152 contexts)"
   ```

---

## Documentation Updates

### 1. my-context-workflow Skill

**Add to Troubleshooting section**:
```markdown
### Check Which Context Home is Active

Verify context home location:
```bash
my-context which
# or
my-context which --short  # Path only
```

All context commands show context home in output header.
```

**Update Setup section**:
```markdown
### Environment Variable

My-context requires the `MY_CONTEXT_HOME` environment variable.

**Check active context home**:
```bash
my-context which
```

**Set context home**:
```bash
# Global contexts
export MY_CONTEXT_HOME=~/.my-context-runtime

# Project-local contexts
export MY_CONTEXT_HOME=.my-context-project
```

**Tip**: Add to shell profile (~/.bashrc) to persist across sessions.
```

### 2. Command Help Text

Update help for affected commands:

```bash
$ my-context show --help
Show context details

All output includes the active context home location for reference.
Use 'my-context which' to see only context home details.

Usage:
  my-context show [name] [flags]
...
```

### 3. README

Add to troubleshooting section:
```markdown
## Troubleshooting

### "I don't see my context"

Check which context home is active:

```bash
my-context which
```

You may have contexts in multiple homes:
- `~/.my-context/` (default)
- `~/.my-context-runtime/` (user-global for cross-project work)
- `.my-context-project/` (project-local)

Switch context home:
```bash
export MY_CONTEXT_HOME=~/.my-context-runtime/
my-context list
```
```

---

## Success Metrics

### Quantitative

1. **Reduced Support Questions**: "Why don't I see my context?" drops from X to <X/4
2. **User Surveys**: 90%+ can identify their context home after feature release
3. **Command Usage**: `my-context which` used 100+ times/month (shows need)

### Qualitative

1. **User Feedback**: "I always know which context home I'm using"
2. **Onboarding**: New users understand context homes faster
3. **Debugging**: Troubleshooting "missing context" becomes trivial

### Before/After Comparison

| Metric | Before | After (Target) |
|--------|--------|----------------|
| "Missing context" support tickets | 10/month | <3/month |
| Users aware of context home | ~20% | ~95% |
| Time to troubleshoot wrong home | 10+ min | <1 min |

---

## Rollout Plan

### Phase 1: Core Implementation (Week 1)
- Implement `GetContextHomePath()` helper
- Add headers to `show`, `list`, `start`, `history`
- Update unit tests

### Phase 2: `which` Command (Week 1)
- Implement `my-context which` command
- Add `--short` and `--json` flags
- Integration tests

### Phase 3: Documentation (Week 2)
- Update my-context-workflow skill
- Update README troubleshooting
- Update command help text

### Phase 4: Release (Week 2)
- Version bump (v2.3.0)
- Changelog entry
- Release notes highlighting feature
- Update skill in production

---

## Alternatives Considered

### Alternative 1: Implicit Detection Only
**Idea**: Don't show context home, just improve error messages

**Pros**: No UI changes, simpler implementation

**Cons**: Still opaque, users must debug via trial and error

**Rejected**: Doesn't solve visibility problem

### Alternative 2: Always Show Full Path
**Idea**: Show absolute path instead of `~/` abbreviation

**Pros**: Unambiguous, copy-pasteable

**Cons**: Verbose, cluttered output

**Rejected**: `~/` is standard shell convention, more readable

### Alternative 3: Colorize Context Home
**Idea**: Use colors to differentiate context homes

**Pros**: Visual distinction between homes

**Cons**: Doesn't work in all terminals, accessibility issues

**Rejected**: Text-based indication is universal

---

## Dependencies

### Code Dependencies
- None (uses existing context core logic)

### Feature Dependencies
- None (standalone feature)

### Testing Dependencies
- Existing test infrastructure

---

## Future Enhancements

### Enhancement 1: Multi-Home Listing
**Idea**: `my-context list --all-homes` to list contexts across all known homes

**Use Case**: User has contexts in multiple homes, wants unified view

**Complexity**: Medium (must discover all homes, merge/dedupe)

**Priority**: P2 (nice to have)

### Enhancement 2: Context Home Switching
**Idea**: `my-context use <home>` to switch active context home

**Use Case**: Quick switching without manual export

**Complexity**: Medium (must update shell env, may require shell integration)

**Priority**: P2 (convenience)

### Enhancement 3: Visual Indicator in Prompt
**Idea**: Official `my-context prompt` command for PS1 integration

**Use Case**: Always-visible context home in shell prompt

**Complexity**: Low (shell function wrapper)

**Priority**: P3 (user can already do this manually)

---

## Sign-off

### Approvals Required

- [ ] Product Owner - Confirm priority and user stories
- [ ] Tech Lead - Review implementation approach
- [ ] UX - Approve output format and messaging
- [ ] QA - Approve testing plan

### Estimated Effort

- **Development**: 2-3 days
- **Testing**: 1 day
- **Documentation**: 1 day
- **Total**: 4-5 days

### Target Release

- **Version**: v2.3.0
- **Sprint**: Next sprint (2-week cycle)
- **Dependencies**: None (can start immediately)

---

**Status**: Ready for Implementation
**Next Step**: Create implementation branch and begin Phase 1
