package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	bots := map[int]*bot{}
	outputs := map[int]*output{}
	for _, line := range lines {
		words := strings.Split(line, " ")
		if strings.HasPrefix(line, "value") {
			botID, _ := strconv.Atoi(words[5])
			bots[botID] = newBot(botID)
			continue
		}
		botID, _ := strconv.Atoi(words[1])
		bots[botID] = newBot(botID)
		if words[5] == "bot" {
			bid, _ := strconv.Atoi(words[6])
			bots[bid] = newBot(bid)
		}
		if words[10] == "bot" {
			bid, _ := strconv.Atoi(words[11])
			bots[bid] = newBot(bid)
		}
	}

	for _, line := range lines {
		words := strings.Split(line, " ")
		// pull out the value lines
		if strings.HasPrefix(line, "value") {
			continue
		}

		botID, _ := strconv.Atoi(words[1])
		low, _ := strconv.Atoi(words[6])
		high, _ := strconv.Atoi(words[11])
		//bot 17 gives low to bot 93 and high to bot 109
		//bot 118 gives low to output 7 and high to bot 96
		if words[5] == "bot" {
			bots[botID].giverLow = bots[low]
		}
		if words[5] == "output" {
			if _, ok := outputs[low]; !ok {
				outputs[low] = newOutput()
			}
			bots[botID].giverLow = outputs[low]
		}
		if words[10] == "bot" {
			bots[botID].giverHigh = bots[high]
		}
		if words[10] == "output" {
			if _, ok := outputs[high]; !ok {
				outputs[high] = newOutput()
			}
			bots[botID].giverHigh = outputs[high]
		}
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "value") {
			words := strings.Split(line, " ")
			value, _ := strconv.Atoi(words[1])
			botID, _ := strconv.Atoi(words[5])
			bots[botID].give(value)
		}
	}
	//	fmt.Println(outputs)
}

type bot struct {
	id        int
	chips     [2]int
	giverLow  giver
	giverHigh giver
	bots      map[int]*bot
}

func newBot(id int) *bot {
	return &bot{
		id:    id,
		chips: [2]int{},
	}
}

func (b *bot) give(chip int) {
	if b.chips[0] == 0 {
		b.chips[0] = chip
	} else if b.chips[1] == 0 {
		b.chips[1] = chip
	} else {
		panic("too many chips?")
	}

	if b.chips[0] != 0 && b.chips[1] != 0 {
		low := b.chips[0]
		high := b.chips[1]
		if low > high {
			low, high = high, low
		}
		//		fmt.Printf("bot %d is giving away %d %d\n", b.id, low, high)
		if low == 17 && high == 61 {
			fmt.Println(b.id)
		}
		b.chips[0] = 0
		b.giverLow.give(low)
		b.chips[1] = 0
		b.giverHigh.give(high)
	}
}

type giver interface {
	give(int)
}

type output struct {
	data []int
}

func newOutput() *output {
	return &output{
		data: make([]int, 0),
	}
}
func (o *output) give(i int) {
	o.data = append(o.data, i)
}
func (o *output) String() string {
	return fmt.Sprint(o.data)
}
