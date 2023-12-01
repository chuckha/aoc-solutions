package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
incorrect guesses pt 1
-6113: incorrect
24829: incorrect
-1031: incorrect
4017: incorrect

1725: too low
2533: too low
6387: correct!
*/

const decryptionKey = 811589153

type indexedNumber struct {
	idx int
	num int
}

func (i indexedNumber) Equal(b indexedNumber) bool {
	return i.num == b.num && i.idx == b.idx
}
func (i indexedNumber) Data() indexedNumber {
	return i
}

const defaultDebug = false

func main() {

	debug := flag.Bool("debug", defaultDebug, "enable debug logging")
	flag.Parse()

	input := internal.ReadInput()
	nums := inputToInt(input, decryptionKey)
	var root *internal.CircularLinkedListV2[indexedNumber]
	var cur *internal.CircularLinkedListV2[indexedNumber]

	zeroIdx := -11
	// max := math.MinInt
	// min := math.MaxInt
	for i, x := range nums {
		if x == 0 {
			zeroIdx = i
		}
		if root == nil {
			root = internal.NewCircularLinkedListV2(indexedNumber{idx: i, num: x})
			cur = root
			continue
		}

		cur = cur.InsertAfter(indexedNumber{idx: i, num: x})
	}
	if zeroIdx == -11 {
		panic("didn't find 0")
	}

	logger := internal.Debug(*debug)
	// fmt.Println(max, min, len(nums))
	// os.Exit(0)
	logger.Println("Initial arrangement")
	logger.Println(root)
	logger.Println()
	//	fmt.Println(root)
	for z := 0; z < 10; z++ {
		for i, y := range nums {
			if y == 0 {
				logger.Println("0 does not move:")
				logger.Println(root)
				logger.Println()
				continue
			}
			find := indexedNumber{idx: i, num: y}
			d := root.Find(find)
			if d == nil {
				panic(fmt.Sprintf("could not find %d", y))
			}
			if d.Data.Equal(root.Data) {
				root = d.Next
			}
			if y < 0 {
				d.Remove().Backwards((-y) % (len(nums) - 1)).InsertBefore(find)
				// insertBeforeHere := d.Backwards((-y) % len(nums))
				// // if insertBeforeHere == d {
				// // 	fmt.Println("landed on myself", y)
				// // 	insertBeforeHere = insertBeforeHere.Backwards(1)
				// // }

				// logger.Printf("%d moves between %d and %d:\n", y, insertBeforeHere.Prev.Data, insertBeforeHere.Data)
				// insertBeforeHere.InsertBefore(find)
				// d.Remove()
				//			x.Backwards((-y) % len(nums)).InsertBefore(y)
			} else {
				d.Remove().Prev.Forward(y % (len(nums) - 1)).InsertAfter(find)
				// insertAfterHere := d.Forward(y % len(nums))
				// // if insertAfterHere == d {
				// // 	fmt.Println("landed on myself", y)
				// // 	insertAfterHere = insertAfterHere.Forward(1)
				// // }

				// logger.Printf("%d moves between %d and %d:\n", y, insertAfterHere.Data, insertAfterHere.Next.Data)
				// insertAfterHere.InsertAfter(find)
				// d.Remove()
				//			x.Forward(y % len(nums)).InsertAfter(y)
			}
			logger.Println(root)
			logger.Println()
		}
	}
	if !root.Verify() {
		panic("invalid list")
	}
	logger.Println("validated list.")

	// groove coordinates
	zero := root.Find(indexedNumber{idx: zeroIdx, num: 0})
	if zero == nil {
		panic("there is no 0")
	}
	a := zero.Forward(1000)
	b := a.Forward(1000)
	c := b.Forward(1000)
	fmt.Println(a.Data, b.Data, c.Data, a.Data.num+b.Data.num+c.Data.num)
}

func inputToInt(in []string, key int) []int {
	out := make([]int, len(in))
	var err error
	for i := 0; i < len(in); i++ {
		out[i], err = strconv.Atoi(in[i])
		if err != nil {
			panic(err)
		}
		out[i] *= key
	}
	return out
}

func testing() {
	n := internal.NewCircularLinkedList(0)
	n.InsertBefore(1)
	n.InsertBefore(2)
	n.InsertBefore(4)
	n.InsertAfter(5)
	x := n.Forward(100)
	m := x.Remove()
	y := m.Forward(911)
	z := y.Remove()
	fmt.Println(z)
}
