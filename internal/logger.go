package internal

import "fmt"

type Debug bool

func (d Debug) Println(a ...interface{}) {
	if d {
		fmt.Println(a...)
	}
}

func (d Debug) Printf(f string, a ...interface{}) {
	if d {
		fmt.Printf(f, a...)
	}
}
