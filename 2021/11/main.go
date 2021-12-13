package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	maxY := len(lines) - 1
	maxX := len(lines[0]) - 1
	octos := newOctos(maxX, maxY)

	for y, line := range lines {
		for x, o := range line {
			nrg, _ := strconv.Atoi(string(o))
			octos.data[point{x, y}] = &octo{x: x, y: y, energy: nrg}
		}
	}
	for i := 0; ; i++ {
		octos.iterate()
		if octos.allFlashed() {
			fmt.Println(i + 1)
			break
		}
		octos.resetFlashed()
	}
	//	fmt.Println(octos.flashes)
}

type octos struct {
	maxX, maxY int
	flashes    int
	data       map[point]*octo
}

func newOctos(maxX, maxY int) *octos {
	return &octos{
		maxX: maxX,
		maxY: maxY,
		data: make(map[point]*octo),
	}
}

func (o *octos) allFlashed() bool {
	for j := 0; j <= o.maxY; j++ {
		for i := 0; i <= o.maxX; i++ {
			if !o.data[point{i, j}].flashed {
				return false
			}
		}
	}
	return true
}

func (o *octos) iterate() {
	for j := 0; j <= o.maxY; j++ {
		for i := 0; i <= o.maxX; i++ {
			o.inc(i, j)
		}
	}
}

func (o *octos) resetFlashed() {
	for j := 0; j <= o.maxY; j++ {
		for i := 0; i <= o.maxX; i++ {
			o.data[point{i, j}].resetFlashed()
		}
	}
}

func (o *octos) inc(x, y int) {
	if o.data[point{x, y}].inc() == 1 {
		o.flashes++
		o.flash(x, y)
	}
}

func (o *octos) flash(x, y int) {
	for _, n := range o.neighbors(x, y) {
		o.inc(n.x, n.y)
	}
}

func (o *octos) neighbors(x, y int) []*octo {
	out := []*octo{}
	if x != 0 {
		out = append(out, o.data[point{x - 1, y}])
		if y != 0 {
			out = append(out, o.data[point{x - 1, y - 1}])
		}
		if y != o.maxY {
			out = append(out, o.data[point{x - 1, y + 1}])
		}
	}
	if x != o.maxX {
		out = append(out, o.data[point{x + 1, y}])
		if y != 0 {
			out = append(out, o.data[point{x + 1, y - 1}])
		}
		if y != o.maxY {
			out = append(out, o.data[point{x + 1, y + 1}])
		}
	}
	if y != 0 {
		out = append(out, o.data[point{x, y - 1}])
	}
	if y != o.maxY {
		out = append(out, o.data[point{x, y + 1}])
	}
	return out
}
func (o *octos) String() string {
	outs := []string{}
	for j := 0; j <= o.maxY; j++ {
		row := []string{}
		for i := 0; i <= o.maxX; i++ {
			row = append(row, fmt.Sprintf("%d", o.data[point{i, j}].energy))
		}
		outs = append(outs, strings.Join(row, ""))
	}
	return strings.Join(outs, "\n")
}

type point struct {
	x, y int
}
type octo struct {
	x, y    int
	flashed bool
	energy  int
}

func (o *octo) resetFlashed() {
	o.flashed = false
}

func (o *octo) inc() int {
	if o.flashed {
		return 0
	}
	o.energy += 1
	if o.energy > 9 {
		o.energy = 0
		o.flashed = true
		return 1
	}
	return 0
}

/*
sum(unique factors) * 10

1 -> 1          sum = 1   * 10
2 -> 1, 2           = 3
3 -> 1, 3           = 4
4 -> 1, 2, 4        = 7
5 -> 1, 5           = 6
6 -> 1, 2, 3, 6     = 12
7 -> 1, 7           = 8
8 -> 1, 2, 4, 8     = 15
9 -> 1, 3, 9        = 13
*/
