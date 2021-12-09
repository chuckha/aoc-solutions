package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
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
	heightMap := map[point]*data{}
	rows := len(lines)
	cols := len(lines[0])
	for y, line := range lines {
		for x, num := range strings.Split(line, "") {
			d, _ := strconv.Atoi(num)
			heightMap[point{x, y}] = &data{
				val:   d,
				point: point{x, y},
			}
		}
	}
	riskLevel := 0

	basins := []point{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			data := heightMap[point{x, y}]
			// check up
			if y != 0 {
				if data.val >= heightMap[point{x, y - 1}].val {
					continue
				}
			}
			// check right
			if x != cols-1 {
				if data.val >= heightMap[point{x + 1, y}].val {
					continue
				}
			}
			// check down
			if y != rows-1 {
				if data.val >= heightMap[point{x, y + 1}].val {
					continue
				}
			}
			// check left
			if x != 0 {
				if data.val >= heightMap[point{x - 1, y}].val {
					continue
				}
			}
			basins = append(basins, point{x, y})
			riskLevel += heightMap[point{x, y}].val + 1
		}
	}
	basinSizes := []int{}
	for _, b := range basins {
		basinSizes = append(basinSizes, count(cols-1, rows-1, heightMap, b.x, b.y))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(basinSizes)))
	fmt.Println(basinSizes[0] * basinSizes[1] * basinSizes[2])
}

type point struct {
	x, y int
}

func getUnvisitedNeighbor(maxX, maxY int, heightMap map[point]*data, x, y int) []*data {
	neighbors := []*data{}
	if x != 0 {
		neighbors = append(neighbors, heightMap[point{x - 1, y}])
	}
	if x != maxX {
		neighbors = append(neighbors, heightMap[point{x + 1, y}])
	}
	if y != maxY {
		neighbors = append(neighbors, heightMap[point{x, y + 1}])
	}
	if y != 0 {
		neighbors = append(neighbors, heightMap[point{x, y - 1}])
	}
	filtered := []*data{}
	for _, neighbor := range neighbors {
		if neighbor.visited {
			continue
		}
		if neighbor.val == 9 {
			continue
		}
		filtered = append(filtered, neighbor)
	}
	return filtered
}

func count(maxX, maxY int, heightMap map[point]*data, x, y int) int {
	neighbors := getUnvisitedNeighbor(maxX, maxY, heightMap, x, y)
	q := internal.NewQueue[*data]()
	for _, n := range neighbors {
		q.Enqueue(n)
	}

	c := 1
	cur := heightMap[point{x, y}]
	for !q.Empty() {
		nabe := q.Dequeue()
		if nabe.visited {
			continue
		}
		if nabe.val > cur.val {
			c++
			nabe.visited = true
		}
		neighbors := getUnvisitedNeighbor(maxX, maxY, heightMap, nabe.point.x, nabe.point.y)
		for _, n := range neighbors {
			q.Enqueue(n)
		}
		nabe.visited = true
	}
	return c
}

// get unvisited neighbors
// for nabe in unvisitedneighbors
//

type data struct {
	val     int
	visited bool
	point   point
	isbasin bool
}

func (d *data) String() string {
	return fmt.Sprintf("%d %v", d.val, d.point)
}
