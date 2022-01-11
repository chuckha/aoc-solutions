package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	// generation 185 = 35211
	// every generation adds 194
	fmt.Println((50000000000-185)*194 + 35211)
	os.Exit(0)
	lines := internal.ReadInput()
	initialState := strings.Split(lines[0], " ")[2]
	rules := map[string]string{}
	rlz := []string{}
	for _, line := range lines[1:] {
		rlz = append(rlz, line)
		items := strings.Split(line, " => ")
		rules[items[0]] = items[1]
	}
	sort.Sort(sort.StringSlice(rlz))
	for _, r := range rlz {
		fmt.Println(r)
	}
	plants := newPots(rules, initialState)
	fmt.Println(0, plants)
	generations := 186
	for i := 1; i <= generations; i++ {
		plants = plants.generation()
		fmt.Println(i, plants, plants.count())
	}
	fmt.Println(plants.count())
}

type pots struct {
	data     map[int]string
	rules    rules
	min, max int
}
type rules struct {
	data map[string]string
}

func (r rules) lookup(in string) string {
	out, ok := r.data[in]
	if !ok {
		return "."
	}
	return out
}

func newPots(r map[string]string, initialState string) *pots {
	plants := &pots{
		rules: rules{r},
		data:  make(map[int]string),
		min:   0, max: len(initialState),
	}
	for i, c := range initialState {
		plants.set(i, string(c))
	}
	return plants
}
func (p *pots) generation() *pots {
	newp := newPots(p.rules.data, "")
	for i := p.min - 2; i <= p.max+2; i++ {
		newp.set(i, p.rules.lookup(p.surrounding(i)))
	}
	return newp
}
func (p *pots) String() string {
	var out strings.Builder
	for i := p.min - 2; i <= p.max+2; i++ {
		out.WriteString(p.at(i))
	}
	return out.String()
}
func (p *pots) set(i int, s string) {
	if s == "." {
		return
	}
	// otherwise try and update the min/max
	if i < p.min {
		p.min = i
	}
	if i > p.max {
		p.max = i
	}
	p.data[i] = s
}
func (p *pots) at(i int) string {
	item, ok := p.data[i]
	if !ok {
		return "."
	}
	return item
}
func (p *pots) surrounding(i int) string {
	return p.at(i-2) + p.at(i-1) + p.at(i) + p.at(i+1) + p.at(i+2)
}
func (p *pots) count() int {
	count := 0
	for k := range p.data {
		count += k
	}
	return count
}
