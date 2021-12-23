package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type status int

const (
	off status = iota
	on
)

func main() {
	restrictSize := false

	lines := internal.ReadInput()
	actions := []action{}
	for _, line := range lines {
		actions = append(actions, parseLine(line, restrictSize))
	}
	// sum := 0
	// for _, action := range actions {
	// 	if action.action == "on" {
	// 		sum += action.affectedCubes()
	// 		continue
	// 	}
	// 	sum -= action.affectedCubes()
	// }
	// fmt.Println(sum)

	cubes := []*cuboid{}
	for _, a := range actions {
		cubes = append(cubes, newCube(a.xspan, a.yspan, a.zspan, a.status))
	}
	// if there is no item in the queue
	// add the cube to the queue
	// get the next cube to add
	// check it against every cube in the queue
	// if there is an overlap
	// split the cube already in the queue and add the pieces back to the queue
	q := internal.NewQueue[*cuboid]()
	q.Enqueue(cubes[0])
	for i := 1; i < len(cubes); i++ {
		newQueue := internal.NewQueue[*cuboid]()
		for !q.Empty() {
			cur := q.Dequeue()
			cuboids := cur.split(cubes[i])
			for _, c := range cuboids {
				newQueue.Enqueue(c)
			}
		}
		newQueue.Enqueue(cubes[i])
		q = newQueue
	}
	sum := 0
	for _, c := range q.Internal() {
		//		fmt.Println(c)
		if c.status == on {
			sum += c.volume()
		}
	}
	fmt.Println(sum)
}

type span struct {
	min, max int
}

func (s span) String() string {
	return fmt.Sprintf("%d..%d", s.min, s.max)
}
func (s span) overlap(s2 span) bool {
	// | s  |  |  s2  |
	if s.max < s2.min {
		return false
	}
	// | s2  |  | s  |
	if s2.max < s.min {
		return false
	}
	return true
}

type cuboid struct {
	xspan, yspan, zspan span
	status              status
}

func (c *cuboid) String() string {
	return fmt.Sprintf("[%d] x=%v,y=%v,z=%v (%d)", c.status, c.xspan, c.yspan, c.zspan, c.volume())
}

func newCube(xspan, yspan, zspan span, status status) *cuboid {
	return &cuboid{
		xspan:  xspan,
		yspan:  yspan,
		zspan:  zspan,
		status: status,
	}
}

func (c *cuboid) overlap(c2 *cuboid) bool {
	return c.xspan.overlap(c2.xspan) && c.yspan.overlap(c2.yspan) && c.zspan.overlap(c2.zspan)
}

func (c *cuboid) volume() int {
	// the +1 is the fact that the ranges are inclusive
	return (c.xspan.max - c.xspan.min + 1) * (c.yspan.max - c.yspan.min + 1) * (c.zspan.max - c.zspan.min + 1)
}

// split should split cuboid c and return c splits with c2 intact
func (c *cuboid) split(c2 *cuboid) []*cuboid {
	if !c.overlap(c2) {
		return []*cuboid{c}
	}
	out := make([]*cuboid, 0)
	// add the above cuboid
	if c.yspan.max > c2.yspan.max {
		yspan := span{c2.yspan.max + 1, c.yspan.max}
		out = append(out, newCube(c.xspan, yspan, c.zspan, c.status))
	}

	// add the below cuboid
	if c.yspan.min < c2.yspan.min {
		yspan := span{c.yspan.min, c2.yspan.min - 1}
		out = append(out, newCube(c.xspan, yspan, c.zspan, c.status))
	}

	// add the right cuboid
	if c.xspan.max > c2.xspan.max {
		xspan := span{c2.xspan.max + 1, c.xspan.max}
		yspan := span{max(c.yspan.min, c2.yspan.min), min(c.yspan.max, c2.yspan.max)}
		out = append(out, newCube(xspan, yspan, c.zspan, c.status))
	}

	// add the left cuboid
	if c.xspan.min < c2.xspan.min {
		xspan := span{c.xspan.min, c2.xspan.min - 1}
		yspan := span{max(c.yspan.min, c2.yspan.min), min(c.yspan.max, c2.yspan.max)}
		out = append(out, newCube(xspan, yspan, c.zspan, c.status))
	}

	// add the back cuboid
	if c.zspan.max > c2.zspan.max {
		xspan := span{max(c2.xspan.min, c.xspan.min), min(c2.xspan.max, c.xspan.max)}
		yspan := span{max(c2.yspan.min, c.yspan.min), min(c2.yspan.max, c.yspan.max)}
		zspan := span{c2.zspan.max + 1, c.zspan.max}
		out = append(out, newCube(xspan, yspan, zspan, c.status))
	}

	// add the front cuboid
	if c.zspan.min < c2.zspan.min {
		xspan := span{max(c2.xspan.min, c.xspan.min), min(c2.xspan.max, c.xspan.max)}
		yspan := span{max(c2.yspan.min, c.yspan.min), min(c2.yspan.max, c.yspan.max)}
		zspan := span{c.zspan.min, c2.zspan.min - 1}
		out = append(out, newCube(xspan, yspan, zspan, c.status))
	}

	return out
	/*
		on x=10..12,y=10..12,z=10..12
		on x=11..13,y=11..13,z=11..13
	*/
	/*
		output should be 3 cuboids
		the bottom face
		x=10..12,y=10..10,z=10..12
		the left face
		x=10..10,y=10..12,z=10..12
		the front face
		x=11..12,y=11..12,z=10..10
	*/
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseLine(line string, restrict bool) action {
	words := strings.Split(line, " ")
	out := action{status: off}
	if words[0] == "on" {
		out.status = on
	}
	vars := strings.Split(words[1], ",")
	minx, maxx := parseVar(vars[0])
	miny, maxy := parseVar(vars[1])
	minz, maxz := parseVar(vars[2])
	if restrict {
		if minx < -50 {
			minx = -50
		}
		if miny < -50 {
			miny = -50
		}
		if minz < -50 {
			minz = -50
		}
		if maxx > 50 {
			maxx = 50
		}
		if maxy > 50 {
			maxy = 50
		}
		if maxz > 50 {
			maxz = 50
		}
	}
	out.xspan.min = minx
	out.yspan.min = miny
	out.zspan.min = minz
	out.xspan.max = maxx
	out.yspan.max = maxy
	out.zspan.max = maxz
	return out
}

func parseVar(in string) (int, int) {
	vs := strings.Split(in, "=")
	nums := strings.Split(vs[1], "..")
	min, _ := strconv.Atoi(nums[0])
	max, _ := strconv.Atoi(nums[1])
	return min, max
}

type action struct {
	status              status
	xspan, yspan, zspan span
}
