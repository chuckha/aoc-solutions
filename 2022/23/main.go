package main

import (
	"fmt"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
4068
17142 too high (this was after they all came to a halt)

pt2
967: too low

*/

func main() {
	input := internal.ReadInput()
	dirs := internal.NewQueue[string]()
	dirs.Enqueue("N", "S", "W", "E")

	ground := newGround(dirs)
	id := 0
	for y, line := range input {
		for x, c := range line {
			if c == '.' {
				continue
			}
			pt := internal.Point{x, y}
			e := newElf(id, pt)
			id++
			ground.layout.Set(pt, e)
			ground.elves = append(ground.elves, e)
		}
	}
	fmt.Println(ground.layout)
	i := 1
	for !ground.round() {
		fmt.Printf("== Round %d (%d)==\n", i, ground.emptyTiles())
		i++
	}
	fmt.Printf("== Round %d (%d)==\n", i, ground.emptyTiles())
	//	fmt.Println(ground.layout)
	fmt.Println(ground.emptyTiles())

}

type elf struct {
	id            int
	location      internal.Point
	consideration internal.Point
	doMove        bool
}

func newElf(id int, loc internal.Point) *elf {
	return &elf{
		id:       id,
		location: loc,
		doMove:   true,
	}
}

func (e *elf) String() string {
	return "#"
}

type ground struct {
	layout     *internal.GridV3[*elf]
	elves      []*elf
	directions *internal.Queue[string]
	// round variables
	considerations map[internal.Point]*elf
}

func (g *ground) round() bool {
	g.startRound()
	g.consider()
	if g.move() {
		return true
	}
	g.finish()
	return false
}

func (g *ground) startRound() {
	g.considerations = make(map[internal.Point]*elf)
	for _, e := range g.elves {
		e.consideration = internal.Point{}
		e.doMove = false
	}
}

func (g *ground) finish() {
	nextInspectionOrder(g.directions)
}

func (g *ground) consider() {
	order := inspectionOrder(g.directions)
	//	fmt.Println(order)
	for _, elf := range g.elves {
		if g.surroundedByEmpty(elf.location) {
			// elf does not move
			continue
		}
		g.considerElf(elf, order)
	}
}

func (g *ground) considerElf(e *elf, order []string) {
	for _, o := range order {
		if g.checkEmpty(e.location, o) {
			newPoint := e.location.Direction(o)
			// if someone's already considered it, don't move and tell them not to move
			if oldElf, ok := g.considerations[newPoint]; ok {
				oldElf.doMove = false
				return
			}
			e.doMove = true
			e.consideration = newPoint
			// fmt.Printf("elf at %v moving to %v\n", e.location, newPoint)
			g.considerations[newPoint] = e
			return
		}
	}
}

func (g *ground) checkEmpty(pt internal.Point, dir string) bool {
	switch dir {
	case "N":
		return !g.layout.In(pt.Direction("N")) && !g.layout.In(pt.Direction("NE")) && !g.layout.In(pt.Direction("NW"))
	case "S":
		return !g.layout.In(pt.Direction("S")) && !g.layout.In(pt.Direction("SE")) && !g.layout.In(pt.Direction("SW"))
	case "W":
		return !g.layout.In(pt.Direction("W")) && !g.layout.In(pt.Direction("NW")) && !g.layout.In(pt.Direction("SW"))
	case "E":
		return !g.layout.In(pt.Direction("E")) && !g.layout.In(pt.Direction("NE")) && !g.layout.In(pt.Direction("SE"))
	}
	panic("bad direction")
}

func (g *ground) surroundedByEmpty(pt internal.Point) bool {
	for _, dir := range pt.Surrounding() {
		if g.layout.In(dir) {
			return false
		}
	}
	return true
}

func (g *ground) move() bool {
	count := 0
	for _, e := range g.elves {
		if e.doMove {
			count++
			g.layout.Clear(e.location)
			e.location = e.consideration
			g.layout.Set(e.location, e)
		}
	}
	return count == 0
}

func (g *ground) emptyTiles() int {
	count := 0
	for y := g.layout.Min.Y; y <= g.layout.Max.Y; y++ {
		for x := g.layout.Min.X; x <= g.layout.Max.X; x++ {
			if !g.layout.In(internal.Point{x, y}) {
				count++
			}
		}
	}
	return count
}

func newGround(dirs *internal.Queue[string]) *ground {
	return &ground{
		layout:     internal.NewGridV3[*elf](),
		elves:      make([]*elf, 0),
		directions: dirs,
	}
}

func inspectionOrder(q *internal.Queue[string]) []string {
	out := []string{}
	for i := 0; i < len(q.Internal()); i++ {
		item := q.Dequeue()
		out = append(out, item)
		q.Enqueue(item)
	}
	return out
}

func nextInspectionOrder(q *internal.Queue[string]) {
	q.Enqueue(q.Dequeue())
}
