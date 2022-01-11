package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	w := newWell()
	for _, line := range lines {
		for _, p := range parseLine(line) {
			w.add(p, "#")
		}
	}
	//	w.add(point{500, 0}, "+")
	fmt.Println(w.miny, w.maxy)
	// for i := 0; i < 20000; i++ {
	// 	w.walk(point{500, 0})
	// }
	// for w.walk(point{500, 0}) {
	// }
	w.walk2(point{500, 0})
	//	w.walk2(point{500, 0})

	fmt.Println(w)
	fmt.Println(w.countWater())
}

// 22720 too low
// 31955 too high

type well struct {
	data       map[point]string
	minx, maxx int
	miny, maxy int
}

func (w *well) countWater() int {
	count := 0
	fmt.Println("counting", w.miny, w.maxy)
	for j := 3; j <= w.maxy; j++ {
		for i := w.minx; i <= w.maxx; i++ {
			item := w.At(point{i, j})
			if item == "~" {
				count++
			}
		}
	}
	return count
}

func newWell() *well {
	return &well{
		data: make(map[point]string),
		minx: math.MaxInt,
		miny: math.MaxInt,
		maxx: math.MinInt,
		maxy: math.MinInt,
	}
}

// walk2 until we're done
func (w *well) walk2(start point) {
	if w.At(start) == "+" {
		start = start.down()
	}
	// fall down until you hit a #. mark each one as "|"
	cur := start
	for w.At(cur) != "#" {
		w.add(cur, "|")
		cur = cur.down()
		if cur.y > w.maxy {
			return
		}
	}
	cur = cur.up()
	// am i on a floor?
	for w.floor(cur) {
		w.fillFloor(cur)
		// go up
		cur = cur.up()
		// walk left and right filling in with water
	}
	for _, e := range w.fillEdge(cur) {
		w.walk2(e)
	}
	for w.floor(cur) {
		w.fillFloor(cur)
		cur = cur.up()
	}
	for _, e := range w.fillEdge(cur) {
		w.walk2(e)
	}
}

func (w *well) fillEdge(cur point) []point {
	w.add(cur, "|")
	out := []point{}
	left := cur.left()
	right := cur.right()
	// until we can go down
	for w.At(left) == "." && w.At(left.down()) != "." {
		if w.At(left) == "." {
			w.add(left, "|")
			left = left.left()
		}
	}
	if w.At(left) == "." {
		out = append(out, left)
	}
	for w.At(right) == "." && w.At(right.down()) != "." {
		if w.At(right) == "." {
			w.add(right, "|")
			right = right.right()
		}
	}
	if w.At(right) == "." {
		out = append(out, right)
	}
	return out
}

func (w *well) fillFloor(cur point) {
	w.add(cur, "~")
	left := cur.left()
	right := cur.right()
	for w.At(left) == "." || w.At(left) == "|" || w.At(right) == "." || w.At(right) == "|" {
		if w.At(left) == "." || w.At(left) == "|" {
			w.add(left, "~")
			left = left.left()
		}
		if w.At(right) == "." || w.At(right) == "|" {
			w.add(right, "~")
			right = right.right()
		}
	}
}

func (w *well) At(p point) string {
	if item, ok := w.data[p]; ok {
		return item
	}
	return "."
}
func (w *well) floor(p point) bool {
	left := false
	for i := p.x - 1; i >= w.minx; i-- {
		if w.At(point{i, p.y}) == "#" {
			left = true
			break
		}
		if w.At(point{i, p.y + 1}) == "|" || w.At(point{i, p.y + 1}) == "." {
			return false
		}
	}
	right := false
	for i := p.x + 1; i <= w.maxx; i++ {
		if w.At(point{i, p.y}) == "#" {
			right = true
			break
		}
		if w.At(point{i, p.y + 1}) == "|" || w.At(point{i, p.y + 1}) == "." {
			return false
		}
	}
	return left && right
}

func (w *well) add(p point, s string) {
	if w.At(p) == "#" && s != "#" {
		panic("you little shit")
	}
	if p.x < w.minx {
		w.minx = p.x
	}
	if p.x > w.maxx {
		w.maxx = p.x
	}
	if p.y < w.miny {
		w.miny = p.y
	}
	if p.y > w.maxy {
		w.maxy = p.y
	}
	w.data[p] = s
}
func (w *well) String() string {
	var out strings.Builder
	for j := w.miny; j <= w.maxy; j++ {
		for i := w.minx; i <= w.maxx; i++ {
			out.WriteString(w.At(point{i, j}))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func parseLine(line string) []point {
	fmt.Print(line, ":")
	coords := strings.Split(line, ", ")
	var xstart, xend, ystart, yend int
	for _, c := range coords {
		items := strings.Split(c, "=")
		coordRange := strings.Split(items[1], "..")
		if len(coordRange) == 1 {
			if items[0] == "x" {
				xstart, _ = strconv.Atoi(coordRange[0])
				xend = xstart
			} else {
				ystart, _ = strconv.Atoi(coordRange[0])
				yend = ystart
			}
		}
		if len(coordRange) == 2 {
			start, _ := strconv.Atoi(coordRange[0])
			end, _ := strconv.Atoi(coordRange[1])
			if items[0] == "x" {
				xstart = start
				xend = end
			} else {
				ystart = start
				yend = end
			}
		}
	}
	return makePoints(xstart, xend, ystart, yend)
}

type point struct {
	x, y int
}

func (p point) up() point {
	return point{p.x, p.y - 1}
}
func (p point) down() point {
	return point{p.x, p.y + 1}
}
func (p point) left() point {
	return point{p.x - 1, p.y}
}
func (p point) right() point {
	return point{p.x + 1, p.y}
}
func (p point) dlr() []point {
	return []point{p.down(), p.left(), p.right()}
}

func makePoints(xstart, xend, ystart, yend int) []point {
	out := []point{}
	for y := ystart; y <= yend; y++ {
		for x := xstart; x <= xend; x++ {
			out = append(out, point{x, y})
		}
	}
	return out
}
