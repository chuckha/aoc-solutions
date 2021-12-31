package main

import (
	"fmt"
	"math"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	grid := internal.NewGridFromInput(lines)

	fmt.Println(grid)
	nums := numbers(grid)
	// pull out 0
	numslice := []string{}
	for k := range nums {
		if k == "0" {
			continue
		}
		numslice = append(numslice, k)
	}
	costMaps := map[internal.Point]map[internal.Point]int{}
	for _, v := range nums {
		costMaps[v] = grid.Dijkstra(v)
		// for j := 0; j < grid.Height; j++ {
		// 	for i := 0; i < grid.Length; i++ {
		// 		if costMaps[v][internal.Point{i, j}] == math.MaxInt {
		// 			fmt.Print("###")
		// 			continue
		// 		}
		// 		fmt.Printf("%02d ", costMaps[v][internal.Point{i, j}])
		// 	}
		// 	fmt.Println()
		// }
		// fmt.Println(strings.Repeat("---", grid.Length))
	}

	// generate the shortest paths to each node from each other node
	// which is really just running Djikstra with various starting points and using the costs to calculate the total paths

	min := math.MaxInt
	for _, c := range heapsAlgo(numslice, len(numslice)) {
		// always start at 0
		path := append([]string{"0"}, c...)
		path = append(path, "0")
		cost := 0
		for i := 0; i < len(path)-1; i++ {
			// /			fmt.Println("nums", nums[path[i+1]])
			cost += costMaps[nums[path[i]]][nums[path[i+1]]]
		}
		if cost < min {
			min = cost
		}
		fmt.Printf("cost of this path: %v is %d\n", path, cost)
		// for each of these,
		// for each node
		//    find the shortest path to this node
		//     from this node find the shortest path to the next node
	}
	fmt.Println(min)
}

func solve(g *internal.Grid) {
	// visit 0, 1, 2, 3, 4
	// staring at 0 visit every other permutation of the other numbers
}

func numbers(g *internal.Grid) map[string]internal.Point {
	out := map[string]internal.Point{}
	for j := 0; j < g.Height; j++ {
		for i := 0; i < g.Length; i++ {
			if g.At(i, j) != "." && g.At(i, j) != "#" {
				out[g.Data[internal.Point{i, j}]] = internal.Point{i, j}
			}
		}
	}
	return out
}

func heapsAlgo(items []string, k int) [][]string {
	if k == 1 {
		next := make([]string, len(items))
		copy(next, items)
		return [][]string{next}
		//		return [][]string{items}
	}
	out := make([][]string, 0)
	for i := 0; i < k; i++ {
		out = append(out, heapsAlgo(items, k-1)...)
		if k%2 == 0 {
			items[0], items[k-1] = items[k-1], items[0]
		} else {
			items[i], items[k-1] = items[k-1], items[i]
		}
	}
	return out
}
