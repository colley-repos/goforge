// Package commands defines the Command data type, CommandQueue, and resolution interfaces.
//
// Commands are pure data objects — they describe an intent (move, attack, wait, etc.)
// but never execute themselves. The CommandQueue stores them, and the PacingController
// decides when to resolve them.
package commands

import (
	"fmt"
)

// Type identifies the kind of command.
type Type int

const (
	TypeMove     Type = iota // Move entity to a target position
	TypeAttack               // Attack a target entity
	TypeWait                 // Skip / end turn without acting
	TypeAbility              // Use a special ability
	TypeOverwatch            // Enter overwatch stance
	TypeInteract             // Interact with environment
	TypeCustom               // Game-specific extension point
)

// String returns the command type name.
func (t Type) String() string {
	switch t {
	case TypeMove:
		return "Move"
	case TypeAttack:
		return "Attack"
	case TypeWait:
		return "Wait"
	case TypeAbility:
		return "Ability"
	case TypeOverwatch:
		return "Overwatch"
	case TypeInteract:
		return "Interact"
	case TypeCustom:
		return "Custom"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}

// Command is a pure data object representing a game action.
// Commands are created by input handlers or AI, queued in the CommandQueue,
// and resolved when the PacingController allows it.
type Command struct {
	// Type of command
	Type Type
	// SourceEntity is the entity performing the action
	SourceEntity uint64
	// TargetEntity is the entity being acted upon (0 = no target)
	TargetEntity uint64
	// TargetX, TargetY for positional commands (move destination, ability target)
	TargetX float64
	TargetY float64
	// Path for move commands (list of waypoints)
	Path [][2]float64
	// AbilityID for ability commands
	AbilityID string
	// APCost is how many action points this command costs
	APCost int
	// Payload for game-specific extension data
	Payload map[string]any
}

// NewMove creates a move command with a path of waypoints.
func NewMove(source uint64, path [][2]float64, apCost int) Command {
	var targetX, targetY float64
	if len(path) > 0 {
		last := path[len(path)-1]
		targetX = last[0]
		targetY = last[1]
	}
	return Command{
		Type:         TypeMove,
		SourceEntity: source,
		TargetX:      targetX,
		TargetY:      targetY,
		Path:         path,
		APCost:       apCost,
	}
}

// NewAttack creates an attack command targeting another entity.
func NewAttack(source, target uint64, apCost int) Command {
	return Command{
		Type:         TypeAttack,
		SourceEntity: source,
		TargetEntity: target,
		APCost:       apCost,
	}
}

// NewWait creates a wait/skip command.
func NewWait(source uint64) Command {
	return Command{
		Type:         TypeWait,
		SourceEntity: source,
		APCost:       0,
	}
}

// NewOverwatch creates an overwatch command (costs all remaining AP).
func NewOverwatch(source uint64, apCost int) Command {
	return Command{
		Type:         TypeOverwatch,
		SourceEntity: source,
		APCost:       apCost,
	}
}

// NewAbility creates an ability command.
func NewAbility(source uint64, abilityID string, targetX, targetY float64, apCost int) Command {
	return Command{
		Type:         TypeAbility,
		SourceEntity: source,
		AbilityID:    abilityID,
		TargetX:      targetX,
		TargetY:      targetY,
		APCost:       apCost,
	}
}
