package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	f := &floor{
		data: make(map[point]string),
		maxx: len(lines[0]) - 1,
		maxy: len(lines) - 1,
	}
	for y, line := range lines {
		for x, c := range line {
			f.data[point{x, y}] = string(c)
		}
	}
	count := 1
	for {
		if f.move() {
			break
		}
		count++
	}
	fmt.Println(count)
}

type point struct {
	x, y int
}
type floor struct {
	data       map[point]string
	maxx, maxy int
}

func (f *floor) String() string {
	var out strings.Builder
	for j := 0; j <= f.maxy; j++ {
		for i := 0; i <= f.maxx; i++ {
			out.WriteString(f.data[point{i, j}])
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (f *floor) move() bool {
	newfloor := make(map[point]string)
	for j := 0; j <= f.maxy; j++ {
		for i := 0; i <= f.maxx; i++ {
			if f.data[point{i, j}] == ">" {
				if i == f.maxx {
					if f.data[point{0, j}] == "." {
						newfloor[point{0, j}] = ">"
						newfloor[point{i, j}] = "."
						continue
					}
				}
				if f.data[point{i + 1, j}] == "." {
					newfloor[point{i + 1, j}] = ">"
					newfloor[point{i, j}] = "."
					continue
				}
			}
			// don't overwrite any data we've already written
			if newfloor[point{i, j}] == "" {
				newfloor[point{i, j}] = f.data[point{i, j}]
			}
		}
	}

	for j := 0; j <= f.maxy; j++ {
		for i := 0; i <= f.maxx; i++ {
			if f.data[point{i, j}] == "v" {
				if j == f.maxy {
					if f.data[point{i, 0}] != "v" && newfloor[point{i, 0}] == "." {
						newfloor[point{i, 0}] = "v"
						newfloor[point{i, j}] = "."
						continue
					}
				}
				if newfloor[point{i, j + 1}] == "." {
					newfloor[point{i, j + 1}] = "v"
					newfloor[point{i, j}] = "."
					continue
				}
			}
			// don't overwrite any data we've already written
			if newfloor[point{i, j}] == "" {
				newfloor[point{i, j}] = f.data[point{i, j}]
			}
		}
	}
	out := same(f.data, newfloor)
	f.data = newfloor
	return out
}

func same(m1, m2 map[point]string) bool {
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}
