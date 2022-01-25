package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()[0]
	for i := 0; i <= 99; i++ {
		for j := 0; j <= 99; j++ {
			if runThing(input, i, j) == 19690720 {
				fmt.Println(i, j, "answer", 100*i+j)
				return
			}
		}
	}
}

func runThing(input string, noun, verb int) int {
	allNums := strings.Split(input, ",")
	nums := []int{}
	for _, n := range allNums {
		x, _ := strconv.Atoi(n)
		nums = append(nums, x)
	}
	nums[1] = noun
	nums[2] = verb
	var a, b int
	var op string
	for i := 0; i < len(nums); i++ {
		switch i % 4 {
		case 0:
			switch nums[i] {
			case 1:
				op = "add"
			case 2:
				op = "mul"
			case 99:
				return nums[0]
			}
		case 1:
			a = nums[nums[i]]
		case 2:
			b = nums[nums[i]]
		case 3:
			switch op {
			case "add":
				nums[nums[i]] = a + b
			case "mul":
				nums[nums[i]] = a * b
			}
			a = 0
			b = 0
		}
	}
	panic("oops")
}

// 190687 too low

type add struct {
	reg1, reg2, reg3 string
}

func (a add) Run(assem *internal.Assem) {
	assem.SetReg(a.reg3, assem.GetVal(a.reg1)+assem.GetVal(a.reg2))
}

type mul struct {
	reg1, reg2, reg3 string
}

func (m mul) Run(assem *internal.Assem) {
	assem.SetReg(m.reg3, assem.GetVal(m.reg1)*assem.GetVal(m.reg2))
}

type done struct{}

func (d done) Run(assem *internal.Assem) {
	fmt.Println(assem.GetVal("0"))
	os.Exit(0)
}
