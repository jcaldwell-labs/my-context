package storage

import (
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
)

// Backend defines the interface for context storage implementations
type Backend interface {
	// Initialization and validation
	Init() error
	Health() error
	Close() error

	// Context operations
	CreateContext(ctx *models.ContextWithMetadata) error
	GetContext(name string) (*models.ContextWithMetadata, error)
	ListContexts() ([]*models.ContextWithMetadata, error)
	UpdateContext(ctx *models.ContextWithMetadata) error
	DeleteContext(name string) error
	ArchiveContext(name string) error

	// Notes operations
	AddNote(contextName, timestamp, content string) error
	GetNotes(contextName string) ([]Note, error)
	GetNotesByTimestamp(contextName string, after, before time.Time) ([]Note, error)

	// Files operations
	AddFile(contextName, timestamp, path string) error
	GetFiles(contextName string) ([]File, error)
	GetFilesByTimestamp(contextName string, after, before time.Time) ([]File, error)

	// Touch operations
	AddTouch(contextName, timestamp string) error
	GetTouches(contextName string) ([]Touch, error)

	// Transitions operations
	LogTransition(transition *Transition) error
	GetTransitions(limit int) ([]Transition, error)
	GetTransitionsByContext(contextName string) ([]Transition, error)

	// State operations
	GetActiveContext() (string, error)
	SetActiveContext(contextName string) error
	ClearActiveContext() error

	// Search and query operations
	GetContextsByLabel(label string) ([]*models.ContextWithMetadata, error)
	GetContextsByParent(parent string) ([]*models.ContextWithMetadata, error)
	SearchContextsByName(query string) ([]*models.ContextWithMetadata, error)
	SearchNotes(query string) (map[string][]Note, error)

	// Statistics and reporting
	GetContextCount() (int, error)
	GetContextStats(contextName string) (ContextStats, error)
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
	TotalContexts  int
	TotalNotes     int
	TotalFiles     int
	ActiveContext  *string
	OldestContext  *models.ContextWithMetadata
	NewestContext  *models.ContextWithMetadata
	ArchiveCount   int
}
