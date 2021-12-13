package internal

import (
	"math/rand"
	"time"
)

func Shuffle[T any](in []T) {
	rand.Seed(time.Now().Unix())
	for i := len(in) - 1; i >= 0; i-- {
		j := rand.Int() % (i + 1)
		in[i], in[j] = in[j], in[i]
	}
}
