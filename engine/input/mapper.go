// Package input provides semantic action mapping for player input.
//
// Rather than game systems checking raw key codes, they check for semantic actions.
// The InputMapper translates raw inputs (from the renderer) into actions.
// This enables remapping, touch-to-keyboard unification, and testability.
package input

// Action represents a semantic game action (not a raw key).
type Action int

// Common actions used across most game types.
// Games should define their own actions starting from ActionCustom.
const (
	ActionNone       Action = iota
	ActionConfirm           // Accept / select
	ActionCancel            // Back / deselect
	ActionMoveUp            // Move/navigate up
	ActionMoveDown          // Move/navigate down
	ActionMoveLeft          // Move/navigate left
	ActionMoveRight         // Move/navigate right
	ActionAttack            // Primary attack
	ActionAbility           // Use ability
	ActionOverwatch         // Enter overwatch
	ActionEndTurn           // End current turn
	ActionPause             // Pause / unpause
	ActionZoomIn            // Camera zoom in
	ActionZoomOut           // Camera zoom out
	ActionScreenshot        // Take screenshot
	ActionCustom     = 1000 // Game-specific actions start here
)

// RawInput represents a single input event from the renderer.
type RawInput struct {
	// Source identifies the input device: "keyboard", "mouse", "touch", "gamepad"
	Source string
	// Code is the device-specific key/button code
	Code int
	// Pressed is true on press, false on release
	Pressed bool
	// X, Y for positional inputs (mouse position, touch point)
	X float64
	Y float64
}

// Binding maps a raw input to a semantic action.
type Binding struct {
	Source string
	Code   int
	Action Action
}

// Mapper translates raw inputs into semantic actions.
type Mapper struct {
	bindings map[bindingKey]Action
	// activeActions tracks currently pressed actions
	activeActions map[Action]bool
	// justPressed tracks actions that became active this frame
	justPressed map[Action]bool
	// justReleased tracks actions that became inactive this frame
	justReleased map[Action]bool
}

type bindingKey struct {
	source string
	code   int
}

// NewMapper creates a mapper with no bindings.
func NewMapper() *Mapper {
	return &Mapper{
		bindings:      make(map[bindingKey]Action),
		activeActions: make(map[Action]bool),
		justPressed:   make(map[Action]bool),
		justReleased:  make(map[Action]bool),
	}
}

// Bind associates a raw input with a semantic action. Overwrites any existing binding.
func (m *Mapper) Bind(source string, code int, action Action) {
	m.bindings[bindingKey{source, code}] = action
}

// BindAll registers multiple bindings at once.
func (m *Mapper) BindAll(bindings []Binding) {
	for _, b := range bindings {
		m.Bind(b.Source, b.Code, b.Action)
	}
}

// Unbind removes a binding.
func (m *Mapper) Unbind(source string, code int) {
	delete(m.bindings, bindingKey{source, code})
}

// BeginFrame clears per-frame state. Call at the start of each frame.
func (m *Mapper) BeginFrame() {
	clear(m.justPressed)
	clear(m.justReleased)
}

// ProcessInput translates a raw input into an action and updates state.
// Returns the mapped action (ActionNone if no binding exists).
func (m *Mapper) ProcessInput(raw RawInput) Action {
	action, ok := m.bindings[bindingKey{raw.Source, raw.Code}]
	if !ok {
		return ActionNone
	}

	if raw.Pressed {
		if !m.activeActions[action] {
			m.justPressed[action] = true
		}
		m.activeActions[action] = true
	} else {
		if m.activeActions[action] {
			m.justReleased[action] = true
		}
		m.activeActions[action] = false
	}

	return action
}

// IsActive returns true if the action is currently held down.
func (m *Mapper) IsActive(action Action) bool {
	return m.activeActions[action]
}

// IsJustPressed returns true if the action was pressed THIS frame.
func (m *Mapper) IsJustPressed(action Action) bool {
	return m.justPressed[action]
}

// IsJustReleased returns true if the action was released THIS frame.
func (m *Mapper) IsJustReleased(action Action) bool {
	return m.justReleased[action]
}

// ActiveActions returns all currently active actions.
func (m *Mapper) ActiveActions() []Action {
	result := make([]Action, 0, len(m.activeActions))
	for a, active := range m.activeActions {
		if active {
			result = append(result, a)
		}
	}
	return result
}
