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

			x := float64(j) * (boardConfig.HorizontalSpace)
			y := float64(i) * (boardConfig.VerticalSpace)
			if i%2 != 0 {
				x += (boardConfig.HorizontalSpace) / 2

				if j == boardConfig.NCols-1 {
					continue
				}
			}

			peg.Position = utils.Point{x, y}
			peg.Radius = getRadius(&pegConfig, &boardConfig, y, x)
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

func getRadius(pegConfig *utils.PegConfig, boardConfig *utils.BoardConfig, row, column float64) float64 {
	yMiddle := boardConfig.VerticalSpace * float64(boardConfig.NRows/2)
	xMiddle := boardConfig.HorizontalSpace * float64(boardConfig.NCols/2)

	switch pegConfig.Distribution {
	case utils.PegUniformDist:
		return pegConfig.MinRadius

	// Horizontal distributions
	case utils.PegLogarithmicDistHorizontal:
		x := math.Abs(column - xMiddle)
		amplitude := (pegConfig.MaxRadius - pegConfig.MinRadius) / math.Log(xMiddle)
		return pegConfig.MinRadius + amplitude*math.Log(x+1)

	case utils.PegGaussianDistHorizontal:
		centerFactor := float64(pegConfig.CenterFactor) * boardConfig.HorizontalSpace
		x := column - xMiddle
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(x-centerFactor, 2))
		return (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss + pegConfig.MinRadius

	case utils.PegInverseGaussianDistHorizontal:
		centerFactor := float64(pegConfig.CenterFactor) * boardConfig.HorizontalSpace
		x := column - xMiddle
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(x-centerFactor, 2))
		return pegConfig.MaxRadius - (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss

	case utils.PegSineDistHorizontal:
		sin := math.Pow(math.Sin(column*pegConfig.DeltaFactor), 2)
		return pegConfig.MinRadius + (pegConfig.MaxRadius-pegConfig.MinRadius)*sin

	// Vertical distributions
	case utils.PegLogarithmicDistVertical:
		y := math.Abs(row - yMiddle)
		amplitude := (pegConfig.MaxRadius - pegConfig.MinRadius) / math.Log(yMiddle)
		return pegConfig.MinRadius + amplitude*math.Log(y+1)

	case utils.PegGaussianDistVertical:
		centerFactor := float64(pegConfig.CenterFactor) * boardConfig.VerticalSpace
		y := row - yMiddle
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(y-centerFactor, 2))
		return (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss + pegConfig.MinRadius

	case utils.PegInverseGaussianDistVertical:
		centerFactor := float64(pegConfig.CenterFactor) * boardConfig.VerticalSpace
		y := row - yMiddle
		gauss := math.Exp(-pegConfig.DeltaFactor * math.Pow(y-centerFactor, 2))
		return pegConfig.MaxRadius - (pegConfig.MaxRadius-pegConfig.MinRadius)*gauss

	case utils.PegSineDistVertical:
		sin := math.Pow(math.Sin(row*pegConfig.DeltaFactor), 2)
		return pegConfig.MinRadius + (pegConfig.MaxRadius-pegConfig.MinRadius)*sin

	// Other distributions
	case utils.SphericDist:
		distanceX := math.Abs(column - xMiddle)
		distanceY := math.Abs(row - yMiddle)
		distance := math.Sqrt(math.Pow(distanceX, 2) + math.Pow(distanceY, 2))
		if distance <= pegConfig.DeltaFactor {
			return pegConfig.MaxRadius
		}
		return pegConfig.MinRadius

	case utils.SphericGaussianDist:
		return pegConfig.MinRadius + (pegConfig.MaxRadius-pegConfig.MinRadius)*rand.NormFloat64()

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
