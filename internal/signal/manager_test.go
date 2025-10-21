package signal

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestSignalsDir(t *testing.T) string {
	tempDir := t.TempDir()
	signalsDir := filepath.Join(tempDir, "signals")
	return signalsDir
}

func TestNewManager(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)

	manager, err := NewManager(signalsDir)
	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, signalsDir, manager.signalsDir)

	// Verify directory was created
	assert.DirExists(t, signalsDir)
}

func TestCreateSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Create a signal
	signal, err := manager.CreateSignal("test-signal")
	require.NoError(t, err)
	assert.NotNil(t, signal)
	assert.Equal(t, "test-signal", signal.Name)
	assert.True(t, signal.Exists())

	// Verify file was created
	expectedPath := filepath.Join(signalsDir, "test-signal.signal")
	assert.FileExists(t, expectedPath)
	assert.Equal(t, expectedPath, signal.Path)
}

func TestCreateDuplicateSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Create first signal
	_, err = manager.CreateSignal("duplicate")
	require.NoError(t, err)

	// Try to create duplicate
	_, err = manager.CreateSignal("duplicate")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestGetSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Create a signal first
	createdSignal, err := manager.CreateSignal("get-test")
	require.NoError(t, err)

	// Get the signal
	retrievedSignal, err := manager.GetSignal("get-test")
	require.NoError(t, err)
	assert.NotNil(t, retrievedSignal)
	assert.Equal(t, "get-test", retrievedSignal.Name)
	assert.Equal(t, createdSignal.Path, retrievedSignal.Path)
}

func TestGetNonexistentSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	_, err = manager.GetSignal("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestListSignals(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Create multiple signals
	_, err = manager.CreateSignal("signal1")
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	_, err = manager.CreateSignal("signal2")
	require.NoError(t, err)

	// List signals
	signals, err := manager.ListSignals()
	require.NoError(t, err)
	assert.Len(t, signals, 2)

	// Check that both signals are present
	names := make([]string, len(signals))
	for i, s := range signals {
		names[i] = s.Name
	}
	assert.Contains(t, names, "signal1")
	assert.Contains(t, names, "signal2")
}

func TestClearSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Create and then clear a signal
	_, err = manager.CreateSignal("to-clear")
	require.NoError(t, err)

	err = manager.ClearSignal("to-clear")
	assert.NoError(t, err)

	// Verify it's gone
	assert.False(t, manager.SignalExists("to-clear"))
}

func TestClearNonexistentSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	err = manager.ClearSignal("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestWaitForSignal(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Start waiting in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		manager.CreateSignal("wait-test")
	}()

	// Wait for the signal
	err = manager.WaitForSignal("wait-test", 2*time.Second)
	assert.NoError(t, err)
	assert.True(t, manager.SignalExists("wait-test"))
}

func TestWaitForSignalTimeout(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Wait for a signal that never comes
	err = manager.WaitForSignal("never-appears", 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestSignalExists(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Signal doesn't exist initially
	assert.False(t, manager.SignalExists("exists-test"))

	// Create signal
	_, err = manager.CreateSignal("exists-test")
	require.NoError(t, err)

	// Now it exists
	assert.True(t, manager.SignalExists("exists-test"))

	// Clear it
	err = manager.ClearSignal("exists-test")
	require.NoError(t, err)

	// No longer exists
	assert.False(t, manager.SignalExists("exists-test"))
}

func TestClearAllSignals(t *testing.T) {
	signalsDir := setupTestSignalsDir(t)
	manager, err := NewManager(signalsDir)
	require.NoError(t, err)

	// Create multiple signals
	_, err = manager.CreateSignal("clear1")
	require.NoError(t, err)
	_, err = manager.CreateSignal("clear2")
	require.NoError(t, err)
	_, err = manager.CreateSignal("clear3")
	require.NoError(t, err)

	// Verify they exist
	signals, err := manager.ListSignals()
	require.NoError(t, err)
	assert.Len(t, signals, 3)

	// Clear all
	err = manager.ClearAllSignals()
	assert.NoError(t, err)

	// Verify all are gone
	signals, err = manager.ListSignals()
	require.NoError(t, err)
	assert.Len(t, signals, 0)
}
