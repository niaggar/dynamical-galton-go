package model

import (
	"go-galtonboard/entities"
	"go-galtonboard/utils"
	"math"
)

type DefaultModel struct {
	Runge RungeKutta
}

func NewDefaultModel() *DefaultModel {
	return &DefaultModel{
		Runge: RungeKutta{
			f1: dPosition,
			f2: dVelocity,
		},
	}
}

func (dm *DefaultModel) UpdateBall(particle *entities.Particle, t, dt float64) {
	state := dm.Runge.RungeKutta4(dt, t, &particle.Position, &particle.Velocity, &particle.Acceleration)
	particle.Position = *state.Position
	particle.Velocity = *state.Velocity
}

func (dm *DefaultModel) UpdatePeg(particle *entities.Particle, t, dt float64, displacement *utils.PegDisplacement) {
	if !displacement.Displacement {
		return
	}

	px := particle.Position[0]
	py := particle.Position[1]
	npx := displacement.AmplitudeX * math.Cos(displacement.FrequencyX*t)
	npy := displacement.AmplitudeY * math.Sin(displacement.FrequencyY*t)

	newDx := particle.PrevUpdateD[0] - npx
	newDy := particle.PrevUpdateD[1] - npy

	particle.PrevUpdateD = [2]float64{npx, npy}
	particle.Position = [2]float64{px + newDx, py + newDy}
}

func (dm *DefaultModel) ResolveCollision(ball *entities.Particle, peg *entities.Particle) {
	dx := ball.Position[0] - peg.Position[0]
	dy := ball.Position[1] - peg.Position[1]
	sumRadius := ball.Radius + peg.Radius

	vxa := ball.Velocity[0]
	vya := ball.Velocity[1]
	alpha0 := peg.Damping

	hip := math.Sqrt(dx*dx + dy*dy)
	sineAngle := dy / hip
	cosineAngle := dx / hip

	vTangent := -vxa*sineAngle + vya*cosineAngle
	vRadial := -alpha0 * (vxa*cosineAngle + vya*sineAngle)
	vxNew := vRadial*cosineAngle - vTangent*sineAngle
	vyNew := vRadial*sineAngle + vTangent*cosineAngle
	newX := sumRadius*cosineAngle + peg.Position[0]
	newY := sumRadius*sineAngle + peg.Position[1]

	ball.Position = [2]float64{newX, newY}
	ball.Velocity = [2]float64{vxNew, vyNew}
}

func dPosition(t float64, position, velocity, acceleration *utils.Point) *utils.Point {
	return velocity
}

func dVelocity(t float64, position, velocity, acceleration *utils.Point) *utils.Point {
	return acceleration
}
