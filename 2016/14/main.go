package main

import (
	"crypto/md5"
	"fmt"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	salt := internal.ReadInput()[0]

	keyCount := 0

	hashes := map[int]*hash{}

	for i := 0; ; i++ {
		// we've already generated this hash
		if _, ok := hashes[i]; !ok {
			hashes[i] = hashit(salt, i)
		}

		// make the next 1000 hashes
		if hashes[i].triple != "" {
			for j := i + 1; j <= i+1000; j++ {
				if _, ok := hashes[j]; !ok {
					hashes[j] = hashit(salt, j)
				}
				for c := range hashes[j].pentas {
					if c == hashes[i].triple {
						//		fmt.Println("inc key count", i, j)
						keyCount++
						break
					}
				}
			}

		}
		if keyCount == 64 {
			fmt.Println(i)
			break
		}
	}
}

func hashit(salt string, i int) *hash {
	data := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%d", salt, i))))
	for x := 0; x < 2016; x++ {
		data = fmt.Sprintf("%x", md5.Sum([]byte(data)))
	}
	return newHash(i, data)
}

type hash struct {
	index  int
	data   string
	triple string
	pentas map[string]struct{}
}

func newHash(i int, d string) *hash {
	t, f := findInterestingData(d)
	return &hash{
		index:  i,
		data:   d,
		triple: t,
		pentas: f,
	}
}

func findInterestingData(h string) (string, map[string]struct{}) {
	triple := ""
	pentas := map[string]struct{}{}
	cur := h[0]
	count := 1
	for i := 1; i < len(h); i++ {
		if h[i] == cur {
			count++
			if count == 3 && triple == "" {
				triple = string(cur)
			}
			if count == 5 {
				pentas[string(cur)] = struct{}{}
			}
			continue
		}
		cur = h[i]
		count = 1
	}
	return triple, pentas
}
