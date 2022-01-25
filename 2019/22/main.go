package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal/input"
)

// This is the only one i really couldn't do. I understood the concepts behind it, nearly got there myself but
// then had to look up the answers as i didn't care enough at this point.

/*
// 116288486925079 too high
//  44646848106000 too high
*/

func main() {
	instructions := input.GetInput(2019, 22)
	packSize := 19
	packSize = 10007
	packSize = 119315717514047
	n := 101741582076661
	x := 2020
	fmt.Println(k(instructions, packSize, 1, x))
	fmt.Println(k(instructions, packSize, n, x))
}

func k(instructions []string, packSize int, n, x int) int {
	bps := big.NewInt(int64(packSize))
	// 	og := 2020
	// //	og = 13360576774089
	// 	indexToTrack := og
	// 	fmt.Println(f(instructions, indexToTrack, packSize))
	a, b, _, _ := coefs(instructions, packSize)
	fmt.Println("a", a, "b", b)
	bx := big.NewInt(int64(x))
	bn := big.NewInt(int64(n))

	//	pow(A, n, D)*X + (pow(A, n, D)-1) * modinv(A-1, D) * B) % D
	first := big.NewInt(0)
	first.Exp(a, bn, bps).Mul(a, bx)
	second := big.NewInt(0)
	second.Exp(a, bn, bps).Sub(second, big.NewInt(1))
	third := big.NewInt(0)
	fourth := big.NewInt(0)
	fourth.Sub(a, big.NewInt(1))
	third.ModInverse(fourth, bps)
	second.Mul(second, third).Mul(second, b)
	first.Add(first, second).Mod(first, bps)
	return int(first.Int64())
}

func g(i int, ps, a, b *big.Int) int {
	in := big.NewInt(int64(i))
	out := in.Mul(in, a).Add(in, b).Mod(in, ps)
	return int(out.Int64())
}

func coefs(instructions []string, packsize int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	x := 2020
	y := f(instructions, x, packsize)
	z := f(instructions, y, packsize)
	bx := big.NewInt(int64(x))
	by := big.NewInt(int64(y))
	bd := big.NewInt(int64(packsize))
	bz := big.NewInt(int64(z))
	ba := big.NewInt(0)
	bb := big.NewInt(0)
	bs := big.NewInt(0)
	inv := big.NewInt(0)
	ba.Sub(by, bz)
	bs.Sub(bx, by)
	inv.ModInverse(bs, bd)
	ba.Mul(ba, inv).Mod(ba, bd)
	bq := big.NewInt(0)
	bq.Mul(ba, bx)
	bb.Sub(by, bq).Mod(bb, bd)
	return ba, bb, bx, by
	//	(ba.Sub(by, bz).Mul(by, bx.Sub(bx, by).ModInverse(bx, bd))).Mod(ba, bd)
	//	bb.Sub(by, ba.Mul(ba, bx)).Mod(by, bd)
	//	b := (y - a*x) % packsize
	//	fmt.Println(ba, bb)
	// X = 2020
	// Y = f(X) = 13360576774089
	// Z = f(Y) = 7835275759222
	// A = (Y-Z) * modinv(X-Y, D) % D
	// B = (Y-A*X) % D

}

// func g(i, ps int) int {
// 	a := 1556132394918066323
// 	b := -45751964910590
// 	math.big

// }

func f(instructions []string, indexToTrack, packSize int) int {
	for i := len(instructions) - 1; i >= 0; i-- {
		line := instructions[i]
		words := strings.Split(line, " ")
		if strings.HasPrefix(line, "deal with increment") {
			inc, _ := strconv.Atoi(words[3])
			indexToTrack = rdeal(indexToTrack, inc, packSize)
			//				fmt.Println("after deal", inc, "tracking", indexToTrack)
			continue
		}
		if strings.HasPrefix(line, "cut") {
			idx, _ := strconv.Atoi(words[1])
			indexToTrack = rcut(indexToTrack, idx, packSize)
			//				fmt.Println("after cut", idx, "tracking:", indexToTrack)
			continue
		}
		if strings.HasPrefix(line, "deal into new") {
			indexToTrack = rreverse(indexToTrack, packSize)
			//				fmt.Println("after reverse", indexToTrack)
			continue
		}
	}
	return indexToTrack
	// og2 := 2020
	// indexToTrack = og2
	// double := map[int]struct{}{}
	// for i := 0; i < 1000; i++ {
	// 	for _, line := range instructions {
	// 		words := strings.Split(line, " ")
	// 		if strings.HasPrefix(line, "deal with increment") {
	// 			inc, _ := strconv.Atoi(words[3])
	// 			indexToTrack = deal2(indexToTrack, inc, packSize)
	// 			//			fmt.Println("after deal", inc, "tracking", indexToTrack)
	// 			continue
	// 		}
	// 		if strings.HasPrefix(line, "cut") {
	// 			idx, _ := strconv.Atoi(words[1])
	// 			indexToTrack = cut2(indexToTrack, idx, packSize)
	// 			//			fmt.Println("after cut", idx, "tracking:", indexToTrack)
	// 			continue
	// 		}
	// 		if strings.HasPrefix(line, "deal into new") {
	// 			indexToTrack = reverse2(indexToTrack, packSize)
	// 			//	/	fmt.Println("after reverse", indexToTrack)
	// 			continue
	// 		}
	// 	}
	// 	fmt.Println(og2, "moved to", indexToTrack)
	// 	if _, ok := double[indexToTrack]; ok {
	// 		fmt.Println(i, indexToTrack)
	// 		break
	// 	}
	// 	double[indexToTrack] = struct{}{}
	// }

}

func deal2(idx, n, packSize int) int {
	return (idx * n) % packSize
}
func deal3(idx, n, packSize *big.Int) int {
	return int(idx.Mul(idx, n).Mod(idx, packSize).Int64())
}

// if the cut is negative, convert it to a positive?
// return idx + cutspot
// 7890123456

func cut3(idx, cutSpot, packSize int) int {
	return (idx + cutSpot + packSize) % packSize
}

// cut -3
func cut2(idx, cutSpot, packSize int) int {
	if cutSpot < 0 {
		cutSpot = packSize + cutSpot
	}
	if idx < cutSpot {
		return idx + (packSize - cutSpot)
	}
	return idx - cutSpot
}
func reverse2(idx, n int) int {
	return n - idx - 1
}

// cut -8 reversed with cut 8
// 5678901234
//
func rdeal(idx, n, packSize int) int {
	// find modular inverse of n = 13; g = idx
	inc := big.NewInt(int64(n))
	packSizex := big.NewInt(int64(packSize))
	inv := inc.ModInverse(inc, packSizex)
	//	fmt.Println("reverseing deal", n, "with", inv, "idx", idx)
	f := big.NewInt(int64(idx))
	ps := big.NewInt(int64(packSize))
	return int(f.Mul(f, inv).Mod(f, ps).Int64())
}

// a cut 8 will be reversed with a cut -8
// a cut -2 will be reversed with a cut 2
func rcut(idx, cutSpot, packSize int) int {
	//	fmt.Println("reversing", "cut", cutSpot, "with", -cutSpot)
	return (idx + cutSpot + packSize) % packSize
}
func rreverse(idx, n int) int {
	//	fmt.Println("reversing reverse is reverse")
	return n - idx - 1 // it's the same.
}

func main2() {
	instructions := input.GetInput(2019, 22)

	size := 10007
	size = 100
	cards := make([]int, size)
	for i := 0; i < size; i++ {
		cards[i] = i
	}

	for _, line := range instructions {
		words := strings.Split(line, " ")
		if strings.HasPrefix(line, "deal with increment") {
			inc, _ := strconv.Atoi(words[3])
			cards = deal(cards, inc)
			continue
		}
		if strings.HasPrefix(line, "cut") {
			idx, _ := strconv.Atoi(words[1])
			cards = cut(cards, idx)
			continue
		}
		if strings.HasPrefix(line, "deal into new") {
			cards = reverse(cards)
			continue
		}
	}
	for i, c := range cards {
		fmt.Println(i, c)
	}
	if len(cards) > 3000 {
		for i, v := range cards {
			if v == 2019 {
				fmt.Println("part 1", i)
			}
		}
	}
}

// 5495 too low
// 6952 too high

func reverse(in []int) []int {
	out := make([]int, len(in))
	j := 0
	for i := len(in) - 1; i >= 0; i-- {
		out[j] = in[i]
		j++
	}
	return out
}

func cut(in []int, idx int) []int {
	if idx < 0 {
		idx = len(in) + idx
	}
	return append(in[idx:], in[:idx]...)
}

func deal(in []int, inc int) []int {
	out := make([]int, len(in))
	j := 0
	for i := 0; i < len(in); i++ {
		out[j] = in[i]
		j += inc
		j = j % len(in)
	}
	return out
}

// incrementer
// 119315717514047 packsize
// 101741582076661 shuffles
