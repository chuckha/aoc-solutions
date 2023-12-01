package strings

func ForEachChar(s string, f func(rune)) {
	for _, r := range s {
		f(r)
	}
}
