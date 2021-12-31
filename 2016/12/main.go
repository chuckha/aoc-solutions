package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	a := &assem{
		instructions: lines,
		c:            1,
	}

	c := 0
	for !a.run() {
		c++
		//		fmt.Println(a)
	}
	fmt.Println(c)
	fmt.Println(a.a)
}

type assem struct {
	a, b, c, d   int
	instructions []string
	cur          int
}

func (a *assem) run() bool {
	line := a.instructions[a.cur]
	//	fmt.Println(line)
	asm := strings.Split(line, " ")
	switch asm[0] {
	case "cpy":
		var v int
		if asm[1][0] >= 'a' && asm[1][0] <= 'd' {
			v = a.getVal(asm[1])
		} else {
			v, _ = strconv.Atoi(asm[1])
		}
		a.setReg(asm[2], v)
	case "inc":
		a.setReg(asm[1], a.getVal(asm[1])+1)
	case "dec":
		a.setReg(asm[1], a.getVal(asm[1])-1)

	case "jnz":
		var v int
		if asm[1][0] >= 'a' && asm[1][0] <= 'd' {
			v = a.getVal(asm[1])
		} else {
			v, _ = strconv.Atoi(asm[1])
		}
		if v == 0 {
			break
		}
		y, _ := strconv.Atoi(asm[2])
		a.cur += y
		return a.cur >= len(a.instructions)
	}
	a.cur++
	return a.cur >= len(a.instructions)
}

func (a *assem) setReg(reg string, val int) {
	switch reg {
	case "a":
		a.a = val
	case "b":
		a.b = val
	case "c":
		a.c = val
	case "d":
		a.d = val
	default:
		panic(fmt.Sprintf("set reg called with %s %d", reg, val))
	}
}
func (a *assem) getVal(reg string) int {
	switch reg {
	case "a":
		return a.a
	case "b":
		return a.b
	case "c":
		return a.c
	case "d":
		return a.d
	default:
		panic(fmt.Sprintf("get val called with %s", reg))
	}
}
