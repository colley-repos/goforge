// Package events provides a typed, synchronous event bus.
//
// Events are appended to a slice during the frame and drained synchronously
// at a defined point in the game loop. This preserves determinism — no channels,
// no goroutines, no race conditions.
//
// Usage:
//
//	bus := events.NewBus()
//	events.Subscribe[DamageEvent](bus, func(e DamageEvent) {
//	    fmt.Printf("%s took %d damage\n", e.TargetName, e.Amount)
//	})
//	events.Publish(bus, DamageEvent{TargetName: "Goblin", Amount: 5})
//	bus.Drain() // handlers fire here, in order
package events

import (
	"fmt"
	"reflect"
)

// Bus is the central event dispatcher. Systems publish events to it,
// and subscribers receive them when Drain() is called.
type Bus struct {
	// subscribers maps event type name → slice of handler wrappers
	subscribers map[string][]handlerWrapper
	// pending events waiting to be drained
	pending []pendingEvent
}

// handlerWrapper stores a type-erased handler function.
type handlerWrapper struct {
	fn func(any)
}

// pendingEvent stores a type-erased event waiting to be dispatched.
type pendingEvent struct {
	typeName string
	value    any
}

// NewBus creates a new event bus with no subscribers.
func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string][]handlerWrapper),
		pending:     make([]pendingEvent, 0, 64),
	}
}

// Subscribe registers a handler for events of type T.
// Handlers are called in subscription order during Drain().
func Subscribe[T any](bus *Bus, handler func(T)) {
	typeName := typeKey[T]()
	wrapped := handlerWrapper{
		fn: func(v any) {
			handler(v.(T))
		},
	}
	bus.subscribers[typeName] = append(bus.subscribers[typeName], wrapped)
}

// Publish queues an event for delivery on the next Drain().
// Events are NOT delivered immediately — this preserves frame-synchronous determinism.
func Publish[T any](bus *Bus, event T) {
	typeName := typeKey[T]()
	bus.pending = append(bus.pending, pendingEvent{
		typeName: typeName,
		value:    event,
	})
}

// Drain processes all pending events, calling registered handlers in order.
// Events published during drain are queued for the NEXT drain (no infinite loops).
// Returns the number of events processed.
func (b *Bus) Drain() int {
	// Snapshot current pending slice, reset for any events published during drain
	toProcess := b.pending
	b.pending = make([]pendingEvent, 0, cap(toProcess))

	for _, evt := range toProcess {
		handlers, ok := b.subscribers[evt.typeName]
		if !ok {
			continue
		}
		for _, h := range handlers {
			h.fn(evt.value)
		}
	}

	return len(toProcess)
}

// PendingCount returns how many events are waiting to be drained.
func (b *Bus) PendingCount() int {
	return len(b.pending)
}

// Clear removes all pending events without processing them.
func (b *Bus) Clear() {
	b.pending = b.pending[:0]
}

// Reset removes all subscribers and pending events.
func (b *Bus) Reset() {
	b.subscribers = make(map[string][]handlerWrapper)
	b.pending = b.pending[:0]
}

// typeKey returns a unique string key for a Go type, used for event dispatch.
func typeKey[T any]() string {
	var zero T
	return fmt.Sprintf("%v", reflect.TypeOf(zero))
}
