# Contract: Export Command

**Command**: `export` (alias: `e`)  
**Purpose**: Generate markdown representation of context data for sharing and archival

---

## Signature

```bash
my-context export <context-name> [--to <path>] [--all] [--json]
my-context e <context-name> [--to <path>] [--all] [--json]
```

---

## Arguments

**Positional**:
- `<context-name>`: Name of context to export (optional if --all used)

**Flags**:
- `--to <path>`: Output file path (default: `./{context_name}.md`)
- `--all`: Export all contexts to separate files
- `--json`: Output JSON instead of markdown (optional, for programmatic use)

---

## Behavior

### Single Context Export

**Input**: `my-context export "ps-cli: Phase 1"`

**Actions**:
1. Verify context exists in `~/.my-context/`
2. Read `meta.json`, `notes.log`, `files.log`, `touch.log`
3. Generate markdown according to data-model.md ExportDocument format
4. Write to `./ps-cli_Phase_1.md` (sanitize filename)
5. Print success message with file path

**Output**:
```
Exported context "ps-cli: Phase 1" to ./ps-cli_Phase_1.md
```

**Exit Code**: 0 (success)

### Export with Custom Path

**Input**: `my-context export "Phase 1" --to contexts/phase-1-summary.md`

**Actions**:
1. Create parent directories if needed (`mkdir -p contexts/`)
2. Check if file exists at target path
3. If exists: Prompt for overwrite confirmation (y/N)
4. Write markdown to specified path
5. Print success message

**Output**:
```
Exported context "Phase 1" to contexts/phase-1-summary.md
```

**Exit Code**: 0 (success), 2 (user cancelled overwrite)

### Export All Contexts

**Input**: `my-context export --all --to reports/`

**Actions**:
1. List all contexts (including archived)
2. Create output directory if needed
3. For each context:
   - Generate sanitized filename
   - Export to `reports/{context_name}.md`
   - Print progress
4. Print summary

**Output**:
```
Exporting 15 contexts to reports/...
  ✓ ps-cli_Phase_1.md
  ✓ ps-cli_Phase_2.md
  ...
Exported 15 contexts to reports/
```

**Exit Code**: 0 (success)

---

## Error Handling

### Context Not Found

**Input**: `my-context export "NonExistent"`

**Output**:
```
Error: Context "NonExistent" not found
Run 'my-context list --all' to see available contexts.
```

**Exit Code**: 1

### File Write Error

**Input**: `my-context export "Test" --to /root/protected.md` (no permission)

**Output**:
```
Error: Cannot write to /root/protected.md: permission denied
```

**Exit Code**: 2

### No Context Name with --all Missing

**Input**: `my-context export` (no args)

**Output**:
```
Error: Context name required or use --all to export all contexts
Usage: my-context export <context-name> [--to <path>] [--all]
```

**Exit Code**: 1

---

## File Format

### Markdown Output

See `data-model.md` for complete ExportDocument format.

**Key Requirements**:
- Timestamps in local timezone with ISO 8601 format
- Duration in human-readable format ("2h 15m", "3d 4h")
- Notes with time prefixes (`` `HH:MM` ``)
- File paths with "Added:" timestamps
- Footer with export timestamp and tool version

### JSON Output (with --json flag)

```json
{
  "command": "export",
  "timestamp": "2025-10-05T20:00:00Z",
  "data": {
    "context": {
      "name": "ps-cli: Phase 1",
      "start_time": "2025-10-04T10:00:00Z",
      "end_time": "2025-10-04T18:00:00Z",
      "status": "stopped",
      "is_archived": false
    },
    "notes": [
      {"timestamp": "2025-10-04T10:15:00Z", "text": "Initial planning"},
      ...
    ],
    "files": [
      {"timestamp": "2025-10-04T11:00:00Z", "path": "/path/to/file.go"},
      ...
    ],
    "activity": {
      "touch_count": 12,
      "last_touch": "2025-10-04T17:45:00Z"
    },
    "export_path": "./ps-cli_Phase_1.md"
  }
}
```

---

## Validation Rules

1. Context name must be non-empty string
2. Output path must be valid filesystem path
3. If --to omitted, use current directory (must be writable)
4. Sanitize context name for filename (replace special chars with underscores)
5. Confirm overwrite if target file exists (unless --force flag added in future)

---

## Integration Points

**Reads from**:
- `internal/core/context.go`: LoadContext()
- `internal/core/storage.go`: ReadLog()

**Writes to**:
- Filesystem: User-specified or default path

**Depends on**:
- Context existence check
- Markdown formatter (new in `internal/output/markdown.go`)

---

## Testing Checklist

- [ ] Export single context with default path
- [ ] Export single context with custom --to path
- [ ] Export all contexts to directory
- [ ] Export non-existent context (error handling)
- [ ] Export with special characters in context name (filename sanitization)
- [ ] Export to existing file (overwrite prompt)
- [ ] Export to read-only directory (permission error)
- [ ] Export with --json flag (JSON output format)
- [ ] Export context with no notes/files (empty sections)
- [ ] Export active vs stopped vs archived contexts (all work)
- [ ] Verify markdown renders correctly in GitHub/VS Code
- [ ] Verify timestamps in export are in local timezone

---

**Contract Version**: 1.0  
**Last Updated**: 2025-10-05
