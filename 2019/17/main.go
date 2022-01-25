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
	input := input.GetInput(2019, 17)[0]
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}

	s := &arcade{
		screen: internal.NewGridV2WithDefaultChar(" "),
		robot:  newRobot(nums),
	}

	for !s.robot.brain.done {
		s.robot.act(s)
	}
	fmt.Println(s)
	//	fmt.Println(getAlignment(scaffoldIntersections(s.screen)))
	path := findPath(s.screen)
	simple := simplifyPath(path)
	fmt.Println(simple)
	combos := buildMvmtMap(simple)
	replacements, mainControl := tryCombos(combos, simple)
	brainInput := []byte{}
	for _, c := range mainControl {
		brainInput = append(brainInput, byte(c))
	}
	brainInput = append(brainInput, byte('\n'))
	cs := make([]string, 3)
	for k, v := range replacements {
		switch v {
		case "A":
			cs[0] = k
		case "B":
			cs[1] = k
		case "C":
			cs[2] = k
		}
	}
	fmt.Println("CS", cs)
	for _, s := range cs {
		for _, c := range s {
			brainInput = append(brainInput, byte(c))
		}
		brainInput = append(brainInput, byte('\n'))
	}
	brainInput = append(brainInput, byte('n'), byte('\n'))
	in := make([]int, len(brainInput))
	for i := 0; i < len(in); i++ {
		in[i] = int(brainInput[i])
	}
	fmt.Println("sending", in)
	s = &arcade{
		screen: internal.NewGridV2WithDefaultChar(" "),
		robot:  newRobot(nums),
		input:  in,
	}
	s.robot.brain.data[0] = 2
	for !s.robot.brain.done {
		s.robot.act(s)
	}
	fmt.Println(s.robot.brain.output)
}

func tryCombos(combos []string, path []string) (map[string]string, string) {
	smallest := strings.Repeat("22222", 100)
	out := map[string]string{}
	for i := 0; i < len(combos)-2; i++ {
		for j := i + 1; j < len(combos)-1; j++ {
			for k := j + 1; k < len(combos); k++ {
				mvmtMap := map[string]string{
					combos[i]: "A",
					combos[j]: "B",
					combos[k]: "C",
				}
				m := reduce(mvmtMap, path)
				if len(m) < len(smallest) {
					out = mvmtMap
					smallest = m
				}
			}
		}
	}
	return out, smallest
}

func buildMvmtMap(path []string) []string {
	//	joined := strings.Join(path, ",")
	// build combos of size <= 20
	combos := []string{}
	start := 0
	for start <= len(path)-4 {
		cur := []string{}
		for i := start; i < len(path); i += 2 {
			cur = append(cur, path[i], path[i+1])
			mvmt := strings.Join(cur, ",")
			if len(mvmt) <= 20 {
				combos = append(combos, mvmt)
				continue
			}
			break
		}
		start += 2
	}
	return combos
}

func reduce(movement map[string]string, path []string) string {
	p := strings.Join(path, ",")
	for k, v := range movement {
		p = strings.ReplaceAll(p, k, v)
	}
	return p
}

func simplifyPath(in []string) []string {
	// remove 0s
	for i := 0; i < len(in); i++ {
		if in[i] == "0" {
			// special case 0
			in = append(in[:i], in[i+1:]...)
			i = 0
			continue
		}
	}
	// simplify turns
	count := 0
	for i := 0; i < len(in); i++ {
		if in[i] == "R" {
			count++
			if count == 3 {
				in = append(append(in[:i-2], "L"), in[i+1:]...)
				i = 0
				continue
			}
		} else {
			count = 0
		}
	}
	return in
}

func findPath(g *internal.GridV2) []string {
	start := find(g, "^")
	end := findEnd(g)
	path := []string{}
	t := &turtle{pos: start, facing: north}
	visited := map[internal.Point]bool{}
	for {
		//		fmt.Println(t.pos)
		// go forward as many times as possible
		count := 0
		for {
			t.forward()
			count++
			if t.pos == end {
				break
			}
			if g.Data[t.pos] != "#" || visited[t.pos] {
				count--
				t.backup()
				break
			}
			left := internal.Point{t.pos.X - 1, t.pos.Y}
			right := internal.Point{t.pos.X + 1, t.pos.Y}
			if t.facing == north || t.facing == south {
				if g.Data[left] == "#" && g.Data[right] == "#" {
					continue
				}
			}
			top := internal.Point{t.pos.X, t.pos.Y - 1}
			bottom := internal.Point{t.pos.X, t.pos.Y + 1}
			if t.facing == east || t.facing == west {
				if g.Data[top] == "#" && g.Data[bottom] == "#" {
					continue
				}
			}
			visited[t.pos] = true
		}
		path = append(path, fmt.Sprintf("%d", count))
		count = 0
		if t.pos == end {
			break
		}
		t.turnRight()
		path = append(path, "R")
	}
	return path
}

type turtle struct {
	pos    internal.Point
	facing dir
}

func (t *turtle) turnRight() {
	switch t.facing {
	case north:
		t.facing = east
	case east:
		t.facing = south
	case south:
		t.facing = west
	case west:
		t.facing = north
	}
}
func (t *turtle) forward() {
	switch t.facing {
	case north:
		t.pos.Y -= 1
	case east:
		t.pos.X += 1
	case south:
		t.pos.Y += 1
	case west:
		t.pos.X -= 1
	}
}
func (t *turtle) backup() {
	switch t.facing {
	case north:
		t.pos.Y += 1
	case east:
		t.pos.X -= 1
	case south:
		t.pos.Y -= 1
	case west:
		t.pos.X += 1
	}
}

func findEnd(g *internal.GridV2) internal.Point {
	for p, v := range g.Data {
		if v == "#" {
			count := 0
			for _, n := range p.Neighbors() {
				if g.Data[n] == "." {
					count++
				}
			}
			if count == 3 {
				return p
			}
		}
	}
	return internal.Point{-1, -1}
}

func find(g *internal.GridV2, search string) internal.Point {
	for p, v := range g.Data {
		if v == search {
			return p
		}
	}
	return internal.Point{}
}

func getAlignment(intersections []internal.Point) int {
	out := 0
	for _, p := range intersections {
		out += alignmentParameter(p)
	}
	return out
}

func alignmentParameter(p internal.Point) int {
	return p.X * p.Y
}

func scaffoldIntersections(g *internal.GridV2) []internal.Point {
	out := []internal.Point{}
	for p, v := range g.Data {
		if v == "#" {
			found := true
			for _, n := range p.Neighbors() {
				if g.At(n) != "#" {
					found = false
					break
				}
			}
			if found {
				out = append(out, p)
			}
		}
	}
	return out
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
	if r.brain.done {
		return
	}
	status := r.brain.output
	s.print(status)
}

type arcade struct {
	screen  *internal.GridV2
	pointer internal.Point
	robot   *robot
	input   []int
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

func (s *arcade) String() string {
	return s.screen.String()
}

func (s *arcade) print(status int) {
	switch status {
	case 10:
		s.pointer.X = 0
		s.pointer.Y = s.pointer.Y + 1
	default:
		s.screen.Set(s.pointer, string(byte(status)))
		s.pointer.X = s.pointer.X + 1
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
