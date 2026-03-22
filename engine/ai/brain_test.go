package ai

import (
	"testing"

	"github.com/colley-repos/goforge/engine/commands"
)

func TestUtilityAISelectsHighestScore(t *testing.T) {
	ai := &UtilityAI{
		ScoreFunc: func(ctx any, eid uint64) []ScoredAction {
			return []ScoredAction{
				{Command: commands.NewWait(eid), Score: 10},
				{Command: commands.NewAttack(eid, 2, 1), Score: 50},
				{Command: commands.NewMove(eid, [][2]float64{{1, 1}}, 1), Score: 30},
			}
		},
	}

	cmd := ai.Decide(nil, 1)
	if cmd.Type != commands.TypeAttack {
		t.Errorf("expected Attack (highest score), got %s", cmd.Type)
	}
}

func TestUtilityAIEmptyActionsWaits(t *testing.T) {
	ai := &UtilityAI{
		ScoreFunc: func(ctx any, eid uint64) []ScoredAction {
			return nil
		},
	}

	cmd := ai.Decide(nil, 5)
	if cmd.Type != commands.TypeWait {
		t.Errorf("expected Wait fallback, got %s", cmd.Type)
	}
}

func TestBrainFuncAdapter(t *testing.T) {
	var brain Brain = BrainFunc(func(ctx any, eid uint64) commands.Command {
		return commands.NewOverwatch(eid, 2)
	})

	cmd := brain.Decide(nil, 7)
	if cmd.Type != commands.TypeOverwatch {
		t.Errorf("expected Overwatch, got %s", cmd.Type)
	}
}
