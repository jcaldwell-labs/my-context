# Feature Specification: Duplicate Detection & Resume Suggestion

**Feature ID**: MCF-002
**Priority**: P0 (Critical)
**Status**: Proposed
**Created**: 2025-10-22
**Category**: User Experience / Workflow Optimization

---

## Problem Statement

Users create duplicate contexts (`_2`, `_3`, `_4` suffixes) instead of resuming existing stopped contexts, fragmenting work history and making it difficult to find relevant past work.

### Evidence

**Observed Duplicates** (from real usage data):
```
002__Tech_Debt_Resolution
002__Tech_Debt_Resolution_2
002__Tech_Debt_Resolution_2_2
002__Tech_Debt_Resolution_2_2_2

006-sprint-006-large
006-sprint-006-large_2
006-sprint-006-large_3
006-sprint-006-large_4

ps-cli-impl-2025-10-16-payment-auth-v2
ps-cli-impl-2025-10-16-payment-auth-v2_2
```

**Impact Data**:
- 176+ total contexts observed
- 20+ contexts have numeric suffixes (`_2`, `_3`, `_4`)
- ~11% of contexts are likely unintentional duplicates
- Work history fragmented across multiple contexts with similar names

### Root Causes

1. **Silent Auto-Numbering**: Creating duplicate name appends `_2` without warning
2. **Resume Not Discoverable**: Users don't know resume command exists/when to use it
3. **No Search Prompt**: No reminder to search before creating new context
4. **Similar-Name Detection**: Tool doesn't detect semantically similar names

---

## Goals

1. **Prevent Duplicates**: Warn users before creating similar-named contexts
2. **Encourage Resume**: Suggest resuming instead of creating new when appropriate
3. **Edu

cate Workflow**: Teach search-first, resume-when-possible pattern
4. **Preserve Flexibility**: Still allow creating new contexts when intentional

### Non-Goals

- Automatic deduplication (too risky, user should decide)
- Fuzzy matching across all 176 contexts (performance concern)
- Context merging tools (separate feature)

---

## User Stories

### Story 1: Duplicate Detection
```
As a developer
When I try to create a context with a name similar to an existing context
Then I should be warned and shown the similar contexts
So that I can decide whether to resume or create new
```

### Story 2: Resume Suggestion
```
As a developer
When I start a context similar to a recently-stopped context
Then I should be prompted to resume the existing context
So that I continue work in the same context instead of fragmenting history
```

### Story 3: Search Encouragement
```
As a developer
When I create a new context
Then I should be reminded to search for existing contexts first
So that I build the habit of checking before creating
```

---

## Proposed Solution

### Solution A: Interactive Duplicate Detection (Recommended)

Detect similar context names at creation time and prompt user for action.

#### Matching Algorithm

**Match Types**:
1. **Exact Match** (case-insensitive) - `"feat-auth"` == `"feat-auth"`
2. **Suffix Match** - `"feat-auth"` matches existing `"feat-auth_2"`
3. **Fuzzy Match** (70%+ similarity) - `"Tech Debt Resolution"` ~ `"Tech-Debt-Res"`
4. **Keyword Match** - Significant keywords overlap (e.g., "payment", "auth")

**Priority**:
- Exact match: Always warn
- Recently stopped (< 7 days): Higher priority
- Same project prefix: Higher priority
- Suffix match: Medium priority
- Fuzzy match: Lower priority

#### User Experience Flow

**Scenario 1: Exact Match (Stopped Context)**
```bash
$ my-context start "feat-oauth-integration"

⚠️  Context "feat-oauth-integration" already exists (stopped 2h ago)

Options:
  1. Resume existing context
  2. Create new with suffix (feat-oauth-integration_2)
  3. Cancel

Choice [1/2/3]: █
```

**Scenario 2: Similar Name (Fuzzy Match)**
```bash
$ my-context start "tech-debt-work"

💡 Similar contexts found:

  1. Tech_Debt_Resolution_2_2 (stopped 3 days ago)
     Last note: "Fixed memory leak in parser"

  2. 002__Tech_Debt_Resolution (stopped 2 weeks ago)
     Last note: "Addressed code review feedback"

Options:
  r) Resume #1: Tech_Debt_Resolution_2_2
  n) Create new: tech-debt-work
  s) Search for more contexts
  c) Cancel

Choice [r/n/s/c]: █
```

**Scenario 3: No Match (Proceed Normally)**
```bash
$ my-context start "feat-new-feature"

✓ Started: feat-new-feature
```

#### Non-Interactive Mode

For scripting/automation, support flags to bypass prompts:

```bash
# Force create new (skip prompts)
my-context start "feat-auth" --force

# Always resume if exists
my-context start "feat-auth" --resume-if-exists

# Fail if duplicate (CI/CD use case)
my-context start "feat-auth" --fail-if-exists
```

### Solution B: Pre-Start Search Hint (Supplementary)

Before creating new context, show hint to search first.

**Implementation**:
```bash
$ my-context start

💡 Tip: Search for existing contexts before creating new:
   my-context list --search <keyword>
   my-context list --project <name>

Context name: █
```

Or with quick search:
```bash
$ my-context start

Search for existing contexts (optional): auth█

Found 3 contexts matching "auth":
  1. feat-oauth-integration (stopped 2h ago)
  2. impl-auth-service (stopped 1 day ago)
  3. debug-auth-failure (stopped 1 week ago)

Resume one of these? [1-3/n]:
```

### Solution C: Resume Command Enhancement

Make resume more discoverable and powerful.

#### `my-context resume` Improvements

**1. Interactive Resume (no args)**
```bash
$ my-context resume

Recently stopped contexts (resumable):

  1. feat-oauth-integration (stopped 2h ago)
     Last note: "Completed JWT implementation"

  2. code-review-PR-5598 (stopped 8h ago)
     Last note: "Requested changes on auth logic"

  3. impl-auth-service (stopped 1 day ago)
     Last note: "Started database schema design"

Select context to resume [1-3/c for cancel]: █
```

**2. Fuzzy Resume**
```bash
$ my-context resume auth

Multiple contexts match "auth":
  1. feat-oauth-integration (stopped 2h ago)
  2. impl-auth-service (stopped 1 day ago)
  3. debug-auth-failure (stopped 1 week ago)

Select [1-3/c]: █
```

**3. Smart `--last` Flag**
```bash
# Resume most recently stopped
$ my-context resume --last

Resumed: feat-oauth-integration (stopped 2h ago)
```

**4. `--recent` Flag**
```bash
# Show recently stopped contexts
$ my-context resume --recent

Stopped in last 24 hours:
  1. feat-oauth-integration (2h ago)
  2. code-review-PR-5598 (8h ago)

Stopped in last week:
  3. impl-auth-service (1 day ago)
  4. debug-auth-failure (3 days ago)

Select [1-4/c]: █
```

---

## Implementation Details

### Code Changes

#### 1. Similarity Detection Function

**Location**: `internal/core/similarity.go`

```go
package core

import (
	"strings"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// SimilarityResult represents a context match
type SimilarityResult struct {
	Context    *Context
	Score      float64  // 0.0 to 1.0
	MatchType  string   // "exact", "suffix", "fuzzy", "keyword"
	LastNote   string
	StoppedAgo string
}

// FindSimilarContexts finds contexts similar to given name
func FindSimilarContexts(name string, threshold float64) ([]SimilarityResult, error) {
	allContexts, err := ListAllContexts()
	if err != nil {
		return nil, err
	}

	var results []SimilarityResult

	for _, ctx := range allContexts {
		score, matchType := calculateSimilarity(name, ctx.Name)

		if score >= threshold {
			result := SimilarityResult{
				Context:    ctx,
				Score:      score,
				MatchType:  matchType,
				LastNote:   getLastNote(ctx),
				StoppedAgo: formatDuration(ctx.StoppedAt),
			}
			results = append(results, result)
		}
	}

	// Sort by score (highest first), then by recency
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Context.StoppedAt.After(results[j].Context.StoppedAt)
		}
		return results[i].Score > results[j].Score
	})

	return results, nil
}

// calculateSimilarity returns similarity score and match type
func calculateSimilarity(name1, name2 string) (float64, string) {
	// Normalize names
	n1 := strings.ToLower(strings.TrimSpace(name1))
	n2 := strings.ToLower(strings.TrimSpace(name2))

	// Exact match
	if n1 == n2 {
		return 1.0, "exact"
	}

	// Suffix match (e.g., "feat-auth" vs "feat-auth_2")
	if strings.HasPrefix(n2, n1+"_") || strings.HasPrefix(n1, n2+"_") {
		return 0.95, "suffix"
	}

	// Levenshtein distance (fuzzy match)
	distance := levenshtein.DistanceForStrings([]rune(n1), []rune(n2), levenshtein.DefaultOptions)
	maxLen := max(len(n1), len(n2))
	fuzzyScore := 1.0 - (float64(distance) / float64(maxLen))

	if fuzzyScore >= 0.7 {
		return fuzzyScore, "fuzzy"
	}

	// Keyword match (significant words in common)
	keywordScore := calculateKeywordOverlap(n1, n2)
	if keywordScore >= 0.6 {
		return keywordScore, "keyword"
	}

	return 0.0, "none"
}

// calculateKeywordOverlap calculates percentage of significant keywords in common
func calculateKeywordOverlap(name1, name2 string) float64 {
	// Split on common delimiters
	words1 := extractKeywords(name1)
	words2 := extractKeywords(name2)

	// Count common keywords
	common := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 && len(w1) > 3 { // Only significant words (>3 chars)
				common++
				break
			}
		}
	}

	totalKeywords := max(len(words1), len(words2))
	if totalKeywords == 0 {
		return 0.0
	}

	return float64(common) / float64(totalKeywords)
}

// extractKeywords extracts significant words from context name
func extractKeywords(name string) []string {
	// Split on delimiters: -, _, :, space
	words := strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return r == '-' || r == '_' || r == ':' || r == ' '
	})

	// Filter out common/stop words and dates
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true,
	}

	var keywords []string
	for _, word := range words {
		// Skip stop words, short words, and dates
		if !stopWords[word] && len(word) > 2 && !isDate(word) {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// isDate checks if string looks like a date
func isDate(s string) bool {
	// Simple check for YYYY, MM, DD patterns
	if len(s) == 4 || len(s) == 2 {
		if _, err := strconv.Atoi(s); err == nil {
			return true
		}
	}
	return false
}
```

#### 2. Interactive Prompt Function

**Location**: `internal/commands/start.go`

```go
// promptForDuplicateResolution shows similar contexts and gets user choice
func promptForDuplicateResolution(name string, similar []core.SimilarityResult) (string, error) {
	if len(similar) == 0 {
		return "create", nil
	}

	// Check for exact match
	exactMatch := similar[0].MatchType == "exact"

	fmt.Println()
	if exactMatch {
		fmt.Printf("⚠️  Context \"%s\" already exists (%s)\n\n", similar[0].Context.Name, similar[0].StoppedAgo)
	} else {
		fmt.Println("💡 Similar contexts found:")
		fmt.Println()
	}

	// Show top 3 matches
	maxShow := min(len(similar), 3)
	for i := 0; i < maxShow; i++ {
		result := similar[i]
		fmt.Printf("  %d. %s (%s)\n", i+1, result.Context.Name, result.StoppedAgo)
		if result.LastNote != "" {
			fmt.Printf("     Last note: \"%s\"\n", truncate(result.LastNote, 60))
		}
		fmt.Println()
	}

	fmt.Println("Options:")
	if exactMatch {
		fmt.Println("  1) Resume existing context")
		fmt.Printf("  2) Create new with suffix (%s_2)\n", name)
		fmt.Println("  3) Cancel")
		fmt.Println()
	} else {
		fmt.Printf("  r) Resume #1: %s\n", similar[0].Context.Name)
		fmt.Printf("  n) Create new: %s\n", name)
		fmt.Println("  s) Search for more contexts")
		fmt.Println("  c) Cancel")
		fmt.Println()
	}

	// Get user choice
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choice: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "1", "r":
		// Resume first match
		return "resume:" + similar[0].Context.Name, nil
	case "2", "n":
		// Create new
		return "create", nil
	case "s":
		// Show search results
		return "search", nil
	case "3", "c", "":
		// Cancel
		return "cancel", nil
	default:
		// Check if numeric choice for resuming specific context
		if idx, err := strconv.Atoi(choice); err == nil && idx >= 1 && idx <= maxShow {
			return "resume:" + similar[idx-1].Context.Name, nil
		}
		fmt.Println("Invalid choice. Cancelling.")
		return "cancel", nil
	}
}
```

#### 3. Update `start` Command

**Location**: `internal/commands/start.go`

```go
func runStart(cmd *cobra.Command, args []string) error {
	// ... existing validation ...

	name := args[0]

	// Check flags
	force, _ := cmd.Flags().GetBool("force")
	resumeIfExists, _ := cmd.Flags().GetBool("resume-if-exists")
	failIfExists, _ := cmd.Flags().GetBool("fail-if-exists")

	// Skip similarity check if force flag
	if !force {
		// Find similar contexts
		similar, err := core.FindSimilarContexts(name, 0.7)
		if err != nil {
			return err
		}

		if len(similar) > 0 {
			// Handle flags
			if failIfExists {
				return fmt.Errorf("context similar to \"%s\" already exists: %s", name, similar[0].Context.Name)
			}

			if resumeIfExists && similar[0].MatchType == "exact" {
				// Resume exact match
				return core.ResumeContext(similar[0].Context.Name)
			}

			// Interactive prompt
			action, err := promptForDuplicateResolution(name, similar)
			if err != nil {
				return err
			}

			switch {
			case strings.HasPrefix(action, "resume:"):
				// Resume selected context
				contextName := strings.TrimPrefix(action, "resume:")
				return core.ResumeContext(contextName)

			case action == "search":
				// Show full search results
				showSearchResults(name)
				return nil

			case action == "cancel":
				fmt.Println("Cancelled.")
				return nil

			case action == "create":
				// Continue with creation
			}
		}
	}

	// ... existing create context logic ...
}

func init() {
	StartCmd.Flags().Bool("force", false, "Skip duplicate detection, create anyway")
	StartCmd.Flags().Bool("resume-if-exists", false, "Resume if exact match exists")
	StartCmd.Flags().Bool("fail-if-exists", false, "Fail if similar context exists (CI/CD)")
}
```

#### 4. Enhanced Resume Command

**Location**: `internal/commands/resume.go`

```go
func runResume(cmd *cobra.Command, args []string) error {
	last, _ := cmd.Flags().GetBool("last")
	recent, _ := cmd.Flags().GetBool("recent")

	// Interactive mode (no args)
	if len(args) == 0 && !last {
		return interactiveResume(recent)
	}

	// Resume last stopped
	if last {
		ctx, err := core.GetLastStoppedContext()
		if err != nil {
			return err
		}
		return core.ResumeContext(ctx.Name)
	}

	// Fuzzy match on provided name
	name := args[0]
	matches, err := core.FindSimilarContexts(name, 0.6)
	if err != nil {
		return err
	}

	// Filter to stopped contexts only
	stoppedMatches := filterStopped(matches)

	if len(stoppedMatches) == 0 {
		return fmt.Errorf("no stopped contexts match \"%s\"", name)
	}

	if len(stoppedMatches) == 1 {
		// Single match, resume directly
		return core.ResumeContext(stoppedMatches[0].Context.Name)
	}

	// Multiple matches, let user choose
	return selectAndResumeContext(stoppedMatches)
}

func interactiveResume(recent bool) error {
	var contexts []core.SimilarityResult

	if recent {
		// Show contexts stopped in last 7 days
		contexts, _ = core.GetRecentlyStoppedContexts(7 * 24 * time.Hour)
	} else {
		// Show last 10 stopped contexts
		contexts, _ = core.GetRecentlyStoppedContexts(0)
		if len(contexts) > 10 {
			contexts = contexts[:10]
		}
	}

	if len(contexts) == 0 {
		fmt.Println("No recently stopped contexts to resume.")
		return nil
	}

	fmt.Println("Recently stopped contexts (resumable):")
	fmt.Println()

	for i, result := range contexts {
		fmt.Printf("  %d. %s (stopped %s)\n", i+1, result.Context.Name, result.StoppedAgo)
		if result.LastNote != "" {
			fmt.Printf("     Last note: \"%s\"\n", truncate(result.LastNote, 60))
		}
		fmt.Println()
	}

	fmt.Print("Select context to resume [1-", len(contexts), "/c for cancel]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "c" || choice == "" {
		fmt.Println("Cancelled.")
		return nil
	}

	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(contexts) {
		return fmt.Errorf("invalid choice: %s", choice)
	}

	return core.ResumeContext(contexts[idx-1].Context.Name)
}

func init() {
	ResumeCmd.Flags().Bool("last", false, "Resume most recently stopped context")
	ResumeCmd.Flags().Bool("recent", false, "Show recently stopped contexts (last 7 days)")
}
```

---

## Testing Plan

### Unit Tests

#### Test: Similarity Detection
```go
func TestFindSimilarContexts(t *testing.T) {
	tests := []struct{
		name           string
		searchName     string
		existingNames  []string
		expectedMatches int
		expectedFirst   string
		expectedType    string
	}{
		{
			name: "Exact match",
			searchName: "feat-auth",
			existingNames: []string{"feat-auth", "feat-api", "debug-auth"},
			expectedMatches: 1,
			expectedFirst: "feat-auth",
			expectedType: "exact",
		},
		{
			name: "Suffix match",
			searchName: "feat-auth",
			existingNames: []string{"feat-auth_2", "feat-api"},
			expectedMatches: 1,
			expectedFirst: "feat-auth_2",
			expectedType: "suffix",
		},
		{
			name: "Fuzzy match",
			searchName: "tech-debt-work",
			existingNames: []string{"Tech_Debt_Resolution", "feat-api"},
			expectedMatches: 1,
			expectedFirst: "Tech_Debt_Resolution",
			expectedType: "fuzzy",
		},
		{
			name: "No match",
			searchName: "feat-completely-new",
			existingNames: []string{"feat-auth", "debug-api"},
			expectedMatches: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test contexts
			setupTestContexts(tt.existingNames)

			results, err := FindSimilarContexts(tt.searchName, 0.7)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMatches, len(results))

			if tt.expectedMatches > 0 {
				assert.Equal(t, tt.expectedFirst, results[0].Context.Name)
				assert.Equal(t, tt.expectedType, results[0].MatchType)
			}
		})
	}
}
```

### Integration Tests

#### Test: Interactive Duplicate Prompt
```bash
#!/bin/bash
# Test: Duplicate detection prompts for existing context

my-context start "test-feat"
my-context stop

# Try to create same name
echo "1" | my-context start "test-feat"

# Should resume, not create _2
active=$(my-context show --json | jq -r '.context.name')
if [ "$active" = "test-feat" ]; then
  echo "✓ Resumed existing context"
else
  echo "✗ Failed to resume, created: $active"
  exit 1
fi
```

#### Test: Fuzzy Matching
```bash
#!/bin/bash
# Test: Fuzzy matching suggests similar contexts

my-context start "Tech-Debt-Resolution"
my-context stop

# Try similar name
echo "1" | my-context start "tech-debt-work"

# Should have prompted to resume "Tech-Debt-Resolution"
# (Verify via mock or output parsing)
```

### Manual Testing Scenarios

1. **Exact Duplicate**
   ```bash
   my-context start "feat-auth"
   my-context stop
   my-context start "feat-auth"  # Should prompt to resume
   ```

2. **Similar Name**
   ```bash
   my-context start "impl-2025-10-22-oauth"
   my-context stop
   my-context start "oauth-implementation"  # Should suggest resume
   ```

3. **Force Create**
   ```bash
   my-context start "feat-auth"
   my-context stop
   my-context start "feat-auth" --force  # Should create feat-auth_2
   ```

4. **Interactive Resume**
   ```bash
   my-context start "feat-a" && my-context stop
   my-context start "feat-b" && my-context stop
   my-context resume  # Should show list, let user choose
   ```

---

## Documentation Updates

### 1. my-context-workflow Skill

**Add "Resume vs Start Decision Guide" section**:

```markdown
## Resume vs Start: Decision Guide

### When to Resume

Resume an existing stopped context when:
- Continuing work on the same feature/task
- After interruption (meeting, bug fix, break)
- Context stopped within last few days
- Work is incomplete

Example:
```bash
# Yesterday
my-context start "feat-oauth-integration"
# ... work ...
my-context stop

# Today (continue same work)
my-context resume "feat-oauth-integration"
# or
my-context resume --last
```

### When to Start New

Create a new context when:
- New feature/task unrelated to previous work
- Significant direction change
- Previous context represents completed work
- Starting new sprint/phase

Example:
```bash
# Completed sprint 4
my-context archive "sprint-4-*"

# Start sprint 5
my-context start "sprint-5-planning"
```

### Search First Workflow

Before creating new context, search for existing:

```bash
# Search by keyword
my-context list --search "oauth"

# Search by project
my-context list --project ps-cli

# If found, resume instead of creating new
my-context resume "feat-oauth-integration"
```

### Duplicate Detection

My-context will warn if you try to create a similar-named context:

```bash
$ my-context start "feat-auth"

⚠️  Context "feat-auth" already exists (stopped 2h ago)

Options:
  1. Resume existing context
  2. Create new with suffix (feat-auth_2)
  3. Cancel

Choice: 1
```

**Best Practice**: Choose option 1 (Resume) unless you have specific reason to create duplicate.
```

**Add to Troubleshooting section**:
```markdown
### Why do I have duplicate contexts (_2, _3, _4)?

You created a context with a name that already existed. My-context automatically appends `_2`, `_3`, etc. to avoid conflicts.

**To avoid duplicates**:
1. Search before creating:
   ```bash
   my-context list --search <keyword>
   ```

2. Resume existing contexts:
   ```bash
   my-context resume --last
   my-context resume <name>
   ```

3. Use duplicate detection (enabled by default in v2.3.0+)

**To consolidate duplicates**:
Export and merge manually, then archive old duplicates:
```bash
my-context export "Tech_Debt_Resolution"
my-context export "Tech_Debt_Resolution_2"
# ... merge notes manually ...
my-context archive "Tech_Debt_Resolution_2"
```
```

### 2. README

Add to Usage section:
```markdown
## Avoiding Duplicates

My-context detects similar context names and prompts you to resume instead:

```bash
$ my-context start "feat-auth"
⚠️  Context "feat-auth" already exists (stopped 2h ago)

Options:
  1. Resume existing context
  2. Create new with suffix (feat-auth_2)
  3. Cancel
```

**Tip**: Search before creating to find existing contexts:
```bash
my-context list --search auth
my-context resume --last
```
```

### 3. Command Help

Update `start` help:
```bash
$ my-context start --help
Start a new context

Checks for similar existing contexts and prompts to resume instead of
creating duplicates. Use --force to skip this check.

Usage:
  my-context start <name> [flags]

Flags:
  --force                Skip duplicate detection
  --resume-if-exists     Resume if exact match exists
  --fail-if-exists       Fail if similar context exists (CI/CD)
  --project string       Project prefix for context name
  --labels string        Comma-separated labels

Examples:
  my-context start "feat-oauth-integration"
  my-context start "auth-work" --project ps-cli
  my-context start "hotfix-bug-123" --force  # Skip duplicate check
```

Update `resume` help:
```bash
$ my-context resume --help
Resume a stopped context

Resume by name, pattern, or interactively select from recent contexts.

Usage:
  my-context resume [name|pattern] [flags]

Flags:
  --last      Resume most recently stopped context
  --recent    Show recently stopped contexts (last 7 days)

Examples:
  my-context resume "feat-oauth-integration"  # Exact name
  my-context resume oauth                      # Fuzzy match
  my-context resume --last                     # Most recent
  my-context resume                            # Interactive selection
```

---

## Success Metrics

### Quantitative

1. **Duplicate Reduction**: Contexts with `_2`, `_3` suffixes drop from 11% to <3%
2. **Resume Usage**: `resume` command usage increases from ~5% to ~40% of `start` usage
3. **User Surveys**: 80%+ users understand when to resume vs start new

### Qualitative

1. **User Feedback**: "I no longer accidentally create duplicate contexts"
2. **Onboarding**: New users discover resume workflow within first week
3. **Workflow**: Users adopt search-first, resume-when-possible pattern

### Before/After Comparison

| Metric | Before | After (Target) |
|--------|--------|----------------|
| Duplicate contexts (_2, _3) | 20+ (11%) | <10 (3%) |
| Resume command usage | ~5% of starts | ~40% of starts |
| Users finding duplicate detection helpful | N/A | 90%+ |
| Fragmented work history complaints | Common | Rare |

---

## Rollout Plan

### Phase 1: Core Implementation (Week 1)
- Implement similarity detection algorithm
- Add interactive prompts to `start` command
- Add `--force` flag for non-interactive use

### Phase 2: Resume Enhancements (Week 1)
- Implement interactive resume (no args)
- Add `--recent` flag
- Add fuzzy matching for resume

### Phase 3: Testing (Week 2)
- Unit tests for similarity algorithm
- Integration tests for prompts
- Manual testing of user flows

### Phase 4: Documentation (Week 2)
- Update my-context-workflow skill
- Update README and command help
- Create migration guide for existing users

### Phase 5: Release (Week 2-3)
- Version bump (v2.3.0)
- Changelog entry
- Release notes highlighting feature
- Update skill in production

---

## Alternatives Considered

### Alternative 1: Automatic Resume (No Prompt)
**Idea**: Always resume if exact match exists, no prompt

**Pros**: Simplest UX, no interruption

**Cons**: Takes away user control, may resume wrong context

**Rejected**: Prompting gives users control and education

### Alternative 2: Fuzzy Match Only on Flag
**Idea**: Only check for duplicates with `--check-duplicates` flag

**Pros**: No behavior change for existing users

**Cons**: Feature not discoverable, duplicates continue

**Rejected**: Feature should be default for maximum impact

### Alternative 3: Post-Creation Merge Tool
**Idea**: Let users create duplicates, provide tool to merge later

**Pros**: Flexible, works with existing duplicates

**Cons**: Reactive, not proactive; complex merge logic

**Deferred**: Implement as separate feature (MCF-007) after this

---

## Dependencies

### Code Dependencies
- Levenshtein distance library (already in project)
- Existing context listing logic

### Feature Dependencies
- None (standalone feature)

### Testing Dependencies
- Interactive testing framework for prompts

---

## Future Enhancements

### Enhancement 1: Machine Learning Similarity
**Idea**: Use ML to better detect semantically similar contexts

**Example**: "feat-auth" and "implement-authentication" recognized as similar

**Complexity**: High (requires training data, ML model)

**Priority**: P3 (nice to have, current algorithm sufficient)

### Enhancement 2: Duplicate Consolidation Tool
**Idea**: `my-context merge <ctx1> <ctx2>` to consolidate duplicates

**Use Case**: Clean up existing duplicates from before this feature

**Complexity**: Medium (note merging, timestamp handling)

**Priority**: P2 (helps with migration)

### Enhancement 3: Global Duplicate Detection
**Idea**: Check for duplicates across all context homes

**Use Case**: User has "feat-auth" in both `~/.my-context/` and `.my-context-project/`

**Complexity**: High (performance, cross-home operations)

**Priority**: P3 (rare use case)

---

## Sign-off

### Approvals Required

- [ ] Product Owner - Confirm UX flow and prompts
- [ ] Tech Lead - Review similarity algorithm approach
- [ ] UX - Approve prompt wording and options
- [ ] QA - Approve testing plan

### Estimated Effort

- **Development**: 4-5 days (similarity detection, prompts, resume enhancements)
- **Testing**: 2 days
- **Documentation**: 1 day
- **Total**: 7-8 days

### Target Release

- **Version**: v2.3.0
- **Sprint**: Next sprint (2-week cycle)
- **Dependencies**: None (can start immediately)

---

**Status**: Ready for Implementation
**Next Step**: Create implementation branch and begin Phase 1
