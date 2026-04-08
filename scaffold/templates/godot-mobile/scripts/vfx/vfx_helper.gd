## VFXHelper — static factory for procedural 3D VFX.
##
## All functions take a container Node3D and add tween-driven mesh children.
## No scene files needed. No persistent nodes. Tweens clean up automatically
## because they're owned by the spawned mesh.
##
## Pattern: static functions, procedural meshes, self-cleaning tweens.
## Never instance VFXHelper — call its functions directly.
##
## Usage:
##   VFXHelper.spawn_ring(vfx_container, position, 2.0, Color.RED)
##   VFXHelper.spawn_bolt(vfx_container, from_pos, to_pos, Color.CYAN)
class_name VFXHelper


## Expanding ring at a world position, laid flat on the ground.
static func spawn_ring(container: Node3D, center: Vector3, radius: float, color: Color) -> void:
	var ring: MeshInstance3D = MeshInstance3D.new()
	ring.name = "RingVFX"
	var torus: TorusMesh = TorusMesh.new()
	torus.inner_radius = radius * 0.9
	torus.outer_radius = radius
	torus.rings = 24
	torus.ring_segments = 16
	var mat: StandardMaterial3D = StandardMaterial3D.new()
	mat.albedo_color = Color(color.r, color.g, color.b, 0.6)
	mat.emission_enabled = true
	mat.emission = color
	mat.emission_energy_multiplier = 2.0
	mat.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
	mat.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	torus.surface_set_material(0, mat)
	ring.mesh = torus
	container.add_child(ring)
	ring.global_position = Vector3(center.x, 0.15, center.z)
	ring.rotation_degrees.x = 90.0
	ring.scale = Vector3(0.1, 0.1, 0.1)
	var tween: Tween = ring.create_tween()
	tween.tween_property(ring, "scale", Vector3(1.0, 1.0, 1.0), 0.3)
	tween.tween_property(mat, "albedo_color:a", 0.0, 0.5)
	tween.tween_callback(ring.queue_free)


## Coloured bolt beam from one world position to another.
static func spawn_bolt(container: Node3D, from: Vector3, to: Vector3, color: Color) -> void:
	var beam: MeshInstance3D = MeshInstance3D.new()
	beam.name = "BoltVFX"
	var box: BoxMesh = BoxMesh.new()
	var dist: float = from.distance_to(to)
	box.size = Vector3(0.15, 0.15, dist)
	var mat: StandardMaterial3D = StandardMaterial3D.new()
	mat.albedo_color = Color(color.r, color.g, color.b, 0.7)
	mat.emission_enabled = true
	mat.emission = color
	mat.emission_energy_multiplier = 3.0
	mat.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
	mat.transparency = BaseMaterial3D.TRANSPARENCY_ALPHA
	box.surface_set_material(0, mat)
	beam.mesh = box
	container.add_child(beam)
	beam.global_position = (from + to) * 0.5
	beam.look_at(to)
	var tween: Tween = beam.create_tween()
	tween.tween_property(mat, "albedo_color:a", 0.0, 0.4)
	tween.tween_callback(beam.queue_free)


## Omni light flash at a position — energy decays to 0 over duration.
static func spawn_light_flash(
	container: Node3D,
	pos: Vector3,
	color: Color,
	energy: float = 3.0,
	range_m: float = 5.0,
	duration: float = 0.5,
) -> void:
	var light: OmniLight3D = OmniLight3D.new()
	light.light_color = color
	light.light_energy = energy
	light.omni_range = range_m
	container.add_child(light)
	light.global_position = pos
	var tween: Tween = light.create_tween()
	tween.tween_property(light, "light_energy", 0.0, duration)
	tween.tween_callback(light.queue_free)


## Floating 3D icon above a unit's head. Floats up, spins, fades out.
## icon_mesh: a pre-loaded Mesh (OBJ, GLB primitive, etc.)
## color:     emission tint applied via material_override
static func spawn_overhead_icon(
	unit: Node3D,
	icon_mesh: Mesh,
	color: Color,
	height: float = 2.8,
) -> void:
	var root: Node3D = Node3D.new()
	root.name = "OverheadIcon"
	unit.add_child(root)
	root.position = Vector3(0.0, height, 0.0)
	root.scale = Vector3(1.2, 1.2, 1.2)

	if icon_mesh:
		var mi: MeshInstance3D = MeshInstance3D.new()
		mi.mesh = icon_mesh
		var mat: StandardMaterial3D = StandardMaterial3D.new()
		mat.albedo_color = color
		mat.emission_enabled = true
		mat.emission = color
		mat.emission_energy_multiplier = 2.0
		mat.shading_mode = BaseMaterial3D.SHADING_MODE_UNSHADED
		mi.material_override = mat  # NOTE: use material_override, NOT surface_set_material
		root.add_child(mi)

	var tw: Tween = root.create_tween()
	tw.set_parallel(true)
	tw.tween_property(root, "position:y", root.position.y + 1.0, 1.5).set_ease(Tween.EASE_OUT)
	tw.tween_property(root, "rotation_degrees:y", 360.0, 1.5)
	tw.set_parallel(false)
	tw.tween_interval(0.4)
	tw.tween_property(root, "scale", Vector3.ZERO, 0.5).set_ease(Tween.EASE_IN)
	tw.tween_callback(root.queue_free)
