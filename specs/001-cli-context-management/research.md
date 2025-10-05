# Research: CLI Context Management System

**Feature**: 001-cli-context-management  
**Date**: 2025-10-04

## Technology Stack Research

### Programming Language Selection

**Question**: Which language provides best cross-platform CLI support with zero runtime dependencies?

**Options Evaluated**:
1. **Go** ⭐ SELECTED
2. Python
3. Rust
4. Shell Script (bash)

**Decision Matrix**:

| Criteria | Go | Python | Rust | Bash |
|----------|----|----|------|------|
| Single binary | ✅ Yes | ❌ No (runtime) | ✅ Yes | ❌ No (interpreter) |
| Cross-platform | ✅ Excellent | ⚠️ Good | ✅ Excellent | ❌ Poor (cmd vs bash) |
| Startup time | ✅ <5ms | ❌ 50-100ms | ✅ <5ms | ✅ <5ms |
| File I/O stdlib | ✅ Excellent | ✅ Good | ✅ Good | ⚠️ Limited |
| JSON support | ✅ Built-in | ✅ Built-in | ⚠️ Needs serde | ❌ jq required |
| Learning curve | ✅ Gentle | ✅ Gentle | ❌ Steep | ✅ Easy |
| Testing | ✅ Built-in | ✅ pytest | ✅ Built-in | ⚠️ bats |

**Final Decision**: Go 1.21+

**Rationale**:
- Zero runtime dependencies (critical for user adoption)
- Native Windows/Linux/macOS support with single codebase
- Excellent standard library for file operations, JSON, time handling
- Fast compilation and execution meets <10ms performance goal
- Strong testing story with table-driven tests
- Cross-compilation built-in: `GOOS=windows go build` produces .exe

---

### CLI Framework Research

**Question**: How to structure 8 subcommands with single-letter aliases and help generation?

**Options Evaluated**:
1. **Cobra** (spf13/cobra) ⭐ SELECTED
2. urfave/cli
3. Standard library flag package

**Comparison**:

**Cobra**:
- ✅ Industry standard (kubectl, hugo, gh CLI all use it)
- ✅ Subcommand tree structure matches our needs exactly
- ✅ `Aliases: []string{"s"}` for single-letter shortcuts
- ✅ Auto-generated help and usage text
- ✅ Persistent flags (`--json`) inherited by all subcommands
- ✅ 39k+ GitHub stars, active maintenance
- ⚠️ Adds ~2MB to binary (acceptable for our use case)

**urfave/cli**:
- ⚠️ Less intuitive subcommand API
- ❌ No built-in alias support (would need custom wrapper)
- ✅ Smaller binary footprint
- ❌ Less idiomatic for complex CLI apps

**Standard library**:
- ❌ Manual subcommand routing
- ❌ Manual help text generation
- ❌ No built-in flag inheritance
- ✅ Zero dependencies

**Decision**: Cobra

**Implementation Notes**:
```go
rootCmd.PersistentFlags().BoolP("json", "j", false, "Output as JSON")
startCmd.Aliases = []string{"s"}
stopCmd.Aliases = []string{"p"}
// etc.
```

---

### Storage Architecture Research

**Question**: How to store context data for portability, performance, and Unix philosophy compliance?

**Options Evaluated**:
1. **Subdirectories + text files** ⭐ SELECTED
2. SQLite database
3. Single JSON file per context
4. Git branches (original idea)

**Analysis**:

**Subdirectories + Text Files**:
```
~/.my-context/
├── state.json                    # Active context pointer
├── transitions.log               # Central transition log
└── Bug_fix/                      # Context subdirectory
    ├── meta.json                 # Context metadata
    ├── notes.log                 # Append-only notes
    ├── files.log                 # Append-only file refs
    └── touch.log                 # Append-only timestamps
```

**Pros**:
- ✅ Constitutional compliance: plain text, greppable
- ✅ Append-only logs = O(1) write performance
- ✅ No corruption risk (atomic file writes)
- ✅ Easy backup: `cp -r ~/.my-context ~/backup`
- ✅ Versioning: `git init ~/.my-context`
- ✅ Direct editing: `vim ~/.my-context/Bug_fix/notes.log`

**Cons**:
- ⚠️ List operation requires directory scan (acceptable for <1000 contexts)

**SQLite**:
- ❌ Binary format violates portability principle
- ❌ Corruption risk with concurrent writes
- ❌ Requires CGo for cross-compilation (complicates builds)
- ✅ Fast queries

**Single JSON per context**:
- ❌ Poor append performance (read entire file, modify, write)
- ❌ Larger memory footprint for big contexts
- ✅ Simple structure

**Git branches**:
- ❌ Git dependency not needed for file management
- ❌ Complex branch management API
- ❌ Overkill for simple versioning needs
- ❌ Harder to implement cross-platform

**Decision**: Subdirectories + Text Files

**File Format Specifications**:

**meta.json** (atomic read/write):
```json
{
  "name": "Bug fix",
  "start_time": "2025-10-04T14:00:00Z",
  "end_time": "2025-10-04T16:30:00Z",
  "status": "stopped"
}
```

**notes.log** (newline-delimited, pipe-separated):
```
2025-10-04T14:05:00Z|Fixed authentication bug in login handler
2025-10-04T14:30:00Z|Updated tests to cover edge case
```

**files.log** (newline-delimited, pipe-separated, POSIX paths):
```
2025-10-04T14:05:00Z|src/auth/login.go
2025-10-04T14:30:00Z|tests/auth/login_test.go
```

**touch.log** (newline-delimited, timestamp only):
```
2025-10-04T14:15:00Z
2025-10-04T15:00:00Z
2025-10-04T15:45:00Z
```

**transitions.log** (central file, pipe-separated):
```
2025-10-04T14:00:00Z|NULL|Bug fix|start
2025-10-04T16:30:00Z|Bug fix|NULL|stop
2025-10-04T16:35:00Z|NULL|Code review|start
```

---

### Path Normalization Research

**Question**: How to handle Windows backslash vs POSIX forward slash paths?

**Strategy**: Store in POSIX, convert on I/O

**Go Standard Library Functions**:
- `filepath.Clean(path)`: Normalizes path (removes .., redundant separators)
- `filepath.ToSlash(path)`: Converts to forward slashes for storage
- `filepath.FromSlash(path)`: Converts to OS-native separators for display
- `filepath.Abs(path)`: Resolves relative to absolute path

**Implementation Pattern**:
```go
// Input normalization (user provides path)
func NormalizePath(userPath string) string {
    cleaned := filepath.Clean(userPath)
    absolute, _ := filepath.Abs(cleaned)
    return filepath.ToSlash(absolute)  // Store as POSIX
}

// Output conversion (display to user)
func DisplayPath(storedPath string) string {
    return filepath.FromSlash(storedPath)  // Convert to OS native
}
```

**Cross-Platform Test Cases**:
- Windows input: `C:\Users\Dev\file.txt` → Store: `C:/Users/Dev/file.txt`
- WSL input: `/mnt/c/Users/Dev/file.txt` → Store: `/mnt/c/Users/Dev/file.txt`
- Git-bash input: `/c/Users/Dev/file.txt` → Store: `/c/Users/Dev/file.txt`

**Benefits**:
- Consistent storage format enables context sharing across shells
- POSIX format is more universal (Linux, macOS, WSL native)
- Windows tools understand forward slashes (cmd, PowerShell)

---

### Duplicate Name Resolution Algorithm

**Question**: How to generate _2, _3 suffixes efficiently?

**Algorithm**:
```
function ResolveContextName(desiredName string) string:
    if !DirectoryExists(desiredName):
        return desiredName
    
    counter = 2
    while DirectoryExists(desiredName + "_" + counter):
        counter++
    
    return desiredName + "_" + counter
```

**Performance**: O(n) where n = number of duplicates
- Typical case: n < 5 (most people don't create 5+ "Bug fix" contexts)
- Worst case: n = 100 (still <1ms with directory stat calls)

**Edge Cases**:
- Name "Fix_2" already exists → Creates "Fix_2_2" (treats _2 as part of name)
- Empty name "" → Reject with error "Context name required"
- Name ">200 chars" → Truncate to 200 chars before suffix logic

---

### JSON Output Format Design

**Question**: What JSON structure supports both human and machine parsing?

**Design Principles**:
1. Include command name and timestamp for debugging
2. Wrap data in `data` field for extensibility
3. Use ISO 8601 timestamps (UTC)
4. Null values explicit (not omitted)

**Example Outputs**:

**my-context show --json**:
```json
{
  "command": "show",
  "timestamp": "2025-10-04T15:30:00Z",
  "data": {
    "context": {
      "name": "Bug fix",
      "start_time": "2025-10-04T14:00:00Z",
      "end_time": null,
      "status": "active",
      "duration_seconds": 5400,
      "notes": [
        {"timestamp": "2025-10-04T14:05:00Z", "text": "Fixed auth bug"}
      ],
      "files": [
        {"timestamp": "2025-10-04T14:05:00Z", "path": "src/auth/login.go"}
      ],
      "touches": [
        {"timestamp": "2025-10-04T14:15:00Z"}
      ]
    }
  }
}
```

**my-context list --json**:
```json
{
  "command": "list",
  "timestamp": "2025-10-04T15:30:00Z",
  "data": {
    "contexts": [
      {
        "name": "Bug fix",
        "start_time": "2025-10-04T14:00:00Z",
        "end_time": null,
        "status": "active",
        "duration_seconds": 5400,
        "note_count": 3,
        "file_count": 2,
        "touch_count": 5
      }
    ]
  }
}
```

**Error Output** (still valid JSON):
```json
{
  "command": "note",
  "timestamp": "2025-10-04T15:30:00Z",
  "error": {
    "code": 1,
    "message": "No active context. Start a context with 'my-context start <name>'"
  }
}
```

---

## Performance Considerations

**Command Response Time Goals**: <10ms

**Optimization Strategies**:
1. **Lazy loading**: Only read files needed for command
   - `show`: Read single context directory
   - `list`: Stat directories, read meta.json only
   - `history`: Read transitions.log only

2. **Append-only logs**: No file rewriting for notes/files/touch
   - O(1) write performance
   - Sequential I/O (fast on all storage)

3. **Small state file**: state.json is <1KB
   - Cached in memory after first read
   - Updated only on start/stop

4. **Binary size**: Keep under 5MB
   - Cobra adds ~2MB
   - Standard library only for core logic
   - No heavy dependencies

**Measured Benchmarks** (expected):
- `my-context start "name"`: 3-5ms
- `my-context note "text"`: 2-3ms (append only)
- `my-context show`: 5-8ms (read 4 files)
- `my-context list`: 10-20ms (scan ~100 contexts)

---

## Security & Error Handling

**Threat Model**: Single-user local tool (not multi-user or networked)

**Security Considerations**:
1. **File permissions**: Use 0700 for ~/.my-context directory (owner-only)
2. **Path traversal**: Validate context names don't contain `..` or `/`
3. **Input sanitization**: Reject control characters in notes/names

**Error Handling Strategy**:

**Exit Codes**:
- `0`: Success
- `1`: User error (bad arguments, no active context)
- `2`: System error (I/O failure, permissions)

**Error Output** (human-readable, stderr):
```
Error: No active context
Start a context with: my-context start <name>
```

**Error Output** (JSON, stdout when --json):
```json
{"error": {"code": 1, "message": "No active context"}}
```

---

## Cross-Platform Testing Strategy

**Test Matrix**:
- Windows (cmd.exe, PowerShell, git-bash)
- WSL (Ubuntu)
- Linux (Ubuntu 20.04+)
- macOS (11+)

**Test Cases**:
1. **Path handling**: Windows `C:\` vs POSIX `/mnt/c/`
2. **Shell wrappers**: Executable works in all shells
3. **State persistence**: Context survives shell restarts
4. **Concurrent access**: Two shells reading same state
5. **Unicode names**: Context names with emoji/Chinese characters

**CI/CD Approach**:
- GitHub Actions matrix: windows-latest, ubuntu-latest, macos-latest
- Integration tests run on all platforms
- Build artifacts for each OS

---

## Dependencies Final List

**Direct Dependencies** (go.mod):
```
github.com/spf13/cobra v1.8.0
github.com/spf13/viper v1.18.0  // For MY_CONTEXT_HOME env var
github.com/stretchr/testify v1.8.4  // For test assertions
```

**Rationale**:
- Cobra: Essential for CLI structure
- Viper: Handles env var with fallback to default
- Testify: Makes test assertions readable

**Zero Runtime Dependencies**: Go compiles to static binary with no external requirements.

---

## Summary

**Key Decisions**:
1. ✅ Go 1.21+ for implementation
2. ✅ Cobra for CLI framework
3. ✅ Subdirectories + text files for storage
4. ✅ POSIX paths internally, convert on I/O
5. ✅ Suffix sequence _2, _3 for duplicates
6. ✅ JSON output via --json flag
7. ✅ Exit codes 0/1/2 for status/user-error/system-error

**No Unresolved Questions**: All technical decisions finalized.

**Ready for Phase 1**: Data model and contract design can proceed.
