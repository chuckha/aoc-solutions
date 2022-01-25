package internal

import (
	"fmt"
	"strings"
)

type stringer interface {
	String() string
}

type GridV3[T stringer] struct {
	Data     map[Point]T
	Min, Max Point
}

func NewGridV3[T stringer]() *GridV3[T] {
	return &GridV3[T]{
		Data: map[Point]T{},
	}
}

func (g *GridV3[T]) Set(p Point, t T) {
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
	g.Data[p] = t
}
func (g *GridV3[T]) At(p Point) T {
	if item, ok := g.Data[p]; ok {
		return item
	}
	panic("does not exist")
}

func (g *GridV3[T]) String() string {
	var out strings.Builder
	for j := g.Min.Y; j <= g.Max.Y; j++ {
		for i := g.Min.X; i <= g.Max.X; i++ {
			out.WriteString(fmt.Sprintf("%v", g.At(Point{i, j})))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (g *GridV3[T]) EntryCount() int {
	return len(g.Data)
}

func (g *GridV3[T]) Copy() *GridV3[T] {
	gc := NewGridV3[T]()
	for k, v := range g.Data {
		gc.Data[k] = v
	}
	gc.Max = g.Max
	gc.Min = g.Min
	return gc
}
