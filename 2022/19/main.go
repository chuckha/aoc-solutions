package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chuckha/aoc-solutions/internal"
)

/*

guesses

pt1:
1470 too low

pt2:
7448 too low
10880 too low

14 c 3 o
ratio: 4:1
4 clay to 1 ore
*/

/*
Blueprint 1: Each ore robot costs 4 ore.  Each clay robot costs 2 ore.  Each obsidian robot costs 3 ore and 14 clay.  Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore.  Each clay robot costs 3 ore.  Each obsidian robot costs 3 ore and 8 clay.   Each geode robot costs 3 ore and 12 obsidian.
*/

/*



 */

var (
	ore      resource = "ore"
	clay     resource = "clay"
	obsidian resource = "obsidian"
	geode    resource = "geode"
)

const maxTime = 32

func main() {

	input := internal.ReadInput()
	product := 1

	if len(input) > 3 {
		input = input[:3]
	}

	for _, line := range input {
		now := time.Now()
		bp := newBlueprint(line)
		inv := newInventory()

		// sim := newSim(bp, newInventory(), make([]decision, 0))
		// sim = sim.run([]action{
		// 	//			wait, wait, buyClayRobot, wait, buyClayRobot, wait, buyClayRobot, wait, wait, wait, buyObsidianRobot, buyClayRobot, wait, wait, buyObsidianRobot, wait, wait, buyGeodeRobot, wait, wait, buyGeodeRobot, wait, wait, wait,
		// 	wait, wait, buyClayRobot, wait,
		// })
		// fmt.Println(sim)
		// sim = sim.run([]action{buyClayRobot})
		// fmt.Println(sim)
		// sim.undoLastMove()
		// fmt.Println(sim)
		// // fmt.Println(sim.value())
		// os.Exit(0)
		// we have one in our backpack
		sim := newSim(bp, inv, make([]decision, 0))
		max, geodes := sim.minimaxNoCopy()
		fmt.Println(sim.decisionsLog)
		fmt.Println("high score", max)
		fmt.Println("cache hits", sim.cache)
		product *= geodes
		fmt.Printf("Finished blueprint %d in %v with a geode count of %s\n", bp.id, time.Since(now), internal.Red(fmt.Sprintf("%d", geodes)))
	}
	fmt.Println(product)

}

type betterCache map[int]int

func (b betterCache) get(i int) int {
	val, ok := b[i]
	if !ok {
		return math.MinInt
	}
	return val
}

type cell string

func (c cell) String() string {
	return string(c)
}

type robot resource

type resource string

type blueprint struct {
	id                int
	rscCollectorCosts map[resource]cost
	maximums          map[resource]int
	minimums          map[resource]int
	goals             map[resource]int
}

func (bp blueprint) String() string {
	return fmt.Sprintf("Blueprint %d", bp.id)
}

func newBlueprint(line string) blueprint {
	costs := make(map[resource]cost)
	re := regexp.MustCompile(`Blueprint (\d*):`)
	matches := re.FindStringSubmatch(line)
	id, _ := strconv.Atoi(matches[1])

	re = regexp.MustCompile(`Each ore robot costs (\d) ore.`)
	matches = re.FindStringSubmatch(line)
	oreBotCost, _ := strconv.Atoi(matches[1])

	re = regexp.MustCompile(`Each clay robot costs (\d) ore.`)
	matches = re.FindStringSubmatch(line)
	clayBotCost, _ := strconv.Atoi(matches[1])

	re = regexp.MustCompile(`Each obsidian robot costs (\d) ore and (\d+) clay.`)
	matches = re.FindStringSubmatch(line)
	obsOreCost, _ := strconv.Atoi(matches[1])
	obsClayCost, _ := strconv.Atoi(matches[2])

	re = regexp.MustCompile(`Each geode robot costs (\d) ore and (\d+) obsidian.`)
	matches = re.FindStringSubmatch(line)
	geodeOreCost, _ := strconv.Atoi(matches[1])
	geodeObsCost, _ := strconv.Atoi(matches[2])

	costs[ore] = cost{
		ore: oreBotCost,
	}
	costs[clay] = cost{
		ore: clayBotCost,
	}
	costs[obsidian] = cost{
		ore:  obsOreCost,
		clay: obsClayCost,
	}
	costs[geode] = cost{
		ore:      geodeOreCost,
		obsidian: geodeObsCost,
	}

	maximums := make(map[resource]int)
	for rsc, cost := range costs {
		if rsc == ore {
			continue
		}
		if cost.ore > maximums[ore] {
			maximums[ore] = cost.ore
		}
		if cost.clay > maximums[clay] {
			maximums[clay] = cost.clay
		}
		if cost.obsidian > maximums[obsidian] {
			maximums[obsidian] = cost.obsidian
		}
	}
	maximums[geode] = math.MaxInt

	goals := make(map[resource]int)
	goals[obsidian] = costs[geode].obsidian - costs[geode].ore
	goals[clay] = costs[obsidian].clay / 2
	goals[ore] = costs[clay].ore

	minimums := make(map[resource]int)
	minimums[geode] = 1
	minimums[obsidian] = 1
	minimums[clay] = 1
	minimums[ore] = 1

	fmt.Printf("Blueprint %d: Maximums: %v Goals: %v\n", id, maximums, goals)
	return blueprint{
		id:                id,
		rscCollectorCosts: costs,
		maximums:          maximums,
		goals:             goals,
		minimums:          minimums,
	}
}

type cost struct {
	ore      int
	clay     int
	obsidian int
}

type inventory struct {
	collectors map[resource]int
	resources  map[resource]int
	building   map[resource]int
}

func newInventory() *inventory {
	return &inventory{
		collectors: map[resource]int{
			ore: 1,
		},
		resources: map[resource]int{},
		building:  make(map[resource]int),
	}
}

func (i *inventory) String() string {
	out := []string{"=== collectors ==="}
	for c, cnt := range i.collectors {
		out = append(out, fmt.Sprintf("%v (%d)", c, cnt))
	}
	out = append(out, "--- resources ---")
	for c, cnt := range i.resources {
		out = append(out, fmt.Sprintf("%v (%d)", c, cnt))
	}
	return strings.Join(out, "\n")
}

func (i *inventory) addCollector(r resource, n int) {
	i.collectors[r] += n
}
func (i *inventory) addResource(r resource, n int) {
	i.resources[r] += n
}
func (i *inventory) canAfford(c cost) bool {
	return i.resources[ore] >= c.ore && i.resources[obsidian] >= c.obsidian && i.resources[clay] >= c.clay
}

func (i *inventory) spend(rsc resource, c cost) {
	i.resources[ore] -= c.ore
	i.resources[obsidian] -= c.obsidian
	i.resources[clay] -= c.clay
	i.building[rsc]++
}

func (i *inventory) resolveBuilding() {
	for b, cnt := range i.building {
		i.addCollector(b, cnt)
	}
	i.building = make(map[resource]int)
}

func (i *inventory) copy() *inventory {
	collectors := make(map[resource]int)
	resources := make(map[resource]int)
	building := make(map[resource]int)
	for k, v := range i.collectors {
		collectors[k] = v
	}
	for k, v := range i.resources {
		resources[k] = v
	}
	for k, v := range i.building {
		building[k] = v
	}

	return &inventory{
		collectors: collectors,
		resources:  resources,
		building:   building,
	}
}

type decision struct {
	action action
	time   int
}

type simulation struct {
	bp           blueprint
	inv          *inventory
	time         int
	decisionsLog []decision
	alpha        int
	alphaRate    int
	betaRate     int
	gammaRate    int
	deltaRate    int
	cache        map[int]int
}

func newSim(bp blueprint, inv *inventory, dl []decision) *simulation {
	return &simulation{
		bp:           bp,
		inv:          inv,
		decisionsLog: dl,
		cache:        make(betterCache),
	}
}

func (s *simulation) resolveTick() {
	for rsc, num := range s.inv.collectors {
		s.inv.addResource(rsc, num)
	}
	s.inv.resolveBuilding()
	s.time++
}

func (s *simulation) String() string {
	out := []string{fmt.Sprintf("%v (%d)", s.bp, s.time)}
	for _, decision := range s.decisionsLog {
		out = append(out, fmt.Sprintf("\t%v @ %v", decision.action, decision.time))
	}
	out = append(out, fmt.Sprintf("%v", s.inv))
	return strings.Join(out, "\n")
}

type action string

const (
	buyOreRobot      action = "buy ore robot"
	wait             action = "wait"
	buyClayRobot     action = "buy clay robot"
	buyObsidianRobot action = "buy obsidian robot"
	buyGeodeRobot    action = "buy geode robot"
)

func (s *simulation) undoLastMove() {
	s.undo(s.decisionsLog[len(s.decisionsLog)-1].action)
	s.decisionsLog = s.decisionsLog[:len(s.decisionsLog)-1]
}

func (s *simulation) undo(act action) {
	switch act {
	case buyOreRobot:
		// unbuy the robot/ remove from inventory
		s.inv.collectors[ore]--
		cost := s.bp.rscCollectorCosts[ore]
		s.inv.resources[clay] += cost.clay
		s.inv.resources[obsidian] += cost.obsidian
		s.inv.resources[ore] += cost.ore
	case buyClayRobot:
		s.inv.collectors[clay]--
		cost := s.bp.rscCollectorCosts[clay]
		s.inv.resources[clay] += cost.clay
		s.inv.resources[obsidian] += cost.obsidian
		s.inv.resources[ore] += cost.ore
	case buyObsidianRobot:
		s.inv.collectors[obsidian]--
		cost := s.bp.rscCollectorCosts[obsidian]
		s.inv.resources[clay] += cost.clay
		s.inv.resources[obsidian] += cost.obsidian
		s.inv.resources[ore] += cost.ore
	case buyGeodeRobot:
		s.inv.collectors[geode]--
		cost := s.bp.rscCollectorCosts[geode]
		s.inv.resources[clay] += cost.clay
		s.inv.resources[obsidian] += cost.obsidian
		s.inv.resources[ore] += cost.ore
	case wait:
	}
	// remove a resource for every collector
	for rsc, count := range s.inv.collectors {
		s.inv.resources[rsc] -= count
	}
	// rewind time
	s.time--
}

func (s *simulation) makeNoCopy(act action) {
	switch act {
	case buyOreRobot:
		s.inv.spend(ore, s.bp.rscCollectorCosts[ore])
	case buyClayRobot:
		s.inv.spend(clay, s.bp.rscCollectorCosts[clay])
	case buyObsidianRobot:
		s.inv.spend(obsidian, s.bp.rscCollectorCosts[obsidian])
	case buyGeodeRobot:
		s.inv.spend(geode, s.bp.rscCollectorCosts[geode])
	case wait:
	}
	s.resolveTick()
	s.decisionsLog = append(s.decisionsLog, decision{action: act, time: s.time})
}

// decisions return a list of possible decisions a simulation can make
// it should never be in a position to make more than one action.
func (s *simulation) decisions() map[action]struct{} {
	decisions := map[action]struct{}{}
	// always wait at the last minute
	if s.time == maxTime-1 {
		return map[action]struct{}{wait: {}}
	}

	// for _, rsc := range []resource{geode, ore, obsidian, clay} {
	// 	if s.inv.collectors[rsc] < s.bp.goals[rsc] {
	// 		if s.inv.canAfford(s.bp.rscCollectorCosts[rsc]) {
	// 			decisions[actionFactory(rsc)] = struct{}{}
	// 		}
	// 	}
	// }

	for _, rsc := range []resource{geode, obsidian, ore, clay} {
		if _, ok := decisions[actionFactory(rsc)]; ok {
			continue
		}
		if s.time >= 20 && (rsc == clay || rsc == ore) {
			continue
		}
		if s.time >= 28 && rsc == obsidian {
			continue
		}
		if s.inv.collectors[rsc] >= s.bp.maximums[rsc] {
			continue
		}
		if s.inv.canAfford(s.bp.rscCollectorCosts[rsc]) {
			decisions[actionFactory(rsc)] = struct{}{}
		}
	}
	decisions[wait] = struct{}{}
	return decisions
}

func actionFactory(rsc resource) action {
	switch rsc {
	case ore:
		return buyOreRobot
	case clay:
		return buyClayRobot
	case obsidian:
		return buyObsidianRobot
	case geode:
		return buyGeodeRobot
	default:
		panic("bad resource")
	}
}

func (s *simulation) value() int {
	//	Blueprint 1: Each ore robot costs 4 ore.  Each clay robot costs 2 ore.  Each obsidian robot costs 3 ore and 14 clay.  Each geode robot costs 2 ore and 7 obsidian.
	points := 10
	desiredClay := s.bp.rscCollectorCosts[obsidian].clay / s.bp.rscCollectorCosts[obsidian].ore
	actualClay := s.inv.collectors[clay] / s.inv.collectors[ore]
	points -= internal.Abs(desiredClay - actualClay)
	desiredObsidian := s.bp.rscCollectorCosts[geode].obsidian / s.bp.rscCollectorCosts[geode].ore
	actualObsidian := s.inv.collectors[obsidian] / s.inv.collectors[ore]
	points -= internal.Abs(desiredObsidian - actualObsidian)
	return points + s.inv.resources[geode]*10
}

func (s *simulation) valuev2() int {
	//	Blueprint 1: Each ore robot costs 4 ore.  Each clay robot costs 2 ore.  Each obsidian robot costs 3 ore and 14 clay.  Each geode robot costs 2 ore and 7 obsidian.
	// points := 0
	// desiredOre := s.bp.goals[ore]
	// actualOre := s.inv.collectors[ore]
	// // fmt.Println("desired ore", desiredOre, "actual", actualOre)
	// points += 10 - internal.Abs(desiredOre-actualOre)
	// // fmt.Println("points after ore", points, desiredOre, actualOre)
	// desiredClay := s.bp.goals[clay] / desiredOre
	// actualClay := s.inv.collectors[clay] / desiredOre
	// points += 10 - internal.Abs(desiredClay-actualClay)
	// // fmt.Println("points after clay", points, desiredClay, actualClay)
	// desiredObsidian := s.bp.goals[obsidian]
	// actualObsidian := s.inv.collectors[obsidian]
	// points += 10 - internal.Abs(desiredObsidian-actualObsidian)
	// fmt.Println("points after obsidian", points, desiredObsidian, actualObsidian)
	//	return points + s.inv.resources[geode]
	return s.inv.resources[geode]
}

func (s *simulation) key() string {
	return fmt.Sprintf("%v %v", s.inv.collectors, s.inv.resources)
}

// func minimax(sim *simulation, cache map[string]int) (float64, *simulation) {
// 	if sim.time == 24 {
// 		return sim.value(), sim
// 	}
// 	//	fmt.Println(sim.key())
// 	// if we've seen this set up before but with a lower turn number, return a negative value
// 	if t, ok := cache[sim.key()]; ok && t < sim.time {
// 		//		fmt.Println("cache hit", sim.key())
// 		return -100, nil
// 	}
// 	//	fmt.Println(sim)
// 	cache[sim.key()] = sim.time
// 	max := float64(0)
// 	var s4 *simulation
// 	decisions := sim.decisions()
// 	fmt.Printf("(%d) making %d decisions\n", sim.time, len(decisions))
// 	for _, decision := range sim.decisions() {
// 		s2 := sim.make(decision)
// 		val, s3 := minimax(s2, cache)
// 		//		fmt.Println(val, sim.inv.collectors, sim.inv.resources)
// 		if val > max {
// 			//			fmt.Println("best", val, s3)
// 			max = val
// 			s4 = s3
// 		}
// 	}
// 	fmt.Println("max", max)
// 	return max, s4
// }

func (s *simulation) geodeUpperBound() int {
	// if we bought a geode robot every turn from now to the end of the game
	// if there are 10 minutes left, we buy 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10 n*(n+1)/2
	n := maxTime - s.time
	return s.inv.resources[geode] + (s.inv.collectors[geode] * n) + (n*(n+1))/2
}

func (sim *simulation) minimaxNoCopy() (int, int) {
	if sim.time == maxTime {
		return sim.valuev2(), sim.inv.resources[geode]
	}

	if sim.geodeUpperBound() <= sim.alpha {
		return 0, sim.inv.resources[geode]
	}

	max := 0
	maxGeode := 0
	for decision := range sim.decisions() {
		sim.makeNoCopy(decision)
		val, geodes := sim.minimaxNoCopy()
		//	fmt.Println(val, sim.inv.collectors, sim.inv.resources)
		if val > max {
			max = val
			maxGeode = geodes
		}
		if val > sim.alpha {
			sim.alpha = val
			sim.alphaRate = sim.inv.collectors[geode]
			sim.betaRate = sim.inv.collectors[obsidian]
			sim.gammaRate = sim.inv.collectors[clay]
			sim.deltaRate = sim.inv.collectors[ore]
			fmt.Println(sim.alpha, sim.alphaRate, sim.betaRate, sim.gammaRate, sim.deltaRate, maxGeode)
			fmt.Println(sim.decisionsLog)
		}
		sim.undoLastMove()
	}
	return max, maxGeode
}

// our value at 24 = geode count
// our value at 23 is the maximum of (wait, buy something)
// our value at 22 is the maximum of (wait or buy soemthing)

// goal 1 is to achieve the ore/clay ratio of robots
// after which we switch to goal 2
// goal 2 is to achieve the geode robot ratio
// after which we switch to goal 3
// goal 3 -> buy geode robots

// Goal 1: == obsidian robot cost (x ore, y clay)
// 		1 ore robot & y/x clay robots
// goal 2: == geode robot cost (x ore, z obsidian)
// 		1 ore robot & z/x obsidian robots

// Blueprint 1:
// t | o | or | c | cr | ob | obr | geo | geor
// 0   0 | 1
// 1   1 | 1
// 2   0 | 1  | 0 | 1
// 3   1 | 1  |
// 	f(t) = ore_robots + clay_robots + obs_robots + geode_robots

//   Each ore robot costs 4 ore.    => 1 ore robot == 4 ore + 1 time
//   Each clay robot costs 2 ore.   => 1 clay robot == 2 ore + 1 time
//   Each obsidian robot costs 3 ore and 14 clay. => 1 obsidian robot ==
//   Each geode robot costs 2 ore and 7 obsidian.

//   ratio: 7 obsidian robots : 2 ore robots
//   3 ores robots : 14 clay robots

//   28 ore + 14 time == 14 clay robots
//   28 time (with 1 ore robot)
//   15 time (with 2 ore robot (+1 time))

/*

given a time we can produce the number of resources it will have at any time with no further purchases
resource(robots, t) = (maxTime - t) * robots[rsc]
- what decisions can we make?
if amount of resources are able to buy something, branch into buy it and not
what is the maximum ore we need to optimize the amount of geodes? don't explore any space with more ore than this
what is the m aximum amount of obsidian robots we want? don't explore any space with morethan that
what is the maximum amount of clay we want? don't explore any space with more than that






*/
