## HUDBase — base class for all heads-up display implementations.
##
## Mirrors Unreal's HUD class: owned and instantiated by GameMode.
## It renders the UI overlay — ability bars, resource bars, modals.
## HUD has NO game logic. It reads state and displays it.
##
## Rules:
##   - Never mutate game state from the HUD. Emit signals, let controllers handle it.
##   - All layout constants belong here or in a shared style resource, not inline.
##   - Modal stacking: use push_modal / pop_modal so screens layer correctly.
##
## Usage:
##   Extend this class and build your concrete HUD scene on top.
##   Wire the HUD to GameState signals so it updates reactively.
##
## Example:
##   class_name BossHUD extends HUDBase
##
##   func _ready() -> void:
##       super._ready()
##       GameState.boss_hp_changed.connect(_on_boss_hp_changed)
##
##   func _on_boss_hp_changed(value: float) -> void:
##       $BossHPBar.value = value * 100.0
class_name HUDBase
extends CanvasLayer

## Emitted when the player requests a pause.
signal pause_requested

## Stack of currently open modal scenes.
var _modal_stack: Array[Control] = []


func _ready() -> void:
	layer = 10  # HUD renders above world geometry


## Push a modal onto the stack. The modal is added as a child and shown.
## Modals are popped in reverse order (LIFO).
func push_modal(modal: Control) -> void:
	_modal_stack.push_back(modal)
	add_child(modal)


## Pop and free the top modal.
func pop_modal() -> void:
	if _modal_stack.is_empty():
		return
	var top: Control = _modal_stack.pop_back()
	top.queue_free()


## Show a result modal with the given outcome string.
## Override to customize the victory/defeat presentation.
func show_result_modal(_outcome: String) -> void:
	pass


## Flash a brief notification message. Override to use your toast system.
func show_toast(_message: String) -> void:
	pass
