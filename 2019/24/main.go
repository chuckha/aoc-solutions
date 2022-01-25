package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

func main() {
	lines := input.GetInput(2019, 24)
	grid := internal.NewGridV2FromLines(lines)

	// depth map
	grids := &grids{
		data: make(map[int]*internal.GridV2),
	}
	grids.data[0] = grid
	for i := 0; i < 10; i++ {
		grids = nextRecursiveItr(grids)
	}
	fmt.Println(grids.countBugs())
}

type grids struct {
	data map[int]*internal.GridV2
}

func (g *grids) String() string {
	var out strings.Builder
	mind := math.MaxInt
	maxd := math.MinInt
	for d := range g.data {
		if d > maxd {
			maxd = d
		}
		if d < mind {
			mind = d
		}
	}
	for i := mind; i <= maxd; i++ {
		d := g.GetDepth(i)
		out.WriteString(fmt.Sprintf("depth: %d\n", i))
		out.WriteString(d.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (g *grids) copy() *grids {
	out := &grids{
		data: make(map[int]*internal.GridV2),
	}
	for k, v := range g.data {
		out.data[k] = v.Copy()
	}
	return out
}

func (g *grids) GetDepth(i int) *internal.GridV2 {
	if _, ok := g.data[i]; !ok {
		g.data[i] = newCell()
	}
	return g.data[i]
}

func (g *grids) countBugs() int {
	sum := 0
	for _, g := range g.data {
		for _, v := range g.Data {
			if v == "#" {
				sum++
			}
		}
	}
	return sum
}

func nextRecursiveItr(grids *grids) *grids {
	middle := internal.Point{2, 2}
	//	for each point in every depth, calculate the next iteration of it
	// when a new depth is required ,initialize is at empty with the same size
	// initialize the above and below grids if we to check them
	for depth, grid := range grids.data {
		if needsToCheckAbove(grid) {
			grids.GetDepth(depth - 1)
		}
		if needsToCheckBelow(grid) {
			grids.GetDepth(depth + 1)
		}
	}
	c := grids.copy()

	for depth, grid := range grids.data {
		for p := range grid.Data {
			if p == middle {
				continue
			}
			count := 0
			for _, n := range p.Neighbors() {
				// go down if necessary
				if n == middle {
					count += countSubGrid(grids.GetDepth(depth+1), p)
					continue
				}

				// otherwise it's inside the grid so we can simply count it
				if item, ok := grid.Data[n]; ok {
					if item == "#" {
						count++
					}
					continue
				}
				// otherwise we need to go up and count the level above
				c := countAboveGrid(grids.GetDepth(depth-1), n)
				count += c
			}
			if grid.At(p) == "." {
				if count == 1 || count == 2 {
					c.GetDepth(depth).Set(p, "#")
				}
				continue
			}
			if grid.At(p) == "#" {
				if count != 1 {
					c.GetDepth(depth).Set(p, ".")
				}
				continue
			}
		}
	}
	return c
}

func countAboveGrid(g *internal.GridV2, belowPoint internal.Point) int {
	sum := 0
	// count the left and right
	if belowPoint.X < g.Min.X {
		if g.At(internal.Point{1, 2}) == "#" {
			sum++
		}
	}
	if belowPoint.X > g.Max.X {
		if g.At(internal.Point{3, 2}) == "#" {
			sum++
		}
	}
	//count the top and bottom
	if belowPoint.Y < g.Min.Y {
		if g.At(internal.Point{2, 1}) == "#" {
			sum++
		}
	}
	if belowPoint.Y > g.Max.Y {
		if g.At(internal.Point{2, 3}) == "#" {
			sum++
		}
	}
	return sum
}

func countSubGrid(g *internal.GridV2, abovePoint internal.Point) int {
	// hardcode this
	sum := 0
	if abovePoint == (internal.Point{2, 1}) {
		for i := g.Min.X; i <= g.Max.X; i++ {
			if g.At(internal.Point{i, g.Min.Y}) == "#" {
				sum++
			}
		}
	}
	if abovePoint == (internal.Point{1, 2}) {
		for j := g.Min.Y; j <= g.Max.Y; j++ {
			if g.At(internal.Point{g.Min.X, j}) == "#" {
				sum++
			}
		}
	}
	if abovePoint == (internal.Point{2, 3}) {
		for i := g.Min.X; i <= g.Max.X; i++ {
			if g.At(internal.Point{i, g.Max.Y}) == "#" {
				sum++
			}
		}
	}
	if abovePoint == (internal.Point{3, 2}) {
		for j := g.Min.Y; j <= g.Max.Y; j++ {
			if g.At(internal.Point{g.Max.X, j}) == "#" {
				sum++
			}
		}
	}
	return sum
}

func part1(grid *internal.GridV2) {
	layouts := map[string]struct{}{}
	done := false
	for !done {
		layouts[grid.String()] = struct{}{}
		grid = nextIter(grid)
		_, done = layouts[grid.String()]
	}
	fmt.Println(grid)
	fmt.Println(biodiversity(grid))
}

func nextIter(g *internal.GridV2) *internal.GridV2 {
	c := g.Copy()
	for p := range g.Data {
		count := 0
		for _, n := range p.Neighbors() {
			if item, ok := g.Data[n]; ok {
				if item == "#" {
					count++
				}
			}
		}

		// An empty space becomes infested with a bug if exactly one or two bugs are adjacent to it.
		if g.At(p) == "." {
			if count == 1 || count == 2 {
				c.Set(p, "#")
			}
			continue
		}
		if g.At(p) == "#" {
			if count != 1 {
				c.Set(p, ".")
			}
			continue
		}
	}
	return c
}

func biodiversity(g *internal.GridV2) int {
	sum := 0
	for j := g.Min.Y; j <= g.Max.Y; j++ {
		for i := g.Min.X; i <= g.Max.X; i++ {
			if g.At(internal.Point{i, j}) == "#" {
				pwr := (i + g.Max.X*j) + j
				sum += (1 << pwr)
			}
		}
	}
	return sum
}

func newCell() *internal.GridV2 {
	grid := internal.NewGridV2()
	for j := 0; j < 5; j++ {
		for i := 0; i < 5; i++ {
			grid.Set(internal.Point{i, j}, ".")
		}
	}
	grid.Set(internal.Point{2, 2}, "?")
	return grid
}

func needsToCheckAbove(g *internal.GridV2) bool {
	for p, v := range g.Data {
		if v == "#" {
			if p.X == g.Max.X || p.X == g.Min.X || p.Y == g.Min.Y || p.Y == g.Max.Y {
				return true
			}
		}
	}
	return false
}

func needsToCheckBelow(g *internal.GridV2) bool {
	if g.At(internal.Point{2, 1}) == "#" {
		return true
	}
	if g.At(internal.Point{1, 2}) == "#" {
		return true
	}
	if g.At(internal.Point{3, 2}) == "#" {
		return true
	}
	if g.At(internal.Point{2, 3}) == "#" {
		return true
	}
	return false
}

func canSkip(g *internal.GridV2) bool {
	for _, v := range g.Data {
		if v == "#" {
			return false
		}
	}
	return true
}
