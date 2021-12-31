package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	//input := 10
	input := 1350

	maxx := 50
	maxy := 50
	start := point{1, 1}
	//	end := point{7,4}
	end := point{31, 39}

	mz := &maze{points: make(map[point]string), maxx: maxx, maxy: maxy}
	for j := 0; j < maxx; j++ {
		for i := 0; i < maxy; i++ {
			mz.points[point{i, j}] = calcWall(i, j, input)
		}
	}
	fmt.Println(mz)
	costs := mz.solve(start, end)
	// for j := 0; j < 10; j++ {
	// 	for i := 0; i < 10; i++ {
	// 		if mz.points[point{i, j}] == "#" {
	// 			fmt.Printf("  #  ")
	// 			continue
	// 		}
	// 		if costs[point{i, j}] == math.MaxInt {
	// 			fmt.Printf("  $  ")
	// 			continue
	// 		}
	// 		fmt.Printf(" %03d ", costs[point{i, j}])
	// 	}
	// 	fmt.Println()
	// }
	fmt.Println(costs[end])
	count := 0
	for _, c := range costs {
		if c <= 50 {
			count++
		}
	}
	fmt.Println(count)
}

type maze struct {
	points     map[point]string
	maxx, maxy int
}

func (m *maze) String() string {
	var out strings.Builder
	for y := 0; y <= m.maxy; y++ {
		for x := 0; x <= m.maxx; x++ {
			out.WriteString(m.points[point{x, y}])
		}
		out.WriteString("\n")
	}
	return out.String()
}

type point struct {
	x, y int
}

func calcWall(x, y, fave int) string {
	v := x*x + 3*x + 2*x*y + y + y*y + fave
	bin := strconv.FormatInt(int64(v), 2)
	ones := strings.Count(bin, "1")
	if ones%2 == 0 {
		return "."
	}
	return "#"
}

func (m *maze) solve(start, end point) map[point]int {
	visited := map[point]struct{}{}
	cost := map[point]int{}
	for p := range m.points {
		cost[p] = math.MaxInt
	}
	cost[start] = 0
	q := internal.NewQueue[point]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		curcost := cost[cur]
		visited[cur] = struct{}{}
		for _, ne := range neighbors(cur) {
			// not in the maze
			if _, ok := m.points[ne]; !ok {
				continue
			}
			// it's a wall so we can't go there
			if m.points[ne] == "#" {
				continue
			}
			// we'ce already visited this spot
			if _, ok := visited[ne]; ok {
				continue
			}
			// otherwise we go to it
			if curcost+1 < cost[ne] {
				cost[ne] = curcost + 1
			}
			q.Enqueue(ne)
		}
	}
	return cost
}

func neighbors(p point) []point {
	return []point{
		{p.x, p.y - 1},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
		{p.x - 1, p.y},
	}
}
