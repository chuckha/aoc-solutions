package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	g := internal.InitGrid()
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWZYZ"
	count := 0
	justLetters := map[internal.Point]string{}
	for _, line := range lines {
		p := toPoint(line)
		if p.X > g.Length {
			g.Length = p.X + 1
		}
		if p.Y > g.Height {
			g.Height = p.Y + 1
		}
		g.Data[p] = string(letters[count])
		justLetters[p] = string(letters[count])
		count++
	}
	for j := g.MinY; j < g.Height; j++ {
		for i := g.MinX; i < g.Length; i++ {
			distances := []int{}
			for p := range justLetters {
				dist := (internal.Point{i, j}).ManhattanDistance(p)
				distances = append(distances, dist)
			}
			if sum(distances) < 10000 {
				g.Data[internal.Point{i, j}] = "#"
			}

			// part 1
			// smallests := []thing{}
			// for p, l := range justLetters {
			// 	dist := (internal.Point{i, j}).ManhattanDistance(p)
			// 	smallests = append(smallests, thing{dist, l})
			// }
			// sort.Sort(things(smallests))
			// if smallests[0].dist == smallests[1].dist {
			// 	// its a tie
			// 	g.Data[internal.Point{i, j}] = "."
			// } else {
			// 	g.Data[internal.Point{i, j}] = smallests[0].letter
			// }

			// find manhattan distance to every point in g
			// find the smallest manhattan distance to every point in g
			// set the pointin g = to that one
		}
	}
	fmt.Println(g.On())

	// PART 1
	// edgeLetters := map[string]struct{}{}
	// for j := g.MinY; j < g.Height; j++ {
	// 	edgeLetters[g.At(0, j)] = struct{}{}
	// 	edgeLetters[g.At(g.Length-1, j)] = struct{}{}
	// }
	// for i := g.MinX; i < g.Length; i++ {
	// 	edgeLetters[g.At(i, 0)] = struct{}{}
	// 	edgeLetters[g.At(i, g.Height-1)] = struct{}{}
	// }
	// counts := map[string]int{}
	// for _, letter := range g.Data {
	// 	counts[letter]++
	// }
	// for l := range edgeLetters {
	// 	delete(counts, l)
	// }
	// thingsagain := []thing{}
	// for l, c := range counts {
	// 	thingsagain = append(thingsagain, thing{c, l})
	// }
	// sort.Sort(things(thingsagain))
	// // 4690 is too high (V)
	// fmt.Println(g)
	// for _, thing := range thingsagain {
	// 	fmt.Println(thing)
	// }
	// END PART 1

}
func sum(in []int) int {
	s := 0
	for _, v := range in {
		s += v
	}
	return s
}

type thing struct {
	dist   int
	letter string
}
type things []thing

func (t things) Len() int           { return len(t) }
func (t things) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t things) Less(i, j int) bool { return t[i].dist < t[j].dist }

func toPoint(line string) internal.Point {
	split := strings.Split(line, ", ")
	x, _ := strconv.Atoi(split[0])
	y, _ := strconv.Atoi(split[1])
	return internal.Point{x, y}
}
