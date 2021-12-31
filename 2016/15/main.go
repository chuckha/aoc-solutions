package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	// Disc #1 has 17 positions; at time=0, it is at position 15.
	// Disc #2 has 3 positions; at time=0, it is at position 2.
	// Disc #3 has 19 positions; at time=0, it is at position 4.
	// Disc #4 has 13 positions; at time=0, it is at position 2.
	// Disc #5 has 7 positions; at time=0, it is at position 2.
	// Disc #6 has 5 positions; at time=0, it is at position 0.
	lines := internal.ReadInput()

	discs := []*disc{}
	for _, line := range lines {
		discs = append(discs, parseInput(line))
	}
	discs = append(discs, parseInput("Disc #7 has 11 positions; at time=0, it is at position 0."))
	for i := 0; ; i++ {
		// check win condition
		all := true
		for t, d := range discs {
			all = all && d.positionAtTime(i+1+t) == 0
			if !all {
				break
			}
		}
		if all {
			fmt.Println(i)
			break
		}
	}
}

type disc struct {
	num       int
	positions int
	startpos  int
}

func (d *disc) positionAtTime(i int) int {
	return (d.startpos + i) % d.positions
}

func parseInput(line string) *disc {
	words := strings.Split(line, " ")
	num := strings.TrimPrefix(words[1], "#")
	n, _ := strconv.Atoi(num)
	maxpos, _ := strconv.Atoi(words[3])
	startpos, _ := strconv.Atoi(strings.TrimSuffix(words[11], "."))
	return &disc{
		num:       n,
		positions: maxpos,
		startpos:  startpos,
	}
}
