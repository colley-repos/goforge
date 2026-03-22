// Package combat provides interfaces and default implementations for combat resolution.
//
// Combat resolvers are pure math — no side effects, no state mutation, no UI knowledge.
// They receive attacker/target data and return auditable result structs.
package combat

import (
	"math"
	"math/rand"
)

// Stats holds combat-relevant attributes for an entity.
type Stats struct {
	Accuracy     int // base hit chance (0-100)
	Dodge        int // evasion chance (0-100)
	Damage       int // base damage per hit
	AttackRange  int // max attack range in grid/world units
	Armor        int // flat damage reduction
	CritChance   int // critical hit chance (0-100)
	CritMultiply float64 // critical damage multiplier (e.g., 1.5)
}

// Modifiers are situational adjustments applied to a combat calculation.
type Modifiers struct {
	CoverPenalty    int     // penalty from target's cover (e.g., -20 for half cover)
	DistancePenalty int     // penalty from distance (e.g., -3 per tile)
	ExtraPenalty    int     // additional penalty (overwatch, suppression, etc.)
	ExtraBonus      int     // bonus accuracy (flanking, height advantage, etc.)
	DamageBonus     int     // flat bonus damage
	DamageMultiply  float64 // multiplicative damage (default 1.0)
}

// Result is the auditable outcome of a combat resolution.
type Result struct {
	Hit          bool
	Critical     bool
	Damage       int
	HitChance    int // calculated hit chance (clamped)
	RawHitChance int // hit chance before clamping
	Roll         int // the random roll (1-100)
	AttackerStats Stats
	TargetStats   Stats
	Modifiers    Modifiers
}

// Resolver computes a combat outcome.
type Resolver interface {
	// Resolve calculates the result of an attack.
	Resolve(attacker, target Stats, mods Modifiers, distance float64) Result
}

// DefaultResolver implements the standard GoForge hit/damage formula:
//
//	hitChance = accuracy - dodge - (distance * distPenalty) - coverPenalty - extraPenalty + extraBonus
//	clamped to [MinHitChance, MaxHitChance]
//	roll 1-100; hit if roll <= hitChance
//	damage = baseDamage - armor + damageBonus, multiplied by damageMult and critMult
type DefaultResolver struct {
	MinHitChance     int // default: 5
	MaxHitChance     int // default: 95
	DistancePenalty  int // per-unit distance penalty, default: 3
	RNG              *rand.Rand
}

// NewDefaultResolver creates a resolver with sane defaults.
// Pass a seeded rand.Rand for deterministic results (testing/replay), or nil for random.
func NewDefaultResolver(rng *rand.Rand) *DefaultResolver {
	if rng == nil {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	return &DefaultResolver{
		MinHitChance:    5,
		MaxHitChance:    95,
		DistancePenalty: 3,
		RNG:             rng,
	}
}

// Resolve computes a combat result.
func (r *DefaultResolver) Resolve(attacker, target Stats, mods Modifiers, distance float64) Result {
	distPenalty := int(math.Round(distance)) * r.DistancePenalty

	rawHitChance := attacker.Accuracy - target.Dodge - distPenalty -
		mods.CoverPenalty - mods.ExtraPenalty + mods.ExtraBonus

	hitChance := clamp(rawHitChance, r.MinHitChance, r.MaxHitChance)

	roll := r.RNG.Intn(100) + 1 // 1-100

	hit := roll <= hitChance

	result := Result{
		Hit:          hit,
		HitChance:    hitChance,
		RawHitChance: rawHitChance,
		Roll:         roll,
		AttackerStats: attacker,
		TargetStats:   target,
		Modifiers:    mods,
	}

	if hit {
		// Calculate damage
		baseDamage := attacker.Damage - target.Armor + mods.DamageBonus
		if baseDamage < 1 {
			baseDamage = 1 // minimum 1 damage on hit
		}

		damageMult := mods.DamageMultiply
		if damageMult == 0 {
			damageMult = 1.0
		}

		// Check critical
		critRoll := r.RNG.Intn(100) + 1
		critChance := attacker.CritChance
		if critChance > 0 && critRoll <= critChance {
			result.Critical = true
			critMult := attacker.CritMultiply
			if critMult == 0 {
				critMult = 1.5
			}
			damageMult *= critMult
		}

		result.Damage = int(math.Round(float64(baseDamage) * damageMult))
		if result.Damage < 1 {
			result.Damage = 1
		}
	}

	return result
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
