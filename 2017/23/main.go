package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	fmt.Println(answer())
}
func runActual() {
	lines := internal.ReadInput()
	instructions := []internal.Instruction{}
	for _, line := range lines {
		instructions = append(instructions, parseLine(line))
	}
	a := internal.NewAssem(instructions)
	a.Registers["a"] = 1
	counts := map[string]int{}
	for !a.Run() {
		//		fmt.Println(a.Registers)
		switch v := a.Instructions[a.InstIdx].(type) {
		case internal.JumpNotZero:
			//			fmt.Println(a.Registers)
			counts[fmt.Sprintf("jnz %v %v", v.X, v.Y)]++
		}
	}
	fmt.Println(a.Registers)
	sortPrintCounts(counts)
}

// how many non prime numbers are between 107900 and 124900

func answer() int {
	count := 0
	for i := 107900; i <= 124900; i += 17 {
		if !isPrime(i) {
			count++
		}
	}
	return count
}

func isPrime(n int) bool {
	sqrt := int(math.Sqrt(float64(n)))
	for i := 2; i < sqrt+1; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func test2() {
	bstart := 107900
	binc := 17
	dinc := 1
	einc := 1
	c := 124900
	dstart := 2
	estart := 2
	h := 0
	countJNZ8 := 0
	countJNZ13 := 0
	countJNZ23 := 0

	// 127 is prime
	bstart = 110
	c = 161

	//	sum(107900, c, 17) for each one of these [107900, 107917, ...]
	// 		do d -> each of those sum(107900-2+107917-2+...)
	// 			(2 -> 107900, sum(2 * ))

	// jnz -23
	for b := bstart; b != c; b += binc {
		fmt.Println("b", b)
		if b > c {
			panic("not going to exit")
		}
		f := 1
		countJNZ23++
		// jnz -13
		for d := dstart; d < b; d += dinc {
			countJNZ13++
			// jnz -8
			for e := estart; e < b; e += einc {
				countJNZ8++
				if e*d == b {
					fmt.Println(e, d, b)
					f = 0
				}
			}
		}
		if f == 0 {
			h++ // H is a count of non prime numbers between b and c incremented by 17
		}
	}
	fmt.Println("JNZ8", countJNZ8, "JNZ13", countJNZ13, "JNZ23", countJNZ23, "h", h)
}

func test() {
	bstart := 107900
	c := 124900
	dstart := 2
	estart := 2
	f := 1
	h := 0
	count1, count2, count3 := 0, 0, 0
	for b := bstart; b < c; b += 17 {
		f = 1 // jnz -23 loop
		count1++
		for d := dstart; d < b; d++ {
			estart = 2
			count2++
			for e := estart; e < b; e++ {
				count3++
				if e*d-b == 0 {
					f = 0
				}
			}
		}
		if f == 0 {
			h++
		}
	}
	fmt.Println("count1", count1, "count2", count2, "count3", count3, "h", h)
}

func sortPrintCounts(counts map[string]int) {

	pairs := []pair{}
	for k, v := range counts {
		pairs = append(pairs, pair{k, v})
	}
	sort.Sort(sort.Reverse(pairz(pairs)))
	for _, k := range pairs {
		fmt.Println(k)
	}
}

type pair struct {
	k string
	v int
}
type pairz []pair

func (p pairz) Len() int           { return len(p) }
func (p pairz) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairz) Less(i, j int) bool { return p[i].v < p[j].v }

func parseLine(line string) internal.Instruction {
	words := strings.Split(line, " ")
	switch words[0] {
	// case "snd":
	// 	return Send{words[1]}
	case "set":
		return Set{words[1], words[2]}
	// case "add":
	// 	return Add{words[1], words[2]}
	case "mul":
		return Multiply{words[1], words[2]}
	// case "mod":
	// 	return Modulo{words[1], words[2]}
	// case "rcv":
	// 	return Recover{words[1]}
	// case "jgz":
	// 	return JumpGreaterThanZero{words[1], words[2]}
	case "jnz":
		return internal.JumpNotZero{words[1], words[2]}
	case "sub":
		return Subtract{words[1], words[2]}
	default:
		panic("bad instruction")
	}
}

type Send struct {
	X string
}

func (s Send) Run(a *internal.Assem) {
	a.ExtraMemory["sound"] = a.GetValOrNum(s.X)
	a.InstIdx++
}

func (s Send) String() string {
	return fmt.Sprintf("snd %v", s.X)
}

type Subtract struct {
	X, Y string
}

func (s Subtract) Run(a *internal.Assem) {
	a.SetReg(s.X, a.GetValOrNum(s.X)-a.GetValOrNum(s.Y))
	a.InstIdx++
}

func (s Subtract) String() string {
	return fmt.Sprintf("sub %v %v", s.X, s.Y)
}

type Set struct {
	X, Y string
}

func (s Set) Run(a *internal.Assem) {
	a.SetReg(s.X, a.GetValOrNum(s.Y))
	a.InstIdx++
}
func (s Set) String() string {
	return fmt.Sprintf("set %v %v", s.X, s.Y)
}

type Add struct {
	X, Y string
}

func (a Add) Run(asm *internal.Assem) {
	asm.SetReg(a.X, asm.GetValOrNum(a.X)+asm.GetValOrNum(a.Y))
	asm.InstIdx++
}
func (a Add) String() string {
	return fmt.Sprintf("add %v %v", a.X, a.Y)
}

type Multiply struct {
	X, Y string
}

func (m Multiply) Run(a *internal.Assem) {
	a.SetReg(m.X, a.GetValOrNum(m.X)*a.GetValOrNum(m.Y))
	a.InstIdx++
}
func (m Multiply) String() string {
	return fmt.Sprintf("mul %v %v", m.X, m.Y)
}

type Modulo struct {
	X, Y string
}

func (m Modulo) Run(a *internal.Assem) {
	a.SetReg(m.X, a.GetValOrNum(m.X)%a.GetValOrNum(m.Y))
	a.InstIdx++
}
func (m Modulo) String() string {
	return fmt.Sprintf("mod %v %v", m.X, m.Y)
}

type Recover struct {
	X string
}

func (r Recover) Run(a *internal.Assem) {
	if a.GetValOrNum(r.X) != 0 {
		fmt.Println(a.ExtraMemory["sound"])
		os.Exit(0)
	}
	a.InstIdx++
}
func (r Recover) String() string {
	return fmt.Sprintf("rcv %v ", r.X)
}

type JumpGreaterThanZero struct {
	X, Y string
}

func (j JumpGreaterThanZero) Run(a *internal.Assem) {
	if a.GetValOrNum(j.X) > 0 {
		a.InstIdx += a.GetValOrNum(j.Y)
		return
	}
	a.InstIdx++
}

func (j JumpGreaterThanZero) String() string {
	return fmt.Sprintf("jgz %v %v", j.X, j.Y)
}
