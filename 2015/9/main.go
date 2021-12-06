package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	nodeMap := make(nodes)

	for _, line := range lines {
		pieces := strings.Split(line, " = ")
		dist, _ := strconv.Atoi(pieces[1])
		nodeNames := strings.Split(pieces[0], " to ")
		left := nodeMap.getOrCreate(nodeNames[0])
		right := nodeMap.getOrCreate(nodeNames[1])
		edge := newEdge(dist, nodeNames[0], nodeNames[1])
		left.addEdge(edge)
		right.addEdge(edge)
	}
	keys := []string{}
	for k := range nodeMap {
		keys = append(keys, k)
	}
	out := allPermutations([]string{}, keys)
	shortest := 0
	//	shortestPerm := []string{}
	for _, perm := range out {
		permDistance := 0
		for i := 0; i < len(perm)-1; i++ {
			permDistance += nodeMap.distance(perm[i], perm[i+1])
		}
		if permDistance > shortest {
			shortest = permDistance
			//			shortestPerm = perm
		}
	}
	fmt.Println(shortest)
}

type nodes map[string]*node

func (n nodes) getOrCreate(name string) *node {
	n1, ok := n[name]
	if !ok {
		n1 = newNode(name)
		n[name] = n1
	}
	return n1
}

func (n nodes) distance(from, to string) int {
	node := n[from]
	for _, e := range node.edges {
		if (e.from == from && e.to == to) || (e.to == from && e.from == to) {
			return e.distance
		}
	}
	return 99999999
}

type edge struct {
	distance int
	from     string
	to       string
}

func newEdge(d int, from, to string) *edge {
	return &edge{d, from, to}
}

type node struct {
	city  string
	edges []*edge
}

func (n *node) String() string {
	out := []string{fmt.Sprintf("%s to:", n.city)}
	for _, e := range n.edges {
		if e.from == n.city {
			out = append(out, fmt.Sprintf("(%s: %d)", e.to, e.distance))
		}
		if e.to == n.city {
			out = append(out, fmt.Sprintf("(%s: %d)", e.from, e.distance))
		}
	}
	return strings.Join(out, " ")
}

func (n *node) addEdge(edge *edge) {
	n.edges = append(n.edges, edge)
}

func newNode(city string, edges ...*edge) *node {
	return &node{
		city:  city,
		edges: edges,
	}
}

func allPermutations(fixed, unfixed []string) [][]string {
	if len(unfixed) == 0 {
		swap := make([]string, len(fixed))
		copy(swap, fixed)
		return [][]string{swap}
	}

	out := [][]string{}
	for i := range unfixed {
		swap := make([]string, len(unfixed))
		copy(swap, unfixed)
		swap[0], swap[i] = swap[i], swap[0]
		out = append(out, allPermutations(append(fixed, swap[0]), swap[1:])...)
	}
	return out
}
