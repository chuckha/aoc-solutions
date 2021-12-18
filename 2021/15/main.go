package main

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	tmx := len(lines[0])
	tmy := len(lines)

	game := newGame(tmx*5, tmy*5)
	for j, line := range lines {
		for i, c := range line {
			risk, _ := strconv.Atoi(string(c))
			game.data[point{i + 1, j + 1}] = risk
		}
	}

	// duplicate right
	for w := 1; w < 5; w++ {
		for j := 1; j <= tmy; j++ {
			for i := 1; i <= tmx; i++ {
				newx := tmx*w + i
				old := game.data[point{newx - tmx, j}]
				game.data[point{newx, j}] = (old % 9) + 1
			}
		}
	}
	// duplicate down
	for h := 1; h < 5; h++ {
		for i := 1; i <= tmx; i++ {
			for j := 1; j <= tmy; j++ {
				newy := tmy*h + j
				old := game.data[point{i, newy - tmy}]
				game.data[point{i, newy}] = (old % 9) + 1
			}
		}
	}
	// second row duplicate right
	for h := 1; h < 5; h++ {
		for w := 1; w < 5; w++ {
			for j := 1; j <= tmy; j++ {
				for i := 1; i <= tmx; i++ {
					newx := tmx*w + i
					old := game.data[point{newx - tmx, tmy*h + j}]
					game.data[point{newx, tmy*h + j}] = (old % 9) + 1
				}
			}
		}
	}
	//	fmt.Println(game)
	solve3(1, 1, game)
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

// y = 0; x = 1 =>
//
// 0 1   len = 2; height = 3 x = 0 y = 0 => y * len + x
// 2 3   // x =
// 4 5
// {1,1}, {2, 1}
// {1,2}, {2, 2}
// {1,3}, {2, 3}
// {1,1}, {2, 1}, {1,2}, {2, 2}, {1,3}, {2, 3}

/// {1,1}, {1,2}, {1,3}, {1,4}
//  {1,2}, {3,4}
//  x-1+y-1 (0),                  4 x-1 + y-1 ( *g.maxX)
// [{1,1}, {1,2}, {1, 3}, {1, 4}, {1, 2}]

func solve3(startx, starty int, g *game) {
	costs := map[point]int{}
	items := map[point]*Item{}
	pq := make(PriorityQueue, (g.maxX)*(g.maxY))
	for y := 1; y <= g.maxY; y++ {
		for x := 1; x <= g.maxX; x++ {
			costs[point{x, y}] = math.MaxInt
			tmp := (y-1)*g.maxY + (x - 1)
			item := &Item{
				point:    point{x, y},
				priority: math.MaxInt,
				index:    tmp,
			}
			pq[tmp] = item
			items[point{x, y}] = item
		}
	}
	pq[0] = &Item{
		point:    point{1, 1},
		priority: 0,
		index:    0,
	}
	costs[point{1, 1}] = 0
	heap.Init(&pq)
	for pq.Len() > 0 {
		lowesti := heap.Pop(&pq)
		lowest := lowesti.(*Item)
		for _, n := range g.neighbors(lowest.point) {
			ten := costs[lowest.point] + g.data[n]
			if ten < costs[n] {
				costs[n] = ten
				pq.update(items[n], ten)
			}
		}
	}
	fmt.Println(items[point{g.maxX, g.maxY}].priority)
}

func solve2(startx, starty int, g *game) {
	q := internal.NewQueue[point]()
	lowestRisks := map[point]int{}
	for y := 1; y <= g.maxY; y++ {
		for x := 1; x <= g.maxX; x++ {
			lowestRisks[point{x, y}] = math.MaxInt
		}
	}
	lowestRisks[point{1, 1}] = 0
	q.Enqueue(point{1, 1})
	for !q.Empty() {
		node := q.Dequeue()
		neighbors := g.neighbors(node)
		for _, n := range neighbors {
			cost := lowestRisks[node] + g.data[n]
			if cost < lowestRisks[n] {
				lowestRisks[n] = cost
				q.Enqueue(n)
			}
		}
	}
	fmt.Println(lowestRisks[point{g.maxX, g.maxY}])
}

func solve(startx, starty int, g *game) {
	start := point{startx, starty}
	visited := map[point]struct{}{}
	unvisited := map[point]int{}
	costs := map[point]int{}
	for y := 1; y <= g.maxY; y++ {
		for x := 1; x <= g.maxX; x++ {
			costs[point{x, y}] = math.MaxInt
			unvisited[point{x, y}] = math.MaxInt
		}
	}
	costs[start] = 0
	unvisited[start] = 0
	pt := start
	count := 0
	for len(unvisited) > 0 {
		pt = lowestUnvisited(unvisited)
		if count%100 == 0 {
			fmt.Println("togo", len(unvisited))
		}
		for _, n := range g.neighbors(pt) {
			if _, ok := visited[n]; ok {
				continue
			}
			tentativeCost := costs[pt] + g.data[n]
			if tentativeCost < unvisited[n] {
				unvisited[n] = tentativeCost
				costs[n] = tentativeCost
			}
		}
		delete(unvisited, pt)
		visited[pt] = struct{}{}
		count++
	}
	fmt.Println(costs[point{g.maxX, g.maxY}])
}

func lowestUnvisited(points map[point]int) point {
	var low point
	lowest := math.MaxInt
	for point, value := range points {
		if value < lowest {
			low = point
			lowest = value
		}
	}
	return low
}

type game struct {
	data       map[point]int
	maxX, maxY int
}

func newGame(maxx, maxy int) *game {
	return &game{
		data: map[point]int{},
		maxX: maxx,
		maxY: maxy,
	}
}

func (g *game) neighbors(p point) []point {
	out := []point{}
	if p.x != 0 {
		out = append(out, point{p.x - 1, p.y})
	}
	if p.x != g.maxX {
		out = append(out, point{p.x + 1, p.y})
	}
	if p.y != 0 {
		out = append(out, point{p.x, p.y - 1})
	}
	if p.y != g.maxY {
		out = append(out, point{p.x, p.y + 1})
	}
	//	fmt.Println(out)
	return out
}
func (g *game) String() string {
	out := ""
	for j := 1; j <= g.maxY; j++ {
		for i := 1; i <= g.maxX; i++ {
			out += fmt.Sprint(g.data[point{i, j}])
		}
		out += "\n"
	}
	return out
}

// An Item is something we manage in a priority queue.
type Item struct {
	point    point
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}
