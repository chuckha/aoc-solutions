package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		data := scanner.Text()
		data = strings.TrimSpace(data)
		if data == "" {
			continue
		}
		lines = append(lines, data)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	reindeers := map[string]*reindeer{}
	for _, line := range lines {
		words := strings.Split(line, " ")
		speed, _ := strconv.Atoi(words[3])
		travel, _ := strconv.Atoi(words[6])
		rest, _ := strconv.Atoi(words[13])
		reindeers[words[0]] = newReindeer(words[0], speed, travel, rest)
	}
	maxTime := 2503

	// for _, r := range reindeers {
	// 	fmt.Println(r)
	// }

	for i := 0; i < maxTime; i++ {
		for _, r := range reindeers {
			r.move()
		}
		var winner *reindeer
		for _, r := range reindeers {
			if winner == nil {
				winner = r
				continue
			}
			if r.distance > winner.distance {
				winner = r
			}
		}
		for _, r := range reindeers {
			if r.distance == winner.distance {
				r.points++
			}
			//			fmt.Printf("%d %s: %d\n", r.time, r.name, r.points)
		}
	}

	biggest := 0
	for _, r := range reindeers {
		if r.points > biggest {
			biggest = r.points
		}
	}
	fmt.Println(biggest)

	/* part 1
	distances := []int{}
	for _, r := range reindeers {
		intervalTime := r.restingTime + r.travelTime
		intervals := maxTime / intervalTime
		remaining := maxTime % intervalTime
		distance := intervals * r.travelTime * r.maxSpeed
		if r.travelTime < remaining {
			distance += r.travelTime * r.maxSpeed
		} else {
			distance += remaining * r.maxSpeed
		}
		distances = append(distances, distance)
	}
	farthest := 0
	for _, d := range distances {
		if d > farthest {
			farthest = d
		}
	}
	fmt.Println(farthest)
	*/
}

type reindeer struct {
	name        string
	maxSpeed    int // km/s
	travelTime  int
	restingTime int
	points      int
	distance    int
	time        int
}

func newReindeer(name string, speed, travel, rest int) *reindeer {
	return &reindeer{
		name:        name,
		maxSpeed:    speed,
		travelTime:  travel,
		restingTime: rest,
		points:      0,
	}
}

func (r *reindeer) move() {
	interval := r.travelTime + r.restingTime
	loc := r.time % interval
	if loc < r.travelTime {
		//fmt.Printf("%s: (%04d) interval: %d, loc: %02d MOVING\n", r.name, r.time, interval, loc)
		r.distance += r.maxSpeed
	} else {
		//fmt.Printf("%s: (%04d) interval: %d, loc: %02d RESTING\n", r.name, r.time, interval, loc)
	}
	r.time += 1
}
