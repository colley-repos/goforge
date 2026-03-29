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
