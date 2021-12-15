package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < len(lines[0]); i++ {
		counts := make([]int, 26)
		for _, line := range lines {
			counts[strings.Index(letters, string(line[i]))] += 1
		}
		max := 0
		for i, c := range counts {
			if c < counts[max] {
				max = i
			}
		}
		fmt.Print(string(letters[max]))
	}
	fmt.Println()
}
