package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/chuckha/aoc-solutions/internal"
	"golang.org/x/exp/constraints"
)

func main() {
	packeta := make([]interface{}, 0)
	packetb := make([]interface{}, 0)
	json.Unmarshal([]byte("[[2]]"), &packeta)
	json.Unmarshal([]byte("[[6]]"), &packetb)

	lines := internal.ReadInput()

	input := [][]interface{}{}
	for _, line := range lines {
		parsed := make([]interface{}, 0)
		json.Unmarshal([]byte(line), &parsed)
		input = append(input, parsed)
	}
	input = append(input, packeta, packetb)
	sort.Sort(data(input))
	prod := 1
	for i, d := range input {
		if reflect.DeepEqual(packeta, d) || reflect.DeepEqual(packetb, d) {
			prod *= (i + 1)
		}
	}
	fmt.Println(prod)
}

type data [][]interface{}

func (d data) Len() int { return len(d) }
func (d data) Less(i, j int) bool {
	out := compLists(d[i], d[j])
	//	fmt.Println("got", out)
	return out == -1
}
func (d data) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// l < r == -1 (can return early (success))
// l == r == 0 (continue comparing)
// l > r == 1 (can return early (fail))
func compLists(left, right []interface{}) int {
	//	fmt.Printf("comparing %v %v\n", left, right)
	for i := 0; i < len(left); i++ {
		// if right runs out of values first, they are in the wrong order
		if i >= len(right) {
			return 1
		}

		lvalue := left[i]
		rvalue := right[i]

		switch v := lvalue.(type) {
		case float64:
			switch v2 := rvalue.(type) {
			case float64:
				c := cmp(v, v2)
				switch c {
				case -1, 1:
					return c
				case 0:
					// next element
					continue
				default:
					panic("bad return")
				}
			case []interface{}:
				// exactly one is a list
				c := compLists(convert(v), v2)
				switch c {
				case -1, 1:
					return c
				case 0:
					continue
				default:
					panic("bad return")
				}
			default:
				panic(fmt.Sprintf("unhandled type on right (L:int) %T", v))
			}
		case []interface{}:
			switch v2 := rvalue.(type) {
			case float64:
				return compLists(v, convert(v2))
			case []interface{}:
				c := compLists(v, v2)
				switch c {
				case -1, 1:
					return c
				case 0:
					continue
				default:
					panic("bad return")
				}
			default:
				panic(fmt.Sprintf("unhandled type on right (L:[]interface{}]) %T", v))
			}
		default:
			panic(fmt.Sprintf("unhandled type on left %T, %v", v, v))
		}
	}
	// if the left runs out of items first, they are in the right order
	if len(left) < len(right) {
		return -1
	}
	return 0
}

// l < r == -1 (can return early (success))
// l == r == 0 (continue comparing)
// l > r == 1 (can return early (fail))
func cmp[T constraints.Ordered](l, r T) int {
	//	fmt.Println("comparing", l, r)
	switch {
	case l < r:
		return -1
	case l == r:
		return 0
	default:
		return 1
	}
}

func convert[T any](t T) []interface{} {
	return []interface{}{t}
}
