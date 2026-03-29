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
