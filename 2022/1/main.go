package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	allElves := ChunkedInput(inputToElf)
	//	fmt.Println(allElves)
	sort.Sort(elves(allElves))
	fmt.Println(allElves[len(allElves)-1].TotalCalories + allElves[len(allElves)-2].TotalCalories + allElves[len(allElves)-3].TotalCalories)
}

type Elf struct {
	FoodItems     []int
	TotalCalories int
}

type elves []*Elf

func (e elves) Len() int           { return len(e) }
func (e elves) Less(i, j int) bool { return e[i].TotalCalories < e[j].TotalCalories }
func (e elves) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (e *Elf) String() string {
	return fmt.Sprintf("items (total): %d (%d)\n", len(e.FoodItems), e.TotalCalories)
}

func NewElf(foods ...int) *Elf {
	sum := 0
	for _, f := range foods {
		sum += f
	}
	return &Elf{
		FoodItems:     foods,
		TotalCalories: sum,
	}
}

func inputToElf(in []string) *Elf {
	cals := make([]int, len(in))
	for i := 0; i < len(in); i++ {
		val, err := strconv.Atoi(in[i])
		if err != nil {
			panic(err)
		}
		cals[i] = val
	}
	return NewElf(cals...)
}

func ChunkedInput[T any](build func([]string) T) []T {
	scanner := bufio.NewScanner(os.Stdin)
	chunk := []string{}
	out := make([]T, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			out = append(out, build(chunk))
			chunk = []string{}
			continue
		}
		chunk = append(chunk, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return out
}
