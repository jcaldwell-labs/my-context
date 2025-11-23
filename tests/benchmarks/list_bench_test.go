package benchmarks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// BenchmarkListWith1000Contexts measures list performance with 1000 contexts
// Constitution requirement: <1s for 1000 contexts
func BenchmarkListWith1000Contexts(b *testing.B) {
	// Setup: Create temporary home directory
	tmpHome := b.TempDir()
	os.Setenv("MY_CONTEXT_HOME", tmpHome)
	defer os.Unsetenv("MY_CONTEXT_HOME")

	// Create 1000 test contexts
	for i := 0; i < 1000; i++ {
		contextName := "TestContext_" + string(rune('0'+i%10)) + string(rune('0'+i/10%10)) + string(rune('0'+i/100%10)) + string(rune('0'+i/1000))
		contextDir := filepath.Join(tmpHome, contextName)
		// Write minimal meta.json for each context
		os.MkdirAll(contextDir, 0o755)
		os.WriteFile(filepath.Join(contextDir, "meta.json"), []byte(`{"name":"test","status":"stopped","is_archived":false}`), 0o644)
	}

	// Benchmark: List all contexts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.ListContextsFiltered(core.ContextFilter{
			Limit: 0, // No limit (show all)
		})
		if err != nil {
			b.Fatalf("ListContextsFiltered failed: %v", err)
		}
	}
	b.StopTimer()

	// Verify performance: Average should be <1s
	avgNs := b.Elapsed().Nanoseconds() / int64(b.N)
	if avgNs > 1_000_000_000 { // 1 second in nanoseconds
		b.Errorf("Performance target FAILED: Average time %v ns (%.3f s) exceeds 1s target",
			avgNs, float64(avgNs)/1_000_000_000.0)
	}
}

// BenchmarkListWithDefaultLimit measures list performance with default limit (10 contexts)
func BenchmarkListWithDefaultLimit(b *testing.B) {
	// Setup
	tmpHome := b.TempDir()
	os.Setenv("MY_CONTEXT_HOME", tmpHome)
	defer os.Unsetenv("MY_CONTEXT_HOME")

	// Create 100 contexts (more than default limit)
	for i := 0; i < 100; i++ {
		contextName := "TestContext_" + string(rune('0'+i))
		contextDir := filepath.Join(tmpHome, contextName)
		os.MkdirAll(contextDir, 0o755)
		os.WriteFile(filepath.Join(contextDir, "meta.json"), []byte(`{"name":"test","status":"stopped","is_archived":false}`), 0o644)
	}

	// Benchmark: List with default limit (10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.ListContextsFiltered(core.ContextFilter{
			Limit: 10,
		})
		if err != nil {
			b.Fatalf("ListContextsFiltered failed: %v", err)
		}
	}
}

// BenchmarkListWithProjectFilter measures list performance with project filtering
func BenchmarkListWithProjectFilter(b *testing.B) {
	// Setup
	tmpHome := b.TempDir()
	os.Setenv("MY_CONTEXT_HOME", tmpHome)
	defer os.Unsetenv("MY_CONTEXT_HOME")

	// Create 1000 contexts with mixed project names
	for i := 0; i < 1000; i++ {
		projectName := "project-" + string(rune('A'+i%10))
		contextName := projectName + ": Phase " + string(rune('0'+i%10))
		contextDir := filepath.Join(tmpHome, contextName)
		os.MkdirAll(contextDir, 0o755)
		os.WriteFile(filepath.Join(contextDir, "meta.json"), []byte(`{"name":"test","status":"stopped","is_archived":false}`), 0o644)
	}

	// Benchmark: List with project filter
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := core.ListContextsFiltered(core.ContextFilter{
			Project: "project-A",
			Limit:   0, // No limit
		})
		if err != nil {
			b.Fatalf("ListContextsFiltered failed: %v", err)
		}
	}
}
