extends Node
class_name GameManager

## GameManager autoload — lightweight session-level bookkeeping.
##
## GameManager is NOT the game loop. It holds references that multiple
## systems need to find each other: the active GameMode, player list,
## session flags.
##
## Rule: if logic belongs to a match, put it in GameMode or GameState.
##       If it's a reference that survives scene changes, put it here.
##
## Add as Autoload named "GameManager" in Project Settings.

## The currently active GameMode node.
var current_mode: GameMode = null

## Whether a match is currently running.
var match_active: bool = false

## Session-level save data (gold, upgrades, etc.).
## Populate after loading from SaveManager.
var save_data: Dictionary = {}


func _ready() -> void:
	add_to_group("game_manager")


## Called by GameMode when a match begins.
func on_match_start(mode: GameMode) -> void:
	current_mode = mode
	match_active = true


## Called by GameMode when a match ends.
func on_match_end() -> void:
	match_active = false


## Convenience accessor — returns the active GameState, or null.
func get_game_state() -> GameState:
	if current_mode:
		return current_mode.game_state
	return null
