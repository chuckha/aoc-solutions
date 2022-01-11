package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	code()
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
	fmt.Println(comp.registers[0])
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

// otherwise known as sum
func code() {
	reg4 := 10551326
	end := int(math.Sqrt(10551326)) + 1
	x := reg4 + 1
	for i := 2; i < end; i++ {
		for reg4%i == 0 && reg4 > 0 {
			reg4 = reg4 / i
			x += i + reg4
			fmt.Println(reg4, i)
		}
	}

	fmt.Println(x)
}

//5275665　too low

// 	reg0, reg1, reg2, reg4, reg5 := 0, 0, 0, 0, 0
// 	reg4 = 926
// 	if reg0 == 1 {
// 		reg4 = 10551326
// 		reg0 = 0
// 	}
// 	x := 0
// 	for i := 0 ; i < reg4; i++ {
// 		for j := 0; j <= reg4; j++ {
// 			if i*j==reg4 {
// 				x+=i
// 			}
// 		}
// 	}
// 	reg2 = 1
// 	for {
// 		reg5 = 1 // needs to get up to reg4 (926 or 10551326)
// 		for {
// 			reg1 = reg2 * reg5
// 			if reg1 == reg4 {
// 				reg0 = reg0 + reg2 // this is the final value to inspect
// 			}
// 			reg5 = reg5 + 1
// 			if reg5 <= reg4 {
// 				continue
// 			}
// 			reg2 = reg2 + 1
// 			if reg2 > reg4 {
// 				return
// 			}
// 			break
// 		}
// 	}
// }

/* can always replace reg3 with immediate value of instr number
#ip 3
0  addi 3 16 3  <- reg 3 + 16 -> 3 (skips to line 17 (16 + 1))
1  seti 1 3 2   <- 1 -> 2 (start of a loop)
2  seti 1 0 5   <- 1 -> 5 (start of an inner loop)
3  mulr 2 5 1   <- reg2 * reg 5 -> 1
4  eqrr 1 4 1   <- if reg 1 == reg 4 ? 1 -> 1 : 0 -> 1 (if reg 1 == reg 4; do not skip instr 5; else skip instr 5)
5  addr 1 3 3   <- reg 1 + reg 3 -> 3 (skips reg1 instructions, if reg 1 == 0; skips no instructions)
6  addi 3 1 3   <- reg 3 + 1 -> 3 (skips one instruction)
7  addr 2 0 0   <- reg 2 + reg 0 -> 0
8  addi 5 1 5   <- reg 5 + 1 -> 5
9  gtrr 5 4 1   <- if reg 5 > reg 4 ? 1 -> 1 : 0 -> 1
10 addr 3 1 3   <- reg 3 + reg 1 -> 3
11 seti 2 2 3   <- 2 -> 3 (loops back to instruction 2+1 (inner loop))
12 addi 2 1 2   <- reg 2 + 1 -> 2
13 gtrr 2 4 1   <- if reg 2 > reg 4 ? 1 -> 1 : 0 -> 1
14 addr 1 3 3   <- reg 1 + reg 3 -> 3 (skips reg1 instructions)
15 seti 1 1 3   <- 1 -> 3 (loops back to line 1+1; inner loop)
16 mulr 3 3 3   <- reg 3 * reg 3 -> 3 (jumps 3*3 instructions forward (presumably ending the program, reg3 must always be 16))
17 addi 4 2 4   <- reg 4 + 2 -> 4
18 mulr 4 4 4   <- reg 4 * reg 4 -> 4
19 mulr 3 4 4   <- reg 3 * reg 4 -> 4
20 muli 4 11 4  <- reg 4 * 11 -> 4
21 addi 1 4 1   <- reg 1 + 4 -> 1
22 mulr 1 3 1   <- reg 1 * reg 3 -> 1
23 addi 1 2 1   <- reg 1 + 2 -> 1
24 addr 4 1 4   <- reg 4 + reg 1 -> 4
25 addr 3 0 3   <- if reg 0 > 0; skip line 26
26 seti 0 2 3   <- loops back to the start (+1)
27 setr 3 6 1   <- starts building a big number
28 mulr 1 3 1       <- reg 1 * reg 3 -> 1
29 addr 3 1 1       <- reg 3 + reg 1 -> 1
30 mulr 3 1 1       <- reg 3 * reg 1 -> 1
31 muli 1 14 1      <- reg 1 * 14 -> 1
32 mulr 1 3 1       <- reg 3 * reg 1 -> 1
33 addr 4 1 4       <- reg 4 + reg 1 -> 4
34 seti 0 6 0   <- set 0 reg to 0 so the next time line 25 hits, skip building the big number
35 seti 0 9 3   <- loops back to the start of the program
*/
