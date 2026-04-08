## UnitBase3D — base class for all 3D units (heroes, enemies, bosses).
##
## Provides the standard unit state machine, health, animation, and signal layer.
## All unit-specific behavior lives in subclasses — this is structure only.
##
## State machine pattern (learned from RaidBoss):
##   - Enum states with explicit transitions via set_state()
##   - Per-state _process handlers avoid "check 6 booleans" anti-pattern
##   - WAITING state for gate conditions (phase transitions, spawn locks)
##
## Usage:
##   Extend UnitBase3D for any character in the game.
##   Override _process_*() methods for per-state behavior.
##   Connect to signals for external reactions (damage, death, state changes).
class_name UnitBase3D
extends CharacterBody3D

# ---------------------------------------------------------------------------
# Signals
# ---------------------------------------------------------------------------

signal died(unit: UnitBase3D)
signal damaged(amount: float, unit: UnitBase3D)
signal state_changed(old_state: State, new_state: State)

# ---------------------------------------------------------------------------
# State machine
# ---------------------------------------------------------------------------

enum State {
	IDLE,      ## Standing still, no orders.
	MOVING,    ## Travelling to a destination.
	FIGHTING,  ## Engaged in combat.
	WAITING,   ## Paused by external gate (phase transition, spawn lock).
	DEAD,      ## Death animation playing, no further processing.
}

var current_state: State = State.IDLE

# ---------------------------------------------------------------------------
# Stats — override in subclass or set via data resource
# ---------------------------------------------------------------------------

@export var max_hp: float = 100.0
var current_hp: float = max_hp

@export var move_speed: float = 4.0

# ---------------------------------------------------------------------------
# Node references — set in subclass _ready() after add_child
# ---------------------------------------------------------------------------

var _anim_player: AnimationPlayer = null

# ---------------------------------------------------------------------------
# Lifecycle
# ---------------------------------------------------------------------------

func _ready() -> void:
	current_hp = max_hp
	_on_ready()


## Override in subclass to do post-ready wiring.
func _on_ready() -> void:
	pass


func _process(delta: float) -> void:
	match current_state:
		State.IDLE:     _process_idle(delta)
		State.MOVING:   _process_moving(delta)
		State.FIGHTING: _process_fighting(delta)
		State.WAITING:  _process_waiting(delta)
		State.DEAD:     pass  # no processing after death


## Override per-state handlers in subclass.
func _process_idle(_delta: float) -> void:
	pass

func _process_moving(_delta: float) -> void:
	pass

func _process_fighting(_delta: float) -> void:
	pass

func _process_waiting(_delta: float) -> void:
	pass

# ---------------------------------------------------------------------------
# State transitions
# ---------------------------------------------------------------------------

func set_state(new_state: State) -> void:
	if new_state == current_state:
		return
	var old: State = current_state
	current_state = new_state
	state_changed.emit(old, new_state)
	_on_state_changed(old, new_state)


## Called after every state transition. Override for entry/exit logic.
func _on_state_changed(_old: State, _new: State) -> void:
	pass

# ---------------------------------------------------------------------------
# Health
# ---------------------------------------------------------------------------

func take_damage(amount: float) -> void:
	if current_state == State.DEAD:
		return
	current_hp = maxf(0.0, current_hp - amount)
	damaged.emit(amount, self)
	if current_hp <= 0.0:
		die()


func heal(amount: float) -> void:
	if current_state == State.DEAD:
		return
	current_hp = minf(max_hp, current_hp + amount)


func die() -> void:
	if current_state == State.DEAD:
		return
	set_state(State.DEAD)
	died.emit(self)
	_on_die()


## Override to play death animation, spawn loot, etc.
func _on_die() -> void:
	pass


func is_alive() -> bool:
	return current_state != State.DEAD


func hp_percent() -> float:
	return current_hp / max_hp

# ---------------------------------------------------------------------------
# Animation
# ---------------------------------------------------------------------------

## Play an animation by name. Safe to call before _anim_player is assigned.
func play_anim(anim_name: String, loop: bool = false) -> void:
	if not _anim_player:
		return
	if not _anim_player.has_animation(anim_name):
		push_warning("UnitBase3D: animation '%s' not found on %s" % [anim_name, name])
		return
	_anim_player.get_animation(anim_name).loop_mode = (
		Animation.LOOP_LINEAR if loop else Animation.LOOP_NONE
	)
	_anim_player.play(anim_name)
