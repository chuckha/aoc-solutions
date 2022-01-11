package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	code3(1)
	os.Exit(0)
	lines := internal.ReadInput()
	instPtrIdx, _ := strconv.Atoi(strings.Split(lines[0], " ")[1])
	instructions := []instruction{}
	for _, line := range lines[1:] {
		words := strings.Split(line, " ")
		one, _ := strconv.Atoi(words[1])
		two, _ := strconv.Atoi(words[2])
		thr, _ := strconv.Atoi(words[3])
		instructions = append(instructions, instruction{
			inst:  words[0],
			input: [4]int{-1, one, two, thr},
		})
	}
	comp := &assem{
		instPtrIdx:   0,
		registers:    [6]int{},
		instReg:      instPtrIdx,
		instructions: instructions,
	}
	//	comp.registers[0] = 1
	for !comp.run() {
	}
	fmt.Println(comp.registers[0], comp.numExecuted)
}

type instruction struct {
	inst  string
	input [4]int
}

type assem struct {
	registers    [6]int
	instReg      int
	instPtrIdx   int
	instructions []instruction
	numExecuted  int
}

func (a *assem) set(in [6]int) {
	a.registers = in
}

func (a *assem) reset() {
	a.registers = [6]int{}
}
func (a *assem) run() bool {
	a.instPtrIdx = a.getReg(a.instReg)
	if a.instPtrIdx >= len(a.instructions) {
		return true
	}
	fmt.Println("before", a.registers)
	fmt.Println(a.instructions[a.instPtrIdx])
	a.runInstruction(a.currentInstruction())
	a.registers[a.instReg] += 1
	fmt.Println("after", a.registers)
	return a.instPtrIdx >= len(a.instructions)
}

func (a *assem) currentInstruction() instruction {
	return a.instructions[a.instPtrIdx]
}

func (a *assem) runInstruction(inst instruction) {
	switch inst.inst {
	case "addi":
		a.addi(inst.input)
	case "addr":
		a.addr(inst.input)
	case "muli":
		a.muli(inst.input)
	case "mulr":
		a.mulr(inst.input)
	case "bani":
		a.bani(inst.input)
	case "banr":
		a.banr(inst.input)
	case "bori":
		a.bori(inst.input)
	case "borr":
		a.borr(inst.input)
	case "eqir":
		a.eqir(inst.input)
	case "eqri":
		a.eqri(inst.input)
	case "eqrr":
		a.eqrr(inst.input)
	case "gtir":
		a.gtir(inst.input)
	case "gtri":
		a.gtri(inst.input)
	case "gtrr":
		a.gtrr(inst.input)
	case "setr":
		a.setr(inst.input)
	case "seti":
		a.seti(inst.input)
	default:
		panic("unknown instruction")
	}
	a.numExecuted++
}
func (a *assem) getReg(i int) int {
	return a.registers[i]
}
func (a *assem) setReg(i, v int) {
	//	fmt.Println("setting register", v, "to", i, "updating from", a.registers)
	a.registers[v] = i
	//	fmt.Println("after setting register", a.registers)
}
func (a *assem) addr(input [4]int) {
	a.setReg(a.getReg(input[1])+a.getReg(input[2]), input[3])
}
func (a *assem) addi(input [4]int) {
	a.setReg(a.getReg(input[1])+input[2], input[3])
}
func (a *assem) mulr(input [4]int) {
	a.setReg(a.getReg(input[1])*a.getReg(input[2]), input[3])
}
func (a *assem) muli(input [4]int) {
	a.setReg(a.getReg(input[1])*input[2], input[3])
}
func (a *assem) banr(input [4]int) {
	a.setReg(a.getReg(input[1])&a.getReg(input[2]), input[3])
}
func (a *assem) bani(input [4]int) {
	a.setReg(a.getReg(input[1])&input[2], input[3])
}
func (a *assem) borr(input [4]int) {
	a.setReg(a.getReg(input[1])|a.getReg(input[2]), input[3])
}
func (a *assem) bori(input [4]int) {
	a.setReg(a.getReg(input[1])|input[2], input[3])
}
func (a *assem) setr(input [4]int) {
	a.setReg(a.getReg(input[1]), input[3])
}
func (a *assem) seti(input [4]int) {
	a.setReg(input[1], input[3])
}
func (a *assem) gtir(input [4]int) {
	if input[1] > a.getReg(input[2]) {
		a.setReg(1, input[3])
		return
	}
	a.setReg(0, input[3])
}
func (a *assem) gtri(input [4]int) {
	if a.getReg(input[1]) > input[2] {
		a.setReg(1, input[3])
		return
	}
	a.setReg(0, input[3])
}
func (a *assem) gtrr(input [4]int) {
	if a.getReg(input[1]) > a.getReg(input[2]) {
		a.setReg(1, input[3])
		return
	}
	a.setReg(0, input[3])
}
func (a *assem) eqir(input [4]int) {
	if input[1] == a.getReg(input[2]) {
		a.setReg(1, input[3])
		return
	}
	a.setReg(0, input[3])
}
func (a *assem) eqri(input [4]int) {
	if a.getReg(input[1]) == input[2] {
		a.setReg(1, input[3])
		return
	}
	a.setReg(0, input[3])
}
func (a *assem) eqrr(input [4]int) {
	if a.getReg(input[1]) == a.getReg(input[2]) {
		a.setReg(1, input[3])
		return
	}
	a.setReg(0, input[3])
}

// func code() {
// 	reg0, reg1, reg2, reg3, reg4, reg5 := 0, 0, 0, 0, 0, 0
// 	for {
// 		if 123&456 == 72 {
// 			break
// 		}
// 	}
// 	for {
// 		reg4 = 65536
// 		reg2 = 6718165
// 		for {
// 			reg3 = reg4 & 255
// 			for {
// 				reg2 = reg2 + reg3
// 				reg2 = reg2 & 16777215
// 				reg2 = reg2 * 65899
// 				reg2 = reg2 & 16777215
// 				if 256 > reg4 {
// 					if reg2 == reg0 {
// 						return
// 					}
// 				}
// 			}
// 			if reg2 == reg0 {
// 				return
// 			}
// 			reg4 = 256
// 		}
// 	}
// }

func code3(reg0 int) {
	test := map[int]struct{}{}
	reg2 := 0
	reg4 := 0
	reg3 := 0
	reg1 := 0
	for {
		if 123&456 == 72 {
			break
		}
	}
	reg2 = 0
	for {
		reg4 = reg2 | 65536
		reg2 = 6718165
		for {
			reg3 = reg4 & 255
			reg2 = reg2 + reg3
			reg2 = reg2 & 16777215
			reg2 = reg2 * 65899
			reg2 = reg2 & 16777215
			if 256 > reg4 {
				break
			}
			reg3 = 0
			for {
				reg1 = reg3 + 1
				reg1 = reg1 * 256
				if reg1 > reg4 {
					break
				}
				reg3 = reg3 + 1
			}
			reg4 = reg3
		}
		fmt.Println("part1", reg2)
		l1 := len(test)
		test[reg2] = struct{}{}
		if len(test) == l1 {
			fmt.Println("part2 is the number above this number", reg2)
			os.Exit(0)
		}
		if reg2 == reg0 {
			break
		}
	}
}

// 10837056 too high

/*
#ip 5           <- [0,0,0,0,0,0]
0  seti 123 0 2         <- sets 123 to register 2 [0,0,123,0,0,0]
1  bani 2 456 2         <- does 123&456 and stores it in register 2
2  eqri 2 72 2          <- ensures the & got the right value and stores 1 or 0 in 2
3  addr 2 5 5           <- reg5 = 3, if reg2 == 0 (not equal 123&456=72), do not skip one, otherwise continue with execution
4  seti 0 0 5           <-   loop the initial instruction
5  seti 0 5 2           <- set register 2 to 0 to clear out the test
6  bori 2 65536 4       <- 0|65536=65536-> reg4 (start of a loop)
7  seti 6718165 9 2     <- 6718165 -> reg2
8  bani 4 255 3         <- 65536&255 = 0 -> reg3 (START OF LOOP FROM INSTR 28)
9  addr 2 3 2           <- reg2 + reg3 -> reg2
10 bani 2 16777215 2    <- reg2 & 16777215 -> reg2
11 muli 2 65899 2       <- reg2 * 65899 -> reg2
12 bani 2 16777215 2    <- reg2 & 16777215 -> reg2
13 gtir 256 4 3         <- if 256 > reg4 ? 1 -> reg3 : 0 -> reg3
14 addr 3 5 5           <- reg3 + reg5 -> reg5 (skip one or not)
15 addi 5 1 5           <- reg5 + 1 -> reg5 (skip an instruction)
16 seti 27 8 5          <- 27 -> reg5
17 seti 0 4 3           <- 0 -> reg3
18 addi 3 1 1           <- reg3 + 1 -> reg1 (START OF LOOP from inst 25)
19 muli 1 256 1         <- reg1 * 256 -> reg1
20 gtrr 1 4 1           <- reg1 > reg4 ? 1 -> reg1 : 0 -> reg1
21 addr 1 5 5           <- reg1 + reg5 -> reg5 (skip one or not)
22 addi 5 1 5           <- reg5 + 1 -> reg5 (skip one)
23 seti 25 8 5          <- 25 -> reg5 (go to instruction 26)
24 addi 3 1 3           <- reg3 +1 -> reg3
25 seti 17 3 5          <- 17 -> reg5 (go to instruction 18)
26 setr 3 6 4           <- reg3 -> reg4
27 seti 7 9 5           <- 7 -> reg5 (go to instruction 8)
28 eqrr 2 0 3           <- entry point to loop from instruction 16, reg2==reg0 ? 1 -> reg3 : 0 -> reg3
29 addr 3 5 5           <- reg3 + reg5 -> reg5 (skip 1 or not) Exit point of script
30 seti 5 1 5           <- 5 -> reg5
*/
