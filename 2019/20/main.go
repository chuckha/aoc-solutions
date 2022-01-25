package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

// 2552 is too low (part 2)

const (
	inner = "inner"
	outer = "outer"
)

func main() {
	lines := input.GetRawInput(2019, 20)
	g := graphFromInput(lines)
	last := solve(g)
	fmt.Println(last.steps)
	fmt.Println(printPath(last))
}

type graph struct {
	depth      int
	nodes      map[nodeKey]*node
	start, end internal.Point
}

type location struct {
	depth int
	nk    nodeKey
	path  []pathNode
	steps int
}
type pathNode struct {
	nk   nodeKey
	cost int
}

func (l location) copy() location {
	np := make([]pathNode, len(l.path))
	copy(np, l.path)
	return location{
		depth: l.depth,
		nk:    l.nk,
		path:  np,
		steps: l.steps,
	}
}

func printPath(final location) string {
	var out strings.Builder
	depth := 0
	cost := 0
	for i := 0; i < len(final.path)-1; i++ {
		out.WriteString(fmt.Sprintf(" [%d](%d) %s -> ", cost, depth, final.path[i].nk.name))
		if final.path[i].nk.name != final.path[i+1].nk.name {
			if final.path[i+1].nk.kind == inner {
				depth++
			}
			if final.path[i+1].nk.kind == outer {
				depth--
			}
		}
		cost += final.path[i].cost
	}
	return out.String()
}

func solve(maze *graph) location {
	q := internal.NewQueue[location]()
	startnk := nodeKey{name: "AA", kind: outer}
	start := location{depth: 0, nk: startnk, path: []pathNode{{nk: startnk, cost: 0}}}
	q.Enqueue(start)
	shortest := location{steps: math.MaxInt}
	answers := 0
	for !q.Empty() {
		cur := q.Dequeue()
		//		fmt.Println(cur, len(q.Internal()))
		if cur.depth == 0 && cur.nk.name == "ZZ" {
			answers++
			fmt.Println("found solution for", cur.steps)
			if cur.steps < shortest.steps {
				shortest = cur
			}
			if answers > 10 {
				return shortest
			}
			continue
		}
		// ignore depths over 50
		if cur.depth > 50 {
			continue
		}
		// get the node from the maze
		node := maze.nodes[cur.nk]
		for edge, cost := range node.edges {
			//			fmt.Println("?", edge)
			if cur.depth > 0 {
				switch edge.name {
				case "AA", "ZZ":
					continue
				}
			}
			//			fmt.Println("!", edge)
			next := cur.copy()
			next.path = append(next.path, pathNode{nk: edge, cost: cost})
			next.nk.name = edge.name
			next.steps = cur.steps + cost + 1
			switch edge.kind {
			case inner:
				next.depth = cur.depth + 1
				next.nk.kind = outer
			case outer:
				if cur.depth == 0 {
					if edge.name == "ZZ" {
						next.steps -= 1 // no recursing at the end
						q.Enqueue(next)
					}
					continue
				}
				next.depth = cur.depth - 1
				next.nk.kind = inner
			}

			q.Enqueue(next)
		}
	}
	return shortest
}

func (g *graph) deeper() *graph {
	newg := newGraph(g.start, g.end)
	for k, v := range g.nodes {
		newg.nodes[k] = v.copy()
	}
	newg.depth = g.depth + 1
	return newg
}

func newGraph(start, end internal.Point) *graph {
	return &graph{
		nodes: make(map[nodeKey]*node),
		start: start,
		end:   end,
	}
}

type nodeKey struct {
	name string
	kind string
}

type node struct {
	key   nodeKey
	edges map[nodeKey]int
}

func (n *node) copy() *node {
	out := newNode(n.key)
	for k, v := range n.edges {
		out.edges[k] = v
	}
	return out
}

func newNode(nk nodeKey) *node {
	return &node{
		key:   nk,
		edges: make(map[nodeKey]int),
	}
}

func graphFromInput(lines []string) *graph {
	// get a list of named nodes and points
	grid := internal.NewGridV2FromLines(lines)
	fmt.Println(grid)
	grid.DefaultChar = " "
	// read the graph top to bottom, left to right.
	letters := map[internal.Point]struct{}{}
	innerPortals := map[internal.Point]string{}
	outerPortals := map[internal.Point]string{}
	start := internal.Point{}
	end := internal.Point{}
	for j := grid.Min.Y; j <= grid.Max.Y; j++ {
		for i := grid.Min.X; i <= grid.Max.X; i++ {
			p := internal.Point{i, j}
			letter := grid.At(p)
			right := grid.At(p.Right())
			down := grid.At(p.Down())
			if _, ok := letters[p]; ok {
				continue
			}
			if letter >= "A" && letter <= "Z" {
				letters[p] = struct{}{}
				if right >= "A" && right <= "Z" {
					letters[p.Right()] = struct{}{}
					label := letter + right
					if grid.At(p.Left()) == "." {
						if p.Right().Right().X > grid.Max.X {
							outerPortals[p.Left()] = label
						} else {
							innerPortals[p.Left()] = label
						}
					}
					if grid.At(p.Right().Right()) == "." {
						if p.Left().X < grid.Min.X {
							outerPortals[p.Right().Right()] = label
						} else {
							innerPortals[p.Right().Right()] = label
						}
					}
				}
				if down >= "A" && down <= "Z" {
					letters[p.Down()] = struct{}{}
					label := letter + down
					if grid.At(p.Up()) == "." {
						if p.Down().Down().Y > grid.Max.Y {
							outerPortals[p.Up()] = label
						} else {
							innerPortals[p.Up()] = label
						}
					}
					if grid.At(p.Down().Down()) == "." {
						if p.Up().Y < grid.Min.Y {
							outerPortals[p.Down().Down()] = label
						} else {
							innerPortals[p.Down().Down()] = label
						}
					}
				}
			}
		}
	}
	for k, v := range outerPortals {
		if v == "AA" {
			start = k
		}
		if v == "ZZ" {
			end = k
		}
	}
	//	delete(outerPortals, start)
	//	delete(outerPortals, end)
	g := newGraph(start, end)
	nk, n := buildNode(start, "AA", outer, grid, innerPortals, outerPortals)
	g.nodes[nk] = n
	nk, n = buildNode(end, "ZZ", outer, grid, innerPortals, outerPortals)
	g.nodes[nk] = n
	for point, portal := range innerPortals {
		nk, n := buildNode(point, portal, inner, grid, innerPortals, outerPortals)
		g.nodes[nk] = n
	}
	for point, portal := range outerPortals {
		nk, n := buildNode(point, portal, outer, grid, innerPortals, outerPortals)
		g.nodes[nk] = n
	}
	return g
}

func buildNode(start internal.Point, name, kind string, grid *internal.GridV2, innerPortals, outerPortals map[internal.Point]string) (nodeKey, *node) {
	nk := nodeKey{name, kind}
	startNode := newNode(nk)
	costs := cost(start, grid)
	for p, name := range innerPortals {
		if p == start {
			continue
		}
		cost, ok := costs[p]
		if !ok {
			continue
		}
		startNode.edges[nodeKey{name, inner}] = cost
	}
	for p, name := range outerPortals {
		if p == start {
			continue
		}
		cost, ok := costs[p]
		if !ok {
			continue
		}
		startNode.edges[nodeKey{name, outer}] = cost
	}
	return nk, startNode
}

func cost(start internal.Point, data *internal.GridV2) map[internal.Point]int {
	costs := map[internal.Point]int{}
	costs[start] = 0
	visited := map[internal.Point]struct{}{}
	q := internal.NewQueue[internal.Point]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		if _, ok := visited[cur]; ok {
			continue
		}
		for _, p := range cur.Neighbors() {
			if data.At(p) == " " || data.At(p) == "#" || (data.At(p) >= "A" && data.At(p) <= "Z") {
				continue
			}
			cost, ok := costs[p]
			if !ok {
				cost = math.MaxInt
			}
			if costs[cur]+1 < cost {
				costs[p] = costs[cur] + 1
			}
			q.Enqueue(p)
		}
		visited[cur] = struct{}{}
	}
	return costs
}

// node is a name/inner name/outer outer -> inner + depth //// inner -> outer - depth

// func old(maze *maze) {
// 	costs := maze.dijkstra()
// 	for j := maze.raw.Min.Y; j <= maze.raw.Max.Y; j++ {
// 		for i := maze.raw.Min.X; i <= maze.raw.Max.X; i++ {
// 			if cost, ok := costs[internal.Point{i, j}]; ok {
// 				if cost >= 100 {
// 					fmt.Print(string(internal.Red("*")))
// 					continue
// 				}
// 				if cost >= 10 {
// 					fmt.Print(string(internal.Green(fmt.Sprintf("%d", cost%10))))
// 					continue
// 				}
// 				fmt.Printf("%d", cost)
// 				continue
// 			}
// 			fmt.Print("#")
// 		}
// 		fmt.Print("\n")
// 	}
// 	fmt.Println(costs[maze.end])
// 	// for _, n := range maze.nodes {
// 	// 	fmt.Println(n)
// 	// }
// }

// type label struct {
// 	name string
// 	pos  internal.Point
// 	done bool
// }

// func newMazeFromLines(lines []string) *maze {
// 	grid := internal.NewGridV2FromLines(lines)
// 	maze := newMaze(grid, 0)
// 	labels := []label{}
// 	for p, c := range grid.Data {
// 		if c != "." {
// 			continue
// 		}
// 		n := maze.GetOrCreateNode(p)
// 		for _, neighbor := range p.Neighbors() {
// 			nabe := grid.Data[neighbor]
// 			if nabe == "#" || nabe == " " || nabe == "" {
// 				continue
// 			}
// 			if nabe == "." {
// 				newNeighbor := maze.GetOrCreateNode(neighbor)
// 				n.addNeighbor(newNeighbor)
// 				continue
// 			}
// 			// otherwise it's a letter
// 			for _, n2 := range neighbor.Neighbors() {
// 				val := grid.Data[n2]
// 				if n2 == p || val == "." || val == "#" || val == " " || val == "" {
// 					continue
// 				}
// 				// it's the other letter
// 				label := maze.buildLabel(neighbor, n2, nabe, val)
// 				label.pos = p
// 				labels = append(labels, label)
// 				break
// 			}
// 		}
// 	}
// 	for i, label := range labels {
// 		if label.done {
// 			continue
// 		}
// 		// every label but AA and ZZ have a matching label
// 		for j, l2 := range labels {
// 			if i == j {
// 				continue
// 			}
// 			if label.name == l2.name {
// 				oneNode := maze.GetOrCreateNode(label.pos)
// 				twoNode := maze.GetOrCreateNode(l2.pos)
// 				oneNode.addNeighbor(twoNode)
// 				twoNode.addNeighbor(oneNode)
// 				label.done = true
// 				l2.done = true
// 				labels[i] = label
// 				labels[j] = l2
// 				continue
// 			}
// 		}
// 	}
// 	for _, label := range labels {
// 		if !label.done {
// 			if label.name == "AA" {
// 				maze.start = label.pos
// 				continue
// 			}
// 			maze.end = label.pos
// 		}
// 	}
// 	return maze
// }

// func (m *maze) dijkstra() map[internal.Point]int {
// 	costs := map[internal.Point]int{}
// 	for _, n := range m.nodes {
// 		costs[n.pos] = math.MaxInt
// 	}
// 	visited := map[internal.Point]bool{}
// 	costs[m.start] = 0
// 	q := internal.NewQueue[*Node]()
// 	q.Enqueue(m.nodes[m.start])
// 	for !q.Empty() {
// 		cur := q.Dequeue()
// 		if _, ok := visited[cur.pos]; ok {
// 			continue
// 		}
// 		for _, n := range cur.Neighbors {
// 			if costs[cur.pos]+1 < costs[n.pos] {
// 				costs[n.pos] = costs[cur.pos] + 1
// 			}
// 			q.Enqueue(n)
// 		}
// 		visited[cur.pos] = true
// 	}
// 	return costs
// }

// // p1 is always the first one found, p2 will be the second found
// func (m *maze) buildLabel(p1, p2 internal.Point, letter1, letter2 string) label {
// 	// the second letter is above the first one up
// 	if p1.Y-1 == p2.Y {
// 		return label{
// 			name: letter2 + letter1,
// 		}
// 	}
// 	// right
// 	if p1.X+1 == p2.X {
// 		return label{
// 			name: letter1 + letter2,
// 		}
// 	}
// 	// down
// 	if p1.Y+1 == p2.Y {
// 		return label{
// 			name: letter1 + letter2,
// 		}
// 	}
// 	// left
// 	if p1.X-1 == p2.X {
// 		return label{
// 			name: letter2 + letter1,
// 		}
// 	}
// 	fmt.Printf("%v, %v, %q, %q\n", p1, p2, letter1, letter2)
// 	panic("uh bad news")
// }

// type mazes struct {
// 	data map[int]*maze
// }

// type maze struct {
// 	raw        *internal.GridV2
// 	nodes      map[internal.Point]*Node
// 	start, end internal.Point // implicitly only at depth 0
// 	rawWarps   map[internal.Point]internal.Point
// 	depth      int
// }

// func newMaze(raw *internal.GridV2, depth int) *maze {
// 	return &maze{
// 		raw:      raw,
// 		nodes:    make(map[internal.Point]*Node),
// 		rawWarps: make(map[internal.Point]internal.Point),
// 		depth:    depth,
// 	}
// }

// func (m *maze) warps() map[internal.Point]internal.Point {
// 	out := make(map[internal.Point]internal.Point)
// 	if depth == 0 {

// 	}
// }

// func (m *maze) GetOrCreateNode(p internal.Point) *Node {
// 	if n, ok := m.nodes[p]; ok {
// 		return n
// 	}
// 	m.nodes[p] = NewNode(p)
// 	return m.nodes[p]
// }

// type Node struct {
// 	pos       internal.Point
// 	Neighbors map[internal.Point]*Node
// }

// func NewNode(p internal.Point) *Node {
// 	return &Node{
// 		pos: p,
// 		// all neighbors are 1 step away
// 		Neighbors: make(map[internal.Point]*Node),
// 	}
// }

// func (n *Node) addNeighbor(n2 *Node) {
// 	n.Neighbors[n2.pos] = n2
// }

// func (n *Node) String() string {
// 	var out strings.Builder
// 	out.WriteString(n.pos.String() + "\n")
// 	for p := range n.Neighbors {
// 		out.WriteString("\t" + p.String() + "\n")
// 	}
// 	return out.String()
// }
