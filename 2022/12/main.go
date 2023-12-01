package main

import (
	"fmt"
	"sort"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()
	g := internal.NewGridV3[Height]()
	for j, line := range input {
		for i, ch := range line {
			pt := internal.Point{i, j}
			g.Set(pt, Height(ch))
		}
	}
	start := find(g, 'S')
	// replace the s with an 'a' because we have to start from every 'a' position
	g.Set(start, 'a')
	end := find(g, 'E')
	starts := findAll(g, 'a')
	startCosts := startCosts{}
	for _, start := range starts {
		costs := g.Dijkstra(start, canReach, cost)
		startCosts = append(startCosts, startCost{start: start, cost: costs[end]})
	}
	sort.Sort(startCosts)
	fmt.Println(startCosts[0], startCosts[1])
}

type startCosts []startCost

func (s startCosts) Len() int           { return len(s) }
func (s startCosts) Less(i, j int) bool { return s[i].cost < s[j].cost }
func (s startCosts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type startCost struct {
	start internal.Point
	cost  int
}

func cost(a, b Height) int {
	return 1
}

func canReach(a, b Height) bool {
	return int(normalize(a))-int(normalize(b)) >= -1
}

func normalize(a Height) Height {
	switch a {
	case 'S':
		return 'a'
	case 'E':
		return 'z'
	default:
		return a
	}
}

func find(g *internal.GridV3[Height], i Height) internal.Point {
	p := internal.Point{}
	g.EachV2(func(pt internal.Point, t Height) {
		if t == i {
			p = pt
		}
	})
	return p
}

func findAll(g *internal.GridV3[Height], i Height) []internal.Point {
	out := []internal.Point{}
	g.EachV2(func(pt internal.Point, t Height) {
		if t == i {
			out = append(out, pt)
		}
	})
	return out
}

type Height byte

func (h Height) String() string { return string(h) }

type Node struct {
	H     string
	Point internal.Point
}

func (n *Node) String() string {
	return n.H
}
