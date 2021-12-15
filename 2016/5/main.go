package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	doorID := lines[0]
	counter := 0
	password := make([]string, 8)
	foundCount := 0
	for foundCount < 8 {
		code := fmt.Sprintf("%s%d", doorID, counter)
		hash := fmt.Sprintf("%x", md5.Sum([]byte(code)))
		if strings.HasPrefix(hash, "00000") {
			if hash[5] >= 48 && hash[5] <= 55 {
				loc, _ := strconv.Atoi(string(hash[5]))
				if password[loc] != "" {
					counter++
					continue
				}
				password[loc] = string(hash[6])
				foundCount++
			}
		}
		counter++
	}
	fmt.Println(strings.Join(password, ""))
}
