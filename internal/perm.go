package internal

func AllPermutations[T any](fixed, unfixed []T) [][]T {
	if len(unfixed) == 0 {
		swap := make([]T, len(fixed))
		copy(swap, fixed)
		return [][]T{swap}
	}

	out := [][]T{}
	for i := range unfixed {
		swap := make([]T, len(unfixed))
		copy(swap, unfixed)
		swap[0], swap[i] = swap[i], swap[0]
		out = append(out, AllPermutations(append(fixed, swap[0]), swap[1:])...)
	}
	return out
}

func NonRecursivePermutations(data []string, size int) [][]string {
	ptrs := make([]int, size)
	for i := 0; i < size; i++ {
		ptrs[i] = i
	}
	out := make([][]string, 0)
	collect := []string{}
	for _, ptr := range ptrs {
		collect = append(collect, data[ptr])
	}
	out = append(out, collect)

	for !Inc(ptrs, len(data)-1) {
		collect := []string{}
		for _, ptr := range ptrs {
			collect = append(collect, data[ptr])
		}
		out = append(out, collect)
	}
	return out
}

func Inc(ptrs []int, max int) bool {
	// find the first point that is the max
	var first int
	found := false
	for i := len(ptrs) - 1; i >= 0; i-- {
		if ptrs[i] == max {
			continue
		}
		first = i
		found = true
		break
	}
	// if found is false we can inc the last one without issue
	if !found {
		return true
	}
	ptrs[first]++
	for i := first + 1; i < len(ptrs); i++ {
		ptrs[i] = 0
	}
	return false
}

func AllCombinations(fixed []string, unfixed []string, size int) [][]string {
	if size == 0 {
		return [][]string{fixed}
	}
	out := [][]string{}
	for i := range unfixed[:len(unfixed)-size+1] {
		ff := make([]string, len(fixed))
		copy(ff, fixed)
		uf := make([]string, len(unfixed))
		copy(uf, unfixed)
		ff = append(ff, uf[i])
		uf = uf[i+1:]
		out = append(out, AllCombinations(ff, uf, size-1)...)
	}
	return out
}
