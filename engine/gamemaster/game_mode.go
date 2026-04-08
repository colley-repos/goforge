package gamemaster

// GameMode defines the rules, class factories, and lifecycle hooks for a game type.
//
// This mirrors Unreal Engine's GameMode philosophy: swap the GameMode to completely
// change the game type without touching GameMaster, the ECS world, or any systems.
//
// GameMode answers:
//   - What are the win/loss conditions?
//   - Who gets spawned, where, and when?
//   - What default systems does this game type need?
//   - How does the match start and end?
//
// Concrete game modes should embed BaseGameMode and override only what they need.
//
// Usage:
//
//	type MyGameMode struct{ gamemaster.BaseGameMode }
//
//	func (m *MyGameMode) OnMatchStart(gm *GameMaster) {
//	    // spawn entities, set up encounter, subscribe to events
//	}
//
//	gm := gamemaster.New(gamemaster.Config{
//	    Mode:   &MyGameMode{},
//	    Pacing: pacing.NewRealTime(0),
//	    ...
//	})
type GameMode interface {
	// OnMatchStart is called once, after GameMaster is fully initialized.
	// Spawn initial entities, subscribe to events, wire encounter logic here.
	OnMatchStart(gm *GameMaster)

	// OnMatchEnd is called when the match concludes.
	// winner is a team/player ID, or -1 for a draw.
	OnMatchEnd(gm *GameMaster, winner int)

	// OnPlayerJoin is called when a player joins the session.
	// Spawn their pawn and player state here.
	OnPlayerJoin(gm *GameMaster, playerID int)

	// OnPlayerLeave is called when a player leaves.
	// Clean up their entities and state here.
	OnPlayerLeave(gm *GameMaster, playerID int)

	// ShouldRespawn returns true if the player should be respawned.
	// Called by the respawn system after a pawn death event.
	ShouldRespawn(gm *GameMaster, playerID int) bool

	// SpawnConfig returns the spawn parameters for a player.
	// Called by OnPlayerJoin and the respawn system.
	GetSpawnConfig(gm *GameMaster, playerID int) SpawnConfig

	// DefaultSystems returns systems to register before any game-specific ones.
	// Use this for systems that belong to the game type, not the engine.
	DefaultSystems() []System
}

// SpawnConfig carries the spawn parameters for a single player.
type SpawnConfig struct {
	// GridX, GridY are the spawn tile coordinates (for grid-based games).
	GridX, GridY int
	// Team is the faction/team index. 0 = player, 1 = enemy, -1 = neutral.
	Team int
	// Facing is the initial facing direction in degrees (0 = north).
	Facing float64
}

// BaseGameMode provides no-op implementations of all GameMode methods.
// Embed this in your concrete GameMode and override only what you need.
//
//	type MyMode struct {
//	    gamemaster.BaseGameMode
//	}
//	func (m *MyMode) OnMatchStart(gm *gamemaster.GameMaster) {
//	    // only OnMatchStart is customised; everything else is a no-op
//	}
type BaseGameMode struct{}

func (BaseGameMode) OnMatchStart(*GameMaster)                        {}
func (BaseGameMode) OnMatchEnd(*GameMaster, int)                     {}
func (BaseGameMode) OnPlayerJoin(*GameMaster, int)                   {}
func (BaseGameMode) OnPlayerLeave(*GameMaster, int)                  {}
func (BaseGameMode) ShouldRespawn(_ *GameMaster, _ int) bool         { return false }
func (BaseGameMode) GetSpawnConfig(_ *GameMaster, _ int) SpawnConfig { return SpawnConfig{} }
func (BaseGameMode) DefaultSystems() []System                        { return nil }
