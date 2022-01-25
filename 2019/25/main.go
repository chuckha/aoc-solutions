package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

const (
	north = "north"
	south = "south"
	west  = "west"
	east  = "east"
	inv   = "inv"
)

func HyperCube() []string {
	return []string{east, east, east, north, north, east, north, "take hypercube", south, west, south, south, west, west, west}
}
func Cake() []string {
	return []string{north, north, east, east, "take cake", west, west, south, south}
}
func KleinBottle() []string {
	return []string{east, east, east, north, north, east, "take klein bottle", west, south, south, west, west, west}
}
func DarkMatter() []string { return []string{east, east, east, "take dark matter", west, west, west} }
func EasterEgg() []string {
	return []string{south, west, west, "take easter egg", east, east, north}
}
func WeightRoom() []string {
	return []string{east, east, east, north, north, east, north, north, inv, west}
}
func FuelCell() []string { return []string{south, west, "take fuel cell", east, north} }

// get eaten by a grue ;(
func Photons() []string {
	return []string{east, south, east, east, "take photons", west, west, north, west}
}

// MoltenLava makes you melt
func MoltenLava() []string { return []string{east, south, east, "take molten lava", west, north, west} }

func Ornament() []string { return []string{east, "take ornament", west} }
func Hologram() []string { return []string{east, east, "take hologram", west, west} }
func main() {
	input := input.GetInput(2019, 25)[0]
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}

	// ptrs := []int{0, 1, 2, 3, 6}
	// for i := 0; i < 1000; i++ {
	// 	internal.Inc(ptrs, 6)
	// 	fmt.Println(ptrs)
	// }

	// s := &arcade{
	// 	robot: newRobot(nums),
	// 	input: cmd(north, "take escape pod"),
	// }
	// for i := 0; i < 3; i++ {
	// 	s.gameoutput = &gameOutput{}
	// 	for s.robot.act(s) {
	// 	}
	// 	fmt.Println(s.gameoutput.asRoom(s.loc).Describe())
	// }
	// fmt.Println(s.gameoutput.asRoom(s.loc).Describe())

	//testAll(nums)
	//	data(nums)
	winGame(nums)
}

func winGame(nums []int) {
	code := 1090617344
	cmds := []string{
		east, east, east, north, north, "1090617344",
	}
	s := &arcade{
		robot: newRobot(nums),
		input: cmd(cmds...),
	}
	//	r := s.gameoutput.asRoom(s.loc)
	for i := 0; i < len(cmds)+1; i++ {
		s.gameoutput = &gameOutput{}
		for s.robot.act(s) {
		}
		fmt.Println(s.gameoutput.data)
	}
	fmt.Println(code)
}

func testAll(nums []int) {
	objects := map[string][]string{
		"fuel":        FuelCell(),
		"ornament":    Ornament(),
		"dark matter": DarkMatter(),
		"hyper cube":  HyperCube(),
		"cake":        Cake(),
		"bottle":      KleinBottle(),
		"egg":         EasterEgg(),
		"hologram":    Hologram(),
	}
	keys := []string{}
	for k := range objects {
		keys = append(keys, k)
	}
	allCombos := [][]string{}
	for i := 1; i <= len(objects); i++ {
		allCombos = append(allCombos, internal.AllCombinations([]string{}, keys, i)...)
	}
	for _, combo := range allCombos {
		cmds := []string{}
		for _, c := range combo {
			cmds = append(cmds, objects[c]...)
		}
		cmds = append(cmds, WeightRoom()...)

		// cmds := []string{}
		// cmds = append(cmds, FuelCell()...)    // L
		// cmds = append(cmds, Ornament()...)    // L
		// cmds = append(cmds, DarkMatter()...)  // L
		// cmds = append(cmds, HyperCube()...)   // L
		// cmds = append(cmds, Cake()...)        // L
		// cmds = append(cmds, KleinBottle()...) // L
		// cmds = append(cmds, EasterEgg()...)   // H
		// cmds = append(cmds, Hologram()...)    // H
		// cmds = append(cmds, WeightRoom()...)
		s := &arcade{
			robot: newRobot(nums),
			input: cmd(cmds...),
		}
		//	r := s.gameoutput.asRoom(s.loc)
		for i := 0; i < len(cmds)+1; i++ {
			s.gameoutput = &gameOutput{}
			for s.robot.act(s) {
			}
			if s.gameoutput.name == "== Pressure-Sensitive Floor ==" {
				fmt.Println(combo)
				fmt.Println(s.gameoutput.data)
			}
		}
	}
	//	fmt.Println(s.gameoutput.asRoom(s.loc).Describe())
	//	r = s.gameoutput.asRoom(s.loc)
	//fmt.Println(r.Describe())
	// s.input = cmd(south)
	// s.robot.brain.curInput = 0
	// s.gameoutput = &gameOutput{}
	// for s.robot.act(s) {
	// }

}

func data(nums []int) {
	s := &arcade{
		gameoutput: &gameOutput{},
		robot:      newRobot(nums),
		input:      cmd(),
	}
	run(s)
}

func (a *arcade) getRoom() *room {
	for a.robot.act(a) {
	}
	return a.gameoutput.asRoom(a.loc)
}

func (a *arcade) copy() *arcade {
	out := &arcade{
		loc:        a.loc,
		robot:      a.robot.copy(),
		gameoutput: &gameOutput{},
	}
	return out
}
func (a *arcade) move(dir string) {
	switch dir {
	case north:
		a.loc.Y = a.loc.Y - 1
	case east:
		a.loc.X = a.loc.X + 1
	case south:
		a.loc.Y = a.loc.Y + 1
	case west:
		a.loc.X = a.loc.X - 1
	}
}

func run(initial *arcade) {
	q := internal.NewQueue[*arcade]()
	q.Enqueue(initial)
	visited := map[internal.Point]struct{}{}
	items := map[string]internal.Point{}
	grid := internal.NewGridV2()
	for !q.Empty() {
		cur := q.Dequeue()
		grid.Set(cur.loc, "#")
		r := cur.getRoom()
		for _, item := range r.items {
			items[item] = cur.loc
		}
		fmt.Println(r.Describe())
		//		fmt.Println(r.name)
		for _, door := range r.doors {
			//			fmt.Println("door", door)
			next := cur.copy()
			next.move(door)
			if _, ok := visited[next.loc]; ok {
				continue
			}
			//			fmt.Println("door", door)
			next.input = cmd(door)
			q.Enqueue(next)
		}
		visited[cur.loc] = struct{}{}
		//		fmt.Println(cur.loc)
		fmt.Println(grid)
	}
	for k, v := range items {
		fmt.Println(k, v)
	}
}

/*
== Observatory ==
There are a few telescopes; they're all bolted down, though.

Doors here lead:
- north
- west

Items here:
- klein bottle

Command?
*/

type room struct {
	name  string
	desc  string
	pos   internal.Point
	doors []string
	items []string
}

func (r *room) Describe() string {
	var out strings.Builder
	out.WriteString(r.name)
	out.WriteString("\n")
	if len(strings.TrimSpace(r.desc)) > 0 {
		out.WriteString(r.desc)
		out.WriteString("\n")
	}
	if len(r.doors) > 0 {
		out.WriteString("Doors ")
		out.WriteString(fmt.Sprintf("%v\n", r.doors))
	}
	if len(r.items) > 0 {
		out.WriteString("Items ")
		out.WriteString(fmt.Sprintf("%v\n", r.items))
	}
	return out.String()
}

func (r *room) String() string {
	for _, c := range r.name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			return string(c)
		}
	}
	panic("bad name, " + r.name)
}

type item struct {
	name string
	pos  internal.Point
}

func cake() []int {
	return cmd(north, north, east, east, "take cake", west, west, south, south)
}
func ornament() []int {
	return cmd(east, "take ornament", west)
}
func hologram() []int {
	return cmd(east, east, "take hologram", west, west)
}

func cmd(commands ...string) []int {
	out := []int{}
	for _, c := range commands {
		split := strings.Split(c, "")
		for i := range split {
			out = append(out, int(split[i][0]))
		}
		out = append(out, int(byte('\n')))
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

func (r *robot) act(s *arcade) bool {
	//	fmt.Println("using", s.input, "for the robot")
	r.brain.inputs = s.input
	r.brain.run()
	out := r.brain.output
	return s.print(out)
}

type arcade struct {
	loc        internal.Point
	robot      *robot
	input      []int
	gameoutput *gameOutput
}

func reverse(dir string) string {
	switch dir {
	case north:
		return south
	case east:
		return west
	case south:
		return north
	case west:
		return east
	}
	panic("bad dir")
}

func (s *arcade) print(x int) bool {
	//	fmt.Printf(string(byte(x)))
	s.gameoutput.addChar(string(byte(x)))
	return s.gameoutput.state != done
}
func intSliceToString(in []int) string {
	out := ""
	for _, x := range in {
		out += string(byte(x))
	}
	return out
}

type gameOutput struct {
	name    string
	desc    string
	doors   string
	items   string
	state   int
	nlcount int

	data string
}

func (g *gameOutput) asRoom(pos internal.Point) *room {
	items := make([]string, 0)
	for _, i := range strings.Split(g.items, "\n") {
		clean := strings.TrimSpace(strings.TrimPrefix(string(i), "- "))
		if len(clean) == 0 {
			continue
		}
		items = append(items, clean)
	}
	doors := make([]string, 0)
	for _, d := range strings.Split(g.doors, "\n") {
		clean := strings.TrimSpace(strings.TrimPrefix(string(d), "- "))
		if len(clean) == 0 {
			continue
		}
		doors = append(doors, clean)
	}
	return &room{
		name:  strings.TrimSpace(g.name),
		desc:  strings.TrimSpace(g.desc),
		pos:   pos,
		doors: doors,
		items: items,
	}
}

func (g *gameOutput) String() string {
	var out strings.Builder
	out.WriteString("name: " + g.name + "\n")
	if len(g.desc) > 0 {
		out.WriteString("desc: " + g.desc + "\n")
	}
	if len(g.doors) > 0 {
		out.WriteString("doors:" + "\n")
		for _, d := range strings.Split(g.doors, "\n") {
			out.WriteString(fmt.Sprintf("- %s\n", d))
		}
	}
	if len(g.items) > 0 {
		out.WriteString("items:" + "\n")
		for _, d := range strings.Split(g.items, "\n") {
			out.WriteString(fmt.Sprintf("- %s\n", d))
		}
	}
	return out.String()
}

const (
	newlines = iota
	readName
	readDesc
	readDoorsHeader
	readDoors
	readItemsHeader
	readItems
	readCommand
	done
	takingOutput
)

func (g *gameOutput) addChar(x string) {
	g.data += string(x)
	if g.state != takingOutput && g.state != readCommand && strings.HasPrefix(strings.TrimSpace(g.data), "You take") {
		g.nlcount = 0
		g.state = takingOutput
	}
	if g.state != readCommand && strings.HasSuffix(strings.TrimSpace(g.data), "Command?") {
		g.state = readCommand
	}
	switch g.state {
	case newlines:
		if x == "\n" {
			g.nlcount += 1
		}
		if g.nlcount == 3 {
			g.state = readName
			g.nlcount = 0
			return
		}
	case readName:
		if x == "\n" {
			g.state = readDesc
			return
		}
		g.name += x
	case readDesc:
		if x == "\n" {
			g.nlcount += 1
		} else {
			g.nlcount = 0
		}
		if g.nlcount == 2 {
			g.state = readDoorsHeader
			g.nlcount = 0
			return
		}
		g.desc += x
	case readDoorsHeader:
		if x == "\n" {
			g.nlcount += 1
		}
		if g.nlcount == 1 {
			g.state = readDoors
			g.nlcount = 0
			return
		}
	case readDoors:
		if x == "\n" {
			g.nlcount += 1
		} else {
			g.nlcount = 0
		}
		if g.nlcount == 2 {
			g.state = readItemsHeader
			g.nlcount = 0
			return
		}
		g.doors += x
	case readItemsHeader:
		if x == "\n" {
			g.nlcount += 1
		}
		if g.nlcount == 1 {
			if strings.HasSuffix(g.data, "Command?\n") {
				g.state = done
				return
			}
			g.state = readItems
			g.nlcount = 0
			return
		}
	case readItems:
		if x == "\n" {
			g.nlcount += 1
		} else {
			g.nlcount = 0
		}
		if g.nlcount == 2 {
			g.state = readCommand
			g.nlcount = 0
			return
		}
		g.items += x
	case readCommand:
		if x == "\n" {
			g.state = done
		}
	case done:
		panic("reinsatniate the game output struct my guy")
	case takingOutput:
		if x == "\n" {
			g.nlcount += 1
		}
		if g.nlcount == 2 {
			g.state = readCommand
			g.nlcount = 0
			return
		}
	}
}

type robot struct {
	brain *opcodeComputer
}

func (r *robot) copy() *robot {
	out := &robot{
		brain: r.brain.copy(),
	}
	return out
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
	outputSet    bool
}

func (o *opcodeComputer) copy() *opcodeComputer {
	newdata := make(map[int]int)
	for k, v := range o.data {
		newdata[k] = v
	}
	out := &opcodeComputer{
		data:         newdata,
		cur:          o.cur,
		relativeBase: o.relativeBase,
		debug:        o.debug,
		silent:       o.silent,
	}
	return out
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
	o.outputSet = false
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
func (o *opcodeComputer) setRegisterUsingInput(reg int) bool {
	if o.curInput == len(o.inputs) {
		return false
	}
	o.data[reg] = o.inputs[o.curInput]
	o.curInput++
	return true
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
		o.outputSet = true
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
		if !o.setRegisterUsingInput(vals[0]) {
			return nil
		}
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
