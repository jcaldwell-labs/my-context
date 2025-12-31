package unit

import (
	"os"
	"testing"

	"github.com/jefferycaldwell/my-context-copilot/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestExtractPartition(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expected      string
		description   string
	}{
		{
			name:        "db:partition syntax",
			envValue:    "db:adventure-engine",
			expected:    "adventure-engine",
			description: "Should extract partition name from db: prefix",
		},
		{
			name:        "database:partition syntax",
			envValue:    "database:my-project",
			expected:    "my-project",
			description: "Should extract partition name from database: prefix",
		},
		{
			name:        "pg:partition syntax",
			envValue:    "pg:payment-service",
			expected:    "payment-service",
			description: "Should extract partition name from pg: prefix",
		},
		{
			name:        "db without partition",
			envValue:    "db",
			expected:    "",
			description: "Should return empty string for db without partition",
		},
		{
			name:        "database without partition",
			envValue:    "database",
			expected:    "",
			description: "Should return empty string for database without partition",
		},
		{
			name:        "file path",
			envValue:    "/home/user/.my-context",
			expected:    "",
			description: "Should return empty string for file paths",
		},
		{
			name:        "empty",
			envValue:    "",
			expected:    "",
			description: "Should return empty string when not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("MY_CONTEXT_HOME", tt.envValue)
			defer os.Unsetenv("MY_CONTEXT_HOME")

			// Test extraction
			result := core.ExtractPartition()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestSanitizePartitionName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name:        "hyphens to underscores",
			input:       "adventure-engine",
			expected:    "adventure_engine",
			description: "Should convert hyphens to underscores",
		},
		{
			name:        "spaces to underscores",
			input:       "my project",
			expected:    "my_project",
			description: "Should convert spaces to underscores",
		},
		{
			name:        "uppercase to lowercase",
			input:       "MyProject",
			expected:    "myproject",
			description: "Should convert uppercase to lowercase",
		},
		{
			name:        "special characters removed",
			input:       "project!@#$%",
			expected:    "project",
			description: "Should remove special characters",
		},
		{
			name:        "mixed transformations",
			input:       "Payment Service v2.1!",
			expected:    "payment_service_v21",
			description: "Should apply all transformations",
		},
		{
			name:        "starts with number",
			input:       "123project",
			expected:    "p_123project",
			description: "Should prefix with p_ if starts with number",
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "public",
			description: "Should return 'public' for empty string",
		},
		{
			name:        "only special characters",
			input:       "!@#$%",
			expected:    "public",
			description: "Should return 'public' if nothing remains after sanitization",
		},
		{
			name:        "valid schema name",
			input:       "valid_schema_123",
			expected:    "valid_schema_123",
			description: "Should keep valid schema names unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.SanitizePartitionName(tt.input)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestGetPartitionSchema(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expected    string
		description string
	}{
		{
			name:        "adventure-engine partition",
			envValue:    "db:adventure-engine",
			expected:    "adventure_engine",
			description: "Should return sanitized schema name for partition",
		},
		{
			name:        "My Project partition",
			envValue:    "db:My Project!",
			expected:    "my_project",
			description: "Should sanitize and return schema name",
		},
		{
			name:        "default (no partition)",
			envValue:    "db",
			expected:    "public",
			description: "Should return 'public' for default database",
		},
		{
			name:        "file-based mode",
			envValue:    "/home/user/.my-context",
			expected:    "public",
			description: "Should return 'public' for file-based mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("MY_CONTEXT_HOME", tt.envValue)
			defer os.Unsetenv("MY_CONTEXT_HOME")

			// Test get schema
			result := core.GetPartitionSchema()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestDetectBackendType(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expected    core.BackendType
		description string
	}{
		{
			name:        "db shorthand",
			envValue:    "db",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for 'db'",
		},
		{
			name:        "db:partition syntax",
			envValue:    "db:my-partition",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for 'db:partition'",
		},
		{
			name:        "database shorthand",
			envValue:    "database",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for 'database'",
		},
		{
			name:        "database:partition syntax",
			envValue:    "database:my-partition",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for 'database:partition'",
		},
		{
			name:        "pg shorthand",
			envValue:    "pg",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for 'pg'",
		},
		{
			name:        "pg:partition syntax",
			envValue:    "pg:my-partition",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for 'pg:partition'",
		},
		{
			name:        "postgresql:// URL",
			envValue:    "postgresql://user:pass@localhost/mydb",
			expected:    core.BackendTypePostgres,
			description: "Should detect PostgreSQL for full URL",
		},
		{
			name:        "file path",
			envValue:    "/home/user/.my-context",
			expected:    core.BackendTypeFile,
			description: "Should detect file backend for paths",
		},
		{
			name:        "empty (default)",
			envValue:    "",
			expected:    core.BackendTypeFile,
			description: "Should default to file backend when not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("MY_CONTEXT_HOME", tt.envValue)
				defer os.Unsetenv("MY_CONTEXT_HOME")
			} else {
				os.Unsetenv("MY_CONTEXT_HOME")
			}

			// Test backend detection
			result := core.DetectBackendType()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestGetPostgresConnectionString(t *testing.T) {
	tests := []struct {
		name             string
		envValue         string
		expectedContains []string
		shouldError      bool
		description      string
	}{
		{
			name:     "db default",
			envValue: "db",
			expectedContains: []string{
				"host=localhost",
				"port=5432",
				"dbname=dev_state",
				"search_path=public",
			},
			shouldError: false,
			description: "Should return default connection with public schema",
		},
		{
			name:     "db:partition",
			envValue: "db:adventure-engine",
			expectedContains: []string{
				"host=localhost",
				"search_path=adventure_engine",
			},
			shouldError: false,
			description: "Should set search_path to sanitized partition name",
		},
		{
			name:     "database:partition with special chars",
			envValue: "database:My Project!",
			expectedContains: []string{
				"search_path=my_project",
			},
			shouldError: false,
			description: "Should sanitize partition name in search_path",
		},
		{
			name:        "file path",
			envValue:    "/home/user/.my-context",
			shouldError: true,
			description: "Should error for non-database backend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("MY_CONTEXT_HOME", tt.envValue)
			defer os.Unsetenv("MY_CONTEXT_HOME")

			// Test connection string generation
			connStr, err := core.GetPostgresConnectionString()

			if tt.shouldError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
				for _, expected := range tt.expectedContains {
					assert.Contains(t, connStr, expected, "Connection string should contain: "+expected)
				}
			}
		})
	}
}
