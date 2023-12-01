package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const (
	register = "x"
)

func main() {
	input := internal.ReadInput()

	inst := make([]internal.Instruction, 0)
	for _, line := range input {
		if strings.HasPrefix(line, "noop") {
			inst = append(inst, internal.Noop{})
		}
		if strings.HasPrefix(line, "addx") {
			parts := strings.Split(line, " ")
			x, _ := strconv.Atoi(parts[1])
			inst = append(inst, internal.Noop{})
			inst = append(inst, internal.AddX{X: register, Val: x})
		}
	}

	assm := internal.NewAssem(nil)
	assm.SetReg(register, 1)
	run(assm, inst)
}

func run(a *internal.Assem, inst []internal.Instruction) {
	if len(inst) < 220 {
		panic("instruction list is too short")
	}
	screen := newScreen(40)
	for i := 0; i < len(inst); i++ {
		screen.drawPixel(sprite(a.GetVal(register)))
		instruction := inst[i]
		instruction.Run(a)
	}
	fmt.Println(screen)
}

func signal(in map[int]int) int {
	sum := 0
	for k, v := range in {
		sum += k * v
	}
	return sum
}

type screen struct {
	width int
	data  string
	ptr   int
}

func newScreen(w int) *screen {
	return &screen{
		width: w,
		data:  "",
	}
}

func (s *screen) drawPixel(sprite sprite) {
	switch sprite.contains(s.ptr % 40) {
	case true:
		s.data += string(internal.Yellow("#"))
	case false:
		s.data += "."
	}
	s.ptr++
	if s.ptr%s.width == 0 {
		s.data += "\n"
	}
}

func (s *screen) String() string {
	return s.data
}

type sprite int

func (s sprite) contains(x int) bool {
	return int(s-1) == x || int(s) == x || int(s+1) == x
}
