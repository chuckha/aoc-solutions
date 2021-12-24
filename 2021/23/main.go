package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type item string

/*
#############
#AA.A.....AD#
###D#B#.#.###
  #D#B#C#.#
  #D#B#C#.#
  #C#B#C#.#
  #########

*/
func main() {
	lines := internal.ReadRawInput()
	extra1 := "  #D#C#B#A#"
	extra2 := "  #D#B#A#C#"
	newlines := make([]string, 0)
	newlines = append(newlines, lines[0:3]...)
	newlines = append(newlines, extra1, extra2)
	newlines = append(newlines, lines[3:]...)
	cave := map[point]string{}
	for y, line := range newlines {
		for x, c := range line {
			cave[point{x, y}] = string(c)
		}
	}
	igs := initGameState(cave)
	fmt.Println(stackSolve(igs))
}

type gamestate struct {
	data     map[point]string
	maxy     int
	roomsize int
	cost     int
	moves    []move
}

func (gs gamestate) impossible() bool {
	/*
		#############
		#AC........A#
		###D#.#.#D###
		  #D#B#.#A#
		  #D#B#.#C#
		  #C#B#C#B#
		  #########
	*/
	// home column A
	// home column D
	// if there is a D in column A and an A blocking a D from getting to its column, the board is impossible
	//
	// are out of place letters capable of reaching their hc or are they blocked by a letter in a column they are in?
	// for each room
	//   get the top letter
	//   if the letter is in its final place, continue
	//   if the top one has no viable moves, continue
	//   and to get to its hc it
	//
	hcs := gs.homeCoordinate("A")
	for _, hc := range hcs {
		if gs.data[hc] == "D" && gs.data[point{4, 1}] == "A" {
			return true
		}
	}
	hcs = gs.homeCoordinate("D")
	for _, hc := range hcs {
		if gs.data[hc] == "A" && gs.data[point{8, 1}] == "D" {
			return true
		}
	}
	return false
}

/*
#############
#.....C....B#
###D#A#.#D###
  #D#C#.#A#
  #D#B#.#C#
  #C#A#B#B#
  #########

*/
func (gs gamestate) validMovesForPoint(from point) []move {
	possibleMoves := make([]move, 0)
	// add valid hallway moves
	possibleDestinations := []point{{1, 1}, {11, 1}, {2, 1}, {10, 1}, {4, 1}, {6, 1}, {8, 1}}
	for _, p := range possibleDestinations {
		// if amphi is in a hallway, it can't move again until it's going home
		if from.y == 1 {
			continue
		}
		pathlen := gs.validPath(from, p)
		if pathlen == -1 {
			continue
		}
		if gs.data[p] == "." {
			newm := move{
				from: from,
				to:   p,
				cost: cost(gs.data[from]) * pathlen,
			}
			possibleMoves = append(possibleMoves, newm)
		}
	}

	//	fmt.Println("printing possible moves", possibleMoves)
	return possibleMoves
}

func (gs gamestate) isInPlace(p point) bool {
	if p.y == 1 {
		return false
	}
	letter := gs.data[p]
	// assumption, home coordinate returns points in order from highest to lowest
	hc := gs.homeCoordinate(letter)
	idx := internal.Search(p, hc)
	if idx == -1 {
		// this means the letter is in the wrong column
		return false
	}
	for i := idx; i < len(hc); i++ {
		if gs.data[hc[i]] != letter {
			return false
		}
	}
	return true
}

func (gs gamestate) String() string {
	var out strings.Builder
	for y := 0; y <= gs.maxy; y++ {
		for x := 0; x <= 12; x++ {
			out.WriteString(gs.data[point{x, y}])
		}
		out.WriteString("\n")
	}
	return out.String()
}

func initGameState(state map[point]string) gamestate {
	maxy := 0
	for p := range state {
		if p.y > maxy {
			maxy = p.y
		}
	}
	return gamestate{
		data:     state,
		maxy:     maxy,
		roomsize: maxy - 2,
		cost:     0,
	}
}

func (gs gamestate) homeCoordinate(letter string) []point {
	heh := "   A B C D"
	i := strings.Index(heh, letter)
	out := []point{}
	for k := 0; k < gs.roomsize; k++ {
		out = append(out, point{i, 2 + k})
	}
	return out
}

func (gs gamestate) move(m move) gamestate {
	newgs := gs.copy()
	tmp := newgs.data[m.from]
	newgs.data[m.from] = "."
	newgs.data[m.to] = tmp
	newgs.cost = gs.cost + m.cost
	newgs.moves = append(newgs.moves, m)
	return newgs
}

func (gs gamestate) copy() gamestate {
	data := map[point]string{}
	for k, v := range gs.data {
		data[k] = v
	}
	moves := make([]move, len(gs.moves))
	copy(moves, gs.moves)
	return gamestate{
		data:     data,
		cost:     gs.cost,
		maxy:     gs.maxy,
		roomsize: gs.roomsize,
		moves:    moves,
	}
}

func (gs gamestate) validPath(from, to point) int {
	q := internal.NewQueue[[]point]()
	q.Enqueue([]point{from})
	visited := map[point]struct{}{}
	visited[from] = struct{}{}
	for !q.Empty() {
		cur := q.Dequeue()
		for _, n := range cur[len(cur)-1].neighbors() {
			if n == to {
				return len(cur)
			}
			if _, ok := visited[n]; ok {
				continue
			}
			if gs.data[n] == "." {
				np := make([]point, len(cur))
				copy(np, cur)
				np = append(np, n)
				q.Enqueue(np)
				visited[n] = struct{}{}
			}
		}
	}
	return -1
}

type moveSlice []move

func (m moveSlice) Len() int      { return len(m) }
func (m moveSlice) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m moveSlice) Less(i, j int) bool {
	return m[i].cost < m[j].cost
}

func (gs gamestate) moveToFinalPlace(p point) (move, bool) {
	// get the letter we're inspecting
	letter := gs.data[p]
	// find the home coordinates of that letter (a -> x=3, b => x=5...)
	hcs := gs.homeCoordinate(letter)
	for i := len(hcs) - 1; i >= 0; i-- {
		val := gs.data[hcs[i]]
		// if there's a letter that's not the right letter in this stack (a B in A's stack)
		if val != letter && val != "." {
			// bail, we can't move this point to its final place
			return move{}, false
		}
		// otherwise if it's the letter keep going
		if val == letter {
			continue
		}
		// and if it's a . and the path to it is valid
		plen := gs.validPath(p, hcs[i])
		if val == "." && plen != -1 {
			return move{
				from: p,
				to:   hcs[i],
				cost: cost(gs.data[p]) * plen,
			}, true
		}
	}
	return move{}, false
}

func (gs gamestate) validMoves() []move {
	moves := make([]move, 0)
	for _, amphi := range gs.allAmphipods() {
		if gs.isInPlace(amphi) {
			continue
		}
		// if any amphi can go to final place, do it
		if m, ok := gs.moveToFinalPlace(amphi); ok {
			//			fmt.Println("m", m, "amphi", amphi)
			return []move{m}
		}
		moves = append(moves, gs.validMovesForPoint(amphi)...)

	}
	sort.Sort(moveSlice(moves))
	return moves
}

func (gs gamestate) allAmphipods() []point {
	out := []point{}
	for k, v := range gs.data {
		switch v {
		case "A", "B", "C", "D":
			out = append(out, k)
		}
	}
	return out
}

func (gs gamestate) solved() bool {
	for _, c := range []string{"A", "B", "C", "D"} {
		hcs := gs.homeCoordinate(c)
		for _, hc := range hcs {
			if gs.data[hc] != c {
				return false
			}
		}
	}
	return true
}

func stackSolve(initgs gamestate) int {
	s := internal.NewStack[gamestate]()
	s.Push(initgs)
	min := 99999999
	for !s.Empty() {
		gs := s.Pop()
		fmt.Println("starting with")
		fmt.Println(gs)
		validMoves := gs.validMoves()
		// fmt.Println("valid moves for above board")
		// for _, m := range validMoves {
		// 	fmt.Println(m)
		// }
		for _, move := range validMoves {
			// skip cycles
			if len(gs.moves) > 0 {
				if move.to == gs.moves[len(gs.moves)-1].from &&
					move.from == gs.moves[len(gs.moves)-1].to {
					continue
				}
			}
			newgs := gs.move(move)
			if newgs.impossible() {
				// fmt.Println("impossible")
				// fmt.Println(newgs)
				continue
			}
			fmt.Println("moved", move)
			fmt.Println(newgs)
			if newgs.solved() {
				if newgs.cost < min {
					min = newgs.cost
					fmt.Println("new minimum", min)
				}
				continue
			}
			// if the current cost is more than our existing min, ignore it
			if newgs.cost > min {
				continue
			}
			s.Push(newgs)
		}
	}
	return min
}

func solve(initgs gamestate) int {
	var _solve func(gs gamestate) int
	min := 9999999
	gamestates := map[string]struct{}{}
	_solve = func(gs gamestate) int {
		if gs.solved() {
			fmt.Println(gs)
			return gs.cost
		}
		if gs.cost > min {
			return min
		}
		//		fmt.Println("$$$", gs.cost, "$$$")
		validMoves := gs.validMoves()

		fmt.Println("base board")
		fmt.Println(gs)
		// for _, m := range validMoves {
		// 	fmt.Println("\t", m)
		// }

		fmt.Println("next moves")
		for _, move := range validMoves {
			// ignore cycles
			if len(gs.moves) > 0 && move.to == gs.moves[len(gs.moves)-1].from {
				continue
			}
			newgs := gs.move(move)
			if _, ok := gamestates[newgs.String()]; ok {
				continue
			}
			fmt.Println(move)
			fmt.Println(newgs)
			gamestates[newgs.String()] = struct{}{}
			cost := _solve(newgs)
			if cost < min {
				min = cost
			}
		}
		return min
	}
	return _solve(initgs)
}

type move struct {
	cost int
	from point
	to   point
}

func (m move) String() string {
	return fmt.Sprintf("from: %v to: %v (%d)", m.from, m.to, m.cost)
}

func cost(a string) int {
	switch a {
	case "A":
		return 1
	case "B":
		return 10
	case "C":
		return 100
	case "D":
		return 1000
	default:
		panic("ehhh")
	}
}

// ----------------
type cavern struct {
	data       map[point]item
	maxx, maxy int
}

type point struct {
	x, y int
}

func (p point) neighbors() []point {
	return []point{
		{p.x + 1, p.y},
		{p.x - 1, p.y},
		{p.x, p.y + 1},
		{p.x, p.y - 1},
	}
}

func (p point) inFrontOfRoom() bool {
	switch p.x {
	case 3, 5, 7, 9:
		return p.y == 1
	default:
		return false
	}
}

/*

#############
#...........#
###D#A#C#D###
  #D#C#B#A#
  #D#B#A#C#
  #C#A#B#B#
  #########

#############
#AC.A......A#
###D#.#.#D###
  #D#B#.#A#
  #D#B#.#C#
  #C#B#C#B#
  #########
*/
