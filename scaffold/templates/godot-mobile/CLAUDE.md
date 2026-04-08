# [Project Name]

[One-line description of the game.]

Full design: `docs/design/game_concept.md`

## Critical Findings

**ALWAYS UPDATE THIS SECTION** with discoveries that prevent future mistakes.
Delete the placeholders below and replace with real findings as you work.

- **Godot executable**: `[path to your Godot binary]`
- **Run game**: `"[godot]" --path "[project]" --scene res://scenes/Main.tscn`
- **Headless import**: `"[godot]" --headless --path "[project]" --import`
- **Error check**: `"[godot]" --path "[project]" --quit-after 300`
- **Grid unit**: [define your world grid size, e.g. "KayKit grid = 4 units"]
- **Collision layers**: [define collision layers here, e.g. "player=2, enemies=4, walls=1"]
- **Animation gotchas**: [add as discovered]
- **Android build**: [add build command if targeting Android]

## Project Structure

```
godot/           Godot 4 project (GDScript, the actual game)
docs/            Architecture, design docs, guides
tools/           Dev tooling
```

## Architecture

This project uses the **GameMode framework** from GoForge:

- `GameMode` — root scene, owns all other framework nodes, defines game rules
- `PlayerController` — translates input into game actions (no input in units)
- `GameState` — shared match truth (phase, score, entity lists)
- `PlayerState` — per-player data (gold, XP, loadout)
- `HUDBase` — UI overlay, reads state, never mutates it

Define your concrete mode in `scripts/modes/[YourMode].gd`, set its exported
scene slots in the Inspector, and add it as the root of your main scene.

## Godot Rules

GDScript only. Type hints everywhere. Signals over direct calls.
One root node per scene. Autoloads minimal (GameManager, SFXManager).
Input Map actions, never raw keycodes.
Never `set_bone_pose_rotation()` — use AnimationPlayer.

## Code Quality

No duplicate code. `push_error()` / `push_warning()` for errors — never silent
failures, never `assert()` in shipped code.

## When Adding Lessons

Update **Critical Findings** above when you discover something future sessions
need. Reusable patterns → `docs/`. Project-specific gotchas → this file.
