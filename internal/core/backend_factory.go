package core

import (
	"fmt"

	"github.com/jefferycaldwell/my-context-copilot/pkg/storage"
	"github.com/jefferycaldwell/my-context-copilot/pkg/storage/postgres"
)

// GetBackend returns the appropriate storage backend based on MY_CONTEXT_HOME
func GetBackend() (storage.Backend, error) {
	backendType := DetectBackendType()

	switch backendType {
	case BackendTypePostgres:
		// Get PostgreSQL connection string
		connStr, err := GetPostgresConnectionString()
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres connection string: %w", err)
		}

		// Create and initialize PostgreSQL backend
		backend := postgres.NewPostgresBackend(connStr)
		if err := backend.Init(); err != nil {
			return nil, fmt.Errorf("failed to initialize postgres backend: %w", err)
		}

		return backend, nil

	case BackendTypeFile:
		// File-based backend (existing implementation)
		// For now, return an error since file backend isn't fully implemented yet
		return nil, fmt.Errorf("file backend not yet implemented in new storage interface")

	default:
		return nil, fmt.Errorf("unknown backend type: %s", backendType)
	}
}
