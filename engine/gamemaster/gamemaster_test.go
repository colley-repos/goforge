package gamemaster

import (
	"testing"

	"github.com/colley-repos/goforge/engine/commands"
	"github.com/colley-repos/goforge/engine/events"
	"github.com/colley-repos/goforge/engine/pacing"
)

func TestGameMasterTickFlow(t *testing.T) {
	resolver := commands.NewDispatchResolver()
	resolver.RegisterFunc(commands.TypeWait, func(cmd commands.Command) commands.Result {
		return commands.Result{Command: cmd, Success: true}
	})

	gm := New(Config{
		Pacing:   pacing.NewTurnBased([]int{0, 1}),
		Resolver: resolver,
	})

	// Track events
	var resolvedCount int
	events.Subscribe[events.CommandResolved](gm.EventBus, func(e events.CommandResolved) {
		resolvedCount++
	})

	// Track system calls
	var systemCalled int
	gm.AddSystem(func(ctx *TickContext) {
		systemCalled++
	})

	// Enqueue a command
	gm.CommandQueue.Enqueue(commands.NewWait(1))

	// Tick without end-turn → should NOT resolve (turn-based, still collecting)
	gm.Tick(1.0 / 60.0)
	if resolvedCount != 0 {
		t.Fatalf("should not resolve during collection phase, got %d", resolvedCount)
	}
	if systemCalled != 1 {
		t.Fatalf("system should have been called once, got %d", systemCalled)
	}

	// End turn → next tick should resolve
	gm.Pacing.RequestEndTurn()
	gm.Tick(1.0 / 60.0)

	// Events are drained in the same tick they're published
	// But the event was published during resolve, then drained at step 4
	// The subscriber was called during drain
	if resolvedCount != 1 {
		t.Fatalf("expected 1 resolved event, got %d", resolvedCount)
	}
}

func TestGameMasterRealTimeResolves(t *testing.T) {
	resolver := commands.NewDispatchResolver()
	resolver.RegisterFunc(commands.TypeAttack, func(cmd commands.Command) commands.Result {
		return commands.Result{Command: cmd, Success: true}
	})

	gm := New(Config{
		Pacing:   pacing.NewRealTime(0),
		Resolver: resolver,
	})

	var resolved int
	events.Subscribe[events.CommandResolved](gm.EventBus, func(e events.CommandResolved) {
		resolved++
	})

	gm.CommandQueue.Enqueue(commands.NewAttack(1, 2, 1))
	gm.Tick(0.016)

	if resolved != 1 {
		t.Fatalf("real-time should resolve immediately, got %d resolved", resolved)
	}
}

func TestGameMasterTickCount(t *testing.T) {
	gm := New(Config{
		Pacing: pacing.NewRealTime(0),
	})

	if gm.TickCount() != 0 {
		t.Fatalf("initial tick count should be 0")
	}

	gm.Tick(0.016)
	gm.Tick(0.016)
	gm.Tick(0.016)

	if gm.TickCount() != 3 {
		t.Fatalf("expected 3, got %d", gm.TickCount())
	}
}
