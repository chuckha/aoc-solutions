package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

// 267272 part 1 is too low

func main() {

	lines := input.GetInput(2019, 14)
	input, produced := newRules2(lines)
	// for k, rule := range input {
	// 	fmt.Println(k, rule)
	// }
	// for k, p := range produced {
	// 	fmt.Println(k, p)
	// }
	r := &rules{
		lookup:     input,
		amountMade: produced,
	}
	start := 1000000
	end := 10000000
	next := (start + end) / 2
	cpf := r.cost(1)
	for {
		fmt.Println(next)
		cost := r.cost(next)
		fmt.Println(cost)
		if cost > 1000000000000 {
			end = next
		} else if cost < 1000000000000-cpf {
			start = next
		} else {
			break
		}
		next = (start + end) / 2
	}
	fmt.Println("cost", next)
}

// i need 1 A
// you only get stacks of 10 so that will be 10 ORE and you will have 10A
// i need 1 B
// ok, that's 1 ORE
// i need 1c
// that's 7 a and 1 b
// i need 2 fuel
// ok that's
// i need 2 A
// ok that's 10 ORE

type rules struct {
	lookup     map[string][]amount
	amountMade map[string]int
}

// return cost and leftovers
func (r *rules) cost(fuel int) int {
	//	look up how much  the thing costs
	q := internal.NewQueue[amount]()
	q.Enqueue(amount{chem: "FUEL", quantity: fuel})
	leftovers := map[string]int{}
	total := 0
	for !q.Empty() {
		cur := q.Dequeue()
		// fmt.Println("creating", cur)
		// fmt.Println("leftovers before creating", leftovers)
		if cur.chem == "ORE" {
			total += cur.quantity
			continue
		}
		// not ore
		// this is the cost of x items which means we need to keep track of leftovers
		cost := r.lookup[cur.chem]
		amountOfChem := r.amountMade[cur.chem] // amount of chemical that will be produced

		// apply leftovers
		if leftovers[cur.chem] > cur.quantity {
			leftovers[cur.chem] -= cur.quantity
			continue
		}
		needed := cur.quantity - leftovers[cur.chem]
		leftovers[cur.chem] = 0 // true because we already handled the case of having more than enough leftovers

		// if we produce less than we need, find the multiplier
		multiplier := 1
		for multiplier*amountOfChem < needed {
			multiplier += 1
		}
		// get the total amount produced
		amountProduced := multiplier * amountOfChem
		// if we produce more than we need, store the remainder as leftovers
		if amountProduced > needed {
			leftovers[cur.chem] = amountProduced - needed
		}
		for _, c := range cost {
			item := amount{chem: c.chem, quantity: c.quantity * multiplier}
			q.Enqueue(item)
		}
	}
	return total
	// push item on queue
	// pop item, get cost of item, push all on queue
	// eventually, you get ores
	// add them up

	//
}

func calculateCost(rules map[amount][]amount, item amount) int {
	if item.chem == "ORE" {
		return item.quantity
	}
	sum := 0
	for _, cost := range rules[item] {
		sum += calculateCost(rules, cost)
	}
	return sum
}
func newRules2(lines []string) (map[string][]amount, map[string]int) {
	rules := map[string][]amount{}
	amountMade := map[string]int{}
	for _, line := range lines {
		sides := strings.Split(line, " => ")
		making := strings.Split(sides[1], " ")
		lqty, _ := strconv.Atoi(making[0])
		amountMade[making[1]] = lqty
		required := []amount{}
		left := strings.Split(sides[0], ", ")
		for _, l := range left {
			k := strings.Split(l, " ")
			qty, _ := strconv.Atoi(k[0])
			required = append(required, amount{quantity: qty, chem: k[1]})
		}
		rules[making[1]] = required
	}
	return rules, amountMade
}

func newRules(lines []string) map[amount][]amount {
	rules := map[amount][]amount{}
	for _, line := range lines {
		sides := strings.Split(line, " => ")
		making := strings.Split(sides[1], " ")
		lqty, _ := strconv.Atoi(making[0])

		required := []amount{}
		left := strings.Split(sides[0], ", ")
		for _, l := range left {
			k := strings.Split(l, " ")
			qty, _ := strconv.Atoi(k[0])
			required = append(required, amount{quantity: qty, chem: k[1]})
		}
		rules[amount{quantity: lqty, chem: making[1]}] = required
	}
	return rules
}

type amount struct {
	quantity int
	chem     string
}
