// Package gamemaster provides the root game loop orchestrator.
//
// The GameMaster owns the ECS world, event bus, command queue, and pacing controller.
// It wires all systems together and drives the main simulation loop.
// This is the equivalent of Dystopia's Main.gd — a wiring script, not gameplay logic.
package gamemaster

import (
	"github.com/colley-repos/goforge/engine/commands"
	"github.com/colley-repos/goforge/engine/events"
	ecslib "github.com/colley-repos/goforge/engine/ecs"
	"github.com/colley-repos/goforge/engine/input"
	"github.com/colley-repos/goforge/engine/pacing"
)

// System is a function that operates on the game state each tick.
// Systems are registered with the GameMaster and called in order.
type System func(ctx *TickContext)

// TickContext provides systems with access to the game state for one tick.
type TickContext struct {
	World        *ecslib.World
	EventBus     *events.Bus
	CommandQueue *commands.Queue
	Pacing       pacing.Controller
	InputMapper  *input.Mapper
	DeltaTime    float64 // seconds since last tick
	TickCount    uint64
}

// GameMaster is the root orchestrator for a GoForge game.
type GameMaster struct {
	World        *ecslib.World
	EventBus     *events.Bus
	CommandQueue *commands.Queue
	Pacing       pacing.Controller
	InputMapper  *input.Mapper
	Resolver     commands.Resolver

	systems   []System
	tickCount uint64
}

// Config holds initialization parameters for the GameMaster.
type Config struct {
	Pacing      pacing.Controller
	Resolver    commands.Resolver
	InputMapper *input.Mapper
}

// New creates a new GameMaster with the given configuration.
func New(cfg Config) *GameMaster {
	world := ecslib.NewWorld()

	gm := &GameMaster{
		World:        world,
		EventBus:     events.NewBus(),
		CommandQueue: commands.NewQueue(),
		Pacing:       cfg.Pacing,
		InputMapper:  cfg.InputMapper,
		Resolver:     cfg.Resolver,
		systems:      make([]System, 0, 16),
	}

	if gm.InputMapper == nil {
		gm.InputMapper = input.NewMapper()
	}

	return gm
}

// AddSystem registers a system to be called each tick, in registration order.
func (gm *GameMaster) AddSystem(sys System) {
	gm.systems = append(gm.systems, sys)
}

// Tick advances the simulation by one step.
// This is called by the renderer's update loop (or manually in tests).
//
// Order of operations:
//  1. Input mapper begins frame
//  2. Run registered systems (AI, game logic, etc.)
//  3. If pacing says resolve → drain command queue through resolver
//  4. Drain event bus (deliver all queued events)
//  5. Advance pacing controller
func (gm *GameMaster) Tick(dt float64) {
	gm.tickCount++

	ctx := &TickContext{
		World:        gm.World,
		EventBus:     gm.EventBus,
		CommandQueue: gm.CommandQueue,
		Pacing:       gm.Pacing,
		InputMapper:  gm.InputMapper,
		DeltaTime:    dt,
		TickCount:    gm.tickCount,
	}

	// 1. Begin input frame
	gm.InputMapper.BeginFrame()

	// 2. Run systems
	for _, sys := range gm.systems {
		sys(ctx)
	}

	// 3. Resolve commands if pacing allows
	if gm.Pacing.ShouldResolve() && gm.Resolver != nil {
		cmds := gm.CommandQueue.Drain()
		for _, cmd := range cmds {
			result := gm.Resolver.Resolve(cmd)
			events.Publish(gm.EventBus, events.CommandResolved{
				CommandType: cmd.Type.String(),
				EntityID:    cmd.SourceEntity,
				Success:     result.Success,
			})
		}
		if len(cmds) > 0 {
			events.Publish(gm.EventBus, events.AllCommandsResolved{})
		}
	}

	// 4. Drain events
	gm.EventBus.Drain()

	// 5. Advance pacing
	gm.Pacing.Advance()
}

// ProcessInput feeds a raw input through the input mapper.
// Typically called by the renderer when it detects key/mouse/touch events.
func (gm *GameMaster) ProcessInput(raw input.RawInput) input.Action {
	return gm.InputMapper.ProcessInput(raw)
}

// TickCount returns how many ticks have elapsed.
func (gm *GameMaster) TickCount() uint64 {
	return gm.tickCount
}
