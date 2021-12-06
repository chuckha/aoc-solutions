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
	agg := map[point]int{}

	for _, line := range lines {
		p1, p2 := newPointsFromLine(line)
		// part 1 ignore diagonal lines
		// if p1.x != p2.x && p1.y != p2.y {
		// 	continue
		// }
		for _, point := range p1.pointsUntil(p2) {
			agg[*point]++
		}
	}
	sum := 0
	for _, v := range agg {
		if v >= 2 {
			sum++
		}
	}
	fmt.Println(sum)
}

func newPointsFromLine(line string) (*point, *point) {
	points := strings.Split(line, " -> ")
	coordA := strings.Split(points[0], ",")
	coordB := strings.Split(points[1], ",")
	x1, _ := strconv.Atoi(coordA[0])
	y1, _ := strconv.Atoi(coordA[1])
	x2, _ := strconv.Atoi(coordB[0])
	y2, _ := strconv.Atoi(coordB[1])
	return &point{x1, y1}, &point{x2, y2}
}

type point struct {
	x, y int
}

func (p *point) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}
func (p *point) pointsUntil(p2 *point) []*point {
	out := []*point{p}
	cur := p
	for {
		newPoint := &point{x: cur.x, y: cur.y}
		if p.x > p2.x {
			newPoint.x -= 1
		}
		if p.x < p2.x {
			newPoint.x += 1
		}
		if p.y > p2.y {
			newPoint.y -= 1
		}
		if p.y < p2.y {
			newPoint.y += 1
		}
		out = append(out, newPoint)
		cur = newPoint
		if cur.x == p2.x && cur.y == p2.y {
			return out
		}
	}
}
