package internal

import (
	"fmt"
)

type Point struct {
	X, Y int
}

func (p Point) Direction(s string) Point {
	switch s {
	case "N", "up":
		return p.Up()
	case "NE":
		return p.Up().Right()
	case "NW":
		return p.Up().Left()
	case "S", "down":
		return p.Down()
	case "SE":
		return p.Down().Right()
	case "SW":
		return p.Down().Left()
	case "W", "left":
		return p.Left()
	case "E", "right":
		return p.Right()
	}
	panic("bad direction")
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
	return Abs(p.X-p2.X) + Abs(p.Y-p2.Y)
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (p Point) Until(p2 Point) []Point {
	if p.X != p2.X && p.Y != p2.Y {
		panic("requires two points to share at least one coordinate")
	}
	out := []Point{}
	switch {
	case p.X != p2.X:
		start := p.X
		end := p2.X
		if start > end {
			start, end = end, start
		}
		for i := start; i <= end; i++ {
			out = append(out, Point{X: i, Y: p.Y})
		}
	case p.Y != p2.Y:
		start := p.Y
		end := p2.Y
		if start > end {
			start, end = end, start
		}
		for i := start; i <= end; i++ {
			out = append(out, Point{X: p.X, Y: i})
		}
	}
	return out
}
