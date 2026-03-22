// Package renderer defines the Renderer interface that display backends implement.
//
// CRITICAL: This file lives in engine/ but the engine NEVER imports any concrete renderer.
// Renderers import the engine — not the reverse. This interface exists here so that
// the GameMaster can accept a Renderer without knowing which one it is.
package renderer

import (
	"github.com/colley-repos/goforge/engine/input"
)

// Renderer is the interface implemented by display backends (Ebitengine, Terminal, Raylib, etc.).
// The engine interacts with the outside world exclusively through this interface.
type Renderer interface {
	// Init starts the renderer (opens window, initializes graphics context).
	Init(cfg Config) error

	// PollInput returns raw input events since the last call.
	// The engine's InputMapper will translate these to semantic actions.
	PollInput() []input.RawInput

	// BeginFrame prepares for a new frame of rendering.
	BeginFrame()

	// EndFrame presents the rendered frame to the display.
	EndFrame()

	// Shutdown cleans up renderer resources.
	Shutdown()

	// ShouldQuit returns true when the user has requested to close the window.
	ShouldQuit() bool
}

// Config holds renderer-agnostic display settings.
type Config struct {
	Title      string
	Width      int
	Height     int
	Fullscreen bool
	VSync      bool
	// TPS is ticks per second for the game loop (0 = renderer default)
	TPS int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Title:  "GoForge Game",
		Width:  1280,
		Height: 720,
		VSync:  true,
		TPS:    60,
	}
}
