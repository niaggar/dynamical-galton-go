package model

import (
	"go-galtonboard/entities"
)

type PhysicsModel interface {
	UpdateBall(particle *entities.Particle, t, dt float64)
	UpdatePeg(particle *entities.Particle, t, dt float64)
	ResolveCollision(particle *entities.Particle, peg *entities.Particle)
}
