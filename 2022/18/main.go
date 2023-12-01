package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()

	sqs := map[pt]struct{}{}
	for _, line := range input {
		sqs[newunit(line)] = struct{}{}
	}
	seen := map[pt]struct{}{}
	min, max := extremes(sqs)
	fmt.Println(min, max)
	air := min.fill(seen, sqs, min, max)

	total := 0
	for _, sq := range air {
		total += sq.contributing(sqs)
	}
	fmt.Println(total)
}

type pt struct {
	x, y, z int
}

func (p pt) contributing(lookup map[pt]struct{}) int {
	touching := 0
	for _, n := range p.neighbors() {
		if _, ok := lookup[n]; ok {
			touching++
		}
	}
	return touching
}

func (p pt) neighbors() []pt {
	return []pt{
		{p.x - 1, p.y, p.z},
		{p.x + 1, p.y, p.z},

		{p.x, p.y - 1, p.z},
		{p.x, p.y + 1, p.z},

		{p.x, p.y, p.z - 1},
		{p.x, p.y, p.z + 1},
	}
}

func newunit(line string) pt {
	pts := strings.Split(line, ",")
	x, _ := strconv.Atoi(pts[0])
	y, _ := strconv.Atoi(pts[1])
	z, _ := strconv.Atoi(pts[2])
	return pt{x, y, z}
}

func (p pt) fill(seen, lookup map[pt]struct{}, min, max pt) []pt {
	out := []pt{p}
	seen[p] = struct{}{}
	for _, n := range p.neighbors() {
		if n.x < min.x || n.x > max.x || n.y < min.y || n.y > max.y || n.z < min.z || n.z > max.z {
			continue
		}
		if _, ok := lookup[n]; ok {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		out = append(out, n.fill(seen, lookup, min, max)...)
	}
	return out
}

func extremes(sqs map[pt]struct{}) (pt, pt) {
	minx := math.MaxInt
	miny := math.MaxInt
	minz := math.MaxInt
	maxx := math.MinInt
	maxy := math.MinInt
	maxz := math.MinInt
	for sq := range sqs {
		if sq.x < minx {
			minx = sq.x
		}
		if sq.x > maxx {
			maxx = sq.x
		}
		if sq.y < miny {
			miny = sq.y
		}
		if sq.y > maxy {
			maxy = sq.y
		}
		if sq.z < minz {
			minz = sq.z
		}
		if sq.z > maxz {
			maxz = sq.z
		}
	}
	return pt{minx - 1, miny - 1, minz - 1}, pt{maxx + 1, maxy + 1, maxz + 1}
}
