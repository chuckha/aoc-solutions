package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	/*
		The first floor contains a polonium generator, a thulium generator, a thulium-compatible microchip, a promethium generator, a ruthenium generator, a ruthenium-compatible microchip, a cobalt generator, and a cobalt-compatible microchip.
		The second floor contains a polonium-compatible microchip and a promethium-compatible microchip.
		The third floor contains nothing relevant.
		The fourth floor contains nothing relevant.
	*/
	lines := internal.ReadInput()
	floorData := make([]floor2, 0)
	for _, line := range lines {
		floorData = append(floorData, parseItems2(line))
	}
	// part 2
	floorData[0] = append(floorData[0], "ELG", "ELM", "DIG", "DIM")
	building := newBuilding2(floorData)
	//	building.validMoves()
	fmt.Println(queueSolve(building))
}

func parseItems2(line string) floor2 {
	words := strings.Split(line, " ")
	components := []string{}
	for i, word := range words {
		word = strings.Trim(word, ",.")
		if word == "generator" {
			components = append(components, strings.ToUpper(string(words[i-1][:2]))+"G")
		}
		if word == "microchip" {
			broken := strings.Split(words[i-1], "-")
			components = append(components, strings.ToUpper(string(broken[0][:2]))+"M")
		}
	}
	return components
}

type building2 struct {
	floors   []floor2
	elevator int
	cost     int
}

func (b building2) componentCount() int {
	sum := 0
	for _, f := range b.floors {
		sum += len(f)
	}
	return sum
}

func (b building2) String() string {
	var out strings.Builder
	for i := len(b.floors) - 1; i >= 0; i-- {
		sort.Sort(sort.StringSlice(b.floors[i]))
		floor := []string{fmt.Sprintf("F%d", i+1)}
		if b.elevator == (i) {
			floor = append(floor, "E")
		} else {
			floor = append(floor, ".")
		}
		for _, c := range b.floors[i] {
			floor = append(floor, c)
		}
		for i := len(floor) - 1; i <= b.componentCount(); i++ {
			floor = append(floor, " . ")
		}
		out.WriteString(strings.Join(floor, "   "))
		out.WriteString("\n")
	}
	return out.String()
}

func queueSolve(b building2) int {
	q := internal.NewQueue[building2]()
	q.Enqueue(b)
	// states are unique with a min cost
	seen := map[string]struct{}{}
	min := 9999999
	for !q.Empty() {
		cur := q.Dequeue()
		// if _, ok := seen[cur.String()]; !ok {
		// 	seen[cur.String()] = cur
		// }
		validMoves := cur.validMoves()
		for _, move := range validMoves {
			newB := cur.applyMove(move)
			// make the move
			if _, ok := seen[newB.String()]; ok {
				continue
			}
			if newB.cost > min {
				continue
			}
			if newB.solved() {
				if newB.cost < min {
					min = newB.cost
				}
				continue
			}
			// fmt.Println(move)
			// fmt.Println(newB)
			// fmt.Println()
			seen[newB.String()] = struct{}{}
			q.Enqueue(newB)
		}
		fmt.Println(len(q.Internal()))
	}
	return min
}

func ssolve(b building2) int {
	seen := map[string]struct{}{}
	var solve func(building2) int

	solve = func(bb building2) int {
		if bb.solved() {
			return 0
		}
		min := 99999999
		seen[bb.String()] = struct{}{}
		for _, m := range bb.validMoves() {
			newb := bb.applyMove(m)
			if _, ok := seen[newb.String()]; ok {
				continue
			}
			//			fmt.Println(newb, m)
			solution := 1 + solve(newb)
			if newb.solved() {
				if solution < min {
					min = solution
				}
			}
		}
		return min
	}
	return solve(b)
}

func (b building2) solved() bool {
	return len(b.floors[3]) == b.componentCount()
}

func (b building2) valid() bool {
	for _, floor := range b.floors {
		//		fmt.Println("floor", i)
		if !floor.valid() {
			//			fmt.Println(b)
			return false
		}
	}
	return true
}

func (b building2) validMoves() []move {
	out := make([]move, 0)
	currentFloor := b.elevator
	combos := combos(b.floors[currentFloor])
	///	fmt.Println("all combos:", combos)
	for _, combo := range combos {
		if b.elevator < 3 {
			// move the combo up
			m := move{
				elevatorFrom: b.elevator,
				elevatorTo:   b.elevator + 1,
				components:   combo,
			}
			newb := b.applyMove(m)
			if newb.valid() {
				out = append(out, m)
			}
		}
		if b.elevator >= 1 {
			// if the floor below is completely empty, don't go down
			if len(b.floors[b.elevator-1]) == 0 {
				continue
			}
			// never move two things down
			if len(combo) == 2 {
				continue
			}
			// move the combo down
			m := move{
				elevatorFrom: b.elevator,
				elevatorTo:   b.elevator - 1,
				components:   combo,
			}
			newb := b.applyMove(m)
			if newb.valid() {
				out = append(out, m)
			}
		}
		// apply a move to a building
	}
	//	fmt.Println("valid moves", out)
	return out
}

func (b building2) equal(b2 building2) bool {
	if b.elevator != b2.elevator {
		return false
	}
	if len(b.floors) != len(b2.floors) {
		return false
	}
	for i, f := range b.floors {
		if !internal.EqualSlice(f, b2.floors[i]) {
			return false
		}
	}
	return true
}

func (b building2) copy() building2 {
	f2 := make([]floor2, len(b.floors))
	for i, f := range b.floors {
		f2[i] = make([]string, len(f))
		copy(f2[i], f)
	}
	return building2{
		elevator: b.elevator,
		floors:   f2,
		cost:     b.cost,
	}
}

func (b building2) applyMove(m move) building2 {
	// fmt.Println()
	// fmt.Println(b)
	// fmt.Println(m)
	newb := b.copy()
	newb.cost = newb.cost + 1
	newb.elevator = m.elevatorTo
	for _, c := range m.components {
		//		fmt.Println("before appending to the new floor", newb.floors[m.elevatorTo])
		newb.floors[m.elevatorTo] = append(newb.floors[m.elevatorTo], c)
		//		fmt.Println("after appending", newb.floors[m.elevatorTo])
		idx := internal.Search(c, newb.floors[m.elevatorFrom])
		if idx == -1 {
			panic("very bad")
			//			fmt.Println("should not be here")
			//			fmt.Println(c, newb.floors[m.elevatorFrom])
			//			continue
		}
		//		fmt.Println("before removing from the old floor", newb.floors[m.elevatorFrom])
		newb.floors[m.elevatorFrom] = append(newb.floors[m.elevatorFrom][:idx], newb.floors[m.elevatorFrom][idx+1:]...)
		//		fmt.Println("after removing from the old floor", newb.floors[m.elevatorFrom])
	}
	// fmt.Println(newb)
	// fmt.Println()
	return newb
}

type floor2 []string

func (f floor2) valid() bool {
	hasGenerator := false
	for _, c := range f {
		if strings.HasSuffix(c, "G") {
			hasGenerator = true
			break
		}
	}
	// if there is no generator, the floor is safe
	if !hasGenerator {
		//		fmt.Println("valid because it has no generator")
		return true
	}

	// there is at least one generator
	chips := internal.Set[string]{}
	gens := internal.Set[string]{}
	for _, component := range f {
		// find all chips on floor
		if strings.HasSuffix(component, "M") {
			chips.Insert(component[0:2], component[0:2])
		}
		if strings.HasSuffix(component, "G") {
			gens.Insert(component[0:2], component[0:2])
		}
	}
	// it's valid if there are no chips at all
	if len(chips) == 0 {
		//		fmt.Println("valid because there are no chips")
		return true
	}

	gensWithMatchingChips := len(gens.Intersect(chips))
	if gensWithMatchingChips != len(chips) {
		//		fmt.Println("not all chips have matching generators")
		return false
	}
	//	fmt.Println(chips, gens)
	//	fmt.Println("valid because all chips have matching generators")
	return true
	// if every generator has a matching chip, then it's a safe floor
}

func matchPrefix(a, b string) bool {
	return a[0:2] == b[0:2]
}

func newBuilding2(floorData []floor2) building2 {
	return building2{
		floors:   floorData,
		elevator: 0,
	}
}

type move struct {
	elevatorFrom, elevatorTo int
	components               []string
}

// combos returns every len1, len2 combination of input list
func combos(input1 []string) [][]string {
	out := [][]string{}

	// just double elements of the list
	for i := 0; i < len(input1)-1; i++ {
		for j := i + 1; j < len(input1); j++ {
			out = append(out, []string{input1[i], input1[j]})
		}
	}
	// Just single element of the list
	for _, s := range input1 {
		out = append(out, []string{s})
	}
	return out
}
