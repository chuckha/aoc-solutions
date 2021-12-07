package main

import (
	"fmt"
	"strings"
)

func main() {
	cur := "1321131112"
	for i := 0; i < 50; i++ {
		cur = join(split(cur))
	}
	fmt.Println(len(cur))
}

type count struct {
	occurances int
	number     string
}

func (c count) String() string {
	return fmt.Sprintf("%d%s", c.occurances, c.number)
}

func join(counts []count) string {
	var b strings.Builder
	for _, c := range counts {
		b.WriteString(c.String())
	}
	return b.String()
}

func split(in string) []count {
	var cur rune
	i := 0
	out := []count{}
	for j, c := range in {
		if j == len(in)-1 {
			if c == cur {
				i++
				out = append(out, count{i, string(cur)})
				return out
			}
			if cur != 0 {
				out = append(out, count{i, string(cur)})
			}
			out = append(out, count{1, string(c)})
			return out
		}
		if cur == 0 {
			i++
			cur = c
			continue
		}
		if c == cur {
			i++
			continue
		}
		out = append(out, count{i, string(cur)})
		i = 1
		cur = c
	}

	return out
}
