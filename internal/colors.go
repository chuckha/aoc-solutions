package internal

import "fmt"

func Black(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;40m%s\033[0m", in))
}

func Green(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;42m%s\033[0m", in))
}

func Red(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;41m%s\033[0m", in))
}

func Blue(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;44m%s\033[0m", in))
}

func Yellow(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;43m%s\033[0m", in))
}

func Pink(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;45m%s\033[0m", in))
}
