package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

// so probably....
// 7, over 2, 11 over 2, , 12 over 2, 11 over 2, 11 over 2, 12 over 2 11 over 2
/// some annoying fraction

// 4 over 2, 3 over 2, 4 over 2, 3 over 2, 3 over 2, 3 over 2, 4 over 2, 3 over 2, 3 over 2, 4 over 2, 3 over 2, 3 over 2
type r struct {
	low, high int
}

func newr(x int) r {
	return r{x, x + 100}
}

func main() {
	input := input.GetInput(2019, 19)[0]
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}
	s := &arcade{
		screen: internal.NewGridV2WithDefaultChar("."),
	}

	ys := 1348
	xs := 1564
	// for changing {
	s = &arcade{
		screen: internal.NewGridV2WithDefaultChar("."),
	}
	yr := newr(ys)
	xr := newr(xs)
	for y := yr.low; y < yr.high; y++ {
		for x := xr.low; x < xr.high; x++ {
			s.robot = newRobot(nums)
			s.input = []int{x, y}
			s.robot.act(s)
		}
	}
	fmt.Println(s)
	// count := 0
	// for i := 0; i < 100; i++ {
	// 	if s.screen.Data[internal.Point{i, 0}] == "#" {
	// 		count++
	// 	}
	// }
	// 		if count > 95 {
	// 			if s.screen.Data[internal.Point{99, 0}] == "#" {
	// 				xs = xs + 1
	// 				changing = true
	// 			}
	// 			if s.screen.Data[internal.Point{0, 99}] == "#" {
	// 				ys = ys - 1
	// 				changing = true
	// 			}
	// 			if s.screen.Data[internal.Point{99, 0}] == "." {
	// 				xs = xs - 1
	// 				changing = true
	// 			}
	// 			if s.screen.Data[internal.Point{0, 99}] == "." {
	// 				ys = ys - 1
	// 				changing = true
	// 			}
	// 			fmt.Println(s)
	// 			fmt.Println(xs, ys)
	// 			changing = true
	// 			continue
	// 		}

	// 		// find the point farthest right and down and use that as the next input
	// 		m := r{}
	// 		for y := 0; y < 100; y++ {
	// 			max := math.MinInt
	// 			for x := 0; x < 100; x++ {
	// 				p := internal.Point{x, y}
	// 				v := s.screen.Data[p]
	// 				if v == "#" {
	// 					if x > max {
	// 						m.low = x
	// 						m.high = y
	// 					}
	// 				}
	// 			}
	// 		}

	// 		ys = ys + m.high
	// 		xs = xs + m.low
	// 		changing = true
	// 		fmt.Println(s)
	// 		fmt.Println(xs, ys)

	// 		continue
	// 		// go down
	// 		if s.screen.Data[internal.Point{99, 99}] == "." {
	// 			ys += 100
	// 			changing = true
	// 		}
	// 		if s.screen.Data[internal.Point{0, 99}] == "." {
	// 			xs += 100
	// 			changing = true
	// 		}
	// //	}
	// 	fmt.Println(ys, xs)
	// 	fmt.Println(s)
	// for y := 0; y < 100; y++ {
	// 	min := math.MaxInt
	// 	max := math.MinInt
	// 	for x := 0; x < 100; x++ {
	// 		p := internal.Point{x, y}
	// 		v := s.screen.Data[p]
	// 		if v == "#" {
	// 			if x > max {
	// 				max = x
	// 			}
	// 			if x < min {
	// 				min = x
	// 			}
	// 			//				fmt.Print(p)
	// 		}
	// 	}
	// 	//		fmt.Println(internal.Point{min, y}, internal.Point{max, y})
	// 	//		fmt.Println(max - min)
	// }
	// fmt.Println(s)
}

func find(g *internal.GridV2, search string) internal.Point {
	for p, v := range g.Data {
		if v == search {
			return p
		}
	}
	return internal.Point{}
}

func newRobot(nums []int) *robot {
	indexData := map[int]int{}
	for i, v := range nums {
		indexData[i] = v
	}
	comp := &opcodeComputer{data: indexData, silent: true, debug: false}
	return &robot{
		brain: comp,
	}
}

func (r *robot) act(s *arcade) {
	r.brain.inputs = s.input
	r.brain.run()
	status := r.brain.output
	s.print(status)
}

type arcade struct {
	screen  *internal.GridV2
	pointer internal.Point
	robot   *robot
	input   []int
	x, y    int
}

func (s *arcade) String() string {
	return s.screen.String()
}

func (s *arcade) print(status int) {
	switch status {
	case 1:
		s.screen.Set(internal.Point{s.x, s.y}, "#")
		s.x++
	case 0:
		s.screen.Set(internal.Point{s.x, s.y}, ".")
		s.x++
	}
	if s.x == 100 {
		s.y++
		s.x = 0
	}
}

type robot struct {
	brain *opcodeComputer
}

type opcodeComputer struct {
	data         map[int]int
	cur          int
	inputs       []int
	curInput     int
	output       int
	relativeBase int
	debug        bool
	silent       bool
	done         bool
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
func (o *opcodeComputer) addToRelativeBase(i int) {
	o.relativeBase += i
}
func (o *opcodeComputer) setInstructionPointer(i int) {
	o.cur = i
}
func (o *opcodeComputer) setRegisterUsingInput(reg int) {
	o.data[reg] = o.inputs[o.curInput]
	o.curInput++
}
func (o *opcodeComputer) setRegisterWithValue(reg, val int) {
	o.data[reg] = val
}
func (o *opcodeComputer) error(msg string) {
	fmt.Println("unexpected error:", msg)
	os.Exit(1)
}
func (o *opcodeComputer) getRegister(reg int) int {
	if reg < 0 {
		panic("cannot read negative register")
	}
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
		return doInput(op)
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
	case 9:
		return adjustRelativeBase(op)
	case 99:
		o.done = true
		return nil
	default:
		o.error(fmt.Sprintf("unknown operator %d", op.op))
		return nil
	}
}
func adjustRelativeBase(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(1)
		if o.debug {
			fmt.Println("updating relative base, raw input:", vals[0])
		}
		switch oc.modes[0] {
		case 0:
			vals[0] = o.getRegister(vals[0])
		case 2:
			vals[0] = o.getRegister(vals[0] + o.relativeBase)
		}
		if o.debug {
			fmt.Println("updating relative base, resolved: ", vals[0], "(", oc, ")")
		}
		o.addToRelativeBase(vals[0])
		if o.debug {
			fmt.Println("new relative base:", o.relativeBase)
		}
		return opcodeRead
	}
}
func jumpIfTrue(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(2)
		if o.debug {
			fmt.Println("jump-if-true raw input", vals)
		}
		for i := 0; i < 2; i++ {
			switch oc.modes[i] {
			case 0:
				vals[i] = o.getRegister(vals[i])
			case 2:
				vals[i] = o.getRegister(vals[i] + o.relativeBase)
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
			switch oc.modes[i] {
			case 0:
				vals[i] = o.getRegister(vals[i])
			case 2:
				vals[i] = o.getRegister(vals[i] + o.relativeBase)
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
			switch oc.modes[i] {
			case 0:
				vals[i] = o.getRegister(vals[i])
			case 2:
				vals[i] = o.getRegister(vals[i] + o.relativeBase)
			}
		}
		if oc.modes[2] == 2 {
			vals[2] += o.relativeBase
		}
		lessThan := 0
		if vals[0] < vals[1] {
			lessThan = 1
		}
		if o.debug {
			if lessThan == 1 {
				fmt.Println("setting register", vals[2], "to", lessThan, "becasue", vals[0], "<", vals[1])
			} else {
				fmt.Println("setting register", vals[2], "to", lessThan, "becasue", vals[0], ">=", vals[1])
			}
		}
		o.setRegisterWithValue(vals[2], lessThan)
		return opcodeRead
	}
}
func equals(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(3)
		for i := 0; i < 2; i++ {
			switch oc.modes[i] {
			case 0:
				vals[i] = o.getRegister(vals[i])
			case 2:
				vals[i] = o.getRegister(vals[i] + o.relativeBase)
			}
		}
		if oc.modes[2] == 2 {
			vals[2] += o.relativeBase
		}
		isEqual := 0
		if vals[0] == vals[1] {
			isEqual = 1
		}
		if o.debug {
			if isEqual == 1 {
				fmt.Println("setting register", vals[2], "to", isEqual, "becasue", vals[0], "==", vals[1])
			} else {
				fmt.Println("setting register", vals[2], "to", isEqual, "becasue", vals[0], "!=", vals[1])
			}
		}
		o.setRegisterWithValue(vals[2], isEqual)
		return opcodeRead
	}
}

func doOutput(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(1)
		switch oc.modes[0] {
		case 0:
			if o.debug {
				fmt.Printf("reading register %d for output\n", vals[0])
			}
			vals[0] = o.getRegister(vals[0])
		case 2:
			if o.debug {
				fmt.Printf("reading register %d for output\n", vals[0]+o.relativeBase)
			}
			vals[0] = o.getRegister(vals[0] + o.relativeBase)
		}
		o.output = vals[0]
		if !o.silent {
			fmt.Printf("output value: %d\n", vals[0])
		}
		return nil
	}
}

func doInput(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(1)
		if o.debug {
			fmt.Println("setting register", vals[0], "to input value", o.inputs[o.curInput])
		}
		switch oc.modes[0] {
		case 2:
			vals[0] = vals[0] + o.relativeBase
		}
		o.setRegisterUsingInput(vals[0])
		return opcodeRead
	}
}

// reading 3 params
func addParams(oc opcode) stateFn {
	return func(o *opcodeComputer) stateFn {
		vals := o.readInputs(3)
		if o.debug {
			fmt.Println("add raw input", vals)
		}
		for i := 0; i < 2; i++ {
			switch oc.modes[i] {
			case 0:
				vals[i] = o.getRegister(vals[i])
			case 2:
				vals[i] = o.getRegister(vals[i] + o.relativeBase)
			}
		}
		if oc.modes[2] == 2 {
			vals[2] += o.relativeBase
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
			switch oc.modes[i] {
			case 0:
				vals[i] = o.getRegister(vals[i])
			case 2:
				vals[i] = o.getRegister(vals[i] + o.relativeBase)
			}
		}
		if oc.modes[2] == 2 {
			vals[2] += o.relativeBase
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

//1187721666102244
//9223372036854775807
