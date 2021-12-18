package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	/*
		The first floor contains a polonium generator, a thulium generator, a thulium-compatible microchip, a promethium generator, a ruthenium generator, a ruthenium-compatible microchip, a cobalt generator, and a cobalt-compatible microchip.
		The second floor contains a polonium-compatible microchip and a promethium-compatible microchip.
		The third floor contains nothing relevant.
		The fourth floor contains nothing relevant.
	*/
	lines := internal.ReadInput()
	mcs := map[string]*microchip{}
	gns := map[string]*generator{}
	for _, line := range lines {
		gens, chips := parseItems(line)
		for k, v := range gens {
			gns[k] = v
		}
		for k, v := range chips {
			mcs[k] = v
		}
	}
	floors := map[int]*floor{}
	for i := 1; i <= 4; i++ {
		floors[i] = newFloor(i)
		for _, v := range mcs {
			if v.floor == i {
				floors[i].addMicrochip(v)
			}
		}
		for _, v := range gns {
			if v.floor == i {
				floors[i].addGenerator(v)
			}
		}
	}
	ele := &elevator{floor: 1}
	bldg := &building{
		floors:   floors,
		elevator: ele,
	}
	fmt.Println(bldg)
}

func parseItems(line string) (map[string]*generator, map[string]*microchip) {
	words := strings.Split(line, " ")
	floor := 0
	switch {
	case strings.HasPrefix(line, "The first"):
		floor = 1
	case strings.HasPrefix(line, "The second"):
		floor = 2
	case strings.HasPrefix(line, "The third"):
		floor = 3
	case strings.HasPrefix(line, "The fourth"):
		floor = 4
	default:
		panic("unknown line: " + line)
	}
	microchips := map[string]*microchip{}
	generators := map[string]*generator{}
	for i, word := range words {
		word = strings.Trim(word, ",.")
		if word == "generator" {
			generators[words[i-1]] = newGenerator(words[i-1], floor)
		}
		if word == "microchip" {
			broken := strings.Split(words[i-1], "-")
			microchips[broken[0]] = newMicrochip(broken[0], floor)
		}
	}
	return generators, microchips
}

type building struct {
	floors   map[int]*floor
	elevator *elevator
}

func (b *building) String() string {
	out := fmt.Sprintf("Valid? **%v**\n", b.valid())
	keys := []int{}
	for i := range b.floors {
		keys = append(keys, i)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	for _, i := range keys {
		floor := b.floors[i]
		if i == b.elevator.floor {
			out += "ele->"
		}
		out += fmt.Sprintf("%v\n", floor)
	}
	return out
}
func (b *building) valid() bool {
	for _, floor := range b.floors {
		if !floor.valid() {
			return false
		}
	}
	return true
}
func (b *building) solve() {
	// for each microchip	
}

type floor struct {
	number     int
	microchips map[string]*microchip
	generators map[string]*generator
}

func newFloor(level int) *floor {
	return &floor{
		number:     level,
		microchips: make(map[string]*microchip),
		generators: make(map[string]*generator),
	}
}
func (f *floor) addMicrochip(m *microchip) {
	f.microchips[m.kind] = m
}
func (f *floor) removeMicrochip(kind string) *microchip {
	out := f.microchips[kind]
	delete(f.microchips, kind)
	return out
}
func (f *floor) addGenerator(g *generator) {
	f.generators[g.kind] = g
}
func (f *floor) removeGenerator(kind string) *generator {
	out := f.generators[kind]
	delete(f.generators, kind)
	return out
}

func (f *floor) valid() bool {
	// if a microchip is on a floor with a non-compatible generator, fail
	// for every microchip
	// look at every generator. if ther
	for _, m := range f.microchips {
		if _, ok := f.generators[m.kind]; ok {
			continue
		}
		for _, g := range f.generators {
			if g.kind != m.kind {
				return false
			}
		}
	}
	return true
}

func (f *floor) String() string {
	return fmt.Sprintf("(%d) 👾 %v ⚡️ %v", f.number, f.microchips, f.generators)
}

type elevator struct {
	floor int
}

type microchip struct {
	kind  string
	floor int
}

func newMicrochip(k string, f int) *microchip {
	return &microchip{kind: k, floor: f}
}

func (m *microchip) String() string {
	return fmt.Sprintf("(%d) M%s", m.floor, strings.ToUpper(string(m.kind[:2])))
}

type generator struct {
	kind  string
	floor int
}

func newGenerator(k string, f int) *generator {
	return &generator{
		kind:  k,
		floor: f,
	}
}
func (g *generator) String() string {
	return fmt.Sprintf("(%d) G%s", g.floor, strings.ToUpper(string(g.kind[:2])))
}

/*

simple solution
kafka if reasonable
consumer group

depends on how the logs get spread out -- if the

worker

* preserve ordering



*/