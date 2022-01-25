package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	orbits := map[string][]string{}
	orbiting := map[string]string{}
	allNames := map[string]struct{}{}

	for _, line := range lines {
		names := strings.Split(line, ")")
		allNames[names[0]] = struct{}{}
		allNames[names[1]] = struct{}{}
		if _, ok := orbits[names[0]]; !ok {
			orbits[names[0]] = make([]string, 0)
		}
		orbits[names[0]] = append(orbits[names[0]], names[1])
		orbiting[names[1]] = names[0]
	}
	leaves := []string{}
	for n := range allNames {
		if _, ok := orbits[n]; !ok {
			leaves = append(leaves, n)
		}
	}
	my := path("YOU", orbiting)
	end := path("SAN", orbiting)
	minimumOrbitalTransfers := 0
	for i := 0; i < len(my); i++ {
		if my[i] != end[i] {
			fmt.Println("top node", my[i])
			minimumOrbitalTransfers = len(my) - i + len(end) - i - 2
			fmt.Println("steps to YOU", len(my)-i)
			fmt.Println("steps to SAN", len(end)-i)
			break
		}
	}

	fmt.Println(minimumOrbitalTransfers)

	// sum := 0
	// for node := range allNames {
	// 	orbits := count(node, orbiting)
	// 	sum += orbits
	// }
	// fmt.Println(sum)
	//	fmt.Println(count("COM", orbits))
}

// com -> []string
// for each string count the orbits

func count(node string, orbiting map[string]string) int {
	fmt.Println(node)
	if node == "COM" {
		return 0
	}
	return 1 + count(orbiting[node], orbiting)
}

func path(start string, orbiting map[string]string) []string {
	if start == "COM" {
		return []string{"COM"}
	}
	return append(path(orbiting[start], orbiting), start)
}

// 1707 too low (part 1)
// 335 too high (part 2)
