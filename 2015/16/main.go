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
	inputData := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	sues := map[int]*sue{}
	for _, line := range lines {
		data := map[string]int{}
		words := strings.Split(line, " ")
		num, _ := strconv.Atoi(strings.TrimSuffix(words[1], ":"))
		first := strings.TrimSuffix(words[2], ":")
		firstCount, _ := strconv.Atoi(strings.TrimSuffix(words[3], ","))
		data[first] = firstCount
		second := strings.TrimSuffix(words[4], ":")
		secondCount, _ := strconv.Atoi(strings.TrimSuffix(words[5], ","))
		data[second] = secondCount
		third := strings.TrimSuffix(words[6], ":")
		thirdCount, _ := strconv.Atoi(strings.TrimSuffix(words[7], ","))
		data[third] = thirdCount
		sues[num] = newSue(num, data)
	}
	for _, sue := range sues {
		if sue.isValid(inputData) {
			fmt.Println(sue.number)
		}
	}
}

type sue struct {
	number      int
	children    int
	cats        int
	samoyeds    int
	pomeranians int
	akitas      int
	vizslas     int
	goldfish    int
	trees       int
	cars        int
	perfumes    int
}

func (s *sue) isValid(conditions map[string]int) bool {
	for key, num := range conditions {
		if s.getAttr(key) == -1 {
			continue
		}
		if (key == "cats" || key == "trees") && s.getAttr(key) > num {
			continue
		}
		if (key == "pomeranians" || key == "goldfish") && s.getAttr(key) < num {
			continue
		}
		if s.getAttr(key) != num {
			return false
		}
	}
	return true
}

func (s *sue) getAttr(attr string) int {
	switch attr {
	case "number":
		return s.number
	case "children":
		return s.children
	case "cats":
		return s.cats
	case "samoyeds":
		return s.samoyeds
	case "pomeranians":
		return s.pomeranians
	case "akitas":
		return s.akitas
	case "vizslas":
		return s.vizslas
	case "goldfish":
		return s.goldfish
	case "trees":
		return s.trees
	case "cars":
		return s.cars
	case "perfumes":
		return s.perfumes
	}
	panic(fmt.Sprintf("well fuck %s", attr))
}

func newSue(num int, data map[string]int) *sue {
	s := &sue{
		number:      num,
		children:    -1,
		cats:        -1,
		samoyeds:    -1,
		pomeranians: -1,
		akitas:      -1,
		vizslas:     -1,
		goldfish:    -1,
		trees:       -1,
		cars:        -1,
		perfumes:    -1,
	}
	if n, ok := data["children"]; ok {
		s.children = n
	}
	if n, ok := data["cats"]; ok {
		s.cats = n
	}
	if n, ok := data["samoyeds"]; ok {
		s.samoyeds = n
	}
	if n, ok := data["pomeranians"]; ok {
		s.pomeranians = n
	}
	if n, ok := data["akitas"]; ok {
		s.akitas = n
	}
	if n, ok := data["vizslas"]; ok {
		s.vizslas = n
	}
	if n, ok := data["goldfish"]; ok {
		s.goldfish = n
	}
	if n, ok := data["trees"]; ok {
		s.trees = n
	}
	if n, ok := data["cars"]; ok {
		s.cars = n
	}
	if n, ok := data["perfumes"]; ok {
		s.perfumes = n
	}
	return s
}
