// Package storage defines the storage backend interface for my-context.
// It provides abstractions for persisting contexts, notes, files, and state
// across different storage implementations (file-based, PostgreSQL, etc.).
package storage

import (
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// Backend defines the interface for context storage implementations.
//
// Implementations must be safe for concurrent use from multiple goroutines.
// All methods should return descriptive errors wrapping the underlying cause.
//
// Error Handling:
//   - Return nil on success
//   - Return wrapped errors with context (e.g., "failed to create context: %w")
//   - Return specific errors for "not found" conditions
//
// Thread Safety:
//   - All methods must be safe for concurrent access
//   - Implementations should use appropriate locking or transactions
//
// Available Implementations:
//   - File-based: stores data in ~/.my-context/ as JSON/log files
//   - PostgreSQL: stores data in PostgreSQL with schema-based partitioning
type Backend interface {
	// Init initializes the backend connection and creates necessary schema/tables.
	// Must be called before any other operations. Returns an error if initialization fails.
	Init() error

	// Health checks if the backend is healthy and accessible.
	// Returns nil if healthy, error otherwise.
	Health() error

	// Close releases any resources held by the backend.
	// Should be called when the backend is no longer needed.
	Close() error

	// CreateContext creates a new context and sets it as active.
	// Automatically stops any previously active context.
	// Returns an error if the context already exists or validation fails.
	CreateContext(ctx *models.ContextWithMetadata) error

	// GetContext retrieves a context by its exact name.
	// Returns an error if the context is not found.
	GetContext(name string) (*models.ContextWithMetadata, error)

	// ListContexts returns all contexts, ordered by start time (newest first).
	// Returns an empty slice if no contexts exist.
	ListContexts() ([]*models.ContextWithMetadata, error)

	// UpdateContext updates an existing context's metadata and status.
	// Returns an error if the context does not exist.
	UpdateContext(ctx *models.ContextWithMetadata) error

	// DeleteContext permanently removes a context and all associated data.
	// Returns an error if the context does not exist.
	DeleteContext(name string) error

	// ArchiveContext marks a context as archived without deleting it.
	// Archived contexts are hidden from default listings.
	ArchiveContext(name string) error

	// AddNote adds a timestamped note to a context.
	// The timestamp should be in RFC3339 format.
	AddNote(contextName, timestamp, content string) error

	// GetNotes retrieves all notes for a context, ordered by timestamp (oldest first).
	// Returns an empty slice if no notes exist.
	GetNotes(contextName string) ([]Note, error)

	// GetNotesByTimestamp retrieves notes within a time range.
	GetNotesByTimestamp(contextName string, after, before time.Time) ([]Note, error)

	// AddFile associates a file path with a context.
	// The timestamp should be in RFC3339 format.
	AddFile(contextName, timestamp, path string) error

	// GetFiles retrieves all file associations for a context.
	// Returns an empty slice if no files exist.
	GetFiles(contextName string) ([]File, error)

	// GetFilesByTimestamp retrieves file associations within a time range.
	GetFilesByTimestamp(contextName string, after, before time.Time) ([]File, error)

	// AddTouch records an activity timestamp for a context.
	// Used to track when a context was last accessed.
	AddTouch(contextName, timestamp string) error

	// GetTouches retrieves all touch timestamps for a context.
	GetTouches(contextName string) ([]Touch, error)

	// LogTransition records a context state transition (start, stop, switch).
	LogTransition(transition *Transition) error

	// GetTransitions retrieves recent transitions across all contexts.
	// The limit parameter controls the maximum number of results.
	GetTransitions(limit int) ([]Transition, error)

	// GetTransitionsByContext retrieves transitions for a specific context.
	GetTransitionsByContext(contextName string) ([]Transition, error)

	// GetActiveContext returns the name of the currently active context.
	// Returns an empty string if no context is active.
	GetActiveContext() (string, error)

	// SetActiveContext sets the active context by name.
	SetActiveContext(contextName string) error

	// ClearActiveContext removes the active context marker.
	ClearActiveContext() error

	// GetContextsByLabel retrieves contexts with a specific label in their metadata.
	GetContextsByLabel(label string) ([]*models.ContextWithMetadata, error)

	// GetContextsByParent retrieves child contexts of a parent context.
	GetContextsByParent(parent string) ([]*models.ContextWithMetadata, error)

	// SearchContextsByName searches for contexts matching a query string.
	SearchContextsByName(query string) ([]*models.ContextWithMetadata, error)

	// SearchNotes searches note content across all contexts.
	// Returns a map of context name to matching notes.
	SearchNotes(query string) (map[string][]Note, error)

	// GetContextCount returns the total number of contexts.
	GetContextCount() (int, error)

	// GetContextStats retrieves statistics for a single context.
	GetContextStats(contextName string) (ContextStats, error)

	// GetGlobalStats retrieves aggregate statistics across all contexts.
	GetGlobalStats() (GlobalStats, error)
}

// Note represents a note entry
type Note struct {
	ID        int64
	ContextID int64
	Timestamp time.Time
	Content   string
}

// File represents a file tracking entry
type File struct {
	ID        int64
	ContextID int64
	Timestamp time.Time
	Path      string
}

// Touch represents a context touch (access timestamp)
type Touch struct {
	ID        int64
	ContextID int64
	Timestamp time.Time
}

// Transition represents a context state transition
type Transition struct {
	ID              int64
	Timestamp       time.Time
	PreviousContext *string
	NewContext      *string
	TransitionType  string // "start", "switch", "stop"
}

// ContextStats contains statistics for a single context
type ContextStats struct {
	ContextName string
	NotesCount  int
	FilesCount  int
	Duration    time.Duration
	LastUpdated time.Time
}

// GlobalStats contains statistics across all contexts
type GlobalStats struct {
	TotalContexts int
	TotalNotes    int
	TotalFiles    int
	ActiveContext *string
	OldestContext *models.ContextWithMetadata
	NewestContext *models.ContextWithMetadata
	ArchiveCount  int
}
