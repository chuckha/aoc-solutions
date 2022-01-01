package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type direction string

const (
	north = direction("north")
	south = direction("south")
	east  = direction("east")
	west  = direction("west")
)

func main() {
	lines := internal.ReadInput()
	iterations := 10000000
	nodes := internal.NewGridFromInput(lines)
	center := nodes.Center()
	c := &cleaner{
		direction: north,
		pos:       center,
		g:         nodes,
	}
	for i := 0; i < iterations; i++ {
		c.act()
	}
	fmt.Println(c.infections)
}

type cleaner struct {
	direction  direction
	pos        internal.Point
	g          *internal.Grid
	infections int
}

func (c *cleaner) act() {
	node := c.g.At(c.pos.X, c.pos.Y)
	switch node {
	case ".":
		c.turnLeft()
		c.g.Set(c.pos.X, c.pos.Y, "W")
	case "W":
		c.g.Set(c.pos.X, c.pos.Y, "#")
		c.infections++
	case "#":
		c.turnRight()
		c.g.Set(c.pos.X, c.pos.Y, "F")
	case "F":
		c.reverse()
		c.g.Set(c.pos.X, c.pos.Y, ".")
	default:
		panic("uh oh: " + node)
	}
	c.forward()
}
func (c *cleaner) turnRight() {
	switch c.direction {
	case north:
		c.direction = east
	case east:
		c.direction = south
	case south:
		c.direction = west
	case west:
		c.direction = north
	}
}
func (c *cleaner) turnLeft() {
	switch c.direction {
	case north:
		c.direction = west
	case east:
		c.direction = north
	case south:
		c.direction = east
	case west:
		c.direction = south
	}
}
func (c *cleaner) reverse() {
	switch c.direction {
	case north:
		c.direction = south
	case east:
		c.direction = west
	case south:
		c.direction = north
	case west:
		c.direction = east
	}
}

func (c *cleaner) forward() {
	switch c.direction {
	case north:
		c.pos.Y -= 1
	case east:
		c.pos.X += 1
	case south:
		c.pos.Y += 1
	case west:
		c.pos.X -= 1
	}
	if c.pos.X > c.g.Length {
		c.g.Length = c.pos.X
		c.g.Length++
	}
	if c.pos.Y > c.g.Height {
		c.g.Height = c.pos.Y
		c.g.Height++
	}
	if c.pos.X < c.g.MinX {
		c.g.MinX = c.pos.X
		c.g.Length++
	}
	if c.pos.Y < c.g.MinY {
		c.g.MinY = c.pos.Y
		c.g.Height++
	}
}

func (c *cleaner) String() string {
	var out strings.Builder
	for j := c.g.MinY; j < c.g.Height; j++ {
		for i := c.g.MinX; i < c.g.Height; i++ {
			if i == c.pos.X && j == c.pos.Y {
				out.WriteString("[")
				out.WriteString(c.g.At(i, j))
				out.WriteString("]")
				continue
			}
			out.WriteString(" ")
			out.WriteString(c.g.At(i, j))
			out.WriteString(" ")
		}
		out.WriteString("\n")
	}
	return out.String()
}
