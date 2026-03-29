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
