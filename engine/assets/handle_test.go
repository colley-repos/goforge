// engine/assets/handle_test.go
package assets

import "testing"

func TestAssetHandle_RefCounting(t *testing.T) {
	tests := []struct {
		name      string
		acquires  int
		releases  int
		wantRefs  int32
		wantInUse bool
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
