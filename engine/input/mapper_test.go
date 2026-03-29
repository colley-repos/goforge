package input

import "testing"

func TestMapperBindAndProcess(t *testing.T) {
	m := NewMapper()
	m.Bind("keyboard", 87, ActionMoveUp)   // W
	m.Bind("keyboard", 83, ActionMoveDown) // S
	m.Bind("keyboard", 32, ActionConfirm)  // Space

	m.BeginFrame()

	// Press W
	action := m.ProcessInput(RawInput{Source: "keyboard", Code: 87, Pressed: true})
	if action != ActionMoveUp {
		t.Errorf("expected MoveUp, got %d", action)
	}
	if !m.IsActive(ActionMoveUp) {
		t.Error("MoveUp should be active")
	}
	if !m.IsJustPressed(ActionMoveUp) {
		t.Error("MoveUp should be just pressed")
	}

	// Unknown key
	action = m.ProcessInput(RawInput{Source: "keyboard", Code: 999, Pressed: true})
	if action != ActionNone {
		t.Errorf("unknown key should map to None, got %d", action)
	}

	// New frame — justPressed should clear
	m.BeginFrame()
	if m.IsJustPressed(ActionMoveUp) {
		t.Error("justPressed should clear on BeginFrame")
	}
	if !m.IsActive(ActionMoveUp) {
		t.Error("MoveUp should still be active (held)")
	}

	// Release W
	action = m.ProcessInput(RawInput{Source: "keyboard", Code: 87, Pressed: false})
	if action != ActionMoveUp {
		t.Errorf("release should still return the mapped action")
	}
	if m.IsActive(ActionMoveUp) {
		t.Error("MoveUp should no longer be active")
	}
	if !m.IsJustReleased(ActionMoveUp) {
		t.Error("MoveUp should be just released")
	}
}

func TestMapperBindAll(t *testing.T) {
	m := NewMapper()
	m.BindAll([]Binding{
		{Source: "keyboard", Code: 65, Action: ActionMoveLeft},
		{Source: "keyboard", Code: 68, Action: ActionMoveRight},
		{Source: "gamepad", Code: 0, Action: ActionConfirm},
	})

	m.BeginFrame()

	a := m.ProcessInput(RawInput{Source: "keyboard", Code: 65, Pressed: true})
	if a != ActionMoveLeft {
		t.Errorf("expected MoveLeft, got %d", a)
	}

	a = m.ProcessInput(RawInput{Source: "gamepad", Code: 0, Pressed: true})
	if a != ActionConfirm {
		t.Errorf("expected Confirm from gamepad, got %d", a)
	}
}

func TestMapperActiveActions(t *testing.T) {
	m := NewMapper()
	m.Bind("keyboard", 1, ActionMoveUp)
	m.Bind("keyboard", 2, ActionAttack)

	m.BeginFrame()
	m.ProcessInput(RawInput{Source: "keyboard", Code: 1, Pressed: true})
	m.ProcessInput(RawInput{Source: "keyboard", Code: 2, Pressed: true})

	active := m.ActiveActions()
	if len(active) != 2 {
		t.Fatalf("expected 2 active actions, got %d", len(active))
	}
}
