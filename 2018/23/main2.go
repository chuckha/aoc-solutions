package main

// import (
// 	"fmt"
// 	"math"
// 	"sort"
// 	"strconv"
// 	"strings"

// 	"github.com/chuckha/aoc-solutions/internal/input"
// )

// func main2() {
// 	lines := internal.ReadInput()

// 	nbs := make([]nanobot, 0)
// 	for _, line := range lines {
// 		nbs = append(nbs, parseLine(line))
// 	}
// 	sort.Sort(nanobots(nbs))
// 	count := 0
// 	strongest := nbs[0]
// 	buckets := newBuckets()
// 	//	exit := 0
// 	for _, nb := range nbs {
// 		if strongest.pos.dist(nb.pos) <= strongest.scanRadius {
// 			count++
// 		}
// 		// part 2
// 		//		fmt.Println("about to add", rg, "from", nb)
// 		//		fmt.Println(nb)
// 		fmt.Println(buckets.root)
// 		// exit++
// 		// if exit > 8 {
// 		// 	os.Exit(0)
// 		// }
// 	}
// 	fmt.Println("part one:", findLargestBucket(buckets.root))
// }

// type buckets struct {
// 	root  *internal.DLL[bucket]
// 	debug bool
// }

// func newBuckets() *buckets {
// 	return &buckets{debug: true}
// }
// func (b *buckets) Debugln(line ...any) {
// 	if b.debug {
// 		fmt.Println(line...)
// 	}
// }

// func (b *buckets) addRange(pr pointRange) {
// 	b.Debugln("adding range..", pr)
// 	if b.root == nil || b.root.Len() == 0 {
// 		bkt := newBucket(pr, 1)
// 		b.root = internal.NewDLL(bkt)
// 		return
// 	}
// 	cur := b.root
// 	for ; cur != nil; cur = cur.Next {
// 		// falls completely within a bucket
// 		if cur.Data.lowLowHighHigh(pr) {
// 			b.Debugln("range falls entirely within a bucket")
// 			bkts := cur.Data.middleSection(pr.xspan.low, pr.xspan.high)
// 			cur.Data = bkts[0]
// 			for _, bucket := range bkts[1:] {
// 				cur = cur.InsertAfter(bucket)
// 			}
// 			return
// 		}
// 		// low is in the middle of this range
// 		if cur.Data.lowLowMidHigh(pr) {
// 			b.Debugln("low of pr in bucket")
// 			bkts := cur.Data.splitAbove(pr.xspan.low)
// 			cur.Data = bkts[0]
// 			cur.InsertAfter(bkts[1])
// 			remainingRange := pr.shrink(bkts[1].pr.xspan.high+1, 0, 0)
// 			b.addRange(remainingRange)
// 			return
// 		}
// 		// low is before current
// 		// if cur.Data.pr.xspan.high == pr.xspan.low && cur.Data.pr.xspan.low < pr.xspan.low && pr.xspan.high > cur.Data.pr.xspan.high {
// 		// 	b.Debugln("range starts at bucket start but then exanpds past bucket")
// 		// 	bkts := cur.Data.splitAbove(cur.Data.pr.xspan.high)
// 		// 	cur.Data = bkts[0]
// 		// 	cur.InsertAfter(bkts[1])
// 		// 	remainingRange := pr.shrink(pr.xspan.low+1, 0, 0)
// 		// 	b.addRange(remainingRange)
// 		// 	return
// 		// }
// 		// range is entirely behind lowest bucket
// 		if b.root.Data.entirelyAbove(pr) {
// 			cur = cur.InsertBefore(newBucket(pr, 1))
// 			b.root = cur
// 			return
// 		}
// 		if b.root.Data.highLowHighHigh(pr) {
// 			bkts := b.root.Data.lowSplitBelow(pr)
// 			cur.Data = bkts[0]
// 			for _, bucket := range bkts[1:] {
// 				cur = cur.InsertAfter(bucket)
// 			}
// 			return
// 		}
// 		// the ran low is inside the range but high is above
// 		// if cur.Data.lowx < pr.xmin && pr.xmin < cur.Data.highx && pr.xmax >= cur.Data.highx {
// 		// 	bkts := cur.Data.splitAbove(cur.Data.highx)
// 		// 	fmt.Println("backets", bkts)
// 		// 	return
// 		// }
// 		// a superset encompasing the whole bucket

// 		if cur.Data.sameLowLowHigh(pr) {
// 			cur.Data.count++
// 			remainingRange := pr.shrink(cur.Data.pr.xspan.high+1, 0, 0)
// 			b.addRange(remainingRange)
// 			return
// 		}
// 		// a subset starting at exactly the start
// 		if cur.Data.sameLowHighHigh(pr) {
// 			bkts := cur.Data.splitBelow(pr.xspan.high)
// 			cur.Data = bkts[0]
// 			cur.InsertAfter(bkts[1])
// 			return
// 		}
// 		// An exact range match, just add one to the bucket
// 		if cur.Data.sameLowSameHigh(pr) {
// 			cur.Data.count = cur.Data.count + 1
// 			return
// 		}
// 		// catch a single point
// 	}
// 	// at the last range and nothing has happened, therefore the new range is entirely
// 	// larger than all buckets
// 	// get to the end
// 	cur = b.root
// 	for cur.Next != nil {
// 		cur = cur.Next
// 	}
// 	cur.InsertAfter(newBucket(pr, 1))
// }

// type bucket struct {
// 	pr    pointRange
// 	count int
// }

// func (b bucket) entirelyAbove(pr pointRange) bool {
// 	return b.pr.xspan.low > pr.xspan.high
// }
// func (b bucket) lowLowHighHigh(pr pointRange) bool {
// 	return b.pr.xspan.low < pr.xspan.low && b.pr.xspan.high > pr.xspan.high // &&
// 	// b.pr.yspan.low < pr.yspan.low && b.pr.yspan.high > pr.yspan.high
// }
// func (b bucket) lowLowMidHigh(pr pointRange) bool {
// 	return b.pr.xspan.low < pr.xspan.low && b.pr.xspan.high > pr.xspan.low && b.pr.xspan.high < pr.xspan.high
// }
// func (b bucket) sameLowLowHigh(pr pointRange) bool {
// 	return b.pr.xspan.low == pr.xspan.low && b.pr.xspan.high < pr.xspan.high
// }
// func (b bucket) sameLowHighHigh(pr pointRange) bool {
// 	return b.pr.xspan.low == pr.xspan.low && b.pr.xspan.high > pr.xspan.high
// }
// func (b bucket) sameLowSameHigh(pr pointRange) bool {
// 	return b.pr.xspan.low == pr.xspan.low && b.pr.xspan.high == pr.xspan.high
// }
// func (b bucket) highLowHighHigh(pr pointRange) bool {
// 	return b.pr.xspan.low > pr.xspan.low && b.pr.xspan.high > pr.xspan.high
// }

// func (b bucket) String() string {
// 	return fmt.Sprintf("[%v](%d)", b.pr, b.count)
// }

// // 1,50 split at 5 1->5 5->50
// func (b bucket) splitAbove(at int) [2]bucket {
// 	low := pointRange{xspan: newXSpan(b.pr.xspan.low, at-1)}
// 	high := pointRange{xspan: newXSpan(at, b.pr.xspan.high)}
// 	return [2]bucket{
// 		newBucket(low, b.count),
// 		newBucket(high, b.count+1),
// 	}
// }
// func (b bucket) splitBelow(at int) [2]bucket {
// 	low := pointRange{xspan: newXSpan(b.pr.xspan.low, at)}
// 	high := pointRange{xspan: newXSpan(at+1, b.pr.xspan.high)}
// 	return [2]bucket{
// 		newBucket(low, b.count+1),
// 		newBucket(high, b.count),
// 	}
// }
// func (b bucket) lowSplitBelow(pr pointRange) [3]bucket {
// 	low := pointRange{xspan: newXSpan(pr.xspan.low, b.pr.xspan.low-1)}
// 	mid := pointRange{xspan: newXSpan(b.pr.xspan.low, pr.xspan.high)}
// 	high := pointRange{xspan: newXSpan(pr.xspan.high+1, b.pr.xspan.high)}
// 	return [3]bucket{
// 		newBucket(low, b.count),
// 		newBucket(mid, b.count+1),
// 		newBucket(high, b.count),
// 	}
// }

// func (b bucket) middleSection(lowx, highx int) [3]bucket {
// 	low := pointRange{xspan: newXSpan(b.pr.xspan.low, lowx-1)}
// 	middle := pointRange{xspan: newXSpan(lowx, highx)}
// 	high := pointRange{xspan: newXSpan(highx+1, b.pr.xspan.high)}
// 	return [3]bucket{
// 		newBucket(low, b.count),
// 		newBucket(middle, b.count+1),
// 		newBucket(high, b.count),
// 	}
// }

// func newBucket(pr pointRange, count int) bucket {
// 	return bucket{
// 		pr:    pr,
// 		count: count,
// 	}
// }
// func findLargestBucket(cur *internal.DLL[bucket]) bucket {
// 	largest := bucket{}

// 	for cur != nil {
// 		if cur.Data.count > largest.count {
// 			largest = cur.Data
// 		}
// 		cur = cur.Next
// 	}
// 	return largest
// }

// // 612 is too low (off by one)

// type nanobots []nanobot

// func (n nanobots) Len() int           { return len(n) }
// func (n nanobots) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
// func (n nanobots) Less(i, j int) bool { return n[i].scanRadius > n[j].scanRadius }

// type nanobot struct {
// 	pos        point
// 	scanRadius int
// }

// // func (n nanobot) pointRange() pointRange {
// // 	return n.pos.ranges(n.scanRadius)
// // }

// func (n nanobot) String() string {
// 	return fmt.Sprintf("%d,%d,%d; r=%d", n.pos.x, n.pos.y, n.pos.z, n.scanRadius)
// }

// func parseLine(line string) nanobot {
// 	words := strings.Split(line, ", ")
// 	position := strings.Split(words[0], "=")
// 	positions := strings.Split(strings.Trim(position[1], "<>"), ",")
// 	x, _ := strconv.Atoi(positions[0])
// 	y, _ := strconv.Atoi(positions[1])
// 	z, _ := strconv.Atoi(positions[2])
// 	rad := strings.Split(words[1], "=")
// 	radius, _ := strconv.Atoi(rad[1])
// 	return nanobot{
// 		pos:        point{x, y, z},
// 		scanRadius: radius,
// 	}
// }

// type point struct {
// 	x, y, z int
// }

// func (p point) dist(p2 point) int {
// 	return int(math.Abs(float64(p.x-p2.x))) + int(math.Abs(float64(p.y-p2.y))) + int(math.Abs(float64(p.z-p2.z)))
// }

// func (p point) ranges(radius int) pointRange {
// 	return pointRange{
// 		xspan: newXSpan(p.x-radius, p.x+radius),
// 		yspan: newYSpan(p.y-radius, p.y+radius),
// 		zspan: newZSpan(p.z-radius, p.z+radius),
// 	}
// }

// type pointRange struct {
// 	xspan, yspan, zspan span
// }

// func (p pointRange) shrink(x, y, z int) pointRange {
// 	return pointRange{
// 		xspan: newXSpan(x, p.xspan.high),
// 	}
// }
// func (p pointRange) String() string {
// 	return fmt.Sprintf("%v,%v,%v", p.xspan, p.yspan, p.zspan)
// }

// type span struct {
// 	axis      string
// 	low, high int
// }

// func newXSpan(low, high int) span {
// 	return newSpan("x", low, high)
// }
// func newYSpan(low, high int) span {
// 	return newSpan("y", low, high)
// }
// func newZSpan(low, high int) span {
// 	return newSpan("z", low, high)
// }
// func newSpan(axis string, low, high int) span {
// 	return span{axis, low, high}
// }

// func (s span) String() string {
// 	if s.low == s.high {
// 		return fmt.Sprintf("%s=%d", s.axis, s.low)
// 	}
// 	return fmt.Sprintf("%s=%d..%d", s.axis, s.low, s.high)
// }

// /*
// pos=<10,12,12>, r=2
// pos=<12,14,12>, r=2
// pos=<16,12,12>, r=4
// pos=<14,14,14>, r=6
// pos=<50,50,50>, r=200
// pos=<10,10,10>, r=5

// 8..12, 10..14, 10..14
// 10
// */

// /*
// pos=<10,12,12>, r=2
// pos=<12,14,12>, r=2
// pos=<16,12,12>, r=4
// pos=<14,14,14>, r=6
// pos=<50,50,50>, r=200
// pos=<10,10,10>, r=5
// */
