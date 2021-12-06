package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	b, err := ioutil.ReadFile("./2015/8/tmpfile.txt")
	if err != nil {
		panic(err)
	}
	lines := bytes.Split(b, []byte("\n"))

	actualSizeSum := 0
	encodedSizeSum := 0
	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		encoded := fmt.Sprintf("%q", line)
		actualSizeSum += len(line)
		encodedSizeSum += len([]byte(encoded))
	}
	fmt.Println(encodedSizeSum - actualSizeSum)
}

func length(in []byte) int {
	in = in[1 : len(in)-1]
	in = bytes.Replace(in, []byte{92, 92}, []byte("*"), -1)
	in = bytes.Replace(in, []byte{92, 34}, []byte("*"), -1)
	size := len(in)
	size -= (bytes.Count(in, []byte{92, 120}) * 3)
	return size
}
