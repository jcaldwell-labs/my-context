package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/jefferycaldwell/my-context-copilot/internal/models"
	"github.com/spf13/cobra"
)

var (
	archivePattern         string
	archiveDryRun          bool
	archiveCompletedBefore string
	archiveAllStopped      bool
)

func NewArchiveCmd(jsonOutput *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "archive [context-name]",
		Aliases: []string{"a"},
		Short:   "Archive completed contexts",
		Long: `Archive contexts to hide them from default list views while preserving all data.

Single context mode:
  Archive a specific context by name.

Bulk mode (when using flags):
  Archive multiple contexts matching patterns, dates, or all stopped contexts.

Archived contexts are hidden from default 'list' output but can be viewed with 'list --archived'.
Contexts must be stopped before archiving. All notes, files, and activity are preserved.

Examples:
  # Single context
  my-context archive "ps-cli: Phase 1"

  # Bulk operations
  my-context archive --pattern "old-*"
  my-context archive --all-stopped --dry-run
  my-context archive --completed-before 2024-01-01
  my-context archive --pattern "temp-*" --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if bulk mode flags are used
			isBulkMode := archivePattern != "" || archiveDryRun || archiveCompletedBefore != "" || archiveAllStopped

			if isBulkMode {
				return runBulkArchive()
			}

			// Single context mode (original behavior)
			return runSingleArchive(args)
		},
	}

	// Add bulk operation flags
	cmd.Flags().StringVar(&archivePattern, "pattern", "", "Archive contexts matching glob pattern (e.g., 'old-*')")
	cmd.Flags().BoolVar(&archiveDryRun, "dry-run", false, "Show what would be archived without actually archiving")
	cmd.Flags().StringVar(&archiveCompletedBefore, "completed-before", "", "Archive contexts completed before date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&archiveAllStopped, "all-stopped", false, "Archive all stopped contexts")

	return cmd
}

// runSingleArchive handles single context archiving (original behavior)
func runSingleArchive(args []string) error {
	// Validate: need context name
	if len(args) == 0 {
		return fmt.Errorf("context name required")
	}

	contextName := args[0]

	// Check if context exists
	ctx, _, _, _, err := core.GetContext(contextName)
	if err != nil {
		return fmt.Errorf("context not found: %s", contextName)
	}

	// Check if already archived
	if ctx.IsArchived {
		return fmt.Errorf("context %q is already archived", contextName)
	}

	// Prevent archiving active context
	activeCtx, err := core.GetActiveContext()
	if err == nil && activeCtx != nil && activeCtx.ActiveContext != nil && *activeCtx.ActiveContext == ctx.Name {
		return fmt.Errorf("cannot archive active context %q - stop it first with 'my-context stop'", contextName)
	}

	// Archive the context
	if err := core.ArchiveContext(contextName); err != nil {
		return fmt.Errorf("failed to archive context: %w", err)
	}

	fmt.Printf("Archived context: %s\n", contextName)
	return nil
}

// runBulkArchive handles bulk archiving operations
func runBulkArchive() error {
	// Step 1: Collect contexts based on flags
	contexts, err := collectContextsForBulkArchive()
	if err != nil {
		return fmt.Errorf("failed to collect contexts: %w", err)
	}

	if len(contexts) == 0 {
		fmt.Println("No contexts found matching the specified criteria.")
		return nil
	}

	// Step 2: Apply safety limit
	safetyLimit := getEnvInt("MC_BULK_LIMIT", 100)
	if len(contexts) > safetyLimit {
		return fmt.Errorf("too many contexts (%d) - exceeds safety limit of %d. Use a more specific pattern or adjust MC_BULK_LIMIT", len(contexts), safetyLimit)
	}

	// Step 3: Show dry-run preview or confirmation
	if archiveDryRun {
		return showBulkDryRun(contexts)
	}

	// Step 4: Get user confirmation
	confirmed, err := promptBulkConfirmation(contexts)
	if err != nil {
		return fmt.Errorf("confirmation failed: %w", err)
	}

	if !confirmed {
		fmt.Println("Bulk archive cancelled.")
		return nil
	}

	// Step 5: Execute bulk archive
	return executeBulkArchive(contexts)
}

// collectContextsForBulkArchive gathers contexts based on the specified flags
func collectContextsForBulkArchive() ([]*models.Context, error) {
	var contexts []*models.Context

	// Start with all stopped contexts
	allContexts, err := core.ListContexts()
	if err != nil {
		return nil, err
	}

	// Filter to stopped contexts only
	for _, ctx := range allContexts {
		if ctx.Status == "stopped" && !ctx.IsArchived {
			contexts = append(contexts, ctx)
		}
	}

	// Apply pattern filter if specified
	if archivePattern != "" {
		contexts, err = filterContextsByPattern(contexts, archivePattern)
		if err != nil {
			return nil, err
		}
	}

	// Apply date filter if specified
	if archiveCompletedBefore != "" {
		beforeTime, err := ParseDateString(archiveCompletedBefore)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		contexts = filterContextsByStopDate(contexts, beforeTime)
	}

	return contexts, nil
}

// filterContextsByPattern filters contexts using glob-style pattern matching
func filterContextsByPattern(contexts []*models.Context, pattern string) ([]*models.Context, error) {
	// Reuse the pattern matching logic from resume command
	var matches []*models.Context
	patternParts := strings.Split(pattern, "*")

	for _, ctx := range contexts {
		if MatchesPattern(ctx.Name, patternParts) {
			matches = append(matches, ctx)
		}
	}

	return matches, nil
}

// filterContextsByStopDate filters contexts stopped before the given time
func filterContextsByStopDate(contexts []*models.Context, before time.Time) []*models.Context {
	var filtered []*models.Context
	for _, ctx := range contexts {
		if ctx.EndTime != nil && ctx.EndTime.Before(before) {
			filtered = append(filtered, ctx)
		}
	}
	return filtered
}

// ParseDateString parses a YYYY-MM-DD date string
func ParseDateString(dateStr string) (time.Time, error) {
	// Parse as YYYY-MM-DD and set to end of day
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	// Set to end of day (23:59:59)
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, parsed.Location()), nil
}

// showBulkDryRun displays what would be archived without actually doing it
func showBulkDryRun(contexts []*models.Context) error {
	fmt.Printf("DRY RUN: Would archive %d contexts:\n", len(contexts))
	for i, ctx := range contexts {
		fmt.Printf("  %d. %s", i+1, ctx.Name)
		if ctx.EndTime != nil {
			fmt.Printf(" (completed: %s)", ctx.EndTime.Format("2006-01-02"))
		}
		fmt.Println()
	}
	fmt.Printf("\nUse --dry-run=false to actually archive these contexts.\n")
	return nil
}

// promptBulkConfirmation asks user to confirm bulk archive operation
func promptBulkConfirmation(contexts []*models.Context) (bool, error) {
	fmt.Printf("Archive %d contexts? This action cannot be undone.\n", len(contexts))

	// Show first 10 contexts
	maxDisplay := 10
	if len(contexts) <= maxDisplay {
		for i, ctx := range contexts {
			fmt.Printf("  %d. %s\n", i+1, ctx.Name)
		}
	} else {
		for i := 0; i < maxDisplay; i++ {
			fmt.Printf("  %d. %s\n", i+1, contexts[i].Name)
		}
		fmt.Printf("  ... and %d more\n", len(contexts)-maxDisplay)
	}

	fmt.Print("Proceed with archiving? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// executeBulkArchive performs the actual bulk archiving
func executeBulkArchive(contexts []*models.Context) error {
	fmt.Printf("Archiving %d contexts...\n", len(contexts))

	successCount := 0
	failedCount := 0
	var errors []string

	for _, ctx := range contexts {
		err := core.ArchiveContext(ctx.Name)
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("%s: %v", ctx.Name, err))
			fmt.Printf("❌ Failed to archive: %s (%v)\n", ctx.Name, err)
		} else {
			successCount++
			fmt.Printf("✅ Archived: %s\n", ctx.Name)
		}
	}

	// Summary
	fmt.Printf("\nBulk archive complete: %d successful, %d failed\n", successCount, failedCount)

	if len(errors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, errMsg := range errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	return nil
}

// getEnvInt reads an integer environment variable with a default value
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// MatchesPattern checks if a context name matches a glob pattern (copied from resume.go)
func MatchesPattern(name string, patternParts []string) bool {
	if len(patternParts) == 0 {
		return false // Empty pattern parts means no pattern was split, so no match
	}

	// Handle the special case of just "*" which should match everything
	if len(patternParts) == 2 && patternParts[0] == "" && patternParts[1] == "" {
		return true
	}

	// Handle simple cases
	if len(patternParts) == 1 {
		// No wildcards
		return name == patternParts[0]
	}

	// Check prefix
	if !strings.HasPrefix(name, patternParts[0]) {
		return false
	}

	remainingName := name[len(patternParts[0]):]

	// Check suffix
	lastPart := patternParts[len(patternParts)-1]
	if !strings.HasSuffix(remainingName, lastPart) {
		return false
	}

	// For multiple wildcards, we do a simple substring check
	for i := 1; i < len(patternParts)-1; i++ {
		if !strings.Contains(remainingName, patternParts[i]) {
			return false
		}
	}

	return true
}
