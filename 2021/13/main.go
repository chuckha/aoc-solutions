package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	gatheringPoints := true
	folds := []fold{}
	instructions := newPaper()
	for _, line := range lines {
		if strings.HasPrefix(line, "fold along") {
			gatheringPoints = false
		}

		if !gatheringPoints {
			words := strings.Split(line, "=")
			along, _ := strconv.Atoi(words[1])
			folds = append(folds, fold{
				dir:   string(words[0][len(words[0])-1]),
				along: along,
			})
			continue
		}
		points := strings.Split(line, ",")
		x, _ := strconv.Atoi(points[0])
		if x > instructions.maxX {
			instructions.maxX = x
		}
		y, _ := strconv.Atoi(points[1])
		if y > instructions.maxY {
			instructions.maxY = y
		}
		instructions.fill(point{x, y}, "#")
	}
	for _, fold := range folds {
		if fold.dir == "y" {
			instructions.foldY(fold.along)
			continue
		}
		instructions.foldX(fold.along)
	}
	fmt.Println(instructions)
	//	fmt.Println("----")
	//	fmt.Println(instructions.numberVisible())
	//	fmt.Println(instructions)
}

type fold struct {
	dir   string
	along int
}

type point struct {
	x, y int
}

type paper struct {
	maxX, maxY int
	data       map[point]string
}

func newPaper() *paper {
	return &paper{
		data: make(map[point]string),
	}
}

func (p *paper) fill(point point, s string) {
	p.data[point] = s
}

func (p *paper) has(point point) bool {
	_, ok := p.data[point]
	return ok
}
func (p *paper) numberVisible() int {
	c := 0
	for j := 0; j <= p.maxY; j++ {
		for i := 0; i <= p.maxX; i++ {
			if _, ok := p.data[point{i, j}]; ok {
				c++
			}
		}
	}
	return c
}

func (p *paper) String() string {
	lines := []string{}
	for j := 0; j <= p.maxY; j++ {
		line := make([]string, p.maxX+1)
		for i := 0; i <= p.maxX; i++ {
			if _, ok := p.data[point{i, j}]; ok {
				line[i] = "#"
				continue
			}
			line[i] = "."
		}
		lines = append(lines, strings.Join(line, ""))
	}
	return strings.Join(lines, "\n")
}

func (p *paper) foldY(y int) {
	folded := newPaper()
	folded.maxX = p.maxX
	folded.maxY = y - 1
	for j := y + 1; j <= p.maxY; j++ {
		for i := 0; i <= p.maxX; i++ {
			if _, ok := p.data[point{i, p.maxY - j}]; ok {
				folded.data[point{i, p.maxY - j}] = "#"
			}
			if _, ok := p.data[point{i, j}]; ok {
				folded.data[point{i, p.maxY - j}] = "#"
			}
		}
	}
	p.data = folded.data
	p.maxX = folded.maxX
	p.maxY = folded.maxY
}

func (p *paper) foldX(x int) {
	folded := newPaper()
	folded.maxX = x - 1
	folded.maxY = p.maxY
	for j := 0; j <= p.maxY; j++ {
		for i := x + 1; i <= p.maxX; i++ {
			if _, ok := p.data[point{p.maxX - i, j}]; ok {
				folded.data[point{p.maxX - i, j}] = "#"
			}
			if _, ok := p.data[point{i, j}]; ok {
				folded.data[point{p.maxX - i, j}] = "#"
			}
		}
	}
	p.data = folded.data
	p.maxX = folded.maxX
	p.maxY = folded.maxY

}
