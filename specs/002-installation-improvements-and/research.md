# Research: Installation & Usability Improvements

**Feature**: 002-installation-improvements-and  
**Date**: 2025-10-09  
**Phase**: Phase 0 (Outline & Research)

## Overview

This document captures research findings and technical decisions for Sprint 2 multi-platform distribution and CLI enhancements. All decisions respond to Sprint 1 retrospective feedback and user requests.

---

## Decision 1: Static Binary Build Strategy

**Context**: Users reported needing to build from source on WSL, creating installation friction.

**Decision**: Use Go's native cross-compilation with CGO_ENABLED=0 for static binaries

**Rationale**:
- Go 1.5+ includes native cross-compilation for all major platforms
- CGO_ENABLED=0 produces statically-linked binaries with zero runtime dependencies
- Static binaries eliminate common installation issues (missing libc versions, dynamic linking errors)
- Binary sizes remain small (<10MB) due to Go's efficient compilation

**Alternatives Considered**:
- Docker-based builds: Rejected - adds complexity and requires Docker installation
- CGo with cross-compilation: Rejected - requires platform-specific toolchains
- Separate compilation per platform: Rejected - manual process, error-prone

**Implementation**:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o my-context-linux-amd64
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o my-context-windows-amd64.exe
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o my-context-darwin-amd64
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o my-context-darwin-arm64
```

**References**:
- Go cross-compilation: https://go.dev/doc/install/source#environment
- Static linking benefits: Common practice in Go CLI tools (kubectl, docker CLI, gh)

---

## Decision 2: User-Level Installation (No Sudo)

**Context**: Requiring sudo/admin for installation creates enterprise environment blockers.

**Decision**: Install binaries to user-specific directories without elevated privileges

**Rationale**:
- `~/.local/bin/` (Unix/macOS) and `%USERPROFILE%\bin\` (Windows) are standard user paths
- Users can modify their own PATH without admin rights
- Aligns with modern security practices (principle of least privilege)
- Works in restricted corporate environments

**Alternatives Considered**:
- System-wide install (`/usr/local/bin`, `C:\Program Files`): Rejected - requires sudo/admin
- Homebrew/Chocolatey exclusive: Rejected - not universally available, especially on WSL
- Manual PATH instructions only: Rejected - poor UX, error-prone

**Implementation**:
- Unix: `~/.local/bin/` with automatic `.bashrc`/`.zshrc` PATH addition
- Windows: `%USERPROFILE%\bin\` with `setx PATH` (cmd) or `[Environment]::SetEnvironmentVariable` (PowerShell)
- Detect existing installations and offer upgrade path

---

## Decision 3: Project Name Extraction Pattern

**Context**: Users naturally organize contexts by project (e.g., "ps-cli: Phase 1", "ps-cli: Phase 2")

**Decision**: Parse project name as text before first colon, case-insensitive matching

**Rationale**:
- Colon separator is clear visual boundary
- Users already using this pattern (observed in Sprint 1 usage logs)
- Simple string operation (`strings.SplitN(name, ":", 2)[0]`)
- Case-insensitive matching prevents "PS-CLI" vs "ps-cli" confusion

**Alternatives Considered**:
- Structured metadata fields: Rejected - requires explicit project setting, breaks existing contexts
- Regex patterns: Rejected - overkill for simple parsing, harder to debug
- Prefix-based (underscore, dash): Rejected - colon is more readable and already in use

**Implementation**:
```go
func ExtractProjectName(contextName string) string {
    parts := strings.SplitN(contextName, ":", 2)
    if len(parts) > 1 {
        return strings.TrimSpace(parts[0])
    }
    return contextName // No colon = full name is project
}
```

**Edge Cases**:
- Multiple colons: Use only first as delimiter
- No colon: Treat entire name as project
- Leading/trailing whitespace: Trim automatically

---

## Decision 4: Archive vs Delete Implementation

**Context**: Users need both "hide completed work" (archive) and "remove test contexts" (delete)

**Decision**: Archive as metadata flag (is_archived:true in meta.json), Delete as physical removal

**Rationale**:
- Archive preserves data for future reference while hiding from default views
- Delete provides true cleanup for mistaken/test contexts
- JSON metadata change is backward compatible (field is optional with omitempty tag)
- Physical deletion respects user's explicit intent (with confirmation prompt)

**Alternatives Considered**:
- Archive as directory move (.archive/ subdirectory): Rejected - complicates path resolution, breaks existing references
- Soft delete with deleted flag: Rejected - no use case for recovering "deleted" contexts
- Single "remove" command with --permanent flag: Rejected - easy to lose data accidentally

**Implementation**:
- Archive: Add `is_archived: true` to meta.json, filter in ListContexts
- Delete: Remove entire context directory, preserve transitions.log history
- Both require confirmation, both prevent operation on active context

---

## Decision 5: Export Format (Markdown)

**Context**: Users want to share context summaries with team members

**Decision**: Generate markdown files with hierarchical structure (headers, lists, timestamps)

**Rationale**:
- Markdown is human-readable in any text editor
- Renders nicely in GitHub, GitLab, Slack, Notion, etc.
- Portable format (no special viewers required)
- Easy to grep/search with standard tools

**Alternatives Considered**:
- JSON export: Rejected - not human-friendly for sharing
- HTML export: Rejected - requires browser, harder to diff/version
- Plain text: Rejected - loses structure, harder to scan

**Format**:
```markdown
# Context: [Name]

**Started**: [Timestamp]
**Ended**: [Timestamp or "Active"]
**Duration**: [Human-readable duration]

## Notes

- [Timestamp] Note content
- [Timestamp] Note content

## Files

- [Timestamp] /path/to/file

## Activity

- [Timestamp] Touch event
```

---

## Decision 6: List Default Limit (10 contexts)

**Context**: Users with 50+ contexts reported overwhelming list output

**Decision**: Show last 10 contexts by default, add --all flag for complete list

**Rationale**:
- 10 contexts fit on single screen without scrolling
- Most users work on 3-5 active contexts at a time
- --all flag provides escape hatch for full access
- --limit <n> allows custom thresholds

**Alternatives Considered**:
- Pagination: Rejected - adds complexity, Unix tools don't paginate by default
- Show all with scroll prompt: Rejected - breaks piping to grep/awk
- Configuration file setting: Rejected - violates zero-config principle

**Implementation**:
- Sort contexts by start_time descending (most recent first)
- Display message: "Showing 10 of 50 contexts. Use --all to see all."
- --limit flag overrides default

---

## Decision 7: Backward Compatibility Strategy

**Context**: Sprint 1 users have existing contexts in ~/.my-context/ that must continue working

**Decision**: Additive JSON fields with omitempty tags, no schema migrations required

**Rationale**:
- Adding `is_archived` field with `omitempty` means old contexts without the field are treated as not archived
- Go's JSON unmarshaling handles missing fields gracefully (zero value = false for bools)
- No need for schema version numbers or migration scripts
- Installation scripts explicitly preserve ~/.my-context/ during upgrades

**Alternatives Considered**:
- Schema versioning with migrations: Rejected - overkill for adding one optional field
- Separate storage for Sprint 2 data: Rejected - fragments user data, complicates queries
- Rebuild contexts on upgrade: Rejected - risky, could lose user data

**Compatibility Testing**:
- Load Sprint 1 meta.json files (without is_archived field)
- Verify all existing commands work unchanged
- Verify new features work on old contexts

---

## Decision 8: GitHub Actions for Release Builds

**Context**: Manual multi-platform builds are error-prone and time-consuming

**Decision**: Use GitHub Actions workflow for automated release builds on git tag push

**Rationale**:
- GitHub-hosted runners support all target platforms
- Automatic binary uploads to GitHub Releases
- SHA256 checksum generation for security verification
- Triggered by semantic version tags (v2.0.0, v2.1.0, etc.)

**Alternatives Considered**:
- GoReleaser: Considered - adds dependency, may be overkill for simple builds
- Manual builds with scripts: Current state - error-prone, no audit trail
- Travis CI / Circle CI: Rejected - GitHub Actions is free for public repos, integrated

**Workflow Triggers**:
- Push tags matching `v*.*.*` pattern
- Build all 4 platform binaries in parallel
- Generate and upload SHA256 checksums
- Create GitHub Release with binaries and changelog

---

## Decision 9: Installation Script Language Choices

**Context**: Different platforms have different default shells

**Decision**: 
- install.sh: POSIX-compliant shell script (bash/zsh/sh compatible)
- install.bat: Windows batch file for cmd.exe
- install.ps1: PowerShell script for modern Windows
- curl-install.sh: Bash script for one-liner installation

**Rationale**:
- POSIX shell works on Linux, macOS, WSL, Git Bash
- Batch files work on all Windows versions (legacy support)
- PowerShell is modern Windows standard but not ubiquitous
- Providing all three maximizes compatibility

**Alternatives Considered**:
- PowerShell-only for Windows: Rejected - not available on all corporate Windows installs
- Python installer: Rejected - requires Python installation, defeats purpose
- Go-based installer: Rejected - requires Go installation or pre-built binary (chicken-egg problem)

---

## Open Questions

None. All technical decisions finalized based on spec requirements and Sprint 1 retrospective findings.

---

## Performance Targets

Based on expected usage patterns:

| Operation | Target | Rationale |
|-----------|--------|-----------|
| List 1000 contexts | <1s | Linear scan of JSON files, Go's file I/O is fast |
| Export 500 notes | <1s | Read + format + write, mostly I/O bound |
| Search 1000 contexts | <1s | Grep-like substring matching in memory |
| Start context | <50ms | Create directory + write JSON |
| Archive context | <50ms | Update single JSON field |

**Testing Strategy**: Create synthetic datasets (1000 contexts, 500 notes each) for performance validation in T038.

---

## Security Considerations

1. **Checksum Verification**: SHA256 checksums for downloaded binaries protect against tampering
2. **User-Level Install**: No sudo required eliminates privilege escalation risks
3. **Confirmation Prompts**: Delete command requires explicit confirmation to prevent accidents
4. **Path Injection**: Installation scripts validate paths before modifying shell rc files
5. **Data Preservation**: Explicit guarantees that upgrades preserve ~/.my-context/ data

---

## References

- Sprint 1 Retrospective: `../../../SPRINT-01-RETROSPECTIVE.md`
- Go Cross-Compilation: https://go.dev/doc/install/source#environment
- Semantic Versioning: https://semver.org/
- GitHub Actions: https://docs.github.com/en/actions
- Windows PATH Management: https://docs.microsoft.com/en-us/windows/win32/procthread/environment-variables

---

**Research Complete**: All NEEDS CLARIFICATION items resolved. Ready for Phase 1 (Design & Contracts).
