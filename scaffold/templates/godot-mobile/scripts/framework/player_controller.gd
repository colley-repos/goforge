## PlayerController — translates player input into game actions.
##
## Mirrors Unreal's PlayerController: it owns the input layer.
## The pawn/character has NO input code — the controller feeds it.
##
## This separation means:
##   - AI can drive any pawn by replacing the PlayerController with an AIController
##   - Input remapping lives in one place, not scattered across units
##   - Spectator mode = swap controller, not the pawn
##
## Usage:
##   Extend this class, connect to input signals, call actions on pawn/game state.
##
## Example:
##   class_name BossPlayerController extends PlayerController
##
##   func _unhandled_input(event: InputEvent) -> void:
##       if event.is_action_pressed("ability_1"):
##           _game_mode.game_state.queue_boss_ability("cleave")
class_name PlayerController
extends Node

## Reference to the GameMode that owns this controller.
## Set automatically by GameMode._spawn_player_controller().
var _game_mode: GameMode

## The pawn this controller is currently possessing.
var possessed_pawn: Node


## Possess a pawn — this controller now drives it.
func possess(pawn: Node) -> void:
	possessed_pawn = pawn
	on_possess(pawn)


## Called after possess(). Override to wire input to the new pawn.
func on_possess(_pawn: Node) -> void:
	pass


## Release the possessed pawn.
func unpossess() -> void:
	if possessed_pawn:
		on_unpossess(possessed_pawn)
	possessed_pawn = null


## Called before unpossess clears the reference.
func on_unpossess(_pawn: Node) -> void:
	pass
