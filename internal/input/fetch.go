package input

import (
	"bytes"
	"context"
	"flag"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc/app"
)

var (
	stdin bool
)

func init() {
	flag.BoolVar(&stdin, "stdin", false, "read input from stdin")
	flag.Parse()
}

func GetInput(year, day int) []string {
	if stdin {
		return internal.ReadInput()
	}
	app, _ := app.NewApplication()
	input, err := app.GetInput(context.Background(), year, day)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(bytes.TrimSpace(input)), "\n")
	return internal.CleanInput(lines)
}

func GetRawInput(year, day int) []string {
	if stdin {
		return internal.ReadRawInput()
	}
	app, _ := app.NewApplication()
	input, err := app.GetInput(context.Background(), year, day)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")
	return lines

}
