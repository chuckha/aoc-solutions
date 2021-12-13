package main

import (
	"fmt"
	"sort"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	s := internal.NewStack[string]()
	sum := 0
	validLines := []string{}
	for _, line := range lines {
		isValid := true
		for _, c := range line {
			switch c {
			case '{', '[', '(', '<':
				s.Push(string(c))
			case '}', ']', ')', '>':
				popped := s.Pop()
				if out := match(popped, string(c)); out != "" {
					sum += score(out)
					isValid = false
				}
			default:
				panic("unknown character")
			}
		}
		if isValid {
			validLines = append(validLines, line)
		}
	}

	scores := make([]int, len(validLines))
	for i, line := range validLines {
		s2 := internal.NewStack[string]()
		for _, c := range line {
			switch c {
			case '{', '[', '(', '<':
				s2.Push(string(c))
			case '}', ']', ')', '>':
				s2.Pop()
			default:
				panic("unknown character")
			}
		}
		scores[i] = closingScore(s2)
	}
	sort.Sort(sort.IntSlice(scores))
	fmt.Println(scores[(len(scores)-1)/2])
}

func closingScore(s *internal.Stack[string]) int {
	score := 0
	for !s.Empty() {
		score = score * 5
		switch s.Pop() {
		case "(":
			score += 1
		case "[":
			score += 2
		case "{":
			score += 3
		case "<":
			score += 4
		default:
			panic("unknown char")
		}
	}
	return score
}

func match(open, close string) string {
	switch open {
	case "(":
		if close != ")" {
			return close
		}
	case "[":
		if close != "]" {
			return close
		}
	case "{":
		if close != "}" {
			return close
		}
	case "<":
		if close != ">" {
			return close
		}
	}
	return ""
}

func score(in string) int {
	switch in {
	case ")":
		return 3
	case "]":
		return 57
	case "}":
		return 1197
	case ">":
		return 25137
	}
	panic("oh shiii")
}
