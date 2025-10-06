# Research: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Date**: 2025-10-05  
**Status**: Complete

---

## 1. Multi-Platform Go Builds

### Decision: Use CGO_ENABLED=0 for Static Linking
**Rationale**: Static binaries eliminate runtime dependencies, ensuring the tool runs on clean systems without requiring libc or other libraries.

**Command Pattern**:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my-context-linux ./cmd/my-context/
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o my-context.exe ./cmd/my-context/
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o my-context-darwin-amd64 ./cmd/my-context/
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o my-context-darwin-arm64 ./cmd/my-context/
```

### Decision: GitHub Actions for Automated Builds
**Rationale**: GitHub Actions provides free CI/CD for open source, supports matrix builds for multiple platforms, and integrates directly with GitHub Releases.

**Alternatives Considered**:
- Makefile: Good for local builds, but requires Make installed on Windows
- Shell scripts only: Portable but no CI/CD integration
- **Selected**: GitHub Actions + shell scripts for local dev

**Implementation Notes**:
- Use `actions/checkout@v3` and `actions/setup-go@v4`
- Matrix strategy for platforms: `[linux, windows, darwin] x [amd64, arm64]` (exclude windows/arm64)
- Generate SHA256 checksums with `sha256sum` (Linux) or `shasum -a 256` (macOS)
- Upload artifacts to GitHub Releases with `actions/upload-artifact@v3`

### Decision: Skip macOS Signing for Sprint 2
**Rationale**: Code signing requires Apple Developer certificate ($99/year) and notarization adds complexity. Users can bypass Gatekeeper with "right-click > Open" workflow.

**Future Consideration**: Add signing in Sprint 3+ if adoption grows on macOS.

---

## 2. Installation Patterns for CLI Tools

### Decision: User-Specific Installation (No sudo Required)
**Rationale**: Avoids permission issues, works in corporate environments with restricted sudo, aligns with modern CLI tool practices (rustup, cargo, npm global).

**Installation Locations**:
- **Linux/macOS/WSL**: `~/.local/bin/my-context` (added to PATH in shell rc files)
- **Windows**: `%USERPROFILE%\bin\my-context.exe` (added to user PATH via registry)

**Alternatives Considered**:
- System-wide install (/usr/local/bin): Requires sudo, fails in restricted environments
- Current directory only: User must manually manage PATH
- **Selected**: User-specific with automatic PATH modification

### Decision: Backup-and-Replace Upgrade Strategy
**Rationale**: Prevents data loss if upgrade fails mid-operation. Old binary preserved as `.backup` until new one confirmed working.

**Upgrade Flow**:
1. Detect existing installation (`which my-context`)
2. Test new binary (`./my-context-new --version`)
3. Backup old binary (`mv my-context my-context.backup`)
4. Install new binary (`mv my-context-new my-context`)
5. Verify (`my-context --version`)
6. Remove backup if successful

### Decision: Shell-Specific PATH Modification
**Rationale**: Different shells use different rc files. Detection ensures PATH persists across sessions.

**Shell Detection Logic**:
```bash
if [ -n "$BASH_VERSION" ]; then
    RC_FILE="$HOME/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    RC_FILE="$HOME/.zshrc"
else
    RC_FILE="$HOME/.profile"  # Fallback
fi
```

**Windows PATH Modification**:
- PowerShell: `[Environment]::SetEnvironmentVariable("Path", $newPath, "User")`
- cmd.exe: `setx PATH "%PATH%;%USERPROFILE%\bin"`

---

## 3. Project Name Parsing Strategy

### Decision: Extract Text Before First Colon, Trim Whitespace
**Rationale**: Matches observed user pattern ("project: phase - description"). Simple, predictable, no ambiguity.

**Parsing Logic**:
```go
func ExtractProjectName(contextName string) string {
    parts := strings.SplitN(contextName, ":", 2)
    return strings.TrimSpace(parts[0])
}
```

**Edge Cases**:
- Multiple colons: `"project: sub: phase"` → project = "project"
- No colon: `"standalone"` → project = "standalone" (full name)
- Leading/trailing spaces: `" project : phase "` → project = "project" (trimmed)
- Unicode: `"проект: фаза"` → Supported (Go strings are UTF-8)

### Decision: Case-Insensitive Filtering
**Rationale**: User convenience. Typing `--project PS-CLI` should match `"ps-cli: Phase 1"`.

**Implementation**: `strings.EqualFold(projectName, extractedProject)`

**Alternatives Considered**:
- Case-sensitive exact match: Too strict, poor UX
- Fuzzy matching: Too complex, unexpected results
- **Selected**: Case-insensitive exact match

---

## 4. Markdown Export Format

### Decision: Human-Readable Structure with Headers and Lists
**Rationale**: Markdown should be readable without rendering. Headers provide clear sections, lists are scannable.

**Format Template**:
```markdown
# Context: {name}

**Started**: {ISO 8601 timestamp}
**Ended**: {ISO 8601 timestamp or "Active"}
**Duration**: {human-readable: "2h 15m" or "Active"}
**Status**: {active | stopped | archived}

## Notes ({count})

- `{HH:MM}` {note_text}
- `{HH:MM}` {note_text}

## Associated Files ({count})

- {file_path}
  Added: {timestamp}

## Activity

- {touch_count} touch events
- Last activity: {most recent touch timestamp}

---
*Exported: {export_timestamp}*
```

### Decision: Include All Metadata
**Rationale**: Export should be self-contained. Timestamps enable reconstruction, duration shows productivity, status indicates completion.

**Alternatives Considered**:
- Content-only (no timestamps): Loses context, can't verify when work happened
- Separate metadata file: More complex, harder to share
- **Selected**: Single markdown file with embedded metadata

### Decision: Overwrite Prompt (No Automatic Timestamp Suffix)
**Rationale**: Explicit user control prevents accidental data loss. Automatic timestamp suffixes create file clutter and make export paths unpredictable.

**Implementation**:
- If target file exists, prompt: `File exists: {path}. Overwrite? (y/N):`
- User accepts (y): Overwrite file
- User declines (n/Enter): Cancel export, exit code 2
- Future consideration: Add `--force` flag to skip prompt (Sprint 3+)

**Alternatives Considered**:
- Automatic timestamp suffix: Creates unpredictable filenames, harder to script
- Always overwrite: Too dangerous, no user control
- **Selected**: Explicit confirmation prompt
### Decision: Markdown Compatibility with GitHub/VS Code/Obsidian
**Rationale**: These are the three most popular markdown viewers/editors. Test format renders correctly in all three.

**Compatibility Notes**:
- Use standard markdown (no extensions)
- Inline code for timestamps (`` `HH:MM` ``) renders as monospace
- Tables avoided (not needed for this data structure)
- Bold for field names (`**Started**:`) is universally supported

---

## 5. Archive vs Delete Semantics

### Decision: Archive as Metadata Flag (Not Directory Move)
**Rationale**: Simpler implementation, less I/O, easier to unarchive in future Sprint. Moving directories risks breaking absolute paths in other systems.

**Implementation**:
- Add `"is_archived": true` to `meta.json`
- Default list view filters out archived contexts
- `--archived` flag shows only archived contexts
- Archive status does not affect data accessibility (files remain in place)

**Alternatives Considered**:
- Move to `~/.my-context/archived/` subdirectory: More visible separation, but complicates path handling
- Soft delete (is_deleted flag): Requires separate "purge" command, adds complexity
- **Selected**: Metadata flag for simplicity and flexibility

### Decision: Prevent Archiving Active Context
**Rationale**: Archiving implies completion. Active context is by definition incomplete.

**Error Message**: `"Cannot archive active context. Run 'my-context stop' first."`

### Decision: Delete Requires Confirmation (Unless --force)
**Rationale**: Deletion is irreversible. Confirmation prevents accidents. --force flag supports scripting.

**Confirmation Flow**:
```
$ my-context delete "Test Context"
⚠️  This will permanently delete context "Test Context" and all associated data.
Continue? (y/N): _
```

### Decision: Preserve Transitions Log on Delete
**Rationale**: Historical record should remain intact. Transitions log documents when context was created/stopped, which has audit value even after context deleted.

**Implementation**: Delete operation removes `~/.my-context/{context_dir}/` but does NOT modify `~/.my-context/transitions.log`. Deleted context name remains visible in history output.

---

## 6. List Pagination Approaches

### Decision: Limit-Based Pagination (Not Cursor-Based)
**Rationale**: CLI tools favor simplicity over scalability. Limit/offset is intuitive, cursor-based is overkill for <10,000 contexts.

**Default Behavior**: Show last 10 contexts (newest first)

**Implementation**:
```go
func ListContexts(limit int, showAll bool) []Context {
    contexts := loadAllContexts()
    sortByStartTime(contexts, descending=true)
    if !showAll && len(contexts) > limit {
        return contexts[:limit]
    }
    return contexts
}
```

### Decision: Newest First Sorting
**Rationale**: Users care most about recent work. Newest-first matches `git log` and other CLI tools.

**Alternatives Considered**:
- Oldest first: Logical for historical review, but poor UX for active work
- Alphabetical: Doesn't reflect temporal workflow
- **Selected**: Newest first (most recent = most relevant)

### Decision: Display Message When Results Truncated
**Rationale**: Users should know more contexts exist.

**Message Format**: `Showing 10 of 50 contexts. Use --all to see all.`

**Implementation**: Count total contexts before filtering, compare to displayed count.

---

## 7. Backward Compatibility Testing

### Decision: Graceful Degradation for Missing Fields
**Rationale**: Sprint 1 contexts lack `is_archived` field. Reading old meta.json should not error.

**Implementation Strategy**:
```go
type Context struct {
    Name      string    `json:"name"`
    StartTime time.Time `json:"start_time"`
    Status    string    `json:"status"`
    IsArchived bool     `json:"is_archived,omitempty"` // omitempty = optional
}

// On read, missing field defaults to false (Go zero value)
```

### Decision: No Migration Script Required
**Rationale**: Optional field with sensible default (false) means old data works immediately with new code.

**Validation**: Add integration test that loads Sprint 1 meta.json (without is_archived) and verifies it works.

### Decision: Semantic Versioning for Data Format
**Rationale**: If future Sprint requires breaking data changes, version field in meta.json enables migration detection.

**Future Consideration**: Add `"version": "1.0"` to meta.json in Sprint 3+. Sprint 2 focuses on additive changes only.

**Test Strategy**:
1. Create contexts with Sprint 1 binary
2. Upgrade to Sprint 2 binary
3. Verify all Sprint 1 contexts remain accessible
4. Add new Sprint 2 features (archive, export)
5. Verify Sprint 2 data works with Sprint 2 binary

---

## Summary of Key Decisions

| Area | Decision | Rationale |
|------|----------|-----------|
| **Builds** | Static binaries via CGO_ENABLED=0 | No runtime dependencies |
| **CI/CD** | GitHub Actions with matrix builds | Free, integrated, multi-platform |
| **Install Location** | ~/.local/bin (user-specific) | No sudo, works everywhere |
| **Upgrade** | Backup-and-replace | Safe, reversible |
| **Project Parsing** | Text before first colon, case-insensitive | Matches user pattern |
| **Export Format** | Markdown with metadata | Readable, portable, standard |
| **Archive** | Metadata flag in meta.json | Simple, reversible |
| **Delete** | Confirmation required (--force skips) | Prevent accidents |
| **Pagination** | Limit-based, newest first | Simple, matches user expectations |
| **Backward Compat** | Optional fields with zero-value defaults | No migration needed |

---

## Risk Assessment

**Low Risk**:
- Multi-platform builds (well-documented Go feature)
- Project name parsing (simple string operation)
- Markdown export (standard format)

**Medium Risk**:
- PATH modification across shells (many edge cases, needs thorough testing)
- Windows installation (cmd.exe vs PowerShell differences)
- Backward compatibility (needs comprehensive integration tests)

**Mitigation**:
- Comprehensive installation tests on real VMs (Windows 10, Ubuntu 22.04, macOS 13)
- Test matrix: {cmd, PowerShell, bash, zsh} x {fresh install, upgrade from Sprint 1}
- Document manual installation steps in TROUBLESHOOTING.md for edge cases

---

*Research complete. Ready for Phase 1 (Design & Contracts).*
