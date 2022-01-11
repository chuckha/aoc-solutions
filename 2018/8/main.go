package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	vals := lines[0]
	d := []int{}
	for _, s := range strings.Split(vals, " ") {
		n, _ := strconv.Atoi(s)
		d = append(d, n)
	}
	parser := &data{
		input: d,
		cur:   0,
	}
	t := parse(parser)

	fmt.Println("part one:", t.metadataSum())
	fmt.Println("part two:", t.value())
}

type data struct {
	input []int
	cur   int
}

func (d *data) getNumChildren() int {
	numChildren := d.input[d.cur]
	d.cur++
	return numChildren
}
func (d *data) getNumMetadata() int {
	numMetadata := d.input[d.cur]
	d.cur++
	return numMetadata
}
func (d *data) getMetadata(n int) []int {
	metadata := d.input[d.cur : d.cur+n]
	d.cur += n
	return metadata
}

func parse(d *data) *tree {
	numChildren := d.getNumChildren()
	numMetadata := d.getNumMetadata()
	t := &tree{
		numChildren: numChildren,
		numMetadata: numMetadata,
		children:    make([]*tree, 0),
	}
	for i := 0; i < numChildren; i++ {
		t.children = append(t.children, parse(d))
	}
	t.metadata = d.getMetadata(numMetadata)
	return t
}

type tree struct {
	numChildren int
	numMetadata int
	children    []*tree
	metadata    []int
}

func (t *tree) metadataSum() int {
	sum := sum(t.metadata)
	for _, c := range t.children {
		sum += c.metadataSum()
	}
	return sum
}

func (t *tree) value() int {
	if len(t.children) == 0 {
		return t.metadataSum()
	}
	sum := 0
	for _, m := range t.metadata {
		if m == 0 {
			continue
		}
		if m >= len(t.children)+1 {
			continue
		}
		sum += t.children[m-1].value()
	}
	return sum
}

func sum(in []int) int {
	out := 0
	for _, n := range in {
		out += n
	}
	return out
}
