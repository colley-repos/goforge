## GameMode — the wiring harness for a game type.
##
## Mirrors Unreal Engine's GameMode philosophy:
##   Swap the GameMode → swap the entire game type.
##   Same engine, completely different rules.
##
## GameMode owns:
##   - Which classes fill each role (HUD, controller, game state, player state)
##   - Match lifecycle (start, end, player join/leave)
##   - Spawn logic (where, what team)
##
## Usage:
##   1. Extend GameMode in your concrete mode script.
##   2. Export your scene assignments in the Inspector.
##   3. Set this node as your root scene — it spawns everything else.
##   4. Override only the lifecycle methods you need.
##
## Example:
##   class_name RaidBossMode extends GameMode
##
##   func match_start() -> void:
##       game_state.spawn_heroes(4)
##       hud.show_encounter_ui()
class_name GameMode
extends Node

# ---------------------------------------------------------------------------
# Slot declarations — assign these in the Inspector or override get_*_scene()
# ---------------------------------------------------------------------------

## The scene that represents what the player embodies or controls.
## In a strategy game: a camera pawn. In an action game: a character.
@export var default_pawn_scene: PackedScene

## The HUD scene. Instantiated once and added as a CanvasLayer child.
@export var hud_scene: PackedScene

## The PlayerController scene. Handles input → action translation.
@export var player_controller_scene: PackedScene

## The GameState scene. Owns shared truth visible to all players.
@export var game_state_scene: PackedScene

## The PlayerState scene. One instance per player, owns per-player data.
@export var player_state_scene: PackedScene

# ---------------------------------------------------------------------------
# Live references — set during _ready(), read by game logic
# ---------------------------------------------------------------------------

var game_state: GameState
var hud: HUDBase
var player_controller: PlayerController
var default_pawn: Node

# ---------------------------------------------------------------------------
# Lifecycle — _ready() wires everything, then calls match_start()
# ---------------------------------------------------------------------------

func _ready() -> void:
	_spawn_game_state()
	_spawn_player_controller()
	_spawn_hud()
	_spawn_default_pawn()
	match_start()


# ---------------------------------------------------------------------------
# Override these in your concrete GameMode
# ---------------------------------------------------------------------------

## Called once after all slots are spawned and wired.
## Spawn enemies, start encounter timers, show intro sequence here.
func match_start() -> void:
	pass


## Called when the match concludes.
## winner: "player" | "enemy" | "draw"
func match_end(winner: String) -> void:
	pass


## Called when a player joins. Spawn their pawn and player state.
func on_player_join(player_id: int) -> void:
	pass


## Called when a player leaves. Clean up their entities.
func on_player_leave(player_id: int) -> void:
	pass


## Return true if this player should respawn after death.
func should_respawn(_player_id: int) -> bool:
	return false


## Return the world-space spawn position for a player.
func get_spawn_position(_player_id: int) -> Vector3:
	return Vector3.ZERO


# ---------------------------------------------------------------------------
# Internal wiring — do not call directly
# ---------------------------------------------------------------------------

func _spawn_game_state() -> void:
	if game_state_scene:
		game_state = game_state_scene.instantiate() as GameState
		add_child(game_state)


func _spawn_player_controller() -> void:
	if player_controller_scene:
		player_controller = player_controller_scene.instantiate() as PlayerController
		add_child(player_controller)


func _spawn_hud() -> void:
	if hud_scene:
		hud = hud_scene.instantiate() as HUDBase
		add_child(hud)


func _spawn_default_pawn() -> void:
	if default_pawn_scene:
		default_pawn = default_pawn_scene.instantiate()
		add_child(default_pawn)
