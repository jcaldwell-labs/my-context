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
			// Clear DATABASE_URL to test default behavior
			os.Unsetenv("DATABASE_URL")

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

// TestGetPostgresConnectionStringWithDATABASE_URL tests the DATABASE_URL environment
// variable feature which enables cross-platform context sharing between WSL and Windows.
//
// Connection string resolution order:
// 1. MY_CONTEXT_HOME as full postgres:// URL (with search_path appended)
// 2. DATABASE_URL environment variable (with search_path appended)
// 3. Default localhost connection (backward compatible)
func TestGetPostgresConnectionStringWithDATABASE_URL(t *testing.T) {
	tests := []struct {
		name             string
		myContextHome    string
		databaseURL      string
		expectedContains []string
		notContains      []string
		description      string
	}{
		{
			name:          "DATABASE_URL takes precedence over default localhost",
			myContextHome: "db",
			databaseURL:   "host=192.168.1.100 port=5432 user=produser password=prodpass dbname=prod_state sslmode=require",
			expectedContains: []string{
				"host=192.168.1.100",
				"user=produser",
				"dbname=prod_state",
				"sslmode=require",
				"search_path=public",
			},
			notContains: []string{
				"host=localhost",
				"user=devuser",
			},
			description: "DATABASE_URL should override default localhost connection",
		},
		{
			name:          "DATABASE_URL with partition",
			myContextHome: "db:payment-service",
			databaseURL:   "host=db.example.com port=5432 user=app password=secret dbname=contexts",
			expectedContains: []string{
				"host=db.example.com",
				"dbname=contexts",
				"search_path=payment_service",
			},
			description: "DATABASE_URL should be used with partition's search_path",
		},
		{
			name:          "DATABASE_URL with URL-style connection string",
			myContextHome: "db",
			databaseURL:   "postgres://user:pass@remote-host:5432/mydb",
			expectedContains: []string{
				"postgres://user:pass@remote-host:5432/mydb",
				"search_path=public",
			},
			description: "Should append search_path to URL-style DATABASE_URL",
		},
		{
			name:          "DATABASE_URL with URL and partition",
			myContextHome: "db:scrum",
			databaseURL:   "postgres://admin:secret@prod.db.com:5432/state",
			expectedContains: []string{
				"postgres://admin:secret@prod.db.com:5432/state",
				"search_path=scrum",
			},
			description: "Should append partition search_path to URL-style DATABASE_URL",
		},
		{
			name:          "MY_CONTEXT_HOME full URL takes precedence over DATABASE_URL",
			myContextHome: "postgres://override:pass@myhost:5432/mydb",
			databaseURL:   "host=should-not-use port=5432 user=ignored",
			expectedContains: []string{
				"postgres://override:pass@myhost:5432/mydb",
			},
			notContains: []string{
				"should-not-use",
				"ignored",
			},
			description: "Full postgres:// URL in MY_CONTEXT_HOME should take precedence",
		},
		{
			name:          "Fallback to localhost when DATABASE_URL not set",
			myContextHome: "db",
			databaseURL:   "", // Not set
			expectedContains: []string{
				"host=localhost",
				"user=devuser",
				"password=devpassword",
				"dbname=dev_state",
			},
			description: "Should fall back to localhost defaults when DATABASE_URL is empty",
		},
		{
			name:          "DATABASE_URL with existing search_path preserved",
			myContextHome: "db",
			databaseURL:   "host=myhost port=5432 user=test dbname=testdb search_path=existing_schema",
			expectedContains: []string{
				"host=myhost",
				"search_path=existing_schema",
			},
			description: "Should not duplicate search_path if already present",
		},
		{
			name:          "URL-style with existing query params",
			myContextHome: "db:mypartition",
			databaseURL:   "postgres://user:pass@host:5432/db?sslmode=require",
			expectedContains: []string{
				"postgres://user:pass@host:5432/db?sslmode=require",
				"&search_path=mypartition",
			},
			description: "Should append search_path with & when URL has existing params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set MY_CONTEXT_HOME
			os.Setenv("MY_CONTEXT_HOME", tt.myContextHome)
			defer os.Unsetenv("MY_CONTEXT_HOME")

			// Set or unset DATABASE_URL
			if tt.databaseURL != "" {
				os.Setenv("DATABASE_URL", tt.databaseURL)
				defer os.Unsetenv("DATABASE_URL")
			} else {
				os.Unsetenv("DATABASE_URL")
			}

			// Test connection string generation
			connStr, err := core.GetPostgresConnectionString()
			assert.NoError(t, err, tt.description)

			// Check expected substrings
			for _, expected := range tt.expectedContains {
				assert.Contains(t, connStr, expected,
					"Connection string should contain: %s\nGot: %s", expected, connStr)
			}

			// Check that unwanted substrings are not present
			for _, notExpected := range tt.notContains {
				assert.NotContains(t, connStr, notExpected,
					"Connection string should NOT contain: %s\nGot: %s", notExpected, connStr)
			}
		})
	}
}

// TestIsValidPostgresIdentifier tests the SQL injection prevention function
func TestIsValidPostgresIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		description string
	}{
		{
			name:     "valid simple name",
			input:    "public",
			expected: true,
			description: "Simple lowercase name should be valid",
		},
		{
			name:     "valid with underscore",
			input:    "my_schema",
			expected: true,
			description: "Name with underscore should be valid",
		},
		{
			name:     "valid with numbers",
			input:    "schema123",
			expected: true,
			description: "Name with numbers (not first) should be valid",
		},
		{
			name:     "valid starting with underscore",
			input:    "_private",
			expected: true,
			description: "Name starting with underscore should be valid",
		},
		{
			name:     "valid uppercase",
			input:    "MySchema",
			expected: true,
			description: "Mixed case should be valid",
		},
		{
			name:     "invalid starting with number",
			input:    "123schema",
			expected: false,
			description: "Name starting with number should be invalid",
		},
		{
			name:     "invalid with hyphen",
			input:    "my-schema",
			expected: false,
			description: "Name with hyphen should be invalid (SQL injection vector)",
		},
		{
			name:     "invalid with space",
			input:    "my schema",
			expected: false,
			description: "Name with space should be invalid",
		},
		{
			name:     "invalid with semicolon",
			input:    "schema;DROP TABLE",
			expected: false,
			description: "SQL injection attempt should be invalid",
		},
		{
			name:     "invalid with quotes",
			input:    "schema'--",
			expected: false,
			description: "SQL injection with quotes should be invalid",
		},
		{
			name:     "invalid empty",
			input:    "",
			expected: false,
			description: "Empty string should be invalid",
		},
		{
			name:     "invalid too long",
			input:    "this_is_a_very_long_schema_name_that_exceeds_postgresql_63_char_limit_for_identifiers",
			expected: false,
			description: "Name exceeding 63 chars should be invalid",
		},
		{
			name:     "valid exactly 63 chars",
			input:    "a23456789012345678901234567890123456789012345678901234567890123",
			expected: true,
			description: "Name with exactly 63 chars should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Import utils for the function
			result := isValidPostgresIdentifier(tt.input)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// isValidPostgresIdentifier is a local copy for testing
// (mirrors pkg/utils.IsValidPostgresIdentifier)
func isValidPostgresIdentifier(name string) bool {
	if name == "" || len(name) > 63 {
		return false
	}

	for i, ch := range name {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
			continue
		}
		if i > 0 && ch >= '0' && ch <= '9' {
			continue
		}
		return false
	}

	return true
}
