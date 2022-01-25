package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

// 38420 is too low (part 1)

func main() {
	c := internal.AllPermutations([]string{}, []string{"5", "6", "7", "8", "9"})

	// fmt.Println(os.Args, len(os.Args))
	// phaseOrder := make([]int, 5)
	// if len(os.Args) != 6 {
	// 	panic("bad number of args")
	// }
	// for i := 1; i < 6; i++ {
	// 	phaseOrder[i-1], _ = strconv.Atoi(os.Args[i])
	// }
	input := internal.ReadInput()[0]
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}

	max := math.MinInt
	com := []string{}
	for _, combo := range c {
		phaseOrder := make([]int, 5)
		for i := 0; i < 5; i++ {
			phaseOrder[i], _ = strconv.Atoi(combo[i])
		}
		out := runThing(nums, phaseOrder)
		if out > max {
			max = out
			com = combo
		}
	}
	fmt.Println(com, max)
}

func runThing(nums, phaseOrder []int) int {
	// build amplifiers
	amplis := make([]*opcodeComputer, 0)
	for i := 0; i < 5; i++ {
		o := make([]int, len(nums))
		copy(o, nums)
		amplis = append(amplis, &opcodeComputer{data: o, debug: false, silent: true})
	}

	nextInput := 0
	curAmpli := 0
	for {
		ampli := amplis[curAmpli]
		ampli.inputs = []int{phaseOrder[curAmpli], nextInput}
		ampli.run()
		if ampli.done {
			break
		}
		nextInput = ampli.output
		curAmpli += 1
		curAmpli %= 5
	}
	return amplis[len(amplis)-1].output
}

// 38420 is too low

type opcodeComputer struct {
	data     []int
	cur      int
	inputs   []int
	curInput int
	output   int
	debug    bool
	silent   bool
	done     bool
}

type stateFn func(o *opcodeComputer) stateFn

func (o *opcodeComputer) readInputs(n int) []int {
	out := []int{}
	for i := 0; i < n; i++ {
		v, more := o.read()
		if !more {
			o.error("no more input")
			return out
		}
		out = append(out, v)
	}
	return out
}

func (o *opcodeComputer) run() {
	// if o.done {
	// 	panic("trying to run a done amplifier")
	// }
	for state := opcodeRead(o); state != nil; {
		state = state(o)
	}
}
func (o *opcodeComputer) read() (int, bool) {
	if o.cur >= len(o.data) {
		return 0, false
	}
	out := o.data[o.cur]
	o.cur++
	return out, true
}
func (o *opcodeComputer) setInstructionPointer(i int) {
	o.cur = i
}
func (o *opcodeComputer) setRegisterUsingInput(reg int) {
	o.data[reg] = o.inputs[o.curInput]
	if o.curInput == 0 {
		o.curInput++
	}
}
func (o *opcodeComputer) setRegisterWithValue(reg, val int) {
	o.data[reg] = val
}
func (o *opcodeComputer) error(msg string) {
	fmt.Println("unexpected error:", msg)
	os.Exit(1)
}
func (o *opcodeComputer) getRegister(reg int) int {
	return o.data[reg]
}

func opcodeRead(o *opcodeComputer) stateFn {
	oc, more := o.read()
	if !more {
		return nil
	}
	op := parseOpcode(oc)
	switch op.op {
	case 1:
		return addParams(op)
	case 2:
		return mulParams(op)
	case 3:
		return doInput
	case 4:
		return doOutput(op)
	case 5:
		return jumpIfTrue(op)
	case 6:
		return jumpIfFalse(op)
	case 7:
		return lessThan(op)
	case 8:
		return equals(op)
	case 99:
		o.done = true
		return nil
	default:
		o.error(fmt.Sprintf("unknown operator %d", op.op))
		return nil
	}
}
func jumpIfTrue(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(2)
		if o.debug {
			fmt.Println("jump-if-true raw input", vals)
		}
		for i := 0; i < 2; i++ {
			if oc.modes[i] == 0 {
				vals[i] = o.getRegister(vals[i])
			}
		}
		if vals[0] != 0 {
			o.setInstructionPointer(vals[1])
		}
		if o.debug {
			fmt.Println("jump-if-true resolved ", vals, "(", oc, ")")
			if vals[0] != 0 {
				fmt.Println("jumping to", vals[1])
			} else {
				fmt.Println("not jumping", vals[0], "==", 0)
			}
		}
		return opcodeRead
	}
}
func jumpIfFalse(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(2)
		if o.debug {
			fmt.Println("jump-if-false raw input", vals)
		}
		for i := 0; i < 2; i++ {
			if oc.modes[i] == 0 {
				vals[i] = o.getRegister(vals[i])
			}
		}
		if vals[0] == 0 {
			o.setInstructionPointer(vals[1])
		}
		if o.debug {
			fmt.Println("jump-if-false resolved ", vals, "(", oc, ")")
			if vals[0] == 0 {
				fmt.Println("jumping to", vals[1])
			} else {
				fmt.Println("not jumping", vals[0], "!=", 0)
			}
		}
		return opcodeRead
	}
}
func lessThan(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(3)
		for i := 0; i < 2; i++ {
			if oc.modes[i] == 0 {
				vals[i] = o.getRegister(vals[i])
			}
		}
		lessThan := 0
		if vals[0] < vals[1] {
			lessThan = 1
		}
		o.setRegisterWithValue(vals[2], lessThan)
		return opcodeRead
	}
}
func equals(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(3)
		for i := 0; i < 2; i++ {
			if oc.modes[i] == 0 {
				vals[i] = o.getRegister(vals[i])
			}
		}
		isEqual := 0
		if vals[0] == vals[1] {
			isEqual = 1
		}
		o.setRegisterWithValue(vals[2], isEqual)
		return opcodeRead
	}
}

func doOutput(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(1)
		if oc.modes[0] == 0 {
			value := o.getRegister(vals[0])
			o.output = value
			if !o.silent {
				fmt.Printf("Register %d value: %d\n", vals[0], value)
			}
			return nil
		}
		if !o.silent {
			fmt.Printf("Value: %d\n", vals[0])
		}
		o.output = vals[0]
		return nil
	}
}

func doInput(o *opcodeComputer) stateFn {
	vals := o.readInputs(1)
	if o.debug {
		fmt.Println("setting register", vals[0], "to input value", o.inputs[o.curInput])
	}
	o.setRegisterUsingInput(vals[0])
	return opcodeRead
}

// reading 3 params
func addParams(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(3)
		if o.debug {
			fmt.Println("add raw input", vals)
		}
		for i := 0; i < 2; i++ {
			if oc.modes[i] == 0 {
				vals[i] = o.getRegister(vals[i])
			}
		}
		o.setRegisterWithValue(vals[2], vals[0]+vals[1])
		if o.debug {
			fmt.Println("add resolved ", vals, "(", oc, ")")
			fmt.Println("setting register", vals[2], "to", vals[0], "+", vals[1], "=", vals[0]+vals[1])
		}
		return opcodeRead
	}
}

func mulParams(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(3)
		if o.debug {
			fmt.Println("mul raw input", vals)
		}
		for i := 0; i < 2; i++ {
			if oc.modes[i] == 0 {
				vals[i] = o.getRegister(vals[i])
			}
		}
		o.setRegisterWithValue(vals[2], vals[0]*vals[1])
		if o.debug {
			fmt.Println("mul resolved ", vals, "(", oc, ")")
			fmt.Println("setting register", vals[2], "to", vals[0], "*", vals[1], "=", vals[0]*vals[1])
		}
		return opcodeRead
	}
}

type opcode struct {
	op    int
	modes []int
}

func parseOpcode(in int) opcode {
	oo := in
	oo %= 10000
	oo %= 1000
	oo %= 100
	o := (oo/10)*10 + in%10
	return opcode{
		op:    o,
		modes: modes(in),
	}
}

func modes(in int) []int {
	if in < 100 {
		return []int{0, 0, 0}
	}
	a := in / 10000
	in %= 10000
	b := in / 1000
	in %= 1000
	c := in / 100
	return []int{c, b, a}
}

// 294 is too low

func combos(items []string) [][]string {
	out := [][]string{}
	// len(items) rotations
	for k := 0; k < len(items); k++ {
		og := make([]string, len(items))
		copy(og, items)
		out = append(out, og)
		// for each rotation
		for i := 0; i < len(items)-1; i++ {
			for j := i + 1; j < len(items); j++ {
				out = append(out, swap(items, i, j))
			}
		}
		items = rotate(items)
	}
	return out
}

func rotate(items []string) []string {
	out := make([]string, len(items))
	for i := 1; i < len(items); i++ {
		out[i-1] = items[i]
	}
	out[len(out)-1] = items[0]
	return out
}

func swap(items []string, i, j int) []string {
	out := make([]string, len(items))
	for k := 0; k < len(items); k++ {
		if k == i {
			out[k] = items[j]
			continue
		}
		if k == j {
			out[k] = items[i]
			continue
		}
		out[k] = items[k]
	}
	return out
}
