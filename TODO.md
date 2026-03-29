# GoForge Roadmap

## Phase 0 — Engine Foundation ✅

Core systems: ECS, events, commands, pacing (turn-based/realtime/RTwP), input
mapper, spatial grid, combat resolver, AI brains, state manager, config loader,
renderer interface. All with mandatory unit tests.

---

## Phase 1 — Asset Pipeline

Implement the asset system documented in `docs/ASSET_PIPELINE.md`.

### P1.1 — AssetProvider Interface & Manager
- [ ] `AssetProvider` interface: `Fetch(id, format) -> ([]byte, error)`
- [ ] `AssetManager` — wraps provider, routes requests, tracks loaded assets
- [ ] `AssetHandle` — lightweight reference returned to callers, ref-counted
- [ ] `AssetManifest` — data-driven registry mapping IDs to paths, variants, tags, LOD levels
- [ ] Variant metadata: each manifest entry supports multiple material/palette variants per mesh

### P1.2 — Local Filesystem Provider
- [ ] Read PNG (2D) and GLB (3D) from configured directory
- [ ] File watcher for hot-reload during development
- [ ] Directory structure convention: `meshes/`, `textures/`, `audio/`

### P1.3 — Meshy Provider
- [ ] REST client for Meshy API (text-to-3D, image-to-3D)
- [ ] Async workflow: submit → poll → download GLB
- [ ] Rate limiting and retry logic
- [ ] Test mode with mock responses (zero API cost in CI)

### P1.4 — Cache Layer
- [ ] Content-addressed storage (SHA256 of asset ID → filename)
- [ ] Configurable cache directory and max size
- [ ] Wraps any provider transparently
- [ ] Cache eviction: LRU when max size exceeded

### P1.5 — Asset Garbage Collection
- [ ] Reference counting on `AssetHandle` — increment on acquire, decrement on release
- [ ] Grace period: assets at ref_count=0 survive for configurable duration (default 30s)
  before being eligible for collection — prevents thrashing between rooms/scenes
- [ ] Memory budget system: per-category ceilings (texture_mb, mesh_mb, audio_mb)
  configurable per platform profile (desktop, mobile, web)
- [ ] GC sweep modes:
  - **Periodic**: sweep every N seconds, collect expired zero-ref assets
  - **Pressure**: triggered when any budget category exceeds 80% — aggressive mode,
    grace period drops to 0, evict lowest-priority assets first (furthest from
    camera, lowest LOD, oldest access time)
  - **Scene transition**: bulk unload — collect all zero-ref assets immediately
    on scene/level change
- [ ] Priority scoring for eviction: `score = distance_weight * cam_distance + age_weight * time_since_access + lod_weight * lod_level`
- [ ] GPU resource cleanup callback: provider-specific unload hook (free textures, VBOs)
- [ ] Metrics emission: GC events published to event bus for debugging/profiling
  (`AssetLoaded`, `AssetEvicted`, `GCSweep` with counts and freed bytes)
- [ ] Unit tests: ref counting lifecycle, grace period expiry, budget pressure triggers,
  eviction ordering, scene transition bulk unload

---

## Phase 2 — Renderer Implementations

### P2.1 — Ebitengine 2D Renderer
- [ ] Implement `Renderer` interface with Ebitengine backend
- [ ] Read `SpriteDescriptor` from ECS world, render as 2D sprites
- [ ] Whitebox fallback: render shape + color when no asset loaded
- [ ] Input forwarding: keyboard, mouse, touch, gamepad → engine InputMapper
- [ ] Cross-platform: Windows, Linux, macOS, Web/WASM, Android, iOS

### P2.2 — Terminal Renderer
- [ ] Text-based rendering for headless development and CI
- [ ] ASCII/Unicode grid display
- [ ] Keyboard input only

### P2.3 — Raylib 3D Renderer (Future)
- [ ] 3D rendering with Raylib-go bindings
- [ ] GLB model loading via asset pipeline
- [ ] Material/palette variant support at render level

---

## Phase 3 — Universal Game Infrastructure

Features that **every game requires** regardless of genre. GoForge should
scaffold these so new projects don't rebuild them from scratch. All
implementations are genre-agnostic — they provide structure and hooks, not
game-specific content.

### P3.1 — UI Framework Interface
- [ ] `UIProvider` interface: engine declares UI needs, renderer fulfills them
- [ ] Layout primitives: panel, label, button, list, progress bar, image
- [ ] Anchor/alignment system (center, top-left, fill, etc.)
- [ ] Theme/skin support: colors, fonts, spacing defined in config — swap themes
  without code changes
- [ ] Input focus management: which UI element receives keyboard/gamepad input
- [ ] UI events on the event bus: `ButtonPressed`, `SliderChanged`, `MenuOpened`
- [ ] Engine-side: UI is pure data (layout tree of components). Renderer-side:
  concrete drawing. Maintains headless-first principle.

### P3.2 — Screen / Scene Manager
- [ ] `Screen` interface: `OnEnter()`, `OnExit()`, `OnUpdate()`, `OnDraw()`
- [ ] Screen stack with push/pop (pause menu pushes over gameplay, pops back)
- [ ] Transition hooks: fade, slide, cut (renderer implements the visual effect)
- [ ] Predefined screen types (game provides content, GoForge provides lifecycle):
  - `SplashScreen` — logo/branding, auto-advance after duration
  - `MainMenuScreen` — title, menu items, delegates to sub-screens
  - `SettingsScreen` — reads/writes config values, binds to UI controls
  - `PauseScreen` — overlay, resume/settings/quit options
  - `LoadingScreen` — progress bar driven by asset loading / scene setup
  - `GameplayScreen` — the actual game loop (delegates to GameMaster.Tick)
  - `GameOverScreen` — results display, retry/quit options

### P3.3 — Audio Manager Interface
- [ ] `AudioProvider` interface: `PlaySFX(id)`, `PlayMusic(id)`, `StopMusic()`,
  `SetVolume(channel, float)`
- [ ] Channels: master, music, sfx, ui, voice (each independently adjustable)
- [ ] Audio triggered by events: subscribe to `AttackResolved` → play hit/miss SFX
- [ ] Music state machine: menu_theme → gameplay_theme → boss_theme, with crossfade
- [ ] Engine-side is pure data (play requests). Renderer/audio backend does the work.
- [ ] Headless mode: audio calls are no-ops (no crash, no dependency)

### P3.4 — Settings & Persistence
- [ ] `SettingsManager` — reads/writes user preferences to JSON file
- [ ] Default categories: display (resolution, fullscreen, vsync), audio (volume
  per channel), controls (key bindings from InputMapper), accessibility
  (text size, colorblind mode flag, screen reader flag)
- [ ] Settings schema defined in config — new games add game-specific settings
  without modifying GoForge core
- [ ] Auto-save on change, load on startup
- [ ] Platform-appropriate save location (AppData, ~/.config, etc.)

### P3.5 — Loading & Progress System
- [ ] `LoadingContext` — tracks N tasks, reports completion percentage
- [ ] Asset loading tasks auto-register with LoadingContext
- [ ] Custom tasks: scene setup, procedural generation, network connect
- [ ] LoadingScreen reads progress from LoadingContext
- [ ] Async loading: tasks run on background goroutines, main thread polls progress

### P3.6 — Localization Interface
- [ ] `LocaleProvider` interface: `T(key) -> string`
- [ ] String table format: JSON or CSV, one file per locale
- [ ] Fallback chain: requested locale → default locale → raw key
- [ ] All UI text goes through `T()` — no hardcoded strings
- [ ] Pluralization and parameterized strings: `T("enemies_remaining", count)`

---

## Phase 4 — Scaffold CLI

### P4.1 — Project Generator
- [ ] `goforge new <project> --template <type>` CLI tool
- [ ] Templates: `tactics`, `platformer`, `roguelike`, `blank`
- [ ] Each template includes: main.go, go.mod, config/, assets/ structure,
  screen scaffolding (main menu → gameplay → game over), placeholder UI,
  README with getting-started instructions
- [ ] Template includes all Phase 3 screens pre-wired with placeholder content

### P4.2 — Asset CLI
- [ ] `goforge asset add <id> --source meshy --prompt "..."` — generate and register
- [ ] `goforge asset variant <id> --palette <file>` — add color variant
- [ ] `goforge asset list --tag <tag>` — query manifest
- [ ] `goforge asset audit` — find unreferenced assets, missing LODs, oversized textures

---

## Phase 5 — Examples & Validation

### P5.1 — Whitebox Tactics Demo
- [ ] Port Dystopia core loop onto GoForge engine
- [ ] Demonstrate: turn-based pacing, grid combat, AI brains, save/load
- [ ] Uses all Phase 3 screens (splash → menu → gameplay → game over)

### P5.2 — Whitebox Platformer Demo
- [ ] Real-time pacing, velocity/gravity components, collision
- [ ] Demonstrates genre-agnostic engine with completely different game type

### P5.3 — Whitebox Roguelike Demo
- [ ] Turn-based pacing, procgen rooms, permadeath, inventory
- [ ] Demonstrates asset variant system (palette-based enemy recoloring)

---

## Phase 6 — Platform Builds & CI

- [ ] CI matrix: Windows, Linux, macOS, Web/WASM, Android, iOS
- [ ] Build scripts per platform (see `docs/PLATFORM_BUILDS.md`)
- [ ] Automated test runs on all platforms
- [ ] Asset budget validation in CI (warn if over platform ceiling)

---

## Future Considerations

- Networking layer (client-server, peer-to-peer via commands)
- Replay system (deterministic command log playback)
- Modding support (user-authored asset packs, script hooks)
- Analytics / telemetry interface (opt-in, privacy-respecting)
- Accessibility framework (remappable controls already exist via InputMapper;
  add screen reader hooks, colorblind simulation, subtitle system)
- **In-game asset browser / prefab editor**: Click any prop in the scene,
  browse available meshes from the asset library by category, click to swap.
  Include scale/rotate gizmo. Save layouts as presets. Essential for any game
  using large asset libraries (Synty, Meshy, etc.). Dramatically speeds up
  level design iteration vs. editing scene files by hand.
  *(Lesson from Dystopia: manual asset placement via code is painfully slow
  when you have thousands of meshes to choose from.)*
