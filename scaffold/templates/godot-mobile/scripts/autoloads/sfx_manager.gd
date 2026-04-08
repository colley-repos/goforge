extends Node
class_name SFXManager

## SFX Manager — preloads all sounds at startup, plays by name.
##
## Key pattern: load ONCE in _ready(), never at play time.
## On-demand loading causes audio hitches on mobile.
##
## Usage:
##   1. Add entries to SFX_MAP pointing to res://assets/sfx/*.ogg
##   2. Call SFXManager.get_instance(self).play("my_sound")
##   3. Use play_random() for variety (attack variants, footsteps, etc.)
##
## Add as an Autoload named SFXManager, OR find via group "sfx_manager".

## Map of logical name → asset path. Add all game sounds here.
const SFX_MAP: Dictionary = {
	# "ui_tap":   "res://assets/sfx/ui_tap.ogg",
	# "hit":      "res://assets/sfx/hit.ogg",
}

var _players: Dictionary = {}  # name -> AudioStreamPlayer
var _loaded: bool = false


func _ready() -> void:
	add_to_group("sfx_manager")
	for sfx_name: String in SFX_MAP:
		var path: String = SFX_MAP[sfx_name]
		if ResourceLoader.exists(path):
			var stream: AudioStream = load(path) as AudioStream
			if stream:
				var player: AudioStreamPlayer = AudioStreamPlayer.new()
				player.name = sfx_name
				player.stream = stream
				player.volume_db = -6.0
				add_child(player)
				_players[sfx_name] = player
	_loaded = true


func play(sfx_name: String, volume_db: float = -6.0) -> void:
	if not _loaded:
		return
	var player: AudioStreamPlayer = _players.get(sfx_name) as AudioStreamPlayer
	if player:
		player.volume_db = volume_db
		player.play()


func play_random(sfx_names: Array[String], volume_db: float = -6.0) -> void:
	if sfx_names.is_empty():
		return
	play(sfx_names[randi() % sfx_names.size()], volume_db)


## Static helper — find the SFX manager from any node in the tree.
static func get_instance(from_node: Node) -> SFXManager:
	var nodes: Array[Node] = from_node.get_tree().get_nodes_in_group("sfx_manager")
	if nodes.size() > 0:
		return nodes[0] as SFXManager
	return null
