package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()
	sum := 0
	for _, line := range input {
		num := convert(line)
		sum += num
	}
	fmt.Println(convertBack(sum))
}

func convert(in string) int {
	digits := strings.Split(in, "")
	snafud := make([]int, len(digits))
	for i, d := range digits {
		switch d {
		case "-":
			snafud[i] = -1
		case "=":
			snafud[i] = -2
		default:
			x, err := strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
			snafud[i] = x
		}
	}
	return snafuToDec(reverse(snafud))
}

func reverse[T any](in []T) []T {
	out := make([]T, len(in))
	j := len(out) - 1
	for i := 0; i < len(in); i++ {
		out[j] = in[i]
		j--
	}
	return out
}

func snafuToDec(in []int) int {
	sum := 0
	for i := 0; i < len(in); i++ {
		base := exp(5, i)
		sum += base * in[i]
	}
	return sum
}

func exp(base, pow int) int {
	out := 1
	for i := 0; i < pow; i++ {
		out = out * base
	}
	return out
}

func convertBack(in int) string {
	//(n ) / (5^i) == x % 5 == ith digit
	i := 0
	out := []string{}
	carry := false
	for {
		thing := (in / exp(5, i))
		d := thing % 5
		if carry {
			d += 1
			carry = false
		}
		if d == 5 {
			carry = true
			d = 0
		}
		if d == 3 || d == 4 {
			carry = true
		}
		out = append(out, digitMap[d])
		if thing == 0 && !carry {
			break
		}
		i++
	}
	return strings.Join(reverse(out), "")
}

/*

7 =>
12

7 - 2 = 5
5 / 5 = 1

*/

var digitMap = map[int]string{
	0: "0",
	1: "1",
	2: "2",
	3: "=",
	4: "-",
}

/*

the front digit will always be (number+2) / 5 (0-4)
the next digit will always be number - (front digit * 5)
// the digits in the middle??
the 0th digit will always be the (number +2) converted to string
the 0th digit will always be the number minus the closest multiple of 5 (converted to =,-,0,1,2)

if num % 5 == 0 (final digit is 0)
if num % 5 == 1 (final digit is 1)
if num % 5 == 2 (final digit is 2)
if num % 5 == 3 (final digit is =)
if num % 5 == 4 (final digit is -)


(n ) / (5^i) == x % 5 == ith digit

5 / 5^0 = 0 % 5 = 0
if 5 / (5^0) == 0 { done}
5 / 5^1 = 1

17 / 5^0 = 17 % 5 = 2
17 / 5^1 = 3


19 / 5^0 = 19 % 5 = 4 == - (if this is a - or =, must carry 1 over to the next )
19 / 5^1 = 3 % 5 == 3 (since the above was a - or =, add 1) == 4
19 / 5^2 = 0

-2  =
-1  -
0   0
1   1
2   2
3  1=
4  1-
5  10
6  11
7  12
8  2=
9  2-
10 20
11 21
12 22
13 3=  // find the closest multiple of 5, (15) (13 - 15 = -2)
14 3-
15 30
16 31
17 32
18 4=
19 4-
20 40
21 41
22 42
23 10=
24 10-
25 100 y = 2
26 101
27 102 // closest multiple of 5, (25) =10,  27 - 25 == 2
28 11=

201 1 // closest multiple of 5 (200) xxx, 201-200 = 1;

(n*(5 ^ y) + 2) / 5 == first digit



(n+2 / 5) = first digit
12345 / 5^0 = 12345 % 5 == 0
12345 / 5^1 = 2469 % 5 == 4 (=) + carry
12345 / 5^2 = 493 % 5 == 3 + carry = 4 = (=)

15625 3125 625 125 25 5 1

0
*/
