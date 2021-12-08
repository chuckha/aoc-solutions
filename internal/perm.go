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
