package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	instr := make([]instruction, 0)
	for _, line := range lines {
		instr = append(instr, parse(line))
	}

	for i := 0; ; i++ {
		a := &assem{
			instrs:      instr,
			a:           i,
			clockSignal: []int{},
		}
		if i == 175 {
			for len(a.clockSignal) < 1000 {
				a.run()
			}
			fmt.Println(a.clockSignal)
			os.Exit(0)
		}

		for len(a.clockSignal) < 10 {
			//if a.cur >= 13 && a.cur <= 20 {
			// fmt.Println(a.instrs[a.cur])
			// fmt.Printf("%03d | %03d | %03d | %03d\n", a.a, a.b, a.c, a.d)
			//			}
			//			time.Sleep(1 * time.Millisecond)
			//			fmt.Println(a.a, a.b, a.c, a.d, a.cur, a.instrs[a.cur])
			if a.run() {
				fmt.Println("breaking", i)
				break
			}
		}
		if validClockSignal(a.clockSignal) {
			fmt.Println(i)
		}
	}
}

func validClockSignal(in []int) bool {
	if len(in) == 0 {
		return false
	}
	for i := 0; i < len(in); i++ {
		if i%2 == 0 {
			if in[i] == 0 {
				continue
			}
			return false
		}
		if i%2 != 0 {
			if in[i] == 1 {
				continue
			}
			return false
		}
	}
	return true
}

type assem struct {
	a, b, c, d  int
	instrs      []instruction
	cur         int
	clockSignal []int
}

func parse(line string) instruction {
	words := strings.Split(line, " ")
	switch words[0] {
	case "cpy":
		return copy{words[1], words[2]}
	// case "tgl":
	// 	return tgl{words[1]}
	case "inc":
		return increase{words[1]}
	case "dec":
		return decrease{words[1]}
	case "jnz":
		return jumpNotZero{words[1], words[2]}
	case "out":
		return out{words[1]}
	default:
		panic("bad instruction")
	}
}

type instruction interface {
	run(a *assem)
}

type out struct {
	x string
}

func (o out) run(a *assem) {
	v := a.getValOrNum(o.x)
	a.clockSignal = append(a.clockSignal, v)
	a.cur++
}

func (o out) String() string {
	return fmt.Sprintf("out %v", o.x)
}

type tgl struct {
	x string
}

func (t tgl) String() string {
	return fmt.Sprintf("tgl %v", t.x)
}

func (t tgl) run(a *assem) {
	v := a.getValOrNum(t.x)
	instrNum := a.cur + v
	if instrNum >= len(a.instrs) {
		a.cur++
		return
	}
	instr := a.instrs[instrNum]
	switch i := instr.(type) {
	case jumpNotZero:
		a.instrs[instrNum] = copy{i.x, i.y}
	case increase:
		a.instrs[instrNum] = decrease{i.x}
	case decrease:
		a.instrs[instrNum] = increase{i.x}
	case tgl:
		a.instrs[instrNum] = increase{i.x}
	case copy:
		a.instrs[instrNum] = jumpNotZero{i.x, i.y}
	default:
		panic(fmt.Sprintf("oh no, bad type %T", instr))
	}
	a.cur++
}

type jumpNotZero struct {
	x, y string
}

func (j jumpNotZero) run(a *assem) {
	v := a.getValOrNum(j.x)
	if v == 0 {
		a.cur++
		return
	}
	h := a.getValOrNum(j.y)
	a.cur += h
}

func (j jumpNotZero) String() string {
	return fmt.Sprintf("jnz %v %v", j.x, j.y)
}

type increase struct {
	x string
}

func (i increase) run(a *assem) {
	a.setReg(i.x, a.getVal(i.x)+1)
	a.cur++
}
func (i increase) String() string {
	return fmt.Sprintf("inc %v ", i.x)
}

type decrease struct {
	x string
}

func (d decrease) run(a *assem) {
	a.setReg(d.x, a.getVal(d.x)-1)
	a.cur++
}
func (d decrease) String() string {
	return fmt.Sprintf("dec %v ", d.x)
}

type copy struct {
	x, y string
}

func (c copy) run(a *assem) {
	a.setReg(c.y, a.getValOrNum(c.x))
	a.cur++
}
func (c copy) String() string {
	return fmt.Sprintf("cpy %v %v", c.x, c.y)
}

func (a *assem) run() bool {
	instr := a.instrs[a.cur]
	instr.run(a)
	return a.cur >= len(a.instrs)
}

func (a *assem) getValOrNum(v string) int {
	if v[0] >= 'a' && v[0] <= 'd' {
		return a.getVal(v)
	}
	n, _ := strconv.Atoi(v)
	return n
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
