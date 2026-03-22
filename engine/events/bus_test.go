package events

import (
	"testing"
)

func TestPublishSubscribeDrain(t *testing.T) {
	bus := NewBus()

	type TestEvent struct {
		Value int
	}

	var received []int
	Subscribe[TestEvent](bus, func(e TestEvent) {
		received = append(received, e.Value)
	})

	// Publish does not deliver immediately
	Publish(bus, TestEvent{Value: 1})
	Publish(bus, TestEvent{Value: 2})
	Publish(bus, TestEvent{Value: 3})

	if len(received) != 0 {
		t.Fatalf("events should not fire before Drain, got %d", len(received))
	}
	if bus.PendingCount() != 3 {
		t.Fatalf("expected 3 pending, got %d", bus.PendingCount())
	}

	// Drain delivers in order
	n := bus.Drain()
	if n != 3 {
		t.Fatalf("expected Drain to return 3, got %d", n)
	}
	if len(received) != 3 {
		t.Fatalf("expected 3 received, got %d", len(received))
	}
	for i, v := range received {
		if v != i+1 {
			t.Errorf("received[%d] = %d, want %d", i, v, i+1)
		}
	}

	// After drain, pending is empty
	if bus.PendingCount() != 0 {
		t.Fatalf("pending should be 0 after drain, got %d", bus.PendingCount())
	}
}

func TestMultipleSubscribers(t *testing.T) {
	bus := NewBus()

	type Ping struct{ ID int }

	var log1, log2 []int
	Subscribe[Ping](bus, func(e Ping) { log1 = append(log1, e.ID) })
	Subscribe[Ping](bus, func(e Ping) { log2 = append(log2, e.ID) })

	Publish(bus, Ping{ID: 42})
	bus.Drain()

	if len(log1) != 1 || log1[0] != 42 {
		t.Errorf("subscriber 1: got %v, want [42]", log1)
	}
	if len(log2) != 1 || log2[0] != 42 {
		t.Errorf("subscriber 2: got %v, want [42]", log2)
	}
}

func TestDifferentEventTypes(t *testing.T) {
	bus := NewBus()

	type Alpha struct{ A int }
	type Beta struct{ B string }

	var gotA int
	var gotB string
	Subscribe[Alpha](bus, func(e Alpha) { gotA = e.A })
	Subscribe[Beta](bus, func(e Beta) { gotB = e.B })

	Publish(bus, Alpha{A: 99})
	Publish(bus, Beta{B: "hello"})
	bus.Drain()

	if gotA != 99 {
		t.Errorf("Alpha: got %d, want 99", gotA)
	}
	if gotB != "hello" {
		t.Errorf("Beta: got %q, want %q", gotB, "hello")
	}
}

func TestEventsPublishedDuringDrainQueuedForNext(t *testing.T) {
	bus := NewBus()

	type Step struct{ N int }

	var order []int

	// First subscriber publishes a new event during drain
	Subscribe[Step](bus, func(e Step) {
		order = append(order, e.N)
		if e.N == 1 {
			Publish(bus, Step{N: 2})
		}
	})

	Publish(bus, Step{N: 1})
	bus.Drain()

	// Only N=1 should have fired; N=2 is pending for next drain
	if len(order) != 1 || order[0] != 1 {
		t.Fatalf("first drain: got %v, want [1]", order)
	}
	if bus.PendingCount() != 1 {
		t.Fatalf("expected 1 pending after first drain, got %d", bus.PendingCount())
	}

	bus.Drain()
	if len(order) != 2 || order[1] != 2 {
		t.Fatalf("second drain: got %v, want [1, 2]", order)
	}
}

func TestClear(t *testing.T) {
	bus := NewBus()

	type Msg struct{ Text string }
	Publish(bus, Msg{Text: "dropped"})
	bus.Clear()

	if bus.PendingCount() != 0 {
		t.Fatalf("Clear should empty pending, got %d", bus.PendingCount())
	}
}

func TestReset(t *testing.T) {
	bus := NewBus()

	type Msg struct{ Text string }
	var called bool
	Subscribe[Msg](bus, func(e Msg) { called = true })

	bus.Reset()

	Publish(bus, Msg{Text: "after reset"})
	bus.Drain()

	if called {
		t.Fatal("Reset should remove subscribers")
	}
}
