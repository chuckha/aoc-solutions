package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type dir string

const (
	up    dir = "up"
	right dir = "right"
	down  dir = "down"
	left  dir = "left"
)

// 21107 too low (part 1)
// 4235104155 too high

func main() {
	input := internal.ReadInput()[0]
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}

	s := &ship{
		shell:         internal.NewGridV2(),
		robot:         newRobot(nums),
		paintedPanels: make(map[internal.Point]bool),
	}
	s.paint(1)
	for !s.robot.brain.done {
		s.act()
	}
	fmt.Println(s)
	fmt.Println(len(s.paintedPanels))
}

// 1223 is too low (part 1)
// 2253 is too high (part 1)

func newRobot(nums []int) *robot {
	indexData := map[int]int{}
	for i, v := range nums {
		indexData[i] = v
	}
	comp := &opcodeComputer{data: indexData, silent: true, debug: false}
	return &robot{
		position: internal.Point{0, 0},
		facing:   up,
		brain:    comp,
	}
}

func (r *robot) act(s *ship) {
	r.brain.inputs = []int{s.cameraInput()}
	r.brain.run()
	if r.brain.done {
		return
	}
	color := r.brain.output
	s.paint(color)
	r.brain.run()
	turn := r.brain.output
	r.turn(turn)
	r.move()
}

type ship struct {
	shell         *internal.GridV2
	robot         *robot
	paintedPanels map[internal.Point]bool
}

func (s *ship) act() {
	s.robot.act(s)
}

func (s *ship) paint(i int) {
	paint := "."
	if i == 1 {
		paint = "#"
	}
	s.paintedPanels[s.robot.position] = true
	s.shell.Set(s.robot.position, paint)
}

func (s *ship) cameraInput() int {
	color := s.shell.At(s.robot.position)
	if color == "#" {
		return 1
	}
	return 0
}

func (s *ship) String() string {
	cur := s.shell.At(s.robot.position)
	s.shell.Set(s.robot.position, s.robot.String())
	out := s.shell.String()
	s.shell.Set(s.robot.position, cur)
	return out
}

type robot struct {
	position internal.Point
	facing   dir
	brain    *opcodeComputer
}

func (r *robot) String() string {
	switch r.facing {
	case up:
		return "^"
	case right:
		return ">"
	case down:
		return "v"
	case left:
		return "<"
	}
	panic("lol")
}
func (r *robot) turn(i int) {
	switch i {
	// left
	case 0:
		switch r.facing {
		case up:
			r.facing = left
		case right:
			r.facing = up
		case down:
			r.facing = right
		case left:
			r.facing = down
		}
		// right
	case 1:
		switch r.facing {
		case up:
			r.facing = right
		case right:
			r.facing = down
		case down:
			r.facing = left
		case left:
			r.facing = up
		}
	}
}
func (r *robot) move() {
	switch r.facing {
	case up:
		r.position.Y--
	case right:
		r.position.X++
	case down:
		r.position.Y++
	case left:
		r.position.X--
	}
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
func (o *opcodeComputer) addToRelativeBase(i int) {
	o.relativeBase += i
}
func (o *opcodeComputer) setInstructionPointer(i int) {
	o.cur = i
}
func (o *opcodeComputer) setRegisterUsingInput(reg int) {
	o.data[reg] = o.inputs[o.curInput]
	//	if o.curInput == 0 {
	//		o.curInput++
	//	}
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
