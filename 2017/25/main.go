package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	startStateLine := strings.Split(lines[0], " ")
	startState := strings.TrimSuffix(startStateLine[3], ".")
	stepsLine := strings.Split(lines[1], " ")
	steps, _ := strconv.Atoi(stepsLine[5])
	lines = lines[2:]
	states := map[string]state{}
	for i := 0; i < len(lines)/9; i++ {
		state := parseLines(lines[i*9 : (i+1)*9])
		states[state.name] = state
	}
	t := newTape(startState, states)
	for i := 0; i < steps; i++ {
		t.run()
	}
	fmt.Println(t.checksum())
}

func parseLines(lines []string) state {
	nameline := strings.Split(lines[0], " ")
	name := strings.TrimSuffix(nameline[2], ":")
	wzvline := strings.Fields(lines[2])
	wzv, _ := strconv.Atoi(strings.TrimSuffix(wzvline[4], "."))
	movezline := strings.Fields(lines[3])
	movez := left
	if movezline[6] == "right." {
		movez = right
	}
	zerostateline := strings.Fields(lines[4])
	zs := strings.TrimSuffix(zerostateline[4], ".")
	wovline := strings.Fields(lines[6])
	wov, _ := strconv.Atoi(strings.TrimSuffix(wovline[4], "."))
	moveoline := strings.Fields(lines[7])
	moveo := left
	if moveoline[6] == "right." {
		moveo = right
	}
	onestateline := strings.Fields(lines[8])
	os := strings.TrimSuffix(onestateline[4], ".")

	return state{
		name:          name,
		writeZeroVal:  wzv,
		moveZero:      movez,
		nextZeroState: zs,
		writeOneVal:   wov,
		moveOne:       moveo,
		nextOneState:  os,
	}
}

type direction string

const (
	right direction = "right"
	left  direction = "left"
)

type state struct {
	name          string
	moveZero      direction
	moveOne       direction
	writeZeroVal  int
	writeOneVal   int
	nextZeroState string
	nextOneState  string
}

func (s *state) run(t *tape) {
	if t.curval() == 0 {
		t.write(s.writeZeroVal)
		t.move(s.moveZero)
		t.nextState(s.nextZeroState)
		return
	}
	if t.curval() == 1 {
		t.write(s.writeOneVal)
		t.move(s.moveOne)
		t.nextState(s.nextOneState)
	}
}

type tape struct {
	cur      int
	data     map[int]int
	states   map[string]state
	curstate state
}

func (t *tape) checksum() int {
	sum := 0
	for _, v := range t.data {
		if v == 1 {
			sum++
		}
	}
	return sum
}

func newTape(startState string, states map[string]state) *tape {
	return &tape{
		cur:      0,
		data:     make(map[int]int),
		states:   states,
		curstate: states[startState],
	}
}

func (t *tape) curval() int {
	return t.data[t.cur]
}
func (t *tape) write(v int) {
	t.data[t.cur] = v
}
func (t *tape) move(d direction) {
	if d == right {
		t.cur++
	}
	if d == left {
		t.cur--
	}
}
func (t *tape) nextState(s string) {
	t.curstate = t.states[s]
}
func (t *tape) run() {
	t.curstate.run(t)
}
