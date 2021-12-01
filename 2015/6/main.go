package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	off = iota
	on
	toggle
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	l := &lights{}
	for _, line := range lines {
		action, x1, y1, x2, y2 := parseLine(line)
		l.act(action, x1, y1, x2, y2)
		fmt.Println(l.lit())
		os.Exit(0)
	}
}

func parseLine(in string) (int, int, int, int, int) {
	action := -1
	if in[:6] == "toggle" {
		action = toggle
		in = in[7:]
	}
	if in[:8] == "turn off" {
		action = off
		in = in[9:]
	}
	if in[:7] == "turn on" {
		action = on
		in = in[8:]
	}
	parts := strings.Split(in, " through ")
	x1, y1 := numPair(parts[0])
	x2, y2 := numPair(parts[1])
	return action, x1, y1, x2, y2
}

func numPair(pair string) (int, int) {
	nums := strings.Split(pair, ",")
	x, _ := strconv.Atoi(nums[0])
	y, _ := strconv.Atoi(nums[1])
	return x, y
}

type lights [1000][1000]int

func (l *lights) lights(action, x, y int) {
	if action == on {
		l[y][x] = on
		return
	}
	if action == off {
		l[y][x] = off
		return
	}
	// toggle
	if l[y][x] == on {
		l[y][x] = off
		return
	}
	l[y][x] = on
}

func (l *lights) act(action, x1, y1, x2, y2 int) {
	for i := x1; i <= x2; i++ {
		for j := y1; j <= y2; j++ {
			l.lights(action, i, j)
		}
	}
}

func (l *lights) lit() int {
	c := 0
	for _, row := range l {
		for _, light := range row {
			if light == on {
				c++
			}
		}
	}
	return c
}
