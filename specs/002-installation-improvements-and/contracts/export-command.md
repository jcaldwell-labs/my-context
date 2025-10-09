# Command Contract: export

**Alias**: `e`  
**Purpose**: Export context details to markdown format for sharing and archival

## Syntax

```bash
my-context export <context-name> [--to <path>] [--force] [--json]
my-context export --all [--to <directory>] [--force] [--json]
```

## Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| context-name | string | ✅ (unless --all) | Name of context to export |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --to | | string | `./<context-name>.md` | Output file path or directory (for --all) |
| --all | | boolean | false | Export all contexts to separate files |
| --force | -f | boolean | false | Overwrite existing files without confirmation |
| --json | -j | boolean | false | Output JSON format instead of markdown |

## Behavior

### Single Context Export

**Input**: `my-context export "ps-cli: Phase 1"`

**Process**:
1. Validate context exists in ~/.my-context/
2. Check if output file exists
3. If exists and not --force, prompt: "File exists. Overwrite? (y/N)"
4. Read meta.json, notes.log, files.log, touches.log
5. Convert timestamps from UTC to local timezone
6. Format as markdown with headers, lists, timestamps
7. Write atomically to output file
8. Print success message with file path

**Output** (stdout):
```
Exported context "ps-cli: Phase 1" to ./ps-cli-Phase-1.md
```

**Exit Codes**:
- 0: Success
- 1: Error (context not found, read failure, write failure)
- 2: User cancellation (declined overwrite)

### All Contexts Export

**Input**: `my-context export --all --to ./exports/`

**Process**:
1. List all contexts in ~/.my-context/
2. Create output directory if doesn't exist
3. For each context:
   - Generate filename: `<sanitized-name>.md`
   - Export using single context logic
   - Track success/failure
4. Print summary

**Output** (stdout):
```
Exported 15 contexts to ./exports/
  ✓ ps-cli-Phase-1.md
  ✓ ps-cli-Phase-2.md
  ...
  ✗ invalid-name.md (Error: ...)
```

## Output Format

### Markdown Format (Default)

```markdown
# Context: ps-cli: Phase 1 - Foundation

**Started**: October 5, 2025 at 2:30 PM PDT
**Ended**: October 5, 2025 at 5:45 PM PDT
**Duration**: 3h 15m

**Exported**: October 9, 2025 at 7:30 PM PDT

---

## Notes

- **2:35 PM** Started research phase for context management improvements
- **2:50 PM** Decision: Use Cobra framework for CLI
- **3:15 PM** Implemented core context switching logic

Total: 45 notes

---

## Files

- **2:40 PM** /home/user/project/internal/core/context.go
- **3:20 PM** /home/user/project/cmd/my-context/main.go

Total: 12 files

---

## Activity

- **2:35 PM** Touch
- **3:10 PM** Touch
- **4:20 PM** Touch

Total: 8 touches

---

*Exported from my-context v2.0.0*
```

### JSON Format (--json)

```json
{
  "context": {
    "name": "ps-cli: Phase 1 - Foundation",
    "start_time": "2025-10-05T14:30:00Z",
    "end_time": "2025-10-05T17:45:00Z",
    "duration_seconds": 11700,
    "is_archived": false
  },
  "notes": [
    {"timestamp": "2025-10-05T14:35:00Z", "content": "Started research phase..."},
    {"timestamp": "2025-10-05T14:50:00Z", "content": "Decision: Use Cobra framework..."}
  ],
  "files": [
    {"timestamp": "2025-10-05T14:40:00Z", "path": "/home/user/project/internal/core/context.go"}
  ],
  "touches": [
    {"timestamp": "2025-10-05T14:35:00Z"}
  ],
  "export_metadata": {
    "exported_at": "2025-10-09T19:30:00Z",
    "version": "v2.0.0"
  }
}
```

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Context not found | Error: "Context 'name' not found" (exit 1) |
| Output file exists, no --force | Prompt for confirmation, exit 2 if declined |
| Output directory doesn't exist | Create parent directories automatically |
| Context has no notes | Show "(none)" in Notes section |
| Active context (no end_time) | Show "Currently Active" and calculate duration from start to now |
| Invalid output path | Error: "Cannot write to path: <reason>" (exit 1) |
| Empty context name | Error: "Context name required" (exit 1) |

## Validation

**Pre-Conditions**:
- Context must exist in ~/.my-context/
- If not --force, must handle existing file confirmation
- Output path must be writable

**Post-Conditions**:
- Output file exists with complete markdown/JSON content
- File is valid markdown (parseable) or valid JSON
- Original context data unchanged
- Success message printed to stdout

## Examples

```bash
# Export single context to default location
my-context export "ps-cli: Phase 1"

# Export to specific file
my-context export "ps-cli: Phase 1" --to ~/reports/phase1.md

# Export all contexts to directory
my-context export --all --to ./context-exports/

# Overwrite without prompting
my-context export "Test Context" --force

# JSON output for scripting
my-context export "ps-cli: Phase 1" --json > context.json
```

## Implementation Notes

- Sanitize context names for filenames (replace `/`, `\`, `:` with `-`)
- Use local timezone for human-readable timestamps
- Write files atomically (write to temp, then rename)
- Handle Unicode in context names and notes
- Preserve newlines and special characters in note content

