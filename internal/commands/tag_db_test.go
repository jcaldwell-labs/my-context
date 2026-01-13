package commands

import (
	"fmt"
	"testing"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/models"
	"github.com/jefferycaldwell/my-context-copilot/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockBackend is a test implementation of storage.Backend for tag tests
type mockBackend struct {
	contexts           map[string]*models.ContextWithMetadata
	getContextError    error
	updateContextError error
	listContextsError  error
}

func newMockBackend() *mockBackend {
	return &mockBackend{
		contexts: make(map[string]*models.ContextWithMetadata),
	}
}

func (m *mockBackend) addContext(name string, labels []string) {
	ctx := &models.ContextWithMetadata{
		Name:      name,
		StartTime: time.Now(),
		Status:    "active",
		Metadata: models.ContextMetadata{
			Labels: labels,
		},
	}
	m.contexts[name] = ctx
}

// Implement storage.Backend interface (only methods needed for testing)
func (m *mockBackend) Init() error                                          { return nil }
func (m *mockBackend) Health() error                                        { return nil }
func (m *mockBackend) Close() error                                         { return nil }
func (m *mockBackend) CreateContext(ctx *models.ContextWithMetadata) error  { return nil }
func (m *mockBackend) DeleteContext(name string) error                      { return nil }
func (m *mockBackend) ArchiveContext(name string) error                     { return nil }
func (m *mockBackend) AddNote(contextName, timestamp, content string) error { return nil }
func (m *mockBackend) GetNotes(contextName string) ([]storage.Note, error)  { return nil, nil }
func (m *mockBackend) GetNotesByTimestamp(contextName string, after, before time.Time) ([]storage.Note, error) {
	return nil, nil
}
func (m *mockBackend) AddFile(contextName, timestamp, path string) error   { return nil }
func (m *mockBackend) GetFiles(contextName string) ([]storage.File, error) { return nil, nil }
func (m *mockBackend) GetFilesByTimestamp(contextName string, after, before time.Time) ([]storage.File, error) {
	return nil, nil
}
func (m *mockBackend) AddTouch(contextName, timestamp string) error           { return nil }
func (m *mockBackend) GetTouches(contextName string) ([]storage.Touch, error) { return nil, nil }
func (m *mockBackend) LogTransition(transition *storage.Transition) error     { return nil }
func (m *mockBackend) GetTransitions(limit int) ([]storage.Transition, error) { return nil, nil }
func (m *mockBackend) GetTransitionsByContext(contextName string) ([]storage.Transition, error) {
	return nil, nil
}
func (m *mockBackend) GetActiveContext() (string, error)         { return "", nil }
func (m *mockBackend) SetActiveContext(contextName string) error { return nil }
func (m *mockBackend) ClearActiveContext() error                 { return nil }
func (m *mockBackend) GetContextsByLabel(label string) ([]*models.ContextWithMetadata, error) {
	return nil, nil
}
func (m *mockBackend) GetContextsByParent(parent string) ([]*models.ContextWithMetadata, error) {
	return nil, nil
}
func (m *mockBackend) SearchContextsByName(query string) ([]*models.ContextWithMetadata, error) {
	return nil, nil
}
func (m *mockBackend) SearchNotes(query string) (map[string][]storage.Note, error) { return nil, nil }
func (m *mockBackend) GetContextCount() (int, error)                               { return 0, nil }
func (m *mockBackend) GetContextStats(contextName string) (storage.ContextStats, error) {
	return storage.ContextStats{}, nil
}
func (m *mockBackend) GetGlobalStats() (storage.GlobalStats, error) {
	return storage.GlobalStats{}, nil
}

func (m *mockBackend) GetContext(name string) (*models.ContextWithMetadata, error) {
	if m.getContextError != nil {
		return nil, m.getContextError
	}
	ctx, ok := m.contexts[name]
	if !ok {
		return nil, fmt.Errorf("context not found: %s", name)
	}
	return ctx, nil
}

func (m *mockBackend) ListContexts() ([]*models.ContextWithMetadata, error) {
	if m.listContextsError != nil {
		return nil, m.listContextsError
	}
	result := make([]*models.ContextWithMetadata, 0, len(m.contexts))
	for _, ctx := range m.contexts {
		result = append(result, ctx)
	}
	return result, nil
}

func (m *mockBackend) UpdateContext(ctx *models.ContextWithMetadata) error {
	if m.updateContextError != nil {
		return m.updateContextError
	}
	m.contexts[ctx.Name] = ctx
	return nil
}

// Compile-time check
var _ storage.Backend = (*mockBackend)(nil)

// Tests for addTagsDB

func TestAddTagsDB_AddsNewTags(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"existing"})

	added, err := addTagsDB(backend, "test-context", []string{"new-tag", "another"})

	require.NoError(t, err)
	assert.Equal(t, []string{"new-tag", "another"}, added)

	ctx := backend.contexts["test-context"]
	assert.Contains(t, ctx.Metadata.Labels, "existing")
	assert.Contains(t, ctx.Metadata.Labels, "new-tag")
	assert.Contains(t, ctx.Metadata.Labels, "another")
}

func TestAddTagsDB_SkipsDuplicates(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"existing", "tag1"})

	added, err := addTagsDB(backend, "test-context", []string{"tag1", "new-tag", "EXISTING"})

	require.NoError(t, err)
	// Only "new-tag" should be added (existing and tag1 already exist, case-insensitive)
	assert.Equal(t, []string{"new-tag"}, added)
}

func TestAddTagsDB_CaseInsensitive(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"Bug"})

	added, err := addTagsDB(backend, "test-context", []string{"bug", "BUG", "Feature"})

	require.NoError(t, err)
	// Only Feature should be added
	assert.Equal(t, []string{"Feature"}, added)
}

func TestAddTagsDB_ContextNotFound(t *testing.T) {
	backend := newMockBackend()

	_, err := addTagsDB(backend, "nonexistent", []string{"tag"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context not found")
}

func TestAddTagsDB_UpdateError(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{})
	backend.updateContextError = fmt.Errorf("database error")

	_, err := addTagsDB(backend, "test-context", []string{"new-tag"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update context")
}

// Tests for removeTagsDB

func TestRemoveTagsDB_RemovesTags(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"tag1", "tag2", "tag3"})

	removed, err := removeTagsDB(backend, "test-context", []string{"tag1", "tag3"})

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"tag1", "tag3"}, removed)

	ctx := backend.contexts["test-context"]
	assert.Equal(t, []string{"tag2"}, ctx.Metadata.Labels)
}

func TestRemoveTagsDB_CaseInsensitive(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"Bug", "Feature"})

	removed, err := removeTagsDB(backend, "test-context", []string{"bug", "FEATURE"})

	require.NoError(t, err)
	assert.Len(t, removed, 2)
}

func TestRemoveTagsDB_NonexistentTag(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"tag1"})

	removed, err := removeTagsDB(backend, "test-context", []string{"nonexistent"})

	require.NoError(t, err)
	assert.Empty(t, removed)
}

func TestRemoveTagsDB_ContextNotFound(t *testing.T) {
	backend := newMockBackend()

	_, err := removeTagsDB(backend, "nonexistent", []string{"tag"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context not found")
}

// Tests for getContextTagsDB

func TestGetContextTagsDB_ReturnsTags(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{"tag1", "tag2"})

	tags, err := getContextTagsDB(backend, "test-context")

	require.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)
}

func TestGetContextTagsDB_EmptyTags(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("test-context", []string{})

	tags, err := getContextTagsDB(backend, "test-context")

	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestGetContextTagsDB_ContextNotFound(t *testing.T) {
	backend := newMockBackend()

	_, err := getContextTagsDB(backend, "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context not found")
}

// Tests for getAllTagsDB

func TestGetAllTagsDB_ReturnsTagCounts(t *testing.T) {
	backend := newMockBackend()
	backend.addContext("ctx1", []string{"bug", "urgent"})
	backend.addContext("ctx2", []string{"bug", "feature"})
	backend.addContext("ctx3", []string{"feature"})

	tagCounts, err := getAllTagsDB(backend)

	require.NoError(t, err)
	assert.Equal(t, 2, tagCounts["bug"])
	assert.Equal(t, 2, tagCounts["feature"])
	assert.Equal(t, 1, tagCounts["urgent"])
}

func TestGetAllTagsDB_EmptyContexts(t *testing.T) {
	backend := newMockBackend()

	tagCounts, err := getAllTagsDB(backend)

	require.NoError(t, err)
	assert.Empty(t, tagCounts)
}

func TestGetAllTagsDB_ListContextsError(t *testing.T) {
	backend := newMockBackend()
	backend.listContextsError = fmt.Errorf("database error")

	_, err := getAllTagsDB(backend)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list contexts")
}
