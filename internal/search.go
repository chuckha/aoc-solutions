package internal

func Search[T comparable](needle T, haystack []T) int {
	for i := range haystack {
		if needle == haystack[i] {
			return i
		}
	}
	return -1
}
