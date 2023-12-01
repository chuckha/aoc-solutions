package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

/*

This is messy. The idea was to find a pattern. In my specific input the pattern was found by looking for 4 0s in a row.
This repeated every 1720 rocks placed and started after the first 849 rocks got placed.

Then it was a matter of not being dumb which took a very long time.

*/

/*
	guesses
	1586627906014 too low
	1586627906147 too low
	1586627906921
*/

func main() {
	input := internal.ReadInput()
	js := &jetstream{data: strings.Split(input[0], ""), ptr: 0}
	c := newCave(7, js)
	// heightlookup is the cycle height lookup, will only go up to 1720
	heightLookup := map[int]int{}
	// m := map[int]state{}
	start := 0
	// zeroesInARow := 0

	//	fmt.Println((100000 - 1322) / 2729)
	// os.Exit(0)
	//	nearest := 581395348

	// iteration: 84 rock placed: 144480 height: 229179

	// // height = 2729 * iter + 1322 = 100000 -1322 =
	// cycleLength := 1720
	// cycleStart := 849 // (at height -1322)// weird
	prev := 0
	//	previ := 0
	heightOffset := 0
	for {
		c.tick()

		// if c.stoppedRockCount != start && between(c.stoppedRockCount, cycleStart, cycleStart+cycleLength) {
		// 	cur := -c.rocks.Min.Y
		// 	fmt.Println(c.stoppedRockCount, cur)
		// 	start = c.stoppedRockCount
		// }

		if c.stoppedRockCount != start {
			if between(c.stoppedRockCount, 2569, 4289) {
				cur := -c.rocks.Min.Y
				heightLookup[heightOffset] = cur - 1322 - 2729
				heightOffset++
			}
			cur := -c.rocks.Min.Y
			fmt.Println(c.stoppedRockCount, cur, prev-cur)
			start = c.stoppedRockCount
			prev = cur
		}
		/*
			(numberOfRocks-849) / 1720 + heightLookup[(numberOfRocks-849)%1720]
			5 9449 14967 -2729
			6 11169 17696 -2729

			rocks: 4979  height: 7847 0
		*/

		// }
		if c.stoppedRockCount == 5000 {
			fmt.Println(heightLookup)
			// how to calculate number of rocks on a cycle
			//			fmt.Println(1720*36 + 849)
			/// how to calculate height given  a cycle number
			//			fmt.Println(36*2729 + 1322)
			// x := 1000000000000
			y := 1000000000000
			cn := cycleNumber(y)
			fmt.Println("cycle number", cn)
			h := height(cn, y, heightLookup)
			fmt.Println(h)
			break
		}

		// This actually shows a really strong pattern (!!)

		// if c.stoppedRockCount != start {
		// 	cur := -c.rocks.Min.Y
		// 	if cur-prev == 0 {
		// 		zeroesInARow++
		// 	}
		// 	if cur-prev != 0 {
		// 		zeroesInARow = 0
		// 	}
		// 	if zeroesInARow == 4 {
		// 		fmt.Printf("got 4 zeroes in a row from %d to %d (%d)\n", c.stoppedRockCount-4, c.stoppedRockCount, c.stoppedRockCount-previ)
		// 		previ = c.stoppedRockCount
		// 	}
		// 	start = c.stoppedRockCount
		// 	prev = cur
		// }
		// if c.stoppedRockCount == 10091*4*5+10 {
		// 	break
		// }

		// if c.stoppedRockCount != start && c.stoppedRockCount%len(input[0]) == 0 {
		// 	// just rockified a new rock
		// 	use := c.shapePtr - 1
		// 	if c.shapePtr == 0 {
		// 		use = len(shapes) - 1
		// 	}
		// 	if c.shapePtr == 1 {
		// 		use = len(shapes)
		// 	}
		// 	m[c.stoppedRockCount] = state{
		// 		latestShape: shapes[use-1],
		// 		height:      -c.rocks.Min.Y,
		// 		rockCount:   c.stoppedRockCount,
		// 	}
		// 	//			c.trim(c.rocks.Min.Y + 100000)

		// 	start = c.stoppedRockCount
		// }
		// if len(m) == 130 {
		// 	doit := states{}
		// 	for _, v := range m {
		// 		doit = append(doit, v)
		// 	}
		// 	sort.Sort(doit)
		// 	for i := 0; i < len(doit)-1; i++ {
		// 		fmt.Println(i, doit[i].latestShape, doit[i], doit[i+1].height-doit[i].height)
		// 	}
		// 	return
		// }
	}
}

func cycleNumber(numberOfRocks int) int {
	return ((numberOfRocks - 849) / 1720)
}

func rocks(cycleNumber int) int {
	return cycleNumber*1720 + 849
}

func height(cycleNumber int, numberOfRocks int, hl map[int]int) int {
	fmt.Println("lookup number", numberOfRocks-rocks(cycleNumber), hl[numberOfRocks-rocks(cycleNumber)])
	return cycleNumber*2729 + 1322 + hl[numberOfRocks-rocks(cycleNumber)]
}

type state struct {
	latestShape shape
	rockCount   int
	height      int
}
type states []state

func (a states) Len() int           { return len(a) }
func (a states) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a states) Less(i, j int) bool { return a[i].rockCount < a[j].rockCount }

// 0
// 1 -> 31970 - 2 * 15959 = 52
// 2 ->

type cave struct {
	shapePtr         int
	jetstream        *jetstream
	rocks            *internal.GridV3[cell]
	falling          *internal.GridV3[cell]
	stoppedRockCount int
	width            int
}

func newCave(width int, js *jetstream) *cave {
	rocks := internal.NewGridV3[cell]()
	for i := 0; i < width; i++ {
		rocks.Set(internal.Point{i, 0}, cell("-"))
	}
	falling := internal.NewGridV3[cell]()
	return &cave{
		jetstream: js,
		rocks:     rocks,
		falling:   falling,
		width:     width,
	}
}

func (c *cave) expectedHeight() int {
	fullSets := c.stoppedRockCount / 5
	counts := map[shape]int{
		Minus:  fullSets,
		Plus:   fullSets,
		L:      fullSets,
		Bar:    fullSets,
		Square: fullSets,
	}
	for i := c.stoppedRockCount % 5; i > 0; i++ {
		counts[Minus]++
		i--
		if i == 0 {
			break
		}
		counts[Plus]++
		i--
		if i == 0 {
			break
		}
		counts[L]++
		i--
		if i == 0 {
			break
		}
		counts[Bar]++
		i--
		if i == 0 {
			break
		}
	}
	return counts[Minus] + counts[Plus]*3 + counts[L]*3 + counts[Bar]*4 + counts[Square]*2
}

func (c *cave) trim(x int) {
	for pt := range c.rocks.Data {
		if pt.Y > x {
			delete(c.rocks.Data, pt)
		}
	}
}

func (c *cave) tick() {
	if len(c.falling.Data) == 0 {
		c.insertShape()
	}

	nf := internal.NewGridV3[cell]()
	switch c.jetstream.next() {
	case "<":
		if c.falling.Min.X == 0 {
			for pt := range c.falling.Data {
				nf.Set(pt, cell("@"))
			}
			break
		}
		skip := false
		for pt := range c.falling.Data {
			if c.rocks.In(pt.Left()) {
				skip = true
			}
		}
		if skip {
			for pt := range c.falling.Data {
				nf.Set(pt, cell("@"))
			}
			break
		}
		for pt := range c.falling.Data {
			nf.Set(pt.Left(), cell("@"))
		}
	case ">":
		if c.falling.Max.X == c.width-1 {
			for pt := range c.falling.Data {
				nf.Set(pt, cell("@"))
			}
			break
		}
		skip := false
		for pt := range c.falling.Data {
			if c.rocks.In(pt.Right()) {
				skip = true
			}
		}
		if skip {
			for pt := range c.falling.Data {
				nf.Set(pt, cell("@"))
			}
			break
		}
		for pt := range c.falling.Data {
			nf.Set(pt.Right(), cell("@"))
		}
	default:
		panic("bad jetstream next")
	}
	c.falling = nf

	// solidify if possible
	for pt := range nf.Data {
		if c.rocks.In(pt.Down()) {
			c.solidify()
			c.insertShape()
			return
		}
	}

	nf2 := internal.NewGridV3[cell]()
	// no collisions, move everything down
	for pt := range nf.Data {
		nf2.Set(pt.Down(), cell("@"))
	}
	c.falling = nf2
}

func (c *cave) insertShape() {
	c.insert(shapes[c.shapePtr])
	c.shapePtr++
	if c.shapePtr >= len(shapes) {
		c.shapePtr = 0
	}
}

func (c *cave) solidify() {
	for pt := range c.falling.Data {
		c.rocks.Set(pt, cell("#"))
	}
	c.stoppedRockCount++
	c.falling = internal.NewGridV3[cell]()
}

func (c *cave) String() string {
	return c.falling.Layer(c.rocks).String()
}

func (c *cave) height() int {
	return c.rocks.Min.Y
}

// 	min := math.MaxInt
// 	for p := range c.rocks.Data {
// 		if p.Y < min {
// 			min = p.Y
// 		}
// 	}
// 	return min
// }

func (c *cave) insert(shape shape) {
	height := c.height() - 4
	switch shape {
	case Minus:
		c.falling.Set(internal.Point{2, height}, cell("@"))
		c.falling.Set(internal.Point{3, height}, cell("@"))
		c.falling.Set(internal.Point{4, height}, cell("@"))
		c.falling.Set(internal.Point{5, height}, cell("@"))
	case Plus:
		c.falling.Set(internal.Point{3, height - 2}, cell("@"))
		c.falling.Set(internal.Point{2, height - 1}, cell("@"))
		c.falling.Set(internal.Point{3, height - 1}, cell("@"))
		c.falling.Set(internal.Point{4, height - 1}, cell("@"))
		c.falling.Set(internal.Point{3, height}, cell("@"))
	case L:
		c.falling.Set(internal.Point{4, height - 2}, cell("@"))
		c.falling.Set(internal.Point{4, height - 1}, cell("@"))
		c.falling.Set(internal.Point{4, height}, cell("@"))
		c.falling.Set(internal.Point{3, height}, cell("@"))
		c.falling.Set(internal.Point{2, height}, cell("@"))
	case Bar:
		c.falling.Set(internal.Point{2, height - 3}, cell("@"))
		c.falling.Set(internal.Point{2, height - 2}, cell("@"))
		c.falling.Set(internal.Point{2, height - 1}, cell("@"))
		c.falling.Set(internal.Point{2, height}, cell("@"))
	case Square:
		c.falling.Set(internal.Point{2, height - 1}, cell("@"))
		c.falling.Set(internal.Point{3, height - 1}, cell("@"))
		c.falling.Set(internal.Point{2, height}, cell("@"))
		c.falling.Set(internal.Point{3, height}, cell("@"))
	}
}

type cell string

func (c cell) String() string {
	return string(c)
}

type shape string

var (
	Minus  shape = "-"
	Plus   shape = "+"
	L      shape = "⅃"
	Bar    shape = "|"
	Square shape = "□"

	shapes = []shape{Minus, Plus, L, Bar, Square}
)

type jetstream struct {
	data []string
	ptr  int
}

func (j *jetstream) next() string {
	next := j.data[j.ptr]
	j.ptr++
	if j.ptr >= len(j.data) {
		j.ptr = 0
	}
	return next
}

/*
minus, plus, ell, tall, square

minus piece will contribute from 0 to 1 to the height
plus piece will contribute from 0 to 3 to the height
ell piece will contribute from 0 to 3 to the height
tall piece will contribute from 0 to 4 to the height
square piece will contribute from 0 to 2 to the height

1 + 3 + 3 + 4 + 2 = 13 * 2022


*/

func between(val, low, high int) bool {
	return low <= val && high >= val
}
