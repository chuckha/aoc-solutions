package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

// target area: x=241..275, y=-75..-49

func main() {
	input := internal.ReadInput()
	words := strings.Split(input[0], " ")
	xData := strings.Split(words[2], "=")
	xVals := strings.Split(strings.TrimSuffix(xData[1], ","), "..")
	xmin, _ := strconv.Atoi(xVals[0])
	xmax, _ := strconv.Atoi(xVals[1])
	yData := strings.Split(words[3], "=")
	yVals := strings.Split(yData[1], "..")
	ymin, _ := strconv.Atoi(yVals[0])
	ymax, _ := strconv.Atoi(yVals[1])

	//	fmt.Println(xmin, xmax, ymin, ymax)

	largesty := 0
	goodVels := []vel{}
	for y := -400; y < 400; y++ {
		for x := -400; x < 400; x++ {
			success, biggesty := position(xmin, xmax, ymin, ymax, 0, x, 0, y, 0)
			if success {
				goodVels = append(goodVels, vel{x, y})
				if biggesty > largesty {
					largesty = biggesty
				}
			}
		}
	}
	fmt.Println(len(goodVels))
}

func position(minx, maxx, miny, maxy, ix, ivx, iy, ivy, biggesty int) (bool, int) {
	if inside(minx, maxx, miny, maxy, ix, iy) {
		return true, biggesty
	}
	if ix > maxx || iy < miny {
		return false, 0
	}
	newx := ix + ivx
	newy := iy + ivy
	if ivx > 0 {
		ivx -= 1
	}
	ivy -= 1
	if iy > biggesty {
		biggesty = iy
	}
	return position(minx, maxx, miny, maxy, newx, ivx, newy, ivy, biggesty)
}
func inside(minx, maxx, miny, maxy, x, y int) bool {
	return x >= minx && x <= maxx && y >= miny && y <= maxy
}

//x=20..30, y=-10..-5

/*
starts at 0,0
vel x
vel y
acc x (if vel x >0; -1; if vel <0; 1 else 0)
acc y -1

intial velocity = 7,2
0:  0, 0 v=7,2;acc=-1,-1
1:  7, 2 v=6,1;acc=-1,-1
2:  13,3 v=5,0;acc=-1,-1
3:  18,3 v=4,-1;acc=-1,-1
4:	22,2 v=3,-2; acc=-1,-1
5:	25,0 v=2,-3; acc=-1,-1
6:  27,-3 v=1,-4; acc=-1,-1
7:  28,-7 v=0,-5; acc=0,-1
t| x | y   vx|vy vxi=7;vyi=2
0  0 | 0    7|2
1  7 | 2	6|1
2  13| 3	5|0
3  18| 3	4|-1
4  22| 2	3|-2
5  25| 0	2|-3
6  27|-3	1|-4
7  28|-7	0|-5

x = vxi + t
func position(ix, ivx, iy, ivy int) {
	newx := ix + ivx
	newy := iy + ivy
	if ivx > 0 {
		ivx -= 1
	}
	ivy -= 1
	position(newx, ivx, newy, ivy)
}


vy1-y^2 = (vx1)+x^2-ai
6 + 0 = 6
6 + 1 = 7
6 + 4 = 10
6 + 9 = 15

x=20..30, y=-10..-5

y = -x^2


*/

type vel struct {
	x, y int
}
