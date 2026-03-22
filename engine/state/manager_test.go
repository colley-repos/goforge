package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotSetGet(t *testing.T) {
	snap := NewSnapshot(100, 5)

	type AgentData struct {
		Name string `json:"name"`
		HP   int    `json:"hp"`
	}

	err := snap.Set("agent_alpha", AgentData{Name: "Alpha", HP: 10})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	var loaded AgentData
	err = snap.Get("agent_alpha", &loaded)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if loaded.Name != "Alpha" || loaded.HP != 10 {
		t.Errorf("got %+v, want Alpha/10", loaded)
	}

	// Missing key
	err = snap.Get("missing", &loaded)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestUndoStack(t *testing.T) {
	m := NewManager(t.TempDir(), 3)

	// Push 4 snapshots with max depth 3
	m.PushUndo(NewSnapshot(1, 0))
	m.PushUndo(NewSnapshot(2, 0))
	m.PushUndo(NewSnapshot(3, 1))
	m.PushUndo(NewSnapshot(4, 1))

	if m.UndoDepth() != 3 {
		t.Fatalf("expected depth 3 (capped), got %d", m.UndoDepth())
	}

	// Pop should return most recent first
	snap := m.PopUndo()
	if snap.TickNumber != 4 {
		t.Errorf("expected tick 4, got %d", snap.TickNumber)
	}

	snap = m.PopUndo()
	if snap.TickNumber != 3 {
		t.Errorf("expected tick 3, got %d", snap.TickNumber)
	}

	// Oldest (tick 1) was evicted, so next should be tick 2
	snap = m.PopUndo()
	if snap.TickNumber != 2 {
		t.Errorf("expected tick 2, got %d", snap.TickNumber)
	}

	// Empty
	if m.PopUndo() != nil {
		t.Error("should return nil when empty")
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir, 10)

	snap := NewSnapshot(42, 3)
	snap.Set("score", 1500)

	err := m.Save("test_save.json", snap)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// File should exist
	if _, err := os.Stat(filepath.Join(dir, "test_save.json")); err != nil {
		t.Fatalf("save file not found: %v", err)
	}

	// Load it back
	loaded, err := m.Load("test_save.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.TickNumber != 42 || loaded.TurnNumber != 3 {
		t.Errorf("loaded snap: tick=%d turn=%d, want 42/3", loaded.TickNumber, loaded.TurnNumber)
	}

	var score int
	loaded.Get("score", &score)
	if score != 1500 {
		t.Errorf("loaded score: %d, want 1500", score)
	}
}
