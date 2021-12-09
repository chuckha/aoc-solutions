package main

import (
	"bufio"
	"fmt"
	"os"
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
	lights := map[point]*light{}
	maxY := len(lines)
	maxX := len(lines[0])
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				lights[point{x, y}] = &light{1}
				continue
			}
			lights[point{x, y}] = &light{0}
		}
	}
	turnCornersOn(lights, maxX, maxY)
	lightsCopy := make(map[point]*light, len(lights))
	for k, v := range lights {
		lightsCopy[k] = &light{v.on}
	}
	iterations := 100
	for i := 0; i < iterations; i++ {
		// printLights(lights, maxX, maxY)
		// fmt.Println()
		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				onNeighbors := 0
				self := lights[point{x, y}]
				selfCopy := lightsCopy[point{x, y}]
				if n, ok := lights[point{x, y - 1}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x + 1, y - 1}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x + 1, y}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x + 1, y + 1}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x, y + 1}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x - 1, y + 1}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x - 1, y}]; ok {
					onNeighbors += n.on
				}
				if n, ok := lights[point{x - 1, y - 1}]; ok {
					onNeighbors += n.on
				}
				if self.on == 1 {
					if onNeighbors != 2 && onNeighbors != 3 {
						selfCopy.on = 0
					}
				}
				if self.on == 0 {
					if onNeighbors == 3 {
						selfCopy.on = 1
					}
				}
			}
		}
		// turn on the four corners again

		lights, lightsCopy = lightsCopy, lights
		for k, v := range lights {
			lightsCopy[k] = &light{v.on}
		}
		turnCornersOn(lights, maxX, maxY)
	}
	// printLights(lights, maxX, maxY)
	// fmt.Println()

	count := 0
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			count += lights[point{x, y}].on
		}
	}
	fmt.Println(count)
}

type point struct {
	x, y int
}

type light struct {
	on int
}

func printLights(lights map[point]*light, maxX, maxY int) {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if lights[point{x, y}].on == 1 {
				fmt.Printf("#")
				continue
			}
			fmt.Printf(".")
		}
		fmt.Println()
	}
}

func turnCornersOn(lights map[point]*light, maxX, maxY int) {
	lights[point{0, 0}] = &light{1}
	lights[point{0, maxY - 1}] = &light{1}
	lights[point{maxX - 1, maxY - 1}] = &light{1}
	lights[point{maxX - 1, 0}] = &light{1}
}
