# Build & Test Patterns

Workflow patterns for Godot 4 projects bootstrapped from GoForge patterns.
These are engine-agnostic where possible, with Godot-specific notes marked.

---

## Testing Workflow

### Unit Tests (Godot)

Run headless unit tests after every change:

```bash
godot --headless --script res://test/run_unit_tests.gd
```

**Requirements:**
- The test runner script must extend `SceneTree` or `MainLoop`, not `Node`.
  `--headless --script` mode requires a MainLoop-compatible entry point.
- Use `Engine.time_scale` to accelerate tween-dependent tests. Tweens run
  in real time by default, making tests slow or flaky.

### Gameplay Testing (Godot)

Always use `--path` to test gameplay, never `--editor`:

```bash
# CORRECT — runs the game
godot --path "path/to/project"

# WRONG — opens the editor, not the game
godot --editor "path/to/project"
```

If the game uses dynamic node creation in `_ready()`, the editor preview
will show an empty scene. This is expected -- run the game to see content.

### Unit Tests (Go / GoForge Engine)

```bash
# Run all engine tests
go test ./engine/...

# Run with coverage
go test ./engine/... -coverprofile=cover.out
go tool cover -html=cover.out

# Run a specific package
go test ./engine/combat/ -v -run TestResolveAttack
```

---

## Visual Debugging: Screenshot Iteration Loop

For visual issues (rendering, layout, animation), use a screenshot-based
iteration loop:

1. Run the game with the suspect feature active
2. Take a screenshot (in-game capture or OS screenshot)
3. Analyze the screenshot to identify the issue
4. Make a targeted code change
5. Re-run and screenshot again
6. Compare before/after

This is faster than trying to reason about visual problems from code alone.
Screenshots reveal Z-fighting, incorrect transforms, missing textures, and
layout issues that are invisible in source code.

---

## Pre-Commit Checklist

Before committing changes to a Godot project:

1. **Run unit tests** (`--headless --script`)
2. **Run the game** (`--path`) and verify basic gameplay loop works
3. **Check console output** for parse errors -- GDScript parse failures are
   silent in the game window but visible in stdout/stderr
4. **Verify initialization order** -- if you refactored code, search for
   setup/init calls and confirm they still execute after entity spawning
5. **Test edge cases** -- empty states, zero values, missing resources

---

## Common Failure Modes

### Silent Script Load Failure

**Symptom:** A node has no behavior. No error in the game window.

**Cause:** A GDScript parse error prevents the script from loading. The node
exists but has no script attached at runtime.

**Fix:** Check console output for parse errors. Common causes:
- `:=` type inference on cross-object property access
- Missing `preload()` for a class_name dependency
- Syntax error introduced during refactoring

### Silent Setup Ordering Bug

**Symptom:** A system appears to work but does nothing (e.g., turns don't
advance, AI doesn't act).

**Cause:** A setup function was called before its dependencies were ready
(e.g., pacing controller initialized before agents were spawned).

**Fix:** Search for setup/init calls and verify execution order. Add runtime
assertions: `assert(agents.size() > 0, "Setup called before agents spawned")`

### Tween Tests Timing Out

**Symptom:** Tests that involve tweens hang or produce inconsistent results.

**Cause:** Tweens run in real time. A 1-second tween takes 1 real second in
tests.

**Fix:** Set `Engine.time_scale = 100.0` at the start of tween-dependent
tests, reset to `1.0` in teardown.
