// не удовлятворяет условию, т.к. использует struct,
// но дополняет способы реализовать опциональные параметры
package main

import "fmt"

type SandglassOptions struct {
	size  int
	char  rune
	color int
}

func funcWithOpts(options *SandglassOptions) {
	if options.size == 0 {
		options.size = 15
	}

	if options.char == 0 {
		options.char = 'X'
	}

	//color имеет нужное дефолтное значение и так

	fmt.Println(options.size, string(options.char), options.color)
}

func main() {
	funcWithOpts(&SandglassOptions{size: 2, char: 'R', color: 5})
}
