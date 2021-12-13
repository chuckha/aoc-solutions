package main

import (
	"fmt"
)

func main() {
	//	lines := internal.ReadInput()

	infinitePaper := newPaper()
	codePoint := point{3019, 3010}
	for {
		infinitePaper.inc()
		if v, ok := infinitePaper.data[codePoint]; ok {
			fmt.Println(v)
			break
		}
	}
}

type point struct {
	x, y int
}

type paper struct {
	data         map[point]int
	maxX, maxY   int
	lastX, lastY int
}

func newPaper() *paper {
	return &paper{
		data: map[point]int{
			{1, 1}: 20151125,
		},
		maxX:  1,
		maxY:  1,
		lastX: 1,
		lastY: 1,
	}
}

func (p *paper) inc() {
	nextX := p.lastX + 1
	nextY := p.lastY - 1
	if nextY == 0 {
		p.maxY += 1
		nextY = p.maxY
		nextX = 1
	}
	if nextX > p.maxX {
		p.maxX = nextX
	}
	p.data[point{nextX, nextY}] = (p.data[point{p.lastX, p.lastY}] * 252533) % 33554393
	p.lastX = nextX
	p.lastY = nextY
}
func (p *paper) String() string {
	out := ""
	for j := 1; j <= p.maxY; j++ {
		for i := 1; i <= p.maxX; i++ {
			out += fmt.Sprintf(" %08d ", p.data[point{i, j}])
		}
		out += "\n"
	}
	return out
}
