package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	lights := []*light{}
	for _, line := range lines {
		lights = append(lights, parseLine(line))
	}
	i := 1
	for {
		fmt.Printf("----------%d------------\n", i)
		move(lights)
		i++
	}
}

func move(lights []*light) {
	minx, miny := math.MaxInt, math.MaxInt
	maxx, maxy := math.MinInt, math.MinInt
	for _, l := range lights {
		l.move()
		if l.position.x < minx {
			minx = l.position.x
		}
		if l.position.x > maxx {
			maxx = l.position.x
		}
		if l.position.y < miny {
			miny = l.position.y
		}
		if l.position.y > maxy {
			maxy = l.position.y
		}
	}
	if maxx-minx < 100 && maxy-miny < 100 {
		data := map[point]string{}
		for j := miny; j <= maxy; j++ {
			for i := minx; i <= maxx; i++ {
				data[point{i, j}] = "."
			}
		}
		for _, l := range lights {
			data[l.position] = "#"
		}
		var out strings.Builder
		for j := miny; j <= maxy; j++ {
			for i := minx; i <= maxx; i++ {
				out.WriteString(data[point{i, j}])
			}
			out.WriteString("\n")
		}
		fmt.Println(out.String())
		time.Sleep(100 * time.Millisecond)
		//		fmt.Printf("[%d,%d]->[%d,%d]\n", minx, miny, maxx, maxy)
	}
}

// position=< 10255, -50258> velocity=<-1,  5>

func parseLine(line string) *light {
	words := strings.Split(line, "<")
	//["position="," 10255, -50258> velocity=", "-1, 5>"]
	first := strings.Fields(strings.TrimSpace(words[1]))
	// ["10255,", "-50258> veolcity="]
	xpos, _ := strconv.Atoi(strings.TrimSuffix(first[0], ","))
	next := strings.Split(first[1], ">")
	ypos, _ := strconv.Atoi(next[0])
	second := strings.Split(strings.TrimSpace(words[2]), ",")
	// ["-1", " 5>"]
	vx, _ := strconv.Atoi(second[0])
	vy, _ := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(second[1], ">")))
	return &light{
		position: point{
			x: xpos, y: ypos,
		},
		vx: vx, vy: vy,
	}
}

type light struct {
	position point
	vx, vy   int
}

func (l *light) move() {
	l.position = l.position.move(l.vx, l.vy)
}

type point struct {
	x, y int
}

func (p point) move(vx, vy int) point {
	return point{
		x: p.x + vx,
		y: p.y + vy,
	}
}

// 10085 too low
