# GoForge Conventions

## Package & File Naming

| Element          | Convention          | Example                          |
|------------------|---------------------|----------------------------------|
| Packages         | lowercase single    | `events`, `pacing`, `combat`     |
| Files            | snake_case          | `command_queue.go`, `grid.go`    |
| Interfaces       | PascalCase, -er     | `Renderer`, `Resolver`           |
| Structs          | PascalCase          | `Command`, `CombatResult`        |
| Functions        | PascalCase (export) | `NewGrid()`, `ResolveAttack()`   |
| Constants        | PascalCase          | `MaxHP`, `DefaultAccuracy`       |
| Private          | camelCase           | `drainEvents()`, `tickCount`     |
| Test files       | `_test.go` suffix   | `event_bus_test.go`              |

## Architecture Rules

### 1. Headless-First
The `engine/` module MUST NOT import any rendering, audio, or platform-specific package.
Renderers import the engine — never the reverse.

### 2. Interface Boundaries
Every major system exposes an interface. Concrete implementations live alongside the interface.
Systems depend on interfaces, not concrete types.

### 3. One Package = One Responsibility
Each package under `engine/` owns exactly one system boundary:
- `ecs/`       — entity-component storage and queries
- `events/`    — typed event bus
- `commands/`  — command data + queue + resolution
- `pacing/`    — timing authority (when commands resolve)
- `input/`     — raw input → semantic actions
- `combat/`    — pure math resolvers
- `spatial/`   — grids, navmesh, spatial queries
- `ai/`        — decision-making interfaces
- `state/`     — snapshots, save/load, undo
- `config/`    — configuration loading
- `gamemaster/`— root orchestrator, system wiring

### 4. Events Over Direct Calls
Systems communicate through the event bus. Direct function calls are allowed only
for same-package internals or explicit dependency injection from `gamemaster`.

### 5. Commands Are Data
All game actions (move, attack, use ability) are `Command` structs.
Commands never execute themselves. The `CommandQueue` resolves them
when the `PacingController` says so.

### 6. Pure Math Resolvers
Combat, economy, and physics resolvers are pure functions:
- No side effects (don't mutate game state directly)
- Return auditable result structs
- Testable without any game loop running

### 7. Dependency Injection from Root
`gamemaster` creates systems and wires dependencies explicitly.
No global mutable state. No init() side effects. No service locators.

## Error Handling

- Return `error` from any operation that can fail
- Never panic in library code
- Use `fmt.Errorf("package: context: %w", err)` for wrapping
- Log warnings with `slog` at the call site, not deep in libraries

## Testing — MANDATORY

Unit tests are **not optional**. Every package MUST ship with tests. Untested
code is unfinished code.

### Hard Rules

1. **No code without tests.** Every `.go` file that contains exported logic MUST
   have a corresponding `_test.go` file in the same package. The only exceptions
   are pure type/constant definitions (e.g., `components.go` with only structs).
2. **Tests gate all merges.** `go test ./engine/...` MUST pass before any commit
   reaches `master`. CI enforces this; local workflow should too.
3. **New features require test-first or test-alongside.** Write the test before
   or concurrently with the implementation. Never "come back later."
4. **Bug fixes require regression tests.** Every bug fix commit includes a test
   that would have caught the bug.

### Test Quality Standards

- **Table-driven tests preferred.** Use `[]struct{ name string; ... }` test
  tables for combinatorial inputs.
- **Test interfaces, not implementations.** If there's an interface, write tests
  against it so all implementations are covered.
- **Deterministic by default.** No `rand` without a fixed seed. No `time.Now()`
  without injection. Tests must produce identical results on every run.
- **No test pollution.** Tests must not write to shared state, global vars, or
  the filesystem outside `t.TempDir()`. Each test is independent.
- **Descriptive names.** `TestBFSPath_BlockedByWall` not `TestBFS2`. The name
  should describe the scenario being verified.

### Coverage Expectations

- **Minimum 80% line coverage** per package is the target. Critical packages
  (combat, commands, pacing, spatial) should aim for 90%+.
- Run coverage locally: `go test ./engine/... -coverprofile=cover.out`
- View report: `go tool cover -html=cover.out`

### What to Test

| Package        | Must test                                              |
|----------------|-------------------------------------------------------|
| `events`       | Pub/sub, drain ordering, multi-subscriber, clear       |
| `commands`     | Queue FIFO, dispatch routing, factory output shapes    |
| `pacing`       | Phase transitions, faction cycling, pause/unpause      |
| `input`        | Bind/unbind, frame lifecycle, multi-action             |
| `combat`       | Hit/miss boundaries, clamp, cover, crit, modifiers     |
| `spatial`      | Grid bounds, BFS reachable, BFS path, cover, blocking  |
| `ai`           | Score selection, fallback, adapter pattern             |
| `gamemaster`   | Tick order, pacing gating, system registration         |
| `state`        | Snapshot round-trip, undo LIFO, save/load file I/O     |
| `config`       | Load, save, defaults creation, missing file handling   |

## Git

- `go.work` is gitignored (developer-local)
- CI uses `GOWORK=off`
- Each module has its own `go.mod` with proper import paths
- Commits: imperative mood, 50-char subject, body if needed
