package logic

import (
	"bufio"
	"errors"
	"fmt"
	"go-galtonboard/entities"
	"go-galtonboard/utils"
	"log"
	"os"
)

type Exporter struct {
	path   string
	file   *os.File
	writer *bufio.Writer
}

func NewExporter(path string) *Exporter {
	return &Exporter{
		path: path,
	}
}

func (e *Exporter) CreateFile(name string) {
	fileName := e.path + name + "-0.csv"
	for i := 0; fileExist(fileName); i++ {
		fileName = e.path + name + fmt.Sprintf("-%d", i) + ".csv"
	}

	file, err := os.Create(fileName)
	writer := bufio.NewWriterSize(file, 128*1024*4)

	if err != nil {
		panic(err)
	}

	e.file = file
	e.writer = writer
}

func (e *Exporter) CloseFile() {
	err := e.writer.Flush()
	if err != nil {
		log.Fatal("Error flushing the writer")
	}
	err = e.file.Close()
	if err != nil {
		log.Fatal("Error closing the file")
	}
}

func (e *Exporter) Write(content string) {
	_, err := e.writer.WriteString(content)

	if err != nil {
		panic(err)
	}
}

func (e *Exporter) WritePath(particles, pegs []*entities.Particle, borders []*utils.Point) {
	total := len(particles) + len(borders) + len(pegs)
	e.Write(getExportHeader(total))

	for i, sphere := range particles {
		e.Write(getExportPath(i, sphere))
	}

	temp := len(particles)
	for i, sphere := range pegs {
		e.Write(getExportPath(i+temp, sphere))
	}

	temp = len(particles) + len(pegs)
	for i, border := range borders {
		e.Write(getExportPathBorders(i+temp, border))
	}
}

func (e *Exporter) WriteHistogram(counts []int) {
	e.Write(getExportHistogram(counts))
}

func getExportHistogram(counts []int) string {
	content := ""
	for i := 0; i < len(counts); i++ {
		content += fmt.Sprintf("%d\t%d\n", i+1, counts[i])
	}

	return content
}

func getExportPath(number int, sphere *entities.Particle) string {
	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		number,
		sphere.Type,
		sphere.Position[0],
		sphere.Position[1],
		sphere.Radius,
	)

	return content
}

func getExportPathBorders(number int, point *utils.Point) string {
	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		number,
		2,
		point[0],
		point[1],
		0.5,
	)

	return content
}

func getExportHeader(total int) string {
	content := fmt.Sprintf("%d\naver\n", total)
	return content
}

func fileExist(name string) bool {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}
