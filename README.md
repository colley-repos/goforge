# GoForge

Genre-agnostic, headless-first game engine scaffold for Go.

## Philosophy

GoForge codifies proven game architecture patterns into a reusable foundation:

- **Headless-first** — the engine has zero rendering dependencies. Renderers are thin pluggable clients that import the engine, never the reverse.
- **Pacing-agnostic** — turn-based, real-time, or real-time-with-pause. One interface swap, zero logic changes.
- **Command-driven** — all game actions are data objects queued and resolved by a timing authority. Enables deterministic replay, networking, and undo.
- **Event-bus decoupled** — systems communicate through typed events, never direct references.
- **Whitebox-first** — reach a playable prototype with zero art assets. Colored shapes and text prove your game loop before committing to asset pipelines.
- **ECS foundation** — entities, components, and queries via Ark. No inheritance hierarchies.

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
