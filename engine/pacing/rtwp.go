package pacing

// RTwP implements Controller for Real-Time with Pause games.
// Commands are collected continuously but only resolve when unpaused.
// Auto-pause can be triggered on configurable events.
type RTwP struct {
	phase    Phase
	paused   bool
	tickCount  int
	turnNumber int
	ticksPerTurn int
}

// NewRTwP creates a real-time-with-pause pacing controller.
func NewRTwP(ticksPerTurn int) *RTwP {
	return &RTwP{
		phase:        PhaseCollecting,
		paused:       true, // start paused
		ticksPerTurn: ticksPerTurn,
	}
}

func (r *RTwP) Phase() Phase        { return r.phase }
func (r *RTwP) ShouldCollect() bool  { return r.phase == PhaseCollecting }
func (r *RTwP) ShouldResolve() bool  { return r.phase == PhaseCollecting && !r.paused }
func (r *RTwP) CurrentFaction() int   { return -1 }
func (r *RTwP) TurnNumber() int       { return r.turnNumber }
func (r *RTwP) SetPhase(phase Phase) { r.phase = phase }
func (r *RTwP) RequestEndTurn()      {} // no-op

func (r *RTwP) RequestPause() {
	r.paused = !r.paused
}

// IsPaused returns the current pause state.
func (r *RTwP) IsPaused() bool {
	return r.paused
}

func (r *RTwP) Advance() {
	if r.phase != PhaseCollecting || r.paused {
		return
	}
	r.tickCount++
	if r.ticksPerTurn > 0 && r.tickCount%r.ticksPerTurn == 0 {
		r.turnNumber++
	}
}
