package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const length = 10

func main() {
	input := internal.ReadInput()
	grid := internal.NewGridV3[*Cell]()
	// starts at 0,0
	grid.Set(internal.Point{0, 0}, &Cell{start: true, visited: true})
	pts := []internal.Point{}
	for i := 0; i < length; i++ {
		pts = append(pts, internal.Point{0, 0})
	}
	for _, line := range input {
		instructions := parse(line)
		//		fmt.Println("running instructions", instructions)
		for _, inst := range instructions {
			pts[0] = Move(pts[0], inst)
			for i := 1; i < len(pts); i++ {
				pts[i] = moveTowards(pts[i-1], pts[i])
			}
			if !grid.In(pts[0]) {
				grid.Set(pts[0], &Cell{})
			}
			if !grid.In(pts[len(pts)-1]) {
				grid.Set(pts[len(pts)-1], &Cell{visited: true})
			}
			if grid.In(pts[len(pts)-1]) {
				grid.At(pts[len(pts)-1]).visited = true
			}
		}
		// printme := grid.Copy()
		// printme.Set(pts[0], &Cell{data: "H"})
		// for i := 1; i < len(pts); i++ {
		// 	if printme.In(pts[i]) && printme.At(pts[i]).data != "" {
		// 		continue
		// 	}
		// 	printme.Set(pts[i], &Cell{data: fmt.Sprintf("%d", i)})
		// }
		// fmt.Println(printme)
	}
	count := 0
	for _, c := range grid.Data {
		if c.visited {
			count++
		}
	}
	fmt.Println(count)
}

// one off in column, move tail towards head
func moveTowards(head, tail internal.Point) internal.Point {
	if head.ManhattanDistance(tail) == 1 {
		return tail
	}
	if internal.Abs(head.X-tail.X) == 1 && internal.Abs(head.Y-tail.Y) == 1 {
		return tail
	}
	if head.X > tail.X {
		tail.X += 1
	}
	if head.Y > tail.Y {
		tail.Y += 1
	}
	if head.X < tail.X {
		tail.X -= 1
	}
	if head.Y < tail.Y {
		tail.Y -= 1
	}
	return tail
}

type Cell struct {
	start   bool
	visited bool
	data    string
}

func (c *Cell) String() string {
	switch {
	case c.data != "":
		return c.data
	case c.start:
		return "s"
	case c.visited:
		return "#"
	default:
		return "."
	}
}

func parse(in string) []string {
	parts := strings.Split(in, " ")
	num, _ := strconv.Atoi(parts[1])
	out := make([]string, num)
	for i := range out {
		out[i] = parts[0]
	}
	return out
}

func Move(p internal.Point, dir string) internal.Point {
	switch dir {
	case "L":
		return p.Left()
	case "U":
		return p.Up()
	case "R":
		return p.Right()
	case "D":
		return p.Down()
	default:
		panic("unsupportedd move direction")
	}
}
