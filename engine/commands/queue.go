package commands

// Queue stores commands and provides access for resolution.
// Commands are appended by input handlers / AI and drained by the resolver.
type Queue struct {
	commands []Command
}

// NewQueue creates an empty command queue.
func NewQueue() *Queue {
	return &Queue{
		commands: make([]Command, 0, 16),
	}
}

// Enqueue adds a command to the end of the queue.
func (q *Queue) Enqueue(cmd Command) {
	q.commands = append(q.commands, cmd)
}

// Peek returns the next command without removing it, or nil if empty.
func (q *Queue) Peek() *Command {
	if len(q.commands) == 0 {
		return nil
	}
	cmd := q.commands[0]
	return &cmd
}

// Dequeue removes and returns the next command, or nil if empty.
func (q *Queue) Dequeue() *Command {
	if len(q.commands) == 0 {
		return nil
	}
	cmd := q.commands[0]
	q.commands = q.commands[1:]
	return &cmd
}

// Drain removes and returns all queued commands.
func (q *Queue) Drain() []Command {
	result := q.commands
	q.commands = make([]Command, 0, cap(result))
	return result
}

// Len returns the number of queued commands.
func (q *Queue) Len() int {
	return len(q.commands)
}

// Clear removes all queued commands without returning them.
func (q *Queue) Clear() {
	q.commands = q.commands[:0]
}
