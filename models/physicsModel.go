package model

import (
	"go-galtonboard/entities"
	"go-galtonboard/utils"
)

type PhysicsModel interface {
	UpdateBall(particle *entities.Particle, t, dt float64)
	UpdatePeg(particle *entities.Particle, t, dt float64, displacement *utils.PegDisplacement)
	ResolveCollision(particle *entities.Particle, peg *entities.Particle)
}
