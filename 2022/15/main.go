package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()

	num := 3364986

	for i := num - 1; i < num+2; i++ {
		spans := spanslice{}
		for _, line := range input {
			sensor := lineToPoint(line)
			// take the beacon dist
			if sensor.beaconDist() < sensor.distanceToY(i) {
				continue
			}
			distToBeacon := sensor.beaconDist()
			distToY := sensor.distanceToY(i)

			lateral := distToBeacon - distToY
			left := internal.Point{sensor.loc.X - lateral, i}
			right := internal.Point{sensor.loc.X + lateral, i}
			spans = append(spans, span{low: left.X, high: right.X})
		}

		sort.Sort(spans)
		//		fmt.Println(spans)
		x := 0
		for {
			if x+1 >= len(spans) {
				break
			}
			if spans[x].contains(spans[x+1]) {
				//			fmt.Println("before", spans, len(spans), spans[:i+1])
				spans = append(spans[:x+1], spans[x+2:]...)
				//			fmt.Println("after", spans, len(spans))
				continue
			}
			if spans[x].connected(spans[x+1]) {
				spans[x].high = spans[x+1].high
				spans = append(spans[:x+1], spans[x+2:]...)
				continue
			}
			x++
		}
		if i%100000 == 0 {
			fmt.Println(i, spans)
		}
		if len(spans) > 1 {
			fmt.Println((spans[0].high+1)*4000000 + i)
			fmt.Println(i, spans)
		}
	}
}

func lineToPoint(line string) sensor {
	//Sensor at x=3907621, y=2895218: closest beacon is at x=3790542, y=2949630
	re := regexp.MustCompile(`x=(-?\d*), y=(-?\d*)`)
	out := re.FindAllStringSubmatch(line, -1)
	locx, _ := strconv.Atoi(out[0][1])
	locy, _ := strconv.Atoi(out[0][2])
	beaconx, _ := strconv.Atoi(out[1][1])
	becaony, _ := strconv.Atoi(out[1][2])
	return sensor{
		loc:    internal.Point{X: locx, Y: locy},
		beacon: internal.Point{X: beaconx, Y: becaony},
	}
}

type sensor struct {
	loc    internal.Point
	beacon internal.Point
}

func (s sensor) beaconDist() int {
	return s.loc.ManhattanDistance(s.beacon)
}

func (s sensor) distanceToY(y int) int {
	return s.loc.ManhattanDistance(internal.Point{X: s.loc.X, Y: y})
}

type span struct {
	low, high int
}

func (s span) contains(s2 span) bool {
	return s.low <= s2.low && s.high >= s2.high
}

func (s span) connected(s2 span) bool {
	return s.high >= s2.low && s2.high >= s.high && s.low <= s2.low
}

type spanslice []span

func (s spanslice) Len() int { return len(s) }
func (s spanslice) Less(i, j int) bool {
	if s[i].low < s[j].low {
		return true
	}
	if s[i].low == s[j].low {
		if s[i].high < s[j].high {
			return true
		}
	}
	return false
}
func (s spanslice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
