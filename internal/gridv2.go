package internal

import (
	"strings"
)

type GridV2 struct {
	Data        map[Point]string
	Min, Max    Point
	DefaultChar string
}

func NewGridV2() *GridV2 {
	return &GridV2{
		Data:        map[Point]string{},
		DefaultChar: ".",
	}
}
func NewGridV2WithDefaultChar(s string) *GridV2 {
	g := NewGridV2()
	g.DefaultChar = s
	return g
}

func NewGridV2FromLines(lines []string) *GridV2 {
	grid := NewGridV2()
	for j, line := range lines {
		for i, c := range line {
			grid.Set(Point{i, j}, string(c))
		}
	}
	return grid
}

func (g *GridV2) Set(p Point, s string) {
	if p.X < g.Min.X {
		g.Min.X = p.X
	}
	if p.X > g.Max.X {
		g.Max.X = p.X
	}
	if p.Y < g.Min.Y {
		g.Min.Y = p.Y
	}
	if p.Y > g.Max.Y {
		g.Max.Y = p.Y
	}
	g.Data[p] = s
}
func (g *GridV2) At(p Point) string {
	if item, ok := g.Data[p]; ok {
		return item
	}
	return g.DefaultChar
}

func (g *GridV2) String() string {
	var out strings.Builder
	for j := g.Min.Y; j <= g.Max.Y; j++ {
		for i := g.Min.X; i <= g.Max.X; i++ {
			out.WriteString(g.At(Point{i, j}))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (g *GridV2) EntryCount() int {
	return len(g.Data)
}

func (g *GridV2) Copy() *GridV2 {
	gc := NewGridV2()
	for k, v := range g.Data {
		gc.Data[k] = v
	}
	gc.Max = g.Max
	gc.Min = g.Min
	gc.DefaultChar = g.DefaultChar
	return gc
}
