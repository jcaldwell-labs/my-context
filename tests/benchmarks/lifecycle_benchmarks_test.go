package benchmarks

import (
	"os"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// BenchmarkNoteWarningOverhead benchmarks the performance impact of note warnings
func BenchmarkNoteWarningOverhead(b *testing.B) {
	// Setup temporary directory for testing
	tempDir := b.TempDir()
	originalHome := os.Getenv("MY_CONTEXT_HOME")
	defer func() {
		os.Setenv("MY_CONTEXT_HOME", originalHome)
	}()
	os.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	if err := core.EnsureContextHome(); err != nil {
		b.Fatalf("Failed to create context home: %v", err)
	}

	// Start a context
	_, _, err := core.CreateContext("benchmark-context")
	if err != nil {
		b.Fatalf("Failed to create context: %v", err)
	}

	// Add many notes to reach warning thresholds
	for i := 0; i < 60; i++ {
		if _, err := core.AddNote("Benchmark note for warning overhead test"); err != nil {
			b.Fatalf("Failed to add note %d: %v", i+1, err)
		}
	}

	// Benchmark adding one more note (which should trigger warnings)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := core.AddNote("Benchmark note with warning"); err != nil {
			b.Fatalf("Failed to add benchmark note: %v", err)
		}
	}

	// Target: <100ms per note add operation
	b.ReportMetric(float64(b.Elapsed().Nanoseconds()/int64(b.N))/1000000, "ms/op")
}

// BenchmarkBulkArchiveOperation benchmarks archive operation performance
func BenchmarkBulkArchiveOperation(b *testing.B) {
	// Note: This benchmark measures the archive operation itself
	// In practice, bulk archiving 100 contexts should take <5s total
	// We simulate this by measuring individual operations

	// Setup temporary directory for testing
	tempDir := b.TempDir()
	originalHome := os.Getenv("MY_CONTEXT_HOME")
	defer func() {
		os.Setenv("MY_CONTEXT_HOME", originalHome)
	}()
	os.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	if err := core.EnsureContextHome(); err != nil {
		b.Fatalf("Failed to create context home: %v", err)
	}

	// Pre-create multiple contexts for realistic bulk operation simulation
	contextsCreated := 0
	const maxContexts = 10 // Test with fewer contexts for benchmark stability

	for i := 0; i < maxContexts && contextsCreated < maxContexts; i++ {
		ctxName := "bulk-archive-test-" + string(rune('a'+i))
		if _, _, err := core.CreateContext(ctxName); err == nil {
			core.StopContext()
			contextsCreated++
		}
	}

	if contextsCreated == 0 {
		b.Fatal("Failed to create any test contexts")
	}

	// Get list of contexts to archive
	allContexts, err := core.ListContexts()
	if err != nil {
		b.Fatalf("Failed to list contexts: %v", err)
	}

	stoppedContexts := make([]string, 0, contextsCreated)
	for _, ctx := range allContexts {
		if ctx.Status == "stopped" && !ctx.IsArchived {
			stoppedContexts = append(stoppedContexts, ctx.Name)
		}
	}

	// Benchmark archiving operations
	b.ResetTimer()
	operations := 0

	for i := 0; i < b.N && operations < len(stoppedContexts); i++ {
		ctxName := stoppedContexts[operations]
		if err := core.ArchiveContext(ctxName); err != nil {
			b.Fatalf("Failed to archive context %s: %v", ctxName, err)
		}
		operations++

		// Wrap around if we run out of contexts
		if operations >= len(stoppedContexts) {
			operations = 0
			// Recreate contexts for next cycle
			for _, name := range stoppedContexts {
				if _, _, err := core.CreateContext(name+"-reset"); err == nil {
					core.StopContext()
				}
			}
		}
	}

	// Report results
	perOpTime := float64(b.Elapsed().Nanoseconds()/int64(b.N)) / 1000000 // ms
	b.ReportMetric(perOpTime, "ms/op")

	// Target: <50ms per operation (allowing <5s for 100 contexts)
	if perOpTime > 50.0 {
		b.Logf("Archive operation took %.2fms per context - for 100 contexts this would be %.2fs",
			perOpTime, perOpTime*100/1000)
	}
}

// BenchmarkResumeLast benchmarks the resume --last operation
func BenchmarkResumeLast(b *testing.B) {
	// Setup temporary directory for testing
	tempDir := b.TempDir()
	originalHome := os.Getenv("MY_CONTEXT_HOME")
	defer func() {
		os.Setenv("MY_CONTEXT_HOME", originalHome)
	}()
	os.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	if err := core.EnsureContextHome(); err != nil {
		b.Fatalf("Failed to create context home: %v", err)
	}

	// Create and stop several contexts to have history
	for i := 0; i < 10; i++ {
		name := "resume-last-test-" + string(rune('a'+i))
		_, _, err := core.CreateContext(name)
		if err != nil {
			b.Fatalf("Failed to create context %d: %v", i, err)
		}
		core.StopContext()
		time.Sleep(1 * time.Millisecond) // Ensure different stop times
	}

	// Benchmark GetMostRecentStopped (used by resume --last)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.GetMostRecentStopped()
		if err != nil {
			b.Fatalf("Failed to get most recent stopped: %v", err)
		}
	}

	// Target: <500ms
	b.ReportMetric(float64(b.Elapsed().Nanoseconds()/int64(b.N))/1000000, "ms/op")
}

// BenchmarkFindRelatedContexts benchmarks related context discovery
func BenchmarkFindRelatedContexts(b *testing.B) {
	// Setup temporary directory for testing
	tempDir := b.TempDir()
	originalHome := os.Getenv("MY_CONTEXT_HOME")
	defer func() {
		os.Setenv("MY_CONTEXT_HOME", originalHome)
	}()
	os.Setenv("MY_CONTEXT_HOME", tempDir)

	// Ensure context home is created
	if err := core.EnsureContextHome(); err != nil {
		b.Fatalf("Failed to create context home: %v", err)
	}

	// Create contexts with various prefixes
	prefixes := []string{"project-a", "project-b", "feature-x", "bug-fix"}
	for _, prefix := range prefixes {
		for i := 0; i < 5; i++ {
			name := prefix + ": task-" + string(rune('a'+i))
			_, _, err := core.CreateContext(name)
			if err != nil {
				b.Fatalf("Failed to create context %s: %v", name, err)
			}
			core.StopContext()
		}
	}

	testContext := "project-a: test-context"
	_, _, err := core.CreateContext(testContext)
	if err != nil {
		b.Fatalf("Failed to create test context: %v", err)
	}
	core.StopContext()

	// Benchmark finding related contexts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.FindRelatedContexts(testContext)
		if err != nil {
			b.Fatalf("Failed to find related contexts: %v", err)
		}
	}

	b.ReportMetric(float64(b.Elapsed().Nanoseconds()/int64(b.N))/1000000, "ms/op")
}

// BenchmarkPatternMatching benchmarks glob-style pattern matching
func BenchmarkPatternMatching(b *testing.B) {
	patterns := []string{
		"project-*",
		"*-feature",
		"bug-fix-*",
		"test-*-context",
		"*",
	}

	// Note: contexts variable not used in this benchmark
	_ = []string{
		"project-alpha",
		"project-beta",
		"feature-auth",
		"bug-fix-123",
		"test-integration-context",
		"standalone-context",
	}

	// Benchmark pattern matching performance
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pattern := range patterns {
			_, err := core.FindContextsByPattern(pattern)
			if err != nil {
				b.Fatalf("Failed pattern matching for %q: %v", pattern, err)
			}
		}
	}

	b.ReportMetric(float64(b.Elapsed().Nanoseconds()/int64(b.N))/1000000, "ms/op")
}
