package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const duration = 30

func main() {
	input := internal.ReadInput()
	valveData := make(map[string]int)
	connections := make(map[string][]string)
	for _, line := range input {
		re := regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=(\d*);`)
		matches := re.FindStringSubmatch(line)
		fr, _ := strconv.Atoi(matches[2])
		valveData[matches[1]] = fr
		parts := strings.Split(line, " valves ")
		if len(parts) == 1 {
			parts = strings.Split(line, " valve ")
		}
		valves := strings.Split(parts[1], ", ")
		connections[matches[1]] = valves
	}
	g := newGame(26, valveData, connections)
	gates := []string{}
	remaining := make(internal.SetV2)
	for gate, v := range g.flowRates {
		if v <= 0 {
			continue
		}
		gates = append(gates, gate)
		remaining[gate] = struct{}{}
	}
	optional := 0
	if len(gates)%2 == 1 {
		optional = 1
	}
	out := internal.AllCombinations([]string{}, gates, len(gates)/2)
	out2 := internal.AllCombinations([]string{}, gates, len(gates)/2+optional)

	biggest := 0
	rounds := [2]internal.SetV3{}
	for _, o := range out {
		for _, k := range out2 {
			first := internal.SetV2{}
			second := internal.SetV2{}
			for _, k := range o {
				first[k] = struct{}{}
				second[k] = struct{}{}
			}
			third := internal.SetV2{}
			for _, j := range k {
				first[j] = struct{}{}
				third[j] = struct{}{}
			}
			if len(first) == len(gates) {
				max1, v1 := g.cost("AA", second, make(internal.SetV3), 0)
				max2, v2 := g.cost("AA", third, make(internal.SetV3), 0)
				if max1+max2 > biggest {
					biggest = max1 + max2
					rounds[0] = v1
					rounds[1] = v2
					fmt.Println(biggest, rounds)
				}
			}
		}
	}

	//	max1, v1 := g.cost("AA", remaining1, make(internal.SetV3), 0)
	//	max2, v2 := g.cost("AA", remaining2, make(internal.SetV3), 0)

	//	fmt.Println(max1+max2, v1, v2)
	//g.combosv2()

	// for _, combo := range g.combos() {
	// 	fr := g.totalFlowRateTwo(2, combo)
	// 	if fr > max {
	// 		fmt.Println(combo, fr)
	// 		max = fr
	// 	}
	// }
	//	fmt.Println(g.totalFlowRateTwo(2, []string{"JJ", "DD", "HH", "BB", "CC", "EE"}))
	//	fmt.Println(g.totalFlowRateTwo(2, []string{"DD", "JJ", "HH", "BB", "CC", "EE"}))
	//	fmt.Println(g.totalFlowRateTwo(2, []string{"DD", "BB", "HH", "JJ", "CC", "EE"}))

	// f := g.totalFlowRateTwo(1, []string{"DD", "JJ", "BB", "HH", "CC", "EE"})
	//	f := g.totalFlowRateTwo(1, []string{"DD", "BB", "JJ", "HH", "EE", "CC"})

	//	fmt.Println(f)
	//
	// max := math.MinInt
	//
	//	for _, c := range g.combos() {
	//		v := g.totalFlowRate(c)
	//		if v > max {
	//			fmt.Println(c, v)
	//			max = v
	//		}
	//	}
	//
	// fmt.Println(max)
}

func (g *game) cost(pos string, remaining internal.SetV2, visited internal.SetV3, time int) (int, internal.SetV3) {
	if len(remaining) == 0 {
		return 0, visited
	}

	max := 0
	out := visited
	for k := range remaining {
		timeLapse := time + g.shortPaths[pos][k].cost + 1
		if timeLapse > g.duration {
			continue
		}
		subCost, newV := g.cost(k, remaining.Remove(k), visited.Add(k, timeLapse), timeLapse)
		cost := g.flowRates[k]*(g.duration-(timeLapse)) + subCost
		if cost > max {
			out = newV
			max = cost
		}
	}
	return max, out
}

type game struct {
	data       map[string][]string
	flowRates  map[string]int
	shortPaths map[string]map[string]path
	duration   int
}

func newGame(duration int, flowRates map[string]int, conn map[string][]string) *game {
	shortestPaths := make(map[string]map[string]path)
	game := &game{
		data:      conn,
		flowRates: flowRates,
		duration:  duration,
	}
	for k := range conn {
		shortestPaths[k] = game.paths(k)
	}
	game.shortPaths = shortestPaths
	return game
}

type path struct {
	path []string
	cost int
}

func (g *game) paths(start string) map[string]path {
	visited := map[string]struct{}{}
	costs := map[string]path{}
	for p := range g.data {
		costs[p] = path{path: make([]string, 0), cost: math.MaxInt}
	}
	costs[start] = path{cost: 0, path: make([]string, 0)}
	q := internal.NewQueue[string]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		if _, ok := visited[cur]; ok {
			continue
		}
		for _, n := range g.data[cur] {
			cost := 1
			if costs[cur].cost+cost < costs[n].cost {
				np := make([]string, len(costs[cur].path))
				copy(np, costs[cur].path)
				np = append(np, n)
				costs[n] = path{
					cost: costs[cur].cost + cost,
					path: np,
				}
				q.Enqueue(n)
			}
		}
		visited[cur] = struct{}{}
	}
	return costs
}

type debug bool

func (d debug) Println(args ...any) {
	if d {
		fmt.Println(args...)
	}
}

func (g *game) totalFlowRateTwo(pyrs int, ordering []string) int {
	d := debug(false)
	durations := map[int]int{
		1: 30,
		2: 26,
	}
	g.duration = durations[pyrs]
	players := []*player{}
	for i := 0; i < pyrs; i++ {
		players = append(players, newPlayer())
	}
	openedGates := []string{}
	flow := 0
	time := 0
	placedidx := 0
	for {
		playersNeedingDest := 0
		for _, p := range players {
			if p.needsDestination() {
				playersNeedingDest++
			}
		}
		if playersNeedingDest == len(players) && placedidx == len(ordering) {
			finish := g.duration - time
			d.Println("ending!", time, g.duration, finish, openedGates)
			for _, o := range openedGates {
				flow += g.flowRates[o] * finish
			}
			return flow
		}

		d.Println("----- time", time)
		timeToTick := 0
		gatesToOpen := []string{}
		for _, player := range players {
			if placedidx == len(ordering) {
				continue
			}
			if player.needsDestination() {
				player.dest = ordering[placedidx]
				player.dist = g.shortPaths[player.pos][player.dest].cost
				placedidx++
			}
		}
		for _, player := range players {
			d.Println(player)
		}

		closestToDestIdx := -2
		for i, player := range players {
			if player.needsDestination() {
				continue
			}
			if closestToDestIdx == -2 {
				closestToDestIdx = i
				continue
			}
			if player.dist < players[closestToDestIdx].dist {
				closestToDestIdx = i
			}
		}
		distance := players[closestToDestIdx].dist
		timeToTick = distance
		for _, player := range players {
			// there is no movement
			if distance == 0 {
				continue
			}
			if player.needsDestination() {
				continue
			}
			d.Println(g.shortPaths[player.pos][player.dest].path)
			player.pos = g.shortPaths[player.pos][player.dest].path[distance-1]
			player.dist -= distance
		}

		for _, player := range players {
			if player.needsDestination() {
				continue
			}
			if player.dist == 0 {
				gatesToOpen = append(gatesToOpen, player.dest)
				player.reset()
				continue
			}
		}
		optionalMinute := 0
		if len(gatesToOpen) > 0 {
			optionalMinute = 1
		}
		for _, player := range players {
			if optionalMinute == 0 {
				continue
			}
			if player.needsDestination() {
				continue
			}
			if player.dist == 0 {
				continue
			}
			player.pos = g.shortPaths[player.pos][player.dest].path[0]
			player.dist -= 1
		}

		for _, o := range openedGates {
			flow += g.flowRates[o] * (timeToTick + optionalMinute)
		}
		time += timeToTick + optionalMinute
		d.Println("opened gates at time", time, gatesToOpen)
		openedGates = append(openedGates, gatesToOpen...)
		if time > g.duration {
			fmt.Println("too long?")
			return -1
		}
	}
	/*
				for each player
			if they don't have any desination
				assign them one of the valid destinations

		find the player closest to their destination
			move both playes that far

		for each player
			if they are at their gate
				add one to the time elapsed (closest distance)
				add gate to list of gates to open
		calculate flow released
		actually open gate(s)

	*/
}

func (g *game) combos() [][]string {
	gatesToOpen := []string{}
	for gate, rate := range g.flowRates {
		if rate <= 0 {
			continue
		}
		gatesToOpen = append(gatesToOpen, gate)
	}
	return internal.AllPermutations([]string{}, gatesToOpen)
}

func (g *game) openableGates() []string {
	gatesToOpen := []string{}
	for gate, rate := range g.flowRates {
		if rate <= 0 {
			continue
		}
		gatesToOpen = append(gatesToOpen, gate)
	}
	return gatesToOpen
}

func (g *game) values(remainingTime int, start string, openableGates []string) map[string]int {
	out := map[string]int{}
	for _, gate := range openableGates {
		dist := g.shortPaths[start][gate].cost + 1
		out[gate] = g.flowRates[gate] * (remainingTime - dist)
	}
	return out
}

func max(totalRates map[string]int) string {
	max := math.MinInt
	key := ""
	for k, v := range totalRates {
		if v > max {
			max = v
			key = k
		}
	}
	return key
}

func (g *game) combosv2() [][]string {
	d := debug(true)
	player := newPlayer()
	time := 30

	opened := map[string]struct{}{}
	openableGates := g.openableGates()
	for k, v := range g.flowRates {
		if v <= 0 {
			continue
		}
		d.Println(k, ":", v)
	}
	// sort gates by value

	// open it if it's on the way (or near by) IF
	// 		getting to it and opening it (dist+1) is less than < ((dist * 2) + 1) * target flow
	//      sideTarget * (dist off path + 1) < (distance off current path * 2 + 1) * target flow
	for {
		if time == 0 {
			return nil
		}
		costs := g.values(time, player.pos, openableGates)
		d.Println(costs)
		best := max(costs)
		d.Println("best is", best)
		player.dest = best
		player.dist = g.shortPaths[player.pos][player.dest].cost

		// walk along a path toward the highest value
		for _, p := range g.shortPaths[player.pos][player.dest].path {
			time -= 1
			player.pos = p

			// on the path
			if internal.Search(p, openableGates) != -1 && p != player.dest {
				d.Println("found", p, "on the way to ", player.dest)
				d.Println("should I open it?")
				if g.flowRates[p] < g.flowRates[player.dest] {
					d.Println("yes, yes we should")
				}
			}
			nc := g.values(time, player.pos, openableGates)
			// check the distance to every other openable gate
			for _, gate := range openableGates {
				if gate == p {
					continue
				}
				distOffPath := g.shortPaths[p][gate].cost
				if (distOffPath+1)*g.flowRates[gate] > (distOffPath*2+1)*g.flowRates[player.dest] {
					fmt.Println("this is a good deviation", gate)
				}
			}
			d.Println(nc)
		}
		opened[player.pos] = struct{}{}
		idx := internal.Search(player.pos, openableGates)
		openableGates = append(openableGates[:idx], openableGates[idx+1:]...)
	}
	return nil
	// at every step, check every path's value, potentially update
}

// func (g *game) combos() [][]string {
// 	out := make([][]string, 0)

// 	q := internal.NewQueue[state]()
// 	start := state{actor1: "AA", actor2: "AA", opened: make([]string, 0), time: 0, fr: 0}

// 	q.Enqueue(start)
// 	for !q.Empty() {
// 		cur := q.Dequeue()
// 		if cur.time == g.duration {
// 			out = append(out, cur.opened)
// 			continue
// 		}

// 		vs := g.remainingSortedValves(cur)
// 		if len(vs) == 0 {
// 			remaining := g.duration - cur.time
// 			for _, o := range cur.opened {
// 				cur.fr += g.flowRates[o] * remaining
// 			}
// 			out = append(out, cur.opened)
// 			continue
// 		}
// 		for i := 0; i < len(vs)-1; i++ {
// 			//		for _, v := range vs {
// 			s := cur.copy()
// 			// figure out which is further

// 			dist := g.shortPaths[s.actor1][v.name].cost + 1
// 			for _, o := range s.opened {
// 				s.fr += g.flowRates[o] * dist
// 			}
// 			s.opened = append(s.opened, v.name)
// 			s.time += dist
// 			s.actor1 = v.name
// 			if s.time > g.duration {
// 				continue
// 			}
// 			q.Enqueue(s)
// 		}
// 	}
// 	return out
// }

type state struct {
	actor1 string
	actor2 string
	opened []string
	time   int
	fr     int
}

func (s state) copy() state {
	out := make([]string, len(s.opened))
	copy(out, s.opened)
	return state{
		actor1: s.actor1,
		actor2: s.actor2,
		opened: out,
		time:   s.time,
		fr:     s.fr,
	}
}

// remaining sorted tells you the next one to go to

func (g *game) remainingSortedValves(s state) valves {
	out := make(valves, 0)
	for n, f := range g.flowRates {
		if f == 0 {
			continue
		}
		if internal.Search(n, s.opened) >= 0 {
			continue
		}
		dist := g.shortPaths[s.actor1][n].cost
		value := f * (g.duration - s.time - (dist + 1))
		if value <= 0 {
			continue
		}
		out = append(out, valve{name: n, total: value})
	}
	sort.Sort(sort.Reverse(out))
	return out
}

type valve struct {
	name  string
	total int
}
type valves []valve

func (a valves) Len() int           { return len(a) }
func (a valves) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a valves) Less(i, j int) bool { return a[i].total < a[j].total }

type player struct {
	pos  string
	dest string
	dist int
}

func (p *player) reset() {
	p.pos = p.dest
	p.dist = -1
	p.dest = ""
}

func newPlayer() *player {
	return &player{
		pos:  "AA",
		dest: "",
		dist: -1,
	}
}
func (p *player) String() string {
	return fmt.Sprintf("POS: %q, DEST: %q, DIST: %d", p.pos, p.dest, p.dist)
}

func (p *player) needsDestination() bool {
	return p.dist == -1
}
