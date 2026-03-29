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
