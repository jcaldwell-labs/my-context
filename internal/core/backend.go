package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// BackendType represents the type of storage backend
type BackendType string

const (
	BackendTypeFile     BackendType = "file"
	BackendTypePostgres BackendType = "postgres"
)

// DetectBackendType determines the backend type from MY_CONTEXT_HOME
func DetectBackendType() BackendType {
	home := os.Getenv("MY_CONTEXT_HOME")

	// Check for database connection strings
	if strings.HasPrefix(home, "postgresql://") ||
		strings.HasPrefix(home, "postgres://") ||
		home == "db" ||
		home == "database" ||
		home == "pg" ||
		strings.HasPrefix(home, "db:") ||
		strings.HasPrefix(home, "database:") ||
		strings.HasPrefix(home, "pg:") {
		return BackendTypePostgres
	}

	// Default to file-based
	return BackendTypeFile
}

// GetPostgresConnectionString returns the PostgreSQL connection string
// with search_path set to the appropriate schema for partitioning.
//
// Connection string resolution order:
// 1. MY_CONTEXT_HOME as full postgres:// URL (with search_path appended for partition)
// 2. DATABASE_URL environment variable (with search_path appended for partition)
// 3. Default localhost connection (for backward compatibility)
func GetPostgresConnectionString() (string, error) {
	home := os.Getenv("MY_CONTEXT_HOME")
	schema := GetPartitionSchema()

	// If full connection string in MY_CONTEXT_HOME, use it with search_path
	if strings.HasPrefix(home, "postgresql://") || strings.HasPrefix(home, "postgres://") {
		return appendSearchPath(home, schema), nil
	}

	// For shorthand syntax (db, db:partition), check DATABASE_URL first
	if home == "db" || home == "database" || home == "pg" ||
		strings.HasPrefix(home, "db:") ||
		strings.HasPrefix(home, "database:") ||
		strings.HasPrefix(home, "pg:") {

		// Try DATABASE_URL environment variable first
		if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
			return appendSearchPath(dbURL, schema), nil
		}

		// Fall back to default localhost connection for local development.
		// WARNING: These are development-only defaults. In production, always
		// set DATABASE_URL with proper credentials. These defaults exist for:
		// - Local Docker development (postgres-dev container)
		// - Quick testing without configuration
		// Do NOT use these credentials in production environments.
		connStr := fmt.Sprintf(
			"host=localhost port=5432 user=devuser password=devpassword dbname=dev_state sslmode=disable search_path=%s",
			schema,
		)

		return connStr, nil
	}

	return "", fmt.Errorf("MY_CONTEXT_HOME=%q is not recognized as a database backend. "+
		"Expected: 'db', 'db:<partition>', 'postgres://...', or a file path", home)
}

// appendSearchPath adds search_path parameter to a connection string
func appendSearchPath(connStr, schema string) string {
	// Handle URL-style connection strings (postgres://...)
	if strings.HasPrefix(connStr, "postgres://") || strings.HasPrefix(connStr, "postgresql://") {
		if strings.Contains(connStr, "?") {
			return connStr + "&search_path=" + schema
		}
		return connStr + "?search_path=" + schema
	}

	// Handle key=value style connection strings
	if strings.Contains(connStr, "search_path=") {
		// Already has search_path, don't add another
		return connStr
	}
	return connStr + " search_path=" + schema
}

// IsUsingDatabase returns true if MY_CONTEXT_HOME points to a database
func IsUsingDatabase() bool {
	return DetectBackendType() == BackendTypePostgres
}

// ExtractPartition extracts the partition name from MY_CONTEXT_HOME
// Examples:
//   - "db:adventure-engine" -> "adventure-engine"
//   - "database:my-project" -> "my-project"
//   - "db" -> "" (no partition, use default)
func ExtractPartition() string {
	home := os.Getenv("MY_CONTEXT_HOME")

	// Extract partition from db:partition syntax
	if strings.HasPrefix(home, "db:") {
		return strings.TrimPrefix(home, "db:")
	}
	if strings.HasPrefix(home, "database:") {
		return strings.TrimPrefix(home, "database:")
	}
	if strings.HasPrefix(home, "pg:") {
		return strings.TrimPrefix(home, "pg:")
	}

	// No partition specified
	return ""
}

// SanitizePartitionName converts a partition name to a valid PostgreSQL schema name
// Examples:
//   - "adventure-engine" -> "adventure_engine"
//   - "My Project!" -> "my_project"
//   - "payment-service-v2" -> "payment_service_v2"
func SanitizePartitionName(name string) string {
	if name == "" {
		return "public"
	}

	// Convert to lowercase
	sanitized := strings.ToLower(name)

	// Replace hyphens with underscores
	sanitized = strings.ReplaceAll(sanitized, "-", "_")

	// Replace spaces with underscores
	sanitized = strings.ReplaceAll(sanitized, " ", "_")

	// Remove any characters that aren't alphanumeric or underscore
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	sanitized = reg.ReplaceAllString(sanitized, "")

	// Ensure it doesn't start with a number (PostgreSQL requirement)
	if len(sanitized) > 0 && sanitized[0] >= '0' && sanitized[0] <= '9' {
		sanitized = "p_" + sanitized
	}

	// If empty after sanitization, use default
	if sanitized == "" {
		return "public"
	}

	return sanitized
}

// GetPartitionSchema returns the PostgreSQL schema name for the current partition
func GetPartitionSchema() string {
	partition := ExtractPartition()
	return SanitizePartitionName(partition)
}
