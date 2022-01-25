package main

import (
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	sum := 0
	for _, line := range lines {
		x, _ := strconv.Atoi(line)
		sum += fuelSum(x)
	}
	fmt.Println(sum)
}

func fuelSum(start int) int {
	out := 0
	for start >= 0 {
		fuel := (start / 3) - 2
		if fuel <= 0 {
			return out
		}
		out += fuel
		start = fuel
	}
	return out
}
