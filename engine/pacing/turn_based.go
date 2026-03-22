package pacing

// TurnBased implements Controller for turn-based games.
// Commands are collected during the active faction's turn
// and resolved when RequestEndTurn() is called.
type TurnBased struct {
	phase          Phase
	factions       []int // faction IDs in turn order
	currentIndex   int
	turnNumber     int
	pendingResolve bool
}

// NewTurnBased creates a turn-based pacing controller.
// factions is the list of faction IDs in turn order (e.g., [0, 1] for player vs enemy).
func NewTurnBased(factions []int) *TurnBased {
	return &TurnBased{
		phase:    PhaseCollecting,
		factions: factions,
	}
}

func (tb *TurnBased) Phase() Phase           { return tb.phase }
func (tb *TurnBased) ShouldCollect() bool     { return tb.phase == PhaseCollecting }
func (tb *TurnBased) ShouldResolve() bool     { return tb.pendingResolve }
func (tb *TurnBased) CurrentFaction() int      { return tb.factions[tb.currentIndex] }
func (tb *TurnBased) TurnNumber() int          { return tb.turnNumber }
func (tb *TurnBased) SetPhase(phase Phase)    { tb.phase = phase }
func (tb *TurnBased) RequestPause()           {} // no-op in turn-based

func (tb *TurnBased) RequestEndTurn() {
	if tb.phase == PhaseCollecting {
		tb.pendingResolve = true
		tb.phase = PhaseResolving
	}
}

func (tb *TurnBased) Advance() {
	switch tb.phase {
	case PhaseResolving:
		// After resolution, move to next faction
		tb.pendingResolve = false
		tb.currentIndex = (tb.currentIndex + 1) % len(tb.factions)
		if tb.currentIndex == 0 {
			tb.turnNumber++
		}
		tb.phase = PhaseCollecting
	}
}
