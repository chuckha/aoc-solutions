package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()[0]
	width := 25
	height := 6
	layers := buildLayers(input, width, height)
	fewestZeroes := layerWithFewestZeroes(layers)
	fmt.Println(part1(fewestZeroes))
	fmt.Println(part2(layers))
}

type layer struct {
	grid   *internal.Grid
	counts map[string]int
	depth  int
}

// l is bottom layer
func (l *layer) merge(l2 *layer) {
	for pt, pixel := range l2.grid.Data {
		if pixel != "2" {
			l.grid.Data[pt] = pixel
		}
	}
}

func (l *layer) String() string {
	var out strings.Builder
	for j := l.grid.MinY; j <= l.grid.Height; j++ {
		for i := l.grid.MinX; i <= l.grid.Length; i++ {
			switch l.grid.Data[internal.Point{i, j}] {
			case "0":
				out.Write(internal.Black("0"))
			case "1":
				out.Write(internal.Green("1"))
			case "2":
				panic("nah b")
			}
		}
		out.WriteString("\n")
	}
	return out.String()
}

func buildLayers(input string, w, h int) []*layer {
	layers := []*layer{}

	layerDepth := 0
	for k := 0; k < len(input); k += (w * h) {
		layer := &layer{grid: internal.NewGrid(w, h, " "), depth: layerDepth, counts: make(map[string]int)}
		for j := 0; j < h; j++ {
			for i := 0; i < w; i++ {
				item := input[k+(i+j*w)]
				layer.grid.Set(i, j, string(item))
				layer.counts[string(item)]++
			}
		}
		layerDepth++
		layers = append(layers, layer)
	}
	return layers
}

func layerWithFewestZeroes(layers []*layer) *layer {
	min := layers[0]
	for _, l := range layers {
		if l.counts["0"] < min.counts["0"] {
			min = l
		}
	}
	return min
}

// number of one digits multipled by the number of 2 digits
func part1(l *layer) int {
	return l.counts["1"] * l.counts["2"]
}

func part2(layers []*layer) *layer {
	stack := internal.NewStack[*layer]()
	for _, l := range layers {
		stack.Push(l)
	}
	bottom := stack.Pop()
	for !stack.Empty() {
		bottom.merge(stack.Pop())
	}
	return bottom
}
