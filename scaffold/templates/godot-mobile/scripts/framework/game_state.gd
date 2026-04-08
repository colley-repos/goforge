## GameState — shared truth visible to all players and systems.
##
## Mirrors Unreal's GameState: objective facts about the match.
## Score, timer, wave number, entity lists, phase — anything that
## every system needs to read. No input handling, no rendering.
##
## Separation from PlayerState:
##   GameState  = shared between all players (wave 3, boss HP 40%)
##   PlayerState = per-player data (gold, XP, selected loadout)
##
## Usage:
##   Extend this class, add your match-specific fields.
##   Read it from any system via GameMode.game_state.
##
## Example:
##   class_name RaidGameState extends GameState
##
##   var boss_hp: float = 1.0
##   var hero_count: int = 4
##   var current_phase: int = 0
class_name GameState
extends Node

## Whether the match is currently in progress.
var match_active: bool = false

## Elapsed match time in seconds.
var match_time: float = 0.0

## Current match phase index (0-based). Games with phases increment this.
var current_phase: int = 0


func _process(delta: float) -> void:
	if match_active:
		match_time += delta


## Start the match timer and mark the match active.
func begin_match() -> void:
	match_active = true
	match_time = 0.0


## Stop the match timer.
func end_match() -> void:
	match_active = false


## Advance to the next phase.
func next_phase() -> void:
	current_phase += 1
