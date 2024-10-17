package entities

import (
	"go-galtonboard/utils"
	"log"
	"math"
)

type Cell struct {
	ParticlesIds []int
	PegsIds      []int
}

type Mesh struct {
	Rows    int
	Columns int
	Width   float64
	Height  float64
	dWidth  float64
	dHeight float64
	Cells   []Cell
}

func NewMesh(rows, columns int, width, height float64) *Mesh {
	mesh := Mesh{
		Rows:    rows,
		Columns: columns,
		Cells:   make([]Cell, columns*rows),
		Width:   width,
		Height:  height,
		dWidth:  width / float64(columns),
		dHeight: height / float64(rows),
	}

	return &mesh
}

func (m *Mesh) AddParticleToCell(x, y float64, particleType int, particleId int) {
	row := int(math.Ceil(y / m.dHeight))
	column := int(math.Ceil(x / m.dWidth))

	if row >= m.Rows {
		row = m.Rows - 1
	}

	if column >= m.Columns {
		column = m.Columns - 1
	}

	cellIndex := column*m.Rows + row
	if cellIndex >= len(m.Cells) || cellIndex < 0 {
		log.Fatal("AddParticleToCell: cellIndex out of bounds. Particle position: (", x, y, ") Row: ", row, " Column: ", column)
	}
	if particleType == utils.Peg {
		m.Cells[cellIndex].PegsIds = append(m.Cells[cellIndex].PegsIds, particleId)
	} else {
		m.Cells[cellIndex].ParticlesIds = append(m.Cells[cellIndex].ParticlesIds, particleId)
	}
}

func (m *Mesh) ClearMesh() {
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Columns; j++ {
			m.Cells[j*m.Rows+i].ParticlesIds = nil
		}
	}
}

func (m *Mesh) GetCell(row, column int) *Cell {
	if row < 0 || row >= m.Rows || column < 0 || column >= m.Columns {
		return nil
	}

	return &m.Cells[column*m.Rows+row]
}

func (m *Mesh) GetCellWithNeighbors(row, column int) []*Cell {
	var neighbors []*Cell

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if row+i < 0 || row+i >= m.Rows || column+j < 0 || column+j >= m.Columns {
				continue
			}

			neighbors = append(neighbors, &m.Cells[(column+j)*m.Rows+(row+i)])
		}
	}

	return neighbors
}
