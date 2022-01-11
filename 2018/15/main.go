package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type ground string
type characterKind string

const (
	empty      ground        = "."
	wall       ground        = "#"
	elfChar    characterKind = "elf"
	goblinChar characterKind = "goblin"
)

func main() {
	lines := internal.ReadInput()

	for i := 3; ; i++ {
		battle := newGame(i)
		for y, line := range lines {
			for x, c := range line {
				battle.addFloor(x, y, string(c))
				battle.addChar(x, y, string(c))
			}
		}
		for battle.round() {
			if battle.hasElfLosses() {
				fmt.Println("BATTLE IS OVER, AN ELF DIED ON ROUND", battle.roundCount)
				break
			}
			//			fmt.Printf("Round %d\n", battle.roundCount)
			fmt.Println(battle)
		}
		//		fmt.Println(battle.roundCount, battle.livingHitPoints())
		if battle.hasElfLosses() {
			continue
		}
		fmt.Println("no elf loses with power of", i)
		fmt.Println(battle.roundCount * battle.livingHitPoints())
		break
	}
}

// 50700 too low

func newGame(elfPower int) *game {
	return &game{
		landscape:  newGrid(wall),
		characters: newGrid[*character](nil),
		stateLayer: newGrid(" "),
		elfPower:   elfPower,
	}
}

type game struct {
	landscape  *grid[ground]
	characters *grid[*character]
	stateLayer *grid[string]
	roundCount int
	elfPower   int
}

func (g *game) hasElfLosses() bool {
	for _, c := range g.characters.data {
		if c.Kind == elfChar {
			if c.isDead() {
				return true
			}
		}
	}
	return false
}

func (g *game) addFloor(x, y int, d string) {
	if d == "#" {
		return
	}
	g.landscape.add(x, y, empty)
}

func (g *game) addChar(x, y int, d string) {
	switch d {
	case "E":
		g.characters.add(x, y, newElf(x, y, g.elfPower))
	case "G":
		g.characters.add(x, y, newGoblin(x, y))
	}
}

func (g *game) String() string {
	var out strings.Builder
	// landscape will be the bigger of the two
	for j := g.landscape.miny; j <= g.landscape.maxy+1; j++ {
		for i := g.landscape.minx; i <= g.landscape.maxx+1; i++ {
			if s := g.stateLayer.At(i, j); s != " " {
				out.WriteString(s)
				continue
			}
			if c := g.characters.At(i, j); c != nil && !c.isDead() {
				out.WriteString(c.String())
				continue
			}
			out.WriteString(string(g.landscape.At(i, j)))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (g *game) movePhase(c *character) {
	// if g.roundCount == 23 {
	// 	fmt.Println(c.Pos, c.Kind)
	// }
	// if any neighbors are enemies, done with move phase
	for _, n := range c.Pos.neighbors() {
		c2 := g.characters.At(n.x, n.y)
		if c2 == nil || c2.isDead() {
			continue
		}
		if c.isEnemy(c2) {
			return
		}
	}
	// identify all valid targets
	targets := g.targets(c)
	if len(targets) == 0 {
		return
	}
	// if g.roundCount == 23 {
	// 	fmt.Println("targets", targets)
	// }
	inRange := []point{}
	// identify all open squares in range of each target
	for _, t := range targets {
		inRange = append(inRange, g.findInRange(t)...)
	}
	sort.Sort(points(inRange))
	// if g.roundCount == 23 {
	// 	fmt.Println("inRange", inRange)
	// }

	// identify the reachable, in-range targets
	costs := g.djikstra(c.Pos)
	// if g.roundCount == 23 {
	// 	g.printCostsMap(costs)
	// }
	reachable := g.filterUnreachable(costs, inRange)
	if len(reachable) == 0 {
		return //no reachable squares
	}
	// find the closest in-range target
	closest := lowestCostPoint(reachable, costs)
	// find the shortest path to this point
	shortestPathCosts := g.djikstra(closest)
	// if g.roundCount == 23 {
	// 	g.printCostsMap(shortestPathCosts)
	// }
	moveTo := lowestCostPoint(c.Pos.neighbors(), shortestPathCosts)
	g.moveCharacter(c, moveTo)
	//		g.printCostsMap(costs)
}
func (g *game) attackPhase(c *character) {
	lowestHP := math.MaxInt
	var lowest *character
	for _, n := range c.Pos.neighbors() {
		neighbor := g.characters.At(n.x, n.y)
		if neighbor == nil || neighbor.isDead() {
			continue
		}
		if c.isEnemy(neighbor) {
			if neighbor.HitPoints < lowestHP {
				lowest = neighbor
				lowestHP = neighbor.HitPoints
			}
		}
	}
	if lowest != nil {
		c.attack(lowest)
	}
}

func (g *game) round() bool {
	//	for each alive character in reading order
	for _, c := range g.aliveCharactersInReadOrder() {
		if c.isDead() {
			continue
		}
		if len(g.targets(c)) == 0 {
			return false
		}
		g.movePhase(c)
		g.attackPhase(c)
	}
	g.roundCount++
	return true
}

func (g *game) printCostsMap(costs map[point]int) {
	for j := g.landscape.miny; j <= g.landscape.maxy; j++ {
		for i := g.landscape.minx; i <= g.landscape.maxx; i++ {
			if costs[point{i, j}] == math.MaxInt {
				fmt.Printf("## ")
				continue
			}
			fmt.Printf("%02d ", costs[point{i, j}])
		}
		fmt.Println()
	}
	fmt.Println()
}

func lowestCostPoint(pts []point, costs map[point]int) point {
	closestCost := math.MaxInt
	closest := point{}
	for _, p := range pts {
		if v, ok := costs[p]; ok {
			if v < closestCost {
				closestCost = v
				closest = p
			}
		}
	}
	return closest
}

// findInRange returns a list of points that are in range to a target
func (g *game) findInRange(c *character) []point {
	out := []point{}
	for _, p := range c.Pos.neighbors() {
		if g.open(p) {
			out = append(out, p)
		}
	}
	return out
}

func (g *game) livingHitPoints() int {
	sum := 0
	for _, c := range g.characters.data {
		if c.isDead() {
			continue
		}
		fmt.Println(c.Kind, c.Pos, c.HitPoints)
		sum += c.HitPoints
	}
	return sum
}

func (g *game) filterUnreachable(costs map[point]int, inRange []point) []point {
	out := []point{}
	for _, p := range inRange {
		if v, ok := costs[p]; !ok || v == math.MaxInt {
			continue
		}
		out = append(out, p)
	}
	return out
}

func (g *game) open(p point) bool {
	_, ok := g.landscape.data[p]
	if !ok {
		return false
	}
	c := g.characters.At(p.x, p.y)
	return c == nil || c.isDead()
}

func (g *game) targets(c *character) []*character {
	if c.Kind == elfChar {
		return g.findAliveCharacters(goblinChar)
	}
	return g.findAliveCharacters(elfChar)
}

func (g *game) findAliveCharacters(kind characterKind) []*character {
	out := make([]*character, 0)
	// read order
	for j := g.characters.miny; j <= g.characters.maxy; j++ {
		for i := g.characters.minx; i <= g.characters.maxx; i++ {
			c := g.characters.At(i, j)
			if c == nil || c.isDead() {
				continue
			}
			if c.Kind == kind {
				out = append(out, c)
			}
		}
	}
	return out
}

// find the cost of moving from one point to another
func (g *game) djikstra(start point) map[point]int {
	visited := map[point]bool{}
	costs := map[point]int{}
	for j := g.landscape.miny; j <= g.landscape.maxy; j++ {
		for i := g.landscape.minx; i <= g.landscape.maxx; i++ {
			costs[point{i, j}] = math.MaxInt
		}
	}
	costs[start] = 0
	q := internal.NewQueue[point]()
	q.Enqueue(start)
	for !q.Empty() {
		cur := q.Dequeue()
		for _, n := range cur.neighbors() {
			if visited[cur] {
				continue
			}
			if !g.open(n) {
				continue
			}
			if costs[cur]+1 < costs[n] {
				costs[n] = costs[cur] + 1
			}
			q.Enqueue(n)
		}
		visited[cur] = true
	}
	return costs
}

func (g *game) moveCharacter(c *character, p point) {
	delete(g.characters.data, c.Pos)
	c.moveTo(p)
	g.characters.add(p.x, p.y, c)
}

func (g *game) aliveCharactersInReadOrder() []*character {
	out := []*character{}
	for j := g.characters.miny; j <= g.characters.maxy; j++ {
		for i := g.characters.minx; i <= g.characters.maxx; i++ {
			c := g.characters.At(i, j)
			if c == nil || c.isDead() {
				continue
			}
			out = append(out, c)
		}
	}
	return out
}

func newGrid[T any](defaultData T) *grid[T] {
	return &grid[T]{
		data:        make(map[point]T),
		defaultData: defaultData,
	}
}

type grid[T any] struct {
	data                   map[point]T
	minx, maxx, miny, maxy int
	defaultData            T
}

func (g *grid[T]) add(x, y int, t T) {
	if x < g.minx {
		g.minx = x
	}
	if x > g.maxx {
		g.maxx = x
	}
	if y < g.miny {
		g.miny = y
	}
	if y > g.maxy {
		g.maxy = y
	}
	g.data[point{x, y}] = t
}

func (g *grid[T]) At(x, y int) T {
	item, ok := g.data[point{x, y}]
	if !ok {
		return g.defaultData
	}
	return item
}

func (g *grid[T]) String() string {
	var out strings.Builder
	for y := g.miny - 1; y <= g.maxy+1; y++ {
		for x := g.minx - 1; x <= g.maxx+1; x++ {
			out.WriteString(fmt.Sprintf("%v", g.At(x, y)))
		}
		out.WriteString("\n")
	}
	return out.String()
}

type point struct {
	x, y int
}

// very specifically in read order
func (p point) neighbors() []point {
	return []point{
		{p.x, p.y - 1}, {p.x - 1, p.y},
		{p.x + 1, p.y}, {p.x, p.y + 1},
	}
}

type Stringer interface {
	String() string
}

type character struct {
	HitPoints   int
	AttackPower int
	Pos         point
	Kind        characterKind
	HasMoved    bool
}

func (c *character) String() string {
	if c == nil {
		return " "
	}
	if c.Kind == goblinChar {
		return "G"
	}
	return "E"
}

func (c *character) isDead() bool {
	return c.HitPoints <= 0
}

func (c *character) moveTo(p point) {
	c.Pos = p
}

func (c *character) attack(c2 *character) {
	c2.HitPoints -= c.AttackPower
}

func (c *character) isEnemy(c2 *character) bool {
	return c.Kind != c2.Kind
}

func newElf(x, y, attackPower int) *character {
	return &character{
		HitPoints:   200,
		AttackPower: attackPower,
		Pos:         point{x, y},
		Kind:        elfChar,
	}
}

func newGoblin(x, y int) *character {
	return &character{
		HitPoints:   200,
		AttackPower: 3,
		Pos:         point{x, y},
		Kind:        goblinChar,
	}
}

type points []point

func (p points) Len() int      { return len(p) }
func (p points) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p points) Less(i, j int) bool {
	if p[i].y < p[j].y {
		return true
	}
	if p[i].y > p[j].y {
		return false
	}
	return p[i].x < p[j].x
}

// #######
// #G..#E#
// #E#E.E#
// #G.##.#
// #...#E#
// #...E.#
// #######

// #######
// #E..EG#
// #.#G.E#
// #E.##E#
// #G..#.#
// #..E#.#
// #######

// #######
// #E.G#.#
// #.#G..#
// #G.#.G#
// #G..#.#
// #...E.#
// #######

// #######
// #.E...#
// #.#..G#
// #.###.#
// #E#G#G#
// #...#G#
// #######

// #########
// #G......#
// #.E.#...#
// #..##..G#
// #...##..#
// #...#...#
// #.G...G.#
// #.....G.#
// #########
