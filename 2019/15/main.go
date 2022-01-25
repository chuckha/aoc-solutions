package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

type dir int

const (
	north dir = 1
	east  dir = 4
	south dir = 2
	west  dir = 3
)

// 21107 too low (part 1)
// 4235104155 too high

func main() {
	input := input.GetInput(2019, 15)[0]
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}

	s := &arcade{
		screen:      internal.NewGridV2WithDefaultChar(" "),
		robot:       newRobot(nums),
		movingState: newState(),
	}

	s.updateInput(1)
	for i := 0; i < 20000; i++ {
		s.robot.act(s)
	}
	costs := cost(s.screen, internal.Point{0, 0})
	fmt.Println("part 1", costs[s.oxygenSensor])
	costs2 := cost(s.screen, s.oxygenSensor)
	fmt.Println("part 2", farthestPoint(costs2))
	fmt.Println(s)
}

func farthestPoint(costs map[internal.Point]int) int {
	biggest := 0
	for _, c := range costs {
		if c != math.MaxInt && c > biggest {
			biggest = c
		}
	}
	return biggest
}

// part 2 394 is too high (but curiously the right answer for someone else)

func newRobot(nums []int) *robot {
	indexData := map[int]int{}
	for i, v := range nums {
		indexData[i] = v
	}
	comp := &opcodeComputer{data: indexData, silent: true, debug: false}
	return &robot{
		brain:  comp,
		facing: up,
	}
}

func doublePrint(g *internal.GridV2, costs map[internal.Point]int) string {
	var out strings.Builder
	for j := g.Min.Y; j <= g.Max.Y; j++ {
		for i := g.Min.X; i <= g.Max.X; i++ {
			if c, ok := costs[internal.Point{i, j}]; ok {
				if c == math.MaxInt {
					out.WriteString("???")
					continue
				}
				out.WriteString(fmt.Sprintf("%03d", c))
				continue
			}
			out.WriteString("###")
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (r *robot) act(s *arcade) {
	r.brain.inputs = []int{s.input}
	r.brain.run()
	if r.brain.done {
		return
	}
	status := r.brain.output
	s.print(status)
}

type arcade struct {
	screen       *internal.GridV2
	robot        *robot
	movement     int
	input        int
	movingState  *state
	oxygenSensor internal.Point
}

func cost(g *internal.GridV2, start internal.Point) map[internal.Point]int {
	costs := map[internal.Point]int{}
	for k := range g.Data {
		costs[k] = math.MaxInt
	}
	visited := map[internal.Point]bool{}
	costs[start] = 0
	q := internal.NewQueue[internal.Point]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		for _, n := range cur.Neighbors() {
			item, ok := g.Data[n]
			if !ok {
				continue
			}
			if item == "#" {
				continue
			}
			if costs[cur]+1 < costs[n] {
				costs[n] = costs[cur] + 1
			}
			if visited[n] {
				continue
			}
			q.Enqueue(n)
		}
		visited[cur] = true
	}
	return costs
}

type state struct {
	init      bool
	stepRight bool
	checkLeft bool
}

func newState() *state {
	return &state{
		init:      true,
		stepRight: false,
		checkLeft: false,
	}
}

func (s *arcade) updateInput(output int) {
	// get destintation
	// if i'm not at the destination
	// set the input to move me toward the destination
	// if i'm at the next destination, start going toward the next destination

	//	fmt.Println(s.movingState, s.robot, output)
	if s.movingState.init && output == 0 {
		s.movingState.init = false
		s.movingState.stepRight = true
	}
	if s.movingState.init {
		s.input = s.robot.forward()
		return
	}

	if s.movingState.stepRight {
		s.input = s.robot.goRight()
		s.movingState.stepRight = false
		s.movingState.checkLeft = true
		return
	}
	if s.movingState.checkLeft {
		if output == 0 {
			s.movingState.stepRight = true
			s.movingState.checkLeft = false
			s.updateInput(output)
			return
		}
		s.input = s.robot.goLeft()
		return
	}
	panic("bad state in the state machine")
	// go forward until you hit a wall (init)

	// turn right
	// go forward
	// turn left
	// go forward
	// if you hit a wall,
	//		turn right
	//		go forward
	//		turn left
	// if you don't hit a wall
	// 		turn left
	//		go forward
	// 		turn right, go forward

}

func (s *arcade) String() string {
	s.screen.Set(s.robot.pos, "D")
	s.screen.Set(s.oxygenSensor, "O")
	s.screen.Set(internal.Point{0, 0}, "S")
	return s.screen.String()
}

func (s *arcade) print(status int) {
	switch status {
	case 0: // hit a wall in the requestesd direction
		switch s.input {
		case int(north):
			s.screen.Set(internal.Point{s.robot.pos.X, s.robot.pos.Y - 1}, "#")
		case int(east):
			s.screen.Set(internal.Point{s.robot.pos.X + 1, s.robot.pos.Y}, "#")
		case int(south):
			s.screen.Set(internal.Point{s.robot.pos.X, s.robot.pos.Y + 1}, "#")
		case int(west):
			s.screen.Set(internal.Point{s.robot.pos.X - 1, s.robot.pos.Y}, "#")
		}
	case 1: // has moved one step in the requested direction
		switch s.input {
		case int(north):
			s.robot.pos.Y -= 1
		case int(east):
			s.robot.pos.X += 1
		case int(south):
			s.robot.pos.Y += 1
		case int(west):
			s.robot.pos.X -= 1
		}
		s.screen.Set(s.robot.pos, ".")

	case 2: // it is now on the oxygen system
		switch s.input {
		case int(north):
			s.robot.pos.Y -= 1
		case int(east):
			s.robot.pos.X += 1
		case int(south):
			s.robot.pos.Y += 1
		case int(west):
			s.robot.pos.X -= 1
		}
		s.screen.Set(s.robot.pos, ".")
		s.oxygenSensor = s.robot.pos
		//		fmt.Println("Found oxygen sensor at", s.robot.pos)
	}
	s.updateInput(status)
}

const (
	up    = "up"
	right = "right"
	down  = "down"
	left  = "left"
)

type robot struct {
	pos    internal.Point
	brain  *opcodeComputer
	facing string
}

func (r *robot) forward() int {
	switch r.facing {
	case up:
		return int(north)
	case right:
		return int(east)
	case down:
		return int(south)
	case left:
		return int(west)
	}
	panic("ljkds")
}
func (r *robot) goRight() int {
	switch r.facing {
	case up:
		r.facing = right
		return int(east)
	case right:
		r.facing = down
		return int(south)
	case down:
		r.facing = left
		return int(west)
	case left:
		r.facing = up
		return int(north)
	}
	panic("bad news bears")
}
func (r *robot) goLeft() int {
	switch r.facing {
	case up:
		r.facing = left
		return int(west)
	case right:
		r.facing = up
		return int(north)
	case down:
		r.facing = right
		return int(east)
	case left:
		r.facing = down
		return int(south)
	}
	panic("bad news bears")
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
	// if o.curInput == 0 {
	// 	o.curInput++
	// }
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
