package combat

import (
	"math/rand"
	"testing"
)

func TestDefaultResolverHit(t *testing.T) {
	// Seeded RNG for deterministic test
	rng := rand.New(rand.NewSource(42))
	r := NewDefaultResolver(rng)

	attacker := Stats{Accuracy: 80, Damage: 5, CritChance: 0}
	target := Stats{Dodge: 10, Armor: 1}
	mods := Modifiers{}

	result := r.Resolve(attacker, target, mods, 1.0)

	// With seed 42, first roll should be deterministic
	expectedRaw := 80 - 10 - 3 // 67
	if result.RawHitChance != expectedRaw {
		t.Errorf("raw hit chance: got %d, want %d", result.RawHitChance, expectedRaw)
	}
	if result.HitChance != 67 {
		t.Errorf("clamped hit chance: got %d, want 67", result.HitChance)
	}

	// Result should have valid roll
	if result.Roll < 1 || result.Roll > 100 {
		t.Errorf("roll out of range: %d", result.Roll)
	}
}

func TestDefaultResolverClamp(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	r := NewDefaultResolver(rng)

	// Very high accuracy → capped at 95
	attacker := Stats{Accuracy: 200, Damage: 5}
	target := Stats{Dodge: 0}
	result := r.Resolve(attacker, target, Modifiers{}, 0)
	if result.HitChance != 95 {
		t.Errorf("should clamp to 95, got %d", result.HitChance)
	}

	// Very low accuracy → floored at 5
	attacker2 := Stats{Accuracy: 1, Damage: 5}
	target2 := Stats{Dodge: 90}
	result2 := r.Resolve(attacker2, target2, Modifiers{}, 10)
	if result2.HitChance != 5 {
		t.Errorf("should clamp to 5, got %d", result2.HitChance)
	}
}

func TestDefaultResolverCoverPenalty(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	r := NewDefaultResolver(rng)

	attacker := Stats{Accuracy: 75, Damage: 5}
	target := Stats{Dodge: 10}
	mods := Modifiers{CoverPenalty: 20} // half cover

	result := r.Resolve(attacker, target, mods, 2.0)

	// 75 - 10 - 6 - 20 = 39
	if result.RawHitChance != 39 {
		t.Errorf("expected raw 39, got %d", result.RawHitChance)
	}
}

func TestDefaultResolverMinDamage(t *testing.T) {
	// Force a hit by using high accuracy and a favorable seed
	rng := rand.New(rand.NewSource(0))
	r := NewDefaultResolver(rng)

	attacker := Stats{Accuracy: 95, Damage: 1}
	target := Stats{Dodge: 0, Armor: 100} // armor exceeds damage

	// Run multiple times to find a hit
	for i := 0; i < 100; i++ {
		result := r.Resolve(attacker, target, Modifiers{}, 0)
		if result.Hit && result.Damage < 1 {
			t.Errorf("hit should deal at least 1 damage, got %d", result.Damage)
		}
	}
}

func TestDefaultResolverModifiers(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	r := NewDefaultResolver(rng)

	attacker := Stats{Accuracy: 50, Damage: 5}
	target := Stats{Dodge: 10}
	mods := Modifiers{
		ExtraBonus:   15,
		ExtraPenalty: 5,
	}

	result := r.Resolve(attacker, target, mods, 1.0)

	// 50 - 10 - 3 - 5 + 15 = 47
	if result.RawHitChance != 47 {
		t.Errorf("expected raw 47, got %d", result.RawHitChance)
	}
}
