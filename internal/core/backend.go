package core

import (
	"fmt"
	"net/url"
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

	// For shorthand syntax (db, db:partition), require DATABASE_URL
	if home == "db" || home == "database" || home == "pg" ||
		strings.HasPrefix(home, "db:") ||
		strings.HasPrefix(home, "database:") ||
		strings.HasPrefix(home, "pg:") {

		// Use DATABASE_URL environment variable
		if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
			return appendSearchPath(dbURL, schema), nil
		}

		// No DATABASE_URL provided: fail fast instead of using hardcoded credentials.
		// This avoids accidentally connecting with insecure development defaults
		// in production environments.
		return "", fmt.Errorf("DATABASE_URL environment variable must be set when using MY_CONTEXT_HOME=%q. "+
			"Example: export DATABASE_URL='host=localhost port=5432 user=myuser dbname=mydb sslmode=disable'", home)
	}

	return "", fmt.Errorf("MY_CONTEXT_HOME=%q is not recognized as a database backend. "+
		"Use 'db', 'db:<partition>', 'postgres://...', or a directory path for file storage", home)
}

// appendSearchPath adds search_path parameter to a connection string
func appendSearchPath(connStr, schema string) string {
	// Handle URL-style connection strings (postgres://...)
	if strings.HasPrefix(connStr, "postgres://") || strings.HasPrefix(connStr, "postgresql://") {
		// URL-encode the schema name for URL-style connection strings
		encodedSchema := url.QueryEscape(schema)
		if strings.Contains(connStr, "?") {
			return connStr + "&search_path=" + encodedSchema
		}
		return connStr + "?search_path=" + encodedSchema
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
