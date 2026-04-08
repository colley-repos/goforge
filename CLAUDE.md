# AI Agent Memory -- GoForge Project

## What Is GoForge?

Genre-agnostic, headless-first game engine scaffold. See README.md for overview,
CONVENTIONS.md for coding standards, TODO.md for roadmap.

## Cross-Project Learning

GoForge is the upstream bootstrap. Downstream projects (Dystopia, etc.) propagate
reusable lessons back here. When working on GoForge, check these docs for
accumulated wisdom:

- `docs/GODOT_GOTCHAS.md` -- GDScript and Godot 4 engine pitfalls
- `docs/ARCHITECTURE_PATTERNS.md` -- Proven game architecture patterns
- `docs/BUILD_TEST_PATTERNS.md` -- Build, test, and debugging workflows
- `docs/ASSET_PIPELINE.md` -- Asset management, assessment, animation pipeline
- `docs/ARCHITECTURE.md` -- System boundary diagram and data flow

## Key Rules

1. **GameMode-first**: Every new project defines a GameMode that declares which
   class fills each role (controller, HUD, game state, player state). This is the
   Unreal Engine framework pattern adapted for Go + Godot. See
   `scaffold/templates/godot-mobile/scripts/framework/game_mode.gd`.
2. **Headless-first**: `engine/` must NEVER import rendering, audio, or
   platform-specific packages. Renderers import the engine, not the reverse.
3. **Tests are mandatory**: Every package ships with tests. See CONVENTIONS.md
   for full testing standards.
4. **Commands are data**: All game actions are Command structs, resolved by a
   timing authority. No direct state mutation from input handlers.
5. **Events over direct calls**: Systems communicate through the typed event bus.
6. **Pure math resolvers**: Combat/economy/physics resolvers are pure functions
   that return result structs. No side effects.

## When Adding Lessons

If you discover a reusable pattern or gotcha while working on GoForge or any
downstream project:

1. Determine if it's project-specific or transferable
2. If transferable, add it to the appropriate GoForge doc:
   - Engine gotcha -> `docs/GODOT_GOTCHAS.md`
   - Architecture pattern -> `docs/ARCHITECTURE_PATTERNS.md`
   - Build/test workflow -> `docs/BUILD_TEST_PATTERNS.md`
   - Asset pipeline -> `docs/ASSET_PIPELINE.md`
3. If it's a new feature idea, add it to `TODO.md` under Future Considerations
4. Update this file's doc index if you create new files
