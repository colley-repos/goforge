package commands

import "testing"

func TestQueueEnqueueDequeue(t *testing.T) {
	q := NewQueue()

	if q.Len() != 0 {
		t.Fatalf("new queue should be empty, got len %d", q.Len())
	}

	q.Enqueue(NewMove(1, [][2]float64{{5, 5}}, 1))
	q.Enqueue(NewAttack(1, 2, 1))
	q.Enqueue(NewWait(1))

	if q.Len() != 3 {
		t.Fatalf("expected len 3, got %d", q.Len())
	}

	// Peek doesn't remove
	peeked := q.Peek()
	if peeked == nil || peeked.Type != TypeMove {
		t.Fatalf("Peek: expected Move, got %v", peeked)
	}
	if q.Len() != 3 {
		t.Fatalf("Peek should not remove, len = %d", q.Len())
	}

	// Dequeue removes in order
	cmd := q.Dequeue()
	if cmd.Type != TypeMove {
		t.Errorf("first dequeue: expected Move, got %s", cmd.Type)
	}
	cmd = q.Dequeue()
	if cmd.Type != TypeAttack {
		t.Errorf("second dequeue: expected Attack, got %s", cmd.Type)
	}
	cmd = q.Dequeue()
	if cmd.Type != TypeWait {
		t.Errorf("third dequeue: expected Wait, got %s", cmd.Type)
	}

	// Empty after all dequeued
	if q.Dequeue() != nil {
		t.Error("dequeue from empty queue should return nil")
	}
}

func TestQueueDrain(t *testing.T) {
	q := NewQueue()
	q.Enqueue(NewMove(1, nil, 1))
	q.Enqueue(NewAttack(2, 3, 1))

	drained := q.Drain()
	if len(drained) != 2 {
		t.Fatalf("expected 2 drained, got %d", len(drained))
	}
	if q.Len() != 0 {
		t.Fatalf("queue should be empty after drain, got %d", q.Len())
	}
}

func TestDispatchResolver(t *testing.T) {
	dr := NewDispatchResolver()

	dr.RegisterFunc(TypeMove, func(cmd Command) Result {
		return Result{Command: cmd, Success: true, Details: map[string]any{"moved": true}}
	})
	dr.RegisterFunc(TypeAttack, func(cmd Command) Result {
		return Result{Command: cmd, Success: true, Details: map[string]any{"damage": 5}}
	})

	// Known types resolve
	moveResult := dr.Resolve(NewMove(1, [][2]float64{{3, 4}}, 1))
	if !moveResult.Success {
		t.Error("move should succeed")
	}

	attackResult := dr.Resolve(NewAttack(1, 2, 1))
	if !attackResult.Success {
		t.Error("attack should succeed")
	}

	// Unregistered type fails gracefully
	waitResult := dr.Resolve(NewWait(1))
	if waitResult.Success {
		t.Error("wait should fail (no handler registered)")
	}
	if waitResult.Reason == "" {
		t.Error("failure should include reason")
	}
}

func TestCommandFactories(t *testing.T) {
	move := NewMove(10, [][2]float64{{1, 2}, {3, 4}}, 2)
	if move.Type != TypeMove || move.SourceEntity != 10 || move.APCost != 2 {
		t.Errorf("NewMove: %+v", move)
	}
	if move.TargetX != 3 || move.TargetY != 4 {
		t.Errorf("NewMove target should be last waypoint: got (%v,%v)", move.TargetX, move.TargetY)
	}

	attack := NewAttack(1, 2, 1)
	if attack.Type != TypeAttack || attack.TargetEntity != 2 {
		t.Errorf("NewAttack: %+v", attack)
	}

	wait := NewWait(5)
	if wait.Type != TypeWait || wait.APCost != 0 {
		t.Errorf("NewWait: %+v", wait)
	}

	ow := NewOverwatch(7, 3)
	if ow.Type != TypeOverwatch || ow.APCost != 3 {
		t.Errorf("NewOverwatch: %+v", ow)
	}

	ability := NewAbility(1, "fireball", 5, 5, 2)
	if ability.Type != TypeAbility || ability.AbilityID != "fireball" {
		t.Errorf("NewAbility: %+v", ability)
	}
}
