package events

// Common event types used across GoForge games.
// Games should define their own domain-specific events too.

// CommandQueued is emitted when a command is added to the queue.
type CommandQueued struct {
	CommandType string
	EntityID    uint64
}

// CommandResolved is emitted when a command has been executed.
type CommandResolved struct {
	CommandType string
	EntityID    uint64
	Success     bool
}

// AllCommandsResolved is emitted when the command queue is fully drained.
type AllCommandsResolved struct{}

// TurnStarted is emitted at the beginning of a turn.
type TurnStarted struct {
	TurnNumber int
	FactionID  int
}

// TurnEnded is emitted at the end of a turn.
type TurnEnded struct {
	TurnNumber int
	FactionID  int
}

// PhaseChanged is emitted when the pacing controller changes phase.
type PhaseChanged struct {
	OldPhase string
	NewPhase string
}

// EntityDamaged is emitted when an entity takes damage.
type EntityDamaged struct {
	EntityID uint64
	Amount   int
	NewHP    int
}

// EntityDied is emitted when an entity's HP reaches zero.
type EntityDied struct {
	EntityID uint64
}

// EntityMoved is emitted when an entity changes position.
type EntityMoved struct {
	EntityID uint64
	FromX    float64
	FromY    float64
	ToX      float64
	ToY      float64
}

// UIMessage carries a string for the HUD/combat log to display.
type UIMessage struct {
	Text     string
	Category string // "info", "combat", "system"
}

// GameOver signals that the game has ended.
type GameOver struct {
	Victory bool
	Reason  string
}
