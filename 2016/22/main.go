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

	nodes := map[point]data{}
	maxx := 0
	maxy := 0
	for _, line := range lines[2:] {
		pt, d := parse(line)
		if pt.x > maxx {
			maxx = pt.x
		}
		if pt.y > maxy {
			maxy = pt.y
		}
		nodes[pt] = d
	}

	pairs := [][2]data{}

	for p, node := range nodes {
		for p2, node2 := range nodes {
			if p == p2 {
				continue
			}
			if node.empty() {
				continue
			}
			if node.used < node2.avail {
				pairs = append(pairs, [2]data{node, node2})
			}
		}
	}

	fmt.Println("part one:", len(pairs))
	emptyNode := point{}
	for y := 0; y <= maxy; y++ {
		for x := 0; x <= maxx; x++ {
			// if nodes[point{x, y}].size > 500 {
			// 	fmt.Print("#")
			// 	continue
			// }
			if nodes[point{x, y}].used == 0 {
				emptyNode = point{x, y}
				//fmt.Print("_")
				continue
			}
			// if x == maxx && y == 0 {
			// 	fmt.Print("G")
			// 	continue
			// }
			// fmt.Print(".")
		}
		//		fmt.Println()
	}
	sum := 0
	// no wall
	costs := shortestPath(maxx, maxy, nodes, emptyNode, point{-1, 1})
	sum += costs[point{maxx, 0}]
	oldGoal := point{maxx, 0}
	goalStart := point{maxx - 1, 0}
	for goalStart != (point{0, 0}) {
		nodes[goalStart] = nodes[oldGoal]
		nodes[oldGoal] = nodes[emptyNode]

		n1costs := shortestPath(maxx, maxy, nodes, oldGoal, goalStart)
		printCosts(maxx, maxy, n1costs)
		fmt.Println()
		oldGoal = point{oldGoal.x - 1, 0}
		goalStart = point{goalStart.x - 1, 0}
		sum += n1costs[goalStart] + 1 // +1 because we move it over the old goal
	}
	fmt.Println(sum)
}

func printCosts(maxx, maxy int, costs map[point]int) {
	for y := 0; y <= maxy; y++ {
		for x := 0; x <= maxx; x++ {
			if costs[point{x, y}] == math.MaxInt {
				fmt.Print(" ## ")
				continue
			}
			fmt.Printf(" %02d ", costs[point{x, y}])
		}
		fmt.Println()
	}
}

// /dev/grid/node-x0-y0     89T   65T    24T   73%
func parse(line string) (point, data) {
	words := strings.Fields(line)
	fs := strings.Split(words[0], "/")
	coord := strings.Split(fs[3], "-")
	xval := strings.Split(coord[1], "x")
	yval := strings.Split(coord[2], "y")
	xx, _ := strconv.Atoi(xval[1])
	yy, _ := strconv.Atoi(yval[1])
	pt := point{
		x: xx, y: yy,
	}
	size, _ := strconv.Atoi(strings.TrimSuffix(words[1], "T"))
	used, _ := strconv.Atoi(strings.TrimSuffix(words[2], "T"))
	avail, _ := strconv.Atoi(strings.TrimSuffix(words[3], "T"))
	percent, _ := strconv.Atoi(strings.TrimSuffix(words[4], "%"))

	return pt, data{
		size:    size,
		used:    used,
		avail:   avail,
		percent: percent,
	}
}

type point struct {
	x, y int
}
type data struct {
	size, used, avail, percent int
}

func (d data) empty() bool {
	return d.used == 0
}
func (d data) String() string {
	return fmt.Sprintf("Size: %d Used: %d, Avail: %d, Percent: %d%%", d.size, d.used, d.avail, d.percent)
}

func shortestPath(maxx, maxy int, nodes map[point]data, from, wall point) map[point]int {
	visited := map[point]struct{}{}
	costs := map[point]int{}
	for y := 0; y <= maxy; y++ {
		for x := 0; x <= maxx; x++ {
			costs[point{x, y}] = math.MaxInt
		}
	}
	costs[from] = 0
	q := internal.NewQueue[point]()
	q.Enqueue(from)
	for !q.Empty() {
		cur := q.Dequeue()
		if _, ok := visited[cur]; ok {
			continue
		}
		for _, n := range cur.neighbors() {
			if _, ok := nodes[n]; !ok {
				continue
			}
			if nodes[from].avail < nodes[n].used || n == wall {
				costs[n] = math.MaxInt
				continue
			}
			if costs[cur]+1 < costs[n] {
				costs[n] = costs[cur] + 1
			}
			q.Enqueue(n)
		}
		visited[cur] = struct{}{}
	}
	return costs
}

func (p point) neighbors() []point {
	return []point{
		{x: p.x, y: p.y - 1},
		{x: p.x - 1, y: p.y},
		{x: p.x, y: p.y + 1},
		{x: p.x + 1, y: p.y},
	}
}
