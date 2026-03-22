// Package pacing defines the PacingController interface and built-in implementations.
//
// The PacingController is the SOLE system that changes when swapping between
// turn-based, real-time, or real-time-with-pause. Every other system in GoForge
// is pacing-agnostic.
package pacing

// Phase represents the current game phase.
type Phase int

const (
	PhaseSetup      Phase = iota // Before the game starts
	PhaseCollecting              // Accepting commands (player's turn, or always in RT)
	PhaseResolving               // Processing queued commands
	PhaseWaiting                 // Waiting for animations/transitions
	PhaseGameOver                // Game has ended
)

// String returns the phase name.
func (p Phase) String() string {
	switch p {
	case PhaseSetup:
		return "Setup"
	case PhaseCollecting:
		return "Collecting"
	case PhaseResolving:
		return "Resolving"
	case PhaseWaiting:
		return "Waiting"
	case PhaseGameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

// Controller is the interface that governs game timing.
// It decides when commands are collected and when they are resolved.
// Implementations: TurnBased, RealTime, RTwP (real-time with pause).
type Controller interface {
	// Phase returns the current game phase.
	Phase() Phase

	// ShouldCollect returns true if the system should accept new commands.
	ShouldCollect() bool

	// ShouldResolve returns true if queued commands should be resolved now.
	ShouldResolve() bool

	// Advance moves the pacing state forward (e.g., end turn, advance phase).
	// Called by the game master each tick.
	Advance()

	// RequestEndTurn signals that the current actor wants to end their turn.
	// Only meaningful in turn-based mode; no-op in real-time.
	RequestEndTurn()

	// RequestPause toggles pause state.
	// Only meaningful in RTwP mode; no-op in turn-based.
	RequestPause()

	// SetPhase forces a phase transition (for game-over, etc.).
	SetPhase(phase Phase)

	// CurrentFaction returns the faction ID of the currently active faction.
	// Returns -1 if not applicable (real-time mode).
	CurrentFaction() int

	// TurnNumber returns the current turn count.
	TurnNumber() int
}
