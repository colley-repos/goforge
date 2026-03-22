package pacing

import "testing"

func TestTurnBasedFlow(t *testing.T) {
	tb := NewTurnBased([]int{0, 1}) // player=0, enemy=1

	// Starts collecting for faction 0
	if tb.Phase() != PhaseCollecting {
		t.Fatalf("expected Collecting, got %s", tb.Phase())
	}
	if tb.CurrentFaction() != 0 {
		t.Fatalf("expected faction 0, got %d", tb.CurrentFaction())
	}
	if tb.TurnNumber() != 0 {
		t.Fatalf("expected turn 0, got %d", tb.TurnNumber())
	}
	if !tb.ShouldCollect() {
		t.Fatal("should be collecting")
	}
	if tb.ShouldResolve() {
		t.Fatal("should NOT be resolving yet")
	}

	// End turn triggers resolution
	tb.RequestEndTurn()
	if tb.Phase() != PhaseResolving {
		t.Fatalf("expected Resolving, got %s", tb.Phase())
	}
	if !tb.ShouldResolve() {
		t.Fatal("should resolve after end turn")
	}

	// Advance moves to next faction
	tb.Advance()
	if tb.CurrentFaction() != 1 {
		t.Fatalf("expected faction 1, got %d", tb.CurrentFaction())
	}
	if tb.Phase() != PhaseCollecting {
		t.Fatalf("expected Collecting, got %s", tb.Phase())
	}

	// End faction 1's turn, advance → wraps to faction 0, turn increments
	tb.RequestEndTurn()
	tb.Advance()
	if tb.CurrentFaction() != 0 {
		t.Fatalf("expected faction 0 again, got %d", tb.CurrentFaction())
	}
	if tb.TurnNumber() != 1 {
		t.Fatalf("expected turn 1, got %d", tb.TurnNumber())
	}
}

func TestRealTimeAlwaysResolves(t *testing.T) {
	rt := NewRealTime(60) // turn increments every 60 ticks

	if !rt.ShouldCollect() {
		t.Fatal("real-time should always collect")
	}
	if !rt.ShouldResolve() {
		t.Fatal("real-time should always resolve")
	}
	if rt.CurrentFaction() != -1 {
		t.Fatalf("real-time has no faction, got %d", rt.CurrentFaction())
	}

	// Advance 60 times → turn increments
	for i := 0; i < 60; i++ {
		rt.Advance()
	}
	if rt.TurnNumber() != 1 {
		t.Fatalf("expected turn 1 after 60 ticks, got %d", rt.TurnNumber())
	}
}

func TestRTwPPauseToggle(t *testing.T) {
	r := NewRTwP(0)

	// Starts paused
	if !r.IsPaused() {
		t.Fatal("RTwP should start paused")
	}
	if r.ShouldResolve() {
		t.Fatal("should NOT resolve while paused")
	}
	if !r.ShouldCollect() {
		t.Fatal("should still collect while paused")
	}

	// Unpause
	r.RequestPause()
	if r.IsPaused() {
		t.Fatal("should be unpaused after toggle")
	}
	if !r.ShouldResolve() {
		t.Fatal("should resolve when unpaused")
	}

	// Pause again
	r.RequestPause()
	if !r.IsPaused() {
		t.Fatal("should be paused after second toggle")
	}
}

func TestSetPhaseGameOver(t *testing.T) {
	tb := NewTurnBased([]int{0, 1})
	tb.SetPhase(PhaseGameOver)

	if tb.Phase() != PhaseGameOver {
		t.Fatalf("expected GameOver, got %s", tb.Phase())
	}
	if tb.ShouldCollect() {
		t.Fatal("should not collect in game over")
	}
}
