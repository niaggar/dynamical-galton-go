package logic

import (
	"go-galtonboard/entities"
	"go-galtonboard/models"
	"go-galtonboard/utils"
	"log"
)

type Engine struct {
	Configs utils.Configs

	Particles []*entities.Particle
	Pegs      []*entities.Particle
	Border    []*utils.Point

	Model model.PhysicsModel

	HorizontalMax float64
	HorizontalMin float64
	VerticalMax   float64
	VerticalMin   float64

	PathExporter      *Exporter
	HistogramExporter *Exporter

	HistogramCount []int
}

// NewEngine returns a new logic with the given values.
func NewEngine(config utils.Configs, route string) *Engine {
	pegs, borders := entities.NewPegs(config.PegConfig, config.BoardConfig)
	particles := entities.NewParticles(config.ParticleConfig, borders[0])

	var (
		pathExporter, histogramExporter *Exporter
	)

	if config.SaveConfig.SavePaths {
		pathExporter = NewExporter(route)
		pathExporter.CreateFile("paths")
	}

	if config.SaveConfig.SaveHistogram {
		histogramExporter = NewExporter(route)
		histogramExporter.CreateFile("histogram")
	}

	return &Engine{
		Configs:           config,
		Particles:         particles,
		Pegs:              pegs,
		Border:            borders,
		Model:             model.NewDefaultModel(),
		PathExporter:      pathExporter,
		HistogramExporter: histogramExporter,
		HorizontalMax:     borders[1][0],
		HorizontalMin:     borders[3][0],
		VerticalMax:       borders[1][1],
		VerticalMin:       borders[3][1],
		HistogramCount:    make([]int, config.BoardConfig.NCols-1),
	}
}

// Run runs the logic.
func (e *Engine) Run() {
	t := 0.0
	dtt := e.Configs.EngineConfig.Dt / float64(e.Configs.EngineConfig.SubSteps)

	for i := 0; i < e.Configs.EngineConfig.MaxSteps; i++ {
		isStopped := e.ValidateStop()
		if isStopped {
			break
		}

		for j := 0; j < e.Configs.EngineConfig.SubSteps; j++ {
			e.applyForces()
			e.updateBodies(t, dtt)
			e.validateCollisions()
			e.validateConstraints()
			t += dtt
		}

		if e.Configs.SaveConfig.SavePaths {
			e.PathExporter.WritePath(e.Particles, e.Pegs, e.Border)
		}
	}

	if e.Configs.SaveConfig.SavePaths {
		e.PathExporter.CloseFile()
	}

	if e.Configs.SaveConfig.SaveHistogram {
		e.HistogramExporter.WriteHistogram(e.HistogramCount)
		e.HistogramExporter.CloseFile()

		totalSumHistogram := 0
		for _, count := range e.HistogramCount {
			totalSumHistogram += count
		}

		log.Printf("Total sum histogram: %d", totalSumHistogram)
	}
}

func (e *Engine) ValidateStop() bool {
	nStopped := 0
	for _, p := range e.Particles {
		if p.IsStopped {
			nStopped++
		}
	}

	if nStopped == len(e.Particles) {
		log.Println("All particles stopped")
		return true
	}

	return false
}

func (e *Engine) applyForces() {
	for _, p := range e.Particles {
		if p.IsStopped {
			continue
		}

		p.Acceleration = utils.Point{0, -9.81}
	}
}

func (e *Engine) updateBodies(t, dt float64) {
	for _, p := range e.Particles {
		if p.IsStopped {
			continue
		}

		e.Model.UpdateBall(p, t, dt)
	}

	for _, p := range e.Pegs {
		e.Model.UpdatePeg(p, t, dt)
	}
}

func (e *Engine) validateCollisions() {
	for _, p := range e.Particles {
		if p.IsStopped {
			continue
		}

		for _, peg := range e.Pegs {
			distanceSquare := utils.DistanceSquare(&p.Position, &peg.Position)
			if distanceSquare < (p.Radius+peg.Radius)*(p.Radius+peg.Radius) {
				e.Model.ResolveCollision(p, peg)
			}
		}
	}
}

func (e *Engine) validateConstraints() {
	for _, p := range e.Particles {
		if p.IsStopped {
			continue
		}

		if p.Position[0]-p.Radius < e.HorizontalMin {
			p.Position[0] = e.HorizontalMin + p.Radius
			p.Velocity[0] = -p.Velocity[0] * p.Damping
		}

		if p.Position[0]+p.Radius > e.HorizontalMax {
			p.Position[0] = e.HorizontalMax - p.Radius
			p.Velocity[0] = -p.Velocity[0] * p.Damping
		}

		if p.Position[1]-p.Radius < e.VerticalMin {
			p.Position[1] = e.VerticalMin + p.Radius
			p.Velocity[1] = -p.Velocity[1] * p.Damping
			p.IsStopped = true

			x := p.Position[0] - e.HorizontalMin
			col := int(x / e.Configs.BoardConfig.HorizontalSpace)
			e.HistogramCount[col]++
		}

		if p.Position[1]+p.Radius > e.VerticalMax {
			p.Position[1] = e.VerticalMax - p.Radius
			p.Velocity[1] = -p.Velocity[1] * p.Damping
		}
	}
}
