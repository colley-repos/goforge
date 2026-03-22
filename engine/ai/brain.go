// Package ai provides interfaces for AI decision-making.
//
// AI brains receive the game state and an entity, and return a Command
// describing what that entity should do. This keeps AI completely decoupled
// from input handling and rendering.
package ai

import (
	"github.com/colley-repos/goforge/engine/commands"
)

// Brain is the interface for AI decision-makers.
// Implementations receive game context and return a command for the entity.
type Brain interface {
	// Decide returns the command this entity should execute.
	// context is a game-specific state object; the brain casts it as needed.
	Decide(context any, entityID uint64) commands.Command
}

// BrainFunc is an adapter to allow ordinary functions as Brain implementations.
type BrainFunc func(context any, entityID uint64) commands.Command

// Decide calls the function.
func (f BrainFunc) Decide(context any, entityID uint64) commands.Command {
	return f(context, entityID)
}

// ScoredAction pairs a command with a utility score for UtilityAI.
type ScoredAction struct {
	Command commands.Command
	Score   float64
}

// UtilityAI selects the action with the highest utility score.
// Games provide a scoring function that evaluates all possible actions.
type UtilityAI struct {
	// ScoreFunc evaluates possible actions and returns scored options.
	ScoreFunc func(context any, entityID uint64) []ScoredAction
}

// Decide picks the highest-scoring action.
func (u *UtilityAI) Decide(context any, entityID uint64) commands.Command {
	scored := u.ScoreFunc(context, entityID)
	if len(scored) == 0 {
		return commands.NewWait(entityID)
	}

	best := scored[0]
	for _, s := range scored[1:] {
		if s.Score > best.Score {
			best = s
		}
	}
	return best.Command
}
