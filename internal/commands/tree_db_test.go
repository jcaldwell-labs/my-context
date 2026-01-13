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

// treeTestBackend is a test implementation of storage.Backend for tree tests
type treeTestBackend struct {
	contexts                 map[string]*models.ContextWithMetadata
	getContextError          error
	getContextsByParentError error
	listContextsError        error
}

func newTreeTestBackend() *treeTestBackend {
	return &treeTestBackend{
		contexts: make(map[string]*models.ContextWithMetadata),
	}
}

func (m *treeTestBackend) addContext(name string, parent string) {
	ctx := &models.ContextWithMetadata{
		Name:      name,
		StartTime: time.Now(),
		Status:    "active",
		Metadata: models.ContextMetadata{
			Parent: parent,
		},
	}
	m.contexts[name] = ctx
}

// Implement storage.Backend interface
func (m *treeTestBackend) Init() error                                          { return nil }
func (m *treeTestBackend) Health() error                                        { return nil }
func (m *treeTestBackend) Close() error                                         { return nil }
func (m *treeTestBackend) CreateContext(ctx *models.ContextWithMetadata) error  { return nil }
func (m *treeTestBackend) DeleteContext(name string) error                      { return nil }
func (m *treeTestBackend) ArchiveContext(name string) error                     { return nil }
func (m *treeTestBackend) AddNote(contextName, timestamp, content string) error { return nil }
func (m *treeTestBackend) GetNotes(contextName string) ([]storage.Note, error)  { return nil, nil }
func (m *treeTestBackend) GetNotesByTimestamp(contextName string, after, before time.Time) ([]storage.Note, error) {
	return nil, nil
}
func (m *treeTestBackend) AddFile(contextName, timestamp, path string) error   { return nil }
func (m *treeTestBackend) GetFiles(contextName string) ([]storage.File, error) { return nil, nil }
func (m *treeTestBackend) GetFilesByTimestamp(contextName string, after, before time.Time) ([]storage.File, error) {
	return nil, nil
}
func (m *treeTestBackend) AddTouch(contextName, timestamp string) error           { return nil }
func (m *treeTestBackend) GetTouches(contextName string) ([]storage.Touch, error) { return nil, nil }
func (m *treeTestBackend) LogTransition(transition *storage.Transition) error     { return nil }
func (m *treeTestBackend) GetTransitions(limit int) ([]storage.Transition, error) { return nil, nil }
func (m *treeTestBackend) GetTransitionsByContext(contextName string) ([]storage.Transition, error) {
	return nil, nil
}
func (m *treeTestBackend) GetActiveContext() (string, error)         { return "", nil }
func (m *treeTestBackend) SetActiveContext(contextName string) error { return nil }
func (m *treeTestBackend) ClearActiveContext() error                 { return nil }
func (m *treeTestBackend) GetContextsByLabel(label string) ([]*models.ContextWithMetadata, error) {
	return nil, nil
}
func (m *treeTestBackend) UpdateContext(ctx *models.ContextWithMetadata) error { return nil }
func (m *treeTestBackend) SearchContextsByName(query string) ([]*models.ContextWithMetadata, error) {
	return nil, nil
}
func (m *treeTestBackend) SearchNotes(query string) (map[string][]storage.Note, error) {
	return nil, nil
}
func (m *treeTestBackend) GetContextCount() (int, error) { return 0, nil }
func (m *treeTestBackend) GetContextStats(contextName string) (storage.ContextStats, error) {
	return storage.ContextStats{}, nil
}
func (m *treeTestBackend) GetGlobalStats() (storage.GlobalStats, error) {
	return storage.GlobalStats{}, nil
}

func (m *treeTestBackend) GetContext(name string) (*models.ContextWithMetadata, error) {
	if m.getContextError != nil {
		return nil, m.getContextError
	}
	ctx, ok := m.contexts[name]
	if !ok {
		return nil, fmt.Errorf("context not found: %s", name)
	}
	return ctx, nil
}

func (m *treeTestBackend) ListContexts() ([]*models.ContextWithMetadata, error) {
	if m.listContextsError != nil {
		return nil, m.listContextsError
	}
	result := make([]*models.ContextWithMetadata, 0, len(m.contexts))
	for _, ctx := range m.contexts {
		result = append(result, ctx)
	}
	return result, nil
}

func (m *treeTestBackend) GetContextsByParent(parent string) ([]*models.ContextWithMetadata, error) {
	if m.getContextsByParentError != nil {
		return nil, m.getContextsByParentError
	}
	var result []*models.ContextWithMetadata
	for _, ctx := range m.contexts {
		if ctx.Metadata.Parent == parent {
			result = append(result, ctx)
		}
	}
	return result, nil
}

// Compile-time check
var _ storage.Backend = (*treeTestBackend)(nil)

// Tests for getContextTreeDB

func TestGetContextTreeDB_SingleNode(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("root", "")

	tree, err := getContextTreeDB(backend, "root")

	require.NoError(t, err)
	assert.Equal(t, "root", tree.Name)
	assert.Empty(t, tree.Children)
}

func TestGetContextTreeDB_WithChildren(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("root", "")
	backend.addContext("child1", "root")
	backend.addContext("child2", "root")

	tree, err := getContextTreeDB(backend, "root")

	require.NoError(t, err)
	assert.Equal(t, "root", tree.Name)
	assert.Len(t, tree.Children, 2)

	childNames := []string{tree.Children[0].Name, tree.Children[1].Name}
	assert.ElementsMatch(t, []string{"child1", "child2"}, childNames)
}

func TestGetContextTreeDB_DeepHierarchy(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("root", "")
	backend.addContext("child", "root")
	backend.addContext("grandchild", "child")

	tree, err := getContextTreeDB(backend, "root")

	require.NoError(t, err)
	assert.Equal(t, "root", tree.Name)
	require.Len(t, tree.Children, 1)
	assert.Equal(t, "child", tree.Children[0].Name)
	require.Len(t, tree.Children[0].Children, 1)
	assert.Equal(t, "grandchild", tree.Children[0].Children[0].Name)
}

func TestGetContextTreeDB_ContextNotFound(t *testing.T) {
	backend := newTreeTestBackend()

	_, err := getContextTreeDB(backend, "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context not found")
}

func TestGetContextTreeDB_GetChildrenError(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("root", "")
	backend.getContextsByParentError = fmt.Errorf("database error")

	_, err := getContextTreeDB(backend, "root")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get children")
}

// Tests for cycle detection

func TestGetContextTreeDB_DetectsCycle(t *testing.T) {
	backend := newTreeTestBackend()
	// Create a cycle: A -> B -> C -> A
	backend.addContext("A", "C") // A's parent is C, but we'll query starting from A
	backend.addContext("B", "A")
	backend.addContext("C", "B")

	// Start from A, it should detect cycle when it reaches A again
	_, err := getContextTreeDB(backend, "A")

	// The cycle will manifest when building the tree:
	// A has children [B] (since B.parent = A)
	// B has children [C] (since C.parent = B)
	// C has children [A] (since A.parent = C) <- cycle detected!
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cycle detected")
}

func TestGetContextTreeDB_SeparateTreesNoFalseCycle(t *testing.T) {
	// Regression test: ensure separate hierarchies don't cause false cycle detection
	backend := newTreeTestBackend()
	backend.addContext("tree1-root", "")
	backend.addContext("tree1-child", "tree1-root")
	backend.addContext("tree2-root", "")
	backend.addContext("tree2-child", "tree2-root")

	// Build first tree
	tree1, err := getContextTreeDB(backend, "tree1-root")
	require.NoError(t, err)
	assert.Equal(t, "tree1-root", tree1.Name)

	// Build second tree - should not trigger false cycle from first tree
	tree2, err := getContextTreeDB(backend, "tree2-root")
	require.NoError(t, err)
	assert.Equal(t, "tree2-root", tree2.Name)
}

// Tests for getRootContextsDB

func TestGetRootContextsDB_ReturnsRoots(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("root1", "")
	backend.addContext("root2", "")
	backend.addContext("child1", "root1")

	roots, err := getRootContextsDB(backend)

	require.NoError(t, err)
	assert.Len(t, roots, 2)
	assert.ElementsMatch(t, []string{"root1", "root2"}, roots)
}

func TestGetRootContextsDB_NoContexts(t *testing.T) {
	backend := newTreeTestBackend()

	roots, err := getRootContextsDB(backend)

	require.NoError(t, err)
	assert.Empty(t, roots)
}

func TestGetRootContextsDB_ListContextsError(t *testing.T) {
	backend := newTreeTestBackend()
	backend.listContextsError = fmt.Errorf("database error")

	_, err := getRootContextsDB(backend)

	assert.Error(t, err)
}

// Tests for getChildrenDB

func TestGetChildrenDB_ReturnsChildren(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("parent", "")
	backend.addContext("child1", "parent")
	backend.addContext("child2", "parent")
	backend.addContext("other", "")

	children, err := getChildrenDB(backend, "parent")

	require.NoError(t, err)
	assert.Len(t, children, 2)
	assert.ElementsMatch(t, []string{"child1", "child2"}, children)
}

func TestGetChildrenDB_NoChildren(t *testing.T) {
	backend := newTreeTestBackend()
	backend.addContext("leaf", "")

	children, err := getChildrenDB(backend, "leaf")

	require.NoError(t, err)
	assert.Empty(t, children)
}

func TestGetChildrenDB_Error(t *testing.T) {
	backend := newTreeTestBackend()
	backend.getContextsByParentError = fmt.Errorf("database error")

	_, err := getChildrenDB(backend, "parent")

	assert.Error(t, err)
}
