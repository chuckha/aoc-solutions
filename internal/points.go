package internal

import "math"

type Point struct {
	X, Y int
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
