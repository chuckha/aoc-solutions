package internal

import (
	"fmt"
	"math"
	"strings"
)

type stringer interface {
	String() string
}

type GridV3[T stringer] struct {
	Data          map[Point]T
	Min, Max      Point
	DefaultOutput string
}

func NewGridV3[T stringer]() *GridV3[T] {
	return &GridV3[T]{
		Data:          map[Point]T{},
		Min:           Point{X: math.MaxInt, Y: math.MaxInt},
		Max:           Point{X: math.MinInt, Y: math.MinInt},
		DefaultOutput: ".",
	}
}

func (g *GridV3[T]) Clear(p Point) {
	delete(g.Data, p)
}

func (g *GridV3[T]) In(p Point) bool {
	_, ok := g.Data[p]
	return ok
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
	panic(fmt.Sprintf("does not exist %v", p))
}

func (g *GridV3[T]) String() string {
	var out strings.Builder
	for j := g.Min.Y; j <= g.Max.Y; j++ {
		for i := g.Min.X; i <= g.Max.X; i++ {
			if !g.In(Point{i, j}) {
				out.WriteString(g.DefaultOutput)
				continue
			}
			out.WriteString(fmt.Sprintf("%v", g.At(Point{i, j})))
		}
		out.WriteString("\n")
	}
	return out.String()
}

// Layer is mostly used for printing
func (g *GridV3[T]) Layer(top *GridV3[T]) *GridV3[T] {
	out := g.Copy()
	min := Point{X: g.Min.X, Y: g.Min.Y}
	max := Point{X: g.Max.X, Y: g.Max.Y}
	if top.Min.X < min.X {
		min.X = top.Min.X
	}
	if top.Min.Y < min.Y {
		min.Y = top.Min.Y
	}
	if top.Max.X > max.X {
		max.X = top.Max.X
	}
	if top.Max.Y > max.Y {
		max.Y = top.Max.Y
	}
	for p, t := range top.Data {
		out.Set(p, t)
	}
	return out
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
	gc.DefaultOutput = g.DefaultOutput
	return gc
}

func (g *GridV3[T]) CopyAsString() *GridV3[Cell] {
	gc := NewGridV3[Cell]()
	for k, v := range g.Data {
		gc.Data[k] = Cell(v.String())
	}
	gc.Max = g.Max
	gc.Min = g.Min
	gc.DefaultOutput = g.DefaultOutput
	return gc
}

type Cell string

func (c Cell) String() string { return string(c) }

func (g *GridV3[T]) All(fn func(pt Point)) {
	for pt := range g.Data {
		fn(pt)
	}
}

func (g *GridV3[T]) Each(fn func(T)) {
	g.EachV2(func(_ Point, t T) { fn(t) })
}

func (g *GridV3[T]) EachV2(fn func(pt Point, t T)) {
	for j := g.Min.Y; j <= g.Max.Y; j++ {
		for i := g.Min.X; i <= g.Max.X; i++ {
			p := Point{i, j}
			if !g.In(p) {
				continue
			}
			fn(p, g.At(p))
		}
	}
}

func (g *GridV3[T]) Dijkstra(start Point, connected func(c, n T) bool, cost func(c, n T) int) map[Point]int {
	visited := map[Point]struct{}{}
	costs := map[Point]int{}
	for p := range g.Data {
		costs[p] = math.MaxInt
	}
	costs[start] = 0
	q := NewQueue[Point]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		if _, ok := visited[cur]; ok {
			continue
		}
		for _, n := range cur.Neighbors() {
			if _, ok := g.Data[n]; !ok {
				continue
			}
			//			fmt.Println(cur, g.At(cur), n, g.At(n), connected(g.At(cur), g.At(n)))
			if !connected(g.At(cur), g.At(n)) {
				continue
			}
			//			fmt.Println(g.At(cur), g.At(n))
			cost := cost(g.At(cur), g.At(n))
			//			fmt.Println(cost)
			if costs[cur]+cost < costs[n] {
				costs[n] = costs[cur] + cost
			}
			q.Enqueue(n)
		}
		visited[cur] = struct{}{}
	}
	return costs
}
