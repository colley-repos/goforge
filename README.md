# GoForge

Genre-agnostic game scaffold for Go + Godot projects.

## Philosophy

GoForge codifies proven game architecture patterns into a reusable foundation:

- **GameMode-first** — every game starts by defining its GameMode: which class
  fills each role (controller, HUD, game state, player state). Swap the mode,
  swap the game type. Inspired by Unreal Engine's framework.
- **Headless-first** — the Go engine has zero rendering dependencies. Renderers
  are thin pluggable clients that import the engine, never the reverse.
- **Pacing-agnostic** — turn-based, real-time, or real-time-with-pause. One
  interface swap, zero logic changes.
- **Command-driven** — all game actions are data objects queued and resolved by
  a timing authority. Enables deterministic replay, networking, and undo.
- **Event-bus decoupled** — systems communicate through typed events, never
  direct references.
- **Whitebox-first** — reach a playable prototype with zero art assets. Colored
  shapes and text prove your game loop before committing to asset pipelines.
- **ECS foundation** — entities, components, and queries via Ark. No inheritance
  hierarchies.

## Two Project Paths

### Go-native
Headless simulation engine + pluggable 2D/3D renderer. Best for: PC/web games,
turn-based, simulation-heavy, or games where you want deterministic replay/networking.

```
engine/         Pure Go game engine — zero graphics imports
renderer/       Pluggable display backends (Ebitengine 2D, terminal, Raylib 3D)
```

Start here: `examples/whitebox-tactics`

### Godot-native
Godot 4 as the full game engine with optional Go headless simulation layer.
Best for: mobile games, 3D real-time action, short session design (3-5 min).

```
scaffold/templates/godot-mobile/    Framework base classes for Godot projects
  scripts/framework/game_mode.gd    GameMode — the wiring harness
  scripts/framework/player_controller.gd
  scripts/framework/game_state.gd
  scripts/framework/player_state.gd
  scripts/framework/hud_base.gd
  scripts/autoloads/sfx_manager.gd  Preloaded SFX singleton
  scripts/units/unit_base_3d.gd     State machine base for all 3D units
  scripts/vfx/vfx_helper.gd        Procedural VFX static factory
  CLAUDE.md                         Project template with Critical Findings section
```

Start here: copy `scaffold/templates/godot-mobile/` into your Godot project.

## Structure

```
engine/         Pure Go game engine — zero graphics imports
renderer/       Pluggable display backends (Ebitengine 2D, terminal, Raylib 3D)
assets/         Asset pipeline with swappable providers (Meshy, local, self-hosted)
scaffold/       CLI tool: goforge new <project> --template <type>
examples/       Genre-proving demos (tactics, platformer, roguelike)
docs/           Architecture docs, conventions, build guides
```

## Quick Start

```bash
# Run the tactics example
cd examples/whitebox-tactics
go run .
```

## Requirements

- Go 1.23+ (for generics and range-over-func)
- No other dependencies for the engine module

## License

MIT
