package commands

// Result is returned by a Resolver after processing a command.
// It carries the outcome and an audit trail for debugging, networking, and replay.
type Result struct {
	// Command that was resolved
	Command Command
	// Success indicates whether the command was valid and executed
	Success bool
	// Reason explains failure (empty on success)
	Reason string
	// Details holds resolver-specific outcome data (e.g., combat results)
	Details map[string]any
}

// Resolver processes commands and returns results.
// Games implement this interface to define how commands mutate game state.
type Resolver interface {
	// Resolve processes a single command and returns the result.
	// The resolver may read and mutate the ECS world.
	Resolve(cmd Command) Result
}

// ResolverFunc is an adapter to allow ordinary functions as Resolvers.
type ResolverFunc func(cmd Command) Result

// Resolve calls the function.
func (f ResolverFunc) Resolve(cmd Command) Result {
	return f(cmd)
}

// DispatchResolver routes commands to type-specific handlers.
// Register a handler per command type; unhandled types return failure.
type DispatchResolver struct {
	handlers map[Type]Resolver
}

// NewDispatchResolver creates a resolver with no handlers registered.
func NewDispatchResolver() *DispatchResolver {
	return &DispatchResolver{
		handlers: make(map[Type]Resolver),
	}
}

// Register associates a Resolver with a command type.
func (d *DispatchResolver) Register(t Type, r Resolver) {
	d.handlers[t] = r
}

// RegisterFunc is a convenience method to register a function as a handler.
func (d *DispatchResolver) RegisterFunc(t Type, fn func(Command) Result) {
	d.handlers[t] = ResolverFunc(fn)
}

// Resolve dispatches to the registered handler for the command's type.
func (d *DispatchResolver) Resolve(cmd Command) Result {
	handler, ok := d.handlers[cmd.Type]
	if !ok {
		return Result{
			Command: cmd,
			Success: false,
			Reason:  "no handler registered for command type: " + cmd.Type.String(),
		}
	}
	return handler.Resolve(cmd)
}
