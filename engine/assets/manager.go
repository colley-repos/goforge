package assets

import (
	"fmt"
	"sync"
)

// AssetManager coordinates asset loading through a provider, tracks loaded
// assets via handles, and provides query access to the manifest.
type AssetManager struct {
	provider AssetProvider
	manifest *AssetManifest

	mu      sync.Mutex
	handles map[string]*AssetHandle
	data    map[string][]byte
}

// NewManager creates an AssetManager with the given provider and manifest.
func NewManager(provider AssetProvider, manifest *AssetManifest) *AssetManager {
	return &AssetManager{
		provider: provider,
		manifest: manifest,
		handles:  make(map[string]*AssetHandle),
		data:     make(map[string][]byte),
	}
}

// Fetch loads an asset synchronously. If already loaded, returns the existing
// handle. Each call increments the handle's ref count.
func (m *AssetManager) Fetch(id string) (*AssetHandle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Already loaded — bump ref count and return
	if h, ok := m.handles[id]; ok {
		h.Acquire()
		return h, nil
	}

	// Must be in the manifest
	entry, ok := m.manifest.Assets[id]
	if !ok {
		return nil, fmt.Errorf("assets: %q not in manifest", id)
	}

	// Fetch from provider
	raw, err := m.provider.Fetch(id, entry.Format)
	if err != nil {
		return nil, fmt.Errorf("assets: fetch %q: %w", id, err)
	}

	h := NewHandle(id, entry.Format)
	h.Acquire()
	m.handles[id] = h
	m.data[id] = raw
	return h, nil
}

// FetchAsync loads an asset in a background goroutine. Returns a channel
// that receives exactly one AssetResult when the fetch completes.
func (m *AssetManager) FetchAsync(id string) <-chan AssetResult {
	ch := make(chan AssetResult, 1)
	go func() {
		handle, err := m.Fetch(id)
		result := AssetResult{ID: id, Err: err}
		if handle != nil {
			m.mu.Lock()
			result.Data = m.data[id]
			m.mu.Unlock()
		}
		ch <- result
	}()
	return ch
}

// GetData returns the raw bytes for a loaded asset.
// Returns nil, false if the asset has not been fetched.
func (m *AssetManager) GetData(id string) ([]byte, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.data[id]
	return d, ok
}

// Release decrements the ref count on a handle.
func (m *AssetManager) Release(h *AssetHandle) {
	h.Release()
}

// Manifest returns the asset manifest for querying.
func (m *AssetManager) Manifest() *AssetManifest {
	return m.manifest
}
