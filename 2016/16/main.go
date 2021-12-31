package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	state := internal.ReadInput()[0]
	maxSize := 272
	maxSize = 35651584 // part 2
	for len(state) < maxSize {
		state = dragon(state)
	}

	fmt.Println(checksum(state[:maxSize]))
}

func dragon(in string) string {
	var b strings.Builder
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] == '0' {
			b.WriteString("1")
			continue
		}
		b.WriteString("0")
	}
	return in + "0" + b.String()
}

func checksum(in string) string {
	for len(in)%2 == 0 {
		in = pair(in)
	}
	return in
}

func pair(in string) string {
	var b strings.Builder
	for i := 0; i < len(in)-1; i += 2 {
		if in[i] == in[i+1] {
			b.WriteString("1")
			continue
		}
		b.WriteString("0")
	}
	return b.String()
}
