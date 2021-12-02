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
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	pos := &position2{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		inst, dist := parseLine(line)
		pos.move(inst, dist)
	}
	fmt.Println(pos.x * pos.depth)
}

func parseLine(line string) (string, int) {
	spl := strings.Split(strings.TrimSpace(line), " ")
	i, _ := strconv.Atoi(spl[1])
	return spl[0], i
}

type position2 struct {
	x, aim, depth int
}

func (p *position2) move(instruction string, dist int) {
	switch instruction {
	case "forward":
		p.x += dist
		p.depth += (p.aim * dist)
	case "down":
		p.aim += dist
	case "up":
		p.aim -= dist
	default:
		panic("unknown instruction: " + instruction)
	}
}

// part 1
type position struct {
	x, y int
}

func (p *position) move(instruction string, dist int) {
	switch instruction {
	case "forward":
		p.x += dist
	case "down":
		p.y += dist
	case "up":
		p.y -= dist
	default:
		panic("unknown instruction: " + instruction)
	}
}
