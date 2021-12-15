package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	p := newPad()
	for _, line := range lines {
		for _, dir := range strings.Split(line, "") {
			p.move(dir)
		}
		fmt.Printf(p.digit())
	}
	fmt.Println()
}

type point struct {
	x, y int
}
type pad struct {
	data map[point]string
	loc  point
}

func newPad() *pad {
	return &pad{
		data: map[point]string{
			{0, 2}: "5",
			{1, 1}: "2",
			{1, 2}: "6",
			{1, 3}: "A",
			{2, 0}: "1",
			{2, 1}: "3",
			{2, 2}: "7",
			{2, 3}: "B",
			{2, 4}: "D",
			{3, 1}: "4",
			{3, 2}: "8",
			{3, 3}: "C",
			{4, 2}: "9",
		},
		loc: point{1, 1},
	}
}

func (p *pad) move(dir string) {
	newX := p.loc.x
	newY := p.loc.y
	switch dir {
	case "U":
		newY -= 1
	case "R":
		newX += 1
	case "D":
		newY += 1
	case "L":
		newX -= 1
	}
	if _, ok := p.data[point{newX, newY}]; ok {
		p.loc.x = newX
		p.loc.y = newY
		return
	}
}

func (p *pad) digit() string {
	return p.data[p.loc]
}
