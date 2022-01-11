package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	points := []point{}
	for _, line := range lines {
		points = append(points, pointFromLine(line))
	}

	// create buckets of points
	// for each new point, check each bucket

	// at the end, check every bucket against every other point in every other bucket and move buckets around
	buckets := []*bucket{}
	for i, p := range points {
		added := false
		var firstMatch *bucket
		toMerge := []*bucket{}
		unmerged := []*bucket{}
		for _, bucket := range buckets {
			if bucket == nil {
				panic("how did it come to this")
			}
			if !added && bucket.contains(p) {
				firstMatch = bucket
				bucket.add(p)
				added = true
				continue
			}
			if added && bucket.contains(p) {
				toMerge = append(toMerge, bucket)
				continue
			}
			unmerged = append(unmerged, bucket)
		}
		if firstMatch != nil {
			// fmt.Println("before merge", len(firstMatch.points))
			// sum := 0
			// for _, b := range toMerge {
			// 	sum += len(b.points)
			// }
			// fmt.Println("adding", sum, "points")
			for _, b := range toMerge {
				firstMatch.merge(b)
			}
			// fmt.Println("after merge", len(firstMatch.points))
			buckets = append([]*bucket{firstMatch}, unmerged...)
		} else {
			fmt.Println(i, "didn't match any other buckets")
			nb := newBucket()
			nb.add(p)
			buckets = append(buckets, nb)
		}
	}
	sum := 0
	for _, b := range buckets {
		sum += len(b.points)
	}
	if sum != len(points) {
		fmt.Println(sum, len(points))
		panic("missed one")
	}
	for i := 0; i < len(buckets)-1; i++ {
		for j := i + 1; j < len(buckets); j++ {
			for _, p := range buckets[i].points {
				for _, p2 := range buckets[j].points {
					if p.dist(p2) <= 3 {
						fmt.Println("u wot m8", p, p2)
					}
				}
			}
		}
	}
	fmt.Println(len(buckets))
}

type point struct {
	x, y, z, w int
}

func pointFromLine(line string) point {
	nums := strings.Split(line, ",")
	x, _ := strconv.Atoi(nums[0])
	y, _ := strconv.Atoi(nums[1])
	z, _ := strconv.Atoi(nums[2])
	w, _ := strconv.Atoi(nums[3])
	return point{x, y, z, w}
}

func (p point) dist(p2 point) int {
	return abs(p.x-p2.x) + abs(p.y-p2.y) + abs(p.z-p2.z) + abs(p.w-p2.w)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// 459 too high
// 568 too high

type bucket struct {
	points []point
}

func newBucket() *bucket {
	return &bucket{points: make([]point, 0)}
}

func (b *bucket) contains(p point) bool {
	for _, point := range b.points {
		if p.dist(point) <= 3 {
			return true
		}
	}
	return false
}

func (b *bucket) add(p point) {
	b.points = append(b.points, p)
}

func (b *bucket) merge(b2 *bucket) {
	b.points = append(b.points, b2.points...)
}
