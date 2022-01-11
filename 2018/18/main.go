package main

import (
	"fmt"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	g := internal.NewGridFromInput(lines)
	for i := 0; i < 1000000; i++ {
		g = turn(g)
		//		fmt.Println(g)
		trees, ly := resourceCounts(g)
		fmt.Println(i, i%35, trees*ly)
	}
	//	fmt.Println(resourceCounts(g))
}

func resourceCounts(g *internal.Grid) (int, int) {
	treeCount, lumberyardCount := 0, 0
	for _, c := range g.Data {
		if c == "|" {
			treeCount++
		}
		if c == "#" {
			lumberyardCount++
		}
	}
	return treeCount, lumberyardCount
}

func turn(g *internal.Grid) *internal.Grid {
	out := g.Copy()
	for k, v := range g.Data {
		openCount := 0
		lumberYardCount := 0
		treeCount := 0
		for _, n := range k.Surrounding() {
			switch g.At(n.X, n.Y) {
			case ".":
				openCount++
			case "#":
				lumberYardCount++
			case "|":
				treeCount++
			case "":
			default:
				panic("bad input: " + g.At(n.X, n.Y))
			}
		}
		if v == "." {
			if treeCount >= 3 {
				out.Set(k.X, k.Y, "|")
			}
		}
		if v == "|" {
			if lumberYardCount >= 3 {
				out.Set(k.X, k.Y, "#")
			}
		}
		if v == "#" {
			if lumberYardCount >= 1 && treeCount >= 1 {
				continue
			}
			out.Set(k.X, k.Y, ".")
		}
	}
	return out
}
