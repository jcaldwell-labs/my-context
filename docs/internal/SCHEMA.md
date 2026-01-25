# Database Schema Documentation

This document describes the PostgreSQL database schema used by my-context, including fields available for external integrations.

## Overview

The my-context database schema consists of five main tables:

1. **contexts** - Context metadata and lifecycle
2. **context_notes** - Timestamped notes with optional metadata
3. **context_files** - File associations
4. **context_touches** - Activity timestamps
5. **state** - Key-value state store

## Schema Definitions

### contexts

Primary table for context information.

```sql
CREATE TABLE contexts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active' NOT NULL,
    project VARCHAR(255),
    started_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    stopped_at TIMESTAMPTZ,
    touch_count INTEGER DEFAULT 0 NOT NULL,
    last_touch_at TIMESTAMPTZ,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_activity_at TIMESTAMPTZ,
    completion_status VARCHAR(20) DEFAULT 'in_progress',
    resume_notes TEXT,
    CONSTRAINT valid_completion_status CHECK (
        completion_status IN ('in_progress', 'paused', 'complete', 'archived')
    )
);

-- Indexes
CREATE INDEX idx_contexts_status ON contexts(status);
CREATE INDEX idx_contexts_name ON contexts(name);
CREATE INDEX idx_contexts_project ON contexts(project);
CREATE INDEX idx_contexts_started_at ON contexts(started_at DESC);
CREATE INDEX idx_contexts_name_search ON contexts USING GIN (to_tsvector('english', name));
```

**Key Fields:**
- `metadata` - JSONB for arbitrary context metadata (labels, parent, created_by)
- `completion_status` - Tracks context lifecycle (in_progress, paused, complete, archived)
- `project` - Optional project grouping

### context_notes

Timestamped notes with support for external tool metadata.

```sql
CREATE TABLE context_notes (
    id SERIAL PRIMARY KEY,
    context_id INTEGER NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    noted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    note_type VARCHAR(50),  -- Available for external tool metadata
    tags TEXT[]              -- Array of tags/labels for categorization
);

-- Indexes
CREATE INDEX idx_notes_context ON context_notes(context_id);
CREATE INDEX idx_notes_noted_at ON context_notes(noted_at DESC);
CREATE INDEX idx_notes_search ON context_notes USING GIN (to_tsvector('english', note));
```

**Key Fields for External Integrations:**

#### note_type (VARCHAR(50))

This field is **reserved for external tool metadata**. My-context core does not use this field, making it safe for external tools to store structured metadata.

**Usage Patterns:**

1. **Simple string identifier:**
   ```sql
   note_type = 'timetracker:session'
   ```

2. **JSON metadata (recommended):**
   ```sql
   note_type = '{"source": "shell:ses-123", "session_id": "ses-xyz", "backdated": false}'
   ```

3. **Structured string:**
   ```sql
   note_type = 'tool=timetracker;session=ses-123;type=checkpoint'
   ```

**Best Practices:**
- Start with tool identifier (e.g., `timetracker:`, `external:`) to avoid collisions
- Keep under 50 characters or use abbreviated JSON
- Use consistent format across your tool for parsing
- Consider using `tags` array for longer metadata

#### tags (TEXT[])

PostgreSQL text array for categorization and filtering.

**Usage Examples:**

```sql
-- Single tag
tags = ARRAY['bug-fix']

-- Multiple tags
tags = ARRAY['sprint-42', 'backend', 'api']

-- External tool metadata tags
tags = ARRAY['tt:session-123', 'tt:ticket-PROJ-456', 'backdated']
```

**Querying tags:**

```sql
-- Find notes with specific tag
SELECT * FROM context_notes WHERE 'bug-fix' = ANY(tags);

-- Find notes with any of multiple tags
SELECT * FROM context_notes WHERE tags && ARRAY['sprint-42', 'backend'];

-- Find notes with external tool tags
SELECT * FROM context_notes WHERE tags @> ARRAY['tt:session-123'];
```

**Best Practices:**
- Use consistent naming conventions (e.g., `tool:value` prefix)
- Keep individual tags under 50 characters
- Use array operations for efficient querying
- Consider indexing: `CREATE INDEX idx_notes_tags ON context_notes USING GIN(tags);`

### context_files

File associations with optional metadata.

```sql
CREATE TABLE context_files (
    id SERIAL PRIMARY KEY,
    context_id INTEGER NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    file_type VARCHAR(50),    -- Optional: MIME type or extension
    file_size BIGINT          -- Optional: File size in bytes
);

-- Indexes
CREATE INDEX idx_files_context ON context_files(context_id);
CREATE INDEX idx_files_added_at ON context_files(added_at DESC);
CREATE INDEX idx_files_path_search ON context_files USING GIN (to_tsvector('english', file_path));
```

**Key Fields:**
- `file_type` - MIME type or extension (e.g., `text/plain`, `.go`)
- `file_size` - File size at time of association

### context_touches

Activity timestamps for context access tracking.

```sql
CREATE TABLE context_touches (
    id SERIAL PRIMARY KEY,
    context_id INTEGER NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
    touched_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexes
CREATE INDEX idx_touches_context ON context_touches(context_id);
CREATE INDEX idx_touches_touched_at ON context_touches(touched_at DESC);
```

**Usage:**
- Record whenever context is accessed or modified
- Updated by `my-context touch` command
- Used for "last activity" tracking

### state

Generic key-value store for application state.

```sql
CREATE TABLE state (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_state_key ON state(key);
```

**Reserved Keys:**
- `active_context` - Name of currently active context

**Example:**
```sql
-- Get active context
SELECT value->>'name' FROM state WHERE key = 'active_context';

-- Set active context
INSERT INTO state (key, value) 
VALUES ('active_context', '{"name": "My_Context"}')
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP;
```

## External Tool Guidelines

### Schema Compatibility

When integrating external tools:

1. **Use existing fields first** - Leverage `note_type`, `tags`, and `metadata` JSONB columns
2. **Coordinate additions** - Propose new columns via GitHub issue before adding
3. **Follow conventions** - Use tool-prefixed identifiers to avoid collisions
4. **Test migrations** - Verify schema changes don't break existing tools

### Safe Schema Changes

These changes are safe and won't break external integrations:

- Adding new columns with `DEFAULT` values
- Adding new indexes
- Creating new tables (not referenced by existing tools)
- Extending VARCHAR limits

### Breaking Schema Changes

These changes require coordination with external tool maintainers:

- Removing columns (especially `note_type`, `tags`, `metadata`)
- Renaming columns
- Changing column types (TEXT → VARCHAR, JSONB → JSON, etc.)
- Modifying foreign key constraints
- Changing index types (B-tree → GIN, etc.)

### Metadata Field Reservation

| Field | Owner | Purpose |
|-------|-------|---------|
| `context_notes.note_type` | External tools | Tool-specific metadata |
| `context_notes.tags` | Shared | Categorization/filtering |
| `contexts.metadata` | my-context | Core context metadata (labels, parent, created_by) |
| `context_files.file_type` | Shared | File classification |
| `context_files.file_size` | Shared | File metadata |

## Example Queries

### External Tool Integration

```sql
-- Add note with external metadata
INSERT INTO context_notes (context_id, note, noted_at, note_type, tags)
SELECT c.id, 'Work completed', NOW(), 
       '{"source": "timetracker", "session_id": "ses-123"}',
       ARRAY['timetracker', 'session-checkpoint']
FROM contexts c
WHERE c.name = 'My_Context';

-- Query notes from specific external tool
SELECT n.note, n.noted_at, n.note_type
FROM context_notes n
JOIN contexts c ON n.context_id = c.id
WHERE c.name = 'My_Context'
  AND n.note_type LIKE 'timetracker%'
ORDER BY n.noted_at DESC;

-- Find all contexts touched by external tool
SELECT DISTINCT c.name, c.status
FROM contexts c
JOIN context_notes n ON n.context_id = c.id
WHERE n.tags @> ARRAY['timetracker'];
```

### Bulk Operations

```sql
-- Import notes from external source
INSERT INTO context_notes (context_id, note, noted_at, note_type, tags)
SELECT 
    c.id,
    external.note_text,
    external.timestamp,
    'imported:' || external.source_tool,
    ARRAY[external.source_tool, 'imported']
FROM external_notes_staging external
JOIN contexts c ON c.name = external.context_name;

-- Export notes to external format
SELECT 
    c.name AS context_name,
    n.note,
    n.noted_at,
    n.note_type,
    n.tags
FROM context_notes n
JOIN contexts c ON n.context_id = c.id
WHERE n.noted_at > NOW() - INTERVAL '24 hours'
ORDER BY c.name, n.noted_at;
```

## Schema Evolution

### Version History

- **v1.0.0** - Initial schema with basic tables
- **v2.0.0** - Added `note_type` and `tags` to context_notes
- **v3.0.0** - Added partition support via schemas

### Future Considerations

Potential schema enhancements under discussion:

- Custom indexes on `note_type` for external tool performance
- Additional metadata columns on `context_files`
- Event/webhook table for real-time integration notifications
- Audit log table for tracking external tool writes

## Related Documentation

- [External Integrations Guide](EXTERNAL-INTEGRATIONS.md)
- [PostgreSQL Backend Code](../../pkg/storage/postgres/postgres_backend.go)
- [Storage Interface](../../pkg/storage/storage_interface.go)
- [Data Models](../../pkg/models/context.go)
