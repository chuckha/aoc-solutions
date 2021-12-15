package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	g := internal.NewGrid(50, 6, ".")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		switch {
		case strings.HasPrefix(line, "rect"):
			dims := strings.Split(parts[1], "x")
			x, _ := strconv.Atoi(dims[0])
			y, _ := strconv.Atoi(dims[1])
			g.Rect(x, y)
		case strings.HasPrefix(line, "rotate row"):
			dims := strings.Split(parts[2], "=")
			y, _ := strconv.Atoi(dims[1])
			by, _ := strconv.Atoi(parts[4])
			g.RotateRow(y, by)
		case strings.HasPrefix(line, "rotate column"):
			dims := strings.Split(parts[2], "=")
			x, _ := strconv.Atoi(dims[1])
			by, _ := strconv.Atoi(parts[4])
			g.RotateCol(x, by)
		default:
			panic("bad instruction")
		}
	}
	fmt.Println(g.On())
	fmt.Println(g)
}
