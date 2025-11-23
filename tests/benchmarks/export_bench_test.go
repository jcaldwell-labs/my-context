package benchmarks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
)

// BenchmarkExportWith500Notes measures export performance with 500 notes
// Constitution requirement: <1s for 500 notes
func BenchmarkExportWith500Notes(b *testing.B) {
	// Setup: Create temporary home directory
	tmpHome := b.TempDir()
	os.Setenv("MY_CONTEXT_HOME", tmpHome)
	defer os.Unsetenv("MY_CONTEXT_HOME")

	// Create context with 500 notes
	contextName := "BenchmarkContext"
	contextDir := filepath.Join(tmpHome, contextName)
	os.MkdirAll(contextDir, 0755)

	// Write meta.json
	meta := `{"name":"` + contextName + `","active":false,"started":"2025-10-09T10:00:00-04:00","stopped":"2025-10-09T11:00:00-04:00"}`
	os.WriteFile(filepath.Join(contextDir, "meta.json"), []byte(meta), 0644)

	// Write 500 notes to notes.log
	var notesBuilder strings.Builder
	for i := 0; i < 500; i++ {
		notesBuilder.WriteString("2025-10-09T10:")
		notesBuilder.WriteString(string(rune('0' + (i/60)%6)))
		notesBuilder.WriteString(string(rune('0' + (i%60)/10)))
		notesBuilder.WriteString(":")
		notesBuilder.WriteString(string(rune('0' + (i%60)%10)))
		notesBuilder.WriteString("-04:00|Note number ")
		notesBuilder.WriteString(string(rune('0' + i/100)))
		notesBuilder.WriteString(string(rune('0' + (i%100)/10)))
		notesBuilder.WriteString(string(rune('0' + i%10)))
		notesBuilder.WriteString(" with some descriptive text about the work being done\n")
	}
	os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte(notesBuilder.String()), 0644)

	// Write some files to files.log
	filesContent := "2025-10-09T10:00:00-04:00|/path/to/file1.go\n2025-10-09T10:01:00-04:00|/path/to/file2.go\n"
	os.WriteFile(filepath.Join(contextDir, "files.log"), []byte(filesContent), 0644)

	// Write some touches to touch.log
	touchContent := "2025-10-09T10:00:00-04:00\n2025-10-09T10:15:00-04:00\n2025-10-09T10:30:00-04:00\n"
	os.WriteFile(filepath.Join(contextDir, "touch.log"), []byte(touchContent), 0644)

	// Create temporary output directory
	tmpOutput := b.TempDir()

	// Benchmark: Export context
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tmpOutput, "export_"+string(rune('0'+i%10))+".md")
		_, err := core.ExportContext(contextName, outputPath, false)
		if err != nil {
			b.Fatalf("ExportContext failed: %v", err)
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

// BenchmarkExportSmallContext measures export performance for typical use case
func BenchmarkExportSmallContext(b *testing.B) {
	// Setup
	tmpHome := b.TempDir()
	os.Setenv("MY_CONTEXT_HOME", tmpHome)
	defer os.Unsetenv("MY_CONTEXT_HOME")

	// Create context with 10 notes (typical)
	contextName := "SmallContext"
	contextDir := filepath.Join(tmpHome, contextName)
	os.MkdirAll(contextDir, 0755)

	meta := `{"name":"` + contextName + `","active":false,"started":"2025-10-09T10:00:00-04:00","stopped":"2025-10-09T10:10:00-04:00"}`
	os.WriteFile(filepath.Join(contextDir, "meta.json"), []byte(meta), 0644)

	notesContent := "2025-10-09T10:00:00-04:00|First note\n" +
		"2025-10-09T10:01:00-04:00|Second note\n" +
		"2025-10-09T10:02:00-04:00|Third note\n" +
		"2025-10-09T10:03:00-04:00|Fourth note\n" +
		"2025-10-09T10:04:00-04:00|Fifth note\n" +
		"2025-10-09T10:05:00-04:00|Sixth note\n" +
		"2025-10-09T10:06:00-04:00|Seventh note\n" +
		"2025-10-09T10:07:00-04:00|Eighth note\n" +
		"2025-10-09T10:08:00-04:00|Ninth note\n" +
		"2025-10-09T10:09:00-04:00|Tenth note\n"
	os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte(notesContent), 0644)

	tmpOutput := b.TempDir()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tmpOutput, "small_export_"+string(rune('0'+i%10))+".md")
		_, err := core.ExportContext(contextName, outputPath, false)
		if err != nil {
			b.Fatalf("ExportContext failed: %v", err)
		}
	}
}

// BenchmarkExportWithFiles measures export performance with many file associations
func BenchmarkExportWithFiles(b *testing.B) {
	// Setup
	tmpHome := b.TempDir()
	os.Setenv("MY_CONTEXT_HOME", tmpHome)
	defer os.Unsetenv("MY_CONTEXT_HOME")

	// Create context with 100 notes and 50 file associations
	contextName := "FileHeavyContext"
	contextDir := filepath.Join(tmpHome, contextName)
	os.MkdirAll(contextDir, 0755)

	meta := `{"name":"` + contextName + `","active":false,"started":"2025-10-09T10:00:00-04:00","stopped":"2025-10-09T11:00:00-04:00"}`
	os.WriteFile(filepath.Join(contextDir, "meta.json"), []byte(meta), 0644)

	// 100 notes
	var notesBuilder strings.Builder
	for i := 0; i < 100; i++ {
		notesBuilder.WriteString("2025-10-09T10:00:00-04:00|Note ")
		notesBuilder.WriteString(string(rune('0' + i/10)))
		notesBuilder.WriteString(string(rune('0' + i%10)))
		notesBuilder.WriteString("\n")
	}
	os.WriteFile(filepath.Join(contextDir, "notes.log"), []byte(notesBuilder.String()), 0644)

	// 50 file associations
	var filesBuilder strings.Builder
	for i := 0; i < 50; i++ {
		filesBuilder.WriteString("2025-10-09T10:00:00-04:00|/path/to/project/src/file")
		filesBuilder.WriteString(string(rune('0' + i/10)))
		filesBuilder.WriteString(string(rune('0' + i%10)))
		filesBuilder.WriteString(".go\n")
	}
	os.WriteFile(filepath.Join(contextDir, "files.log"), []byte(filesBuilder.String()), 0644)

	tmpOutput := b.TempDir()

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tmpOutput, "file_heavy_export_"+string(rune('0'+i%10))+".md")
		_, err := core.ExportContext(contextName, outputPath, false)
		if err != nil {
			b.Fatalf("ExportContext failed: %v", err)
		}
	}
}
