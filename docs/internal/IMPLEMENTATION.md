# Implementation Guide

This guide walks through implementing all 42 tasks from `specs/001-cli-context-management/tasks.md`.

## Current Status

✅ **T001-T003 Complete**: Setup finished
- Go module initialized: `github.com/jefferycaldwell/my-context-copilot`
- Dependencies installed: cobra v1.10.1, viper v1.21.0, testify v1.11.1
- Directory structure created

✅ **Go Installation**: Fixed PATH issue - Go 1.25.1 now available

⏭️ **Next**: Phase 2 - Write tests first (TDD approach)

---

## Implementation Phases

### Phase 1: Setup (T001-T003) - READY TO RUN

Once Go is installed:

```bash
# T001: Initialize module
go mod init github.com/yourusername/my-context-copilot

# T002: Install dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/stretchr/testify@latest
go mod tidy

# T003: Already done! ✅
```

---

### Phase 2: Tests First - TDD (T004-T013)

These tests will initially **fail** - that's expected and correct!

**Priority Order:**
1. Start with T004 (start command test)
2. Then T005-T011 (other command tests)
3. Finally T012-T013 (cross-platform tests)

**Example - T004: Start Command Test**
Create `tests/integration/start_test.go`:
```go
package integration_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestStartCommand_CreatesNewContext(t *testing.T) {
    // Test will fail until implementation exists
    // This is TDD - write test first!
}
```

Run to verify it fails:
```bash
go test ./tests/integration/start_test.go
```

Repeat for T005-T013.

---

### Phase 3: Models (T014-T019) - Parallel

Once tests are written, create all 6 model files simultaneously:

**T014: Context Model**
```bash
# Create internal/models/context.go
# Define Context struct with JSON tags
```

**T015-T019**: Similar for Note, FileAssociation, TouchEvent, ContextTransition, AppState

Run tests after models:
```bash
go test ./... # Some import errors expected, models need core logic
```

---

### Phase 4: Core Logic (T020-T022) - Sequential

**MUST be done in order:**

1. **T020: Storage** (`internal/core/storage.go`)
   - CreateDir, ReadJSON, WriteJSON, AppendLog, ReadLog
   - NormalizePath functions

2. **T021: State** (`internal/core/state.go`) - depends on T020
   - GetActiveContext, SetActiveContext, ClearActiveContext

3. **T022: Context Core** (`internal/core/context.go`) - depends on T020+T021
   - CreateContext, StopContext, AddNote, AddFile, etc.

After T022, run tests:
```bash
go test ./... # More tests should pass now
```

---

### Phase 5: Output Formatters (T023-T024) - Parallel

Both can be done simultaneously:

**T023**: Human-readable output (`internal/output/human.go`)
**T024**: JSON output (`internal/output/json.go`)

---

### Phase 6: Commands (T025-T033) - Parallel

All 9 command files can be created in parallel:

**T025**: Root command - `cmd/my-context/main.go`
**T026-T033**: 8 subcommands in `internal/commands/*.go`

After commands, run full test suite:
```bash
go test ./... # Most/all tests should pass now!
```

---

### Phase 7: Integration (T034-T036)

**T034**: Add environment config to main.go
**T035**: Create error types in `internal/core/errors.go`
**T036**: Enhance help text in all commands

Build and test:
```bash
go build -o my-context.exe ./cmd/my-context/
./my-context.exe --help
./my-context.exe start "Test" # First real command!
```

---

### Phase 8: Build & Distribution (T037-T039)

**T037**: Create `scripts/build.sh`
```bash
#!/bin/bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/my-context.exe ./cmd/my-context/
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/my-context-linux ./cmd/my-context/
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/my-context-macos ./cmd/my-context/
```

**T038**: Create `scripts/install.sh`
**T039**: Create `.github/workflows/build.yml`

---

### Phase 9: Polish (T040-T042) - Parallel

**T040**: Add unit tests in `tests/unit/`
**T041**: Run performance benchmarks from `specs/001-cli-context-management/quickstart.md`
**T042**: Update README.md with usage examples

Final validation:
```bash
go test ./... -cover # Should have >80% coverage
./scripts/build.sh # Build all platforms
ls -lh dist/ # Check binary sizes (<5MB)
```

---

## Quick Reference

### Run specific test
```bash
go test ./tests/integration/start_test.go -v
```

### Run all tests with coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Build for current platform
```bash
go build -o my-context.exe ./cmd/my-context/
```

### Check for errors before running tests
```bash
go build ./...
```

---

## Task Checklist

Copy to track progress:

**Phase 1: Setup**
- [ ] T001: go mod init
- [ ] T002: Install dependencies
- [x] T003: Directory structure

**Phase 2: Tests (Write First!)**
- [ ] T004: start_test.go
- [ ] T005: stop_test.go
- [ ] T006: note_test.go
- [ ] T007: file_test.go
- [ ] T008: touch_test.go
- [ ] T009: show_test.go
- [ ] T010: list_test.go
- [ ] T011: history_test.go
- [ ] T012: paths_test.go
- [ ] T013: json_test.go

**Phase 3: Models**
- [ ] T014: context.go
- [ ] T015: note.go
- [ ] T016: file_association.go
- [ ] T017: touch_event.go
- [ ] T018: transition.go
- [ ] T019: state.go

**Phase 4: Core (Sequential!)**
- [ ] T020: storage.go
- [ ] T021: state.go (management)
- [ ] T022: context.go (operations)

**Phase 5: Output**
- [ ] T023: human.go
- [ ] T024: json.go

**Phase 6: Commands**
- [ ] T025: main.go (root)
- [ ] T026: start.go
- [ ] T027: stop.go
- [ ] T028: note.go
- [ ] T029: file.go
- [ ] T030: touch.go
- [ ] T031: show.go
- [ ] T032: list.go
- [ ] T033: history.go

**Phase 7: Integration**
- [ ] T034: Environment config
- [ ] T035: Error types
- [ ] T036: Help text

**Phase 8: Build**
- [ ] T037: build.sh
- [ ] T038: install.sh
- [ ] T039: CI/CD workflow

**Phase 9: Polish**
- [ ] T040: Unit tests
- [ ] T041: Performance validation
- [ ] T042: Update README

---

## Need Help?

- **Design docs**: See `specs/001-cli-context-management/`
- **Quick commands**: See `HERE.md`
- **Task details**: See `specs/001-cli-context-management/tasks.md`
- **Constitution**: See `.specify/memory/constitution.md`

---

**Next Step**: Install Go, then run T001-T002 commands from SETUP.md
