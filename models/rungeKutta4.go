package model

import (
	"go-galtonboard/utils"
)

type RungeState struct {
	Position *utils.Point
	Velocity *utils.Point
}

type DiffEq func(t float64, position, velocity, acceleration *utils.Point) *utils.Point

type RungeKutta struct {
	f1, f2 DiffEq
}

func (rk *RungeKutta) RungeKutta4(dt, t float64, position, velocity, acceleration *utils.Point) RungeState {
	var (
		k11, k21, k12, k22, k13, k23, k14, k24 *utils.Point
	)

	k11 = utils.Mul(rk.f1(t, position, velocity, acceleration), dt)
	k21 = utils.Mul(rk.f2(t, position, velocity, acceleration), dt)

	k12 = utils.Mul(rk.f1(t+dt/2.0, utils.Add(position, utils.Mul(k11, 0.5)), utils.Add(velocity, utils.Mul(k21, 0.5)), acceleration), dt)
	k22 = utils.Mul(rk.f2(t+dt/2.0, utils.Add(position, utils.Mul(k11, 0.5)), utils.Add(velocity, utils.Mul(k21, 0.5)), acceleration), dt)

	k13 = utils.Mul(rk.f1(t+dt/2.0, utils.Add(position, utils.Mul(k12, 0.5)), utils.Add(velocity, utils.Mul(k22, 0.5)), acceleration), dt)
	k23 = utils.Mul(rk.f2(t+dt/2.0, utils.Add(position, utils.Mul(k12, 0.5)), utils.Add(velocity, utils.Mul(k22, 0.5)), acceleration), dt)

	k14 = utils.Mul(rk.f1(t+dt, utils.Add(position, k13), utils.Add(velocity, k23), acceleration), dt)
	k24 = utils.Mul(rk.f2(t+dt, utils.Add(position, k13), utils.Add(velocity, k23), acceleration), dt)

	posT := utils.Add(position, utils.Mul(utils.Add(utils.Add(k11, utils.Mul(k12, 2.0)), utils.Add(utils.Mul(k13, 2.0), k14)), 1.0/6.0))
	velT := utils.Add(velocity, utils.Mul(utils.Add(utils.Add(k21, utils.Mul(k22, 2.0)), utils.Add(utils.Mul(k23, 2.0), k24)), 1.0/6.0))

	return RungeState{
		Position: posT,
		Velocity: velT,
	}
}
