package pacing

// RealTime implements Controller for real-time games.
// Commands are always being collected and resolved immediately.
type RealTime struct {
	phase      Phase
	tickCount  int
	turnNumber int // increments every N ticks (for UI display)
	ticksPerTurn int
}

// NewRealTime creates a real-time pacing controller.
// ticksPerTurn controls how often the "turn number" increments (for UI/wave tracking).
// Use 0 to never increment.
func NewRealTime(ticksPerTurn int) *RealTime {
	return &RealTime{
		phase:        PhaseCollecting,
		ticksPerTurn: ticksPerTurn,
	}
}

func (rt *RealTime) Phase() Phase           { return rt.phase }
func (rt *RealTime) ShouldCollect() bool     { return rt.phase == PhaseCollecting }
func (rt *RealTime) ShouldResolve() bool     { return rt.phase == PhaseCollecting } // always resolve
func (rt *RealTime) CurrentFaction() int      { return -1 }                          // all factions act simultaneously
func (rt *RealTime) TurnNumber() int          { return rt.turnNumber }
func (rt *RealTime) SetPhase(phase Phase)    { rt.phase = phase }
func (rt *RealTime) RequestEndTurn()         {}                                      // no-op in real-time
func (rt *RealTime) RequestPause()           {}                                      // no-op (use RTwP for pause)

func (rt *RealTime) Advance() {
	if rt.phase != PhaseCollecting {
		return
	}
	rt.tickCount++
	if rt.ticksPerTurn > 0 && rt.tickCount%rt.ticksPerTurn == 0 {
		rt.turnNumber++
	}
}
