// Package state provides game state snapshot and persistence.
//
// StateManager captures the full ECS world state for save/load,
// undo/redo, and future networking (state diff/sync).
package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Snapshot is a serializable capture of game state at a point in time.
// Games extend this with their own fields.
type Snapshot struct {
	// TickNumber when the snapshot was taken
	TickNumber uint64 `json:"tick_number"`
	// TurnNumber when the snapshot was taken
	TurnNumber int `json:"turn_number"`
	// Data holds game-specific serialized state
	Data map[string]json.RawMessage `json:"data"`
}

// NewSnapshot creates an empty snapshot for the given tick.
func NewSnapshot(tick uint64, turn int) *Snapshot {
	return &Snapshot{
		TickNumber: tick,
		TurnNumber: turn,
		Data:       make(map[string]json.RawMessage),
	}
}

// Set stores a named piece of data in the snapshot.
func (s *Snapshot) Set(key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("state: marshal %q: %w", key, err)
	}
	s.Data[key] = raw
	return nil
}

// Get deserializes a named piece of data from the snapshot.
func (s *Snapshot) Get(key string, target any) error {
	raw, ok := s.Data[key]
	if !ok {
		return fmt.Errorf("state: key %q not found", key)
	}
	return json.Unmarshal(raw, target)
}

// Manager handles save/load operations and undo history.
type Manager struct {
	// SaveDir is the directory where save files are written
	SaveDir string
	// undoStack holds recent snapshots for undo
	undoStack []*Snapshot
	// maxUndo is the maximum undo depth
	maxUndo int
}

// NewManager creates a state manager with the given save directory and undo depth.
func NewManager(saveDir string, maxUndo int) *Manager {
	return &Manager{
		SaveDir: saveDir,
		maxUndo: maxUndo,
	}
}

// PushUndo adds a snapshot to the undo stack.
func (m *Manager) PushUndo(snap *Snapshot) {
	m.undoStack = append(m.undoStack, snap)
	if len(m.undoStack) > m.maxUndo {
		m.undoStack = m.undoStack[1:]
	}
}

// PopUndo removes and returns the most recent undo snapshot, or nil.
func (m *Manager) PopUndo() *Snapshot {
	if len(m.undoStack) == 0 {
		return nil
	}
	last := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]
	return last
}

// UndoDepth returns how many undo snapshots are available.
func (m *Manager) UndoDepth() int {
	return len(m.undoStack)
}

// Save writes a snapshot to disk as JSON.
func (m *Manager) Save(filename string, snap *Snapshot) error {
	if err := os.MkdirAll(m.SaveDir, 0o755); err != nil {
		return fmt.Errorf("state: mkdir: %w", err)
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("state: marshal: %w", err)
	}

	path := filepath.Join(m.SaveDir, filename)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("state: write: %w", err)
	}

	return nil
}

// Load reads a snapshot from disk.
func (m *Manager) Load(filename string) (*Snapshot, error) {
	path := filepath.Join(m.SaveDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("state: read: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("state: unmarshal: %w", err)
	}

	return &snap, nil
}
