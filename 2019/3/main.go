package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	first := move(strings.Split(lines[0], ","))
	second := move(strings.Split(lines[1], ","))
	merged := merge(first, second)
	fmt.Println(smallestSteps(merged))
}

func smallestSteps(a *grid) int {
	smallest := math.MaxInt
	for _, v := range a.data {
		if v.steps < smallest {
			smallest = v.steps
		}
	}
	return smallest
}

func merge(a, b *grid) *grid {
	merged := newGrid()
	for j := min(a.min.y, b.min.y); j <= max(a.max.y, b.max.y); j++ {
		for i := min(a.min.x, b.min.x); i <= max(a.max.x, b.max.x); i++ {
			if i == 0 && j == 0 {
				continue
			}
			if fi, ok := a.data[point{i, j}]; ok {
				if si, ok := b.data[point{i, j}]; ok {
					merged.add(point{i, j}, data{out: "X", steps: fi.steps + si.steps})
				}
			}
		}
	}
	return merged
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func move(mvmt []string) *grid {
	g := newGrid()
	pos := point{0, 0}
	count := 0
	g.add(pos, data{out: "O", steps: count})
	for _, m := range mvmt {
		dist, _ := strconv.Atoi(m[1:])
		switch m[0] {
		case 'U':
			for i := 0; i < dist; i++ {
				count++
				pos = point{pos.x, pos.y - 1}
				if pos == (point{0, 0}) {
					continue
				}
				g.add(pos, data{out: "|", steps: count})
			}
		case 'R':
			for i := 0; i < dist; i++ {
				count++
				pos = point{pos.x + 1, pos.y}
				if pos == (point{0, 0}) {
					continue
				}
				g.add(pos, data{out: "-", steps: count})
			}
		case 'D':
			for i := 0; i < dist; i++ {
				count++
				pos = point{pos.x, pos.y + 1}
				if pos == (point{0, 0}) {
					continue
				}
				g.add(pos, data{out: "|", steps: count})
			}
		case 'L':
			for i := 0; i < dist; i++ {
				count++
				pos = point{pos.x - 1, pos.y}
				if pos == (point{0, 0}) {
					continue
				}
				g.add(pos, data{out: "-", steps: count})
			}
		}
		g.add(pos, data{out: "+", steps: count})
	}
	return g
}

func closestIntersectionDistance(g *grid) int {
	smallestDistance := math.MaxInt
	for p, c := range g.data {
		if c.out == "X" {
			if p.dist(point{0, 0}) < smallestDistance {
				smallestDistance = p.dist(point{0, 0})
			}
		}
	}
	return smallestDistance
}

type data struct {
	out   string
	steps int
}

type grid struct {
	min, max point
	data     map[point]data
}

func newGrid() *grid {
	return &grid{
		data: make(map[point]data),
	}
}

func (g *grid) add(p point, d data) {
	g._updateBounds(p)
	g.data[p] = d
}

func (g *grid) String() string {
	var out strings.Builder
	for j := g.min.y; j <= g.max.y; j++ {
		for i := g.min.x; i <= g.max.x; i++ {
			if item, ok := g.data[point{i, j}]; ok {
				out.WriteString(item.out)
				continue
			}
			out.WriteString(".")
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (g *grid) _updateBounds(p point) {
	if p.x < g.min.x {
		g.min.x = p.x
	}
	if p.x > g.max.x {
		g.max.x = p.x
	}
	if p.y < g.min.y {
		g.min.y = p.y
	}
	if p.y > g.max.y {
		g.max.y = p.y
	}
}

type point struct {
	x, y int
}

func (p point) dist(b point) int {
	return abs(p.x-b.x) + abs(p.y-b.y)
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
