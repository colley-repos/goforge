# Reusable Game Architecture Patterns

Genre-agnostic patterns proven in downstream GoForge projects. These are
structural templates, not implementations -- adapt them to your game's
specific needs.

---

## 1. Data-Driven Content via Resource Files

**Problem:** Hardcoded weapon stats, armor values, class definitions, and
mission parameters scattered across code files. Every balance change requires
a code edit, rebuild, and redeploy.

**Solution:** Define all content as data resources (`.tres` in Godot, JSON/TOML
in Go, YAML in other engines). Code reads the data; code never contains
balance numbers.

```
# Directory structure
data/
  weapons/
    assault_rifle.tres
    shotgun.tres
  armor/
    kevlar_vest.tres
    tactical_plate.tres
  classes/
    assault.tres
    medic.tres
    sniper.tres
  missions/
    tutorial_01.tres
    first_contact.tres
```

**Benefits:**
- Balance changes without recompilation
- Designers can edit data files without touching code
- Easy to add content (new weapon = new file, no code changes)
- Version control shows balance changes as data diffs, not code diffs

**Anti-pattern:** A single giant `GameData` struct/dictionary with everything
hardcoded. This is the #1 scaling bottleneck in hobby projects.

---

## 2. Equipment System with Slot Architecture

**Pattern:** An `EquipmentSlots` class/struct that manages named slots
(weapon, armor, gadget_1, gadget_2, etc.) with typed getters.

```
EquipmentSlots:
  slots: Dict[SlotName, ItemResource]

  get_weapon() -> WeaponResource
  get_armor() -> ArmorResource
  get_gadget(index) -> GadgetResource

  equip(slot, item) -> old_item
  unequip(slot) -> item

  get_total_stat_modifier(stat_name) -> float
```

**Key insight:** Use property getters that compute derived stats by iterating
all equipped items. This provides backward compatibility -- code that reads
`agent.accuracy` still works even after you add equipment modifiers, because
the getter sums base + equipment bonuses.

---

## 3. Status Effect / Injury System

**Pattern:** Roll-table-driven status effects with Resource definitions and
stat penalties applied via a `recalculate()` method.

```
StatusEffect (Resource):
  name: String
  stat_penalties: Dict[StatName, float]  # e.g. {"accuracy": -15, "movement": -1}
  duration_turns: int  # -1 = permanent until healed
  severity: enum (LIGHT, MODERATE, SEVERE, CRITICAL)

InjuryTable (Resource):
  entries: Array[{weight: float, effect: StatusEffect}]

  roll() -> StatusEffect  # Weighted random selection

Agent:
  base_stats: Dict[StatName, float]
  active_effects: Array[StatusEffect]

  recalculate():
    for each stat:
      effective_value = base_stat + sum(effect.penalties[stat] for effect in active_effects)
      effective_value = clamp(effective_value, min, max)
```

**Key insight:** Never modify base stats directly. Always keep base stats
pristine and compute effective stats by layering modifiers. This makes it
trivial to add/remove effects without tracking what the "original" value was.

---

## 4. Use-Based Skill Leveling

**Pattern:** Skills improve through use, not through XP allocation.

```
SkillTracker:
  uses: Dict[SkillName, int]
  levels: Dict[SkillName, int]

  record_use(skill):
    uses[skill] += 1
    threshold = levels[skill] * 3  # or whatever curve
    if uses[skill] >= threshold:
      levels[skill] += 1
      uses[skill] = 0

  get_bonus(skill) -> int:
    return levels[skill] * 2  # flat bonus per level
```

**Benefits:**
- Emergent specialization: agents get better at what they actually do
- No skill point allocation UI needed
- Feels organic to the player ("my sniper got better at sniping because I
  used them as a sniper")

**Tuning knobs:** threshold curve (linear, quadratic, logarithmic), bonus
per level, max level cap.

---

## 5. Combat Math as Pure Static Functions

**Pattern:** Separate combat math (pure functions, no side effects) from
combat resolution (applies results, mutates state, emits events).

```
# Pure math -- testable without any game loop
CombatMath:
  static calculate_hit_chance(attacker_accuracy, target_evasion, cover_bonus, range_penalty) -> float
  static calculate_damage(weapon_damage, armor_reduction, crit_multiplier) -> int
  static calculate_crit_chance(base_crit, skill_bonus, target_vulnerability) -> float

# Resolution -- has side effects
CombatResolver:
  resolve_attack(attacker, target, weapon):
    hit_chance = CombatMath.calculate_hit_chance(...)
    roll = rng.randf()
    if roll < hit_chance:
      damage = CombatMath.calculate_damage(...)
      target.apply_damage(damage)
      event_bus.emit(AttackHit{attacker, target, damage})
    else:
      event_bus.emit(AttackMissed{attacker, target})
```

**Benefits:**
- Math functions are trivially unit-testable (table-driven tests)
- Resolver is thin and auditable
- Easy to add new combat modifiers without touching resolution logic
- Deterministic: same inputs always produce same outputs

---

## 6. AI Brain Pattern (Strategy Pattern)

**Pattern:** Base brain interface with specialized subclasses that produce
different behavior without changing the AI framework.

```
Brain (interface):
  evaluate(agent, world_state) -> Command

AggressiveBrain:
  evaluate(agent, world_state):
    target = find_nearest_enemy(agent, world_state)
    if in_range(agent, target):
      return AttackCommand(agent, target)
    else:
      return MoveCommand(agent, toward(target))

DefensiveBrain:
  evaluate(agent, world_state):
    if agent.health < threshold:
      cover = find_nearest_cover(agent, world_state)
      return MoveCommand(agent, cover)
    target = find_weakest_enemy(agent, world_state)
    return AttackCommand(agent, target)

CautiousBrain:
  evaluate(agent, world_state):
    threats = find_visible_enemies(agent, world_state)
    if len(threats) > 2:
      return OverwatchCommand(agent)
    ...
```

**Key insight:** The brain returns a Command (data), not an action. The
command goes through the same CommandQueue as player commands. This means
AI and player actions are processed identically -- crucial for replays,
networking, and undo.

---

## 7. Procedural Audio Engine (GDScript-Specific)

**Pattern:** Full synthesizer engine using `AudioStreamGenerator` -- no
sample files needed for music or procedural sound effects.

```
AudioStreamGenerator setup:
  - Create AudioStreamPlayer with AudioStreamGenerator
  - Set mix_rate (44100 Hz standard)
  - Get AudioStreamGeneratorPlayback
  - Push frames via push_frame(Vector2(left, right))

Synth components:
  Oscillator: sine, square, sawtooth, triangle waveforms
  Envelope: ADSR (attack, decay, sustain, release)
  Filter: low-pass, high-pass, band-pass
  Effects: delay, reverb (via feedback buffers)
  Sequencer: pattern-based note triggering

Use cases:
  - Ambient music generation (no audio files needed)
  - Procedural sound effects (explosions, UI clicks, footsteps)
  - Adaptive music that responds to game state
```

**Benefits:**
- Zero audio asset dependencies for prototyping
- Adaptive music that reacts to game state in real-time
- Tiny build size (no audio samples)
- Unique soundscape for every playthrough

**Caveat:** This is a significant engineering investment. Use it for
prototyping or as a creative choice, not as a replacement for professional
audio in a shipped product (unless that's your aesthetic).

---

---

## 8. GameMode / Framework Pattern (Unreal-Inspired)

**Problem:** Every new game project rebuilds the same wiring: who spawns the
player, who owns the HUD, who tracks score, how does match start and end?
The first weeks of a project are spent on plumbing, not gameplay.

**Solution:** A `GameMode` class that declares the *slots* a game type needs
and fills them. Swap the GameMode → swap the entire game type. Same engine,
completely different rules.

```
GameMode (root scene / orchestrator)
  ├── PlayerController  ← input → actions (no input code in units or HUD)
  ├── HUD               ← UI overlay, reads state, never mutates it
  ├── GameState         ← shared truth (phase, timer, score, entity lists)
  └── PlayerState       ← per-player data (gold, XP, loadout, deaths)
```

Key insight: **the unit doesn't handle input. The controller does.**
This is the same reason Unreal separates PlayerController from Pawn — it
lets you possess any pawn, drive it from AI instead of a player, or
spectate without changing the pawn at all.

**In GDScript:**

```gdscript
class_name MyGameMode extends GameMode

func match_start() -> void:
    game_state.spawn_enemies(wave_data)
    hud.show_encounter_ui()

func match_end(winner: String) -> void:
    hud.show_result_modal(winner)
```

**Benefits:**
- New game type = new GameMode subclass, nothing else changes
- Every project starts with HUD, state, and controller wired correctly
- Testing: swap in a TestGameMode with no HUD, reduced enemy count, keyboard cheats

---

## 9. Director / Encounter Phase Pattern

**Problem:** A boss encounter or wave system has complex timing: abilities
fire on cooldowns, phases trigger at HP thresholds, multiple subsystems
need to coordinate without knowing about each other.

**Solution:** A `Director` singleton (autoload) that owns the encounter
state machine. Units have AI for *local* decisions (movement, targeting).
The Director makes *global* decisions (when to start phase 2, which hero
to taunt, when to call reinforcements). They communicate via signals only.

```
Director (autoload singleton)
  ├── Phase state machine (phase_1, phase_2, enrage)
  ├── Ability cooldown timers
  ├── Global targeting rules ("taunt the tank", "slam the cluster")
  └── Signals: ability_fired, phase_changed, encounter_ended

Unit AI (attached to each unit)
  ├── Local decisions: move to target, play attack animation
  ├── Responds to: Director.ability_fired → react to AoE
  └── Never calls Director directly
```

**Key rules:**
- Director → units via signals. Units never call Director.
- Director knows *what* to do and *when*. Unit AI knows *how* to move.
- Combat math (damage, hit chance) lives in a separate pure resolver.

**Pacing:** Director runs ability cooldowns on `_process()`. When a
cooldown expires, Director picks a target using global knowledge (who
has least HP, who is in a cluster), fires the ability signal, and
resets the timer. No unit needs to know this logic exists.

---

## 10. Unit State Machine Pattern

**Problem:** Units accumulate boolean flags (`is_attacking`, `is_moving`,
`is_stunned`, `is_dead`) that interact in undocumented ways. "Why isn't
the unit moving?" requires checking 6 booleans.

**Solution:** A single typed enum state. One state at a time. Explicit
transitions with an `_on_state_changed()` hook.

```gdscript
enum State { IDLE, MOVING, FIGHTING, WAITING, DEAD }

var current_state: State = State.IDLE

func _process(delta: float) -> void:
    match current_state:
        State.IDLE:     _process_idle(delta)
        State.MOVING:   _process_moving(delta)
        State.FIGHTING: _process_fighting(delta)
        State.WAITING:  pass  # gate condition, do nothing
        State.DEAD:     pass  # no processing

func set_state(new_state: State) -> void:
    if new_state == current_state:
        return
    var old := current_state
    current_state = new_state
    _on_state_changed(old, new_state)
```

**`WAITING` state:** Any time an external system needs to pause a unit
(phase gate, spawn lock, cutscene), call `unit.set_state(State.WAITING)`.
The unit stops processing, plays an idle animation, and resumes when the
gate clears. Never add a boolean for this.

---

## 11. Combat / AI Separation

**Problem:** Combat logic (damage calc, cooldowns, ability selection)
bleeds into unit AI scripts. Changing an ability requires touching every
unit that uses it.

**Solution:** Two distinct responsibilities, never mixed:

| Layer | Owns | Example |
|-------|------|---------|
| `CombatDirector` (autoload) | encounter rules, ability timing, global targeting | "fire cleave at the three clustered heroes" |
| `UnitAI` (per-unit node) | local movement, animation, reaction to combat events | "move toward target, play attack anim" |
| `CombatResolver` (static) | pure math — damage, hit chance, crit | `CombatResolver.calculate_damage(atk, def)` |

Units never call the Director. Director signals units. Resolver is called
by whoever needs a number — Director, units, or UI.

---

## General Principles

1. **Separate data from behavior.** Content (stats, definitions, tables) lives
   in data files. Code reads data and implements mechanics.

2. **Separate pure math from side effects.** Math functions are stateless and
   testable. Resolution functions apply results and emit events.

3. **Use the Command pattern for all game actions.** Player actions, AI actions,
   and scripted events all produce the same Command objects. This enables
   replay, undo, networking, and testing.

4. **Prefer composition over inheritance.** An entity is a bag of components,
   not a deep class hierarchy. A "medic" is an entity with HealAbility +
   LowAccuracy components, not a MedicClass extending SoldierClass extending
   UnitClass extending Entity.

5. **Make the implicit explicit.** If a system depends on another system being
   initialized first, document it and enforce it (assert, error check, or
   API design that makes wrong ordering impossible).
