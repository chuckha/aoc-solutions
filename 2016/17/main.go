package main

import (
	"crypto/md5"
	"fmt"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	state := internal.ReadInput()[0]
	m := &maze{
		data: make(map[point]struct{}),
		maxx: 3,
		maxy: 3,
	}
	fmt.Println(m.solve(point{0, 0, state})[point{x: 3, y: 3}])
}

func (m *maze) solve(start point) map[point]int {
	// let's say you can't visit a room more than 5 times
	cost := map[point]int{}
	// the first time we see a room, it must cost maxint to get there
	// end when x, y == 3,3
	cost[start] = 0
	q := internal.NewQueue[point]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		fmt.Println(cur)
		for _, ne := range neighbors(cur) {
			// Ignore the out of bound neighbors...
			if ne.x > m.maxx || ne.y > m.maxy || ne.x < 0 || ne.y < 0 {
				continue
			}
			if _, ok := cost[ne]; !ok {
				cost[ne] = cost[cur] + 1
			}
			if ne.x == m.maxx && ne.y == m.maxy {
				if cost[point{x: m.maxx, y: m.maxy}] < cost[ne] {
					cost[point{x: m.maxx, y: m.maxy}] = cost[ne]
				}
				continue
			}
			q.Enqueue(ne)
		}
	}
	return cost
}

func neighbors(p point) []point {
	out := []point{}
	doors := md5hash(p.code)[:4]
	for i, c := range doors {
		switch c {
		case 'b', 'c', 'd', 'e', 'f':
			out = append(out, dir(i, p))
		default:
			// locked
		}
	}
	return out
}

func dir(i int, p point) point {
	switch i {
	case 0:
		return point{x: p.x, y: p.y - 1, code: p.code + "U"}
	case 1:
		return point{x: p.x, y: p.y + 1, code: p.code + "D"}
	case 2:
		return point{x: p.x - 1, y: p.y, code: p.code + "L"}
	case 3:
		return point{x: p.x + 1, y: p.y, code: p.code + "R"}
	default:
		panic("ehhhhhh")
	}
}

func md5hash(in string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(in)))
}

type maze struct {
	data       map[point]struct{}
	maxx, maxy int
}

type point struct {
	x, y int
	code string
}
