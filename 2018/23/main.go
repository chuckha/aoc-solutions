package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	nbs := make([]nanobot, 0)
	for _, line := range lines {
		nbs = append(nbs, parseLine(line))
	}
	// find the largest bounding cube
	var min, max point
	for _, nb := range nbs {
		// if nb.pos.x-nb.scanRadius < minmax[0] {
		// 	minmax[0] = nb.pos.x - nb.scanRadius
		// }
		// if nb.pos.x+nb.scanRadius > minmax[1] {
		// 	minmax[1] = nb.pos.x + nb.scanRadius
		// }

		// if nb.pos.y-nb.scanRadius < minmax[2] {
		// 	minmax[2] = nb.pos.y - nb.scanRadius
		// }
		// if nb.pos.y+nb.scanRadius > minmax[3] {
		// 	minmax[3] = nb.pos.y + nb.scanRadius
		// }

		// if nb.pos.z-nb.scanRadius < minmax[4] {
		// 	minmax[4] = nb.pos.z - nb.scanRadius
		// }
		// if nb.pos.z+nb.scanRadius > minmax[5] {
		// 	minmax[5] = nb.pos.z + nb.scanRadius
		// }
		if nb.pos.x < min.x {
			min.x = nb.pos.x
		}
		if nb.pos.x > max.x {
			max.x = nb.pos.x
		}

		if nb.pos.y < min.y {
			min.y = nb.pos.y
		}
		if nb.pos.y > max.y {
			max.y = nb.pos.y
		}

		if nb.pos.z < min.z {
			min.z = nb.pos.z
		}
		if nb.pos.z > max.z {
			max.z = nb.pos.z
		}
	}
	// for i := 0; i < 30; i++ {

	// 	// fmt.Println("--", c.corners())
	// 	// fmt.Println("--", c.size())
	// 	c.split()
	// 	largest := 0
	// 	li := -1
	// 	for i, sc := range c.subCubes {
	// 		// fmt.Println(sc.center())
	// 		touching := sc.touchingNanobots(nbs)
	// 		if touching > largest {
	// 			li = i
	// 			largest = touching
	// 		}
	// 		if touching == largest {
	// 			if sc.center().dist(point{0, 0, 0}) < c.subCubes[li].center().dist(point{0, 0, 0}) {
	// 				li = i
	// 				largest = touching
	// 			}
	// 		}
	// 		// fmt.Println(i, touching)
	// 	}
	// 	// fmt.Println("picking", li)
	// 	c = c.subCubes[li]
	// 	if c.size() < 8 {
	// 		break
	// 	}
	// }
	// fmt.Println(c, c.size())
	// fmt.Printf("x: %d -> %d\n", c.minx, c.maxx)
	// fmt.Printf("y: %d -> %d\n", c.miny, c.maxy)
	// fmt.Printf("z: %d -> %d\n", c.minz, c.maxz)
	// for _, p := range c.allPoints() {
	// 	fmt.Println(p, c.touchingNanobots(nbs), p.dist(point{0, 0, 0}))
	// }
	// os.Exit(0)
	priority := make(pq2, 0)
	heap.Init(&priority)
	//	pq := NewPriorityQueue[*cube]()
	bigCube := newCube(min, max)
	fmt.Println("min/max corners", min, max)
	fmt.Println("center", bigCube.center)
	count := bigCube.touchingNanobots(nbs)
	bigCube.count = count

	heap.Push(&priority, &d3{
		//		size:  bigCube.volume(),
		size:  bigCube.volume2(),
		count: bigCube.count,
		dist:  bigCube.center.dist(point{0, 0, 0}),
		cube:  bigCube,
	})
	//	pq.Add(c, data{size: c.size(), count: 0, dist: 0})
	smallCubes := []*cube{}
	bestSmallScore := 0
	//	for !pq.Empty() {
	i := 0
	for priority.Len() != 0 {
		//		fmt.Println("remaining items", priority.Len(), len(smallCubes))
		//		fmt.Println("best small score", bestSmallScore)
		cur := heap.Pop(&priority).(*d3).cube
		fmt.Printf("Next Cube: sz=%d, count=%d, %v,%v, queue=%d\n", cur.volume2(), cur.count, cur.min, cur.max, len(priority))
		i++

		for _, c := range cur.split() {
			count := c.touchingNanobots(nbs)
			c.count = count
			vol := c.volume2()
			// protec from integer overflow
			//if vol > 0 && vol <= 8 && c.count >= bestSmallScore {
			if vol.Cmp(big.NewInt(8)) <= 0 {
				if c.count > bestSmallScore {
					fmt.Println("new best", c.count)
					bestSmallScore = c.count
					smallCubes = []*cube{c}
					continue
				}
				if c.count == bestSmallScore {
					smallCubes = append(smallCubes, c)
					continue
				}
			}
			if c.count >= bestSmallScore && c.count > 0 {
				//				if vol > 8 || vol < 0 {
				if vol.Cmp(big.NewInt(8)) == 1 {
					//				fmt.Println(c, vol)
					heap.Push(&priority, &d3{
						//				pq.Add(c, data{
						cube:  c,
						size:  vol,
						count: c.count,
						dist:  c.center.dist(point{0, 0, 0}),
					})
					continue
				}
			}
		}
	}
	//	fmt.Println("done?", len(smallCubes), smallCubes)
	biggestCount := math.MinInt
	bestPoint := point{0, 0, 0}
	for _, cube := range smallCubes {
		for _, p := range cube.allPoints() {
			nbc := p.intersecting(nbs)
			if nbc > biggestCount {
				biggestCount = nbc
				bestPoint = p
			}
			if nbc == biggestCount && p.dist(point{0, 0, 0}) < bestPoint.dist(point{0, 0, 0}) {
				bestPoint = p
			}
		}
	}
	fmt.Println(bestPoint)
	fmt.Println(bestPoint.intersecting(nbs))
	fmt.Println(bestPoint.dist(point{0, 0, 0}))

	nb := biggestCount
	np := bestPoint
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			for k := 0; k < 100; k++ {
				p := point{bestPoint.x + i, bestPoint.y + j, bestPoint.z + k}
				count := p.intersecting(nbs)
				if count > nb {
					nb = count
					np = p
				}
				if count == nb && p.dist(point{0, 0, 0}) < np.dist(point{0, 0, 0}) {
					nb = count
					np = p
				}
			}
		}
	}
	fmt.Println(nb, np)
	fmt.Println(np.intersecting(nbs))
}

// too low  101599536
// too low  101599537
//          101599540
// ntra     106678601
// too high 106678603

type cube struct {
	min, center, max point
	count            int
}

func (c *cube) isPoint() bool {
	return c.min == c.max
}

func (c *cube) volume() int {
	return abs(c.max.x-c.min.x) * abs(c.max.y-c.min.y) * abs(c.max.z-c.min.z)
}
func (c *cube) volume2() *big.Int {
	x := big.NewInt(int64(c.max.x - c.min.x))
	x.Mul(x, big.NewInt(int64(c.max.y-c.min.y)))
	return x.Mul(x, big.NewInt(int64(c.max.z-c.min.z)))
}
func newCube(min, max point) *cube {
	return &cube{
		min:    min,
		center: point{(max.x + min.x) / 2, (max.y + min.y) / 2, (max.z + min.z) / 2},
		max:    max,
	}
}

func (c *cube) corners() []point {
	return []point{
		c.min,
		{c.min.x, c.min.y, c.max.z},
		{c.min.x, c.max.y, c.min.z},
		{c.min.x, c.max.y, c.max.z},
		{c.max.x, c.min.y, c.min.z},
		{c.max.x, c.min.y, c.max.z},
		{c.max.x, c.max.y, c.min.z},
		c.max,
	}
}

func (c *cube) split() []*cube {
	if c.isPoint() {
		panic("don't split a point")
	}
	out := []*cube{}
	for _, p := range c.corners() {
		minP := point{min(c.center.x, p.x), min(c.center.y, p.y), min(c.center.z, p.z)}
		maxP := point{max(c.center.x, p.x), max(c.center.y, p.y), max(c.center.z, p.z)}
		out = append(out, newCube(minP, maxP))
	}
	return out
}

func (c *cube) overlap(n nanobot) bool {
	if c.min.x <= n.pos.x && c.max.x >= n.pos.x &&
		c.min.y <= n.pos.y && c.max.y >= n.pos.y &&
		c.min.z <= n.pos.z && c.max.z >= n.pos.z {
		return true
	}
	d := 0
	if n.pos.x < c.min.x {
		d += c.min.x - n.pos.x
	}
	if n.pos.x > c.max.x {
		d += n.pos.x - c.max.x
	}
	if n.pos.y < c.min.y {
		d += c.min.y - n.pos.y
	}
	if n.pos.y > c.max.y {
		d += n.pos.y - c.max.y
	}
	if n.pos.z < c.min.z {
		d += c.min.z - n.pos.z
	}
	if n.pos.z > c.max.z {
		d += n.pos.z - c.max.z
	}
	return d <= n.scanRadius

	// for _, corner := range c.corners() {
	// 	if corner.dist(n.pos) <= n.scanRadius {
	// 		return true
	// 	}
	// }
	return false
}

func (c *cube) allPoints() []point {
	out := []point{}
	for x := c.min.x; x <= c.max.x; x++ {
		for y := c.min.y; y <= c.max.y; y++ {
			for z := c.min.z; z <= c.max.z; z++ {
				out = append(out, point{x, y, z})
			}
		}
	}
	return out
}

func (c *cube) touchingNanobots(nbs []nanobot) int {
	out := 0
	for _, nb := range nbs {
		if c.overlap(nb) {
			out++
		}
	}
	return out
}

func (c *cube) String() string {
	return fmt.Sprintf("%v, %v", c.min, c.max)
}

type nanobot struct {
	pos        point
	scanRadius int
}

func (n nanobot) String() string {
	return fmt.Sprintf("%d,%d,%d; r=%d", n.pos.x, n.pos.y, n.pos.z, n.scanRadius)
}

func parseLine(line string) nanobot {
	words := strings.Split(line, ", ")
	position := strings.Split(words[0], "=")
	positions := strings.Split(strings.Trim(position[1], "<>"), ",")
	x, _ := strconv.Atoi(positions[0])
	y, _ := strconv.Atoi(positions[1])
	z, _ := strconv.Atoi(positions[2])
	rad := strings.Split(words[1], "=")
	radius, _ := strconv.Atoi(rad[1])
	return nanobot{
		pos:        point{x, y, z},
		scanRadius: radius,
	}
}

type point struct {
	x, y, z int
}

func (p point) dist(p2 point) int {
	return abs(p.x-p2.x) + abs(p.y-p2.y) + abs(p.z-p2.z)
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.x, p.y, p.z)
}

func (p point) intersecting(nbs []nanobot) int {
	count := 0
	for _, nb := range nbs {
		if p.dist(nb.pos) <= nb.scanRadius {
			count++
		}
	}
	return count
}

type d2 struct {
	count int
	size  int
	dist  int
	cube  *cube
}
type d3 struct {
	count int
	size  *big.Int
	dist  int
	cube  *cube
}
type pq []*d2
type pq2 []*d3

func (p pq) Len() int      { return len(p) }
func (p pq) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p pq) Less(i, j int) bool {
	if p[i].size < p[j].size {
		return true
	}
	if p[i].size > p[j].size {
		return false
	}
	if p[i].count > p[j].count {
		return true
	}
	if p[i].count < p[j].count {
		return false
	}
	return p[i].dist < p[j].dist
}

func (p *pq) Push(x interface{}) {
	*p = append(*p, x.(*d2))
}

func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*p = old[0 : n-1]
	return item
}

func (p pq2) Len() int      { return len(p) }
func (p pq2) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p pq2) Less(i, j int) bool {
	if p[i].count > p[j].count {
		return true
	}
	if p[i].count < p[j].count {
		return false
	}
	if p[i].dist < p[j].dist {
		return true
	}
	if p[i].dist > p[j].dist {
		return false
	}
	if p[i].size.Cmp(p[j].size) == 1 {
		return true
	}
	if p[i].size.Cmp(p[j].size) == -1 {
		return false
	}
	return true
}

func (p *pq2) Push(x interface{}) {
	*p = append(*p, x.(*d3))
}

func (p *pq2) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*p = old[0 : n-1]
	return item
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
type data struct {
	count int
	size  int
	dist  int
}

type Item[T any] struct {
	Data     T
	Priority data
}

type PriorityQueue[T any] struct {
	Items []Item[T]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		Items: make([]Item[T], 0),
	}
}

func (p *PriorityQueue[T]) Add(item T, priority data) {
	p.Items = append(p.Items, Item[T]{Data: item, Priority: priority})
}

func (p *PriorityQueue[T]) Empty() bool {
	return len(p.Items) == 0
}

func (p *PriorityQueue[T]) Len() int {
	return len(p.Items)
}

func (p *PriorityQueue[T]) Pull() T {
	minPrio := p.Items[0].Priority
	min := p.Items[0].Data
	minidx := 0
	for i, item := range p.Items {
		if item.Priority.count < minPrio.count {
			minPrio = item.Priority
			min = item.Data
			minidx = i
			continue
		}
		if item.Priority.count > minPrio.count {
			continue
		}
		if item.Priority.size < item.Priority.size {
			minPrio = item.Priority
			min = item.Data
			minidx = i
			continue
		}
		if item.Priority.size > item.Priority.size {
			continue
		}
		if item.Priority.dist < minPrio.dist {
			minPrio = item.Priority
			min = item.Data
			minidx = i
			continue
		}
	}
	p.Items = append(p.Items[:minidx], p.Items[minidx+1:]...)
	return min
}
*/
