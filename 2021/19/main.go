package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type logger struct {
	debug bool
}

func (l *logger) Println(a ...interface{}) {
	if !l.debug {
		return
	}
	fmt.Println(a...)
}
func (l *logger) Printf(f string, a ...interface{}) {
	if !l.debug {
		return
	}
	fmt.Printf(f, a...)
}

type orientation int

const (
	zero orientation = iota
	one
	two
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	eleven
	twelve
	thirteen
	fourteen
	fifteen
	sixteen
	seventeen
	eighteen
	nineteen
	twewnty
	twentyone
	twentytwo
	twentythree
)

func main() {
	lines := internal.ReadInput()

	scanners := []*scanner{}
	var s *scanner
	for _, line := range lines {
		if strings.HasPrefix(line, "---") {
			if s != nil {
				scanners = append(scanners, s)
			}
			s = newScanner()
			s.name = strings.TrimSpace(strings.Trim(line, "---"))
			continue
		}
		coords := strings.Split(line, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])
		s.beacons = append(s.beacons, point{x, y, z})
	}
	if s != nil {
		scanners = append(scanners, s)
	}

	// set 0 as the origin
	scanners[0].setAbsolute(point{0, 0, 0}, zero)
	relativeToOrigin := map[point]*scanner{
		{0, 0, 0}: scanners[0],
	}

	squeue := internal.NewQueue[*scanner]()
	squeue.Enqueue(scanners[0])

	for !squeue.Empty() {
		scanner := squeue.Dequeue()
		for j := 0; j < len(scanners); j++ {
			// skip if we've already done it
			if relativeToOrigin[scanners[j].location].name == scanners[j].name {
				continue
			}
			set := overlappingBeacons(scanner, scanners[j])
			if !filterSize(set) {
				continue
			}
			fmt.Println(scanner.name, "to", scanners[j].name)
			pairs := assignPairs(set)
			origin := calculateOrigin(pairs)
			for _, v := range origin {
				relativeToOrigin[v.location] = scanners[j]
				scanners[j].setAbsolute(v.location, v.orientation)
			}
			squeue.Enqueue(scanners[j])
		}
	}
	out := uniquePoints(relativeToOrigin)
	fmt.Println(len(out))
	fmt.Println(largestManhattanDistance(relativeToOrigin))
	//	fmt.Println(relativeToOrigin)

	// this is a bit too much...don't need all the calculations
	// This is the actual code to calculate overlapping spots
	// for i := 0; i < len(scanners)-1; i++ {
	// 	for j := i + 1; j < len(scanners); j++ {
	// 		set := overlappingBeacons(scanners[i], scanners[j])
	// 		if !filterSize(set) {
	// 			continue
	// 		}
	// 		// // Print out the overlapping beacons
	// 		// pairs := assignPairs(set)
	// 		// origin := calculateOrigin(pairs)
	// 		// for _, v := range origin {
	// 		// 	fmt.Println(v)
	// 		// }
	// 		fmt.Printf("-- %s vs %s --\n", scanners[i].name, scanners[j].name)
	// 		for k := range set {
	// 			fmt.Println(k)
	// 		}
	// 		fmt.Println("--")
	// 	}
	// }

	// // find scanner origin with rotation
	// b := overlappingBeacons(scanners[0], scanners[1])
	// pairs := assignPairs(b)
	// origin := calculateOrigin(pairs)
	// for _, v := range origin {
	// 	fmt.Println(v)
	// }
	// for _, v := range b {
	// 	fmt.Println(v.segmentA, v.segmentB)
	// }

	// // finding origin of another scanner
	// possibleOrigins1 := internal.Set[scanner]{}
	// pointa := point{-485, -357, 347}
	// pointb := point{553, 889, -390}
	// for i, p := range rotations(pointb) {
	// 	possibleOrigin := pointa.add(p.inverse())
	// 	s := scanner{}
	// 	s.orientation = orientation(i)
	// 	s.location = possibleOrigin
	// 	possibleOrigins1.Insert(possibleOrigin.String(), s)
	// }
	// possibleOrigins2 := internal.Set[scanner]{}
	// pointa1 := point{544, -627, -890}
	// pointb1 := point{-476, 619, 847}
	// for i, p := range rotations(pointb1) {
	// 	possibleOrigin := pointa1.add(p.inverse())
	// 	s := scanner{}
	// 	s.orientation = orientation(i)
	// 	s.location = possibleOrigin
	// 	possibleOrigins2.Insert(possibleOrigin.String(), s)
	// }
	// for _, v := range possibleOrigins1.Intersect(possibleOrigins2) {
	// 	fmt.Println(v.location, v.orientation)
	// }

	// {point0, [11 points] from scanner0} {point0, [11 points] from scanner1}
	// where the distance to each point is ==

	// [a] [b,c,d,e,f,g] [A], [B,C,D,E,F,G]
	// [a,b][c,d,e,f,g] [A,C] [B,D,E,F,G]
	// pick one point (a) from scanner0
	//		for each other point in scanner0 (b,c,d,e,f,g,...)
	//			pick a point in scanner1 (A)
	//			for each other point in scanner1 (B,C,D,E,F,G,...)
	//				if dist(a,b) == dist(A,B)
	//					pick another point from scanner0
	//					same checks as above
	//				else
	//					continue

	//	pick two points from scanner0
	//		pick one point from scanner1
	//		for each other point in scanner1
	//			if the distance of two pionts from scanner 1  == distance of points from scanner 0
	//				pick another point from scanner 0

	// for _, b := range scanners[0].beacons {
	// 	for _, p := range rotations(b) {
	// 		fmt.Println(p)
	// 	}
	// }
}

// scanner that contains beacons
// scanner has 24 sets each containing every orientation of the beacon set

func (p point) inverse() point {
	return point{-p.x, -p.y, -p.z}
}

func rotations(p point) []point {
	return []point{
		{p.x, p.y, p.z},   // facing +x up: +y		// zero
		{p.x, -p.z, p.y},  // facing +x up -z
		{p.x, -p.y, -p.z}, // facing +x up: -y
		{p.x, p.z, -p.y},  // facing +x up: +z

		{p.y, p.z, p.x},   // facing +y up: +z		// four
		{p.y, -p.x, p.z},  // facing +y up: -x
		{p.y, -p.z, -p.x}, // facing +y up: -z
		{p.y, p.x, -p.z},  // facing +y up: +x

		{p.z, p.x, p.y},   // facing +z up: +x		// eight
		{p.z, -p.y, p.x},  // facing +z up: -y
		{p.z, -p.x, -p.y}, // facing +z up: -x
		{p.z, p.y, -p.x},  // facing +z up: +y

		{-p.x, -p.z, -p.y}, // facing -x up -z		// twelve
		{-p.x, p.y, -p.z},  // facing -x up: +y
		{-p.x, p.z, p.y},   // facing -x up +z
		{-p.x, -p.y, p.z},  // facing -x up -y

		{-p.y, -p.x, -p.z}, //facing -y up -x		// sixteen
		{-p.y, p.z, -p.x},  //facing -y up +z
		{-p.y, p.x, p.z},   //facing -y up +x
		{-p.y, -p.z, p.x},  //facing -y up -z

		{-p.z, -p.y, -p.x}, // facing -z up -y		// twenty
		{-p.z, p.x, -p.y},  // facing -z up +x
		{-p.z, p.y, p.x},   // facing -z up +y
		{-p.z, -p.x, p.y},  // facing -z up -x
	}
}

type scanner struct {
	name            string
	orientation     orientation
	location        point
	beacons         []point
	absoluteBeacons []point
}

func (s *scanner) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("--- %s (%d) ---\n", s.name, s.orientation))
	if len(s.absoluteBeacons) != 0 {
		for _, b := range s.absoluteBeacons {
			out.WriteString(fmt.Sprintf("%v\n", b))
		}
		return out.String()
	}
	for _, b := range s.beacons {
		out.WriteString(fmt.Sprintf("%v\n", b))
	}
	return out.String()
}

func (s *scanner) setAbsolute(loc point, orientation orientation) {
	s.location = loc
	s.orientation = orientation
	absoluteBeacons := make([]point, len(s.beacons))
	for i, b := range s.beacons {
		absoluteBeacons[i] = b.absolute(loc, orientation)
	}
	s.absoluteBeacons = absoluteBeacons
}

func newScanner() *scanner {
	return &scanner{
		beacons: make([]point, 0),
	}
}

type point struct {
	x, y, z int
}

func (p point) add(p2 point) point {
	return point{p.x + p2.x, p.y + p2.y, p.z + p2.z}
}

func (p point) absolute(origin point, orientation orientation) point {
	return origin.add(rotations(p)[orientation])
}

func (p point) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.x, p.y, p.z)
}

func (p point) dist(p2 point) float64 {
	return math.Sqrt(
		math.Pow(float64(p2.x)-float64(p.x), 2) +
			math.Pow(float64(p2.y)-float64(p.y), 2) +
			math.Pow(float64(p2.z)-float64(p.z), 2))
}

//	 ux: inputs scannera, scannerb
// recursion: input list of beacons from scanner a, list of beacons from scanner b
//            input list of beacon distances for scanner a == list of becacon distances for scanner b
// base case:
// 		if len(matched beacons) == 12
//
// [a] [b,c,d,e,f,g] [A], [B,C,D,E,F,G]
// [a,b][c,d,e,f,g] [A,C] [B,D,E,F,G]
// pick one point (a) from scanner0
//		for each other point in scanner0 (b,c,d,e,f,g,...)
//			pick a point in scanner1 (A)
//			for each other point in scanner1 (B,C,D,E,F,G,...)
//				if dist(a,b) == dist(A,B)
//					pick another point from scanner0
//					same checks as above
//				else
//					continue

func pointDistances(points []point) []distanceData {
	out := []distanceData{}
	for i := 0; i < len(points)-1; i++ {
		for j := i + 1; j < len(points); j++ {
			out = append(out, distanceData{
				a:    points[i],
				b:    points[j],
				dist: points[i].dist(points[j]),
			})
		}
	}
	sort.Sort(distanceDatas(out))
	return out
}

func matchBeacons(scannerA, scannerB *scanner) {
	availableA := scannerA.beacons
	availableB := scannerB.beacons
	for i := 0; i < len(availableA)-1; i++ {
		for j := i + 1; j < len(availableA); j++ {
			for k := 0; k < len(availableB)-1; k++ {
				for l := k + 1; l < len(availableB); l++ {
					a1 := availableA[i]
					a2 := availableA[j]
					b1 := availableB[k]
					b2 := availableB[l]
					if a1.dist(a2) == b1.dist(b2) {
						fmt.Println(a1, a2, b1, b2)
					}
				}
			}
		}
	}
}

func copyList(in []point) []point {
	cp := make([]point, len(in))
	copy(cp, in)
	return cp
}

type distanceData struct {
	a, b point
	dist float64
}

type segment struct {
	endpointA, endpointB point
}

type equalDistancePairs struct {
	segmentA, segmentB segment
	distance           float64
}

type doublePoint struct {
	a, b point
}
type distanceDatas []distanceData

func (d distanceDatas) Len() int {
	return len(d)
}
func (d distanceDatas) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
func (d distanceDatas) Less(i, j int) bool {
	return d[i].dist < d[j].dist
}

type points []point

func (p points) Len() int {
	return len(p)
}
func (p points) Less(i, j int) bool {
	return p[i].x < p[j].x
}
func (p points) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func findEqualDistances(a, b []distanceData) []equalDistancePairs {
	out := make([]equalDistancePairs, 0)
	i := 0
	j := 0
	for {
		if i == len(a) && j == len(b) {
			return out
		}
		if a[i].dist == b[j].dist {
			out = append(out, equalDistancePairs{
				segmentA: segment{
					endpointA: a[i].a,
					endpointB: a[i].b,
				},
				segmentB: segment{
					endpointA: b[j].a,
					endpointB: b[j].b,
				},
				distance: a[i].dist,
			})
		}
		if a[i].dist > b[j].dist {
			// this means that all of the remaining items in b
			// are smaller/bigger than a
			if j == len(b)-1 {
				return out
			}
			j++
			continue
		}
		if i == len(a)-1 {
			return out
		}
		i++
	}
}

func overlappingBeacons(a, b *scanner) internal.Set[equalDistancePairs] {
	da := pointDistances(a.absoluteBeacons)
	db := pointDistances(b.beacons)
	distanceSet := make(internal.Set[equalDistancePairs])
	for _, d := range findEqualDistances(da, db) {
		distanceSet.Insert(d.segmentA.endpointA.String(), d)
		distanceSet.Insert(d.segmentA.endpointB.String(), d)
	}
	return distanceSet
}

func filterSize(in internal.Set[equalDistancePairs]) bool {
	return len(in) >= 12
}

func assignPairs(set internal.Set[equalDistancePairs]) map[point]point {
	log := &logger{false}
	// map point a to potentials in point b
	possibilities := map[point]map[point]int{}
	for _, v := range set {
		if _, ok := possibilities[v.segmentA.endpointA]; !ok {
			possibilities[v.segmentA.endpointA] = map[point]int{}
		}
		possibilities[v.segmentA.endpointA][v.segmentB.endpointA]++
		possibilities[v.segmentA.endpointA][v.segmentB.endpointB]++
		if _, ok := possibilities[v.segmentA.endpointB]; !ok {
			possibilities[v.segmentA.endpointB] = map[point]int{}
		}
		possibilities[v.segmentA.endpointB][v.segmentB.endpointA]++
		possibilities[v.segmentA.endpointB][v.segmentB.endpointB]++
	}
	out := map[point]point{}

	for k, v := range possibilities {
		max := 0
		var actual point
		equals := 1
		log.Println("poss", k, v)
		for k2, v2 := range v {
			if v2 == max {
				equals++
			}
			if v2 > max {
				max = v2
				actual = k2
			}
		}
		log.Println(equals, len(v))
		if equals == len(v) {
			continue
		}
		out[k] = actual
	}
	return out
}

func calculateOrigin(in map[point]point) internal.Set[scanner] {
	if len(in) == 0 {
		return nil
	}
	sets := []internal.Set[scanner]{}
	for pointa, pointb := range in {
		possibleOrigins := internal.Set[scanner]{}
		for i, p := range rotations(pointb) {
			possibleOrigin := pointa.add(p.inverse())
			s := scanner{}
			s.orientation = orientation(i)
			s.location = possibleOrigin
			possibleOrigins.Insert(possibleOrigin.String(), s)
		}
		sets = append(sets, possibleOrigins)
	}
	out := sets[0]
	for i := 1; i < len(sets); i++ {
		out = out.Intersect(sets[i])
	}
	return out
}

func uniquePoints(scanners map[point]*scanner) []point {
	uniquePoints := internal.Set[point]{}
	for _, s := range scanners {
		for _, p := range s.absoluteBeacons {
			uniquePoints.Insert(p.String(), p)
		}
	}
	up := make(points, 0)
	//	fmt.Println(len(uniquePoints))
	for _, v := range uniquePoints {
		up = append(up, v)
	}
	sort.Sort(up)
	return up
}

func largestManhattanDistance(scanners map[point]*scanner) int {
	max := 0
	for _, s := range scanners {
		for _, s1 := range scanners {
			d := s.location.manhattanDistance(s1.location)
			if d > max {
				max = d
			}
		}
	}
	return max
}

func (p point) manhattanDistance(p2 point) int {
	return int(math.Abs(float64(p.x-p2.x))) + int(math.Abs(float64(p.y-p2.y))) + int(math.Abs(float64(p.z-p2.z)))
}

/*
-618,-824,-621
-537,-823,-458
-447,-329,318
404,-588,-901
544,-627,-890
528,-643,409
-661,-816,-575
390,-675,-793
423,-701,434
-345,-311,381
459,-707,401
-485,-357,347

{-485 -357 347} {553 889 -390}



{{-618 -824 -621} {423 -701 434}} {{686 422 578} {-355 545 -477}}
{{528 -643 409} {-537 -823 -458}} {{605 423 415} {-460 603 -452}}
{{-447 -329 318} {544 -627 -890}} {{515 917 -361} {-476 619 847}}
{{404 -588 -901} {-485 -357 347}} {{-336 658 858} {553 889 -390}}
{{-485 -357 347} {544 -627 -890}} {{-476 619 847} {553 889 -390}}
{{528 -643 409} {-661 -816 -575}} {{-460 603 -452} {729 430 532}}
{{528 -643 409} {-661 -816 -575}} {{-460 603 -452} {729 430 532}}
{{-661 -816 -575} {459 -707 401}} {{729 430 532} {-391 539 -444}}
{{390 -675 -793} {-485 -357 347}} {{-322 571 750} {553 889 -390}}
{{-345 -311 381} {544 -627 -890}} {{-476 619 847} {413 935 -424}}
{{528 -643 409} {-618 -824 -621}} {{686 422 578} {-460 603 -452}}
{{-485 -357 347} {544 -627 -890}} {{-476 619 847} {553 889 -390}}

// scanenr a == 0,0,0 w/ orientation x,y,z
// the point {-485 -357 347} is also known as {553 889 -390} relative to some other point & orientation
// look at each of the 24 possible scanner b values and find those positions relative to scanner a by
//
a: {-485 -357 347}
b: {553 889 -390}

a: {544 -627 -890}
b: {-476 619 847}

// given a position, and another way to express that position based on a different origin
// return all possilbe 24 origin points
//
*/

/*
input pos: (686,422,578) ->  (-686,422,-578)
should be : (-618,-824,-621)

scanner1 loc (68,-1246,-43)
scanenr1 orientation: {-p.x, p.y, -p.z},

68+-686
-1246+244
-43+-578
*/
