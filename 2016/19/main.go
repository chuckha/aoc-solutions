package main

import (
	"fmt"
	"math"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	//	num := internal.ReadInput()[0]

	// g := game{
	// 	start: 1,
	// 	round: 0,
	// 	last:  3014387,
	// 	total: 3014387,
	// }
	// for g.total > 1 {
	// 	g = g.steal()
	// }
	// fmt.Println(g.atPos(0))
	max := 3014387
	//	max = 7
	if max%2 == 0 {
		panic("only works for odd numbers")
	}
	root := internal.NewNode(1)
	node := root
	for i := 2; i <= max; i++ {
		node = node.Insert(i)
	}
	size := max
	node = root
	node = node.Forward(size / 2)
	// we know it's odd so let's just write the odd code
	for size > 1 {
		node = node.Remove()
		size--
		node = node.Forward(2)
		node = node.Remove()
		size--
		node = node.Forward(1)
	}
	fmt.Println(node)
}

/*

1 2 3 4 5
size: 5
pointer: 1, index: 0
index halfway around is idx+(size/2) % size = 0 + 2 % 5 = 2

1 2 4 5
size: 4
pointer: 2
index: 1
halfway around: 1 + 2 % 4 = 3

1 2 4
size: 3
pointer: 4
index: 2
halfway round = 2 + 1 % 3 = 0

2 4

1 2 3 4 5 6 7 8 9
i: 0 s/2 = 4
halfway = 4; remove: 5

1 2 3 4 6 7 8 9
i: 1 s/2 = 4
halfway = 5; remove: 7

1 2 3 4 6 8 9
i: 2 s/2 = 3
halfway = 5; remove: 8

<-> inflection point when idx + s/2 > size
every step reduces size and increases index

size = n - index

size == 9
i = 0
inflection point is at 3 steps
i=1,size/2=4
i=2,size/2=3
i=3,size/2=3

when i = size/2

on the i = size/2 selection it wraps back to 0
so on 9, size/2 = 4
0, 1, 2, (3)
so on 7, size/2 = 3
0, 1, (2)


i: 0
size: n
i: 1
size: n-1
i: 2
size: n-2

i flips back to 0 when

1 2 3 4 6 9
i: 3 s/2 = 3
halfway = 0; remove 1

1 2

|     |      | when do we start removing things from the first half?
total size: s
first half s/2
second half s/2

starting at s/2 re


1 2 3 4 5 6 7
s: 7
v: 1
i: 0
halfway around: 0 + 3 % 6 = 3

1 2 3 5 6 7
s: 6
v: 2
i: 1
halfway round: 1 + 3 % 6 = 4

1 2 3 5 7
s: 5
v: 3
i: 2
halfway round: 2 + 2 % 5 = 0

1 2 3 5
1 3 5
1 5
5

after size/2 = 4 selections we start over at 0/1
until then, we remove the half way element and then +2 from that +3

rm, fwd2, rm, fwd1, rm, fwd2, rm, fwd2, rm, fwd2, rm

1 2 3 4 5 6 7 8 9 10 11 | s = 11 r=0
1 2 3 4 5   7 8 9 10 11 | s = 10 r=1
1 2 3 4 5   7   9 10 11 | s = 9  r=2
1 2 3 4 5   7     10 11 | s = 8  r=3
1 2 3 4 5   7     10    | s = 7  r=4
1   3 4 5   7     10    | s = 6  r=5
1   3   5   7     10    | s = 5  r=6
1   3       7     10    | s = 4  r=7
1   3       7           | s = 3  r=8
    3       7           | s = 2  r=9
    3                   | s = 1  r=10

1 2 3 4 5 6 7 8 9 10  | s = 10
1 2 3 4 5   7 8 9 10  | s = 9
1 2 3 4 5     8 9 10  | s = 8
1 2 3 4 5     8   10  | s = 7
1 2 3 4 5     8       | s = 6
1   3 4 5     8
1     4 5     8
1     4       8
1     4
1

first remove s/2-1 items  as the s/2 item will cycle back

1 2 3 4 5 6 7  8 9 10 11 12 13 14 15

8 9 10 11 12 13 14 15
  9 10 11 12 13 14 15
  9    11 12 13 14 15
  9       12 13 14 15
  9       12    14 15
  9       12       15

1 2 3 4 5 6 7 8 9 10 11 12 13 14

8 9 10 11 12 13 14
  9 10 11 12 13 14
    10 11 12 13 14
    10    12 13 14

a b c
c
idx = 1


a b
idx = 0/1
winner = idx


1 2 3 4 5 6
2 3 5 6 1
3 6 1 2
6 2 3
3 6
3

1 2 3 4 5
2 4 5 1
4 1 2
2 4
2



1 2 3 4 5 | s = 5 r = 0
1 2   4 5 | s = 4 r = 1
1 2   4   | s = 3 r = 2
  2   4   | s = 2 r = 3
  2       | s = 1 r = 4

1 2 4 5
size / 2

if size is odd, s/2+1 is the infleciton point
if size is even s/2 is the inflection point




1 2 3 4 5
1 2 0 4 5
1 2 0 4 0
0 2 0 4 0
0 2 0 0 0

idx + len/2 == item to remove

if len == 0
1 2 4 5


1 2 3 4 5 6
1 2 3 0 5 6
1 2 3 0 0 6
0 2 3 0 0 6
0 0 3 0 0 6
0 0 3 0 0 0

1 2 3 4 5 6 7
1 2 3 0 5 6 7
1 2 3 0 5 0 7
1 2 3 0 5 0 0
1 0 3 0 5 0 0
1 0 0 0 5 0 0
0 0 0 0 5 0 0


1 2 4 5
1 + 1/2len = 3


1 2 3 4 5
idx = 1 + 2 == remove 3
len = 5
pick idx + len/2 and remove it
2

1 2 4 5
idx = 2
idx = (idx + 1 ) % len = 1
len = 4
pick idx + len/2 and remove it

1 2 4
idx = (3 + 1 + 2) %
idx = (idx + 1 ) % len = 2
len = 3
pick (idx + len/2) % len and remove it

2 4
idx = (idx + 1 ) % len = 0
len = 2
pick (idx + len/2) % len and remove it
2




1 2 3 4 5 6
idx = 0
len = 6
pick (idx + len/2) % len and remove it

1 2 3 5 6
idx = (idx + 1 ) % len
len = 5

1 2 3 4 5
7 11





1 2 3 4 5
0 1 2 3 4

idx = 0 -> 2/3
	  1 -> 3/4
	  2 -> 4/0
	  3 -> 0/1
	  4 -> 1/2
index + total / 2 =
0 -> 2.5
1 -> 3
2 ->
1 2 4 5
0 1 2 3
idx = 1

*/

type game struct {
	start int
	round int
	last  int
	total int
}

func (g game) atPos(n int) int {
	if n == 0 {
		return g.start
	}
	if n == g.total-1 {
		return g.last
	}
	if g.round == 0 {
		return n + 1
	}
	coef := math.Pow(2, float64(g.round))
	return g.start + int(coef)*n
}

func (g game) steal() game {
	out := game{
		round: g.round + 1,
		total: g.total / 2,
	}
	if g.total%2 == 0 {
		out.start = g.start
		out.last = g.atPos(g.total - 2)
		return out
	}
	out.start = g.atPos(2)
	out.last = g.last
	return out
}

/*

        12
    11      1
  10          2
9	           3
  8          4
    7	    5
         6


0  1 2 3  4 5 6 7 8  9 10
1  2 3 4  5 6 7 8 9 10 11
3  5 7 9 11
7 11
7
round = 0
0 = 1 = 1 + 2 * 0
1 = 2 = 1 + 2 * 1
2 = 3 = 1 + 2 * 2
3 = 4 = 1 + 2 * 3
4 = 5

round = 1
0 = 3 = 3 + 2 *1 * 0
1 = 5 = 3 + 2 *1 * 1
2 = 7 = 3 + 2 *1 * 2
3 = 9 = 3 + 2 *1 * 3
4 = 11

0 1 2 3 4  5 6 7 8  9 10 11
1 2 3 4 5  6 7 8 9 10 11 12
1 3 5 7 9 11
1 5 9

round = 2
0 = 1 = 1
1 = 5 = 1 + 2 *2 * 1
2 = 9 = 1 + 2 *2 * 2
first number == 1
len = 11
last number == 11

after round 1
because the previous round was odd in length, this one will start with the second odd number of the previous round
it will contain every other element in the previous list
first number == 3
len = floor(11/2) = 5
last number == 11
all numbers are the odd ones

after round 2
first number

(i) i + (2 * round) (...)
(i + 2*round)

1 2 3 4 5 6 7 8 9 10 11
3 5 7 9 11 round = 1, start = 3 = 3 + 2 * 0, 5 = start + 2 * 1, 7 = start + 2 * 2 = start + 2 * 3
7 11 round = 2, start = 7, start + 4
7

123456789 10
13579
59
5

123456789
3579
37
3

12345678
1357
15
1

1234567
357
7

123456
135
5

12345
35
3

1234
13
1

5/2 = 2
3014387/2 =1507193/2 = 753596/2=376798/2 =188399/2=94199/2=47099/2=23549/2=11774/2=5887/2=2943/2=1471/2=735/2=367/2=183/2=91/2=45/2=22/2=11/2=5/2=2/2=1


when there are 5 left, the middle number wins
when there are 4 left,
when there are 3 left, the last number wins
when there are 2 left, the first number wins


1: i						1234567
2: i*2+1					357
3: i*?+3					7



1
5
9
13
17
21
25
29
33
37
41
45
49
53
57
61
65
69
73
77
81
85
89
93
97

round = 3
9 25 33 41 49 57 65 73　81　89　97

0 = 9 = 9
1 = 17 = 9 + 2 * 2 * 2
2 = 25 = 9 + 2 * 2
3 = 33 = 9 + 8 * 3





*/
