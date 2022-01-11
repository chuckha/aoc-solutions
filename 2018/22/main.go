package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const (
	nothing = "neither"
	torch   = "torch"
	climb   = "climbing gear"
)

func main() {
	input := internal.ReadInput()
	words1 := strings.Split(input[0], " ")
	id, _ := strconv.Atoi(words1[1])
	words2 := strings.Split(input[1], " ")
	target := strings.Split(words2[1], ",")
	x, _ := strconv.Atoi(target[0])
	y, _ := strconv.Atoi(target[1])
	m := newMaze(id, point{x, y})
	// m.target = point{6, 50}
	//	fmt.Println(m)
	//	fmt.Println(m.risk())
	shortestPath := m.djikstraSorta()
	fmt.Println(m.PrintWithPath(shortestPath))
	//	fmt.Println(costs[point{0, 0}], costs[m.target])
	// for y := 0; y <= m.target.y; y++ {
	// 	for x := 0; x <= m.target.x; x++ {
	// 		fmt.Println(x, y, costs[point{x, y}])
	// 	}
	// }
}

type maze struct {
	data map[point]string

	inputDepth    int
	target        point
	erosionCache  map[point]int
	geologicCache map[point]int
	kindCache     map[point]string
}

type color func(string) []byte

func (m *maze) PrintWithPath(path []state) string {
	ll := map[point]state{}
	bigx := 0
	bigy := 0
	for _, p := range path {
		ll[p.pos] = p
		if p.pos.x > bigx {
			bigx = p.pos.x
		}
		if p.pos.y > bigy {
			bigy = p.pos.y
		}
	}
	var out strings.Builder
	c := green
	for j := 0; j <= bigy; j++ {
		for i := 0; i <= bigx; i++ {
			if s, ok := ll[point{i, j}]; ok {
				if s.equipedItem == torch {
					c = red
				}
				if s.equipedItem == climb {
					c = green
				}
				if s.equipedItem == nothing {
					c = blue
				}
				out.Write(c(m.kind(point{i, j})))
				continue
			}
			out.WriteString(m.kind(point{i, j}))
		}
		out.WriteString("\n")
	}
	return out.String()

}

type state struct {
	pos         point
	equipedItem string
}

func newState(pos point, item string) state {
	return state{
		pos:         pos,
		equipedItem: item,
	}
}

func (s state) changeItem(newItem string) state {
	return state{
		pos:         s.pos,
		equipedItem: newItem,
	}
}

func (s state) update(p point, item string) state {
	return state{
		pos:         p,
		equipedItem: item,
	}
}

func (m *maze) cost(cur, next state) int {
	equipmentCost := 0
	if cur.equipedItem != next.equipedItem {
		equipmentCost = 7
	}
	return 1 + equipmentCost
}

func (m *maze) possibleStates(cur state, next point) []state {
	if next == m.target {
		return []state{
			{pos: m.target, equipedItem: torch},
		}
	}
	curKind := m.kind(cur.pos)
	nextKind := m.kind(next)
	out := make([]state, 0)
	switch curKind {
	case ".":
		switch nextKind {
		case ".":
			return []state{
				cur.update(next, torch),
				cur.update(next, climb),
			}
		case "|":
			return []state{cur.update(next, torch)}
		case "=":
			return []state{cur.update(next, climb)}
		}
	case "|":
		switch nextKind {
		case ".":
			return []state{cur.update(next, torch)}
		case "|":
			return []state{
				cur.update(next, torch),
				cur.update(next, nothing),
			}
		case "=":
			return []state{cur.update(next, nothing)}
		}
	case "=":
		switch nextKind {
		case ".":
			return []state{cur.update(next, climb)}
		case "|":
			return []state{cur.update(next, nothing)}
		case "=":
			return []state{
				cur.update(next, nothing),
				cur.update(next, climb),
			}
		}

	}
	switch nextKind {
	case ".":
		out = append(out, cur.update(next, climb))
		out = append(out, cur.update(next, torch))
	case "|":
		out = append(out, cur.update(next, torch))
		out = append(out, cur.update(next, nothing))
	case "=":
		out = append(out, cur.update(next, climb))
		out = append(out, cur.update(next, nothing))
	}
	return out
}

type movement struct {
	pos  point
	item string
}

func newMaze(inputDepth int, target point) *maze {
	return &maze{
		data:          map[point]string{},
		inputDepth:    inputDepth,
		target:        target,
		erosionCache:  make(map[point]int),
		geologicCache: make(map[point]int),
		kindCache:     make(map[point]string),
	}
}

func (m *maze) String() string {
	var out strings.Builder
	for j := 0; j <= m.target.y; j++ {
		for i := 0; i <= m.target.x; i++ {
			out.WriteString(m.kind(point{i, j}))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func red(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;41m%s\033[0m", in))
}
func green(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;42m%s\033[0m", in))
}
func blue(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;44m%s\033[0m", in))
}

func (m *maze) risk() int {
	out := 0
	for j := 0; j <= m.target.y; j++ {
		for i := 0; i <= m.target.x; i++ {
			switch m.kind(point{i, j}) {
			case "=":
				out += 1
			case "|":
				out = out + 2
			}
		}
	}
	return out
}

func (m *maze) kind(p point) string {
	if item, ok := m.kindCache[p]; ok {
		return item
	}
	el := m.erosionLevel(p)
	switch el % 3 {
	case 0:
		m.kindCache[p] = "."
		return "."
	case 1:
		m.kindCache[p] = "="
		return "="
	case 2:
		m.kindCache[p] = "|"
		return "|"
	}
	panic("no bueno")
}

func (m *maze) erosionLevel(p point) int {
	if item, ok := m.erosionCache[p]; ok {
		return item
	}
	m.erosionCache[p] = (m.geologicIndex(p) + m.inputDepth) % 20183
	return m.erosionCache[p]
}

func (m *maze) geologicIndex(p point) int {
	if item, ok := m.geologicCache[p]; ok {
		return item
	}
	if p == (point{0, 0}) {
		m.geologicCache[p] = 0
		return m.geologicCache[p]
	}
	if p == m.target {
		m.geologicCache[p] = 0
		return m.geologicCache[p]
	}
	if p.y == 0 {
		m.geologicCache[p] = p.x * 16807
		return m.geologicCache[p]
	}
	if p.x == 0 {
		m.geologicCache[p] = p.y * 48271
		return m.geologicCache[p]
	}
	m.geologicCache[p] = m.erosionLevel(p.left()) * m.erosionLevel(p.up())
	return m.geologicCache[p]
}

type point struct {
	x, y int
}

func (p point) left() point {
	return point{p.x - 1, p.y}
}
func (p point) up() point {
	return point{p.x, p.y - 1}
}
func (p point) neighbors() []point {
	return []point{
		{p.x, p.y - 1}, {p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y + 1},
	}
}
func (p point) invalid() bool {
	return p.x < 0 || p.y < 0
}

func (p point) dist(p2 point) int {
	return int(math.Abs(float64(p.x-p2.x))) + int(math.Abs(float64(p.y-p2.y)))
}

func (m *maze) hueristic(state state) int {
	return state.pos.dist(m.target)
}

func (m *maze) djikstraSorta() []state {
	initialState := state{
		pos:         point{0, 0},
		equipedItem: torch,
	}
	frontier := internal.NewPriorityQueue[state]()
	cameFrom := map[state]state{initialState: {}}
	frontier.Add(initialState, 0)
	costSoFar := map[state]int{initialState: 0}
	for !frontier.Empty() {
		cur := frontier.Pull()
		if cur.pos == m.target {
			break
		}
		for _, n := range cur.pos.neighbors() {
			if n.invalid() {
				continue
			}

			// only generate possible states
			possibleStates := m.possibleStates(cur, n)
			for _, state := range possibleStates {
				//				fmt.Println("possible state", state)
				cost := m.cost(cur, state)
				//				fmt.Println("cur", cur, "state", state, "cost", cost)
				newCost := cost + costSoFar[cur]
				currentCost, ok := costSoFar[state]
				if !ok || newCost < currentCost {
					costSoFar[state] = newCost
					priority := newCost + m.hueristic(state)
					//					fmt.Printf("%v(%s)->%v(%s): %d with prio %d\n", cur.pos, cur.equipedItem, state.pos, state.equipedItem, cost, priority)
					frontier.Add(state, priority)
					cameFrom[state] = cur
				}
				// lowestState := possibleStates[0]
				// lowestCost := m.cost(cur, lowestState)
				// for _, s := range possibleStates[1:] {
				// 	nextCost := m.cost(cur, s)
				// 	if nextCost < lowestCost {
				// 		lowestCost = nextCost
				// 		lowestState = s
				// 	}
				// }
				// newCost := lowestCost + costSoFar[cur.pos]
				// ///			fmt.Println(lowestState, lowestCost)
				// currentCost, ok := costSoFar[n]
				// if !ok || newCost < currentCost {
				// 	costSoFar[n] = newCost
				// 	frontier.Add(lowestState)
				// 	cameFrom[lowestState] = cur
				// }
			}
		}
	}
	fmt.Println("did we finish?")
	//	fmt.Println(cost)
	//	fmt.Println(cameFrom)
	//	os.Exit(0)
	//	fmt.Println(cameFrom)
	finalState := state{
		pos:         m.target,
		equipedItem: torch,
	}

	fmt.Println("cost to final state", costSoFar[finalState])
	stack := internal.NewStack[state]()
	for {
		stack.Push(finalState)
		finalState = cameFrom[finalState]
		if finalState.pos.x == 0 && finalState.pos.y == 0 {
			break
		}
	}
	stack.Push(finalState)

	out := []state{}
	for !stack.Empty() {
		out = append(out, stack.Pop())
	}
	//	return out
	return out
}

type data struct {
	pos  point
	item string
}

type statePrioritizer struct {
	target point
}

func (s *statePrioritizer) Prioritize(state state) int {
	additionalCost := 0
	if state.equipedItem != torch {
		additionalCost = 7
	}
	return state.pos.dist(s.target) + additionalCost
}

/*
The region at 0,0 (the mouth of the cave) has a geologic index of 0.
The region at the coordinates of the target has a geologic index of 0.
If the region's Y coordinate is 0, the geologic index is its X coordinate times 16807.
If the region's X coordinate is 0, the geologic index is its Y coordinate times 48271.
Otherwise, the region's geologic index is the result of multiplying the erosion levels of the regions at X-1,Y and X,Y-1.
*/

// 736 too low?
// 982 too low
// 984 got it
// 1079 too high
