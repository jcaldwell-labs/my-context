package watch

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jefferycaldwell/my-context-copilot/pkg/utils"
)

// Monitor handles context watching and change detection
type Monitor struct {
	contextDir string
	lastMtime  time.Time
	useInotify bool
	mu         sync.RWMutex
}

// NewMonitor creates a new context monitor
func NewMonitor(contextDir string) (*Monitor, error) {
	// Get initial modification time
	mtime, err := utils.GetModTime(contextDir)
	if err != nil {
		if utils.FileExists(contextDir) {
			return nil, fmt.Errorf("failed to get context directory mtime: %w", err)
		}
		// Directory doesn't exist yet, use zero time
		mtime = time.Time{}
	}

	// Enable inotify on Linux for better performance
	useInotify := runtime.GOOS == "linux"

	return &Monitor{
		contextDir: contextDir,
		lastMtime:  mtime,
		useInotify: useInotify,
	}, nil
}

// CheckForChanges checks if the context has been modified since last check
func (m *Monitor) CheckForChanges() (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	currentMtime, err := utils.GetModTime(m.contextDir)
	if err != nil {
		return false, fmt.Errorf("failed to get current mtime: %w", err)
	}

	hasChanged := currentMtime.After(m.lastMtime)
	if hasChanged {
		m.lastMtime = currentMtime
	}

	return hasChanged, nil
}

// CheckForNewNotes checks specifically for new notes in the context
func (m *Monitor) CheckForNewNotes() (bool, error) {
	return m.CheckForChanges() // For now, any change means potential new notes
}

// CheckForNewNotesWithPattern checks for new notes matching a pattern
func (m *Monitor) CheckForNewNotesWithPattern(pattern string) (bool, error) {
	hasChanged, err := m.CheckForNewNotes()
	if err != nil || !hasChanged {
		return hasChanged, err
	}

	// If we have a pattern, we need to check if the new content matches
	if pattern == "" {
		return true, nil
	}

	// For pattern matching, we'd need to read the notes file and check
	// This is a simplified implementation - in practice, you'd want to
	// track the last read position and only check new content
	return m.notesMatchPattern(pattern)
}

// notesMatchPattern checks if recent notes match the given pattern
func (m *Monitor) notesMatchPattern(pattern string) (bool, error) {
	// This is a simplified implementation
	// In a full implementation, you'd need to:
	// 1. Read the notes.log file
	// 2. Parse entries added since last check
	// 3. Check if any match the pattern

	// For now, we return true if pattern is empty (matches everything)
	// or if we have a basic implementation
	if pattern == "" {
		return true, nil
	}

	// Validate the pattern by creating a matcher
	_, err := NewNotePatternMatcher(pattern)
	if err != nil {
		return false, err
	}

	// Placeholder: in a real implementation, you'd read the notes file
	// and check recent entries. For now, return true to indicate
	// the pattern is valid and would match if notes existed.
	return true, nil
}

// WatchOptions contains options for watching
type WatchOptions struct {
	NewNotesOnly bool
	Pattern      string
	ExecCommand  string
	Interval     time.Duration
}

// WatchResult represents the result of a watch operation
type WatchResult struct {
	HasChanges   bool
	Matched      bool
	Executed     bool
	Error        error
}

// Watcher manages continuous watching with execution
type Watcher struct {
	monitor *Monitor
	options WatchOptions
	stopCh  chan struct{}
	mu      sync.Mutex
	running bool
}

// NewWatcher creates a new watcher
func NewWatcher(contextDir string, options WatchOptions) (*Watcher, error) {
	if options.Interval == 0 {
		options.Interval = 5 * time.Second // Default interval
	}

	monitor, err := NewMonitor(contextDir)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		monitor: monitor,
		options: options,
		stopCh:  make(chan struct{}),
		running: false,
	}, nil
}

// Start begins the watching process
func (w *Watcher) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return fmt.Errorf("watcher is already running")
	}

	w.running = true
	go w.watchLoop()
	return nil
}

// Stop stops the watching process
func (w *Watcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return nil
	}

	w.running = false
	close(w.stopCh)
	return nil
}

// IsRunning returns true if the watcher is currently running
func (w *Watcher) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// watchLoop is the main watching goroutine
func (w *Watcher) watchLoop() {
	ticker := time.NewTicker(w.options.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result := w.checkAndExecute()
			if result.Error != nil {
				// In a real implementation, you'd want to log this
				// For now, we'll continue watching
			}

		case <-w.stopCh:
			return
		}
	}
}

// checkAndExecute checks for changes and executes command if needed
func (w *Watcher) checkAndExecute() WatchResult {
	var hasChanges bool
	var err error

	if w.options.NewNotesOnly {
		hasChanges, err = w.monitor.CheckForNewNotesWithPattern(w.options.Pattern)
	} else {
		hasChanges, err = w.monitor.CheckForChanges()
	}

	if err != nil || !hasChanges {
		return WatchResult{
			HasChanges: hasChanges,
			Matched:    hasChanges,
			Executed:   false,
			Error:      err,
		}
	}

	// Execute command if specified
	if w.options.ExecCommand != "" {
		execErr := w.executeCommand(w.options.ExecCommand)
		return WatchResult{
			HasChanges: true,
			Matched:    true,
			Executed:   execErr == nil,
			Error:      execErr,
		}
	}

	return WatchResult{
		HasChanges: true,
		Matched:    true,
		Executed:   false,
		Error:      nil,
	}
}

// executeCommand runs the specified command
func (w *Watcher) executeCommand(cmd string) error {
	// Simple command execution - split by spaces for basic support
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	command := exec.Command(parts[0], parts[1:]...)
	return command.Run()
}
