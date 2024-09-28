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

func (dm *DefaultModel) UpdatePeg(particle *entities.Particle, t, dt float64) {

}

func (dm *DefaultModel) ResolveCollision(ball *entities.Particle, peg *entities.Particle) {
	dx := ball.Position[0] - peg.Position[0]
	dy := ball.Position[1] - peg.Position[1]
	sumRadius := ball.Radius + peg.Radius
	impactAngle := math.Atan2(dy, dx)

	vxa := ball.Velocity[0]
	vya := ball.Velocity[1]
	alpha0 := peg.Damping

	//sineAngle := math.Sin(impactAngle)
	//cosineAngle := math.Cos(impactAngle)

	vTangen := -vxa*math.Sin(impactAngle) + vya*math.Cos(impactAngle)
	vRadial := -alpha0 * (vxa*math.Cos(impactAngle) + vya*math.Sin(impactAngle))

	vxNew := vRadial*math.Cos(impactAngle) - vTangen*math.Sin(impactAngle)
	vyNew := vRadial*math.Sin(impactAngle) + vTangen*math.Cos(impactAngle)

	newX := sumRadius*math.Cos(impactAngle) + peg.Position[0]
	newY := sumRadius*math.Sin(impactAngle) + peg.Position[1]

	ball.Position = [2]float64{newX, newY}
	ball.Velocity = [2]float64{vxNew, vyNew}
}

func dPosition(t float64, position, velocity, acceleration *utils.Point) *utils.Point {
	return velocity
}

func dVelocity(t float64, position, velocity, acceleration *utils.Point) *utils.Point {
	return acceleration
}
