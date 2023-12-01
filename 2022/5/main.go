package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	state := newState()
	input := internal.ReadRealRawInput()
	init := true
	for _, line := range input {
		if init {
			state.readInitialState(line)
			if line == "" {
				state.initialize()
				init = false
			}
			continue
		}
		if line == "" {
			break
		}
		state.do(parseInstruction(line))
	}
	for i := 1; i <= len(state.containers); i++ {
		fmt.Print(state.containers[i].Peek())
	}
	fmt.Println()
}

type instruction struct {
	num, from, to int
}

func parseInstruction(line string) instruction {
	parts := strings.Split(line, " ")
	i := instruction{}
	i.num, _ = strconv.Atoi(parts[1])
	i.from, _ = strconv.Atoi(parts[3])
	i.to, _ = strconv.Atoi(parts[5])
	return i
}

type state struct {
	initial    map[int][]string
	containers map[int]*internal.Stack[string]
}

func (s *state) initialize() {
	for idx := range s.initial {
		s.containers[idx] = internal.NewStack[string]()
	}
	for idx, crates := range s.initial {
		for i := len(crates) - 1; i >= 0; i-- {
			s.containers[idx].Push(crates[i])
		}
	}
}

func (s *state) do(inst instruction) {
	collect := []string{}
	for i := 0; i < inst.num; i++ {
		val := s.containers[inst.from].Pop()
		collect = append(collect, val)
		//		s.containers[inst.to].Push(val)
	}
	for i := inst.num - 1; i >= 0; i-- {
		s.containers[inst.to].Push(collect[i])
	}
}

func newState() *state {
	return &state{
		initial:    make(map[int][]string),
		containers: make(map[int]*internal.Stack[string]),
	}
}

// readInitialState reads in a line at a time (ends with the first blank line)
func (s *state) readInitialState(line string) {
	containerIndex := 1
	// read 3 chars, (maybe) sep repeat
	for i := 0; i < len(line); {
		// if we're past the first one, read in a separator
		if i > 0 {
			i += 1
		}
		cur := line[i : i+3]
		i += 3
		if cur[0] == '[' {
			box := string(cur[1])
			s.initial[containerIndex] = append(s.initial[containerIndex], box)
		}
		containerIndex++
		// read in 3
		// if it starts with a ' ', it will be empty, read in 2 more charactesr and then read in a separator
		// if it starts with a '[', read in the next number and a ']'
	}
}
