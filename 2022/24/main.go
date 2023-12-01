package main

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
298: too high
230: correct

pt 2:
1424: too high

// OOPS BUG: the bests[i] keep track of the sum, so the last value is the actual sum of the entire trip.
*/

func main() {
	maze := newWinterMaze()
	for y, line := range internal.ReadInput() {
		for x, c := range line {
			if c == '.' {
				continue
			}
			if c == '#' {
				maze.addWall(internal.Point{x, y})
				continue
			}
			maze.addBlizzard(newBlizzard(internal.Point{x, y}, c))
		}
	}

	conditions := []cond{
		{start: maze.entrance(), goal: maze.exit()},
		{start: maze.exit(), goal: maze.entrance()},
		{start: maze.entrance(), goal: maze.exit()},
	}
	bests := make([]int, len(conditions))

	for i, cond := range conditions {
		goal := cond.goal
		fmt.Println("New Goal: ", goal, bests)
		bests[i] = math.MaxInt
		q := make(internal.FPQ[*player], 0)
		heap.Init(&q)
		startPlayer := initPlayer(cond.start)
		if i > 0 {
			startPlayer = initPlayerWithtime(cond.start, bests[i-1])
		}
		heap.Push(&q, &internal.HeapItem[*player]{
			Value:    startPlayer,
			Priority: startPlayer.pos.ManhattanDistance(goal),
		})
		visited := make(map[id]struct{})
		for q.Len() > 0 {
			p := heap.Pop(&q).(*internal.HeapItem[*player]).Value
			// if bests[i] == 481 {
			// 	fmt.Println(q.Len())
			// 	fmt.Println(p.pos, p.time)
			// }
			visited[id{p.time, p.pos}] = struct{}{}
			//		fmt.Println("player priority", p.time)
			// fmt.Println(p.path, p.pos, p.time)

			if p.pos.ManhattanDistance(goal)+p.time >= bests[i] {
				continue
			}

			if p.pos == goal {
				if p.time < bests[i] {
					bests[i] = p.time
					// fmt.Println(p.path)
					// fmt.Println("items in heap", q.Len())
					fmt.Println("best time so far", bests[i])
				}
				//			fmt.Println("an exit", p.path, p.time)
				continue
			}

			if _, ok := maze.timing[p.time+1]; !ok {
				for !ok {
					maze.tick()
					_, ok = maze.timing[p.time+1]
				}
			}

			m := maze.timing[p.time+1]
			options := decisions(maze.layout, m, p.pos)
			// fmt.Println(q.Len(), p.pos.ManhattanDistance(goal), options)
			// fmt.Println(maze.Print(p, m))

			for _, decision := range options {
				np := p.copy()
				np.move(decision)
				if _, ok := visited[id{np.time, np.pos}]; ok {
					continue
				}
				// fmt.Println(p.time, decision, np.path, np.pos)
				// fmt.Println(maze.Print(np, m))
				heap.Push(&q, &internal.HeapItem[*player]{
					Value:    np,
					Priority: np.pos.ManhattanDistance(goal),
				})
			}
			// time.Sleep(1 * time.Second)
		}
		// end of for loop (q.Len() > 0)
	}
	sum := 0
	for _, b := range bests {
		sum += b
	}
	fmt.Println(bests, sum)
}

type cond struct {
	start internal.Point
	goal  internal.Point
}

type id struct {
	time int
	pos  internal.Point
}

type decision string

var cardinalDecisions = []decision{downD, rightD, upD, leftD}

const (
	waitD  decision = "wait"
	leftD  decision = "left"
	upD    decision = "up"
	rightD decision = "right"
	downD  decision = "down"
)

// decisions returns a valid list of decisions for a player to make
// the maze it takes in has already made the moves
// 1. if there is no blizzard on current pos you can wait
// 2. if there is no blizzard around you , you can go in that direction
func decisions(layout *internal.GridV3[internal.Cell], blizzard *internal.GridV3[*blizzard], player internal.Point) []decision {
	out := make([]decision, 0)
	if !blizzard.In(player) {
		out = append(out, waitD)
	}
	for _, dir := range cardinalDecisions {
		pt := player.Direction(string(dir))
		if pt.Y < layout.Min.Y {
			continue
		}
		if pt.Y > layout.Max.Y {
			continue
		}
		// this direction is a wall
		if layout.In(pt) {
			continue
		}
		// there is a blizzard there
		if blizzard.In(pt) {
			continue
		}
		out = append(out, dir)
	}
	return out
}

type player struct {
	path []internal.Point
	pos  internal.Point
	time int
}

func initPlayer(pt internal.Point) *player {
	return &player{
		path: make([]internal.Point, 0),
		pos:  pt,
		time: 0,
	}
}

func initPlayerWithtime(pt internal.Point, t int) *player {
	return &player{
		path: make([]internal.Point, 0),
		pos:  pt,
		time: t,
	}
}

func (p *player) String() string {
	return "E"
}

func (p *player) move(dir decision) {
	if dir != waitD {
		p.pos = p.pos.Direction(string(dir))
	}
	p.path = append(p.path, p.pos)
	p.time++
}

func (p *player) copy() *player {
	p2 := make([]internal.Point, len(p.path))
	copy(p2, p.path)
	return &player{
		path: p2,
		pos:  p.pos,
		time: p.time,
	}
}

type winterMaze struct {
	ticktimer    int
	player       *player
	layout       *internal.GridV3[internal.Cell]
	curBlizzards *internal.GridV3[*blizzard]
	blizzards    []*blizzard
	timing       map[int]*internal.GridV3[*blizzard]
}

func (w *winterMaze) entrance() internal.Point {
	for x := w.layout.Min.X; x <= w.layout.Max.X; x++ {
		pt := internal.Point{x, w.layout.Min.Y}
		if !w.layout.In(pt) {
			return pt
		}
	}
	panic("no entry point")
}
func (w *winterMaze) exit() internal.Point {
	for x := w.layout.Min.X; x <= w.layout.Max.X; x++ {
		pt := internal.Point{x, w.layout.Max.Y}
		if !w.layout.In(pt) {
			return pt
		}
	}
	panic("no exit point")
}

func newWinterMaze() *winterMaze {
	return &winterMaze{
		ticktimer:    0,
		layout:       internal.NewGridV3[internal.Cell](),
		curBlizzards: internal.NewGridV3[*blizzard](),
		blizzards:    make([]*blizzard, 0),
		timing:       make(map[int]*internal.GridV3[*blizzard]),
	}
}

func (w *winterMaze) String() string {
	cp2 := w.curBlizzards.CopyAsString()
	cp := w.layout.CopyAsString()
	layered := cp.Layer(cp2)
	track := map[internal.Point]int{}
	for _, bliz := range w.blizzards {
		track[bliz.pt]++
	}
	for pt, ct := range track {
		if ct > 1 {
			layered.Set(pt, internal.Cell(fmt.Sprintf("%d", ct)))
		}
	}
	return layered.String()
}

func (w *winterMaze) Print(p *player, m *internal.GridV3[*blizzard]) string {
	cp2 := m.CopyAsString()
	cp := w.layout.CopyAsString()
	layered := cp.Layer(cp2)
	track := map[internal.Point]int{}
	for _, bliz := range w.blizzards {
		track[bliz.pt]++
	}
	for pt, ct := range track {
		if ct > 1 {
			layered.Set(pt, internal.Cell(fmt.Sprintf("%d", ct)))
		}
	}
	layered.Set(p.pos, internal.Cell(internal.Red(p.String())))
	return layered.String()
}

func (w *winterMaze) addWall(pt internal.Point) {
	w.layout.Set(pt, internal.Cell("#"))
}

func (w *winterMaze) addBlizzard(b *blizzard) {
	w.curBlizzards.Set(b.pt, b)
	w.blizzards = append(w.blizzards, b)
}

func (w *winterMaze) tick() {
	w.ticktimer++
	for _, b := range w.blizzards {
		w.curBlizzards.Clear(b.pt)
		w.moveForward(b)
	}
	for _, b := range w.blizzards {
		w.curBlizzards.Set(b.pt, b)
	}
	w.timing[w.ticktimer] = w.curBlizzards.Copy()
}

func (w *winterMaze) isWall(pt internal.Point) bool {
	return w.layout.In(pt)
}

type blizzard struct {
	dir dir
	pt  internal.Point
}

func (w *winterMaze) moveForward(b *blizzard) {
	switch b.dir {
	case up:
		b.pt = b.pt.Up()
	case right:
		b.pt = b.pt.Right()
	case down:
		b.pt = b.pt.Down()
	case left:
		b.pt = b.pt.Left()
	}
	if !w.isWall(b.pt) {
		return
	}
	switch b.dir {
	case up:
		b.pt = internal.Point{b.pt.X, w.layout.Max.Y - 1}
	case right:
		b.pt = internal.Point{w.layout.Min.X + 1, b.pt.Y}
	case down:
		b.pt = internal.Point{b.pt.X, w.layout.Min.Y + 1}
	case left:
		b.pt = internal.Point{w.layout.Max.X - 1, b.pt.Y}
	}
}

func newBlizzard(pt internal.Point, c rune) *blizzard {
	bliz := &blizzard{pt: pt}
	switch c {
	case '>':
		bliz.dir = right
	case 'v':
		bliz.dir = down
	case '<':
		bliz.dir = left
	case '^':
		bliz.dir = up
	}
	return bliz
}

func (b *blizzard) String() string {
	return string(b.dir)
}

type dir string

const (
	none  dir = " "
	right dir = ">"
	down  dir = "v"
	left  dir = "<"
	up    dir = "^"
)
