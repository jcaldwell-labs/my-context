# External Integrations

This guide documents known external tools and integrations that work with my-context through the PostgreSQL database backend.

## Overview

My-context supports external integrations through:

1. **Direct PostgreSQL Access** - Tools can read/write to the same database
2. **Shared Schema** - All tools use the same table structures
3. **Extended Metadata** - Additional fields for tool-specific data

**Important:** External integrations bypass my-context's validation and business logic. This is intentional for performance but requires careful coordination.

## TimeTracker Integration

**Location:** `/home/cdev/projects/c/timetracker` (private project)  
**Maintainer:** cdev  
**Bridge Script:** `scripts/tt-context-bridge`

### Architecture

TimeTracker is a C/ncurses TUI application for time tracking that integrates bidirectionally with my-context:

```
┌─────────────────────┐     ┌─────────────────────┐
│  TimeTracker TUI    │     │  my-context CLI     │
│  (ncurses C app)    │     │  (Go CLI)           │
└─────────┬───────────┘     └─────────┬───────────┘
          │                           │
          │ Unix Socket               │ Direct
          ↓                           ↓
┌─────────────────────┐     ┌─────────────────────┐
│  tt-context-bridge  │────→│  PostgreSQL         │
│  (bash script)      │←────│  (shared database)  │
└─────────────────────┘     └─────────────────────┘
```

### Bridge Commands

The `tt-context-bridge` script provides several synchronization modes:

| Command | Description | Use Case |
|---------|-------------|----------|
| `sync` | Import contexts from PostgreSQL (last 15 min) | Initial sync after starting timetracker |
| `export` | Push timetracker notes to PostgreSQL | Save timetracker session notes to my-context |
| `watch` | Bidirectional real-time sync (30s interval) | Continuous synchronization during work session |
| `report` | Generate daily/ticket/standup summaries | Create reports from combined data |

### Extended Note Metadata

TimeTracker adds additional metadata to notes through the `context_notes` table's `note_type` and `tags` columns. For complete schema documentation, see [Database Schema Documentation](../internal/SCHEMA.md).

#### Database Schema

The `context_notes` table supports extended metadata:

```sql
CREATE TABLE context_notes (
    id SERIAL PRIMARY KEY,
    context_id INTEGER NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    noted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    note_type VARCHAR(50),        -- Used by timetracker for metadata
    tags TEXT[]                   -- Array of metadata tags
);
```

**Key Fields:**
- `note_type` - Reserved for external tool metadata (my-context core does not use this field)
- `tags` - PostgreSQL array for categorization and filtering

#### Metadata Fields

TimeTracker stores metadata in the `note_type` field as JSON or structured text:

| Field | Description | Example |
|-------|-------------|---------|
| `source` | Origin identifier | `shell:ses-123:TICKET-456` |
| `session_id` | Terminal session ID | `ses-20260125-143022` |
| `backdated` | Note created with `--ago` flag | `true` / `false` |
| `no_context` | From "shadow session" (no active context) | `true` / `false` |

**Example note_type value:**
```json
{
  "source": "shell:ses-123:PROJ-789",
  "session_id": "ses-20260125-143022",
  "backdated": false,
  "no_context": false
}
```

### Direct Database Access

The bridge uses `psql` directly rather than the my-context CLI:

**Benefits:**
- Faster sync operations (no CLI startup overhead)
- Bulk insert/update capabilities
- Direct SQL query optimization

**Trade-offs:**
- Bypasses my-context validation hooks
- Notes may appear without normal my-context flow
- Must maintain schema compatibility manually

**SQL Examples:**

```sql
-- Import recent notes from my-context to timetracker
SELECT c.name, n.note, n.noted_at, n.note_type, n.tags
FROM contexts c
JOIN context_notes n ON c.id = n.context_id
WHERE n.noted_at > NOW() - INTERVAL '15 minutes'
ORDER BY n.noted_at DESC;

-- Export timetracker notes to my-context
INSERT INTO context_notes (context_id, note, noted_at, note_type, tags)
SELECT c.id, $1, $2, $3, $4
FROM contexts c
WHERE c.name = $5;
```

### Concurrency Considerations

Both my-context and timetracker may write to the same contexts simultaneously:

**Concurrent Write Scenarios:**
- User adds note via `my-context note "..."` while timetracker bridge is syncing
- Multiple timetracker sessions on different machines
- Background bridge `watch` mode while actively using my-context

**PostgreSQL Handles:**
- Row-level locking for INSERT/UPDATE operations
- Transaction isolation (READ COMMITTED default)
- Automatic conflict resolution at database level

**Best Practices:**
- Use timetracker bridge `watch` mode for continuous sync
- Avoid manual `psql` writes during active sessions
- Test schema changes against timetracker before deployment

### State Management

Each tool maintains its own active context state:

| Tool | State Location | Format |
|------|----------------|--------|
| my-context | PostgreSQL `state` table | JSON: `{"key": "active_context", "value": {"name": "..."}}` |
| timetracker | `~/.timetracker/context-state` | JSON file with session info |

**Important:** Active context state is NOT synchronized between tools. Each tool independently tracks which context is active in its own session.

## Schema Coordination

### Before Schema Changes

**Required:** Coordinate with timetracker maintainer before making schema changes to:

- `contexts` table
- `context_notes` table (especially `note_type` and `tags` columns)
- `context_files` table
- `state` table

**Process:**
1. Propose schema change in GitHub issue
2. Tag timetracker maintainer (@cdev)
3. Test migration with timetracker bridge
4. Deploy schema changes to dev environment first
5. Verify timetracker compatibility before production

### Schema Change Examples

**Safe Changes:**
- Adding new columns with defaults to `contexts`
- Adding new indexes
- Creating new tables (not referenced by bridge)

**Breaking Changes:**
- Removing or renaming `note_type` or `tags` columns
- Changing column types (TEXT → VARCHAR, etc.)
- Modifying foreign key constraints

## Integration Testing

### Manual Testing

Test timetracker bridge compatibility:

```bash
# Setup test environment
export MY_CONTEXT_HOME=db
export DATABASE_URL="postgres://user:pass@localhost:5432/test_db"

# Create test context
my-context start "Integration test"
my-context note "Test note from my-context"

# Verify bridge can read
tt-context-bridge sync

# Verify bridge can write (from timetracker)
# (Requires timetracker running)

# Verify my-context can read bridge notes
my-context show
```

### Automated Testing

**Future Enhancement:** Consider adding integration tests that:

1. Simulate bridge writes with extended metadata
2. Verify my-context can read/display bridge-created notes
3. Test concurrent writes from both tools
4. Validate schema compatibility

**Test Location:** `tests/integration/external_integration_test.go` (not yet implemented)

## Risk Assessment

| Area | Risk Level | Mitigation |
|------|------------|------------|
| Schema changes | Medium | Coordinate with timetracker before migrations |
| Note format | Low | Bridge uses standard note fields with optional extensions |
| Performance | Low | Bridge is read-heavy, minimal write impact |
| Data integrity | Low | PostgreSQL handles concurrency with row-level locking |
| State sync | Low | Each tool maintains separate active context state |

## Troubleshooting

### Bridge Sync Failures

**Problem:** Bridge fails to sync recent notes

**Solution:**
1. Check PostgreSQL connection: `psql $DATABASE_URL -c "SELECT 1"`
2. Verify schema exists: `my-context which`
3. Check note timestamps: `my-context show --json | jq .notes`

### Duplicate Notes

**Problem:** Same note appears multiple times

**Solution:**
- Check if bridge `watch` mode is running multiple times
- Verify bridge script logic for duplicate detection
- Query database directly: `SELECT * FROM context_notes WHERE note LIKE '%text%'`

### Missing Extended Metadata

**Problem:** Notes from timetracker don't show source information

**Solution:**
- Verify `note_type` column exists: `\d context_notes` in psql
- Check bridge script is setting metadata correctly
- Query raw metadata: `SELECT note_type, tags FROM context_notes WHERE id = ...`

## Contributing

### Adding New Integrations

To integrate a new tool with my-context:

1. **Use PostgreSQL Backend** - Set `MY_CONTEXT_HOME=db`
2. **Follow Schema Conventions** - Use existing table structures
3. **Document Metadata** - Add section to this guide if extending schema
4. **Test Compatibility** - Verify your tool works with my-context operations
5. **Submit PR** - Add integration documentation

### Documentation Updates

When updating this guide:

- Keep schema examples synchronized with actual code
- Update risk assessments as integration evolves
- Add troubleshooting entries for common issues
- Link to related my-context documentation

## Related Documentation

- [Database Schema Documentation](../internal/SCHEMA.md) — Complete schema reference with external tool guidelines
- [PostgreSQL Backend Setup](SETUP.md#postgresql-backend) — Database configuration
- [PostgreSQL Backend Code](../../pkg/storage/postgres/postgres_backend.go) — Implementation details
- [Troubleshooting Guide](TROUBLESHOOTING.md) — Common issues and solutions
- [Contributing Guide](../../CONTRIBUTING.md) — Development workflow

## Future Enhancements

Potential improvements for external integrations:

- [ ] Webhook/event system for real-time notifications
- [ ] REST API for structured integration
- [ ] Integration SDK/library
- [ ] Validation hooks for external writes
- [ ] Cross-tool active context synchronization
- [ ] Standardized metadata schema for external tools
