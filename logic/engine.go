package logic

import (
	"go-galtonboard/entities"
	"go-galtonboard/models"
	"go-galtonboard/utils"
	"log"
	"sync"
)

type Engine struct {
	Configs utils.Configs

	Particles []*entities.Particle
	Pegs      []*entities.Particle
	Border    []*utils.Point

	Model model.PhysicsModel
	Mesh  entities.Mesh

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
		Mesh:              *entities.NewMesh(config.BoardConfig.NRows, config.BoardConfig.NCols, borders[1][0], borders[1][1]),
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

	for i := 0; i < len(e.Pegs); i++ {
		p := e.Pegs[i]
		e.Mesh.AddParticleToCell(p.Position[0], p.Position[1], utils.Peg, i)
	}

	for i := 0; i < e.Configs.EngineConfig.MaxSteps; i++ {
		isStopped := e.ValidateStop()
		if isStopped {
			break
		}

		for j := 0; j < e.Configs.EngineConfig.SubSteps; j++ {
			e.applyForces()
			e.updateBodies(t, dtt)
			//e.validateCollisions()
			//e.validateConstraints()
			e.updateMesh()
			e.validateCollisionsMesh()
			e.validateConstraintsMesh()

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

func (e *Engine) updateMesh() {
	e.Mesh.ClearMesh()
	for i := 0; i < len(e.Particles); i++ {
		p := e.Particles[i]
		e.Mesh.AddParticleToCell(p.Position[0], p.Position[1], utils.Particle, i)
	}
}

func (e *Engine) updateBodies(t, dt float64) {
	for _, p := range e.Particles {
		if p.IsStopped {
			continue
		}

		e.Model.UpdateBall(p, t, dt)
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

func (e *Engine) validateCollisionsMesh() {
	rows := e.Configs.BoardConfig.NRows
	cols := e.Configs.BoardConfig.NCols

	sliceCount := e.Configs.EngineConfig.ThreadCount
	sliceWidth := rows / sliceCount
	sliceHeight := cols / sliceCount

	maxWidth := 0
	maxHeight := 0
	var wg sync.WaitGroup
	for i := 0; i < sliceCount; i++ {
		for j := 0; j < sliceCount; j++ {
			wg.Add(1)
			go e.solveCollisionThreaded(i*sliceWidth, (i+1)*sliceWidth, j*sliceHeight, (j+1)*sliceHeight, &wg)

			maxWidth = (i + 1) * sliceWidth
			maxHeight = (j + 1) * sliceHeight
		}
	}

	if maxWidth < rows {
		wg.Add(1)
		go e.solveCollisionThreaded(maxWidth, rows, 0, cols, &wg)
	}

	if maxHeight < cols {
		wg.Add(1)
		go e.solveCollisionThreaded(0, maxWidth, maxHeight, cols, &wg)
	}
	wg.Wait()
}

func (e *Engine) solveCollisionThreaded(iStart, iEnd, jStart, jEnd int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := iStart; i < iEnd; i++ {
		for j := jStart; j < jEnd; j++ {
			e.processCell(e.Mesh.GetCell(i, j), i, j)
		}
	}
}

func (e *Engine) processCell(c *entities.Cell, i, j int) {
	for _, pId := range c.ParticlesIds {
		if e.Particles[pId].IsStopped {
			continue
		}

		e.checkAtomCellCollisions(pId, c)
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i-1, j))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i+1, j))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i, j-1))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i, j+1))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i-1, j-1))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i-1, j+1))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i+1, j-1))
		e.checkAtomCellCollisions(pId, e.Mesh.GetCell(i+1, j+1))
	}
}

func (e *Engine) checkAtomCellCollisions(particleId int, c *entities.Cell) {
	if c == nil {
		return
	}

	p := e.Particles[particleId]
	for _, pegId := range c.PegsIds {
		peg := e.Pegs[pegId]
		distanceSquare := utils.DistanceSquare(&p.Position, &peg.Position)
		if distanceSquare < (p.Radius+peg.Radius)*(p.Radius+peg.Radius) {
			e.Model.ResolveCollision(p, peg)
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

func (e *Engine) validateConstraintsMesh() {
	rows := e.Configs.BoardConfig.NRows
	cols := e.Configs.BoardConfig.NCols

	for j := 0; j < cols; j++ {
		e.processCellConstraints(e.Mesh.GetCell(0, j))
		e.processCellConstraints(e.Mesh.GetCell(rows-1, j))
	}

	for i := 1; i < rows; i++ {
		e.processCellConstraints(e.Mesh.GetCell(i, 0))
		e.processCellConstraints(e.Mesh.GetCell(i, cols-1))
	}
}

func (e *Engine) processCellConstraints(cell *entities.Cell) {
	if cell == nil {
		return
	}

	for _, pId := range cell.ParticlesIds {
		p := e.Particles[pId]
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
