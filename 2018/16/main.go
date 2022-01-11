package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

// 294 too low
// 310 too low
// 682 not the right answer
func main() {
	lines := internal.ReadRawInput()
	codes := []inputCode{}
	testCode := [][4]int{}
	start := 0
	for i := 0; i < len(lines)-2; i += 3 {
		if !strings.HasPrefix(lines[i], "Before:") {
			start = i
			break
		}
		b4 := readBeforeLine(lines[i])
		code := readNumberLine(lines[i+1])
		after := readAfterLine(lines[i+2])
		codes = append(codes, inputCode{b4, code, after})
	}
	for i := start; i < len(lines); i++ {
		nums := strings.Split(lines[i], " ")
		add := [4]int{}
		for j := 0; j < 4; j++ {
			add[j], _ = strconv.Atoi(nums[j])
		}
		testCode = append(testCode, add)
	}
	// for _, c := range codes {
	// 	fmt.Println(c)
	// }
	solver := &solver{
		possibilities: make(map[int]*possibility),
	}
	for _, c := range codes {
		a := &assem{}
		poss := a.gamut(c.before, c.in, c.after)
		solver.update(c.in[0], poss)
	}
	// for k, v := range solver.possibilities {
	// 	fmt.Println(k)
	// 	fmt.Println(v)
	// }
	a := &assem{}
	for _, c := range testCode {
		a.getCode(solver, c[0])(c)
	}
	fmt.Println(a.registers[0])
}

type inputCode struct {
	before, in, after [4]int
}

func readBeforeLine(line string) [4]int {
	beforew := strings.Split(line, ": ")
	nums := strings.Trim(beforew[1], "[]")
	items := strings.Split(nums, ", ")
	out := [4]int{}
	for i := 0; i < 4; i++ {
		out[i], _ = strconv.Atoi(items[i])
	}
	return out
}
func readAfterLine(line string) [4]int {
	beforew := strings.Split(line, ":  ")
	nums := strings.Trim(beforew[1], "[]")
	items := strings.Split(nums, ", ")
	out := [4]int{}
	for i := 0; i < 4; i++ {
		out[i], _ = strconv.Atoi(items[i])
	}
	return out
}

func readNumberLine(line string) [4]int {
	out := [4]int{}
	nums := strings.Split(line, " ")
	for i := 0; i < 4; i++ {
		out[i], _ = strconv.Atoi(nums[i])
	}
	return out
}

func opcodeRunner(input []int) {}

type assem struct {
	registers [4]int
}

func (a *assem) set(in [4]int) {
	a.registers = in
}

func (a *assem) getCode(s *solver, i int) opcode {
	name := extractKey(s.possibilities[i].can)
	switch name {
	case "addi":
		return a.addi
	case "addr":
		return a.addr
	case "muli":
		return a.muli
	case "mulr":
		return a.mulr
	case "bani":
		return a.bani
	case "banr":
		return a.banr
	case "bori":
		return a.bori
	case "borr":
		return a.borr
	case "eqir":
		return a.eqir
	case "eqri":
		return a.eqri
	case "eqrr":
		return a.eqrr
	case "gtir":
		return a.gtir
	case "gtri":
		return a.gtri
	case "gtrr":
		return a.gtrr
	case "setr":
		return a.setr
	case "seti":
		return a.seti
	default:
		panic("unknown instruction")
	}
}

type opcode func([4]int)

func (a *assem) gamut(before, in, expected [4]int) *possibility {
	//	fmt.Println(before, in, expected)
	opcodes := map[string]opcode{
		"addi": a.addi, "addr": a.addr,
		"muli": a.muli, "mulr": a.mulr,
		"bani": a.bani, "banr": a.banr,
		"bori": a.bori, "borr": a.borr,
		"eqir": a.eqir, "eqri": a.eqri, "eqrr": a.eqrr,
		"gtir": a.gtir, "gtri": a.gtri, "gtrr": a.gtrr,
		"setr": a.setr, "seti": a.seti,
	}
	out := newPossibility()
	for k := range opcodes {
		out.can[k] = struct{}{}
	}
	for k, c := range opcodes {
		a.reset()
		// if k == "eqrr" {
		// 	fmt.Println("----", k, "----")
		// 	fmt.Println(in, "<- in")
		// 	fmt.Println(before, "<- before")
		// }
		a.set(before)
		c(in)
		// if k == "eqrr" {
		// 	fmt.Println(a.registers, "<- actual")
		// 	fmt.Println(expected, "<- expected")
		// 	//			fmt.Println("checking", expected, a.registers)
		// }
		if !a.checkOutput(expected) {
			out.notPossible(k)
		}
	}
	//	fmt.Println(before, in, expected)
	//	fmt.Println(out)
	return out
}

func (a *assem) reset() {
	a.registers = [4]int{}
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

func (a *assem) checkOutput(expected [4]int) bool {
	return expected[0] == a.registers[0] &&
		expected[1] == a.registers[1] &&
		expected[2] == a.registers[2] &&
		expected[3] == a.registers[3]
}

// run before on every opcode
// compare to after
// if more than 3 add 1

type solver struct {
	possibilities map[int]*possibility
}
type possibility struct {
	can, cannot map[string]struct{}
}

func newPossibility() *possibility {
	return &possibility{
		can:    make(map[string]struct{}),
		cannot: make(map[string]struct{}),
	}
}
func (p *possibility) notPossible(item string) {
	delete(p.can, item)
	p.cannot[item] = struct{}{}
}
func (p *possibility) String() string {
	var out strings.Builder
	out.WriteString("can:\n")
	for k := range p.can {
		out.WriteString("\t" + k)
	}
	out.WriteString("\n")
	out.WriteString("cannot:\n")
	for k := range p.cannot {
		out.WriteString("\t" + k)
	}
	return out.String()
}

// if a number can be one, add it to the set
func (s *solver) add(opcodeNum int, opcodeName string) {
	if _, ok := s.possibilities[opcodeNum]; !ok {
		s.possibilities[opcodeNum] = newPossibility()
	}
	if _, ok := s.possibilities[opcodeNum].cannot[opcodeName]; ok {
		//		fmt.Printf("opcode (%d) %q is already proven to not be this\n", opcodeNum, opcodeName)
		return // don't do anything if it's already disproven
	}
	s.possibilities[opcodeNum].can[opcodeName] = struct{}{}
}
func (s *solver) remove(opcodeNum int, opcodeName string) {
	delete(s.possibilities[opcodeNum].can, opcodeName)
	s.possibilities[opcodeNum].cannot[opcodeName] = struct{}{}
}
func (s *solver) update(code int, p *possibility) {
	for can := range p.can {
		s.add(code, can)
	}
	for cannot := range p.cannot {
		s.remove(code, cannot)
	}
	if len(s.possibilities[code].can) == 1 {
		key := extractKey(s.possibilities[code].can)
		for k, v := range s.possibilities {
			if k == code {
				continue
			}
			v.notPossible(key)
		}
	}
}

func extractKey(in map[string]struct{}) string {
	if len(in) != 1 {
		panic("did real bad")
	}
	for k := range in {
		return k
	}
	panic("you did bad")
}
