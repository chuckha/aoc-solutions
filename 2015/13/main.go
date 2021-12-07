package main

import (
	"bufio"
	"strconv"

	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	people := map[string]*person{}
	for _, line := range lines {
		info := strings.Split(line, " ")
		person, ok := people[info[0]]
		if !ok {
			person = newPerson(info[0])
			people[info[0]] = person
		}
		multiplier := 1
		if info[2] == "lose" {
			multiplier = -1
		}
		amount, _ := strconv.Atoi(info[3])
		amount = amount * multiplier
		near := strings.TrimSuffix(info[10], ".")
		person.addRule(near, amount)
	}
	me := &person{
		name:      "me",
		happiness: 0,
		rules:     map[string]int{},
	}
	for _, person := range people {
		person.addRule("me", 0)
		me.addRule(person.name, 0)
	}
	people["me"] = me
	allPeople := make([]string, 0)
	for name := range people {
		allPeople = append(allPeople, name)
	}
	biggest := 0
	for _, permutation := range allPermutations([]string{}, allPeople) {
		h := totalHappiness(people, permutation)
		if h > biggest {
			biggest = h
		}
	}
	fmt.Println(biggest)
}

func totalHappiness(data map[string]*person, perm []string) int {
	sum := 0
	for i := 0; i < len(perm); i++ {
		var left, right string
		if i == 0 {
			left = perm[len(perm)-1]
		} else {
			left = perm[i-1]
		}
		if i == len(perm)-1 {
			right = perm[0]
		} else {
			right = perm[i+1]
		}
		person := data[perm[i]]
		sum += person.rules[left]
		sum += person.rules[right]
	}
	return sum
}

type person struct {
	name      string
	happiness int
	rules     map[string]int
}

func (p *person) String() string {
	out := fmt.Sprintf("%s (%d)\n", p.name, p.happiness)
	for name, amt := range p.rules {
		out += fmt.Sprintf("\t%s -> %d\n", name, amt)
	}
	return out
}

func newPerson(name string) *person {
	return &person{
		name:      name,
		happiness: 0,
		rules:     make(map[string]int),
	}
}

func (p *person) addRule(name string, happiness int) {
	p.rules[name] = happiness
}

func allPermutations(fixed, unfixed []string) [][]string {
	if len(unfixed) == 0 {
		swap := make([]string, len(fixed))
		copy(swap, fixed)
		return [][]string{swap}
	}

	out := [][]string{}
	for i := range unfixed {
		swap := make([]string, len(unfixed))
		copy(swap, unfixed)
		swap[0], swap[i] = swap[i], swap[0]
		out = append(out, allPermutations(append(fixed, swap[0]), swap[1:])...)
	}
	return out
}
