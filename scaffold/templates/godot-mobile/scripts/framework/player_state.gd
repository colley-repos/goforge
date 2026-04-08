## PlayerState — per-player data for one session.
##
## Mirrors Unreal's PlayerState: data that belongs to one player,
## not the whole match. Score, XP, gold, loadout, deaths.
##
## Separation from GameState:
##   PlayerState = "my data" (gold, kills, selected abilities)
##   GameState   = "match data" (wave, timer, boss HP)
##
## Usage:
##   Extend this class, add per-player fields.
##   One PlayerState is created per player by GameMode.on_player_join().
##
## Example:
##   class_name RaidPlayerState extends PlayerState
##
##   var gold: int = 0
##   var selected_abilities: Array[String] = []
class_name PlayerState
extends Node

## The player ID this state belongs to (matches PlayerController.player_id).
var player_id: int = 0

## Player display name.
var player_name: String = "Player"

## Number of times this player has died in the current match.
var death_count: int = 0


## Record a death and increment the counter.
func record_death() -> void:
	death_count += 1
