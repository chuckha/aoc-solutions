package main

import (
	"fmt"
	"strings"
)

func main() {
	cur := "1"
	for i := 0; i < 5; i++ {
		cur = gen(cur)
		fmt.Println(cur)
	}
	fmt.Println(cur)
}

// 1 -> 11
// 11 -> 21
// 21 -> 1211
func gen(in string) string {
	fmt.Println("going in", in)
	memo := map[string]string{}

	if len(in) == 1 {
		val := make(in)
		memo[in] = val
		return val
	}
	out := []string{}
	hold := []byte{}
	for i := 0; i < len(in); i++ {
		if len(hold) == 0 {
			hold = append(hold, in[i])
			continue
		}
		if in[i] == hold[len(hold)-1] {
			hold = append(hold, in[i])
			continue
		}
		// do something with hold
		key := string(hold)
		val, ok := memo[key]
		if !ok {
			val = gen(key)
			memo[key] = val
		}
		out = append(out, val)
		// clear hold
		hold = []byte{}
	}
	fmt.Println("before flush", hold)
	// flush hold
	key := string(hold)
	val, ok := memo[key]
	if !ok {
		val = gen(key)
		memo[key] = val
	}
	out = append(out, val)

	return strings.Join(out, "")
}

func make(in string) string {
	if in == "1" {
		return "11"
	}
	if in == "2" {
		return "12"
	}
	if in == "3" {
		return "13"
	}
	panic(fmt.Sprintf("unknown, %s", in))
}
