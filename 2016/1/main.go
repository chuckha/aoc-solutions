package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	unparsed := strings.Split(lines[0], ", ")

	instructions := []instruction{}
	for _, inst := range unparsed {
		num, _ := strconv.Atoi(inst[1:])
		in := instruction{
			dir: string(inst[0]),
			num: num,
		}
		instructions = append(instructions, in)
	}
	b := &board{
		data: map[point]struct{}{},
		loc:  point{0, 0},
		dir:  "N",
	}
	for _, inst := range instructions {
		b.move(inst)
	}
}

type instruction struct {
	dir string
	num int
}

type point struct {
	x, y int
}

func (p point) dist() int {
	if p.x < 0 {
		p.x = p.x * -1
	}
	if p.y < 0 {
		p.y = p.y * -1
	}
	return p.x + p.y
}

type board struct {
	data map[point]struct{}
	loc  point
	dir  string
}

func (b *board) move(inst instruction) {
	switch b.dir {
	case "N":
		switch inst.dir {
		case "R":
			b.dir = "E"
		case "L":
			b.dir = "W"
		}
	case "E":
		switch inst.dir {
		case "R":
			b.dir = "S"
		case "L":
			b.dir = "N"
		}
	case "S":
		switch inst.dir {
		case "R":
			b.dir = "W"
		case "L":
			b.dir = "E"
		}
	case "W":
		switch inst.dir {
		case "R":
			b.dir = "N"
		case "L":
			b.dir = "S"
		}
	}
	for i := 0; i < inst.num; i++ {
		switch b.dir {
		case "N":
			b.loc.y += 1
		case "E":
			b.loc.x += 1
		case "S":
			b.loc.y -= 1
		case "W":
			b.loc.x -= 1
		}
		if _, ok := b.data[b.loc]; ok {
			fmt.Println(b.loc)
			fmt.Println(b.loc.dist())
			os.Exit(0)
		}
		b.data[b.loc] = struct{}{}
	}
}
