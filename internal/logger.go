package internal

import "fmt"

type Logger struct {
	Debug bool
}

func (l *Logger) Println(a ...interface{}) {
	if l.Debug {
		fmt.Println(a...)
	}
}

func (l *Logger) Printf(f string, a ...interface{}) {
	if l.Debug {
		fmt.Printf(f, a...)
	}
}
