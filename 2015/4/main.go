package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	key := []byte(strings.TrimSpace(lines[0]))
	for i := 0; ; i++ {
		out := fmt.Sprintf("%x", md5.Sum(append(key, []byte(strconv.Itoa(i))...)))
		if strings.HasPrefix(out, "000000") {
			fmt.Println(i)
			os.Exit(0)
		}
	}

	//part 1
	// key := []byte(strings.TrimSpace(lines[0]))
	// for i := 0; ; i++ {
	// 	out := fmt.Sprintf("%x", md5.Sum(append(key, []byte(strconv.Itoa(i))...)))
	// 	if strings.HasPrefix(out, "00000") {
	// 		fmt.Println(i)
	// 		os.Exit(0)
	// 	}
	// }
}
