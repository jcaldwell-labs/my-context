package core

import (
	"testing"
)

func TestDetectBackendType(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     BackendType
	}{
		// Database backends
		{"shorthand db", "db", BackendTypePostgres},
		{"shorthand database", "database", BackendTypePostgres},
		{"shorthand pg", "pg", BackendTypePostgres},
		{"partition syntax db:name", "db:adventure-engine", BackendTypePostgres},
		{"partition syntax database:name", "database:my-project", BackendTypePostgres},
		{"partition syntax pg:name", "pg:scrum", BackendTypePostgres},
		{"full postgres URL", "postgres://user:pass@localhost/db", BackendTypePostgres},
		{"full postgresql URL", "postgresql://user:pass@localhost/db", BackendTypePostgres},

		// File backends
		{"empty string", "", BackendTypeFile},
		{"relative path", "./contexts", BackendTypeFile},
		{"absolute path", "/home/user/.my-context", BackendTypeFile},
		{"home expansion", "~/.my-context", BackendTypeFile},
		{"random string", "some-random-value", BackendTypeFile},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MY_CONTEXT_HOME", tt.envValue)
			got := DetectBackendType()
			if got != tt.want {
				t.Errorf("DetectBackendType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPartition(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{"no partition db", "db", ""},
		{"no partition database", "database", ""},
		{"no partition pg", "pg", ""},
		{"partition with db:", "db:adventure-engine", "adventure-engine"},
		{"partition with database:", "database:my-project", "my-project"},
		{"partition with pg:", "pg:scrum", "scrum"},
		{"partition with special chars", "db:my-awesome-project!", "my-awesome-project!"},
		{"full URL no partition", "postgres://localhost/db", ""},
		{"empty", "", ""},
		{"file path", "/home/user/.my-context", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MY_CONTEXT_HOME", tt.envValue)
			got := ExtractPartition()
			if got != tt.want {
				t.Errorf("ExtractPartition() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSanitizePartitionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty returns public", "", "public"},
		{"lowercase passthrough", "adventure", "adventure"},
		{"hyphens to underscores", "adventure-engine", "adventure_engine"},
		{"spaces to underscores", "my project", "my_project"},
		{"uppercase to lowercase", "MyProject", "myproject"},
		{"special chars removed", "my-project!", "my_project"},
		{"complex name", "My Project! v2.0", "my_project_v20"},
		{"starts with number", "123test", "p_123test"},
		{"only special chars", "!@#$%", "public"},
		{"mixed", "Payment-Service-v2", "payment_service_v2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizePartitionName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizePartitionName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetPartitionSchema(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{"no partition uses public", "db", "public"},
		{"partition is sanitized", "db:adventure-engine", "adventure_engine"},
		{"complex partition", "db:My Project!", "my_project"},
		{"file backend uses public", "/home/user/.my-context", "public"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MY_CONTEXT_HOME", tt.envValue)
			got := GetPartitionSchema()
			if got != tt.want {
				t.Errorf("GetPartitionSchema() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsUsingDatabase(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     bool
	}{
		{"db shorthand", "db", true},
		{"full postgres URL", "postgres://localhost/db", true},
		{"file path", "/home/user/.my-context", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MY_CONTEXT_HOME", tt.envValue)
			got := IsUsingDatabase()
			if got != tt.want {
				t.Errorf("IsUsingDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendSearchPath(t *testing.T) {
	tests := []struct {
		name    string
		connStr string
		schema  string
		want    string
	}{
		{
			name:    "URL without params",
			connStr: "postgres://localhost/mydb",
			schema:  "myschema",
			want:    "postgres://localhost/mydb?search_path=myschema",
		},
		{
			name:    "URL with existing params",
			connStr: "postgres://localhost/mydb?sslmode=disable",
			schema:  "myschema",
			want:    "postgres://localhost/mydb?sslmode=disable&search_path=myschema",
		},
		{
			name:    "key=value style",
			connStr: "host=localhost dbname=mydb",
			schema:  "myschema",
			want:    "host=localhost dbname=mydb search_path=myschema",
		},
		{
			name:    "already has search_path key=value",
			connStr: "host=localhost dbname=mydb search_path=existing",
			schema:  "newschema",
			want:    "host=localhost dbname=mydb search_path=existing", // unchanged
		},
		{
			name:    "schema with special chars gets encoded",
			connStr: "postgres://localhost/mydb",
			schema:  "my schema",
			want:    "postgres://localhost/mydb?search_path=my+schema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendSearchPath(tt.connStr, tt.schema)
			if got != tt.want {
				t.Errorf("appendSearchPath(%q, %q) = %q, want %q", tt.connStr, tt.schema, got, tt.want)
			}
		})
	}
}

func TestGetPostgresConnectionString(t *testing.T) {
	tests := []struct {
		name        string
		home        string
		dbURL       string
		wantErr     bool
		wantContain string
	}{
		{
			name:        "full URL in MY_CONTEXT_HOME",
			home:        "postgres://user:pass@localhost/mydb",
			dbURL:       "",
			wantErr:     false,
			wantContain: "postgres://user:pass@localhost/mydb?search_path=public",
		},
		{
			name:        "shorthand with DATABASE_URL",
			home:        "db",
			dbURL:       "postgres://user:pass@localhost/mydb",
			wantErr:     false,
			wantContain: "search_path=public",
		},
		{
			name:        "shorthand without DATABASE_URL fails",
			home:        "db",
			dbURL:       "",
			wantErr:     true,
			wantContain: "",
		},
		{
			name:        "partition with DATABASE_URL",
			home:        "db:mypartition",
			dbURL:       "postgres://localhost/mydb",
			wantErr:     false,
			wantContain: "search_path=mypartition",
		},
		{
			name:        "invalid backend type fails",
			home:        "/some/path",
			dbURL:       "",
			wantErr:     true,
			wantContain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MY_CONTEXT_HOME", tt.home)
			if tt.dbURL != "" {
				t.Setenv("DATABASE_URL", tt.dbURL)
			} else {
				// Ensure DATABASE_URL is unset
				t.Setenv("DATABASE_URL", "")
			}

			got, err := GetPostgresConnectionString()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostgresConnectionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.wantContain != "" {
				if got != tt.wantContain && !contains(got, tt.wantContain) {
					t.Errorf("GetPostgresConnectionString() = %q, want containing %q", got, tt.wantContain)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
