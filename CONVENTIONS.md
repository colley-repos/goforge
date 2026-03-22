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

## Testing

- Every package has `_test.go` files
- Table-driven tests preferred
- Test interfaces, not implementations
- `go test ./engine/...` must always pass
- Aim for deterministic tests — no random unless seeded

## Git

- `go.work` is gitignored (developer-local)
- CI uses `GOWORK=off`
- Each module has its own `go.mod` with proper import paths
- Commits: imperative mood, 50-char subject, body if needed
