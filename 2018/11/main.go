package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	serialIn := internal.ReadInput()[0]
	serial, _ := strconv.Atoi(serialIn)
	gridSize := 300

	g := newGrid(gridSize, serial)
	// // serial: 57
	// fmt.Println(g.variousSizePowers[1][point{122, 79}] == -5)
	// // serial: 39
	// fmt.Println(g.variousSizePowers[1][point{217, 196}] == 0)
	// // serial: 71
	// fmt.Println(g.variousSizePowers[1][point{101, 153}] == 4)
	// fmt.Println(g.powerSize3(90, 269, 16))

	// for i := 2; i < 17; i++ {
	// 	g.powerSize2(i)
	// }
	// for i := 1; i < 17; i++ {
	// 	fmt.Printf("size: %d, value: %d\n", i, g.variousSizePowers[i][point{90, 269}])
	// }
	// points := []point{}
	// for k := range g.variousSizePowers[16] {
	// 	points = append(points, k)
	// }
	// sort.Sort(pts(points))
	// for _, p := range points {
	// 	if g.variousSizePowers[16][p] < 112 {
	// 		continue
	// 	}
	// 	fmt.Println("here we go", p)
	// 	for i := 1; i < 17; i++ {
	// 		fmt.Println(i, p, g.variousSizePowers[i][p])
	// 	}
	// }
	// os.Exit(0)
	biggest := []power{}
	for i := 2; i <= 300; i++ {
		pwr := g.powerSize2(i)
		fmt.Println(pwr)
		//		print the values around this pwr
		// fmt.Println(point{232, 251}, i, g.variousSizePowers[i][point{232, 251}])
		// for j := 0; j < pwr.size; j++ {
		// 	for i := 0; i < pwr.size; i++ {
		// 		fmt.Printf("%02d ", g.data[point{232 + i, 251 + j}])
		// 	}
		// 	fmt.Println()
		// }
		biggest = append(biggest, pwr)
	}

	// tests for 18
	// fmt.Println("33,45,2", g.variousSizePowers[2][point{33, 45}] == 14)
	// fmt.Println("33,45,3", g.variousSizePowers[3][point{33, 45}] == 29)
	// fmt.Println("33,45,4", g.variousSizePowers[4][point{33, 45}] == 12)
	// fmt.Println("32,44,5", g.variousSizePowers[5][point{32, 44}] == 18)

	// tests for 42
	// fmt.Println("21,61,2", g.variousSizePowers[2][point{21, 61}] == 13)
	// fmt.Println("21,61,3", g.variousSizePowers[3][point{21, 61}] == 30)
	// fmt.Println("21,61,4", g.variousSizePowers[4][point{21, 61}] == 27)
	// fmt.Println("20,60,5", g.variousSizePowers[5][point{20, 60}] == 32)
	sort.Sort(powers(biggest))
	fmt.Println(biggest[0])
}

type grid struct {
	data map[point]int
	// size->map [point] starting point == power
	// then just add the edges to that number to make the next one
	variousSizePowers map[int]map[point]int
	minx, maxx        int
	miny, maxy        int
}

func (g *grid) powerSize2(size int) power {
	fmt.Println("--------------------------", size, "------------------------------")
	if _, ok := g.variousSizePowers[size]; ok {
		panic("wut")
	}
	g.variousSizePowers[size] = make(map[point]int)
	//	out := []power{}
	max := math.MinInt
	maxpower := power{}
	fmt.Println("from", g.miny, "to", g.maxy-size+1, "and x goes", g.minx, "to", g.maxx-size+1)
	for j := g.miny; j <= g.maxy-size+1; j++ {
		for i := g.minx; i <= g.maxx-size+1; i++ {
			cur := point{i, j}
			// get the square from the previous size at {i,j}
			edgePoints := cur.edges(size)
			// if size == 2 {
			// 	fmt.Println("top left", i, j)
			// 	for _, e := range edgePoints {
			// 		fmt.Println(e)
			// 	}
			// }
			subsquareVal, ok := g.variousSizePowers[size-1][cur]
			if !ok {
				panic(fmt.Sprintf("%v isn't in %d", cur, size-1))
			}
			if len(edgePoints) != (size-1)*2+1 {
				panic("too many or few edge points")
			}
			for _, p := range edgePoints {
				item, ok := g.data[p]
				if !ok {
					panic("not good")
				}
				subsquareVal += item
			}
			if _, ok := g.variousSizePowers[size][cur]; ok {
				panic("another baddie")
			}
			g.variousSizePowers[size][cur] = subsquareVal
			if subsquareVal > max {
				max = subsquareVal
				maxpower = power{x: cur.x, y: cur.y, size: size, val: subsquareVal}
				//				fmt.Println(maxpower)
			}
			//			out = append(out, power{x: i, y: j, size: size, val: subsquareVal})
		}
	}
	return maxpower
}

type powers []power

func (p powers) Len() int           { return len(p) }
func (p powers) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p powers) Less(i, j int) bool { return p[i].val > p[j].val }

type power struct {
	x, y int
	size int
	val  int
}

func newGrid(size, serial int) *grid {
	g := &grid{
		minx: 1, maxx: size,
		miny: 1, maxy: size,
		data:              make(map[point]int),
		variousSizePowers: make(map[int]map[point]int),
	}
	for j := g.miny; j <= g.maxy; j++ {
		for i := g.minx; i <= g.maxx; i++ {
			g.data[point{i, j}] = pwr(i, j, serial)
		}
	}
	g.variousSizePowers[1] = g.data
	return g
}

// 233,251, 42
// rackID: 243
// PL: 61035
// pl
func pwr(x, y, serial int) int {
	rackID := x + 10
	powerLevel := rackID * y
	powerLevel += serial
	powerLevel *= rackID
	return ((powerLevel / 100) % 10) - 5
	//	return hundreds(powerLevel) - 5
}

type point struct {
	x, y int
}

// (2) 1,1 -> 1,2, 2,1,   2,2
// (3) 1,1 -> 3,1, 3,2 1,3 2,3   3,3

// just get the edges for the new square
func (p point) edges(size int) []point {
	if size == 1 {
		return []point{}
	}
	out := []point{}
	for i := 0; i < size-1; i++ {
		out = append(out, point{p.x + i, p.y + size - 1})
	}
	for j := 0; j < size-1; j++ {
		out = append(out, point{p.x + size - 1, p.y + j})
	}
	out = append(out, point{p.x + size - 1, p.y + size - 1})
	return out
}

func (p point) square(size int) []point {
	out := []point{}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			out = append(out, point{p.x + i, p.y + j})
		}
	}
	return out
}

func hundreds(in int) int {
	for in > 10000000 {
		in = in - 10000000
	}
	for in > 1000000 {
		in = in - 1000000
	}
	for in > 100000 {
		in = in - 100000
	}
	for in > 10000 {
		in = in - 10000
	}
	for in >= 1000 {
		in = in - 1000
	}
	return in / 100
}
func hundos(in int) int {
	return (in / 100) % 10
}

// 12345 / 10

// 1,24,189 wrong answer
// 182,255,12 wrong answer
// 182,255,13 wrong
// 237 284 11 91

type pts []point

func (p pts) Len() int      { return len(p) }
func (p pts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p pts) Less(i, j int) bool {
	if p[i].x < p[j].x {
		return true
	}
	if p[i].x > p[j].x {
		return false
	}
	return p[i].y < p[j].y
}
