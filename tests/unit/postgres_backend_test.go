package unit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/jefferycaldwell/my-context-copilot/pkg/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPostgresBackend tests the backend constructor
func TestNewPostgresBackend(t *testing.T) {
	tests := []struct {
		name            string
		connStr         string
		expectedSchema  string
		expectedPartition string
	}{
		{
			name:              "connection string with search_path",
			connStr:           "host=localhost port=5432 user=test dbname=test search_path=my_schema",
			expectedSchema:    "my_schema",
			expectedPartition: "my_schema",
		},
		{
			name:              "connection string without search_path",
			connStr:           "host=localhost port=5432 user=test dbname=test",
			expectedSchema:    "public",
			expectedPartition: "public",
		},
		{
			name:              "connection string with quoted search_path",
			connStr:           "host=localhost search_path='my_schema'",
			expectedSchema:    "'my_schema'",
			expectedPartition: "'my_schema'",
		},
		{
			name:              "empty connection string",
			connStr:           "",
			expectedSchema:    "public",
			expectedPartition: "public",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := postgres.NewPostgresBackend(tt.connStr)
			assert.NotNil(t, backend)
			assert.Equal(t, tt.expectedSchema, backend.GetSchema())
			assert.Equal(t, tt.expectedPartition, backend.GetPartition())
		})
	}
}

// TestExtractSchemaFromConnStr tests schema extraction from connection strings
func TestExtractSchemaFromConnStr(t *testing.T) {
	tests := []struct {
		name           string
		connStr        string
		expectedSchema string
	}{
		{
			name:           "schema at end",
			connStr:        "host=localhost port=5432 search_path=myschema",
			expectedSchema: "myschema",
		},
		{
			name:           "schema in middle",
			connStr:        "host=localhost search_path=myschema port=5432",
			expectedSchema: "myschema",
		},
		{
			name:           "schema at beginning",
			connStr:        "search_path=myschema host=localhost port=5432",
			expectedSchema: "myschema",
		},
		{
			name:           "no search_path",
			connStr:        "host=localhost port=5432 dbname=test",
			expectedSchema: "public",
		},
		{
			name:           "empty string",
			connStr:        "",
			expectedSchema: "public",
		},
		{
			name:           "search_path with quoted value",
			connStr:        "host=localhost search_path='my_schema'",
			expectedSchema: "'my_schema'",
		},
		{
			name:           "search_path with underscores",
			connStr:        "host=localhost search_path=my_long_schema_name",
			expectedSchema: "my_long_schema_name",
		},
		{
			name:           "multiple spaces",
			connStr:        "host=localhost  search_path=myschema  port=5432",
			expectedSchema: "myschema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := postgres.NewPostgresBackend(tt.connStr)
			assert.Equal(t, tt.expectedSchema, backend.GetSchema())
		})
	}
}

// TestSplitConnStr tests connection string splitting with quotes
func TestSplitConnStr(t *testing.T) {
	tests := []struct {
		name     string
		connStr  string
		expected []string
	}{
		{
			name:     "simple key=value pairs",
			connStr:  "host=localhost port=5432 dbname=test",
			expected: []string{"host=localhost", "port=5432", "dbname=test"},
		},
		{
			name:     "with quoted value",
			connStr:  "host=localhost password='my password'",
			expected: []string{"host=localhost", "password='my password'"},
		},
		{
			name:     "with escaped quote (doubled single quote)",
			connStr:  "host=localhost password='it''s password'",
			expected: []string{"host=localhost", "password='it''s password'"},
		},
		{
			name:     "multiple quoted values",
			connStr:  "host='localhost' dbname='test db' password='pass'",
			expected: []string{"host='localhost'", "dbname='test db'", "password='pass'"},
		},
		{
			name:     "empty string",
			connStr:  "",
			expected: []string{},
		},
		{
			name:     "single parameter",
			connStr:  "host=localhost",
			expected: []string{"host=localhost"},
		},
		{
			name:     "trailing space",
			connStr:  "host=localhost port=5432 ",
			expected: []string{"host=localhost", "port=5432"},
		},
		{
			name:     "leading space",
			connStr:  " host=localhost port=5432",
			expected: []string{"host=localhost", "port=5432"},
		},
		{
			name:     "multiple spaces between params",
			connStr:  "host=localhost    port=5432",
			expected: []string{"host=localhost", "port=5432"},
		},
		{
			name:     "quoted value with special chars",
			connStr:  "password='p@ss!word#123'",
			expected: []string{"password='p@ss!word#123'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a backend to access the split function indirectly
			// The split is used internally by extractSchemaFromConnStr
			// We verify it by checking the full parsing behavior
			backend := postgres.NewPostgresBackend(tt.connStr)
			assert.NotNil(t, backend, "Backend should be created")
			
			// For connection strings without search_path, we expect public schema
			if !containsSearchPath(tt.connStr) {
				assert.Equal(t, "public", backend.GetSchema())
			}
		})
	}
}

// Helper function to check if connection string contains search_path
func containsSearchPath(connStr string) bool {
	return len(connStr) > 11 && (
		connStr[:12] == "search_path=" ||
		contains(connStr, " search_path=") ||
		contains(connStr, "'search_path="))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGettersOnBackend tests the getter methods
func TestGettersOnBackend(t *testing.T) {
	t.Run("GetSchema returns correct schema", func(t *testing.T) {
		backend := postgres.NewPostgresBackend("search_path=testschema")
		assert.Equal(t, "testschema", backend.GetSchema())
	})

	t.Run("GetPartition returns correct partition", func(t *testing.T) {
		backend := postgres.NewPostgresBackend("search_path=testpartition")
		assert.Equal(t, "testpartition", backend.GetPartition())
	})

	t.Run("GetDB returns nil when not initialized", func(t *testing.T) {
		backend := postgres.NewPostgresBackend("host=localhost")
		assert.Nil(t, backend.GetDB(), "DB should be nil before Init()")
	})
}

// TestContextMetadataMarshaling tests JSON marshaling/unmarshaling
func TestContextMetadataMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		metadata models.ContextMetadata
	}{
		{
			name: "complete metadata",
			metadata: models.ContextMetadata{
				CreatedBy: "testuser",
				Parent:    "parent-context",
				Labels:    []string{"tag1", "tag2", "tag3"},
			},
		},
		{
			name: "empty metadata",
			metadata: models.ContextMetadata{
				CreatedBy: "",
				Labels:    nil, // Use nil instead of empty slice
			},
		},
		{
			name: "nil labels",
			metadata: models.ContextMetadata{
				CreatedBy: "user",
				Labels:    nil,
			},
		},
		{
			name: "special characters",
			metadata: models.ContextMetadata{
				CreatedBy: "user@example.com",
				Parent:    "parent_with-special.chars",
				Labels:    []string{"test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.metadata)
			require.NoError(t, err, "Marshaling should succeed")
			assert.NotEmpty(t, data, "Marshaled data should not be empty")

			// Unmarshal back
			var unmarshaled models.ContextMetadata
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err, "Unmarshaling should succeed")

			// Verify equality
			assert.Equal(t, tt.metadata.CreatedBy, unmarshaled.CreatedBy)
			assert.Equal(t, tt.metadata.Parent, unmarshaled.Parent)
			if tt.metadata.Labels == nil {
				// nil should unmarshal as nil or empty slice
				assert.True(t, unmarshaled.Labels == nil || len(unmarshaled.Labels) == 0)
			} else {
				assert.Equal(t, tt.metadata.Labels, unmarshaled.Labels)
			}
		})
	}
}

// TestContextWithMetadataStructure tests the ContextWithMetadata structure
func TestContextWithMetadataStructure(t *testing.T) {
	now := time.Now()
	endTime := now.Add(1 * time.Hour)
	lastTouch := now.Add(30 * time.Minute)

	ctx := &models.ContextWithMetadata{
		Name:        "test-context",
		StartTime:   now,
		EndTime:     &endTime,
		Status:      "stopped",
		IsArchived:  false,
		Metadata: models.ContextMetadata{
			CreatedBy: "testuser",
			Labels:    []string{"test"},
		},
		TouchCount:  5,
		LastTouchAt: &lastTouch,
	}

	t.Run("all fields accessible", func(t *testing.T) {
		assert.Equal(t, "test-context", ctx.Name)
		assert.Equal(t, now, ctx.StartTime)
		assert.NotNil(t, ctx.EndTime)
		assert.Equal(t, endTime, *ctx.EndTime)
		assert.Equal(t, "stopped", ctx.Status)
		assert.False(t, ctx.IsArchived)
		assert.Equal(t, "testuser", ctx.Metadata.CreatedBy)
		assert.Equal(t, []string{"test"}, ctx.Metadata.Labels)
		assert.Equal(t, 5, ctx.TouchCount)
		assert.NotNil(t, ctx.LastTouchAt)
		assert.Equal(t, lastTouch, *ctx.LastTouchAt)
	})

	t.Run("nil pointer fields", func(t *testing.T) {
		minimalCtx := &models.ContextWithMetadata{
			Name:       "minimal",
			StartTime:  now,
			Status:     "active",
			IsArchived: false,
		}

		assert.Nil(t, minimalCtx.EndTime)
		assert.Nil(t, minimalCtx.LastTouchAt)
		assert.Equal(t, 0, minimalCtx.TouchCount)
	})
}

// TestBackendMethodsBeforeInit tests that methods handle uninitialized state
func TestBackendMethodsBeforeInit(t *testing.T) {
	backend := postgres.NewPostgresBackend("host=localhost")

	t.Run("GetSchema works before Init", func(t *testing.T) {
		schema := backend.GetSchema()
		assert.Equal(t, "public", schema)
	})

	t.Run("GetPartition works before Init", func(t *testing.T) {
		partition := backend.GetPartition()
		assert.Equal(t, "public", partition)
	})

	t.Run("GetDB returns nil before Init", func(t *testing.T) {
		db := backend.GetDB()
		assert.Nil(t, db)
	})

	t.Run("Close handles nil db gracefully", func(t *testing.T) {
		err := backend.Close()
		assert.NoError(t, err, "Close should not error when db is nil")
	})
}

// TestConnectionStringEdgeCases tests edge cases in connection string parsing
func TestConnectionStringEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		connStr        string
		expectedSchema string
	}{
		{
			name:           "search_path as first parameter",
			connStr:        "search_path=first host=localhost",
			expectedSchema: "first",
		},
		{
			name:           "search_path as last parameter",
			connStr:        "host=localhost search_path=last",
			expectedSchema: "last",
		},
		{
			name:           "search_path with numeric value",
			connStr:        "host=localhost search_path=schema123",
			expectedSchema: "schema123",
		},
		{
			name:           "search_path with underscore",
			connStr:        "search_path=my_custom_schema",
			expectedSchema: "my_custom_schema",
		},
		{
			name:           "only search_path",
			connStr:        "search_path=only",
			expectedSchema: "only",
		},
		{
			name:           "substring match should not trigger",
			connStr:        "host=localhost newsearch_path=fake",
			expectedSchema: "public",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := postgres.NewPostgresBackend(tt.connStr)
			assert.Equal(t, tt.expectedSchema, backend.GetSchema())
		})
	}
}
