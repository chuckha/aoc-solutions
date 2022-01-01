package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	// key is left value
	ns := []*node{}
	for _, line := range lines {
		n := nodeFromLine(line)
		ns = append(ns, n)
	}
	out := longest([]*node{}, ns)
	fmt.Println(out)
	fmt.Println(nodes(out).strength())
}

// if i make every single possible connection...
func test() {
	/*
		for every node
			add a node
	*/
}

// give me a node and the new set with that node removed
// find the largest path given the current path and the node bag copy

type nodes []*node

func (n nodes) copy() nodes {
	out := make([]*node, len(n))
	for i := range out {
		out[i] = n[i].copy()
	}
	return out
}

func (n nodes) strength() int {
	sum := 0
	for _, v := range n {
		sum += v.strength()
	}
	return sum
}
func (n nodes) next(remaining []*node) (*node, bool) {
	if len(n) == 0 {
		// pick the 0 with the highest value
		strongest := 0
		var sn *node
		for _, no := range remaining {
			if !no.canStart() {
				continue
			}
			if no.strength() > strongest {
				strongest = no.strength()
				sn = no
			}
		}
		return sn, true
	}
	latest := n[len(n)-1]

	strongestVal := 0
	var strongest *node
	ok := false
	for _, node := range remaining {
		if !latest.canMatch(node) {
			continue
		}
		if node.strength() > strongestVal {
			ok = true
			strongest = node
			strongestVal = node.strength()
		}
	}
	return strongest, ok
}
func (n nodes) remove(node *node) nodes {
	out := n.copy()
	for i, no := range n {
		if no.equal(node) {
			return append(out[:i], out[i+1:]...)
		}
	}
	panic("big uh oh")
}
func (n nodes) canAdd(node *node) bool {
	if len(n) == 0 {
		return node.canStart()
	}
	latest := n[len(n)-1]
	///	fmt.Println("latest", latest, node)
	return latest.canMatch(node)
}

// 617 too low
// 3217 too high
func strongest(path []*node, remaining []*node) []*node {

	// if there are no more components, the bridge cannot continue
	if len(remaining) == 0 {
		return path
	}
	// if we cannot continue at all, just end here
	if _, ok := nodes(path).next(remaining); !ok {
		return path
	}
	strongestval := nodes(path).strength()
	strongestpath := path
	// check every node possiblity
	for _, node := range remaining {
		if !nodes(path).canAdd(node) {
			continue
		}
		newn := node.copy()
		// copy remaining
		newr := nodes(remaining).copy()
		// remove the node from the copy
		newr = newr.remove(node)
		// copy the path
		newp := nodes(path).copy()
		if len(newp) == 0 {
			newn.start()
		} else {
			newp[len(newp)-1].connect(newn)
		}
		// add it to the new path
		newp = append(newp, newn)
		// if we can add more, recurse
		bestPath := strongest(newp, newr)
		if nodes(bestPath).strength() > strongestval {
			strongestval = nodes(bestPath).strength()
			strongestpath = bestPath
			fmt.Println("new strongest", strongestval)
		}
	}
	return strongestpath
}

func longest(path []*node, remaining []*node) []*node {
	// if there are no more components, the bridge cannot continue
	if len(remaining) == 0 {
		return path
	}
	// if we cannot continue at all, just end here
	if _, ok := nodes(path).next(remaining); !ok {
		return path
	}
	longestPath := path
	// check every node possiblity
	for _, node := range remaining {
		if !nodes(path).canAdd(node) {
			continue
		}
		newn := node.copy()
		// copy remaining
		newr := nodes(remaining).copy()
		// remove the node from the copy
		newr = newr.remove(node)
		// copy the path
		newp := nodes(path).copy()
		if len(newp) == 0 {
			newn.start()
		} else {
			newp[len(newp)-1].connect(newn)
		}
		// add it to the new path
		newp = append(newp, newn)
		// if we can add more, recurse
		bestPath := longest(newp, newr)
		if len(bestPath) == len(longestPath) {
			if nodes(bestPath).strength() > nodes(longestPath).strength() {
				longestPath = bestPath
			}
		}
		if len(bestPath) > len(longestPath) {
			longestPath = bestPath
		}
	}
	return longestPath
}

type node struct {
	left         int
	right        int
	leftPlugged  bool
	rightPlugged bool
}

func (n node) strength() int {
	return n.left + n.right
}
func (n *node) canMatch(n2 *node) bool {
	// if left is open
	if !n.leftPlugged {
		if !n2.leftPlugged {
			if n.left == n2.left {
				return true
			}
		}
		if !n2.rightPlugged {
			if n.left == n2.right {
				return true
			}
		}
	}
	// if right is open
	if !n.rightPlugged {
		if !n2.leftPlugged {
			if n.right == n2.left {
				return true
			}
		}
		if !n2.rightPlugged {
			if n.right == n2.right {
				return true
			}
		}
	}
	return false
}
func (n *node) connect(n2 *node) bool {
	// if left is open
	if !n.leftPlugged {
		if !n2.leftPlugged {
			if n.left == n2.left {
				n.leftPlugged = true
				n2.leftPlugged = true
				return true
			}
		}
		if !n2.rightPlugged {
			if n.left == n2.right {
				n.leftPlugged = true
				n2.rightPlugged = true
				return true
			}
		}
	}
	// if right is open
	if !n.rightPlugged {
		if !n2.leftPlugged {
			if n.right == n2.left {
				n.rightPlugged = true
				n2.leftPlugged = true
				return true
			}
		}
		if !n2.rightPlugged {
			if n.right == n2.right {
				n.rightPlugged = true
				n2.rightPlugged = true
				return true
			}
		}
	}
	return false
}
func (n node) equal(n2 *node) bool {
	return n.left == n2.left && n.right == n2.right
}
func (n *node) canStart() bool {
	return (!n.leftPlugged && n.left == 0) || (!n.rightPlugged && n.right == 0)
}
func (n *node) start() {
	if n.left == 0 {
		n.leftPlugged = true
		return
	}
	if n.right == 0 {
		n.rightPlugged = true
	}
}
func (n *node) String() string {
	left := "<"
	if n.leftPlugged {
		left = "x"
	}
	right := ">"
	if n.rightPlugged {
		right = "x"
	}
	return fmt.Sprintf("%d %s-%s %d", n.left, left, right, n.right)
}
func (n *node) copy() *node {
	return &node{
		left:         n.left,
		right:        n.right,
		leftPlugged:  n.leftPlugged,
		rightPlugged: n.rightPlugged,
	}
}
func nodeFromLine(line string) *node {
	ports := strings.Split(line, "/")
	left, _ := strconv.Atoi(ports[0])
	right, _ := strconv.Atoi(ports[1])
	return &node{
		left:  left,
		right: right,
	}
}
