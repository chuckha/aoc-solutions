package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"golang.org/x/exp/constraints"
)

// A, X = rock
// B, Y = paper
// C, Z = scissors

var selectPoints = map[string]int{"R": 1, "P": 2, "S": 3}

func main() {
	sum := 0
	for _, line := range internal.ReadInput() {

		parts := strings.Split(line, " ")
		ptgo, pts := pickToGetOutcome(parts[0], parts[1])
		sum += pts + selectPoints[ptgo]
	}
	fmt.Println(sum)
}

// X = lose
// Y = draw
// Z = win

type result struct {
	move string
	pts  int
}

const (
	lose = 0
	draw = 3
	win  = 6
)

func pickToGetOutcome(them, me string) (string, int) {
	lookup := map[string]map[string]result{
		"A": {
			"X": {"S", lose},
			"Y": {"R", draw},
			"Z": {"P", win},
		},
		"B": {
			"X": {"R", lose},
			"Y": {"P", draw},
			"Z": {"S", win},
		},
		"C": {
			"X": {"P", lose},
			"Y": {"S", draw},
			"Z": {"R", win},
		},
	}
	r, ok := lookup[them][me]
	if !ok {
		panic("bad lookup")
	}
	return r.move, r.pts
}

func winPoints(them, me string) int {
	switch them {
	case "A":
		switch me {
		case "X":
			return 3
		case "Y":
			return 6
		case "Z":
			return 0
		}
	case "B":
		switch me {
		case "X":
			return 0
		case "Y":
			return 3
		case "Z":
			return 6
		}
	case "C":
		switch me {
		case "X":
			return 6
		case "Y":
			return 0
		case "Z":
			return 3
		}
	}
	panic("do not get here")
}

type Node[T constraints.Ordered] struct {
	Value T
	Left  *Node[T]
	Right *Node[T]
}

func (n *Node[T]) Insert(item T) {
	if item <= n.Value {
		if n.Left == nil {
			n.Left = &Node[T]{Value: item}
			return
		}
		n.Left.Insert(item)
		return
	}
	if n.Right == nil {
		n.Right = &Node[T]{Value: item}
		return
	}
	n.Right.Insert(item)
}
