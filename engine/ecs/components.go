package ecs

// Position represents a 2D position in world space.
// For grid-based games, X and Y are tile coordinates.
// For free-movement games, they are world units.
type Position struct {
	X float64
	Y float64
}

// Velocity represents directional speed in world units per tick.
type Velocity struct {
	DX float64
	DY float64
}

// Health tracks current and maximum hit points.
type Health struct {
	Current int
	Max     int
}

// Faction identifies which team/side an entity belongs to.
type Faction struct {
	ID   int
	Name string
}

// SpriteDescriptor is a renderer-agnostic description of how to display an entity.
// The engine never imports rendering code — this struct carries enough info
// for any renderer to decide what to draw.
type SpriteDescriptor struct {
	// AssetID references an asset in the AssetProvider (empty = use shape)
	AssetID string
	// Shape is the fallback whitebox shape: "rect", "circle", "diamond", "cube"
	Shape string
	// Color as RGBA hex string (e.g., "#FF0000FF")
	Color string
	// Width and Height in logical units
	Width  float64
	Height float64
	// Layer for draw ordering (higher = on top)
	Layer int
}

// InputControllable marks an entity as player-controllable.
type InputControllable struct {
	// PlayerID identifies which player controls this entity (for co-op)
	PlayerID int
}

// AIControllable marks an entity as AI-controlled.
type AIControllable struct {
	// BrainID identifies which AI behavior to use
	BrainID string
}

// Collidable marks an entity as participating in collision detection.
type Collidable struct {
	// Solid means other entities cannot move through this
	Solid bool
	// Radius for circular collision (0 = use sprite dimensions)
	Radius float64
}

// GridPosition represents a discrete tile position on a grid.
// Separate from Position to allow both grid-snapped and free movement.
type GridPosition struct {
	Col int
	Row int
}

// ActionPoints tracks remaining actions per turn (for turn-based games).
type ActionPoints struct {
	Current int
	Max     int
}

// Named gives an entity a display name (for UI, logs, debugging).
type Named struct {
	DisplayName string
}
