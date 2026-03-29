# P1.1 Asset Provider Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement the core asset system interfaces and manager for GoForge's asset pipeline.

**Architecture:** Sync `AssetProvider` interface with one method (`Fetch`). `AssetManager` wraps a provider, adds async fetch via goroutines, tracks `AssetHandle` refs, and queries `AssetManifest`. Manifest is JSON-driven, loaded from disk. No concrete providers yet — tests use a mock.

**Tech Stack:** Go 1.24, standard library only (encoding/json, sync/atomic, os). Module: `github.com/colley-repos/goforge/engine`.

**Conventions:** See `CONVENTIONS.md`. Table-driven tests, `fmt.Errorf("assets: context: %w", err)` for errors, PascalCase exports, camelCase private, `_test.go` in same package. No globals, no init().

---

### Task 1: AssetFormat Enum

**Files:**
- Create: `engine/assets/format.go`
- Test: `engine/assets/format_test.go`

**Step 1: Write the failing test**

```go
// engine/assets/format_test.go
package assets

import "testing"

func TestAssetFormat_String(t *testing.T) {
	tests := []struct {
		name   string
		format AssetFormat
		want   string
	}{
		{"PNG", FormatPNG, "PNG"},
		{"GLB", FormatGLB, "GLB"},
		{"OGG", FormatOGG, "OGG"},
		{"WAV", FormatWAV, "WAV"},
		{"JSON", FormatJSON, "JSON"},
		{"unknown", AssetFormat(99), "AssetFormat(99)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.format.String(); got != tt.want {
				t.Errorf("AssetFormat.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd engine && go test ./assets/ -run TestAssetFormat -v`
Expected: FAIL — package doesn't exist yet.

**Step 3: Write minimal implementation**

```go
// engine/assets/format.go
package assets

import "fmt"

// AssetFormat identifies the file format of an asset.
type AssetFormat int

const (
	// FormatPNG is a 2D texture in PNG format.
	FormatPNG AssetFormat = iota
	// FormatGLB is a 3D model in binary glTF format.
	FormatGLB
	// FormatOGG is audio in OGG Vorbis format.
	FormatOGG
	// FormatWAV is audio in WAV format.
	FormatWAV
	// FormatJSON is a JSON data file.
	FormatJSON
)

var formatNames = [...]string{
	FormatPNG:  "PNG",
	FormatGLB:  "GLB",
	FormatOGG:  "OGG",
	FormatWAV:  "WAV",
	FormatJSON: "JSON",
}

// String returns the human-readable name of the format.
func (f AssetFormat) String() string {
	if int(f) >= 0 && int(f) < len(formatNames) {
		return formatNames[f]
	}
	return fmt.Sprintf("AssetFormat(%d)", f)
}
```

**Step 4: Run test to verify it passes**

Run: `cd engine && go test ./assets/ -run TestAssetFormat -v`
Expected: PASS

**Step 5: Commit**

```bash
git add engine/assets/format.go engine/assets/format_test.go
git commit -m "feat(assets): add AssetFormat enum with String()"
```

---

### Task 2: AssetProvider Interface & AssetResult

**Files:**
- Create: `engine/assets/provider.go`

No test file — this is a pure interface + struct definition with no logic.

**Step 1: Write the implementation**

```go
// engine/assets/provider.go
package assets

// AssetProvider abstracts an asset source. Implementations include local
// filesystem, Meshy REST API, and cache layers. Each is a simple sync
// interface — async orchestration lives in AssetManager.
type AssetProvider interface {
	// Fetch retrieves raw asset bytes by ID and format.
	// Returns the bytes or an error if the asset cannot be found or fetched.
	Fetch(id string, format AssetFormat) ([]byte, error)
}

// AssetResult carries the outcome of an async fetch.
type AssetResult struct {
	ID   string
	Data []byte
	Err  error
}
```

**Step 2: Run existing tests to verify nothing broke**

Run: `cd engine && go test ./assets/ -v`
Expected: PASS (format tests still pass, new file compiles)

**Step 3: Commit**

```bash
git add engine/assets/provider.go
git commit -m "feat(assets): add AssetProvider interface and AssetResult"
```

---

### Task 3: AssetHandle with Ref Counting

**Files:**
- Create: `engine/assets/handle.go`
- Create: `engine/assets/handle_test.go`

**Step 1: Write the failing test**

```go
// engine/assets/handle_test.go
package assets

import "testing"

func TestAssetHandle_RefCounting(t *testing.T) {
	tests := []struct {
		name       string
		acquires   int
		releases   int
		wantRefs   int32
		wantInUse  bool
	}{
		{"new handle starts at zero", 0, 0, 0, false},
		{"one acquire", 1, 0, 1, true},
		{"acquire then release", 1, 1, 0, false},
		{"multiple acquires", 3, 0, 3, true},
		{"partial release", 3, 1, 2, true},
		{"full release", 3, 3, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandle("test-asset", FormatPNG)
			for range tt.acquires {
				h.Acquire()
			}
			for range tt.releases {
				h.Release()
			}
			if got := h.Refs(); got != tt.wantRefs {
				t.Errorf("Refs() = %d, want %d", got, tt.wantRefs)
			}
			if got := h.InUse(); got != tt.wantInUse {
				t.Errorf("InUse() = %v, want %v", got, tt.wantInUse)
			}
		})
	}
}

func TestAssetHandle_ReleaseFloor(t *testing.T) {
	h := NewHandle("test", FormatGLB)
	// Release with no acquires should not go negative
	h.Release()
	if refs := h.Refs(); refs != 0 {
		t.Errorf("Refs() = %d after release-without-acquire, want 0", refs)
	}
}

func TestAssetHandle_Fields(t *testing.T) {
	h := NewHandle("sword-mesh", FormatGLB)
	if h.ID != "sword-mesh" {
		t.Errorf("ID = %q, want %q", h.ID, "sword-mesh")
	}
	if h.Format != FormatGLB {
		t.Errorf("Format = %v, want %v", h.Format, FormatGLB)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd engine && go test ./assets/ -run TestAssetHandle -v`
Expected: FAIL — `NewHandle` undefined.

**Step 3: Write minimal implementation**

```go
// engine/assets/handle.go
package assets

import "sync/atomic"

// AssetHandle is a lightweight reference to a loaded asset.
// Callers acquire and release handles; the manager uses ref counts
// to determine when assets can be garbage collected (P1.5).
type AssetHandle struct {
	ID     string
	Format AssetFormat
	refs   atomic.Int32
}

// NewHandle creates a handle with zero references.
func NewHandle(id string, format AssetFormat) *AssetHandle {
	return &AssetHandle{ID: id, Format: format}
}

// Acquire increments the reference count.
func (h *AssetHandle) Acquire() {
	h.refs.Add(1)
}

// Release decrements the reference count, floored at zero.
func (h *AssetHandle) Release() {
	for {
		cur := h.refs.Load()
		if cur <= 0 {
			return
		}
		if h.refs.CompareAndSwap(cur, cur-1) {
			return
		}
	}
}

// Refs returns the current reference count.
func (h *AssetHandle) Refs() int32 {
	return h.refs.Load()
}

// InUse returns true if the handle has at least one reference.
func (h *AssetHandle) InUse() bool {
	return h.refs.Load() > 0
}
```

**Step 4: Run test to verify it passes**

Run: `cd engine && go test ./assets/ -run TestAssetHandle -v`
Expected: PASS

**Step 5: Commit**

```bash
git add engine/assets/handle.go engine/assets/handle_test.go
git commit -m "feat(assets): add AssetHandle with atomic ref counting"
```

---

### Task 4: AssetManifest — Types and Load/Save

**Files:**
- Create: `engine/assets/manifest.go`
- Create: `engine/assets/manifest_test.go`

**Step 1: Write the failing test for load/save round-trip**

```go
// engine/assets/manifest_test.go
package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManifest_LoadSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")

	original := &AssetManifest{
		Assets: map[string]ManifestEntry{
			"house-01": {
				Path:   "meshes/SM_Bld_House_01.glb",
				Format: FormatGLB,
				Tags:   []string{"building", "synty"},
				Variants: []VariantEntry{
					{Name: "desert", Path: "meshes/SM_Bld_House_01_desert.glb"},
				},
				LODLevels: []string{"meshes/SM_Bld_House_01_LOD1.glb"},
				Assessment: &AssetAssessment{
					Role:       "full_cover",
					Dimensions: [3]float64{4.0, 6.0, 4.0},
					IsSkeletal: false,
					SourcePack: "synty-polygon-city",
					Notes:      "Two-story house",
				},
			},
		},
	}

	if err := SaveManifest(path, original); err != nil {
		t.Fatalf("SaveManifest: %v", err)
	}

	// File should exist
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("manifest file not found: %v", err)
	}

	loaded, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}

	entry, ok := loaded.Assets["house-01"]
	if !ok {
		t.Fatal("expected house-01 in loaded manifest")
	}
	if entry.Path != "meshes/SM_Bld_House_01.glb" {
		t.Errorf("Path = %q, want %q", entry.Path, "meshes/SM_Bld_House_01.glb")
	}
	if entry.Format != FormatGLB {
		t.Errorf("Format = %v, want %v", entry.Format, FormatGLB)
	}
	if len(entry.Tags) != 2 || entry.Tags[0] != "building" {
		t.Errorf("Tags = %v, want [building synty]", entry.Tags)
	}
	if len(entry.Variants) != 1 || entry.Variants[0].Name != "desert" {
		t.Errorf("Variants = %v, want [{desert ...}]", entry.Variants)
	}
	if entry.Assessment == nil || entry.Assessment.Role != "full_cover" {
		t.Errorf("Assessment.Role = %v, want full_cover", entry.Assessment)
	}
	if entry.Assessment.Dimensions != [3]float64{4.0, 6.0, 4.0} {
		t.Errorf("Dimensions = %v, want [4 6 4]", entry.Assessment.Dimensions)
	}
}

func TestManifest_LoadMissing(t *testing.T) {
	_, err := LoadManifest("/nonexistent/manifest.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestManifest_Empty(t *testing.T) {
	m := NewManifest()
	if len(m.Assets) != 0 {
		t.Errorf("new manifest should be empty, got %d entries", len(m.Assets))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd engine && go test ./assets/ -run TestManifest -v`
Expected: FAIL — types undefined.

**Step 3: Write minimal implementation**

```go
// engine/assets/manifest.go
package assets

import (
	"encoding/json"
	"fmt"
	"os"
)

// AssetManifest is a data-driven registry mapping asset IDs to their
// metadata, file paths, variants, and assessment data.
type AssetManifest struct {
	Assets map[string]ManifestEntry `json:"assets"`
}

// ManifestEntry describes a single asset in the manifest.
type ManifestEntry struct {
	// Path is relative to the asset root directory (e.g. "meshes/house.glb")
	Path string `json:"path"`
	// Format is the asset's file format
	Format AssetFormat `json:"format"`
	// Tags for searching/filtering (e.g. "building", "cover", "synty")
	Tags []string `json:"tags,omitempty"`
	// Variants are alternate palette/material versions of this asset
	Variants []VariantEntry `json:"variants,omitempty"`
	// LODLevels are paths to lower level-of-detail versions
	LODLevels []string `json:"lod_levels,omitempty"`
	// Assessment holds classification data (nil until assessed)
	Assessment *AssetAssessment `json:"assessment,omitempty"`
}

// VariantEntry describes an alternate version of an asset (palette swap, etc.).
type VariantEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// AssetAssessment holds visual classification data for an asset.
// See docs/ASSET_PIPELINE.md for the assessment workflow.
type AssetAssessment struct {
	// Role classifies the asset: "half_cover", "full_cover", "prop", "obstacle", etc.
	Role string `json:"role"`
	// Dimensions are width, height, depth in world units
	Dimensions [3]float64 `json:"dimensions"`
	// IsSkeletal is true for rigged/animated meshes that need a Skeleton3D
	IsSkeletal bool `json:"is_skeletal"`
	// SourcePack identifies where the asset came from (e.g. "synty-polygon-city")
	SourcePack string `json:"source_pack"`
	// Notes are free-text observations from assessment
	Notes string `json:"notes,omitempty"`
}

// NewManifest creates an empty manifest.
func NewManifest() *AssetManifest {
	return &AssetManifest{Assets: make(map[string]ManifestEntry)}
}

// LoadManifest reads a manifest from a JSON file.
func LoadManifest(path string) (*AssetManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("assets: read manifest %q: %w", path, err)
	}
	var m AssetManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("assets: parse manifest %q: %w", path, err)
	}
	if m.Assets == nil {
		m.Assets = make(map[string]ManifestEntry)
	}
	return &m, nil
}

// SaveManifest writes a manifest to a JSON file.
func SaveManifest(path string, m *AssetManifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("assets: marshal manifest: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("assets: write manifest %q: %w", path, err)
	}
	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `cd engine && go test ./assets/ -run TestManifest -v`
Expected: PASS

**Step 5: Commit**

```bash
git add engine/assets/manifest.go engine/assets/manifest_test.go
git commit -m "feat(assets): add AssetManifest with load/save and assessment types"
```

---

### Task 5: AssetManifest — Query Methods

**Files:**
- Modify: `engine/assets/manifest.go`
- Modify: `engine/assets/manifest_test.go`

**Step 1: Write the failing tests**

Append to `engine/assets/manifest_test.go`:

```go
func newTestManifest() *AssetManifest {
	return &AssetManifest{
		Assets: map[string]ManifestEntry{
			"house-01": {
				Path: "meshes/house.glb", Format: FormatGLB,
				Tags: []string{"building", "synty"},
				Assessment: &AssetAssessment{Role: "full_cover"},
			},
			"crate-01": {
				Path: "meshes/crate.glb", Format: FormatGLB,
				Tags: []string{"cover", "synty"},
				Assessment: &AssetAssessment{Role: "half_cover"},
			},
			"tree-01": {
				Path: "meshes/tree.glb", Format: FormatGLB,
				Tags: []string{"prop", "nature"},
				Assessment: &AssetAssessment{Role: "prop"},
			},
			"unassessed": {
				Path: "meshes/mystery.glb", Format: FormatGLB,
				Tags: []string{"synty"},
			},
		},
	}
}

func TestManifest_ByTag(t *testing.T) {
	m := newTestManifest()

	tests := []struct {
		tag  string
		want int
	}{
		{"synty", 3},
		{"building", 1},
		{"nature", 1},
		{"nonexistent", 0},
	}
	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			got := m.ByTag(tt.tag)
			if len(got) != tt.want {
				t.Errorf("ByTag(%q) returned %d entries, want %d", tt.tag, len(got), tt.want)
			}
		})
	}
}

func TestManifest_ByRole(t *testing.T) {
	m := newTestManifest()

	tests := []struct {
		role string
		want int
	}{
		{"full_cover", 1},
		{"half_cover", 1},
		{"prop", 1},
		{"nonexistent", 0},
	}
	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			got := m.ByRole(tt.role)
			if len(got) != tt.want {
				t.Errorf("ByRole(%q) returned %d entries, want %d", tt.role, len(got), tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd engine && go test ./assets/ -run "TestManifest_By" -v`
Expected: FAIL — `ByTag` and `ByRole` undefined.

**Step 3: Write minimal implementation**

Append to `engine/assets/manifest.go`:

```go
// ByTag returns all manifest entries that have the given tag.
func (m *AssetManifest) ByTag(tag string) []ManifestEntry {
	var results []ManifestEntry
	for _, entry := range m.Assets {
		for _, t := range entry.Tags {
			if t == tag {
				results = append(results, entry)
				break
			}
		}
	}
	return results
}

// ByRole returns all assessed manifest entries with the given role.
// Entries without an assessment are skipped.
func (m *AssetManifest) ByRole(role string) []ManifestEntry {
	var results []ManifestEntry
	for _, entry := range m.Assets {
		if entry.Assessment != nil && entry.Assessment.Role == role {
			results = append(results, entry)
		}
	}
	return results
}
```

**Step 4: Run test to verify it passes**

Run: `cd engine && go test ./assets/ -run "TestManifest_By" -v`
Expected: PASS

**Step 5: Run all tests**

Run: `cd engine && go test ./assets/ -v`
Expected: ALL PASS

**Step 6: Commit**

```bash
git add engine/assets/manifest.go engine/assets/manifest_test.go
git commit -m "feat(assets): add ByTag and ByRole query methods to manifest"
```

---

### Task 6: AssetManager — Core

**Files:**
- Create: `engine/assets/manager.go`
- Create: `engine/assets/manager_test.go`

**Step 1: Write the failing test — mock provider + sync fetch**

```go
// engine/assets/manager_test.go
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
			"sword": []byte("sword-mesh-data"),
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
```

**Step 2: Run test to verify it fails**

Run: `cd engine && go test ./assets/ -run TestManager -v`
Expected: FAIL — `NewManager` undefined.

**Step 3: Write minimal implementation**

```go
// engine/assets/manager.go
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
```

**Step 4: Run test to verify it passes**

Run: `cd engine && go test ./assets/ -run TestManager -v`
Expected: PASS

**Step 5: Commit**

```bash
git add engine/assets/manager.go engine/assets/manager_test.go
git commit -m "feat(assets): add AssetManager with sync/async fetch and ref tracking"
```

---

### Task 7: AssetManager — Async Fetch Test

**Files:**
- Modify: `engine/assets/manager_test.go`

**Step 1: Write the async test**

Append to `engine/assets/manager_test.go`:

```go
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
```

**Step 2: Run test to verify it passes**

Run: `cd engine && go test ./assets/ -run "TestManager_FetchAsync" -v`
Expected: PASS (implementation already supports this)

**Step 3: Run full test suite**

Run: `cd engine && go test ./assets/ -v`
Expected: ALL PASS

**Step 4: Commit**

```bash
git add engine/assets/manager_test.go
git commit -m "test(assets): add async fetch tests for AssetManager"
```

---

### Task 8: Full Suite Verification

**Step 1: Run all engine tests**

Run: `cd engine && go test ./... -v`
Expected: ALL PASS — both existing engine tests and new assets tests.

**Step 2: Run coverage**

Run: `cd engine && go test ./assets/ -coverprofile=cover.out && go tool cover -func=cover.out`
Expected: 80%+ line coverage on assets package.

**Step 3: No commit needed — verification only.**

---

### Task 9: Update TODO.md

**Files:**
- Modify: `TODO.md`

**Step 1: Mark P1.1 items as complete**

Change the P1.1 checklist items from `- [ ]` to `- [x]`:

```
- [x] `AssetProvider` interface: `Fetch(id, format) -> ([]byte, error)`
- [x] `AssetManager` — wraps provider, routes requests, tracks loaded assets
- [x] `AssetHandle` — lightweight reference returned to callers, ref-counted
- [x] `AssetManifest` — data-driven registry mapping IDs to paths, variants, tags, LOD levels
- [x] Variant metadata: each manifest entry supports multiple material/palette variants per mesh
```

**Step 2: Commit**

```bash
git add TODO.md
git commit -m "docs: mark P1.1 asset provider interface tasks complete"
```
