package assets

import (
	"fmt"
	"testing"
)

// mockProvider returns canned responses for known asset IDs.
type mockProvider struct {
	assets map[string][]byte
}

func (m *mockProvider) Fetch(id string, format AssetFormat) ([]byte, error) {
	data, ok := m.assets[id]
	if !ok {
		return nil, fmt.Errorf("mock: asset %q not found", id)
	}
	return data, nil
}

func newTestManager() *AssetManager {
	provider := &mockProvider{
		assets: map[string][]byte{
			"sword":  []byte("sword-mesh-data"),
			"shield": []byte("shield-mesh-data"),
		},
	}
	manifest := &AssetManifest{
		Assets: map[string]ManifestEntry{
			"sword":  {Path: "meshes/sword.glb", Format: FormatGLB},
			"shield": {Path: "meshes/shield.glb", Format: FormatGLB},
		},
	}
	return NewManager(provider, manifest)
}

func TestManager_Fetch(t *testing.T) {
	m := newTestManager()

	handle, err := m.Fetch("sword")
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if handle.ID != "sword" {
		t.Errorf("ID = %q, want %q", handle.ID, "sword")
	}
	if handle.Format != FormatGLB {
		t.Errorf("Format = %v, want %v", handle.Format, FormatGLB)
	}
	if handle.Refs() != 1 {
		t.Errorf("Refs = %d, want 1 (auto-acquired on fetch)", handle.Refs())
	}
}

func TestManager_FetchDuplicate(t *testing.T) {
	m := newTestManager()

	h1, _ := m.Fetch("sword")
	h2, _ := m.Fetch("sword")

	if h1 != h2 {
		t.Error("expected same handle for duplicate fetch")
	}
	if h1.Refs() != 2 {
		t.Errorf("Refs = %d, want 2 after two fetches", h1.Refs())
	}
}

func TestManager_FetchUnknown(t *testing.T) {
	m := newTestManager()

	_, err := m.Fetch("nonexistent")
	if err == nil {
		t.Error("expected error for unknown asset ID")
	}
}

func TestManager_FetchNotInManifest(t *testing.T) {
	m := newTestManager()

	_, err := m.Fetch("unlisted")
	if err == nil {
		t.Error("expected error for asset not in manifest")
	}
}

func TestManager_GetData(t *testing.T) {
	m := newTestManager()

	m.Fetch("sword")
	data, ok := m.GetData("sword")
	if !ok {
		t.Fatal("expected data for fetched asset")
	}
	if string(data) != "sword-mesh-data" {
		t.Errorf("data = %q, want %q", string(data), "sword-mesh-data")
	}

	_, ok = m.GetData("unfetched")
	if ok {
		t.Error("expected no data for unfetched asset")
	}
}

func TestManager_Release(t *testing.T) {
	m := newTestManager()

	handle, _ := m.Fetch("sword")
	m.Release(handle)

	if handle.Refs() != 0 {
		t.Errorf("Refs = %d after release, want 0", handle.Refs())
	}
}

func TestManager_FetchAsync(t *testing.T) {
	m := newTestManager()

	ch := m.FetchAsync("shield")
	result := <-ch

	if result.Err != nil {
		t.Fatalf("FetchAsync: %v", result.Err)
	}
	if result.ID != "shield" {
		t.Errorf("ID = %q, want %q", result.ID, "shield")
	}
	if string(result.Data) != "shield-mesh-data" {
		t.Errorf("Data = %q, want %q", string(result.Data), "shield-mesh-data")
	}
}

func TestManager_FetchAsyncError(t *testing.T) {
	m := newTestManager()

	ch := m.FetchAsync("nonexistent")
	result := <-ch

	if result.Err == nil {
		t.Error("expected error for unknown asset")
	}
}
