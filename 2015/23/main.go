package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	instructions := make(map[int]instruction)
	registers := map[string]uint{"a": 1, "b": 0}
	for i, line := range lines {
		parts := strings.Split(line, " ")
		offset := 0
		if len(parts) == 3 {
			offset, _ = strconv.Atoi(parts[2])
		}
		if strings.HasPrefix(line, "jmp") {
			offset, _ = strconv.Atoi(parts[1])
			instructions[i] = instruction{
				kind:     parts[0],
				index:    i,
				offset:   offset,
				register: "",
			}
			continue
		}
		instructions[i] = instruction{
			kind:     parts[0],
			index:    i,
			register: strings.TrimSuffix(parts[1], ","),
			offset:   offset,
		}
	}
	//	fmt.Println(instructions)
	for i := 0; i < len(instructions); i++ {
		registerVal := registers[instructions[i].register]
		switch instructions[i].kind {
		case "hlf":
			registerVal = registerVal / 2
		case "tpl":
			registerVal = registerVal * 3
		case "inc":
			registerVal = registerVal + 1
		case "jmp":
			i = i + instructions[i].offset - 1
			continue
		case "jie":
			if registerVal%2 == 0 {
				i = i + instructions[i].offset - 1
				continue
			}
		case "jio":
			if registerVal == 1 {
				i = i + instructions[i].offset - 1
				continue
			}
		default:
			panic("oh no")
		}
		registers[instructions[i].register] = registerVal
	}
	fmt.Println(registers["b"])
}

type machine struct{}

type instruction struct {
	kind     string
	index    int
	register string
	offset   int
}

func (i instruction) run(regVal int) {
	switch i.kind {
	}
}
