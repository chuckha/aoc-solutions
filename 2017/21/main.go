package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	rules := map[string]string{}
	for _, line := range lines {
		rule := parseLine(line)
		rules[rule.from] = rule.to
	}

	// {0,0} => #./..
	// map[point]map[point]string

	start := strings.Split(`.#.
..#
###`, "\n")
	startingGrid := internal.NewGridFromInput(start)
	fmt.Println(startingGrid)
	fmt.Println("------")

	iterations := 18
	for i := 0; i < iterations; i++ {
		fmt.Println(startingGrid)
		fmt.Println("----------")

		if startingGrid.Length%2 == 0 {
			fmt.Println("mod 2 == 0:", startingGrid.Length)
			m := splitN(startingGrid, 2)
			// for p, g := range m {
			// 	fmt.Printf("**** %v ****\n", p)
			// 	fmt.Println(g)
			// }
			// fmt.Println("*******")
			for j := 0; j < startingGrid.Height/2; j++ {
				for i := 0; i < startingGrid.Length/2; i++ {
					for _, ol := range allCombos(m[internal.Point{i, j}]) {
						if out, ok := rules[ol]; ok {
							m[internal.Point{i, j}] = gridFromOneLine(out)
							break
						}
					}
				}
			}
			startingGrid = join(m)
			continue
		}

		if startingGrid.Length%3 == 0 {
			fmt.Println("mod 3 == 0:", startingGrid.Length)
			m := splitN(startingGrid, 3)
			// for _, g := range m {
			// 	fmt.Println("********")
			// 	fmt.Println(g)
			// }
			// fmt.Println("*******")
			for j := 0; j < startingGrid.Height/3; j++ {
				for i := 0; i < startingGrid.Length/3; i++ {
					for _, ol := range allCombos(m[internal.Point{i, j}]) {
						if out, ok := rules[ol]; ok {
							m[internal.Point{i, j}] = gridFromOneLine(out)
							break
						}
					}
				}
			}
			startingGrid = join(m)
			continue
		}

	}
	fmt.Println(startingGrid)
	fmt.Println("----------")
	fmt.Println(startingGrid.On())

}

func substitute(rules map[string]string, input string) string {
	if out, ok := rules[input]; ok {
		return out
	}
	panic("not found input: " + input)
}

func join(in map[internal.Point]*internal.Grid) *internal.Grid {
	maxx, maxy := 0, 0
	size := 0
	for p, g := range in {
		size = g.Length
		if p.X > maxx {
			maxx = p.X
		}
		if p.Y > maxy {
			maxy = p.Y
		}
	}
	out := &internal.Grid{
		Length: size * (maxx + 1),
		Height: size * (maxy + 1),
		Data:   make(map[internal.Point]string),
	}
	for p, g := range in {
		for p2 := range g.Data {
			out.Data[internal.Point{p.X*size + p2.X, p.Y*size + p2.Y}] = g.At(p2.X, p2.Y)
		}
	}
	return out
}

type rule struct {
	from string
	to   string
}

func parseLine(line string) rule {
	w := strings.Split(line, " => ")
	return rule{w[0], w[1]}
}

func allCombos(g *internal.Grid) []string {
	out := []string{}
	for i := 0; i < 4; i++ {
		out = append(out, g.OneLine())
		g.Rotate()
	}
	g.Flip()
	for i := 0; i < 4; i++ {
		out = append(out, g.OneLine())
		g.Rotate()
	}
	return out
}

func gridFromOneLine(in string) *internal.Grid {
	return internal.NewGridFromInput(strings.Split(in, "/"))
}

func split(in *internal.Grid) map[internal.Point]*internal.Grid {
	out := map[internal.Point]*internal.Grid{}
	if in.Length != 4 {
		panic("bad size in splitting grid")
	}
	out[internal.Point{0, 0}] = in.SubGrid(0, 1, 0, 1)
	out[internal.Point{0, 1}] = in.SubGrid(0, 1, 2, 3)
	out[internal.Point{1, 0}] = in.SubGrid(2, 3, 0, 1)
	out[internal.Point{1, 1}] = in.SubGrid(2, 3, 2, 3)
	return out
}

func splitN(in *internal.Grid, n int) map[internal.Point]*internal.Grid {
	out := map[internal.Point]*internal.Grid{}
	if in.Length%n != 0 {
		panic("bad size in splitting grid")
	}
	size := in.Length / n
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			out[internal.Point{i, j}] = in.SubGrid(i*n, i*n+(n-1), j*n, j*n+(n-1))
		}
	}
	return out
}
