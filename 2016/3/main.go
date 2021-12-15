package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	count := 0
	cols := [3][]int{}
	for _, line := range lines {
		sides := strings.Fields(line)
		a, _ := strconv.Atoi(sides[0])
		cols[0] = append(cols[0], a)
		b, _ := strconv.Atoi(sides[1])
		cols[1] = append(cols[1], b)
		c, _ := strconv.Atoi(sides[2])
		cols[2] = append(cols[2], c)
	}
	for _, col := range cols {
		for i := 0; i < len(col); i += 3 {
			sides := col[i : i+3]
			fmt.Println(sides[0], sides[1], sides[2])
			if isValidTriangle(sides[0], sides[1], sides[2]) {
				count += 1
			}
		}
	}
	fmt.Println(count)
}

func isValidTriangle(a, b, c int) bool {
	return a+b > c && a+c > b && b+c > a
}
