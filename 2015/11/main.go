package main

import "fmt"

const (
	aLet = 97
	iLet = 105
	oLet = 111
	lLet = 108
	zLet = 122
)

func main() {
	// invalid := "iol"
	// fmt.Println(invalid[0], invalid[1], invalid[2])
	// fmt.Println("should be true", hasStraight("abc"))
	// fmt.Println("should be false", hasStraight("aaa"))
	// fmt.Println("should be true", hasStraight("ssssssssssssstu"))

	pw := "vzbxkghb"
	x := 0
	for i := 0; i < 10000000; i++ {
		if hasStraight(pw) && hasTwoPair(pw) {
			x++
			if x == 2 {
				fmt.Println(pw)
			}
		}
		pw = inc(pw)
	}
}

func inc(password string) string {
	out := make([]byte, len(password))
	copy(out, password)
	reset := false
	for i := 0; i < len(password); i++ {
		if reset {
			out[i] = aLet
		}
		if password[i] == iLet || password[i] == oLet || password[i] == lLet {
			out[i] += 1
			reset = true
		}
	}

	for i := len(password) - 1; i >= 0; i-- {
		next := int(password[i] + 1)
		if next > zLet {
			out[i] = aLet
			continue
		}
		if next == iLet || next == oLet || next == lLet {
			out[i] = password[i] + 2
			break
		}
		out[i] = password[i] + 1
		break
	}
	return string(out)
}

func hasStraight(in string) bool {
	prev := 0
	inc := 1
	for i := 0; i < len(in); i++ {
		if inc == 3 {
			return true
		}
		if int(in[i])-prev == 1 {
			inc++
			prev = int(in[i])
			continue
		}
		inc = 1
		prev = int(in[i])
	}
	return inc == 3
}

func hasTwoPair(in string) bool {
	var pairCountKind byte
	pairCount := 0
	for i := 0; i < len(in)-1; i++ {
		if in[i] == in[i+1] {
			if pairCountKind == in[i] {
				continue
			}
			pairCountKind = in[i]
			pairCount += 1
		}
		if pairCount == 2 {
			return true
		}
	}
	return false
}
