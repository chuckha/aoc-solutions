package internal

import (
	"math"
	"strings"
)

type Grid struct {
	Data           map[Point]string
	MinX, MinY     int
	Length, Height int
	DefaultChar    string
}

func (g *Grid) String() string {
	rows := []string{}
	for j := g.MinY; j < g.Height; j++ {
		var row strings.Builder
		for i := g.MinX; i < g.Length; i++ {
			row.WriteString(g.Data[Point{i, j}])
		}
		rows = append(rows, row.String())
	}
	return strings.Join(rows, "\n")
}

func (g *Grid) At(i, j int) string {
	if v, ok := g.Data[Point{i, j}]; ok {
		return v
	}
	return g.DefaultChar
}

func NewGridFromInput(lines []string) *Grid {
	width := len(lines[0])
	height := len(lines)
	g := &Grid{
		Data:        map[Point]string{},
		Length:      width,
		Height:      height,
		DefaultChar: ".",
	}
	for y, line := range lines {
		for x, c := range line {
			g.Data[Point{x, y}] = string(c)
		}
	}
	return g
}

func NewGrid(length, height int, defaultChar string) *Grid {
	g := &Grid{
		Data:   map[Point]string{},
		Length: length,
		Height: height,
	}
	for j := 0; j < g.Height; j++ {
		for i := 0; i < g.Length; i++ {
			g.Data[Point{i, j}] = defaultChar
		}
	}
	return g
}

// RotateRow shifts a row to the right at position j by 'by' places
func (g *Grid) RotateRow(j, by int) {
	newPoints := map[Point]string{}
	for i := 0; i < g.Length; i++ {
		newI := (i + by) % g.Length
		newPoints[Point{newI, j}] = g.Data[Point{i, j}]
	}
	for k, v := range newPoints {
		g.Data[k] = v
	}
}

// RotateCol shifts a column down at position i by 'by' places
func (g *Grid) RotateCol(i, by int) {
	newPoints := map[Point]string{}
	for j := 0; j < g.Height; j++ {
		newJ := (j + by) % g.Height
		newPoints[Point{i, newJ}] = g.Data[Point{i, j}]
	}
	for k, v := range newPoints {
		g.Data[k] = v
	}
}

// Rect fills in a rectangle starting wiht the top left with the "on" character
func (g *Grid) Rect(x, y int) {
	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			g.Data[Point{i, j}] = "#"
		}
	}
}

// On counts the number of grid points that are on
func (g *Grid) On() int {
	count := 0
	for j := g.MinY; j < g.Height; j++ {
		for i := g.MinX; i < g.Length; i++ {
			if g.Data[Point{i, j}] == "#" {
				count++
			}
		}
	}
	return count
}

// Rotate rotates the grid 90 degrees to the right (clockwise)
func (g *Grid) Rotate() {
	newdata := make(map[Point]string)
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Length; x++ {
			newdata[Point{(g.Length - 1) - y, x}] = g.At(x, y)
		}
	}
	g.Height, g.Length = g.Length, g.Height
	g.Data = newdata
}

func (g *Grid) Flip() {
	newdata := make(map[Point]string)
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Length; x++ {
			newdata[Point{(g.Length - 1) - x, y}] = g.At(x, y)
		}
	}
	g.Data = newdata
}

// One line returns the string value of the grid where rows are separated by /
func (g *Grid) OneLine() string {
	return strings.Replace(g.String(), "\n", "/", -1)
}

func (g *Grid) SubGrid(minx, maxx, miny, maxy int) *Grid {
	out := &Grid{
		Length: maxx + 1 - minx,
		Height: maxy + 1 - miny,
		Data:   make(map[Point]string),
	}
	for y := g.MinY; y <= maxy-miny; y++ {
		for x := g.MinX; x <= maxx-minx; x++ {
			out.Data[Point{x, y}] = g.At(x+minx, y+miny)
		}
	}
	return out
}

// Dijkstra
// cost assumed to be 1
// wall assumed to be #
// open spot assumed to be non # character
func (g *Grid) Dijkstra(start Point) map[Point]int {
	visited := map[Point]struct{}{}
	costs := map[Point]int{}
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Length; x++ {
			costs[Point{x, y}] = math.MaxInt
		}
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
			if g.Data[n] == "#" {
				continue
			}
			if costs[cur]+1 < costs[n] {
				costs[n] = costs[cur] + 1
			}
			q.Enqueue(n)
		}
		visited[cur] = struct{}{}
	}
	return costs
}

func (g *Grid) Center() Point {
	return Point{g.Length / 2, g.Height / 2}
}

func (g *Grid) Set(i, j int, val string) {
	g.Data[Point{i, j}] = val
}
