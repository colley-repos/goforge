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
