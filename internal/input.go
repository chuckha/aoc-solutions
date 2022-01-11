package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadInput() []string {
	lines := ReadRawInput()
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return lines
}

func ReadRawInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		data := scanner.Text()
		if data == "" {
			continue
		}
		lines = append(lines, data)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return lines

}

func ReadRealRawInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return lines

}
