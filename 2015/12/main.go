package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	msi := map[string]interface{}{}
	json.Unmarshal([]byte(lines[0]), &msi)
	fmt.Println(findNumber(msi))
}

func findNumber(item interface{}) float64 {
	sum := 0.0
	switch v := item.(type) {
	case map[string]interface{}:
		for _, val := range v {
			if val == "red" {
				return 0
			}
		}
		for _, v := range v {
			sum += findNumber(v)
		}
		return sum
	case []interface{}:
		for _, v := range v {
			sum += findNumber(v)
		}
	case int:
		return float64(v)
	case float64:
		return v
	case string:
		return 0
	default:
		panic(fmt.Sprintf("unhandled type %T", v))
	}
	return sum
}

// if it's a map, recurse
// if it's a list for each, recurse
// if it's number, add it
