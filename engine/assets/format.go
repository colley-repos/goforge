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
