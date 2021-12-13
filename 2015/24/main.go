package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	weights := make([]int, len(lines))
	for i, line := range lines {
		w, _ := strconv.Atoi(line)
		weights[i] = w
	}
	singleGroupWeight := sum(weights) / 4
	sort.Sort(sort.Reverse(sort.IntSlice(weights)))
	minimizeRecSplit(singleGroupWeight, []int{}, weights)
	// allSolutions := [][]int{}
	// out, ok := greedySplit(singleGroupWeight, []int{}, []int{}, []int{}, weights)
	// if !ok {
	// 	panic("not ok")
	// }
	// allSolutions = append(allSolutions, out[0])
	// for i := 0; i < len(weights); i++ {
	// 	for j := i + 1; j < len(weights); j++ {
	// 		fix := []int{}
	// 		weightCopy := make([]int, len(weights))
	// 		copy(weightCopy, weights)
	// 		fix = append(fix, weights[i])
	// 		fix = append(fix, weightCopy[j])
	// 		weightCopy = append(weightCopy[:i], append(weightCopy[i+1:j], weightCopy[j+1:]...)...)
	// 		fmt.Println(fix, weightCopy)
	// 		out, good := greedySplit(singleGroupWeight, fix, []int{}, []int{}, weightCopy)
	// 		if !good {
	// 			continue
	// 		}
	// 		allSolutions = append(allSolutions, out[0])
	// 	}
	// 	// fix1 := []int{}
	// 	// weightCopy := make([]int, len(weights))
	// 	// copy(weightCopy, weights)
	// 	// fix1 = append(fix1, weights[i])
	// 	// weightCopy = append(weightCopy[i+1:], weightCopy[:i]...)
	// 	// fmt.Println(fix1, weightCopy)
	// 	// out, good := greedySplit(singleGroupWeight, fix1, []int{}, []int{}, weightCopy)
	// 	// if !good {
	// 	// 	continue
	// 	// }
	// 	// allSolutions = append(allSolutions, out[0])
	// 	// allSolutions = append(allSolutions, out[1])
	// 	// allSolutions = append(allSolutions, out[2])
	// }
	// sort.Sort(iislice(allSolutions))
	// for _, sol := range allSolutions {
	// 	fmt.Println(sol)
	// }
	// fmt.Println(smallestGroupQE(allSolutions))
	//	fmt.Println(weights)

	//	fmt.Println(singleGroups(singleGroupWeight, []int{}, weights))
	//	fmt.Println(prod(greedyGroups[0]), prod(greedyGroups[1]), prod(greedyGroups[2]))
}

type iislice [][]int

func (s iislice) Len() int { return len(s) }
func (s iislice) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
func (s iislice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sum(ws []int) int {
	out := 0
	for _, w := range ws {
		out += w
	}
	return out
}

func prod(ws []int) int {
	out := 1
	for _, w := range ws {
		out *= w
	}
	return out
}

// 29728298883
// 29728298883 (bleh)
func smallestGroupQE(sol [][]int) int {
	sort.Sort(iislice(sol))
	prods := []int{}
	minlen := len(sol[0])
	for i := 0; i < len(sol); i++ {
		if len(sol[i]) > minlen {
			continue
		}
		prods = append(prods, prod(sol[i]))
	}
	sort.Sort(sort.IntSlice(prods))
	return prods[0]
}

func greedySplit(groupWeight int, fix1, fix2, fix3, allWeights []int) ([3][]int, bool) {
	sgw1 := groupWeight
	sgw2 := groupWeight
	sgw3 := groupWeight
	for _, i := range fix1 {
		sgw1 -= i
	}
	for _, i := range fix2 {
		sgw2 -= i
	}
	for _, i := range fix3 {
		sgw3 -= i
	}
	greedyGroups := [3][]int{fix1, fix2, fix3}
	for _, w := range allWeights {
		if w <= sgw1 {
			sgw1 -= w
			greedyGroups[0] = append(greedyGroups[0], w)
			continue
		}
		if w <= sgw2 {
			sgw2 -= w
			greedyGroups[1] = append(greedyGroups[1], w)
			continue
		}
		sgw3 -= w
		greedyGroups[2] = append(greedyGroups[2], w)
	}
	return greedyGroups, sgw1 == 0 && sgw2 == 0 && sgw3 == 0
}

func singleGroups(singleGroupWeight int, weights, allWeights []int) [][]int {
	fmt.Println(weights, allWeights)
	if singleGroupWeight == 0 {
		return [][]int{weights}
	}
	if singleGroupWeight < 0 {
		return [][]int{}
	}
	out := [][]int{}
	for i, w := range allWeights {
		neww := make([]int, len(weights))
		copy(neww, weights)
		newaw := make([]int, len(allWeights))
		copy(newaw, allWeights)
		newaw = append(newaw[:i], newaw[i+1:]...)
		out = append(out, singleGroups(singleGroupWeight-w, append(neww, w), newaw)...)
	}
	return out
}

func minimizeRecSplit(w int, f, aw []int) [][]int {
	var recSplit func(int, []int, []int) [][]int
	minSoFar := 9999
	minProd := 999999999999
	recSplit = func(singleGroupWeight int, fixed, allWeights []int) [][]int {
		if singleGroupWeight == 0 {
			p := prod(fixed)
			fmt.Println(p, singleGroupWeight, fixed, allWeights)
			if p < minProd {
				minProd = p
			}
			if len(fixed) < minSoFar {
				minSoFar = len(fixed)
			}
			//			fmt.Println(minSoFar, singleGroupWeight, fixed, allWeights)
			return [][]int{fixed}
		}
		if singleGroupWeight < 0 {
			return [][]int{}
		}
		out := [][]int{}
		for i, w := range allWeights {
			if singleGroupWeight-w < 0 {
				continue
			}
			newfixed := make([]int, len(fixed))
			copy(newfixed, fixed)
			newfixed = append(newfixed, allWeights[i])
			newWeights := make([]int, len(allWeights))
			copy(newWeights, allWeights)
			newWeights = append(newWeights[:i], newWeights[i+1:]...)
			if len(newfixed) > minSoFar {
				continue
			}
			if prod(newfixed) >= minProd {
				continue
			}
			out = append(out, recSplit(singleGroupWeight-w, newfixed, newWeights)...)
		}
		return out
	}
	return recSplit(w, f, aw)
}

func recSplit(singleGroupWeight int, fixed, allWeights []int) [][]int {
	if singleGroupWeight == 0 {
		fmt.Println(singleGroupWeight, fixed, allWeights)
		return [][]int{fixed}
	}
	if singleGroupWeight < 0 {
		return [][]int{}
	}
	out := [][]int{}
	for i, w := range allWeights {
		if singleGroupWeight-w < 0 {
			continue
		}
		newfixed := make([]int, len(fixed))
		copy(newfixed, fixed)
		newfixed = append(newfixed, allWeights[i])
		newWeights := make([]int, len(allWeights))
		copy(newWeights, allWeights)
		newWeights = append(newWeights[:i], newWeights[i+1:]...)
		out = append(out, recSplit(singleGroupWeight-w, newfixed, newWeights)...)
	}
	return out
}
