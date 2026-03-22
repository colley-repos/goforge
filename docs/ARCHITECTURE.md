# GoForge Architecture

## System Boundary Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      GameMaster                              │
│  (owns ECS World, EventBus, wires all systems)              │
│                                                              │
│  ┌──────────┐  ┌───────────────┐  ┌──────────────────┐     │
│  │  Input    │→│  CommandQueue  │→│  PacingController  │     │
│  │  Mapper   │  │  (data store) │  │  (timing gate)     │     │
│  └──────────┘  └───────────────┘  └──────────────────┘     │
│       ↑              │                     │                 │
│  Raw input     Commands enqueued     When to resolve?        │
│  from Renderer       │                     │                 │
│                      ▼                     ▼                 │
│              ┌───────────────┐    ┌─────────────────┐       │
│              │ CommandResolver│    │ CombatResolver   │       │
│              │ (dispatches)  │───→│ (pure math)      │       │
│              └───────────────┘    └─────────────────┘       │
│                      │                                       │
│                      ▼                                       │
│              ┌───────────────┐                               │
│              │   EventBus    │ ← All state changes emit     │
│              │ (typed, sync) │   events for UI/audio/net    │
│              └───────────────┘                               │
│                      │                                       │
│              ┌───────────────┐    ┌─────────────────┐       │
│              │   Spatial     │    │     AI Brain     │       │
│              │ (Grid/NavMesh)│    │ (Utility/BTree)  │       │
│              └───────────────┘    └─────────────────┘       │
│                                                              │
│              ┌───────────────┐                               │
│              │ StateManager  │  ← Snapshots, save/load      │
│              └───────────────┘                               │
└─────────────────────────────────────────────────────────────┘
        ↕ (Renderer interface — engine never imports renderer)
┌───────────────────────────────┐
│         Renderer              │
│  (Ebitengine / Terminal /     │
│   Raylib / Web / Godot)       │
│                               │
│  - Reads ECS world for display│
│  - Feeds raw input to engine  │
│  - Implements Renderer iface  │
└───────────────────────────────┘
        ↕ (AssetProvider interface)
┌───────────────────────────────┐
│       Asset Pipeline          │
│  (Meshy / Local / S3 / etc)  │
│                               │
│  - Fetches assets on demand   │
│  - Content-addressed cache    │
│  - Format: GLB (3D), PNG (2D)│
└───────────────────────────────┘
```

## Data Flow

1. **Renderer** polls raw input (keyboard, mouse, touch, gamepad)
2. **InputMapper** translates raw input → semantic `Action` values
3. **Actions** create `Command` data objects, enqueued in `CommandQueue`
4. **PacingController** decides when to drain the queue (end-of-turn, immediately, on-unpause)
5. **CommandResolver** processes each command, calling pure-math resolvers as needed
6. **Results** are applied to the ECS world and emitted as typed events
7. **EventBus** distributes events to all subscribers (UI, audio, network, AI)
8. **Renderer** reads updated ECS world and renders the frame

## Key Principle: The Pacing Swap

The `PacingController` is the **sole system** that changes when you switch between:
- **Turn-based**: collect during player phase, resolve on "end turn"
- **Real-time**: resolve immediately, cooldown-gated
- **RTwP**: collect while paused, resolve on unpause

Every other system is pacing-agnostic. This is the core architectural insight
from the Dystopia project that GoForge generalizes.
