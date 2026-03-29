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
