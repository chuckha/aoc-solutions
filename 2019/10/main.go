package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

func main() {
	lines := input.GetInput(2019, 10)
	grid := internal.NewGridFromInput(lines)
	points := []internal.Point{}
	origin := internal.Point{23, 19}
	for p, v := range grid.Data {
		if p == origin {
			continue
		}
		if v != "#" {
			continue
		}
		points = append(points, p)
	}

	fmt.Println("part 1", findNumAsteroidsVisibleFromBestSpaceStation(grid))
	expl(points, origin)
}

func expl(pts []internal.Point, origin internal.Point) {
	count := 0
	for len(pts) != 0 {
		fmt.Println(pts)
		round := splitUp(pts, origin)
		//		fmt.Println("next round", round)
		// (9,34) == too low
		for _, p := range round {
			count++
			fmt.Printf("the %dth asteroid to be destroyed is at %v\n", count, p)
			for i, p2 := range pts {
				if p2 == p {
					pts = append(pts[:i], pts[i+1:]...)
				}
			}
		}
	}
}

type closest struct {
	pts []internal.Point
	og  internal.Point
}

func (a closest) Len() int      { return len(a.pts) }
func (a closest) Swap(i, j int) { a.pts[i], a.pts[j] = a.pts[j], a.pts[i] }
func (a closest) Less(i, j int) bool {
	return a.pts[i].ManhattanDistance(a.og) < a.pts[j].ManhattanDistance(a.og)
}

type angles struct {
	pts []internal.Point
	og  internal.Point
}

func (a angles) Len() int      { return len(a.pts) }
func (a angles) Swap(i, j int) { a.pts[i], a.pts[j] = a.pts[j], a.pts[i] }
func (a angles) Less(i, j int) bool {
	if angle(a.og, a.pts[i]) > angle(a.og, a.pts[j]) {
		return true
	}
	if angle(a.og, a.pts[i]) < angle(a.og, a.pts[j]) {
		return false
	}
	return a.og.ManhattanDistance(a.pts[i]) < a.og.ManhattanDistance(a.pts[j])
}

type angles2 struct {
	pts []internal.Point
	og  internal.Point
}

func (a angles2) Len() int      { return len(a.pts) }
func (a angles2) Swap(i, j int) { a.pts[i], a.pts[j] = a.pts[j], a.pts[i] }
func (a angles2) Less(i, j int) bool {
	if angle(a.og, a.pts[i]) < angle(a.og, a.pts[j]) {
		return true
	}
	if angle(a.og, a.pts[i]) > angle(a.og, a.pts[j]) {
		return false
	}
	return a.og.ManhattanDistance(a.pts[i]) < a.og.ManhattanDistance(a.pts[j])
}
func angle(origin, pt2 internal.Point) float64 {
	a := abs(pt2.X - origin.X)
	a = a * a
	b := abs(pt2.Y - origin.Y)
	b = b * b
	c := math.Sqrt(float64(a*a + b*b))
	//	theta := math.Asin(float64(b) / float64(c))
	return round(math.Sin(float64(b)/float64(c)), 0.000000000001)

}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func splitUp(in []internal.Point, origin internal.Point) []internal.Point {
	oneRound := []internal.Point{}
	out := []internal.Point{}
	//	first  return all points above
	above := []internal.Point{}
	for _, point := range in {
		if point.X == origin.X && point.Y < origin.Y {
			above = append(above, point)
		}
	}
	c := closest{
		pts: above,
		og:  origin,
	}
	sort.Sort(c)
	out = append(out, c.pts...)
	if len(c.pts) > 0 {
		oneRound = append(oneRound, c.pts[0])
	}
	// quadrant 1
	q1 := []internal.Point{}
	for _, point := range in {
		if point.X > origin.X && point.Y < origin.Y {
			q1 = append(q1, point)
		}
	}
	a := angles{
		pts: q1,
		og:  origin,
	}
	sort.Sort(a)
	for _, p := range a.pts {
		fmt.Println(p, angle(origin, p))
	}
	var last float64
	for _, p := range a.pts {
		if angle(origin, p) == last {
			continue
		}
		last = angle(origin, p)
		oneRound = append(oneRound, p)
	}
	// get right
	right := []internal.Point{}
	for _, point := range in {
		if point.X > origin.X && point.Y == origin.Y {
			right = append(right, point)
		}
	}
	c = closest{
		pts: right,
		og:  origin,
	}
	sort.Sort(c)
	if len(c.pts) > 0 {
		oneRound = append(oneRound, c.pts[0])
	}
	// q4
	q4 := []internal.Point{}
	for _, point := range in {
		if point.X > origin.X && point.Y > origin.Y {
			q4 = append(q4, point)
		}
	}
	a2 := angles2{
		pts: q4,
		og:  origin,
	}
	sort.Sort(a2)
	last = 0.0
	for _, p := range a2.pts {
		if angle(origin, p) == last {
			continue
		}
		last = angle(origin, p)
		oneRound = append(oneRound, p)
	}
	// down
	down := []internal.Point{}
	for _, point := range in {
		if point.X == origin.X && point.Y > origin.Y {
			down = append(down, point)
		}
	}
	c = closest{
		pts: down,
		og:  origin,
	}
	sort.Sort(c)
	if len(c.pts) > 0 {
		oneRound = append(oneRound, c.pts[0])
	}
	// q3
	q3 := []internal.Point{}
	for _, point := range in {
		if point.X < origin.X && point.Y > origin.Y {
			q3 = append(q3, point)
		}
	}
	a = angles{
		pts: q3,
		og:  origin,
	}
	sort.Sort(a)
	last = 0.0
	for _, p := range a.pts {
		if angle(origin, p) == last {
			continue
		}
		last = angle(origin, p)
		oneRound = append(oneRound, p)
	}
	// left
	left := []internal.Point{}
	for _, point := range in {
		if point.X < origin.X && point.Y == origin.Y {
			left = append(left, point)
		}
	}
	c = closest{
		pts: left,
		og:  origin,
	}
	sort.Sort(c)
	if len(c.pts) > 0 {
		oneRound = append(oneRound, c.pts[0])
	}
	// q2
	q2 := []internal.Point{}
	for _, point := range in {
		if point.X < origin.X && point.Y < origin.Y {
			q2 = append(q2, point)
		}
	}
	a2 = angles2{
		pts: q2,
		og:  origin,
	}
	sort.Sort(a2)
	last = 0.0
	for _, p := range a2.pts {
		if angle(origin, p) == last {
			continue
		}
		last = angle(origin, p)
		oneRound = append(oneRound, p)
	}

	return oneRound
}

func explode(g *internal.Grid) {
	p := internal.Point{8, 3}
	as := &asteroidStats{
		data:   make([]asteroidStat, 0),
		origin: p,
	}
	for p2, v := range g.Data {
		if p2 == p {
			continue
		}
		if v != "#" {
			continue
		}
		as.data = append(as.data, newStat(p, p2))
	}
	sort.Sort(as)
	for i, s := range as.data {
		fmt.Println(i+1, s)
	}
	// for i, p := range round {
	// 	fmt.Println(i+1, p)
	// }
	fmt.Println(g)
}
func remove(g *internal.Grid, toRemove []internal.Point) {
	for _, r := range toRemove {
		g.Set(r.X, r.Y, ".")
	}
}

func findNumAsteroidsVisibleFromBestSpaceStation(g *internal.Grid) int {
	// for each point in the grid, iterate through every other point in the grid
	data := map[internal.Point]int{}
	biggest := math.MinInt
	pt := internal.Point{}
	for cur := range g.Data {
		if g.At(cur.X, cur.Y) != "#" {
			continue
		}
		checkedSlopes := map[internal.Point]bool{}
		checked := map[internal.Point]bool{}
		cp := g.Copy()
		count := 0
		for p := range g.Data {
			if p == cur {
				continue
			}
			dx, dy := xychange(cur, p)
			if _, ok := checkedSlopes[internal.Point{dx, dy}]; ok {
				continue
			}
			x := cur.X
			y := cur.Y

			seenOne := false
			// add dx,dy
			for x >= g.MinX && x <= g.Length && y >= g.MinY && y <= g.Height {
				x += dx
				y += dy
				if checked[internal.Point{x, y}] {
					continue
				}
				checked[internal.Point{x, y}] = true
				if g.At(x, y) == "#" {
					if !seenOne {
						cp.Set(x, y, string(internal.Green("#")))
						count++
						seenOne = true
						continue
					}
					if seenOne {
						cp.Set(x, y, string(internal.Red("#")))
					}
				}
			}
			x = cur.X
			y = cur.Y
			seenOne = false
			for x >= g.MinX && x <= g.Length && y >= g.MinY && y <= g.Height {
				x -= dx
				y += dy
				if checked[internal.Point{x, y}] {
					continue
				}
				checked[internal.Point{x, y}] = true
				if g.At(x, y) == "#" {
					if !seenOne {
						count++
						cp.Set(x, y, string(internal.Green("#")))
						seenOne = true
						continue
					}
					if seenOne {
						cp.Set(x, y, string(internal.Red("#")))
					}
				}
			}
			x = cur.X
			y = cur.Y
			seenOne = false
			for x >= g.MinX && x <= g.Length && y >= g.MinY && y <= g.Height {
				x += dx
				y -= dy
				if checked[internal.Point{x, y}] {
					continue
				}
				checked[internal.Point{x, y}] = true
				if g.At(x, y) == "#" {

					if !seenOne {
						count++
						cp.Set(x, y, string(internal.Green("#")))
						seenOne = true
						continue
					}
					if seenOne {
						cp.Set(x, y, string(internal.Red("#")))
					}
				}
			}
			x = cur.X
			y = cur.Y
			seenOne = false
			for x >= g.MinX && x <= g.Length && y >= g.MinY && y <= g.Height {
				x -= dx
				y -= dy
				if checked[internal.Point{x, y}] {
					continue
				}
				checked[internal.Point{x, y}] = true
				if g.At(x, y) == "#" {
					if !seenOne {
						count++
						cp.Set(x, y, string(internal.Green("#")))
						seenOne = true
						continue
					}
					if seenOne {
						cp.Set(x, y, string(internal.Red("#")))
					}
				}
			}
			data[cur] = count
			checkedSlopes[internal.Point{dx, dy}] = true
			checkedSlopes[internal.Point{-dx, dy}] = true
			checkedSlopes[internal.Point{dx, -dy}] = true
			checkedSlopes[internal.Point{-dx, -dy}] = true
		}
		if count > biggest {
			pt = cur
			biggest = count
		}
		//		fmt.Println(colorPrint(cp, cur, count))
	}
	fmt.Println(pt)
	return biggest
}

func colorPrint(g *internal.Grid, start internal.Point, count int) string {
	var out strings.Builder
	g.Set(start.X, start.Y, string(internal.Yellow(fmt.Sprintf("%d", count))))
	for j := g.MinY; j <= g.Height; j++ {
		for i := g.MinX; i <= g.Length; i++ {
			out.WriteString(g.At(i, j))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func xychange(p1, p2 internal.Point) (int, int) {
	xchange := p2.X - p1.X
	ychange := p2.Y - p1.Y
	f := gcd(xchange, ychange)
	return xchange / f, ychange / f
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type asteroidStat struct {
	theta float64
	pos   internal.Point
}

func deg2rad(in float64) float64 {
	return in * math.Pi / 180
}

//    /
//   / (9,0) -> (3, 1) a = 3, b = 1, c = sqrt(a*a +b*b), angle == tan(b/a)
// i actually want an angle from oirign (point)
func newStat(origin, dest internal.Point) asteroidStat {
	a := abs(dest.X - origin.X)
	a = a * a
	b := abs(dest.Y - origin.Y)
	b = b * b
	c := math.Sqrt(float64(a*a + b*b))
	//	theta := math.Asin(float64(b) / float64(c))
	theta := math.Sin(float64(b) / float64(c))
	sign := 1.0
	if dest.Y > origin.Y {
		sign = -1.0
	}
	// first quadrant is positive
	// if dest.X < origin.X {
	// 	sign = -1
	// }
	// if dest.X <= origin.X && dest.Y > origin.Y {
	// 	sign = -1
	// }
	return asteroidStat{
		theta: sign * theta,
		pos:   dest,
	}
}

func (a asteroidStat) String() string {
	return fmt.Sprintf("%v %v", a.theta, a.pos)
}

type asteroidStats struct {
	origin internal.Point
	data   []asteroidStat
}

func (a *asteroidStats) Len() int      { return len(a.data) }
func (a *asteroidStats) Swap(i, j int) { a.data[i], a.data[j] = a.data[j], a.data[i] }
func (a *asteroidStats) Less(i, j int) bool {
	if a.data[i].theta > a.data[j].theta {
		return true
	}
	if a.data[i].theta < a.data[j].theta {
		return false
	}
	return a.data[i].pos.X > a.data[j].pos.X
}

/*
.#....#####...#..
##...##.#####..##
##...#...#.#####.
..#.....#...###..
..#.#.....#....##
*/
