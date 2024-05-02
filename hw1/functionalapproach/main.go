package main

import "fmt"

type Opt int

const (
	size Opt = iota
	char
	color
)

type OptValue interface{}

type OptsMap map[Opt]OptValue

func funcWithOpts(optMaps ...OptsMap) {
	summaryMap := OptsMap{size: 15, char: 'X', color: 0} // default values
	for _, optMap := range optMaps {
		for opt, optValue := range optMap {
			summaryMap[opt] = optValue
		}
	}

	fmt.Println(summaryMap[size].(int), string(summaryMap[char].(rune)), summaryMap[color].(int))
}

func WithSize(v int) OptsMap {
	return OptsMap{size: v}
}

func WithChar(v rune) OptsMap {
	return OptsMap{char: v}
}

func WithColor(v int) OptsMap {
	return OptsMap{color: v}
}

func main() {
	funcWithOpts(WithColor(1), WithSize(3), WithChar('X'))
}
