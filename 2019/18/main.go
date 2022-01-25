package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

// NEW IDEA

// turn it into a graph.
// figure out state to represent a path through the graph (steps & keys in posession & where in the graph the state is)

// part 1: 5712 too high
// part 1: 5538 too high
// part 1: 5540 too high
// part 1: 5534 too high
// part 2: 2504 too high

func main() {
	lines := input.GetInput(2019, 18)
	part2(lines)
	m := newMazeFromLines(lines)
	g := newGraph()

	for p, k := range m.keys() {
		costs, _ := m.movementCost(p)
		n := newNode(k, p)
		for p2, k2 := range m.keys() {
			if p2 == p {
				continue
			}
			g.addNode(n)
			n.addEdge(k2, costs[p2])
		}
	}

	for i, p := range m.players {
		costs, _ := m.movementCost(p)
		n := newNode(fmt.Sprintf("%d", i), p)
		g.addNode(n)
		for p, k := range m.keys() {
			n.addEdge(k, costs[p])
		}
		gs := m.searchMaze2(p)
		// add constraints
		for key, gates := range gs.keys {
			if len(gates) == 0 {
				continue
			}
			g.addConstraints(key, gates...)
		}
	}

	//	fmt.Println(g)
	fmt.Println(boundedPath(g))
	//	shortestPath(g, initState)
	//	fmt.Println(g.StringWithState(initState))
	//	fmt.Println(m)
	//	c, c2 := thing(m)
	//	fmt.Println(lowerBounds(m, c, c2))
	//	tryingPQ(m)
	//tryingStack(m)
	//	costs := thing(m)
	//	fmt.Println(lowerBounds(m, costs))
}

func part2(lines []string) {
	// find the line
	lineIdx := 0
	idx := 0
	for i, line := range lines {
		lineIdx = i
		idx = strings.Index(line, ".@.")
		if idx != -1 {
			break
		}
	}
	lines[lineIdx-1] = (lines[lineIdx-1][:idx] + "@#@" + lines[lineIdx-1][idx+3:])
	lines[lineIdx] = lines[lineIdx][:idx] + "###" + lines[lineIdx][idx+3:]
	lines[lineIdx+1] = lines[lineIdx+1][:idx] + "@#@" + lines[lineIdx+1][idx+3:]
}

func boundedPath(g *graph) state {
	items := make(internal.FastPQ[state], 0)
	seen := map[string]int{}
	heap.Init(&items)
	players := make([]string, 0)
	for i := 0; i < 4; i++ {
		players = append(players, fmt.Sprintf("%d", i))
	}
	initState := state{
		players:       players,
		keysCollected: make(map[string]int),
		steps:         0,
	}
	heap.Push(&items, &internal.FastItem[state]{
		Value:    initState,
		Priority: 0,
		Found:    len(initState.keysCollected),
	})
	bound := math.MaxInt
	shortests := state{}
	for items.Len() > 0 {
		//		fmt.Println("queue size:", items.Len())
		cur := items.Pop().(*internal.FastItem[state]).Value
		if cur.steps > bound {
			continue
		}
		//		fmt.Println(cur)

		if g.solved(cur) {
			fmt.Println("solved", cur, "items remaining", items.Len())
			if cur.steps < bound {
				bound = cur.steps
				shortests = cur
			}
			continue
		}
		//solved 83: c->a->i->f->g->d->b->h->e items remaining 16
		for player, nodes := range g.accessibleKeys(cur) {
			for node, dist := range nodes {
				newState := cur.copy()
				newState = newState.collectKey(player, node, dist)
				if newState.steps >= bound {
					continue
				}
				stateKey := newState.key()
				if cost, ok := seen[stateKey]; ok {
					if newState.steps >= cost {
						//					fmt.Println("throwing away", stateKey, newState)
						continue
					}
				}
				seen[stateKey] = newState.steps
				heap.Push(&items, &internal.FastItem[state]{
					Value:    newState,
					Priority: newState.steps,
					Found:    len(newState.keysCollected),
				})
			}
		}
		// for each edge
		// make a new state with that edge in it
	}
	return shortests
}

type state struct {
	players       []string
	steps         int
	keysCollected map[string]int
}

func (s state) key() string {
	ks := make([]string, len(s.keysCollected))
	i := 0
	for k := range s.keysCollected {
		ks[i] = k
		i++
	}
	if len(ks) == 1 {
		return ks[0]
	}
	sort.Sort(sort.StringSlice(ks[1 : len(ks)-1]))
	return strings.Join(ks, "")
}

func (s state) collectKey(playerIdx int, key string, dist int) state {
	s.players[playerIdx] = key
	s.steps += dist
	s.keysCollected[key] = len(s.keysCollected)
	return s
}

func (s state) String() string {
	out := make([]string, len(s.keysCollected))
	for k, v := range s.keysCollected {
		out[v] = k
	}
	return fmt.Sprintf("%v: %v", s.steps, strings.Join(out, "->"))
}

func (s state) copy() state {
	kc := make(map[string]int)
	for k, v := range s.keysCollected {
		kc[k] = v
	}
	players := make([]string, len(s.players))
	copy(players, s.players)
	return state{
		players:       players,
		steps:         s.steps,
		keysCollected: kc,
	}
}

type graph struct {
	nodes       map[string]*node
	constraints map[string]map[string]struct{}
}

func (g *graph) solved(s state) bool {
	return len(s.keysCollected) == len(g.nodes)-len(s.players)
}

func newGraph() *graph {
	return &graph{
		nodes: make(map[string]*node),
		// node => required keys
		constraints: make(map[string]map[string]struct{}),
	}
}

func (g *graph) accessibleKeys(s state) map[int]map[string]int {
	out := make(map[int]map[string]int)
	for i, player := range s.players {
		// range over all keys that the state can get to
		for node, dist := range g.nodes[player].edges {
			if dist == math.MaxInt {
				continue
			}
			// skip the key if we've already collected it
			if _, ok := s.keysCollected[node]; ok {
				continue
			}
			// look at the node's constraints
			hasAllConstraints := true
			for c := range g.constraints[node] {
				// if the state hasn't collected the required then it cannot access the node
				if _, ok := s.keysCollected[c]; !ok {
					hasAllConstraints = false
					break
				}
			}
			if hasAllConstraints {
				if _, ok := out[i]; !ok {
					out[i] = make(map[string]int)
				}
				out[i][node] = dist
			}
		}
	}
	return out
}

func (g *graph) addConstraints(key string, behind ...string) {
	if _, ok := g.constraints[key]; !ok {
		g.constraints[key] = make(map[string]struct{}, 0)
	}
	for _, b := range behind {
		g.constraints[key][strings.ToLower(b)] = struct{}{}
	}
}

func (g *graph) String() string {
	var out strings.Builder
	for k, n := range g.nodes {
		out.WriteString(fmt.Sprintf("%s -> \n", k))
		for e, c := range n.edges {
			if gates, ok := g.constraints[e]; ok {
				gatelist := make([]string, 0)
				for gate := range gates {
					gatelist = append(gatelist, gate)
				}
				out.WriteString(fmt.Sprintf("\t %s (%d) (GATED %v)\n", e, c, gatelist))
				continue
			}
			out.WriteString(fmt.Sprintf("\t %s (%d)\n", e, c))
		}
	}
	return out.String()
}

func (g *graph) StringWithState(s state) string {
	var out strings.Builder
	for k, n := range g.nodes {
		out.WriteString(fmt.Sprintf("%s -> \n", k))

		for e, c := range n.edges {
			if gates, ok := g.constraints[e]; ok {
				gatelist := make([]string, 0)
				for gate := range gates {
					gatelist = append(gatelist, gate)
				}
				out.WriteString(fmt.Sprintf("\t %s (%d) (GATED %v)\n", e, c, gatelist))
				continue
			}
			out.WriteString(fmt.Sprintf("\t %s (%d)\n", e, c))
		}
	}
	return out.String()
}

func (g *graph) addNode(n *node) {
	g.nodes[n.name] = n
}

type node struct {
	name  string
	pos   internal.Point
	edges map[string]int
}

func newNode(name string, pos internal.Point) *node {
	return &node{
		name:  name,
		pos:   pos,
		edges: make(map[string]int),
	}
}

func (n *node) addEdge(name string, cost int) {
	n.edges[name] = cost
}

// returns keys needed to get key k
func (g gameState) keysRequiredForKey(k string) map[string]struct{} {
	requiredKeys := map[string]struct{}{}
	doors := g.gatedKeys[k]
	for _, door := range doors {
		key := strings.ToLower(door)
		accessible := false
		for _, k := range g.accessibleKeys {
			if k == key {
				accessible = true
			}
		}
		requiredKeys[key] = struct{}{}
		if accessible {
			continue
		}
		for k := range g.keysRequiredForKey(key) {
			requiredKeys[k] = struct{}{}
		}
	}
	return requiredKeys
}

func (g gameState) deepestKey() (string, []string) {
	if len(g.gatedKeys) == 0 {
		return g.accessibleKeys[0], []string{}
	}
	big := []string{}
	key := ""
	for k, behind := range g.gatedKeys {
		if len(behind) > len(big) {
			big = behind
			key = k
		}
	}
	return key, big
}

// func debugCosts(m *maze) {
// 	costs, _ := thing(m)
// 	for k, v := range costs {
// 		fmt.Println("=========", k, "=========")
// 		for j := m.walls.Min.Y; j <= m.walls.Max.Y; j++ {
// 			for i := m.walls.Min.X; i <= m.walls.Max.X; i++ {
// 				if i, ok := v[internal.Point{i, j}]; ok {
// 					fmt.Printf("%02d", i)
// 					continue
// 				}
// 				fmt.Print("##")
// 			}
// 			fmt.Println()
// 		}
// 		fmt.Println()
// 	}
// }

// if i don't want shortest...
// shortest is not right
// shortest should take into account gates
// but this still isn't right...
// yields lower bounds higher than the shortest...
// func lowerBounds(m *maze, costs map[internal.Point]map[internal.Point]int, cameFroms map[internal.Point]map[internal.Point]internal.Point) int {
// 	lol := m.copy()
// 	cost := 0
// 	for !lol.solved() {
// 		state := lol.getGameState()
// 		shortest := math.MaxInt
// 		key := ""
// 		keypoints := lol.pointsOfKeys()
// 		for _, k2 := range state.accessibleKeys {
// 			p := keypoints[k2]
// 			// for p1, k := range lol.objects.Data {
// 			// 	if k < "a" || k > "z" {
// 			// 		continue
// 			// 	}
// 			if costs[p][lol.player] < shortest {
// 				shortest = costs[p][lol.player]
// 				key = k2
// 			}
// 			//			}
// 		}
// 		if key == "" {
// 			panic("empty key")
// 		}
// 		cost += lol.collectKey(key, costs[lol.player], cameFroms[lol.player])
// 	}
// 	return cost
// }

// func fuckitLowerBound(m *maze, costs map[internal.Point]map[internal.Point]int, cameFroms map[internal.Point]map[internal.Point]internal.Point) int {
// 	lol := m.copy()
// 	cost := 0
// 	for !lol.solved() {
// 		state := lol.getGameState()
// 		key := ""
// 		for _, k2 := range state.accessibleKeys {
// 			key = k2
// 			break
// 		}
// 		if key == "" {
// 			panic("empty key")
// 		}
// 		cost += lol.collectKey(key, costs[lol.player], cameFroms[lol.player])
// 	}
// 	return cost
// }

// func lowerBoundsWithoutGates(m *maze, costs map[internal.Point]map[internal.Point]int, cameFroms map[internal.Point]map[internal.Point]internal.Point) int {
// 	lol := m.copy()
// 	cost := 0
// 	for !lol.solved() {
// 		shortest := math.MaxInt
// 		key := ""
// 		for p1, k := range lol.objects.Data {
// 			if k < "a" || k > "z" {
// 				continue
// 			}
// 			if costs[p1][lol.player] < shortest {
// 				shortest = costs[p1][lol.player]
// 				key = k
// 			}
// 			//			}
// 		}
// 		if key == "" {
// 			panic("empty key")
// 		}
// 		cost += lol.collectKey(key, costs[lol.player], cameFroms[lol.player])
// 	}
// 	return cost
// }

func (m *maze) pointsOfDoors() map[string]internal.Point {
	out := map[string]internal.Point{}
	for p, v := range m.objects.Data {
		if v >= "A" && v <= "Z" {
			out[v] = p
		}
	}
	return out

}
func (m *maze) pointsOfKeys() map[string]internal.Point {
	out := map[string]internal.Point{}
	for p, v := range m.objects.Data {
		if v >= "a" && v <= "z" {
			out[v] = p
		}
	}
	return out
}

// func thing(m *maze) (map[internal.Point]map[internal.Point]int, map[internal.Point]map[internal.Point]internal.Point) {
// 	costs := map[internal.Point]map[internal.Point]int{}
// 	cameFroms := map[internal.Point]map[internal.Point]internal.Point{}
// 	for p := range m.objects.Data {
// 		c, cf := m.movementCost(p)
// 		costs[p] = c
// 		cameFroms[p] = cf
// 	}
// 	costs[m.player], cameFroms[m.player] = m.movementCost(m.player)
// 	return costs, cameFroms
// }

func (m *maze) keys() map[internal.Point]string {
	out := map[internal.Point]string{}
	for p, k := range m.objects.Data {
		if k < "a" || k > "z" {
			continue
		}
		out[p] = k
	}
	return out
}

// a bad solution of picking the closest key: 6118
// an even worse solution is always picking the longest key: 7890

// func trying(initial *maze) {
// 	costs := thing(initial)
// 	bound := 6128
// 	mq := internal.NewQueue[*maze]()
// 	mq.Enqueue(initial)
// 	for !mq.Empty() {
// 		fmt.Println("Queue", len(mq.Internal()))
// 		cur := mq.Dequeue()
// 		if cur.solved() {
// 			fmt.Println("found a solution for", cur.totalMovement)
// 			if cur.totalMovement < bound {
// 				bound = cur.totalMovement
// 			}
// 			continue
// 		}
// 		gs := cur.getGameState()
// 		for _, key := range gs.accessibleKeys {
// 			mz := cur.copy()
// 			mz.collectKey(key)
// 			lb := lowerBounds(mz, costs)
// 			//			fmt.Println(lb, "vs", bound)
// 			if lb < bound {
// 				mq.Enqueue(mz)
// 			}
// 		}
// 	}
// 	fmt.Println("lowest", bound)
// }

// func tryingStack(initial *maze) {
// 	costs := thing(initial)
// 	bound := 5538
// 	mq := internal.NewStack[*maze]()
// 	mq.Push(initial)
// 	for !mq.Empty() {
// 		cur := mq.Pop()
// 		if cur.solved() {
// 			//			fmt.Println("found a solution for", cur.totalMovement)
// 			fmt.Println(cur.totalMovement, cur.keysCollected)
// 			fmt.Println("Queue", mq.Depth())
// 			if cur.totalMovement < bound {
// 				bound = cur.totalMovement
// 			}
// 			continue
// 		}
// 		gs := cur.getGameState()
// 		for _, key := range gs.accessibleKeys {
// 			mz := cur.copy()
// 			mz.collectKey(key)
// 			lb := lowerBounds(mz, costs)
// 			//			fmt.Println(mz.keysCollected)
// 			//			fmt.Println(lb+mz.totalMovement, "vs", bound)
// 			if (lb + mz.totalMovement) < bound {
// 				mq.Push(mz)
// 			}
// 		}
// 	}
// 	fmt.Println("lowest", bound)
// }

// func tryingPQ(initial *maze) {
// 	costs, cameFroms := thing(initial)
// 	bound := 5538
// 	bound = math.MaxInt
// 	path := []string{}
// 	mq := make(internal.FastPQ[*maze], 0)
// 	heap.Init(&mq)
// 	fi := &internal.FastItem[*maze]{
// 		Value:    initial,
// 		Priority: fuckitLowerBound(initial, costs, cameFroms),
// 		Found:    len(initial.keysCollected),
// 	}
// 	//	seen := map[string]bool{}
// 	pruned := 0
// 	heap.Push(&mq, fi)
// 	for mq.Len() > 0 {
// 		fmt.Println("Queue", mq.Len())
// 		cur := heap.Pop(&mq).(*internal.FastItem[*maze]).Value
// 		//		fmt.Println(cur.keysCollected, cur.totalMovement+lowerBounds(cur, costs, cameFroms))
// 		if cur.solved() {
// 			fmt.Println("found a solution for", cur.totalMovement)
// 			//fmt.Println(cur.totalMovement, cur.keysCollected)
// 			if cur.totalMovement < bound {
// 				bound = cur.totalMovement
// 				path = cur.keysCollected
// 			}
// 			continue
// 		}
// 		gs := cur.getGameState()
// 		for _, key := range gs.accessibleKeys {
// 			mz := cur.copy()
// 			//			fmt.Println("before", mz.keysCollected)
// 			mz.collectKey(key, costs[mz.player], cameFroms[mz.player])
// 			//			fmt.Println("after", mz.keysCollected)
// 			lb := fuckitLowerBound(mz, costs, cameFroms)
// 			//			fmt.Println(mz.keysCollected)
// 			//			fmt.Println(lb+mz.totalMovement, "vs", bound)
// 			fmt.Println(mz.keysCollected, lb+mz.totalMovement, "<", bound)
// 			if (lb + mz.totalMovement) < bound {
// 				fi := &internal.FastItem[*maze]{
// 					Value:    mz,
// 					Priority: lb + mz.totalMovement,
// 					Found:    len(mz.keysCollected),
// 				}
// 				// seenKey := strings.Join(mz.keysCollected, "")
// 				// if _, ok := seen[seenKey]; !ok {
// 				//					seen[seenKey] = true
// 				heap.Push(&mq, fi)
// 				// } else {
// 				// 	pruned++
// 				// }
// 			}
// 		}
// 	}
// 	fmt.Println("pruned", pruned)
// 	fmt.Println("lowest", bound)
// 	fmt.Println("path", path)

// }

type maze struct {
	walls         *internal.GridV2
	objects       *internal.GridV2
	players       []internal.Point
	totalMovement int
	keysCollected []string
}

func (m *maze) solved() bool {
	for _, v := range m.objects.Data {
		if v >= "a" && v <= "z" {
			return false
		}
	}
	return true
}

func (m *maze) copy() *maze {
	collected := make([]string, len(m.keysCollected))
	copy(collected, m.keysCollected)
	players := make([]internal.Point, len(m.players))
	copy(players, m.players)
	return &maze{
		walls:         m.walls.Copy(),
		objects:       m.objects.Copy(),
		players:       players,
		totalMovement: m.totalMovement,
		keysCollected: collected,
	}
}

func (m *maze) printCosts(costs map[internal.Point]int) string {
	var out strings.Builder
	for j := m.walls.Min.Y; j <= m.walls.Max.Y; j++ {
		for i := m.walls.Min.X; i <= m.walls.Max.X; i++ {
			p := internal.Point{i, j}
			if _, ok := m.walls.Data[p]; ok {
				out.WriteString("###")
				continue
			}
			out.WriteString(fmt.Sprintf("%03d", costs[p]))
		}
		out.WriteString("\n")
	}
	return out.String()

}

func (m *maze) movementCost(start internal.Point) (map[internal.Point]int, map[internal.Point]internal.Point) {
	costs := map[internal.Point]int{}
	cameFrom := map[internal.Point]internal.Point{
		start: {-1, -1},
	}
	for j := m.walls.Min.Y; j <= m.walls.Max.Y; j++ {
		for i := m.walls.Min.X; i <= m.walls.Max.X; i++ {
			if m.walls.At(internal.Point{i, j}) == "." {
				costs[internal.Point{i, j}] = math.MaxInt
			}
		}
	}
	visited := map[internal.Point]bool{}
	q := internal.NewQueue[internal.Point]()
	costs[start] = 0
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		if visited[cur] {
			continue
		}
		for _, n := range cur.Neighbors() {
			if visited[n] {
				continue
			}
			if m.walls.At(n) == "#" {
				continue
			}
			if costs[cur]+1 < costs[n] {
				costs[n] = costs[cur] + 1
				cameFrom[n] = cur
			}
			q.Enqueue(n)
		}
		visited[cur] = true
	}
	return costs, cameFrom
}

func (m *maze) String() string {
	var out strings.Builder
	for j := m.walls.Min.Y; j <= m.walls.Max.Y; j++ {
		for i := m.walls.Min.X; i <= m.walls.Max.X; i++ {
			p := internal.Point{i, j}
			printedPlayer := false
			for _, p2 := range m.players {
				if p == p2 {
					out.WriteString("@")
					printedPlayer = true
				}
			}
			if printedPlayer {
				continue
			}
			if _, ok := m.walls.Data[p]; ok {
				out.WriteString("#")
				continue
			}
			if item, ok := m.objects.Data[p]; ok {
				out.WriteString(item)
				continue
			}
			out.WriteString(".")
		}
		out.WriteString("\n")
	}
	return out.String()
}

func newMazeFromLines(lines []string) *maze {
	g := internal.NewGridV2FromLines(lines)
	walls := internal.NewGridV2()
	m := &maze{}
	for k, v := range g.Data {
		if v == "@" {
			m.players = append(m.players, k)
		}
		if v == "#" {
			walls.Set(k, v)
		}
	}
	objects := internal.NewGridV2()
	for k, v := range g.Data {
		if v != "#" && v != "." && v != "@" {
			objects.Set(k, v)
		}
	}
	m.objects = objects
	m.walls = walls
	return m
}

type gameState struct {
	accessibleKeys  []string
	accessibleDoors []string
	gatedKeys       map[string][]string
	gatedDoors      map[string][]string
}
type gameState2 struct {
	keys map[string][]string
}

func (g gameState) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("Accessible keys: %v\n", g.accessibleKeys))
	out.WriteString(fmt.Sprintf("Accessible doors: %v\n", g.accessibleDoors))
	for k, v := range g.gatedKeys {
		out.WriteString(fmt.Sprintf("%s is gated behind %v\n", k, v))
	}
	for k, v := range g.gatedDoors {
		out.WriteString(fmt.Sprintf("%s is gated behind %v\n", k, v))
	}
	return out.String()
}

// func (m *maze) collectKey(k string, costs map[internal.Point]int, cameFrom map[internal.Point]internal.Point) int {
// 	//	_, cameFrom := m.movementCost(m.player)
// 	keyPoints := m.pointsOfKeys()
// 	doorPoints := m.pointsOfDoors()
// 	keyPoint := keyPoints[k]
// 	// doors don't always exist
// 	delete(m.objects.Data, keyPoint)
// 	s := internal.NewStack[string]()
// 	s.Push(k)
// 	doorPoint, ok := doorPoints[strings.ToUpper(k)]
// 	if ok {
// 		delete(m.objects.Data, doorPoint)
// 	}
// 	m.player = keyPoint
// 	p := keyPoint
// 	for p != (internal.Point{-1, -1}) {
// 		p = cameFrom[p]
// 		if w, ok := m.objects.Data[p]; ok {
// 			s.Push(w)
// 			delete(m.objects.Data, p)
// 			doorPoint, ok := doorPoints[strings.ToUpper(w)]
// 			if ok {
// 				delete(m.objects.Data, doorPoint)
// 			}
// 		}
// 	}
// 	for !s.Empty() {
// 		m.keysCollected = append(m.keysCollected, s.Pop())
// 	}
// 	m.totalMovement += costs[keyPoint]
// 	return costs[keyPoint]
// }

// func (m *maze) getGameState() gameState {
// 	return m.searchMaze(m.player, make(map[internal.Point]bool))
// }

func (m *maze) searchMaze(start internal.Point, visited map[internal.Point]bool) gameState {
	out := gameState{
		accessibleKeys:  make([]string, 0),
		accessibleDoors: make([]string, 0),
		gatedKeys:       map[string][]string{},
		gatedDoors:      map[string][]string{},
	}
	q := internal.NewQueue[internal.Point]()
	for p := range m.walls.Data {
		visited[p] = true
	}
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		if visited[cur] {
			continue
		}
		if m.objects.At(cur) >= "a" && m.objects.At(cur) <= "z" {
			out.accessibleKeys = append(out.accessibleKeys, m.objects.At(cur))
		}
		if m.objects.At(cur) >= "A" && m.objects.At(cur) <= "Z" {
			out.accessibleDoors = append(out.accessibleDoors, m.objects.At(cur))
		}
		for _, n := range cur.Neighbors() {
			if visited[n] {
				continue
			}
			if m.objects.At(n) >= "A" && m.objects.At(n) <= "Z" {
				curDoor := m.objects.At(n)
				out.accessibleDoors = append(out.accessibleDoors, m.objects.At(n))
				// find the point beyond the door
				for _, n2 := range n.Neighbors() {
					if n2 == cur {
						continue
					}
					if _, ok := m.walls.Data[n2]; ok {
						continue
					}
					// found it
					visited := map[internal.Point]bool{
						n: true,
					}
					test := m.searchMaze(n2, visited)
					for _, k := range test.accessibleKeys {
						if _, ok := out.gatedKeys[k]; !ok {
							out.gatedKeys[k] = make([]string, 0)
						}
						out.gatedKeys[k] = append(out.gatedKeys[k], curDoor)
					}
					for _, d := range test.accessibleDoors {
						if _, ok := out.gatedDoors[d]; !ok {
							out.gatedDoors[d] = make([]string, 0)
						}
						out.gatedDoors[d] = append(out.gatedDoors[d], curDoor)
					}
					for k, v := range test.gatedKeys {
						out.gatedKeys[k] = append(out.gatedKeys[k], curDoor)
						out.gatedKeys[k] = append(out.gatedKeys[k], v...)
					}
					for d, v := range test.gatedDoors {
						out.gatedDoors[d] = append(out.gatedDoors[d], curDoor)
						out.gatedDoors[d] = append(out.gatedDoors[d], v...)
					}
				}
				continue
			}
			q.Enqueue(n)
		}
		visited[cur] = true
	}
	return out
}

type state3 struct {
	pos    internal.Point
	behind []string
}

func (s state3) copy() state3 {
	b := make([]string, len(s.behind))
	copy(b, s.behind)
	return state3{
		pos:    s.pos,
		behind: b,
	}
}

func (m *maze) searchMaze2(start internal.Point) gameState2 {
	out := gameState2{
		keys: map[string][]string{},
	}
	visited := map[internal.Point]bool{}
	q := internal.NewQueue[state3]()
	for p := range m.walls.Data {
		visited[p] = true
	}
	q.Enqueue(state3{start, make([]string, 0)})
	for !q.Empty() {
		cur := q.Dequeue()
		if visited[cur.pos] {
			continue
		}
		switch v := m.objects.At(cur.pos); {
		case v >= "A" && v <= "Z":
			cur.behind = append(cur.behind, v)
		case v >= "a" && v <= "z":
			out.keys[v] = make([]string, len(cur.behind))
			copy(out.keys[v], cur.behind)
		case v == ".":
		default:
			panic("hello???")
		}
		// if we're on an empty space, check neighbors that we haven't visited an keep going
		// if we're on a key, mark it as a key and keep going
		// if we're on a door, everything past this point is gated, so mark the state as being behind a door
		for _, n := range cur.pos.Neighbors() {
			ns := cur.copy()
			ns.pos = n
			q.Enqueue(ns)
		}
		visited[cur.pos] = true
	}
	return out
}

/*
#######
#a.#Cd#
##@#@##
#######
##@#@##
#cB#Ab#
#######
*/
