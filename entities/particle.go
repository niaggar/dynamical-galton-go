package entities

import (
	"go-galtonboard/utils"
	"log"
	"math"
	"math/rand/v2"
)

// Particle represents a particle
type Particle struct {
	Position     utils.Point
	Velocity     utils.Point
	Acceleration utils.Point
	Damping      float64
	Radius       float64
	Type         int
	IsStopped    bool
}

// NewParticles returns a new particle with the given values.
func NewParticles(config utils.ParticleConfig, startPoint *utils.Point) []*Particle {
	particles := make([]*Particle, config.NParticles)

	for i := 0; i < config.NParticles; i++ {
		randomVx := config.InitDeltaVx * (2*rand.Float64() - 1)
		randomVy := config.InitDeltaVy * rand.Float64()

		randomX := config.InitDeltaX * (2*rand.Float64() - 1)
		randomY := config.InitDeltaY * rand.Float64()

		particles[i] = &Particle{
			Position: utils.Point{startPoint[0] + randomX, startPoint[1] + randomY},
			Velocity: utils.Point{randomVx, randomVy},
			Damping:  1,
			Radius:   config.Radius,
			Type:     utils.Particle,
		}
	}

	return particles
}

// NewPegs returns a new peg with the given values.
func NewPegs(pegConfig utils.PegConfig, boardConfig utils.BoardConfig) ([]*Particle, []*utils.Point) {
	pegs := make([]*Particle, 0)
	border := make([]*utils.Point, 5)

	for i := 0; i < boardConfig.NRows; i++ {
		for j := 0; j < boardConfig.NCols; j++ {
			peg := &Particle{}
			radius := getRadius(&pegConfig, &boardConfig, i, j)

			x := float64(j) * (boardConfig.HorizontalSpace)
			y := float64(i) * (boardConfig.VerticalSpace)
			if i%2 != 0 {
				x += (boardConfig.HorizontalSpace) / 2

				if j == boardConfig.NCols-1 {
					continue
				}
			}

			peg.Position = utils.Point{x, y}
			peg.Radius = radius
			peg.Damping = pegConfig.Damping
			peg.Type = utils.Peg
			pegs = append(pegs, peg)
		}
	}

	width := boardConfig.HorizontalSpace * float64(boardConfig.NCols-1)
	height := boardConfig.VerticalSpace * float64(boardConfig.NRows+1)

	border[0] = &utils.Point{width / 2, height}
	border[1] = &utils.Point{width, height}
	border[2] = &utils.Point{width, 0}
	border[3] = &utils.Point{0, 0}
	border[4] = &utils.Point{0, height}

	return pegs, border
}

func getRadius(pegConfig *utils.PegConfig, boardConfig *utils.BoardConfig, row, column int) float64 {
	switch pegConfig.Distribution {
	case utils.PegUniformDist:
		return pegConfig.MinRadius

	// Horizontal distributions
	case utils.PegLogarithmicDistHorizontal:
		x := column - boardConfig.NCols/2
		if x == 0 {
			x = 1
		} else if x < 0 {
			x -= 1
		} else {
			x += 1
		}

		amplitude := (pegConfig.MaxRadius - pegConfig.MinRadius) / math.Log(float64(boardConfig.NCols/2))
		return pegConfig.MinRadius + amplitude*math.Log(float64(x))

	case utils.PegGaussianDistHorizontal:
		x := column - boardConfig.NCols/2
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(float64(x-pegConfig.CenterFactor), 2))
		return (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss + pegConfig.MinRadius

	case utils.PegInverseGaussianDistHorizontal:
		x := column - boardConfig.NCols/2
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(float64(x-pegConfig.CenterFactor), 2))
		return pegConfig.MaxRadius - (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss

	case utils.PegSineDistHorizontal:
		sin := math.Pow(math.Sin(float64(column)*pegConfig.DeltaFactor), 2)
		return pegConfig.MinRadius + (pegConfig.MaxRadius-pegConfig.MinRadius)*sin

	// Vertical distributions
	case utils.PegLogarithmicDistVertical:
		y := row - boardConfig.NRows/2
		if y == 0 {
			y = 1
		} else if y < 0 {
			y -= 1
		} else {
			y += 1
		}

		amplitude := (pegConfig.MaxRadius - pegConfig.MinRadius) / math.Log(float64(boardConfig.NCols/2))
		return pegConfig.MinRadius + amplitude*math.Log(float64(y))

	case utils.PegGaussianDistVertical:
		y := row - boardConfig.NRows/2
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(float64(y-pegConfig.CenterFactor), 2))
		return (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss + pegConfig.MinRadius

	case utils.PegInverseGaussianDistVertical:
		y := row - boardConfig.NRows/2
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(float64(y-pegConfig.CenterFactor), 2))
		return pegConfig.MaxRadius - (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss

	case utils.PegSineDistVertical:
		sin := math.Pow(math.Sin(float64(row)*pegConfig.DeltaFactor), 2)
		return pegConfig.MinRadius + (pegConfig.MaxRadius-pegConfig.MinRadius)*sin

	default:
		log.Fatal("Invalid pegs distribution in the config file. Valid values are: \n" +
			"HORIZONTAL DISTRIBUTIONS\n" +
			"\t0: Uniform distribution\n" +
			"\t1: Logarithmic distribution\n" +
			"\t2: Gaussian distribution\n" +
			"\t3: Inverse Gaussian distribution\n" +
			"\t4: Sine distribution\n" +
			"VERTICAL DISTRIBUTIONS\n" +
			"\t5: Logarithmic distribution\n" +
			"\t6: Gaussian distribution\n" +
			"\t7: Inverse Gaussian distribution\n" +
			"\t8: Sine distribution")
	}

	return 0
}
