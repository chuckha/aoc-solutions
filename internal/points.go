package internal

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y int
}

func (p Point) Right() Point {
	return Point{p.X + 1, p.Y}
}
func (p Point) Down() Point {
	return Point{p.X, p.Y + 1}
}
func (p Point) Left() Point {
	return Point{p.X - 1, p.Y}
}
func (p Point) Up() Point {
	return Point{p.X, p.Y - 1}
}

func (p Point) Neighbors() []Point {
	return []Point{
		{X: p.X, Y: p.Y - 1},
		{X: p.X + 1, Y: p.Y},
		{X: p.X, Y: p.Y + 1},
		{X: p.X - 1, Y: p.Y},
	}
}

func (p Point) Surrounding() []Point {
	return []Point{
		{X: p.X - 1, Y: p.Y - 1}, {X: p.X, Y: p.Y - 1}, {X: p.X + 1, Y: p.Y - 1},
		{X: p.X - 1, Y: p.Y}, {X: p.X + 1, Y: p.Y},
		{X: p.X - 1, Y: p.Y + 1}, {X: p.X, Y: p.Y + 1}, {X: p.X + 1, Y: p.Y + 1},
	}
}

func (p Point) ManhattanDistance(p2 Point) int {
	return int(math.Abs(float64(p.X-p2.X))) + int(math.Abs(float64(p.Y-p2.Y)))
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}
