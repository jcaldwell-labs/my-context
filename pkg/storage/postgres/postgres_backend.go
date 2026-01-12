package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/jefferycaldwell/my-context-copilot/pkg/storage"
	"github.com/jefferycaldwell/my-context-copilot/pkg/utils"
)

// Backend implements storage.Backend for PostgreSQL
type Backend struct {
	db        *sql.DB
	connStr   string
	schema    string // PostgreSQL schema name for partitioning
	partition string // Friendly partition name (e.g., "adventure-engine")
}

// NewPostgresBackend creates a new PostgreSQL backend
func NewPostgresBackend(connStr string) *Backend {
	// Extract schema from connection string search_path parameter
	schema := extractSchemaFromConnStr(connStr)
	if schema == "" {
		schema = "public" // Default schema if not specified
	}

	// Extract partition name from schema (reverse of sanitization)
	// For display purposes: adventure_engine -> adventure-engine
	partition := schema
	if schema != "public" {
		// Keep sanitized name for now, could add reverse mapping later
		partition = schema
	}

	return &Backend{
		connStr:   connStr,
		schema:    schema,
		partition: partition,
	}
}

// extractSchemaFromConnStr parses the search_path parameter from connection string
func extractSchemaFromConnStr(connStr string) string {
	// Look for search_path=schema_name in connection string
	// Format: "host=... search_path=my_schema ..."
	parts := splitConnStr(connStr)
	for _, part := range parts {
		if len(part) > 12 && part[:12] == "search_path=" {
			return part[12:]
		}
	}
	return "public"
}

// splitConnStr splits connection string by spaces, respecting quoted values.
// Handles escaped quotes in PostgreSQL connection strings:
// - Doubled single quotes (”) are treated as escaped quotes within a value
// - Example: key='value”s data' contains the value: value's data
func splitConnStr(connStr string) []string {
	var parts []string
	var current string
	inQuote := false
	runes := []rune(connStr)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if char == '\'' {
			// Check for escaped quote (doubled single quote)
			if inQuote && i+1 < len(runes) && runes[i+1] == '\'' {
				// Escaped quote - add both quotes and skip the next one
				current += "''"
				i++
				continue
			}
			// Toggle quote state
			inQuote = !inQuote
			current += string(char)
		} else if char == ' ' && !inQuote {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// Init initializes the database connection and schema
func (b *Backend) Init() error {
	db, err := sql.Open("postgres", b.connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for CLI tool usage.
	// Conservative defaults to avoid connection exhaustion.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close() // Clean up on error
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	b.db = db

	// Create schema (includes security validation)
	if err := b.createSchema(); err != nil {
		db.Close() // Clean up on error to prevent resource leak
		b.db = nil
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Health checks if the database is healthy
func (b *Backend) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return b.db.PingContext(ctx)
}

// Close closes the database connection
func (b *Backend) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

// GetSchema returns the PostgreSQL schema name for this backend
func (b *Backend) GetSchema() string {
	return b.schema
}

// GetPartition returns the friendly partition name for this backend
func (b *Backend) GetPartition() string {
	return b.partition
}

// GetDB returns the underlying database connection for cross-partition queries
func (b *Backend) GetDB() *sql.DB {
	return b.db
}

// scanContextRow scans a context row from the database and returns a ContextWithMetadata.
// This helper reduces code duplication across listing functions.
// Expects columns: name, status, started_at, stopped_at, project, completion_status, metadata, touch_count, last_touch_at
func scanContextRow(scanner interface{ Scan(...any) error }) (*models.ContextWithMetadata, error) {
	var (
		name             string
		status           string
		startedAt        time.Time
		stoppedAt        sql.NullTime
		project          sql.NullString
		completionStatus sql.NullString
		metadataJSON     []byte
		touchCount       int
		lastTouchAt      sql.NullTime
	)

	err := scanner.Scan(&name, &status, &startedAt, &stoppedAt, &project, &completionStatus, &metadataJSON, &touchCount, &lastTouchAt)
	if err != nil {
		return nil, err
	}

	var metadata models.ContextMetadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			// Log error but continue - don't fail for one bad record
			fmt.Fprintf(os.Stderr, "warning: failed to parse metadata for context %s: %v\n", name, err)
		}
	}

	ctx := &models.ContextWithMetadata{
		Name:       name,
		StartTime:  startedAt,
		Status:     status,
		IsArchived: completionStatus.Valid && completionStatus.String == "archived",
		Metadata:   metadata,
		TouchCount: touchCount,
	}

	if stoppedAt.Valid {
		ctx.EndTime = &stoppedAt.Time
	}

	if lastTouchAt.Valid {
		ctx.LastTouchAt = &lastTouchAt.Time
	}

	return ctx, nil
}

// createSchema creates the database schema if it doesn't exist
func (b *Backend) createSchema() error {
	// Security: Validate schema name before using in SQL to prevent SQL injection.
	// The schema name comes from user input via MY_CONTEXT_HOME=db:partition-name
	// and is sanitized by SanitizePartitionName(), but we add defense-in-depth here.
	if !utils.IsValidPostgresIdentifier(b.schema) {
		return fmt.Errorf("invalid schema name %q: must contain only letters, digits, and underscores", b.schema)
	}

	// Create schema if it doesn't exist (skip for public schema)
	// Use pq.QuoteIdentifier for safe schema quoting in addition to validation
	if b.schema != "public" {
		_, err := b.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", pq.QuoteIdentifier(b.schema)))
		if err != nil {
			return fmt.Errorf("failed to create schema %s: %w", b.schema, err)
		}
	}

	// Check if tables already exist in this schema
	var exists bool
	err := b.db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = $1 AND tablename = 'contexts')", b.schema).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if tables exist: %w", err)
	}

	if exists {
		// Tables already exist in this schema, skip table creation
		return nil
	}

	schema := `
	-- Contexts table
	CREATE TABLE IF NOT EXISTS contexts (
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
		CONSTRAINT valid_completion_status CHECK (completion_status IN ('in_progress', 'paused', 'complete', 'archived'))
	);

	CREATE INDEX IF NOT EXISTS idx_contexts_status ON contexts(status);
	CREATE INDEX IF NOT EXISTS idx_contexts_name ON contexts(name);
	CREATE INDEX IF NOT EXISTS idx_contexts_project ON contexts(project);
	CREATE INDEX IF NOT EXISTS idx_contexts_started_at ON contexts(started_at DESC);
	CREATE INDEX IF NOT EXISTS idx_contexts_name_search ON contexts USING GIN (to_tsvector('english', name));

	-- Notes table
	CREATE TABLE IF NOT EXISTS context_notes (
		id SERIAL PRIMARY KEY,
		context_id INTEGER NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
		note TEXT NOT NULL,
		noted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		note_type VARCHAR(50),
		tags TEXT[]
	);

	CREATE INDEX IF NOT EXISTS idx_notes_context ON context_notes(context_id);
	CREATE INDEX IF NOT EXISTS idx_notes_noted_at ON context_notes(noted_at DESC);
	CREATE INDEX IF NOT EXISTS idx_notes_search ON context_notes USING GIN (to_tsvector('english', note));

	-- Files table
	CREATE TABLE IF NOT EXISTS context_files (
		id SERIAL PRIMARY KEY,
		context_id INTEGER NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
		file_path TEXT NOT NULL,
		added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
		file_type VARCHAR(50),
		file_size BIGINT
	);

	CREATE INDEX IF NOT EXISTS idx_files_context ON context_files(context_id);
	CREATE INDEX IF NOT EXISTS idx_files_added_at ON context_files(added_at DESC);
	CREATE INDEX IF NOT EXISTS idx_files_path_search ON context_files USING GIN (to_tsvector('english', file_path));

	-- State table (generic key-value store)
	CREATE TABLE IF NOT EXISTS state (
		id SERIAL PRIMARY KEY,
		key VARCHAR(255) UNIQUE NOT NULL,
		value JSONB NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_state_key ON state(key);
	`

	_, err = b.db.Exec(schema)
	return err
}

// CreateContext creates a new context
func (b *Backend) CreateContext(ctx *models.ContextWithMetadata) error {
	// Validate context
	if err := ctx.Validate(); err != nil {
		return fmt.Errorf("context validation failed: %w", err)
	}

	// Start a transaction
	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Rollback is safe to ignore; commit success/failure is what matters

	// Stop any currently active context
	_, err = tx.Exec(`
		UPDATE contexts
		SET status = 'stopped', stopped_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE status = 'active'
	`)
	if err != nil {
		return fmt.Errorf("failed to stop active contexts: %w", err)
	}

	// Marshal metadata to JSON
	metadataJSON, err := json.Marshal(ctx.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Insert new context
	_, err = tx.Exec(`
		INSERT INTO contexts (
			name,
			status,
			started_at,
			stopped_at,
			project,
			metadata,
			completion_status,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`,
		ctx.Name,
		ctx.Status,
		ctx.StartTime,
		nil, // stopped_at is NULL for new contexts
		nil, // project (we'll extract from name if it contains ":")
		metadataJSON,
		"in_progress", // default completion_status
	)
	if err != nil {
		return fmt.Errorf("failed to insert context: %w", err)
	}

	// Set as active context
	activeJSON, _ := json.Marshal(ctx.Name)
	_, err = tx.Exec(`
		INSERT INTO state (key, value, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE
		SET value = $2, updated_at = CURRENT_TIMESTAMP
	`, "active_context", activeJSON)
	if err != nil {
		return fmt.Errorf("failed to set active context: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetContext retrieves a context by name
func (b *Backend) GetContext(name string) (*models.ContextWithMetadata, error) {
	query := `
		SELECT
			name,
			status,
			started_at,
			stopped_at,
			project,
			completion_status,
			metadata,
			touch_count,
			last_touch_at
		FROM contexts
		WHERE name = $1
	`

	var (
		contextName      string
		status           string
		startedAt        time.Time
		stoppedAt        sql.NullTime
		project          sql.NullString
		completionStatus sql.NullString
		metadataJSON     []byte
		touchCount       int
		lastTouchAt      sql.NullTime
	)

	err := b.db.QueryRow(query, name).Scan(
		&contextName,
		&status,
		&startedAt,
		&stoppedAt,
		&project,
		&completionStatus,
		&metadataJSON,
		&touchCount,
		&lastTouchAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("context '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query context: %w", err)
	}

	// Parse metadata JSON
	var metadata models.ContextMetadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, fmt.Errorf("failed to parse context metadata: %w", err)
		}
	}

	// Build context object
	ctx := &models.ContextWithMetadata{
		Name:       contextName,
		StartTime:  startedAt,
		Status:     status,
		IsArchived: completionStatus.Valid && completionStatus.String == "archived",
		Metadata:   metadata,
		TouchCount: touchCount,
	}

	if stoppedAt.Valid {
		ctx.EndTime = &stoppedAt.Time
	}

	if lastTouchAt.Valid {
		ctx.LastTouchAt = &lastTouchAt.Time
	}

	return ctx, nil
}

// ListContexts lists all contexts
func (b *Backend) ListContexts() ([]*models.ContextWithMetadata, error) {
	query := `
		SELECT
			name,
			status,
			started_at,
			stopped_at,
			project,
			completion_status,
			metadata,
			touch_count,
			last_touch_at
		FROM contexts
		ORDER BY started_at DESC
	`

	rows, err := b.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query contexts: %w", err)
	}
	defer rows.Close()

	var contexts []*models.ContextWithMetadata
	for rows.Next() {
		ctx, err := scanContextRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan context row: %w", err)
		}
		contexts = append(contexts, ctx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating context rows: %w", err)
	}

	return contexts, nil
}

// UpdateContext updates an existing context
func (b *Backend) UpdateContext(ctx *models.ContextWithMetadata) error {
	// Validate context
	if err := ctx.Validate(); err != nil {
		return fmt.Errorf("context validation failed: %w", err)
	}

	// Marshal metadata
	metadataJSON, err := json.Marshal(ctx.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Determine completion status based on context state
	completionStatus := "in_progress"
	if ctx.IsArchived {
		completionStatus = "archived"
	} else if ctx.Status == "stopped" {
		completionStatus = "paused"
	}

	// Update context
	_, err = b.db.Exec(`
		UPDATE contexts
		SET
			status = $1,
			stopped_at = $2,
			metadata = $3,
			completion_status = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE name = $5
	`,
		ctx.Status,
		ctx.EndTime,
		metadataJSON,
		completionStatus,
		ctx.Name,
	)

	if err != nil {
		return fmt.Errorf("failed to update context: %w", err)
	}

	return nil
}

// DeleteContext deletes a context
func (b *Backend) DeleteContext(name string) error {
	// Delete will CASCADE to notes, files due to foreign keys
	result, err := b.db.Exec(`DELETE FROM contexts WHERE name = $1`, name)
	if err != nil {
		return fmt.Errorf("failed to delete context: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("context '%s' not found", name)
	}

	return nil
}

// ArchiveContext archives a context
func (b *Backend) ArchiveContext(name string) error {
	// Update context to archived status
	_, err := b.db.Exec(`
		UPDATE contexts
		SET
			completion_status = 'archived',
			updated_at = CURRENT_TIMESTAMP
		WHERE name = $1
	`, name)

	if err != nil {
		return fmt.Errorf("failed to archive context: %w", err)
	}

	return nil
}

// AddNote adds a note to a context
//
//nolint:dupl // AddNote and AddFile have similar structure but operate on different tables
func (b *Backend) AddNote(contextName, timestamp, content string) error {
	// Parse timestamp
	noteTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Get context ID
	var contextID int64
	err = b.db.QueryRow(`SELECT id FROM contexts WHERE name = $1`, contextName).Scan(&contextID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("context '%s' not found", contextName)
	}
	if err != nil {
		return fmt.Errorf("failed to query context: %w", err)
	}

	// Start transaction
	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Rollback is safe to ignore; commit success/failure is what matters

	// Insert note
	_, err = tx.Exec(`
		INSERT INTO context_notes (context_id, note, noted_at)
		VALUES ($1, $2, $3)
	`, contextID, content, noteTime)
	if err != nil {
		return fmt.Errorf("failed to insert note: %w", err)
	}

	// Update last_activity_at
	_, err = tx.Exec(`
		UPDATE contexts
		SET last_activity_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, noteTime, contextID)
	if err != nil {
		return fmt.Errorf("failed to update last_activity_at: %w", err)
	}

	// Commit
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetNotes retrieves all notes for a context
func (b *Backend) GetNotes(contextName string) ([]storage.Note, error) {
	query := `
		SELECT n.id, n.context_id, n.noted_at, n.note
		FROM context_notes n
		JOIN contexts c ON n.context_id = c.id
		WHERE c.name = $1
		ORDER BY n.noted_at ASC
	`

	rows, err := b.db.Query(query, contextName)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []storage.Note
	for rows.Next() {
		var note storage.Note
		err := rows.Scan(&note.ID, &note.ContextID, &note.Timestamp, &note.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}

	return notes, nil
}

// GetNotesByTimestamp retrieves notes within a time range
func (b *Backend) GetNotesByTimestamp(contextName string, after, before time.Time) ([]storage.Note, error) {
	return nil, fmt.Errorf("GetNotesByTimestamp: not implemented")
}

// AddFile adds a file reference to a context
//
//nolint:dupl // AddFile and AddNote have similar structure but operate on different tables
func (b *Backend) AddFile(contextName, timestamp, path string) error {
	// Parse timestamp
	addTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Get context ID
	var contextID int64
	err = b.db.QueryRow(`SELECT id FROM contexts WHERE name = $1`, contextName).Scan(&contextID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("context '%s' not found", contextName)
	}
	if err != nil {
		return fmt.Errorf("failed to query context: %w", err)
	}

	// Start transaction
	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Rollback is safe to ignore; commit success/failure is what matters

	// Insert file reference
	_, err = tx.Exec(`
		INSERT INTO context_files (context_id, file_path, added_at)
		VALUES ($1, $2, $3)
	`, contextID, path, addTime)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	// Update last_activity_at
	_, err = tx.Exec(`
		UPDATE contexts
		SET last_activity_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, addTime, contextID)
	if err != nil {
		return fmt.Errorf("failed to update last_activity_at: %w", err)
	}

	// Commit
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetFiles retrieves all files for a context
func (b *Backend) GetFiles(contextName string) ([]storage.File, error) {
	query := `
		SELECT f.id, f.context_id, f.added_at, f.file_path
		FROM context_files f
		JOIN contexts c ON f.context_id = c.id
		WHERE c.name = $1
		ORDER BY f.added_at ASC
	`

	rows, err := b.db.Query(query, contextName)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []storage.File
	for rows.Next() {
		var file storage.File
		err := rows.Scan(&file.ID, &file.ContextID, &file.Timestamp, &file.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating files: %w", err)
	}

	return files, nil
}

// GetFilesByTimestamp retrieves files within a time range
func (b *Backend) GetFilesByTimestamp(contextName string, after, before time.Time) ([]storage.File, error) {
	return nil, fmt.Errorf("GetFilesByTimestamp: not implemented")
}

// AddTouch records a context access
func (b *Backend) AddTouch(contextName, timestamp string) error {
	// Parse timestamp
	touchTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Update context touch tracking
	_, err = b.db.Exec(`
		UPDATE contexts
		SET
			touch_count = touch_count + 1,
			last_touch_at = $1,
			last_activity_at = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE name = $2
	`, touchTime, contextName)

	if err != nil {
		return fmt.Errorf("failed to update touch: %w", err)
	}

	return nil
}

// GetTouches retrieves all touches for a context
func (b *Backend) GetTouches(contextName string) ([]storage.Touch, error) {
	return nil, fmt.Errorf("GetTouches: not implemented")
}

// LogTransition logs a context transition
func (b *Backend) LogTransition(transition *storage.Transition) error {
	return fmt.Errorf("LogTransition: not implemented")
}

// GetTransitions retrieves recent transitions
func (b *Backend) GetTransitions(limit int) ([]storage.Transition, error) {
	return nil, fmt.Errorf("GetTransitions: not implemented")
}

// GetTransitionsByContext retrieves transitions for a specific context
func (b *Backend) GetTransitionsByContext(contextName string) ([]storage.Transition, error) {
	return nil, fmt.Errorf("GetTransitionsByContext: not implemented")
}

// GetActiveContext retrieves the currently active context
func (b *Backend) GetActiveContext() (string, error) {
	var valueJSON []byte
	err := b.db.QueryRow(`
		SELECT value FROM state WHERE key = $1
	`, "active_context").Scan(&valueJSON)

	if err == sql.ErrNoRows {
		// No active context set
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to query active context: %w", err)
	}

	// Parse JSON value (it might be stored as {"name": "context_name"} or just "context_name")
	var result struct {
		Name string `json:"name"`
	}

	// Try to unmarshal as object first
	if err := json.Unmarshal(valueJSON, &result); err == nil && result.Name != "" {
		return result.Name, nil
	}

	// Try as plain string
	var name string
	if err := json.Unmarshal(valueJSON, &name); err == nil {
		return name, nil
	}

	// If all else fails, return empty (no active context)
	return "", nil
}

// SetActiveContext sets the active context
func (b *Backend) SetActiveContext(contextName string) error {
	valueJSON, err := json.Marshal(contextName)
	if err != nil {
		return fmt.Errorf("failed to marshal context name: %w", err)
	}

	_, err = b.db.Exec(`
		INSERT INTO state (key, value, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE
		SET value = $2, updated_at = CURRENT_TIMESTAMP
	`, "active_context", valueJSON)

	if err != nil {
		return fmt.Errorf("failed to set active context: %w", err)
	}

	return nil
}

// ClearActiveContext clears the active context
func (b *Backend) ClearActiveContext() error {
	_, err := b.db.Exec(`DELETE FROM state WHERE key = $1`, "active_context")
	if err != nil {
		return fmt.Errorf("failed to clear active context: %w", err)
	}
	return nil
}

// GetContextsByLabel retrieves contexts with a specific label
//
//nolint:dupl // GetContextsByLabel and GetContextsByParent share query structure but filter differently
func (b *Backend) GetContextsByLabel(label string) ([]*models.ContextWithMetadata, error) {
	// Use JSONB containment operator to find contexts with this label
	// Cast to text explicitly to avoid type inference issues
	query := `
		SELECT
			name,
			status,
			started_at,
			stopped_at,
			project,
			completion_status,
			metadata,
			touch_count,
			last_touch_at
		FROM contexts
		WHERE metadata->'labels' ? $1::text
		ORDER BY started_at DESC
	`

	rows, err := b.db.Query(query, label)
	if err != nil {
		return nil, fmt.Errorf("failed to query contexts by label: %w", err)
	}
	defer rows.Close()

	var contexts []*models.ContextWithMetadata
	for rows.Next() {
		ctx, err := scanContextRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan context: %w", err)
		}
		contexts = append(contexts, ctx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating contexts: %w", err)
	}

	return contexts, nil
}

// GetContextsByParent retrieves child contexts of a parent
//
//nolint:dupl // GetContextsByParent and GetContextsByLabel share query structure but filter differently
func (b *Backend) GetContextsByParent(parent string) ([]*models.ContextWithMetadata, error) {
	// Query contexts where metadata.parent = parent
	query := `
		SELECT
			name,
			status,
			started_at,
			stopped_at,
			project,
			completion_status,
			metadata,
			touch_count,
			last_touch_at
		FROM contexts
		WHERE metadata->>'parent' = $1
		ORDER BY started_at DESC
	`

	rows, err := b.db.Query(query, parent)
	if err != nil {
		return nil, fmt.Errorf("failed to query contexts by parent: %w", err)
	}
	defer rows.Close()

	var contexts []*models.ContextWithMetadata
	for rows.Next() {
		ctx, err := scanContextRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan context: %w", err)
		}
		contexts = append(contexts, ctx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating contexts: %w", err)
	}

	return contexts, nil
}

// SearchContextsByName searches contexts by name
func (b *Backend) SearchContextsByName(query string) ([]*models.ContextWithMetadata, error) {
	return nil, fmt.Errorf("SearchContextsByName: not implemented")
}

// SearchNotes searches notes across all contexts
func (b *Backend) SearchNotes(query string) (map[string][]storage.Note, error) {
	// Use full-text search with GIN index
	sqlQuery := `
		SELECT c.name, n.id, n.context_id, n.noted_at, n.note
		FROM context_notes n
		JOIN contexts c ON n.context_id = c.id
		WHERE to_tsvector('english', n.note) @@ plainto_tsquery('english', $1)
		ORDER BY n.noted_at DESC
	`

	rows, err := b.db.Query(sqlQuery, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	// Group notes by context name
	notesByContext := make(map[string][]storage.Note)

	for rows.Next() {
		var (
			contextName string
			note        storage.Note
		)

		err := rows.Scan(&contextName, &note.ID, &note.ContextID, &note.Timestamp, &note.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		notesByContext[contextName] = append(notesByContext[contextName], note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notes: %w", err)
	}

	return notesByContext, nil
}

// GetContextCount returns the total number of contexts
func (b *Backend) GetContextCount() (int, error) {
	var count int
	err := b.db.QueryRow("SELECT COUNT(*) FROM contexts").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count contexts: %w", err)
	}
	return count, nil
}

// GetContextStats retrieves statistics for a single context
func (b *Backend) GetContextStats(contextName string) (storage.ContextStats, error) {
	return storage.ContextStats{}, fmt.Errorf("GetContextStats: not implemented")
}

// GetGlobalStats retrieves global statistics
func (b *Backend) GetGlobalStats() (storage.GlobalStats, error) {
	return storage.GlobalStats{}, fmt.Errorf("GetGlobalStats: not implemented")
}

// Cross-Partition Query Methods
// These methods query across all partitions (schemas) for global operations

// PartitionInfo holds information about a partition and its contents
type PartitionInfo struct {
	Schema       string
	ContextCount int
	Contexts     []*models.ContextWithMetadata
}

// ListAllPartitions returns all partitions (schemas) that contain my-context data
func (b *Backend) ListAllPartitions() ([]string, error) {
	query := `
		SELECT schemaname
		FROM pg_tables
		WHERE tablename = 'contexts'
		  AND schemaname NOT IN ('pg_catalog', 'information_schema')
		ORDER BY schemaname
	`

	rows, err := b.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query partitions: %w", err)
	}
	defer rows.Close()

	var partitions []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, fmt.Errorf("failed to scan partition: %w", err)
		}
		partitions = append(partitions, schema)
	}

	return partitions, rows.Err()
}

// ListContextsAcrossPartitions lists contexts from all partitions
func (b *Backend) ListContextsAcrossPartitions() (map[string][]*models.ContextWithMetadata, error) {
	partitions, err := b.ListAllPartitions()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*models.ContextWithMetadata)

	for _, partition := range partitions {
		// Security: Validate partition name before using in SQL (defense-in-depth)
		if !utils.IsValidPostgresIdentifier(partition) {
			continue // Skip invalid partition names
		}

		// Query contexts from this partition using pq.QuoteIdentifier for safe schema quoting
		query := fmt.Sprintf(`
			SELECT
				name,
				status,
				started_at,
				stopped_at,
				project,
				completion_status,
				metadata,
				touch_count,
				last_touch_at
			FROM %s.contexts
			ORDER BY started_at DESC
		`, pq.QuoteIdentifier(partition))

		rows, err := b.db.Query(query)
		if err != nil {
			// Skip partitions we can't query
			continue
		}

		var contexts []*models.ContextWithMetadata
		for rows.Next() {
			ctx, err := scanContextRow(rows)
			if err != nil {
				rows.Close()
				continue
			}
			contexts = append(contexts, ctx)
		}
		rows.Close()

		result[partition] = contexts
	}

	return result, nil
}

// SearchContextsAcrossPartitions searches for contexts by name across all partitions
func (b *Backend) SearchContextsAcrossPartitions(query string) (map[string][]*models.ContextWithMetadata, error) {
	partitions, err := b.ListAllPartitions()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*models.ContextWithMetadata)

	for _, partition := range partitions {
		// Security: Validate partition name before using in SQL (defense-in-depth)
		if !utils.IsValidPostgresIdentifier(partition) {
			continue // Skip invalid partition names
		}

		// Search contexts in this partition using full-text search and pq.QuoteIdentifier for safe schema quoting
		searchQuery := fmt.Sprintf(`
			SELECT
				name,
				status,
				started_at,
				stopped_at,
				project,
				completion_status,
				metadata,
				touch_count,
				last_touch_at
			FROM %s.contexts
			WHERE to_tsvector('english', name) @@ plainto_tsquery('english', $1)
			ORDER BY started_at DESC
		`, pq.QuoteIdentifier(partition))

		rows, err := b.db.Query(searchQuery, query)
		if err != nil {
			// Skip partitions we can't query
			continue
		}

		var contexts []*models.ContextWithMetadata
		for rows.Next() {
			ctx, err := scanContextRow(rows)
			if err != nil {
				rows.Close()
				continue
			}
			contexts = append(contexts, ctx)
		}
		rows.Close()

		if len(contexts) > 0 {
			result[partition] = contexts
		}
	}

	return result, nil
}
