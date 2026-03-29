# P1.1 Design: AssetProvider Interface & Manager

## Decision: Approach C — Sync Interface, Async Wrapper

Simple synchronous `AssetProvider` interface. The `AssetManager` adds async
orchestration on top. Cache (P1.4) wraps any provider as a decorator. Providers
stay trivially testable with one method.

## Core Types (`engine/assets/`)

### AssetFormat

```go
type AssetFormat int
const (
    FormatPNG  AssetFormat = iota // 2D textures
    FormatGLB                     // 3D models
    FormatOGG                     // audio
    FormatWAV                     // audio
    FormatJSON                    // data files
)
```

### AssetProvider Interface

```go
type AssetProvider interface {
    Fetch(id string, format AssetFormat) ([]byte, error)
}
```

One method. Local filesystem blocks briefly, Meshy does submit-poll-download
internally, cache checks disk first then delegates. All the same interface.

### AssetResult

```go
type AssetResult struct {
    Data []byte
    Err  error
}
```

Returned via channel from `AssetManager.FetchAsync()`.

### AssetHandle

```go
type AssetHandle struct {
    ID     string
    Format AssetFormat
    refs   int32 // atomic ref count
}
```

Lightweight reference returned to callers. Ref counting enables garbage
collection in P1.5.

## AssetManifest

JSON-driven registry loaded from `manifest.json` in the asset root.

```go
type AssetManifest struct {
    Assets map[string]ManifestEntry
}

type ManifestEntry struct {
    Path       string            // relative to asset root: "meshes/SM_Bld_House_01.glb"
    Format     AssetFormat
    Tags       []string          // searchable: "building", "cover", "synty"
    Variants   []VariantEntry    // palette/material swaps
    LODLevels  []string          // paths to LOD variants
    Assessment *AssetAssessment  // nil until assessed
}

type VariantEntry struct {
    Name string // "desert", "snow"
    Path string // "meshes/SM_Bld_House_01_desert.glb"
}

type AssetAssessment struct {
    Role       string     // "half_cover", "full_cover", "prop", etc.
    Dimensions [3]float64 // width, height, depth
    IsSkeletal bool
    SourcePack string     // "synty-polygon-city", "meshy"
    Notes      string
}
```

Aligns with `docs/ASSET_PIPELINE.md` assessment workflow. Assets are assessed
externally (browsed from source library at `F:\GameDev\Assets\Synty-Godot`),
selected, migrated into the project's local `assets/` directory, then registered
in the manifest.

## AssetManager

Coordinator that wraps a provider, tracks handles, provides async fetch.

```go
type AssetManager struct {
    provider AssetProvider
    manifest *AssetManifest
    handles  map[string]*AssetHandle
    data     map[string][]byte // loaded asset bytes, keyed by ID
}

func (m *AssetManager) Fetch(id string) (*AssetHandle, error)
func (m *AssetManager) FetchAsync(id string) <-chan AssetResult
func (m *AssetManager) Acquire(h *AssetHandle)
func (m *AssetManager) Release(h *AssetHandle)
func (m *AssetManager) ByTag(tag string) []ManifestEntry
func (m *AssetManager) ByRole(role string) []ManifestEntry
```

## Asset Directory Structure

```
assets/
  meshes/        # .glb files
  textures/      # .png files
  audio/         # .ogg/.wav files
  manifest.json
```

Flat per-category. Assets migrated in from source libraries (Synty, Meshy, etc.)
after assessment and selection. Runtime provider reads only from this directory.

## Package Layout

```
engine/assets/
    format.go        -- AssetFormat enum, String()
    provider.go      -- AssetProvider interface, AssetResult
    handle.go        -- AssetHandle with atomic ref counting
    manifest.go      -- AssetManifest, ManifestEntry, load/save/query
    manager.go       -- AssetManager orchestrator
    manager_test.go  -- tests against mock provider
    manifest_test.go -- manifest load/save/query tests
    handle_test.go   -- ref counting tests
```

No concrete providers in this phase. Local filesystem (P1.2) and Meshy (P1.3)
come later as separate packages.

## Testing Strategy

- Mock provider returning canned bytes for known IDs, errors for unknown
- Manifest: load from JSON, query by tag/role, round-trip save/load
- Handle: acquire increments, release decrements, ref count lifecycle
- Manager: sync fetch, async fetch, duplicate fetch returns same handle,
  unknown ID errors
- Table-driven, deterministic, no filesystem (embedded test data)
