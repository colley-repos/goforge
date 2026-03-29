# Godot 4 Gotchas

Hard-won lessons from downstream projects. These apply to any Godot 4 game,
not just a specific genre. Propagated here so every GoForge-derived project
benefits without re-discovering them.

---

## GDScript Type System

### 1. `:=` Type Inference Fails Silently on Cross-Object Access

```gdscript
# BAD — entire script silently fails to parse, no visible error
var half := other_node.TILE_SIZE * 0.5

# GOOD — always use explicit types when accessing another object's properties
var half: float = other_node.TILE_SIZE * 0.5
```

**Why it matters:** The script fails to load entirely. No error in the game
window. The only clue is a parse error buried in the console output. This
caused a multi-hour debug session where a GridRenderer silently disappeared.

**Rule:** ALWAYS use explicit type annotations (`var x: Type = ...`) when the
right-hand side involves properties from other objects.

### 2. `:=` Fails on Untyped Return Values

```gdscript
# BAD — fails if some_func() return type isn't explicitly declared
var x := some_func()

# GOOD
var x: Array[Vector2i] = some_func()
```

Even `var x := 0` is fine, but `var x := other.method()` is dangerous unless
the method has an explicit `-> Type` return annotation.

### 3. `class_name` Doesn't Auto-Resolve Without the Editor

`class_name MyClass` declarations create global scope entries, but these are
NOT available when running headless/console without the editor's class database.

```gdscript
# BAD — breaks in headless/CI
var obj = MyClass.new()

# GOOD — explicit preload always works
const MyClass = preload("res://path/to/my_class.gd")
var obj = MyClass.new()
```

### 4. God Object Splitting: Parse Errors Kill Entire Scripts

When refactoring a large script into smaller files, a parse error in ANY part
of a new file silently kills the ENTIRE script. The node using that script
appears to have no behavior at all.

**Rule:** After every file split, run the game and verify the node still works.
Check console for parse errors. Do NOT batch multiple splits before testing.

---

## Rendering & Shaders

### 5. `SCREEN_TEXTURE` Removed in Godot 4

```glsl
// BAD — Godot 3 syntax, does not exist in Godot 4
texture(SCREEN_TEXTURE, SCREEN_UV);

// GOOD — Godot 4 replacement
uniform sampler2D screen_texture : hint_screen_texture, filter_linear_mipmap;
// Then use: texture(screen_texture, SCREEN_UV);
```

### 6. Tonemap Enum Constants Not Accessible by Name

`Environment.TONE_MAP_ACES` and similar named constants are not accessible as
enum members on the Environment resource in GDScript.

```gdscript
# BAD
env.tonemap_mode = Environment.TONE_MAP_ACES

# GOOD — use integer values
# 0 = Linear, 1 = Reinhard, 2 = Filmic, 3 = ACES
env.tonemap_mode = 3
```

### 7. PlaneMesh Z-Fighting with Ground Planes

PlaneMesh tiles placed at the same Y height as a ground plane are invisible
due to Z-fighting. This is a common trap when building tile highlights or
selection indicators.

**Solution:** Use separate overlay meshes with a slight Y offset, or use a
different rendering approach (e.g., decals, viewport-based overlays) for
tile highlights. Do NOT try to fix this by changing tile material colors on
the ground mesh itself.

---

## Animation

### 8. "Animation Not Playing" Is Usually a Timing Problem

If an animation appears to not play during tween movement, the animation IS
playing -- it's just finishing too fast to see.

```gdscript
# If move_step_duration is 0.15s, a walk animation plays for 0.15s
# which is visually imperceptible

# Fix: match animation speed to movement speed
anim_player.speed_scale = base_anim_duration / move_step_duration
```

**Rule:** When debugging "animation not playing", check `move_step_duration`
or equivalent timing first, not animation tracks or bone names.

### 9. AnimationLibrary Name Prefixing

When adding an AnimationLibrary to an AnimationPlayer, all animation names
become prefixed with the library name:

```gdscript
anim_player.add_animation_library("shooter", lib)
anim_player.play("shooter/idle")  # NOT just "idle"
```

### 10. Never Use `set_bone_pose_rotation()` for Posing

Programmatic bone rotation fights the skeleton's rest pose processing and
produces subtle/broken results. Always use AnimationPlayer with
AnimationLibrary for character posing.

### 11. FBX Files Cannot Be Loaded at Runtime

Godot's FBX importer only runs at editor import time. You CANNOT `load()` a
`.fbx` file at runtime. Use pre-compiled `.res` AnimationLibrary files instead.

---

## Scene & Node Lifecycle

### 12. `set_script()` + `add_child()` Ordering

`set_script()` on a node followed by `add_child()` triggers `_ready()`.
Properties set after `add_child()` are NOT visible in `_ready()`.

```gdscript
# BAD — _ready() sees default values
var node = Node.new()
node.set_script(my_script)
add_child(node)
node.my_property = 42  # Too late, _ready() already ran

# GOOD — set properties between set_script() and add_child()
var node = Node.new()
node.set_script(my_script)
node.my_property = 42
add_child(node)  # _ready() sees my_property = 42
```

### 13. Custom UIDs in .tscn Files

Do NOT use custom `uid=` values in `.tscn` scene files. Godot generates its
own UIDs. Custom UIDs cause import conflicts.

```
# BAD
[gd_scene uid="uid://main_scene" ...]

# GOOD — omit uid entirely, let Godot assign one
[gd_scene ...]
```

### 14. Initialization Order: Setup After Spawn

Any system that depends on spawned entities (e.g., a pacing controller that
needs a list of agents) MUST be initialized AFTER the entities are spawned.

This is a silent bug: the game appears to work but the dependent system has
empty data and silently does nothing. It can happen when code refactoring
accidentally moves the setup call into a helper function that runs before
spawning.

**Rule:** After any refactor, search for setup/init calls and verify they
still execute in the correct order relative to entity spawning.

---

## Running & Testing

### 15. `--path` vs `--editor`

- `--path <dir>` runs the game (what you want for testing)
- `--editor` opens the editor (NOT for testing gameplay)
- Editor preview is empty if all nodes are created dynamically in `_ready()`

### 16. `--headless --script` Requires MainLoop

`--headless --script` mode requires the script to be a MainLoop/SceneTree
class, not a regular Node. Do NOT use it to validate individual scripts.
