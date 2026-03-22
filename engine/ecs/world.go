// Package ecs provides entity-component-system primitives wrapping the Ark ECS library.
// It defines common component types used across all GoForge games.
package ecs

import (
	"github.com/mlange-42/ark/ecs"
)

// World wraps the Ark ECS world, providing the central entity store.
type World = ecs.World

// NewWorld creates a new ECS world with default configuration.
func NewWorld() *World {
	return ecs.NewWorld()
}

// Entity is a handle to an entity in the world.
type Entity = ecs.Entity
