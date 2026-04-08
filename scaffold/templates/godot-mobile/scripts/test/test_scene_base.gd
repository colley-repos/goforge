## TestSceneBase — base for isolated subsystem test scenes.
##
## Every complex subsystem (combat, boss encounter, unit AI) gets its own
## test scene that instantiates only what it needs. This is the Godot
## equivalent of a unit test harness.
##
## Pattern (from RaidBoss BossTest):
##   - Set game phase BEFORE instantiating HUD — HUD reads phase to build layout
##   - Spawn SFXManager yourself if the test scene doesn't have the main autoload
##   - Use keyboard shortcuts to trigger abilities/states without UI
##   - Print to Output, not to the HUD
##
## Usage:
##   Extend this class. Override _setup() and _register_shortcuts().
##   Run directly: Godot --path <project> --scene res://scenes/test/MyTest.tscn
class_name TestSceneBase
extends Node

## Shortcut map: key scancode → label shown in console on startup.
var _shortcuts: Dictionary = {}


func _ready() -> void:
	_setup()
	_register_shortcuts()
	_print_shortcuts()


## Override — instantiate the subsystem under test here.
## Rule: set all state (phase, config) BEFORE adding children.
func _setup() -> void:
	pass


## Override — call _bind_shortcut() for each test action.
func _register_shortcuts() -> void:
	pass


## Bind a key to a label for the startup print.
func _bind_shortcut(key: Key, label: String) -> void:
	_shortcuts[key] = label


func _unhandled_key_input(event: InputEventKey) -> void:
	if not event.pressed:
		return
	_on_key(event.keycode)


## Override — handle key presses for test controls.
func _on_key(_key: Key) -> void:
	pass


func _print_shortcuts() -> void:
	print("=== %s Test Controls ===" % name)
	for key: Key in _shortcuts:
		print("  [%s] %s" % [OS.get_keycode_string(key), _shortcuts[key]])
	print("========================")
