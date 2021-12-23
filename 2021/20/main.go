package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	iea := lines[0]
	image := newInputImage(iea, 0, 0, 0, 0)
	for y, line := range lines[1:] {
		for x, c := range line {
			if c == '#' {
				image.add(point{x, y})
			}
		}
	}
	//	fmt.Println(iea)
	for i := 0; i < 50; i++ {
		image = image.iterate()
		if i%2 == 1 {
			image.shrink(18)
		}
	}
	//	fmt.Println(image)
	fmt.Println(len(image.data))
}

type imageEnhancementAlgorithm struct {
	data string
}

func newIEA(in string) *imageEnhancementAlgorithm {
	return &imageEnhancementAlgorithm{in}
}

type inputImage struct {
	algo       string
	data       map[point]string
	minx, miny int
	maxx, maxy int
}

func newInputImage(algo string, minx, miny, maxx, maxy int) *inputImage {
	return &inputImage{
		algo: algo,
		data: map[point]string{},
		minx: minx, miny: miny, maxx: maxx, maxy: maxy,
	}
}

func (n *inputImage) iterate() *inputImage {
	out := newInputImage(n.algo, 0, 0, 0, 0)
	for j := n.miny - 10; j <= n.maxy+10; j++ {
		for i := n.minx - 10; i <= n.maxx+10; i++ {
			if n.outputPixel(point{i, j}) == "#" {
				out.add(point{i, j})
			}
		}
	}
	return out
}

func (n *inputImage) shrink(num int) {
	for j := n.miny; j <= n.maxy; j++ {
		for i := n.minx; i <= n.maxx; i++ {
			if j < n.miny+num || i < n.minx+num {
				delete(n.data, point{i, j})
			}
			if j > n.maxy-num || i > n.maxx-num {
				delete(n.data, point{i, j})
			}
		}
	}
	n.minx += num
	n.miny += num
	n.maxx -= num
	n.maxy -= num
}

func (i *inputImage) add(p point) {
	if p.x < i.minx {
		i.minx = p.x
	}
	if p.y < i.miny {
		i.miny = p.y
	}
	if p.x > i.maxx {
		i.maxx = p.x
	}
	if p.y > i.maxy {
		i.maxy = p.y
	}
	i.data[p] = "#"
}

func (n *inputImage) String() string {
	var out strings.Builder
	for j := n.miny - 1; j <= n.maxy+1; j++ {
		for i := n.minx - 1; i <= n.maxx+1; i++ {
			if _, ok := n.data[point{i, j}]; ok {
				out.WriteString("#")
				continue
			}
			out.WriteString(".")
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (i *inputImage) outputPixel(p point) string {
	return i.getPixelValueAtIndex(binToDec(i.outputPixelRaw(p)))
}

func (i *inputImage) outputPixelRaw(p point) string {
	var out strings.Builder
	for _, n := range p.neighbors() {
		if _, ok := i.data[n]; ok {
			out.WriteString("1")
			continue
		}
		out.WriteString("0")
	}
	return out.String()
}

func binToDec(in string) int {
	idx, _ := strconv.ParseInt(in, 2, 64)
	return int(idx)
}

func (i *inputImage) getPixelValueAtIndex(idx int) string {
	return string(i.algo[idx])
}

type point struct {
	x, y int
}

func (p point) neighbors() []point {
	return []point{
		{p.x - 1, p.y - 1}, {p.x, p.y - 1}, {p.x + 1, p.y - 1},
		{p.x - 1, p.y}, {p.x, p.y}, {p.x + 1, p.y},
		{p.x - 1, p.y + 1}, {p.x, p.y + 1}, {p.x + 1, p.y + 1},
	}
}
